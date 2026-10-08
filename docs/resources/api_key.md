# API Key

The FusionAuth APIs are primarily secured using API keys. In order to retrieve, create, update, or delete an API key, an API key set as a keyManager with equal or greater permissions than the API to be created must be used. A "tenant-scoped" API key can retrieve, create, update or delete an API key for the same tenant. This page describes APIs that are used to manage API keys.

[API Key](https://fusionauth.io/docs/v1/tech/apis/api-keys/)

## Example Usage

```hcl
resource "fusionauth_api_key" "example" {
  tenant_id   = "94f751c5-4883-4684-a817-6b106778edec"
  description = "my super secret key"
  key         = "super-secret-key"
  permissions_endpoints {
    endpoint = "/api/application"
    get      = true
    delete   = true
    patch    = true
    post     = true
    put      = true
  }
}
```

### Write-only Key

With Terraform 1.11 or later, the key can be set with `key_wo` so that it is not stored in the plan or state.

```hcl
resource "fusionauth_api_key" "example" {
  description    = "my super secret key"
  key_wo         = var.api_key
  key_wo_version = 1
}
```

~> **Note:** An API key cannot be changed once it is created. Changing `key_wo_version` replaces the API key with one using the current value of `key_wo`. Setting `key_wo_version` on an existing or imported key does not change the key, so `key_wo` must match it.

## Argument Reference

* `description` - (Optional) Description of the key.
* `expiration_instant` - (Optional) The expiration instant of this API key. Using an expired API key for API Authentication will result in a 401 response code.
* `ip_access_control_list_id` - (Optional) The Id of the IP Access Control List limiting access to this API key.
* `key` - (Optional) API key string. When you create an API key the key is defaulted to a secure random value but the API key is simply a string, so you may call it super-secret-key if you’d like. However a long and random value makes a good API key in that it is unique and difficult to guess. Conflicts with `key_wo`.
* `key_id` - (Optional) The Id to use for the new Form. If not specified a secure random UUID will be generated.
* `key_wo` - (Optional) API key string, as a write-only argument that is not stored in state. Conflicts with `key`. Requires Terraform 1.11 or later.
* `key_wo_version` - (Optional) The version of `key_wo`. Required with `key_wo`. Changing this value replaces the API key.
* `name` - (Optional) The name of the API key. Must be unique. If `retrievable` is `false` then this field is required.
* `permissions_endpoints` - (Optional) Endpoint permissions for this key. Each key of the object is an endpoint, with the value being an array of the HTTP methods which can be used against the endpoint. An Empty permissions_endpoints object mean that this is a super key that authorizes this key for all the endpoints.
  * `endpoint` - (Optional)
    * `delete` - (Optional) HTTP DELETE Verb.
    * `get` - (Optional) HTTP GET Verb.
    * `patch` - (Optional) HTTP PATCH Verb
    * `post` - (Optional) HTTP POST Verb
    * `put` - (Optional) HTTP PUT Verb
* `retrievable` - (Optional) Indicates whether this key is retrievable. If this value is false, the key will not be returned in the API response. This value is cannot be updated once the key is created and will be a replacement operation if updated. If this value is set to `false` then the `name` field is required.
* `tenant_id` - (Optional) The unique Id of the Tenant. This value is required if the key is meant to be tenant scoped. Tenant scoped keys can only be used to access users and other tenant scoped objects for the specified tenant. This value is read-only once the key is created.
