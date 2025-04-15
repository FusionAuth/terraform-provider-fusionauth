package fusionauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func createTenant(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	tenant, diags := buildTenant(data)
	if diags != nil {
		return diags
	}

	t := fusionauth.TenantRequest{
		Tenant:     tenant,
		WebhookIds: handleStringSlice("webhook_ids", data),
	}
	client.FAClient.TenantId = ""

	if srcTenant, ok := data.GetOk("source_tenant_id"); ok {
		t.SourceTenantId = srcTenant.(string)
	}

	// Add validation for SCIM configuration fields
	if err := validateSCIMConfiguration(data); err != nil {
		return diag.FromErr(err)
	}

	var tid string
	if t, ok := data.GetOk("tenant_id"); ok {
		tid = t.(string)
	}
	resp, faErrs, err := client.FAClient.CreateTenant(tid, t)
	if err != nil {
		return diag.Errorf("CreateTenant err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.Tenant.Id)
	return buildResourceDataFromTenant(resp.Tenant, data)
}

func readTenant(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveTenant(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return buildResourceDataFromTenant(resp.Tenant, data)
}

func updateTenant(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	tenant, diags := buildTenant(data)
	if diags != nil {
		return diags
	}

	t := fusionauth.TenantRequest{
		Tenant:     tenant,
		WebhookIds: handleStringSlice("webhook_ids", data),
	}

	// Add validation for SCIM configuration fields
	if err := validateSCIMConfiguration(data); err != nil {
		return diag.FromErr(err)
	}

	resp, faErrs, err := client.FAClient.UpdateTenant(data.Id(), t)
	if err != nil {
		return diag.Errorf("UpdateTenant err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return buildResourceDataFromTenant(resp.Tenant, data)
}

func deleteTenant(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	resp, faErrs, err := client.FAClient.DeleteTenant(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceTenantV0() *schema.Resource {
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

func resourceTenantUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
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

// validateSCIMConfiguration validates that required SCIM fields are provided when SCIM is enabled
func validateSCIMConfiguration(d *schema.ResourceData) error {
	// Check if SCIM is enabled
	scimEnabled := false
	if scimServerConfig, ok := d.GetOk("scim_server_configuration"); ok && len(scimServerConfig.([]interface{})) > 0 {
		config := scimServerConfig.([]interface{})[0].(map[string]interface{})
		if enabled, ok := config["enabled"].(bool); ok && enabled {
			scimEnabled = true
		}
	}

	// If SCIM is enabled, check for required lambda fields
	if scimEnabled {
		lambdaConfig, hasLambdaConfig := d.GetOk("lambda_configuration")
		if !hasLambdaConfig || len(lambdaConfig.([]interface{})) == 0 {
			return fmt.Errorf("lambda_configuration is required when scim_server_configuration.enabled is true")
		}

		lambdaConfigMap := lambdaConfig.([]interface{})[0].(map[string]interface{})
		requiredFields := []string{
			"scim_enterprise_user_request_converter_id",
			"scim_enterprise_user_response_converter_id",
			"scim_group_request_converter_id",
			"scim_group_response_converter_id",
			"scim_user_request_converter_id",
			"scim_user_response_converter_id",
		}

		for _, field := range requiredFields {
			if _, ok := lambdaConfigMap[field].(string); !ok || lambdaConfigMap[field].(string) == "" {
				return fmt.Errorf("%s is required in lambda_configuration when scim_server_configuration.enabled is true", field)
			}
		}
	}

	return nil
}
