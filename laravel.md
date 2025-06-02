# Конфигурация под битрикс

nginx + php (8.2, 8.3, 8.4) + mysql|postgres|sqlite + node 23 версии
а так же redis|memcached

## Шаги публикации docker-compose.yml

Шаги:

0. Выбор фреймворка - laravel
1. Выбрать версию php
2. Выберите базу данных - mysql|postgres|sqlite
3. Выбрать версию базы данных
4. Выберите сервер для кеширования - redis|memcached|Пропуск
5. Выполните пункты установщика laravel

После этого в директории где выполнялась команда появиться docker-compose.yml файл с настроенными сервисами.

В директории site уже будет развернут laravel проект и скачаны npm пакеты

## Cron

Крон не устанавливается. В laravel используйте Task Scheduling - https://laravel.com/docs/12.x/scheduling
В supervisor настроен запуск задач и queue worker

## Почта

msmtp клиент не устанавливается.
Используйте сервис mailhog для тестирования отправки почты.

В файле .env laravel проект найдите эти переменные и установите такие же значения

```
MAIL_MAILER=smtp
MAIL_HOST=mailhog
MAIL_PORT=1025
```

## Символические ссылки

При добавлении сайта автоматически добавляются символическая ссылка на storage

```
/var/www/storage/app/public /var/www/public/storage 
```

## Vite config 

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

