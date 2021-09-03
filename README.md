# addon-flow-operator

Operator responsible for managing AddOn resources in OSD

## Getting Started

### Go version 1.16

The initial scaffolding is done using `kubebuilder` which doesn't work with go1.17. Here is how to download and setup go1.16 for this project:

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
