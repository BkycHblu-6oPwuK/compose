FROM php:8.2-fpm

ARG USERGROUP
ENV SUPERVISOR_PHP_USER="docky"

RUN groupadd --force -g $USERGROUP docky
RUN useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP docky

COPY nginx/certs/ /usr/local/share/ca-certificates/
COPY nginx/certs/ /etc/nginx/certs

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
    unzip

RUN apt-get install -y \ 
    jpegoptim \
    optipng \
    ghostscript \
    msmtp \
    nginx \
    ca-certificates \
    wget \
    cron \
    supervisor \
    libmagickwand-dev imagemagick \
    gifsicle
    
RUN update-ca-certificates && \
    pecl install imagick && \
    pecl install xdebug && \
    docker-php-ext-enable imagick && \
    docker-php-ext-enable xdebug && \
    docker-php-ext-configure gd --with-freetype --with-jpeg --with-webp --with-xpm && \
    docker-php-ext-install pdo_mysql mysqli bcmath gd zip

RUN pecl install memcache && docker-php-ext-enable memcache

RUN pecl install igbinary && \
    pecl install -D 'enable-redis-igbinary="yes"' redis && \
    docker-php-ext-enable igbinary redis

ENV COMPOSER_ALLOW_SUPERUSER=1
RUN curl -sS https://getcomposer.org/installer | php -- \
    --filename=composer \
    --install-dir=/usr/local/bin

COPY app/msmtprc /home/docky/.msmtprc
RUN chmod 600 /home/docky/.msmtprc
RUN chown docky:docky /home/docky/.msmtprc

COPY app/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY app/cron/docky.txt /var/spool/cron/crontabs/docky
COPY app/cron/root.txt /var/spool/cron/crontabs/root
RUN echo "" >> /var/spool/cron/crontabs/docky
RUN echo "" >> /var/spool/cron/crontabs/root
RUN chown docky:docky /var/spool/cron/crontabs/docky
RUN chmod 0600 /var/spool/cron/crontabs/docky /var/spool/cron/crontabs/root
RUN touch /var/log/cron.log
RUN service cron restart
COPY bin/app_entrypoint.sh /usr/local/bin/entrypoint.sh
RUN sed -i 's/\r$//' /usr/local/bin/entrypoint.sh && \
    chmod +x /usr/local/bin/entrypoint.sh
COPY app/simlinks /usr/simlinks
COPY bin/create_simlink.sh /usr/local/bin/create_simlink.sh
RUN sed -i 's/\r$//' /usr/local/bin/create_simlink.sh && \
    chmod +x /usr/local/bin/create_simlink.sh

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /var/www

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]