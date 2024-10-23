package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: createWebhook,
		ReadContext:   readWebhook,
		UpdateContext: updateWebhook,
		DeleteContext: deleteWebhook,
		Schema: map[string]*schema.Schema{
			"tenant_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The Ids of the Tenants that this Webhook should be associated with. If no Ids are specified and the global field is false, this Webhook will not be used.",
			},
			"connect_timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The connection timeout in milliseconds used when FusionAuth sends events to the Webhook.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the Webhook. This is used for display purposes only.",
			},
			"events_enabled": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audit_log_create": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When an audit log is created",
						},
						"event_log_create": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When an event log is created",
						},
						"jwt_public_key_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a JWT RSA Public / Private keypair may have been changed",
						},
						"jwt_refresh": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When an access token is refreshed using a refresh token",
						},
						"jwt_refresh_token_revoke": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a JWT Refresh Token is revoked",
						},
						"kickstart_success": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When kickstart has successfully completed",
						},
						"user_action": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user action is triggered",
						},
						"user_bulk_create": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When multiple users are created in bulk (i.e. during an import)",
						},
						"user_create": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user is created",
						},
						"user_create_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user create transaction has completed",
						},
						"user_deactivate": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user is deactivated",
						},
						"user_delete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user is deleted",
						},
						"user_delete_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user delete transaction has completed",
						},
						"user_email_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user updates their email address",
						},
						"user_email_verified": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user verifies their email address",
						},
						"user_identity_provider_link": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user is linked to an identity provider",
						},
						"user_identity_provider_unlink": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a link to an identity provider is removed",
						},
						"user_login_id_duplicate_create": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a request to create a user with a login Id (email or username) which is already in use has been received",
						},
						"user_login_id_duplicate_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a request to update a user and change their login Id (email or username) to one that is already in use has been received",
						},
						"user_login_failed": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user fails a login request",
						},
						"user_login_new_device": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user begins a login request with a new device",
						},
						"user_login_success": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user completes a login request",
						},
						"user_login_suspicious": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user logs in and is considered to be a potential threat",
						},
						"user_password_breach": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When Reactor detects a user is using a potentially breached password (requires an activated license)",
						},
						"user_password_reset_send": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a forgot password email has been sent to a user",
						},
						"user_password_reset_start": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When the process to reset a user password has started",
						},
						"user_password_reset_success": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user has successfully reset their password",
						},
						"user_password_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user has updated their password",
						},
						"user_reactivate": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user is reactivated",
						},
						"user_registration_create": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user registration is created",
						},
						"user_registration_create_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user registration create transaction has completed",
						},
						"user_registration_delete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user registration is deleted",
						},
						"user_registration_delete_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user registration delete transaction has completed",
						},
						"user_registration_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user registration is updated",
						},
						"user_registration_update_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user registration update transaction has completed",
						},
						"user_registration_verified": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user completes registration verification",
						},
						"user_two_factor_method_add": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user has added a two-factor method",
						},
						"user_two_factor_method_remove": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user has removed a two-factor method",
						},
						"user_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user is updated",
						},
						"user_update_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user update transaction has completed",
						},
					},
				},
			},
			"global": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not this Webhook is used for all events or just for specific Applications.",
			},
			"headers": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that contains headers that are sent as part of the HTTP request for the events.",
			},
			"http_authentication_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The HTTP basic authentication password that is sent as part of the HTTP request for the events.",
			},
			"http_authentication_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The HTTP basic authentication username that is sent as part of the HTTP request for the events.",
			},
			"read_timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The read timeout in milliseconds used when FusionAuth sends events to the Webhook.",
			},
			"signature_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates if the Webhook request should be signed.",
							RequiredWith: []string{
								"signature_configuration.0.signing_key_id",
							},
						},
						"signing_key_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "The Id of the key used to sign the Webhook request.",
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},
			"ssl_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An SSL certificate in PEM format that is used to establish the a SSL (TLS specifically) connection to the Webhook.",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The fully qualified URL of the Webhook’s endpoint that will accept the event requests from FusionAuth.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func buildWebhook(data *schema.ResourceData) fusionauth.Webhook {
	wh := fusionauth.Webhook{
		TenantIds:      handleStringSlice("tenant_ids", data),
		ConnectTimeout: data.Get("connect_timeout").(int),
		Description:    data.Get("description").(string),
		EventsEnabled:  buildEventsEnabled("events_enabled", data),
		Global:         data.Get("global").(bool),
		// Headers:                    data.Get("headers").(map[string]string),
		HttpAuthenticationPassword: data.Get("http_authentication_password").(string),
		HttpAuthenticationUsername: data.Get("http_authentication_username").(string),
		ReadTimeout:                data.Get("read_timeout").(int),
		SslCertificate:             data.Get("ssl_certificate").(string),
		Url:                        data.Get("url").(string),
		SignatureConfiguration:     buildSignatureConfiguration(data),
	}

	if i, ok := data.GetOk("headers"); ok {
		wh.Headers = intMapToStringMap(i.(map[string]interface{}))
	}

	return wh
}

