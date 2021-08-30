# Application Role Resource

This Resource is used to create a role for an Application.

[Application Roles API](https://fusionauth.io/docs/v1/tech/apis/applications)

## Example Usage

```hcl
resource "fusionauth_application_role" "my_app_admin_role" {
  application_id = fusionauth_application.my_app.id
  description    = ""
  is_default     = false
  is_super_role  = true
  name           = "admin"
}
```

## Argument Reference

* `application_id` - (Required) ID of the application that this role is for.
* `description` - (Optional) A description for the role.
* `name` - (Required) The name of the Role.
* `is_default` - (Optional) Whether or not the Role is a default role. A default role is automatically assigned to a user during registration if no roles are provided.
* `is_super_role` - (Optional) Whether or not the Role is a considered to be a super user role. This is a marker to indicate that it supersedes all other roles. FusionAuth will attempt to enforce this contract when using the web UI, it is not enforced programmatically when using the API.