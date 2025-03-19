package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGenericConnector() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGenericConnectorRead,
		Schema: map[string]*schema.Schema{
			// Data Source Parameters
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The unique Id of the Generic Connector to retrieve.",
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The case-insensitive string to search for in the Generic Connector name.",
			},
			// Data Source Attributes
			"authentication_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fully qualified URL used to send an HTTP request to authenticate the user.",
			},
			"connect_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The connect timeout for the HTTP connection, in milliseconds. Value must be greater than 0.",
			},
			"data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An object that can hold any information about the Generic Connector that should be persisted. Please review the limits on data field types as you plan for and build your custom data schema. Must be a JSON string.",
			},
			"debug": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines if debug should be enabled to create an event log to assist in debugging integration errors.",
			},
			"headers": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "An object that can hold HTTPHeader key and value pairs.",
			},
			"http_authentication_password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The basic authentication password to use for requests to the Connector.",
			},
			"http_authentication_username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The basic authentication username to use for requests to the Connector.",
			},
			"read_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The read timeout for the HTTP connection, in milliseconds. Value must be greater than 0.",
			},
			"ssl_certificate_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Id of an existing Key. The X509 certificate is used for client certificate authentication in requests to the Connector.",
			},
		},
	}
}

func dataSourceGenericConnectorRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	var searchTerm string
	var resp *GenericConnectorResponse
	var err error

	if entityID, ok := data.GetOk("id"); ok {
		searchTerm = entityID.(string)
		resp, _, err = RetrieveConnector(ctx, client.FAClient, searchTerm)
	} else if name, ok := data.GetOk("name"); ok {
		searchTerm = name.(string)
		resp, _, err = RetrieveConnectors(ctx, client.FAClient)
	} else {
		return diag.Errorf("Either 'id' or 'name' must be specified")
	}

	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	var connector fusionauth.GenericConnectorConfiguration
	if resp.Connector.Id != "" {
		connector = resp.Connector
	} else if len(resp.Connectors) > 0 {
		found := false
		for _, c := range resp.Connectors {
			if c.Name == searchTerm && c.Type == "Generic" {
				connector = c
				found = true
				break
			}
		}
		if !found {
			return diag.Errorf("Couldn't find Generic Connector with name '%s'", searchTerm)
		}
	} else {
		return diag.Errorf("No Generic Connectors found matching '%s'", searchTerm)
	}

	data.SetId(connector.Id)

	if diags := setConnectorFields(data, connector); diags != nil {
		return diags
	}

	return nil
}

// Helper function to set all connector fields in the schema data
func setConnectorFields(data *schema.ResourceData, connector fusionauth.GenericConnectorConfiguration) diag.Diagnostics {
	fields := map[string]interface{}{
		"id":                           connector.Id,
		"name":                         connector.Name,
		"authentication_url":           connector.AuthenticationURL,
		"connect_timeout":              connector.ConnectTimeout,
		"debug":                        connector.Debug,
		"headers":                      connector.Headers,
		"http_authentication_username": connector.HttpAuthenticationUsername,
		"http_authentication_password": connector.HttpAuthenticationPassword,
		"read_timeout":                 connector.ReadTimeout,
		"ssl_certificate_key_id":       connector.SslCertificateKeyId,
	}

	dataJSON, diags := mapStringInterfaceToJSONString(connector.Data)
	if diags != nil {
		return diags
	}
	fields["data"] = dataJSON

	for key, value := range fields {
		if err := data.Set(key, value); err != nil {
			return diag.Errorf("error setting connector.%s: %s", key, err.Error())
		}
	}

	return nil
}

func RetrieveConnectors(ctx context.Context, client fusionauth.FusionAuthClient) (*GenericConnectorResponse, *fusionauth.Errors, error) {
	return makeConnectorRequest(ctx, client, "", GenericConnectorRequest{}, http.MethodGet)
}
