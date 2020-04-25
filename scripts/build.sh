#!/bin/bash
set -e
set -x

DIR=`pwd`
NAME=`basename ${DIR}`
SHA=`git rev-parse --short HEAD`
VERSION=${VERSION:-$SHA}

GOOS=linux GOARCH=amd64 go build .

docker build -t hatamiarash7/${NAME}:${VERSION} .
docker push hatamiarash7/${NAME}:${VERSION}

rm arvand-exporter
