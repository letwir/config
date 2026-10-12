# Method

実行環境は常に `.harness/HARDWARE.md` の `active_profile` を参照し、指定されたプロファイル（`android-termux-bash` では Bash / ネイティブ aarch64 ツール、`windows-pwsh` では PowerShell 7 / Windows ツール）に適応したシェルを使用する。
rgで対象検索、lsdで構成確認、dustで容量確認。実在する仕組みを先に確認し、事実・推論・未確認と完了条件を分ける。既存変更を保持し、調査結果を実装済みと扱わない。非アクティブなプロファイルのシェル・パス記法は使用しない。
