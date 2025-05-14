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

* `name` - (Required) The unique name of the Tenant.

---

* `access_control_configuration` - (Optional)
  * `ui_ip_access_control_list_id` - (Optional) The Id of the IP Access Control List limiting access to all applications in this tenant.
* `captcha_configuration` - (Optional)
  * `captcha_method` - (Optional) The type of captcha method to use. This field is required when tenant.captchaConfiguration.enabled is set to true.
  * `enabled` - (Optional) Whether captcha configuration is enabled.
  * `secret_key` - (Optional) The secret key for this captcha method. This field is required when tenant.captchaConfiguration.enabled is set to true.
  * `site_key` - (Optional) The site key for this captcha method. This field is required when tenant.captchaConfiguration.enabled is set to true.
  * `threshold` - (Optional) The numeric threshold which separates a passing score from a failing one. This value only applies if using either the Google v3 or HCaptcha Enterprise method, otherwise this value is ignored.
* `connector_policy` - (Optional) A list of Connector policies. Users will be authenticated against Connectors in order. Each Connector can be included in this list at most once and must exist.
  * `connector_id` - (Optional) The identifier of the Connector to which this policy refers.
  * `domains` - (Optional) A list of email domains to which this connector should apply. A value of ["*"] indicates this connector applies to all users.
  * `migrate` - (Optional) If true, the user’s data will be migrated to FusionAuth at first successful authentication; subsequent authentications will occur against the FusionAuth datastore. If false, the Connector’s source will be treated as authoritative.
