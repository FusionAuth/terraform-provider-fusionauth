package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceApplicationRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApplicationRoleRead,
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
				ForceNew:    true,
				Description: "ID of the application that this role is for.",
			},
		},
	}
}

func dataSourceApplicationRoleRead(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	aid := data.Get("application_id").(string)
	resp, err := client.FAClient.RetrieveApplication(aid)
	if err != nil {
		return err
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return err
	}
	name := data.Get("name").(string)
	var role *fusionauth.ApplicationRole

	for i := range resp.Application.Roles {
		if name == resp.Application.Roles[i].Name {
			role = &resp.Application.Roles[i]
		}
	}

	if role == nil {
		return fmt.Errorf("couldn't find role %s in application %s", name, aid)
	}
	data.SetId(role.Id)
	return nil
}
