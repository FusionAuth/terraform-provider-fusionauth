# Theme Resource

UI login themes can be configured to enable custom branding for your FusionAuth login workflow. Themes are configured per Tenant or optionally by Application.

[Themes API](https://fusionauth.io/docs/v1/tech/apis/themes)

## Example Usage

```hcl
resource "fusionauth_theme" "mytheme" {
  default_messages                               = "[#ftl/]"
  name                                           = "my theme"
  stylesheet                                     = "/* stylish */"
  account_edit                                   = "[#ftl/]"
  account_index                                  = "[#ftl/]"
  account_two_factor_disable                     = "[#ftl/]"
  account_two_factor_enable                      = "[#ftl/]"
  account_two_factor_index                       = "[#ftl/]"
  account_webauthn_add                           = "[#ftl/]"
  account_webauthn_delete                        = "[#ftl/]"
  account_webauthn_index                         = "[#ftl/]"
  confirmation_required                          = "[#ftl/]"
  email_complete                                 = "[#ftl/]"
  email_sent                                     = "[#ftl/]"
  email_verification_required                    = "[#ftl/]"
  email_verify                                   = "[#ftl/]"
  helpers                                        = "[#ftl/]"
  index                                          = "[#ftl/]"
  oauth2_authorize                               = "[#ftl/]"
  oauth2_authorized_not_registered               = "[#ftl/]"
  oauth2_child_registration_not_allowed          = "[#ftl/]"
  oauth2_child_registration_not_allowed_complete = "[#ftl/]"
  oauth2_complete_registration                   = "[#ftl/]"
  oauth2_consent                                 = "[#ftl/]"
  oauth2_device                                  = "[#ftl/]"
  oauth2_device_complete                         = "[#ftl/]"
  oauth2_error                                   = "[#ftl/]"
  oauth2_logout                                  = "[#ftl/]"
  oauth2_passwordless                            = "[#ftl/]"
  oauth2_register                                = "[#ftl/]"
  oauth2_start_idp_link                          = "[#ftl/]"
  oauth2_two_factor                              = "[#ftl/]"
  oauth2_two_factor_methods                      = "[#ftl/]"
  oauth2_two_factor_enable                       = "[#ftl/]"
  oauth2_two_factor_enable_complete              = "[#ftl/]"
  oauth2_wait                                    = "[#ftl/]"
  oauth2_webauthn                                = "[#ftl/]"
  oauth2_webauthn_reauth                         = "[#ftl/]"
  oauth2_webauthn_reauth_enable                  = "[#ftl/]"
  password_change                                = "[#ftl/]"
  password_complete                              = "[#ftl/]"
  password_forgot                                = "[#ftl/]"
  password_sent                                  = "[#ftl/]"
  registration_complete                          = "[#ftl/]"
  registration_sent                              = "[#ftl/]"
  registration_verification_required             = "[#ftl/]"
  registration_verify                            = "[#ftl/]"
  samlv2_logout                                  = "[#ftl/]"
  unauthorized                                   = "[#ftl/]"

  # Deprecated Properties
  email_send                                     = "[#ftl/]"
  registration_send                              = "[#ftl/]"
}
```

## Argument Reference

* `source_theme_id` - (Optional) The optional Id of an existing Theme to make a copy of. If present, the defaultMessages, localizedMessages, templates, and stylesheet from the source Theme will be copied to the new Theme.
* `default_messages` - (Optional) A properties file formatted String containing at least all of the message keys defined in the FusionAuth shipped messages file. 

~> **Note:** `default_messages` Is Required if not copying an existing Theme.

* `localized_messages` - (Optional) A Map of localized versions of the messages. The key is the Locale and the value is a properties file formatted String.
* `name` - (Required) A unique name for the Theme.
* `stylesheet` - (Optional) A CSS stylesheet used to style the templates.
* `account_edit` - (Optional) A FreeMarker template that is rendered when the user requests the /account/edit path. This page contains a form that enables authenticated users to update their profile.
* `account_index` - (Optional) A FreeMarker template that is rendered when the user requests the /account path. This is the self-service account landing page. An authenticated user may use this as a starting point for operations such as updating their profile or configuring multi-factor authentication.
* `account_two_factor_disable` - (Optional) A FreeMarker template that is rendered when the user requests the /account/two-factor/disable path. This page contains a form that accepts a verification code used to disable a multi-factor authentication method.
* `account_two_factor_enable` - (Optional) A FreeMarker template that is rendered when the user requests the /account/two-factor/enable path. This page contains a form that accepts a verification code used to enable a multi-factor authentication method. Additionally, this page contains presentation of recovery codes when a user enables multi-factor authentication for the first time.
* `account_two_factor_index` - (Optional) A FreeMarker template that is rendered when the user requests the /account/two-factor path. This page displays an authenticated user’s configured multi-factor authentication methods. Additionally, it provides links to enable and disable a method.
* `account_webauthn_add` - (Optional) A FreeMarker template that is rendered when the user requests the /account/webauthn/add path. This page contains a form that allows a user to register a new WebAuthn passkey.
* `account_webauthn_delete` - (Optional) A FreeMarker template that is rendered when the user requests the /account/webauthn/delete path. This page contains a form that allows a user to delete a WebAuthn passkey.
* `account_webauthn_index` - (Optional) A FreeMarker template that is rendered when the user requests the /account/webauthn/ path. This page displays an authenticated user’s registered WebAuthn passkeys. Additionally, it provides links to delete an existing passkey and register a new passkey.
* `confirmation_required` - (Optional) A FreeMarker template that is rendered when the user requests the /confirmation-required path. This page is displayed when a user attempts to complete an email based workflow that did not begin in the same browser. For example, if the user starts a forgot password workflow, and then opens the link in a separate browser the user will be shown this panel.
* `email_complete` - (Optional) A FreeMarker template that is rendered when the user requests the /email/complete path. This page is used after a user has verified their email address by clicking the URL in the email. After FusionAuth has updated their user object to indicate that their email was verified, the browser is redirected to this page.
* `email_sent` - (Optional) A FreeMarker template that is rendered when the user requests the /email/sent path. This page is used after a user has asked for the verification email to be resent. This can happen if the URL in the email expired and the user clicked it. In this case, the user can provide their email address again and FusionAuth will resend the email. After the user submits their email and FusionAuth re-sends a verification email to them, the browser is redirected to this page.
* `email_verification_required` - (Optional) A FreeMarker template that is rendered when the user requests the /email/verification-required path. This page is rendered when a user is required to verify their email address prior to being allowed to proceed with login. This occurs when Unverified behavior is set to Gated in email verification settings on the Tenant.
* `email_verify` - (Optional) A FreeMarker template that is rendered when the user requests the /email/verify path. This page is rendered when a user clicks the URL from the verification email and the verificationId has expired. FusionAuth expires verificationId after a period of time (which is configurable). If the user has a URL from the verification email that has expired, this page will be rendered and the error will be displayed to the user.
* `helpers` - (Optional) A FreeMarker template that contains all of the macros and templates used by the rest of the login Theme FreeMarker templates. This allows you to configure the general layout of your UI configuration and login theme without having to copy and paste HTML into each of the templates.
* `index` - (Optional) A FreeMarker template that is rendered when the user requests the / path. This is the root landing page. This page is available to unauthenticated users and will be displayed whenever someone navigates to the FusionAuth host’s root page. Prior to version 1.27.0, navigating to this URL would redirect to /admin and would subsequently render the FusionAuth admin login page.
* `oauth2_authorize` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/authorize path. This is the main login page for FusionAuth and is used for all interactive OAuth2 and OpenID Connect workflows.
* `oauth2_authorized_not_registered` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/authorized-not-registered path. This page is rendered when a user is not registered and the Application configuration requires registration before FusionAuth will complete the redirect.
* `oauth2_child_registration_not_allowed` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/child-registration-not-allowed path. This page contains a form where a child must provide their parent’s email address to ask their parent to create an account for them in a Consent workflow.
* `oauth2_child_registration_not_allowed_complete` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/child-registration-not-allowed-complete path. This page is rendered is rendered after a child provides their parent’s email address for parental consent in a Consent workflow.
* `oauth2_complete_registration` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/complete-registration path. This page contains a form that is used for users that have accounts but might be missing required fields.
* `oauth2_consent` - (Optional) A FreeMarker template that is rendered when a third party application requests scopes from the user.
* `oauth2_device` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/device path. This page contains a form for accepting an end user’s short code for the interactive portion of the OAuth Device Authorization Grant workflow.
* `oauth2_device_complete` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/device-complete path. This page contains a complete message indicating the device authentication has completed.
* `oauth2_error` - (Optional) This page is used if the user starts or is in the middle of the OAuth workflow and any type of error occurs. This could be caused by the user messing with the URL or internally some type of information wasn’t passed between the OAuth endpoints correctly. For example, if you are federating login to an external IdP and that IdP does not properly echo the state parameter, FusionAuth’s OAuth workflow will break and this page will be displayed.
* `oauth2_logout` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/logout page. This page is used if the user initiates a logout. This page causes the user to be logged out of all associated applications via a front-channel mechanism before being redirected.
* `oauth2_passwordless` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/passwordless path. This page is rendered when the user starts the passwordless login workflow. The page renders the form where the user types in their email address.
* `oauth2_register` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/register path. This page is used to register or sign up the user for the application when self-service registration is enabled.
* `oauth2_start_idp_link` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/start-idp-link path. This page is used if the Identity Provider is configured to have a pending link. The user is presented with the option to link their account with an existing FusionAuth user account.
* `oauth2_two_factor` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/two-factor path. This page is used if the user has two-factor authentication enabled and they need to type in their code again. FusionAuth will properly handle the processing on the back end. This page contains the form that the user will put their code into.
* `oauth2_two_factor_methods` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/two-factor-methods path. This page contains a form providing a user with their configured multi-factor authentication options that they may use to complete the authentication challenge.
* `oauth2_two_factor_enable` - (Optional) A FreeMarker template that contains the OAuth2 two-factor enable form.
* `oauth2_two_factor_enable_complete` - (Optional) A FreeMarker template that contains the OAuth2 two-factor enable complete form.
* `oauth2_wait` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/wait path. This page is rendered when FusionAuth is waiting for an external provider to complete an out of band authentication request. For example, during a HYPR login this page will be displayed until the user completes authentication.
* `oauth2_webauthn` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/webauthn path. This page contains a form where a user can enter their loginId (username or email address) to authenticate with one of their registered WebAuthn passkeys. This page uses the WebAuthn bootstrap workflow.
* `oauth2_webauthn_reauth` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/webauthn-reauth path. This page contains a form that lists the WebAuthn passkeys currently available for re-authentication. A user can select one of the listed passkeys to authenticate using the corresponding passkey and user account.
* `oauth2_webauthn_reauth_enable` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/webauthn-reauth-enable path. This page contains two forms. One allows the user to select one of their existing WebAuthn passkeys to use for re-authentication. The other allows the user to register a new WebAuthn passkey for re-authentication.
* `password_change` - (Optional) A FreeMarker template that is rendered when the user requests the /password/change path. This page is used if the user is required to change their password or if they have requested a password reset. This page contains the form that allows the user to provide a new password.
* `password_complete` - (Optional) A FreeMarker template that is rendered when the user requests the /password/complete path. This page is used after the user has successfully updated their password, or reset it. This page should instruct the user that their password was updated and that they need to login again.
* `password_forgot` - (Optional) A FreeMarker template that is rendered when the user requests the /password/forgot path. This page is used when a user starts the forgot password workflow. This page renders the form where the user types in their email address.
* `password_sent` - (Optional) A FreeMarker template that is rendered when the user requests the /password/sent path. This page is used when a user has submitted the forgot password form with their email. FusionAuth does not indicate back to the user if their email address was valid in order to prevent malicious activity that could reveal valid email addresses. Therefore, this page should indicate to the user that if their email was valid, they will receive an email shortly with a link to reset their password.
* `registration_complete` - (Optional) A FreeMarker template that is rendered when the user requests the /registration/complete path. This page is used after a user has verified their email address for a specific application (i.e. a user registration) by clicking the URL in the email. After FusionAuth has updated their registration object to indicate that their email was verified, the browser is redirected to this page.
* `registration_sent` - (Optional) A FreeMarker template that is rendered when the user requests the /registration/sent path. This page is used after a user has asked for the application specific verification email to be resent. This can happen if the URL in the email expired and the user clicked it. In this case, the user can provide their email address again and FusionAuth will resend the email. After the user submits their email and FusionAuth re-sends a verification email to them, the browser is redirected to this page.
* `registration_verification_required` - (Optional) A FreeMarker template that is rendered when the user requests the /registration/verification-required path. This page is rendered when a user is required to verify their registration prior to being allowed to proceed with the registration flow. This occurs when Unverified behavior is set to Gated in registration verification settings on the Application.
* `registration_verify` - (Optional) A FreeMarker template that is rendered when the user requests the /registration/verify path. This page is used when a user clicks the URL from the application specific verification email and the verificationId has expired. FusionAuth expires verificationId after a period of time (which is configurable). If the user has a URL from the verification email that has expired, this page will be rendered and the error will be displayed to the user.
* `samlv2_logout` - (Optional) A FreeMarker template that is rendered when the user requests the /samlv2/logout path. This page is used if the user initiates a SAML logout. This page causes the user to be logged out of all associated applications via a front-channel mechanism before being redirected.
* `unauthorized` - (Optional) An optional FreeMarker template that contains the unauthorized page.

### Deprecated Theme Properties
* `email_send` - (Optional) A FreeMarker template that is rendered when the user requests the /email/send page. This page is used after a user has asked for the verification email to be resent. This can happen if the URL in the email expired and the user clicked it. In this case, the user can provide their email address again and FusionAuth will resend the email. After the user submits their email and FusionAuth re-sends a verification email to them, the browser is redirected to this page.
* `registration_send` - (Optional) A FreeMarker template that is rendered when the user requests the /registration/send page. This page is used after a user has asked for the application specific verification email to be resent. This can happen if the URL in the email expired and the user clicked it. In this case, the user can provide their email address again and FusionAuth will resend the email. After the user submits their email and FusionAuth re-sends a verification email to them, the browser is redirected to this page.
