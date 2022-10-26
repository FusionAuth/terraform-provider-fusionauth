package fusionauth

import (
	"context"
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

	resp, faErrs, err := client.FAClient.UpdateApplication(data.Id(), ar)
	if err != nil {
		return diag.Errorf("UpdateApplication err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
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
