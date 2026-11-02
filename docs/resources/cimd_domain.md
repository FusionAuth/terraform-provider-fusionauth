# CIMD Domain Resource

A CIMD (Client ID Metadata Document) domain represents one trusted `client_id` domain, or the wildcard `*` for open mode, within a tenant. OAuth clients that use a `client_id` URL on this domain are registered dynamically from their metadata document.

Note: An Enterprise plan is required to utilize CIMD. Available since FusionAuth `1.70.0`.

[CIMD Domains API](https://fusionauth.io/docs/apis/cimd-domains)

## Example Usage

```hcl
resource "fusionauth_cimd_domain" "example" {
  tenant_id                = fusionauth_tenant.example.id
  domain                   = "https://mcp.example.com"
  name                     = "Example MCP clients"
  authorized_resource_uris = ["https://api.example.com/mcp"]
  consent_mode             = "RememberDecision"
  grace_period_in_seconds  = 3600
  max_cache_ttl_in_seconds = 86400
  registration_lambda_id   = fusionauth_lambda.cimd_registration.id
  data = jsonencode({
    owner = "platform"
  })
  provided_scope_policy {
    address {
      enabled = false
    }
    email {
      enabled  = true
      required = true
    }
  }
}
```

## Argument Reference

* `domain` - (Required) The origin that CIMD clients register from, given as `scheme://host[:port]`. FusionAuth normalizes this value before storing it and uses it to validate the `client_id` URL presented on an OAuth request. Use the literal value `*` to accept registrations from any origin. A Tenant may contain either a single `*` domain or any number of exact-match domains, but never both. The origin must use `https`, and its host must contain at least one `.`. A path, query string, fragment, or user info component is not permitted, and neither are IP addresses or single-label hosts such as `localhost`. Subdomain wildcards such as `https://*.example.com` are not supported. Changing this value forces a new resource to be created.
* `tenant_id` - (Required) The Id of the Tenant this CIMD domain belongs to. Changing this value forces a new resource to be created.

---

* `allow_localhost_redirect_uris` - (Optional) Whether CIMD clients on this domain may use a `localhost` redirect URI, such as `http://localhost:3334/callback`. The URI must also appear in the client's metadata document `redirect_uris`. When `false`, FusionAuth rejects any `localhost` redirect URI, but still accepts loopback IP addresses such as `http://127.0.0.1/callback`. [RFC 8252 Section 7.3](https://datatracker.ietf.org/doc/html/rfc8252#section-7.3) recommends loopback IP addresses over `localhost`, because the name can resolve to an interface other than loopback. Enable this only for CIMD clients that require it. Defaults to `false`.
* `authorized_resource_uris` - (Optional) The resource URIs that clients registered on this CIMD domain may request using the [RFC 8707](https://datatracker.ietf.org/doc/html/rfc8707) `resource` parameter. Every CIMD client must send exactly one `resource` parameter, and FusionAuth rejects the request with `invalid_target` unless it appears in this list. Clients on a domain with an empty list cannot log in.
* `cimd_domain_id` - (Optional) The Id to use for the new CIMD domain. If not specified a secure random UUID will be generated.
* `consent_mode` - (Optional) Controls whether a user is prompted to consent to the OAuth scopes requested by this domain's clients. The possible values are:
  * `AlwaysPrompt` - Always prompt the user for consent.
  * `RememberDecision` - Remember previous consents; only prompt if the choice expires or if the requested or required scopes have changed. The duration of this persisted choice is controlled by the Tenant's `external_identifier_configuration.remember_oauth_scope_consent_choice_time_to_live_in_seconds` value.

  FusionAuth remembers a decision separately for each client registration and `resource`. Changing the consent mode to `AlwaysPrompt` deletes every remembered decision for the domain's clients.

  A value of `NeverPrompt` will be rejected. A CIMD client is always treated as a third-party client, and once a consent decision is skipped entirely there is no way to re-prompt the client's user. Defaults to `AlwaysPrompt`.
* `data` - (Optional) A JSON string that can hold any information about the CIMD domain that should be persisted.
* `debug` - (Optional) Whether FusionAuth writes the standard OAuth debug Event Log entries, such as the authorization code exchange, for this domain's clients. FusionAuth also writes a Debug entry when it creates a client registration, refreshes one after its cache expires, or adds a new resource to an existing one. The authorization code exchange entry lists the requested resources. FusionAuth always writes an Error entry for a failed client registration, whether or not this is enabled. It always writes an Information entry when it rate limits a registration. Defaults to `false`.
* `grace_period_in_seconds` - (Optional) How long, in seconds, FusionAuth will continue to honor an expired cached metadata document when a re-fetch of that document fails. Must be greater than or equal to `0`. Defaults to `0`.
* `jwt_populate_lambda_id` - (Optional) The unique Id of the lambda used to populate the claims of the access token and Id token issued to this domain's clients. A single lambda is applied to both tokens. The lambda must be a JWT populate lambda.
* `max_cache_ttl_in_seconds` - (Optional) The ceiling, in seconds, applied to the `Cache-Control: max-age` value returned when fetching a client's metadata document. When the document specifies a smaller value, that smaller value is used instead. The `Expires` header is not honored. Must be at least `900`, the non-configurable 15 minute minimum cache duration. Defaults to `86400`.
* `name` - (Optional) The name of the CIMD domain, used to identify it in the FusionAuth admin UI. Leave this blank to default to the value of `domain`.
* `provided_scope_policy` - (Optional) Configures the OAuth scopes provided by FusionAuth for this CIMD domain. An omitted scope keeps the FusionAuth default of enabled and not required.
  * `address` - (Optional)
    * `enabled` - (Optional) Whether the `address` OAuth scope provided by FusionAuth is enabled for this CIMD domain. A disabled scope is treated as an unknown scope when a client requests it. Defaults to `true`.
    * `required` - (Optional) Whether consent to the `address` OAuth scope provided by FusionAuth is required for this CIMD domain when present on the OAuth request. Defaults to `false`.
  * `email` - (Optional)
    * `enabled` - (Optional) Whether the `email` OAuth scope provided by FusionAuth is enabled for this CIMD domain. A disabled scope is treated as an unknown scope when a client requests it. Defaults to `true`.
    * `required` - (Optional) Whether consent to the `email` OAuth scope provided by FusionAuth is required for this CIMD domain when present on the OAuth request. Defaults to `false`.
  * `phone` - (Optional)
    * `enabled` - (Optional) Whether the `phone` OAuth scope provided by FusionAuth is enabled for this CIMD domain. A disabled scope is treated as an unknown scope when a client requests it. Defaults to `true`.
    * `required` - (Optional) Whether consent to the `phone` OAuth scope provided by FusionAuth is required for this CIMD domain when present on the OAuth request. Defaults to `false`.
  * `profile` - (Optional)
    * `enabled` - (Optional) Whether the `profile` OAuth scope provided by FusionAuth is enabled for this CIMD domain. A disabled scope is treated as an unknown scope when a client requests it. Defaults to `true`.
    * `required` - (Optional) Whether consent to the `profile` OAuth scope provided by FusionAuth is required for this CIMD domain when present on the OAuth request. Defaults to `false`.
* `registration_lambda_id` - (Optional) The unique Id of the [CIMD Client Registration lambda](https://fusionauth.io/docs/extend/code/lambdas/cimd-client-registration) used to accept, reject, or modify a client registration for this domain. The lambda must be a CIMD Client Registration lambda (`CimdClientRegistration`). When this is not set, no registration lambda runs and every client from this origin is allowed to register. When `domain` is `*`, a registration lambda is the only way to filter which clients may register.

Custom OAuth scopes for a CIMD domain are managed with the [`fusionauth_cimd_domain_oauth_scope`](cimd_domain_oauth_scope.md) resource.

To enable an identity provider for the clients of a CIMD domain, use the CIMD domain Id as the `application_id` of the identity provider's `application_configuration`, with `create_registration = false`.

## Import

In Terraform v1.5.0 and later, use an `import` block to import CIMD domains using the CIMD domain ID. For example:

```hcl
import {
  to = fusionauth_cimd_domain.name
  id = "cimd_domain_id"
}
```

Using terraform import, import CIMD domains using the CIMD domain ID. For example:

```shell
terraform import fusionauth_cimd_domain.name cimd_domain_id
```
