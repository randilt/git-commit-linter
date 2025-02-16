#!/bin/bash

set -e

# Config
GITHUB_USER="randilt"
REPO_NAME="git-commit-linter"
BINARY_NAME="git-commit-linter"
INSTALL_DIR="/usr/local/bin"
TMP_DIR=$(mktemp -d)
LATEST_RELEASE_URL="https://api.github.com/repos/$GITHUB_USER/$REPO_NAME/releases/latest"

# Color output
RED='\033[0;31m'
YELLOW='\033[0;33m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# Print step message
print_step() {
    echo -e "${BLUE}==>${NC} $1"
}

# Print success message
print_success() {
    echo -e "${GREEN}==>${NC} $1"
}

# Print error message
print_error() {
    echo -e "${RED}Error:${NC} $1"
}

# Print warn message
print_warn() {
    echo -e "${YELLOW}Warning:${NC} $1"
}

# Cleanup on exit
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

# Warn that sudo access may be required in middle of script
if [ "$EUID" -ne 0 ]; then
    print_warn "This installation may require sudo access to install $BINARY_NAME to $INSTALL_DIR"
fi

# Check if curl or wget is available
if command -v curl >/dev/null 2>&1; then
    DOWNLOAD_CMD="curl -L"
    API_CMD="curl -s"
elif command -v wget >/dev/null 2>&1; then
    DOWNLOAD_CMD="wget -O -"
    API_CMD="wget -qO-"
else
    print_error "CURL or WGET not found. Please install either and try again."
    exit 1
fi

# Detect OS and architecture
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
    "Darwin")
        OS="Darwin"
        ;;
    "Linux")
        OS="Linux"
        ;;
    *)
        print_error "Unsupported operating system: $OS"
        exit 1
        ;;
esac

case "$ARCH" in
    "x86_64"|"amd64")
        ARCH="x86_64"
        ;;
    "arm64"|"aarch64")
        ARCH="arm64"
        ;;
    *)
        print_error "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

print_step "Detected OS: $OS, Architecture: $ARCH"

# Get the latest release download URL
print_step "Fetching latest release information..."
RELEASE_DATA=$($API_CMD "$LATEST_RELEASE_URL")
DOWNLOAD_URL=$(echo "$RELEASE_DATA" | grep -o "https://github.com/$GITHUB_USER/$REPO_NAME/releases/download/[^\"]*${OS}_${ARCH}.tar.gz")

if [ -z "$DOWNLOAD_URL" ]; then
    print_error "Could not find download URL for your system"
    exit 1
fi

# Download and extract
print_step "Downloading latest release..."
cd "$TMP_DIR"
$DOWNLOAD_CMD "$DOWNLOAD_URL" | tar xz


# print_step "Debugging directory contents..."
# echo "Content of TMP_DIR ($TMP_DIR):"
# ls -la "$TMP_DIR"

EXTRACTED_DIR="$TMP_DIR/${BINARY_NAME}_${OS}_${ARCH}"
# if [ -d "$EXTRACTED_DIR" ]; then
#     echo "Content of $EXTRACTED_DIR:"
#     ls -la "$EXTRACTED_DIR"
# else
#     print_error "Repository directory not found in $TMP_DIR"
#     exit 1
# fi

# if [ ! -f "$EXTRACTED_DIR/$BINARY_NAME" ]; then
#     print_error "Binary not found in the downloaded archive"
#     exit 1
# fi


# Install binary
print_step "Installing $BINARY_NAME to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$EXTRACTED_DIR/$BINARY_NAME" "$INSTALL_DIR/"
else
    sudo mv "$EXTRACTED_DIR/$BINARY_NAME" "$INSTALL_DIR/"
fi

# Make binary executable
if [ -w "$INSTALL_DIR/$BINARY_NAME" ]; then
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

print_success "$BINARY_NAME has been installed successfully!"
print_step "You can now use it by running: $BINARY_NAME"

# Verify installation
if command -v $BINARY_NAME >/dev/null 2>&1; then
    print_success "Installation verified successfully!"
    $BINARY_NAME version
else
    print_error "Installation seems to have failed. Please check if $INSTALL_DIR is in your PATH"
    exit 1
fi