# Twitch Identity Provider Resource

The Twitch identity provider type will use the Twitch OAuth v2.0 login API. It will also provide a Login with Twitch button on FusionAuth’s login page that will direct a user to the Twitch login page.

This identity provider will call Twitch’s API to load the user’s email and preferred_username and use those as email and username to lookup or create a user in FusionAuth depending on the linking strategy configured for this identity provider. Additional claims returned by Twitch can be used to reconcile the user to FusionAuth by using a Twitch Reconcile Lambda.

FusionAuth will also store the Twitch refresh_token returned from the Twitch API in the link between the user and the identity provider. This token can be used by an application to make further requests to Twitch APIs on behalf of the user.

[Twitch Identity Provider APIs](https://fusionauth.io/docs/v1/tech/apis/identity-providers/twitch/)

## Example Usage

```hcl
resource "fusionauth_idp_twitch" "twitch" {
  application_configuration {
    application_id      = fusionauth_application.my_app.id
    create_registration = true
    enabled             = true
  }
  button_text   = "Login with Twitch"
  client_id     = "0eb1ce3c-2fb1-4ae9-b361-d49fc6e764cc"
  client_secret = "693s000cbn66k0mxtqzr_c_NfLy3~6_SEA"
  scope         = "Xboxlive.signin Xboxlive.offline_access"
}
```

## Argument Reference

* `idp_id` - (Optional) The ID to use for the new identity provider. If not specified a secure random UUID will be generated.
* `application_configuration` - (Optional) The configuration for each Application that the identity provider is enabled for.
    - `application_id` - (Optional) ID of the Application to apply this configuration to.
    - `button_text` - (Optional) This is an optional Application specific override for the top level button text.
    - `client_id` - (Optional) This is an optional Application specific override for the top level client_id.
    - `client_secret` - (Optional) This is an optional Application specific override for the top level client_secret.
    - `create_registration` - (Optional) Determines if a UserRegistration is created for the User automatically or not. If a user doesn’t exist in FusionAuth and logs in through an identity provider, this boolean controls whether or not FusionAuth creates a registration for the User in the Application they are logging into.
    - `enabled` - (Optional) Determines if this identity provider is enabled for the Application specified by the applicationId key.
    - `scope` - (Optional)This is an optional Application specific override for the top level scope.
* `button_text` - (Required) The top-level button text to use on the FusionAuth login page for this Identity Provider.
* `client_id` - (Required) TThe top-level Xbox client id for your Application. This value is retrieved from the Xbox developer website when you setup your Xbox developer account.
* `client_secret` - (Required) The top-level client secret to use with the Xbox Identity Provider when retrieving the long-lived token. This value is retrieved from the Xbox developer website when you setup your Xbox developer account.
* `debug` - (Optional) Determines if debug is enabled for this provider. When enabled, each time this provider is invoked to reconcile a login an Event Log will be created.
* `enabled` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
* `lambda_reconcile_id` - (Optional) The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user.
* `linking_strategy` - (Optional) The linking strategy to use when creating the link between the {idp_display_name} Identity Provider and the user.
* `scope` - (Optional) The top-level scope that you are requesting from Xbox.
* `tenant_configuration` - (Optional) The configuration for each Tenant that limits the number of links a user may have for a particular identity provider.
    - `tenant_id` - (Optional) The unique Id of the tenant that this configuration applies to.
    - `limit_user_link_count_enabled` - (Optional) When enabled, the number of identity provider links a user may create is enforced by maximumLinks.
    - `limit_user_link_count_maximum_links` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
