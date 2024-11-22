FROM php:7.4-fpm

ARG USERGROUP
ARG DOCKER_PATH
ENV SUPERVISOR_PHP_USER="appuser"

RUN groupadd --force -g $USERGROUP appuser
RUN useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP appuser

COPY ${DOCKER_PATH}/nginx/certs/ /usr/local/share/ca-certificates/
COPY ${DOCKER_PATH}/nginx/certs/ /etc/nginx/certs

RUN echo "postfix postfix/mailname string localhost" | debconf-set-selections && \
    echo "postfix postfix/main_mailer_type string 'Internet Site'" | debconf-set-selections && \
    apt-get update && apt-get install -y \
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
    postfix \
    nginx \
    ca-certificates \
    wget \
#cron    cron \
#cron    supervisor \
#imagick    libmagickwand-dev imagemagick \
    gifsicle && \
    update-ca-certificates && \
#imagick    pecl install imagick && \
#xdebug    pecl install xdebug-3.1.6 && \
#imagick    docker-php-ext-enable imagick && \
#xdebug    docker-php-ext-enable xdebug && \
    docker-php-ext-configure gd --with-freetype --with-jpeg --with-webp --with-xpm && \
    docker-php-ext-install pdo_mysql mysqli bcmath gd zip

ENV COMPOSER_ALLOW_SUPERUSER=1
RUN curl -sS https://getcomposer.org/installer | php -- \
    --filename=composer \
    --install-dir=/usr/local/bin

#cron COPY ${DOCKER_PATH}/app/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
#cron COPY ${DOCKER_PATH}/app/cron/appuser.txt /var/spool/cron/crontabs/appuser
#cron COPY ${DOCKER_PATH}/app/cron/root.txt /var/spool/cron/crontabs/root
#cron RUN echo "" >> /var/spool/cron/crontabs/appuser
#cron RUN echo "" >> /var/spool/cron/crontabs/root
#cron RUN chown appuser:appuser /var/spool/cron/crontabs/appuser
#cron RUN chmod 0600 /var/spool/cron/crontabs/appuser /var/spool/cron/crontabs/root
#cron RUN touch /var/log/cron.log
#cron RUN service cron restart
#cron COPY ${DOCKER_PATH}/bin/app_entrypoint.sh /usr/local/bin/entrypoint.sh
#cron RUN sed -i 's/\r$//' /usr/local/bin/entrypoint.sh
#cron RUN chmod +x /usr/local/bin/entrypoint.sh

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

#можно закомментировать если используется cron
CMD service nginx start && php-fpm 

WORKDIR /var/www

#ssh RUN mkdir ~/.ssh && ln -s /run/secrets/host_ssh_key ~/.ssh/id_rsa

#cron ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]