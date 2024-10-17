# Tenant Resource

A FusionAuth Tenant is a named object that represents a discrete namespace for Users, Applications and Groups. A user is unique by email address or username within a tenant.

Tenants may be useful to support a multi-tenant application where you wish to use a single instance of FusionAuth but require the ability to have duplicate users across the tenants in your own application. In this scenario a user may exist multiple times with the same email address and different passwords across tenants.

Tenants may also be useful in a test or staging environment to allow multiple users to call APIs and create and modify users without possibility of collision.

[Tenants API](https://fusionauth.io/docs/v1/tech/apis/tenants)

## Example Usage

```hcl

resource "fusionauth_tenant" "example" {
  name = "Playtronics Co."
  connector_policy {
    connector_id = "b57b3d0f-f7a4-4831-a838-549717362ea8"
    domains      = ["*"]
    migrate      = false
  }
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
    additional_headers = {
      "HeaderName1" = "HeaderValue1"
      "HeaderName2" = "HeaderValue2"
    }
  }
  event_configuration {
    enabled          = false
    event            = "jwt.public-key.update"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "jwt.refresh-token.revoke"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "jwt.refresh"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "group.create"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "group.create.complete"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "group.delete"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "group.delete.complete"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "group.member.add"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "group.member.add.complete"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "group.member.remove"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "group.member.remove.complete"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "group.member.update"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "group.member.update.complete"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "group.update"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "group.update.complete"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.action"
    transaction_type = "None"
  }
  event_configuration {
    event            = "user.bulk.create"
    enabled          = false
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.create"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.create.complete"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.deactivate"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.delete"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.delete.complete"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.email.update"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.email.verified"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.identity-provider.link"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.identity-provider.unlink"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.loginId.duplicate.create"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.loginId.duplicate.update"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.login.failed"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.login.new-device"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.login.success"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.login.suspicious"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.password.breach"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.password.reset.send"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.password.reset.start"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.password.reset.success"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.password.update"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.reactivate"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.registration.create"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.registration.create.complete"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.registration.delete"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.registration.delete.complete"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.registration.update"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.registration.update.complete"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.registration.verified"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.two-factor.method.add"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.two-factor.method.remove"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.update"
    transaction_type = "None"
  }
  event_configuration {
    enabled          = false
    event            = "user.update.complete"
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
    passwordless_login_time_to_live_in_seconds = 600
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
    two_factor_one_time_code_id_time_to_live_in_seconds = 60
    two_factor_trust_id_time_to_live_in_seconds         = 2592000
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
  form_configuration {
    admin_user_form_id = "e92751a5-25f4-4bca-ad91-66cdf67725d2"
  }
  http_session_max_inactive_interval = 3600
  issuer                             = "https://example.com"
  jwt_configuration {
    access_token_key_id                   = fusionauth_key.accesstoken.id
    id_token_key_id                       = fusionauth_key.idtoken.id
    refresh_token_time_to_live_in_minutes = 43200
    time_to_live_in_seconds               = 3600
  }
  login_configuration {
    require_authentication = true
  }
  maximum_password_age {
    days    = 180
    enabled = false
  }
  minimum_password_age {
    enabled = false
    seconds = 30
  }
  oauth_configuration {
    client_credentials_access_token_populate_lambda_id = fusionauth_lambda.client_jwt_populate.id
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
  rate_limit_configuration {
    failed_login {
      enabled                = true
      limit                  = 5
      time_period_in_seconds = 60
    }
    forgot_password {
      enabled                = false
      limit                  = 5
      time_period_in_seconds = 60
    }
    send_email_verification {
      enabled                = false
      limit                  = 5
      time_period_in_seconds = 60
    }
    send_passwordless {
      enabled                = false
      limit                  = 5
      time_period_in_seconds = 60
    }
    send_registration_verification {
      enabled                = false
      limit                  = 5
      time_period_in_seconds = 60
    }
    send_two_factor {
      enabled                = false
      limit                  = 5
      time_period_in_seconds = 60
    }
  }
  registration_configuration {
    blocked_domains = ["example.com"]
  }
  captcha_configuration {
    enabled         = true
    captcha_method  = "GoogleRecaptchaV3"
    site_key        = "captcha_site_key"
    secret_key      = "captcha_secret_key"
    threshold       = 0.5
  }
  theme_id = fusionauth_theme.example_theme.id
  user_delete_policy {
    unverified_enabled                  = false
    unverified_number_of_days_to_retain = 30
  }
}
```

## Argument Reference
* `source_tenant_id` - (Optional) The optional Id of an existing Tenant to make a copy of. If present, the tenant.id and tenant.name values of the request body will be applied to the new Tenant, all other values will be copied from the source Tenant to the new Tenant.
* `webhook_ids` - (Optional) An array of Webhook Ids. For Webhooks that are not already configured for All Tenants, specifying an Id on this request will indicate the associated Webhook should handle events for this tenant.
* `tenant_id` - (Optional) The Id to use for the new Tenant. If not specified a secure random UUID will be generated.
* `access_control_configuration` - (Optional)
    - `ui_ip_access_control_list_id` - (Optional) The Id of the IP Access Control List limiting access to all applications in this tenant.
* `captcha_configuration` - (Optional)
    - `enabled` - (Optional) Whether captcha configuration is enabled.
    - `captcha_method` - (Optional) The type of captcha method to use. This field is required when tenant.captchaConfiguration.enabled is set to true.
    - `secret_key` - (Optional) The secret key for this captcha method. This field is required when tenant.captchaConfiguration.enabled is set to true.
    - `site_key` - (Optional) The site key for this captcha method. This field is required when tenant.captchaConfiguration.enabled is set to true.
    - `threshold` - (Optional) The numeric threshold which separates a passing score from a failing one. This value only applies if using either the Google v3 or HCaptcha Enterprise method, otherwise this value is ignored.
* `connector_policy` - (Optional) A list of Connector policies. Users will be authenticated against Connectors in order. Each Connector can be included in this list at most once and must exist.
    - `connector_id` - (Optional) The identifier of the Connector to which this policy refers.
    - `domains` - (Optional) A list of email domains to which this connector should apply. A value of ["*"] indicates this connector applies to all users.
    - `migrate` - (Optional) If true, the user’s data will be migrated to FusionAuth at first successful authentication; subsequent authentications will occur against the FusionAuth datastore. If false, the Connector’s source will be treated as authoritative.
* `data` - (Optional) An object that can hold any information about the Tenant that should be persisted.
* `email_configuration` - (Required)
    - `additional_headers` - (Optional) The additional SMTP headers to be added to each outgoing email. Each SMTP header consists of a name and a value.
    - `email_update_email_template_id` - (Optional) The Id of the Email Template that is used when a user is sent a forgot password email.
    - `email_verified_email_template_id` - (Optional) The Id of the Email Template used to verify user emails.
    - `host` - (Required) The host name of the SMTP server that FusionAuth will use.
    - `implicit_email_verification_allowed` - (Optional) When set to true, this allows email to be verified as a result of completing a similar email based workflow such as change password. When seto false, the user must explicitly complete the email verification workflow even if the user has already completed a similar email workflow such as change password.
    - `login_id_in_use_on_create_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when another user attempts to create an account with their login Id.
    - `login_id_in_use_on_update_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when another user attempts to create an account with their login Id.
    - `login_new_device_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when they log in on a new device.
    - `login_suspicious_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when a suspicious login occurs.
    - `password_reset_success_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when they have completed a 'forgot password' workflow and their password habeen reset.
    - `password_update_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when they have completed a 'forgot password' workflow and their password has been rese
    - `default_from_name` - (Optional) The default From Name used in sending emails when a from name is not provided on an individual email template. This is the display name part of the email address ( i.e. Jared Dunn <jared@piedpiper.com>).
    - `default_from_email` - (Optional) The default email address that emails will be sent from when a from address is not provided on an individual email template. This is the address part email address (i.e. Jared Dunn <jared@piedpiper.com>).
    - `forgot_password_email_template_id` - (Optional) The Id of the Email Template that is used when a user is sent a forgot password email.
    - `password` - (Optional) An optional password FusionAuth will use to authenticate with the SMTP server.
    - `passwordless_email_template_id` - (Optional) The Id of the Passwordless Email Template.
    - `port` - (Required) The port of the SMTP server that FusionAuth will use.
    - `properties` - (Optional) Additional Email Configuration in a properties file formatted String.
    - `security` - (Optional) The type of security protocol FusionAuth will use when connecting to the SMTP server.
    - `set_password_email_template_id` - (Optional) The Id of the Email Template that is used when a user had their account created for them and they must set their password manually and they are sent an email to set their password.
    - `username` - (Optional) An optional username FusionAuth will to authenticate with the SMTP server.
    - `verification_email_template_id` - (Optional) The Id of the Email Template that is used to send the verification emails to users. These emails are used to verify that a user’s email address ivalid. If either the verifyEmail or verifyEmailWhenChanged fields are true this field is required.
    - `verification_strategy` - (Optional) The process by which the user will verify their email address. Possible values are `ClickableLink` or `FormField`.
    - `verify_email` - (Optional) Whether the user’s email addresses are verified when the registers with your application.
    - `verify_email_when_changed` - (Optional) Whether the user’s email addresses are verified when the user changes them.
    - `two_factor_method_add_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when a MFA method has been added to their account.
    - `two_factor_method_remove_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when a MFA method has been removed from their account.
    - `unverified` - (Optional)
        - `allow_email_change_when_gated` - (Optional) When this value is set to true, the user is allowed to change their email address when they are gated because they haven’t verified their email address.
        - `behavior` = (Optional) The behavior when detecting breaches at time of user login
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
    - `email_verification_one_time_code_generator` - (Required)
        - `length` - (Required) The length of the secure generator used for generating the email verification one time code.
        - `type` - (Optional) The type of the secure generator used for generating the email verification one time code.
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
    - `registration_verification_one_time_code_generator` - (Required)
        - `length` - (Required) The length of the secure generator used for generating the registration verification one time code.
        - `type` - (Optional) The type of the secure generator used for generating the registration verification one time code.
    - `registration_verification_id_time_to_live_in_seconds` - (Required) The time in seconds until a registration verification Id is no longer valid and cannot be used by the Verify Registration API. Value must be greater than 0.
    - `saml_v2_authn_request_id_ttl_seconds` - (Optional) The time in seconds that a SAML AuthN request will be eligible for use to authenticate with FusionAuth.
    - `setup_password_id_generator` - (Required)
        - `length` - (Required) The length of the secure generator used for generating the change password Id.
        - `type` - (Required) The type of the secure generator used for generating the change password Id.
    - `setup_password_id_time_to_live_in_seconds` - (Required) The time in seconds until a setup password Id is no longer valid and cannot be used by the Change Password API. Value must be greater than 0.
    - `two_factor_id_time_to_live_in_seconds` - (Required) The time in seconds until a two factor Id is no longer valid and cannot be used by the Two Factor Login API. Value must be greater than 0.
    - `trust_token_time_to_live_in_seconds` - (Optional) The number of seconds before the Trust Token is no longer valid to complete a request that requires trust. Value must be greater than 0.
    - `pending_account_link_time_to_live_in_seconds` - (Optional) The number of seconds before the pending account link identifier is no longer valid to complete an account link request. Value must be greater than 0.
    - `two_factor_one_time_code_id_time_to_live_in_seconds` - (Optional) The number of seconds before the Two-Factor One Time Code used to enable or disable a two-factor method is no longer valid. Must be greater than 0.
    - `two_factor_trust_id_time_to_live_in_seconds` - (Required) The time in seconds until an issued Two Factor trust Id is no longer valid and the User will be required to complete Two Factor authentication during the next authentication attempt. Value must be greater than 0.
    - `two_factor_one_time_code_id_generator` - (Required)
        - `length` - (Required) TThe length of the secure generator used for generating the the two factor code Id.
        - `type` - (Optional) The type of the secure generator used for generating the two factor one time code Id.
* `failed_authentication_configuration` - (Optional)
    - `action_duration` - (Required) The duration of the User Action. This value along with the actionDurationUnit will be used to set the duration of the User Action. Value must be greater than 0.
    - `action_duration_unit` - (Optional) The unit of time associated with a duration.
    - `reset_count_in_seconds` - (Optional) The length of time in seconds before the failed authentication count will be reset. Value must be greater than 0.
    - `too_many_attempts` - (Optional) The number of failed attempts considered to be too many. Once this threshold is reached the specified User Action will be applied to the user for the duration specified. Value must be greater than 0.
    - `action_cancel_policy_on_password_reset` - (Optional) Indicates whether you want the user to be able to self-service unlock their account prior to the action duration by completing a password reset workflow.
    - `email_user` - (Optional) Indicates you would like to email the user when the user’s account is locked due to this action being taken. This requires the User Action specified by the tenant.failedAuthenticationConfiguration.userActionId to also be configured for email. If the User Action is not configured to be able to email the user, this configuration will be ignored.
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
* `form_configuration` - (Optional)
    - `admin_user_form_id` - (Optional) The unique Id of the form to use for the Add and Edit User form when used in the FusionAuth admin UI.
* `http_session_max_inactive_interval` - (Optional) Time in seconds until an inactive session will be invalidated. Used when creating a new session in the FusionAuth OAuth frontend.
* `issuer` - (Required) The named issuer used to sign tokens, this is generally your public fully qualified domain.
* `jwt_configuration` - (Required)
    - `access_token_key_id` - (Optional) The unique id of the signing key used to sign the access token. Required prior to `1.30.0`.
    - `id_token_key_id` - (Optional) The unique id of the signing key used to sign the Id token. Required prior to `1.30.0`.
    - `refresh_token_expiration_policy` - (Optional) The refresh token expiration policy.
    - `refresh_token_revocation_policy_on_login_prevented` - (Optional) When enabled, the refresh token will be revoked when a user action, such as locking an account based on a number of failed login attempts, prevents user login.
    - `refresh_token_revocation_policy_on_password_change` - (Optional) When enabled, the refresh token will be revoked when a user changes their password."
    - `refresh_token_sliding_window_maximum_time_to_live_in_minutes` - (Optional) The maximum lifetime of a refresh token when using a refresh token expiration policy of SlidingWindowWithMaximumLifetime. Value must be greater than 0.
    - `refresh_token_time_to_live_in_minutes` - (Required) The length of time in minutes a Refresh Token is valid from the time it was issued. Value must be greater than 0.
    - `refresh_token_usage_policy` - (Optional) The refresh token usage policy.
    - `time_to_live_in_seconds` - (Required) The length of time in seconds this JWT is valid from the time it was issued. Value must be greater than 0.
* `login_configuration`
    - `require_authentication` - (Optional) Indicates whether to require an API key for the Login API when an `applicationId` is not provided. When an `applicationId` is provided to the Login API call, the application configuration will take precedence. In almost all cases, you will want to this to be `true`.
* `logout_url` - (Optional) The logout redirect URL when sending the user’s browser to the /oauth2/logout URI of the FusionAuth Front End. This value is only used when a logout URL is not defined in your Application.
* `maximum_password_age` - (Optional)
    - `days` - (Optional) The password maximum age in days. The number of days after which FusionAuth will require a user to change their password. Required when systemConfiguration.maximumPasswordAge.enabled is set to true.
    - `enabled` - (Optional) Indicates that the maximum password age is enabled and being enforced.
* `minimum_password_age` - (Optional)
    - `seconds` - (Optional) The password minimum age in seconds. When enabled FusionAuth will not allow a password to be changed until it reaches this minimum age. Required when systemConfiguration.minimumPasswordAge.enabled is set to true.
    - `enabled` - (Optional) Indicates that the minimum password age is enabled and being enforced.
* `multi_factor_configuration` - (Optional)
    - `login_policy` - (Optional)  When set to `Enabled` and a user has one or more two-factor methods configured, the user will be required to complete a two-factor challenge during login. When set to `Disabled`, even when a user has configured one or more two-factor methods, the user will not be required to complete a two-factor challenge during login. When the login policy is to `Required`, a two-factor challenge will be required during login. If a user does not have configured two-factor methods, they will not be able to log in.
    - `authenticator` - (Optional)
        * `enabled` - (Optional) When enabled, users may utilize an authenticator application to complete a multi-factor authentication request. This method uses TOTP (Time-Based One-Time Password) as defined in RFC 6238 and often uses an native mobile app such as Google Authenticator.
    - `email` - (Optional)
        * `enabled` - (Optional) When enabled, users may utilize an email address to complete a multi-factor authentication request.
        * `template_id` - (Optional) The Id of the email template that is used when notifying a user to complete a multi-factor authentication request.
    - `sms` - (Optional)
        * `enabled` - (Optional) When enabled, users may utilize a mobile phone number to complete a multi-factor authentication request.
        * `messenger_id` - (Optional) The messenger that is used to deliver a SMS multi-factor authentication request.
        * `template_id` - (Optional) The Id of the SMS template that is used when notifying a user to complete a multi-factor authentication request.
* `name` - (Required) The unique name of the Tenant.
* `oauth_configuration` - (Optional)
    - `client_credentials_access_token_populate_lambda_id` - (Optional) The Id of a lambda that will be called to populate the JWT during a client credentials grant. **Note:** A paid edition of FusionAuth is required to utilize client credentials grant.
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
    - `max_length` - (Optional) The maximum length of a password when a new user is created or a user requests a password change. This value must be greater than 0 and less than or equal to 256. When `passwordEncryptionConfiguration.encryptionScheme` is equal to `bcrypt`, the maximum will be limited to 50.
    - `min_length` - (Optional) The minimum length of a password when a new user is created or a user requests a password change.
    - `remember_previous_passwords` - (Optional)
        - `count` - (Optional) The number of previous passwords to remember. Value must be greater than 0.
        - `enabled` - (Optional) Whether to prevent a user from using any of their previous passwords.
    - `required_mixed_case` - (Optional) Whether to force the user to use at least one uppercase and one lowercase character.
    - `require_non_alpha` - (Optional) Whether to force the user to use at least one non-alphanumeric character.
    - `require_number` - (Optional) Whether to force the user to use at least one number.
    - `validate_on_login` - (Optional) When enabled the user’s password will be validated during login. If the password does not meet the currently configured validation rules the user will be required to change their password.
* `rate_limit_configuration` - (Optional)
    - `failed_login` - (Optional)
      - `enabled` -  (Optional) Whether rate limiting is enabled for failed login.
      - `limit` -  (Optional) The number of times a user can fail to login within the configured `time_period_in_seconds` duration. If a Failed authentication action has been configured then it will take precedence.
      - `time_period_in_seconds` - (Optional) The duration for the number of times a user can fail login before being rate limited.
    - `forgot_password` - (Optional)
      - `enabled` - (Optional) Whether rate limiting is enabled for forgot password.
      - `limit` - (Optional) The number of times a user can request a forgot password email within the configured `time_period_in_seconds` duration.            
      - `time_period_in_seconds` - (Optional) The duration for the number of times a user can request a forgot password email before being rate limited.          
    - `send_email_verification` - (Optional) 
      - `enabled` - (Optional) Whether rate limiting is enabled for send email verification.
      - `limit` - (Optional) The number of times a user can request a verification email within the configured `time_period_in_seconds` duration.                 
      - `time_period_in_seconds` - (Optional) The duration for the number of times a user can request a verification email before being rate limited. 
    - `send_passwordless` - (Optional)
      - `enabled` - (Optional) Whether rate limiting is enabled for send passwordless.
      - `limit` - (Optional) The number of times a user can request a passwordless login email within the configured `time_period_in_seconds` duration.                
      - `time_period_in_seconds` - (Optional) The duration for the number of times a user can request a passwordless login email before being rate limited.
    - `send_registration_verification` - (Optional)
      - `enabled` - (Optional) Whether rate limiting is enabled for send registration verification.
      - `limit` - (Optional) The number of times a user can request a registration verification email within the configured `time_period_in_seconds` duration.                 
      - `time_period_in_seconds` - (Optional) The duration for the number of times a user can request a registration verification email before being rate limited.
    - `send_two_factor` - (Optional)
      - `enabled` - (Optional) Whether rate limiting is enabled for send two factor.
      - `limit` - (Optional) The number of times a user can request a two-factor code by email or SMS within the configured `time_period_in_seconds` duration.   
      - `time_period_in_seconds` - (Optional) The duration for the number of times a user can request a two-factor code by email or SMS before being rate limited.
* `registration_configuration` - (Optional)
    - `blocked_domains` - (Optional) A list of unique domains that are not allowed to register when self service is enabled.
* `theme_id` - (Required) The unique Id of the theme to be used to style the login page and other end user templates.
* `username_configuration` - (Optional)
    - `unique` - (Optional) Indicates that users without a verified email address will be permanently deleted after tenant.userDeletePolicy.unverified.numberOfDaysToRetain days.
        * `enabled` - (Optional) When true, FusionAuth will handle username collisions by generating a random suffix.
        * `number_of_digits` - (Optional) The maximum number of digits to use when building a unique suffix for a username. A number will be randomly selected and will be 1 or more digits up to this configured value in length. For example, if this value is 5, the suffix will be a number between 00001 and 99999, inclusive.
        * `separator` - (Optional) A single character to use as a separator from the requested username and a unique suffix that is added when a duplicate username is detected. This value can be a single non-alphanumeric ASCII character.
        * `strategy` - (Optional) When enabled the user’s password will be validated during login. If the password does not meet the currently configured validation rules the user will be required to change their password.
* `user_delete_policy` - (Optional)
    - `unverified_enabled` - (Optional) Indicates that users without a verified email address will be permanently deleted after tenant.userDeletePolicy.unverified.numberOfDaysToRetain days.
    - `unverified_number_of_days_to_retain` - (Optional)
