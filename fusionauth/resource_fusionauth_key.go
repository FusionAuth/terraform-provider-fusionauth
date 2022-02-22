package fusionauth

import (
	"context"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: createKey,
		ReadContext: func(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			return keyRead(data, buildResourceDataFromKey, i)
		},
		UpdateContext: func(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			return keyUpdate(data, buildKey, i)
		},
		DeleteContext: keyDelete,
		Schema: map[string]*schema.Schema{
			"key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The Id to use for the new Key. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
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
			"kid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id used in the JWT header to identify the key used to generate the signature",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func createKey(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	l := buildKey(data)

	var keyID string
	if a, ok := data.GetOk("key_id"); ok {
		keyID = a.(string)
	}

	resp, faErrs, err := client.FAClient.GenerateKey(keyID, fusionauth.KeyRequest{
		Key: l,
	})
	if err != nil {
		return diag.Errorf("CreateKey err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.Key.Id)
	return buildResourceDataFromKey(data, resp.Key)
}

func buildResourceDataFromKey(data *schema.ResourceData, res fusionauth.Key) diag.Diagnostics {
	if err := data.Set("key_id", res.Id); err != nil {
		return diag.Errorf("key.key_id: %s", err.Error())
	}
	if err := data.Set("algorithm", res.Algorithm); err != nil {
		return diag.Errorf("key.algorithm: %s", err.Error())
	}
	if err := data.Set("name", res.Name); err != nil {
		return diag.Errorf("key.name: %s", err.Error())
	}
	if err := data.Set("length", res.Length); err != nil {
		return diag.Errorf("key.length: %s", err.Error())
	}
	if err := data.Set("kid", res.Kid); err != nil {
		return diag.Errorf("key.kid: %s", err.Error())
	}

	return nil
}
