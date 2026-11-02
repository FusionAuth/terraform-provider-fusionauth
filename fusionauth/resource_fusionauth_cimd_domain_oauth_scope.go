package fusionauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newCIMDDomainOAuthScope() *schema.Resource {
	return &schema.Resource{
		CreateContext: createCIMDDomainOAuthScope,
		ReadContext:   readCIMDDomainOAuthScope,
		UpdateContext: updateCIMDDomainOAuthScope,
		DeleteContext: deleteCIMDDomainOAuthScope,
		Schema: map[string]*schema.Schema{
			"cimd_domain_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The Id of the CIMD domain to which this scope belongs.",
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
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "An object that can hold any information about the OAuth Scope that should be persisted. Must be a JSON string.",
				DiffSuppressFunc: diffSuppressJSON,
				ValidateFunc:     validation.StringIsJSON,
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
				Description: "A description of the OAuth Scope for internal use.",
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
			StateContext: importCIMDDomainOAuthScope,
		},
	}
}

func buildCIMDDomainOAuthScope(data *schema.ResourceData) (fusionauth.CimdDomainOAuthScopeRequest, diag.Diagnostics) {
	resourceData, diags := jsonStringToMapStringInterface(data.Get("data").(string))
	if diags != nil {
		return fusionauth.CimdDomainOAuthScopeRequest{}, diags
	}

	return fusionauth.CimdDomainOAuthScopeRequest{
		Scope: fusionauth.CimdDomainOAuthScope{
			Data:                  resourceData,
			DefaultConsentDetail:  data.Get("default_consent_detail").(string),
			DefaultConsentMessage: data.Get("default_consent_message").(string),
			Description:           data.Get("description").(string),
			DomainId:              data.Get("cimd_domain_id").(string),
			Name:                  data.Get("name").(string),
			Required:              data.Get("required").(bool),
		},
	}, nil
}

// cimdDomainMissing reports the 400 the CIMD scope APIs return, instead of a 404, once the domain is gone.
func cimdDomainMissing(statusCode int, faErrs *fusionauth.Errors) bool {
	if statusCode != http.StatusBadRequest || faErrs == nil {
		return false
	}
	for _, e := range faErrs.FieldErrors["domainId"] {
		if e.Code == "[invalid]domainId" {
			return true
		}
	}
	return false
}

func createCIMDDomainOAuthScope(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	req, diags := buildCIMDDomainOAuthScope(data)
	if diags != nil {
		return diags
	}

	var scopeID string
	if si, ok := data.GetOk("scope_id"); ok {
		scopeID = si.(string)
	}

	resp, faErrs, err := client.FAClient.CreateCimdDomainOAuthScope(req.Scope.DomainId, scopeID, req)
	if err != nil {
		return diag.Errorf("CreateCimdDomainOAuthScope err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.Scope.Id)
	if err := data.Set("scope_id", resp.Scope.Id); err != nil {
		return diag.Errorf("scope.scope_id: %s", err.Error())
	}
	return nil
}

func importCIMDDomainOAuthScope(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	idParts := strings.Split(d.Id(), ":")
	if len(idParts) != 2 {
		return nil, fmt.Errorf("unexpected format of ID (%q), expected cimd_domain_id:scope_id", d.Id())
	}

	d.SetId(idParts[1])
	if err := d.Set("cimd_domain_id", idParts[0]); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func readCIMDDomainOAuthScope(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	resp, faErrs, err := client.FAClient.RetrieveCimdDomainOAuthScope(data.Get("cimd_domain_id").(string), data.Id())
	if err != nil {
		return diag.Errorf("RetrieveCimdDomainOAuthScope err: %v", err)
	}
	if resp.StatusCode == http.StatusNotFound || cimdDomainMissing(resp.StatusCode, faErrs) {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	s := resp.Scope
	dataMapping := map[string]interface{}{
		"cimd_domain_id":          s.DomainId,
		"data":                    s.Data,
		"default_consent_detail":  s.DefaultConsentDetail,
		"default_consent_message": s.DefaultConsentMessage,
		"description":             s.Description,
		"name":                    s.Name,
		"required":                s.Required,
		"scope_id":                s.Id,
	}

	return setResourceData("cimd_domain_oauth_scope", data, dataMapping)
}

func updateCIMDDomainOAuthScope(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	req, diags := buildCIMDDomainOAuthScope(data)
	if diags != nil {
		return diags
	}

	resp, faErrs, err := client.FAClient.UpdateCimdDomainOAuthScope(req.Scope.DomainId, data.Id(), req)
	if err != nil {
		return diag.Errorf("UpdateCimdDomainOAuthScope err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteCIMDDomainOAuthScope(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	resp, faErrs, err := client.FAClient.DeleteCimdDomainOAuthScope(data.Get("cimd_domain_id").(string), data.Id())
	if err != nil {
		return diag.Errorf("DeleteCimdDomainOAuthScope err: %v", err)
	}
	if resp.StatusCode == http.StatusNotFound || cimdDomainMissing(resp.StatusCode, faErrs) {
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
