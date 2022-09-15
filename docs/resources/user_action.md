# User Action Resource


[User Actions API](https://fusionauth.io/docs/v1/tech/apis/user-actions/)

## Example Usage

```hcl
resource "fusionauth_user_action" "example" {
  name          = "Lock user"
  temporal      = true
  prevent_login = true
}
```

## Argument Reference
* `name` - (Required) The name of this User Action.
* `user_action_id` - (Optional) The id of this User Action.
* `cancel_email_template_id` - (Optional) The Id of the Email Template that is used when User Actions are canceled.
* `end_email_template_id` - (Optional) The Id of the Email Template that is used when User Actions expired automatically (end).
* `include_email_in_event_json` - (Optional) Whether to include the email information in the JSON that is sent to the Webhook when a user action is taken.
* `localized_names` - (Optional) A mapping of localized names for this User Action. The key is the Locale and the value is the name of the User Action for that language.
* `options` - (Optional)
    - `name` - (Required) The name of this User Action Option.
    - `localized_names` - (Optional) A mapping of localized names for this User Action Option. The key is the Locale and the value is the name of the User Action Option for that language.
* `prevent_login` - (Optional) Whether or not this User Action will prevent user login. When this value is set to true the action must also be marked as a time based action. See `temporal`.
* `send_end_event` - (Optional) Whether or not FusionAuth will send events to any registered Webhooks when this User Action expires.
* `start_email_template_id` - (Optional) The Id of the Email Template that is used when User Actions are started (created).
* `temporal` - (Optional) Whether or not this User Action is time-based (temporal).
* `user_emailing_enabled` - (Optional) Whether or not email is enabled for this User Action.
* `user_notifications_enabled` - (Optional) Whether or not user notifications are enabled for this User Action.
