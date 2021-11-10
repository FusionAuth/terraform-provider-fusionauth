# Entity Type Permission Resource

Permissions are defined on an Entity Type. These are arbitrary strings which can fit the business domain. A Permission
could be `read`, `write`, or `file-lawsuit`.

[Entity Type Permission API](https://fusionauth.io/docs/v1/tech/apis/entity-management/entity-types/#create-an-entity-type-permission)

## Example Usage

```hcl
resource "fusionauth_entity_type_permission" "file_lawsuit" {
  entity_type_id = fusionauth_entity_type.company.id
  data           = jsonencode({
    foo = "bar"
  })
  description    = "Enables the ability to file lawsuits"
  is_default     = true
  name           = "file-lawsuit"
}
```

## Argument Reference

* `entity_type_id` - (Required) The ID of the Entity Type.
* `permission_id` - (Optional) The ID to use for the new permission. If not specified a secure random UUID will be
  generated.
* `data` - (Optional) An object that can hold any information about the Permission that should be persisted. Must be a
  JSON string.
* `description` - (Optional) The description of the Permission.
* `is_default` - (Optional) Whether or not the Permission is a default permission. A default permission is automatically
  granted to an entity of this type if no permissions are provided in a grant request.
* `name` - (Required) The name of the Permission.
