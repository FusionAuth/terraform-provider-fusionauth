package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newApplicationOAuthScope() *schema.Resource {
	return &schema.Resource{
		CreateContext: createApplicationOAuthScope,
		ReadContext:   readApplicationOAuthScope,
		UpdateContext: updateApplicationOAuthScope,
		DeleteContext: deleteApplicationOAuthScope,
		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Id of the application to which this scope belongs.",
				ValidateFunc: validation.IsUUID,
			},
			"scope_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The Id to use for the new OAuth Scope. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that can hold any information about the OAuth Scope that should be persisted.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"default_consent_detail": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The default detail to display on the OAuth consent screen if one cannot be found in the theme.",
			},
			"default_consent_message": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The default message to display on the OAuth consent screen if one cannot be found in the theme.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the OAuth Scope. This is used for display purposes only.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the OAuth Scope. This is the value that will be used to request the scope in OAuth workflows.",
			},
			"required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if the OAuth Scope is required when requested in an OAuth workflow.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func buildApplicationOAuthScope(data *schema.ResourceData) fusionauth.ApplicationOAuthScopeRequest {
	var aid string
	if ai, ok := data.Get("application_id").(string); ok {
		aid = ai
	}

	var sid string
	if si, ok := data.Get("scope_id").(string); ok {
		sid = si
	}

	oas := fusionauth.ApplicationOAuthScopeRequest{
		Scope: fusionauth.ApplicationOAuthScope{
			ApplicationId:         aid,
			Id:                    sid,
			Data:                  data.Get("data").(map[string]interface{}),
			DefaultConsentMessage: data.Get("default_consent_message").(string),
			DefaultConsentDetail:  data.Get("default_consent_detail").(string),
			Description:           data.Get("description").(string),
			Name:                  data.Get("name").(string),
			Required:              data.Get("required").(bool),
		},
	}

	return oas
}

func createApplicationOAuthScope(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	oas := buildApplicationOAuthScope(data)

	var scopeID string
	if sid, ok := data.GetOk("scope_id"); ok {
		scopeID = sid.(string)
	}

	resp, faErrs, err := client.FAClient.CreateOAuthScope(oas.Scope.ApplicationId, scopeID, oas)
	if err != nil {
		return diag.Errorf("CreateApplicationOAuthScope err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.Scope.Id)
	return nil
}

func readApplicationOAuthScope(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	aid := data.Get("application_id").(string)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveOAuthScope(aid, id)
	if err != nil {
		return diag.Errorf("RetrieveApplicationOAuthScope err: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	oas := resp.Scope
	if err := data.Set("application_id", oas.ApplicationId); err != nil {
		return diag.Errorf("scope.application_id: %s", err.Error())
	}
	if err := data.Set("scope_id", oas.Id); err != nil {
		return diag.Errorf("scope.scope_id: %s", err.Error())
	}
	if err := data.Set("data", oas.Data); err != nil {
		return diag.Errorf("scope.data: %s", err.Error())
	}
	if err := data.Set("default_consent_detail", oas.DefaultConsentDetail); err != nil {
		return diag.Errorf("scope.default_consent_detail: %s", err.Error())
	}
	if err := data.Set("default_consent_message", oas.DefaultConsentMessage); err != nil {
		return diag.Errorf("scope.default_consent_message: %s", err.Error())
	}
	if err := data.Set("description", oas.Description); err != nil {
		return diag.Errorf("scope.description: %s", err.Error())
	}
	if err := data.Set("name", oas.Name); err != nil {
		return diag.Errorf("scope.name: %s", err.Error())
	}
	if err := data.Set("required", oas.Required); err != nil {
		return diag.Errorf("scope.required: %s", err.Error())
	}

	return nil
}

func updateApplicationOAuthScope(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	oas := buildApplicationOAuthScope(data)
	id := data.Id()
	resp, faErrs, err := client.FAClient.UpdateOAuthScope(oas.Scope.ApplicationId, id, oas)

	if err != nil {
		return diag.Errorf("UpdateApplicationOAuthScope err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteApplicationOAuthScope(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	aid := data.Get("application_id").(string)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteOAuthScope(aid, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
