#!/bin/bash

set -exvo pipefail -o nounset

GO_1_19="/opt/go/1.19.3/bin"

if [ -d  "${GO_1_19}" ]; then
     PATH="${GO_1_19}:${PATH}"
fi
