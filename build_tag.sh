#!/bin/bash

set -exvo pipefail -o nounset

# Run build_tag.sh in a sandbox to enforce golang & gcc bindings (CGO_ENABLED=1):
# - sync secrets from app-interface:
#   -> secrets: /resources/jenkins/global/secrets.yaml
#   -> for job `gh-build-tag`: /resources/jenkins/global/templates.yaml
IMAGE_TEST=addon-metadata-operator
docker build -t ${IMAGE_TEST} -f Dockerfile.ci .
docker_run_args=(
    --rm
    --privileged
    -e CGO_ENABLED=1
    # github API token to post release
    -e "GITHUB_TOKEN=${GITHUB_TOKEN}"
    -v /var/run/docker.sock:/var/run/docker.sock
    -v $(pwd):/go/src/github.com/mt-sre/addon-metadata-operator
    -w /go/src/github.com/mt-sre/addon-metadata-operator
    # goreleaser-cross version
    -e "VERSION=v1.17.6"
    release --rm-dist
)
docker run "${docker_run_args[@]}" "${IMAGE_TEST}" -c goreleaser-cross.sh
