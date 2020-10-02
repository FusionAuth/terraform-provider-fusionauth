package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func createApplication(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	ar := fusionauth.ApplicationRequest{
		Application: buildApplication(data),
	}

	oldTenantID := client.FAClient.TenantId
	client.FAClient.TenantId = ar.Application.TenantId
	defer func() {
		client.FAClient.TenantId = oldTenantID
	}()

	if hooks, ok := data.GetOk("webhooks"); ok {
		ar.WebhookIds = hooks.([]string)
	}

	var aid string
	if a, ok := data.GetOk("application_id"); ok {
		aid = a.(string)
	}

	resp, faErrs, err := client.FAClient.CreateApplication(aid, ar)
	if err != nil {
		return fmt.Errorf("CreateApplication errors: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	data.SetId(resp.Application.Id)
	if err := buildResourceDataFromApplication(resp.Application, data); err != nil {
		return err
	}

	return nil
}

func readApplication(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, err := client.FAClient.RetrieveApplication(id)
	if err != nil {
		return err
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return err
	}

	return buildResourceDataFromApplication(resp.Application, data)
}

func updateApplication(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	ar := fusionauth.ApplicationRequest{
		Application: buildApplication(data),
	}

	if hooks, ok := data.GetOk("webhooks"); ok {
		ar.WebhookIds = hooks.([]string)
	}

	resp, faErrs, err := client.FAClient.UpdateApplication(data.Id(), ar)
	if err != nil {
		return fmt.Errorf("UpdateApplication err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	return nil
}

func deleteApplication(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	resp, faErrs, err := client.FAClient.DeleteApplication(data.Id())
	if err != nil {
		return err
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	return nil
}
