# Apple Identity Provider Resource

The Apple identity provider type will use the Sign in with Apple APIs and will provide a Sign with Apple button on FusionAuth’s login page that will either redirect to an Apple sign in page or leverage native controls when using Safari on macOS or iOS. Additionally, this identity provider will call Apples’s /auth/token API to load additional details about the user and store them in FusionAuth.

FusionAuth will also store the Apple refresh_token that is returned from the /auth/token endpoint in the UserRegistration object inside the tokens Map. This Map stores the tokens from the various identity providers so that you can use them in your application to call their APIs.


[Apple Identity Providers API](https://fusionauth.io/docs/v1/tech/apis/identity-providers/apple/#create-the-apple-identity-provider)

## Example Usage

```hcl
resource "fusionauth_idp_apple" "apple" {
  application_configuration {
    application_id      = "1c212e59-0d0e-6b1a-ad48-f4f92793be32"
    create_registration = true
    enabled             = true
  }
  button_text = "Sign in with Apple"
  debug       = false
  enabled     = true
  key_id      = "2f81529c-4d39-4ce2-982e-cf5fbb1325f6"
  scope       = "email name"
  services_id = "com.piedpiper.webapp"
  team_id     = "R4NQ1P4UEB"
}
```

## Argument Reference

* `application_configuration` - (Optional) The configuration for each Application that the identity provider is enabled for.
    - `application_id` - (Optional) ID of the Application to apply this configuration to.
    - `button_text` - (Optional) This is an optional Application specific override for the top level button text.
    - `create_registration` - (Optional) Determines if a UserRegistration is created for the User automatically or not. If a user doesn’t exist in FusionAuth and logs in through an identity provider, this boolean controls whether or not FusionAuth creates a registration for the User in the Application they are logging into.
    - `enabled` - (Optional) Determines if this identity provider is enabled for the Application specified by the applicationId key.
    - `key_id` - (Optional) This is an optional Application specific override for the top level keyId.
    - `scope` - (Optional) This is an optional Application specific override for for the top level scope.
    - `services_id` - (Optional) This is an optional Application specific override for for the top level servicesId.
    - `team_id` - (Optional) This is an optional Application specific override for for the top level teamId.
* `button_text` - (Required) The top-level button text to use on the FusionAuth login page for this Identity Provider.
* `debug` - (Optional) Determines if debug is enabled for this provider. When enabled, each time this provider is invoked to reconcile a login an Event Log will be created.
* `enabled` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
* `key_id` - (Required) The unique Id of the private key downloaded from Apple and imported into Key Master that will be used to sign the client secret.
* `lambda_reconcile_id` - (Optional) The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user.
* `scope` - (Optional) The top-level space separated scope that you are requesting from Apple.
* `services_id` - (Required) The unique Id of the private key downloaded from Apple and imported into Key Master that will be used to sign the client secret.
* `team_id` - (Required) The Apple App ID Prefix, or Team ID found in your Apple Developer Account which has been configured for Sign in with Apple.
* `linking_strategy` - (Optional) The linking strategy to use when creating the link between the {idp_display_name} Identity Provider and the user.
* `tenant_configuration` - (Optional) The configuration for each Tenant that limits the number of links a user may have for a particular identity provider.
    - `tenant_id` - (Optional) The unique Id of the tenant that this configuration applies to.
    - `limit_user_link_count_enabled` - (Optional) When enabled, the number of identity provider links a user may create is enforced by maximumLinks.
    - `limit_user_link_count_maximum_links` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
