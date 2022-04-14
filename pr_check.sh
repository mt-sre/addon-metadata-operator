#!/bin/bash

set -exvo pipefail -o nounset

GO_1_16="/opt/go/1.16.15/bin"

if [ -d  "${GO_1_16}" ]; then
     PATH="${GO_1_16}:${PATH}"
fi

# pre-emptively install go-sqlite3 to ensure amalgamated libsqlite3
# source is present for compilation.
go install github.com/mattn/go-sqlite3

./mage run-hooks && ./mage test
