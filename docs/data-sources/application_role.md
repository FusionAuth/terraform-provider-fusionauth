# Application Role Resource

This Resource is used to create a role for an Application.

[Application Roles API](https://fusionauth.io/docs/v1/tech/apis/applications)

## Example Usage

```hcl
data "fusionauth_application_role" "admin" {
  application_id = data.fusionauth_application.FusionAuth.id
  name           = "admin"
}
```

## Argument Reference

* `application_id` - (Required) ID of the application that this role is for.
* `name` - (Required) The name of the Role.