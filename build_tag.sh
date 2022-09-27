#!/bin/bash

set -exvo pipefail -o nounset

./mage release:cli
