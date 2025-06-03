FROM debian:12-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    default-libmysqlclient-dev \
    libexpat1-dev \
    wget \
    ca-certificates \
    gnupg \
    apt-utils \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

RUN wget http://sphinxsearch.com/files/sphinx-2.2.11-release.tar.gz && \
    tar xzf sphinx-2.2.11-release.tar.gz && \
    cd sphinx-2.2.11-release && \
    ./configure && \
    make && \
    make install

RUN mkdir -p /var/run/sphinx && \
    rm -rf /sphinx-2.2.11-release* && \
    mkdir -p /var/lib/sphinx/data

COPY sphinx/sphinx.conf /usr/local/etc/sphinx.conf

WORKDIR /var/lib/sphinx/data

CMD ["searchd", "--nodetach"]