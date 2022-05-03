# Google Identity Provider Resource

The Google identity provider type will use the Google OAuth v2.0 login API. it will provide a Login with Google button on FusionAuth’s login page that will leverage the Google login pop-up dialog. Additionally, this identity provider will call Google’s /oauth2/v3/tokeninfo API to load additional details about the user and store them in FusionAuth.

The email address returned by the Google Token info API will be used to create or lookup the existing user. Additional claims returned by Google can be used to reconcile the User to FusionAuth by using a Google Reconcile Lambda. Unless you assign a reconcile lambda to this provider, on the email address will be used from the available claims returned by Google.

FusionAuth will also store the Google access_token that is returned from the login pop-up dialog in the UserRegistration object inside the tokens Map. This Map stores the tokens from the various identity providers so that you can use them in your application to call their APIs.

[Google Identity Providers API](https://fusionauth.io/docs/v1/tech/apis/identity-providers/google#create-the-google-identity-provider)

## Example Usage

```hcl
resource "fusionauth_idp_google" "google" {
    application_configuration {
        application_id = fusionauth_application.myapp.id
        create_registration = true
        enabled = true
    }
    button_text = "Login with Google"
    debug = false
    client_id = "254311943570-8e2i2hds0qdnee4124socceeh2q2mtjl.apps.googleusercontent.com"
    client_secret ="BRr7x7xz_-cXxIFznBDIdxF1"
    scope = "profile"
}
```

## Argument Reference

* `application_configuration` - (Optional) The configuration for each Application that the identity provider is enabled for.
    - `application_id` - (Optional) ID of the Application to apply this configuration to.
    - `button_text` - (Optional) This is an optional Application specific override for the top level button text.
    - `client_id` - (Optional) This is an optional Application specific override for the top level client id.
    - `client_secret` - (Optional) This is an optional Application specific override for the top level client secret.
    - `create_registration` - (Optional) Determines if a UserRegistration is created for the User automatically or not. If a user doesn’t exist in FusionAuth and logs in through an identity provider, this boolean controls whether or not FusionAuth creates a registration for the User in the Application they are logging into.
    - `enabled` - (Optional) Determines if this identity provider is enabled for the Application specified by the applicationId key.
    - `scope` - (Optional) This is an optional Application specific override for for the top level scope.
* `button_text` - (Required) The top-level button text to use on the FusionAuth login page for this Identity Provider.
* `client_id` - (Required) The top-level Google client id for your Application. This value is retrieved from the Google developer website when you setup your Google developer account.
* `client_secret` - (Optional) The top-level client secret to use with the Google Identity Provider when retrieving the long-lived token. This value is retrieved from the Google developer website when you setup your Google developer account.
* `debug` - (Optional) Determines if debug is enabled for this provider. When enabled, each time this provider is invoked to reconcile a login an Event Log will be created.
* `domains` - (Optional) This is an optional list of domains that this OpenID Connect provider should be used for. This converts the FusionAuth login form to a domain-based login form. This type of form first asks the user for their email. FusionAuth then uses their email to determine if an OpenID Connect identity provider should be used. If an OpenID Connect provider should be used, the browser is redirected to the authorization endpoint of that identity provider. Otherwise, the password field is revealed on the form so that the user can login using FusionAuth.
* `enabled` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
* `lambda_reconcile_id` - (Optional) The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user.
* `scope` - (Optional) The top-level scope that you are requesting from Google.
* `linking_strategy` - (Optional) The linking strategy to use when creating the link between the {idp_display_name} Identity Provider and the user.
* `login_method` - (Optional) The login method to use for this Identity Provider.
* `tenant_configuration` - (Optional) The configuration for each Tenant that limits the number of links a user may have for a particular identity provider.
    - `tenant_id` - (Optional) The unique Id of the tenant that this configuration applies to.
    - `limit_user_link_count_enabled` - (Optional) When enabled, the number of identity provider links a user may create is enforced by maximumLinks.
    - `limit_user_link_count_maximum_links` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
