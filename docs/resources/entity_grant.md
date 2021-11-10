# Entity Grant Resource

Entities can have Grants. Grants are relationships between a target Entity and one of two other types:

* A Recipient Entity
* A User.

Grants can have zero or more Permissions associated with them.

[Entity Grant API](https://fusionauth.io/docs/v1/tech/apis/entity-management/grants)

## Example Usage

### Example User Grant Request

```hcl
resource "fusionauth_entity_grant" "david" {
  data        = jsonencode({
    expiresAt = 1695361142909
  })
  permissions = [
    fusionauth_entity_type_permission.read.name, // note: permissions must be bound by name
    fusionauth_entity_type_permission.write.name,
    fusionauth_entity_type_permission.sue.name,
    fusionauth_entity_type_permission.file_lawsuit.name,
  ],
  user_id     = fusionauth_user.david.id
}
```

### Example Grant Request

```hcl
resource "fusionauth_entity_grant" "raviga" {
  data                = jsonencode({
    foo = "bar"
  })
  permissions         = [
    fusionauth_entity_type_permission.read.name, // note: permissions must be bound by name
    fusionauth_entity_type_permission.write.name,
    fusionauth_entity_type_permission.sue.name,
    fusionauth_entity_type_permission.file_lawsuit.name,
  ],
  recipient_entity_id = fusionauth_entity.raviga.id
}
```

## Argument Reference

* `?` - (Optional)
* `name` - (Required) The name of the Entity.
