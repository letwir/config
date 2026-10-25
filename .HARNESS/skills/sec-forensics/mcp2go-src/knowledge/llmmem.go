package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// KnowledgeSearchResult represents JSON output from llm-mem search.
type KnowledgeSearchResult struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Category  string   `json:"category"`
	L2Summary string   `json:"l2_summary"`
	L1Points  []string `json:"l1_points"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
}

func getLlmMemExecutablePath() string {
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		return "llm-mem.exe"
	}
	target := filepath.Join(userProfile, ".gemini", "config", "plugins", "ltw-skills-plugin", "skills", "llm-memory", "llm-mem.exe")
	if _, err := os.Stat(target); err == nil {
		return target
	}
	return "llm-mem.exe"
}

// FetchSecurityKnowledgeComplex queries PostgreSQL via llm-mem CLI for CVE/CTF intelligence.
// [Chaotic×ojou-sama] 調査開始前にPostgresから最新のCVE・攻撃シナリオをフェッチいたしますわ！
func FetchSecurityKnowledgeComplex(query, category string) ([]KnowledgeSearchResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	binPath := getLlmMemExecutablePath()
	args := []string{"search", "-json", "-limit", "20"}
	if category != "" {
		args = append(args, "-cat", category)
	}
	if query != "" {
		args = append(args, "-q", query)
	}

	cmd := exec.CommandContext(ctx, binPath, args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("FetchSecurityKnowledge: llm-mem execution failed: %w", err)
	}

	trimmed := strings.TrimSpace(string(out))
	if len(trimmed) == 0 || trimmed == "[]" {
		return []KnowledgeSearchResult{}, nil
	}

	var results []KnowledgeSearchResult
	if err := json.Unmarshal([]byte(trimmed), &results); err != nil {
		return nil, fmt.Errorf("FetchSecurityKnowledge: failed to parse llm-mem JSON output: %w", err)
	}

	return results, nil
}

// IngestSecurityKnowledgeComplex registers or updates CVE/CTF intelligence into PostgreSQL via llm-mem.
// [Chaotic×ojou-sama] 最新の脅威知見をPostgresのllm-mem知識ベースへ刻み込みますわ！
func IngestSecurityKnowledgeComplex(title, content, category, l1, l2 string, tags []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	binPath := getLlmMemExecutablePath()
	tagStr := strings.Join(tags, ",")
	args := []string{
		"add",
		"-title", title,
		"-content", content,
		"-cat", category,
		"-l1", l1,
		"-l2", l2,
		"-tags", tagStr,
	}

	cmd := exec.CommandContext(ctx, binPath, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("IngestSecurityKnowledge: llm-mem add failed: %s: %w", string(out), err)
	}

	return nil
}

// SeedInitialKnowledgeComplex populates baseline Pillow CVEs and CTF attack vectors into PostgreSQL.
// [Chaotic×ojou-sama] 既知のPillow致命的脆弱性とCTF攻撃手法の初期ナレッジを一括投入いたしますわ！
func SeedInitialKnowledgeComplex() error {
	seedEntries := []struct {
		title    string
		cat      string
		content  string
		l1       string
		l2       string
		tags     []string
	}{
		{
			title:   "CVE-2026-59205: Pillow Heap Out-of-Bounds Write in ImageCmsTransform.apply()",
			cat:     "sec_cve",
			content: "Pillow < 12.3.0 では、ICC カラープロファイル出力モードの不一致により ImageCmsTransform.apply() 内でヒープ境界外書き込み（Heap OOB Write）が発生し、潜在的なリモートコード実行 (RCE) またはクラッシュ (DoS) を引き起こす。\n攻撃手法: FLAC や画像の埋め込みカラープロファイルを細工し、デコード時に異常なメモリレイアウトをトリガーする。\n対策: pillow>=12.3.0 への更新、Image.open() でのフォーマット制限。",
			l1:      "- Heap OOB Write in ImageCmsTransform.apply()\n- 攻撃起点: 悪意あるICCプロファイル\n- 対策: pillow>=12.3.0",
			l2:      "Pillow < 12.3.0 の ICC プロファイル処理不備による Heap OOB Write (RCE/DoS) 脆弱性。",
			tags:    []string{"cve", "pillow", "rce", "heap-overflow"},
		},
		{
			title:   "CVE-2026-54058: Pillow Out-of-Bounds Read in McIdas AREA mmap path",
			cat:     "sec_cve",
			content: "Pillow < 12.3.0 の McIdasImagePlugin では、攻撃者が stride=1 かつ巨大な xsize を指定したヘッダを細工すると、mmap 経由の読み出し時に linesize チェックが欠落し、隣接ヒープメモリ（最大数十KB）が画像ピクセルとして読み出され情報漏洩する。\n対策: pillow>=12.3.0、McIdas プラグインの無効化。",
			l1:      "- McIdas AREA ヘッダの stride 不正値によるヒープ境界外読み出し\n- 影響: 最大57KBの隣接メモリ漏洩\n- 対策: pillow>=12.3.0",
			l2:      "Pillow < 12.3.0 の McIdas 画像ヘッダ stride 不備による隣接ヒープメモリ漏洩脆弱性。",
			tags:    []string{"cve", "pillow", "info-disclosure", "oob-read"},
		},
		{
			title:   "CTF-DEFCON: Unicode Tag Smuggling (Invisible ASCII Smuggling)",
			cat:     "sec_immunity",
			content: "DEFCON AI Village で実証された攻撃。Unicode の Tag Plane (U+E0000..U+E007F) を用いると、人間には不可視の文字として画面に表示されるが、LLM のトークナイザーは通常の ASCII 命令として解釈・実行する。\n対策: U+E0000..U+E007F の静的サニタイズ・検知。",
			l1:      "- Unicode Tag Plane (U+E0000..U+E007F) による不可視プロンプトインジェクション\n- 対策: 入力時の Tag プレーン文字除去・検知",
			l2:      "Unicode Tag 領域を用いた人間不可視・LLM実行可能な ASCII Smuggling 攻撃。",
			tags:    []string{"ctf", "defcon", "prompt-injection", "unicode-smuggling"},
		},
		{
			title:   "CTF-SECCON: Cloud Metadata SSRF & Loopback Abuse",
			cat:     "sec_immunity",
			content: "Web / Agent が外部 URL をフェッチする際、169.254.169.254 (AWS/GCP/Azure メタデータ) や 127.0.0.1 (ローカル管理ポート) を指定されると、IAM ロールクレデンシャルや内部 API が窃取される。\n対策: プライベート IP / メタデータホストへのリクエスト遮断 (Egress Firewall / URL Validator)。",
			l1:      "- 169.254.169.254 や localhost への不正フェッチ\n- 対策: URL スキーム & IP ホワイトリスト検証",
			l2:      "クラウドメタデータや内部ループバック宛の SSRF 攻撃手法と防御策。",
			tags:    []string{"ctf", "seccon", "ssrf", "cloud-security"},
		},
	}

	for _, s := range seedEntries {
		if err := IngestSecurityKnowledgeComplex(s.title, s.content, s.cat, s.l1, s.l2, s.tags); err != nil {
			return fmt.Errorf("SeedInitialKnowledge: failed on %s: %w", s.title, err)
		}
	}
	return nil
}
