FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    build-essential \
    libmysqlclient-dev \
    libexpat1-dev \
    wget

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