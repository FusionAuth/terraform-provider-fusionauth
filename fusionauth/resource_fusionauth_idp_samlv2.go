package fusionauth

import (
	"context"
	"encoding/json"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type SAMLIdentityProviderBody struct {
	IdentityProvider fusionauth.SAMLv2IdentityProvider `json:"identityProvider"`
}

type SAMLAppConfig struct {
	ButtonImageURL     string `json:"buttonImageURL,omitempty"`
	ButtonText         string `json:"buttonText,omitempty"`
	CreateRegistration bool   `json:"createRegistration"`
	Enabled            bool   `json:"enabled"`
}

func resourceIDPSAMLv2() *schema.Resource {
	return &schema.Resource{
		CreateContext: createIDPSAMLv2,
		ReadContext:   readIDPSAMLv2,
		UpdateContext: updateIDPSAMLv2,
		DeleteContext: deleteIdentityProvider,
		Schema: map[string]*schema.Schema{
			"idp_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The ID to use for the new identity provider. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
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
					},
				},
			},
			"assertion_configuration": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				Description:      "The assertion configuration for the SAML v2 identity provider.",
				DiffSuppressFunc: suppressBlockDiff,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The destination of the assertion.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alternates": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The array of URLs that FusionAuth will accept as SAML login destinations if the identityProvider.assertionConfiguration.destination.policy setting is AllowAlternates.",
										DefaultFunc: func() (interface{}, error) {
											return []interface{}{}, nil
										},
										Elem: &schema.Schema{Type: schema.TypeString},
									},
									"policy": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The policy to use when performing a destination assertion on the SAML login request. The possible values are: Enabled, Disabled and AllowAlternates.",
										Default:     "Enabled",
										ValidateFunc: validation.StringInSlice([]string{
											"Enabled",
											"Disabled",
											"AllowAlternates",
										}, false),
									},
								},
							},
						},
						"decryption": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The decryption configuration for the SAML v2 identity provider.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Determines if FusionAuth requires encrypted assertions in SAML responses from the identity provider. When true, SAML responses from the identity provider containing unencrypted assertions will be rejected by FusionAuth.",
									},
									"key_transport_decryption_key_id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsUUID,
										Description:  "The Id of the key stored in Key Master that is used to decrypt the symmetric key on the SAML response sent to FusionAuth from the identity provider. The selected Key must contain an RSA private key. Required when `'enabled` is true.",
									},
								},
							},
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
			"email_claim": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the email claim (Attribute in the Assertion element) in the SAML response that FusionAuth uses to uniquely identity the user. If this is not set, the `use_name_for_email` flag must be true.",
			},
			"unique_id_claim": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the unique claim in the SAML response that FusionAuth uses to uniquely link the user. If this is not set, the emailClaim will be used when linking user.",
			},
			"username_claim": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the claim in the SAML response that FusionAuth uses to identify the username. If this is not set, the NameId value will be used to link a user. This property is required when linkingStrategy is set to LinkByUsername or LinkByUsernameForExistingUser.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if this provider is enabled. If it is false then it will be disabled globally.",
			},
			"idp_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The SAML v2 login page of the identity provider.",
			},
			"idp_initiated_configuration": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				Description:      "The configuration for the IdP initiated login.",
				DiffSuppressFunc: suppressBlockDiff,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines if FusionAuth will accept IdP initiated login requests from this SAMLv2 Identity Provider.",
						},
						"issuer": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The EntityId (unique identifier) of the SAML v2 identity provider. This value should be provided to you. Required when `enabled` is true.",
						},
					},
				},
			},
			"key_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The id of the key stored in Key Master that is used to verify the SAML response sent back to FusionAuth from the identity provider. This key must be a verification only key or certificate (meaning that it only has a public key component).",
			},
			"lambda_reconcile_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The id of a SAML reconcile lambda that is applied when the identity provider sends back a successful SAML response.",
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
			"login_hint_configuration": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				Description:      "The configuration for the login hint.",
				DiffSuppressFunc: suppressBlockDiff,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "When enabled and HTTP-Redirect bindings are used, FusionAuth will provide the username or email address when available to the IdP as a login hint using the configured parameter name set by the identityProvider.loginHintConfiguration.parameterName to initiate the AuthN request.",
						},
						"parameter_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "login_hint",
							Description: "The name of the parameter used to pass the username or email as login hint to the IDP when enabled, and HTTP redirect bindings are used to initiate the AuthN request.",
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of this SAML v2 identity provider. This is only used for display purposes.",
			},
			"name_id_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
				Description: "The Name Id is used to facilitate communication about a user between a Service Provider (SP) and Identity Provider (IdP). The SP can specify the preferred format in the AuthN request regarding a user. The identity Provider will attempt to honor this format in the AuthN response. When this parameter is omitted a default value of urn:oasis:names:tc:SAML:2.0:nameid-format:persistent will be used.",
			},
			"post_request": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true the authentication request will use the HTTP POST binding with the identity provider instead of the default Redirect binding which uses the HTTP GET method.",
			},
			"request_signing_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The key pair Id to use to sign the SAML request. Required when `sign_request` is true.",
				ValidateFunc: validation.IsUUID,
			},
			"sign_request": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true authentication requests sent to the identity provider will be signed.",
			},
			"use_name_for_email": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not FusionAuth will use the NameID element value as the email address of the user for reconciliation processing. If this is false, then the `email_claim` property must be set.",
			},
			"xml_signature_canonicalization_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "exclusive",
				Description: "The XML signature canonicalization method used when digesting and signing the SAML request.",
				ValidateFunc: validation.StringInSlice([]string{
					"exclusive",
					"exclusive_with_comments",
					"inclusive",
					"inclusive_with_comments",
				}, false),
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

