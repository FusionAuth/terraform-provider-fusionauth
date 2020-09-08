package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func newApplicationRole() *schema.Resource {
	return &schema.Resource{
		Create: createApplicationRole,
		Read:   readApplicationRole,
		Update: updateApplicationRole,
		Delete: deleteApplicationRole,
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
			State: schema.ImportStatePassthrough,
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

func createApplicationRole(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)

	ar := buildApplicationRole(data)
	aid := data.Get("application_id").(string)

	resp, faErrs, err := client.FAClient.CreateApplicationRole(
		aid, "", fusionauth.ApplicationRequest{Role: ar},
	)

	if err != nil {
		return fmt.Errorf("CreateApplicationRole errors: %v", err)
	}

	if faErrs != nil {
		return fmt.Errorf("CreateApplicationRole errors: %v", faErrs)
	}

	data.SetId(resp.Role.Id)

	return nil
}

func readApplicationRole(data *schema.ResourceData, i interface{}) error {
	return nil
}

func updateApplicationRole(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)

	ar := buildApplicationRole(data)
	aid := data.Get("application_id").(string)

	_, faErrs, err := client.FAClient.UpdateApplicationRole(
		aid, data.Id(), fusionauth.ApplicationRequest{Role: ar},
	)

	if err != nil {
		return fmt.Errorf("CreateApplicationRole errors: %v", err)
	}

	if faErrs != nil {
		return fmt.Errorf("CreateApplicationRole errors: %v", faErrs)
	}

	return nil
}

func deleteApplicationRole(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)

	id := data.Id()
	aid := data.Get("application_id").(string)

	_, faErrs, err := client.FAClient.DeleteApplicationRole(aid, id)
	if err != nil {
		return err
	}

	if faErrs != nil {
		return fmt.Errorf("DeleteApplicationRole errors: %v", faErrs)
	}

	return nil
}
