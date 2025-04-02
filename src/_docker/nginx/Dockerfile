FROM nginx:latest

ARG USERGROUP
ARG DOCKER_PATH

RUN groupadd --force -g $USERGROUP nginxuser
RUN useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP nginxuser

COPY ${DOCKER_PATH}/nginx/certs/ /usr/local/share/ca-certificates/
COPY ${DOCKER_PATH}/nginx/certs/ /etc/nginx/certs
RUN apt-get update && apt-get install -y \
    ca-certificates \
    wget \
    && update-ca-certificates

COPY ${DOCKER_PATH}/bin/nginx_entrypoint.sh /usr/local/bin/entrypoint.sh
RUN sed -i 's/\r$//' /usr/local/bin/entrypoint.sh && \
        chmod +x /usr/local/bin/entrypoint.sh

COPY ${DOCKER_PATH}/files/simlinks /usr/simlinks
COPY ${DOCKER_PATH}/bin/create_simlink.sh /usr/local/bin/create_simlink.sh
RUN sed -i 's/\r$//' /usr/local/bin/create_simlink.sh && \
        chmod +x /usr/local/bin/create_simlink.sh

# для получения сертификатов letscrypt
# RUN apt-get update && apt-get install -y \
#     certbot \
#     python3-certbot-nginx

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]