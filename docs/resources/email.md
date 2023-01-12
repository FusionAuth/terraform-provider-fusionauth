# Email Resource

This resource contains the APIs for managing Email Templates.

[Emails API](https://fusionauth.io/docs/v1/tech/apis/emails)

## Example Usage

```hcl
resource "fusionauth_email" "HelloWorld" {
  name                  = "Hello World"
  default_from_name     = "Welcome Team"
  default_html_template = file("${path.module}/email_templates/HelloWorld.html.ftl")
  default_subject       = "Hello"
  default_text_template = file("${path.module}/email_templates/HelloWorld.txt.ftl")
  from_email            = "welcome@example.com.com"
}
```

## Argument Reference

* `email_id` - (Optional) The Id to use for the new Email Template. If not specified a secure random UUID will be generated.
* `default_from_name` - (Optional) The default From Name used when sending emails. If not provided, and a localized value cannot be determined, the default value for the tenant will be used. This is the display name part of the email address ( i.e. Jared Dunn <jared@piedpiper.com>).
* `default_html_template` - (Required) The default HTML Email Template.
* `default_subject` - (Required) The default Subject used when sending emails.
* `default_text_template` - (Required) The default Text Email Template.
* `from_email` - (Optional) The email address that this email will be sent from. If not provided, the default value for the tenant will be used. This is the address part email address (i.e. Jared Dunn <jared@piedpiper.com>).
* `localized_from_names` - (Optional) The From Name used when sending emails to users who speak other languages. This overrides the default From Name based on the user’s list of preferred languages.
* `localized_html_templates` - (Optional) The HTML Email Template used when sending emails to users who speak other languages. This overrides the default HTML Email Template based on the user’s list of preferred languages.
* `localized_subjects` - (Optional) The Subject used when sending emails to users who speak other languages. This overrides the default Subject based on the user’s list of preferred languages.
* `localized_text_templates` - (Optional) The Text Email Template used when sending emails to users who speak other languages. This overrides the default Text Email Template based on the user’s list of preferred languages.
* `name` - (Required) A descriptive name for the email template (i.e. "April 2016 Coupon Email")