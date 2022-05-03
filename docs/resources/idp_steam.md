# Steam Identity Provider Resource

The Steam identity provider type will use the Steam OAuth login API. It will also provide a Login with Steam button on FusionAuth’s login page that will direct a user to the Steam login page. The Steam login uses the implicit OAuth grant and will return to the callback URL with token and state in the URL Fragment. This is handled by the FusionAuth /oauth2/implicit javascript function to pass those values to the /oauth2/callback endpoint.

This identity provider will call Steam’s API to load the Steam user’s personaname and use that as username to lookup or create a user in FusionAuth depending on the linking strategy configured for this identity provider. However, Steam does not allow access to user emails, so neither email linking strategy can be used and user’s will not be able to login or be created.

FusionAuth will also store the Steam token that is returned from the Steam login in the link between the user and the identity provider. This token can be used by an application to make further requests to Steam APIs on behalf of the user.

[Steam Identity Provider APIs](https://fusionauth.io/docs/v1/tech/apis/identity-providers/steam/ )

## Example Usage

```hcl
resource "fusionauth_idp_steam" "steam" {
  application_configuration {
    application_id      = fusionauth_application.GPS_Insight.id
    create_registration = true
    enabled             = true
  }
  button_text   = "Login with Steam"
  client_id     = "0eb1ce3c-2fb1-4ae9-b361-d49fc6e764cc"
  scope         = "Xboxlive.signin Xboxlive.offline_access"
  web_api_key   = "HG0A97ACKPQ5ZLPU0007BN6674OA25TY"
}
```

## Argument Reference

* `idp_id` - (Optional) The ID to use for the new identity provider. If not specified a secure random UUID will be generated.
* `application_configuration` - (Optional) The configuration for each Application that the identity provider is enabled for.
    - `application_id` - (Optional) ID of the Application to apply this configuration to.
    - `button_text` - (Optional) This is an optional Application specific override for the top level button text.
    - `client_id` - (Optional) This is an optional Application specific override for the top level client_id.
    - `create_registration` - (Optional) Determines if a UserRegistration is created for the User automatically or not. If a user doesn’t exist in FusionAuth and logs in through an identity provider, this boolean controls whether or not FusionAuth creates a registration for the User in the Application they are logging into.
    - `enabled` - (Optional) Determines if this identity provider is enabled for the Application specified by the applicationId key.
    - `scope` - (Optional)This is an optional Application specific override for the top level scope.
    - `web_api_key` - (Optional) This is an optional Application specific override for the top level webAPIKey.
* `button_text` - (Required) The top-level button text to use on the FusionAuth login page for this Identity Provider.
* `client_id` - (Required) The top-level Steam client id for your Application. This value is retrieved from the Steam developer website when you setup your Steam developer account.
* `debug` - (Optional) Determines if debug is enabled for this provider. When enabled, each time this provider is invoked to reconcile a login an Event Log will be created.
* `enabled` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
* `lambda_reconcile_id` - (Optional) The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user.
* `linking_strategy` - (Optional) The linking strategy to use when creating the link between the {idp_display_name} Identity Provider and the user.
* `scope` - (Optional) The top-level scope that you are requesting from Steam.
* `web_api_key` - (Required) The top-level web API key to use with the Steam Identity Provider when retrieving the player summary info. This value is retrieved from the Steam developer website when you setup your Steam developer account.
* `tenant_configuration` - (Optional) The configuration for each Tenant that limits the number of links a user may have for a particular identity provider.
    - `tenant_id` - (Optional) The unique Id of the tenant that this configuration applies to.
    - `limit_user_link_count_enabled` - (Optional) When enabled, the number of identity provider links a user may create is enforced by maximumLinks.
    - `limit_user_link_count_maximum_links` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