func createIDPSAMLv2(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	o := buildIDPSAMLv2(data)

	b, err := json.Marshal(o)
	if err != nil {
		return diag.FromErr(err)
	}

	client := i.(Client)
	bb, err := createIdentityProvider(b, client, data.Get("idp_id").(string))
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
func readIDPSAMLv2(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	b, err := readIdentityProvider(data.Id(), client)
	if err != nil {
		return diag.FromErr(err)
	}

	var ipb SAMLIdentityProviderBody
	_ = json.Unmarshal(b, &ipb)

	return buildResourceDataFromIDPSAMLv2(data, ipb.IdentityProvider)
}

func updateIDPSAMLv2(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	o := buildIDPSAMLv2(data)

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

// * This is a workaround for the fact that the AssertionConfiguration.Destination.Alternates SDK field doesn't seem to handle empty slices correctly.
func handleDestinationAlternates(alternates []interface{}) []string {
	if len(alternates) == 0 {
		return nil
	}

	result := make([]string, 0)
	for _, alt := range alternates {
		if str, ok := alt.(string); ok && str != "" {
			result = append(result, str)
		}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func buildIDPSAMLv2(data *schema.ResourceData) SAMLIdentityProviderBody {
	s := fusionauth.SAMLv2IdentityProvider{
		AssertionConfiguration: fusionauth.SAMLv2AssertionConfiguration{
			Destination: fusionauth.SAMLv2DestinationAssertionConfiguration{
				Alternates: handleDestinationAlternates(data.Get("assertion_configuration.0.destination.0.alternates").([]interface{})),
				Policy:     fusionauth.SAMLv2DestinationAssertionPolicy(data.Get("assertion_configuration.0.destination.0.policy").(string)),
			},
		},
		ButtonImageURL: data.Get("button_image_url").(string),
		ButtonText:     data.Get("button_text").(string),
		BaseSAMLv2IdentityProvider: fusionauth.BaseSAMLv2IdentityProvider{
			AssertionDecryptionConfiguration: fusionauth.SAMLv2AssertionDecryptionConfiguration{
				Enableable:                  buildEnableable("assertion_configuration.0.decryption.0.enabled", data),
				KeyTransportDecryptionKeyId: data.Get("assertion_configuration.0.decryption.0.key_transport_decryption_key_id").(string),
			},
			BaseIdentityProvider: fusionauth.BaseIdentityProvider{
				Debug:      data.Get("debug").(bool),
				Enableable: buildEnableable("enabled", data),
				LambdaConfiguration: fusionauth.ProviderLambdaConfiguration{
					ReconcileId: data.Get("lambda_reconcile_id").(string),
				},
				Name:            data.Get("name").(string),
				Type:            fusionauth.IdentityProviderType_SAMLv2,
				LinkingStrategy: fusionauth.IdentityProviderLinkingStrategy(data.Get("linking_strategy").(string)),
			},
			UniqueIdClaim:     data.Get("unique_id_claim").(string),
			EmailClaim:        data.Get("email_claim").(string),
			UsernameClaim:     data.Get("username_claim").(string),
			KeyId:             data.Get("key_id").(string),
			UseNameIdForEmail: data.Get("use_name_for_email").(bool),
		},
		Domains:     handleStringSlice("domains", data),
		IdpEndpoint: data.Get("idp_endpoint").(string),
		IdpInitiatedConfiguration: fusionauth.SAMLv2IdpInitiatedConfiguration{
			Enableable: buildEnableable("idp_initiated_configuration.0.enabled", data),
			Issuer:     data.Get("idp_initiated_configuration.0.issuer").(string),
		},
		LoginHintConfiguration: fusionauth.LoginHintConfiguration{
			Enableable:    buildEnableable("login_hint_configuration.0.enabled", data),
			ParameterName: data.Get("login_hint_configuration.0.parameter_name").(string),
		},
		NameIdFormat:        data.Get("name_id_format").(string),
		PostRequest:         data.Get("post_request").(bool),
		RequestSigningKeyId: data.Get("request_signing_key").(string),
		SignRequest:         data.Get("sign_request").(bool),
		XmlSignatureC14nMethod: fusionauth.CanonicalizationMethod(
			data.Get("xml_signature_canonicalization_method").(string),
		),
	}

	s.ApplicationConfiguration = buildSAMLv2AppConfig("application_configuration", data)
	s.TenantConfiguration = buildTenantConfiguration(data)

	return SAMLIdentityProviderBody{IdentityProvider: s}
}
func buildResourceDataFromIDPSAMLv2(data *schema.ResourceData, res fusionauth.SAMLv2IdentityProvider) diag.Diagnostics {
	if err := data.Set("assertion_configuration", []map[string]interface{}{
		{
			"destination": []map[string]interface{}{
				{
					"alternates": res.AssertionConfiguration.Destination.Alternates,
					"policy":     res.AssertionConfiguration.Destination.Policy.String(),
				},
			},
			"decryption": []map[string]interface{}{
				{
					"enabled":                         res.AssertionDecryptionConfiguration.Enabled,
					"key_transport_decryption_key_id": res.AssertionDecryptionConfiguration.KeyTransportDecryptionKeyId,
				},
			},
		},
	}); err != nil {
		return diag.Errorf("idpSAMLv2.assertion_configuration: %s", err.Error())
	}

	if err := data.Set("button_image_url", res.ButtonImageURL); err != nil {
		return diag.Errorf("idpSAMLv2.button_image_url: %s", err.Error())
	}
	if err := data.Set("button_text", res.ButtonText); err != nil {
		return diag.Errorf("idpSAMLv2.button_text: %s", err.Error())
	}
	if err := data.Set("debug", res.Debug); err != nil {
		return diag.Errorf("idpSAMLv2.debug: %s", err.Error())
	}
	if err := data.Set("domains", res.Domains); err != nil {
		return diag.Errorf("idpSAMLv2.domains: %s", err.Error())
	}
	if err := data.Set("email_claim", res.EmailClaim); err != nil {
		return diag.Errorf("idpSAMLv2.email_claim: %s", err.Error())
	}
	if err := data.Set("unique_id_claim", res.UniqueIdClaim); err != nil {
		return diag.Errorf("idpSAMLv2.unique_id_claim: %s", err.Error())
	}
	if err := data.Set("username_claim", res.UsernameClaim); err != nil {
		return diag.Errorf("idpSAMLv2.username_claim: %s", err.Error())
	}
	if err := data.Set("enabled", res.Enabled); err != nil {
		return diag.Errorf("idpSAMLv2.enabled: %s", err.Error())
	}
	if err := data.Set("idp_endpoint", res.IdpEndpoint); err != nil {
		return diag.Errorf("idpSAMLv2.idp_endpoint: %s", err.Error())
	}
	if err := data.Set("idp_initiated_configuration", []map[string]interface{}{
		{
			"enabled": res.IdpInitiatedConfiguration.Enabled,
			"issuer":  res.IdpInitiatedConfiguration.Issuer,
		},
	}); err != nil {
		return diag.Errorf("idpSAMLv2.idp_initiated_configuration: %s", err.Error())
	}
	if err := data.Set("key_id", res.KeyId); err != nil {
		return diag.Errorf("idpSAMLv2.key_id: %s", err.Error())
	}
	if err := data.Set("lambda_reconcile_id", res.LambdaConfiguration.ReconcileId); err != nil {
		return diag.Errorf("idpSAMLv2.lambda_reconcile_id: %s", err.Error())
	}
	if err := data.Set("name", res.Name); err != nil {
		return diag.Errorf("idpSAMLv2.name: %s", err.Error())
	}
	if err := data.Set("name_id_format", res.NameIdFormat); err != nil {
		return diag.Errorf("idpSAMLv2.nameIdFormat: %s", err.Error())
	}
	if err := data.Set("post_request", res.PostRequest); err != nil {
		return diag.Errorf("idpSAMLv2.post_request: %s", err.Error())
	}
	if err := data.Set("request_signing_key", res.RequestSigningKeyId); err != nil {
		return diag.Errorf("idpSAMLv2.request_signing_key: %s", err.Error())
	}
	if err := data.Set("sign_request", res.SignRequest); err != nil {
		return diag.Errorf("idpSAMLv2.sign_request: %s", err.Error())
	}
	if err := data.Set("use_name_for_email", res.UseNameIdForEmail); err != nil {
		return diag.Errorf("idpSAMLv2.use_name_for_email: %s", err.Error())
	}
	if err := data.Set("xml_signature_canonicalization_method", res.XmlSignatureC14nMethod); err != nil {
		return diag.Errorf("idpSAMLv2.xml_signature_canonicalization_method: %s", err.Error())
	}
	if err := data.Set("linking_strategy", res.LinkingStrategy); err != nil {
		return diag.Errorf("idpExternalJwt.linking_strategy: %s", err.Error())
	}
	if err := data.Set("login_hint_configuration", []map[string]interface{}{
		{
			"enabled":        res.LoginHintConfiguration.Enabled,
			"parameter_name": res.LoginHintConfiguration.ParameterName,
		},
	}); err != nil {
		return diag.Errorf("idpSAMLv2.login_hint_configuration: %s", err.Error())
	}

	// Since this is coming down as an interface and would end up being map[string]interface{}
	// with one of the values being map[string]interface{}
	b, _ := json.Marshal(res.ApplicationConfiguration)
	m := make(map[string]SAMLAppConfig)
	_ = json.Unmarshal(b, &m)

	ac := make([]map[string]interface{}, 0, len(res.ApplicationConfiguration))
	for k, v := range m {
		ac = append(ac, map[string]interface{}{
			"application_id":      k,
			"button_image_url":    v.ButtonImageURL,
			"button_text":         v.ButtonText,
			"create_registration": v.CreateRegistration,
			"enabled":             v.Enabled,
		})
	}
	if err := data.Set("application_configuration", ac); err != nil {
		return diag.Errorf("idpSAMLv2.application_configuration: %s", err.Error())
	}

	tc := buildTenantConfigurationResource(res.TenantConfiguration)
	if err := data.Set("tenant_configuration", tc); err != nil {
		return diag.Errorf("idpSAMLv2.tenant_configuration: %s", err.Error())
	}

	return nil
}

func buildSAMLv2AppConfig(key string, data *schema.ResourceData) map[string]interface{} {
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
		oc := SAMLAppConfig{
			ButtonImageURL:     ac["button_image_url"].(string),
			ButtonText:         ac["button_text"].(string),
			CreateRegistration: ac["create_registration"].(bool),
			Enabled:            ac["enabled"].(bool),
		}
		m[aid] = oc
	}
	return m
}
