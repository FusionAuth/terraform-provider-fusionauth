package fusionauth

import (
	"context"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceImportedKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: createImportedKey,
		ReadContext: func(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			return keyRead(data, buildResourceDataFromImportedKey, i)
		},
		UpdateContext: func(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			return keyUpdate(data, buildImportedKey, i)
		},
		DeleteContext: keyDelete,
		Schema: map[string]*schema.Schema{
			"key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id to use for the new Key. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"certificate": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				Description:      "The certificate to import. The publicKey will be extracted from the certificate.",
				DiffSuppressFunc: diffSuppressCertKey,
			},
			"kid": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "The Key identifier 'kid'.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Key.",
			},
			"public_key": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				Description:      "The Key public key. Required if importing an RSA or EC key and a certificate is not provided.",
				DiffSuppressFunc: diffSuppressCertKey,
			},
			"private_key": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Sensitive:        true,
				Description:      "The Key private key. Optional if importing an RSA or EC key. If the key is only to be used for token validation, only a public key is necessary and this field may be omitted.",
				DiffSuppressFunc: diffSuppressCertKey,
			},
			"secret": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "The Key secret. This field is required if importing an HMAC key type.",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"EC",
					"RSA",
					"HMAC",
				}, false),
				Description: "The Key type. This field is required if importing an HMAC key type, or if importing a public key / private key pair.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func createImportedKey(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	l := buildImportedKey(data)

	var keyID string
	if a, ok := data.GetOk("key_id"); ok {
		keyID = a.(string)
	}

	resp, faErrs, err := client.FAClient.ImportKey(keyID, fusionauth.KeyRequest{
		Key: l,
	})
	if err != nil {
		return diag.Errorf("CreateKey err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.Key.Id)
	return buildResourceDataFromImportedKey(data, resp.Key)
}

func buildImportedKey(data *schema.ResourceData) fusionauth.Key {
	return fusionauth.Key{
		Algorithm:   fusionauth.KeyAlgorithm(data.Get("algorithm").(string)),
		Certificate: data.Get("certificate").(string),
		Kid:         data.Get("kid").(string),
		Name:        data.Get("name").(string),
		PublicKey:   data.Get("public_key").(string),
		PrivateKey:  data.Get("private_key").(string),
		Secret:      data.Get("secret").(string),
		Type:        fusionauth.KeyType(data.Get("type").(string)),
	}
}

func buildResourceDataFromImportedKey(data *schema.ResourceData, res fusionauth.Key) diag.Diagnostics {
	if err := data.Set("algorithm", res.Algorithm); err != nil {
		return diag.Errorf("key.algorithm: %s", err.Error())
	}
	if err := data.Set("certificate", res.Certificate); err != nil {
		return diag.Errorf("key.certificate: %s", err.Error())
	}
	if err := data.Set("kid", res.Kid); err != nil {
		return diag.Errorf("key.kid: %s", err.Error())
	}
	if err := data.Set("name", res.Name); err != nil {
		return diag.Errorf("key.name: %s", err.Error())
	}
	if err := data.Set("public_key", res.PublicKey); err != nil {
		return diag.Errorf("key.public_key: %s", err.Error())
	}
	if err := data.Set("type", res.Type); err != nil {
		return diag.Errorf("key.type: %s", err.Error())
	}

	return nil
}
