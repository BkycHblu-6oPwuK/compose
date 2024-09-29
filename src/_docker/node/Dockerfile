# Use the official Debian image as the base
FROM ubuntu:22.04

ENV PATH=$PATH:/usr/bin

# Install essential tools
RUN apt-get update && apt-get install -y \
    apt-utils \
    sudo \
    build-essential \
    gnupg2 \
    curl \
    bash

# Set the Node.js version and user/group configurations
ARG NODE_VERSION=20
ARG USERGROUP=1000
ARG NODE_PATH=local/js/vite

# Install Node.js
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash - && \
    apt-get install -y nodejs && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Create a new user
RUN groupadd --force -g $USERGROUP nodeuser && \
    useradd -ms /bin/bash --no-user-group -g $USERGROUP -u $USERGROUP nodeuser

RUN mkdir -p /var/www/${NODE_PATH} && \
    chown -R nodeuser:nodeuser /var/www

# Switch to the new user
USER nodeuser

WORKDIR ${NODE_PATH}
