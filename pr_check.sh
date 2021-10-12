#!/bin/bash

set -exvo pipefail -o nounset

IMAGE_TEST=addon-metadata-operator

docker build -t ${IMAGE_TEST} -f Dockerfile.test .
docker run --rm ${IMAGE_TEST}
make docker-build
