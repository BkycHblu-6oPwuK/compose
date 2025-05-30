FROM nginx:latest

ARG USERGROUP

RUN groupadd --force -g $USERGROUP docky && \
    useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP docky

COPY nginx/certs/ /usr/local/share/ca-certificates/
COPY nginx/conf.d /etc/nginx/conf.d

RUN apt-get update && apt-get install -y \
    ca-certificates \
    wget \
    && update-ca-certificates