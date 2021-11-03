# Entity Type Permission Resource

A resource for permissions that are associated with entity types.
These permissions can then be granted to either entities or user receipients for specific entities.

[Entity Type API - Grants](https://fusionauth.io/docs/v1/tech/apis/entity-management/entity-types/#create-an-entity-type-permission)

## Example

```hcl
resource "fusionauth_entity_type_permission" "switch_toggle" {
  entity_type_id = fusionauth_entity_type.switch.id
  name = "Toggle Switch"
}
```

## Argument Reference

* `name` - (Required) The name of the permission.
* `entity_type_id` - (Required) The ID of the entity type to create a permission on
* `description` - (Optional) long form description of the permission
* `is_default` - (Optional) should the permission be granted by default when a permission grant is created for an entity of the type this permission is associated with.