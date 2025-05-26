FROM golang:1.24-slim

ARG USERGROUP=1000

RUN apt-get update && apt-get install -y \
    bash \
    curl \
    openssl \
    make \
    docker.io \
    docker-compose \
    && rm -rf /var/lib/apt/lists/*
    
RUN groupadd -g $USERGROUP appuser && \
    useradd -m -u $USERGROUP -g $USERGROUP -s /bin/bash appuser
