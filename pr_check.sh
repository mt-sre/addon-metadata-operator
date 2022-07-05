#!/bin/bash

set -exvo pipefail -o nounset

GO_1_17="/opt/go/1.17.8/bin"

if [ -d  "${GO_1_17}" ]; then
     PATH="${GO_1_17}:${PATH}"
fi

# pre-emptively install go-sqlite3 to ensure amalgamated libsqlite3
# source is present for compilation.
go install github.com/mattn/go-sqlite3

./mage run-hooks && ./mage test
