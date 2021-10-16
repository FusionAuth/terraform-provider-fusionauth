package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newEntityTypePermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: createEntityTypePermission,
		ReadContext:   readEntityTypePermission,
		UpdateContext: updateEntityTypePermission,
		DeleteContext: deleteEntityTypePermission,
		Schema: map[string]*schema.Schema{
			"entity_type_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The id of the entity type this permission is for",
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the permission",
				ValidateFunc: validation.NoZeroValues,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The long form description of the permission",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the permission is assigned by default during grants",
			},
		},
	}
}

func createEntityTypePermission(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	entityTypePerm := createEntityTypePermissionFromData(data)

	entityTypeId := data.Get("entity_type_id").(string)

	resp, faErrs, err := client.FAClient.CreateEntityTypePermission(entityTypeId, "", fusionauth.EntityTypeRequest{
		Permission: entityTypePerm,
	})

	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return entityTypePermissionToData(&resp.Permission, data)
}

func createEntityTypePermissionFromData(data *schema.ResourceData) fusionauth.EntityTypePermission {
	var description string
	if v, ok := data.GetOk("description"); ok {
		description = v.(string)
	}
	return fusionauth.EntityTypePermission{
		Name:        data.Get("name").(string),
		Description: description,
		IsDefault:   data.Get("is_default").(bool),
	}
}

func entityTypePermissionToData(perm *fusionauth.EntityTypePermission, data *schema.ResourceData) diag.Diagnostics {
	data.SetId(perm.Id)
	if err := data.Set("name", perm.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("description", perm.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("is_default", perm.IsDefault); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func readEntityTypePermission(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	// We need to find the permission on the main element
	entityTypeId := data.Get("entity_type_id").(string)
	resp, faErrs, err := client.FAClient.RetrieveEntityType(entityTypeId)
	if err != nil {
		return diag.Errorf("RetrieveEntityType err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	for _, perm := range resp.EntityType.Permissions {
		if perm.Id == data.Id() {
			entityTypePermissionToData(&perm, data)
			return nil
		}
	}
	data.SetId("")
	return nil
}

func updateEntityTypePermission(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()
	entityTypeId := data.Get("entity_type_id").(string)
	resp, faErrs, err := client.FAClient.UpdateEntityTypePermission(entityTypeId, id, fusionauth.EntityTypeRequest{
		Permission: createEntityTypePermissionFromData(data),
	})
	if err != nil {
		return diag.Errorf("UpdateEntityTypePermission err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func deleteEntityTypePermission(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()
	entityTypeId := data.Get("entity_type_id").(string)

	resp, faErrs, err := client.FAClient.DeleteEntityTypePermission(entityTypeId, id)
	if err != nil {
		return diag.Errorf("DeleteEntityTypePermission err: %v", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
