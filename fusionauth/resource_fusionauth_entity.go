package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceEntity() *schema.Resource {
	return &schema.Resource{
		Description:   "Entities are instances of a single type. An Entity could be a 'nest device', an 'Email API' or 'Raviga'.",
		CreateContext: createEntity,
		ReadContext:   readEntity,
		UpdateContext: updateEntity,
		DeleteContext: deleteEntity,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "The unique Id of the tenant used to scope this API request.",
				ValidateFunc: validation.IsUUID,
			},
			"entity_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				Description:  "The Id to use for the new Entity. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"data": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "An object that can hold any information about the Entity that should be persisted. Please review the limits on data field types as you plan for and build your custom data schema.  Must be a JSON string.",
				DiffSuppressFunc: diffSuppressJSON,
				ValidateFunc:     validation.StringIsJSON,
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The OAuth 2.0 client Id. If you leave this blank during a POST, the value of the Entity Id will be used.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "The OAuth 2.0 client secret. If you leave this blank during a POST, a secure secret will be generated for you. If you leave this blank during an update, the previous value will be maintained. For both Create and Update you can provide a value and it will be stored.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A descriptive name for the Entity (i.e. 'Raviga' or 'Email Service').",
			},
			"entity_type_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Id of the Entity Type. Types are consulted for permission checks.",
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func createEntity(_ context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	resourceReq, diags := dataToEntityRequest(data)
	if diags != nil {
		return diags
	}

	client := i.(Client)
	oldTenantID := client.FAClient.TenantId
	client.FAClient.TenantId = resourceReq.Entity.TenantId
	defer func() {
		client.FAClient.TenantId = oldTenantID
	}()

	res, faErrs, err := client.FAClient.CreateEntity(resourceReq.Entity.Id, resourceReq)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = checkResponse(res.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return entityResponseToData(data, res)
}

func readEntity(_ context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	client := i.(Client)
	revertTid := clientTenantIDOverride(&client, data)
	defer revertTid()

	res, faErrs, err := client.FAClient.RetrieveEntity(data.Id())
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

	return entityResponseToData(data, res)
}

func updateEntity(_ context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	client := i.(Client)
	revertTid := clientTenantIDOverride(&client, data)
	defer revertTid()

	req, diags := dataToEntityRequest(data)
	if diags != nil {
		return diags
	}

	res, faErrs, err := client.FAClient.UpdateEntity(data.Id(), req)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = checkResponse(res.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return entityResponseToData(data, res)
}

func deleteEntity(_ context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	client := i.(Client)
	revertTid := clientTenantIDOverride(&client, data)
	defer revertTid()

	res, faErrs, err := client.FAClient.DeleteEntity(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if res.StatusCode == http.StatusNotFound {
		// Entity successfully deleted
		data.SetId("")
		return nil
	}

	if err = checkResponse(res.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func dataToEntityRequest(data *schema.ResourceData) (req fusionauth.EntityRequest, diags diag.Diagnostics) {
	resourceData, diags := jsonStringToMapStringInterface(data.Get("data").(string))

	req = fusionauth.EntityRequest{
		Entity: fusionauth.Entity{
			ClientId:     data.Get("client_id").(string),
			ClientSecret: data.Get("client_secret").(string),
			Data:         resourceData,
			Id:           data.Get("entity_id").(string),
			Name:         data.Get("name").(string),
			TenantId:     data.Get("tenant_id").(string),
			Type: fusionauth.EntityType{
				Id: data.Get("entity_type_id").(string),
			},
		},
	}

	return req, diags
}

func entityResponseToData(data *schema.ResourceData, res *fusionauth.EntityResponse) (diags diag.Diagnostics) {
	data.SetId(res.Entity.Id)

	dataMapping := map[string]interface{}{
		"client_id":      res.Entity.ClientId,
		"client_secret":  res.Entity.ClientSecret,
		"data":           res.Entity.Data,
		"entity_id":      res.Entity.Id,
		"entity_type_id": res.Entity.Type.Id,
		"name":           res.Entity.Name,
		"tenant_id":      res.Entity.TenantId,
	}

	return setResourceData("entity", data, dataMapping)
}
