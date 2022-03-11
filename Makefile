# These shell flags are REQUIRED for an early exit in case any program called by make errors!
.SHELLFLAGS=-euo pipefail -c
SHELL := /bin/bash

.PHONY: all test fmt vet clean build tidy docker-build docker-push test-e2e

REPO := quay.io/mtsre/addon-metadata-operator
TAG := $(shell git rev-parse --short HEAD)

# Set the GOBIN environment variable so that dependencies will be installed
# always in the same place, regardless of the value of GOPATH
CACHE := $(PWD)/.cache
export GOBIN := $(CACHE)/bin
export PATH := $(GOBIN):$(PATH)
export KUBECONFIG := $(CACHE)/kubeconfig
PKGS := $(shell go list ./... | grep -v -E '/vendor|/integration')
INTEGRATION_TESTS := $(shell go list ./integration...)
E2E_MTCLI_PATH := $(CACHE)/mtcli

# make prow to NOT expect this project to have vendoring
GOFLAGS=
# required by opm to extract sql-based catalog
export CGO_CFLAGS := -DSQLITE_ENABLE_JSON1

all: build

##@ Development

test: ## Run tests.
	@go test -count=1 -race $(PKGS)


test-e2e: ## Run e2e integration tests
	@CGO_ENABLED=1 go build -a -o $(E2E_MTCLI_PATH) cmd/mtcli/main.go
	@E2E_MTCLI_PATH=$(E2E_MTCLI_PATH) go test -count=1 -race $(INTEGRATION_TESTS)

check: golangci-lint goimports ## Runs all checks.
	@go fmt $(PKGS) $(INTEGRATION_TESTS)
	@go vet $(PKGS) $(INTEGRATION_TESTS)

clean: ## Clean this directory
	@if [ -f "$(KIND)" ]; then $(KIND) delete cluster --name $(KIND_CLUSTER_NAME); fi
	@rm -fr $(CACHE) $(GOBIN) bin/* dist/ || true
	@find . -type d -name "*index_tmp_*" -exec rm -fr {} +

##@ Build

build: tidy generate ## Build binaries
	# Disable cgo for for cross-compilation: https://pkg.go.dev/cmd/cgo
	@CGO_ENABLED=1 go build -a -o bin/mtcli cmd/mtcli/main.go
	@CGO_ENABLED=0 go build -a -o bin/addon-metadata-operator cmd/addon-metadata-operator/main.go

tidy:
	@go mod tidy
	@go mod verify

docker-build: ## Build docker image with the operator.
	@docker build -t $(REPO):$(TAG) -f Dockerfile.build .

docker-push: ## Push docker image with the operator.
	@if [ -z "$(DOCKER_CONF)" ]; then echo "Please set DOCKER_CONF. Exiting." && exit 1; fi
	@docker tag $(REPO):$(TAG) $(REPO):latest
	@docker --config=$(DOCKER_CONF) push $(REPO):$(TAG)
	@docker --config=$(DOCKER_CONF) push $(REPO):latest

##@ CRD and K8S

manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) crd rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

generate: controller-gen manifests ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

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


## Dependencies

CONTROLLER_GEN := $(GOBIN)/controller-gen
controller-gen:
	@$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.7.0)

KUSTOMIZE := $(GOBIN)/kustomize
kustomize:
	@$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v4)

KIND := $(GOBIN)/kind
kind:
	@$(call go-get-tool,$(KIND),sigs.k8s.io/kind)

GOIMPORTS := $(GOBIN)/goimports
goimports:
	@$(call go-get-tool,$(GOIMPORTS),golang.org/x/tools/cmd/goimports)
	@$(GOIMPORTS) -w -l $(shell find . -type f -name "*.go" -not -path "./vendor/*")

GOLANGCI_LINT := $(GOBIN)/golangci-lint
golangci-lint:
	@$(call go-get-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0)
	@echo "Running golangci-lint..."
	@$(GOLANGCI_LINT) run --timeout=10m -E unused,gosimple,staticcheck --skip-dirs-use-default --verbose

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
