#!/bin/bash

set -exvo pipefail -o nounset

# Simple script to download and run goreleaser
# uses config from .goreleaser.yml
curl -sL https://git.io/goreleaser | bash
