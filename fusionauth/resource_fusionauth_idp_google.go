package fusionauth

import (
	"context"
	"encoding/json"
	"log"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type GoogleIdentityProviderBody struct {
	IdentityProvider fusionauth.GoogleIdentityProvider `json:"identityProvider"`
}

type GoogleAppConfig struct {
	ButtonText         string                           `json:"buttonText,omitempty"`
	ClientID           string                           `json:"client_id,omitempty"`
	ClientSecret       string                           `json:"client_secret,omitempty"`
	Properties         GoogleIdentityProviderProperties `json:"properties,omitempty"`
	Scope              string                           `json:"scope,omitempty"`
	CreateRegistration bool                             `json:"createRegistration"`
	Enabled            bool                             `json:"enabled"`
}

type GoogleIdentityProviderProperties struct {
	Api    string `json:"api,omitempty"`
	Button string `json:"button,omitempty"`
}

func newIDPGoogle() *schema.Resource {
	return &schema.Resource{
		CreateContext: createIDPGoogle,
		ReadContext:   readIDPGoogle,
		UpdateContext: updateIDPGoogle,
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
						"client_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for the top level client id.",
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
							Description: "Determines if this identity provider is enabled for the Application specified by the applicationId key.",
						},
						"properties": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "An object to hold configuration parameters for the Google Identity Services API.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "?",
									},
									"button": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "?",
									},
								},
							},
						},
						"scope": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for for the top level scope.",
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
				Description: "The top-level Google client id for your Application. This value is retrieved from the Google developer website when you setup your Google developer account.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The top-level client secret to use with the Google Identity Provider when retrieving the long-lived token. This value is retrieved from the Google developer website when you setup your Google developer account.",
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
			"login_method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UseRedirect",
				ValidateFunc: validation.StringInSlice([]string{
					"UsePopup",
					"UseRedirect",
					"UseVendorJavaScript",
				}, false),
				Description: "The login method to use for this Identity Provider.",
			},
			"properties": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				Description:      "An object to hold configuration parameters for the Google Identity Services API.",
				DiffSuppressFunc: suppressBlockDiff,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api": {
							Type:     schema.TypeString,
							Optional: true,
							Default: `# Omit the data- prefix
auto_prompt=true
auto_select=false
cancel_on_tap_outside=false
context=signin
itp_support=true`,
							Description: "Google Identity Services login API configuration in a properties file formatted String. Any attribute from Google's documentation can be added. Properties can be referenced in templates that support Google login to initialize the API via HTML or JavaScript. The properties specified in this field should not include the data- prefix on the property name. If the `login_method` is set to UsePopup and this value contains the conflicting ux_mode=redirect property, that single property will be replaced with ux_mode=popup.",
						},
						"button": {
							Type:     schema.TypeString,
							Optional: true,
							Default: `# Omit the data- prefix
logo_alignment=left
shape=rectangular
size=large
text=signin_with
theme=outline
type=standard`,
							Description: "Google Identity Services button configuration in a properties file formatted String. Any attribute from Google's documentation can be added. Properties can be referenced in templates that support Google login to render the login button via HTML or JavaScript. The properties specified in this field should not include the data- prefix on the property name.",
						},
					},
				},
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The top-level scope that you are requesting from Google.",
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

func buildIDPGoogle(data *schema.ResourceData) GoogleIdentityProviderBody {
	o := fusionauth.GoogleIdentityProvider{
		ButtonText: data.Get("button_text").(string),
		BaseIdentityProvider: fusionauth.BaseIdentityProvider{
			Debug:      data.Get("debug").(bool),
			Enableable: buildEnableable("enabled", data),
			LambdaConfiguration: fusionauth.ProviderLambdaConfiguration{
				ReconcileId: data.Get("lambda_reconcile_id").(string),
			},
			Type:            fusionauth.IdentityProviderType_Google,
			LinkingStrategy: fusionauth.IdentityProviderLinkingStrategy(data.Get("linking_strategy").(string)),
		},
		ClientId:     data.Get("client_id").(string),
		ClientSecret: data.Get("client_secret").(string),
		Properties: fusionauth.GoogleIdentityProviderProperties{
			Api:    data.Get("properties.0.api").(string),
			Button: data.Get("properties.0.button").(string),
		},
		Scope:       data.Get("scope").(string),
		LoginMethod: fusionauth.IdentityProviderLoginMethod(data.Get("login_method").(string)),
	}

	o.ApplicationConfiguration = buildGoogleAppConfig("application_configuration", data)
	o.TenantConfiguration = buildTenantConfiguration(data)
	return GoogleIdentityProviderBody{IdentityProvider: o}
}

