# Entity Grant Resource

The grant of access between a receiving user or entity and a target entity.
Each grant confers a set of permissions on the receiver to the target.

[Entity Grants - API](https://fusionauth.io/docs/v1/tech/apis/entity-management/grants/)

## Example

```hcl
resource "fusionauth_entity_grant" "robot_toggle_door_switch" {
  grant_entity_id =  fusionauth_entity.light_switch.id
  recipient_entity_id = fusionauth_entity.robot.id
  permissions = [fusionauth_entity_type_permission.toggle.name]
}
```

Note that permissions are associated by the permission name rather than id.

## Argument Reference

* `grant_entity_id` - (Required) The target entity. This is the entity that actions will be taken against
* `recipient_entity_id` - (Optional) The ID of the entity type that should receive the permissions. If this is unset, then user_id should be set instead.
* `user_id` - (Optional) The ID of a user that should receive the permissions. If this is unset, then recipient_entity_id should be set instead.
* `permissions` - (Optional) An array of permission names that should be granted to the recipient (either user or entity).
* `data` - (Optional) Additional data that should be associated with the entity grant.