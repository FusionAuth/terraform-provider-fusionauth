# Generic Connector Data Source

This data source can be used to fetch information about a specific Generic Connector.

[Connectors API](https://fusionauth.io/docs/apis/connectors/generic)

## Example Usage

```hcl
data "fusionauth_generic_connector" "example" {
  id = "75a068fd-e94b-451a-9aeb-3ddb9a3b5987"
}

data "fusionauth_generic_connector" "example" {
  name = "My Generic Connector"
}
```

## Argument Reference

* `id` - (Optional) The unique Id of the Generic Connector to retrieve. If this is not specified, the `name` argument must be specified.
* `name` - (Optional) The case-insensitive string to search for in the Generic Connector name. If this is not specified, the `id` argument must be specified.

## Attributes Reference

* `authentication_url` - The fully qualified URL used to send an HTTP request to authenticate the user.
* `connect_timeout` - The connect timeout for the HTTP connection, in milliseconds. Value must be greater than 0.
* `data` - An object that can hold any information about the Generic Connector that should be persisted. Represented as a JSON string.
* `debug` - Determines if debug should be enabled to create an event log to assist in debugging integration errors.
* `headers` - An object that can hold HTTPHeader key and value pairs.
* `http_authentication_password` - The basic authentication password to use for requests to the Connector.
* `http_authentication_username` - The basic authentication username to use for requests to the Connector.
* `read_timeout` - The read timeout for the HTTP connection, in milliseconds. Value must be greater than 0.
* `ssl_certificate_key_id` - The Id of an existing Key. The X509 certificate is used for client certificate authentication in requests to the Connector.
