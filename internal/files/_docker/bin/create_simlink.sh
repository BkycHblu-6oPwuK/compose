#!/bin/bash

SIMLINK_FILE="/usr/simlinks"

if [ ! -f "$SIMLINK_FILE" ]; then
  echo "Файл $SIMLINK_FILE не найден"
  exit 1
fi

mapfile -t SIMLINK_LINES < "$SIMLINK_FILE"

for SIMLINK in "${SIMLINK_LINES[@]}"; do
  if [ -z "$SIMLINK" ]; then
    continue
  fi

  SOURCE=$(echo "$SIMLINK" | awk '{print $1}')
  TARGET=$(echo "$SIMLINK" | awk '{print $2}')

  if [ ! -d "$SOURCE" ]; then
    echo "Каталог $SOURCE не найден"
    continue
  fi

  if [ -e "$TARGET" ]; then
    echo "Цель $TARGET существует"
    continue
  fi

  if [ -L "$TARGET" ]; then
    LINK_TARGET=$(readlink "$TARGET")
    if [ "$LINK_TARGET" == "$SOURCE" ]; then
      echo "Симлинк $TARGET уже ведет на $SOURCE. Пропускаю создание."
      continue
    else
      echo "Симлинк $TARGET ведет на $LINK_TARGET, удаляю."
      unlink "$TARGET"
    fi
  fi
  
  echo "Создаю симлинк: $SOURCE -> $TARGET"
  ln -s "$SOURCE" "$TARGET"
done
