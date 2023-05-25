package fusionauth

import (
	"context"
	"encoding/json"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type SAMLIDPInitiatedIdentityProviderBody struct {
	IdentityProvider fusionauth.SAMLv2IdPInitiatedIdentityProvider `json:"identityProvider"`
}

type SAMLIDPInitiatedAppConfig struct {
	CreateRegistration bool `json:"createRegistration"`
	Enabled            bool `json:"enabled"`
}

func resourceIDPSAMLv2IdPInitiated() *schema.Resource {
	return &schema.Resource{
		CreateContext: createIDPSAMLv2IdPInitiated,
		ReadContext:   readIDPSAMLv2IdPInitiated,
		UpdateContext: updateIDPSAMLv2IdPInitiated,
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
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if debug is enabled for this provider. When enabled, each time this provider is invoked to reconcile a login an Event Log will be created.",
			},
			"email_claim": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the email claim (Attribute in the Assertion element) in the SAML response that FusionAuth uses to uniquely identity the user. If this is not set, the `use_name_for_email` flag must be true.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if this provider is enabled. If it is false then it will be disabled globally.",
			},
			"issuer": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The EntityId (unique identifier) of the SAML v2 identity provider. This value should be provided to you. Prior to 1.27.1 this value was required to be a URL.",
			},
			"key_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The id of the key stored in Key Master that is used to verify the SAML response sent back to FusionAuth from the identity provider. This key must be a verification only key or certificate (meaning that it only has a public key component).",
				ForceNew:     true,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of this SAML v2 identity provider. This is only used for display purposes.",
			},
			"use_name_for_email": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not FusionAuth will use the NameID element value as the email address of the user for reconciliation processing. If this is false, then the `email_claim` property must be set.",
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

func createIDPSAMLv2IdPInitiated(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	o := buildIDPSAMLv2IdPInitiated(data)

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

func readIDPSAMLv2IdPInitiated(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	b, err := readIdentityProvider(data.Id(), client)
	if err != nil {
		return diag.FromErr(err)
	}

	var ipb SAMLIDPInitiatedIdentityProviderBody
	_ = json.Unmarshal(b, &ipb)

	return buildResourceDataFromIDPSAMLv2IdPInitiated(data, ipb.IdentityProvider)
}

func updateIDPSAMLv2IdPInitiated(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	o := buildIDPSAMLv2IdPInitiated(data)

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

func buildIDPSAMLv2IdPInitiated(data *schema.ResourceData) SAMLIDPInitiatedIdentityProviderBody {
	s := fusionauth.SAMLv2IdPInitiatedIdentityProvider{
		BaseSAMLv2IdentityProvider: fusionauth.BaseSAMLv2IdentityProvider{
			BaseIdentityProvider: fusionauth.BaseIdentityProvider{
				Debug:      data.Get("debug").(bool),
				Enableable: buildEnableable("enabled", data),
				LambdaConfiguration: fusionauth.ProviderLambdaConfiguration{
					ReconcileId: data.Get("lambda_reconcile_id").(string),
				},
				Name:            data.Get("name").(string),
				Type:            fusionauth.IdentityProviderType_SAMLv2IdPInitiated,
				LinkingStrategy: fusionauth.IdentityProviderLinkingStrategy(data.Get("linking_strategy").(string)),
			},
			EmailClaim:        data.Get("email_claim").(string),
			KeyId:             data.Get("key_id").(string),
			UseNameIdForEmail: data.Get("use_name_for_email").(bool),
		},
		Issuer: data.Get("issuer").(string),
	}
	s.ApplicationConfiguration = buildIDPSAMLv2IdPInitiatedAppConfig("application_configuration", data)
	s.TenantConfiguration = buildTenantConfiguration(data)

	return SAMLIDPInitiatedIdentityProviderBody{IdentityProvider: s}
}

func buildResourceDataFromIDPSAMLv2IdPInitiated(data *schema.ResourceData, res fusionauth.SAMLv2IdPInitiatedIdentityProvider) diag.Diagnostics {
	if err := data.Set("debug", res.Debug); err != nil {
		return diag.Errorf("idpSAMLv2IdpInitiated.debug: %s", err.Error())
	}
	if err := data.Set("email_claim", res.EmailClaim); err != nil {
		return diag.Errorf("idpSAMLv2IdpInitiated.email_claim: %s", err.Error())
	}
	if err := data.Set("enabled", res.Enabled); err != nil {
		return diag.Errorf("idpSAMLv2IdpInitiated.enabled: %s", err.Error())
	}
	if err := data.Set("issuer", res.Issuer); err != nil {
		return diag.Errorf("idpSAMLv2IdpInitiated.issuer: %s", err.Error())
	}
	if err := data.Set("key_id", res.KeyId); err != nil {
		return diag.Errorf("idpSAMLv2IdpInitiated.key_id: %s", err.Error())
	}
	if err := data.Set("lambda_reconcile_id", res.LambdaConfiguration.ReconcileId); err != nil {
		return diag.Errorf("idpSAMLv2IdpInitiated.lambda_reconcile_id: %s", err.Error())
	}
	if err := data.Set("name", res.Name); err != nil {
		return diag.Errorf("idpSAMLv2IdpInitiated.name: %s", err.Error())
	}
	if err := data.Set("use_name_for_email", res.UseNameIdForEmail); err != nil {
		return diag.Errorf("idpSAMLv2IdpInitiated.use_name_for_email: %s", err.Error())
	}
	if err := data.Set("linking_strategy", res.LinkingStrategy); err != nil {
		return diag.Errorf("idpSAMLv2IdpInitiated.linking_strategy: %s", err.Error())
	}

	// Since this is coming down as an interface and would end up being map[string]interface{}
	// with one of the values being map[string]interface{}
	b, _ := json.Marshal(res.ApplicationConfiguration)
	m := make(map[string]SAMLIDPInitiatedAppConfig)
	_ = json.Unmarshal(b, &m)

	ac := make([]map[string]interface{}, 0, len(res.ApplicationConfiguration))
	for k, v := range m {
		ac = append(ac, map[string]interface{}{
			"application_id":      k,
			"create_registration": v.CreateRegistration,
			"enabled":             v.Enabled,
		})
	}
	if err := data.Set("application_configuration", ac); err != nil {
		return diag.Errorf("idpSAMLv2IdpInitiated.application_configuration: %s", err.Error())
	}

	tc := buildTenantConfigurationResource(res.TenantConfiguration)
	if err := data.Set("tenant_configuration", tc); err != nil {
		return diag.Errorf("idpSAMLv2IdpInitiated.tenant_configuration: %s", err.Error())
	}

	return nil
}

func buildIDPSAMLv2IdPInitiatedAppConfig(key string, data *schema.ResourceData) map[string]interface{} {
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
		oc := SAMLIDPInitiatedAppConfig{
			CreateRegistration: ac["create_registration"].(bool),
			Enabled:            ac["enabled"].(bool),
		}
		m[aid] = oc
	}
	return m
}
