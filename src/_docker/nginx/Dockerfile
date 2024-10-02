FROM nginx:latest

ARG USERGROUP

RUN groupadd --force -g $USERGROUP nginxuser
RUN useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP nginxuser

# для получения сертификатов letscrypt
# RUN apt-get update && apt-get install -y \
#     certbot \
#     python3-certbot-nginx
