package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newEntityType() *schema.Resource {
	return &schema.Resource{
		CreateContext: createEntityType,
		ReadContext:   readEntityType,
		UpdateContext: updateEntityType,
		DeleteContext: deleteEntityType,
		Schema: map[string]*schema.Schema{
			"entity_type_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The id of the entity type",
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the entity type",
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that can hold any information about the Group that should be persisted.",
			},
		},
	}
}

func readEntityType(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveEntityType(id)
	if err != nil {
		return diag.Errorf("RetrieveEntity err: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return entityTypeToData(&resp.EntityType, data)
}

func createEntityType(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	var id string
	if etid, ok := data.GetOk("entity_type_id"); ok {
		id = etid.(string)
	}

	resp, faErrs, err := client.FAClient.CreateEntityType(id, fusionauth.EntityTypeRequest{
		EntityType: createEntityTypeFromData(data),
	})

	if err != nil {
		return diag.Errorf("CreateEntityType err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return entityTypeToData(&resp.EntityType, data)
}

func updateEntityType(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.UpdateEntityType(id, fusionauth.EntityTypeRequest{
		EntityType: createEntityTypeFromData(data),
	})

	if err != nil {
		return diag.Errorf("UpdateEntity err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteEntityType(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	// Need to determine the behavior when there are existing entities bound to this type
	resp, faErrs, err := client.FAClient.DeleteEntityType(id)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func entityTypeToData(entityType *fusionauth.EntityType, data *schema.ResourceData) diag.Diagnostics {
	data.SetId(entityType.Id)
	if err := data.Set("entity_type_id", entityType.Id); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("name", entityType.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("data", entityType.Data); err != nil {
		return diag.FromErr(err)
	}
	// if err := data.Set("permission", buildPermissionSlice(entityType)); err != nil {
	// 	return diag.FromErr(err)
	// }
	return nil
}

// func buildPermissionSlice(entityType *fusionauth.EntityType) []map[string]interface{} {
// 	perms := make([]map[string]interface{}, len(entityType.Permissions))
// 	for i, value := range entityType.Permissions {
// 		perms[i] = map[string]interface{}{
// 			"permissions_id": value.Id,
// 			"name":           value.Name,
// 			"description":    value.Description,
// 			"is_default":     value.IsDefault,
// 		}
// 	}
// 	return perms
// }

func createEntityTypeFromData(data *schema.ResourceData) fusionauth.EntityType {
	return fusionauth.EntityType{
		Name: data.Get("name").(string),
		Data: data.Get("data").(map[string]interface{}),
		// Permissions: createPermissionsFromData(data),
	}
}

// func createPermissionsFromData(data *schema.ResourceData) []fusionauth.EntityTypePermission {
// 	permissionsDataRaw := data.Get("permission").([]interface{})
// 	permissions := make([]fusionauth.EntityTypePermission, len(permissionsDataRaw))
// 	for i, permRaw := range permissionsDataRaw {
// 		perm := permRaw.(map[string]interface{})

// 		// var id string
// 		// if perm["permission_id"] != nil {
// 		// 	id = perm["permission_id"].(string)
// 		// }

// 		permissions[i] = fusionauth.EntityTypePermission{
// 			Name:        perm["name"].(string),
// 			Description: perm["description"].(string),
// 		}
// 	}
// 	return permissions
// }
