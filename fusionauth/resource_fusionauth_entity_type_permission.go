package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceEntityTypePermission() *schema.Resource {
	return &schema.Resource{
		Description:   "Permissions are defined on an Entity Type. These are arbitrary strings which can fit the business domain. A Permission could be 'read', 'write', or 'file-lawsuit'.",
		CreateContext: createEntityTypePermission,
		ReadContext:   readEntityTypePermission,
		UpdateContext: updateEntityTypePermission,
		DeleteContext: deleteEntityTypePermission,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"entity_type_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The Id of the Entity Type.",
				ValidateFunc: validation.IsUUID,
			},
			"permission_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				Description:  "The Id to use for the new permission. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"data": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "An object that can hold any information about the Permission that should be persisted.  Must be a JSON string.",
				DiffSuppressFunc: diffSuppressJSON,
				ValidateFunc:     validation.StringIsJSON,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Permission.",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not the Permission is a default permission. A default permission is automatically granted to an entity of this type if no permissions are provided in a grant request.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Permission.",
			},
		},
	}
}

func createEntityTypePermission(_ context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	resourceReq, diags := dataToEntityTypePermissionRequest(data)
	if diags != nil {
		return diags
	}

	client := i.(Client)
	entityTypeID := data.Get("entity_type_id").(string)
	res, faErrs, err := client.FAClient.CreateEntityTypePermission(entityTypeID, resourceReq.Permission.Id, resourceReq)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = checkResponse(res.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return entityTypePermissionResponseToData(data, entityTypeID, res)
}

func readEntityTypePermission(_ context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	client := i.(Client)

	// entity type permissions are only returned via an entity type, so we need
	// to grab the entity type and drill down into the linked permissions.
	entityTypeID := data.Get("entity_type_id").(string)
	res, faErrs, err := client.FAClient.RetrieveEntityType(entityTypeID)
	if err != nil {
		return diag.FromErr(err)
	}
	if res.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err = checkResponse(res.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	// Attempt to find the linked entity type permission...
	var isFound bool
	resourceID := data.Id()
	for _, permission := range res.EntityType.Permissions {
		if permission.Id == resourceID {
			// Manually create a single permission response in order to update
			// terraform state.
			localRes := &fusionauth.EntityTypeResponse{
				Permission: permission,
			}
			diags = entityTypePermissionResponseToData(data, entityTypeID, localRes)
			isFound = true
			break
		}
	}

	if !isFound {
		// Couldn't find the permission given the entity type permission :(
		data.SetId("")
		return nil
	}

	return diags
}

func updateEntityTypePermission(_ context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	req, diags := dataToEntityTypePermissionRequest(data)
	if diags != nil {
		return diags
	}

	client := i.(Client)
	entityTypeID := data.Get("entity_type_id").(string)
	res, faErrs, err := client.FAClient.UpdateEntityTypePermission(entityTypeID, data.Id(), req)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = checkResponse(res.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return entityTypePermissionResponseToData(data, entityTypeID, res)
}

func deleteEntityTypePermission(_ context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	client := i.(Client)

	resourceID := data.Id()
	entityTypeID := data.Get("entity_type_id").(string)
	res, faErrs, err := client.FAClient.DeleteEntityTypePermission(entityTypeID, resourceID)
	if err != nil {
		return diag.FromErr(err)
	}

	if res.StatusCode == http.StatusNotFound {
		// Entity Type Permission successfully deleted
		data.SetId("")
		return nil
	}

	if err = checkResponse(res.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func dataToEntityTypePermissionRequest(data *schema.ResourceData) (req fusionauth.EntityTypeRequest, diags diag.Diagnostics) {
	resourceData, diags := jsonStringToMapStringInterface(data.Get("data").(string))

	req = fusionauth.EntityTypeRequest{
		Permission: fusionauth.EntityTypePermission{
			Data:        resourceData,
			Description: data.Get("description").(string),
			Id:          data.Get("permission_id").(string),
			IsDefault:   data.Get("is_default").(bool),
			Name:        data.Get("name").(string),
		},
	}

	return req, diags
}

func entityTypePermissionResponseToData(data *schema.ResourceData, entityTypeID string, res *fusionauth.EntityTypeResponse) (diags diag.Diagnostics) {
	data.SetId(res.Permission.Id)

	dataMapping := map[string]interface{}{
		"entity_type_id": entityTypeID,
		"permission_id":  res.Permission.Id,
		"data":           res.Permission.Data,
		"description":    res.Permission.Description,
		"is_default":     res.Permission.IsDefault,
		"name":           res.Permission.Name,
	}

	return setResourceData("entity_type_permission", data, dataMapping)
}
