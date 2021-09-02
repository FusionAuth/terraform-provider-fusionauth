package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type keyReadFunc func(*schema.ResourceData, fusionauth.Key) diag.Diagnostics
type keyBuildFunc func(*schema.ResourceData) fusionauth.Key

func keyUpdate(data *schema.ResourceData, f keyBuildFunc, i interface{}) diag.Diagnostics {
	client := i.(Client)
	l := f(data)

	resp, faErrs, err := client.FAClient.UpdateKey(data.Id(), fusionauth.KeyRequest{Key: l})
	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func keyDelete(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteKey(id)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func keyRead(data *schema.ResourceData, f keyReadFunc, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveKey(id)
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

	return f(data, resp.Key)
}
