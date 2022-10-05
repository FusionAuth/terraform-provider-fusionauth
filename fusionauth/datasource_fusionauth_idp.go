package fusionauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type IdentityProvidersResponse struct {
	IdentityProviders []IdentityProvider `json:"identityProviders"`
}

type IdentityProvider struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func dataSourceIDP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIDPRead,
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

func dataSourceIDPRead(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	b, err := readIdentityProviders(client)
	if err != nil {
		return diag.FromErr(err)
	}

	var idps IdentityProvidersResponse
	_ = json.Unmarshal(b, &idps)

	n := data.Get("name").(string)
	t := data.Get("type").(string)
	switch t {
	case "Apple", "Facebook", "Google", "HYPR", "Twitter":
		n = t
	}
	var idp *IdentityProvider
	for i := range idps.IdentityProviders {
		if idps.IdentityProviders[i].Name == n && idps.IdentityProviders[i].Type == t {
			idp = &idps.IdentityProviders[i]
		}
	}

	if idp == nil {
		return diag.Errorf("couldn't find identity provider name %s, type %s", n, t)
	}
	data.SetId(idp.ID)
	return nil
}

func readIdentityProviders(client Client) ([]byte, error) {
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

	b, _ := io.ReadAll(resp.Body)

	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return b, err
	}
	return b, nil
}
