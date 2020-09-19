# Tenant Resource

A FusionAuth Tenant is a named object that represents a discrete namespace for Users, Applications and Groups. A user is unique by email address or username within a tenant.

Tenants may be useful to support a multi-tenant application where you wish to use a single instance of FusionAuth but require the ability to have duplicate users across the tenants in your own application. In this scenario a user may exist multiple times with the same email address and different passwords across tenants.

Tenants may also be useful in a test or staging environment to allow multiple users to call APIs and create and modify users without possibility of collision.

[Tenants API](https://fusionauth.io/docs/v1/tech/apis/tenants)

## Example Useage

```hcl

resource "fusionauth_tenant" "example" {
  name = "Playtronics Co."
  email_configuration {
    forgot_password_email_template_id = fusionauth_email.ForgotPassword_Example.id
    host                              = "smtp.sendgrid.net"
    password                          = "password"
    passwordless_email_template_id    = fusionauth_email.PasswordlessLogin_Example.id
    port                              = 587
    security                          = "TLS"
    set_password_email_template_id    = fusionauth_email.SetupPassword_Example.id
    username                          = "username"
    verify_email                      = true
    verify_email_when_changed         = true
  }
  event_configuration {
    event            = "user.delete"
    enabled          = true
    transaction_type = "None"
  }
  event_configuration {
    event            = "user.create"
    enabled          = true
    transaction_type = "None"
  }
  event_configuration {
    event            = "user.update"
    enabled          = true
    transaction_type = "None"
  }
  event_configuration {
    event            = "user.deactivate"
    enabled          = true
    transaction_type = "Any"
  }
  event_configuration {
    event            = "user.bulk.create"
    enabled          = true
    transaction_type = "Any"
  }
  event_configuration {
    event            = "user.reactivate"
    enabled          = true
    transaction_type = "Any"
  }
  event_configuration {
    event            = "jwt.refresh-token.revoke"
    enabled          = true
    transaction_type = "Any"
  }
  event_configuration {
    event            = "jwt.refresh"
    enabled          = true
    transaction_type = "Any"
  }
  event_configuration {
    event            = "jwt.public-key.update"
    enabled          = true
    transaction_type = "Any"
  }
  event_configuration {
    event            = "user.login.success"
    enabled          = true
    transaction_type = "Any"
  }
  event_configuration {
    event            = "user.login.failed"
    enabled          = true
    transaction_type = "Any"
  }
  event_configuration {
    event            = "user.registration.create"
    enabled          = true
    transaction_type = "Any"
  }
  event_configuration {
    event            = "user.registration.update"
    enabled          = true
    transaction_type = "Any"
  }
  event_configuration {
    event            = "user.registration.delete"
    enabled          = true
    transaction_type = "Any"
  }
  event_configuration {
    event            = "user.registration.verified"
    enabled          = true
    transaction_type = "Any"
  }
  event_configuration {
    event            = "user.email.verified"
    enabled          = true
    transaction_type = "Any"
  }
  event_configuration {
    event            = "user.password.breach"
    enabled          = false
    transaction_type = "None"
  }
  external_identifier_configuration {
    authorization_grant_id_time_to_live_in_seconds = 30
    change_password_id_generator {
      length = 32
      type   = "randomBytes"
    }
    change_password_id_time_to_live_in_seconds = 600
    device_code_time_to_live_in_seconds        = 1800
    device_user_code_id_generator {
      length = 6
      type   = "randomAlphaNumeric"
    }
    email_verification_id_generator {
      length = 32
      type   = "randomBytes"
    }
    email_verification_id_time_to_live_in_seconds      = 86400
    external_authentication_id_time_to_live_in_seconds = 300
    one_time_password_time_to_live_in_seconds          = 60
    passwordless_login_generator {
      length = 32
      type   = "randomBytes"
    }
    passwordless_login_time_to_live_in_seconds = 600
    registration_verification_id_generator {
      length = 32
      type   = "randomBytes"
    }
    registration_verification_id_time_to_live_in_seconds = 86400
    setup_password_id_generator {
      length = 32
      type   = "randomBytes"
    }
    setup_password_id_time_to_live_in_seconds   = 86400
    two_factor_id_time_to_live_in_seconds       = 300
    two_factor_trust_id_time_to_live_in_seconds = 2592000
  }
  failed_authentication_configuration {
    action_duration        = 3
    action_duration_unit   = "MINUTES"
    reset_count_in_seconds = 60
    too_many_attempts      = 5
  }
  family_configuration {
    allow_child_registrations     = true
    delete_orphaned_accounts      = false
    delete_orphaned_accounts_days = 30
    enabled                       = true
    maximum_child_age             = 12
    minimum_owner_age             = 21
    parent_email_required         = false
  }
  http_session_max_inactive_interval = 3600
  issuer                             = "https://example.com"
  jwt_configuration {
    access_token_key_id                   = fusionauth_key.accesstoken.id
    id_token_key_id                       = fusionauth_key.idtoken.id
    refresh_token_time_to_live_in_minutes = 43200
    time_to_live_in_seconds               = 3600
  }
  maximum_password_age {
    days    = 180
    enabled = false
  }
  minimum_password_age {
    enabled = false
    seconds = 30
  }
  password_encryption_configuration {
    encryption_scheme                 = "salted-pbkdf2-hmac-sha256"
    encryption_scheme_factor          = 24000
    modify_encryption_scheme_on_login = false
  }
  password_validation_rules {
    max_length = 256
    min_length = 7
    remember_previous_passwords {
      count   = 1
      enabled = false
    }
    required_mixed_case = false
    require_non_alpha   = false
    require_number      = false
    validate_on_login   = false
  }
  theme_id = fusionauth_theme.example_theme.id
  user_delete_policy {
    unverified_enabled                  = false
    unverified_number_of_days_to_retain = 30
  }
}
```

## Argument Reference
* `source_tentant_id` - (Optional) The optional Id of an existing Tenant to make a copy of. If present, the tenant.id and tenant.name values of the request body will be applied to the new Tenant, all other values will be copied from the source Tenant to the new Tenant.
* `tenant_id` - (Optional) The Id to use for the new Tenant. If not specified a secure random UUID will be generated.
* `data` - (Optional) An object that can hold any information about the Tenant that should be persisted.
* `email_configuration` - (Required)
    - `default_from_name` - (Optional) The default From Name used in sending emails when a from name is not provided on an individual email template. This is the display name part of the email address ( i.e. Jared Dunn <jared@piedpiper.com>).
    - `default_from_email` - (Optional) The default email address that emails will be sent from when a from address is not provided on an individual email template. This is the address part email address (i.e. Jared Dunn <jared@piedpiper.com>).
    - `forgot_password_email_template_id` - (Optional) The Id of the Email Template that is used when a user is sent a forgot password email.
    - `host` - (Required) The host name of the SMTP server that FusionAuth will use.
    - `password` - (Optional) An optional password FusionAuth will use to authenticate with the SMTP server.
    - `passwordless_email_template_id` - (Optional) The Id of the Passwordless Email Template.
    - `port` - (Required) The port of the SMTP server that FusionAuth will use.
    - `properties` - (Optional) Additional Email Configuration in a properties file formatted String.
    - `security` - (Optional) The type of security protocol FusionAuth will use when connecting to the SMTP server.
    - `set_password_email_template_id` - (Optional) The Id of the Email Template that is used when a user had their account created for them and they must set their password manually and they are sent an email to set their password.
    - `username` - (Optional) An optional username FusionAuth will to authenticate with the SMTP server.
    - `verification_email_template_id` - () The If of the Email Template that is used to send the verification emails to users. These emails are used to verify that a user’s email address is valid. If either the verifyEmail or verifyEmailWhenChanged fields are true this field is required.
    - `verify_email` - (Optional) Whether the user’s email addresses are verified when the registers with your application.
    - `verify_email_when_changed` - (Optional) Whether the user’s email addresses are verified when the user changes them.
* `event_configuration` - (Optional)
    - `event` - (Optional) The event type
    - `enabled` - (Optional) Whether or not FusionAuth should send these types of events to any configured Webhooks.
    - `transaction_type` - (Optional) The transaction type that FusionAuth uses when sending these types of events to any configured Webhooks.
* `external_identifier_configuration` - (Required)
    - `authorization_grant_id_time_to_live_in_seconds` - (Required) The time in seconds until a OAuth authorization code in no longer valid to be exchanged for an access token. This is essentially the time allowed between the start of an Authorization request during the Authorization code grant and when you request an access token using this authorization code on the Token endpoint.
    - `change_password_id_generator` - (Required)
        - `length` - (Required) The length of the secure generator used for generating the change password Id.
        - `type` - (Required) The type of the secure generator used for generating the change password Id.
    - `change_password_id_time_to_live_in_seconds` - (Required) The time in seconds until a change password Id is no longer valid and cannot be used by the Change Password API. Value must be greater than 0.
    - `device_code_time_to_live_in_seconds` - (Required) The time in seconds until a device code Id is no longer valid and cannot be used by the Token API. Value must be greater than 0.
    - `device_user_code_id_generator` - (Required)
        - `length` - (Required) The length of the secure generator used for generating the change password Id.
        - `type` - (Required) The type of the secure generator used for generating the change password Id.
    - `email_verification_id_generator` - (Required)
        - `length` - (Required) The length of the secure generator used for generating the change password Id.
        - `type` - (Required) The type of the secure generator used for generating the change password Id.
    - `email_verification_id_time_to_live_in_seconds` - (Required) The time in seconds until a email verification Id is no longer valid and cannot be used by the Verify Email API. Value must be greater than 0.
    - `external_authentication_id_time_to_live_in_seconds` - (Required) The time in seconds until an external authentication Id is no longer valid and cannot be used by the Token API. Value must be greater than 0.
    - `one_time_password_time_to_live_in_seconds` - (Required) The time in seconds until a One Time Password is no longer valid and cannot be used by the Login API. Value must be greater than 0.
    - `passwordless_login_generator` - (Required)
        - `length` - (Required) The length of the secure generator used for generating the change password Id.
        - `type` - (Required) The type of the secure generator used for generating the change password Id.
    - `passwordless_login_time_to_live_in_seconds` - (Required) The time in seconds until a passwordless code is no longer valid and cannot be used by the Passwordless API. Value must be greater than 0.
    - `registration_verification_id_generator` - (Required)
        - `length` - (Required) The length of the secure generator used for generating the change password Id.
        - `type` - (Required) The type of the secure generator used for generating the change password Id.
    - `registration_verification_id_time_to_live_in_seconds` - (Required) The time in seconds until a registration verification Id is no longer valid and cannot be used by the Verify Registration API. Value must be greater than 0.
    - `setup_password_id_generator` - (Required)
        - `length` - (Required) The length of the secure generator used for generating the change password Id.
        - `type` - (Required) The type of the secure generator used for generating the change password Id.
    - `setup_password_id_time_to_live_in_seconds` - (Required) The time in seconds until a setup password Id is no longer valid and cannot be used by the Change Password API. Value must be greater than 0.
    - `two_factor_id_time_to_live_in_seconds` - (Required) The time in seconds until a two factor Id is no longer valid and cannot be used by the Two Factor Login API. Value must be greater than 0.
    - `two_factor_trust_id_time_to_live_in_seconds` - (Require) The time in seconds until an issued Two Factor trust Id is no longer valid and the User will be required to complete Two Factor authentication during the next authentication attempt. Value must be greater than 0.
* `failed_authentication_configuration` - (Optional)
    - `action_duration` - (Required) The duration of the User Action. This value along with the actionDurationUnit will be used to set the duration of the User Action. Value must be greater than 0. 
    - `action_duration_unit` - (Optional) The unit of time associated with a duration.
    - `reset_count_in_seconds` - (Optional) The length of time in seconds before the failed authentication count will be reset. Value must be greater than 0.
    - `too_many_attempts` - (Optional) The number of failed attempts considered to be too many. Once this threshold is reached the specified User Action will be applied to the user for the duration specified. Value must be greater than 0.
    - `user_action_id` - (Optional) The Id of the User Action that is applied when the threshold is reached for too many failed authentication attempts.
* `family_configuration` - (Optional)
    - `allow_child_registrations` - (Optional) Whether to allow child registrations.
    - `confirm_child_email_template_id` - (Optional) The unique Id of the email template to use when confirming a child.
    - `delete_orphaned_accounts` - (Optional) Indicates that child users without parental verification will be permanently deleted after tenant.familyConfiguration.deleteOrphanedAccountsDays days.
    - `delete_orphaned_accounts_days` - (Optional) The number of days from creation child users will be retained before being deleted for not completing parental verification. Value must be greater than 0.
    - `enabled` - (Optional) Whether family configuration is enabled.
    - `family_request_email_template_id` - (Optional) The unique Id of the email template to use when a family request is made.
    - `maximum_child_age` - (Optional) The maximum age of a child. Value must be greater than 0.
    - `minimum_owner_age` - (Optional) The minimum age to be an owner. Value must be greater than 0.
    - `parent_email_required` - (Optional) Whether a parent email is required.
    - `parent_registration_email_template_id` - (Optional) The unique Id of the email template to use for parent registration.
* `http_session_max_inactive_interval` - (Optional) Time in seconds until an inactive session will be invalidated. Used when creating a new session in the FusionAuth OAuth frontend.
* `issuer` - (Required) The named issuer used to sign tokens, this is generally your public fully qualified domain.
* `jwt_configuration` - (Required)
    - `access_token_key_id` - (Required) The unique id of the signing key used to sign the access token.
    - `id_token_key_id` - (Required) The unique id of the signing key used to sign the Id token.
    - `refresh_token_expiration_policy` - (Optional) The refresh token expiration policy.
    - `refresh_token_revocation_policy_on_login_prevented` - (Optional) When enabled, the refresh token will be revoked when a user action, such as locking an account based on a number of failed login attempts, prevents user login.
    - `refresh_token_revocation_policy_on_password_change` - (Optional) When enabled, the refresh token will be revoked when a user changes their password."
    - `refresh_token_time_to_live_in_minutes` - (Required) The length of time in minutes a Refresh Token is valid from the time it was issued. Value must be greater than 0.
    - `refresh_token_usage_policy` - (Optional) The refresh token usage policy.
    - `time_to_live_in_seconds` - (Required) The length of time in seconds this JWT is valid from the time it was issued. Value must be greater than 0.
* `logout_url` - (Optional) The logout redirect URL when sending the user’s browser to the /oauth2/logout URI of the FusionAuth Front End. This value is only used when a logout URL is not defined in your Application.
* `maximum_password_age` - (Optional)
    - `days` - (Optional) The password maximum age in days. The number of days after which FusionAuth will require a user to change their password. Required when systemConfiguration.maximumPasswordAge.enabled is set to true.
    - `enabled` - (Optional) Indicates that the maximum password age is enabled and being enforced.
* `minimum_password_age` - (Optional)
    - `seconds` - (Optional) The password minimum age in seconds. When enabled FusionAuth will not allow a password to be changed until it reaches this minimum age. Required when systemConfiguration.minimumPasswordAge.enabled is set to true.
    - `enabled` - (Optional) Indicates that the minimum password age is enabled and being enforced.
* `name` - (Required) The unique name of the Tenant.
* `password_encryption_configuration` - (Optional)
    - `encryption_scheme` - (Optional) The default method for encrypting the User’s password.
    - `encryption_scheme_factor` - (Optional) The factor used by the password encryption scheme. If not provided, the PasswordEncryptor provides a default value. Generally this will be used as an iteration count to generate the hash. The actual use of this value is up to the PasswordEncryptor implementation.
    - `modify_encryption_scheme_on_login` - (Optional) When enabled a user’s hash configuration will be modified to match these configured settings. This can be useful to increase a password hash strength over time or upgrade imported users to a more secure encryption scheme after an initial import.
* `password_validation_rules` - (Optional)
    - `breach_detection` - (Optional)
        - `enabled` - (Optional) Whether to enable Reactor breach detection. Requires an activated license.
        - `match_mode` - (Optional) The level of severity where Reactor will consider a breach.
        - `notify_user_email_template_id` - (Optional) The Id of the email template to use when notifying user of breached password. Required if tenant.passwordValidationRules.breachDetection.onLogin is set to NotifyUser.
        - `on_login` - (Optional) The behavior when detecting breaches at time of user login
    - `max_length` - (Optional) The maximum length of a password when a new user is created or a user requests a password change.
    - `min_length` - (Optional) The minimum length of a password when a new user is created or a user requests a password change.
    - `remember_previous_passwords` - (Optional)
        - `count` - (Optional) The number of previous passwords to remember. Value must be greater than 0.
        - `enabled` - (Optional) Whether to prevent a user from using any of their previous passwords.
    - `required_mixed_case` - (Optional) Whether to force the user to use at least one uppercase and one lowercase character.
    - `require_non_alpha` - (Optional) Whether to force the user to use at least one non-alphanumeric character.
    - `require_number` - (Optional) Whether to force the user to use at least one number.
    - `validate_on_login` - (Optional) When enabled the user’s password will be validated during login. If the password does not meet the currently configured validation rules the user will be required to change their password.
* `theme_id` - (Required) The unique Id of the theme to be used to style the login page and other end user templates.
* `user_delete_policy` - (Optional)
    - `unverified_enabled` - (Optional) Indicates that users without a verified email address will be permanently deleted after tenant.userDeletePolicy.unverified.numberOfDaysToRetain days.
    - `unverified_number_of_days_to_retain` - (Optional)