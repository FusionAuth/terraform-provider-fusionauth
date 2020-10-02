package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApplicationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Application.",
			},
		},
	}
}

func dataSourceApplicationRead(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)

	resp, err := client.FAClient.RetrieveApplications()
	if err != nil {
		return err
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return err
	}
	name := data.Get("name").(string)
	var app *fusionauth.Application

	for i := range resp.Applications {
		if resp.Applications[i].Name == name {
			app = &resp.Applications[i]
		}
	}
	if app == nil {
		return fmt.Errorf("couldn't find application %s", name)
	}
	data.SetId(app.Id)
	return nil
}
