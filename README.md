# addon-metadata-operator

- [addon-metadata-operator](#addon-metadata-operator)
  - [Develop](#develop)
    - [Useful make commands](#useful-make-commands)
    - [Adding validators](#adding-validators)
  - [Release](#release)
    - [mtcli](#mtcli)
  - [License](#license)

Operator responsible for managing AddOn resources in OSD

## Develop

### Useful make commands

Install `pre-commit` hooks:

```bash
./mage hooks:enable
```

Run code checks:

```bash
./mage check
```

Run tests:

```bash
./mage test
```

```bash
./pr_check.sh
```

### Adding validators

See this [doc](docs/adding_validators.md) for more information on adding new validators.

## Release

### mtcli

The external Jenkins integration will run the `build_tag.sh` script after you push a tag. The script runs `goreleaser`, which is responsible for building the binaries and publishing them as a Github Release.

```bash
git tag vX.X.X
git push upstream --tags
```

**A tag can't be built twice, you will need to add a new commit and a new tag referencing the latest commit.**

## License

Apache License 2.0, see [LICENSE](LICENSE).
