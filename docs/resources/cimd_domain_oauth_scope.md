# CIMD Domain OAuth Scope Resource

The CIMD Domain OAuth Scope resource allows you to define the custom OAuth scopes that clients registered through a CIMD domain can request, in addition to the standard OpenID Connect scopes enabled by the domain's `provided_scope_policy`.

Note: An Enterprise plan is required to utilize CIMD. Available since FusionAuth `1.70.0`.

[CIMD Domain OAuth Scopes API](https://fusionauth.io/docs/apis/cimd-domain-scopes)

## Example Usage

```hcl
resource "fusionauth_cimd_domain_oauth_scope" "this" {
  cimd_domain_id = fusionauth_cimd_domain.example.id
  data = jsonencode({
    createdBy = "jared@fusionauth.io"
  })
  default_consent_detail  = "This will provide the requesting client read-only access to your data"
  default_consent_message = "View your data"
  description             = "Provides a CIMD client read-only access to a user's data"
  name                    = "data:read"
  required                = true
}
```

## Argument Reference

* `cimd_domain_id` - (Required) The Id of the CIMD domain to which this scope belongs. Changing this value forces a new resource to be created.
* `name` - (Required) The name of the OAuth Scope. This is the value that will be used to request the scope in OAuth workflows. The name cannot be changed after creation, so changing this value forces a new resource to be created.

---

* `data` - (Optional) A JSON string that can hold any information about the OAuth Scope that should be persisted.
* `default_consent_detail` - (Optional) The default detail to display on the OAuth consent screen if one cannot be found in the theme. [Learn more about setting this value using themes.](https://fusionauth.io/docs/theme/localization#oauth-scope-consent-prompt)
* `default_consent_message` - (Optional) The default message to display on the OAuth consent screen if one cannot be found in the theme. [Learn more about setting this value using themes.](https://fusionauth.io/docs/theme/localization#oauth-scope-consent-prompt)
* `description` - (Optional) A description of the OAuth Scope for internal use.
* `required` - (Optional) Determines if the OAuth Scope is required when requested in an OAuth workflow. Defaults to `false`.
* `scope_id` - (Optional) The Id to use for the new OAuth Scope. If not specified a secure random UUID will be generated.

## Import

In Terraform v1.5.0 and later, use an `import` block to import CIMD domain scopes using the CIMD domain ID and scope ID, separated by a colon. For example:

```hcl
import {
  to = fusionauth_cimd_domain_oauth_scope.name
  id = "cimd_domain_id:scope_id"
}
```

Using terraform import, import CIMD domain scopes using the CIMD domain ID and scope ID. For example:

```shell
terraform import fusionauth_cimd_domain_oauth_scope.name cimd_domain_id:scope_id
```
