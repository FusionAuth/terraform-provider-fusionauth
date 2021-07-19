package fusionauth

import (
	"fmt"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func newTheme() *schema.Resource {
	return &schema.Resource{
		Create: createTheme,
		Read:   readTheme,
		Update: updateTheme,
		Delete: deleteTheme,
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
				DiffSuppressFunc: templateCompare,
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
				DiffSuppressFunc: templateCompare,
			},
			"account_edit": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that contains the account edit page.",
				DiffSuppressFunc: templateCompare,
			},
			"account_index": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that contains the account index page, this is the account landing page.",
				DiffSuppressFunc: templateCompare,
			},
			"account_two_factor_disable": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that contains the two factor disable form.",
				DiffSuppressFunc: templateCompare,
			},
			"account_two_factor_enable": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that contains the two factor enable form.",
				DiffSuppressFunc: templateCompare,
			},
			"account_two_factor_index": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that contains the two factor index page.",
				DiffSuppressFunc: templateCompare,
			},
			"email_complete": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /email/complete page. This page is used after a user has verified their email address by clicking the URL in the email. After FusionAuth has updated their user object to indicate that their email was verified, the browser is redirected to this page.",
				DiffSuppressFunc: templateCompare,
			},
			"email_send": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /email/send page. This page is used after a user has asked for the verification email to be resent. This can happen if the URL in the email expired and the user clicked it. In this case, the user can provide their email address again and FusionAuth will resend the email. After the user submits their email and FusionAuth re-sends a verification email to them, the browser is redirected to this page.",
				DiffSuppressFunc: templateCompare,
			},
			"email_verify": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /email/verify page by clicking the URL from the verification email and the verificationId has expired. FusionAuth expires verificationId after a period of time (which is configurable). If the user has a URL from the verification email that has expired, this page will be rendered and the error will be displayed to the user.",
				DiffSuppressFunc: templateCompare,
			},
			"helpers": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that contains all of the macros and templates used by the rest of the loginTheme FreeMarker templates (i.e. oauth2Authorize). This allows you to configure the general layout of your UI configuration and login theme without having to copy and paste HTML into each of the templates.",
				DiffSuppressFunc: templateCompare,
			},
			"index": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that contains the main index page, this is the root landing page.",
				DiffSuppressFunc: templateCompare,
			},
			"oauth2_authorize": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/authorize page. This is the main login page for FusionAuth and is used for all interactive OAuth and OpenId Connect workflows.",
				DiffSuppressFunc: templateCompare,
			},
			"oauth2_child_registration_not_allowed": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/child-registration-not-allowed page. This is where the child must provide their parent’s email address to ask their parent to create an account for them.",
				DiffSuppressFunc: templateCompare,
			},
			"oauth2_child_registration_not_allowed_complete": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/child-registration-not-allowed-complete page. This is where the browser is taken after the child provides their parent’s email address.",
				DiffSuppressFunc: templateCompare,
			},
			"oauth2_complete_registration": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/complete-registration page. This page is used for users that have accounts but might be missing required fields.",
				DiffSuppressFunc: templateCompare,
			},
			"oauth2_error": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/error page. This page is used if the user starts or is in the middle of the OAuth workflow and any type of error occurs. This could be caused by the user messing with the URL or internally some type of information wasn’t passed between the OAuth endpoints correctly. For example, if you are federating login to an external IdP and that IdP does not properly echo the state parameter, FusionAuth’s OAuth workflow will break and this page will be displayed.",
				DiffSuppressFunc: templateCompare,
			},
			"oauth2_logout": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/logout page. This page is used if the user initiates a logout. This page causes the user to be logged out of all associated applications via a front-channel mechanism before being redirected.",
				DiffSuppressFunc: templateCompare,
			},
			"oauth2_register": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/register page. This page is used for users that need to register (sign-up)",
				DiffSuppressFunc: templateCompare,
			},
			"oauth2_device": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template",
				DiffSuppressFunc: templateCompare,
			},
			"oauth2_device_complete": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template",
				DiffSuppressFunc: templateCompare,
			},
			"oauth2_passwordless": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template",
				DiffSuppressFunc: templateCompare,
			},
			"oauth2_two_factor_methods": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that contains the OAuth2 two-factor option form.",
				DiffSuppressFunc: templateCompare,
			},
			"oauth2_wait": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template",
				DiffSuppressFunc: templateCompare,
			},
			"oauth2_two_factor": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /oauth2/two-factor page. This page is used if the user has two-factor authentication enabled and they need to type in their code again. FusionAuth will properly handle the SMS or authenticator app processing on the back end. This page contains the form that the user will put their code into.",
				DiffSuppressFunc: templateCompare,
			},
			"password_change": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /password/change page. This page is used if the user is required to change their password or if they have requested a password reset. This page contains the form that allows the user to provide a new password.",
				DiffSuppressFunc: templateCompare,
			},
			"password_complete": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /password/complete page. This page is used after the user has successfully updated their password (or reset it). This page should instruct the user that their password was updated and that they need to login again.",
				DiffSuppressFunc: templateCompare,
			},
			"password_forgot": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /password/forgot page. This page is used when a user starts the forgot password workflow. This page renders the form where the user types in their email address.",
				DiffSuppressFunc: templateCompare,
			},
			"password_sent": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /password/sent page. This page is used when a user has submitted the forgot password form with their email. FusionAuth does not indicate back to the user if their email address was valid in order to prevent malicious activity that could reveal valid email addresses. Therefore, this page should indicate to the user that if their email was valid, they will receive an email shortly with a link to reset their password.",
				DiffSuppressFunc: templateCompare,
			},
			"registration_complete": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /registration/complete page. This page is used after a user has verified their email address for a specific application (i.e. a user registration) by clicking the URL in the email. After FusionAuth has updated their registration object to indicate that their email was verified, the browser is redirected to this page.",
				DiffSuppressFunc: templateCompare,
			},
			"registration_send": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /registration/send page. This page is used after a user has asked for the application specific verification email to be resent. This can happen if the URL in the email expired and the user clicked it. In this case, the user can provide their email address again and FusionAuth will resend the email. After the user submits their email and FusionAuth re-sends a verification email to them, the browser is redirected to this page.",
				DiffSuppressFunc: templateCompare,
			},
			"registration_verify": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template that is rendered when the user requests the /registration/verify page by clicking the URL from the application specific verification email and the verificationId has expired. FusionAuth expires verificationId after a period of time (which is configurable). If the user has a URL from the verification email that has expired, this page will be rendered and the error will be displayed to the user.",
				DiffSuppressFunc: templateCompare,
			},
			"samlv2_logout": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "This page is used if the user initiates a SAML logout. This page causes the user to be logged out of all associated applications via a front-channel mechanism before being redirected.",
				DiffSuppressFunc: templateCompare,
			},
			"email_sent": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template",
				DiffSuppressFunc: templateCompare,
			},
			"email_verification_required": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template",
				DiffSuppressFunc: templateCompare,
			},
			"oauth2_authorized_not_registered": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template",
				DiffSuppressFunc: templateCompare,
			},
			"oauth2_start_idp_link": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template",
				DiffSuppressFunc: templateCompare,
			},
			"registration_sent": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template",
				DiffSuppressFunc: templateCompare,
			},
			"registration_verification_required": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "A FreeMarker template",
				DiffSuppressFunc: templateCompare,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			EmailComplete:                     data.Get("email_complete").(string),
			EmailSend:                         data.Get("email_send").(string),
			EmailVerify:                       data.Get("email_verify").(string),
			Helpers:                           data.Get("helpers").(string),
			Index:                             data.Get("index").(string),
			Oauth2Authorize:                   data.Get("oauth2_authorize").(string),
			Oauth2ChildRegistrationNotAllowed: data.Get("oauth2_child_registration_not_allowed").(string),
			Oauth2ChildRegistrationNotAllowedComplete: data.Get("oauth2_child_registration_not_allowed_complete").(string),
			Oauth2CompleteRegistration:                data.Get("oauth2_complete_registration").(string),
			Oauth2Error:                               data.Get("oauth2_error").(string),
			Oauth2Logout:                              data.Get("oauth2_logout").(string),
			Oauth2TwoFactor:                           data.Get("oauth2_two_factor").(string),
			Oauth2TwoFactorMethods:                    data.Get("oauth2_two_factor_methods").(string),
			Oauth2Register:                            data.Get("oauth2_register").(string),
			Oauth2Device:                              data.Get("oauth2_device").(string),
			Oauth2DeviceComplete:                      data.Get("oauth2_device_complete").(string),
			Oauth2Passwordless:                        data.Get("oauth2_passwordless").(string),
			Oauth2Wait:                                data.Get("oauth2_wait").(string),
			PasswordChange:                            data.Get("password_change").(string),
			PasswordComplete:                          data.Get("password_complete").(string),
			PasswordForgot:                            data.Get("password_forgot").(string),
			PasswordSent:                              data.Get("password_sent").(string),
			RegistrationComplete:                      data.Get("registration_complete").(string),
			RegistrationSend:                          data.Get("registration_send").(string),
			RegistrationVerify:                        data.Get("registration_verify").(string),
			Samlv2Logout:                              data.Get("samlv2_logout").(string),
			EmailSent:                                 data.Get("email_sent").(string),
			EmailVerificationRequired:                 data.Get("email_verification_required").(string),
			RegistrationSent:                          data.Get("registration_sent").(string),
			Oauth2AuthorizedNotRegistered:             data.Get("oauth2_authorized_not_registered").(string),
			Oauth2StartIdPLink:                        data.Get("oauth2_start_idp_link").(string),
			RegistrationVerificationRequired:          data.Get("registration_verification_required").(string),
		},
	}

	m := data.Get("localized_messages").(map[string]interface{})
	t.LocalizedMessages = make(map[string]string)

	for k, v := range m {
		t.LocalizedMessages[k] = v.(string)
	}

	return t
}

