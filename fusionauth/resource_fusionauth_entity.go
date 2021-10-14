package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newEntity() *schema.Resource {
	return &schema.Resource{
		CreateContext: createEntity,
		ReadContext:   readEntity,
		UpdateContext: updateEntity,
		DeleteContext: deleteEntity,
		Schema: map[string]*schema.Schema{
			"entity_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The id of the entity",
				ValidateFunc: validation.IsUUID,
			},
			"entity_type_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The id of the entity type",
				ValidateFunc: validation.IsUUID,
			},
			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The tenant the entity lives in",
				ValidateFunc: validation.IsUUID,
			},
			"client_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The client_id for this entity",
				ValidateFunc: validation.IsUUID,
			},
			"client_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The client_secret for this entity",
				ValidateFunc: validation.NoZeroValues,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The parent for this entity",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the entity",
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that can hold any information about the entity that should be persisted.",
			},
		},
	}
}

func createEntity(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	var id string
	if eid, ok := data.GetOk("entity_id"); ok {
		id = eid.(string)
	}
	entity := createEntityFromData(data)

	oldTenantId := client.FAClient.TenantId
	client.FAClient.TenantId = entity.TenantId

	defer func() {
		client.FAClient.TenantId = oldTenantId
	}()

	resp, faErrs, err := client.FAClient.CreateEntity(id, fusionauth.EntityRequest{Entity: entity})

	if err != nil {
		return diag.Errorf("CreateEntity err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return entityToData(&resp.Entity, data)
}

func readEntity(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	oldTenantId := client.FAClient.TenantId
	client.FAClient.TenantId = data.Get("tenant_id").(string)

	defer func() {
		client.FAClient.TenantId = oldTenantId
	}()

	resp, faErrs, err := client.FAClient.RetrieveEntity(id)
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
	return entityToData(&resp.Entity, data)
}

func updateEntity(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	entity := createEntityFromData(data)

	oldTenantId := client.FAClient.TenantId
	client.FAClient.TenantId = entity.TenantId

	defer func() {
		client.FAClient.TenantId = oldTenantId
	}()

	resp, faErrs, err := client.FAClient.UpdateEntity(id, fusionauth.EntityRequest{Entity: entity})
	if err != nil {
		return diag.Errorf("UpdateEntity err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func deleteEntity(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteEntity(id)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func createEntityFromData(data *schema.ResourceData) fusionauth.Entity {
	return fusionauth.Entity{
		ClientId:     data.Get("client_id").(string),
		ClientSecret: data.Get("client_secret").(string),
		Data:         data.Get("data").(map[string]interface{}),
		Name:         data.Get("name").(string),
		ParentId:     data.Get("parent_id").(string),
		TenantId:     data.Get("tenant_id").(string),
		// Is this sufficient for the API to function? There is a whole EntityType here
		Type: fusionauth.EntityType{Id: data.Get("entity_type_id").(string)},
	}
}

func entityToData(entity *fusionauth.Entity, data *schema.ResourceData) diag.Diagnostics {
	data.SetId(entity.Id)
	if err := data.Set("entity_id", entity.Id); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("entity_type_id", entity.Type.Id); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("tenant_id", entity.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("client_id", entity.ClientId); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("client_secret", entity.ClientSecret); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("parent_id", entity.ParentId); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("name", entity.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("data", entity.Data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
