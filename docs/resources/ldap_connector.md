# LDAP Connector Resource

A FusionAuth LDAP Connector is a named object that provides configuration for allowing authentication against external LDAP systems.

[LDAP Connector API](https://fusionauth.io/docs/apis/connectors/ldap)

## Example Usage

```hcl
resource "fusionauth_ldap_connector" "example" {
  authentication_url    = "ldap://localhost:389"
  base_structure        = "dc=example,dc=com"
  connect_timeout       = 1000
  identifying_attribute = "uid"
  lambda_configuration {
    reconcile_id = "5bb52c51-bd10-4afc-84d3-e55c62d94987"
  }
  login_id_attribute = "uid"
  name               = "example"
  read_timeout       = 1000
  requested_attributes = [
    "uid",
    "cn",
    "sn",
    "mail"
  ]
  security_method         = "None"
  system_account_dn       = "cn=admin,dc=example,dc=com"
  system_account_password = "password"
}
```

## Argument Reference

* `authentication_url` - (Required) The fully qualified LDAP URL to authenticate.
* `base_structure` - (Required) The top of the LDAP directory hierarchy. Typically this contains the `dc` (domain component) element.
* `connect_timeout` - (Required) The connect timeout for the HTTP connection, in milliseconds. Value must be greater than `0`.
* `identifying_attribute` - (Required) The entry attribute name which is the first component of the distinguished name of entries in the directory.
* `lambda_configuration` - (Required)
  * `reconcile_id` - (Required) The Id of a Lambda. The lambda is executed after the user authenticates with the connector. This lambda can create a user, registrations, and group memberships in FusionAuth based on attributes returned from the connector.
* `login_id_attribute` - (Required) The entity attribute name which stores the identifier that is used for logging the user in.
* `name` - (Required) The unique LDAP Connector name.
* `read_timeout` - (Required) The read timeout for the HTTP connection, in milliseconds. Value must be greater than `0`.
* `requested_attributes` - (Required) The list of attributes to request from the LDAP server. This is a list of strings.
* `security_method` - (Required) The LDAP security method. Possible values are: `None` (Requests will be made without encryption), `LDAPS` (A secure connection will be made to a secure port over using the LDAPS protocol) or `StartTLS` (An un-secured connection will initially be established, followed by secure connection established using the StartTLS extension).
* `system_account_dn` - (Required) The distinguished name of an entry that has read access to the directory.
* `system_account_password` - (Required) The password of an entry that has read access to the directory.

---

* `data` - (Optional) A JSON string that can hold any information about the Connector that should be persisted.
* `debug` - (Optional) Determines if debug should be enabled to create an event log to assist in debugging integration errors.
