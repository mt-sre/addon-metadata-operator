# These shell flags are REQUIRED for an early exit in case any program called by make errors!
.SHELLFLAGS=-euo pipefail -c
SHELL := /bin/bash

.PHONY: all test fmt vet clean build tidy docker-build docker-push test-e2e release release-snapshot build-release-container

REPO := quay.io/mtsre/addon-metadata-operator
TAG := $(shell git rev-parse --short HEAD)

# Set the GOBIN environment variable so that dependencies will be installed
# always in the same place, regardless of the value of GOPATH
CACHE := $(PWD)/.cache
export GOBIN := $(CACHE)/bin
export PATH := $(GOBIN):$(PATH)
export KUBECONFIG := $(CACHE)/kubeconfig

# make prow to NOT expect this project to have vendoring
GOFLAGS=
# required by opm to extract sql-based catalog
export CGO_CFLAGS := -DSQLITE_ENABLE_JSON1

all: build

##@ Development

test: ## Run tests.
	./mage test:unit


test-e2e: ## Run e2e integration tests
	./mage test:integration

check: ## Runs all checks.
	./mage check

clean: ## Clean this directory
	@if [ -f "$(KIND)" ]; then $(KIND) delete cluster --name $(KIND_CLUSTER_NAME); fi
	./mage clean

##@ Build

build: tidy generate ## Build binaries
	./mage build:cli
	./mage build:operator

tidy:
	./mage check:tidy
	./mage check:verify

docker-build: ## Build docker image with the operator.
	./mage build:operatorimage

docker-push: ## Push docker image with the operator.
	./mage release:pushoperatorimage

##@ Release

release:
	./mage release:cli

release-snapshot:
	./mage release:clisnapshot

##@ CRD and K8S

manifests: ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	./mage generate:manifests

generate: ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	./mage generate:boilerplate

install: manifests kustomize ## Install CRDs to kind cluster
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize ## Uninstall CRDs from kind cluster
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize docker-build ## Deploy controller to kind cluster
	@$(KIND) load docker-image $(REPO):$(TAG) --name $(KIND_CLUSTER_NAME)
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(REPO):$(TAG)
	$(KUSTOMIZE) build config/default | kubectl apply -f -

undeploy: ## Undeploy controller from the kind cluster
	$(KUSTOMIZE) build config/default | kubectl delete -f -


##@ Kind

KIND_CLUSTER_NAME := addon-metadata-operator

kind-create: kind ## Create a plain kind cluster /w secret to allow pulling from Quay.io
	@$(KIND) create cluster --name $(KIND_CLUSTER_NAME) --kubeconfig $(KUBECONFIG) || true
	@echo -e "Ignoring existing cluster error...\n"

kind-integration: kind-create install deploy ## Create a kind cluster to run integration tests (deploy CRDS, operator, etc.)
	@./hack/deploy-console.sh
	#@./hack/deploy-olm.sh - Do we need OLM?
	#@./hack/deploy-hive.sh - Enable to debug SyncSet/SelectorSyncSet


KIND := $(GOBIN)/kind
kind:
	@$(call go-get-tool,$(KIND),sigs.k8s.io/kind)

# go-get-tool will 'go get' any package $2 and install it to $1.
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
