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

func resourceAPIKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: createAPIKey,
		ReadContext:   readAPIKey,
		UpdateContext: updateAPIKey,
		DeleteContext: deleteAPIKey,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The unique Id of the Tenant. This value is required if the key is meant to be tenant scoped. Tenant scoped keys can only be used to access users and other tenant scoped objects for the specified tenant. This value is read-only once the key is created.",
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
			},
			"key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id to use for the new Form. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "API key string. When you create an API key the key is defaulted to a secure random value but the API key is simply a string, so you may call it super-secret-key if you’d like. However a long and random value makes a good API key in that it is unique and difficult to guess.",
				Sensitive:   true,
				ForceNew:    true,
			},
			"ip_access_control_list_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the IP Access Control List limiting access to this API key.",
				ValidateFunc: validation.IsUUID,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the key.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the API key.",
			},
			"permissions_endpoints": {
				Optional:    true,
				Type:        schema.TypeSet,
				Description: "Endpoint permissions for this key. Each key of the object is an endpoint, with the value being an array of the HTTP methods which can be used against the endpoint. An Empty permissions_endpoints object mean that this is a super key that authorizes this key for all the endpoints.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"/api/application",
								"/api/application/oauth-configuration",
								"/api/application/role",
								"/api/application/scope",
								"/api/application/search",
								"/api/cleanspeak/notify",
								"/api/connector",
								"/api/consent",
								"/api/consent/search",
								"/api/email/send",
								"/api/email/template",
								"/api/email/template/preview",
								"/api/email/template/search",
								"/api/entity",
								"/api/entity/grant",
								"/api/entity/grant/search",
								"/api/entity/search",
								"/api/entity/type",
								"/api/entity/type/permission",
								"/api/entity/type/search",
								"/api/form",
								"/api/form/field",
								"/api/group",
								"/api/group/member",
								"/api/group/member/search",
								"/api/group/search",
								"/api/identity-provider",
								"/api/identity-provider/link",
								"/api/identity-provider/link/pending",
								"/api/identity-provider/search",
								"/api/integration",
								"/api/ip-acl",
								"/api/ip-acl/search",
								"/api/jwt/refresh",
								"/api/jwt/vend",
								"/api/key",
								"/api/key/generate",
								"/api/key/import",
								"/api/key/search",
								"/api/lambda",
								"/api/lambda/search",
								"/api/logger",
								"/api/login",
								"/api/message/template",
								"/api/message/template/preview",
								"/api/messenger",
								"/api/passwordless/start",
								"/api/prometheus/metrics",
								"/api/reactor",
								"/api/reactor/metrics",
								"/api/report/daily-active-user",
								"/api/report/login",
								"/api/report/monthly-active-user",
								"/api/report/registration",
								"/api/report/totals",
								"/api/status",
								"/api/system-configuration",
								"/api/system/audit-log",
								"/api/system/audit-log/export",
								"/api/system/audit-log/search",
								"/api/system/event-log",
								"/api/system/event-log/search",
								"/api/system/log/export",
								"/api/system/login-record/export",
								"/api/system/login-record/search",
								"/api/system/reindex",
								"/api/system/version",
								"/api/system/webhook-attempt-log",
								"/api/system/webhook-event-log",
								"/api/system/webhook-event-log/search",
								"/api/tenant",
								"/api/tenant/search",
								"/api/theme",
								"/api/theme/search",
								"/api/two-factor/secret",
								"/api/two-factor/send",
								"/api/two-factor/start",
								"/api/two-factor/status",
								"/api/user",
								"/api/user-action",
								"/api/user-action-reason",
								"/api/user/action",
								"/api/user/bulk",
								"/api/user/change-password",
								"/api/user/comment",
								"/api/user/comment/search",
								"/api/user/consent",
								"/api/user/family",
								"/api/user/family/pending",
								"/api/user/family/request",
								"/api/user/forgot-password",
								"/api/user/import",
								"/api/user/recent-login",
								"/api/user/refresh-token/import",
								"/api/user/registration",
								"/api/user/search",
								"/api/user/two-factor",
								"/api/user/two-factor/recovery-code",
								"/api/user/verify-email",
								"/api/user/verify-registration",
								"/api/webauthn",
								"/api/webauthn/import",
								"/api/webauthn/register/start",
								"/api/webauthn/start",
								"/api/webhook",
								"/api/webhook/search",
								"/oauth2/device/user-code",
							}, false),
						},
						"delete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "HTTP DELETE Verb",
						},
						"get": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "HTTP GET Verb",
						},
						"patch": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "HTTP PATCH Verb",
						},
						"post": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "HTTP POST Verb",
						},
						"put": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "HTTP PUT Verb",
						},
					},
				},
			},
			"retrievable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "Whether the API key is retrievable. If set to false, the API key will not be retrievable after creation.",
			},
			"expiration_instant": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The expiration instant of this API key. Using an expired API key for API Authentication will result in a 401 response code.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// func createAPIKey(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
