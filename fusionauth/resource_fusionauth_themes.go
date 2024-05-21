package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newTheme() *schema.Resource {
	return &schema.Resource{
		CreateContext: createTheme,
		ReadContext:   readTheme,
		UpdateContext: updateTheme,
		DeleteContext: deleteTheme,
		// Ordered based on the documented schema at: https://fusionauth.io/docs/v1/tech/apis/themes/#create-a-theme
		Schema: map[string]*schema.Schema{
			"source_theme_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The optional Id of an existing Theme to make a copy of. If present, the defaultMessages, localizedMessages, templates, and stylesheet from the source Theme will be copied to the new Theme.",
				ValidateFunc: validation.IsUUID,
			},
			"default_messages": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A properties file formatted String containing at least all of the message keys defined in the FusionAuth shipped messages file. Required if not copying an existing Theme.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"localized_messages": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A Map of localized versions of the messages. The key is the Locale and the value is a properties file formatted String.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A unique name for the Theme.",
			},
			"stylesheet": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A CSS stylesheet used to style the templates.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"account_edit": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /account/edit path. This page contains a form that enables authenticated users to update their profile.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"account_index": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /account path. This is the self-service account landing page. An authenticated user may use this as a starting point for operations such as updating their profile or configuring multi-factor authentication.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"account_two_factor_disable": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /account/two-factor/disable path. This page contains a form that accepts a verification code used to disable a multi-factor authentication method.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"account_two_factor_enable": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /account/two-factor/enable path. This page contains a form that accepts a verification code used to enable a multi-factor authentication method. Additionally, this page contains presentation of recovery codes when a user enables multi-factor authentication for the first time.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"account_two_factor_index": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /account/two-factor path. This page displays an authenticated user’s configured multi-factor authentication methods. Additionally, it provides links to enable and disable a method.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"confirmation_required": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /confirmation-required path. This page is displayed when a user attempts to complete an email based workflow that did not begin in the same browser. For example, if the user starts a forgot password workflow, and then opens the link in a separate browser the user will be shown this panel.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"account_webauthn_add": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /account/webauthn/add path. This page contains a form that allows a user to register a new WebAuthn passkey.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"account_webauthn_delete": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /account/webauthn/delete path. This page contains a form that allows a user to delete a WebAuthn passkey.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"account_webauthn_index": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /account/webauthn/ path. This page displays an authenticated user’s registered WebAuthn passkeys. Additionally, it provides links to delete an existing passkey and register a new passkey.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"email_complete": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /email/complete path. This page is used after a user has verified their email address by clicking the URL in the email. After FusionAuth has updated their user object to indicate that their email was verified, the browser is redirected to this page.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"email_sent": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /email/sent path. This page is used after a user has asked for the verification email to be resent. This can happen if the URL in the email expired and the user clicked it. In this case, the user can provide their email address again and FusionAuth will resend the email. After the user submits their email and FusionAuth re-sends a verification email to them, the browser is redirected to this page.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"email_verification_required": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /email/verification-required path. This page is rendered when a user is required to verify their email address prior to being allowed to proceed with login. This occurs when Unverified behavior is set to Gated in email verification settings on the Tenant.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"email_verify": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /email/verify path. This page is rendered when a user clicks the URL from the verification email and the verificationId has expired. FusionAuth expires verificationId after a period of time (which is configurable). If the user has a URL from the verification email that has expired, this page will be rendered and the error will be displayed to the user.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"helpers": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that contains all of the macros and templates used by the rest of the login Theme FreeMarker templates. This allows you to configure the general layout of your UI configuration and login theme without having to copy and paste HTML into each of the templates.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"index": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the / path. This is the root landing page. This page is available to unauthenticated users and will be displayed whenever someone navigates to the FusionAuth host’s root page. Prior to version 1.27.0, navigating to this URL would redirect to /admin and would subsequently render the FusionAuth admin login page.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_authorize": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/authorize path. This is the main login page for FusionAuth and is used for all interactive OAuth2 and OpenID Connect workflows.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_authorized_not_registered": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/authorized-not-registered path. This page is rendered when a user is not registered and the Application configuration requires registration before FusionAuth will complete the redirect.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_child_registration_not_allowed": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/child-registration-not-allowed path. This page contains a form where a child must provide their parent’s email address to ask their parent to create an account for them in a Consent workflow.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_child_registration_not_allowed_complete": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/child-registration-not-allowed-complete path. This page is rendered is rendered after a child provides their parent’s email address for parental consent in a Consent workflow.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_complete_registration": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/complete-registration path. This page contains a form that is used for users that have accounts but might be missing required fields.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_consent": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/consent path. This page contains a form for capturing a user's OAuth scope consent choices. If there are no scopes that require a prompt, the user is redirected automatically.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_device": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/device path. This page contains a form for accepting an end user’s short code for the interactive portion of the OAuth Device Authorization Grant workflow.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_device_complete": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/device-complete path. This page contains a complete message indicating the device authentication has completed.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_error": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "This page is used if the user starts or is in the middle of the OAuth workflow and any type of error occurs. This could be caused by the user messing with the URL or internally some type of information wasn’t passed between the OAuth endpoints correctly. For example, if you are federating login to an external IdP and that IdP does not properly echo the state parameter, FusionAuth’s OAuth workflow will break and this page will be displayed.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_logout": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/logout page. This page is used if the user initiates a logout. This page causes the user to be logged out of all associated applications via a front-channel mechanism before being redirected.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_passwordless": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/passwordless path. This page is rendered when the user starts the passwordless login workflow. The page renders the form where the user types in their email address.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_register": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/register path. This page is used to register or sign up the user for the application when self-service registration is enabled.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_start_idp_link": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/start-idp-link path. This page is used if the Identity Provider is configured to have a pending link. The user is presented with the option to link their account with an existing FusionAuth user account.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_two_factor": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/two-factor path. This page is used if the user has two-factor authentication enabled and they need to type in their code again. FusionAuth will properly handle the processing on the back end. This page contains the form that the user will put their code into.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_two_factor_methods": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/two-factor-methods path. This page contains a form providing a user with their configured multi-factor authentication options that they may use to complete the authentication challenge.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_two_factor_enable": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that contains the OAuth2 two-factor enable form.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_two_factor_enable_complete": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that contains the OAuth2 two-factor enable complete form.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_wait": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/wait path. This page is rendered when FusionAuth is waiting for an external provider to complete an out of band authentication request. For example, during a HYPR login this page will be displayed until the user completes authentication.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_webauthn": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/webauthn path. This page contains a form where a user can enter their loginId (username or email address) to authenticate with one of their registered WebAuthn passkeys. This page uses the WebAuthn bootstrap workflow.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_webauthn_reauth": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/webauthn-reauth path. This page contains a form that lists the WebAuthn passkeys currently available for re-authentication. A user can select one of the listed passkeys to authenticate using the corresponding passkey and user account.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"oauth2_webauthn_reauth_enable": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/webauthn-reauth-enable path. This page contains two forms. One allows the user to select one of their existing WebAuthn passkeys to use for re-authentication. The other allows the user to register a new WebAuthn passkey for re-authentication.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"password_change": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /password/change path. This page is used if the user is required to change their password or if they have requested a password reset. This page contains the form that allows the user to provide a new password.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"password_complete": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /password/complete path. This page is used after the user has successfully updated their password, or reset it. This page should instruct the user that their password was updated and that they need to login again.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"password_forgot": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /password/forgot path. This page is used when a user starts the forgot password workflow. This page renders the form where the user types in their email address.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"password_sent": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /password/sent path. This page is used when a user has submitted the forgot password form with their email. FusionAuth does not indicate back to the user if their email address was valid in order to prevent malicious activity that could reveal valid email addresses. Therefore, this page should indicate to the user that if their email was valid, they will receive an email shortly with a link to reset their password.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"registration_complete": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /registration/complete path. This page is used after a user has verified their email address for a specific application (i.e. a user registration) by clicking the URL in the email. After FusionAuth has updated their registration object to indicate that their email was verified, the browser is redirected to this page.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"registration_sent": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /registration/sent path. This page is used after a user has asked for the application specific verification email to be resent. This can happen if the URL in the email expired and the user clicked it. In this case, the user can provide their email address again and FusionAuth will resend the email. After the user submits their email and FusionAuth re-sends a verification email to them, the browser is redirected to this page.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"registration_verification_required": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /registration/verification-required path. This page is rendered when a user is required to verify their registration prior to being allowed to proceed with the registration flow. This occurs when Unverified behavior is set to Gated in registration verification settings on the Application.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"registration_verify": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /registration/verify path. This page is used when a user clicks the URL from the application specific verification email and the verificationId has expired. FusionAuth expires verificationId after a period of time (which is configurable). If the user has a URL from the verification email that has expired, this page will be rendered and the error will be displayed to the user.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"samlv2_logout": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /samlv2/logout path. This page is used if the user initiates a SAML logout. This page causes the user to be logged out of all associated applications via a front-channel mechanism before being redirected.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"unauthorized": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "An optional FreeMarker template that contains the unauthorized page.",
				DiffSuppressFunc: diffSuppressTemplate,
			},

			// Deprecated Theme Properties.
			"email_send": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Deprecated:       "Use email_sent instead. API endpoint has been migrated from /email/send to /email/sent.",
				Description:      "A FreeMarker template that is rendered when the user requests the /email/send page. This page is used after a user has asked for the verification email to be resent. This can happen if the URL in the email expired and the user clicked it. In this case, the user can provide their email address again and FusionAuth will resend the email. After the user submits their email and FusionAuth re-sends a verification email to them, the browser is redirected to this page.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"registration_send": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Deprecated:       "Use registration_sent instead. API endpoint has been migrated from /registration/send to /registration/sent.",
				Description:      "A FreeMarker template that is rendered when the user requests the /registration/send page. This page is used after a user has asked for the application specific verification email to be resent. This can happen if the URL in the email expired and the user clicked it. In this case, the user can provide their email address again and FusionAuth will resend the email. After the user submits their email and FusionAuth re-sends a verification email to them, the browser is redirected to this page.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func buildTheme(data *schema.ResourceData) fusionauth.Theme {
	t := fusionauth.Theme{
		DefaultMessages: data.Get("default_messages").(string),
		// LocalizedMessages: data.Get("localized_messages").(map[string]string),
		Name:       data.Get("name").(string),
		Stylesheet: data.Get("stylesheet").(string),
		Templates: fusionauth.Templates{
			AccountEdit:                       data.Get("account_edit").(string),
			AccountIndex:                      data.Get("account_index").(string),
			AccountTwoFactorDisable:           data.Get("account_two_factor_disable").(string),
			AccountTwoFactorEnable:            data.Get("account_two_factor_enable").(string),
			AccountTwoFactorIndex:             data.Get("account_two_factor_index").(string),
			ConfirmationRequired:              data.Get("confirmation_required").(string),
			AccountWebAuthnAdd:                data.Get("account_webauthn_add").(string),
			AccountWebAuthnDelete:             data.Get("account_webauthn_delete").(string),
			AccountWebAuthnIndex:              data.Get("account_webauthn_index").(string),
			EmailComplete:                     data.Get("email_complete").(string),
			EmailSend:                         data.Get("email_send").(string),
			EmailSent:                         data.Get("email_sent").(string),
			EmailVerificationRequired:         data.Get("email_verification_required").(string),
			EmailVerify:                       data.Get("email_verify").(string),
			Helpers:                           data.Get("helpers").(string),
			Index:                             data.Get("index").(string),
			Oauth2Authorize:                   data.Get("oauth2_authorize").(string),
			Oauth2AuthorizedNotRegistered:     data.Get("oauth2_authorized_not_registered").(string),
			Oauth2ChildRegistrationNotAllowed: data.Get("oauth2_child_registration_not_allowed").(string),
			Oauth2ChildRegistrationNotAllowedComplete: data.Get("oauth2_child_registration_not_allowed_complete").(string),
			Oauth2CompleteRegistration:                data.Get("oauth2_complete_registration").(string),
			Oauth2Consent:                             data.Get("oauth2_consent").(string),
			Oauth2Device:                              data.Get("oauth2_device").(string),
			Oauth2DeviceComplete:                      data.Get("oauth2_device_complete").(string),
			Oauth2Error:                               data.Get("oauth2_error").(string),
			Oauth2Logout:                              data.Get("oauth2_logout").(string),
			Oauth2Passwordless:                        data.Get("oauth2_passwordless").(string),
			Oauth2Register:                            data.Get("oauth2_register").(string),
			Oauth2StartIdPLink:                        data.Get("oauth2_start_idp_link").(string),
			Oauth2TwoFactor:                           data.Get("oauth2_two_factor").(string),
			Oauth2TwoFactorMethods:                    data.Get("oauth2_two_factor_methods").(string),
			Oauth2TwoFactorEnable:                     data.Get("oauth2_two_factor_enable").(string),
			Oauth2TwoFactorEnableComplete:             data.Get("oauth2_two_factor_enable_complete").(string),
			Oauth2Wait:                                data.Get("oauth2_wait").(string),
			Oauth2WebAuthn:                            data.Get("oauth2_webauthn").(string),
			Oauth2WebAuthnReauth:                      data.Get("oauth2_webauthn_reauth").(string),
			Oauth2WebAuthnReauthEnable:                data.Get("oauth2_webauthn_reauth_enable").(string),
			PasswordChange:                            data.Get("password_change").(string),
			PasswordComplete:                          data.Get("password_complete").(string),
			PasswordForgot:                            data.Get("password_forgot").(string),
			PasswordSent:                              data.Get("password_sent").(string),
			RegistrationComplete:                      data.Get("registration_complete").(string),
			RegistrationSend:                          data.Get("registration_send").(string),
			RegistrationSent:                          data.Get("registration_sent").(string),
			RegistrationVerificationRequired:          data.Get("registration_verification_required").(string),
			RegistrationVerify:                        data.Get("registration_verify").(string),
			Samlv2Logout:                              data.Get("samlv2_logout").(string),
			Unauthorized:                              data.Get("unauthorized").(string),
		},
	}

	if i, ok := data.GetOk("localized_messages"); ok {
		t.LocalizedMessages = intMapToStringMap(i.(map[string]interface{}))
	}

	return t
}

