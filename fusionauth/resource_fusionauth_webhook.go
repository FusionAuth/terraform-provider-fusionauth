package fusionauth

import (
	"context"
	"encoding/json"
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
				Computed:    true,
				Description: "The Ids of the Tenants that this Webhook should be associated with. If no Ids are specified and the global field is false, this Webhook will not be used.",
			},
			"connect_timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The connection timeout in milliseconds used when FusionAuth sends events to the Webhook.",
			},
			"data": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "An object that can hold any information about the Webhook that should be persisted. Please review the limits on data field types as you plan for and build your custom data schema. Must be a JSON string.",
				DiffSuppressFunc: diffSuppressJSON,
				ValidateFunc:     validation.StringIsJSON,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the Webhook. This is used for display purposes only.",
			},
			"events_enabled": {
				Type:             schema.TypeList,
				MaxItems:         1,
				Optional:         true,
				DiffSuppressFunc: suppressBlockDiff,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audit_log_create": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "An audit log was created",
						},
						"event_log_create": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "An event log was created",
						},
						"group_create": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A group is being created",
						},
						"group_create_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A create group request completed",
						},
						"group_delete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A group is being deleted",
						},
						"group_delete_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A group delete request completed",
						},
						"group_member_add": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user is being added to a group",
						},
						"group_member_add_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user add request has completed",
						},
						"group_member_remove": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user is being removed from a group",
						},
						"group_member_remove_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user remove request has completed",
						},
						"group_member_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A groups membership is being updated",
						},
						"group_member_update_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A group member update request has completed",
						},
						"group_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A group is being updated",
						},
						"group_update_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A request to update a group has completed",
						},
						"jwt_public_key_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A configuration occurred that may affect public keys used to verify a JWT signed by FusionAuth",
						},
						"jwt_refresh": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A JWT was refreshed using a refresh token",
						},
						"jwt_refresh_token_revoke": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "One or more refresh tokens were revoked",
						},
						"kickstart_success": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Kickstart completed successfully the system is ready for use",
						},
						"user_action": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "An action was taken on a user, or an existing event may be changing states if the action is time based",
						},
						"user_bulk_create": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "One or more users were created using the Bulk create API",
						},
						"user_create": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user is being created",
						},
						"user_create_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A create user request completed",
						},
						"user_deactivate": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user is being de-activated, this is synonymous with a soft-delete when using the API",
						},
						"user_delete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user is being deleted",
						},
						"user_delete_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user delete request has completed",
						},
						"user_email_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user updated their email address",
						},
						"user_email_verified": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user has verified their email address",
						},
						"user_identity_provider_link": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A link has been established between a user and an identity provider",
						},
						"user_identity_provider_unlink": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "An existing link has been removed between a user and an identify provider",
						},
						"user_login_id_duplicate_create": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user attempted to register using an email address or username of an existing user",
						},
						"user_login_id_duplicate_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user attempted to modify their email address or username to that of an existing user",
						},
						"user_login_failed": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A login request has failed",
						},
						"user_login_new_device": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user has logged in from a new device",
						},
						"user_login_success": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A login request has succeeded",
						},
						"user_login_suspicious": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A suspicious login request has succeeded. This may be due to an impossible travel calculation, or other indicators",
						},
						"user_password_breach": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user's password has been identified as vulnerable due to being found in one or more breached data sets",
						},
						"user_password_reset_send": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user has been sent an email as part of a password reset workflow",
						},
						"user_password_reset_start": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user has started a password reset workflow",
						},
						"user_password_reset_success": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user has completed a password reset workflow",
						},
						"user_password_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user has updated their password",
						},
						"user_reactivate": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user has been re-activated",
						},
						"user_registration_create": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user registration is being created",
						},
						"user_registration_create_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user registration has been created",
						},
						"user_registration_delete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user registration has been deleted",
						},
						"user_registration_delete_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user registration delete request has completed",
						},
						"user_registration_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user registration is being updated",
						},
						"user_registration_update_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user registration update request has completed",
						},
						"user_registration_verified": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user registration has been verified",
						},
						"user_two_factor_method_add": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user has added a two-factor method",
						},
						"user_two_factor_method_remove": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user has removed a two-factor method",
						},
						"user_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user is being updated",
						},
						"user_update_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A user update request has completed",
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
				Default:     "",
				Description: "An SSL certificate in PEM format that is used to establish the a SSL (TLS specifically) connection to the Webhook.",
			},
			"ssl_certificate_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Id of an existing Key. The X.509 certificate is used for client certificate authentication in requests to the Webhook.",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The fully qualified URL of the Webhook’s endpoint that will accept the event requests from FusionAuth.",
			},
			"webhook_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The Id of the Webhook.",
				ForceNew:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceWebhookV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceWebhookUpgradeV0,
				Version: 0,
			},
		},
	}
}

