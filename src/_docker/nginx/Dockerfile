FROM nginx:latest

ARG USERGROUP
ARG DOCKER_PATH

RUN groupadd --force -g $USERGROUP nginxuser
RUN useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP nginxuser

COPY ${DOCKER_PATH}/nginx/certs/fullchainloc.pem /usr/local/share/ca-certificates/
RUN apt-get update && apt-get install -y \
    ca-certificates \
    wget \
    && update-ca-certificates

# для получения сертификатов letscrypt
# RUN apt-get update && apt-get install -y \
#     certbot \
#     python3-certbot-nginx
