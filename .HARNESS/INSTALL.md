# ハーネス環境セットアップ & プロファイル切り替え手順書 (INSTALL.md)

本ドキュメントは、本リポジトリ／ハーネス環境を異なるホスト環境（Android-Termux、Windows、汎用Linux）へ展開または移行する際に、実行環境を自動・手動で検出し、`HARDWARE.md` の環境セレクターを切り替えるための標準手順書ですわ。

---

## 🎯 必須手順概要

ハーネスの環境切り替えは以下の2ステップで完結します：
1. **①: 今の環境を探索。** （ホストOS、アーキテクチャ、シェル、ツールチェインの特定）
2. **②: active_profileを切り替え。** （`HARDWARE.md` のプロファイルセレクターを対象環境へ更新）

---

## ①: 今の環境を探索。

現在のホスト環境が「Android (Termux)」「Windows」「汎用 Linux」のいずれであるかを探索・判定します。

### 1. 自動判別ワンライナー

#### 【Bash / POSIX シェル環境の場合】
以下のコマンドを実行すると、現在の環境と推奨プロファイルが表示されます：
```bash
if [ -n "${TERMUX_VERSION:-}" ] || [[ "${PREFIX:-}" == *"com.termux"* ]]; then
  echo "検出結果: Android Termux (${HOSTTYPE:-$(uname -m)}) -> 推奨プロファイル: android-termux-bash"
elif [[ "$(uname -s)" =~ (CYGWIN|MINGW|MSYS) ]]; then
  echo "検出結果: Windows (POSIX互換レイヤ) -> 推奨プロファイル: windows-pwsh"
elif [ "$(uname -s)" = "Linux" ]; then
  echo "検出結果: Linux Native ($(uname -m)) -> 推奨プロファイル: linux-native-bash"
else
  echo "検出結果: 未知の環境 ($(uname -a))"
fi
```

#### 【PowerShell (Windows) 環境の場合】
PowerShell コンソールで以下のコマンドを実行します：
```powershell
if ($PSVersionTable.PSEdition -eq "Core" -and $IsWindows -or $env:OS -like "*Windows*") {
  Write-Host "検出結果: Windows ($([System.Environment]::OSVersion.Platform)) -> 推奨プロファイル: windows-pwsh"
} elseif ($IsLinux) {
  Write-Host "検出結果: Linux Native ($([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture)) -> 推奨プロファイル: linux-native-bash"
} else {
  Write-Host "検出結果: Windows PowerShell / その他 -> 推奨プロファイル: windows-pwsh"
}
```

---

### 2. 手動詳細探索コマンド

環境を手動で精査する場合は、以下のコマンド群でシステム仕様を確認してください：

| 探索項目 | Bash (Termux / Linux) | PowerShell (Windows) |
| :--- | :--- | :--- |
| **OS / カーネル** | `uname -a` | `[System.Environment]::OSVersion` |
| **CPU アーキテクチャ** | `uname -m` | `$env:PROCESSOR_ARCHITECTURE` |
| **アクティブシェル** | `echo "$SHELL ($0)"` | `$PSVersionTable.PSVersion` |
| **ホームディレクトリ** | `echo "$HOME"` | `$env:USERPROFILE` |
| **一時ディレクトリ** | `echo "${TMPDIR:-/tmp}"` | `$env:TEMP` |
| **Termux 特有環境** | `echo "$PREFIX"` / `echo "$TERMUX_VERSION"` | (N/A) |

#### 判定基準テーブル:
- **`android-termux-bash`**: `$PREFIX` が `/data/data/com.termux/...` を含み、`uname -o` が `Android` または Linux aarch64 上の Termux の場合。一時ディレクトリは必ず `$TMPDIR` (`$PREFIX/tmp`) を使用し、`/tmp` へのアクセスは禁止。
- **`windows-pwsh`**: OSが Windows、PowerShell (`pwsh` または `powershell.exe`) が主シェル、パス区切りにドライブ文字 (`C:\...`) やバックスラッシュを使用する場合。
- **`linux-native-bash`**: 標準的な Debian/Ubuntu/Fedora/Arch Linux 環境、シェルが GNU Bash、一時ディレクトリが `/tmp` の場合。

---

## ②: active_profileを切り替え。

探索結果に基づき、`HARDWARE.md`（SSOT）の `active_profile` を切り替えます。

### 1. 切り替え対象プロファイル名一覧

