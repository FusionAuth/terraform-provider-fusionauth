# Form Resource

A FusionAuth Form is a customizable object that contains one-to-many ordered steps. Each step is comprised of one or more Form Fields.

[Forms API](https://fusionauth.io/docs/v1/tech/apis/forms)

## Example Usage

```hcl
data "fusionauth_form" "default" {
    name = "Default User Self Service provided by FusionAuth"
}
```

## Argument Reference

* `form_id` - (Optional) The unique id of the Form. Either `form_id` or `name` must be specified.
* `name` - (Optional) The name of the Form. Either `form_id` or `name` must be specified.

## Attributes Reference

All the argument attributes are also exported as result attributes.

The following additional attributes are exported:

* `id` - The unique Id of the Form.
* `data` - An object that can hold any information about the Form that should be persisted.
* `name` - The unique name of the Form.
* `steps` - An ordered list of objects containing one or more Form Fields.
* `type` - The form type. The possible values are:
    * `registration` - This form will be used for self service registration.
    * `adminRegistration` - This form be used to customize the add and edit User Registration form in the FusionAuth UI.
    * `adminUser` - This form can be used to customize the add and edit User form in the FusionAuth UI.
    * `selfServiceUser` - This form will be used to for self service user management.
