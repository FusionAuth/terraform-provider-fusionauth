# System Configuration Resource

A registration is the association between a User and an Application that they log into.

[System Configuration API](https://fusionauth.io/docs/v1/tech/apis/system)

## Example Usage

```hcl
resource "fusionauth_system_configuration" "example" {
  audit_log_configuration {
    delete {
      enabled                  = true
      number_of_days_to_retain = 367
    }
  }
  cors_configuration {
    allowed_methods = ["POST", "PUT"]
  }
}
```

## Argument Reference

* `report_timezone` - (Required) The time zone used to adjust the stored UTC time when generating reports. Since reports are usually rolled up hourly, this timezone will be used for demarcating the hours.

---

* `audit_log_configuration` - (Optional)
  * `delete` - (Optional)
    * `enabled` - (Optional) Whether or not FusionAuth should delete the Audit Log based upon this configuration. When true the auditLogConfiguration.delete.numberOfDaysToRetain will be used to identify audit logs that are eligible for deletion. When this value is set to false audit logs will be preserved forever.
    * `number_of_days_to_retain` - (Optional) The number of days to retain the Audit Log.
* `cors_configuration` - (Optional)
  * `allow_credentials` - (Optional) The Access-Control-Allow-Credentials response header values as described by MDN Access-Control-Allow-Credentials.
  * `allowed_headers` - (Optional) The Access-Control-Allow-Headers response header values as described by MDN Access-Control-Allow-Headers.
  * `allowed_methods` - (Optional) The Access-Control-Allow-Methods response header values as described by MDN Access-Control-Allow-Methods.
  * `allowed_origins` - (Optional) The Access-Control-Allow-Origin response header values as described by MDN Access-Control-Allow-Origin. If the wildcard * is specified, no additional domains may be specified.
  * `debug` - (Optional) Whether or not FusionAuth will log debug messages to the event log. This is primarily useful for identifying why the FusionAuth CORS filter is rejecting a request and returning an HTTP response status code of 403.
  * `enabled` - (Optional) Whether the FusionAuth CORS filter will process requests made to FusionAuth.
  * `exposed_headers` - (Optional) The Access-Control-Expose-Headers response header values as described by MDN Access-Control-Expose-Headers.
  * `preflight_max_age_in_seconds` - (Optional) The Access-Control-Max-Age response header values as described by MDN Access-Control-Max-Age.
* `data` - (Optional) An object that can hold any information about the System that should be persisted.
* `event_log_configuration` - (Optional)
  * `number_to_retain` - (Optional) The number of events to retain. Once the the number of event logs exceeds this configured value they will be deleted starting with the oldest event logs.
* `login_record_configuration` - (Optional)
  * `delete` - (Optional)
    * `enabled` - (Optional) Whether or not FusionAuth should delete the login records based upon this configuration. When true the loginRecordConfiguration.delete.numberOfDaysToRetain will be used to identify login records that are eligible for deletion. When this value is set to false login records will be preserved forever.
    * `number_of_days_to_retain` - (Optional) The number of days to retain login records.
* `trusted_proxy_configuration` - (Optional)
  * `trust_policy` - (Optional) This setting is used to resolve the client IP address for use in logging, webhooks, and IP-based access control when an X-Forwarded-For header is provided. Because proxies are free to rewrite the X-Forwarded-For header, an untrusted proxy could write a value that allowed it to bypass IP-based ACLs, or cause an incorrect IP address to be logged or sent to a webhook. Valid values are `All` and `OnlyConfigured`.
  * `trusted` - (Optional) An array of IP addresses, representing the set of trusted upstream proxies. This value will be accepted but ignored when `trust_policy` is set to `All`. Values may be specified as IPv4, or IPv6 format, and ranges of addresses are also accepted in CIDR notation.
* `ui_configuration` - (Optional)
  * `header_color` - (Optional) A hexadecimal color to override the default menu color in the user interface.
  * `logo_url` - (Optional) A URL of a logo to override the default FusionAuth logo in the user interface.
  * `menu_font_color` - (Optional) A hexadecimal color to override the default menu font color in the user interface.
* `usage_data_configuration` - (Optional)
  * `enabled` - (Optional) Whether or not FusionAuth collects and sends usage data to improve the product.
* `webhook_event_log_configuration` - (Optional)
  * `delete` - (Optional)
    * `enabled` - (Optional) Whether or not FusionAuth should delete the webhook event logs based upon this configuration. When true the webhookEventLogConfiguration.delete.numberOfDaysToRetain will be used to identify webhook event logs that are eligible for deletion. When this value is set to false webhook event logs will be preserved forever.
    * `number_of_days_to_retain` - (Optional) The number of days to retain webhook event logs.
