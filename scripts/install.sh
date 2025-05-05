#!/bin/sh

set -e

BINARY_NAME="docky"
INSTALL_DIR="/usr/local/bin"
REPO="BkycHblu-6oPwuK/compose"
BRANCH="main"
GITHUB_URL="https://github.com/$REPO/raw/$BRANCH/build/$BINARY_NAME"

echo "Downloading $BINARY_NAME..."
curl https://raw.githubusercontent.com/BkycHblu-6oPwuK/compose/go/build/docky > "$INSTALL_DIR/$BINARY_NAME" || {
    echo "Failed to download binary"
    exit 1
}

chmod +x "$INSTALL_DIR/$BINARY_NAME"
$BINARY_NAME clean-cache
echo "Successfully installed $BINARY_NAME to $INSTALL_DIR"