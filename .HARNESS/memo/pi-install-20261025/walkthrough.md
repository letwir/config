# Walkthrough — pi-install-20261025

1. Run `sh -n install-pi.sh` to validate POSIX shell syntax.
2. Run `sh install-pi.sh --help` to inspect the supported options.
3. Review mappings without installation using `sh install-pi.sh --family debian --dry-run`, then repeat with `termux` and `fedora`.
4. Confirm dry-run output does not contain `pi install` calls.
5. Run the script normally only on a supported target host when installation is desired; package-manager operations and Pi extension installs are intentionally not part of local verification.
