package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceIPAccessControlList() *schema.Resource {
	return &schema.Resource{
		CreateContext: createIPAccessControlList,
		ReadContext:   readIPAccessControlList,
		UpdateContext: updateIPAccessControlList,
		DeleteContext: deleteIPAccessControlList,
		Schema: map[string]*schema.Schema{
			"ip_access_control_list_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The Id to use for the new IP ACL. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"entries": {
				Required:    true,
				Type:        schema.TypeSet,
				MinItems:    1,
				Description: "A list of IP ranges and the action to apply for each. One and only one entry must have a startIPAddress of * to indicate the default action of the IP ACL.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_ip_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The starting IP (IPv4) for this range.",
						},
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The action to take for this IP Range.",
							ValidateFunc: validation.StringInSlice([]string{
								"Allow",
								"Block",
							}, false),
						},
						"end_ip_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ending IP (IPv4) for this range. The only time this is not required is when start_ip_address is equal to *, in which case this field is ignored. This value must be greater than or equal to the start_ip_address. To define a range of a single IP address, set this field equal to the value for start_ip_address.",
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The unique name of this IP ACL.`,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func createIPAccessControlList(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	acl := buildIPAccessControlList(data)
	var id string
	if a, ok := data.GetOk("ip_access_control_list_id"); ok {
		id = a.(string)
	}

	resp, faErrs, err := client.FAClient.CreateIPAccessControlList(id, fusionauth.IPAccessControlListRequest{
		IpAccessControlList: acl,
	})
	if err != nil {
		return diag.Errorf("createIPAccessControlList err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.IpAccessControlList.Id)
	if err := data.Set("ip_access_control_list_id", resp.IpAccessControlList.Id); err != nil {
		return diag.Errorf("ipAccessControlList.ip_access_control_list_id: %s", err.Error())
	}
	return nil
}

func readIPAccessControlList(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, err := client.FAClient.RetrieveIPAccessControlList(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}
	return buildResourceDataFromIPAccessControlList(data, resp.IpAccessControlList)
}

func updateIPAccessControlList(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	acl := buildIPAccessControlList(data)

	resp, faErrs, err := client.FAClient.UpdateIPAccessControlList(data.Id(), fusionauth.IPAccessControlListRequest{
		IpAccessControlList: acl,
	})
	if err != nil {
		return diag.Errorf("updateIPAccessControlList err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteIPAccessControlList(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	resp, faErrs, err := client.FAClient.DeleteIPAccessControlList(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildIPAccessControlList(data *schema.ResourceData) fusionauth.IPAccessControlList {
	acl := fusionauth.IPAccessControlList{
		Name: data.Get("name").(string),
		Id:   data.Get("ip_access_control_list_id").(string),
	}
	entries := []fusionauth.IPAccessControlEntry{}
	s := data.Get("entries")
	set := s.(*schema.Set)
	l := set.List()
	for _, x := range l {
		e := x.(map[string]interface{})
		entries = append(entries, fusionauth.IPAccessControlEntry{
			StartIPAddress: e["start_ip_address"].(string),
			EndIPAddress:   e["end_ip_address"].(string),
			Action:         fusionauth.IPAccessControlEntryAction(e["action"].(string)),
		})
	}
	acl.Entries = entries
	return acl
}

func buildResourceDataFromIPAccessControlList(data *schema.ResourceData, res fusionauth.IPAccessControlList) diag.Diagnostics {
	if err := data.Set("name", res.Name); err != nil {
		return diag.Errorf("ipAccessControlList.name: %s", err.Error())
	}

	entries := make([]map[string]interface{}, 0, len(res.Entries))
	for _, e := range res.Entries {
		entries = append(entries, map[string]interface{}{
			"start_ip_address": e.StartIPAddress,
			"end_ip_address":   e.EndIPAddress,
			"action":           e.Action,
		})
	}
	if err := data.Set("entries", entries); err != nil {
		return diag.Errorf("ipAccessControlList.entries: %s", err.Error())
	}
	return nil
}