// 	client := i.(Client)
// 	ak := buildAPIKey(data)

// 	oldTenantID := client.FAClient.TenantId
// 	client.FAClient.TenantId = ak.TenantId
// 	defer func() {
// 		client.FAClient.TenantId = oldTenantID
// 	}()
// 	kid := data.Get("key_id").(string)
// 	resp, faErrs, err := client.FAClient.CreateAPIKey(kid, fusionauth.APIKeyRequest{ApiKey: ak})
// 	if err != nil {
// 		return diag.Errorf("createAPIKey errors: %v", err)
// 	}
// 	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	data.SetId(resp.ApiKey.Id)
// 	return buildResourceDataFromAPIKey(data, resp.ApiKey)
// }

func createAPIKey(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	ak := buildAPIKey(data)

	oldTenantID := client.FAClient.TenantId
	client.FAClient.TenantId = ak.TenantId
	defer func() {
		client.FAClient.TenantId = oldTenantID
	}()
	kid := data.Get("key_id").(string)
	resp, faErrs, err := createManualAPIKey(context.Background(), client.FAClient, kid, manualAPIKeyRequest{
		ApiKey: convertToManualAPIKey(ak),
	})
	if err != nil {
		return diag.Errorf("createAPIKey errors: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(resp.ApiKey.Id)
	return buildResourceDataFromAPIKey(data, convertFromManualAPIKey(resp.ApiKey))
}

func readAPIKey(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveAPIKey(id)
	if err != nil {
		return diag.Errorf("readAPIKey errors: %v", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	return buildResourceDataFromAPIKey(data, resp.ApiKey)
}

// func updateAPIKey(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
// 	client := i.(Client)
// 	ak := buildAPIKey(data)

// 	oldTenantID := client.FAClient.TenantId
// 	client.FAClient.TenantId = ak.TenantId
// 	defer func() {
// 		client.FAClient.TenantId = oldTenantID
// 	}()

// 	resp, faErrs, err := client.FAClient.UpdateAPIKey(data.Id(), fusionauth.APIKeyRequest{ApiKey: ak})
// 	if err != nil {
// 		return diag.Errorf("updateAPIKey errors: %v", err)
// 	}
// 	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
// 		return diag.FromErr(err)
// 	}

// 	data.SetId(resp.ApiKey.Id)
// 	return buildResourceDataFromAPIKey(data, resp.ApiKey)
// }

func updateAPIKey(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	ak := buildAPIKey(data)

	oldTenantID := client.FAClient.TenantId
	client.FAClient.TenantId = ak.TenantId
	defer func() {
		client.FAClient.TenantId = oldTenantID
	}()

	resp, faErrs, err := updateManualAPIKey(context.Background(), client.FAClient, data.Id(), manualAPIKeyRequest{
		ApiKey: convertToManualAPIKey(ak),
	})
	if err != nil {
		return diag.Errorf("updateAPIKey errors: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.ApiKey.Id)
	return buildResourceDataFromAPIKey(data, convertFromManualAPIKey(resp.ApiKey))
}

func deleteAPIKey(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	resp, faErrs, err := client.FAClient.DeleteAPIKey(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildAPIKey(data *schema.ResourceData) fusionauth.APIKey {
	ak := fusionauth.APIKey{
		Key:                   data.Get("key").(string),
		TenantId:              data.Get("tenant_id").(string),
		IpAccessControlListId: data.Get("ip_access_control_list_id").(string),
		MetaData: fusionauth.APIKeyMetaData{
			Attributes: map[string]string{
				"description": data.Get("description").(string),
			},
		},
		ExpirationInstant: int64(data.Get("expiration_instant").(int)),
		Name:              data.Get("name").(string),
		Retrievable:       data.Get("retrievable").(bool),
	}

	m := make(map[string][]string)
	s, ok := data.GetOk("permissions_endpoints")
	if !ok {
		return ak
	}
	set, ok := s.(*schema.Set)
	if !ok {
		return ak
	}
	l := set.List()
	for _, x := range l {
		ac := x.(map[string]interface{})
		ep := ac["endpoint"].(string)
		ss := []string{}
		if ac["delete"].(bool) {
			ss = append(ss, "DELETE")
		}
		if ac["get"].(bool) {
			ss = append(ss, "GET")
		}
		if ac["patch"].(bool) {
			ss = append(ss, "PATCH")
		}
		if ac["post"].(bool) {
			ss = append(ss, "POST")
		}
		if ac["put"].(bool) {
			ss = append(ss, "PUT")
		}
		m[ep] = ss
	}
	ak.Permissions.Endpoints = m
	return ak
}

func buildResourceDataFromAPIKey(data *schema.ResourceData, res fusionauth.APIKey) diag.Diagnostics {
	if err := data.Set("tenant_id", res.TenantId); err != nil {
		return diag.Errorf("apiKey.tenant_id: %s", err.Error())
	}
	if err := data.Set("key", res.Key); err != nil {
		return diag.Errorf("apiKey.key: %s", err.Error())
	}
	if desc, ok := res.MetaData.Attributes["description"]; ok {
		if err := data.Set("description", desc); err != nil {
			return diag.Errorf("apiKey.description: %s", err.Error())
		}
	}
	if err := data.Set("tenant_id", res.TenantId); err != nil {
		return diag.Errorf("apiKey.tenant_id: %s", err.Error())
	}

	if err := data.Set("ip_access_control_list_id", res.IpAccessControlListId); err != nil {
		return diag.Errorf("apiKey.ip_access_control_list_id: %s", err.Error())
	}

	pe := make([]map[string]interface{}, 0, len(res.Permissions.Endpoints))
	for ep := range res.Permissions.Endpoints {
		endpointPermission := map[string]interface{}{
			"delete": false,
			"get":    false,
			"patch":  false,
			"post":   false,
			"put":    false,
		}

		for _, s := range res.Permissions.Endpoints[ep] {
			method := strings.ToLower(s)
			if _, exists := endpointPermission[method]; exists {
				endpointPermission[method] = true
			}
		}

		endpointPermission["endpoint"] = ep
		pe = append(pe, endpointPermission)
	}

	if err := data.Set("permissions_endpoints", pe); err != nil {
		return diag.Errorf("apiKey.permissions_endpoints: %s", err.Error())
	}

	if err := data.Set("expiration_instant", res.ExpirationInstant); err != nil {
		return diag.Errorf("apiKey.expiration_instant: %s", err.Error())
	}

	if err := data.Set("name", res.Name); err != nil {
		return diag.Errorf("apiKey.name: %s", err.Error())
	}
	if err := data.Set("retrievable", res.Retrievable); err != nil {
		return diag.Errorf("apiKey.retrievable: %s", err.Error())
	}
	return nil
}

// Workaround for Super API Key support: https://github.com/FusionAuth/terraform-provider-fusionauth/issues/126
// Lint exclusions are to maintain consistency with the FusionAuth SDK APIKey structure
type manualAPIKey struct {
	ExpirationInstant     int64                    `json:"expirationInstant,omitempty"`
	Id                    string                   `json:"id,omitempty"` //nolint:revive,stylecheck
	InsertInstant         int64                    `json:"insertInstant,omitempty"`
	IpAccessControlListId string                   `json:"ipAccessControlListId,omitempty"` //nolint:revive,stylecheck
	Key                   string                   `json:"key,omitempty"`
	KeyManager            bool                     `json:"keyManager"`
	LastUpdateInstant     int64                    `json:"lastUpdateInstant,omitempty"`
	MetaData              manualAPIKeyMetaData     `json:"metaData,omitempty"`
	Name                  string                   `json:"name,omitempty"`
	Permissions           *manualAPIKeyPermissions `json:"permissions,omitempty"`
	Retrievable           bool                     `json:"retrievable"`
	TenantId              string                   `json:"tenantId,omitempty"` //nolint:revive,stylecheck
}

type manualAPIKeyMetaData struct {
	Attributes map[string]string `json:"attributes,omitempty"`
}

type manualAPIKeyPermissions struct {
	Endpoints map[string][]string `json:"endpoints,omitempty"`
}

type manualAPIKeyRequest struct {
	ApiKey manualAPIKey `json:"apiKey,omitempty"` //nolint:revive,stylecheck
}

type manualAPIKeyResponse struct {
	fusionauth.BaseHTTPResponse
	ApiKey  manualAPIKey   `json:"apiKey,omitempty"`  //nolint:revive,stylecheck
	ApiKeys []manualAPIKey `json:"apiKeys,omitempty"` //nolint:revive,stylecheck
}

func createManualAPIKey(ctx context.Context, client fusionauth.FusionAuthClient, keyID string, request manualAPIKeyRequest) (*manualAPIKeyResponse, *fusionauth.Errors, error) {
	return makeAPIKeyRequest(ctx, client, keyID, request, http.MethodPost)
}

func updateManualAPIKey(ctx context.Context, client fusionauth.FusionAuthClient, keyID string, request manualAPIKeyRequest) (*manualAPIKeyResponse, *fusionauth.Errors, error) {
	return makeAPIKeyRequest(ctx, client, keyID, request, http.MethodPut)
}

func convertFromManualAPIKey(manualAPIKey manualAPIKey) fusionauth.APIKey {
	apiKey := fusionauth.APIKey{
		ExpirationInstant:     manualAPIKey.ExpirationInstant,
		Id:                    manualAPIKey.Id,
		InsertInstant:         manualAPIKey.InsertInstant,
		IpAccessControlListId: manualAPIKey.IpAccessControlListId,
		Key:                   manualAPIKey.Key,
		KeyManager:            manualAPIKey.KeyManager,
		LastUpdateInstant:     manualAPIKey.LastUpdateInstant,
		MetaData: fusionauth.APIKeyMetaData{
			Attributes: manualAPIKey.MetaData.Attributes,
		},
		Name:        manualAPIKey.Name,
		Retrievable: manualAPIKey.Retrievable,
		TenantId:    manualAPIKey.TenantId,
	}

	if manualAPIKey.Permissions != nil && len(manualAPIKey.Permissions.Endpoints) > 0 {
		apiKey.Permissions = fusionauth.APIKeyPermissions{
			Endpoints: make(map[string][]string, len(manualAPIKey.Permissions.Endpoints)),
		}
		for endpoint, methods := range manualAPIKey.Permissions.Endpoints {
			apiKey.Permissions.Endpoints[endpoint] = append([]string(nil), methods...)
		}
	}

	return apiKey
}

func convertToManualAPIKey(apiKey fusionauth.APIKey) manualAPIKey {
	manualAPIKey := manualAPIKey{
		ExpirationInstant:     apiKey.ExpirationInstant,
		Id:                    apiKey.Id,
		InsertInstant:         apiKey.InsertInstant,
		IpAccessControlListId: apiKey.IpAccessControlListId,
		Key:                   apiKey.Key,
		KeyManager:            apiKey.KeyManager,
		LastUpdateInstant:     apiKey.LastUpdateInstant,
		MetaData: manualAPIKeyMetaData{
			Attributes: apiKey.MetaData.Attributes,
		},
		Name:        apiKey.Name,
		Retrievable: apiKey.Retrievable,
		TenantId:    apiKey.TenantId,
	}

	if len(apiKey.Permissions.Endpoints) > 0 {
		manualAPIKey.Permissions = &manualAPIKeyPermissions{
			Endpoints: make(map[string][]string, len(apiKey.Permissions.Endpoints)),
		}
		for endpoint, methods := range apiKey.Permissions.Endpoints {
			manualAPIKey.Permissions.Endpoints[endpoint] = append([]string(nil), methods...)
		}
	}

	return manualAPIKey
}

func makeAPIKeyRequest(ctx context.Context, client fusionauth.FusionAuthClient, keyID string, request manualAPIKeyRequest, method string) (*manualAPIKeyResponse, *fusionauth.Errors, error) {
	var resp manualAPIKeyResponse
	var errors fusionauth.Errors

	restClient := client.Start(&resp, &errors)
	if method != http.MethodGet {
		restClient.WithJSONBody(request)
	}

	err := restClient.WithUri("/api/api-key").
		WithUriSegment(keyID).
		WithMethod(method).
		Do(ctx)

	if restClient.ErrorRef == nil {
		return &resp, nil, err
	}
	return &resp, &errors, err
}
