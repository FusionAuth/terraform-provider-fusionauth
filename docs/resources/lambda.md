# Lambda Resource

Lambdas are user defined JavaScript functions that may be executed at runtime to perform various functions. Lambdas may be used to customize the claims returned in a JWT, reconcile a SAML v2 response or an OpenID Connect response when using these external identity providers.

[Lambdas API](https://fusionauth.io/docs/v1/tech/apis/lambdas)

## Example Usage

```hcl
resource "fusionauth_lambda" "preferred_username" {
  name    = "Preferred Username"
  type    = "JWTPopulate"
  enabled = true
  body    = <<EOT
// Using the user and registration parameters add additional values to the jwt object.
function populate(jwt, user, registration) {
  jwt.preferred_username = registration.username;
}
  EOT
}
```

## Argument Reference

* `body` - (Required) The lambda function body, a JavaScript function.
* `name` - (Required) The name of the lambda.
* `type` - (Required) The lambda type. The possible values are:
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
  * `MFARequirement`
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

---

* `debug` - (Optional) Whether or not debug event logging is enabled for this Lambda.
* `engine_type` - (Optional) The JavaScript execution engine for the lambda.
* `lambda_id` - (Optional) The Id to use for the new lambda. If not specified a secure random UUID will be generated.
