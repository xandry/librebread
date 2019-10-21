#!/bin/sh

docker build . -t docker.pkg.github.com/vasyahuyasa/librebread/librebread:dev
docker push docker.pkg.github.com/vasyahuyasa/librebread/librebread:dev
