#!/bin/sh

set -e

BINARY_NAME="docky"
INSTALL_DIR="/usr/local/bin"
REPO="BkycHblu-6oPwuK/compose"

echo "Установка $BINARY_NAME..."
curl https://raw.githubusercontent.com/BkycHblu-6oPwuK/compose/main/build/docky > "$INSTALL_DIR/$BINARY_NAME" || {
    echo "Ошибка загрузки $BINARY_NAME"
    exit 1
}

chmod +x "$INSTALL_DIR/$BINARY_NAME"
echo "Установка прошла успешно"
echo "Выполните команду $BINARY_NAME clean-cache для очистки кеша"