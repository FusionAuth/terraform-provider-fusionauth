package fusionauth

import (
	"context"
	"encoding/json"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type FacebookIdentityProviderBody struct {
	IdentityProvider fusionauth.FacebookIdentityProvider `json:"identityProvider"`
}

type FacebookAppConfig struct {
	AppID              string `json:"appId,omitempty"`
	ButtonText         string `json:"buttonText,omitempty"`
	ClientSecret       string `json:"client_secret,omitempty"`
	CreateRegistration bool   `json:"createRegistration"`
	Enabled            bool   `json:"enabled"`
	Fields             string `json:"fields,omitempty"`
	Permissions        string `json:"permissions,omitempty"`
}

func resourceIDPFacebook() *schema.Resource {
	return &schema.Resource{
		CreateContext: createIDPFacebook,
		ReadContext:   readIDPFacebook,
		UpdateContext: updateIDPFacebook,
		DeleteContext: deleteIdentityProvider,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The top-level Facebook appId for your Application. This value is retrieved from the Facebook developer website when you setup your Facebook developer account.",
			},
			"application_configuration": {
				Optional:    true,
				Type:        schema.TypeSet,
				Description: "The configuration for each Application that the identity provider is enabled for.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
							Description:  "ID of the FusionAuth Application to apply this configuration to.",
						},
						"app_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for the top level app id.",
						},
						"button_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for the top level button text.",
						},
						"client_secret": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "This is an optional Application specific override for the top level client secret.",
						},
						"create_registration": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Determines if a UserRegistration is created for the User automatically or not. If a user doesn’t exist in FusionAuth and logs in through an identity provider, this boolean controls whether or not FusionAuth creates a registration for the User in the Application they are logging into.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines if this identity provider is enabled for the Application specified by the application id key.",
						},
						"fields": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for the top level fields.",
						},
						"permissions": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for the top level permissions.",
						},
					},
				},
			},
			"button_text": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The top-level button text to use on the FusionAuth login page for this Identity Provider.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The top-level client secret to use with the Facebook Identity Provider when retrieving the long-lived token. This value is retrieved from the Facebook developer website when you setup your Facebook developer account.",
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if debug is enabled for this provider. When enabled, an Event Log is created each time this provider is invoked to reconcile a login.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if this provider is enabled. If it is false then it will be disabled globally.",
			},
			"fields": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "email",
				Description: "The top-level fields that you are requesting from Facebook.",
			},
			"lambda_reconcile_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user. The specified Lambda Id must be of type FacebookReconcile.",
			},
			"linking_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CreatePendingLink",
					"LinkAnonymously",
					"LinkByEmail",
					"LinkByEmailForExistingUser",
					"LinkByUsername",
					"LinkByUsernameForExistingUser",
				}, false),
				Description: "The linking strategy to use when creating the link between the Facebook Identity Provider and the user.",
			},
			"login_method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UseRedirect",
				ValidateFunc: validation.StringInSlice([]string{
					"UsePopup",
					"UseRedirect",
				}, false),
				Description: "The login method to use for this Identity Provider.",
			},
			"permissions": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "email",
				Description: "The top-level permissions that your application is asking of the user’s Facebook account.",
			},
			"tenant_configuration": {
				Optional:    true,
				Type:        schema.TypeSet,
				Description: "The configuration for each Tenant that limits the number of links a user may have for a particular identity provider.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
						},
						"limit_user_link_count_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "When enabled, the number of identity provider links a user may create is enforced by maximumLinks",
						},
						"limit_user_link_count_maximum_links": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     42,
							Description: "Determines if this provider is enabled. If it is false then it will be disabled globally.",
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func createIDPFacebook(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	fbIDP := buildIDPFacebook(data)
	b, err := json.Marshal(fbIDP)
	if err != nil {
		return diag.FromErr(err)
	}

	client := i.(Client)
	bb, err := createIdentityProvider(b, client, "")
	if err != nil {
		return diag.FromErr(err)
	}

	err = json.Unmarshal(bb, &fbIDP)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fbIDP.IdentityProvider.Id)
	return buildResourceFromIDPFacebook(fbIDP.IdentityProvider, data)
}

func readIDPFacebook(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	b, err := readIdentityProvider(data.Id(), client)
	if err != nil {
		return diag.FromErr(err)
	}

	var ipb FacebookIdentityProviderBody
	err = json.Unmarshal(b, &ipb)
	if err != nil {
		return diag.FromErr(err)
	}

	return buildResourceFromIDPFacebook(ipb.IdentityProvider, data)
}

func updateIDPFacebook(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	fbIDP := buildIDPFacebook(data)
	b, err := json.Marshal(fbIDP)
	if err != nil {
		return diag.FromErr(err)
	}

	client := i.(Client)
	bb, err := updateIdentityProvider(b, data.Id(), client)
	if err != nil {
		return diag.FromErr(err)
	}

	err = json.Unmarshal(bb, &fbIDP)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fbIDP.IdentityProvider.Id)
	return buildResourceFromIDPFacebook(fbIDP.IdentityProvider, data)
}

