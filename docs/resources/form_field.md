# Form Field Resource

A FusionAuth Form Field is an object that can be customized to receive input within a FusionAuth Form.

[Form Field API](https://fusionauth.io/docs/v1/tech/apis/form-fields/)

## Example Usage

```hcl
resource "fusionauth_form_field" "field" {
  data = {
    "leftAddOn" = "send"
  }
  description = "Information about this custom field"
  key         = "user.firstName"
  name        = "Custom first-name Form Field"
  required    = true
  confirm     = true
}
```

## Argument Reference

* `form_field_id` - (Optional) The Id to use for the new Form Field. If not specified a secure random UUID will be generated.
* `confirm` - (Optional) Determines if the user input should be confirmed by requiring the value to be entered twice. If true, a confirmation field is included.
* `consent_id` - (Optional) The Id of an existing Consent. This field will be required when the type is set to consent.
* `control` - (Optional) The Form Field control
* `data` - (Optional) An object that can hold any information about the Form Field that should be persisted.
* `description` - (Optional) A description of the Form Field.
* `key` - (Required) The key is the path to the value in the user or registration object.
* `name` - (Required) The unique name of the Form Field.
* `options` - (Optional) A list of options that are applied to checkbox, radio, or select controls.
* `required` - (Optional) Determines if a value is required to complete the form.
* `type` - (Optional) The data type used to store the value in FusionAuth.
* `validator` - (Optional)
    - `enabled` - (Optional) Determines if user input should be validated.
    - `expression` - (Optional) A regular expression used to validate user input. Must be a valid regular expression pattern.