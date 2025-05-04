#!/bin/sh

#!/bin/sh

set -e

BINARY_NAME="docky"
INSTALL_DIR="/usr/local/bin"
RAW_GITHUB_URL="https://raw.githubusercontent.com/BkycHblu-6oPwuK/compose/go/build/docky"

echo "Downloading $BINARY_NAME..."
curl -sL "$RAW_GITHUB_URL" -o "$INSTALL_DIR/$BINARY_NAME" || {
    echo "Failed to download binary"
    exit 1
}

chmod +x "$INSTALL_DIR/$BINARY_NAME"
echo "Successfully installed $BINARY_NAME"