func createTheme(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)

	req := fusionauth.ThemeRequest{
		Theme: buildTheme(data),
	}

	if srcTheme, ok := data.GetOk("source_theme_id"); ok {
		req.SourceThemeId = srcTheme.(string)
	}

	resp, faErrs, err := client.FAClient.CreateTheme("", req)

	if err != nil {
		return fmt.Errorf("CreateTheme err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	data.SetId(resp.Theme.Id)
	return buildResourceDataFromTheme(resp.Theme, data)
}

func readTheme(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveTheme(id)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	t := resp.Theme

	return buildResourceDataFromTheme(t, data)
}

func updateTheme(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	req := fusionauth.ThemeRequest{
		Theme: buildTheme(data),
	}

	if srcTheme, ok := data.GetOk("source_theme_id"); ok {
		req.SourceThemeId = srcTheme.(string)
	}

	resp, faErrs, err := client.FAClient.UpdateTheme(data.Id(), req)

	if err != nil {
		return fmt.Errorf("UpdateTheme err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	data.SetId(resp.Theme.Id)

	return buildResourceDataFromTheme(resp.Theme, data)
}

func deleteTheme(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteTheme(id)
	if err != nil {
		return err
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}
	return nil
}

func buildResourceDataFromTheme(t fusionauth.Theme, data *schema.ResourceData) error { //nolint:gocognit,gocyclo
	if err := data.Set("default_messages", t.DefaultMessages); err != nil {
		return fmt.Errorf("theme.default_messages: %s", err.Error())
	}
	if err := data.Set("localized_messages", t.LocalizedMessages); err != nil {
		return fmt.Errorf("theme.localized_messages: %s", err.Error())
	}
	if err := data.Set("name", t.Name); err != nil {
		return fmt.Errorf("theme.name: %s", err.Error())
	}
	if err := data.Set("stylesheet", t.Stylesheet); err != nil {
		return fmt.Errorf("theme.stylesheet: %s", err.Error())
	}
	if err := data.Set("account_edit", t.Templates.AccountEdit); err != nil {
		return fmt.Errorf("theme.account_edit: %s", err.Error())
	}
	if err := data.Set("account_index", t.Templates.AccountIndex); err != nil {
		return fmt.Errorf("theme.account_index: %s", err.Error())
	}
	if err := data.Set("account_two_factor_disable", t.Templates.AccountTwoFactorDisable); err != nil {
		return fmt.Errorf("theme.account_two_factor_disable: %s", err.Error())
	}
	if err := data.Set("account_two_factor_enable", t.Templates.AccountTwoFactorEnable); err != nil {
		return fmt.Errorf("theme.account_two_factor_enable: %s", err.Error())
	}
	if err := data.Set("account_two_factor_index", t.Templates.AccountTwoFactorIndex); err != nil {
		return fmt.Errorf("theme.account_two_factor_index: %s", err.Error())
	}
	if err := data.Set("email_complete", t.Templates.EmailComplete); err != nil {
		return fmt.Errorf("theme.email_complete: %s", err.Error())
	}
	if err := data.Set("email_send", t.Templates.EmailSend); err != nil {
		return fmt.Errorf("theme.email_send: %s", err.Error())
	}
	if err := data.Set("email_verify", t.Templates.EmailVerify); err != nil {
		return fmt.Errorf("theme.email_verify: %s", err.Error())
	}
	if err := data.Set("helpers", t.Templates.Helpers); err != nil {
		return fmt.Errorf("theme.helpers: %s", err.Error())
	}
	if err := data.Set("index", t.Templates.Index); err != nil {
		return fmt.Errorf("theme.index: %s", err.Error())
	}
	if err := data.Set("oauth2_authorize", t.Templates.Oauth2Authorize); err != nil {
		return fmt.Errorf("theme.oauth2_authorize: %s", err.Error())
	}
	if err := data.Set("oauth2_child_registration_not_allowed", t.Templates.Oauth2ChildRegistrationNotAllowed); err != nil {
		return fmt.Errorf("theme.oauth2_child_registration_not_allowed: %s", err.Error())
	}
	if err := data.Set("oauth2_child_registration_not_allowed_complete", t.Templates.Oauth2ChildRegistrationNotAllowedComplete); err != nil {
		return fmt.Errorf("theme.oauth2_child_registration_not_allowed_complete: %s", err.Error())
	}
	if err := data.Set("oauth2_complete_registration", t.Templates.Oauth2CompleteRegistration); err != nil {
		return fmt.Errorf("theme.oauth2_complete_registration: %s", err.Error())
	}
	if err := data.Set("oauth2_error", t.Templates.Oauth2Error); err != nil {
		return fmt.Errorf("theme.oauth2_error: %s", err.Error())
	}
	if err := data.Set("oauth2_logout", t.Templates.Oauth2Logout); err != nil {
		return fmt.Errorf("theme.oauth2_logout: %s", err.Error())
	}
	if err := data.Set("oauth2_two_factor", t.Templates.Oauth2TwoFactor); err != nil {
		return fmt.Errorf("theme.oauth2_two_factor: %s", err.Error())
	}
	if err := data.Set("oauth2_two_factor_methods", t.Templates.Oauth2TwoFactorMethods); err != nil {
		return fmt.Errorf("theme.oauth2_two_factor_methods: %s", err.Error())
	}
	if err := data.Set("oauth2_register", t.Templates.Oauth2Register); err != nil {
		return fmt.Errorf("theme.oauth2_register: %s", err.Error())
	}
	if err := data.Set("oauth2_device", t.Templates.Oauth2Device); err != nil {
		return fmt.Errorf("theme.oauth2_device: %s", err.Error())
	}
	if err := data.Set("oauth2_device_complete", t.Templates.Oauth2DeviceComplete); err != nil {
		return fmt.Errorf("theme.oauth2_device_complete: %s", err.Error())
	}
	if err := data.Set("oauth2_passwordless", t.Templates.Oauth2Passwordless); err != nil {
		return fmt.Errorf("theme.oauth2_passwordless: %s", err.Error())
	}
	if err := data.Set("oauth2_wait", t.Templates.Oauth2Wait); err != nil {
		return fmt.Errorf("theme.oauth2_wait: %s", err.Error())
	}
	if err := data.Set("password_change", t.Templates.PasswordChange); err != nil {
		return fmt.Errorf("theme.password_change: %s", err.Error())
	}
	if err := data.Set("password_complete", t.Templates.PasswordComplete); err != nil {
		return fmt.Errorf("theme.password_complete: %s", err.Error())
	}
	if err := data.Set("password_forgot", t.Templates.PasswordForgot); err != nil {
		return fmt.Errorf("theme.password_forgot: %s", err.Error())
	}
	if err := data.Set("password_sent", t.Templates.PasswordSent); err != nil {
		return fmt.Errorf("theme.password_sent: %s", err.Error())
	}
	if err := data.Set("registration_complete", t.Templates.RegistrationComplete); err != nil {
		return fmt.Errorf("theme.registration_complete: %s", err.Error())
	}
	if err := data.Set("registration_send", t.Templates.RegistrationSend); err != nil {
		return fmt.Errorf("theme.registration_send: %s", err.Error())
	}
	if err := data.Set("registration_verify", t.Templates.RegistrationVerify); err != nil {
		return fmt.Errorf("theme.registration_verify: %s", err.Error())
	}
	if err := data.Set("samlv2_logout", t.Templates.Samlv2Logout); err != nil {
		return fmt.Errorf("theme.samlv2_logout: %s", err.Error())
	}

	if err := data.Set("email_sent", t.Templates.EmailSent); err != nil {
		return fmt.Errorf("theme.email_sent: %s", err.Error())
	}
	if err := data.Set("email_verification_required", t.Templates.EmailVerificationRequired); err != nil {
		return fmt.Errorf("theme.email_verification_required: %s", err.Error())
	}
	if err := data.Set("registration_sent", t.Templates.RegistrationSent); err != nil {
		return fmt.Errorf("theme.registration_sent: %s", err.Error())
	}
	if err := data.Set("oauth2_authorized_not_registered", t.Templates.Oauth2AuthorizedNotRegistered); err != nil {
		return fmt.Errorf("theme.oauth2_authorized_not_registered: %s", err.Error())
	}
	if err := data.Set("oauth2_start_idp_link", t.Templates.Oauth2StartIdPLink); err != nil {
		return fmt.Errorf("theme.oauth2_start_idp_link: %s", err.Error())
	}
	if err := data.Set("registration_verification_required", t.Templates.RegistrationVerificationRequired); err != nil {
		return fmt.Errorf("theme.registration_verification_required: %s", err.Error())
	}

	return nil
}