func buildSignatureConfiguration(data *schema.ResourceData) fusionauth.WebhookSignatureConfiguration {
	return fusionauth.WebhookSignatureConfiguration{
		Enableable:   buildEnableable("signature_configuration.0.enabled", data),
		SigningKeyId: data.Get("signature_configuration.0.signing_key_id").(string),
	}
}

func buildEventsEnabled(key string, data *schema.ResourceData) map[fusionauth.EventType]bool {
	prefix := key + ".0."
	return map[fusionauth.EventType]bool{
		fusionauth.EventType_AuditLogCreate:                 data.Get(prefix + "audit_log_create").(bool),
		fusionauth.EventType_EventLogCreate:                 data.Get(prefix + "event_log_create").(bool),
		fusionauth.EventType_JWTPublicKeyUpdate:             data.Get(prefix + "jwt_public_key_update").(bool),
		fusionauth.EventType_JWTRefresh:                     data.Get(prefix + "jwt_refresh").(bool),
		fusionauth.EventType_JWTRefreshTokenRevoke:          data.Get(prefix + "jwt_refresh_token_revoke").(bool),
		fusionauth.EventType_KickstartSuccess:               data.Get(prefix + "kickstart_success").(bool),
		fusionauth.EventType_UserAction:                     data.Get(prefix + "user_action").(bool),
		fusionauth.EventType_UserBulkCreate:                 data.Get(prefix + "user_bulk_create").(bool),
		fusionauth.EventType_UserCreate:                     data.Get(prefix + "user_create").(bool),
		fusionauth.EventType_UserCreateComplete:             data.Get(prefix + "user_create_complete").(bool),
		fusionauth.EventType_UserDeactivate:                 data.Get(prefix + "user_deactivate").(bool),
		fusionauth.EventType_UserDelete:                     data.Get(prefix + "user_delete").(bool),
		fusionauth.EventType_UserDeleteComplete:             data.Get(prefix + "user_delete_complete").(bool),
		fusionauth.EventType_UserEmailUpdate:                data.Get(prefix + "user_email_update").(bool),
		fusionauth.EventType_UserEmailVerified:              data.Get(prefix + "user_email_verified").(bool),
		fusionauth.EventType_UserIdentityProviderLink:       data.Get(prefix + "user_identity_provider_link").(bool),
		fusionauth.EventType_UserIdentityProviderUnlink:     data.Get(prefix + "user_identity_provider_unlink").(bool),
		fusionauth.EventType_UserLoginIdDuplicateOnCreate:   data.Get(prefix + "user_login_id_duplicate_create").(bool),
		fusionauth.EventType_UserLoginIdDuplicateOnUpdate:   data.Get(prefix + "user_login_id_duplicate_update").(bool),
		fusionauth.EventType_UserLoginFailed:                data.Get(prefix + "user_login_failed").(bool),
		fusionauth.EventType_UserLoginNewDevice:             data.Get(prefix + "user_login_new_device").(bool),
		fusionauth.EventType_UserLoginSuccess:               data.Get(prefix + "user_login_success").(bool),
		fusionauth.EventType_UserLoginSuspicious:            data.Get(prefix + "user_login_suspicious").(bool),
		fusionauth.EventType_UserPasswordBreach:             data.Get(prefix + "user_password_breach").(bool),
		fusionauth.EventType_UserPasswordResetSend:          data.Get(prefix + "user_password_reset_send").(bool),
		fusionauth.EventType_UserPasswordResetStart:         data.Get(prefix + "user_password_reset_start").(bool),
		fusionauth.EventType_UserPasswordResetSuccess:       data.Get(prefix + "user_password_reset_success").(bool),
		fusionauth.EventType_UserPasswordUpdate:             data.Get(prefix + "user_password_update").(bool),
		fusionauth.EventType_UserReactivate:                 data.Get(prefix + "user_reactivate").(bool),
		fusionauth.EventType_UserRegistrationCreate:         data.Get(prefix + "user_registration_create").(bool),
		fusionauth.EventType_UserRegistrationCreateComplete: data.Get(prefix + "user_registration_create_complete").(bool),
		fusionauth.EventType_UserRegistrationDelete:         data.Get(prefix + "user_registration_delete").(bool),
		fusionauth.EventType_UserRegistrationDeleteComplete: data.Get(prefix + "user_registration_delete_complete").(bool),
		fusionauth.EventType_UserRegistrationUpdate:         data.Get(prefix + "user_registration_update").(bool),
		fusionauth.EventType_UserRegistrationUpdateComplete: data.Get(prefix + "user_registration_update_complete").(bool),
		fusionauth.EventType_UserRegistrationVerified:       data.Get(prefix + "user_registration_verified").(bool),
		fusionauth.EventType_UserTwoFactorMethodAdd:         data.Get(prefix + "user_two_factor_method_add").(bool),
		fusionauth.EventType_UserTwoFactorMethodRemove:      data.Get(prefix + "user_two_factor_method_remove").(bool),
		fusionauth.EventType_UserUpdate:                     data.Get(prefix + "user_update").(bool),
		fusionauth.EventType_UserUpdateComplete:             data.Get(prefix + "user_update_complete").(bool),
	}
}

