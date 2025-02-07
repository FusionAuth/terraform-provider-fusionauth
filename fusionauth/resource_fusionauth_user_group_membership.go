package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newUserGroupMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: createUserGroupMembership,
		ReadContext:   readUserGroupMembership,
		UpdateContext: updateUserGroupMembership,
		DeleteContext: deleteUserGroupMembership,
		Schema: map[string]*schema.Schema{
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that can hold any information about the User for this membership that should be persisted.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"group_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Id of the Group of this membership.",
				ValidateFunc: validation.IsUUID,
			},
			"membership_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The Id of the Membership.",
				ValidateFunc: validation.IsUUID,
			},
			"user_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Id of the User of this membership.",
				ValidateFunc: validation.IsUUID,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func buildUserGroupMembership(data *schema.ResourceData) fusionauth.MemberRequest {
	mr := fusionauth.MemberRequest{
		Members: map[string][]fusionauth.GroupMember{
			data.Get("group_id").(string): {
				{
					Data:    data.Get("data").(map[string]interface{}),
					GroupId: data.Get("group_id").(string),
					Id:      data.Get("membership_id").(string),
					UserId:  data.Get("user_id").(string),
				},
			},
		},
	}

	return mr
}

func createUserGroupMembership(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	mr := buildUserGroupMembership(data)

	resp, faErrs, err := client.FAClient.CreateGroupMembers(mr)
	if err != nil {
		return diag.Errorf("CreateUserGroupMembership err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.Members[data.Get("group_id").(string)][0].Id)
	return nil
}

func readUserGroupMembership(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
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
	switch resp.Total {
	case 0:
		data.SetId("")
		return nil
	case 1:
		data.SetId(gmsresp[0].Id)
	default:
		return diag.Errorf("Found %d memberships for user %s in group %s", resp.Total, data.Get("user_id").(string), data.Get("group_id").(string))
	}

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

func updateUserGroupMembership(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	mr := buildUserGroupMembership(data)

	resp, faErrs, err := client.FAClient.UpdateGroupMembers(mr)

	if err != nil {
		return diag.Errorf("UpdateUserGroupMembership err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteUserGroupMembership(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	mdr := fusionauth.MemberDeleteRequest{
		MemberIds: []string{data.Id()},
	}

	resp, faErrs, err := client.FAClient.DeleteGroupMembers(mdr)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
