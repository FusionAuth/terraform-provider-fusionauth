# Group Resource

A FusionAuth Group is a named object that optionally contains one to many Application Roles.

When a Group does not contain any Application Roles it can still be utilized to logically associate users. Assigning Application Roles to a group allow it to be used to dynamically manage Role assignment to registered Users. In this second scenario as long as a User is registered to an Application the Group membership will allow them to inherit the corresponding Roles from the Group.

[Groups API](https://fusionauth.io/docs/v1/tech/apis/groups)

## Example Usage

```hcl
resource "fusionauth_group" "my_group" {
  name      = "My Group"
  tenant_id = fusionauth_tenant.my_tenant.id
  role_ids = [
    fusionauth_application_role.admins.id,
  ]
}
```

## Argument Reference

* `group_id` - (Optional) The Id to use for the new Group. If not specified a secure random UUID will be generated.
* `data` - (Optional) An object that can hold any information about the Group that should be persisted.
* `name` - (Required) The name of the Group.
* `role_ids` - (Optional) The Application Roles to assign to this group.
* `tenant_id` - (Required) The unique Id of the tenant used to scope this API request.
