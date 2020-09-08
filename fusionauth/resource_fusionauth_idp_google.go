package fusionauth

import (
	"encoding/json"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

type GoogleIdentityProviderBody struct {
	IdentityProvider fusionauth.GoogleIdentityProvider `json:"identityProvider"`
}

type GoogleAppConfig struct {
	ButtonImageURL     string `json:"buttonImageURL,omitempty"`
	ButtonText         string `json:"buttonText,omitempty"`
	ClientID           string `json:"client_id,omitempty"`
	ClientSecret       string `json:"client_secret,omitempty"`
	Scope              string `json:"scope,omitempty"`
	CreateRegistration bool   `json:"createRegistration,omitempty"`
	Enabled            bool   `json:"enabled,omitempty"`
}

func newIDPGoogle() *schema.Resource {
	return &schema.Resource{
		Create: createIDPGoogle,
		Read:   readIDPGoogle,
		Update: updateIDPGoogle,
		Delete: deleteIdentityProvider,
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user.",
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The top-level scope that you are requesting from Google.",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			Type: fusionauth.IdentityProviderType_Google,
		},
		ClientId:     data.Get("client_id").(string),
		ClientSecret: data.Get("client_secret").(string),
		Scope:        data.Get("scope").(string),
	}

	ac := buildGoogleAppConfig("application_configuration", data)
	o.ApplicationConfiguration = ac
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
		oc := GoogleAppConfig{
			ButtonText:         ac["button_text"].(string),
			CreateRegistration: ac["create_registration"].(bool),
			Enabled:            ac["enabled"].(bool),
			ClientID:           ac["client_id"].(string),
			ClientSecret:       ac["client_secret"].(string),
			Scope:              ac["scope"].(string),
		}
		m[aid] = oc
	}
	return m
}

func createIDPGoogle(data *schema.ResourceData, i interface{}) error {
	o := buildIDPGoogle(data)

	b, err := json.Marshal(o)
	if err != nil {
		return err
	}

	client := i.(Client)

	bb, err := createIdenityProvider(b, client)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bb, &o)
	if err != nil {
		return err
	}

	data.SetId(o.IdentityProvider.Id)
	return nil
}

func readIDPGoogle(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	b, err := readIdenityProvider(data.Id(), client)
	if err != nil {
		return err
	}

	var ipb GoogleIdentityProviderBody
	_ = json.Unmarshal(b, &ipb)

	buildResourceFromIDPGoogle(ipb.IdentityProvider, data)
	return nil
}

func buildResourceFromIDPGoogle(o fusionauth.GoogleIdentityProvider, data *schema.ResourceData) {
	_ = data.Set("button_text", o.ButtonText)
	_ = data.Set("debug", o.Debug)
	_ = data.Set("enabled", o.Enabled)
	_ = data.Set("lambda_reconcile_id", o.LambdaConfiguration.ReconcileId)
	_ = data.Set("client_id", o.ClientId)
	_ = data.Set("client_secret", o.ClientSecret)
	_ = data.Set("scope", o.Scope)

	// Since this is coming down as an interface and would end up being map[string]interface{}
	// with one of the values being map[string]interface{}
	b, _ := json.Marshal(o.ApplicationConfiguration)
	m := make(map[string]GoogleAppConfig)
	_ = json.Unmarshal(b, &m)

	ac := make([]map[string]interface{}, 0, len(o.ApplicationConfiguration))
	for k, v := range m {
		ac = append(ac, map[string]interface{}{
			"application_id":      k,
			"button_image_url":    v.ButtonImageURL,
			"button_text":         v.ButtonText,
			"client_id":           v.ClientID,
			"client_secret":       v.ClientSecret,
			"create_registration": v.CreateRegistration,
			"enabled":             v.Enabled,
			"scope":               v.Scope,
		})
	}
	_ = data.Set("application_configuration", ac)
}

func updateIDPGoogle(data *schema.ResourceData, i interface{}) error {
	o := buildIDPGoogle(data)

	b, err := json.Marshal(o)
	if err != nil {
		return err
	}

	client := i.(Client)
	bb, err := updateIdenityProvider(b, data.Id(), client)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bb, &o)
	if err != nil {
		return err
	}

	data.SetId(o.IdentityProvider.Id)
	return nil
}
