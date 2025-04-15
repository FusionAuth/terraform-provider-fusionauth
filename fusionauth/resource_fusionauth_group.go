package fusionauth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: createGroup,
		ReadContext:   readGroup,
		UpdateContext: updateGroup,
		DeleteContext: deleteGroup,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The Id to use for the new Group. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"data": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "An object that can hold any information about the Group that should be persisted. Please review the limits on data field types as you plan for and build your custom data schema. Must be a JSON string.",
				DiffSuppressFunc: diffSuppressJSON,
				ValidateFunc:     validation.StringIsJSON,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Group.",
			},
			"role_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The Application Roles to assign to this group.",
			},
			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The unique Id of the tenant used to scope this API request.",
				ValidateFunc: validation.IsUUID,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceGroupV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceGroupUpgradeV0,
				Version: 0,
			},
		},
	}
}

func buildGroup(data *schema.ResourceData) fusionauth.GroupRequest {
	var gid string
	if gi, ok := data.Get("group_id").(string); ok {
		gid = gi
	}

	resourceData, _ := jsonStringToMapStringInterface(data.Get("data").(string))

	g := fusionauth.GroupRequest{
		Group: fusionauth.Group{
			Id:       gid,
			Data:     resourceData,
			Name:     data.Get("name").(string),
			TenantId: data.Get("tenant_id").(string),
		},
		RoleIds: handleStringSlice("role_ids", data),
	}

	return g
}

func createGroup(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	g := buildGroup(data)
	oldTenantID := client.FAClient.TenantId
	client.FAClient.TenantId = g.Group.TenantId
	defer func() {
		client.FAClient.TenantId = oldTenantID
	}()
	resp, faErrs, err := client.FAClient.CreateGroup(g.Group.Id, g)
	if err != nil {
		return diag.Errorf("CreateGroup err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.Group.Id)
	return nil
}

func readGroup(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveGroup(id)
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

	t := resp.Group
	if err := data.Set("name", t.Name); err != nil {
		return diag.Errorf("group.name: %s", err.Error())
	}
	if err := data.Set("tenant_id", t.TenantId); err != nil {
		return diag.Errorf("group.tenant_id: %s", err.Error())
	}
	dataJSON, diags := mapStringInterfaceToJSONString(t.Data)
	if diags != nil {
		return diags
	}
	err = data.Set("data", dataJSON)
	if err != nil {
		return diag.Errorf("group.data: %s", err.Error())
	}

	var s []string

	for i := range t.Roles {
		for j := range t.Roles[i] {
			s = append(s, t.Roles[i][j].Id)
		}
	}
	if err := data.Set("role_ids", s); err != nil {
		return diag.Errorf("group.role_ids: %s", err.Error())
	}
	return nil
}

func updateGroup(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	g := buildGroup(data)
	id := data.Id()
	oldTenantID := client.FAClient.TenantId
	client.FAClient.TenantId = g.Group.TenantId
	defer func() {
		client.FAClient.TenantId = oldTenantID
	}()
	resp, faErrs, err := client.FAClient.UpdateGroup(id, g)

	if err != nil {
		return diag.Errorf("UpdateGroup err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteGroup(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteGroup(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGroupV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGroupUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if v, ok := rawState["data"]; ok {
		if dataMap, ok := v.(map[string]interface{}); ok {
			jsonBytes, err := json.Marshal(dataMap)
			if err != nil {
				return nil, err
			}

			rawState["data"] = string(jsonBytes)
		}
	}

	return rawState, nil
}
