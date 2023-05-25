# Email Resource

This data source is used to fetch information about a specific Email Template.

[Emails API](https://fusionauth.io/docs/v1/tech/apis/emails)

## Example Usage

```hcl
data "fusionauth_email" "default_breached_password" {
    name = "[FusionAuth Default] Breached Password Notification"
}
```

## Argument Reference

* `name` - (Required) The name of the Email Template.

## Attributes Reference

All the argument attributes are also exported as result attributes.

* `id` - The Id of the Email Template.
* `default_from_name` - The default From Name used when sending emails.
* `default_html_template` - The default HTML Email Template.
* `default_subject` - The default Subject used when sending emails.
* `default_text_template` - The default Text Email Template.
* `from_email` - The email address that this email will be sent from.
* `localized_from_names` - The From Name used when sending emails to users who speak other languages.
* `localized_html_templates` - The HTML Email Template used when sending emails to users who speak other languages.
* `localized_subjects` - The Subject used when sending emails to users who speak other languages.
* `localized_text_templates` - The Text Email Template used when sending emails to users who speak other languages.
