package fusionauth

import (
	"context"
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Application.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the Tenant that owns the Application.",
			},
			"webauthn_configuration": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The WebAuthnConfiguration for the Application.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bootstrap_workflow_enabled": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"enabled": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"reauthentication_workflow_enabled": {
							Computed: true,
							Type:     schema.TypeBool,
						},
					},
				},
			},
		},
	}
}

func dataSourceApplicationRead(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	resp, err := client.FAClient.RetrieveApplications()
	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}
	name := data.Get("name").(string)
	var app *fusionauth.Application

	for i := range resp.Applications {
		if resp.Applications[i].Name == name {
			app = &resp.Applications[i]
		}
	}
	if app == nil {
		return diag.Errorf("couldn't find application %s", name)
	}
	data.SetId(app.Id)
	data.Set("tenant_id", app.TenantId)
	// Properly structure WebAuthn configuration
	webauthnConfig := []map[string]interface{}{
		{
			"bootstrap_workflow_enabled": app.WebAuthnConfiguration.BootstrapWorkflow.Enabled,
			"enabled": app.WebAuthnConfiguration.Enabled,
			"reauthentication_workflow_enabled": app.WebAuthnConfiguration.ReauthenticationWorkflow.Enabled,
		},
	}

	if err := data.Set("webauthn_configuration", webauthnConfig); err != nil {
		return diag.FromErr(fmt.Errorf("error setting webauthn_configuration: %v", err))
	}
	return nil
}
