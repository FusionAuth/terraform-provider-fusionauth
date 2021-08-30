# Form Resource

A FusionAuth Form is a customizable object that contains one-to-many ordered steps. Each step is comprised of one or more Form Fields.

[Form API](https://fusionauth.io/docs/v1/tech/apis/forms/)

## Example Usage

```hcl
resource "fusionauth_form" "form" {
  data = {
    "description" : "This form customizes the registration experience."
  }
  name = "Custom Registration Form"
  steps {
    fields = ["91909721-7d4f-b110-8f21-cfdee2a1edb8"]
  }
  steps {
    fields = ["8ed89a31-c325-3156-72ed-6e89183af917", "a977cfd4-a9ed-c4cf-650f-f4539268ac38"]
  }
}
```

## Argument Reference

* `form_id` - (Optional) The Id to use for the new Form. If not specified a secure random UUID will be generated.
* `data` - (Optional) An object that can hold any information about the Form Field that should be persisted.
* `name` - (Required) The unique name of the Form Field.
* `steps` - (Required) An ordered list of objects containing one or more Form Fields. A Form must have at least one step defined.
    - `fields` - (Required) An ordered list of Form Field Ids assigned to this step.
* `type` - (Optional) The type of form being created, a form type cannot be changed after the form has been created.