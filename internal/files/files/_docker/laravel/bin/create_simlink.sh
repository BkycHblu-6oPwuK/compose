#!/bin/bash

process_symlinks_file() {
  local SIMLINK_FILE="$1"

  if [ ! -f "$SIMLINK_FILE" ]; then
    return 1
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

    if [ -L "$TARGET" ]; then
      LINK_TARGET=$(readlink "$TARGET")
      if [ "$LINK_TARGET" == "$SOURCE" ]; then
        echo "Симлинк $TARGET уже ведет на $SOURCE. Пропускаю создание."
        continue
      else
        echo "Симлинк $TARGET ведет на $LINK_TARGET, удаляю."
        unlink "$TARGET"
      fi
    elif [ -e "$TARGET" ]; then
      echo "Цель $TARGET существует и не является симлинком. Пропускаю."
      continue
    fi

    echo "Создаю симлинк: $SOURCE -> $TARGET"
    ln -s "$SOURCE" "$TARGET"
  done
}

process_symlinks_file "/usr/symlinks"
process_symlinks_file "/usr/symlinks_extra"