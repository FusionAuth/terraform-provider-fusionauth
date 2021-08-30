# Application Resource

[Applications API](https://fusionauth.io/docs/v1/tech/apis/applications)

## Example Usage

```hcl
data "fusionauth_application" "FusionAuth"{
    name = "FusionAuth"
}
```

## Argument Reference

* `name` - (Required) The name of the Application.