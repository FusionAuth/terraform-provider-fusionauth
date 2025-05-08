# Application OAuth Scope Resource

The Application OAuth Scope resource allows you to define the scopes that an application can request when using OAuth.

[Application OAuth Scope API](https://fusionauth.io/docs/apis/scopes)

## Example Usage

```hcl
resource "fusionauth_application_oauth_scope" "this" {
  application_id = fusionauth_application.this.id
  data = jsonencode({
    createdBy = "jared@fusionauth.io"
  })
  default_consent_detail  = "This will provide the requesting application read-only access to your data"
  default_consent_message = "View your data"
  description             = "Provides an application read-only access to a user's data"
  name                    = "data:read"
  required                = true
}

```

## Argument Reference

* `application_id` - (Required) ID of the application that this role is for.
* `name` - (Required) The name of the Role.

---

* `data` - (Optional) A JSON string that can hold any information about the OAuth Scope that should be persisted.
* `default_consent_detail` - (Optional) "The default detail to display on the OAuth consent screen if one cannot be found in the theme.
* `default_consent_message` - (Optional) The default message to display on the OAuth consent screen if one cannot be found in the theme.
* `description` - (Optional) A description of the OAuth Scope. This is used for display purposes only.
* `required` - (Optional) Determines if the OAuth Scope is required when requested in an OAuth workflow.
* `scope_id` - (Optional) The Id to use for the new OAuth Scope. If not specified a secure random UUID will be generated.

## Import

In Terraform v1.5.0 and later, use an `import` block to import application scopes using the application ID and scope ID, separated by a colon. For example:

```hcl
import {
  to = fusionauth_application_oauth_scope.name
  id = "application_id:scope_id"
}
```

Using terraform import, import application roles using the application ID and scope ID. For example:

```shell
terraform import fusionauth_application_role.name application_id:scope_id
