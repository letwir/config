#!/bin/sh

# Install Pi's system prerequisites. No downloaded installer scripts are run.
set -eu

usage() {
  cat <<'EOF'
Usage: install-pi.sh [--family auto|debian|termux|fedora] [--dry-run] [--help]

Detect Debian/Ubuntu, native Termux, or Fedora and install their prerequisite
packages. --family overrides detection (useful for reviewing package mappings).
--dry-run prints the selected package-manager commands without running them.
Pi extension installs are retained and are not run in dry-run mode.
Browser tooling is optional: install Chromium separately if available for your
platform. Fedora package availability and Termux Rust support can vary.
EOF
}

family=auto
dry_run=no
while [ "$#" -gt 0 ]; do
  case "$1" in
    --family)
      [ "$#" -ge 2 ] || { echo "--family requires a value" >&2; exit 2; }
      family=$2; shift 2 ;;
    --family=*) family=${1#*=}; shift ;;
    --dry-run) dry_run=yes; shift ;;
    --help|-h) usage; exit 0 ;;
    *) echo "Unknown option: $1" >&2; usage >&2; exit 2 ;;
  esac
done

case "$family" in auto|debian|termux|fedora) ;; *) echo "Unsupported family: $family" >&2; exit 2 ;; esac

if [ "$family" = auto ]; then
  if [ -n "${TERMUX_VERSION:-}" ] || { [ -n "${PREFIX:-}" ] && [ -d "$PREFIX" ] && [ "$(basename "$PREFIX")" = usr ] && [ -e "$PREFIX/etc/termux" ]; }; then
    family=termux
  else
    [ -r /etc/os-release ] || { echo "Cannot detect OS: /etc/os-release is unavailable; use --family" >&2; exit 1; }
    # shellcheck disable=SC1091
    . /etc/os-release
    case "${ID:-}" in
      debian|ubuntu) family=debian ;;
      fedora) family=fedora ;;
      *) echo "Unsupported OS ID '${ID:-unknown}'; supported: Debian/Ubuntu, native Termux, Fedora" >&2; exit 1 ;;
    esac
  fi
fi

case "$family" in
  debian)
    manager=apt-get
    packages='build-essential pkg-config ca-certificates curl git jq nodejs npm python3 python3-pip python3-pil golang-go ripgrep w3m rustup'
    refresh='apt-get update'
    install='apt-get install -y'
    rust_setup='rustup default stable' 
    ;;
  termux)
    manager=pkg
    packages='clang make pkg-config ca-certificates curl git jq nodejs python python-pip python-pillow golang ripgrep w3m rust'
    refresh=''
    install='pkg install -y'
    rust_setup='' 
    ;;
  fedora)
    manager=dnf
    packages='gcc gcc-c++ make pkgconf-pkg-config ca-certificates curl git jq nodejs npm python3 python3-pip python3-pillow golang ripgrep w3m rustup'
    refresh=''
    install='dnf install -y'
    rust_setup='rustup default stable' 
    ;;
esac

printf 'Family: %s\nPackages: %s\n' "$family" "$packages"
if [ "$family" = termux ]; then
  echo "Termux Rust: package-managed rust is provided where available; upstream rustup is not supported/recommended on Termux."
fi
if [ "$family" = debian ] || [ "$family" = fedora ]; then
  if [ "$(id -u)" -eq 0 ]; then
    privilege=''
  else
    command -v sudo >/dev/null 2>&1 || { echo "sudo is required for package installation as non-root" >&2; exit 1; }
    privilege='sudo'
  fi
else
  privilege=''
fi
if [ "$dry_run" = yes ]; then
  [ -z "$refresh" ] || printf 'DRY RUN: %s%s\n' "$privilege" "${refresh:+ $refresh}"
  printf 'DRY RUN: %s%s %s\n' "$privilege" "${install:+ $install}" "$packages"
  [ -z "$rust_setup" ] || printf 'DRY RUN: %s (network may be used to install stable toolchain)\n' "$rust_setup"
else
  command -v "$manager" >/dev/null 2>&1 || { echo "Required package manager '$manager' not found for $family" >&2; exit 1; }
  if [ -n "$refresh" ]; then
    # Intentional fixed package-manager command.
    # shellcheck disable=SC2086
    $privilege $refresh
  fi
  # Intentional word splitting of the fixed package list.
  # shellcheck disable=SC2086
  $privilege $install $packages
  if [ -n "$rust_setup" ]; then
    command -v rustup >/dev/null 2>&1 || { echo "rustup package did not provide rustup; install Rust manually" >&2; exit 1; }
    rustup default stable
  fi
  pi install npm:pi-lens
  pi install npm:pi-web-access
  pi install npm:pi-subagents
  pi install npm:@juicesharp/rpiv-todo
  pi install npm:@juicesharp/rpiv-ask-user-question
fi
