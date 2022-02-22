# Key Resource

Cryptographic keys are used in signing and verifying JWTs and verifying responses for third party identity providers. It is more likely you will interact with keys using the FusionAuth UI in the Key Master menu.

[Keys API](https://fusionauth.io/docs/v1/tech/apis/keys)

## Example Usage

```hcl
resource "fusionauth_key" "admin_id" {
  algorithm = "RS256"
  name      = "Id token signing key generated for application Administrator Login"
  length    = 2048
}
```

## Argument Reference

* `key_id` - (Optional) The Id to use for the new key. If not specified a secure random UUID will be generated.
* `algorithm` - (Required) The algorithm used to encrypt the Key. The following values represent algorithms supported by FusionAuth:
    - `ES256` - ECDSA using P-256 curve and SHA-256 hash algorithm
    - `ES384` - ECDSA using P-384 curve and SHA-384 hash algorithm
    - `ES512` - ECDSA using P-521 curve and SHA-512 hash algorithm
    - `RS256` - RSA using SHA-256 hash algorithm
    - `RS384` - RSA using SHA-384 hash algorithm
    - `RS512` - RSA using SHA-512 hash algorithm
    - `HS256` - HMAC using SHA-256 hash algorithm
    - `HS384` - HMAC using SHA-384 hash algorithm
    - `HS512` - HMAC using SHA-512 hash algorithm
* `name` - (Required) The name of the Key.
* `length` - (Optional)

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `kid` - The id used in the JWT header to identify the key used to generate the signature
