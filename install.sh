#!/bin/sh
# cargo-cpp installer
# Usage: sh -c "$(curl -fsSL https://raw.githubusercontent.com/ozacod/cpp-repo-creator/master/install.sh)"

set -e

REPO="ozacod/cpp-repo-creator"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="cargo-cpp"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_banner() {
    echo ""
    echo "${CYAN}   ____                           ____            "
    echo "  / ___|__ _ _ __ __ _  ___      / ___|_ __  _ __  "
    echo " | |   / _\` | '__/ _\` |/ _ \\ ____| |   | '_ \\| '_ \\ "
    echo " | |__| (_| | | | (_| | (_) |____| |___| |_) | |_) |"
    echo "  \\____\\__,_|_|  \\__, |\\___/      \\____|.__/| .__/ "
    echo "                 |___/                 |_|   |_|    ${NC}"
    echo ""
    echo "  ${YELLOW}C++ Project Generator - Like Cargo for Rust!${NC}"
    echo ""
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
            echo "${RED}Error: curl or wget is required${NC}"
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
    
    echo "${CYAN}→ Downloading ${BINARY_NAME} ${VERSION} for ${OS}/${ARCH}...${NC}" >&2
    
    # Create temp directory
    TMP_DIR=$(mktemp -d)
    TMP_FILE="$TMP_DIR/$BINARY_NAME"
    
    if [ "$DOWNLOADER" = "curl" ]; then
        if ! curl -fsSL "$URL" -o "$TMP_FILE"; then
            echo "${RED}Error: Failed to download from $URL${NC}" >&2
            rm -rf "$TMP_DIR"
            exit 1
        fi
    else
        if ! wget -q "$URL" -O "$TMP_FILE"; then
            echo "${RED}Error: Failed to download from $URL${NC}" >&2
            rm -rf "$TMP_DIR"
            exit 1
        fi
    fi
    
    printf '%s' "$TMP_FILE"
}

install_binary() {
    TMP_FILE=$1
    
    echo "${CYAN}→ Installing to ${INSTALL_DIR}...${NC}"
    
    # Make executable
    chmod +x "$TMP_FILE"
    
    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
    else
        echo "${YELLOW}→ Requesting sudo access to install to ${INSTALL_DIR}${NC}"
        sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Cleanup temp directory
    rm -rf "$(dirname "$TMP_FILE")"
}

verify_installation() {
    if command -v "$BINARY_NAME" > /dev/null 2>&1; then
        echo ""
        echo "${GREEN}✓ cargo-cpp installed successfully!${NC}"
        echo ""
        echo "  Version: $("$BINARY_NAME" -v 2>/dev/null || echo "installed")"
        echo "  Location: $(command -v "$BINARY_NAME")"
        echo ""
        echo "${CYAN}Get started:${NC}"
        echo "  ${YELLOW}mkdir my_project && cd my_project${NC}"
        echo "  ${YELLOW}cargo-cpp init${NC}"
        echo "  ${YELLOW}cargo-cpp build${NC}"
        echo ""
        echo "${CYAN}Or use a template:${NC}"
        echo "  ${YELLOW}cargo-cpp init -t web-server${NC}"
        echo ""
        echo "${CYAN}See all commands:${NC}"
        echo "  ${YELLOW}cargo-cpp --help${NC}"
        echo ""
    else
        echo "${RED}Error: Installation failed${NC}"
        exit 1
    fi
}

main() {
    print_banner
    
    # Detect platform
    OS=$(detect_os)
    ARCH=$(detect_arch)
    
    echo "${CYAN}→ Detected: ${OS}/${ARCH}${NC}"
    
    if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
        echo "${RED}Error: Unsupported platform: ${OS}/${ARCH}${NC}"
        echo "Supported platforms:"
        echo "  - Linux (amd64, arm64)"
        echo "  - macOS (amd64, arm64)"
        echo "  - Windows (amd64)"
        exit 1
    fi
    
    # Check dependencies
    check_dependencies
    
    # Get latest version
    VERSION=$(get_latest_version)
    echo "${CYAN}→ Latest version: ${VERSION}${NC}"
    
    # Download binary
    TMP_FILE=$(download_binary "$OS" "$ARCH" "$VERSION")
    
    # Install binary
    install_binary "$TMP_FILE"
    
    # Verify installation
    verify_installation
}

# Run main
main