func createTheme(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	req := fusionauth.ThemeRequest{
		Theme: buildTheme(data),
	}

	if srcTheme, ok := data.GetOk("source_theme_id"); ok {
		req.SourceThemeId = srcTheme.(string)
	}

	resp, faErrs, err := client.FAClient.CreateTheme("", req)

	if err != nil {
		return diag.Errorf("CreateTheme err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.Theme.Id)
	return buildResourceDataFromTheme(resp.Theme, data)
}

func readTheme(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveTheme(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	t := resp.Theme

	return buildResourceDataFromTheme(t, data)
}

func updateTheme(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	req := fusionauth.ThemeRequest{
		Theme: buildTheme(data),
	}

	if srcTheme, ok := data.GetOk("source_theme_id"); ok {
		req.SourceThemeId = srcTheme.(string)
	}

	resp, faErrs, err := client.FAClient.UpdateTheme(data.Id(), req)
	if err != nil {
		return diag.Errorf("UpdateTheme err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.Theme.Id)

	return buildResourceDataFromTheme(resp.Theme, data)
}

func deleteTheme(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteTheme(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func buildResourceDataFromTheme(t fusionauth.Theme, data *schema.ResourceData) diag.Diagnostics { //nolint:gocognit,gocyclo
	if err := data.Set("default_messages", t.DefaultMessages); err != nil {
		return diag.Errorf("theme.default_messages: %s", err.Error())
	}
	if err := data.Set("localized_messages", t.LocalizedMessages); err != nil {
		return diag.Errorf("theme.localized_messages: %s", err.Error())
	}
	if err := data.Set("name", t.Name); err != nil {
		return diag.Errorf("theme.name: %s", err.Error())
	}
	if err := data.Set("stylesheet", t.Stylesheet); err != nil {
		return diag.Errorf("theme.stylesheet: %s", err.Error())
	}
	if err := data.Set("account_edit", t.Templates.AccountEdit); err != nil {
		return diag.Errorf("theme.account_edit: %s", err.Error())
	}
	if err := data.Set("account_index", t.Templates.AccountIndex); err != nil {
		return diag.Errorf("theme.account_index: %s", err.Error())
	}
	if err := data.Set("account_two_factor_disable", t.Templates.AccountTwoFactorDisable); err != nil {
		return diag.Errorf("theme.account_two_factor_disable: %s", err.Error())
	}
	if err := data.Set("account_two_factor_enable", t.Templates.AccountTwoFactorEnable); err != nil {
		return diag.Errorf("theme.account_two_factor_enable: %s", err.Error())
	}
	if err := data.Set("account_two_factor_index", t.Templates.AccountTwoFactorIndex); err != nil {
		return diag.Errorf("theme.account_two_factor_index: %s", err.Error())
	}
	if err := data.Set("confirmation_required", t.Templates.ConfirmationRequired); err != nil {
		return diag.Errorf("theme.confirmation_required: %s", err.Error())
	}
	if err := data.Set("account_webauthn_add", t.Templates.AccountWebAuthnAdd); err != nil {
		return diag.Errorf("theme.account_webauthn_add: %s", err.Error())
	}
	if err := data.Set("account_webauthn_delete", t.Templates.AccountWebAuthnDelete); err != nil {
		return diag.Errorf("theme.account_webauthn_delete: %s", err.Error())
	}
	if err := data.Set("account_webauthn_index", t.Templates.AccountWebAuthnIndex); err != nil {
		return diag.Errorf("theme.account_webauthn_index: %s", err.Error())
	}
	if err := data.Set("email_complete", t.Templates.EmailComplete); err != nil {
		return diag.Errorf("theme.email_complete: %s", err.Error())
	}
	if err := data.Set("email_send", t.Templates.EmailSend); err != nil {
		return diag.Errorf("theme.email_send: %s", err.Error())
	}
	if err := data.Set("email_verify", t.Templates.EmailVerify); err != nil {
		return diag.Errorf("theme.email_verify: %s", err.Error())
	}
	if err := data.Set("helpers", t.Templates.Helpers); err != nil {
		return diag.Errorf("theme.helpers: %s", err.Error())
	}
	if err := data.Set("index", t.Templates.Index); err != nil {
		return diag.Errorf("theme.index: %s", err.Error())
	}
	if err := data.Set("oauth2_authorize", t.Templates.Oauth2Authorize); err != nil {
		return diag.Errorf("theme.oauth2_authorize: %s", err.Error())
	}
	if err := data.Set("oauth2_child_registration_not_allowed", t.Templates.Oauth2ChildRegistrationNotAllowed); err != nil {
		return diag.Errorf("theme.oauth2_child_registration_not_allowed: %s", err.Error())
	}
	if err := data.Set("oauth2_child_registration_not_allowed_complete", t.Templates.Oauth2ChildRegistrationNotAllowedComplete); err != nil {
		return diag.Errorf("theme.oauth2_child_registration_not_allowed_complete: %s", err.Error())
	}
	if err := data.Set("oauth2_complete_registration", t.Templates.Oauth2CompleteRegistration); err != nil {
		return diag.Errorf("theme.oauth2_complete_registration: %s", err.Error())
	}
	if err := data.Set("oauth2_error", t.Templates.Oauth2Error); err != nil {
		return diag.Errorf("theme.oauth2_error: %s", err.Error())
	}
	if err := data.Set("oauth2_logout", t.Templates.Oauth2Logout); err != nil {
		return diag.Errorf("theme.oauth2_logout: %s", err.Error())
	}
	if err := data.Set("oauth2_two_factor", t.Templates.Oauth2TwoFactor); err != nil {
		return diag.Errorf("theme.oauth2_two_factor: %s", err.Error())
	}
	if err := data.Set("oauth2_two_factor_methods", t.Templates.Oauth2TwoFactorMethods); err != nil {
		return diag.Errorf("theme.oauth2_two_factor_methods: %s", err.Error())
	}
	if err := data.Set("oauth2_two_factor_enable", t.Templates.Oauth2TwoFactorEnable); err != nil {
		return diag.Errorf("theme.oauth2_two_factor_enable: %s", err.Error())
	}
	if err := data.Set("oauth2_two_factor_enable_complete", t.Templates.Oauth2TwoFactorEnableComplete); err != nil {
		return diag.Errorf("theme.oauth2_two_factor_enable_complete: %s", err.Error())
	}
	if err := data.Set("oauth2_register", t.Templates.Oauth2Register); err != nil {
		return diag.Errorf("theme.oauth2_register: %s", err.Error())
	}
	if err := data.Set("oauth2_consent", t.Templates.Oauth2Consent); err != nil {
		return diag.Errorf("theme.oauth2_consent: %s", err.Error())
	}
	if err := data.Set("oauth2_device", t.Templates.Oauth2Device); err != nil {
		return diag.Errorf("theme.oauth2_device: %s", err.Error())
	}
	if err := data.Set("oauth2_device_complete", t.Templates.Oauth2DeviceComplete); err != nil {
		return diag.Errorf("theme.oauth2_device_complete: %s", err.Error())
	}
	if err := data.Set("oauth2_passwordless", t.Templates.Oauth2Passwordless); err != nil {
		return diag.Errorf("theme.oauth2_passwordless: %s", err.Error())
	}
	if err := data.Set("oauth2_wait", t.Templates.Oauth2Wait); err != nil {
		return diag.Errorf("theme.oauth2_wait: %s", err.Error())
	}
	if err := data.Set("oauth2_webauthn", t.Templates.Oauth2WebAuthn); err != nil {
		return diag.Errorf("theme.oauth2_webauthn: %s", err.Error())
	}
	if err := data.Set("oauth2_webauthn_reauth", t.Templates.Oauth2WebAuthnReauth); err != nil {
		return diag.Errorf("theme.oauth2_webauthn_reauth: %s", err.Error())
	}
	if err := data.Set("oauth2_webauthn_reauth_enable", t.Templates.Oauth2WebAuthnReauthEnable); err != nil {
		return diag.Errorf("theme.oauth2_webauthn_reauth_enable: %s", err.Error())
	}
	if err := data.Set("password_change", t.Templates.PasswordChange); err != nil {
		return diag.Errorf("theme.password_change: %s", err.Error())
	}
	if err := data.Set("password_complete", t.Templates.PasswordComplete); err != nil {
		return diag.Errorf("theme.password_complete: %s", err.Error())
	}
	if err := data.Set("password_forgot", t.Templates.PasswordForgot); err != nil {
		return diag.Errorf("theme.password_forgot: %s", err.Error())
	}
	if err := data.Set("password_sent", t.Templates.PasswordSent); err != nil {
		return diag.Errorf("theme.password_sent: %s", err.Error())
	}
	if err := data.Set("registration_complete", t.Templates.RegistrationComplete); err != nil {
		return diag.Errorf("theme.registration_complete: %s", err.Error())
	}
	if err := data.Set("registration_send", t.Templates.RegistrationSend); err != nil {
		return diag.Errorf("theme.registration_send: %s", err.Error())
	}
	if err := data.Set("registration_verify", t.Templates.RegistrationVerify); err != nil {
		return diag.Errorf("theme.registration_verify: %s", err.Error())
	}
	if err := data.Set("samlv2_logout", t.Templates.Samlv2Logout); err != nil {
		return diag.Errorf("theme.samlv2_logout: %s", err.Error())
	}

	if err := data.Set("email_sent", t.Templates.EmailSent); err != nil {
		return diag.Errorf("theme.email_sent: %s", err.Error())
	}
	if err := data.Set("email_verification_required", t.Templates.EmailVerificationRequired); err != nil {
		return diag.Errorf("theme.email_verification_required: %s", err.Error())
	}
	if err := data.Set("registration_sent", t.Templates.RegistrationSent); err != nil {
		return diag.Errorf("theme.registration_sent: %s", err.Error())
	}
	if err := data.Set("oauth2_authorized_not_registered", t.Templates.Oauth2AuthorizedNotRegistered); err != nil {
		return diag.Errorf("theme.oauth2_authorized_not_registered: %s", err.Error())
	}
	if err := data.Set("oauth2_start_idp_link", t.Templates.Oauth2StartIdPLink); err != nil {
		return diag.Errorf("theme.oauth2_start_idp_link: %s", err.Error())
	}
	if err := data.Set("registration_verification_required", t.Templates.RegistrationVerificationRequired); err != nil {
		return diag.Errorf("theme.registration_verification_required: %s", err.Error())
	}
	if err := data.Set("unauthorized", t.Templates.Unauthorized); err != nil {
		return diag.Errorf("theme.unauthorized: %s", err.Error())
	}

	return nil
}
