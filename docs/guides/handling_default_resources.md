---
page_title: Handling Default Resources
description: |-
  How to handle resources that are always present in FusionAuth
---

# Handling Default Resources

There are [FusionAuth default configuration elements](https://fusionauth.io/docs/get-started/core-concepts/limitations#default-configuration) present in every FusionAuth instance. If you want to manage changes to these elements via Terraform, you must tell Terraform about them by either importing the resource or setting up a datasource.

### Importing A Resource

To import a resource, you must provide all required attributes. Here's an example for the default tenant:

```hcl
#tag::defaultTenantImport[]
import {
  to = fusionauth_tenant.Default
  id = "Replace-This-With-The-Existing-Default-Tenant-Id"
}

resource "fusionauth_tenant" "Default" {
  lifecycle {
    prevent_destroy = true
  }
  name = "Default"
  issuer = "acme.com"
  theme_id = "00000000-0000-0000-0000-000000000000"
  external_identifier_configuration {
    authorization_grant_id_time_to_live_in_seconds = 30
    change_password_id_generator {
      length = 32
      type   = "randomBytes"
    }
    change_password_id_time_to_live_in_seconds = 600
    device_code_time_to_live_in_seconds        = 300
    device_user_code_id_generator {
      length = 6
      type   = "randomAlphaNumeric"
    }
    email_verification_id_generator {
      length = 32
      type   = "randomBytes"
    }
    email_verification_id_time_to_live_in_seconds      = 86400
    email_verification_one_time_code_generator {
      length = 6
      type   = "randomAlphaNumeric"
    }
    external_authentication_id_time_to_live_in_seconds = 300
    login_intent_time_to_live_in_seconds               = 1800
    one_time_password_time_to_live_in_seconds          = 60
    passwordless_login_generator {
      length = 32
      type   = "randomBytes"
    }
    passwordless_login_time_to_live_in_seconds = 180
    registration_verification_id_generator {
      length = 32
      type   = "randomBytes"
    }
    registration_verification_id_time_to_live_in_seconds = 86400
    registration_verification_one_time_code_generator {
      length = 6
      type   = "randomAlphaNumeric"
    }
    saml_v2_authn_request_id_ttl_seconds = 300
    setup_password_id_generator {
      length = 32
      type   = "randomBytes"
    }
    setup_password_id_time_to_live_in_seconds   = 86400
    two_factor_id_time_to_live_in_seconds       = 300
    two_factor_one_time_code_id_generator {
      length = 6
      type   = "randomDigits"
    }
    two_factor_trust_id_time_to_live_in_seconds = 2592000
  }
  jwt_configuration {
    refresh_token_time_to_live_in_minutes              = 43200
    time_to_live_in_seconds                            = 3600
    refresh_token_revocation_policy_on_login_prevented = true
    refresh_token_revocation_policy_on_password_change = true
    access_token_key_id                                = "00000000-0000-0000-0000-000000000000"
    id_token_key_id                                    = "00000000-0000-0000-0000-000000000000"
  }
  login_configuration {
    require_authentication = true
  }
  email_configuration {
    default_from_email                  = "change-me@example.com"
    default_from_name                   = "FusionAuth"
    host                                = "localhost"
    implicit_email_verification_allowed = true
    port                                = 25
    security                            = "NONE"
    verification_strategy               = "ClickableLink"
    verify_email                        = false
    verify_email_when_changed           = false
  }
}
#end::defaultTenantImport[]
```

You can set some attribute id values to `00000000-0000-0000-0000-000000000000`, then run `terraform plan` to find out the real values. Then update the import statement. This will also display any new tenant default attributes that may have been added over time.

You can do the same for other default resources such as the FusionAuth application or the default theme.

## Data Sources

If you don't need to manage the resource with Terraform, but just want to access its attributes from other places in your Terraform file, you can use a data source.

```hcl
data "fusionauth_tenant" "Default" {
  name = "Default"
}
```

Examples of this include:

* adding applications in the default tenant
* associating a JWT signing key with the FusionAuth application
* setting up an IP ACL to limit access to the FusionAuth application

## Deleting Default Resources

You cannot delete a default resource such as the default tenant or theme. Doing so will cause a Terraform error, since such actions are not allowed by the underlying API.
