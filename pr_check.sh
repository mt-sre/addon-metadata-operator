#!/bin/bash

set -exvo pipefail -o nounset

IMAGE_TEST=addon-metadata-operator

docker build -t ${IMAGE_TEST} -f Dockerfile.ci .
docker run --rm ${IMAGE_TEST} check test test-e2e build

# build addon-metadata-operator containers
make docker-build
