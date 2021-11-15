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
  tenant_id = fusionauth_tenant.default.id
  entity_id = fusionauth_entity.raviga.id
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
  tenant_id = fusionauth_tenant.default.id
  entity_id = fusionauth_entity.raviga.id
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

* `tenant_id` - (Optional) The unique Id of the tenant used to scope this API request.
* `entity_id` - (Required) The Id of the Entity to which access is granted.
* `data` - (Optional) An object that can hold any information about the Grant that should be persisted. Please review the limits on data field types as you plan for and build your custom data schema.  Must be a JSON string.
* `permissions` - (Optional) The set of permissions of this Grant.
* `recipient_entity_id` - (Optional) The Entity Id for which access is granted. If `recipient_entity_id` is not provided, then the `user_id` will be required.
* `user_id` - (Optional) The User Id for which access is granted. If `user_id` is not provided, then the `recipient_entity_id` will be required.
