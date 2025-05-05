# Docker Compose for Bitrix

Настроенная обертка над docker compose для локальных проектов bitrix 

nginx + php (7.4, 8.2, 8.3, 8.4) + mysql + node 23 версии

Для установки можно запустить установочный скрипт с правами суперпользователя:

```bash
curl -sSL https://raw.githubusercontent.com/BkycHblu-6oPwuK/compose/main/scripts/install.sh | sudo sh
```
Либо же можете вручую скачать бинарник (находится в build/docky) и поместить его в ```/usr/local/bin``` и не забудьте дать ему необходимые права (команда ```chmod +x```)

Обновление скрипта происходит точно также, при выполнении команды curl - файл заменяется.

После установки проверьте работу скрипта, можете выполнить команду ```docky --help```

## Публикация docker-compose.yml

Команды выполняются в директории с docker-compose.yml или в любой другой дочерней, но при этом файл docker-compose.yml должен существовать (для всех команд, кроме - init и clean-cache).

Для публикации docker-compose.yml выполните команду:

```bash
docky init
```

Если docker-compose.yml уже существует, то будет предложено создать новый - в таком случае файл будет переименован и после этого вы перейдете к дальшейшей публикации нового файла, а иначе будет произведен выход.
Сам сайт размещается в директории ```site```, которая создается в той же директории где находиться docker-compose.yml

Шаги установки:

1. Выбрать версию php
2. Выбрать версию mysql - 5.7 или 8.0
3. Устанавливать ли node.js - Y или N, если установите Y, то будет создан сервис с node js 23 версии
    1. Если будет устанавливаться node.js, то нужно указать корневую директорию для него, то есть директория содержащая файл package.json. Путь указывается относительно корня сайта - local/js/vite или пустое поле если package.json в корне сайта.
4. Устанавливать ли sphinx - Y или N, если установите Y, то будет создан сервис с shphinx версии 2.2.11

После этого в директории где выполнялась команда появиться docker-compose.yml файл с настроенными сервисами.

## Конфигурация yml через файл .env

В этом файле задается версия php, mysql, node, путь к package json и путь к сайту если необходимо заменить стандартный

```
PHP_VERSION=7.4|8.2|8.3|8.4 # по сути не фактическая версия php - используется для построения пути к dockerfile в ```_docker/app/php-{PHP_VERSION}/dockerfile```, в самом докерфайле задается фактическая версия
MYSQL_VERSION={любая версия доступная на docker hub}
NODE_VERSION={любая версия доступная по ссылке - https://deb.nodesource.com/setup_${NODE_VERSION}.x}
NODE_PATH=/var/www/local/js/vite # здесь путь до package json в контейнере, поэтому указывайте вместе абсолютный путь. /var/www - это DOCUMENT_ROOT сайта в контейнере
SITE_PATH={абсолютный или относительный путь к директории сайта} # указывайте если не хотите размещать сайт в директории site по умолчанию
```

## Публикация докерфайлов и файлов конфигурации

Если вам необходимо внести изменения в докерфайлы или файлы конфигурации, или добавить что то свое, то используйте команду:

```bash
docky publish
```

При этом если директория ./_docker уже существует, то она будет переименована.

## SSL сертификаты для nginx

Сертификаты и ключи копируются в контейнер из /_docker/nginx/certs/ и запись о них уже добавлена в nginx.conf.

Свои сертификаты вы так же можете помещать в /_docker/nginx/certs/ и после этого делать build.

Сейчас там используются временные самописные сертификаты которые будут действительны до ~2051 года.

Сервер одинаково настроен на работу как по http, так и по https.

Более подробно - https://github.com/BkycHblu-6oPwuK/compose/tree/main/src/_docker/nginx

## php

Доступные версии php - 7.4, 8.2, 8.3, 8.4.

Конфигурации для каждых из версий находятся по пути - ``` ./_docker/app/php-{PHP_VERSION}/ ```.

Для изменения версии - измените ее в файле .env, переменая ```PHP_VERSION```

## Xdebug

По умолчанию установлен.

Либо же проверьте установку пакетов в dockerfile 

```dockerfile
RUN pecl install xdebug && \
    docker-php-ext-enable xdebug
```
 
сервис app в docker-compose yml должен выглядеть вот так

