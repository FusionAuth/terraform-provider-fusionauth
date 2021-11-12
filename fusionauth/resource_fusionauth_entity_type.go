package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceEntityType() *schema.Resource {
	return &schema.Resource{
		Description:   "Entity Types categorize Entities. An Entity Type could be 'Device', 'API' or 'Company'.",
		CreateContext: createEntityType,
		ReadContext:   readEntityType,
		UpdateContext: updateEntityType,
		DeleteContext: deleteEntityType,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"entity_type_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				Description:  "The Id to use for the new Entity Type. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"data": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "An object that can hold any information about the Entity Type that should be persisted. Must be a JSON string.",
				DiffSuppressFunc: diffSuppressJSON,
				ValidateFunc:     validation.StringIsJSON,
			},
			"jwt_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "A block to configure JSON Web Token (JWT) options.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates if this application is using the JWT configuration defined here or the global JWT configuration defined by the Tenant. If this is false the signing algorithm configured in the Tenant will be used. If true the signing algorithm defined in this application will be used.",
							RequiredWith: []string{
								"jwt_configuration.0.access_token_key_id",
								"jwt_configuration.0.time_to_live_in_seconds",
							},
						},
						"access_token_key_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "The unique id of the signing key used to sign the access token. Required when enabled is set to true.",
							ValidateFunc: validation.IsUUID,
						},
						"time_to_live_in_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The length of time in seconds the JWT will live before it is expired and no longer valid. Required when enabled is set to true.",
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A descriptive name for the entity type (i.e. 'Customer' or 'Email_Service').",
			},
		},
	}
}

func createEntityType(_ context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	req, diags := dataToEntityTypeRequest(data)
	if diags != nil {
		return diags
	}

	client := i.(Client)
	res, faErrs, err := client.FAClient.CreateEntityType(req.EntityType.Id, req)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = checkResponse(res.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return entityTypeResponseToData(data, res)
}

func readEntityType(_ context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	client := i.(Client)
	res, faErrs, err := client.FAClient.RetrieveEntityType(data.Id())
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

	return entityTypeResponseToData(data, res)
}

func updateEntityType(_ context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	client := i.(Client)
	req, diags := dataToEntityTypeRequest(data)
	if diags != nil {
		return diags
	}

	resp, faErrs, err := client.FAClient.UpdateEntityType(data.Id(), req)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return entityTypeResponseToData(data, resp)
}

func deleteEntityType(_ context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	client := i.(Client)

	resourceID := data.Id()
	resp, faErrs, err := client.FAClient.DeleteEntityType(resourceID)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		// Entity Type successfully deleted
		data.SetId("")
		return nil
	}

	if err = checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func dataToEntityTypeRequest(data *schema.ResourceData) (req fusionauth.EntityTypeRequest, diags diag.Diagnostics) {
	// Build Entity Type Data
	resourceData, diags := jsonStringToMapStringInterface(data.Get("data").(string))

	// Build Entity Type Request
	req = fusionauth.EntityTypeRequest{
		EntityType: fusionauth.EntityType{
			Data:             resourceData,
			Id:               data.Get("entity_type_id").(string),
			JwtConfiguration: dataEntryToEntityJWTConfiguration(data),
			Name:             data.Get("name").(string),
		},
	}

	return req, diags
}

// dataEntryToEntityJWTConfiguration transforms an entities' jwt configuration
// sub-object to a fusionauth entity jwt configuration.
func dataEntryToEntityJWTConfiguration(data *schema.ResourceData) fusionauth.EntityJWTConfiguration {
	return fusionauth.EntityJWTConfiguration{
		Enableable:          buildEnableable("jwt_configuration.0.enabled", data),
		AccessTokenKeyId:    data.Get("jwt_configuration.0.access_token_key_id").(string),
		TimeToLiveInSeconds: data.Get("jwt_configuration.0.time_to_live_in_seconds").(int),
	}
}

func entityTypeResponseToData(data *schema.ResourceData, res *fusionauth.EntityTypeResponse) (diags diag.Diagnostics) {
	data.SetId(res.EntityType.Id)

	dataMapping := map[string]interface{}{
		"entity_type_id":    res.EntityType.Id,
		"data":              res.EntityType.Data,
		"jwt_configuration": flattenEntityJwtConfiguration(res.EntityType.JwtConfiguration),
		"name":              res.EntityType.Name,
	}

	return setResourceData("entity_type", data, dataMapping)
}

func flattenEntityJwtConfiguration(conf fusionauth.EntityJWTConfiguration) []interface{} {
	// jwt_configuration expects a list with a single entry.
	return []interface{}{
		map[string]interface{}{
			"enabled":                 conf.Enabled,
			"access_token_key_id":     conf.AccessTokenKeyId,
			"time_to_live_in_seconds": conf.TimeToLiveInSeconds,
		},
	}
}
