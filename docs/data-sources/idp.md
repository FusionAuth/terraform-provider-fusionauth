# Application Resource

[Identiy Providers API](https://fusionauth.io/docs/v1/tech/apis/identity-providers/)

## Example Useage

```hcl
data "fusionauth_idp" "FusionAuth"{
    name = "Apple"
    type = "Apple"
}
```

## Argument Reference

* `name` - (Optional) The name of the identity provider. This is only used for display purposes. Will be the type for types: `Apple`, `Facebook`, `Google`, `HYPR`, `Twitter`
* `type` - (Optional) The type of the identity provider.