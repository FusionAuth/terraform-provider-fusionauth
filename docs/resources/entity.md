# Entity Resource

Entities are arbitrary objects which can be modeled in FusionAuth. Anything which is not a user but might need
permissions managed by FusionAuth is a possible entity. Examples might include devices, cars, computers, customers,
companies, etc.

FusionAuth’s Entity Management has the following major concepts:

* Entity Types categorize Entities. An Entity Type could be `Device`, `API` or `Company`.
* Permissions are defined on an Entity Type. These are arbitrary strings which can fit the business domain. A Permission
  could be `read`, `write`, or `file-lawsuit`.
* Entities are instances of a single type. An Entity could be a `nest device`, an `Email API` or `Raviga`.
* Entities can have Grants. Grants are relationships between a target Entity and one of two other types: a recipient
  Entity or a User. Grants can have zero or more Permissions associated with them.

You can use the Client Credentials grant to see if an Entity has permission to access another Entity.

[Entity API](https://fusionauth.io/docs/v1/tech/apis/entity-management/entities)

## Example Usage

```hcl
// Create an entity for Raviga...
resource "fusionauth_entity" "raviga" {
  tenant_id      = fusionauth_tenant.default.id
  client_id      = "092dbded-30af-4149-9c61-b578f2c72f59"
  client_secret  = "+fcXet9Iu2kQi61yWD9Tu4ReZ113P6yEAkr32v6WKOQ="
  data           = jsonencode({
    companyType = "Legal"
  })
  entity_type_id = fusionauth_entity_type.company.id
  name           = "Raviga"
}

// With a bound entity type of "Company"...
resource "fusionauth_entity_type" "company" {
  name      = "Company"
  data = jsonencode({
    createdBy = "jared@fusionauth.io"
  })
  jwt_configuration {
    access_token_key_id = "a7516c7c-6234-4021-b0b4-8870c807aeb2"
    enabled = true
    time_to_live_in_seconds = 3600
  }
}

// Which holds a permission to "file-lawsuit"...
resource "fusionauth_entity_permission" "file_lawsuit" {
  entity_type_id = fusionauth_entity_type.company.id
  data        = jsonencode({
    foo = "bar"
  })
  description = "Enables the ability to file lawsuits"
  is_default  = true
  name        = "file-lawsuit"
}
```

## Argument Reference

* `tenant_id` - (Optional) The unique ID of the tenant used to scope this API request.
* `entity_id` - (Optional) The ID to use for the new Entity. If not specified a secure random UUID will be generated.
* `data` - (Optional) An object that can hold any information about the Entity that should be persisted. Please review
  the limits on data field types as you plan for and build your custom data schema. Must be a JSON serialised string.
* `client_id` - (Optional) The OAuth 2.0 client ID. If you leave this blank on create, the value of the Entity ID will
  be used. Must be a UUID.
* `client_secret` - (Optional) The OAuth 2.0 client secret. If you leave this blank on create, a secure secret will be
  generated for you. If you leave this blank during an update, the previous value will be maintained. For both create
  and update you can provide a value and it will be stored.
* `name` - (Required) A descriptive name for the Entity (i.e. "Raviga" or "Email Service").
* `entity_type_id` - (Required) The ID of the Entity Type. Types are consulted for permission checks.

For more information see:
[FusionAuth Entity Management API Overview](https://fusionauth.io/docs/v1/tech/apis/entity-management/)
