package fusionauth

import (
	"context"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApplicationOAuthScope() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationOAuthScopeRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Application.",
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Id of the application to which this scope belongs.",
			},
			"scope_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Id to use for the new OAuth Scope. If not specified a secure random UUID will be generated.",
			},
			"data": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The data to associate with the scope.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"default_consent_detail": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default detail to display on the OAuth consent screen if one cannot be found in the theme.",
			},
			"default_consent_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default message to display on the OAuth consent screen if one cannot be found in the theme.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the scope.",
			},
			"required": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines if the OAuth Scope is required when requested in an OAuth workflow.",
			},
		},
	}
}

func dataSourceApplicationOAuthScopeRead(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	aid := data.Get("application_id").(string)
	resp, err := client.FAClient.RetrieveApplication(aid)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}
	name := data.Get("name").(string)
	var scope *fusionauth.ApplicationOAuthScope

	for i := range resp.Application.Scopes {
		if name == resp.Application.Scopes[i].Name {
			scope = &resp.Application.Scopes[i]
		}
	}

	if scope == nil {
		return diag.Errorf("couldn't find role %s in application oauth scope %s", name, aid)
	}
	data.SetId(scope.Id)

	if err := data.Set("application_id", scope.ApplicationId); err != nil {
		return diag.Errorf("scope.application_id: %s", err.Error())
	}
	if err := data.Set("scope_id", scope.Id); err != nil {
		return diag.Errorf("scope.scope_id: %s", err.Error())
	}
	if err := data.Set("data", scope.Data); err != nil {
		return diag.Errorf("scope.data: %s", err.Error())
	}
	if err := data.Set("default_consent_detail", scope.DefaultConsentDetail); err != nil {
		return diag.Errorf("scope.default_consent_detail: %s", err.Error())
	}
	if err := data.Set("default_consent_message", scope.DefaultConsentMessage); err != nil {
		return diag.Errorf("scope.default_consent_message: %s", err.Error())
	}
	if err := data.Set("description", scope.Description); err != nil {
		return diag.Errorf("scope.description: %s", err.Error())
	}
	if err := data.Set("name", scope.Name); err != nil {
		return diag.Errorf("scope.name: %s", err.Error())
	}
	if err := data.Set("required", scope.Required); err != nil {
		return diag.Errorf("scope.required: %s", err.Error())
	}

	return nil
}
