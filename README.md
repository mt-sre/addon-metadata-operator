# addon-metadata-operator

Operator responsible for managing AddOn resources in OSD

## Developer Guide

### Getting Started

The initial scaffolding is done using [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) which doesn't work with go1.17. Here is how to download and setup go1.16 for this project:

```bash
# https://golang.org/doc/manage-install#installing-multiple
$ go install golang.org/dl/go1.16.7@latest
$ go1.16.7 download

$ cat <<EOF > .envrc
export GOROOT=$(go1.16.7 env GOROOT)
export PATH=${GOROOT}/bin:${PATH}
EOF

$ direnv allow .
```

### `.envrc`

Here is my `.envrc` to be used with [direnv](https://github.com/direnv/direnv). This is to make sure my environment is set properly to work with this repo:

```bash
export GOROOT=$(go1.16.7 env GOROOT)
export GOBIN=${PWD}/.cache/bin
export PATH=${GOROOT}/bin:${GOBIN}:${PATH}
export KUBECONFIG=${PWD}/.cache/kubeconfig
```

### Useful make commands

**`$ make manifests`**

Use this to regenerate the underlying CRD located here: `./config/crd/bases/addons.managed.openshift.io_addonmetadata.yaml` whenever you modify the `AddonMetadata` type.
