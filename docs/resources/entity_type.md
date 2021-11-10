# Entity Type Resource

Entity Types categorize Entities. For example, an Entity Type could be `Device`, `API` or `Company`.

[Entity Type API](https://fusionauth.io/docs/v1/tech/apis/entity-management/entity-types/#create-an-entity-type)

## Example Usage

```hcl
resource "fusionauth_entity_type" "company" {
  name = "Company"
  data = jsonencode({
    createdBy = "jared@fusionauth.io"
  })
  jwt_configuration {
    access_token_key_id     = "a7516c7c-6234-4021-b0b4-8870c807aeb2"
    enabled                 = true
    time_to_live_in_seconds = 3600
  }
}
```

## Argument Reference

* `entity_type_id` - (Optional) The ID to use for the new Entity Type. If not specified a secure random UUID will be
  generated.
* `data` - (Optional) An object that can hold any information about the Entity Type that should be persisted. Must be a
  JSON string.
* `jwt_configuration` - (Optional) A block to configure JSON Web Token (JWT) options.
    - `enabled` - (Optional) Indicates if this application is using the JWT configuration defined here or the global JWT
      configuration defined by the Tenant. If this is false the signing algorithm configured in the Tenant will be used.
      If true the signing algorithm defined in this application will be used.
    - `access_token_key_id` - (Required) The unique ID of the signing key used to sign the access token. Required when
      enabled is set to true.
    - `time_to_live_in_seconds` - (Required) The length of time in seconds the JWT will live before it is expired and no
      longer valid. Required when enabled is set to true.
* `name` - (Required) A descriptive name for the entity type (i.e. `Customer` or `Email_Service`).
