package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGenericMessenger() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGenericMessengerRead,
		Schema: map[string]*schema.Schema{
			// Data Source Parameters
			"messenger_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"messenger_id", "name"},
				Description:  "The unique Id of the Generic Messenger to retrieve.",
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"messenger_id", "name"},
				Description:  "The case-insensitive string to search for in the Generic Messenger name.",
			},
			// Data Source Attributes
			"connect_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The connect timeout for the HTTP connection, in milliseconds. Value must be greater than 0.",
			},
			"data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An object that can hold any information about the Generic Messenger that should be persisted. Please review the limits on data field types as you plan for and build your custom data schema. Must be a JSON string.",
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
				Description: "The basic authentication password to use for requests to the Messenger.",
			},
			"http_authentication_username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The basic authentication username to use for requests to the Messenger.",
			},
			"read_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The read timeout for the HTTP connection, in milliseconds. Value must be greater than 0.",
			},
			"ssl_certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "An SSL certificate. The certificate is used for client certificate authentication in requests to the Messenger.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fully qualified URL used to send an HTTP request.",
			},
		},
	}
}

func dataSourceGenericMessengerRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	var searchTerm string
	var resp *GenericMessengerResponse
	var err error

	if entityID, ok := data.GetOk("messenger_id"); ok {
		searchTerm = entityID.(string)
		resp, _, err = RetrieveMessenger(ctx, client.FAClient, searchTerm)
	} else if name, ok := data.GetOk("name"); ok {
		searchTerm = name.(string)
		resp, _, err = retrieveMessengers(ctx, client.FAClient)
	} else {
		return diag.Errorf("Either 'messenger_id' or 'name' must be specified")
	}

	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	var messenger fusionauth.GenericMessengerConfiguration
	switch {
	case resp.Messenger.Id != "":
		messenger = resp.Messenger
	case len(resp.Messengers) > 0:
		found := false
		for _, m := range resp.Messengers {
			if m.Name == searchTerm && m.Type == fusionauth.MessengerType_Generic {
				messenger = m
				found = true
				break
			}
		}
		if !found {
			return diag.Errorf("Couldn't find Generic Messenger with name '%s'", searchTerm)
		}
	default:
		return diag.Errorf("No Generic Messengers found matching '%s'", searchTerm)
	}

	data.SetId(messenger.Id)

	if diags := setMessengerFields(data, messenger); diags != nil {
		return diags
	}

	return nil
}

// Helper function to set all messenger fields in the schema data
func setMessengerFields(data *schema.ResourceData, messenger fusionauth.GenericMessengerConfiguration) diag.Diagnostics {
	fields := map[string]interface{}{
		"messenger_id":                 messenger.Id,
		"name":                         messenger.Name,
		"connect_timeout":              messenger.ConnectTimeout,
		"debug":                        messenger.Debug,
		"headers":                      messenger.Headers,
		"http_authentication_username": messenger.HttpAuthenticationUsername,
		"http_authentication_password": messenger.HttpAuthenticationPassword,
		"read_timeout":                 messenger.ReadTimeout,
		"ssl_certificate":              messenger.SslCertificate,
		"url":                          messenger.Url,
	}

	dataJSON, diags := mapStringInterfaceToJSONString(messenger.Data)
	if diags != nil {
		return diags
	}
	fields["data"] = dataJSON

	for key, value := range fields {
		if err := data.Set(key, value); err != nil {
			return diag.Errorf("error setting messenger.%s: %s", key, err.Error())
		}
	}

	return nil
}

func retrieveMessengers(ctx context.Context, client fusionauth.FusionAuthClient) (*GenericMessengerResponse, *fusionauth.Errors, error) {
	return makeMessengerRequest(ctx, client, "", GenericMessengerRequest{}, http.MethodGet)
}
