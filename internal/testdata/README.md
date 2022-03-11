# Testdata <!-- omit in TOC -->

- [Description](#description)
- [Paths](#paths)
- [Schemas](#schemas)

## Description

This package contains testdata used to work with both the `mtcli` CLI and the `addon-metadata-operator`.

## Paths

| Path                     | Description                                                                                            |
| ------------------------ | ------------------------------------------------------------------------------------------------------ |
| `metadata_v1/`           | Contains addons in the metadata v1 format.                                                             |
| `metadata_v1/imagesets/` | Contains v1 addons that use the imageset feature.                                                      |
| `metadata_v1/legacy/`    | Contains addons that don't use the imageset feature, and statically set the `indexImage: <...>` field. |
| `metadata_v2/`           | Contains addons in the metadata v2 format, which is CRD based (AddonMetadata, AddonImageSet).          |
| `validators/AMXXXX/`     | Contains resources used to test validators (e.g.: CSV manifests)                                       |
| `bundles/`               | Contains OLM operator bundles: https://olm.operatorframework.io/docs/tasks/creating-operator-bundle/   |

## Schemas

- [Metadata v1 JSON schema](https://github.com/mt-sre/managed-tenants-cli/blob/main/managedtenants/schemas/metadata.schema.yaml)
- [ImageSet v1 JSON schema](https://github.com/mt-sre/managed-tenants-cli/blob/main/managedtenants/schemas/imageset.schema.yaml)
- [Metadata v2 CRD spec](../../api/v1alpha1/addonmetadata_types.go)
- [ImageSet v2 CRD spec](../../api/v1alpha1/addonimageset_types.go)
