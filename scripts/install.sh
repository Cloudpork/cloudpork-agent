#!/bin/bash

# CloudPork Agent Installation Script
# Usage: curl -sSL https://cli.cloudpork.com/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Constants
REPO="Cloudpork/cloudpork-agent"
BINARY_NAME="cloudpork"
INSTALL_DIR="/usr/local/bin"

# Print banner
print_banner() {
    echo -e "${MAGENTA}🐷 CloudPork Agent Installer${NC}"
    echo -e "${CYAN}Cut the pork from your cloud costs${NC}"
    echo ""
}

# Print colored message
log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

# Detect OS and architecture
detect_os() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $ARCH in
        x86_64) ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        armv7l) ARCH="arm" ;;
        *) 
            error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    
    case $OS in
        linux) OS="linux" ;;
        darwin) OS="darwin" ;;
        mingw*|cygwin*|msys*) OS="windows" ;;
        *)
            error "Unsupported OS: $OS"
            exit 1
            ;;
    esac
    
    log "Detected OS: $OS, Architecture: $ARCH"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."
    
    if ! command_exists curl && ! command_exists wget; then
        error "curl or wget is required but not installed."
        exit 1
    fi
    
    if ! command_exists tar; then
        error "tar is required but not installed."
        exit 1
    fi
    
    # Check if Claude Code is installed
    if ! command_exists claude; then
        warn "Claude Code CLI not found. CloudPork Agent requires Claude Code to function."
        echo ""
        echo -e "${CYAN}Install Claude Code first:${NC}"
        echo "  • macOS: brew install claude-ai/tap/claude"
        echo "  • Linux: curl -fsSL https://claude.ai/install.sh | sh"
        echo "  • Windows: Download from https://claude.ai/cli/download"
        echo "  • Web: https://claude.ai/cli"
        echo ""
        echo -e "${YELLOW}Continue installation anyway? (you can install Claude Code later) [y/N]${NC}"
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        log "Claude Code CLI found: $(claude --version 2>/dev/null || echo 'unknown version')"
    fi
}

# Get latest release
get_latest_release() {
    log "Fetching latest release information..."
    
    if command_exists curl; then
        LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command_exists wget; then
        LATEST_RELEASE=$(wget -qO- "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    fi
    
    if [ -z "$LATEST_RELEASE" ]; then
        error "Failed to get latest release information"
        exit 1
    fi
    
    log "Latest release: $LATEST_RELEASE"
}

# Download and install
install_binary() {
    local filename="${BINARY_NAME}_${LATEST_RELEASE#v}_${OS}_${ARCH}"
    
    if [ "$OS" = "windows" ]; then
        filename="${filename}.exe"
    fi
    
    local download_url="https://github.com/$REPO/releases/download/$LATEST_RELEASE/${filename}.tar.gz"
    local temp_dir=$(mktemp -d)
    local temp_file="$temp_dir/cloudpork.tar.gz"
    
    log "Downloading CloudPork Agent..."
    log "URL: $download_url"
    
    if command_exists curl; then
        curl -sL "$download_url" -o "$temp_file"
    elif command_exists wget; then
        wget -q "$download_url" -O "$temp_file"
    fi
    
    if [ ! -f "$temp_file" ]; then
        error "Failed to download CloudPork Agent"
        exit 1
    fi
    
    log "Extracting binary..."
    tar -xzf "$temp_file" -C "$temp_dir"
    
    # Find the binary (it might be in a subdirectory)
    local binary_path
    if [ -f "$temp_dir/$BINARY_NAME" ]; then
        binary_path="$temp_dir/$BINARY_NAME"
    elif [ -f "$temp_dir/${filename}" ]; then
        binary_path="$temp_dir/${filename}"
    else
        # Look for any executable file
        binary_path=$(find "$temp_dir" -name "$BINARY_NAME*" -type f -executable | head -n 1)
    fi
    
    if [ -z "$binary_path" ] || [ ! -f "$binary_path" ]; then
        error "Binary not found in downloaded package"
        exit 1
    fi
    
    log "Installing to $INSTALL_DIR..."
    
    # Check if we need sudo
    if [ ! -w "$INSTALL_DIR" ]; then
        if command_exists sudo; then
            sudo mv "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
            sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
        else
            error "No write permission to $INSTALL_DIR and sudo not available"
            exit 1
        fi
    else
        mv "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Cleanup
    rm -rf "$temp_dir"
    
    log "Installation complete!"
}

# Verify installation
verify_installation() {
    log "Verifying installation..."
    
    if command_exists "$BINARY_NAME"; then
        local version=$($BINARY_NAME version 2>/dev/null | head -n 1 || echo "unknown")
        log "CloudPork Agent installed successfully: $version"
    else
        error "Installation failed - binary not found in PATH"
        exit 1
    fi
}

# Print post-install instructions
print_instructions() {
    echo ""
    echo -e "${GREEN}✅ CloudPork Agent installed successfully!${NC}"
    echo ""
    echo -e "${CYAN}🚀 Getting Started:${NC}"
    echo "  1. Authenticate:    cloudpork auth login"
    echo "  2. Analyze project: cloudpork analyze"
    echo "  3. View help:       cloudpork --help"
    echo ""
    echo -e "${CYAN}📖 Documentation: https://docs.cloudpork.com${NC}"
    echo -e "${CYAN}🌐 Dashboard:      https://cloudpork.com/dashboard${NC}"
    echo ""
    
    if ! command_exists claude; then
        echo -e "${YELLOW}⚠️  Don't forget to install Claude Code CLI for full functionality!${NC}"
        echo ""
    fi
}

# Main installation flow
main() {
    print_banner
    detect_os
    check_prerequisites
    get_latest_release
    install_binary
    verify_installation
    print_instructions
}

# Handle interrupts
trap 'error "Installation interrupted"; exit 1' INT TERM

# Run main function
main "$@"