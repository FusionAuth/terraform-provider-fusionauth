package fusionauth

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func deleteIdentityProvider(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteIdentityProvider(id)
	if err != nil {
		return err
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
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

	b, _ := ioutil.ReadAll(resp.Body)

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
	bb, _ := ioutil.ReadAll(resp.Body)

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
	bb, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("status: %d, response: \n\t%s\nreq body:\n\t%s", resp.StatusCode, string(bb), string(b))
	}
	return bb, nil
}
