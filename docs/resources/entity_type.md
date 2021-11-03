# Entity Type Resource

Entity types serve as archetypes for entities. 
Permissions that can be granted for individual entities are defined on entity types. 
The provider models this as anadditional entity_type_permission resource.
For more information see the official documentation

[Entity Type API](https://fusionauth.io/docs/v1/tech/apis/entity-management/entity-types/)

## Example

```hcl
resource "fusionauth_entity_type" "switch" {
  name = "switch"
}
```

## Argument Reference

* `name` - (Required) The name of the type.
* `entity_type_id` - (Optional) A specific UUID to use as the id for the type
* `data` - (Optional) An object that can hold any information about the type that should be persisted.