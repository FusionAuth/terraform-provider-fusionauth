# User Group Membership Resource

[User Group Membership API](https://fusionauth.io/docs/apis/groups#request-5)

## Example Usage

```hcl
resource "fusionauth_user_group_membership" "this" {
  group_id      = fusionauth_group.this.id
  user_id       = fusionauth_user.this.id
}
```

## Argument Reference

* `group_id` - (Required) The Id of the Group of this membership.
* `user_id` - (Required) "The Id of the User of this membership.

---

* `data` - (Optional) A JSON string that can hold any information about the User for this membership that should be persisted.
* `membership_id` - (Optional) The Id of the User Group Membership. If not provided, a random UUID will be generated.
