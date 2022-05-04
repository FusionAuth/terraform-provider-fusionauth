package testdata

import "fmt"

// MessageProperties supplies the default fusionauth properties file.
func MessageProperties(name string) string {
	if name == "" {
		name = "FusionAuth"
	}

	return fmt.Sprintf(`
#
# Copyright (c) 2019-2021, FusionAuth, All Rights Reserved
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
# either express or implied. See the License for the specific
# language governing permissions and limitations under the License.
#

#
# Text used on the page (inside the HTML). You can create new key-value pairs here and use them in the templates.
#
access-denied=Access denied
account=Account
add-two-factor=Add two-factor
back-to-login=Return to Login
cancel=Cancel
captcha-google-branding=This site is protected by reCAPTCHA and the Google <a href="https://policies.google.com/privacy">Privacy Policy</a> and <a href="https://policies.google.com/terms">Terms of Service</a> apply.
authorized-not-registered=Registration is required to access this application and your account has not been registered for this application. Please complete your registration and try again.
authorized-not-registered-title=Registration Required
cancel-link=Cancel link request
child-registration-not-allowed=We cannot create an account for you. Your parent or guardian can create an account for you. Enter their email address and we will ask them to create your account.
click-here-to-logout=Click here to logout
complete=Complete
complete-registration=Complete registration
configure=Configure
configured=Configured
create-an-account=Create an account
complete-external-login=Complete login on your external device\u2026
completed-link=You have successfully linked your %s account.
completed-links=You have successfully linked your %s and %s account.
confirm=Confirm
device-form-title=Device login
device-login-complete=Successfully connected device
device-title=Connect Your Device
device-link-count-exceeded-next-step=To continue, click the button below. You will be logged out and then redirected here to continue the device login.
device-link-count-exceeded-pending-logout=You are logged in as %s. No additional links may be made to %s.
disable=Disable
done=Done
dont-have-an-account=Don't have an account?
edit=Edit
email-verification-complete=Thank you. Your email has been verified.
email-verification-complete-title=Email verification complete
email-verification-form=Complete the form to request a new verification email.
email-verification-form-title=Email verification
email-verification-sent=We have sent an email to %s with your verification code. Follow the instructions in the email to verify your email address.
email-verification-sent-title=Verification sent
email-verification-required-title=Verification required
email-verification-required-send-another=Send me another email
enabled=Enabled
forgot-password=Forgot your password? Type in your email address in the form below to reset your password.
forgot-password-email-sent=We've sent you an email containing a link that will allow you to reset your password. Once you receive the email follow the instructions to change your password.
forgot-password-email-sent-title=Email sent
forgot-password-title=Forgot password
forgot-your-password=Forgot your password?
help=Help
instructions=Instructions
ip-address=IP address
link-to-existing-user=Link to an existing user
link-to-new-user=Create a new user
device-logged-in-as-not-you=You are logged in as %s. If you continue, the device login will be completed without an additional prompt. If this is not you, click logout before continuing.
link-count-exceeded-next-step=To continue, click the button below. You will be logged out and then redirected here to link to an existing user or create a new user.
link-count-exceeded-next-step-no-registration=To continue, click the button below. You will be logged out and then redirect here to link to an existing user.
link-count-exceeded-pending-logout=You have already linked to %s and no additional links are allowed.
logged-in-as=You are logged in as %s.
login=Login
login-cancel-link=Or, cancel the link request.
logout=Logout
logout-and-continue=Logout and continue\u2026
logging-out=Logging out\u2026
logout-title=Logging out
multi-factor-configuration=Multi-Factor configuration
next=Next
not-configured=Not configured
note=Note:
or=Or
parent-notified=We've sent an email to your parent. They can set up an account for you once they receive it.
parent-notified-title=Parent notified
password-alpha-constraint=Must contain at least one non-alphanumeric character
password-case-constraint=Must contain both upper and lower case characters
password-change-title=Update your password
password-changed=Your password has been updated successfully.
password-changed-title=Password updated
password-constraints-intro=Password must meet the following constraints:
password-length-constraint=Must be between %s and %s characters in length
password-number-constraint=Must contain at least one number
password-previous-constraint=Must not match the previous %s passwords
passwordless-login=Passwordless login
passwordless-button-text=Login with a magic link
pending-link-info=You have successfully authenticated using %s.
pending-link-next-step=To complete this request you may link to an existing user or create a new user.
pending-link-next-step-no-registration=To complete this request you must link to an existing user.
pending-link-login-to-complete=Login to complete your link to %s.
pending-links-login-to-complete=Login to complete your link to %s and %s.
pending-device-link=Continue to complete your link to %s.
pending-device-links=Continue to complete your link to %s and %s.
pending-link-register-to-complete=Register to complete your link to %s.
pending-links-register-to-complete=Register to complete your link to %s and %s.
profile=User Profile
provide-parent-email=Provide parent email
register-cancel-link=Or, cancel the link request.
registration-verification-complete=Thank you. Your registration has been verified.
registration-verification-complete-title=Registration verification complete
registration-verification-form=Complete the form to request a new verification email.
registration-verification-form-title=Registration verification
registration-verification-sent=We have sent an email to %s with your verification code. Follow the instructions in the email to verify your registration address.
registration-verification-sent-title=Verification sent
registration-verification-required-title=Verification required
registration-verification-required-send-another=Send me another email
return-to-login=Return to login
send-another-code=Send another code
send-code-to-phone=Send a code to your mobile phone
set-up=Set up
sms=SMS
sign-in-as-different-user=Sign in as a different user
start-idp-link-title=Link your account
two-factor-challenge=Authentication challenge
two-factor-challenge-options=Authentication challenge
two-factor-recovery-code=Recovery code
two-factor-select-method=Didn't receive a code? Try another option
two-factor-use-one-of-n-recover-codes=Use one of your %d recovery codes
trust-computer=Trust this computer for %s days
unauthorized=Unauthorized
unauthorized-message=You are not authorized to make this request.
unauthorized-message-blocked-ip=The owner of this website (%s) has blocked your IP address.
value=Value
wait-title=Complete login on your external device
waiting=Waiting

# Locale Specific separators, etc
#  - list separator - comma and a space
listSeparator=,\u0020
propertySeparator=:

#
# Success messages displayed at the top of the page. These are hard-coded in the FusionAuth code and the keys cannot be changed. You can
# still change the values though.
#
sent-code=Code successfully sent


#
# Labels for form fields. You can change the key names to anything you like but ensure that you don't change the name of the form fields.
#
birthDate=Birth date
code=Enter your verification or recovery code
email=Email
firstName=First name
fullName=Full name
lastName=Last name
loginId=Email
middleName=Middle name
mobilePhone=Mobile phone
password=Password
passwordConfirm=Confirm password
parentEmail=Parent's email
register=Register
register-step=Step %d of %d
remember-device=Keep me signed in
send=Send
submit=Submit
update=Update
username=Username
userCode=Enter your user code
verify=Verify

#
# Custom Registration forms. These must match the domain names.
#
registration.preferredLanguages=Languages
registration.timezone=Timezone
registration.username=Username
user.birthDate=Birthdate
user.email=Email
user.firstName=First name
user.fullName=Full name
user.imageUrl=Image URL
user.lastName=Last name
user.mobilePhone=Mobile phone
user.middleName=Middle name
user.password=Password
confirm.user.password=Confirm password
user.preferredLanguages=Languages
user.timezone=Timezone
user.username=Username

#
# Self service account management
#
cancel-go-back=Cancel and go back
change-password=Change password
disable-instructions=Disable two-factor
disable-two-factor=Disable two-factor
edit-profile=Edit profile
enable-instructions=Enable two-factor
enable-two-factor=Enable two-factor
go-back=Go back
send-one-time-code=Send a one-time code

#
# Self service two-factor configuration
#
no-two-factor-methods-configured=No methods have been configured
select-two-factor-method=Select a method
two-factor-authentication=Two-factor authentication
two-factor-method=Method

# Form input place holders
{placeholder}two-factor-code=Enter the one-time code

#
# Multi-factor configuration text
#
authenticator=Authenticator app

# Authenticator Enable / Disable
authenticator-disable-step-1=Enter the code from your authenticator app in the verification code field below to disable this two-factor method.
authenticator-enable-step-1=Open your authentication app and add your account by scanning the QR code to the right or by manually entering the Base32 encoded secret <strong>%s</strong>.
authenticator-enable-step-2=Once you have completed the first step, enter the code from your authenticator app in the verification code field below.

# Email Enable / Disable
email-disable-step-1=To disable two-factor using email, click the button to send a one-time use code to %s. Once you receive the code, enter it in the form below.
email-enable-step-1=To enable two-factor using email, enter an email address and click the button to send a one-time use code. Once you receive the code, enter it in the form below.

# SMS Enable / Disable
sms-disable-step-1=To disable two-factor using SMS, click the button to send a one-time use code to %s. Once you receive the code, enter it in the form below.
sms-enable-step-1=Two enable two-factor using SMS, enter a mobile phone and click the button to send a one-time use code. Once you receive the code, enter it in the form below.

authenticator-configuration=Authenticator configuration
verification-code=Verification code

manage-two-factor=Manage two-factor
go-back-to-send=Go back to send

#
# Multi-factor configuration descriptions
#
{description}two-factor-authentication=Two-factor authentication adds an additional layer of security to your account by requiring more than just a password to login. Configure one or more methods to utilize during login.
{description}two-factor-methods-selection=A second step is required to complete sign in. Select one of the following methods to complete login.
{description}two-factor-recovery-code-note=If you no longer have access to the device or application to obtain a verification code, you may use a recovery code to disable this two-factor method. Warning, when you use a recovery code to disable any two-factor method, all two-factor methods will be removed and all of your recovery codes will be cleared.
{description}recovery-codes-1=Because this is the first time you have enabled two-factor, we have generated you %d recovery codes. These codes will not be shown again, so record them right now and store them in a safe place. These codes can be used to complete a two-factor login if you lose your device, and they can be used to disable two-factor authentication as well.
{description}recovery-codes-2=Once you have recorded the codes, click Done to return to two-factor management.

{description}email-verification-required-change-email=Confirm your email address is correct and update it if you mis-typed it during registration. Updating your address will also send you a new email to the new address.
{description}email-verification-required=You must verify your email address before you continue.
{description}email-verification-required-non-interactive=Email verification is configured to be completed outside of this request. Once you have verified your email, retry this request.

{description}registration-verification-required=You must verify your registration before you continue.
{description}registration-verification-required-non-interactive=Registration verification is configured to be completed outside of this request. Once you have verified your registration, retry this request.

#
# Custom Self Service User form sections.
#
# - Names are optional, and if not provided they will be labeled 'Section 1', 'Section 2', etc.
# - The first section label will be omitted unless you specify a named label below. For your convenience, these
#   sections are configured below and commented out as 'Optionally name me!'.
#
# - By default, all section labels will be used for all tenants and all applications that are using this theme.
#
# - If you want a section title that is specific to a tenant in a user form, you may optionally prefix the key with the Tenant Id.
#
#   For example, if the tenant Id is equal to: cbeaf8fe-f4a7-4a27-9f77-c609f1b01856
#
#   [cbeaf8fe-f4a7-4a27-9f77-c609f1b01856]{self-service-form}2=Tenant specific label for section 2
#

# {self-service-form}1=Optionally name me!
# {self-service-form}2=

#
# Custom Admin User and Registration form sections.
#
# - Names are optional, and if not provided they will be labeled 'Section 1', 'Section 2', etc.
# - The first section label on the User and and Registration form in the admin UI will be omitted unless
#   you specify a named label below. For your convenience, these sections are configured below and commented out as 'Optionally name me!'.
#
# - By default, all section labels will be used for all tenants, and all applications respectively.
#
# - If you want a section title that is specific to a tenant in a user form, you may optionally prefix the key with the Tenant Id.
#
#   For example, if the tenant Id is equal to: cbeaf8fe-f4a7-4a27-9f77-c609f1b01856
#
#   [cbeaf8fe-f4a7-4a27-9f77-c609f1b01856]{user-form-section}2=Tenant specific label for section 2
#
# - If you want a section title that is specific to an Application in a registration form, you may optionally prefix the key with the Application Id.
#
#   For example, if the application Id is equal to: de2f91c7-c27a-4ad6-8be2-cfb36996cc89
#
#   [de2f91c7-c27a-4ad6-8be2-cfb36996cc89]{registration-form-section}2=Application specific label for section 2

# {user-form-section}1=Optionally name me!
{user-form-section}2=Options

# {registration-form-section}1=Optionally name me!
{registration-form-section}2=Options

#
# Custom Admin User and Registration tooltips
#
{tooltip}registration.preferredLanguages=Select one or more preferred languages
{tooltip}user.preferredLanguages=Select one or more preferred languages

#
# Custom Registration form validation errors.
#
[confirm]user.password=Confirm password

#
# Default validation errors. Add custom messages by adding field messages.
# For example, to provide a custom message for a string field named user.data.companyName, add the
# following message key: [blank]user.data.companyName=Company name is required
#
[blank]=Required
[blocked]=Not allowed
[confirm]=Confirm
[configured]=Already configured
[couldNotConvert]=Invalid
[doNotMatch]=Values do not match
[duplicate]=Already exists
[empty]=Required
[inUse]=In use
[invalid]=Invalid
[missing]=Required
[mismatch]=Unexpected value
[notEmail]=Invalid email
[tooLong]=Too long
[tooShort]=Too short
[type]=Invalid type

#
# Tooltips. You can change the key names and values to anything you like.
#
{tooltip}remember-device=Check this to stay signed into FusionAuth for the configured duration, do not select this on a public computer or when this device is shared with multiple users
{tooltip}trustComputer=Check this to bypass two-factor authentication for the configured duration, do not select this on a public computer or when this device is shared with multiple users


#
# Validation errors when forms are invalid. The format is [<error-code>]<field-name>. These are hard-coded in the FusionAuth code and the
# keys cannot be changed. You can still change the values though.
#
[invalid]applicationId=The provided application id is invalid.
[blank]code=Required
[invalid]code=Invalid code
[blank]email=Required
[blank]loginId=Required
[blank]methodId=Select a two-factor method
[blank]parentEmail=Required
[blank]password=Required
[blank]user_code=Required
[blank]captcha_token=Required
[invalid]captcha_token=Invalid challenge, try again
[cannotSend]method=A message cannot be sent to an authenticator
[disabled]method=Not enabled
[invalid]user_code=Invalid user code
[notEqual]password=Passwords don't match
[onlyAlpha]password=Password requires a non-alphanumeric character
[previouslyUsed]password=Password has been recently used
[requireNumber]password=Password requires at least one number
[singleCase]password=Password requires upper and lower case characters
[tooYoung]password=Password was changed too recently, try again later
[tooShort]password=Password does not meet the minimum length requirement
[tooLong]password=Password exceeds the maximum length requirement
[blank]passwordConfirm=Required
[missing]user.birthDate=Required
[couldNotConvert]user.birthDate=Invalid
[blank]user.email=Required
[blocked]user.email=Email address not allowed
[notEmail]user.email=Invalid email
[duplicate]user.email=An account already exists for that email
[inactive]user.email=An account already exists for that email but is locked. Contact the administrator for assistance
[blank]user.firstName=Required
[blank]user.fullName=Required
[blank]user.lastName=Required
[blank]user.middleName=Required
[blank]user.mobilePhone=Required
[invalid]user.mobilePhone=Invalid
[blank]user.parentEmail=Required
[blank]user.password=Required
[doNotMatch]user.password=Passwords don't match
[singleCase]user.password=Password must use upper and lowercase characters
[onlyAlpha]user.password=Password must contain a punctuation character
[requireNumber]user.password=Password must contain a number character
[tooShort]user.password=Password does not meet the minimum length requirement
[tooLong]user.password=Password exceeds the maximum length requirement
[blank]user.username=Required
[duplicate]user.username=An account already exists for that username
[inactive]user.username=An account already exists for that username but is locked. Contact the administrator for assistance
[mismatch]email=The requested email does not match where the code was sent
[mismatch]mobilePhone=The requested phone number does not match where the code was sent
[moderationRejected]registration.username=That username is not allowed. Please select a new one
[moderationRejected]user.username=That username is not allowed. Please select a new one

#
# Breached password messages
#
# - ExactMatch        The password and email or username combination was found in a breached data set.
# - SubAddressMatch   The password and email or username, or email sub-address was found in a breached data set.
# - PasswordOnly      The password was found in a breached data set.
# - CommonPassword    The password is one of the most commonly known breached passwords.
#
[breachedExactMatch]password=This password was found in the list of vulnerable passwords, and is no longer secure. Select a different password.
[breachedExactMatch]user.password=This password was found in the list of vulnerable passwords, and is no longer secure. Select a different password.
[breachedSubAddressMatch]password=This password was found in the list of vulnerable passwords, and is no longer secure. Select a different password.
[breachedSubAddressMatch]user.password=This password was found in the list of vulnerable passwords, and is no longer secure. Select a different password.
[breachedPasswordOnly]password=This password was found in the list of vulnerable passwords, and is no longer secure. Select a different password.
[breachedPasswordOnly]user.password=This password was found in the list of vulnerable passwords, and is no longer secure. Select a different password.
[breachedCommonPassword]password=This password is a commonly known vulnerable password. Select a more secure password.
[breachedCommonPassword]user.password=This password is a commonly known vulnerable password. Select a more secure password.

#
# Error messages displayed at the top of the page. These are always inside square brackets. These are hard-coded in the FusionAuth code and
# the keys cannot be changed. You can still change the values though.
#
[APIError]=An unexpected error occurred.
[AdditionalFieldsRequired]=Additional fields are required to complete your registration.
[EmailVerificationEmailUpdated]=Your email address has been updated and another email is on the way.
[EmailVerificationSent]=A verification email is on the way.
[EmailVerificationDisabled]=Email verification functionality is currently disabled. Contact your FusionAuth administrator for assistance.
[ErrorException]=An unexpected error occurred.
[ForgotPasswordDisabled]=Forgot password handling is not enabled. Please contact your system administrator for assistance.
[IdentityProviderDoesNotSupportRedirect]=This identity provider does not support this redirect workflow.
[InvalidChangePasswordId]=Your password reset code has expired or is invalid. Please retry your request.
[InvalidEmail]=FusionAuth was unable to find a user with that email address.
[InvalidIdentityProviderId]=Invalid request. Unable to handle the identity provider login. Please contact your system administrator or support for assistance.
[InvalidLogin]=Invalid login credentials.
[InvalidPasswordlessLoginId]=Your link has expired or is invalid. Please retry your request.
[InvalidVerificationId]=Sorry. The request contains an invalid or expired verification Id. You may need to request another verification to be sent.
[InvalidPendingIdPLinkId]=Your link has expired or is invalid. Please retry your login request.
[LinkCountExceeded]=You have reached the configured link limit of %d for this identity provider.
[LoginPreventedException]=Your account has been locked.
[LoginPreventedExceptionTooManyTwoFactorAttempts]=You have exceeded the number of allowed attempts. Your account has been locked.
[MissingApplicationId]=An applicationId is required and is missing from the request.
[MissingChangePasswordId]=A changePasswordId is required and is missing from the request.
[MissingEmail]=Your email address is required and is missing from the request.
[MissingEmailAddressException]=You must have an email address to utilize passwordless login.
[MissingPendingIdPLinkId]=You must first log into a 3rd party identity provider to complete an account link.
[MissingPKCECodeVerifier]=The code_verifier could not be determined, this request likely did not originate from FusionAuth. Unable to complete this login request.
[MissingVerificationId]=A verification Id was not sent in the request.
[NotFoundException]=The requested OAuth configuration is invalid.
[OAuthv1TokenMismatch]=Invalid request. The token provided on the OAuth v1 callback did not match the one sent during authorization. Unable to handle the identity provider login. Please contact your system administrator or support for assistance.
[Oauthv2Error]=An invalid request was made to the Authorize endpoint. %s
[PasswordlessRequestSent]=An email is on the way.
[PasswordChangeRequired]=You must change your password in order to continue.
[PasswordChangeReasonExpired]=Your password has expired and must be changed.
[PasswordChangeReasonBreached]=Your password was found in the list of vulnerable passwords and must be changed.
[PasswordChangeReasonValidation]=Your password does not meet password validation rules and must be changed.
[PasswordlessDisabled]=Passwordless login is not currently configured.
[PushTwoFactorFailed]=Failed to send a verification code using the configured push service.
[RegistrationVerificationSent]=A verification email is on the way.
[SSOSessionDeletedOrExpired]=You have been logged out of FusionAuth.
[TenantIdRequired]=FusionAuth is unable to determine which tenant to use for this request. Please add the tenantId to the URL as a request parameter.
[TwoFactorTimeout]=You did not complete the two-factor challenge in time. Please complete login again.
[UserAuthorizedNotRegisteredException]=Your account has not been registered for this application.
[UserExpiredException]=Your account has expired. Please contact your system administrator.
[UserLockedException]=Your account has been locked. Please contact your system administrator.
[UserUnauthenticated]=Oops. It looks like you've gotten here by accident. Please return to your application and log in to begin the authorization sequence.
[ExternalAuthenticationExpired]=Your external authentication request has expired, please re-attempt authentication.

# External authentication errors
# - Some of these errors are development time issues. But it is possible they could be shown to an end user depending upon your configuration.
[ExternalAuthenticationException]AppleIdToken=The id_token returned from Apple is invalid or cannot be verified. Unable to complete this login request.
[ExternalAuthenticationException]AppleTokenEndpoint=A request to the Apple Token API has failed. Unable to complete this login request.
[ExternalAuthenticationException]AppleUserObject=Failed to read the user details provided by Apple. Unable to complete this login request.
[ExternalAuthenticationException]EpicGamesAccount=A request to the Epic Games Account API has failed. Unable to complete this login request.
[ExternalAuthenticationException]EpicGamesToken=A request to the Epic Games Token API has failed. Unable to complete this login request.
[ExternalAuthenticationException]FacebookAccessToken=A request to the Facebook Access Token Info API has failed. Unable to complete this login request.
[ExternalAuthenticationException]FacebookMe=A request to the Facebook Me API has failed. Unable to complete this login request.
[ExternalAuthenticationException]FacebookMePicture=A request to the Facebook Picture API has failed. Unable to complete this login request.
[ExternalAuthenticationException]GoogleToken=A request to the Google Token API has failed. Unable to complete this login request.
[ExternalAuthenticationException]GoogleTokenInfo=A request to the Google Token Info API has failed. Unable to complete this login request.
[ExternalAuthenticationException]InvalidApplication=The requested application does not exist or is currently disabled. Unable to complete this login request.
[ExternalAuthenticationException]InvalidIdentityProviderId=The requested identityProviderId is invalid. Unable to complete this login request.
[ExternalAuthenticationException]LinkedInEmail=A request to the LinkedIn Email API has failed. Unable to complete this login request.
[ExternalAuthenticationException]LinkedInMe=A request to the LinkedIn Me API has failed. Unable to complete this login request.
[ExternalAuthenticationException]LinkedInToken=A request to the LinkedIn Token API has failed. Unable to complete this login request.
[ExternalAuthenticationException]MissingEmail=An email address was not provided for the user. This account cannot be used to login, unable to complete this login request.
[ExternalAuthenticationException]MissingUniqueId=A unique identifier was not provided for the user. This account cannot be used to login, unable to complete this login request.
[ExternalAuthenticationException]MissingUser=An authentication request cannot be completed because the user that started the request no longer exists. This account cannot be used to login, unable to complete this login request.
[ExternalAuthenticationException]MissingUsername=A username was not returned by the identity provider. This account cannot be used login, unable to complete this login request.
[ExternalAuthenticationException]NintendoToken=A request to the Nintendo Token API has failed. Unable to complete this login request.
[ExternalAuthenticationException]OpenIDConnectToken=A request to the OpenID Connect Token API has failed. Unable to complete this login request.
[ExternalAuthenticationException]OpenIDConnectUserinfo=A request to the OpenID Connect Userinfo API has failed. Unable to complete this login request.
[ExternalAuthenticationException]SAMLIdPInitiatedIssuerVerificationFailed=The SAML issuer failed validation. Unable to complete this login request.
[ExternalAuthenticationException]SAMLIdPInitiatedResponseSolicited=The SAML AuthNResponse contained an InResponseTo attribute. In an IdP Initiated Login this is un-expected.
[ExternalAuthenticationException]SAMLResponse=The SAML AuthnResponse object could not be parsed or verified. Unable to complete this login request.
[ExternalAuthenticationException]SAMLResponseAudienceNotBeforeVerificationFailed=The SAML audience is not yet available to be confirmed. Unable to complete this request.
[ExternalAuthenticationException]SAMLResponseAudienceNotOnOrAfterVerificationFailed=The SAML audience is no longer eligible to be confirmed. Unable to complete this request.
[ExternalAuthenticationException]SAMLResponseAudienceVerificationFailed=The SAML audience failed validation. Unable to complete this login request.
[ExternalAuthenticationException]SAMLResponseDestinationVerificationFailed=The SAML destination failed validation. Unable to complete this login request.
[ExternalAuthenticationException]SAMLResponseStatus=The SAML AuthnResponse status indicated the request has failed. Unable to complete this login request.
[ExternalAuthenticationException]SAMLResponseSubjectNoOnOrAfterVerificationFailed=The SAML subject is no longer eligible to be confirmed. Unable to complete this login request.
[ExternalAuthenticationException]SAMLResponseSubjectNotBeforeVerificationFailed=The SAML subject is not yet available to be confirmed. Unable to complete this login request.
[ExternalAuthenticationException]SAMLResponseUnexpectedOrReplayed=The SAML response has not been requested or has already been processed. Unable to complete this login request.
[ExternalAuthenticationException]SAMLResponseUnsolicited=The SAML response was unsolicited. Unable to complete this login request.
[ExternalAuthenticationException]SonyPSNToken=A request to the Sony PlayStation Network Token API has failed. Unable to complete this login request.
[ExternalAuthenticationException]SonyPSNUserInfo=A request to the Sony PlayStation Network User Info API has failed. Unable to complete this login request.
[ExternalAuthenticationException]SteamPlayerSummary=A request to the Steam Player summary API has failed. Unable to complete this login request.
[ExternalAuthenticationException]SteamToken=A request to the Steam Token API has failed. Unable to complete this login request.
[ExternalAuthenticationException]TwitchToken=A request to the Twitch Token API has failed. Unable to complete this login request.
[ExternalAuthenticationException]TwitchUserInfo=A request to the Twitch User Info API has failed. Unable to complete this login request.
[ExternalAuthenticationException]TwitterAccessToken=A request to the Twitter Access Token API has failed. Unable to complete this login request.
[ExternalAuthenticationException]TwitterCallbackUnconfirmed=The Twitter callback URL has not been confirmed. Unable to complete this login request.
[ExternalAuthenticationException]TwitterRequestToken=A request to the Twitter Request Token API has failed. Unable to complete this login request.
[ExternalAuthenticationException]TwitterVerifyCredentials=A request to Twitter Verify Credentials API has failed. Unable to complete this login request.
[ExternalAuthenticationException]UserDoesNotExistByEmail=You must first create a user with the same email address in order to complete this login request.
[ExternalAuthenticationException]UserDoesNotExistByUsername=You must first create a user with the same username in order to complete this login request.
[ExternalAuthenticationException]XboxSecurityTokenService=A request to the Xbox Security Token Service API has failed. Unable to complete this login request.
[ExternalAuthenticationException]XboxToken=A request to the Xbox Token API has failed. Unable to complete this login request.
[ExternalAuthenticationException]XboxUserInfo=A request to the Xbox User Info API has failed. Unable to complete this login request.

# OAuth token endpoint and callback errors
[TokenExchangeFailed]=An unexpected error occurred while completing your login attempt. Please attempt the request again.
[TokenExchangeException]=We were unable to complete your login attempt. Please attempt the request again.

# Webhook transaction failure
[WebhookTransactionException]=One or more webhooks returned an invalid response or were unreachable. Based on your transaction configuration, your action cannot be completed.

# Self Service
[SelfServiceFormNotConfigured]=Configuration is incomplete. The FusionAuth administrator must configure a form for this application.
[SelfServiceUserNotRegisteredException]=You are not registered for this application. Not all features will be available.
[TwoFactorAuthenticationMethodDisabled]=Two-factor authentication has been disabled
[TwoFactorAuthenticationMethodEnabled]=Two-factor authentication has been enabled
[TwoFactorSendFailed]=A request to send a one-time code for two-factor configuration code has failed.
[TwoFactorMessageSent]=A one-time use code was sent

# General messages
[UserWillBeLoggedIn]=You will be logged in after you complete this request.

[TrustTokenRequired]=Please complete this step-up authentication request to complete this request.
[TrustTokenExpired]=Your trust expired, please retry.
`, name)
}
