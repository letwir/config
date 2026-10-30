# Knowledge — pi-install-20261025

- `install-pi.sh` now supports Debian/Ubuntu (`apt-get`), native Termux (`pkg`), and Fedora (`dnf`), with `--family` override and `--dry-run`.
- Shared prerequisites include Node/npm, Python/pip, Go, ripgrep, Git, curl, CA certificates, jq, and w3m. Rust uses distro packages plus `rustup default stable` on Debian/Fedora; Termux uses package-managed Rust because upstream rustup is unsupported there.
- Pillow package mappings cover the hatch-pet scripts. Chromium remains optional and is not installed by the script.
- Pi extensions retained: pi-lens, pi-web-access, pi-subagents. Added `@juicesharp/rpiv-todo` and `@juicesharp/rpiv-ask-user-question`.
- Package availability and Rustup support may vary by distribution release; no host installation was performed.
