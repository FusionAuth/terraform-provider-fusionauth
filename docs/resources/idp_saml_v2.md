# SAML v2 Identity Provider Resource

SAML v2 identity providers connect to external SAML v2 login systems. This type of login will optionally provide a Login with …​ button on FusionAuth’s login page. This button is customizable by using different properties of the identity provider.

Optionally, this identity provider can define one or more domains it is associated with. This is useful for allowing employees to log in with their corporate credentials. As long as the company has an identity solution that provides SAML v2, you can leverage this feature. This is referred to as a Domain Based Identity Provider. If you enable domains for an identity provider, the Login with …​ button will not be displayed. Instead, only the email form field will be displayed initially on the FusionAuth login page. Once the user types in their email address, FusionAuth will determine if the user is logging in locally or if they should be redirected to this identity provider. This is determined by extracting the domain from their email address and comparing it to the domains associated with the identity provider.

FusionAuth will locate the user’s email address in the SAML assertion which will be used to create or lookup the existing user. Additional claims from the SAML response can be used to reconcile the User to FusionAuth by using a SAML v2 Reconcile Lambda. Unless you assign a reconcile lambda to this provider, on the email address will be used from the available assertions returned by the SAML v2 identity provider.

[SAML v2 Connect Identity Providers API](https://fusionauth.io/docs/v1/tech/apis/identity-providers/samlv2/)

## Example Usage

```hcl
resource "fusionauth_idp_saml_v2" "Saml" {
  application_configuration {
    application_id      = fusionauth_application.myapp.id
    button_text         = "Login with SAML (app text)"
    create_registration = true
    enabled             = true
  }
  button_text         = "Login with SAML"
  debug               = false
  email_claim         = "email"
  idp_endpoint        = "https://www.example.com/login"
  name                = "My SAML provider"
  post_request        = true
  request_signing_key = "3168129b-91fa-46f4-9676-947f5509fdce"
  sign_request        = true
  use_name_for_email  = true
}
```

## Argument Reference

* `button_text` - (Required) The top-level button text to use on the FusionAuth login page for this Identity Provider.
* `key_id` - (Required) The id of the key stored in Key Master that is used to verify the SAML response sent back to FusionAuth from the identity provider. This key must be a verification only key or certificate (meaning that it only has a public key component).
* `name` - (Required) The name of this OpenID Connect identity provider. This is only used for display purposes.

---

* `application_configuration` - (Optional) The configuration for each Application that the identity provider is enabled for.
  * `application_id` - (Optional) ID of the Application to apply this configuration to.
  * `button_image_url` - (Optional) This is an optional Application specific override for the top level button image URL.
  * `button_text` - (Optional) This is an optional Application specific override for the top level button text.
  * `create_registration` - (Optional) Determines if a UserRegistration is created for the User automatically or not. If a user doesn’t exist in FusionAuth and logs in through an identity provider, this boolean controls whether or not FusionAuth creates a registration for the User in the Application they are logging into.
  * `enabled` - (Optional) Determines if this identity provider is enabled for the Application specified by the applicationId key.
* `assertion_configuration` - (Optional) The configuration for the SAML assertion.
  * `destination` - (Optional) The array of URLs that FusionAuth will accept as SAML login destinations if the `policy` setting is AllowAlternates.
    * `alternates` - (Optional) The alternate destinations of the assertion.
    * `policy` - (Optional) The policy to use when performing a destination assertion on the SAML login request. The possible values are `Enabled`, `Disabled`, and `AllowAlternates`.
  * `decryption` - (Optional) The configuration for the SAML assertion decryption.
    * `enabled` - (Optional) Determines if FusionAuth requires encrypted assertions in SAML responses from the identity provider. When true, SAML responses from the identity provider containing unencrypted assertions will be rejected by FusionAuth.
    * `key_transport_decryption_key_id` - (Optional) The Id of the key stored in Key Master that is used to decrypt the symmetric key on the SAML response sent to FusionAuth from the identity provider. The selected Key must contain an RSA private key. Required when `enabled` is true.
* `button_image_url` - (Optional) The top-level button image (URL) to use on the FusionAuth login page for this Identity Provider.
* `debug` - (Optional) Determines if debug is enabled for this provider. When enabled, each time this provider is invoked to reconcile a login an Event Log will be created.
* `domains` - (Optional) This is an optional list of domains that this OpenID Connect provider should be used for. This converts the FusionAuth login form to a domain-based login form. This type of form first asks the user for their email. FusionAuth then uses their email to determine if an OpenID Connect identity provider should be used. If an OpenID Connect provider should be used, the browser is redirected to the authorization endpoint of that identity provider. Otherwise, the password field is revealed on the form so that the user can login using FusionAuth.
* `email_claim` - (Optional) The name of the email claim (Attribute in the Assertion element) in the SAML response that FusionAuth uses to uniquely identity the user. If this is not set, the `use_name_for_email` flag must be true.
* `enabled` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
* `idp_endpoint` - (Optional) The SAML v2 login page of the identity provider.
* `idp_id` - (Optional) The ID to use for the new identity provider. If not specified a secure random UUID will be generated.
* `idp_initiated_configuration` - (Optional) The configuration for the IdP initiated login.
  * `enabled` - (Optional) Determines if FusionAuth will accept IdP initiated login requests from this SAMLv2 Identity Provider.
  * `issuer` - (Optional)The EntityId (unique identifier) of the SAML v2 identity provider. This value should be provided to you. Required when `enabled` is true.
* `lambda_reconcile_id` - (Optional) The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user.
* `linking_strategy` - (Optional) The linking strategy to use when creating the link between the {idp_display_name} Identity Provider and the user.
* `login_hint_configuration` - (Optional) The configuration for the login hint.
  * `enabled` - (Optional) When enabled and HTTP-Redirect bindings are used, FusionAuth will provide the username or email address when available to the IdP as a login hint using the configured parameter name set by the `parameter_name` to initiate the AuthN request.
  * `parameter_name` - (Optional) The name of the parameter used to pass the username or email as login hint to the IDP when enabled, and HTTP redirect bindings are used to initiate the AuthN request. The default value is `login_hint`. Required when `enabled` is true.
* `name` - (Optional) The name of the provider. This is only used for display purposes. The display name of this provider instance. Required when using a provided `tenant_id` or `identity_provider.tenant_id`.
* `name_id_format` - (Optional) The Name Id is used to facilitate communication about a user between a Service Provider (SP) and Identity Provider (IdP). The SP can specify the preferred format in the AuthN request regarding a user. The identity Provider will attempt to honor this format in the AuthN response. When this parameter is omitted a default value of urn:oasis:names:tc:SAML:2.0:nameid-format:persistent will be used.
* `post_request` - (Optional) When true the authentication request will use the HTTP POST binding with the identity provider instead of the default Redirect binding which uses the HTTP GET method.
* `post_request` - (Optional) Set this value equal to true if you wish to use POST bindings with this OpenID Connect identity provider. The default value of false means that a redirect binding which uses a GET request will be used.
* `request_signing_key` - (Optional) The key pair Id to use to sign the SAML request. Required when `sign_request` is true.
* `sign_request` - (Optional) When true authentication requests sent to the identity provider will be signed.
* `tenant_configuration` - (Optional) The configuration for each Tenant that limits the number of links a user may have for a particular identity provider.
  * `tenant_id` - (Optional) The unique Id of the tenant that this configuration applies to.
    * `limit_user_link_count_enabled` - (Optional) When enabled, the number of identity provider links a user may create is enforced by maximumLinks.
    * `limit_user_link_count_maximum_links` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
* `tenant_id` - (Optional) The unique Id of the Tenant. Providing a value creates an identity provider scoped to the specified tenant, otherwise a global identity provider is created. Tenant-scoped identity providers can only be used to authenticate in the context of the specified tenant. Global identity providers can be used with any tenant. This value cannot be updated after creation and requires recreating the resource to change.
* `unique_id_claim` - (Optional) The name of the unique claim in the SAML response that FusionAuth uses to uniquely link the user. If this is not set, `the email_claim` will be used when linking user.
* `use_name_for_email` - (Optional) Whether or not FusionAuth will use the NameID element value as the email address of the user for reconciliation processing. If this is false, then the `email_claim` property must be set.
* `username_claim` - (Optional) The name of the claim in the SAML response that FusionAuth uses to identify the username. If this is not set, the NameId value will be used to link a user. This property is required when linkingStrategy is set to LinkByUsername or LinkByUsernameForExistingUser.
* `xml_signature_canonicalization_method` - (Optional) The XML signature canonicalization method used when digesting and signing the SAML request.
