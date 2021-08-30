# Registration Resource

A registration is the association between a User and an Application that they log into. 

[Registrations API](https://fusionauth.io/docs/v1/tech/apis/registrations)

## Example Usage

```hcl
resource "fusionauth_registration" "example" {
  user_id        = fusionauth_user.example.id
  application_id = data.fusionauth_application.FusionAuth.id
  roles          = ["admin"]
  username       = "theadmin"
}
```

## Argument Reference

* `user_id` - (Required) The Id of the User that is registering for the Application.
* `application_id` - (Required) The Id of the Application that this registration is for.
* `authentication_token` - (Optional) The authentication token that may be used in place of the User’s password when authenticating against this application represented by this registration. This parameter is ignored if generateAuthenticationToken is set to true and instead the value will be generated.
* `generate_authentication_token` - (Optional) Determines if FusionAuth should generate an Authentication Token for this registration.
* `data` - (Optional) An object that can hold any information about the User for this registration that should be persisted.
* `preferred_languages` - (Optional) An array of locale strings that give, in order, the User’s preferred languages for this registration. These are important for email templates and other localizable text.
* `roles` - (Optional) The list of roles that the User has for this registration.
* `timezone` - (Optional) The User’s preferred timezone for this registration. The string will be in an IANA time zone format.
* `username` - (Optional) The username of the User for this registration only.
* `skip_registration_validation` - (Optional) Indicates to FusionAuth that it should skip registration verification even if it is enabled for the Application.