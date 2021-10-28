#!/bin/bash

set -exvo pipefail -o nounset

# Need to install proper go version because Jenkins runs go < 1.16
echo "[DEBUG] System go version is: $(go version)"
# Set GO111MODULE=off to avoid modifying go.mod or go.sum
export GO111MODULE=off
export PATH="$(go env GOPATH)/bin:${PATH}"
export GOVERSION="1.16.8"

function downloadGo() {
    go get golang.org/dl/go${GOVERSION}
    go${GOVERSION} download
}

function setupEnv() {
    export GOROOT=$(go${GOVERSION} env GOROOT)
    export PATH="${GOROOT}/bin:${PATH}"
}

downloadGo
setupEnv

# Simple script to download and run goreleaser
# uses config from .goreleaser.yml
curl -sL https://git.io/goreleaser | bash