func buildGoogleAppConfig(key string, data *schema.ResourceData) map[string]interface{} {
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
		propertiesList := ac["properties"].([]interface{})
		var properties GoogleIdentityProviderProperties
		if len(propertiesList) > 0 {
			prop := propertiesList[0].(map[string]interface{})
			properties = GoogleIdentityProviderProperties{
				Api:    getStringOrDefault(prop, "api", ""),
				Button: getStringOrDefault(prop, "button", ""),
			}
		}
		oc := GoogleAppConfig{
			ButtonText:         ac["button_text"].(string),
			CreateRegistration: ac["create_registration"].(bool),
			Enabled:            ac["enabled"].(bool),
			ClientID:           ac["client_id"].(string),
			ClientSecret:       ac["client_secret"].(string),
			Properties:         properties,
			Scope:              ac["scope"].(string),
		}
		m[aid] = oc
	}
	return m
}

func getStringOrDefault(m map[string]interface{}, key string, defaultValue string) string {
	log.Printf("TOMTEST: %v, %v, %v", m, key, defaultValue)
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultValue
}

func createIDPGoogle(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	o := buildIDPGoogle(data)

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

func readIDPGoogle(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	b, err := readIdentityProvider(data.Id(), client)
	if err != nil {
		return diag.FromErr(err)
	}

	var ipb GoogleIdentityProviderBody
	_ = json.Unmarshal(b, &ipb)

	return buildResourceFromIDPGoogle(ipb.IdentityProvider, data)
}

func buildResourceFromIDPGoogle(o fusionauth.GoogleIdentityProvider, data *schema.ResourceData) diag.Diagnostics {
	if err := data.Set("button_text", o.ButtonText); err != nil {
		return diag.Errorf("idpGoogle.button_text: %s", err.Error())
	}
	if err := data.Set("debug", o.Debug); err != nil {
		return diag.Errorf("idpGoogle.debug: %s", err.Error())
	}
	if err := data.Set("enabled", o.Enabled); err != nil {
		return diag.Errorf("idpGoogle.enabled: %s", err.Error())
	}
	if err := data.Set("lambda_reconcile_id", o.LambdaConfiguration.ReconcileId); err != nil {
		return diag.Errorf("idpGoogle.lambda_reconcile_id: %s", err.Error())
	}
	if err := data.Set("client_id", o.ClientId); err != nil {
		return diag.Errorf("idpGoogle.client_id: %s", err.Error())
	}
	if err := data.Set("client_secret", o.ClientSecret); err != nil {
		return diag.Errorf("idpGoogle.client_secret: %s", err.Error())
	}
	if err := data.Set("properties", []map[string]interface{}{
		{
			"api":    o.Properties.Api,
			"button": o.Properties.Button,
		},
	}); err != nil {
		return diag.Errorf("idpGoogle.properties: %s", err.Error())
	}
	if err := data.Set("scope", o.Scope); err != nil {
		return diag.Errorf("idpGoogle.scope: %s", err.Error())
	}
	if err := data.Set("linking_strategy", o.LinkingStrategy); err != nil {
		return diag.Errorf("idpGoogle.linking_strategy: %s", err.Error())
	}
	if err := data.Set("login_method", o.LoginMethod); err != nil {
		return diag.Errorf("idpGoogle.login_method: %s", err.Error())
	}
	// Since this is coming down as an interface and would end up being map[string]interface{}
	// with one of the values being map[string]interface{}
	b, _ := json.Marshal(o.ApplicationConfiguration)
	m := make(map[string]GoogleAppConfig)
	_ = json.Unmarshal(b, &m)

	ac := make([]map[string]interface{}, 0, len(o.ApplicationConfiguration))
	for k, v := range m {
		properties := []map[string]interface{}{}
		if v.Properties.Api != "" || v.Properties.Button != "" {
			properties = append(properties, map[string]interface{}{
				"api":    v.Properties.Api,
				"button": v.Properties.Button,
			})
		}
		ac = append(ac, map[string]interface{}{
			"application_id":      k,
			"button_text":         v.ButtonText,
			"client_id":           v.ClientID,
			"client_secret":       v.ClientSecret,
			"create_registration": v.CreateRegistration,
			"enabled":             v.Enabled,
			"properties":          properties,
			"scope":               v.Scope,
		})
	}
	if err := data.Set("application_configuration", ac); err != nil {
		return diag.Errorf("idpGoogle.application_configuration: %s", err.Error())
	}

	tc := buildTenantConfigurationResource(o.TenantConfiguration)
	if err := data.Set("tenant_configuration", tc); err != nil {
		return diag.Errorf("idpGoogle.tenant_configuration: %s", err.Error())
	}

	return nil
}

func updateIDPGoogle(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	o := buildIDPGoogle(data)

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
