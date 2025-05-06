package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceTheme() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceThemeRead,
		Schema: map[string]*schema.Schema{
			// Data Source Parameters
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"theme_id", "name"},
				RequiredWith: []string{"type"},
				Description:  "The case-insensitive string to search for in the Theme name. Must be used together with 'type'.",
			},
			"theme_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ExactlyOneOf:  []string{"theme_id", "name"},
				ConflictsWith: []string{"type"},
				Description:   "The unique Id of the Theme to retrieve.",
				ValidateFunc:  validation.IsUUID,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"name"},
				Description:  "The type of Theme to retrieve. Required when searching by name. The value must be one of the following: `simple` or `advanced`.",
				ValidateFunc: validation.StringInSlice([]string{"simple", "advanced"}, false),
			},
			// Data Source Attributes
			"data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An object that can hold any information about the Theme that should be persisted.",
			},
			"default_messages": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A properties file formatted String containing at least all of the message keys defined in the FusionAuth shipped messages file. Required if not copying an existing Theme.",
			},
			"localized_messages": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A Map of localized versions of the messages. The key is the Locale and the value is a properties file formatted String.",
			},
			"stylesheet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A CSS stylesheet used to style the templates.",
			},
			"account_edit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /account/edit path. This page contains a form that enables authenticated users to update their profile.",
			},
			"account_index": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /account path. This is the self-service account landing page. An authenticated user may use this as a starting point for operations such as updating their profile or configuring multi-factor authentication.",
			},
			"account_two_factor_disable": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /account/two-factor/disable path. This page contains a form that accepts a verification code used to disable a multi-factor authentication method.",
			},
			"account_two_factor_enable": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /account/two-factor/enable path. This page contains a form that accepts a verification code used to enable a multi-factor authentication method. Additionally, this page contains presentation of recovery codes when a user enables multi-factor authentication for the first time.",
			},
			"account_two_factor_index": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /account/two-factor path. This page displays an authenticated user’s configured multi-factor authentication methods. Additionally, it provides links to enable and disable a method.",
			},
			"confirmation_required": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /confirmation-required path. This page is displayed when a user attempts to complete an email based workflow that did not begin in the same browser. For example, if the user starts a forgot password workflow, and then opens the link in a separate browser the user will be shown this panel.",
			},
			"account_webauthn_add": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /account/webauthn/add path. This page contains a form that allows a user to register a new WebAuthn passkey.",
			},
			"account_webauthn_delete": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /account/webauthn/delete path. This page contains a form that allows a user to delete a WebAuthn passkey.",
			},
			"account_webauthn_index": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /account/webauthn/ path. This page displays an authenticated user’s registered WebAuthn passkeys. Additionally, it provides links to delete an existing passkey and register a new passkey.",
			},
			"email_complete": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /email/complete path. This page is used after a user has verified their email address by clicking the URL in the email. After FusionAuth has updated their user object to indicate that their email was verified, the browser is redirected to this page.",
			},
			"email_sent": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /email/sent path. This page is used after a user has asked for the verification email to be resent. This can happen if the URL in the email expired and the user clicked it. In this case, the user can provide their email address again and FusionAuth will resend the email. After the user submits their email and FusionAuth re-sends a verification email to them, the browser is redirected to this page.",
			},
			"email_verification_required": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /email/verification-required path. This page is rendered when a user is required to verify their email address prior to being allowed to proceed with login. This occurs when Unverified behavior is set to Gated in email verification settings on the Tenant.",
			},
			"email_verify": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /email/verify path. This page is rendered when a user clicks the URL from the verification email and the verificationId has expired. FusionAuth expires verificationId after a period of time (which is configurable). If the user has a URL from the verification email that has expired, this page will be rendered and the error will be displayed to the user.",
			},
			"helpers": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that contains all of the macros and templates used by the rest of the login Theme FreeMarker templates. This allows you to configure the general layout of your UI configuration and login theme without having to copy and paste HTML into each of the templates.",
			},
			"index": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the / path. This is the root landing page. This page is available to unauthenticated users and will be displayed whenever someone navigates to the FusionAuth host’s root page. Prior to version 1.27.0, navigating to this URL would redirect to /admin and would subsequently render the FusionAuth admin login page.",
			},
			"oauth2_authorize": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/authorize path. This is the main login page for FusionAuth and is used for all interactive OAuth2 and OpenID Connect workflows.",
			},
			"oauth2_authorized_not_registered": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/authorized-not-registered path. This page is rendered when a user is not registered and the Application configuration requires registration before FusionAuth will complete the redirect.",
			},
			"oauth2_child_registration_not_allowed": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/child-registration-not-allowed path. This page contains a form where a child must provide their parent’s email address to ask their parent to create an account for them in a Consent workflow.",
			},
			"oauth2_child_registration_not_allowed_complete": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/child-registration-not-allowed-complete path. This page is rendered is rendered after a child provides their parent’s email address for parental consent in a Consent workflow.",
			},
			"oauth2_complete_registration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/complete-registration path. This page contains a form that is used for users that have accounts but might be missing required fields.",
			},
			"oauth2_consent": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/consent path. This page contains a form for capturing a user's OAuth scope consent choices. If there are no scopes that require a prompt, the user is redirected automatically.",
			},
			"oauth2_device": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/device path. This page contains a form for accepting an end user’s short code for the interactive portion of the OAuth Device Authorization Grant workflow.",
			},
			"oauth2_device_complete": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/device-complete path. This page contains a complete message indicating the device authentication has completed.",
			},
			"oauth2_error": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This page is used if the user starts or is in the middle of the OAuth workflow and any type of error occurs. This could be caused by the user messing with the URL or internally some type of information wasn’t passed between the OAuth endpoints correctly. For example, if you are federating login to an external IdP and that IdP does not properly echo the state parameter, FusionAuth’s OAuth workflow will break and this page will be displayed.",
			},
			"oauth2_logout": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/logout page. This page is used if the user initiates a logout. This page causes the user to be logged out of all associated applications via a front-channel mechanism before being redirected.",
			},
			"oauth2_passwordless": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/passwordless path. This page is rendered when the user starts the passwordless login workflow. The page renders the form where the user types in their email address.",
			},
			"oauth2_register": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/register path. This page is used to register or sign up the user for the application when self-service registration is enabled.",
			},
			"oauth2_start_idp_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/start-idp-link path. This page is used if the Identity Provider is configured to have a pending link. The user is presented with the option to link their account with an existing FusionAuth user account.",
			},
			"oauth2_two_factor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/two-factor path. This page is used if the user has two-factor authentication enabled and they need to type in their code again. FusionAuth will properly handle the processing on the back end. This page contains the form that the user will put their code into.",
			},
			"oauth2_two_factor_methods": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/two-factor-methods path. This page contains a form providing a user with their configured multi-factor authentication options that they may use to complete the authentication challenge.",
			},
			"oauth2_two_factor_enable": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that contains the OAuth2 two-factor enable form.",
			},
			"oauth2_two_factor_enable_complete": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that contains the OAuth2 two-factor enable complete form.",
			},
			"oauth2_wait": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/wait path. This page is rendered when FusionAuth is waiting for an external provider to complete an out of band authentication request. For example, during a HYPR login this page will be displayed until the user completes authentication.",
			},
			"oauth2_webauthn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/webauthn path. This page contains a form where a user can enter their loginId (username or email address) to authenticate with one of their registered WebAuthn passkeys. This page uses the WebAuthn bootstrap workflow.",
			},
			"oauth2_webauthn_reauth": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/webauthn-reauth path. This page contains a form that lists the WebAuthn passkeys currently available for re-authentication. A user can select one of the listed passkeys to authenticate using the corresponding passkey and user account.",
			},
			"oauth2_webauthn_reauth_enable": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /oauth2/webauthn-reauth-enable path. This page contains two forms. One allows the user to select one of their existing WebAuthn passkeys to use for re-authentication. The other allows the user to register a new WebAuthn passkey for re-authentication.",
			},
			"password_change": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /password/change path. This page is used if the user is required to change their password or if they have requested a password reset. This page contains the form that allows the user to provide a new password.",
			},
			"password_complete": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /password/complete path. This page is used after the user has successfully updated their password, or reset it. This page should instruct the user that their password was updated and that they need to login again.",
			},
			"password_forgot": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /password/forgot path. This page is used when a user starts the forgot password workflow. This page renders the form where the user types in their email address.",
			},
			"password_sent": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /password/sent path. This page is used when a user has submitted the forgot password form with their email. FusionAuth does not indicate back to the user if their email address was valid in order to prevent malicious activity that could reveal valid email addresses. Therefore, this page should indicate to the user that if their email was valid, they will receive an email shortly with a link to reset their password.",
			},
			"phone_complete": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /phone/complete path. This page is used after a user has verified their phone number by clicking the URL in the message. After FusionAuth has updated their user object to indicate that their phone number was verified, the browser is redirected to this page.",
			},
			"phone_sent": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /phone/sent path. This page is used after a user has asked for the verification message to be resent. This can happen if the URL in the message expired and the user clicked it. In this case, the user can provide their phone number again and FusionAuth will resend the message. After the user submits their phone number and FusionAuth re-sends a verification message to them, the browser is redirected to this page.",
			},
			"phone_verification_required": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /phone/verification-required path. This page is rendered when a user is required to verify their phone number prior to being allowed to proceed with login. This occurs when Unverified behavior is set to Gated in identities/phone verification settings on the Tenant.",
			},
			"phone_verify": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /phone/verify path. This page is rendered when a user clicks the URL from the verification message and the verificationId has expired. FusionAuth expires verificationId after a period of time (which is configurable). If the user has a URL from the verification message that has expired, this page will be rendered and the error will be displayed to the user.",
			},
			"registration_complete": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /registration/complete path. This page is used after a user has verified their email address for a specific application (i.e. a user registration) by clicking the URL in the email. After FusionAuth has updated their registration object to indicate that their email was verified, the browser is redirected to this page.",
			},
			"registration_sent": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /registration/sent path. This page is used after a user has asked for the application specific verification email to be resent. This can happen if the URL in the email expired and the user clicked it. In this case, the user can provide their email address again and FusionAuth will resend the email. After the user submits their email and FusionAuth re-sends a verification email to them, the browser is redirected to this page.",
			},
			"registration_verification_required": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /registration/verification-required path. This page is rendered when a user is required to verify their registration prior to being allowed to proceed with the registration flow. This occurs when Unverified behavior is set to Gated in registration verification settings on the Application.",
			},
			"registration_verify": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /registration/verify path. This page is used when a user clicks the URL from the application specific verification email and the verificationId has expired. FusionAuth expires verificationId after a period of time (which is configurable). If the user has a URL from the verification email that has expired, this page will be rendered and the error will be displayed to the user.",
			},
			"samlv2_logout": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A FreeMarker template that is rendered when the user requests the /samlv2/logout path. This page is used if the user initiates a SAML logout. This page causes the user to be logged out of all associated applications via a front-channel mechanism before being redirected.",
			},
			"unauthorized": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An optional FreeMarker template that contains the unauthorized page.",
			},

			// Deprecated Theme Properties.
			"email_send": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "Use email_sent instead. API endpoint has been migrated from /email/send to /email/sent.",
				Description: "A FreeMarker template that is rendered when the user requests the /email/send page. This page is used after a user has asked for the verification email to be resent. This can happen if the URL in the email expired and the user clicked it. In this case, the user can provide their email address again and FusionAuth will resend the email. After the user submits their email and FusionAuth re-sends a verification email to them, the browser is redirected to this page.",
			},
			"registration_send": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "Use registration_sent instead. API endpoint has been migrated from /registration/send to /registration/sent.",
				Description: "A FreeMarker template that is rendered when the user requests the /registration/send page. This page is used after a user has asked for the application specific verification email to be resent. This can happen if the URL in the email expired and the user clicked it. In this case, the user can provide their email address again and FusionAuth will resend the email. After the user submits their email and FusionAuth re-sends a verification email to them, the browser is redirected to this page.",
			},
		},
	}
}

