#!/bin/bash

set -exvo pipefail -o nounset

# Need to install proper go version because Jenkins runs go < 1.16
echo "[DEBUG] System go version is: $(go version)"
export GO111MODULE=on
export PATH="$(go env GOPATH)/bin:${PATH}"
export GOVERSION="1.16.8"

function downloadGo() {
    go get golang.org/dl/go${GOVERSION}@latest
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
