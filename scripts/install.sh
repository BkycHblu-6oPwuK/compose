#!/bin/sh

set -e

BINARY_NAME="docky"
INSTALL_DIR="/usr/local/bin"

echo "Установка $BINARY_NAME..."
TMP_FILE="$(mktemp)"
curl -sSL https://raw.githubusercontent.com/BkycHblu-6oPwuK/docky/main/bin/docky -o "$TMP_FILE" || {
    echo "Ошибка загрузки docky"
    exit 1
}

chmod +x "$TMP_FILE"
sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
echo "Установка прошла успешно"
echo "Выполните команду $BINARY_NAME clean-cache для очистки кеша"