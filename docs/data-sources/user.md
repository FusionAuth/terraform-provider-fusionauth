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
* `data` - A JSON serialised string that can hold any information about the user.
* `email` - The user’s email address.
* `expiry` - The expiration instant of the user’s account. An expired user is not permitted to login.
* `first_name` - The first name of the user.
* `full_name` - The user’s full name.
* `image_url` - The URL that points to an image file that is the user’s profile image.
* `last_name` - The user’s last name.
* `middle_name` - The user’s middle name.
* `mobile_phone` - The user’s mobile phone number.
* `parent_email` - The email address of the user’s parent or guardian.
* `password_change_required` - Indicates that the user’s password needs to be changed during their next login attempt.
* `preferred_languages` - An array of locale strings that give, in order, the user’s preferred languages.
* `timezone` - The user’s preferred timezone.
* `username_status` - The current status of the username. This is used if you are moderating usernames via CleanSpeak.
