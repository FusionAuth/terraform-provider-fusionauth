# Consent Data Source

This data source can be used to fetch information about a specific consent.

[Consents API](https://fusionauth.io/docs/v1/tech/apis/consents)

## Example Usage

```hcl
data "fusionauth_consent" "example" {
  consent_id = "34ee46ad-60a7-44f6-8edf-d818c2209072"
}

data "fusionauth_consent" "by_name" {
  name = "Issue197"
}
```

## Argument Reference

* `consent_id` - (Optional) The unique Id of the Consent to retrieve. This is mutually exclusive with `name`.
* `name` - (Optional) The case-insensitive string to search for in the Consent name. This is mutually exclusive with `consent_id`.

## Attributes Reference

* `consent_email_template_id` - The Id of the Email Template that is used to send confirmation to the end user.
* `country_minimum_age_for_self_consent` - This property optionally overrides the value provided in defaultMinimumAgeForSelfConsent if a more specific value is defined.
* `data` - An object that can hold any information about the Consent that should be persisted. Must be a JSON string.
* `default_minimum_age_for_self_consent` - The default age of self consent used when granting this consent to a user unless a more specific one is provided by the countryMinimumAgeForSelfConsent.
* `email_plus` - Email Plus provides and additional opportunity to notify the giver that consent was provided.
  * `email_template_id` - The Id of the Email Template that is used to send the reminder notice to the consent giver.
  * `enabled` - Set to true when the Email Plus workflow is enabled.
  * `maximum_time_to_send_email_in_hours` - The maximum number of hours to wait until sending the reminder notice the consent giver.
  * `minimum_time_to_send_email_in_hours` - The minimum number of hours to wait until sending the reminder notice the consent giver.
* `multiple_values_allowed` - Set to true if more than one value may be used when granting this consent to a User.
* `values` - The list of values that may be assigned to this consent.
