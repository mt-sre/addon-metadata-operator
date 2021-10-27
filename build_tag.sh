#!/bin/bash

set -exvo pipefail -o nounset

# Need to install proper go version because Jenkins runs go < 1.16
# and I could not get goreleaser to work inside a container
go install golang.org/dl/go1.16.7@latest
go1.16.7 download
export GOROOT=$(go1.16.7 env GOROOT)
export PATH=${GOROOT}/bin:${PATH}

# Simple script to download and run goreleaser
# uses config from .goreleaser.yml
curl -sL https://git.io/goreleaser | bash
