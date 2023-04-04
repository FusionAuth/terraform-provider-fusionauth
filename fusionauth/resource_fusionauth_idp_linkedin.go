package fusionauth

import (
	"context"
	"encoding/json"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type LinkedInIdentityProviderBody struct {
	IdentityProvider fusionauth.LinkedInIdentityProvider `json:"identityProvider"`
}

type LinkedInAppConfig struct {
	ButtonText         string `json:"buttonText,omitempty"`
	ClientID           string `json:"client_id,omitempty"`
	ClientSecret       string `json:"client_secret,omitempty"`
	CreateRegistration bool   `json:"createRegistration"`
	Enabled            bool   `json:"enabled"`
	Scope              string `json:"scope,omitempty"`
}

func resourceIDPLinkedIn() *schema.Resource {
	return &schema.Resource{
		CreateContext: createIDPLinkedIn,
		ReadContext:   readIDPLinkedIn,
		UpdateContext: updateIDPLinkedIn,
		DeleteContext: deleteIdentityProvider,
		Schema: map[string]*schema.Schema{
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
						"button_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for the top level button text.",
						},
						"client_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "This is an optional Application specific override for the top level client ID.",
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
						"scope": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for the top level scope.",
						},
					},
				},
			},
			"button_text": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The top-level button text to use on the FusionAuth login page for this Identity Provider.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The top-level LinkedIn client id for your Application. This value is retrieved from the LinkedIn developer website when you set up your LinkedIn app.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The top-level client secret to use with the LinkedIn Identity Provider when retrieving the long-lived token. This value is retrieved from the LinkedIn developer website when you set up your LinkedIn app.",
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
			"lambda_reconcile_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user. The specified Lambda Id must be of type LinkedInReconcile.",
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
				Description: "The linking strategy to use when creating the link between the LinkedIn Identity Provider and the user.",
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "r_emailaddress r_liteprofile",
				Description: "The top-level scope that you are requesting from LinkedIn.",
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

func createIDPLinkedIn(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	linkedInIDP := buildIDPLinkedIn(data)
	b, err := json.Marshal(linkedInIDP)
	if err != nil {
		return diag.FromErr(err)
	}

	client := i.(Client)
	bb, err := createIdentityProvider(b, client, "")
	if err != nil {
		return diag.FromErr(err)
	}

	err = json.Unmarshal(bb, &linkedInIDP)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(linkedInIDP.IdentityProvider.Id)
	return buildResourceFromIDPLinkedIn(linkedInIDP.IdentityProvider, data)
}

func readIDPLinkedIn(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	b, err := readIdentityProvider(data.Id(), client)
	if err != nil {
		return diag.FromErr(err)
	}

	var ipb LinkedInIdentityProviderBody
	err = json.Unmarshal(b, &ipb)
	if err != nil {
		return diag.FromErr(err)
	}

	return buildResourceFromIDPLinkedIn(ipb.IdentityProvider, data)
}

func updateIDPLinkedIn(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	linkedInIDP := buildIDPLinkedIn(data)
	b, err := json.Marshal(linkedInIDP)
	if err != nil {
		return diag.FromErr(err)
	}

	client := i.(Client)
	bb, err := updateIdentityProvider(b, data.Id(), client)
	if err != nil {
		return diag.FromErr(err)
	}

	err = json.Unmarshal(bb, &linkedInIDP)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(linkedInIDP.IdentityProvider.Id)
	return buildResourceFromIDPLinkedIn(linkedInIDP.IdentityProvider, data)
}

