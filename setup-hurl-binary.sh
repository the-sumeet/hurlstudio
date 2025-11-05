#!/bin/bash

# Setup script for embedding Hurl binary

set -e

echo "🔍 Checking for Hurl installation..."

# Check if hurl is installed
if ! command -v hurl &> /dev/null; then
    echo "Hurl is not installed on your system"
    echo ""
    echo "Please install it first:"
    echo "  macOS:   brew install hurl"
    echo "  Linux:   apt install hurl (or cargo install hurl)"
    echo "  Windows: scoop install hurl"
    exit 1
fi

echo "Found Hurl at: $(which hurl)"
echo "   Version: $(hurl --version)"
echo ""

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize architecture
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
fi

TARGET_DIR="binaries/${OS}-${ARCH}"

echo "Setting up binary for: ${OS}-${ARCH}"
echo ""

# Create directory
mkdir -p "$TARGET_DIR"

# Copy binary
HURL_PATH=$(which hurl)
cp "$HURL_PATH" "$TARGET_DIR/"

echo "Binary copied to: $TARGET_DIR/hurl"
echo ""
echo "Setup complete!"
echo ""
echo "You can now build your app with the embedded Hurl binary."
echo "To support other platforms, download binaries from:"
echo "  https://github.com/Orange-OpenSource/hurl/releases"
