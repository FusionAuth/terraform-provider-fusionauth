# SAML v2 IdP InitiatedIdentity Provider Resource

The SAML v2 IdP Initiated IdP initiated Identity Provider allows an external IdP to send an unsolicited AuthN request when FusionAuth is acting as the Service Provider (or SP).

[SAML v2 IdP Initiated Identity Provider API](https://fusionauth.io/docs/v1/tech/apis/identity-providers/samlv2-idp-initiated/)

## Example Usage

```hcl
resource "fusionauth_idp_saml_v2_idp_initated" "Saml" {
  application_configuration {
    application_id      = fusionauth_application.myapp.id
    create_registration = true
    enabled             = true
  }
  debug               = false
  email_claim         = "email"
  issuer        = "https://www.example.com/login"
  name                = "My SAML provider"
  use_name_for_email  = true
}
```

## Argument Reference

* `application_configuration` - (Optional) The configuration for each Application that the identity provider is enabled for.
  * `application_id` - (Optional) ID of the Application to apply this configuration to.
  * `create_registration` - (Optional) Determines if a UserRegistration is created for the User automatically or not. If a user doesn't exist in FusionAuth and logs in through an identity provider, this boolean controls whether or not FusionAuth creates a registration for the User in the Application they are logging into.
  * `enabled` - (Optional) Determines if this identity provider is enabled for the Application specified by the applicationId key.
* `assertion_configuration` - (Optional) The configuration for the SAML assertion.
  * `decryption` - (Optional) The configuration for the SAML assertion decryption.
    * `enabled` - (Optional) Determines if FusionAuth requires encrypted assertions in SAML responses from the identity provider. When true, SAML responses from the identity provider containing unencrypted assertions will be rejected by FusionAuth.
    * `key_transport_decryption_key_id` - (Optional) The Id of the key stored in Key Master that is used to decrypt the symmetric key on the SAML response sent to FusionAuth from the identity provider. The selected Key must contain an RSA private key. Required when `enabled` is true.
* `debug` - (Optional) Determines if debug is enabled for this provider. When enabled, each time this provider is invoked to reconcile a login an Event Log will be created.
* `email_claim` - (Optional) The name of the email claim (Attribute in the Assertion element) in the SAML response that FusionAuth uses to uniquely identity the user. If this is not set, the `use_name_for_email` flag must be true.
* `enabled` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
* `idp_id` - (Optional) The ID to use for the new identity provider. If not specified a secure random UUID will be generated.
* `issuer` - (Required) The EntityId (unique identifier) of the SAML v2 identity provider. This value should be provided to you. Prior to 1.27.1 this value was required to be a URL.
* `key_id` - (Required) The id of the key stored in Key Master that is used to verify the SAML response sent back to FusionAuth from the identity provider. This key must be a verification only key or certificate (meaning that it only has a public key component).
* `lambda_reconcile_id` - (Optional) The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user.
* `linking_strategy` - (Optional) The linking strategy to use when creating the link between the {idp_display_name} Identity Provider and the user.
* `name` - (Required) The name of this OpenID Connect identity provider. This is only used for display purposes.
* `tenant_configuration` - (Optional) The configuration for each Tenant that limits the number of links a user may have for a particular identity provider.
  * `tenant_id` - (Optional) The unique Id of the tenant that this configuration applies to.
  * `limit_user_link_count_enabled` - (Optional) When enabled, the number of identity provider links a user may create is enforced by maximumLinks.
  * `limit_user_link_count_maximum_links` - (Optional) Determines if this provider is enabled. If it is false then it will be disabled globally.
* `unique_id_claim` - (Optional) The name of the unique claim in the SAML response that FusionAuth uses to uniquely link the user. If this is not set, the `email_claim` will be used when linking user.
* `use_name_for_email` - (Optional) Whether or not FusionAuth will use the NameID element value as the email address of the user for reconciliation processing. If this is false, then the `email_claim` property must be set.
* `username_claim` - (Optional) The name of the claim in the SAML response that FusionAuth uses to identity the username. If this is not set, the NameId value will be used to link a user. This property is required when `linking_strategy` is set to LinkByUsername or LinkByUsernameForExistingUser