func buildIDPLinkedIn(data *schema.ResourceData) LinkedInIdentityProviderBody {
	linkedInIDP := fusionauth.LinkedInIdentityProvider{
		BaseIdentityProvider: fusionauth.BaseIdentityProvider{
			Debug:      data.Get("debug").(bool),
			Enableable: buildEnableable("enabled", data),
			LambdaConfiguration: fusionauth.ProviderLambdaConfiguration{
				ReconcileId: data.Get("lambda_reconcile_id").(string),
			},
			Type:            fusionauth.IdentityProviderType_LinkedIn,
			LinkingStrategy: fusionauth.IdentityProviderLinkingStrategy(data.Get("linking_strategy").(string)),
		},
		ButtonText:   data.Get("button_text").(string),
		ClientId:     data.Get("client_id").(string),
		ClientSecret: data.Get("client_secret").(string),
		Scope:        data.Get("scope").(string),
	}

	// Compute application specific configuration and overrides.
	linkedInIDP.ApplicationConfiguration = buildLinkedInAppConfig("application_configuration", data)
	linkedInIDP.TenantConfiguration = buildTenantConfiguration(data)
	return LinkedInIdentityProviderBody{
		IdentityProvider: linkedInIDP,
	}
}

// buildLinkedInAppConfig transforms the incoming application configuration as
// recorded in Terraform to the format as required by the FusionAuth API.
func buildLinkedInAppConfig(key string, data *schema.ResourceData) map[string]interface{} {
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
		m[aid] = LinkedInAppConfig{
			ButtonText:         ac["button_text"].(string),
			ClientID:           ac["client_id"].(string),
			ClientSecret:       ac["client_secret"].(string),
			CreateRegistration: ac["create_registration"].(bool),
			Enabled:            ac["enabled"].(bool),
			Scope:              ac["scope"].(string),
		}
	}

	return m
}

// buildResourceFromIDPLinkedIn writes changes back to terraform data with the
// provided LinkedIn identity provider response.
func buildResourceFromIDPLinkedIn(res fusionauth.LinkedInIdentityProvider, data *schema.ResourceData) diag.Diagnostics {
	if err := data.Set("button_text", res.ButtonText); err != nil {
		return diag.Errorf("idpLinkedIn.button_text: %s", err.Error())
	}
	if err := data.Set("client_id", res.ClientId); err != nil {
		return diag.Errorf("idpLinkedIn.client_id: %s", err.Error())
	}
	if err := data.Set("client_secret", res.ClientSecret); err != nil {
		return diag.Errorf("idpLinkedIn.client_secret: %s", err.Error())
	}
	if err := data.Set("debug", res.Debug); err != nil {
		return diag.Errorf("idpLinkedIn.debug: %s", err.Error())
	}
	if err := data.Set("enabled", res.Enabled); err != nil {
		return diag.Errorf("idpLinkedIn.enabled: %s", err.Error())
	}
	if err := data.Set("lambda_reconcile_id", res.LambdaConfiguration.ReconcileId); err != nil {
		return diag.Errorf("idpLinkedIn.lambda_reconcile_id: %s", err.Error())
	}
	if err := data.Set("linking_strategy", res.LinkingStrategy); err != nil {
		return diag.Errorf("idpLinkedIn.linking_strategy: %s", err.Error())
	}
	if err := data.Set("scope", res.Scope); err != nil {
		return diag.Errorf("idpLinkedIn.scope: %s", err.Error())
	}

	// Since this is coming down as an interface and would end up being map[string]interface{}
	// with one of the values being map[string]interface{}
	b, _ := json.Marshal(res.ApplicationConfiguration)
	m := make(map[string]LinkedInAppConfig)
	_ = json.Unmarshal(b, &m)

	ac := make([]map[string]interface{}, 0, len(res.ApplicationConfiguration))
	for k, v := range m {
		ac = append(ac, map[string]interface{}{
			"application_id":      k,
			"button_text":         v.ButtonText,
			"client_id":           v.ClientID,
			"client_secret":       v.ClientSecret,
			"create_registration": v.CreateRegistration,
			"enabled":             v.Enabled,
			"scope":               v.Scope,
		})
	}
	if err := data.Set("application_configuration", ac); err != nil {
		return diag.Errorf("idpLinkedIn.application_configuration: %s", err.Error())
	}

	tc := buildTenantConfigurationResource(res.TenantConfiguration)
	if err := data.Set("tenant_configuration", tc); err != nil {
		return diag.Errorf("idpLinkedIn.tenant_configuration: %s", err.Error())
	}

	return nil
}
