# Конфигурация под symfony

nginx + php (8.2, 8.3, 8.4) + mysql|mariadb|postgres|sqlite
а так же redis|memcached

## Шаги публикации docker-compose.yml

Шаги:

0. Выбор фреймворка - symfony
1. Выбрать версию php
2. Выберите базу данных - mysql|mariadb|postgres|sqlite
3. Выбрать версию базы данных
4. Выберите сервер для кеширования - redis|memcached|Пропуск
5. Установка symfony

После этого в директории где выполнялась команда появиться docker-compose.yml файл с настроенными сервисами.

В директории site уже будет развернут symfony проект

## bin/console

запускайте команды с помощью

```bash
docky symfony {arg}
```

## Cron

По умолчанию cron включен

Если необходимо добавить задания, то сделайте публикацию файлов с заданиями

```bash
docky publish --file cron_tasks
```

Запись заданий осуществляйте в:
- `${CONF_PATH}/app/cron/docky` - для пользователя сайта
- `${CONF_PATH}/app/cron/root` - для root пользователя

## Почта

msmtp клиент не устанавливается.
Используйте сервис mailhog для тестирования отправки почты.

- host=mailhog
- port=1025

панель доступна на - localhost:8025

## Node

Для публикации сервиса с node.js выполните команду и укажите точку входа для node.js в файле .env (переменная NODE_PATH)

```bash 
docky publish --service node
```

используйте команды ```docky npm``` чтобы запускать npm команды в контейнере

и добавьте настройки сервера в vite.config.js 

```
    server: {
        host: '0.0.0.0',
        port: 5173,
        open: false,
        cors: {
            origin: '*'
        },
        hmr: {
            host: 'localhost',
        },
    }
```

