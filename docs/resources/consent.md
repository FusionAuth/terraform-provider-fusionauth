# Consent Resource

Consents are used to define and manage user consent for data usage within your applications. They can enforce age restrictions, require parent/guardian approval for minors, and track what users have agreed to.

[Consents API](https://fusionauth.io/docs/v1/tech/apis/consents)

## Example Usage

```hcl
resource "fusionauth_consent" "example" {
  name                                = "Newsletter Signup"
  default_minimum_age_for_self_consent = 16
  country_minimum_age_for_self_consent = {
    "us" = 13
    "uk" = 16
    "de" = 17
  }
  consent_email_template_id = "04881b0a-f543-4218-ae1d-ee009591a4b4"
  multiple_values_allowed   = true
  values = [
    "Email",
    "SMS",
    "Phone"
  ]

  email_plus {
    enabled                             = true
    email_template_id                   = "847de2b9-14ba-45fc-9726-9d9c9fcee4b6"
    minimum_time_to_send_email_in_hours = 24
    maximum_time_to_send_email_in_hours = 48
  }

  data = jsonencode({
    description = "Marketing consent for newsletter"
    internal_id = "MKT-001"
  })
}
```

## Argument Reference

* `name` - (Required) The unique name of the consent.
* `default_minimum_age_for_self_consent` - (Required) The default age of self consent used when granting this consent to a user unless a more specific one is provided by the `country_minimum_age_for_self_consent`.

---

* `consent_email_template_id` - (Optional) The Id of the Email Template that is used to send confirmation to the end user. If omitted, an email will not be sent to the user.
* `consent_id` - (Optional) The Id to use for the new Consent. If not specified a secure random UUID will be generated.
* `country_minimum_age_for_self_consent` - (Optional) This property optionally overrides the value provided in `default_minimum_age_for_self_consent` if a more specific value is defined. This can be useful when the age of self consent varies by country.
* `data` - (Optional) An object that can hold any information about the Consent that should be persisted. Must be a JSON string.
* `email_plus` - (Optional) Email Plus provides an additional opportunity to notify the giver that consent was provided. Configuration block detailed below.
* `multiple_values_allowed` - (Optional) Set to `true` if more than one value may be used when granting this consent to a User. Default is `false`.
* `values` - (Optional) The list of values that may be assigned to this consent. Required when `multiple_values_allowed` is set to `true`.

### Email Plus Configuration Block

* `email_template_id` - (Required) The Id of the Email Template that is used to send the reminder notice to the consent giver.
* `enabled` - (Optional) Set to `true` to enable the Email Plus workflow. Default is `false`.
* `maximum_time_to_send_email_in_hours` - (Optional) The maximum number of hours to wait until sending the reminder notice to the consent giver. Default is `48`.
* `minimum_time_to_send_email_in_hours` - (Optional) The minimum number of hours to wait until sending the reminder notice to the consent giver. Default is `24`.

## Import

Consents can be imported using the consent ID:

```shell
terraform import fusionauth_consent.example 34ee46ad-60a7-44f6-8edf-d818c2209072
```
