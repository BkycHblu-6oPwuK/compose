FROM php:7.4-fpm

ARG USERGROUP
ARG DOCKER_PATH
ENV SUPERVISOR_PHP_USER="appuser"

RUN groupadd --force -g $USERGROUP appuser
RUN useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP appuser

COPY ${DOCKER_PATH}/nginx/certs/ /usr/local/share/ca-certificates/
COPY ${DOCKER_PATH}/nginx/certs/ /etc/nginx/certs

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
    msmtp \
    nginx \
    ca-certificates \
    wget \
    cron \
    supervisor \
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

COPY ${DOCKER_PATH}/app/msmtprc /home/appuser/.msmtprc
RUN chmod 600 /home/appuser/.msmtprc
RUN chown appuser:appuser /home/appuser/.msmtprc

COPY ${DOCKER_PATH}/app/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ${DOCKER_PATH}/app/cron/appuser.txt /var/spool/cron/crontabs/appuser
COPY ${DOCKER_PATH}/app/cron/root.txt /var/spool/cron/crontabs/root
RUN echo "" >> /var/spool/cron/crontabs/appuser
RUN echo "" >> /var/spool/cron/crontabs/root
RUN chown appuser:appuser /var/spool/cron/crontabs/appuser
RUN chmod 0600 /var/spool/cron/crontabs/appuser /var/spool/cron/crontabs/root
RUN touch /var/log/cron.log
RUN service cron restart
COPY ${DOCKER_PATH}/bin/app_entrypoint.sh /usr/local/bin/entrypoint.sh
RUN sed -i 's/\r$//' /usr/local/bin/entrypoint.sh && \
    chmod +x /usr/local/bin/entrypoint.sh
COPY ${DOCKER_PATH}/files/simlinks /usr/simlinks
COPY ${DOCKER_PATH}/bin/create_simlink.sh /usr/local/bin/create_simlink.sh
RUN sed -i 's/\r$//' /usr/local/bin/create_simlink.sh && \
    chmod +x /usr/local/bin/create_simlink.sh

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /var/www

#ssh RUN mkdir ~/.ssh && ln -s /run/secrets/host_ssh_key ~/.ssh/id_rsa

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]