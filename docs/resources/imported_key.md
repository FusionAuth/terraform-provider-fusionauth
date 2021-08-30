# Imported Key Resource

Cryptographic keys are used in signing and verifying JWTs and verifying responses for third party identity providers. It is more likely you will interact with keys using the FusionAuth UI in the Key Master menu. 

[Keys API](https://fusionauth.io/docs/v1/tech/apis/keys)

## Example Usage

```hcl
resource "fusionauth_imported_key" "name" {
  name        = "apple"
  kid         = "8675309"
  private_key = file("./AuthKey_8675309.p8")
}
```

## Argument Reference

* `key_id` - (Optional) The Id to use for the new key. If not specified a secure random UUID will be generated.
* `algorithm` - (Optional) The algorithm used to encrypt the Key. The following values represent algorithms supported by FusionAuth:
    - `ES256` - ECDSA using P-256 curve and SHA-256 hash algorithm
    - `ES384` - ECDSA using P-384 curve and SHA-384 hash algorithm
    - `ES512` - ECDSA using P-521 curve and SHA-512 hash algorithm
    - `RS256` - RSA using SHA-256 hash algorithm
    - `RS384` - RSA using SHA-384 hash algorithm
    - `RS512` - RSA using SHA-512 hash algorithm
    - `HS256` - HMAC using SHA-256 hash algorithm
    - `HS384` - HMAC using SHA-384 hash algorithm
    - `HS512` - HMAC using SHA-512 hash algorithm
* `certificate` - (Optional) The certificate to import. The publicKey will be extracted from the certificate.
* `kid` - (Optional) The Key identifier 'kid'.
* `name` - (Required) The name of the Key.
* `public_key` - (Optional) "The Key public key. Required if importing an RSA or EC key and a certificate is not provided."
* `private_key` - (Optional) The Key private key. Optional if importing an RSA or EC key. If the key is only to be used for token validation, only a public key is necessary and this field may be omitted.
* `secret` - (Optional) The Key secret. This field is required if importing an HMAC key type.
* `type` - (Optional) The Key type. This field is required if importing an HMAC key type, or if importing a public key / private key pair. The possible values are:
    - `EC`
    - `RSA`
    - `HMAC`