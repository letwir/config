# Host Hardware & Execution Environment Specification (HARDWARE.md)

本仕様書は、実行環境・ハードウェア仕様・シェル環境・パス体系・およびツールチェインの正本（SSOT: Single Source of Truth）ですわ。
**本書は環境切り替え記法（Profile Selector）を備えており、`active_profile` の設定値を変更することで、Android-Termux / Bash 環境と Windows / PowerShell 環境（または汎用 Linux 環境）をハーネス全体で即座に切り替えられる構造になっています。**

本書は **HarnessRoot を基準とする `HARDWARE.md`** です。プロセスのカレントディレクトリ（CWD）を基準に解決してはなりません。すべてのエージェント、サブエージェント、スクリプトは、コマンド生成や実行、パス解決の前に必ず本書の **`active_profile`** を確認し、選択されたプロファイルに適合したシェルとツールチェインを選択しなければなりません。

---

<!-- ENVIRONMENT_PROFILE_SWITCH -->
## 🔄 アクティブ環境設定 (Active Profile Selector)

| 項目 | 現在の設定値 | 利用可能なプロファイル |
| :--- | :--- | :--- |
| **`active_profile`** | **`android-termux-bash`** | `android-termux-bash`, `windows-pwsh`, `linux-native-bash` |

```yaml
# Profile Switcher Directive
hardware_environment:
  active_profile: "android-termux-bash"
  fallback_profile: "linux-native-bash"
```

> [!IMPORTANT]
> **エージェントへの指示**:
> コマンド実行前に必ず上記 `active_profile` の値を確認してください。
> - `active_profile: "android-termux-bash"` の場合 ➔ **【プロファイル 1】** に従い、Bash / POSIXパス / Termuxネイティブツールを使用（PowerShell / cmd / Windowsパスは禁止）。
> - `active_profile: "windows-pwsh"` の場合 ➔ **【プロファイル 2】** に従い、PowerShell 7 / Windowsパス / `.exe` バイナリを使用。
> - `active_profile: "linux-native-bash"` の場合 ➔ **【プロファイル 3】** に従い、標準 Linux Bash / POSIXパスを使用。
<!-- /ENVIRONMENT_PROFILE_SWITCH -->

---

## 📋 プロファイル定義一覧 (Profile Definitions)

### 【プロファイル 1】 `android-termux-bash` (現在のアクティブ環境 🟢)

Android 端末上の Termux ユーザランド環境。高電力効率・オンデバイスエージェント運用プロファイル。

