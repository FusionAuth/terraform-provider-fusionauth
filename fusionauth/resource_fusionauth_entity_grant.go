package fusionauth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceEntityGrant() *schema.Resource {
	return &schema.Resource{
		Description:   "Entities can have Grants. Grants are relationships between a target Entity and one of two other types: a recipient Entity or a User. Grants can have zero or more Permissions associated with them.",
		CreateContext: createEntityGrant,
		ReadContext:   readEntityGrant,
		UpdateContext: updateEntityGrant,
		DeleteContext: deleteEntityGrant,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The unique Id of the tenant used to scope this API request.",
				ValidateFunc: validation.IsUUID,
			},
			"entity_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The Id of the Entity to which access is granted.",
				ValidateFunc: validation.IsUUID,
			},
			"data": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "An object that can hold any information about the Grant that should be persisted. Please review the limits on data field types as you plan for and build your custom data schema.  Must be a JSON string.",
				DiffSuppressFunc: diffSuppressJSON,
				ValidateFunc:     validation.StringIsJSON,
			},
			"permissions": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The set of permissions of this Grant.",
			},
			"recipient_entity_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "The Entity Id for which access is granted. If `recipientEntityId` is not provided, then the `userId` will be required.",
				ValidateFunc: validation.IsUUID,
				ExactlyOneOf: []string{
					"recipient_entity_id",
					"user_id",
				},
			},
			"user_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "The User Id for which access is granted. If `userId` is not provided, then the `recipientEntityId` will be required.",
				ValidateFunc: validation.IsUUID,
				ExactlyOneOf: []string{
					"recipient_entity_id",
					"user_id",
				},
			},
		},
	}
}

func createEntityGrant(ctx context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	resourceReq, diags := dataToEntityGrantRequest(data)
	if diags != nil {
		return diags
	}

	client := i.(Client)

	if iTenantID, ok := data.GetOk("tenant_id"); ok {
		// Inject Tenant ID if specified...
		oldTenantID := client.FAClient.TenantId
		client.FAClient.TenantId = iTenantID.(string)
		defer func() {
			client.FAClient.TenantId = oldTenantID
		}()
	}

	entityID := data.Get("entity_id").(string)
	res, faErrs, err := client.FAClient.UpsertEntityGrant(entityID, resourceReq)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = checkResponse(res.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(syntheticGrantID(data))
	// The entity grant api doesn't return a payload on POST/PUT, so we have to
	// search for the newly created/updated grant with a read op.
	return readEntityGrant(ctx, data, i)
}

func syntheticGrantID(data *schema.ResourceData) string {
	entityID := data.Get("entity_id").(string)
	if recipientID, ok := data.GetOk("recipient_entity_id"); ok {
		return fmt.Sprintf("%s_%s", entityID, recipientID.(string))
	}
	userID := data.Get("user_id").(string)
	return fmt.Sprintf("%s_%s", entityID, userID)
}

func readEntityGrant(_ context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	client := i.(Client)

	if iTenantID, ok := data.GetOk("tenant_id"); ok {
		// Inject Tenant ID if specified...
		oldTenantID := client.FAClient.TenantId
		client.FAClient.TenantId = iTenantID.(string)
		defer func() {
			client.FAClient.TenantId = oldTenantID
		}()
	}

	entityID := data.Get("entity_id").(string)
	recipientEntityID := data.Get("recipient_entity_id").(string)
	userID := data.Get("user_id").(string)

	res, faErrs, err := client.FAClient.RetrieveEntityGrant(entityID, recipientEntityID, userID)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(res.Grants) > 0 || res.StatusCode == http.StatusNotFound {
		// A grant was not specifically found (i.e. a list was performed) or
		// the request itself returned not found!
		data.SetId("")
		return nil
	}
	if err = checkResponse(res.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return entityGrantResponseToData(data, res)
}

func updateEntityGrant(ctx context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	resourceReq, diags := dataToEntityGrantRequest(data)
	if diags != nil {
		return diags
	}

	client := i.(Client)

	if iTenantID, ok := data.GetOk("tenant_id"); ok {
		// Inject Tenant ID if specified...
		oldTenantID := client.FAClient.TenantId
		client.FAClient.TenantId = iTenantID.(string)
		defer func() {
			client.FAClient.TenantId = oldTenantID
		}()
	}

	entityID := data.Get("entity_id").(string)
	res, faErrs, err := client.FAClient.UpsertEntityGrant(entityID, resourceReq)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = checkResponse(res.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	// The entity grant api doesn't return a payload on POST/PUT, so we have to
	// search for the newly created/updated grant with a read op.
	return readEntityGrant(ctx, data, i)
}

func deleteEntityGrant(_ context.Context, data *schema.ResourceData, i interface{}) (diags diag.Diagnostics) {
	client := i.(Client)

	entityID := data.Get("entity_id").(string)
	recipientEntityID := data.Get("recipient_entity_id").(string)
	userID := data.Get("user_id").(string)

	resp, faErrs, err := client.FAClient.DeleteEntityGrant(entityID, recipientEntityID, userID)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		// Entity Grant successfully deleted
		data.SetId("")
		return nil
	}

	if err = checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func dataToEntityGrantRequest(data *schema.ResourceData) (req fusionauth.EntityGrantRequest, diags diag.Diagnostics) {
	resourceData, diags := jsonStringToMapStringInterface(data.Get("data").(string))

	req = fusionauth.EntityGrantRequest{
		Grant: fusionauth.EntityGrant{
			Data:              resourceData,
			RecipientEntityId: data.Get("recipient_entity_id").(string),
			UserId:            data.Get("user_id").(string),
			Permissions:       handleStringSlice("permissions", data),
		},
	}

	return req, diags
}

func entityGrantResponseToData(data *schema.ResourceData, res *fusionauth.EntityGrantResponse) (diags diag.Diagnostics) {
	data.SetId(res.Grant.Id)

	dataMapping := map[string]interface{}{
		"tenant_id":           res.Grant.Entity.TenantId,
		"entity_id":           res.Grant.Entity.Id,
		"data":                res.Grant.Data,
		"permissions":         res.Grant.Permissions,
		"recipient_entity_id": res.Grant.RecipientEntityId,
		"user_id":             res.Grant.UserId,
	}

	return setResourceData("entity_grant", data, dataMapping)
}
