FROM debian:12-slim

ARG NODE_VERSION=23
ARG USERGROUP
ENV PATH=$PATH:/usr/bin

RUN apt-get update && apt-get install -y \
    apt-utils \
    sudo \
    build-essential \
    gnupg2 \
    curl \
    bash

RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash - && \
    apt-get install -y nodejs && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN groupadd --force -g $USERGROUP docky && \
    useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP docky

RUN npm install -g pm2

WORKDIR /var/www