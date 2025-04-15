# SMS Message Template Resource

A FusionAuth SMS Message Template is a named object that provides configuration for the content of SMS messages sent through configured messengers.

[Message Templates API](https://fusionauth.io/docs/v1/tech/apis/message-templates/)

## Example Usage

```hcl
resource "fusionauth_sms_message_template" "two_factor" {
  name             = "Two Factor Authentication"
  default_template = "Here's your Two Factor Code: ${code}"

  data = jsonencode({
    updatedBy = "admin@example.com"
  })

  localized_templates = {
    de = "Hier ist Ihr Zwei-Faktoren-Code: ${code}"
    es = "Este es su código de dos factores: ${code}"
    fr = "Voici votre code à deux facteurs: ${code}"
  }
}
```

## Argument Reference

* `name` - (Required) A descriptive name for the Message Template (i.e. "Two Factor Code Message").
* `default_template` - (Required) The default Message Template content that will be used when sending SMS messages.

---

* `message_template_id` - (Optional) The Id to use for the new Message Template. If not specified a secure random UUID will be generated.
* `data` - (Optional) A JSON string that can hold any information about the Message Template that should be persisted. Must be a JSON string.
* `localized_templates` - (Optional) A map of language codes to template strings used when sending messages to users who speak other languages. This overrides the default Message Template based on the user's list of preferred languages.

## Attribute Reference

* `type` - The type of Message Template. This is always 'SMS'.

## Import

SMS Message Templates can be imported using the template ID:

```shell
terraform import fusionauth_sms_message_template.example 00000000-0000-0000-0000-000000000000
