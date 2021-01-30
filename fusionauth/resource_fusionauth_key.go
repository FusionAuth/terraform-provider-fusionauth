package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func newKey() *schema.Resource {
	return &schema.Resource{
		Create: createKey,
		Read: func(data *schema.ResourceData, i interface{}) error {
			return keyRead(data, buildResourceDataFromKey, i)
		},
		Update: func(data *schema.ResourceData, i interface{}) error {
			return keyUpdate(data, buildKey, i)
		},
		Delete: keyDelete,
		Schema: map[string]*schema.Schema{
			"algorithm": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Description: "The algorithm used to encrypt the Key.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Key.",
			},
			"length": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
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

func buildResourceDataFromKey(data *schema.ResourceData, res fusionauth.Key) error {
	if err := data.Set("algorithm", res.Algorithm); err != nil {
		return fmt.Errorf("key.algorithm: %s", err.Error())
	}
	if err := data.Set("name", res.Name); err != nil {
		return fmt.Errorf("key.name: %s", err.Error())
	}
	if err := data.Set("length", res.Length); err != nil {
		return fmt.Errorf("key.length: %s", err.Error())
	}

	return nil
}
