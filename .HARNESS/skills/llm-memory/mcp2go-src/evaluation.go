package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const evaluationRecordType = "eval"

// EvaluationInput は、今回の評価として外部から受け取る値です。
// 数値の採点は呼び出し側が行いますが、前回との差分はGo側で計算します。
type EvaluationInput struct {
	EvaluationID  string             `json:"evaluation_id,omitempty"`
	ComparisonKey string             `json:"comparison_key"`
	Role          string             `json:"role,omitempty"`
	Persona       string             `json:"persona,omitempty"`
	TaskFamily    string             `json:"task_family,omitempty"`
	ContractHash  string             `json:"contract_hash,omitempty"`
	Axes          map[string]float64 `json:"axes"`
	EvidenceRefs  []string           `json:"evidence_refs,omitempty"`
}

type EvaluationRelative struct {
	Delta float64 `json:"delta"`
	Trend string  `json:"trend"`
}

type EvaluationPrevious struct {
	MemoryID string             `json:"memory_id"`
	Values   map[string]float64 `json:"values"`
}

type EvaluationMetadata struct {
	RecordType    string                        `json:"record_type"`
	SchemaVersion string                        `json:"schema_version"`
	EvaluationID  string                        `json:"evaluation_id"`
	ComparisonKey string                        `json:"comparison_key"`
	Role          string                        `json:"role,omitempty"`
	Persona       string                        `json:"persona,omitempty"`
	TaskFamily    string                        `json:"task_family,omitempty"`
	ContractHash  string                        `json:"contract_hash,omitempty"`
	Axes          map[string]float64            `json:"axes"`
	Previous      *EvaluationPrevious           `json:"previous,omitempty"`
	Relative      map[string]EvaluationRelative `json:"relative,omitempty"`
	Deterministic EvaluationDeterministic       `json:"deterministic"`
	EvidenceRefs  []string                      `json:"evidence_refs,omitempty"`
	Reuse         EvaluationReuse               `json:"reuse"`
}

type EvaluationDeterministic struct {
	RangeValid         bool `json:"range_valid"`
	BaselineCompatible bool `json:"baseline_compatible"`
}

type EvaluationReuse struct {
	Eligible bool `json:"eligible"`
}

func (in EvaluationInput) normalize() (EvaluationInput, error) {
	in.ComparisonKey = strings.TrimSpace(in.ComparisonKey)
	if in.ComparisonKey == "" {
		return in, fmt.Errorf("comparison_keyは必須でしてよ")
	}
	if len(in.Axes) == 0 {
		return in, fmt.Errorf("axesは1件以上必要でしてよ")
	}
	if in.EvaluationID == "" {
		in.EvaluationID = fmt.Sprintf("eval-%d", time.Now().UTC().UnixNano())
	}
	for name, value := range in.Axes {
		if strings.TrimSpace(name) == "" {
			return in, fmt.Errorf("評価軸名が空ですわ")
		}
		if value < 0 || value > 1 {
			return in, fmt.Errorf("評価軸 %q の値 %.6f は0以上1以下で指定してくださいまし", name, value)
		}
	}
	return in, nil
}

func previousEvaluation(rec *MemoryRecord) (*EvaluationPrevious, error) {
	if rec == nil {
		return nil, nil
	}
	values := make(map[string]float64)
	rawAxes, ok := rec.Metadata["axes"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("前回評価 %s のaxes形式が不正でしてよ", rec.ID)
	}
	for name, raw := range rawAxes {
		value, ok := raw.(float64)
		if !ok {
			return nil, fmt.Errorf("前回評価 %s の評価軸 %q が数値ではありませんわ", rec.ID, name)
		}
		values[name] = value
	}
	return &EvaluationPrevious{MemoryID: rec.ID, Values: values}, nil
}

