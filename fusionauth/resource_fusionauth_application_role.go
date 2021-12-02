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

	data.SetId(resp.Role.Id)

	return nil
}

func readApplicationRole(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return nil
}

func updateApplicationRole(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	ar := buildApplicationRole(data)
	aid := data.Get("application_id").(string)

	resp, faErrs, err := client.FAClient.UpdateApplicationRole(
		aid, data.Id(), fusionauth.ApplicationRequest{Role: ar},
	)

	if err != nil {
		return diag.Errorf("CreateApplicationRole errors: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteApplicationRole(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	id := data.Id()
	aid := data.Get("application_id").(string)

	resp, faErrs, err := client.FAClient.DeleteApplicationRole(aid, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
