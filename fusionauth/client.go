package fusionauth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Client struct {
	FAClient fusionauth.FusionAuthClient
	Host     string
	APIKey   string
}

func configureClient(_ context.Context, data *schema.ResourceData) (client interface{}, diags diag.Diagnostics) {
	host := data.Get("host").(string)
	apiKey := data.Get("api_key").(string)

	hostURL, err := url.Parse(host)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Fusionauth client",
			Detail:   fmt.Sprintf("Unable to parse the provided Fusionauth hostname to a URL: %s", err),
		})
		return nil, diags
	}

	client = Client{
		Host:   host,
		APIKey: apiKey,
		FAClient: *fusionauth.NewClient(
			&http.Client{
				Timeout: time.Second * 30,
			},
			hostURL,
			apiKey,
		),
	}

	return
}
