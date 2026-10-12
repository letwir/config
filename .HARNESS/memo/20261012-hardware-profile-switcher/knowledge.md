# knowledge / 20261012-hardware-profile-switcher

2026-10-12: ハードウェアおよびシェル環境の定義を「静的なOS決め打ち」にするのではなく、`HARDWARE.md` の冒頭に機械可読な `active_profile` セレクターを配置することで、マルチプラットフォーム（Windows 11 / PowerShell 7、Android-Termux / Bash、汎用 Linux / Bash）を1箇所の値書き換えだけで全体連動・切り替え可能となった。
LRFルール側でも `resolve=active_profile-in-HARDWARE.md; adapt-shell-and-paths=active_profile` と宣言することで、ルールの静的スキーマを変更することなく実行時の環境適応性を最大化できる。
各スキル・ドキュメントでも `active_profile` 別の記法対照表（Cross-Profile Command Mapping）を参照させることで、エージェントが現在のアクティブ環境に応じた正しい構文を選択できるようになる。
