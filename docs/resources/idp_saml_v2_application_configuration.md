# SAML v2 Identity Provider Application Configuration Resource

Manages an application configuration for a SAML v2 identity provider.

[SAML v2 Connect Identity Providers API](https://fusionauth.io/docs/v1/tech/apis/identity-providers/samlv2/)

!> **WARNING:** You should not use the `fusionauth_idp_saml_v2_application_configuration` resource in conjunction with the [`fusionauth_idp_saml_v2`](idp_saml_v2.html) resource with _in-line application configurations_ (using the `application_configuration` argument of `fusionauth_idp_saml_v2`). Doing so may cause configuration conflicts, perpetual differences, and result in configuration being overwritten.

## Example Usage

```hcl
resource "fusionauth_idp_saml_v2_application_configuration" "myapp" {
  idp_id         = fusionauth_idp_saml_v2.Saml.id
  application_id = fusionauth_application.myapp.id

  button_text         = "Login with SAML (app text)"
  create_registration = true
  enabled             = true
}
```

## Argument Reference

* `idp_id` - (Required) ID of the SAML v2 identity provider to apply this configuration to.
* `application_id` - (Required) ID of the Application to apply this configuration to.

---

* `button_image_url` - (Optional) This is an optional Application specific override for the top level button image URL.
* `button_text` - (Optional) This is an optional Application specific override for the top level button text.
* `create_registration` - (Optional) Determines if a UserRegistration is created for the User automatically or not. If a user doesn’t exist in FusionAuth and logs in through an identity provider, this boolean controls whether or not FusionAuth creates a registration for the User in the Application they are logging into.
* `enabled` - (Optional) Determines if this identity provider is enabled for the Application specified by the applicationId key.
