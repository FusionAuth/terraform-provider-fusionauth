package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func newApplicationRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: createApplicationRole,
		ReadContext:   readApplicationRole,
		UpdateContext: updateApplicationRole,
		DeleteContext: deleteApplicationRole,
		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the application that this role is for.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description for the role.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Role.",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not the Role is a default role. A default role is automatically assigned to a user during registration if no roles are provided.",
			},
			"is_super_role": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not the Role is a considered to be a super user role. This is a marker to indicate that it supersedes all other roles. FusionAuth will attempt to enforce this contract when using the web UI, it is not enforced programmatically when using the API.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func buildApplicationRole(data *schema.ResourceData) fusionauth.ApplicationRole {
	return fusionauth.ApplicationRole{
		Description: data.Get("description").(string),
		Name:        data.Get("name").(string),
		IsDefault:   data.Get("is_default").(bool),
		IsSuperRole: data.Get("is_super_role").(bool),
	}
}

func createApplicationRole(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	ar := buildApplicationRole(data)
	aid := data.Get("application_id").(string)

	resp, faErrs, err := client.FAClient.CreateApplicationRole(
		aid, "", fusionauth.ApplicationRequest{Role: ar},
	)

	if err != nil {
		return diag.Errorf("CreateApplicationRole errors: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return applicationRoleToData(data, aid, resp)
}

func readApplicationRole(_ context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	client := i.(Client)

	// application roles are only returned via an application, so we need
	// to grab the application and drill down into the linked roles.
	appID := data.Get("application_id").(string)
	resp, err := client.FAClient.RetrieveApplication(appID)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err = checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	// Attempt to find the linked entity type permission...
	var isFound bool
	resourceID := data.Id()
	for _, role := range resp.Application.Roles {
		if role.Id == resourceID {
			// Manually create a single permission response in order to update
			// terraform state.
			localRes := &fusionauth.ApplicationResponse{
				Role: role,
			}
			diags = applicationRoleToData(data, appID, localRes)
			isFound = true
			break
		}
	}

	if !isFound {
		// Couldn't find the permission given the entity type permission :(
		data.SetId("")
		return nil
	}

	return diags
}

func updateApplicationRole(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	ar := buildApplicationRole(data)
	aid := data.Get("application_id").(string)

	resp, faErrs, err := client.FAClient.UpdateApplicationRole(
		aid, data.Id(), fusionauth.ApplicationRequest{Role: ar},
	)

	if err != nil {
		return diag.Errorf("UpdateApplicationRole errors: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return applicationRoleToData(data, aid, resp)
}

func deleteApplicationRole(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	id := data.Id()
	aid := data.Get("application_id").(string)

	resp, faErrs, err := client.FAClient.DeleteApplicationRole(aid, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func applicationRoleToData(data *schema.ResourceData, applicationID string, res *fusionauth.ApplicationResponse) (diags diag.Diagnostics) {
	data.SetId(res.Role.Id)

	dataMapping := map[string]interface{}{
		"application_id": applicationID,
		"description":    res.Role.Description,
		"is_default":     res.Role.IsDefault,
		"is_super_role":  res.Role.IsSuperRole,
		"name":           res.Role.Name,
	}

	return setResourceData("application_role", data, dataMapping)
}
