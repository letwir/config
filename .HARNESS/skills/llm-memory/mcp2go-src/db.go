// Morphism: RawEnvString → ValidatedDSN ∘ PostgreSQLConnectionPool
// Functor: F(ConfigSpace) ⇒ Category(PostgreSQL_Session)

package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DefaultTimeoutSeconds はDB接続および通常クエリのタイムアウト秒数ですわ
const DefaultTimeoutSeconds = 10

// ParseDSNFromEnvComplex は環境変数からDSN文字列をエレガントに解決いたしますわ
func ParseDSNFromEnvComplex() (string, error) {
	raw := os.Getenv("LLM_MEMORY_DB_URL")
	if raw == "" {
		// プロセス環境変数になければWindowsレジストリ/ユーザー環境変数から探しますの
		raw = os.Getenv("LLM_MEMORY_URL")
	}

	if raw == "" {
		return "", fmt.Errorf("環境変数 LLM_MEMORY_DB_URL が設定されておりませんことよ")
	}

	if strings.Contains(raw, "://") {
		return raw, nil
	}

	// pgpassフォーマット (host:port:dbname:user:password) のパース
	parts := strings.Split(raw, ":")
	if len(parts) < 5 {
		return "", fmt.Errorf("DSNのフォーマットが不正でしてよ: host:port:dbname:user:password 形式が必要ですわ")
	}

	host := parts[0]
	port := parts[1]
	dbname := parts[2]
	user := parts[3]
	pass := strings.Join(parts[4:], ":") // パスワード中にコロンが含まれる可能性を考慮

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=prefer",
		url.QueryEscape(user),
		url.QueryEscape(pass),
		host,
		port,
		dbname,
	)
	return dsn, nil
}

// GetDBPoolHeavy は接続プールを安全に初期化して返却いたしますわ
func GetDBPoolHeavy(ctx context.Context) (*pgxpool.Pool, error) {
	dsn, err := ParseDSNFromEnvComplex()
	if err != nil {
		return nil, fmt.Errorf("DSN解決エラーですわ: %w", err)
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("プール設定のパースに失敗いたしましたわ: %w", err)
	}

	// コネクションプールの上限設定
	config.MaxConns = 10
	config.MinConns = 1
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("DBプール作成に失敗いたしましたわ: %w", err)
	}

	// 疎通確認でPingを打ちますわ
	pingCtx, cancel := context.WithTimeout(ctx, DefaultTimeoutSeconds*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("PostgreSQLへのPing疎通に失敗いたしましたわ: %w", err)
	}

	return pool, nil
}
