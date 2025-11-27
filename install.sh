#!/bin/sh
# forge installer - C++ Project Generator
# Usage: sh -c "$(curl -fsSL https://raw.githubusercontent.com/ozacod/forge/master/install.sh)"

set -e

REPO="ozacod/forge"
BINARY_NAME="forge"

# Determine install directory (prefer user-local, fallback to system with sudo)
get_install_dir() {
    # Check for user-specified directory
    if [ -n "$FORGE_INSTALL_DIR" ]; then
        echo "$FORGE_INSTALL_DIR"
        return
    fi
    
    # Prefer ~/.local/bin (no sudo needed)
    LOCAL_BIN="$HOME/.local/bin"
    if [ -d "$LOCAL_BIN" ] && [ -w "$LOCAL_BIN" ]; then
        echo "$LOCAL_BIN"
        return
    fi
    
    # Check if /usr/local/bin is writable
    if [ -w "/usr/local/bin" ]; then
        echo "/usr/local/bin"
        return
    fi
    
    # Default to ~/.local/bin (will be created)
    echo "$LOCAL_BIN"
}

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_banner() {
    printf "\n"
    printf "%b  ███████╗ ██████╗ ██████╗  ██████╗ ███████╗%b\n" "$CYAN" "$NC"
    printf "%b  ██╔════╝██╔═══██╗██╔══██╗██╔════╝ ██╔════╝%b\n" "$CYAN" "$NC"
    printf "%b  █████╗  ██║   ██║██████╔╝██║  ███╗█████╗  %b\n" "$CYAN" "$NC"
    printf "%b  ██╔══╝  ██║   ██║██╔══██╗██║   ██║██╔══╝  %b\n" "$CYAN" "$NC"
    printf "%b  ██║     ╚██████╔╝██║  ██║╚██████╔╝███████╗%b\n" "$CYAN" "$NC"
    printf "%b  ╚═╝      ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝%b\n" "$CYAN" "$NC"
    printf "\n"
    printf "  %bC++ Project Generator - Forge Your Code!%b\n" "$YELLOW" "$NC"
    printf "\n"
}

detect_os() {
    OS="$(uname -s)"
    case "$OS" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *)       echo "unknown" ;;
    esac
}

detect_arch() {
    ARCH="$(uname -m)"
    case "$ARCH" in
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *)            echo "unknown" ;;
    esac
}

check_dependencies() {
    if ! command -v curl > /dev/null 2>&1; then
        if ! command -v wget > /dev/null 2>&1; then
            printf "%bError: curl or wget is required%b\n" "$RED" "$NC"
            exit 1
        fi
        DOWNLOADER="wget"
    else
        DOWNLOADER="curl"
    fi
}

get_latest_version() {
    if [ "$DOWNLOADER" = "curl" ]; then
        VERSION=$(curl -sI "https://github.com/$REPO/releases/latest" | grep -i "location:" | sed 's/.*tag\///' | tr -d '\r\n')
    else
        VERSION=$(wget -qO- --server-response "https://github.com/$REPO/releases/latest" 2>&1 | grep -i "location:" | sed 's/.*tag\///' | tr -d '\r\n')
    fi
    
    if [ -z "$VERSION" ]; then
        VERSION="v1.0.0"
    fi
    echo "$VERSION"
}

download_binary() {
    OS=$1
    ARCH=$2
    VERSION=$3
    
    # Construct filename
    if [ "$OS" = "windows" ]; then
        FILENAME="${BINARY_NAME}-${OS}-${ARCH}.exe"
    else
        FILENAME="${BINARY_NAME}-${OS}-${ARCH}"
    fi
    
    URL="https://github.com/$REPO/releases/download/$VERSION/$FILENAME"
    
    printf "%b→ Downloading %s %s for %s/%s...%b\n" "$CYAN" "$BINARY_NAME" "$VERSION" "$OS" "$ARCH" "$NC" >&2
    
    # Create temp directory
    TMP_DIR=$(mktemp -d)
    TMP_FILE="$TMP_DIR/$BINARY_NAME"
    
    if [ "$DOWNLOADER" = "curl" ]; then
        if ! curl -fsSL "$URL" -o "$TMP_FILE"; then
            printf "%bError: Failed to download from %s%b\n" "$RED" "$URL" "$NC" >&2
            rm -rf "$TMP_DIR"
            exit 1
        fi
    else
        if ! wget -q "$URL" -O "$TMP_FILE"; then
            printf "%bError: Failed to download from %s%b\n" "$RED" "$URL" "$NC" >&2
            rm -rf "$TMP_DIR"
            exit 1
        fi
    fi
    
    printf '%s' "$TMP_FILE"
}

