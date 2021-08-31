package fusionauth

import (
	"encoding/json"
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

type OpenIDConnectIdentityProviderBody struct {
	IdentityProvider fusionauth.OpenIdConnectIdentityProvider `json:"identityProvider"`
}

type OpenIDAppConfig struct {
	ButtonImageURL     string          `json:"buttonImageURL,omitempty"`
	ButtonText         string          `json:"buttonText,omitempty"`
	OAuth2             OAuth2AppConfig `json:"oauth2,omitempty"`
	CreateRegistration bool            `json:"createRegistration"`
	Enabled            bool            `json:"enabled,omitempty"`
}

type OAuth2AppConfig struct {
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

func newIDPOpenIDConnect() *schema.Resource {
	return &schema.Resource{
		Create: createOpenIDConnect,
		Read:   readOpenIDConnect,
		Update: updateOpenIDConnect,
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
						"button_image_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for the top level button image URL.",
						},
						"button_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for the top level button text.",
						},
						"oauth2_client_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for the top level client id.",
						},
						"oauth2_client_secret": {
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
						"oauth2_scope": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for for the top level scope.",
						},
					},
				},
			},
			"button_image_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The top-level button image (URL) to use on the FusionAuth login page for this Identity Provider.",
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
			"domains": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "This is an optional list of domains that this OpenID Connect provider should be used for. This converts the FusionAuth login form to a domain-based login form. This type of form first asks the user for their email. FusionAuth then uses their email to determine if an OpenID Connect identity provider should be used. If an OpenID Connect provider should be used, the browser is redirected to the authorization endpoint of that identity provider. Otherwise, the password field is revealed on the form so that the user can login using FusionAuth.",
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of this OpenID Connect identity provider. This is only used for display purposes.",
			},
			"oauth2_authorization_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The top-level authorization endpoint for the OpenID Connect identity provider. You can leave this blank if you provide the issuer field, which will be used to make a request to the OpenID Connect .well-known endpoint in order to dynamically resolve the authorization endpoint. If you provide an issuer then this field will be ignored.",
			},
			"oauth2_client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The top-level client id for your Application.",
			},
			"oauth2_client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The top-level client secret to use with the OpenID Connect identity provider.",
			},
			"oauth2_client_authentication_method": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"client_secret_basic",
					"client_secret_post",
					"none",
				}, false),
				Default:     "client_secret_basic",
				Description: "The client authentication method to use with the OpenID Connect identity provider.",
			},
			"oauth2_email_claim": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "email",
				Description: "An optional configuration to modify the expected name of the claim returned by the IdP that contains the email address.",
			},
			"oauth2_issuer": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The top-level issuer URI for the OpenID Connect identity provider. If this is provided, the authorization endpoint, token endpoint and userinfo endpoint will all be resolved using the issuer URI plus /.well-known/openid-configuration.",
			},
			"oauth2_scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The top-level scope that you are requesting from the OpenID Connect identity provider.",
			},
			"oauth2_token_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The top-level token endpoint for the OpenID Connect identity provider. You can leave this blank if you provide the issuer field, which will be used to make a request to the OpenID Connect .well-known endpoint in order to dynamically resolve the token endpoint. If you provide an issuer then this field will be ignored.",
			},
			"oauth2_user_info_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The top-level userinfo endpoint for the OpenID Connect identity provider. You can leave this blank if you provide the issuer field, which will be used to make a request to the OpenID Connect .well-known endpoint in order to dynamically resolve the userinfo endpoint. If you provide an issuer then this field will be ignored.",
			},
			"post_request": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set this value equal to true if you wish to use POST bindings with this OpenID Connect identity provider. The default value of false means that a redirect binding which uses a GET request will be used.",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func buildOpenIDConnect(data *schema.ResourceData) OpenIDConnectIdentityProviderBody {
	o := fusionauth.OpenIdConnectIdentityProvider{
		ButtonImageURL: data.Get("button_image_url").(string),
		ButtonText:     data.Get("button_text").(string),
		BaseIdentityProvider: fusionauth.BaseIdentityProvider{
			Debug:      data.Get("debug").(bool),
			Enableable: buildEnableable("enabled", data),
			LambdaConfiguration: fusionauth.ProviderLambdaConfiguration{
				ReconcileId: data.Get("lambda_reconcile_id").(string),
			},
			Name:            data.Get("name").(string),
			Type:            fusionauth.IdentityProviderType_OpenIDConnect,
			LinkingStrategy: fusionauth.IdentityProviderLinkingStrategy(data.Get("linking_strategy").(string)),
		},
		PostRequest: data.Get("post_request").(bool),
		Domains:     handleStringSlice("domains", data),
		Oauth2: fusionauth.IdentityProviderOauth2Configuration{
			AuthorizationEndpoint: data.Get("oauth2_authorization_endpoint").(string),
			ClientId:              data.Get("oauth2_client_id").(string),
			ClientSecret:          data.Get("oauth2_client_secret").(string),
			ClientAuthenticationMethod: fusionauth.ClientAuthenticationMethod(
				data.Get(
					"oauth2_client_authentication_method",
				).(string)),
			EmailClaim:       data.Get("oauth2_email_claim").(string),
			Issuer:           data.Get("oauth2_issuer").(string),
			Scope:            data.Get("oauth2_scope").(string),
			TokenEndpoint:    data.Get("oauth2_token_endpoint").(string),
			UserinfoEndpoint: data.Get("oauth2_user_info_endpoint").(string),
		},
	}

	ac := buildOpenIDAppConfig("application_configuration", data)
	o.ApplicationConfiguration = ac
	return OpenIDConnectIdentityProviderBody{IdentityProvider: o}
}