func buildWebhook(data *schema.ResourceData) fusionauth.Webhook {
	wh := fusionauth.Webhook{
		TenantIds:                  handleStringSlice("tenant_ids", data),
		ConnectTimeout:             data.Get("connect_timeout").(int),
		Description:                data.Get("description").(string),
		EventsEnabled:              buildEventsEnabled("events_enabled", data),
		Global:                     data.Get("global").(bool),
		HttpAuthenticationPassword: data.Get("http_authentication_password").(string),
		HttpAuthenticationUsername: data.Get("http_authentication_username").(string),
		ReadTimeout:                data.Get("read_timeout").(int),
		SslCertificate:             data.Get("ssl_certificate").(string),
		SslCertificateKeyId:        data.Get("ssl_certificate_key_id").(string),
		Url:                        data.Get("url").(string),
		SignatureConfiguration:     buildSignatureConfiguration(data),
	}

	if i, ok := data.GetOk("data"); ok {
		resourceData, _ := jsonStringToMapStringInterface(i.(string))
		wh.Data = resourceData
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
		fusionauth.EventType_GroupCreate:                    data.Get(prefix + "group_create").(bool),
		fusionauth.EventType_GroupCreateComplete:            data.Get(prefix + "group_create_complete").(bool),
		fusionauth.EventType_GroupDelete:                    data.Get(prefix + "group_delete").(bool),
		fusionauth.EventType_GroupDeleteComplete:            data.Get(prefix + "group_delete_complete").(bool),
		fusionauth.EventType_GroupMemberAdd:                 data.Get(prefix + "group_member_add").(bool),
		fusionauth.EventType_GroupMemberAddComplete:         data.Get(prefix + "group_member_add_complete").(bool),
		fusionauth.EventType_GroupMemberRemove:              data.Get(prefix + "group_member_remove").(bool),
		fusionauth.EventType_GroupMemberRemoveComplete:      data.Get(prefix + "group_member_remove_complete").(bool),
		fusionauth.EventType_GroupMemberUpdate:              data.Get(prefix + "group_member_update").(bool),
		fusionauth.EventType_GroupMemberUpdateComplete:      data.Get(prefix + "group_member_update_complete").(bool),
		fusionauth.EventType_GroupUpdate:                    data.Get(prefix + "group_update").(bool),
		fusionauth.EventType_GroupUpdateComplete:            data.Get(prefix + "group_update_complete").(bool),
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
	webhookID := data.Get("webhook_id").(string)
	resp, faErrs, err := client.FAClient.CreateWebhook(webhookID, fusionauth.WebhookRequest{
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
	dataJSON, diags := mapStringInterfaceToJSONString(l.Data)
	if diags != nil {
		return diags
	}
	err = data.Set("data", dataJSON)
	if err != nil {
		return diag.Errorf("webhook.data: %s", err.Error())
	}
	if err := data.Set("description", l.Description); err != nil {
		return diag.Errorf("webhook.description: %s", err.Error())
	}

	err = data.Set("events_enabled", []map[string]interface{}{
		{
			"audit_log_create":                  l.EventsEnabled[fusionauth.EventType_AuditLogCreate],
			"event_log_create":                  l.EventsEnabled[fusionauth.EventType_EventLogCreate],
			"group_create":                      l.EventsEnabled[fusionauth.EventType_GroupCreate],
			"group_create_complete":             l.EventsEnabled[fusionauth.EventType_GroupCreateComplete],
			"group_delete":                      l.EventsEnabled[fusionauth.EventType_GroupDelete],
			"group_delete_complete":             l.EventsEnabled[fusionauth.EventType_GroupDeleteComplete],
			"group_member_add":                  l.EventsEnabled[fusionauth.EventType_GroupMemberAdd],
			"group_member_add_complete":         l.EventsEnabled[fusionauth.EventType_GroupMemberAddComplete],
			"group_member_remove":               l.EventsEnabled[fusionauth.EventType_GroupMemberRemove],
			"group_member_remove_complete":      l.EventsEnabled[fusionauth.EventType_GroupMemberRemoveComplete],
			"group_member_update":               l.EventsEnabled[fusionauth.EventType_GroupMemberUpdate],
			"group_member_update_complete":      l.EventsEnabled[fusionauth.EventType_GroupMemberUpdateComplete],
			"group_update":                      l.EventsEnabled[fusionauth.EventType_GroupUpdate],
			"group_update_complete":             l.EventsEnabled[fusionauth.EventType_GroupUpdateComplete],
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
	if err := data.Set("ssl_certificate_key_id", l.SslCertificateKeyId); err != nil {
		return diag.Errorf("webhook.ssl_certificate_key_id: %s", err.Error())
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

func resourceWebhookV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceWebhookUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if v, ok := rawState["data"]; ok {
		if dataMap, ok := v.(map[string]interface{}); ok {
			jsonBytes, err := json.Marshal(dataMap)
			if err != nil {
				return nil, err
			}

			rawState["data"] = string(jsonBytes)
		}
	}

	return rawState, nil
}
