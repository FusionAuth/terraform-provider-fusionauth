package fusionauth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type IDPExternalJWTProviderBody struct {
	IdentityProvider fusionauth.ExternalJWTIdentityProvider `json:"identityProvider"`
}

type IDPExternalJWTAppConfig struct {
	CreateRegistration bool `json:"createRegistration"`
	Enabled            bool `json:"enabled"`
}

func resourceIDPExternalJWT() *schema.Resource {
	return &schema.Resource{
		CreateContext: createIDPExternalJWT,
		ReadContext:   readIDPExternalJWT,
		UpdateContext: updateIDPExternalJWT,
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
			"claim_map": {
				Type:         schema.TypeMap,
				Optional:     true,
				Description:  "A map of incoming claims to User fields, User data or Registration data. The key of the map is the incoming claim name from the configured identity provider.",
				ValidateFunc: validateClaimMap,
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if debug is enabled for this provider. When enabled, each time this provider is invoked to reconcile a login an Event Log will be created.",
			},
			"default_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "When configured this key will be used to verify the signature of the JWT when the header key defined by the headerKeyParameter property is not found in the JWT header. In most cases, the JWT header will contain the key identifier and this value will be used to resolve the correct public key or X.509 certificate to verify the signature. This assumes the public key or X.509 certificate has already been imported using the Key API or Key Master in the FusionAuth admin UI.",
				ValidateFunc: validation.IsUUID,
			},
			"domains": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "An array of domains that are managed by this Identity Provider.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if this provider is enabled. If it is false then it will be disabled globally.",
			},
			"header_key_parameter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name header claim that identifies the public key used to verify the signature. In most cases this be kid or x5t.",
			},
			// TODO: Not implemented in client
			// "keys": {
			// 	Type:        schema.TypeMap,
			// 	Optional:    true,
			// 	Description: "A map of public keys used to verify JWT signatures issued from the configured Identity Provider. The key is the key identifier, this may be referred to as the kid or for X.509 certificates the x5t claim may be used.\nThe map may contain one entry with an empty map key. When provided this key will be used when no header claim is provided to indicate which public key should be used to verify the signature. Generally speaking this will only be used when the Identity Provider issues JWTs without a key identifier in the header.",
			// },
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
				Optional:    true,
				Description: "The name of the provider. This is only used for display purposes.",
			},
			"oauth2_authorization_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The authorization endpoint for this Identity Provider. This value is not utilized by FusionAuth is only provided to be returned by the Lookup Identity Provider API response. During integration you may then utilize this value to perform the browser redirect to the OAuth2 authorize endpoint.",
			},
			"oauth2_email_claim": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "email",
				Description: "The name of the claim that contains the user's email address. This will only be used when the linking_stategy is equal to LinkByEmail or LinkByEmailForExistingUser.",
			},
			"oauth2_email_verified_claim": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "email_verified",
				Description: "The name of the claim that identities if the user's email address has been verified. When the linking_stategy is equal to LinkByEmail or LinkByEmailForExistingUser and this claim is present and the value is false a link will not be established and an error will be returned indicating a link cannot be established using an unverified email address.",
			},
			"oauth2_token_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The token endpoint for this Identity Provider. This value is not utilized by FusionAuth is only provided to be returned by the Lookup Identity Provider API response. During integration you may then utilize this value to complete the OAuth2 grant workflow.",
			},
			"oauth2_unique_id_claim": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "sub",
				Description: "The name of the claim that contains the user's unique user Id.",
			},
			"oauth2_username_claim": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "preferred_username",
				Description: "The name of the claim that contains the user's username. This will only be used when the linking_stategy is equal to LinkByUsername or LinkByUsernameForExistingUser.",
			},
			"unique_identity_claim": {
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This field is deprecated and will be removed in a future release. Prefer the use of oauth2_unique_id_claim.",
				Description: "The name of the claim that represents the unique identify of the User. This will generally be email or the name of the claim that provides the email address.",
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
			"tenant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "The unique Id of the Tenant. Providing a value creates an identity provider scoped to the specified tenant, otherwise a global identity provider is created. Tenant-scoped identity providers can only be used to authenticate in the context of the specified tenant. Global identity providers can be used with any tenant. This value cannot be updated after creation and requires recreating the resource to change.",
				ValidateFunc: validation.IsUUID,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func createIDPExternalJWT(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	o := buildIDPExternalJWT(data)

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

func readIDPExternalJWT(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	b, err := readIdentityProvider(data.Id(), client)
	if err != nil {
		if err.Error() == NotFoundError {
			data.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	var ipb IDPExternalJWTProviderBody
	_ = json.Unmarshal(b, &ipb)

	return buildResourceDataFromIDPExternalJWT(data, ipb.IdentityProvider)
}

func updateIDPExternalJWT(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	o := buildIDPExternalJWT(data)

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

func buildIDPExternalJWT(data *schema.ResourceData) IDPExternalJWTProviderBody {
	idp := fusionauth.ExternalJWTIdentityProvider{
		BaseIdentityProvider: fusionauth.BaseIdentityProvider{
			Debug:      data.Get("debug").(bool),
			Enableable: buildEnableable("enabled", data),
			LambdaConfiguration: fusionauth.ProviderLambdaConfiguration{
				ReconcileId: data.Get("lambda_reconcile_id").(string),
			},
			LinkingStrategy: fusionauth.IdentityProviderLinkingStrategy(data.Get("linking_strategy").(string)),
			Name:            data.Get("name").(string),
			TenantId:        data.Get("tenant_id").(string),
			Type:            fusionauth.IdentityProviderType_ExternalJWT,
		},
		DefaultKeyId:       data.Get("default_key_id").(string),
		Domains:            handleStringSlice("domains", data),
		HeaderKeyParameter: data.Get("header_key_parameter").(string),
		// TODO: handle keys
		Oauth2: fusionauth.IdentityProviderOauth2Configuration{
			AuthorizationEndpoint: data.Get("oauth2_authorization_endpoint").(string),
			EmailClaim:            data.Get("oauth2_email_claim").(string),
			EmailVerifiedClaim:    data.Get("oauth2_email_verified_claim").(string),
			TokenEndpoint:         data.Get("oauth2_token_endpoint").(string),
			UniqueIdClaim:         data.Get("oauth2_unique_id_claim").(string),
			UsernameClaim:         data.Get("oauth2_username_claim").(string),
		},
		UniqueIdentityClaim: data.Get("unique_identity_claim").(string),
	}

	if i, ok := data.GetOk("claim_map"); ok {
		idp.ClaimMap = intMapToStringMap(i.(map[string]interface{}))
	}

	idp.ApplicationConfiguration = buildIDPExternalJWTAppConfig("application_configuration", data)
	idp.TenantConfiguration = buildTenantConfiguration(data)
	return IDPExternalJWTProviderBody{IdentityProvider: idp}
}

func buildIDPExternalJWTAppConfig(key string, data *schema.ResourceData) map[string]interface{} {
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
		oc := IDPExternalJWTAppConfig{
			CreateRegistration: ac["create_registration"].(bool),
			Enabled:            ac["enabled"].(bool),
		}
		m[aid] = oc
	}
	return m
}

func buildResourceDataFromIDPExternalJWT(data *schema.ResourceData, res fusionauth.ExternalJWTIdentityProvider) diag.Diagnostics {
	if err := data.Set("claim_map", res.ClaimMap); err != nil {
		return diag.Errorf("idpExternalJwt.claim_map: %s", err.Error())
	}
	if err := data.Set("debug", res.Debug); err != nil {
		return diag.Errorf("idpExternalJwt.debug: %s", err.Error())
	}
	if err := data.Set("default_key_id", res.DefaultKeyId); err != nil {
		return diag.Errorf("idpExternalJwt.default_key_id: %s", err.Error())
	}
	if err := data.Set("domains", res.Domains); err != nil {
		return diag.Errorf("idpExternalJwt.domains: %s", err.Error())
	}
	if err := data.Set("enabled", res.Enabled); err != nil {
		return diag.Errorf("idpExternalJwt.enabled: %s", err.Error())
	}
	if err := data.Set("header_key_parameter", res.HeaderKeyParameter); err != nil {
		return diag.Errorf("idpExternalJwt.header_key_parameter: %s", err.Error())
	}

	// TODO: get keys
	// if err := data.Set("keys", res.); err != nil {
	// 	return diag.Errorf("idpExternalJwt.keys: %s", err.Error())
	// }
	if err := data.Set("lambda_reconcile_id", res.LambdaConfiguration.ReconcileId); err != nil {
		return diag.Errorf("idpExternalJwt.lambda_reconcile_id: %s", err.Error())
	}
	if err := data.Set("linking_strategy", res.LinkingStrategy); err != nil {
		return diag.Errorf("idpExternalJwt.linking_strategy: %s", err.Error())
	}
	if err := data.Set("name", res.Name); err != nil {
		return diag.Errorf("idpExternalJwt.name: %s", err.Error())
	}
	if err := data.Set("oauth2_authorization_endpoint", res.Oauth2.AuthorizationEndpoint); err != nil {
		return diag.Errorf("idpExternalJwt.oauth2_authorization_endpoint: %s", err.Error())
	}
	if err := data.Set("oauth2_email_claim", res.Oauth2.EmailClaim); err != nil {
		return diag.Errorf("idpExternalJwt.oauth2_email_claim: %s", err.Error())
	}
	if err := data.Set("oauth2_email_verified_claim", res.Oauth2.EmailVerifiedClaim); err != nil {
		return diag.Errorf("idpExternalJwt.oauth2_email_verified_claim: %s", err.Error())
	}
	if err := data.Set("oauth2_token_endpoint", res.Oauth2.TokenEndpoint); err != nil {
		return diag.Errorf("idpExternalJwt.oauth2_token_endpoint: %s", err.Error())
	}
	if err := data.Set("oauth2_unique_id_claim", res.Oauth2.UniqueIdClaim); err != nil {
		return diag.Errorf("idpExternalJwt.oauth2_unique_id_claim: %s", err.Error())
	}
	if err := data.Set("oauth2_username_claim", res.Oauth2.UsernameClaim); err != nil {
		return diag.Errorf("idpExternalJwt.oauth2_username_claim: %s", err.Error())
	}
	if err := data.Set("tenant_id", res.TenantId); err != nil {
		return diag.Errorf("idpExternalJwt.tenant_id: %s", err.Error())
	}
	if err := data.Set("unique_identity_claim", res.UniqueIdentityClaim); err != nil {
		return diag.Errorf("idpExternalJwt.unique_identity_claim: %s", err.Error())
	}

	// Since this is coming down as an interface and would end up being map[string]interface{}
	// with one of the values being map[string]interface{}
	b, _ := json.Marshal(res.ApplicationConfiguration)
	m := make(map[string]IDPExternalJWTAppConfig)
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
		return diag.Errorf("idpExternalJwt.application_configuration: %s", err.Error())
	}

	tc := buildTenantConfigurationResource(res.TenantConfiguration)
	if err := data.Set("tenant_configuration", tc); err != nil {
		return diag.Errorf("idpExternalJwt.tenant_configuration: %s", err.Error())
	}

	return nil
}

func validateClaimMap(i interface{}, k string) (warnings []string, errors []error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be map", k))
		return
	}

	for k, v := range m {
		s, ok := v.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected value for %q to be string", k))
			return
		}
		switch s {
		case "birthDate", "firstName", "lastName", "fullName", "middleName", "mobilePhone", "imageUrl", "timezone", "UserData", "RegistrationData":
			continue
		default:
			errors = append(errors, fmt.Errorf("valid values for %q are: %q", k, []string{"birthDate", "firstName", "lastName", "fullName", "middleName", "mobilePhone", "imageUrl", "timezone", "UserData", "RegistrationData"}))
		}
	}
	return warnings, errors
}
