package fusionauth

import (
	"context"
	"net/http"
	"strings"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newCIMDDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: createCIMDDomain,
		ReadContext:   readCIMDDomain,
		UpdateContext: updateCIMDDomain,
		DeleteContext: deleteCIMDDomain,
		Schema: map[string]*schema.Schema{
			"cimd_domain_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The Id to use for the new CIMD domain. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The Id of the Tenant this CIMD domain belongs to. This value cannot be changed after creation.",
				ValidateFunc: validation.IsUUID,
			},
			"allow_localhost_redirect_uris": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether CIMD clients on this domain may use a `localhost` redirect URI, such as `http://localhost:3334/callback`. The URI must also appear in the client's metadata document redirect_uris. When `false`, FusionAuth rejects any `localhost` redirect URI, but still accepts loopback IP addresses such as `http://127.0.0.1/callback`. RFC 8252 Section 7.3 recommends loopback IP addresses over `localhost`, because the name can resolve to an interface other than loopback. Enable this only for CIMD clients that require it.",
			},
			"authorized_resource_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The resource URIs that clients registered on this CIMD domain may request using the RFC 8707 `resource` parameter. Every CIMD client must send exactly one `resource` parameter, and FusionAuth rejects the request with `invalid_target` unless it appears in this list. Clients on a domain with an empty list cannot log in.",
			},
			"consent_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     fusionauth.OAuthScopeConsentMode_AlwaysPrompt.String(),
				Description: "Controls whether a user is prompted to consent to the OAuth scopes requested by this domain's clients. The possible values are: `AlwaysPrompt` - Always prompt the user for consent. `RememberDecision` - Remember previous consents; only prompt if the choice expires or if the requested or required scopes have changed. The duration of this persisted choice is controlled by the Tenant's `external_identifier_configuration.remember_oauth_scope_consent_choice_time_to_live_in_seconds` value. FusionAuth remembers a decision separately for each client registration and `resource`. Changing the consent mode to `AlwaysPrompt` deletes every remembered decision for the domain's clients. A value of `NeverPrompt` will be rejected. A CIMD client is always treated as a third-party client, and once a consent decision is skipped entirely there is no way to re-prompt the client's user.",
				ValidateFunc: validation.StringInSlice([]string{
					fusionauth.OAuthScopeConsentMode_AlwaysPrompt.String(),
					fusionauth.OAuthScopeConsentMode_RememberDecision.String(),
				}, false),
			},
			"data": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "An object that can hold any information about the CIMD domain that should be persisted. Must be a JSON string.",
				DiffSuppressFunc: diffSuppressJSON,
				ValidateFunc:     validation.StringIsJSON,
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether FusionAuth writes the standard OAuth debug Event Log entries, such as the authorization code exchange, for this domain's clients. The authorization code exchange entry lists the requested resources. FusionAuth also writes a Debug entry when it creates a client registration, refreshes one after its cache expires, or adds a new resource to an existing one. FusionAuth always writes an Error entry for a failed client registration, whether or not this is enabled. It always writes an Information entry when it rate limits a registration.",
			},
			"domain": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: diffSuppressCIMDDomain,
				Description:      "The origin that CIMD clients register from, given as `scheme://host[:port]`. FusionAuth normalizes this value before storing it and uses it to validate the `client_id` URL presented on an OAuth request. Use the literal value `*` to accept registrations from any origin. A Tenant may contain either a single `*` domain or any number of exact-match domains, but never both. The origin must use `https`, and its host must contain at least one `.`. A path, query string, fragment, or user info component is not permitted, and neither are IP addresses or single-label hosts such as `localhost`. Subdomain wildcards such as `https://*.example.com` are not supported. This value cannot be changed after creation.",
			},
			"grace_period_in_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				Description:  "How long, in seconds, FusionAuth will continue to honor an expired cached metadata document when a re-fetch of that document fails. Must be greater than or equal to `0`.",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"jwt_populate_lambda_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The unique Id of the lambda used to populate the claims of the access token and Id token issued to this domain's clients. A single lambda is applied to both tokens. The lambda must be a JWT populate lambda.",
				ValidateFunc: validation.IsUUID,
			},
			"max_cache_ttl_in_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      86400,
				Description:  "The ceiling, in seconds, applied to the `Cache-Control: max-age` value returned when fetching a client's metadata document. When the document specifies a smaller value, that smaller value is used instead. The `Expires` header is not honored. Must be at least `900`, the non-configurable 15 minute minimum cache duration.",
				ValidateFunc: validation.IntAtLeast(900),
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the CIMD domain, used to identify it in the FusionAuth admin UI. Leave this blank to default to the value of `domain`.",
			},
			"provided_scope_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     newCIMDDomainProvidedScopePolicy(),
			},
			"registration_lambda_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The unique Id of the CIMD Client Registration lambda used to accept, reject, or modify a client registration for this domain. The lambda must be a CIMD Client Registration lambda. When this is not set, no registration lambda runs and every client from this origin is allowed to register. When `domain` is `*`, a registration lambda is the only way to filter which clients may register.",
				ValidateFunc: validation.IsUUID,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func newCIMDDomainProvidedScopePolicy() *schema.Resource {
	requireable := func(scope string) *schema.Schema {
		return &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     true,
						Description: "Whether the `" + scope + "` OAuth scope provided by FusionAuth is enabled for this CIMD domain. A disabled scope is treated as an unknown scope when a client requests it.",
					},
					"required": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "Whether consent to the `" + scope + "` OAuth scope provided by FusionAuth is required for this CIMD domain when present on the OAuth request.",
					},
				},
			},
		}
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"address": requireable("address"),
			"email":   requireable("email"),
			"phone":   requireable("phone"),
			"profile": requireable("profile"),
		},
	}
}

