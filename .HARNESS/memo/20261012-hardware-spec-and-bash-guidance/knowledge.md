# knowledge / 20261012-hardware-spec-and-bash-guidance

2026-10-12: ホスト環境は Android 16 (Linux 6.12.38-android16 aarch64), Sony XQ-GE44 (Snapdragon 8 Gen 3, 16GB RAM), Termux 0.118.3。既定シェルは GNU Bash (`/data/data/com.termux/files/usr/bin/bash`)。
PowerShell (`pwsh`, `powershell.exe`) および `cmd.exe` は本環境には存在せず、直接実行できない。
また、Android の制約上ルート直下の `/tmp` は権限拒否となるため、一時出力先には必ず `$TMPDIR` (`$PREFIX/tmp`) を使用する必要がある。
LLM向けルールやスキルにおける古いWindows/PowerShell表記（`$env:...`, `rules\*.lrf`, `& ...`）は、エージェントがPowerShellコマンドを生成・実行しようとする原因となるため、SSOTとしての `.harness/HARDWARE.md` を起点に、`AGENTS.md`, `BOOTSTRAP.lrf`, `README-forLLM.md`, `SKILL.md` を環境適応シェル（Bash）へ統一的に誘導することが肝要である。
