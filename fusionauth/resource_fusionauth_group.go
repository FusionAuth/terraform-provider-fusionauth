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
				ForceNew:    true,
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

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
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
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	t := resp.Group
	if err := data.Set("name", t.Name); err != nil {
		return fmt.Errorf("group.name: %s", err.Error())
	}
	if err := data.Set("tenant_id", t.TenantId); err != nil {
		return fmt.Errorf("group.tenant_id: %s", err.Error())
	}
	if err := data.Set("data", t.Data); err != nil {
		return fmt.Errorf("group.data: %s", err.Error())
	}

	var s []string

	for i := range t.Roles {
		for j := range t.Roles[i] {
			s = append(s, t.Roles[i][j].Id)
		}
	}
	if err := data.Set("role_ids", s); err != nil {
		return fmt.Errorf("group.role_ids: %s", err.Error())
	}
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
	resp, faErrs, err := client.FAClient.UpdateGroup(id, g)

	if err != nil {
		return fmt.Errorf("UpdateGroup err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	return nil
}

func deleteGroup(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteGroup(id)
	if err != nil {
		return err
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	return nil
}
