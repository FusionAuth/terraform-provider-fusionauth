# Entity Resource

Entities are instantiations of entity types used for permissions and client credentials.
Entities are associated with an entity type and a tenant.

[Entity API](https://fusionauth.io/docs/v1/tech/apis/entity-management/entities/)

## Example

```hcl
resource "fusionauth_entity" "light_switch" {
  tenant_id = fusionauth_tenant.house.id
  entity_type_id = fusionauth_entity_type.switch.id
  name = "Light Switch"
}
```

## Argument Reference

* `name` - (Required) The name of the type.
* `tenant_id` - (Required) The id of the entity's tenant
* `entity_type_id` - (Required) The ID of the entity's type
* `client_id` - (Optional) the client ID for the entity
* `client_secret` - (Optional) The client secret for the entity
* `data` - (Optional) An object that can hold any information about the entity that should be persisted.