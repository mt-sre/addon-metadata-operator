# Types

## OCM

OCM types can be found inside internal Red Hat codebases. For example, addon subtypes are defined in `uhc-clusters-service/pkg/models/addons.go`.

These types are defined outside of `/apis/<version>/<type_name>_type.go` as we do not control the behavior of upstream OCM API.

We also modify these types slightly to match the interface we want to expose to our tenants, and add validation tags.

## Optional fields guideline

The use of pointers is required for optional fields so we can distinguish between a field being unset or the zero value.

| Type       | Optional Type |
| ---------- | ------------- |
| `string`   | `*string`     |
| `bool`     | `*bool`       |
| `int`      | `*int`        |
| `<struct>` | `*<struct>`   |
| `[]<type>` | `*[]<type>`   |
