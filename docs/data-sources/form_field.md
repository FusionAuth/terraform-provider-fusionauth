# Form Field Resource

A FusionAuth Form Field is an object that can be customized to receive input within a FusionAuth [Form](https://fusionauth.io/docs/v1/tech/apis/forms).

[Form Field API](https://fusionauth.io/docs/v1/tech/apis/form-fields)

## Example Usage

```hcl
data "fusionauth_form_field" "default" {
    name = "Email"
}
```

## Argument Reference

- `form_field_id` - (Optional) The unique id of the Form Field. Either `form_field_id` or `name` must be specified.
- `name` - (Optional) The name of the Form field. Either `form_field_id` or `name` must be specified.

## Attributes Reference

All the argument attributes are also exported as result attributes.

The following additional attributes are exported:

- `id` - The unique Id of the Form Field.
- `confirm` - Determines if the user input should be confirmed by requiring the value to be entered twice.
- consent_id
- control
- `data` - An object that can hold any information about the Form Field that should be persisted.
- description
- key
- `name` - The unique name of the Form Field.
- `options` - A list of options that are applied to checkbox, radio, or select controls.
- `required` - Determines if a value is required to complete the form.
- `type` - The form field type. The possible values are:
  - `bool`
  - `consent`
  - `date`
  - `email`
  - `number`
  - `string`
- `validator`
  - `enabled` - Determines if user input should be validated.
  - `expression` - A regular expression used to validate user input. Must be a valid regular expression pattern.
