# addon-metadata-operator

- [addon-metadata-operator](#addon-metadata-operator)
  - [Developer Guide](#developer-guide)
    - [Getting Started](#getting-started)
    - [`.envrc`](#envrc)
  - [Develop](#develop)
    - [Useful make commands](#useful-make-commands)
  - [Release](#release)
    - [mtcli](#mtcli)
    - [addon-metadata-operator image](#addon-metadata-operator-image)
  - [License](#license)

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

Regenerate manifests and CRDs (`./config/crd/bases/addonsflow.redhat.openshift.io_addonmetadata.yaml` whenever you modify the `AddonMetadata` type:

```bash
$ make manifest
```

You can also use the `Dockerfile.test` for a more stable environment:

```bash
$ ./pr_check.sh
```

### Adding validators

See this [doc](docs/adding_validators.md) for more information on adding new validators.

## Release

### mtcli

The external Jenkins integration will run the `build_tag.sh` script after you push a tag. The script runs `goreleaser`, which is responsible for building the binaries and publishing them as a Github Release.

```bash
$ git tag vX.X.X
$ git push upstream --tags
```

**A tag can't be built twice, you will need to add a new commit and a new tag referencing the latest commit.**

### addon-metadata-operator image

The container image is built on every merge to the `master` branch. The external Jenkins integration will run the `build_deploy.sh` script, effectively building and pushing to: `https://quay.io/repository/app-sre/addon-metadata-operator`.

## License

Apache License 2.0, see [LICENSE](LICENSE).
