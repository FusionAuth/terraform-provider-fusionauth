# LinkedIn Identity Provider Resource

The LinkedIn identity provider type will use OAuth 2.0 to authenticate a user with LinkedIn. It will also provide a
`Login with LinkedIn` button on FusionAuth’s login page that will direct a user to the LinkedIn login page.
Additionally, after successful user authentication, this identity provider will call LinkedIn’s `/v2/me` and
`/v2/emailAddress` APIs to load additional details about the user and store them in FusionAuth.

The email address returned by the LinkedIn `/v2/emailAddress` API will be used to create or look up the existing user.
Additional claims returned by LinkedIn can be used to reconcile the User to FusionAuth by using a LinkedIn Reconcile
lambda. Unless you assign a reconcile lambda to this provider, only the email address will be used from the available
claims returned by LinkedIn.

FusionAuth will also store the LinkedIn `access_token` returned from the login endpoint in the `identityProviderLink`
object. This object is accessible using the Link API.

The `identityProviderLink` object stores the token so that you can use it in your application to call LinkedIn APIs on
behalf of the user if desired.

[LinkedIn Identity Providers API](https://fusionauth.io/docs/v1/tech/apis/identity-providers/linkedin)

## Example Usage

```hcl
resource "fusionauth_idp_linkedin" "linkedin" {
    application_configuration {
        application_id = fusionauth_application.myapp.id
        create_registration = true
        enabled = true
    }
    button_text = "Login with LinkedIn"
    debug = false
    enabled = true
    client_id = "9876543210"
    client_secret = "716a572f917640698cdb99e9d7e64115"
    scope = "r_emailaddress r_liteprofile"
}
```

## Argument Reference

* `application_configuration` - (Optional) The configuration for each Application that the identity provider is enabled for.
    - `application_id` - (Optional) ID of the FusionAuth Application to apply this configuration to.
    - `button_text` - (Optional) This is an optional Application specific override for the top level `button_text`.
    - `app_id` - (Optional) This is an optional Application specific override for the top level `app_id`.
    - `client_secret` - (Optional) This is an optional Application specific override for the top level `client_secret`.
    - `create_registration` - (Optional) Determines if a `UserRegistration` is created for the User automatically or not. If a user doesn’t exist in FusionAuth and logs in through an identity provider, this boolean controls whether or not FusionAuth creates a registration for the User in the Application they are logging into.
    - `enabled` - (Optional) Determines if this identity provider is enabled for the Application specified by the `application_id` property.
    - `fields` - (Optional) This is an optional Application specific override for the top level `fields`.
    - `permissions` - (Optional) This is an optional Application specific override for the top level `permissions`.
* `button_text` - (Required) The top-level button text to use on the FusionAuth login page for this Identity Provider.
* `client_id` - (Required) The top-level LinkedIn client id for your Application. This value is retrieved from the LinkedIn developer website when you set up your LinkedIn app.
* `client_secret` - (Required) The top-level client secret to use with the LinkedIn Identity Provider when retrieving the long-lived token. This value is retrieved from the LinkedIn developer website when you set up your LinkedIn app.
* `debug` - (Optional) Determines if debug is enabled for this provider. When enabled, an Event Log is created each time this provider is invoked to reconcile a login.
* `enabled` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
* `lambda_reconcile_id` - (Optional) The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user.
* `linking_strategy` - (Optional) The linking strategy to use when creating the link between the Facebook Identity Provider and the user.
  The valid values are:
    - `CreatePendingLink` - Do not automatically link, instead return a pending link identifier that can be used to link to an existing user.
    - `LinkAnonymously` - Always create a link based upon the unique Id returned by the identity provider. A username or email is not required and will not be used to link the user. A reconcile lambda will not be used in this configuration.
    - `LinkByEmail` - Link to an existing user based upon email. A user will be created with the email returned by the identity provider if one does not already exist.
    - `LinkByEmailForExistingUser` - Only link to an existing user based upon email. A user will not be created if one does not already exist with email returned by the identity provider.
    - `LinkByUsername` - Link to an existing user based upon username. A user will be created with the username returned by the identity provider if one does not already exist.
    - `LinkByUsernameForExistingUser` - Only link to an existing user based upon username. A user will not be created if one does not already exist with username returned by the identity provider.
* `scope` - (Optional) The top-level scope that you are requesting from LinkedIn.
* `tenant_configuration` - (Optional) The configuration for each Tenant that limits the number of links a user may have for a particular identity provider.
    - `tenant_id` - (Optional) The unique Id of the tenant that this configuration applies to.
    - `limit_user_link_count_enabled` - (Optional) When enabled, the number of identity provider links a user may create is enforced by maximumLinks.
    - `limit_user_link_count_maximum_links` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