func buildIDPFacebook(data *schema.ResourceData) FacebookIdentityProviderBody {
	fbIDP := fusionauth.FacebookIdentityProvider{
		BaseIdentityProvider: fusionauth.BaseIdentityProvider{
			Debug:      data.Get("debug").(bool),
			Enableable: buildEnableable("enabled", data),
			LambdaConfiguration: fusionauth.ProviderLambdaConfiguration{
				ReconcileId: data.Get("lambda_reconcile_id").(string),
			},
			Type:            fusionauth.IdentityProviderType_Facebook,
			LinkingStrategy: fusionauth.IdentityProviderLinkingStrategy(data.Get("linking_strategy").(string)),
		},
		AppId:        data.Get("app_id").(string),
		ButtonText:   data.Get("button_text").(string),
		ClientSecret: data.Get("client_secret").(string),
		Fields:       data.Get("fields").(string),
		LoginMethod:  fusionauth.IdentityProviderLoginMethod(data.Get("login_method").(string)),
		Permissions:  data.Get("permissions").(string),
	}

	// Compute application specific configuration and overrides.
	fbIDP.ApplicationConfiguration = buildFacebookAppConfig("application_configuration", data)
	fbIDP.TenantConfiguration = buildTenantConfiguration(data)
	return FacebookIdentityProviderBody{
		IdentityProvider: fbIDP,
	}
}

// buildFacebookAppConfig transforms the incoming application configuration as
// recorded in Terraform to the format as required by the FusionAuth API.
func buildFacebookAppConfig(key string, data *schema.ResourceData) map[string]interface{} {
	m := make(map[string]interface{})
	s := data.Get(key)
	set, ok := s.(*schema.Set)
	if !ok {
		return m
	}

	l := set.List()
	for _, x := range l {
		ac := x.(map[string]interface{})
		aid := ac["application_id"].(string)
		m[aid] = FacebookAppConfig{
			AppID:              ac["app_id"].(string),
			ButtonText:         ac["button_text"].(string),
			ClientSecret:       ac["client_secret"].(string),
			CreateRegistration: ac["create_registration"].(bool),
			Enabled:            ac["enabled"].(bool),
			Fields:             ac["fields"].(string),
			Permissions:        ac["permissions"].(string),
		}
	}

	return m
}

// buildResourceFromIDPFacebook writes changes back to terraform data with the
// provided facebook identity provider response.
func buildResourceFromIDPFacebook(res fusionauth.FacebookIdentityProvider, data *schema.ResourceData) diag.Diagnostics {
	if err := data.Set("app_id", res.AppId); err != nil {
		return diag.Errorf("idpFacebook.app_id: %s", err.Error())
	}
	if err := data.Set("button_text", res.ButtonText); err != nil {
		return diag.Errorf("idpFacebook.button_text: %s", err.Error())
	}
	if err := data.Set("client_secret", res.ClientSecret); err != nil {
		return diag.Errorf("idpFacebook.client_secret: %s", err.Error())
	}
	if err := data.Set("debug", res.Debug); err != nil {
		return diag.Errorf("idpFacebook.debug: %s", err.Error())
	}
	if err := data.Set("enabled", res.Enabled); err != nil {
		return diag.Errorf("idpFacebook.enabled: %s", err.Error())
	}
	if err := data.Set("fields", res.Fields); err != nil {
		return diag.Errorf("idpFacebook.fields: %s", err.Error())
	}
	if err := data.Set("lambda_reconcile_id", res.LambdaConfiguration.ReconcileId); err != nil {
		return diag.Errorf("idpFacebook.lambda_reconcile_id: %s", err.Error())
	}
	if err := data.Set("linking_strategy", res.LinkingStrategy); err != nil {
		return diag.Errorf("idpFacebook.linking_strategy: %s", err.Error())
	}
	if err := data.Set("login_method", res.LoginMethod); err != nil {
		return diag.Errorf("idpFacebook.login_method: %s", err.Error())
	}
	if err := data.Set("permissions", res.Permissions); err != nil {
		return diag.Errorf("idpFacebook.permissions: %s", err.Error())
	}

	// Since this is coming down as an interface and would end up being map[string]interface{}
	// with one of the values being map[string]interface{}
	b, _ := json.Marshal(res.ApplicationConfiguration)
	m := make(map[string]FacebookAppConfig)
	_ = json.Unmarshal(b, &m)

	ac := make([]map[string]interface{}, 0, len(res.ApplicationConfiguration))
	for k, v := range m {
		ac = append(ac, map[string]interface{}{
			"application_id":      k,
			"app_id":              v.AppID,
			"button_text":         v.ButtonText,
			"client_secret":       v.ClientSecret,
			"create_registration": v.CreateRegistration,
			"enabled":             v.Enabled,
			"fields":              v.Fields,
			"permissions":         v.Permissions,
		})
	}
	if err := data.Set("application_configuration", ac); err != nil {
		return diag.Errorf("idpFacebook.application_configuration: %s", err.Error())
	}

	tc := buildTenantConfigurationResource(res.TenantConfiguration)
	if err := data.Set("tenant_configuration", tc); err != nil {
		return diag.Errorf("idpFacebook.tenant_configuration: %s", err.Error())
	}

	return nil
}
