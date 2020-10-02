package fusionauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

type IdenityProvidersResponse struct {
	IdentityProviders []IdenityProvider `json:"identityProviders"`
}

type IdenityProvider struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func dataSourceIDP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIDPRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the identity provider. This is only used for display purposes. Will be the type for types: `Apple`, `Facebook`, `Google`, `HYPR`, `Twitter`",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the identity provider.",
				ValidateFunc: validation.StringInSlice([]string{
					"Apple",
					"Facebook",
					"Google",
					"HYPR",
					"Twitter",
					"OpenIDConnect",
					"SAMLv2",
				}, false),
			},
		},
	}
}

func dataSourceIDPRead(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	b, err := readIdenityProviders(client)
	if err != nil {
		return err
	}

	var idps IdenityProvidersResponse
	_ = json.Unmarshal(b, &idps)

	n := data.Get("name").(string)
	t := data.Get("type").(string)
	switch t {
	case "Apple", "Facebook", "Google", "HYPR", "Twitter":
		n = t
	}
	var idp *IdenityProvider
	for i := range idps.IdentityProviders {
		if idps.IdentityProviders[i].Name == n && idps.IdentityProviders[i].Type == t {
			idp = &idps.IdentityProviders[i]
		}
	}

	if idp == nil {
		return fmt.Errorf("couldn't find identity provider name %s, type %s", n, t)
	}
	data.SetId(idp.ID)
	return nil
}

func readIdenityProviders(client Client) ([]byte, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/%s", strings.TrimRight(client.Host, "/"), "api/identity-provider"),
		nil,
	)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", client.APIKey)

	hc := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return b, err
	}
	return b, nil
}
