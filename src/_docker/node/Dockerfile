FROM ubuntu:22.04

ENV PATH=$PATH:/usr/bin

RUN apt-get update && apt-get install -y \
    apt-utils \
    sudo \
    build-essential \
    gnupg2 \
    curl \
    bash

ARG NODE_VERSION=20
ARG USERGROUP=1000
ARG NODE_PATH=local/js/vite

RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash - && \
    apt-get install -y nodejs && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN groupadd --force -g $USERGROUP nodeuser && \
    useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP nodeuser

USER nodeuser

WORKDIR ${NODE_PATH}