func calculateEvaluation(in EvaluationInput, previous *EvaluationPrevious) EvaluationMetadata {
	meta := EvaluationMetadata{
		RecordType:    evaluationRecordType,
		SchemaVersion: "eval.v1",
		EvaluationID:  in.EvaluationID,
		ComparisonKey: in.ComparisonKey,
		Role:          in.Role,
		Persona:       in.Persona,
		TaskFamily:    in.TaskFamily,
		ContractHash:  in.ContractHash,
		Axes:          in.Axes,
		Previous:      previous,
		Relative:      make(map[string]EvaluationRelative),
		Deterministic: EvaluationDeterministic{
			RangeValid:         true,
			BaselineCompatible: previous != nil,
		},
		EvidenceRefs: in.EvidenceRefs,
		Reuse:        EvaluationReuse{Eligible: false},
	}
	if previous == nil {
		return meta
	}
	for name, current := range in.Axes {
		old, exists := previous.Values[name]
		if !exists {
			meta.Deterministic.BaselineCompatible = false
			continue
		}
		delta := current - old
		trend := "unchanged"
		if delta > 0 {
			trend = "improved"
		} else if delta < 0 {
			trend = "declined"
		}
		meta.Relative[name] = EvaluationRelative{Delta: delta, Trend: trend}
	}
	for name := range previous.Values {
		if _, exists := in.Axes[name]; !exists {
			meta.Deterministic.BaselineCompatible = false
		}
	}
	return meta
}

func EvaluateAndInsert(ctx context.Context, pool *pgxpool.Pool, in EvaluationInput) (*MemoryRecord, error) {
	var err error
	in, err = in.normalize()
	if err != nil {
		return nil, err
	}
	history, err := ListEvaluationHistoryHeavy(ctx, pool, in.ComparisonKey, 1)
	if err != nil {
		return nil, fmt.Errorf("前回評価の取得に失敗いたしましたわ: %w", err)
	}
	var previous *EvaluationPrevious
	if len(history) > 0 {
		previous, err = previousEvaluation(&history[0])
		if err != nil {
			return nil, err
		}
	}
	meta := calculateEvaluation(in, previous)
	content, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("評価本文のJSON化に失敗いたしましたわ: %w", err)
	}
	tags := []string{"eval"}
	if in.Role != "" {
		tags = append(tags, "role:"+in.Role)
	}
	if in.TaskFamily != "" {
		tags = append(tags, "task:"+in.TaskFamily)
	}
	return InsertMemoryHeavy(ctx, pool, AddMemoryInput{
		Category:  "eval",
		Title:     "eval:" + in.ComparisonKey + ":" + in.EvaluationID,
		ContentL0: string(content),
		Tags:      tags,
		Metadata:  structToMap(meta),
	})
}

func structToMap(value interface{}) map[string]interface{} {
	encoded, _ := json.Marshal(value)
	var result map[string]interface{}
	_ = json.Unmarshal(encoded, &result)
	return result
}

func ListEvaluationHistoryHeavy(ctx context.Context, pool *pgxpool.Pool, comparisonKey string, limit int) ([]MemoryRecord, error) {
	if limit <= 0 {
		limit = 20
	}
	const query = `
		SELECT id, client_id, category, title, content_l0, content_l1, content_l2, tags,
		       current_level, valid_from, valid_to, version, metadata, created_at, updated_at
		FROM v_active_memories
		WHERE category = 'eval'
		  AND metadata->>'record_type' = 'eval'
		  AND metadata->>'comparison_key' = $1
		ORDER BY created_at DESC
		LIMIT $2`
	rows, err := pool.Query(ctx, query, strings.TrimSpace(comparisonKey), limit)
	if err != nil {
		return nil, fmt.Errorf("評価履歴クエリ実行エラーでしてよ: %w", err)
	}
	defer rows.Close()

	var records []MemoryRecord
	for rows.Next() {
		var rec MemoryRecord
		var rawMeta []byte
		if err := rows.Scan(&rec.ID, &rec.ClientID, &rec.Category, &rec.Title,
			&rec.ContentL0, &rec.ContentL1, &rec.ContentL2, &rec.Tags,
			&rec.CurrentLevel, &rec.ValidFrom, &rec.ValidTo, &rec.Version,
			&rawMeta, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, fmt.Errorf("評価履歴レコードのスキャンに失敗いたしましたわ: %w", err)
		}
		if err := json.Unmarshal(rawMeta, &rec.Metadata); err != nil {
			return nil, fmt.Errorf("評価履歴metadataのJSON解析に失敗いたしましたわ: %w", err)
		}
		rec.Status = "ACTIVE"
		records = append(records, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("評価履歴走査エラーでしてよ: %w", err)
	}
	return records, nil
}