* `data` - (Optional) A JSON string that can hold any information about the Tenant that should be persisted.
* `email_configuration` - (Optional) The email configuration for the tenant.
  * `additional_headers` - (Optional) The additional SMTP headers to be added to each outgoing email. Each SMTP header consists of a name and a value.
  * `debug` - (Optional) Determines if debug should be enabled to create an event log to assist in debugging SMTP errors.
  * `default_from_email` - (Optional) The default email address that emails will be sent from when a from address is not provided on an individual email template. This is the address part email address (i.e. Jared Dunn <jared@piedpiper.com>).
  * `default_from_name` - (Optional) The default From Name used in sending emails when a from name is not provided on an individual email template. This is the display name part of the email address ( i.e. Jared Dunn <jared@piedpiper.com>).
  * `email_update_email_template_id` - (Optional) The Id of the Email Template that is used when a user is sent a forgot password email.
  * `email_verified_email_template_id` - (Optional) The Id of the Email Template used to verify user emails.
  * `forgot_password_email_template_id` - (Optional) The Id of the Email Template that is used when a user is sent a forgot password email.
  * `host` - (Optional) The host name of the SMTP server that FusionAuth will use.
  * `implicit_email_verification_allowed` - (Optional) When set to true, this allows email to be verified as a result of completing a similar email based workflow such as change password. When seto false, the user must explicitly complete the email verification workflow even if the user has already completed a similar email workflow such as change password.
  * `login_id_in_use_on_create_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when another user attempts to create an account with their login Id.
  * `login_id_in_use_on_update_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when another user attempts to create an account with their login Id.
  * `login_new_device_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when they log in on a new device.
  * `login_suspicious_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when a suspicious login occurs.
  * `password` - (Optional) An optional password FusionAuth will use to authenticate with the SMTP server.
  * `passwordless_email_template_id` - (Optional) The Id of the Passwordless Email Template.
  * `password_reset_success_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when they have completed a 'forgot password' workflow and their password habeen reset.
  * `password_update_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when they have completed a 'forgot password' workflow and their password has been rese
  * `port` - (Optional) The port of the SMTP server that FusionAuth will use.
  * `properties` - (Optional) Additional Email Configuration in a properties file formatted String.
  * `security` - (Optional) The type of security protocol FusionAuth will use when connecting to the SMTP server.
  * `set_password_email_template_id` - (Optional) The Id of the Email Template that is used when a user had their account created for them and they must set their password manually and they are sent an email to set their password.
  * `two_factor_method_add_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when a MFA method has been added to their account.
  * `two_factor_method_remove_email_template_id` - (Optional) The Id of the Email Template used to send emails to users when a MFA method has been removed from their account.
  * `username` - (Optional) An optional username FusionAuth will to authenticate with the SMTP server.
  * `verification_email_template_id` - (Optional) The Id of the Email Template that is used to send the verification emails to users. These emails are used to verify that a user’s email address ivalid. If either the verifyEmail or verifyEmailWhenChanged fields are true this field is required.
  * `verification_strategy` - (Optional) The process by which the user will verify their email address. Possible values are `ClickableLink` or `FormField`.
  * `verify_email` - (Optional) Whether the user’s email addresses are verified when the registers with your application.
  * `verify_email_when_changed` - (Optional) Whether the user’s email addresses are verified when the user changes them.
  * `unverified` - (Optional)
    * `allow_email_change_when_gated` - (Optional) When this value is set to true, the user is allowed to change their email address when they are gated because they haven’t verified their email address.
    * `behavior` = (Optional) The behavior when detecting breaches at time of user login.
* `event_configuration` - (Optional)
  * `enabled` - (Optional) Whether or not FusionAuth should send these types of events to any configured Webhooks.
  * `event` - (Optional) The event type
  * `transaction_type` - (Optional) The transaction type that FusionAuth uses when sending these types of events to any configured Webhooks.
* `external_identifier_configuration` - (Optional)
  * `authorization_grant_id_time_to_live_in_seconds` - (Optional) The time in seconds until a OAuth authorization code in no longer valid to be exchanged for an access token. This is essentially the time allowed between the start of an Authorization request during the Authorization code grant and when you request an access token using this authorization code on the Token endpoint. Defaults to 30.
  * `change_password_id_generator` - (Optional)
    * `length` - (Optional) The length of the secure generator used for generating the change password Id. Defaults to 32.
    * `type` - (Optional) The type of the secure generator used for generating the change password Id. Defaults to randomBytes.
  * `change_password_id_time_to_live_in_seconds` - (Optional) The time in seconds until a change password Id is no longer valid and cannot be used by the Change Password API. Value must be greater than 0. Defaults to 600.
  * `device_code_time_to_live_in_seconds` - (Optional) The time in seconds until a device code Id is no longer valid and cannot be used by the Token API. Value must be greater than 0. Defaults to 300.
  * `device_user_code_id_generator` - (Optional)
    * `length` - (Optional) The length of the secure generator used for generating the change password Id. Defaults to 6.
    * `type` - (Optional) The type of the secure generator used for generating the change password Id. Defaults to randomAlphaNumeric.
  * `email_verification_id_generator` - (Optional)
    * `length` - (Optional) The length of the secure generator used for generating the change password Id. Defaults to 32.
    * `type` - (Optional) The type of the secure generator used for generating the change password Id. Defaults to randomBytes.
  * `email_verification_id_time_to_live_in_seconds` - (Optional) The time in seconds until a email verification Id is no longer valid and cannot be used by the Verify Email API. Value must be greater than 0.
  * `email_verification_one_time_code_generator` - (Optional)
    * `length` - (Optional) The length of the secure generator used for generating the email verification one time code. Defaults to 6.
    * `type` - (Optional) The type of the secure generator used for generating the email verification one time code. Defaults to randomAlphaNumeric.
  * `external_authentication_id_time_to_live_in_seconds` - (Optional) The time in seconds until an external authentication Id is no longer valid and cannot be used by the Token API. Value must be greater than 0. Defaults to 300.
  * `login_intent_time_to_live_in_seconds` - (Optional ) The time in seconds until a Login Timeout identifier is no longer valid to complete post-authentication steps in the OAuth workflow. Must be greater than 0. Defaults to 1800.
  * `one_time_password_time_to_live_in_seconds` - (Optional) The time in seconds until a One Time Password is no longer valid and cannot be used by the Login API. Value must be greater than 0. Defaults to 60.
  * `passwordless_login_generator` - (Optional)
    * `length` - (Optional) The length of the secure generator used for generating the change password Id. Defaults to 32
    * `type` - (Optional) The type of the secure generator used for generating the change password Id. Defaults to randomBytes.
  * `passwordless_login_one_time_code_generator` - (Optional)
    * `length` - (Optional) The length of the secure generator used for generating the passwordless one-time code login. Defaults to 32
    * `type` - (Optional) The type of the secure generator used for generating the passwordless one-time code login. Defaults to randomBytes.
  * `passwordless_login_time_to_live_in_seconds` - (Optional) The time in seconds until a passwordless code is no longer valid and cannot be used by the Passwordless API. Value must be greater than 0. Defaults to 180.
  * `pending_account_link_time_to_live_in_seconds` - (Optional) The number of seconds before the pending account link identifier is no longer valid to complete an account link request. Value must be greater than 0. Defaults to 3600
  * `phone_verification_id_generator` - (Optional)
    * `length` - (Optional) The length of the secure generator used for generating the the phone verification Id. Defaults to 32
    * `type` - (Optional) The type of the secure generator used for generating the phone verification Id. Defaults to randomBytes.
  * `phone_verification_id_time_to_live_in_seconds` - (Optional) The time in seconds until a phone verification Id is no longer valid and cannot be used by the Verify Phone API. Value must be greater than 0. Defaults to 86400.
  * `phone_verification_one_time_code_generator` - (Optional)
    * `length` - (Optional) The length of the secure generator used for generating the phone verification one time code.. Defaults to 6
    * `type` - (Optional) The type of the secure generator used for generating the phone verification one time code. Defaults to randomAlphaNumeric.
  * `registration_verification_id_generator` - (Optional)
    * `length` - (Optional) The length of the secure generator used for generating the change password Id. Defaults to 32
    * `type` - (Optional) The type of the secure generator used for generating the change password Id. Defaults to randomBytes.
  * `registration_verification_id_time_to_live_in_seconds` - (Optional) The time in seconds until a registration verification Id is no longer valid and cannot be used by the Verify Registration API. Value must be greater than 0.
  * `registration_verification_one_time_code_generator` - (Optional)
    * `length` - (Optional) The length of the secure generator used for generating the registration verification one time code. Defaults to 6.
    * `type` - (Optional) The type of the secure generator used for generating the registration verification one time code. Defaults to randomAlphaNumeric.
  * `remember_oauth_scope_consent_choice_time_to_live_in_seconds` - (Optional) The time in seconds until remembered OAuth scope consent choices are no longer valid, and the User will be prompted to consent to requested OAuth scopes even if they have not changed. Applies only when `application.oauthConfiguration.consentMode` is set to RememberDecision. Value must be greater than 0. Note: An Essentials or Enterprise plan is required to utilize advanced OAuth scopes. Defaults to 2592000.
  * `saml_v2_authn_request_id_ttl_seconds` - (Optional) The time in seconds that a SAML AuthN request will be eligible for use to authenticate with FusionAuth. Defaults to 300.
  * `setup_password_id_generator` - (Optional)
    * `length` - (Optional) The length of the secure generator used for generating the change password Id. Defaults to 32.
    * `type` - (Optional) The type of the secure generator used for generating the change password Id. Defaults to randomBytes.
  * `setup_password_id_time_to_live_in_seconds` - (Optional) The time in seconds until a setup password Id is no longer valid and cannot be used by the Change Password API. Value must be greater than 0.
  * `trust_token_time_to_live_in_seconds` - (Optional) The number of seconds before the Trust Token is no longer valid to complete a request that requires trust. Value must be greater than 0. Defaults to 180
  * `two_factor_id_time_to_live_in_seconds` - (Optional) The time in seconds until a two factor Id is no longer valid and cannot be used by the Two Factor Login API. Value must be greater than 0. Defaults to 300.
  * `two_factor_one_time_code_id_generator` - (Optional)
    * `length` - (Optional) The length of the secure generator used for generating the the two factor code Id. Defaults to 6
    * `type` - (Optional) The type of the secure generator used for generating the two factor one time code Id. Defaults to randomDigits.
  * `two_factor_one_time_code_id_time_to_live_in_seconds` - (Optional) The number of seconds before the Two-Factor One Time Code used to enable or disable a two-factor method is no longer valid. Must be greater than 0. Defaults to 60.
  * `two_factor_trust_id_time_to_live_in_seconds` - (Optional) The time in seconds until an issued Two Factor trust Id is no longer valid and the User will be Optional to complete Two Factor authentication during the next authentication attempt. Value must be greater than 0.
  * `webauthn_authentication_challenge_time_to_live_in_seconds` - (Optional) The time in seconds until a WebAuthn authentication challenge is no longer valid and the User will be required to restart the WebAuthn authentication ceremony by creating a new challenge. This value also controls the timeout for the client-side WebAuthn navigator.credentials.get API call. Value must be greater than 0. Note: A license is required to utilize WebAuthn. Defaults to 180.
  * `webauthn_registration_challenge_time_to_live_in_seconds` - (Optional) The time in seconds until a WebAuthn registration challenge is no longer valid and the User will be required to restart the WebAuthn registration ceremony by creating a new challenge. This value also controls the timeout for the client-side WebAuthn navigator.credentials.create API call. Value must be greater than 0. Note: A license is required to utilize WebAuthn. Defaults to 180.
* `failed_authentication_configuration` - (Optional)
  * `action_duration` - (Required) The duration of the User Action. This value along with the actionDurationUnit will be used to set the duration of the User Action. Value must be greater than 0.
  * `action_duration_unit` - (Optional) The unit of time associated with a duration.
  * `reset_count_in_seconds` - (Optional) The length of time in seconds before the failed authentication count will be reset. Value must be greater than 0.
  * `too_many_attempts` - (Optional) The number of failed attempts considered to be too many. Once this threshold is reached the specified User Action will be applied to the user for the duration specified. Value must be greater than 0.
  * `action_cancel_policy_on_password_reset` - (Optional) Indicates whether you want the user to be able to self-service unlock their account prior to the action duration by completing a password reset workflow.
  * `email_user` - (Optional) Indicates you would like to email the user when the user’s account is locked due to this action being taken. This requires the User Action specified by the tenant.failedAuthenticationConfiguration.userActionId to also be configured for email. If the User Action is not configured to be able to email the user, this configuration will be ignored.
  * `user_action_id` - (Optional) The Id of the User Action that is applied when the threshold is reached for too many failed authentication attempts.
* `family_configuration` - (Optional)
  * `allow_child_registrations` - (Optional) Whether to allow child registrations.
  * `confirm_child_email_template_id` - (Optional) The unique Id of the email template to use when confirming a child.
  * `delete_orphaned_accounts` - (Optional) Indicates that child users without parental verification will be permanently deleted after tenant.familyConfiguration.deleteOrphanedAccountsDays days.
  * `delete_orphaned_accounts_days` - (Optional) The number of days from creation child users will be retained before being deleted for not completing parental verification. Value must be greater than 0.
  * `enabled` - (Optional) Whether family configuration is enabled.
  * `family_request_email_template_id` - (Optional) The unique Id of the email template to use when a family request is made.
  * `maximum_child_age` - (Optional) The maximum age of a child. Value must be greater than 0.
  * `minimum_owner_age` - (Optional) The minimum age to be an owner. Value must be greater than 0.
  * `parent_email_required` - (Optional) Whether a parent email is required.
  * `parent_registration_email_template_id` - (Optional) The unique Id of the email template to use for parent registration.
* `form_configuration` - (Optional)
  * `admin_user_form_id` - (Optional) The unique Id of the form to use for the Add and Edit User form when used in the FusionAuth admin UI.
* `http_session_max_inactive_interval` - (Optional) Time in seconds until an inactive session will be invalidated. Used when creating a new session in the FusionAuth OAuth frontend.
* `issuer` - (Optional) The named issuer used to sign tokens, this is generally your public fully qualified domain.
* `jwt_configuration` - (Optional) The JWT configuration for the tenant.
  * `access_token_key_id` - (Optional) The unique id of the signing key used to sign the access token. Required prior to `1.30.0`.
  * `id_token_key_id` - (Optional) The unique id of the signing key used to sign the Id token. Required prior to `1.30.0`.
  * `refresh_token_expiration_policy` - (Optional) The refresh token expiration policy.
  * `refresh_token_one_time_use_configuration_grace_period_in_seconds` - (Optional) The length of time specified in seconds that a one-time use token can be reused. This value must be greater than 0 and less than 86400 which is equal to 24 hours. Setting this value to 0 effectively disables the grace period which means a one-time token may not be reused. Defaults to 0.
  * `refresh_token_revocation_policy_on_login_prevented` - (Optional) When enabled, the refresh token will be revoked when a user action, such as locking an account based on a number of failed login attempts, prevents user login.
  * `refresh_token_revocation_policy_on_multi_factor_enable` - (Optional) When enabled, all refresh tokens will be revoked when a user enables multi-factor authentication for the first time. This policy will not be applied when adding subsequent multi-factor methods to the user.
  * `refresh_token_revocation_policy_on_one_time_token_reuse` - (Optional) When enabled, if a one-time use refresh token is reused, the token will be revoked. This does not cause all refresh tokens to be revoked, only the reused token is revoked.
  * `refresh_token_revocation_policy_on_password_change` - (Optional) When enabled, the refresh token will be revoked when a user changes their password."
  * `refresh_token_sliding_window_maximum_time_to_live_in_minutes` - (Optional) The maximum lifetime of a refresh token when using a refresh token expiration policy of SlidingWindowWithMaximumLifetime. Value must be greater than 0.
  * `refresh_token_time_to_live_in_minutes` - (Optional) The length of time in minutes a Refresh Token is valid from the time it was issued. Value must be greater than 0.
  * `refresh_token_usage_policy` - (Optional) The refresh token usage policy.
  * `time_to_live_in_seconds` - (Optional) The length of time in seconds this JWT is valid from the time it was issued. Value must be greater than 0.
* `lambda_configuration` - (Optional) Lamnda configuration for this tenant.
  * `login_validation_id` - (Required) The Id of the lambda that will be invoked at the end of a successful login request in order to extend custom validation of a login request.
  * `scim_enterprise_user_request_converter_id` - (Required) The Id of a SCIM User Request lambda that will be used to convert the SCIM Enterprise User request to a FusionAuth User. Note: An Enterprise plan is required to utilize SCIM. Required when `scim_server_configuration.enabled` is true.
  * `scim_enterprise_user_response_converter_id` - (Required) The Id of a SCIM User Response lambda that will be used to convert a FusionAuth Enterprise User to a SCIM Server response. Note: An Enterprise plan is required to utilize SCIM. Required when `scim_server_configuration.enabled` is true.
  * `scim_group_request_converter_id` - (Required) The Id of a SCIM Group Request lambda that will be used to convert the SCIM Group request to a FusionAuth Group. Note: An Enterprise plan is required to utilize SCIM. Required when `scim_server_configuration.enabled` is true.
  * `scim_group_response_converter_id` - (Required) The Id of a SCIM Group Response lambda that will be used to convert a FusionAuth Group to a SCIM Server response. Note: An Enterprise plan is required to utilize SCIM. Required when `scim_server_configuration.enabled` is true.
  * `scim_user_request_converter_id` - (Required) The Id of a SCIM User Request lambda that will be used to convert the SCIM User request to a FusionAuth User. Note: An Enterprise plan is required to utilize SCIM. Required when `scim_server_configuration.enabled` is true.
  * `scim_user_response_converter_id` - (Required) The Id of a SCIM User Response lambda that will be used to convert a FusionAuth User to a SCIM Server response. Note: An Enterprise plan is required to utilize SCIM. Required when `scim_server_configuration.enabled` is true.
* `login_configuration`
  * `require_authentication` - (Optional) Indicates whether to require an API key for the Login API when an `applicationId` is not provided. When an `applicationId` is provided to the Login API call, the application configuration will take precedence. In almost all cases, you will want to this to be `true`.
* `logout_url` - (Optional) The logout redirect URL when sending the user’s browser to the /oauth2/logout URI of the FusionAuth Front End. This value is only used when a logout URL is not defined in your Application.
* `maximum_password_age` - (Optional)
  * `days` - (Optional) The password maximum age in days. The number of days after which FusionAuth will require a user to change their password. Required when systemConfiguration.maximumPasswordAge.enabled is set to true.
  * `enabled` - (Optional) Indicates that the maximum password age is enabled and being enforced.
* `minimum_password_age` - (Optional)
  * `enabled` - (Optional) Indicates that the minimum password age is enabled and being enforced.
  * `seconds` - (Optional) The password minimum age in seconds. When enabled FusionAuth will not allow a password to be changed until it reaches this minimum age. Required when systemConfiguration.minimumPasswordAge.enabled is set to true.
* `multi_factor_configuration` - (Optional)
  * `authenticator` - (Optional)
    * `enabled` - (Optional) When enabled, users may utilize an authenticator application to complete a multi-factor authentication request. This method uses TOTP (Time-Based One-Time Password) as defined in RFC 6238 and often uses an native mobile app such as Google Authenticator.
  * `email` - (Optional)
    * `enabled` - (Optional) When enabled, users may utilize an email address to complete a multi-factor authentication request.
    * `template_id` - (Optional) The Id of the email template that is used when notifying a user to complete a multi-factor authentication request.
  * `login_policy` - (Optional)  When set to `Enabled` and a user has one or more two-factor methods configured, the user will be required to complete a two-factor challenge during login. When set to `Disabled`, even when a user has configured one or more two-factor methods, the user will not be required to complete a two-factor challenge during login. When the login policy is to `Required`, a two-factor challenge will be required during login. If a user does not have configured two-factor methods, they will not be able to log in.
  * `sms` - (Optional)
    * `enabled` - (Optional) When enabled, users may utilize a mobile phone number to complete a multi-factor authentication request.
    * `messenger_id` - (Optional) The messenger that is used to deliver a SMS multi-factor authentication request.
    * `template_id` - (Optional) The Id of the SMS template that is used when notifying a user to complete a multi-factor authentication request.
* `oauth_configuration` - (Optional)
  * `client_credentials_access_token_populate_lambda_id` - (Optional) The Id of a lambda that will be called to populate the JWT during a client credentials grant. **Note:** A paid edition of FusionAuth is required to utilize client credentials grant.
* `password_enabled` - (Optional) Indicates whether the password is enabled for this tenant. This value is used to determine if the password is required when registering a new user or updating an existing user. Defaults to `true`.
* `password_encryption_configuration` - (Optional)
  * `encryption_scheme` - (Optional) The default method for encrypting the User’s password.
  * `encryption_scheme_factor` - (Optional) The factor used by the password encryption scheme. If not provided, the PasswordEncryptor provides a default value. Generally this will be used as an iteration count to generate the hash. The actual use of this value is up to the PasswordEncryptor implementation.
  * `modify_encryption_scheme_on_login` - (Optional) When enabled a user’s hash configuration will be modified to match these configured settings. This can be useful to increase a password hash strength over time or upgrade imported users to a more secure encryption scheme after an initial import.
* `password_validation_rules` - (Optional)
  * `breach_detection` - (Optional)
    * `enabled` - (Optional) Whether to enable Reactor breach detection. Requires an activated license.
    * `match_mode` - (Optional) The level of severity where Reactor will consider a breach.
    * `notify_user_email_template_id` - (Optional) The Id of the email template to use when notifying user of breached password. Required if tenant.passwordValidationRules.breachDetection.onLogin is set to NotifyUser.
    * `on_login` - (Optional) The behavior when detecting breaches at time of user login
  * `max_length` - (Optional) The maximum length of a password when a new user is created or a user requests a password change. This value must be greater than 0 and less than or equal to 256. When `passwordEncryptionConfiguration.encryptionScheme` is equal to `bcrypt`, the maximum will be limited to 50.
  * `min_length` - (Optional) The minimum length of a password when a new user is created or a user requests a password change.
  * `remember_previous_passwords` - (Optional)
    * `count` - (Optional) The number of previous passwords to remember. Value must be greater than 0.
    * `enabled` - (Optional) Whether to prevent a user from using any of their previous passwords.
  * `required_mixed_case` - (Optional) Whether to force the user to use at least one uppercase and one lowercase character.
  * `require_non_alpha` - (Optional) Whether to force the user to use at least one non-alphanumeric character.
  * `require_number` - (Optional) Whether to force the user to use at least one number.
  * `validate_on_login` - (Optional) When enabled the user’s password will be validated during login. If the password does not meet the currently configured validation rules the user will be required to change their password.
* `phone_configuration` - (Optional)
  * `messenger_id` - (Optional) The messenger that is used to deliver SMS messages for phone number verification and passwordless logins. This field is required when any of `tenant.phone_configuration.passwordless_template_id` , `tenant.phone_configuration.verification_complete_template_id` , or `tenant.phone_configuration.verification_template_id` is set.
  * `passwordless_template_id` - (Optional) The Id of the Passwordless Message Template, sent to users when they start a passwordless login.
  * `unverified` - (Optional)
    * `behavior` - (Optional) The desired behavior during login for a user that does not have a verified phone number. The possible values are: `Allow` and `Gated`. Defaults to `Allow`.
  * `verification_complete_template_id` - (Optional) The Id of the Message Template used to notify a user that their phone number has been verified.
  * `verification_strategy` - (Optional) The process by which the user will verify their phone number. The possible values are: `ClickableLink` and `FormField`. Defaults to `ClickableLink`.
  * `verification_template_id` - (Optional) The Id of the Message Template used to send SMS messages to users to verify that their phone number is valid.
  * `verify_phone_number` - (Optional) Whether a user’s phone number is verified when they register with your application. Defaults to `false`.
* `rate_limit_configuration` - (Optional)
  * `failed_login` - (Optional)
    * `enabled` -  (Optional) Whether rate limiting is enabled for failed login.
    * `limit` -  (Optional) The number of times a user can fail to login within the configured `time_period_in_seconds` duration. If a Failed authentication action has been configured then it will take precedence.
    * `time_period_in_seconds` - (Optional) The duration for the number of times a user can fail login before being rate limited.
  * `forgot_password` - (Optional)
    * `enabled` - (Optional) Whether rate limiting is enabled for forgot password.
    * `limit` - (Optional) The number of times a user can request a forgot password email within the configured `time_period_in_seconds` duration.
    * `time_period_in_seconds` - (Optional) The duration for the number of times a user can request a forgot password email before being rate limited.
  * `send_email_verification` - (Optional)
    * `enabled` - (Optional) Whether rate limiting is enabled for send email verification.
    * `limit` - (Optional) The number of times a user can request a verification email within the configured `time_period_in_seconds` duration.
    * `time_period_in_seconds` - (Optional) The duration for the number of times a user can request a verification email before being rate limited.
  * `send_passwordless` - (Optional)
    * `enabled` - (Optional) Whether rate limiting is enabled for send passwordless.
    * `limit` - (Optional) The number of times a user can request a passwordless login email within the configured `time_period_in_seconds` duration.
    * `time_period_in_seconds` - (Optional) The duration for the number of times a user can request a passwordless login email before being rate limited.
  * `send_phone_verification` - (Optional)
    * `enabled` - (Optional) Whether rate limiting is enabled for send phone verification.
    * `limit` - (Optional) The number of times a user can request a phone verification message within the configured `time_period_in_seconds` duration. Value must be greater than 0.
    * `time_period_in_seconds` - (Optional) The duration for the number of times a user can request a phone verification message before being rate limited. Value must be greater than 0.
  * `send_registration_verification` - (Optional)
    * `enabled` - (Optional) Whether rate limiting is enabled for send registration verification.
    * `limit` - (Optional) The number of times a user can request a registration verification email within the configured `time_period_in_seconds` duration.
    * `time_period_in_seconds` - (Optional) The duration for the number of times a user can request a registration verification email before being rate limited.
  * `send_two_factor` - (Optional)
    * `enabled` - (Optional) Whether rate limiting is enabled for send two factor.
    * `limit` - (Optional) The number of times a user can request a two-factor code by email or SMS within the configured `time_period_in_seconds` duration.
    * `time_period_in_seconds` - (Optional) The duration for the number of times a user can request a two-factor code by email or SMS before being rate limited.
* `registration_configuration` - (Optional)
  * `blocked_domains` - (Optional) A list of unique domains that are not allowed to register when self service is enabled.
* `scim_server_configuration` - (Optional)
  * `client_entity_type_id` - (Required) The Entity Type that will be used to represent SCIM Clients for this tenant. Note: An Enterprise plan is required to utilize SCIM. Required when `scim_server_configuration.enabled` is true.
  * `enabled` - (Optional) Whether or not this tenant has the SCIM endpoints enabled. Note: An Enterprise plan is required to utilize SCIM.
  * `schemas` - (Optional) SON formatted as a SCIM Schemas endpoint response. Because the SCIM lambdas may modify the JSON response, ensure the Schema's response matches that generated by the response lambdas. More about Schema definitions. When this parameter is not provided, it will default to EnterpriseUser, Group, and User schema definitions as defined by the SCIM core schemas spec. Note: An Enterprise plan is required to utilize SCIM.
  * `server_entity_type_id` - (Required) The Entity Type that will be used to represent SCIM Servers for this tenant. Note: An Enterprise plan is required to utilize SCIM. Required when `scim_server_configuration.enabled` is true.
* `source_tenant_id` - (Optional) The optional Id of an existing Tenant to make a copy of. If present, the tenant.id and tenant.name values of the request body will be applied to the new Tenant, all other values will be copied from the source Tenant to the new Tenant.
* `sso_configuration` - (Optional)
  * `allow_access_token_bootstrap` - (Optional) When enabled, an SSO session can be created after login by providing an access token as a bearer token in a request to the OAuth2 Authorize endpoint. Defaults to `false`.
  * `device_trust_time_to_live_in_seconds` - (Optional) The number of seconds before a trusted device is reset. When reset, a user is forced to complete captcha during login and complete two factor authentication if applicable.
* `tenant_id` - (Optional) The Id to use for the new Tenant. If not specified a secure random UUID will be generated.
* `theme_id` - (Optional) The unique Id of the theme to be used to style the login page and other end user templates.
* `username_configuration` - (Optional)
  * `unique` - (Optional) Indicates that users without a verified email address will be permanently deleted after tenant.userDeletePolicy.unverified.numberOfDaysToRetain days.
    * `enabled` - (Optional) When true, FusionAuth will handle username collisions by generating a random suffix.
    * `number_of_digits` - (Optional) The maximum number of digits to use when building a unique suffix for a username. A number will be randomly selected and will be 1 or more digits up to this configured value in length. For example, if this value is 5, the suffix will be a number between 00001 and 99999, inclusive.
    * `separator` - (Optional) A single character to use as a separator from the requested username and a unique suffix that is added when a duplicate username is detected. This value can be a single non-alphanumeric ASCII character.
    * `strategy` - (Optional) When enabled the user’s password will be validated during login. If the password does not meet the currently configured validation rules the user will be required to change their password.
* `user_delete_policy` - (Optional)
  * `unverified_enabled` - (Optional) Indicates that users without a verified email address will be permanently deleted after tenant.userDeletePolicy.unverified.numberOfDaysToRetain days.
  * `unverified_number_of_days_to_retain` - (Optional)
* `webauthn_configuration` - (Optional) The WebAuthn configuration for this tenant.
  * `bootstrap_workflow` - (Optional) The Bootstrap Workflow configuration.
    * `authenticator_attachment_preference` - (Optional) Determines the authenticator attachment requirement for WebAuthn passkey registration when using the bootstrap workflow. The possible values are: Any, CrossPlatform and Platform. Note: A license is required to utilize WebAuthn and an Enterprise plan is required to utilize WebAuthn cross-platform authenticators..
    * `enabled` - (Optional) Whether or not this tenant has the WebAuthn bootstrap workflow enabled. The bootstrap workflow is used when the user must "bootstrap" the authentication process by identifying themselves prior to the WebAuthn ceremony and can be used to authenticate from a new device using WebAuthn. Note: A license is required to utilize WebAuthn..
    * `user_verification_requirement` - (Optional) Determines the user verification requirement for WebAuthn passkey registration when using the bootstrap workflow. The possible values are: Discouraged, Preferred and Required. Note: A license is required to utilize WebAuthn..
  * `debug` - (Optional) Determines if debug should be enabled for this tenant to create an event log to assist in debugging WebAuthn errors. Note: A license is required to utilize WebAuthn..
  * `enabled` - (Optional) Whether or not this tenant has WebAuthn enabled globally.. Note: A license is required to utilize WebAuthn..
  * `reauthentication_workflow` - (Optional) The Reauthentication Workflow configuration.
    * `authenticator_attachment_preference` - (Optional) Determines the authenticator attachment requirement for WebAuthn passkey registration when using the reauthentication workflow. The possible values are:: Any, CrossPlatform and Platform. Note: A license is required to utilize WebAuthn and an Enterprise plan is required to utilize WebAuthn cross-platform authenticators..
    * `enabled` - (Optional) Whether or not this tenant has the WebAuthn reauthentication workflow enabled. The reauthentication workflow will automatically prompt a user to authenticate using WebAuthn for repeated logins from the same device. Note: A license is required to utilize WebAuthn..
    * `user_verification_requirement` - (Optional) Determines the user verification requirement for WebAuthn passkey registration when using the bootstrap workflow. The possible values are: Discouraged, Preferred and Required. Note: A license is required to utilize WebAuthn..
  * `relying_party_id` - (Optional) The value this tenant will use for the Relying Party Id in WebAuthn ceremonies. Passkeys can only be used to authenticate on sites using the same Relying Party Id they were registered with. This value must match the browser origin or be a registrable domain suffix of the browser origin. For example, if your domain is auth.piedpiper.com, you could use auth.piedpiper.com or piedpiper.com but not m.auth.piedpiper.com or com. When this parameter is omitted, FusionAuth will use null for the Relying Party Id in passkey creation and request options. A null value in the WebAuthn JavaScript API will use the browser origin. Note: A license is required to utilize WebAuthn.
  * `relying_party_name` - (Optional) The value this tenant will use for the Relying Party name in WebAuthn ceremonies. This value may be displayed by browser or operating system dialogs during WebAuthn ceremonies. When this parameter is omitted, FusionAuth will use the tenant.issuer value. Note: A license is required to utilize WebAuthn.
* `webhook_ids` - (Optional) An array of Webhook Ids. For Webhooks that are not already configured for All Tenants, specifying an Id on this request will indicate the associated Webhook should handle events for this tenant.
