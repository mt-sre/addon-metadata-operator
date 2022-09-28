#!/bin/bash

set -exvo pipefail -o nounset

source "${PWD}/cicd/jenkins_env.sh"

# pre-emptively install go-sqlite3 to ensure amalgamated libsqlite3
# source is present for compilation.
go install github.com/mattn/go-sqlite3

./mage run-hooks && ./mage test
