# Generic Messenger Resource

A FusionAuth Generic Messenger is a named object that provides configuration for sending messages through external HTTP endpoints.

[Generic Messenger API](https://fusionauth.io/docs/v1/tech/apis/messengers/generic/)

## Example Usage

```hcl
resource "fusionauth_generic_messenger" "example" {
  name            = "Generic Messenger"
  url             = "https://www.example.com/webhook"
  connect_timeout = 1000
  read_timeout    = 1000

  data = jsonencode({
    "foo" : "bar"
  })

  debug = false

  headers = {
    "Content-Type" = "application/json"
  }

  http_authentication_username = "user-login"
  http_authentication_password = "password"

  ssl_certificate = <<-EOT
  -----BEGIN CERTIFICATE-----
  MIIDazCCAlOgAwIBAgIUJlq+zz9CO2EIuBWULFPYoBWKDFwwDQYJKoZIhvcNAQEL
  ...
  -----END CERTIFICATE-----
  EOT
}
```

## Argument Reference

* `name` - (Required) The unique Messenger name.
* `url` - (Required) The fully qualified URL used to send an HTTP request.
* `connect_timeout` - (Required) The connect timeout for the HTTP connection, in milliseconds. Value must be greater than 0.
* `read_timeout` - (Required) The read timeout for the HTTP connection, in milliseconds. Value must be greater than 0.

---

* `data` - (Optional) A JSON string that can hold any information about the Generic Messenger that should be persisted.
* `debug` - (Optional) Determines if debug should be enabled to create an event log to assist in debugging integration errors. Defaults to false.
* `headers` - (Optional) An object that can hold HTTPHeader key and value pairs.
* `http_authentication_username` - (Optional) The basic authentication username to use for requests to the Messenger.
* `http_authentication_password` - (Optional) The basic authentication password to use for requests to the Messenger.
* `messenger_id` - (Optional) The Id to use for the new Messenger. If not specified a secure random UUID will be generated.
* `ssl_certificate` - (Optional) An SSL certificate. The certificate is used for client certificate authentication in requests to the Messenger.
