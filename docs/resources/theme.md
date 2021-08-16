# Theme Resource

This Resource is used to create a role for an Application.

[Themes API]https://fusionauth.io/docs/v1/tech/apis/themes)

## Example Useage

```hcl
resource "fusionauth_theme" "mytheme" {
  default_messages                               = "[#ftl/]"
  name                                           = "my theme"
  email_complete                                 = "[#ftl/]"
  email_send                                     = "[#ftl/]"
  email_verify                                   = "[#ftl/]"
  helpers                                        = "[#ftl/]"
  oauth2_child_registration_not_allowed          = "[#ftl/]"
  oauth2_child_registration_not_allowed_complete = "[#ftl/]"
  oauth2_complete_registration                   = "[#ftl/]"
  oauth2_error                                   = "[#ftl/]"
  oauth2_logout                                  = "[#ftl/]"
  oauth2_two_factor                              = "[#ftl/]"
  password_change                                = "[#ftl/]"
  password_complete                              = "[#ftl/]"
  password_forgot                                = "[#ftl/]"
  password_sent                                  = "[#ftl/]"
  registration_complete                          = "[#ftl/]"
  registration_send                              = "[#ftl/]"
  registration_verify                            = "[#ftl/]"
  oauth2_register                                = "[#ftl/]"
  oauth2_device                                  = "[#ftl/]"
  oauth2_device_complete                         = "[#ftl/]"
  oauth2_passwordless                            = "[#ftl/]"
  oauth2_wait                                    = "[#ftl/]"
  oauth2_authorize                               = "[#ftl/]"
}
```

## Argument Reference

* `source_theme_id` - (Optional) The optional Id of an existing Theme to make a copy of. If present, the defaultMessages, localizedMessages, templates, and stylesheet from the source Theme will be copied to the new Theme.
* `default_messages` - (Optional) A properties file formatted String containing at least all of the message keys defined in the FusionAuth shipped messages file. Required if not copying an existing Theme.
* `localized_messages` - (Optional) A Map of localized versions of the messages. The key is the Locale and the value is a properties file formatted String.
* `name` - (Required) A unique name for the Theme.
* `stylesheet` - (Optional) A CSS stylesheet used to style the templates.
* `account_edit` - (Optional) A FreeMarker template that contains the account edit page.
* `account_index` - (Optional) A FreeMarker template that contains the account index page, this is the account landing page.
* `account_two_factor_disable` - (Optional) A FreeMarker template that contains the two factor disable form.
* `account_two_factor_enable` - (Optional) A FreeMarker template that contains the two factor enable form.
* `account_two_factor_index` - (Optional) A FreeMarker template that contains the two factor index page.
* `email_complete` - (Optional) A FreeMarker template that is rendered when the user requests the /email/complete page. This page is used after a user has verified their email address by clicking the URL in the email. After FusionAuth has updated their user object to indicate that their email was verified, the browser is redirected to this page.
* `email_send` - (Optional) A FreeMarker template that is rendered when the user requests the /email/send page. This page is used after a user has asked for the verification email to be resent. This can happen if the URL in the email expired and the user clicked it. In this case, the user can provide their email address again and FusionAuth will resend the email. After the user submits their email and FusionAuth re-sends a verification email to them, the browser is redirected to this page.
* `email_verify` - (Optional) A FreeMarker template that is rendered when the user requests the /email/verify page by clicking the URL from the verification email and the verificationId has expired. FusionAuth expires verificationId after a period of time (which is configurable). If the user has a URL from the verification email that has expired, this page will be rendered and the error will be displayed to the user.
* `helpers` - (Optional) A FreeMarker template that contains all of the macros and templates used by the rest of the loginTheme FreeMarker templates (i.e. oauth2Authorize). This allows you to configure the general layout of your UI configuration and login theme without having to copy and paste HTML into each of the templates.
* `index` - (Optional) A FreeMarker template that contains the main index page, this is the root landing page.
* `oauth2_authorize` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/authorize page. This is the main login page for FusionAuth and is used for all interactive OAuth and OpenId Connect workflows.
* `oauth2_child_registration_not_allowed` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/child-registration-not-allowed page. This is where the child must provide their parent’s email address to ask their parent to create an account for them.
* `oauth2_child_registration_not_allowed_complete` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/child-registration-not-allowed-complete page. This is where the browser is taken after the child provides their parent’s email address.
* `oauth2_complete_registration` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/complete-registration page. This page is used for users that have accounts but might be missing required fields.
* `oauth2_error` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/error page. This page is used if the user starts or is in the middle of the OAuth workflow and any type of error occurs. This could be caused by the user messing with the URL or internally some type of information wasn’t passed between the OAuth endpoints correctly. For example, if you are federating login to an external IdP and that IdP does not properly echo the state parameter, FusionAuth’s OAuth workflow will break and this page will be displayed.
* `oauth2_logout` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/logout page. This page is used if the user initiates a logout. This page causes the user to be logged out of all associated applications via a front-channel mechanism before being redirected.
* `oauth2_register` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/register page. This page is used for users that need to register (sign-up)
* `oauth2_device` - (Optional) A FreeMarker template
* `oauth2_device_complete` - (Optional) A FreeMarker template
* `oauth2_passwordless` - (Optional) A FreeMarker template
* `oauth2_wait` - (Optional) A FreeMarker template
* `oauth2_two_factor` - (Optional) A FreeMarker template that is rendered when the user requests the /oauth2/two-factor page. This page is used if the user has two-factor authentication enabled and they need to type in their code again. FusionAuth will properly handle the SMS or authenticator app processing on the back end. This page contains the form that the user will put their code into.
* `oauth2_two_factor_methods` - (Optional) A FreeMarker template that contains the OAuth2 two-factor option form.
* `password_change` - (Optional) A FreeMarker template that is rendered when the user requests the /password/change page. This page is used if the user is required to change their password or if they have requested a password reset. This page contains the form that allows the user to provide a new password.
* `password_complete` - (Optional) A FreeMarker template that is rendered when the user requests the /password/complete page. This page is used after the user has successfully updated their password (or reset it). This page should instruct the user that their password was updated and that they need to login again.
* `password_forgot` - (Optional) A FreeMarker template that is rendered when the user requests the /password/forgot page. This page is used when a user starts the forgot password workflow. This page renders the form where the user types in their email address.
* `password_sent` - (Optional) A FreeMarker template that is rendered when the user requests the /password/sent page. This page is used when a user has submitted the forgot password form with their email. FusionAuth does not indicate back to the user if their email address was valid in order to prevent malicious activity that could reveal valid email addresses. Therefore, this page should indicate to the user that if their email was valid, they will receive an email shortly with a link to reset their password.
* `registration_complete` - (Optional) A FreeMarker template that is rendered when the user requests the /registration/complete page. This page is used after a user has verified their email address for a specific application (i.e. a user registration) by clicking the URL in the email. After FusionAuth has updated their registration object to indicate that their email was verified, the browser is redirected to this page.
* `registration_send` - (Optional)
* `registration_verify` - (Optional) A FreeMarker template that is rendered when the user requests the /registration/verify page by clicking the URL from the application specific verification email and the verificationId has expired. FusionAuth expires verificationId after a period of time (which is configurable). If the user has a URL from the verification email that has expired, this page will be rendered and the error will be displayed to the user.
* `samlv2_logout` - (Optional)
* `email_sent` - (Optional)
* `email_verification_required` - (Optional)
* `registration_sent` - (Optional)
* `oauth2_authorized_not_registered` - (Optional)
* `oauth2_start_idp_link` - (Optional)
* `registration_verification_required` - (Optional)
* `unauthorized` - (Optional)