| 区分 | 項目 | 値・仕様 |
| :--- | :--- | :--- |
| **ハードウェア** | デバイス | Sony XQ-GE44 (Sony Xperia 1 VI) |
| | SoC / CPU | Qualcomm Snapdragon 8 Gen 3 (8 コア: 1x X4, 5x A720, 2x A520) |
| | CPU アーキテクチャ | `aarch64` (`arm64-v8a`) |
| | メモリ (RAM) | 16 GB (LPDDR5X, 実効 ~15.6 GB) |
| | ストレージ | 449 GB (`/data/user/0`, 空き ~267 GB) |
| **OS / カーネル** | OS | Android 16 (Release 16, targetSdkVersion 28 / API 36) |
| | カーネル | Linux 6.12.38-android16-5-gcaefa63e005e-ab14676683-4k (`Linux-aarch64`) |
| | SELinux コンテキスト | `u:r:untrusted_app_27:s0:c31,c258,c512,c768` |
| **ランタイム** | パッケージ環境 | Termux `0.118.3` (Debian apt 形式) |
| | ホーム (`$HOME`) | `/data/data/com.termux/files/home` (`~`) |
| | プレフィックス (`$PREFIX`) | `/data/data/com.termux/files/usr` |
| | 一時ディレクトリ (`$TMPDIR`) | `/data/data/com.termux/files/usr/tmp` (**重要: `/tmp` は権限拒否**) |
| **シェル & コマンド** | 主シェル | **GNU Bash** (`/data/data/com.termux/files/usr/bin/bash`) / POSIX sh |
| | コマンド実行方式 | `bash -c '<command>'` または直接バイナリ呼出し |
| | スクリプト要件 | `#!/usr/bin/env bash`、`set -euo pipefail`、事前検査 `bash -n` |
| | パス表記 | POSIX 形式 (`/data/...`, `$HOME/...`, スラッシュ `/`) |
| | 環境変数設定 | `export VAR="val"` (`~/.bashrc`) |
| | 動作バイナリ | Linux aarch64 ネイティブ ELF (`rg`, `llm-mem`, `crabcode`, `python3`) |
| | **禁止事項** | ❌ `pwsh`, `powershell.exe`, `cmd.exe`, `.bat`, `.ps1` 直接実行, Windows パス (`C:\...`, `\`), `$env:...` |

---

### 【プロファイル 2】 `windows-pwsh` (Windows 11 / PowerShell 環境 ⚪)

Windows 11 デスクトップ / ワークステーション環境。PowerShell 7 によるスクリプト実行プロファイル。

| 区分 | 項目 | 値・仕様 |
| :--- | :--- | :--- |
| **ハードウェア** | デバイス / CPU | Windows 11 PC (x86_64 / amd64) |
| **OS / カーネル** | OS | Windows 11 Pro / Enterprise |
| **シェル & コマンド** | 主シェル | **PowerShell 7** (`pwsh.exe`) または **Windows PowerShell** (`powershell.exe`) |
| | コマンド実行方式 | `pwsh.exe -Command "<command>"` |
| | スクリプト要件 | PowerShell スクリプト (`.ps1`)、Pester テストフレームワーク |
| | パス表記 | Windows 形式 (`C:\Users\<user>\...`, `%USERPROFILE%`, バックスラッシュ `\`) |
| | 環境変数設定 | `$env:VAR = "val"` または `[System.Environment]::SetEnvironmentVariable(...)` |
| | 動作バイナリ | Windows PE 実行形式 (`rg.exe`, `llm-mem.exe`, `harness-lint.exe`, `python.exe`) |
| | **禁止事項** | ❌ POSIX 専用のハードコードパス (`/data/...`, `/usr/...`)、Termux 専用コマンド |

---

### 【プロファイル 3】 `linux-native-bash` (汎用 Linux 環境 ⚪)

一般的な Debian / Ubuntu / Arch などの Linux ディストリビューション向けプロファイル。

| 区分 | 項目 | 値・仕様 |
| :--- | :--- | :--- |
| **OS / カーネル** | OS | GNU/Linux (x86_64 または aarch64) |
| **シェル & コマンド** | 主シェル | **GNU Bash** (`/bin/bash` または `/usr/bin/bash`) |
| | パス表記 | POSIX 形式 (`/home/<user>/...`, `/tmp`) |
| | 一時ディレクトリ | `/tmp` または `$TMPDIR` |
| | 環境変数設定 | `export VAR="val"` (`~/.bashrc`) |
| | 動作バイナリ | Linux ELF ネイティブバイナリ |

---

## 🔀 プロファイル別コマンド対照表 (Cross-Profile Command Mapping)

エージェントがコマンドを記述・変換する際は、以下の対応表に基づき、`active_profile` に合致した構文を使用してください。

| 操作・概念 | `android-termux-bash` (現アクティブ) | `windows-pwsh` | `linux-native-bash` |
| :--- | :--- | :--- | :--- |
| **シェル実行** | `bash -c '...'` | `pwsh.exe -Command "..."` | `bash -c '...'` |
| **ホームパス** | `$HOME` (`/data/data/com.termux/files/home`) | `$env:USERPROFILE` (`C:\Users\<user>`) | `$HOME` (`/home/<user>`) |
| **環境変数参照** | `$VARIABLE` または `${VARIABLE}` | `$env:VARIABLE` | `$VARIABLE` または `${VARIABLE}` |
| **環境変数設定** | `export VAR="val"` | `$env:VAR = "val"` | `export VAR="val"` |
| **パス区切り** | スラッシュ `/` | バックスラッシュ `\` | スラッシュ `/` |
| **一時ディレクトリ** | `$TMPDIR` (`$PREFIX/tmp`) | `$env:TEMP` | `/tmp` |
| **ファイル検索** | `rg -n "pat" rules/*.lrf` | `rg.exe -n "pat" rules\*.lrf` | `rg -n "pat" rules/*.lrf` |
| **スクリプト実行** | `bash scripts/foo.sh` | `& .\scripts\foo.ps1` | `bash scripts/foo.sh` |
| **構文検査** | `bash -n scripts/foo.sh` | `Pester` / `[scriptblock]::Create(...)` | `bash -n scripts/foo.sh` |
| **メモリ CLI** | `llm-mem status` | `.\llm-mem.exe status` | `llm-mem status` |
| **一覧表示** | `ls -la`, `find` | `Get-ChildItem -Force` | `ls -la`, `find` |

---

## 🛠️ 各地における切り替え追従指針 (Guidance for Other Files)

1. **ルールファイル (`.lrf`)**:
   `BOOTSTRAP.lrf` の `R|host|*` は `hardware=HARDWARE.md relative to HarnessRoot, never process CWD; resolve=active_profile-in-HARDWARE.md` と規定し、静的な決め打ちではなく `active_profile` から動的にシェルとパス規約を解決します。
2. **ナビゲーションガイド (`README-forLLM.md`)**:
   `active_profile` を確認した上で、対応するプロファイルのコマンドブロックを選択して実行するようエージェントに指示します。
3. **スキル定義 (`SKILL.md` / `USAGE.md`)**:
   スニペット提示時は、`active_profile` に応じて Bash 記法または PowerShell 記法が選択されるよう、対照形式または条件分岐形式で記述します。
4. **プロファイルの切り替え手順**:
   環境を移行する際は、本書の `active_profile` を変更するだけで、後続のすべてのエージェントが自動的に新環境のシェル・パス規約に従います。
