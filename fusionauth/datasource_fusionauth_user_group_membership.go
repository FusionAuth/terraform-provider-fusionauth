package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUserGroupMembership() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserGroupMembershipRead,
		Schema: map[string]*schema.Schema{
			"data": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "An object that can hold any information about the User for this membership that should be persisted.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Id of the Group of this membership.",
			},
			"membership_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Id of the Membership.",
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Id of the User of this membership.",
			},
		},
	}
}

func dataSourceUserGroupMembershipRead(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	gmsreq := fusionauth.GroupMemberSearchRequest{
		Search: fusionauth.GroupMemberSearchCriteria{
			GroupId: data.Get("group_id").(string),
			UserId:  data.Get("user_id").(string),
		},
	}

	resp, faErrs, err := client.FAClient.SearchGroupMembers(gmsreq)
	if err != nil {
		return diag.Errorf("RetrieveUserGroupMembership err: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	gmsresp := resp.Members
	if resp.Total == 0 {
		data.SetId("")
		return nil
	}
	if resp.Total == 1 {
		data.SetId(gmsresp[0].Id)
		if err := data.Set("data", gmsresp[0].Data); err != nil {
			return diag.Errorf("Error setting data: %v", err)
		}
		if err := data.Set("group_id", gmsresp[0].GroupId); err != nil {
			return diag.Errorf("Error setting group_id: %v", err)
		}
		if err := data.Set("membership_id", gmsresp[0].Id); err != nil {
			return diag.Errorf("Error setting membership_id: %v", err)
		}
		if err := data.Set("user_id", gmsresp[0].UserId); err != nil {
			return diag.Errorf("Error setting user_id: %v", err)
		}
		return nil
	}
	return diag.Errorf("Found %d memberships for user %s in group %s", resp.Total, data.Get("user_id").(string), data.Get("group_id").(string))
}