| プロファイル名 | 対象環境 | 主シェル | パス形式 | コマンド実行規約 |
| :--- | :--- | :--- | :--- | :--- |
| **`android-termux-bash`** | Sony Xperia 1 VI 等 Android Termux | GNU Bash | POSIX (`/data/...`, スラッシュ `/`) | `bash -c '<cmd>'`, aarch64 ELFバイナリ, `$TMPDIR` 使用 |
| **`windows-pwsh`** | Windows 11 PC / デスクトップ | PowerShell 7 / Win PS | Windows (`C:\...`, `\`) | `pwsh -Command "<cmd>"`, PE `.exe` バイナリ |
| **`linux-native-bash`** | Ubuntu / Debian サーバー・WSL | GNU Bash | POSIX (`/home/...`, スラッシュ `/`) | `bash -c '<cmd>'`, x86_64/arm64 ELFバイナリ, `/tmp` 使用 |

---

### 2. 切り替えの実行方法

#### 【方法 A: Bash からワンライナーで切り替える場合 (sed)】
```bash
# 例: android-termux-bash に切り替える場合
sed -i -E 's/(active_profile: ")[^"]+(")/\1android-termux-bash\2/' ~/.harness/HARDWARE.md
sed -i -E 's/(\*\*`active_profile`\*\* \| \*\*`)[^`]+(`\*\*)/\1android-termux-bash\2/' ~/.harness/HARDWARE.md

# 例: windows-pwsh に切り替える場合
sed -i -E 's/(active_profile: ")[^"]+(")/\1windows-pwsh\2/' ~/.harness/HARDWARE.md
sed -i -E 's/(\*\*`active_profile`\*\* \| \*\*`)[^`]+(`\*\*)/\1windows-pwsh\2/' ~/.harness/HARDWARE.md

# 例: linux-native-bash に切り替える場合
sed -i -E 's/(active_profile: ")[^"]+(")/\1linux-native-bash\2/' ~/.harness/HARDWARE.md
sed -i -E 's/(\*\*`active_profile`\*\* \| \*\*`)[^`]+(`\*\*)/\1linux-native-bash\2/' ~/.harness/HARDWARE.md
```

#### 【方法 B: PowerShell からワンライナーで切り替える場合】
```powershell
# 例: windows-pwsh に切り替える場合
$hwFile = "$env:USERPROFILE\.harness\HARDWARE.md"
(Get-Content $hwFile) -replace 'active_profile: "[^"]+"', 'active_profile: "windows-pwsh"' `
                      -replace '(\*\*`active_profile`\*\* \| \*\*`)[^`]+(`\*\*)', '$1windows-pwsh$2' | Set-Content $hwFile
```

#### 【方法 C: 手動エディタ編集】
`~/.harness/HARDWARE.md`（およびリポジトリ原本 `config/.HARNESS/HARDWARE.md`）の以下のブロックを開き、値を直接書き換えます：

```markdown
| 項目 | 現在の設定値 | 利用可能なプロファイル |
| :--- | :--- | :--- |
| **`active_profile`** | **`android-termux-bash`** | `android-termux-bash`, `windows-pwsh`, `linux-native-bash` |

```yaml
# Profile Switcher Directive
hardware_environment:
  active_profile: "android-termux-bash"
  fallback_profile: "linux-native-bash"
```
```

---

## ③: セットアップと検証 (Setup & Verification)

切り替え完了後、ハーネスの正常動作を確認します。

### 1. プレフライト検査の実行
ルールとスクリプトの整合性を確認します：
```bash
# Bash 環境 (android-termux-bash / linux-native-bash)
bash ~/.harness/scripts/invoke-rule-preflight.sh --task change
```
```powershell
# Windows 環境 (windows-pwsh)
pwsh ~/.harness/scripts/invoke-rule-preflight.ps1 -Task change
```

### 2. スキルバイナリのパス解決規則
各スキルの `SKILL.md` では、プラットフォーム非依存の統合解決関数 `first-existing(...)` が定義されています：
```lrf
bin=first-existing(~/.harness/skills/<skill>/<bin>, ~/.harness/skills/<skill>/<bin>.exe, ./<bin>, ./<bin>.exe)
```
これにより、ELF 形式バイナリ（拡張子なし）と Windows PE 形式（`.exe`）のいずれが存在しても、自動的に現在のホスト環境に存在するバイナリが選択されます。

### 3. エージェント記憶 (llm-mem) の接続確認
```bash
# Bash
llm-mem status
```
```powershell
# PowerShell
& llm-mem.exe status
```
Status が `OK` で PostgreSQL ノードと正常通信できていることを確認してください。
