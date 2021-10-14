package fusionauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newEntityGrant() *schema.Resource {
	return &schema.Resource{
		CreateContext: createEntityGrant,
		ReadContext:   readEntityGrant,
		// UpdateContext: updateEntityGrant,
		DeleteContext: deleteEntityGrant,
		Schema: map[string]*schema.Schema{
			"grant_entity_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The id of the entity that the grant is provided on",
				ValidateFunc: validation.IsUUID,
			},
			"recipient_entity_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The id of the entity receiving the permission",
				ValidateFunc: validation.IsUUID,
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Data associated with the grant",
			},
		},
	}
}

func createEntityGrant(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	grantEntityId := data.Get("grant_entity_id").(string)
	entityGrant := createEntityGrantFromData(data)

	resp, faErrs, err := client.FAClient.UpsertEntityGrant(grantEntityId, fusionauth.EntityGrantRequest{Grant: entityGrant})

	if err != nil {
		return diag.Errorf("UpsertEntityGrant err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(synthesizeEntityGrantId(grantEntityId, entityGrant.RecipientEntityId))

	return nil
}

func createEntityGrantFromData(data *schema.ResourceData) fusionauth.EntityGrant {
	return fusionauth.EntityGrant{
		// TODO: The API supports granting users this way as well.
		// Probably should select 1 or the other rather than assuming recipient_
		RecipientEntityId: data.Get("recipient_entity_id").(string),
	}
}

func readEntityGrant(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	grant_entity_id := data.Get("grant_entity_id").(string)
	recipient_entity_id := data.Get("recipient_entity_id").(string)

	resp, faErrs, err := client.FAClient.RetrieveEntityGrant(grant_entity_id, recipient_entity_id, "")

	if err != nil {
		return diag.Errorf("SearchEntityGrants", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return entityGrantToData(&resp.Grant, data)
}

func synthesizeEntityGrantId(entityId string, recipientEntityId string) string {
	return fmt.Sprintf("%s::%s", entityId, recipientEntityId)
}

func entityGrantToData(entityGrant *fusionauth.EntityGrant, data *schema.ResourceData) diag.Diagnostics {
	data.SetId(synthesizeEntityGrantId(entityGrant.Id, entityGrant.RecipientEntityId))
	if err := data.Set("grant_entity_id", entityGrant.Id); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("recipient_entity_id", entityGrant.RecipientEntityId); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// func updateEntityGrant(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
// 	client := i.(Client)

// 	return nil
// }

func deleteEntityGrant(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	id := data.Id()
	parts := strings.SplitN(id, "::", 2)

	grantEntityId := parts[0]
	recipientEntityId := parts[1]

	resp, faErrs, err := client.FAClient.DeleteEntityGrant(grantEntityId, recipientEntityId, "")
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
