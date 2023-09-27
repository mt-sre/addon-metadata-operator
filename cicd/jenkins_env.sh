#!/bin/bash

set -exvo pipefail -o nounset

GO_1_20="/opt/go/1.20.6/bin"

if [ -d  "${GO_1_20}" ]; then
     PATH="${GO_1_20}:${PATH}"
fi
