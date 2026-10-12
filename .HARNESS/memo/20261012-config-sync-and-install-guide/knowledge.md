# knowledge / 20261012-config-sync-and-install-guide

- topic: マルチプラットフォーム対応ハーネスの環境切り替え (INSTALL.md) とリポジトリ同期
- facts:
  1. `repo/config/.HARNESS` はリポジトリ管理される設定原本であり、実行時SSOTである `~/.harness` との整合性を保つ必要がある。
  2. スキルのバイナリ解決規則は `bin=first-existing(~/.harness/skills/<skill>/<bin>, ~/.harness/skills/<skill>/<bin>.exe, ./<bin>, ./<bin>.exe); missing=>STOP` と記述することで、ELFバイナリ（拡張子なし）とPEバイナリ（`.exe`）の双方が同一定義で自然に解決される。
  3. PowerShell 特有の実行構文 `& $bin` は、POSIXシェル（Bash）では構文エラーまたはバックグラウンド実行となるため、環境非依存の呼出し `$bin` に統一する必要がある。
  4. プラットフォーム切り替えは `INSTALL.md` に定める2ステップ：
     - ①: 今の環境を探索（`TERMUX_VERSION` / `uname` / `$PSVersionTable` によるワンライナー自動検出）
     - ②: active_profileを切り替え（`HARDWARE.md` 内のセレクターを `sed` または PowerShell で更新）
     により、ファイル群の書き換えを伴わずにエージェント挙動全体を制御可能。
- pattern:
  - リポジトリ配布設定（`client-distribution.json` やルール・スキル）内のパスはすべてスラッシュ `/` で統一し、各環境のランタイム側で正規化させる設計が最も移植性が高い。
