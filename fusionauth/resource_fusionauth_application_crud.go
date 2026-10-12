package fusionauth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func createApplication(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	ar := fusionauth.ApplicationRequest{
		Application: buildApplication(data),
	}

	clientSecret, diags := getWriteOnlyString(data, "oauth_client_secret_wo")
	if diags.HasError() {
		return diags
	}
	if clientSecret != "" {
		ar.Application.OauthConfiguration.ClientSecret = clientSecret
	}

	oldTenantID := client.FAClient.TenantId
	client.FAClient.TenantId = ar.Application.TenantId
	defer func() {
		client.FAClient.TenantId = oldTenantID
	}()

	var aid string
	if a, ok := data.GetOk("application_id"); ok {
		aid = a.(string)
	}

	resp, faErrs, err := client.FAClient.CreateApplication(aid, ar)
	if err != nil {
		return diag.Errorf("CreateApplication errors: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.Application.Id)
	return buildResourceDataFromApplication(resp.Application, data)
}

func readApplication(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, err := client.FAClient.RetrieveApplication(id)
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

	return buildResourceDataFromApplication(resp.Application, data)
}

func updateApplication(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	ar := fusionauth.ApplicationRequest{
		Application: buildApplication(data),
	}

	if data.HasChange("oauth_client_secret_wo_version") {
		clientSecret, diags := getWriteOnlyString(data, "oauth_client_secret_wo")
		if diags.HasError() {
			return diags
		}
		if clientSecret != "" {
			ar.Application.OauthConfiguration.ClientSecret = clientSecret
		}
	}

	resp, faErrs, err := client.FAClient.UpdateApplication(data.Id(), ar)
	if err != nil {
		return diag.Errorf("UpdateApplication err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	// Moving from oauth_configuration.client_secret to oauth_client_secret_wo would otherwise
	// leave the old secret in state until the next refresh.
	if data.Get("oauth_client_secret_wo_version").(int) != 0 {
		if oc, ok := data.Get("oauth_configuration").([]interface{}); ok && len(oc) == 1 && oc[0] != nil {
			oc[0].(map[string]interface{})["client_secret"] = ""
			if err := data.Set("oauth_configuration", oc); err != nil {
				return diag.Errorf("application.oauth_configuration: %s", err.Error())
			}
		}
	}

	return nil
}

func deleteApplication(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	resp, faErrs, err := client.FAClient.DeleteApplication(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceApplicationV0() *schema.Resource {
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

func resourceApplicationUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
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
