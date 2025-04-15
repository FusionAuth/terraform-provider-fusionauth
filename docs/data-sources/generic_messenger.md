# Generic Messenger Data Source

This data source can be used to fetch information about a specific Generic Messenger.

[Messengers API](https://fusionauth.io/docs/apis/messengers/generic)

## Example Usage

```hcl
data "fusionauth_generic_messenger" "example" {
  messenger_id = "4dcbdbb0-385a-4980-ab9f-0c575f6815e0"
}

data "fusionauth_generic_messenger" "example" {
  name = "Generic Messenger"
}
```

## Argument Reference

* `messenger_id` - (Optional) The unique Id of the Generic Messenger to retrieve. If this is not specified, the `name` argument must be specified.
* `name` - (Optional) The case-insensitive string to search for in the Generic Messenger name. If this is not specified, the `messenger_id` argument must be specified.

## Attributes Reference

* `connect_timeout` - The connect timeout for the HTTP connection, in milliseconds. Value must be greater than 0.
* `data` - An object that can hold any information about the Generic Messenger that should be persisted. Represented as a JSON string.
* `debug` - Determines if debug should be enabled to create an event log to assist in debugging integration errors.
* `headers` - An object that can hold HTTPHeader key and value pairs.
* `http_authentication_password` - The basic authentication password to use for requests to the Messenger.
* `http_authentication_username` - The basic authentication username to use for requests to the Messenger.
* `read_timeout` - The read timeout for the HTTP connection, in milliseconds. Value must be greater than 0.
* `ssl_certificate` - An SSL certificate. The certificate is used for client certificate authentication in requests to the Messenger.
* `url` - The fully qualified URL used to send an HTTP request.
