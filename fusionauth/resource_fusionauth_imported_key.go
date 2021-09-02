package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceImportedKey() *schema.Resource {
	return &schema.Resource{
		Create: createImportedKey,
		Read: func(data *schema.ResourceData, i interface{}) error {
			return keyRead(data, buildResourceDataFromImportedKey, i)
		},
		Update: func(data *schema.ResourceData, i interface{}) error {
			return keyUpdate(data, buildImportedKey, i)
		},
		Delete: keyDelete,
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
				DiffSuppressFunc: certKeyCompare,
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
				DiffSuppressFunc: certKeyCompare,
			},
			"private_key": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Sensitive:        true,
				Description:      "The Key private key. Optional if importing an RSA or EC key. If the key is only to be used for token validation, only a public key is necessary and this field may be omitted.",
				DiffSuppressFunc: certKeyCompare,
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
			State: schema.ImportStatePassthrough,
		},
	}
}

func createImportedKey(data *schema.ResourceData, i interface{}) error {
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
		return fmt.Errorf("CreateKey err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
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

func buildResourceDataFromImportedKey(data *schema.ResourceData, res fusionauth.Key) error {
	if err := data.Set("algorithm", res.Algorithm); err != nil {
		return fmt.Errorf("key.algorithm: %s", err.Error())
	}
	if err := data.Set("certificate", res.Certificate); err != nil {
		return fmt.Errorf("key.certificate: %s", err.Error())
	}
	if err := data.Set("kid", res.Kid); err != nil {
		return fmt.Errorf("key.kid: %s", err.Error())
	}
	if err := data.Set("name", res.Name); err != nil {
		return fmt.Errorf("key.name: %s", err.Error())
	}
	if err := data.Set("public_key", res.PublicKey); err != nil {
		return fmt.Errorf("key.public_key: %s", err.Error())
	}
	if err := data.Set("type", res.Type); err != nil {
		return fmt.Errorf("key.type: %s", err.Error())
	}

	return nil
}
