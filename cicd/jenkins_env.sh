#!/bin/bash

set -exvo pipefail -o nounset

GO_1_18="/opt/go/1.18.1/bin"

if [ -d  "${GO_1_18}" ]; then
     PATH="${GO_1_18}:${PATH}"
fi
