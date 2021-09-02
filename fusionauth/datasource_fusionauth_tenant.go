package fusionauth

import (
	"context"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTenant() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTenantRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Tenant.",
			},
		},
	}
}

func dataSourceTenantRead(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	resp, err := client.FAClient.RetrieveTenants()
	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}
	name := data.Get("name").(string)
	var t *fusionauth.Tenant

	for i := range resp.Tenants {
		if resp.Tenants[i].Name == name {
			t = &resp.Tenants[i]
		}
	}
	if t == nil {
		return diag.Errorf("couldn't find tenant %s", name)
	}
	data.SetId(t.Id)
	return nil
}