```
app:
    build:
        context: ${DOCKER_PATH}
        dockerfile: ${DOCKER_PATH}/app/php-${PHP_VERSION}/Dockerfile
        args:
            USERGROUP: ${USERGROUP}
    volumes:
        - ${SITE_PATH}:/var/www
        - ${DOCKER_PATH}/app/php-${PHP_VERSION}/php.ini:/usr/local/etc/php/conf.d/php.ini
        - ${DOCKER_PATH}/app/php-${PHP_VERSION}/xdebug.ini:/usr/local/etc/php/conf.d/xdebug.ini # xdebug ini
        - ${DOCKER_PATH}/app/php-fpm.conf:/usr/local/etc/php-fpm.d/zzzzwww.conf
        - ${DOCKER_PATH}/app/nginx:/etc/nginx/conf.d
    ports:
        - 9000:9000
    environment: # переенные для xdebug
        PHP_IDE_CONFIG: serverName=xdebugServer
        XDEBUG_TRIGGER: testTrig
    depends_on:
        - mysql
    networks:
        - compose
    extra_hosts:
        - host.docker.internal:host-gateway # extra host
    container_name: app
```

публикуйте файлы конфигурации и настраивайте ```_docker/app/php-${PHP_VERSION}/xdebug.ini``` по своему усмотрению

## Cron

По умолчанию cron включен и выполняется задание на запуск файла ```/var/www/bitrix/modules/main/tools/cron_events.php```

Если необходимо добавить задания, то сделайте публикацию докерфайлов и файлов конфигурации

```bash
docky publish
```

Запись заданий осуществляйте в:
- `_docker/app/cron/appuser.txt` - для пользователя сайта
- `_docker/app/cron/root.txt` - для root пользователя

выполните команду

```bash
docky build
```

## Nginx в php контейнере

Для работы с сокетами в php кентейнер был установлен nginx который проксирует запросы на контейнер nginx.

nginx.conf для контейнера php лежит в _docker/app/nginx.conf

## Mysql

Версия mysql меняется в файле .env, переменная ```MYSQL_VERSION```, версия любая доступная на docker hub для образа mysql

По умолчанию база данных храниться в томе (volumes) mysql_data, если хотите хранить базу локально в директории, то замените mysql_data на вашу директорию, например - ./tmp/db

Так же в сервис прокидывается файл конфигурации my.cnf, он располагается в ./_docker/mysql/my.cnf - туда вы можете вносить свои правки.

## Node и npm, npx

По умолчанию установлена 23 версия.

Так же по умолчанию контейнер работает на двух портах - ``` 5173 ``` и ``` 5174 ```

При установке можно указать корневую директории для вашего фронтенда, то есть та директория, где расположен ``` package.json ```. Контейнер работает внутри данной директории и соответственно там выполняются команды.

Вы можете изменить версию ``` node ```  и корневую директорию, поменяв значение переменных ``` NODE_VERSION ``` и ``` NODE_PATH ``` в файле .env

Скачивание node идет с адреса - https://deb.nodesource.com/setup_${NODE_VERSION}.x, и часто подверсий там нет, указывайте целые числа или можете перейти по ссылке и проверить доступность версии.

Команды npm/npx выполняются слудующим образом:
```bash
docky npm {arg}
docky npx {arg}
```

## Пакет pm2 в node контейнере

В dockerfile node устанавливается пакет pm2, для работы с ним используйте команду

```bash 
docky pm2 {arg}
```

Запустить сервер node js можно:

1. командой pm2 - ```docky pm2 start server.js```
2. настроить на запуск при запуске контейнеров - для этого добавьте команду (в сервисе node) - ``` command: sh -c "pm2 start /var/www/local/js/vite/server.js --name node-server & tail -f /dev/null" ```, здесь указывается точка входа к скрипту который будет запускаться

## Почта

Для отправки почты настроен SMTP клиент  - ```msmtp```

Для завершения настройки вам необходимо добавить вашу почту (в поля user и from) и пароль в файл:

- `_docker/app/msmtprc`

По умолчанию в этом файле заготовка под почту яндекса, для других сервисов просто сделайте новый блок аккаунта на основе yandex и замените имя аккаунта в строке ```account default```, тогда ваш аккаунт будет по умолчанию использоваться при отправке почты.

Если почта не отправляется или в проверке системы написано что почта не работает, то проверьте логи ```msmtp``` в контейнере, которые находятся в файле ```/home/appuser/msmtp.log```. Вероятнее всего произошла ошибка авторизации или почтовый сервис отклюнил отправку из-за подозрений в спаме.

## Туннелирование локального сайта

