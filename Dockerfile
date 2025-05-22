FROM golang:1.24-alpine

RUN apk add --no-cache \
    make \
    bash \
    curl \
    openssl \
    docker-cli \
    docker-compose