package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func newWebhook() *schema.Resource {
	return &schema.Resource{
		Create: createWebhook,
		Read:   readWebhook,
		Update: updateWebhook,
		Delete: deleteWebhook,
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
			State: schema.ImportStatePassthrough,
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

func createWebhook(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	l := buildWebhook(data)
	resp, faErrs, err := client.FAClient.CreateWebhook("", fusionauth.WebhookRequest{
		Webhook: l,
	})
	if err != nil {
		return fmt.Errorf("CreateWebhook err: %v", err)
	}

	if faErrs != nil {
		return fmt.Errorf("CreateWebhook errors: %v", faErrs)
	}
	data.SetId(resp.Webhook.Id)
	return nil
}

func readWebhook(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, err := client.FAClient.RetrieveWebhook(id)
	if err != nil {
		return err
	}

	l := resp.Webhook
	_ = data.Set("application_ids", l.ApplicationIds)
	_ = data.Set("connect_timeout", l.ConnectTimeout)
	_ = data.Set("description", l.Description)
	_ = data.Set("events_enabled.0.user_action", l.EventsEnabled[fusionauth.EventType_UserAction])
	_ = data.Set("events_enabled.0.user_bulk_create", l.EventsEnabled[fusionauth.EventType_UserBulkCreate])
	_ = data.Set("events_enabled.0.user_create", l.EventsEnabled[fusionauth.EventType_UserCreate])
	_ = data.Set("events_enabled.0.user_email_verified", l.EventsEnabled[fusionauth.EventType_UserEmailVerified])
	_ = data.Set("events_enabled.0.user_update", l.EventsEnabled[fusionauth.EventType_UserUpdate])
	_ = data.Set("events_enabled.0.user_deactivate", l.EventsEnabled[fusionauth.EventType_UserDeactivate])
	_ = data.Set("events_enabled.0.user_reactivate", l.EventsEnabled[fusionauth.EventType_UserReactivate])
	_ = data.Set("events_enabled.0.user_login_success", l.EventsEnabled[fusionauth.EventType_UserLoginSuccess])
	_ = data.Set("events_enabled.0.user_login_failed", l.EventsEnabled[fusionauth.EventType_UserLoginFailed])
	_ = data.Set("events_enabled.0.user_password_breach", l.EventsEnabled[fusionauth.EventType_UserPasswordBreach])
	_ = data.Set("events_enabled.0.user_registration_create", l.EventsEnabled[fusionauth.EventType_UserRegistrationCreate])
	_ = data.Set("events_enabled.0.user_registration_update", l.EventsEnabled[fusionauth.EventType_UserRegistrationUpdate])
	_ = data.Set("events_enabled.0.user_registration_delete", l.EventsEnabled[fusionauth.EventType_UserRegistrationDelete])
	_ = data.Set("events_enabled.0.user_registration_verified", l.EventsEnabled[fusionauth.EventType_UserRegistrationVerified])
	_ = data.Set("events_enabled.0.user_delete", l.EventsEnabled[fusionauth.EventType_UserDelete])
	_ = data.Set("events_enabled.0.jwt_public_key_update", l.EventsEnabled[fusionauth.EventType_JWTPublicKeyUpdate])
	_ = data.Set("events_enabled.0.jwt_refresh", l.EventsEnabled[fusionauth.EventType_JWTRefresh])
	_ = data.Set("events_enabled.0.jwt_refresh_token_revoke", l.EventsEnabled[fusionauth.EventType_JWTRefreshTokenRevoke])
	_ = data.Set("global", l.Global)
	_ = data.Set("headers", l.Headers)
	_ = data.Set("http_authentication_password", l.HttpAuthenticationPassword)
	_ = data.Set("http_authentication_username", l.HttpAuthenticationUsername)
	_ = data.Set("read_timeout", l.ReadTimeout)
	_ = data.Set("ssl_certificate", l.SslCertificate)
	_ = data.Set("url", l.Url)
	return nil
}

func updateWebhook(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	l := buildWebhook(data)

	_, faErrs, err := client.FAClient.UpdateWebhook(data.Id(), fusionauth.WebhookRequest{
		Webhook: l,
	})
	if err != nil {
		return err
	}

	if faErrs != nil {
		return fmt.Errorf("UpdateWebhook errors: %v", faErrs)
	}

	return nil
}

func deleteWebhook(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	_, faErrs, err := client.FAClient.DeleteWebhook(id)
	if err != nil {
		return err
	}

	if faErrs != nil {
		return fmt.Errorf("DeleteWebhook errors: %v", faErrs)
	}

	return nil
}
