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
		UpdateContext: updateEntityGrant,
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
				Optional:     true,
				ForceNew:     true,
				Description:  "The id of the entity receiving the permission",
				ValidateFunc: validation.IsUUID,
			},
			"user_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The id of the user receiving the permission",
				ValidateFunc: validation.IsUUID,
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Data associated with the grant",
			},
			"permissions": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The permissions provided by the grant",
			},
		},
	}
}

func createEntityGrant(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	grantEntityID := data.Get("grant_entity_id").(string)
	dg, resourceIDSuffix, entityGrant := createEntityGrantFromData(data)
	if dg != nil {
		return dg
	}

	resp, faErrs, err := client.FAClient.UpsertEntityGrant(grantEntityID, fusionauth.EntityGrantRequest{Grant: *entityGrant})

	if err != nil {
		return diag.Errorf("UpsertEntityGrant err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fmt.Sprintf("%s::%s", grantEntityID, resourceIDSuffix))

	return nil
}

func createEntityGrantFromData(data *schema.ResourceData) (d diag.Diagnostics, resourceIDSuffix string, entity *fusionauth.EntityGrant) {
	var perms []string
	if setPermsRaw, ok := data.GetOk("permissions"); ok {
		setPerms := setPermsRaw.([]interface{})
		for _, p := range setPerms {
			perms = append(perms, p.(string))
		}
	}

	var recipientID string
	var userID string
	if id, ok := data.GetOk("recipient_entity_id"); ok {
		recipientID = id.(string)
		resourceIDSuffix = fmt.Sprintf("entity::%s", recipientID)
	} else if id, ok := data.GetOk("user_id"); ok {
		userID = id.(string)
		resourceIDSuffix = fmt.Sprintf("user::%s", userID)
	} else {
		return diag.Errorf("Either recipient_entity_id or user_id must be set"), "", nil
	}

	return nil, resourceIDSuffix, &fusionauth.EntityGrant{
		// TODO: The API supports granting users this way as well.
		// Probably should select 1 or the other rather than assuming recipient_
		RecipientEntityId: recipientID,
		UserId:            userID,
		Permissions:       perms,
	}
}

func readEntityGrant(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	grantEntityID := data.Get("grant_entity_id").(string)
	recipientEntityID := data.Get("recipient_entity_id").(string)

	resp, faErrs, err := client.FAClient.RetrieveEntityGrant(grantEntityID, recipientEntityID, "")

	if err != nil {
		return diag.Errorf("SearchEntityGrants '%v'", err)
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

func updateEntityGrant(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	grantEntityID := data.Get("grant_entity_id").(string)
	dg, _, entityGrant := createEntityGrantFromData(data)

	if dg != nil {
		return dg
	}

	resp, faErrs, err := client.FAClient.UpsertEntityGrant(grantEntityID, fusionauth.EntityGrantRequest{Grant: *entityGrant})

	if err != nil {
		return diag.Errorf("UpsertEntityGrant err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func entityGrantToData(entityGrant *fusionauth.EntityGrant, data *schema.ResourceData) diag.Diagnostics {
	// The EntityGrant has an id, but there doesn't appear to be a way to find it later by that id
	// So, we generate a synthetic id containing the grant entity and the recipient identity to use for lookup later
	if entityGrant.RecipientEntityId != "" {
		data.SetId(fmt.Sprintf("%s::entity::%s", entityGrant.Entity.Id, entityGrant.RecipientEntityId))
	} else {
		data.SetId(fmt.Sprintf("%s::user::%s", entityGrant.Entity.Id, entityGrant.UserId))
	}

	if err := data.Set("grant_entity_id", entityGrant.Entity.Id); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("recipient_entity_id", entityGrant.RecipientEntityId); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func deleteEntityGrant(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	id := data.Id()
	parts := strings.SplitN(id, "::", 3)

	if len(parts) < 3 {
		return diag.Errorf("Inexplicable fusionauth_entity_grant id of %s does not match expected pattern", id)
	}

	grantEntityID := parts[0]
	recipientEntityID := parts[2]

	var resp *fusionauth.BaseHTTPResponse
	var faErrs *fusionauth.Errors
	var err error

	if parts[1] == "entity" {
		resp, faErrs, err = client.FAClient.DeleteEntityGrant(grantEntityID, recipientEntityID, "")
	} else {
		return diag.Errorf("Entity grant id is malformed, unrecognized switch type %s", parts[1])
	}

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
