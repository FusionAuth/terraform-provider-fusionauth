# OpenID Connect Identity Provider Resource

OpenID Connect identity providers connect to external OpenID Connect login systems. This type of login will optionally provide a Login with …​ button on FusionAuth’s login page. This button is customizable by using different properties of the identity provider.

Optionally, this identity provider can define one or more domains it is associated with. This is useful for allowing employees to log in with their corporate credentials. As long as the company has an identity solution that provides OpenID Connect, you can leverage this feature. This is referred to as a Domain Based Identity Provider. If you enable domains for an identity provider, the Login with …​ button will not be displayed. Instead, only the email form field will be displayed initially on the FusionAuth login page. Once the user types in their email address, FusionAuth will determine if the user is logging in locally or if they should be redirected to this identity provider. This is determined by extracting the domain from their email address and comparing it to the domains associated with the identity provider.

FusionAuth will also leverage the /userinfo API that is part of the OpenID Connect specification. The email address returned from the Userinfo response will be used to create or lookup the existing user. Additional claims from the Userinfo response can be used to reconcile the User in FusionAuth by using an OpenID Connect Reconcile Lambda. Unless you assign a reconcile lambda to this provider, on the email address will be used from the available claims returned by the OpenID Connect identity provider.

If the external OpenID Connect identity provider returns a refresh token, it will be stored in the UserRegistration object inside the tokens Map. This Map stores the tokens from the various identity providers so that you can use them in your application to call their APIs.

[OpenID Connect Identity Providers API](https://fusionauth.io/docs/v1/tech/apis/identity-providers/openid-connect)

## Example Usage

```hcl
resource "fusionauth_idp_open_id_connect" "OpenID" {
    application_configuration {
        application_id      = fusionauth_application.myapp.id
        create_registration = true
        enabled             = true
    }
    oauth2_authorization_endpoint       = "https://acme.com/oauth2/authorization"
    oauth2_client_id                    = "191c23dc-b772-4558-bd21-dc1cbf74ae21"
    oauth2_client_secret                ="SUsnoP0pWUYfXvWbSe5pvj8Di5nAxOvO"
    oauth2_client_authentication_method = "client_secret_basic"
    oauth2_scope                        = "openid offline_access"
    oauth2_token_endpoint               = "https://acme.com/oauth2/token"
    oauth2_user_info_endpoint           = "https://acme.com/oauth2/userinfo"
    button_text                         = "Login with OpenID Connect"
    debug                               = false
    enabled                             = true
    name                                = "Super Awesome OpenID Connect Provider"
    tenant_configuration {
      tenant_id                           = fusionauth_tenant.example.id
      limit_user_link_count_enabled       = false
      limit_user_link_count_maximum_links = 42
    }
}
```

## Argument Reference

* `idp_id` - (Optional) The ID to use for the new identity provider. If not specified a secure random UUID will be generated.
* `application_configuration` - (Optional) The configuration for each Application that the identity provider is enabled for.
    - `application_id` - (Optional) ID of the Application to apply this configuration to.
    - `button_image_url` - (Optional) This is an optional Application specific override for the top level button image URL.
    - `button_text` - (Optional) This is an optional Application specific override for the top level button text.
    - `oauth2_client_id` - (Optional) This is an optional Application specific override for the top level client id.
    - `oauth2_client_secret` - (Optional) This is an optional Application specific override for the top level client secret.
    - `create_registration` - (Optional) Determines if a UserRegistration is created for the User automatically or not. If a user doesn’t exist in FusionAuth and logs in through an identity provider, this boolean controls whether or not FusionAuth creates a registration for the User in the Application they are logging into.
    - `enabled` - (Optional) Determines if this identity provider is enabled for the Application specified by the applicationId key.
    - `oauth2_scope` - (Optional) This is an optional Application specific override for the top level scope.
* `button_image_url` - (Optional) The top-level button image (URL) to use on the FusionAuth login page for this Identity Provider.
* `button_text` - (Required) The top-level button text to use on the FusionAuth login page for this Identity Provider.
* `debug` - (Optional) Determines if debug is enabled for this provider. When enabled, each time this provider is invoked to reconcile a login an Event Log will be created.
* `domains` - (Optional) This is an optional list of domains that this OpenID Connect provider should be used for. This converts the FusionAuth login form to a domain-based login form. This type of form first asks the user for their email. FusionAuth then uses their email to determine if an OpenID Connect identity provider should be used. If an OpenID Connect provider should be used, the browser is redirected to the authorization endpoint of that identity provider. Otherwise, the password field is revealed on the form so that the user can login using FusionAuth.
* `enabled` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
* `lambda_reconcile_id` - (Optional) The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user.
* `linking_strategy` - (Optional) The linking strategy to use when creating the link between the {idp_display_name} Identity Provider and the user.
* `name` - (Required) The name of this OpenID Connect identity provider. This is only used for display purposes.
* `oauth2_authorization_endpoint` - (Optional) The top-level authorization endpoint for the OpenID Connect identity provider. You can leave this blank if you provide the issuer field, which will be used to make a request to the OpenID Connect .well-known endpoint in order to dynamically resolve the authorization endpoint. If you provide an issuer then this field will be ignored.
* `oauth2_client_id` - (Required) The top-level client id for your Application.
* `oauth2_client_secret` - (Optional) The top-level client secret to use with the OpenID Connect identity provider.
* `oauth2_client_authentication_method` - (Optional) The client authentication method to use with the OpenID Connect identity provider. 
* `oauth2_email_claim` - (Optional) An optional configuration to modify the expected name of the claim returned by the IdP that contains the email address.
* `oauth2_unique_id_claim` - (Optional) An optional configuration to modify the expected name of the claim returned by the IdP that contains the user Id.
* `oauth2_username_claim` - (Optional) An optional configuration to modify the expected name of the claim returned by the IdP that contains the username.
* `oauth2_issuer` - (Optional) The top-level issuer URI for the OpenID Connect identity provider. If this is provided, the authorization endpoint, token endpoint and userinfo endpoint will all be resolved using the issuer URI plus /.well-known/openid-configuration.
* `oauth2_scope` - (Optional) The top-level scope that you are requesting from the OpenID Connect identity provider.
* `oauth2_token_endpoint` - (Optional) The top-level token endpoint for the OpenID Connect identity provider. You can leave this blank if you provide the issuer field, which will be used to make a request to the OpenID Connect .well-known endpoint in order to dynamically resolve the token endpoint. If you provide an issuer then this field will be ignored.
* `oauth2_user_info_endpoint` - (Optional) The top-level userinfo endpoint for the OpenID Connect identity provider. You can leave this blank if you provide the issuer field, which will be used to make a request to the OpenID Connect .well-known endpoint in order to dynamically resolve the userinfo endpoint. If you provide an issuer then this field will be ignored.
* `post_request` - (Optional) Set this value equal to true if you wish to use POST bindings with this OpenID Connect identity provider. The default value of false means that a redirect binding which uses a GET request will be used.
* `tenant_configuration` - (Optional) The configuration for each Tenant that limits the number of links a user may have for a particular identity provider.
    - `tenant_id` - (Optional) The unique Id of the tenant that this configuration applies to.
    - `limit_user_link_count_enabled` - (Optional) When enabled, the number of identity provider links a user may create is enforced by maximumLinks.
    - `limit_user_link_count_maximum_links` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
