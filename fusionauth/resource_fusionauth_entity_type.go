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
			// "permission": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"permissions_id": {
			// 				Type:         schema.TypeString,
			// 				Optional:     true,
			// 				Computed:     true,
			// 				ValidateFunc: validation.IsUUID,
			// 				Description:  "The permissions id",
			// 			},
			// 			"name": {
			// 				Type:        schema.TypeString,
			// 				Required:    true,
			// 				Description: "The permission name",
			// 			},
			// 			"description": {
			// 				Type:        schema.TypeString,
			// 				Optional:    true,
			// 				Description: "The permission description",
			// 			},
			// 			"is_default": {
			// 				Type:        schema.TypeBool,
			// 				Default:     false,
			// 				Optional:    true,
			// 				Description: "Should the permission be applied by default",
			// 			},
			// 		},
			// 	},
			// },
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

	t := resp.EntityType
	if err := data.Set("name", t.Name); err != nil {
		return diag.Errorf("entity.name: %s", err.Error())
	}
	if err := data.Set("data", t.Data); err != nil {
		return diag.Errorf("entity.data: %s", err.Error())
	}
	return nil
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

	data.SetId(resp.EntityType.Id)
	if err := data.Set("name", resp.EntityType.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("data", resp.EntityType.Data); err != nil {
		return diag.FromErr(err)
	}

	return nil
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
	return nil
}

func createEntityTypeFromData(data *schema.ResourceData) fusionauth.EntityType {
	return fusionauth.EntityType{
		Name: data.Get("name").(string),
		Data: data.Get("data").(map[string]interface{}),
	}
}
