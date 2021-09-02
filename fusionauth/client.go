package fusionauth

import (
	"net/http"
	"net/url"
	"time"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Client struct {
	FAClient *fusionauth.FusionAuthClient
	Host     string
	APIKey   string
}

func configureClient(data *schema.ResourceData) (interface{}, error) {
	key := data.Get("api_key").(string)

	parsedURL, err := url.Parse(data.Get("host").(string))
	if err != nil {
		return nil, err
	}

	auth := fusionauth.NewClient(
		&http.Client{
			Timeout: time.Second * 30,
		},
		parsedURL,
		key,
	)

	return Client{
		Host:     data.Get("host").(string),
		APIKey:   key,
		FAClient: auth,
	}, nil
}
