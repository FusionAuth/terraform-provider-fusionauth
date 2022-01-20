# IP Access Control List Resource

An IP ACL (Access Control List) is a list of IP ranges that are either Allowed or Blocked. Along with one entry that defines a start IP address of * (wild) that defines the default behavior when an IP address does not match any other range in the list. This means an IP ACL will have a default action of either Allow or Block. The IP address start and end entries for ranges currently only support IPv4.

An IP ACL may be assigned to an API Key, a Tenant or an Application.

When an IP ACL is assigned to an API key, the IP ACL will restrict the usage of the API key based upon the request originating IP address. If a request is made using an API key with an assigned IP ACL and the IP address is found to be blocked, a 401 status code will be returned. The user of this API key will not be able to tell the difference between an invalid API key and an API key that is blocked due to the IP ACL.

When an IP ACL is assigned to a Tenant or Application, it is used to restrict access to the FusionAuth SSO. This means it will be used to restrict access to endpoints that begin with /oauth2/, /account/, /email/, /password/, /registration/ and any other user accessible themed pages. It will not be used to restrict access to the FusionAuth admin UI except when accessed through SSO, or the FusionAuth API.

If two IP ACLs are assigned one to a Tenant and the other to an Application, the Application IP ACL will take precedence.

The IP address used to test against the IP ACL is resolved by using the first value in the X-Forwarded-For HTTP header. If this header is not found, then the IP address reported by the HTTP Servlet request as the remote address will be used. If you are accessing FusionAuth through a proxy it is important that you trust your edge proxy to set the correct value in the X-Forwarded-For HTTP header. Because this header can be set by any HTTP client, it is only secure or trustworthy when managed by a trusted edge proxy. You should not rely upon this feature alone to restrict access to an API key.


[IP Access Control List APIs](https://fusionauth.io/docs/v1/tech/apis/ip-acl/#create-an-ip-acl)

## Example Usage

```hcl
resource "fusionauth_ip_access_control_list" "acl" {
  entries {
    action           = "Block"
    start_ip_address = "*"
    end_ip_address   = "*"
  }
  entries {
    action           = "Allow"
    start_ip_address = "71.205.92.217"
    end_ip_address   = "76.104.251.50"
  }
  name = "Block all except one range"
}
```

## Argument Reference

* `ip_access_control_list_id` - (Optional) The Id to use for the new IP ACL. If not specified a secure random UUID will be generated.
* `entries` - (Required) A list of IP ranges and the action to apply for each. One and only one entry must have a startIPAddress of * to indicate the default action of the IP ACL.
    - `start_ip_address` - (Required) The starting IP (IPv4) for this range.
    - `action` - (Required) The action to take for this IP Range.
    - `end_ip_address` - (Required) The ending IP (IPv4) for this range. The only time this is not required is when start_ip_address is equal to *, in which case this field is ignored. This value must be greater than or equal to the start_ip_address. To define a range of a single IP address, set this field equal to the value for start_ip_address.
* `name` - (Required) The unique name of this IP ACL.