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
    # goreleaser version
    -e "VERSION=v0.184.0"
    # github API token to post release
    -e "GITHUB_TOKEN=${GITHUB_TOKEN}"
)
docker run "${docker_run_args[@]}" "${IMAGE_TEST}" goreleaser
