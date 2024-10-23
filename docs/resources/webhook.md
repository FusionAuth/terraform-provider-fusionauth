# Webhook Resource

A FusionAuth Webhook is intended to consume JSON events emitted by FusionAuth. Creating a Webhook allows you to tell FusionAuth where you would like to receive these JSON events.

[Webhooks API](https://fusionauth.io/docs/v1/tech/apis/webhooks)

## Example Usage

```hcl
resource "fusionauth_webhook" "example" {
  tenant_ids = [
    "00000000-0000-0000-0000-000000000003",
    fusionauth_tenant.example.id
  ]
  connect_timeout = 1000
  description     = "The standard game Webhook"
  events_enabled {
    user_create = true
    user_delete = false
  }
  global                       = false
  headers                      = { "foo" : "bar", "bar" : "baz" }
  http_authentication_password = "password"
  http_authentication_username = "username"
  read_timeout                 = 2000
  ssl_certificate              = <<EOT
  -----BEGIN CERTIFICATE-----\nMIIDUjCCArugAwIBAgIJANZCTNN98L9ZMA0GCSqGSIb3DQEBBQUAMHoxCzAJBgNV\nBAYTAlVTMQswCQYDVQQIEwJDTzEPMA0GA1UEBxMGZGVudmVyMQ8wDQYDVQQKEwZz\nZXRoLXMxCjAIBgNVBAsTAXMxDjAMBgNVBAMTBWludmVyMSAwHgYJKoZIhvcNAQkB\nFhFzamZkZkBsc2tkamZjLmNvbTAeFw0xNDA0MDkyMTA2MDdaFw0xNDA1MDkyMTA2\nMDdaMHoxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJDTzEPMA0GA1UEBxMGZGVudmVy\nMQ8wDQYDVQQKEwZzZXRoLXMxCjAIBgNVBAsTAXMxDjAMBgNVBAMTBWludmVyMSAw\nHgYJKoZIhvcNAQkBFhFzamZkZkBsc2tkamZjLmNvbTCBnzANBgkqhkiG9w0BAQEF\nAAOBjQAwgYkCgYEAxnQBqyuYvjUE4aFQ6vVZU5RqHmy3KiTg2NcxELIlZztUTK3a\nVFbJoBB4ixHXCCYslujthILyBjgT3F+IhSpPAcrlu8O5LVPaPCysh/SNrGNwH4lq\neiW9Z5WAhRO/nG7NZNa0USPHAei6b9Sv9PxuKCY+GJfAIwlO4/bltIH06/kCAwEA\nAaOB3zCB3DAdBgNVHQ4EFgQUU4SqJEFm1zW+CcLxmLlARrqtMN0wgawGA1UdIwSB\npDCBoYAUU4SqJEFm1zW+CcLxmLlARrqtMN2hfqR8MHoxCzAJBgNVBAYTAlVTMQsw\nCQYDVQQIEwJDTzEPMA0GA1UEBxMGZGVudmVyMQ8wDQYDVQQKEwZzZXRoLXMxCjAI\nBgNVBAsTAXMxDjAMBgNVBAMTBWludmVyMSAwHgYJKoZIhvcNAQkBFhFzamZkZkBs\nc2tkamZjLmNvbYIJANZCTNN98L9ZMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEF\nBQADgYEAY/cJsi3w6R4hF4PzAXLhGOg1tzTDYvol3w024WoehJur+qM0AY6UqtoJ\nneCq9af32IKbbOKkoaok+t1+/tylQVF/0FXMTKepxaMbG22vr4TmN3idPUYYbPfW\n5GkF7Hh96BjerrtiUPGuBZL50HoLZ5aR5oZUMAu7TXhOFp+vZp8=\n-----END CERTIFICATE-----
  EOT
  url                          = "http://mygameserver.local:7001/fusionauth-webhook"

  signature_configuration {
    enabled = true
    signing_key_id = fusionauth_key.webhook_key.id
  }

}
```

## Argument Reference
* `tenant_ids` - (Optional) The Ids of the tenants that this Webhook should be associated with. If no Ids are specified and the global field is false, this Webhook will not be used.
* `connect_timeout` - (Required) The connection timeout in milliseconds used when FusionAuth sends events to the Webhook.
* `description` - (Optional) A description of the Webhook. This is used for display purposes only.
* `signature_configuration` - (Optional) Configuration for webhook signing
    - `enabled` - (Optional) Wether or not webhook signing is enabled
    - `signing_key_id` - (Optional) The UUID key used for signing the Webhook
* `events_enabled` - (Optional) A mapping for the events that are enabled for this Webhook.
    - `audit_log_create` - (Optional) When an audit log is created
    - `event_log_create` - (Optional) When an event log is created
    - `jwt_public_key_update` - (Optional) When a JWT RSA Public / Private keypair may have been changed
    - `jwt_refresh` - (Optional) When an access token is refreshed using a refresh token
    - `jwt_refresh_token_revoke` - (Optional) When a JWT Refresh Token is revoked
    - `kickstart_success` - (Optional) When kickstart has successfully completed
    - `user_action` - (Optional) When a user action is triggered
    - `user_bulk_create` - (Optional) When multiple users are created in bulk (i.e. during an import)
    - `user_create` - (Optional) When a user is created
    - `user_create_complete` - (Optional) When a user create transaction has completed
    - `user_deactivate` - (Optional) When a user is deactivated
    - `user_delete` - (Optional) When a user is deleted
    - `user_delete_complete` - (Optional) When a user delete transaction has completed
    - `user_email_update` - (Optional) When a user updates their email address
    - `user_email_verified` - (Optional) When a user verifies their email address
    - `user_identity_provider_link` - (Optional) When a user is linked to an identity provider
    - `user_identity_provider_unlink` - (Optional) When a link to an identity provider is removed
    - `user_login_id_duplicate_create` - (Optional) When a request to create a user with a login Id (email or username) which is already in use has been received
    - `user_login_id_duplicate_update` - (Optional) When a request to update a user and change their login Id (email or username) to one that is already in use has been received
    - `user_login_failed` - (Optional) When a user fails a login request
    - `user_login_new_device` - (Optional) When a user begins a login request with a new device
    - `user_login_success` - (Optional) When a user completes a login request
    - `user_login_suspicious` - (Optional) When a user logs in and is considered to be a potential threat
    - `user_password_breach` - (Optional) When Reactor detects a user is using a potentially breached password (requires an activated license)
    - `user_password_reset_send` - (Optional) When a forgot password email has been sent to a user
    - `user_password_reset_start` - (Optional) When the process to reset a user password has started
    - `user_password_reset_success` - (Optional) When a user has successfully reset their password
    - `user_password_update` - (Optional) When a user has updated their password
    - `user_reactivate` - (Optional) When a user is reactivated
    - `user_registration_create` - (Optional) When a user registration is created
    - `user_registration_create_complete` - (Optional) When a user registration create transaction has completed
    - `user_registration_delete` - (Optional) When a user registration is deleted
    - `user_registration_delete_complete` - (Optional) When a user registration delete transaction has completed
    - `user_registration_update` - (Optional) When a user registration is updated
    - `user_registration_update_complete` - (Optional) When a user registration update transaction has completed
    - `user_registration_verified` - (Optional) When a user completes registration verification
    - `user_two_factor_method_add` - (Optional) When a user has added a two-factor method
    - `user_two_factor_method_remove` - (Optional) When a user has removed a two-factor method
    - `user_update` - (Optional) When a user is updated
    - `user_update_complete` - (Optional) When a user update transaction has completed
* `global` - (Optional) Whether or not this Webhook is used for all events or just for specific Applications.
* `headers` - (Optional) An object that contains headers that are sent as part of the HTTP request for the events.
* `http_authentication_password` - (Optional) The HTTP basic authentication password that is sent as part of the HTTP request for the events.
* `http_authentication_username` -(Optional) The HTTP basic authentication username that is sent as part of the HTTP request for the events.
* `read_timeout` - (Required) The read timeout in milliseconds used when FusionAuth sends events to the Webhook.
* `ssl_certificate` - (Optional) An SSL certificate in PEM format that is used to establish the a SSL (TLS specifically) connection to the Webhook.
* `url` - (Required) The fully qualified URL of the Webhook’s endpoint that will accept the event requests from FusionAuth.