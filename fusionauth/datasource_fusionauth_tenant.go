package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTenant() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTenantRead,
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

func dataSourceTenantRead(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)

	resp, err := client.FAClient.RetrieveTenants()
	if err != nil {
		return err
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return err
	}
	name := data.Get("name").(string)
	var t *fusionauth.Tenant

	for i := range resp.Tenants {
		if resp.Tenants[i].Name == name {
			t = &resp.Tenants[i]
		}
	}
	if t == nil {
		return fmt.Errorf("couldn't find tenant %s", name)
	}
	data.SetId(t.Id)
	return nil
}
