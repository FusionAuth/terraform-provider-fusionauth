# Twilio Messenger Data Source

This data source can be used to fetch information about a specific Twilio Messenger.

[Messengers API](https://fusionauth.io/docs/apis/messengers/twilio)

## Example Usage

```hcl
data "fusionauth_twilio_messenger" "example" {
  id = "75a068fd-e94b-451a-9aeb-3ddb9a3b5987"
}

data "fusionauth_twilio_messenger" "example" {
  name = "My Twilio Messenger"
}
```

## Argument Reference

* `id` - (Optional) The unique Id of the Twilio Messenger to retrieve. If this is not specified, the `name` argument must be specified.
* `name` - (Optional) The case-insensitive string to search for in the Twilio Messenger name. If this is not specified, the `id` argument must be specified.

## Attributes Reference

* `account_sid` - The Twilio Account ID used when connecting to the Twilio API.
* `auth_token` - The Twilio Auth Token used when connecting to the Twilio API.
* `data` - An object that can hold any information about the Twilio Messenger that should be persisted. Represented as a JSON string.
* `debug` - Determines if debug is enabled to create an event log to assist in debugging messenger errors.
* `from_phone_number` - The configured Twilio phone number used to send messages.
* `messaging_service_sid` - The Twilio message service Id used when using Twilio Copilot to load balance between numbers.
* `url` - The Twilio URL that FusionAuth uses to communicate with the Twilio API.