func createWebhook(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	l := buildWebhook(data)
	resp, faErrs, err := client.FAClient.CreateWebhook("", fusionauth.WebhookRequest{
		Webhook: l,
	})
	if err != nil {
		return diag.Errorf("CreateWebhook err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(resp.Webhook.Id)
	return nil
}

func readWebhook(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, err := client.FAClient.RetrieveWebhook(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	l := resp.Webhook
	if err := data.Set("tenant_ids", l.TenantIds); err != nil {
		return diag.Errorf("webhook.tenant_ids: %s", err.Error())
	}
	if err := data.Set("connect_timeout", l.ConnectTimeout); err != nil {
		return diag.Errorf("webhook.connect_timeout: %s", err.Error())
	}
	if err := data.Set("description", l.Description); err != nil {
		return diag.Errorf("webhook.description: %s", err.Error())
	}

	err = data.Set("events_enabled", []map[string]interface{}{
		{
			"audit_log_create":                  l.EventsEnabled[fusionauth.EventType_AuditLogCreate],
			"event_log_create":                  l.EventsEnabled[fusionauth.EventType_EventLogCreate],
			"jwt_public_key_update":             l.EventsEnabled[fusionauth.EventType_JWTPublicKeyUpdate],
			"jwt_refresh":                       l.EventsEnabled[fusionauth.EventType_JWTRefresh],
			"jwt_refresh_token_revoke":          l.EventsEnabled[fusionauth.EventType_JWTRefreshTokenRevoke],
			"kickstart_success":                 l.EventsEnabled[fusionauth.EventType_KickstartSuccess],
			"user_action":                       l.EventsEnabled[fusionauth.EventType_UserAction],
			"user_bulk_create":                  l.EventsEnabled[fusionauth.EventType_UserBulkCreate],
			"user_create":                       l.EventsEnabled[fusionauth.EventType_UserCreate],
			"user_create_complete":              l.EventsEnabled[fusionauth.EventType_UserCreateComplete],
			"user_deactivate":                   l.EventsEnabled[fusionauth.EventType_UserDeactivate],
			"user_delete":                       l.EventsEnabled[fusionauth.EventType_UserDelete],
			"user_delete_complete":              l.EventsEnabled[fusionauth.EventType_UserDeleteComplete],
			"user_email_update":                 l.EventsEnabled[fusionauth.EventType_UserEmailUpdate],
			"user_email_verified":               l.EventsEnabled[fusionauth.EventType_UserEmailVerified],
			"user_identity_provider_link":       l.EventsEnabled[fusionauth.EventType_UserIdentityProviderLink],
			"user_identity_provider_unlink":     l.EventsEnabled[fusionauth.EventType_UserIdentityProviderUnlink],
			"user_login_id_duplicate_create":    l.EventsEnabled[fusionauth.EventType_UserLoginIdDuplicateOnCreate],
			"user_login_id_duplicate_update":    l.EventsEnabled[fusionauth.EventType_UserLoginIdDuplicateOnUpdate],
			"user_login_failed":                 l.EventsEnabled[fusionauth.EventType_UserLoginFailed],
			"user_login_new_device":             l.EventsEnabled[fusionauth.EventType_UserLoginNewDevice],
			"user_login_success":                l.EventsEnabled[fusionauth.EventType_UserLoginSuccess],
			"user_login_suspicious":             l.EventsEnabled[fusionauth.EventType_UserLoginSuspicious],
			"user_password_breach":              l.EventsEnabled[fusionauth.EventType_UserPasswordBreach],
			"user_password_reset_send":          l.EventsEnabled[fusionauth.EventType_UserPasswordResetSend],
			"user_password_reset_start":         l.EventsEnabled[fusionauth.EventType_UserPasswordResetStart],
			"user_password_reset_success":       l.EventsEnabled[fusionauth.EventType_UserPasswordResetSuccess],
			"user_password_update":              l.EventsEnabled[fusionauth.EventType_UserPasswordUpdate],
			"user_reactivate":                   l.EventsEnabled[fusionauth.EventType_UserReactivate],
			"user_registration_create":          l.EventsEnabled[fusionauth.EventType_UserRegistrationCreate],
			"user_registration_create_complete": l.EventsEnabled[fusionauth.EventType_UserRegistrationCreateComplete],
			"user_registration_delete":          l.EventsEnabled[fusionauth.EventType_UserRegistrationDelete],
			"user_registration_delete_complete": l.EventsEnabled[fusionauth.EventType_UserRegistrationDeleteComplete],
			"user_registration_update":          l.EventsEnabled[fusionauth.EventType_UserRegistrationUpdate],
			"user_registration_update_complete": l.EventsEnabled[fusionauth.EventType_UserRegistrationUpdateComplete],
			"user_registration_verified":        l.EventsEnabled[fusionauth.EventType_UserRegistrationVerified],
			"user_two_factor_method_add":        l.EventsEnabled[fusionauth.EventType_UserTwoFactorMethodAdd],
			"user_two_factor_method_remove":     l.EventsEnabled[fusionauth.EventType_UserTwoFactorMethodRemove],
			"user_update":                       l.EventsEnabled[fusionauth.EventType_UserUpdate],
			"user_update_complete":              l.EventsEnabled[fusionauth.EventType_UserUpdateComplete],
		},
	})
	if err != nil {
		return diag.Errorf("webhook.events_enabled: %s", err.Error())
	}

	if err := data.Set("global", l.Global); err != nil {
		return diag.Errorf("webhook.global: %s", err.Error())
	}
	if err := data.Set("headers", l.Headers); err != nil {
		return diag.Errorf("webhook.headers: %s", err.Error())
	}
	if err := data.Set("http_authentication_password", l.HttpAuthenticationPassword); err != nil {
		return diag.Errorf("webhook.http_authentication_password: %s", err.Error())
	}
	if err := data.Set("http_authentication_username", l.HttpAuthenticationUsername); err != nil {
		return diag.Errorf("webhook.http_authentication_username: %s", err.Error())
	}
	if err := data.Set("read_timeout", l.ReadTimeout); err != nil {
		return diag.Errorf("webhook.read_timeout: %s", err.Error())
	}
	if err := data.Set("ssl_certificate", l.SslCertificate); err != nil {
		return diag.Errorf("webhook.ssl_certificate: %s", err.Error())
	}
	if err := data.Set("url", l.Url); err != nil {
		return diag.Errorf("webhook.url: %s", err.Error())
	}
	return nil
}

func updateWebhook(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	l := buildWebhook(data)

	resp, faErrs, err := client.FAClient.UpdateWebhook(data.Id(), fusionauth.WebhookRequest{
		Webhook: l,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteWebhook(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteWebhook(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
