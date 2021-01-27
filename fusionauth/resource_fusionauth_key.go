package fusionauth

import (
	"fmt"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func newKey() *schema.Resource {
	return &schema.Resource{
		Create: createKey,
		Read:   readKey,
		Update: updateKey,
		Delete: deleteKey,
		Schema: map[string]*schema.Schema{
			"algorithm": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ES256",
					"ES384",
					"ES512",
					"RS256",
					"RS384",
					"RS512",
					"HS256",
					"HS384",
					"HS512",
				}, false),
				Description: "The algorithm used to encrypt the Key. ",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Key.",
			},
			"length": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The length of the RSA or EC certificate. This field is required when generating RSA key types.",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func buildKey(data *schema.ResourceData) fusionauth.Key {
	l := fusionauth.Key{
		Algorithm: fusionauth.KeyAlgorithm(data.Get("algorithm").(string)),
		Name:      data.Get("name").(string),
		Length:    data.Get("length").(int),
	}
	return l
}

func createKey(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	l := buildKey(data)
	resp, faErrs, err := client.FAClient.GenerateKey("", fusionauth.KeyRequest{
		Key: l,
	})
	if err != nil {
		return fmt.Errorf("CreateKey err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	data.SetId(resp.Key.Id)
	return nil
}

func readKey(data *schema.ResourceData, i interface{}) error {
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

	l := resp.Key
	if err := data.Set("algorithm", l.Algorithm); err != nil {
		return fmt.Errorf("key.algorithm: %s", err.Error())
	}
	if err := data.Set("name", l.Name); err != nil {
		return fmt.Errorf("key.name: %s", err.Error())
	}

	return nil
}

func updateKey(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	l := buildKey(data)

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

func deleteKey(data *schema.ResourceData, i interface{}) error {
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