Для туннелирования используется Expose (https://github.com/beyondcode/expose) и для того чтобы поделиться вашим локальным сайтом выполните команду:

```bash
docky share
```

Сайт будет доступен 1 час, после этого команду можно выполнить заново.

доступные флаги expose для проброса через команду docky share:

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

## Создание новых сайтов

0. Выполните команду ```docky down``` или убедитесь что контейнеры остановлены
1. Для создания сайта выполните команду ``` docky create site ```
2. Введите доменное имя сайта
3. Проверьте что все создано (в директории вашего сайта должна была появиться директории с названием введенного вашего доменного имени)
4. Выполните команду ``` docky build ```

## Символические ссылки

При запуске контейнера app запускается скрипт - _docker/bin/create_simlink.sh, он создает ссылки внутри контейнера и соответственно ссылки внутри сайта распространяются и на хост и другие контейнеры.

Ссылки берутся из файла - _docker/app/simlinks. Структура файла должна быть такой:

```
/var/www/bitrix /var/www/<domain>/bitrix
/var/www/local /var/www/<domain>/local
/var/www/upload /var/www/<domain>/upload
```

Соответственно при добавлении сайта автоматически добавляются символические ссылки на каталоги - bitrix, local, upload для введенного домена.

Если же вам нужно дополнительные ссылки добавить, то формируйте все пути относительно структуры контейнера. После этого необходимо выполнить команду:

```bash
docky build
```

## Создание нового домена для основного сайта

0. Выполните команду ```docky down``` или убедитесь что контейнеры остановлены
1. Для создания выполните команду ``` docky create domain ```
2. Введите доменное имя сайта
3. Выполните команду ``` docky build ```

## Добавление записей в hosts 

При создании сайта или домена создается файл hosts в директории с docker-compose.yml

В целом вы можете добавлять в него записи вида:

```
127.0.0.1 new_site
```

После этого можно выполнить команду:

```bash
docky hosts push
```

Но лучше выполнять команду ``` docky create domain ```, если вы хотите чтобы создались сертификаты и конфиги для nginx

И все записи из вашего локального hosts будут добавлены в глобальный (${SYSTEM_DISK}\Windows\System32\drivers\etc\hosts - если wsl, или /etc/hosts - если ubuntu)

(${SYSTEM_DISK} - скрипт через команду powershell попытается найти системный диск)

Если записи уже существуют, то дублирования не будет.

## Описание всех доступных команд

- `init` - Создание docker-compose.yml в текущей директории и создает директорию ```site```
```bash
docky init
```
- `publish` - Публикация файлов конфигурации в директории с docker-compose.yml, доступен флаг ``` service ``` для публикации отдельного сервиса в docker-compose.yml (node или sphinx)
```bash
docky publish
docky publish --service node|sphinx
```
- `clean-cache` - очищает кэш директории скрипта, в ней храняться файлы конфигурации, докерфайлы
```bash
docky clean-cache
```
- `upgrade` - изменяет docker-compose.yml под вторую версия скрипта, при этом старый файл будет переименован и вы всегда можете откатить изменение
```bash
docky upgrade
```
- `create site` - Создание нового сайта в директории сайта (./site/new-site.ru)
```bash
docky create site
```
- `create domain` - Создание нового домена для вашего основного сайта
```bash
docky create domain
```
- `hosts push` - Переносит записи из вашего локального hosts файла в глобальный
```bash
docky hosts push
```
- `share` - Позволяет сделать сайт доступным из интернета
```bash
docky share
```
- `php` - Выполнение команды php в контейнере с php
```bash
docky php -v
```
- `composer` - Выполнение команды composer в контейнере с php
```bash
docky composer install
```
- `npm` - Выполнение команды npm в контейнере с node js если он был установлен
```bash
docky npm install
```
- `npx` - Выполнение команды npx в контейнере с node js если он был установлен
```bash
docky npx create-vite@latest my-app
```
- `pm2` - Выполнение команды pm2 в контейнере с node js если он был установлен. Команда принимает такие же аргументы как и оригинальная pm2
```bash
docky pm2 {arg}
```
- `И все дефолтные команды docker-compose` - Выполнение любой команды docker-compose
```bash
docky up -d
docky down
docky build
```

## Пользователи в контейнерах

- `docky` - в контейнере с php (service app), nginx, node

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
2. В dockerfile вашего app сервиса (_docker/app/php-{PHP_VERSION}/Dockerfile) добавьте установку php модуля redis и igbinary
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
2. В dockerfile вашего app сервиса (_docker/app/php-{PHP_VERSION}/Dockerfile) добавьте установку php модуля memcached или memcache и нескольких пакетов
```dockerfile
# эти пакеты нужны, если вы используете php модуль memcached
RUN apt-get update && apt-get install -y \
    libmemcached-dev \
    zlib1g-dev
# используйте один из php модулей для работы с сервером memcached
RUN pecl install memcache && docker-php-ext-enable memcache
RUN pecl install memcached && docker-php-ext-enable memcached
```
3. Выполните build.

## Переход с первой версии на вторую

Для перехода выполните команду:

```bash
docky upgrade
```

Скрипт попытается поменять docker-compose.yml под вторую версию и создаст файл .env.

После выполнения рекомендуется проверить корректность изменения (файл docker-compose.yml, файл .env)