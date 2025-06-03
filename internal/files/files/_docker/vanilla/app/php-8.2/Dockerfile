FROM php:8.2-fpm

ARG USERGROUP
ENV SUPERVISOR_PHP_USER="docky"

RUN groupadd --force -g $USERGROUP docky
RUN useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP docky

COPY nginx/certs/ /usr/local/share/ca-certificates/

RUN apt-get update && apt-get install -y \
    apt-utils \
    autoconf \
    pkg-config \
    build-essential \
    libpng-dev \
    libzip-dev \
    libfreetype6-dev \
    libjpeg62-turbo-dev \
    libwebp-dev \
    libxpm-dev \
    zip \
    unzip \
    ca-certificates

RUN apt-get install -y libmagickwand-dev imagemagick 

RUN apt-get install -y \
    wget \
    nginx \
    cron \
    supervisor \
    gifsicle \
    jpegoptim \
    optipng \
    ghostscript \
    curl \
    libxml2-dev \
    libonig-dev \
    libsqlite3-dev \
    libpq-dev

RUN update-ca-certificates

RUN pecl install imagick xdebug && \
    docker-php-ext-enable imagick xdebug

RUN docker-php-ext-configure gd --with-freetype --with-jpeg --with-webp --with-xpm && \
    docker-php-ext-install \
    pdo_mysql \
    mysqli \
    pdo_pgsql \
    pgsql \
    pdo_sqlite \
    bcmath \
    gd \
    zip \
    xml \
    opcache

RUN apt-get update && apt-get install -y \
  libmemcached-dev \
  zlib1g-dev
RUN pecl install memcached && docker-php-ext-enable memcached

RUN pecl install igbinary \
    && pecl install -D 'enable-redis-igbinary="yes"' redis \
    && docker-php-ext-enable igbinary redis

ENV COMPOSER_ALLOW_SUPERUSER=1

RUN curl -sS https://getcomposer.org/installer | php -- \
    --filename=composer \
    --install-dir=/usr/local/bin

COPY app/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY app/cron /var/spool/cron/crontabs
RUN echo "" >> /var/spool/cron/crontabs/docky
RUN chown docky:docky /var/spool/cron/crontabs/docky
RUN chmod 0600 /var/spool/cron/crontabs/docky /var/spool/cron/crontabs/root
RUN touch /var/log/cron.log
RUN service cron restart

COPY bin/app_entrypoint.sh /usr/local/bin/entrypoint.sh
RUN sed -i 's/\r$//' /usr/local/bin/entrypoint.sh && \
    chmod +x /usr/local/bin/entrypoint.sh

COPY app/symlinks /usr/symlinks
COPY bin/create_simlink.sh /usr/local/bin/create_simlink.sh
RUN sed -i 's/\r$//' /usr/local/bin/create_simlink.sh && \
    chmod +x /usr/local/bin/create_simlink.sh

COPY app/php-8.2/php.ini /usr/local/etc/php/conf.d/php.ini
COPY app/php-8.2/xdebug.ini /usr/local/etc/php/conf.d/xdebug.ini
COPY app/php-fpm.conf /usr/local/etc/php-fpm.d/zzzzwww.conf
COPY app/nginx /etc/nginx/conf.d

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /var/www

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