func dataSourceThemeRead(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	var searchTerm string
	var theme fusionauth.Theme

	if entityID, ok := data.GetOk("theme_id"); ok {
		searchTerm = entityID.(string)
		res, err, _ := client.FAClient.RetrieveTheme(searchTerm)
		if err != nil {
			return diag.FromErr(err)
		}
		if res.StatusCode == http.StatusNotFound {
			return diag.Errorf("couldn't find theme '%s'", searchTerm)
		}
		if err := checkResponse(res.StatusCode, nil); err != nil {
			return diag.FromErr(err)
		}
		theme = res.Theme
	} else {
		searchRequest := fusionauth.ThemeSearchRequest{
			Search: fusionauth.ThemeSearchCriteria{
				Name: data.Get("name").(string),
				Type: fusionauth.ThemeType(data.Get("type").(string)),
				BaseSearchCriteria: fusionauth.BaseSearchCriteria{
					NumberOfResults: 1,
					StartRow:        0,
				},
			},
		}
		res, err, _ := client.FAClient.SearchThemes(searchRequest)
		if err != nil {
			return diag.FromErr(err)
		}
		if res.StatusCode == http.StatusNotFound || len(res.Themes) == 0 {
			return diag.Errorf("couldn't find theme with name '%s'", data.Get("name").(string))
		}
		if err := checkResponse(res.StatusCode, nil); err != nil {
			return diag.FromErr(err)
		}
		theme = res.Themes[0]
	}

	data.SetId(theme.Id)
	return buildResourceDataFromTheme(theme, data)
}
