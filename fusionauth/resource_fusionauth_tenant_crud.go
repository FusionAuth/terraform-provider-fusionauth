package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func createTenant(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	t := fusionauth.TenantRequest{
		Tenant:         buildTentant(data),
		SourceTenantId: data.Get("source_tentant_id").(string),
	}

	var tid string
	if t, ok := data.GetOk("tenant_id"); ok {
		tid = t.(string)
	}
	resp, faErrs, err := client.FAClient.CreateTenant(tid, t)

	if err != nil {
		return fmt.Errorf("CreateTenant err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	data.SetId(resp.Tenant.Id)
	return nil
}

func readTenant(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveTenant(id)
	if err != nil {
		return err
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	return buildResourceDataFromTenant(resp.Tenant, data)
}

func updateTenant(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)

	t := fusionauth.TenantRequest{
		Tenant:         buildTentant(data),
		SourceTenantId: data.Get("source_tentant_id").(string),
	}

	resp, faErrs, err := client.FAClient.UpdateTenant(data.Id(), t)

	if err != nil {
		return fmt.Errorf("UpdateTenant err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	return nil
}

func deleteTenant(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	resp, faErrs, err := client.FAClient.DeleteTenant(data.Id())
	if err != nil {
		return err
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	return nil
}
