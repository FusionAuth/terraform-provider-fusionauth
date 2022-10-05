package fusionauth

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func deleteIdentityProvider(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteIdentityProvider(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func readIdentityProvider(id string, client Client) ([]byte, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/%s/%s", strings.TrimRight(client.Host, "/"), "api/identity-provider", id),
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
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return nil, err
	}

	b, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, response: \n\t%s", resp.StatusCode, string(b))
	}
	return b, nil
}

func createIdentityProvider(b []byte, client Client, idpID string) ([]byte, error) {
	var u string
	if idpID != "" {
		u = fmt.Sprintf("%s/%s/%s", strings.TrimRight(client.Host, "/"), "api/identity-provider", idpID)
	} else {
		u = fmt.Sprintf("%s/%s", strings.TrimRight(client.Host, "/"), "api/identity-provider")
	}

	req, err := http.NewRequest(
		http.MethodPost,
		u,
		bytes.NewBuffer(b),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", client.APIKey)
	req.Header.Add("Content-Type", "application/json")

	hc := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return nil, err
	}
	bb, _ := io.ReadAll(resp.Body)

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("status: %d, response: \n\t%s\nreq body:\n\t%s", resp.StatusCode, string(bb), string(b))
	}
	return bb, nil
}

func updateIdentityProvider(b []byte, id string, client Client) ([]byte, error) {
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/%s/%s", strings.TrimRight(client.Host, "/"), "api/identity-provider", id),
		bytes.NewBuffer(b),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", client.APIKey)
	req.Header.Add("Content-Type", "application/json")

	hc := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return nil, err
	}
	bb, _ := io.ReadAll(resp.Body)

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("status: %d, response: \n\t%s\nreq body:\n\t%s", resp.StatusCode, string(bb), string(b))
	}
	return bb, nil
}

func buildTenantConfigurationResource(tcm map[string]fusionauth.IdentityProviderTenantConfiguration) []map[string]interface{} {
	t := make([]map[string]interface{}, 0, len(tcm))
	for k, v := range tcm {
		t = append(t, map[string]interface{}{
			"tenant_id":                           k,
			"limit_user_link_count_enabled":       v.LimitUserLinkCount.Enabled,
			"limit_user_link_count_maximum_links": v.LimitUserLinkCount.MaximumLinks,
		})
	}
	return t
}

func buildTenantConfiguration(data *schema.ResourceData) map[string]fusionauth.IdentityProviderTenantConfiguration {
	m := make(map[string]fusionauth.IdentityProviderTenantConfiguration)
	s := data.Get("tenant_configuration")
	set, ok := s.(*schema.Set)
	if !ok {
		return nil
	}

	l := set.List()
	for _, x := range l {
		ac := x.(map[string]interface{})
		aid := ac["tenant_id"].(string)
		oc := fusionauth.IdentityProviderTenantConfiguration{
			LimitUserLinkCount: fusionauth.IdentityProviderLimitUserLinkingPolicy{
				MaximumLinks: ac["limit_user_link_count_maximum_links"].(int),
			},
		}
		oc.LimitUserLinkCount.Enabled = ac["limit_user_link_count_enabled"].(bool)
		m[aid] = oc
	}

	return m
}
