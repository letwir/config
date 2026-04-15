# Rustインストール (nightly)

winget install rustlang.Rustup
# Nightlyツールチェーンのインストール
rustup toolchain install nightly
# デフォルトをNightlyに切り替え
rustup default nightly

cd
curl https://raw.githubusercontent.com/letwir/config/refs/heads/main/rustup/cargo-windows.toml  -o \.cargo\config.toml

Set-ExecutionPolicy Unrestricted -Scope Process; iex (iwr "https://raw.githubusercontent.com/cargo-bins/cargo-binstall/main/install-from-binstall-release.ps1").Content
curl -sS https://starship.rs/install.sh | sh
cargo binstall git-delta bat sd vivid cargo-cache
cargo install fd-find lsd frs ripgrep du-dust hexyl choose lms starship --all-features
cargo cache -a



