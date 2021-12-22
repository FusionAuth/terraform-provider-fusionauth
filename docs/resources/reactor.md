# Reactor Resource

The Reactor is FusionAuth’s license system. Reactor is used to activate features based upon your licensing tier.

[Reactor API](https://fusionauth.io/docs/v1/tech/apis/reactor)

## Example Usage

```hcl
resource "fusionauth_reactor" "reactor" {
  license_id = "xyz"
  license    = "abc"
}
```

## Argument Reference

* `license_id` - (Required) The license Id to activate.
* `license` - (Optional) The Base64 encoded license value. This value is necessary in an air gapped configuration where outbound network access is not available.
