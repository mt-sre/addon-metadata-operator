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

Use `.envrc` with [direnv](https://github.com/direnv/direnv) to make sure your environment is set properly:

```bash
export GOROOT=$(go1.16.7 env GOROOT)
export GOBIN=${PWD}/.cache/bin
export PATH=${GOROOT}/bin:${GOBIN}:${PATH}
export KUBECONFIG=${PWD}/.cache/kubeconfig
```

## Develop

### Useful make commands

Install `pre-commit` hooks:

```bash
$ pre-commit install
```

Run code checks:

```bash
$ make check
```

Run tests:

```bash
$ make test
```

Regenerate manifests and CRDs (`./config/crd/bases/addons.managed.openshift.io_addonmetadata.yaml` whenever you modify the `AddonMetadata` type:

```bash
$ make manifest
```

You can also use the `Dockerfile.test` for a more stable environment:

```bash
$ ./pr_check.sh
```

## Release

## License

Apache License 2.0, see [LICENSE](LICENSE).
