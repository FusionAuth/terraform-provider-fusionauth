package fusionauth

import (
	"context"
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
		Tenant:         tenant,
		SourceTenantId: data.Get("source_tenant_id").(string),
		WebhookIds:     handleStringSlice("webhook_ids", data),
	}
	client.FAClient.TenantId = ""

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
		Tenant:         tenant,
		SourceTenantId: data.Get("source_tenant_id").(string),
		WebhookIds:     handleStringSlice("webhook_ids", data),
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
