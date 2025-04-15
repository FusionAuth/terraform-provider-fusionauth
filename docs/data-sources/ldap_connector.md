# LDAP Connector Data Source

This data source can be used to fetch information about a specific LDAP connector.

[Connectors API](https://fusionauth.io/docs/v1/tech/apis/connectors)

## Example Usage

```hcl
data "fusionauth_ldap_connector" "example" {
  id = "8e98aaf8-1e8c-4b11-9f3a-654c94392fb8"
}

data "fusionauth_ldap_connector" "by_name" {
  name = "Corporate LDAP"
}
```

## Argument Reference

* `id` - (Optional) The unique Id of the LDAP Connector to retrieve. This is mutually exclusive with `name`.
* `name` - (Optional) The case-insensitive string to search for in the LDAP Connector name. This is mutually exclusive with `id`.

## Attributes Reference

* `authentication_url` - The fully qualified LDAP URL to authenticate.
* `base_structure` - The top of the LDAP directory hierarchy. Typically this contains the dc (domain component) element.
* `connect_timeout` - The connect timeout for the HTTP connection, in milliseconds. Value must be greater than 0.
* `data` - An object that can hold any information about the Connector that should be persisted. Must be a JSON string.
* `debug` - Determines if debug should be enabled to create an event log to assist in debugging integration errors.
* `identifying_attribute` - The entry attribute name which is the first component of the distinguished name of entries in the directory.
* `lambda_configuration` - Configuration for lambdas (functions) associated with this connector.
  * `reconcile_id` - The Id of an existing Lambda. The lambda is executed after the user authenticates with the connector.
* `login_id_attribute` - The entity attribute name which stores the identifier that is used for logging the user in.
* `read_timeout` - The read timeout for the HTTP connection, in milliseconds. Value must be greater than 0.
* `requested_attributes` - The list of attributes to request from the LDAP server. This is a list of strings.
* `security_method` - The LDAP security method. Possible values are: None, LDAPS, or StartTLS.
* `system_account_dn` - The distinguished name of the system account used to authenticate to the LDAP server.
* `system_account_password` - The password of an entry that has read access to the directory.
