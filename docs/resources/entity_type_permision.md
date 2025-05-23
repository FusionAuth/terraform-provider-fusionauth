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
* `name` - (Required) The name of the Permission.

---

* `data` - (Optional) A JSON string that can hold any information about the Permission that should be persisted. Must be a JSON string.
* `description` - (Optional) The description of the Permission.
* `is_default` - (Optional) Whether or not the Permission is a default permission. A default permission is automatically granted to an entity of this type if no permissions are provided in a grant request.
* `permission_id` - (Optional) The ID to use for the new permission. If not specified a secure random UUID will be generated.

## Import

In Terraform v1.5.0 and later, use an `import` block to import entity type permission resources using the entity type ID and entity type permission ID, separated by a colon. For example:

```hcl
import {
  to = fusionauth_entity_type_permission.name
  id = "entity_type_id:entity_type_permission_id"
}
```

Using terraform import, import entity resources using the entity type ID and entity type permission id. For example:

```shell
terraform import fusionauth_entity_type_permission.name entity_type_id:entity_type_permission_id
