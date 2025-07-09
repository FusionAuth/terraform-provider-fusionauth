# Webhook Resource

A FusionAuth Webhook is intended to consume JSON events emitted by FusionAuth. Creating a Webhook allows you to tell FusionAuth where you would like to receive these JSON events.

[Webhooks API](https://fusionauth.io/docs/v1/tech/apis/webhooks)

## Example Usage

```hcl
resource "fusionauth_webhook" "example" {
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

* `connect_timeout` - (Required) The connection timeout in milliseconds used when FusionAuth sends events to the Webhook.
* `read_timeout` - (Required) The read timeout in milliseconds used when FusionAuth sends events to the Webhook.
* `url` - (Required) The fully qualified URL of the Webhook’s endpoint that will accept the event requests from FusionAuth.

---

* `data` - (Optional) A JSON string that can hold any information about the Webhook that should be persisted.
* `description` - (Optional) A description of the Webhook. This is used for display purposes only.
* `events_enabled` - (Optional) A mapping for the events that are enabled for this Webhook.
  * `audit_log_create` - (Optional) An audit log was created
  * `event_log_create` - (Optional) An event log was created
  * `identity_verified` - (Optional) When a user's identity is verified
  * `group_create` - (Optional) A group is being created
  * `group_create_complete` - (Optional) A create group request completed
  * `group_delete` - (Optional) A group is being deleted
  * `group_delete_complete` - (Optional) A group delete request completed
  * `group_member_add` - (Optional) A user is being added to a group
  * `group_member_add_complete` - (Optional) A user add request has completed
  * `group_member_remove` - (Optional) A user is being removed from a group
  * `group_member_remove_complete` - (Optional) A user remove request has completed
  * `group_member_update` - (Optional) A groups membership is being updated
  * `group_member_update_complete` - (Optional) A group member update request has completed
  * `group_update` - (Optional) A group is being updated
  * `group_update_complete` - (Optional) A request to update a group has completed
  * `jwt_public_key_update` - (Optional) A configuration occurred that may affect public keys used to verify a JWT signed by FusionAuth
  * `jwt_refresh` - (Optional) A JWT was refreshed using a refresh token
  * `jwt_refresh_token_revoke` - (Optional) One or more refresh tokens were revoked
  * `kickstart_success` - (Optional) Kickstart completed successfully the system is ready for use
  * `user_action` - (Optional) An action was taken on a user, or an existing event may be changing states if the action is time based
  * `user_bulk_create` - (Optional) One or more users were created using the Bulk create API
  * `user_create` - (Optional) A user is being created
  * `user_create_complete` - (Optional) A create user request completed
  * `user_deactivate` - (Optional) A user is being de-activated, this is synonymous with a soft-delete when using the API
  * `user_delete` - (Optional) A user is being deleted
  * `user_delete_complete` - (Optional) A user delete request has completed
  * `user_email_update` - (Optional) A user updated their email address
  * `user_email_verified` - (Optional) A user has verified their email address
  * `user_identity_provider_link` - (Optional) A link has been established between a user and an identity provider
  * `user_identity_provider_unlink` - (Optional) An existing link has been removed between a user and an identify provider
  * `user_login_id_duplicate_create` - (Optional) A user attempted to register using an email address or username of an existing user
  * `user_login_id_duplicate_update` - (Optional) A user attempted to modify their email address or username to that of an existing user
  * `user_login_failed` - (Optional) A login request has failed
  * `user_login_new_device` - (Optional) A user has logged in from a new device
  * `user_login_success` - (Optional) A login request has succeeded
  * `user_login_suspicious` - (Optional) A suspicious login request has succeeded. This may be due to an impossible travel calculation, or other indicators
  * `user_password_breach` - (Optional) A user's password has been identified as vulnerable due to being found in one or more breached data sets
  * `user_password_reset_send` - (Optional) A user has been sent an email as part of a password reset workflow
  * `user_password_reset_start` - (Optional) A user has started a password reset workflow
  * `user_password_reset_success` - (Optional) A user has completed a password reset workflow
  * `user_password_update` - (Optional) A user has updated their password
  * `user_reactivate` - (Optional) A user has been re-activated
  * `user_registration_create` - (Optional) A user registration is being created
  * `user_registration_create_complete` - (Optional) A user registration has been created
  * `user_registration_delete` - (Optional) A user registration has been deleted
  * `user_registration_delete_complete` - (Optional) A user registration delete request has completed
  * `user_registration_update` - (Optional) A user registration is being updated
  * `user_registration_update_complete` - (Optional) A user registration update request has completed
  * `user_registration_verified` - (Optional) A user registration has been verified
  * `user_two_factor_method_add` - (Optional) A user has added a two-factor method
  * `user_two_factor_method_remove` - (Optional) A user has removed a two-factor method
  * `user_update` - (Optional) A user is being updated
  * `user_update_complete` - (Optional) A user update request has completed
* `global` - (Optional) Whether or not this Webhook is used for all events or just for specific Applications.
* `headers` - (Optional) An object that contains headers that are sent as part of the HTTP request for the events.
* `http_authentication_password` - (Optional) The HTTP basic authentication password that is sent as part of the HTTP request for the events.
* `http_authentication_username` -(Optional) The HTTP basic authentication username that is sent as part of the HTTP request for the events.
* `signature_configuration` - (Optional) Configuration for webhook signing.
  * `enabled` - (Optional) Whether or not webhook events are signed.
  * `signing_key_id` - (Optional) The Id of the key used to sign webhook events. Required when `signature_configuration` is set to true.
* `signature_configuration` - (Optional) Configuration for webhook signing
  * `enabled` - (Optional) Wether or not webhook signing is enabled
  * `signing_key_id` - (Optional) The UUID key used for signing the Webhook
* `ssl_certificate` - (Optional) An SSL certificate in PEM format that is used to establish the a SSL (TLS specifically) connection to the Webhook.
* `ssl_certificate_key_id` - (Optional) The Id of an existing Key. The X.509 certificate is used for client certificate authentication in requests to the Webhook.
* `webhook_id` - (Optional) The Id to use for the new Webhook. If not specified a secure random UUID will be generated.

## Attributes Reference

* `tenant_ids` - The list of tenant ids that this Webhook is associated with.
