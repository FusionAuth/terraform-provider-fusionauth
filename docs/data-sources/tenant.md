# Tenant Resource

A FusionAuth Tenant is a named object that represents a discrete namespace for Users, Applications and Groups. A user is unique by email address or username within a tenant.

Tenants may be useful to support a multi-tenant application where you wish to use a single instance of FusionAuth but require the ability to have duplicate users across the tenants in your own application. In this scenario a user may exist multiple times with the same email address and different passwords across tenants.

Tenants may also be useful in a test or staging environment to allow multiple users to call APIs and create and modify users without possibility of collision.

[Tenants API](https://fusionauth.io/docs/v1/tech/apis/tenants)

## Example Usage

```hcl
data "fusionauth_tenant" "default"{
    name = "Default"
}
```

## Argument Reference

* `name` - (Required) The name of the Tenant.
