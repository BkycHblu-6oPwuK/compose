FROM php:8.2-fpm

ARG USERGROUP
ARG DOCKER_PATH
ENV SUPERVISOR_PHP_USER="appuser"

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

ENV COMPOSER_ALLOW_SUPERUSER=1
RUN curl -sS https://getcomposer.org/installer | php -- \
    --filename=composer \
    --install-dir=/usr/local/bin

#cron RUN apt-get update && apt-get install -y \
#cron     cron \
#cron     supervisor
    
#cron COPY ${DOCKER_PATH}/app/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
#cron COPY ${DOCKER_PATH}/app/crontab.txt /etc/cron.d/my-crontab
#cron RUN chmod 0644 /etc/cron.d/my-crontab
#cron RUN echo "" >> /etc/cron.d/my-crontab && crontab /etc/cron.d/my-crontab
#cron RUN touch /var/log/cron.log

# RUN pecl install imagick && docker-php-ext-enable imagick

# RUN pecl install xdebug \
#     && docker-php-ext-enable xdebug

#cron RUN groupadd --force -g $USERGROUP appuser
#cron RUN useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP appuser

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

#cron COPY ${DOCKER_PATH}/bin/app_entrypoint.sh /usr/local/bin/entrypoint.sh
#cron RUN chmod +x /usr/local/bin/entrypoint.sh

WORKDIR /var/www

#RUN mkdir ~/.ssh && ln -s /run/secrets/host_ssh_key ~/.ssh/id_rsa

#cron ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]