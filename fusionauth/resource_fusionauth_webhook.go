package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func newWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: createWebhook,
		ReadContext:   readWebhook,
		UpdateContext: updateWebhook,
		DeleteContext: deleteWebhook,
		Schema: map[string]*schema.Schema{
			"application_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The Ids of the Applications that this Webhook should be associated with. If no Ids are specified and the global field is false, this Webhook will not be used.",
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
						"user_email_verified": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user verifies their email address ",
						},
						"user_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user is updated",
						},
						"user_deactivate": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user is deactivated",
						},
						"user_reactivate": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user is reactivated",
						},
						"user_login_success": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user completes a login request ",
						},
						"user_login_failed": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user fails a login request",
						},
						"user_password_breach": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When Reactor detects a user is using a potentially breached password (requires an activated license)",
						},
						"user_registration_create": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user registration is created",
						},
						"user_registration_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user registration is updated",
						},
						"user_registration_delete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user registration is deleted",
						},
						"user_registration_verified": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user completes registration verification",
						},
						"user_delete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When a user is deleted",
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
		ApplicationIds: handleStringSlice("application_ids", data),
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
	}

	if hi, ok := data.GetOk("headers"); ok {
		h := hi.(map[string]interface{})
		m := make(map[string]string)
		for k, v := range h {
			m[k] = v.(string)
		}
		wh.Headers = m
	}

	return wh
}

func buildEventsEnabled(key string, data *schema.ResourceData) map[fusionauth.EventType]bool {
	prefix := key + ".0."
	return map[fusionauth.EventType]bool{
		fusionauth.EventType_UserAction:               data.Get(prefix + "user_action").(bool),
		fusionauth.EventType_UserBulkCreate:           data.Get(prefix + "user_bulk_create").(bool),
		fusionauth.EventType_UserCreate:               data.Get(prefix + "user_create").(bool),
		fusionauth.EventType_UserEmailVerified:        data.Get(prefix + "user_email_verified").(bool),
		fusionauth.EventType_UserUpdate:               data.Get(prefix + "user_update").(bool),
		fusionauth.EventType_UserDeactivate:           data.Get(prefix + "user_deactivate").(bool),
		fusionauth.EventType_UserReactivate:           data.Get(prefix + "user_reactivate").(bool),
		fusionauth.EventType_UserLoginSuccess:         data.Get(prefix + "user_login_success").(bool),
		fusionauth.EventType_UserLoginFailed:          data.Get(prefix + "user_login_failed").(bool),
		fusionauth.EventType_UserPasswordBreach:       data.Get(prefix + "user_password_breach").(bool),
		fusionauth.EventType_UserRegistrationCreate:   data.Get(prefix + "user_registration_create").(bool),
		fusionauth.EventType_UserRegistrationUpdate:   data.Get(prefix + "user_registration_update").(bool),
		fusionauth.EventType_UserRegistrationDelete:   data.Get(prefix + "user_registration_delete").(bool),
		fusionauth.EventType_UserRegistrationVerified: data.Get(prefix + "user_registration_verified").(bool),
		fusionauth.EventType_UserDelete:               data.Get(prefix + "user_delete").(bool),
		fusionauth.EventType_JWTPublicKeyUpdate:       data.Get(prefix + "jwt_public_key_update").(bool),
		fusionauth.EventType_JWTRefresh:               data.Get(prefix + "jwt_refresh").(bool),
		fusionauth.EventType_JWTRefreshTokenRevoke:    data.Get(prefix + "jwt_refresh_token_revoke").(bool),
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
	if err := data.Set("application_ids", l.ApplicationIds); err != nil {
		return diag.Errorf("webhook.application_ids: %s", err.Error())
	}
	if err := data.Set("connect_timeout", l.ConnectTimeout); err != nil {
		return diag.Errorf("webhook.connect_timeout: %s", err.Error())
	}
	if err := data.Set("description", l.Description); err != nil {
		return diag.Errorf("webhook.description: %s", err.Error())
	}

	err = data.Set("events_enabled", []map[string]interface{}{
		{
			"user_action":                l.EventsEnabled[fusionauth.EventType_UserAction],
			"user_bulk_create":           l.EventsEnabled[fusionauth.EventType_UserBulkCreate],
			"user_create":                l.EventsEnabled[fusionauth.EventType_UserCreate],
			"user_email_verified":        l.EventsEnabled[fusionauth.EventType_UserEmailVerified],
			"user_update":                l.EventsEnabled[fusionauth.EventType_UserUpdate],
			"user_deactivate":            l.EventsEnabled[fusionauth.EventType_UserDeactivate],
			"user_reactivate":            l.EventsEnabled[fusionauth.EventType_UserReactivate],
			"user_login_success":         l.EventsEnabled[fusionauth.EventType_UserLoginSuccess],
			"user_login_failed":          l.EventsEnabled[fusionauth.EventType_UserLoginFailed],
			"user_password_breach":       l.EventsEnabled[fusionauth.EventType_UserPasswordBreach],
			"user_registration_create":   l.EventsEnabled[fusionauth.EventType_UserRegistrationCreate],
			"user_registration_update":   l.EventsEnabled[fusionauth.EventType_UserRegistrationUpdate],
			"user_registration_delete":   l.EventsEnabled[fusionauth.EventType_UserRegistrationDelete],
			"user_registration_verified": l.EventsEnabled[fusionauth.EventType_UserRegistrationVerified],
			"user_delete":                l.EventsEnabled[fusionauth.EventType_UserDelete],
			"jwt_public_key_update":      l.EventsEnabled[fusionauth.EventType_JWTPublicKeyUpdate],
			"jwt_refresh":                l.EventsEnabled[fusionauth.EventType_JWTRefresh],
			"jwt_refresh_token_revoke":   l.EventsEnabled[fusionauth.EventType_JWTRefreshTokenRevoke],
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
