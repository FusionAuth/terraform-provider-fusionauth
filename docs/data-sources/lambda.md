# Lambda Resource

Lambdas are user defined JavaScript functions that may be executed at runtime to perform various functions. Lambdas may be used to customize the claims returned in a JWT, reconcile a SAML v2 response or an OpenID Connect response when using these external identity providers.

[Lambdas API](https://fusionauth.io/docs/v1/tech/apis/lambdas)

## Example Usage

```hcl
data "fusionauth_lambda" "default_google_reconcile" {
    name = "Default Google Reconcile provided by FusionAuth"
    type = "GoogleReconcile"
}
```

## Argument Reference

* `name` - (Optional) The name of the Lambda. At least one of `id` or `name` must be specified.
* `id`   - (Optional) The ID of the Lambda. At least one of `id` or `name` must be specified.
* `type` - (Required) The Lambda type. The possible values are:
  * `AppleReconcile`
  * `ClientCredentialsJWTPopulate`
  * `EpicGamesReconcile`
  * `ExternalJWTReconcile`
  * `FacebookReconcile`
  * `GoogleReconcile`
  * `HYPRReconcile`
  * `JWTPopulate`
  * `LDAPConnectorReconcile`
  * `LinkedInReconcile`
  * `LoginValidation`
  * `NintendoReconcile`
  * `OpenIDReconcile`
  * `SAMLv2Populate`
  * `SAMLv2Reconcile`
  * `SCIMServerGroupRequestConverter`
  * `SCIMServerGroupResponseConverter`
  * `SCIMServerUserRequestConverter`
  * `SCIMServerUserResponseConverter`
  * `SelfServiceRegistrationValidation`
  * `SonyPSNReconcile`
  * `SteamReconcile`
  * `TwitchReconcile`
  * `TwitterReconcile`
  * `UserInfoPopulate`
  * `XboxReconcile`

## Attributes Reference

All of the argument attributes are also exported as result attributes.

The following additional attributes are exported:

* `body`  - The lambda function body, a JavaScript function.
* `debug` - Whether or not debug event logging is enabled for this Lambda.
