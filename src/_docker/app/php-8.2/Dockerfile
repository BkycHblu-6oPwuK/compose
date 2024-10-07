FROM php:8.2-fpm

ARG USERGROUP
ARG DOCKER_PATH
ENV SUPERVISOR_PHP_USER="appuser"

RUN groupadd --force -g $USERGROUP appuser
RUN useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP appuser

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
    jpegoptim \
    optipng \
    ghostscript \
    gifsicle && \
    docker-php-ext-configure gd --with-freetype --with-jpeg --with-webp --with-xpm && \
    docker-php-ext-install pdo_mysql mysqli bcmath gd zip

ENV COMPOSER_ALLOW_SUPERUSER=1
RUN curl -sS https://getcomposer.org/installer | php -- \
    --filename=composer \
    --install-dir=/usr/local/bin

COPY ${DOCKER_PATH}/nginx/certs/ /usr/local/share/ca-certificates/
RUN apt-get update && apt-get install -y \
    ca-certificates \
    wget \
    && update-ca-certificates

RUN apt-get update && apt-get install -y \
    postfix

RUN apt-get update && apt-get install -y nginx

#cron RUN apt-get update && apt-get install -y \
#cron     cron \
#cron     supervisor

#cron COPY ${DOCKER_PATH}/app/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
#cron COPY ${DOCKER_PATH}/app/crontab.txt /tmp/crontab.txt
#cron RUN cat /tmp/crontab.txt >> /etc/crontab && echo "" >> /etc/crontab && rm /tmp/crontab.txt
#cron RUN crontab /etc/crontab && service cron reload
#cron RUN touch /var/log/cron.log
#cron COPY ${DOCKER_PATH}/bin/app_entrypoint.sh /usr/local/bin/entrypoint.sh
#cron RUN sed -i 's/\r$//' /usr/local/bin/entrypoint.sh
#cron RUN chmod +x /usr/local/bin/entrypoint.sh

#imagick RUN pecl install imagick && docker-php-ext-enable imagick

#xdebug RUN pecl install xdebug \
#xdebug     && docker-php-ext-enable xdebug

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

#можно закомментировать если используется cron
CMD service nginx start && php-fpm 

WORKDIR /var/www

#ssh RUN mkdir ~/.ssh && ln -s /run/secrets/host_ssh_key ~/.ssh/id_rsa

#cron ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]