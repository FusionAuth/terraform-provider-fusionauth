# Facebook Identity Provider Resource

The Facebook identity provider type will use the Facebook OAuth login API. It will provide a `Login with Facebook` button on FusionAuth’s login page that will leverage the Facebook login pop-up dialog. Additionally, this identity provider will call Facebook’s Graph API to load additional details about the user and store them in FusionAuth.

The email address returned by the Facebook Graph API will be used to create or lookup the existing user. Additional claims returned by Facebook can be used to reconcile the User to FusionAuth by using a Facebook Reconcile Lambda. Unless you assign a reconcile lambda to this provider, on the `email` address will be used from the available claims returned by Facebook.

When the `picture` field is not requested FusionAuth will also call Facebook’s `/me/picture` API to load the user’s profile image and store it as the `imageUrl` in FusionAuth. When the `picture` field is requested, the user’s profile image will be returned by the `/me` API and a second request to the `/me/picture` endpoint will not be required.

Finally, FusionAuth will call Facebook’s `/oauth/access_token` API to exchange the login token for a long-lived Facebook token. This token is stored in the `UserRegistration` object inside the `tokens` Map. This Map stores the tokens from the various identity providers so that you can use them in your application to call their APIs.

Please note if an `idp_hint` is appended to the OAuth Authorize endpoint, then the interaction behavior will be defaulted to `redirect`, even if popup interaction is explicitly configured.

[Facebook Identity Providers API](https://fusionauth.io/docs/v1/tech/apis/identity-providers/facebook)

## Example Usage

```hcl
resource "fusionauth_idp_facebook" "facebook" {
    application_configuration {
        application_id = fusionauth_application.myapp.id
        create_registration = true
        enabled = true
    }
    button_text = "Login with Facebook"
    debug = false
    enabled = true
    app_id = "9876543210"
    client_secret = "716a572f917640698cdb99e9d7e64115"
    fields = "email"
    permissions = "email"
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
* `app_id` - (Required) The top-level Facebook `appId` for your Application. This value is retrieved from the Facebook developer website when you setup your Facebook developer account.
* `client_secret` - (Required) The top-level client secret, also known as 'App Secret', to use with the Facebook Identity Provider when retrieving the long-lived token. This value is retrieved from the Facebook developer website when you setup your Facebook developer account.
* `debug` - (Optional) Determines if debug is enabled for this provider. When enabled, an Event Log is created each time this provider is invoked to reconcile a login.
* `enabled` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
* `fields` - (Optional) The top-level fields that you are requesting from Facebook.
  Field values are documented at [Facebook Graph API](https://developers.facebook.com/docs/graph-api/using-graph-api/)
* `lambda_reconcile_id` - (Optional) The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user.
* `linking_strategy` - (Optional) The linking strategy to use when creating the link between the Facebook Identity Provider and the user.
  The valid values are:
    - `CreatePendingLink` - Do not automatically link, instead return a pending link identifier that can be used to link to an existing user.
    - `LinkAnonymously` - Always create a link based upon the unique Id returned by the identity provider. A username or email is not required and will not be used to link the user. A reconcile lambda will not be used in this configuration.
    - `LinkByEmail` - Link to an existing user based upon email. A user will be created with the email returned by the identity provider if one does not already exist.
    - `LinkByEmailForExistingUser` - Only link to an existing user based upon email. A user will not be created if one does not already exist with email returned by the identity provider.
    - `LinkByUsername` - Link to an existing user based upon username. A user will be created with the username returned by the identity provider if one does not already exist.
    - `LinkByUsernameForExistingUser` - Only link to an existing user based upon username. A user will not be created if one does not already exist with username returned by the identity provider.
* `login_method` - (Optional) The login method to use for this Identity Provider.
  The valid values are:
    - `UsePopup` - When logging in use a popup window and the Facebook javascript library.
    - `UseRedirect` - When logging in use the Facebook OAuth redirect login flow.
* `permissions` - (Optional) The top-level permissions that your application is asking of the user’s Facebook account.
  Permission values are documented at [Facebook Login API](https://developers.facebook.com/docs/permissions/reference)
* `tenant_configuration` - (Optional) The configuration for each Tenant that limits the number of links a user may have for a particular identity provider.
    - `tenant_id` - (Optional) The unique Id of the tenant that this configuration applies to.
    - `limit_user_link_count_enabled` - (Optional) When enabled, the number of identity provider links a user may create is enforced by maximumLinks.
    - `limit_user_link_count_maximum_links` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
