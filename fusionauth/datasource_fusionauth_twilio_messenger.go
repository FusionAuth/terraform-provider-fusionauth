package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceTwilioMessenger() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTwilioMessengerRead,
		Schema: map[string]*schema.Schema{
			// Data Source Parameters
			"messenger_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"messenger_id", "name"},
				Description:  "The unique Id of the Twilio Messenger to retrieve.",
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"messenger_id", "name"},
				Description:  "The case-insensitive string to search for in the Twilio Messenger name.",
			},
			// Data Source Attributes
			"account_sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Twilio Account ID used when connecting to the Twilio API.",
			},
			"auth_token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The Twilio Auth Token used when connecting to the Twilio API.",
			},
			"data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An object that can hold any information about the Twilio Messenger that should be persisted. Must be a JSON string.",
			},
			"debug": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines if debug should be enabled to create an event log to assist in debugging integration errors.",
			},
			"from_phone_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configured Twilio phone number used to send messages.",
			},
			"messaging_service_sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Twilio message service Id used for Twilio Copilot to load balance between numbers.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Twilio URL that FusionAuth uses to communicate with the Twilio API.",
			},
		},
	}
}

func dataSourceTwilioMessengerRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	var searchTerm string
	var resp *TwilioMessengerResponse
	var err error

	if entityID, ok := data.GetOk("messenger_id"); ok {
		searchTerm = entityID.(string)
		resp, _, err = RetrieveTwilioMessenger(ctx, client.FAClient, searchTerm)
	} else if name, ok := data.GetOk("name"); ok {
		searchTerm = name.(string)
		resp, _, err = retrieveTwilioMessengers(ctx, client.FAClient)
	} else {
		return diag.Errorf("Either 'messenger_id' or 'name' must be specified")
	}

	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	var messenger fusionauth.TwilioMessengerConfiguration
	switch {
	case resp.Messenger.Id != "":
		messenger = resp.Messenger
	case len(resp.Messengers) > 0:
		found := false
		for _, m := range resp.Messengers {
			if m.Name == searchTerm && m.Type == fusionauth.MessengerType_Twilio {
				messenger = m
				found = true
				break
			}
		}
		if !found {
			return diag.Errorf("Couldn't find Twilio Messenger with name '%s'", searchTerm)
		}
	default:
		return diag.Errorf("No Twilio Messengers found matching '%s'", searchTerm)
	}

	data.SetId(messenger.Id)

	if diags := setTwilioMessengerFields(data, messenger); diags != nil {
		return diags
	}

	return nil
}

// Helper function to set all Twilio messenger fields in the schema data
func setTwilioMessengerFields(data *schema.ResourceData, messenger fusionauth.TwilioMessengerConfiguration) diag.Diagnostics {
	fields := map[string]interface{}{
		"messenger_id":          messenger.Id,
		"name":                  messenger.Name,
		"account_sid":           messenger.AccountSID,
		"auth_token":            messenger.AuthToken,
		"debug":                 messenger.Debug,
		"from_phone_number":     messenger.FromPhoneNumber,
		"messaging_service_sid": messenger.MessagingServiceSid,
		"url":                   messenger.Url,
	}

	dataJSON, diags := mapStringInterfaceToJSONString(messenger.Data)
	if diags != nil {
		return diags
	}
	fields["data"] = dataJSON

	for key, value := range fields {
		if err := data.Set(key, value); err != nil {
			return diag.Errorf("error setting twilio_messenger.%s: %s", key, err.Error())
		}
	}

	return nil
}

func retrieveTwilioMessengers(ctx context.Context, client fusionauth.FusionAuthClient) (*TwilioMessengerResponse, *fusionauth.Errors, error) {
	return makeTwilioMessengerRequest(ctx, client, "", TwilioMessengerRequest{}, http.MethodGet)
}