func buildOpenIDAppConfig(key string, data *schema.ResourceData) map[string]interface{} {
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
		oc := OpenIDAppConfig{
			ButtonImageURL:     ac["button_image_url"].(string),
			ButtonText:         ac["button_text"].(string),
			CreateRegistration: ac["create_registration"].(bool),
			Enabled:            ac["enabled"].(bool),
			OAuth2: OAuth2AppConfig{
				ClientID:     ac["oauth2_client_id"].(string),
				ClientSecret: ac["oauth2_client_secret"].(string),
				Scope:        ac["oauth2_scope"].(string),
			},
		}
		m[aid] = oc
	}
	return m
}

func createOpenIDConnect(data *schema.ResourceData, i interface{}) error {
	o := buildOpenIDConnect(data)

	b, err := json.Marshal(o)
	if err != nil {
		return err
	}

	client := i.(Client)
	bb, err := createIdentityProvider(b, client)
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

func readOpenIDConnect(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	b, err := readIdentityProvider(data.Id(), client)
	if err != nil {
		return err
	}

	var ipb OpenIDConnectIdentityProviderBody
	_ = json.Unmarshal(b, &ipb)

	return buildResourceFromOpenIDConnect(ipb.IdentityProvider, data)
}

func buildResourceFromOpenIDConnect(o fusionauth.OpenIdConnectIdentityProvider, data *schema.ResourceData) error {
	if err := data.Set("button_image_url", o.ButtonImageURL); err != nil {
		return fmt.Errorf("idpOpenIDConnect.button_image_url: %s", err.Error())
	}
	if err := data.Set("button_text", o.ButtonText); err != nil {
		return fmt.Errorf("idpOpenIDConnect.button_text: %s", err.Error())
	}
	if err := data.Set("debug", o.Debug); err != nil {
		return fmt.Errorf("idpOpenIDConnect.debug: %s", err.Error())
	}
	if err := data.Set("domains", o.Domains); err != nil {
		return fmt.Errorf("idpOpenIDConnect.domains: %s", err.Error())
	}
	if err := data.Set("enabled", o.Enabled); err != nil {
		return fmt.Errorf("idpOpenIDConnect.enabled: %s", err.Error())
	}
	if err := data.Set("lambda_reconcile_id", o.LambdaConfiguration.ReconcileId); err != nil {
		return fmt.Errorf("idpOpenIDConnect.lambda_reconcile_id: %s", err.Error())
	}
	if err := data.Set("name", o.Name); err != nil {
		return fmt.Errorf("idpOpenIDConnect.name: %s", err.Error())
	}
	if err := data.Set("oauth2_authorization_endpoint", o.Oauth2.AuthorizationEndpoint); err != nil {
		return fmt.Errorf("idpOpenIDConnect.oauth2_authorization_endpoint: %s", err.Error())
	}
	if err := data.Set("oauth2_client_id", o.Oauth2.ClientId); err != nil {
		return fmt.Errorf("idpOpenIDConnect.oauth2_client_id: %s", err.Error())
	}
	if err := data.Set("oauth2_client_secret", o.Oauth2.ClientSecret); err != nil {
		return fmt.Errorf("idpOpenIDConnect.oauth2_client_secret: %s", err.Error())
	}
	if err := data.Set("oauth2_client_authentication_method", o.Oauth2.ClientAuthenticationMethod); err != nil {
		return fmt.Errorf("idpOpenIDConnect.oauth2_client_authentication_method: %s", err.Error())
	}
	if err := data.Set("oauth2_email_claim", o.Oauth2.EmailClaim); err != nil {
		return fmt.Errorf("idpOpenIDConnect.oauth2_email_claim: %s", err.Error())
	}
	if err := data.Set("oauth2_issuer", o.Oauth2.Issuer); err != nil {
		return fmt.Errorf("idpOpenIDConnect.oauth2_issuer: %s", err.Error())
	}
	if err := data.Set("oauth2_scope", o.Oauth2.Scope); err != nil {
		return fmt.Errorf("idpOpenIDConnect.oauth2_scope: %s", err.Error())
	}
	if err := data.Set("oauth2_token_endpoint", o.Oauth2.TokenEndpoint); err != nil {
		return fmt.Errorf("idpOpenIDConnect.oauth2_token_endpoint: %s", err.Error())
	}
	if err := data.Set("oauth2_user_info_endpoint", o.Oauth2.UserinfoEndpoint); err != nil {
		return fmt.Errorf("idpOpenIDConnect.oauth2_user_info_endpoint: %s", err.Error())
	}
	if err := data.Set("linking_strategy", o.LinkingStrategy); err != nil {
		return fmt.Errorf("idpExternalJwt.linking_strategy: %s", err.Error())
	}

	if err := data.Set("post_request", o.PostRequest); err != nil {
		return fmt.Errorf("idpOpenIDConnect.post_request: %s", err.Error())
	}

	// Since this is coming down as an interface and would end up being map[string]interface{}
	// with one of the values being map[string]interface{}
	b, _ := json.Marshal(o.ApplicationConfiguration)
	m := make(map[string]OpenIDAppConfig)
	_ = json.Unmarshal(b, &m)

	ac := make([]map[string]interface{}, 0, len(o.ApplicationConfiguration))
	for k, v := range m {
		ac = append(ac, map[string]interface{}{
			"application_id":       k,
			"button_image_url":     v.ButtonImageURL,
			"button_text":          v.ButtonText,
			"oauth2_client_id":     v.OAuth2.ClientID,
			"oauth2_client_secret": v.OAuth2.ClientSecret,
			"create_registration":  v.CreateRegistration,
			"enabled":              v.Enabled,
			"oauth2_scope":         v.OAuth2.Scope,
		})
	}
	if err := data.Set("application_configuration", ac); err != nil {
		return fmt.Errorf("idpOpenIDConnect.application_configuration: %s", err.Error())
	}

	return nil
}

func updateOpenIDConnect(data *schema.ResourceData, i interface{}) error {
	o := buildOpenIDConnect(data)

	b, err := json.Marshal(o)
	if err != nil {
		return err
	}

	client := i.(Client)
	bb, err := updateIdentityProvider(b, data.Id(), client)
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
