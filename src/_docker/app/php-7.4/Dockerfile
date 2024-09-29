FROM php:7.4-fpm

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
# RUN pecl install xdebug-3.1.6 \
#     && docker-php-ext-enable xdebug

RUN useradd -ms /bin/bash appuser
RUN chown -R appuser:appuser /var/www
USER appuser

# Очистка
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /var/www

#RUN mkdir ~/.ssh && ln -s /run/secrets/host_ssh_key ~/.ssh/id_rsa