# Application OAuth Scope Resource

The Application OAuth Scope resource allows you to define the scopes that an application can request when using OAuth.

[Application OAuth Scope API](https://fusionauth.io/docs/apis/scopes)

## Example Usage

```hcl
data "fusionauth_application_oauth_scope" "this" {
  application_id = data.fusionauth_application.this.id
  name           = "data:read"
}
```

## Argument Reference

* `application_id` - (Required) ID of the application that this role is for.
* `name` - (Required) The name of the Role.

## Attributes Reference

All of the argument attributes are also exported as result attributes.

The following additional attributes are exported:

* `data` - (Optional) An object that can hold any information about the OAuth Scope that should be persisted.
* `default_consent_detail` - (Optional) "The default detail to display on the OAuth consent screen if one cannot be found in the theme.
* `default_consent_message` - (Optional) The default message to display on the OAuth consent screen if one cannot be found in the theme.
* `description` - (Optional) A description of the OAuth Scope. This is used for display purposes only.
* `required` - (Optional) Determines if the OAuth Scope is required when requested in an OAuth workflow.
* `scope_id` - (Optional) The Id to use for the new OAuth Scope. If not specified a secure random UUID will be generated.
