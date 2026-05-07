# Tenant Manager Configuration Resource

Settings that control the behavior of the FusionAuth Tenant Manager.

[Tenant Manager Configuration API](https://fusionauth.io/docs/apis/tenant-manager)

## Example Usage

```hcl
resource "fusionauth_tenant_manager_configuration" "example" {
  brand_name        = "Acme Corp"
  attribute_form_id = "e95e4455-3c5e-4e53-99c4-1d3a5ea4ca9d"

  application_configurations {
    application_id = "a0c7b3b3-1e4e-4e53-99c4-1d3a5ea4ca9d"
  }

  identity_provider_type_configurations {
    type             = "SAMLv2"
    enabled          = true
    linking_strategy = "LinkByEmail"
  }
}
```

## Argument Reference

* `application_configurations` - (Optional) One or more application configurations for the Tenant Manager.
  * `application_id` - (Required) The unique Id of the Application that the Tenant Manager will use.
* `attribute_form_id` - (Optional) The unique Id of the Form to use for collecting additional user attributes during Tenant Manager registration.
* `brand_name` - (Optional) The brand name for the Tenant Manager.
* `identity_provider_type_configurations` - (Optional) One or more identity provider type configurations allowed in the Tenant Manager. Each block corresponds to one identity provider type.
  * `type` - (Required) The identity provider type. Valid values are: `OpenIDConnect`, `SAMLv2`.
  * `enabled` - (Optional) Whether or not this identity provider type is enabled in the Tenant Manager.
  * `default_attribute_mappings` - (Optional) A map of default attribute mappings for this identity provider type.
  * `linking_strategy` - (Optional) The linking strategy for this identity provider type. Defaults to `LinkByEmail`. Valid values are: `LinkByEmail`, `LinkByEmailForExistingUser`, `LinkByUsername`, `LinkByUsernameForExistingUser`.

## Import

This resource can be imported using the id `tenantmanager_cfg`, e.g.

```shell
terraform import fusionauth_tenant_manager_configuration.example tenantmanager_cfg
```
