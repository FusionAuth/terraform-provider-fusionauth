# SMS Message Template Data Source

This data source can be used to fetch information about a specific SMS Message Template.

[SMS Message Templates API](https://fusionauth.io/docs/v1/tech/apis/message-templates)

## Example Usage

```hcl
data "fusionauth_sms_message_template" "example" {
  message_template_id = "f3a91c71-ee0b-476a-92db-0a16f983a47f"
}

data "fusionauth_sms_message_template" "example" {
  name = "Example"
}
```

## Argument Reference

* `message_template_id` - (Optional) The unique Id of the SMS Message Template to retrieve. If this is not specified, the `name` argument must be specified.
* `name` - (Optional) The case-insensitive string to search for in the SMS Message Template name. If this is not specified, the `message_template_id` argument must be specified.

## Attributes Reference

* `data` - An object that can hold any information about the Message Template that should be persisted. Represented as a JSON string.
* `default_template` - The default Message Template.
* `localized_templates` - The Message Template used when sending messages to users who speak other languages. This overrides the default Message Template based on the user's list of preferred languages.
* `type` - The type of Message Template. This is always 'SMS'.
