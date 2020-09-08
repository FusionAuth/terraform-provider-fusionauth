package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func newGroup() *schema.Resource {
	return &schema.Resource{
		Create: createGroup,
		Read:   readGroup,
		Update: updateGroup,
		Delete: deleteGroup,
		Schema: map[string]*schema.Schema{
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that can hold any information about the Group that should be persisted.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Group.",
			},
			"role_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The Application Roles to assign to this group.",
			},
			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The unique Id of the tenant used to scope this API request.",
				ValidateFunc: validation.IsUUID,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func buildGroup(data *schema.ResourceData) fusionauth.GroupRequest {
	g := fusionauth.GroupRequest{
		Group: fusionauth.Group{
			Data:     nil,
			Name:     data.Get("name").(string),
			TenantId: data.Get("tenant_id").(string),
		},
		RoleIds: handleStringSlice("role_ids", data),
	}

	return g
}

func createGroup(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	g := buildGroup(data)
	oldTenantID := client.FAClient.TenantId
	client.FAClient.TenantId = g.Group.TenantId
	defer func() {
		client.FAClient.TenantId = oldTenantID
	}()
	resp, faErrs, err := client.FAClient.CreateGroup("", g)

	if err != nil {
		return fmt.Errorf("CreateGroup err: %v", err)
	}

	if faErrs != nil {
		return fmt.Errorf("CreateGroup errors: %v", faErrs)
	}

	data.SetId(resp.Group.Id)
	return nil
}

func readGroup(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveGroup(id)
	if err != nil {
		return fmt.Errorf("RetrieveGroup err: %v", err)
	}
	if faErrs != nil {
		return fmt.Errorf("RetrieveGroup errors: %v", faErrs)
	}

	t := resp.Group
	_ = data.Set("name", t.Name)
	_ = data.Set("tenant_id", t.TenantId)
	_ = data.Set("data", t.Data)

	var s []string

	for i := range t.Roles {
		for j := range t.Roles[i] {
			s = append(s, t.Roles[i][j].Id)
		}
	}
	_ = data.Set("role_ids", s)
	return nil
}

func updateGroup(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	g := buildGroup(data)
	id := data.Id()
	oldTenantID := client.FAClient.TenantId
	client.FAClient.TenantId = g.Group.TenantId
	defer func() {
		client.FAClient.TenantId = oldTenantID
	}()
	_, faErrs, err := client.FAClient.UpdateGroup(id, g)

	if err != nil {
		return fmt.Errorf("UpdateGroup err: %v", err)
	}

	if faErrs != nil {
		_, faErrs, err := client.FAClient.DeleteGroup(id)
		if err != nil {
			return err
		}

		if faErrs != nil {
			return fmt.Errorf("DeleteGroup errors: %v", faErrs)
		}
		_, faErrs, err = client.FAClient.CreateGroup(id, g)

		if err != nil {
			return fmt.Errorf("CreateGroup err: %v", err)
		}

		if faErrs != nil {
			return fmt.Errorf("CreateGroup errors: %v", faErrs)
		}
	}

	return nil
}

func deleteGroup(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	_, faErrs, err := client.FAClient.DeleteGroup(id)
	if err != nil {
		return err
	}

	if faErrs != nil {
		return fmt.Errorf("DeleteGroup errors: %v", faErrs)
	}

	return nil
}
