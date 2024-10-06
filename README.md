# Bitrix docker-compose

Настроенная обертка над docker-compose для локальных проектов bitrix 

nginx + php (7.4, 8.2) + mysql (5.7, 8.0) + node 20 версии

Установка:

```bash
composer require beeralex/compose:dev-main
```

## Установка docker-compose.yml

Команды выполнять в директории где находится composer.json!!!

Для установки docker-compose.yml выполнить команду

```bash
./vendor/bin/compose install
```

Или можно создать алиас

```bash
alias compose='sh $([ -f compose ] && echo compose || echo vendor/bin/compose)'
```

Тогда команда install будет выполняться слудующим образом:

```bash
compose install
```

Шаги установки:

1. Путь до директории сайта указывается относительно текущей директории (./). Можно оставить пустым и вмонтироваться в контейнеры будет все из текущей директории. Лучше выносить за сайт и указывать например - site, а в директории site уже разворачивать сайт.
2. Выбрать версию php - 7.4 или 8.2
3. Выбрать версию mysql - 5.7 или 8.0
4. Устанавливать ли node.js - Y или N, если установите Y, то будет создан сервис с node js 20 версии
    1. Если будет устанавливаться node.js, то нужно указать корневую директорию для него, то есть директория содержащая файл package.json. Путь указывается относительно корня сайта - local/js/vite или пустое поле если package.json в корне сайта.
5. Устанавливать ли sphinx - Y или N, если установите Y, то будет создан сервис с shphinx версии 2.2.11


После этого в директории где выполнялась команда появиться docker-compose.yml файл с настроенными сервисами.

## Публикация докерфайлов и файлов конфигурации

Если вам необходимо внести изменения в докерфайлы или файлы конфигурации, или добавить что то свое, то используйте команду:

```bash
compose publish
```

Эта команда перенесет всю конфигурацию в текущую директорию и изменит пути до файлов в docker-compose.yml, т.к. теперь все будет браться из ./_docker

## SSL сертификаты для nginx

Сертификаты и ключ копируются в контейнер из /_docker/nginx/certs/ и запись о них уже добавлена в nginx.conf.

Сейчас там используются временные самописные сертификаты которые будут действительны до ~2051 года.

Сервер одинаково настроен на работу как по http, так и по https.

Если нужно сгенерировать свои сертификат и ключ в текущую директорию:
```bash 
openssl req -x509 -nodes -days 9999 -newkey rsa:2048 -keyout privkey.pem -out fullchain.pem
```

Так как пакет расчитан для локальной разработки, здесь не будет рассматриваться letscrypt

## nginx в php контейнере

Для работы с сокетами в php кентейнер был установлен nginx который проксирует запросы на контейнер nginx.

nginx.conf для контейнера php лежит в _docker/app/nginx.conf

## Настройка cron

По умолчанию установка cron выключена. Для включания опубликуйте докерфайлы и файлы конфигурации если еще этого не делали

```bash
compose publish
```

Перейдите в "_docker/app/php-{your-version}/Dockerfile" и раскомментируйте строки начинающиеся с - "#cron"

Запись заданий осуществляйте в "_docker/app/crontab.txt"

выполните команду

```bash
compose build
```
и запустите контейнеры, cron должен начать работать

## Туннелирование локального сайта

Для туннелирования используется Expose (https://github.com/beyondcode/expose) и для того чтобы поделить вашим локальным сайтом выполните команду:

```bash
compose share
```

Сайт будет доступен 1 час, после этого команду можно выполнить заново.

доступные флаги expose для проброса через команду compose share:

- `--auth`
- `--server`
- `--subdomain`
- `--domain`
- `--server-host`
- `--server-port`

Документация expose - https://expose.dev/docs/introduction

## Sphinx (поисковая система)

sphinx (версия 2.2.11) является сервисом в docker-compose.yml (добавляется при установке) и собирается на основе Dockerfile из _docker/sphinx/Dockerfile, где так же лежит и файл конфигурации sphinx.conf.

После запуска контейнеров можно подключаться к sphinx:

```
sphinx:9306 - протокол MySql
sphinx:9312 - стандартный протокой
```

## SSH ключи

ssh ключи пробрасываются с помощью 'secrets' и вам достаточно расскомментировать соответствующие строки (начинающиеся с "#ssh") в следующих файлах:

- `docker-compose.yml`
- `_docker/app/php-{your-version}/Dockerfile`

## xdebug

для работы xdebug раскомментируйте строки начинающиеся с "#xdebug", в следующих файлах:

- `docker-compose.yml`
- `_docker/app/php-{your-version}/Dockerfile`

## Описание всех доступных команд

- `install` - Создание docker-compose.yml в текущей директории
```bash
compose install
```
- `publish` - Публикация файлов конфигурации в текущей директории
```bash
compose publish
```
- `share` - Позволяет сделать сайт доступным из интернета
```bash
compose share
```
- `php` - Выполнение команды php в контейнере с php
```bash
compose php -v
```
- `composer` - Выполнение команды composer в контейнере с php
```bash
compose composer install
```
- `npm` - Выполнение команды npm в контейнере с node js если он был установлен ()
```bash
compose npm install
```
- `И все дефолтные команды docker-compose` - Выполнение любой команды docker-compose
```bash
compose up -d
compose down
compose build
```

## Пользователи в контейнерах

- `appuser` - в контейнере с php (service app)
- `nodeuser` - в контейнере с node js (service node)

## Возможные проблемы

- `bash: compose: Permission denied` - выполнить команду
```bash
chmod +x ./vendor/beeralex/compose/src/bin/compose
```
- `compose: not found или compose: /bin/bash^M: bad interpreter` - выполнить команду
```bash
sed -i 's/\r$//' ./vendor/beeralex/compose/src/bin/compose
```