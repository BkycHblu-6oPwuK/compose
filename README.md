# Bitrix docker-compose

Настроенная обертка над docker-compose для локальных проектов bitrix 

nginx + php (7.4, 8.2, 8.3, 8.4) + mysql (5.7, 8.0) + node 20 версии

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
2. Выбрать версию php
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

Более подробно - https://github.com/BkycHblu-6oPwuK/compose/tree/main/src/_docker/nginx

## nginx в php контейнере

Для работы с сокетами в php кентейнер был установлен nginx который проксирует запросы на контейнер nginx.

nginx.conf для контейнера php лежит в _docker/app/nginx.conf

## cron

По умолчанию cron включен и выполняется задание на запуск файла ```/var/www/bitrix/modules/main/tools/cron_events.php```

Если необходимо добавить задания, то сделайте публикацию докерфайлов и файлов конфигурации

```bash
compose publish
```

Запись заданий осуществляйте в:
- `_docker/app/cron/appuser.txt` - для пользователя сайта
- `_docker/app/cron/root.txt` - для root пользователя

выполните команду

```bash
compose build
```

## Почта

Для отправки почты настроен SMTP клиент  - ```msmtp```

Для завершения настройки вам необходимо добавить вашу почту (в поля user и from) и пароль в файл:

- `_docker/app/msmtprc`

По умолчанию в этом файле заготовка под почту яндекса, для других сервисов просто сделайте новый блок аккаунта на основе yandex и замените имя аккаунта в строке ```account default```, тогда ваш аккаунт будет по умолчанию использоваться при отправке почты.

Если почта не отправляется или в проверке системы написано что почта не работает, то проверьте логи ```msmtp``` в контейнере, которые находятся в файле ```/home/appuser/msmtp.log```. Вероятнее всего произошла ошибка авторизации или почтовый сервис отклюнил отправку из-за подозрений в спаме.

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

## Пакет pm2 в node контейнере

В dockerfile node устанавливается пакет pm2, для работы с ним используйте команду

```bash 
compose pm2 {arg}
```

Запустить сервер node js можно:

1. командой pm2 - ```compose pm2 start server.js```
2. настроить на запуск при запуске контейнеров - для этого нужно раскомментировать строчку ```command``` в docker-compose.yml в секции с сервисом node, при этом command строчкой выше можно удалить. Так же в этой команде проверьте путь до файла с сервером

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
- `npm` - Выполнение команды npm в контейнере с node js если он был установлен
```bash
compose npm install
```
- `pm2` - Выполнение команды pm2 в контейнере с node js если он был установлен. Команда принимает такие же аргументы как и оригинальная pm2
```bash
compose pm2 {arg}
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

## Настройка Redis

Для настройки redis выполните следующие шаги:

1. Добавьте сервис сервера redis в docker-compose.yml
```yml
  redis:
    image: redis
    ports:
      - "6379:6379"
    command: ["redis-server", "--appendonly", "yes"]
    volumes:
      - redis_data:/data
    container_name: redis
    networks:
      - compose
volumes:
  ...
  redis_data:
```
2. В dockerfile вашего app сервиса (_docker/app/php-{your-version}/Dockerfile) добавьте установку php модуля redis и igbinary
```dockerfile
RUN pecl install igbinary && \
    pecl install -D 'enable-redis-igbinary="yes"' redis && \
    docker-php-ext-enable igbinary redis
```
3. Выполните build.

## Настройка memcached

Для настройки memcached выполните следующие шаги:

1. Добавьте сервис сервера memcached в docker-compose.yml
```yml
  memcached:
    image: memcached
    ports:
      - "11211:11211"
    container_name: memcached
    networks:
      - compose
```
2. В dockerfile вашего app сервиса (_docker/app/php-{your-version}/Dockerfile) добавьте установку php модуля memcached или memcache и нескольких пакетов
```dockerfile
RUN apt-get update && apt-get install -y \
    libmemcached-dev \
    zlib1g-dev
# используйте один из php модулей для работы с сервером memcached
RUN pecl install memcache && docker-php-ext-enable memcache
RUN pecl install memcached && docker-php-ext-enable memcached
```
3. Выполните build.

## Возможные проблемы

- `bash: compose: Permission denied` - выполнить команду
```bash
chmod +x ./vendor/beeralex/compose/src/bin/compose
```
- `compose: not found или compose: /bin/bash^M: bad interpreter` - выполнить команду
```bash
sed -i 's/\r$//' ./vendor/beeralex/compose/src/bin/compose
```