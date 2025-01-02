#!/bin/bash

set -exvo pipefail -o nounset

GO_1_23="/opt/go/1.23.1/bin"

if [ -d  "${GO_1_23}" ]; then
     PATH="${GO_1_23}:${PATH}"
fi
