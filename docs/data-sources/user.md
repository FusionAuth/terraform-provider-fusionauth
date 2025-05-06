# User Data Source

This data source can be used to fetch information about a specific user.

[Users API](https://fusionauth.io/docs/v1/tech/apis/users)

## Example Usage

```hcl
# Fetch user by username
data "fusionauth_user" "example" {
  username = "foo@example.com"
}
```

## Argument Reference

* `tenant_id` - (Optional) The Id of the tenant used to scope this API request.
* `user_id` - (Optional) The Id of the user. Either `user_id` or `username` must be specified.
* `username` - (Optional) The username of the user. Either `user_id` or `username` must be specified.

## Attributes Reference

All of the argument attributes are also exported as result attributes.

The following additional attributes are exported:

* `active` - True if the user is active. False if the user has been deactivated. Deactivated users will not be able to login.
* `birth_date` - An ISO-8601 formatted date of the user’s birthdate such as YYYY-MM-DD.
* `data` - A JSON string that can hold any information about the user.
* `email` - The user’s email address.
* `expiry` - The expiration instant of the user’s account. An expired user is not permitted to login.
* `first_name` - The first name of the user.
* `full_name` - The user’s full name.
* `identities` - The list of identities that exist for a User.
  * `display_value` - The display value for the identity. Only used for username type identities. If the unique username feature is not enabled, this value will be the same as user.identities[x].value. Otherwise, it will be the username the User has chosen. For primary username identities, this will be the same value as user.username.
  * `insert_instant` - The instant when the identity was created.
  * `last_login_instant` - The instant when the identity was last used to log in. If a User has multiple identity types (username, email, and phoneNumber), then this value will represent the specific identity they last used to log in. This contrasts with user.lastLoginInstant, which represents the last time any of the User’s identities was used to log in.
  * `last_update_instant` - The instant when the identity was last updated.
  * `moderation_status` - The current status of the username. This is used if you are moderating usernames via CleanSpeak.
  * `type` - he identity type.
  * `value` - The value represented by the identity.
  * `verified` - Whether verification was actually performed on the identity by FusionAuth.
  * `verified_instant` - The instant when verification was performed on the identity.
  * `verified_reason` - The reason the User’s identity was verified or not verified.
* `image_url` - The URL that points to an image file that is the user’s profile image.
* `last_name` - The user’s last name.
* `middle_name` - The user’s middle name.
* `mobile_phone` - The user’s mobile phone number.
* `parent_email` - The email address of the user’s parent or guardian.
* `password_change_required` - Indicates that the user’s password needs to be changed during their next login attempt.
* `phone_number` - The user’s phone number.
* `preferred_languages` - An array of locale strings that give, in order, the user’s preferred languages.
* `timezone` - The user’s preferred timezone.
* `username_status` - The current status of the username. This is used if you are moderating usernames via CleanSpeak.
* `verification_ids` - The list of all verifications that exist for a user. This includes the email and phone identities that a user may have. The values from emailVerificationId and emailVerificationOneTimeCode are legacy fields and will also be present in this list.
  * `verification_id` - A verification Id.
  * `one_time_code` - A one time code that will be paired with the verificationIds[x].id.
  * `type` - The identity type that the verification Id is for. This identity type, along with verificationIds[x].value , matches exactly one identity via user.identities[x].type.
  * `value` - The identity value that the verification Id is for. This identity value, along with verificationIds[x].type , matches exactly one identity via user.identities[x].value.
