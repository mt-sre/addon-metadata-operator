#!/bin/bash

source ${PWD}/hack/_helpers.sh

# If release is missing, assume latest
release="${1:-$(latestGithubRelease 'operator-framework/operator-lifecycle-manager')}"
dockercfg="${HOME}/.docker/config.json"

echo "Installing OLM..."
kubectl apply -f https://github.com/operator-framework/operator-lifecycle-manager/releases/download/${release}/crds.yaml
kubectl apply -f https://github.com/operator-framework/operator-lifecycle-manager/releases/download/${release}/olm.yaml

echo "Waiting for deployment/olm-operator..."
kubectl wait --for=condition=available deployment/olm-operator -n olm --timeout=240s

echo "waiting for deployment/catalog-operator..."
kubectl wait --for=condition=available deployment/catalog-operator -n olm --timeout=240s

# Required so CatalogSource can pull image from private quay repositories
if [ -f "${dockercfg}" ]; then
    echo "Installing quay.io secret..."
    kubectl create secret generic -n olm quaycreds --from-file=.dockerconfigjson="${dockercfg}" --type=kubernetes.io/dockerconfigjson
fi
