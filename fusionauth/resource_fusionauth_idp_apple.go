package fusionauth

import (
	"context"
	"encoding/json"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type AppleIdentityProviderBody struct {
	IdentityProvider fusionauth.AppleIdentityProvider `json:"identityProvider"`
}

type AppleAppConfig struct {
	ButtonText         string `json:"buttonText,omitempty"`
	CreateRegistration bool   `json:"createRegistration"`
	Enabled            bool   `json:"enabled"`
	KeyID              string `json:"keyId,omitempty"`
	Scope              string `json:"scope,omitempty"`
	BundleID           string `json:"bundleId,omitempty"`
	ServicesID         string `json:"servicesId,omitempty"`
	TeamID             string `json:"teamId,omitempty"`
}

func resourceIDPApple() *schema.Resource {
	return &schema.Resource{
		CreateContext: createIDPApple,
		ReadContext:   readIDPApple,
		UpdateContext: updateIDPApple,
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
						},
						"button_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for the top level button text.",
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
							Description: "Determines if this provider is enabled. If it is false then it will be disabled globally.",
						},
						"key_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "This is an optional Application specific override for the top level key_id.",
							ValidateFunc: validation.IsUUID,
						},
						"scope": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for for the top level scope.",
						},
						"bundle_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for for the top level bundleId.",
						},
						"services_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for for the top level servicesId.",
						},
						"team_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Apple App ID Prefix, or Team ID found in your Apple Developer Account which has been configured for Sign in with Apple.",
						},
					},
				},
			},
			"button_text": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The top-level button text to use on the FusionAuth login page for this Identity Provider.",
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if debug is enabled for this provider. When enabled, each time this provider is invoked to reconcile a login an Event Log will be created.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if this provider is enabled. If it is false then it will be disabled globally.",
			},
			"key_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The unique Id of the private key downloaded from Apple and imported into Key Master that will be used to sign the client secret.",
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
			},
			"lambda_reconcile_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user.",
				ValidateFunc: validation.IsUUID,
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
					"Unsupported",
				}, false),
				Description: "The linking strategy to use when creating the link between the {idp_display_name} Identity Provider and the user.",
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The top-level space separated scope that you are requesting from Apple.",
			},
			"bundle_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Apple Bundle Id you have configured in your Apple developer account to uniquely identify your native app",
			},
			"services_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Apple Services identifier found in your Apple Developer Account which has been configured for Sign in with Apple.",
			},
			"team_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Apple App ID Prefix, or Team ID found in your Apple Developer Account which has been configured for Sign in with Apple.",
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

func createIDPApple(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	o := buildIDPApple(data)

	b, err := json.Marshal(o)
	if err != nil {
		return diag.FromErr(err)
	}

	client := i.(Client)

	bb, err := createIdentityProvider(b, client, "")
	if err != nil {
		return diag.FromErr(err)
	}

	err = json.Unmarshal(bb, &o)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(o.IdentityProvider.Id)
	return nil
}

func readIDPApple(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	b, err := readIdentityProvider(data.Id(), client)
	if err != nil {
		return diag.FromErr(err)
	}

	var ipb AppleIdentityProviderBody
	_ = json.Unmarshal(b, &ipb)

	return buildResourceFromIDPApple(ipb.IdentityProvider, data)
}

func updateIDPApple(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	o := buildIDPApple(data)

	b, err := json.Marshal(o)
	if err != nil {
		return diag.FromErr(err)
	}

	client := i.(Client)
	bb, err := updateIdentityProvider(b, data.Id(), client)
	if err != nil {
		return diag.FromErr(err)
	}

	err = json.Unmarshal(bb, &o)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(o.IdentityProvider.Id)
	return nil
}

