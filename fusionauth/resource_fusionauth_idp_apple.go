package fusionauth

import (
	"encoding/json"
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

type AppleIdentityProviderBody struct {
	IdentityProvider fusionauth.AppleIdentityProvider `json:"identityProvider"`
}

type AppleAppConfig struct {
	ButtonText         string `json:"buttonText,omitempty"`
	CreateRegistration bool   `json:"createRegistration"`
	Enabled            bool   `json:"enabled,omitempty"`
	KeyID              string `json:"keyId,omitempty"`
	Scope              string `json:"scope,omitempty"`
	ServicesID         string `json:"servicesId,omitempty"`
	TeamID             string `json:"teamId,omitempty"`
}

func resourceIDPApple() *schema.Resource {
	return &schema.Resource{
		Create: createIDPApple,
		Read:   readIDPApple,
		Update: updateIDPApple,
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
			},
			"lambda_reconcile_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user.",
				ValidateFunc: validation.IsUUID,
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The top-level space separated scope that you are requesting from Apple.",
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
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func createIDPApple(data *schema.ResourceData, i interface{}) error {
	o := buildIDPApple(data)

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

func readIDPApple(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	b, err := readIdenityProvider(data.Id(), client)
	if err != nil {
		return err
	}

	var ipb AppleIdentityProviderBody
	_ = json.Unmarshal(b, &ipb)

	return buildResourceFromIDPApple(ipb.IdentityProvider, data)
}

func updateIDPApple(data *schema.ResourceData, i interface{}) error {
	o := buildIDPApple(data)

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

func buildIDPApple(data *schema.ResourceData) AppleIdentityProviderBody {
	a := fusionauth.AppleIdentityProvider{
		ButtonText: data.Get("button_text").(string),
		BaseIdentityProvider: fusionauth.BaseIdentityProvider{
			Debug:      data.Get("debug").(bool),
			Enableable: buildEnableable("enabled", data),
			LambdaConfiguration: fusionauth.ProviderLambdaConfiguration{
				ReconcileId: data.Get("lambda_reconcile_id").(string),
			},
			Type: fusionauth.IdentityProviderType_Apple,
		},
		KeyId:      data.Get("key_id").(string),
		Scope:      data.Get("scope").(string),
		ServicesId: data.Get("services_id").(string),
		TeamId:     data.Get("team_id").(string),
	}
	ac := buildAppleAppConfig("application_configuration", data)
	a.ApplicationConfiguration = ac
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
			ServicesID:         ac["services_id"].(string),
			TeamID:             ac["team_id"].(string),
		}
		m[aid] = oc
	}
	return m
}

func buildResourceFromIDPApple(o fusionauth.AppleIdentityProvider, data *schema.ResourceData) error {
	if err := data.Set("button_text", o.ButtonText); err != nil {
		return fmt.Errorf("idpApple.button_text: %s", err.Error())
	}
	if err := data.Set("debug", o.Debug); err != nil {
		return fmt.Errorf("idpApple.debug: %s", err.Error())
	}
	if err := data.Set("enabled", o.Enabled); err != nil {
		return fmt.Errorf("idpApple.enabled: %s", err.Error())
	}
	if err := data.Set("key_id", o.KeyId); err != nil {
		return fmt.Errorf("idpApple.key_id: %s", err.Error())
	}
	if err := data.Set("lambda_reconcile_id", o.LambdaConfiguration.ReconcileId); err != nil {
		return fmt.Errorf("idpApple.lambda_reconcile_id: %s", err.Error())
	}
	if err := data.Set("scope", o.Scope); err != nil {
		return fmt.Errorf("idpApple.scope: %s", err.Error())
	}
	if err := data.Set("services_id", o.ServicesId); err != nil {
		return fmt.Errorf("idpApple.services_id: %s", err.Error())
	}
	if err := data.Set("team_id", o.TeamId); err != nil {
		return fmt.Errorf("idpApple.team_id: %s", err.Error())
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
			"services_id":         v.ServicesID,
			"team_id":             v.TeamID,
		})
	}
	if err := data.Set("application_configuration", ac); err != nil {
		return fmt.Errorf("idpApple.application_configuration: %s", err.Error())
	}
	return nil
}
