# knowledge / 20261012-termux-migration-plan

2026-10-12: Termux bash利用可能。crabcode CLIは --model を公開するが、gpt-5.4-mini と gpt-5.3-codex の読取専用起動は現在のChatGPT認証で非対応エラー。モデル一覧は実行可能性を保証しない。別モデルレビュー未成立。

2026-10-12(追記): crabcode の設定探索は layer1 OpenCode global → layer2 crabcode global（`~/.config/crabcode/crabcode.jsonc` ほか）の順で上位が優先（公式docs 2026-10-12取得、XDG_CONFIG_HOME未設定時は `~/.config` が基準）。同一layerに複数候補があると停止して問い合わせる。credentials は `~/.local/state/crabcode/auth.json`。openai/gpt-6-luna は `--print --no-session-persistence --model` で起動可能（親のグローバル既定は変更しない）。
