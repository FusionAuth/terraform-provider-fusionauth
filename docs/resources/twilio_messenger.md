# Twilio Messenger Resource

A FusionAuth Twilio Messenger is a named object that provides configuration for sending messages through the Twilio API.

[Twilio Messenger API](https://fusionauth.io/docs/v1/tech/apis/messengers/twilio/)

## Example Usage

```hcl
resource "fusionauth_twilio_messenger" "example" {
  account_sid       = "983C6FACEBBE4D858570FADD967A9DD7"
  auth_token        = "184C73BE8E44420EBAA0BA147A61B6A9"
  data              = jsonencode({
    foo = "bar"
  })
  debug             = false
  from_phone_number = "555-555-6666"
  name              = "Twilio Messenger"
  url               = "https://api.twilio.com"
}
```

## Argument Reference

* `account_sid` - (Required) The Twilio Account ID to use when connecting to the Twilio API. This can be found in your Twilio dashboard.
* `name` - (Required) The unique Messenger name.
* `url` - (Required) The Twilio URL that FusionAuth will use to communicate with the Twilio API.

---

* `auth_token` - (Optional) The Twilio Auth Token to use when connecting to the Twilio API. This can be found in your Twilio dashboard.
* `data` - (Optional) A JSON string that can hold any information about the Twilio Messenger that should be persisted. Must be a JSON string.
* `debug` - (Optional) If debug is enabled, an event log is created to assist in debugging messenger errors. Defaults to false.
* `from_phone_number` - (Optional) The configured Twilio phone number that will be used to send messages. This can be found in your Twilio dashboard.
* `messenger_id` - (Optional) The Id to use for the new Messenger. If not specified a secure random UUID will be generated.
* `messaging_service_sid` - (Optional) The Twilio message service Id, this is used when using Twilio Copilot to load balance between numbers. This can be found in your Twilio dashboard. If this is set, the fromPhoneNumber will be ignored.
