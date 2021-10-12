#!/bin/bash

set -exvo pipefail -o nounset

export DOCKER_CONF="${PWD}/.docker"
mkdir -p "$DOCKER_CONF"
docker --config="$DOCKER_CONF" login -u="$QUAY_USER" -p="$QUAY_TOKEN" quay.io

make docker-build docker-push