install_binary() {
    TMP_FILE=$1
    INSTALL_DIR=$2
    
    printf "%b→ Installing to %s...%b\n" "$CYAN" "$INSTALL_DIR" "$NC"
    
    # Make executable
    chmod +x "$TMP_FILE"
    
    # Create directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        mkdir -p "$INSTALL_DIR"
    fi
    
    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
    else
        printf "%b→ Requesting sudo access to install to %s%b\n" "$YELLOW" "$INSTALL_DIR" "$NC"
        printf "%b  (Use FORGE_INSTALL_DIR=~/.local/bin to avoid sudo)%b\n" "$YELLOW" "$NC"
        sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Cleanup temp directory
    rm -rf "$(dirname "$TMP_FILE")"
    
    # Check if install dir is in PATH
    case ":$PATH:" in
        *":$INSTALL_DIR:"*) ;;
        *)
            printf "\n%b⚠ Note: %s is not in your PATH%b\n" "$YELLOW" "$INSTALL_DIR" "$NC"
            printf "  Add it to your shell config:\n"
            printf "  %bexport PATH=\"%s:\$PATH\"%b\n" "$YELLOW" "$INSTALL_DIR" "$NC"
            ;;
    esac
}

verify_installation() {
    if command -v "$BINARY_NAME" > /dev/null 2>&1; then
        printf "\n"
        printf "%b✓ forge installed successfully!%b\n" "$GREEN" "$NC"
        printf "\n"
        printf "  Version: %s\n" "$("$BINARY_NAME" -v 2>/dev/null || echo "installed")"
        printf "  Location: %s\n" "$(command -v "$BINARY_NAME")"
        printf "\n"
        printf "%bGet started:%b\n" "$CYAN" "$NC"
        printf "  %bmkdir my_project && cd my_project%b\n" "$YELLOW" "$NC"
        printf "  %bforge init%b\n" "$YELLOW" "$NC"
        printf "  %bforge build%b\n" "$YELLOW" "$NC"
        printf "\n"
        printf "%bOr use a template:%b\n" "$CYAN" "$NC"
        printf "  %bforge init -t web-server%b\n" "$YELLOW" "$NC"
        printf "\n"
        printf "%bSee all commands:%b\n" "$CYAN" "$NC"
        printf "  %bforge --help%b\n" "$YELLOW" "$NC"
        printf "\n"
    else
        printf "%bError: Installation failed%b\n" "$RED" "$NC"
        exit 1
    fi
}

main() {
    print_banner
    
    # Detect platform
    OS=$(detect_os)
    ARCH=$(detect_arch)
    
    printf "%b→ Detected: %s/%s%b\n" "$CYAN" "$OS" "$ARCH" "$NC"
    
    if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
        printf "%bError: Unsupported platform: %s/%s%b\n" "$RED" "$OS" "$ARCH" "$NC"
        printf "Supported platforms:\n"
        printf "  - Linux (amd64, arm64)\n"
        printf "  - macOS (amd64, arm64)\n"
        printf "  - Windows (amd64)\n"
        exit 1
    fi
    
    # Check dependencies
    check_dependencies
    
    # Get latest version
    VERSION=$(get_latest_version)
    printf "%b→ Latest version: %s%b\n" "$CYAN" "$VERSION" "$NC"
    
    # Determine install directory
    INSTALL_DIR=$(get_install_dir)
    
    # Download binary
    TMP_FILE=$(download_binary "$OS" "$ARCH" "$VERSION")
    
    # Install binary
    install_binary "$TMP_FILE" "$INSTALL_DIR"
    
    # Verify installation
    verify_installation
}

# Run main
main

