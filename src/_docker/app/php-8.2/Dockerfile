FROM php:8.2-fpm

ARG USERGROUP
ENV SUPERVISOR_PHP_USER="appuser"

# Установка зависимостей
RUN apt-get update && apt-get install -y \
    apt-utils \
    libpq-dev \
    libpng-dev \
    libzip-dev \
    libfreetype6-dev \
    libjpeg62-turbo-dev \
    libwebp-dev \
    libxpm-dev \
    zip \
    unzip \
    git \
    jpegoptim \
    optipng \
    ghostscript \
    gifsicle && \
    docker-php-ext-configure gd --with-freetype --with-jpeg --with-webp --with-xpm && \
    docker-php-ext-install pdo_mysql mysqli bcmath gd zip

# Install composer
ENV COMPOSER_ALLOW_SUPERUSER=1
RUN curl -sS https://getcomposer.org/installer | php -- \
    --filename=composer \
    --install-dir=/usr/local/bin

# Установка и включение расширения Imagick
# RUN pecl install imagick && docker-php-ext-enable imagick

# Xdebug
# RUN pecl install xdebug \
#     && docker-php-ext-enable xdebug

RUN groupadd --force -g $USERGROUP appuser
RUN useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP appuser
# Очистка
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
USER appuser
WORKDIR /var/www

#RUN mkdir ~/.ssh && ln -s /run/secrets/host_ssh_key ~/.ssh/id_rsa