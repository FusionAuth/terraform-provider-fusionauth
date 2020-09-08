package fusionauth

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func deleteIdentityProvider(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	_, faErrs, err := client.FAClient.DeleteIdentityProvider(id)
	if err != nil {
		return err
	}

	if faErrs != nil {
		return fmt.Errorf("DeleteIdentityProvider errors: %v", faErrs)
	}

	return nil
}

func readIdenityProvider(id string, client Client) ([]byte, error) {
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

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, response: \n\t%s", resp.StatusCode, string(b))
	}
	return b, nil
}

func createIdenityProvider(b []byte, client Client) ([]byte, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%s", strings.TrimRight(client.Host, "/"), "api/identity-provider"),
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

	bb, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("status: %d, response: \n\t%s\nreq body:\n\t%s", resp.StatusCode, string(bb), string(b))
	}
	return bb, nil
}

func updateIdenityProvider(b []byte, id string, client Client) ([]byte, error) {
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

	bb, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("status: %d, response: \n\t%s\nreq body:\n\t%s", resp.StatusCode, string(bb), string(b))
	}
	return bb, nil
}
