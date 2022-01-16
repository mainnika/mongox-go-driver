# syntax = docker/dockerfile:1.3-labs
FROM registry.access.redhat.com/ubi8/ubi

RUN <<EOF cat >> /etc/yum.repos.d/mongo.repo
[mongodb-org]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/8/mongodb-org/4.4/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-4.4.asc
EOF

RUN set -eux \
    && dnf makecache \
	&& dnf install -yq mongodb-org-server golang \
    && dnf clean all

WORKDIR /root/mongox

ENV GOPATH=/root/go

COPY go.mod .
COPY go.sum .

RUN set -eux \
    && go mod download

COPY mongox-testing mongox-testing
COPY mongox mongox

CMD set -eux \
    && nohup mongod --dbpath $(mktemp -d) \
    & go test -timeout 30s -v ./...