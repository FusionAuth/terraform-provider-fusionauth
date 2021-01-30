package fusionauth

import (
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type keyReadFunc func(*schema.ResourceData, fusionauth.Key) error
type keyBuildFunc func(*schema.ResourceData) fusionauth.Key

func keyUpdate(data *schema.ResourceData, f keyBuildFunc, i interface{}) error {
	client := i.(Client)
	l := f(data)

	resp, faErrs, err := client.FAClient.UpdateKey(data.Id(), fusionauth.KeyRequest{
		Key: l,
	})
	if err != nil {
		return err
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	return nil
}

func keyDelete(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteKey(id)
	if err != nil {
		return err
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	return nil
}

func keyRead(data *schema.ResourceData, f keyReadFunc, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveKey(id)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	return f(data, resp.Key)
}
