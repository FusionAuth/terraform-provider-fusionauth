# Generic Connector Resource

A FusionAuth Generic Connector is a named object that provides configuration for allowing authentication against external systems.

[Generic Connector API](https://fusionauth.io/docs/v1/tech/apis/connectors/generic/)

## Example Usage

```hcl
resource "fusionauth_generic_connector" "example" {
  authentication_url           = "http://mygameserver.local:7001/fusionauth-connector"
  connect_timeout              = 1000
  data                         = { "foo" : "bar" }
  debug                        = false
  headers                      = { "foo" : "bar", "bar" : "baz" }
  http_authentication_password = "supersecret"
  http_authentication_username = "me"
  name                         = "my connector"
  read_timeout                 = 2000
  ssl_certificate_key_id       = "00000000-0000-0000-0000-000000000678"
}
```

## Argument Reference
* `authentication_url` - (Required) The fully qualified URL used to send an HTTP request to authenticate the user.
* `connect_timeout` - (Required) The connect timeout for the HTTP connection, in milliseconds. Value must be greater than 0.
* `data` - (Optional) An object that can hold any information about the Connector that should be persisted.
* `debug` - (Optional) Determines if debug should be enabled to create an event log to assist in debugging integration errors. Defaults to false.
* `headers` - (Optional) An object that can hold HTTPHeader key and value pairs.
* `http_authentication_password` - (Optional) The HTTP basic authentication password that is sent as part of the HTTP request for the events.
* `http_authentication_username` -(Optional) The HTTP basic authentication username that is sent as part of the HTTP request for the events.
* `name` - (Required) The unique Connector name.
* `read_timeout` - (Required) The read timeout in milliseconds used when FusionAuth sends events to the Webhook.
* `ssl_certificate_key_id` - (Optional) The Id of an existing [Key](https://fusionauth.io/docs/v1/tech/apis/keys/). The X509 certificate is used for client certificate authentication in requests to the Connector.
