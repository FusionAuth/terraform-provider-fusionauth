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
		return diag.Errorf("RetrieveGroup err: %v", err)
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
		return diag.Errorf("group.name: %s", err.Error())
	}
	if err := data.Set("data", t.Data); err != nil {
		return diag.Errorf("group.data: %s", err.Error())
	}

	return nil
}

func createEntityType(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	name := data.Get("name").(string)

	var id string
	if etid, ok := data.GetOk("entity_type_id"); ok {
		id = etid.(string)
	}

	resp, faErrs, err := client.FAClient.CreateEntityType(id, fusionauth.EntityTypeRequest{
		EntityType: fusionauth.EntityType{
			Name: name,
		},
	})

	if err != nil {
		return diag.Errorf("CreateEntity err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.EntityType.Id)

	return nil
}

func updateEntityType(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()
	name := data.Get("name").(string)

	resp, faErrs, err := client.FAClient.UpdateEntityType(id, fusionauth.EntityTypeRequest{
		EntityType: fusionauth.EntityType{
			Name: name,
		},
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
