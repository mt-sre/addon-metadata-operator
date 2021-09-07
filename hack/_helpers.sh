#!/bin/bash

function latestGithubRelease() {
    repo="${1}"
    curl -fsSL "https://api.github.com/repos/${repo}/releases/latest" | jq -r .tag_name
}