func buildIDPApple(data *schema.ResourceData) AppleIdentityProviderBody {
	a := fusionauth.AppleIdentityProvider{
		ButtonText: data.Get("button_text").(string),
		BaseIdentityProvider: fusionauth.BaseIdentityProvider{
			Debug:      data.Get("debug").(bool),
			Enableable: buildEnableable("enabled", data),
			LambdaConfiguration: fusionauth.ProviderLambdaConfiguration{
				ReconcileId: data.Get("lambda_reconcile_id").(string),
			},
			Type:            fusionauth.IdentityProviderType_Apple,
			LinkingStrategy: fusionauth.IdentityProviderLinkingStrategy(data.Get("linking_strategy").(string)),
		},
		KeyId:      data.Get("key_id").(string),
		Scope:      data.Get("scope").(string),
		BundleId:   data.Get("bundle_id").(string),
		ServicesId: data.Get("services_id").(string),
		TeamId:     data.Get("team_id").(string),
	}

	a.ApplicationConfiguration = buildAppleAppConfig("application_configuration", data)
	a.TenantConfiguration = buildTenantConfiguration(data)

	return AppleIdentityProviderBody{IdentityProvider: a}
}

func buildAppleAppConfig(key string, data *schema.ResourceData) map[string]interface{} {
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
		oc := AppleAppConfig{
			ButtonText:         ac["button_text"].(string),
			CreateRegistration: ac["create_registration"].(bool),
			Enabled:            ac["enabled"].(bool),
			KeyID:              ac["key_id"].(string),
			Scope:              ac["scope"].(string),
			BundleID:           ac["bundle_id"].(string),
			ServicesID:         ac["services_id"].(string),
			TeamID:             ac["team_id"].(string),
		}
		m[aid] = oc
	}
	return m
}

func buildResourceFromIDPApple(o fusionauth.AppleIdentityProvider, data *schema.ResourceData) diag.Diagnostics {
	if err := data.Set("button_text", o.ButtonText); err != nil {
		return diag.Errorf("idpApple.button_text: %s", err.Error())
	}
	if err := data.Set("debug", o.Debug); err != nil {
		return diag.Errorf("idpApple.debug: %s", err.Error())
	}
	if err := data.Set("enabled", o.Enabled); err != nil {
		return diag.Errorf("idpApple.enabled: %s", err.Error())
	}
	if err := data.Set("key_id", o.KeyId); err != nil {
		return diag.Errorf("idpApple.key_id: %s", err.Error())
	}
	if err := data.Set("lambda_reconcile_id", o.LambdaConfiguration.ReconcileId); err != nil {
		return diag.Errorf("idpApple.lambda_reconcile_id: %s", err.Error())
	}
	if err := data.Set("scope", o.Scope); err != nil {
		return diag.Errorf("idpApple.scope: %s", err.Error())
	}
	if err := data.Set("bundle_id", o.BundleId); err != nil {
		return diag.Errorf("idpApple.bundle_id: %s", err.Error())
	}
	if err := data.Set("services_id", o.ServicesId); err != nil {
		return diag.Errorf("idpApple.services_id: %s", err.Error())
	}
	if err := data.Set("team_id", o.TeamId); err != nil {
		return diag.Errorf("idpApple.team_id: %s", err.Error())
	}
	if err := data.Set("linking_strategy", o.LinkingStrategy); err != nil {
		return diag.Errorf("idpApple.linking_strategy: %s", err.Error())
	}

	// Since this is coming down as an interface and would end up being map[string]interface{}
	// with one of the values being map[string]interface{}
	b, _ := json.Marshal(o.ApplicationConfiguration)
	m := make(map[string]AppleAppConfig)
	_ = json.Unmarshal(b, &m)

	ac := make([]map[string]interface{}, 0, len(o.ApplicationConfiguration))
	for k, v := range m {
		ac = append(ac, map[string]interface{}{
			"application_id":      k,
			"button_text":         v.ButtonText,
			"create_registration": v.CreateRegistration,
			"enabled":             v.Enabled,
			"key_id":              v.KeyID,
			"scope":               v.Scope,
			"bundle_id":           v.BundleID,
			"services_id":         v.ServicesID,
			"team_id":             v.TeamID,
		})
	}
	if err := data.Set("application_configuration", ac); err != nil {
		return diag.Errorf("idpApple.application_configuration: %s", err.Error())
	}

	tc := buildTenantConfigurationResource(o.TenantConfiguration)
	if err := data.Set("tenant_configuration", tc); err != nil {
		return diag.Errorf("idpApple.tenant_configuration: %s", err.Error())
	}

	return nil
}