// diffSuppressCIMDDomain ignores the case and default https port that FusionAuth strips when normalizing a domain.
func diffSuppressCIMDDomain(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	return normalizeCIMDDomain(oldValue) == normalizeCIMDDomain(newValue)
}

func normalizeCIMDDomain(domain string) string {
	return strings.TrimSuffix(strings.ToLower(domain), ":443")
}

// buildCIMDDomainRequireable defaults an omitted provided scope to enabled, matching the server default.
func buildCIMDDomainRequireable(key string, data *schema.ResourceData) fusionauth.Requirable {
	if len(data.Get(key).([]interface{})) == 0 {
		return fusionauth.Requirable{Enableable: fusionauth.Enableable{Enabled: true}}
	}
	return buildRequireable(key, data)
}

func buildCIMDDomain(data *schema.ResourceData) (fusionauth.CimdDomainRequest, diag.Diagnostics) {
	resourceData, diags := jsonStringToMapStringInterface(data.Get("data").(string))
	if diags != nil {
		return fusionauth.CimdDomainRequest{}, diags
	}

	return fusionauth.CimdDomainRequest{
		Domain: fusionauth.CimdDomain{
			AllowLocalhostRedirectUris: data.Get("allow_localhost_redirect_uris").(bool),
			AuthorizedResourceUris:     handleStringSliceFromList(data.Get("authorized_resource_uris").([]interface{})),
			ConsentMode:                fusionauth.OAuthScopeConsentMode(data.Get("consent_mode").(string)),
			Data:                       resourceData,
			Debug:                      data.Get("debug").(bool),
			Domain:                     data.Get("domain").(string),
			GracePeriodInSeconds:       data.Get("grace_period_in_seconds").(int),
			JwtPopulateLambdaId:        data.Get("jwt_populate_lambda_id").(string),
			MaxCacheTtlInSeconds:       data.Get("max_cache_ttl_in_seconds").(int),
			Name:                       data.Get("name").(string),
			ProvidedScopePolicy: fusionauth.ProvidedScopePolicy{
				Address: buildCIMDDomainRequireable("provided_scope_policy.0.address", data),
				Email:   buildCIMDDomainRequireable("provided_scope_policy.0.email", data),
				Phone:   buildCIMDDomainRequireable("provided_scope_policy.0.phone", data),
				Profile: buildCIMDDomainRequireable("provided_scope_policy.0.profile", data),
			},
			RegistrationLambdaId: data.Get("registration_lambda_id").(string),
			TenantId:             data.Get("tenant_id").(string),
		},
	}, nil
}

func createCIMDDomain(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	revertTid := clientTenantIDOverride(&client, data)
	defer revertTid()

	req, diags := buildCIMDDomain(data)
	if diags != nil {
		return diags
	}

	var domainID string
	if di, ok := data.GetOk("cimd_domain_id"); ok {
		domainID = di.(string)
	}

	resp, faErrs, err := client.FAClient.CreateCimdDomain(domainID, req)
	if err != nil {
		return diag.Errorf("CreateCimdDomain err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.Domain.Id)
	return buildResourceDataFromCIMDDomain(data, resp.Domain)
}

func readCIMDDomain(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	revertTid := clientTenantIDOverride(&client, data)
	defer revertTid()

	resp, faErrs, err := client.FAClient.RetrieveCimdDomain(data.Id())
	if err != nil {
		return diag.Errorf("RetrieveCimdDomain err: %v", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return buildResourceDataFromCIMDDomain(data, resp.Domain)
}

func updateCIMDDomain(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	revertTid := clientTenantIDOverride(&client, data)
	defer revertTid()

	req, diags := buildCIMDDomain(data)
	if diags != nil {
		return diags
	}

	resp, faErrs, err := client.FAClient.UpdateCimdDomain(data.Id(), req)
	if err != nil {
		return diag.Errorf("UpdateCimdDomain err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return buildResourceDataFromCIMDDomain(data, resp.Domain)
}

func deleteCIMDDomain(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	revertTid := clientTenantIDOverride(&client, data)
	defer revertTid()

	resp, faErrs, err := client.FAClient.DeleteCimdDomain(data.Id())
	if err != nil {
		return diag.Errorf("DeleteCimdDomain err: %v", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildResourceDataFromCIMDDomain(data *schema.ResourceData, d fusionauth.CimdDomain) diag.Diagnostics {
	requireable := func(r fusionauth.Requirable) []map[string]interface{} {
		return []map[string]interface{}{{"enabled": r.Enabled, "required": r.Required}}
	}

	dataMapping := map[string]interface{}{
		"allow_localhost_redirect_uris": d.AllowLocalhostRedirectUris,
		"authorized_resource_uris":      d.AuthorizedResourceUris,
		"cimd_domain_id":                d.Id,
		"consent_mode":                  d.ConsentMode,
		"data":                          d.Data,
		"debug":                         d.Debug,
		"domain":                        d.Domain,
		"grace_period_in_seconds":       d.GracePeriodInSeconds,
		"jwt_populate_lambda_id":        d.JwtPopulateLambdaId,
		"max_cache_ttl_in_seconds":      d.MaxCacheTtlInSeconds,
		"name":                          d.Name,
		"provided_scope_policy": []map[string]interface{}{
			{
				"address": requireable(d.ProvidedScopePolicy.Address),
				"email":   requireable(d.ProvidedScopePolicy.Email),
				"phone":   requireable(d.ProvidedScopePolicy.Phone),
				"profile": requireable(d.ProvidedScopePolicy.Profile),
			},
		},
		"registration_lambda_id": d.RegistrationLambdaId,
		"tenant_id":              d.TenantId,
	}

	return setResourceData("cimd_domain", data, dataMapping)
}
