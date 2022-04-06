#!/bin/bash

set -exvo pipefail -o nounset

function changes_include(){
     changed_files=$(git diff origin/master..HEAD --name-only | grep -c "$1")

     [ "${changed_files}" -gt 0 ]
}

GO_1_16="/opt/go/1.16.15/bin"

if [ -d  "${GO_1_16}" ]; then
     PATH="${GO_1_16}:${PATH}"
fi

# pre-emptively install go-sqlite3 to ensure amalgamated libsqlite3
# source is present for compilation.
go install github.com/mattn/go-sqlite3

make check test test-e2e build

# only build docker image if Dockerfile or Makefile has changed.
if changes_include "Dockerfile.build" || changes_include "Makefile"; then
     make docker-build
fi
