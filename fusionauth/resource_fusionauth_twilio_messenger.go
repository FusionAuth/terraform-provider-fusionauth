package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newTwilioMessenger() *schema.Resource {
	return &schema.Resource{
		CreateContext: createTwilioMessenger,
		ReadContext:   readTwilioMessenger,
		UpdateContext: updateTwilioMessenger,
		DeleteContext: deleteTwilioMessenger,
		Schema: map[string]*schema.Schema{
			"account_sid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Twilio Account ID to use when connecting to the Twilio API. This can be found in your Twilio dashboard.",
			},
			"auth_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The Twilio Auth Token to use when connecting to the Twilio API. This can be found in your Twilio dashboard.",
			},
			"data": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "An object that can hold any information about the Twilio Messenger that should be persisted. Please review the limits on data field types as you plan for and build your custom data schema. Must be a JSON string.",
				DiffSuppressFunc: diffSuppressJSON,
				ValidateFunc:     validation.StringIsJSON,
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If debug is enabled, an event log is created to assist in debugging messenger errors.",
			},
			"from_phone_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The configured Twilio phone number that will be used to send messages. This can be found in your Twilio dashboard.",
			},
			"messenger_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The Id to use for the new Messenger. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"messaging_service_sid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Twilio message service Id, this is used when using Twilio Copilot to load balance between numbers. This can be found in your Twilio dashboard. If this is set, the fromPhoneNumber will be ignored.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique Messenger name.",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Twilio URL that FusionAuth will use to communicate with the Twilio API.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func buildTwilioMessenger(data *schema.ResourceData) fusionauth.TwilioMessengerConfiguration {
	messenger := fusionauth.TwilioMessengerConfiguration{
		BaseMessengerConfiguration: fusionauth.BaseMessengerConfiguration{
			Debug: data.Get("debug").(bool),
			Id:    data.Get("messenger_id").(string),
			Name:  data.Get("name").(string),
			Type:  fusionauth.MessengerType_Twilio,
		},
		AccountSID:          data.Get("account_sid").(string),
		AuthToken:           data.Get("auth_token").(string),
		FromPhoneNumber:     data.Get("from_phone_number").(string),
		MessagingServiceSid: data.Get("messaging_service_sid").(string),
		Url:                 data.Get("url").(string),
	}

	if i, ok := data.GetOk("data"); ok {
		resourceData, _ := jsonStringToMapStringInterface(i.(string))
		messenger.BaseMessengerConfiguration.Data = resourceData
	}

	return messenger
}

func createTwilioMessenger(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	messenger := buildTwilioMessenger(data)
	resp, faErrs, err := CreateTwilioMessenger(ctx, client.FAClient, messenger.Id, TwilioMessengerRequest{Messenger: messenger})
	if err != nil {
		return diag.Errorf("CreateTwilioMessenger err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(resp.Messenger.Id)
	return nil
}

func readTwilioMessenger(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := RetrieveTwilioMessenger(ctx, client.FAClient, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	messenger := resp.Messenger
	if err := data.Set("account_sid", messenger.AccountSID); err != nil {
		return diag.Errorf("messenger.account_sid: %s", err.Error())
	}
	if err := data.Set("auth_token", messenger.AuthToken); err != nil {
		return diag.Errorf("messenger.auth_token: %s", err.Error())
	}
	dataJSON, diags := mapStringInterfaceToJSONString(messenger.Data)
	if diags != nil {
		return diags
	}
	err = data.Set("data", dataJSON)
	if err != nil {
		return diag.Errorf("messenger.data: %s", err.Error())
	}
	if err := data.Set("debug", messenger.Debug); err != nil {
		return diag.Errorf("messenger.debug: %s", err.Error())
	}
	if err := data.Set("from_phone_number", messenger.FromPhoneNumber); err != nil {
		return diag.Errorf("messenger.from_phone_number: %s", err.Error())
	}
	if err := data.Set("messaging_service_sid", messenger.MessagingServiceSid); err != nil {
		return diag.Errorf("messenger.messaging_service_sid: %s", err.Error())
	}
	if err := data.Set("name", messenger.Name); err != nil {
		return diag.Errorf("messenger.name: %s", err.Error())
	}
	if err := data.Set("url", messenger.Url); err != nil {
		return diag.Errorf("messenger.url: %s", err.Error())
	}

	return nil
}

func updateTwilioMessenger(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	messenger := buildTwilioMessenger(data)

	resp, faErrs, err := UpdateTwilioMessenger(ctx, client.FAClient, data.Id(), TwilioMessengerRequest{Messenger: messenger})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteTwilioMessenger(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteMessenger(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

type TwilioMessengerRequest struct {
	Messenger fusionauth.TwilioMessengerConfiguration `json:"messenger,omitempty"`
}

type TwilioMessengerResponse struct {
	fusionauth.BaseHTTPResponse
	Messenger  fusionauth.TwilioMessengerConfiguration   `json:"messenger,omitempty"`
	Messengers []fusionauth.TwilioMessengerConfiguration `json:"messengers,omitempty"`
}

func (b *TwilioMessengerResponse) SetStatus(status int) {
	b.StatusCode = status
}

func CreateTwilioMessenger(ctx context.Context, client fusionauth.FusionAuthClient, messengerID string, request TwilioMessengerRequest) (*TwilioMessengerResponse, *fusionauth.Errors, error) {
	return makeTwilioMessengerRequest(ctx, client, messengerID, request, http.MethodPost)
}

func RetrieveTwilioMessenger(ctx context.Context, client fusionauth.FusionAuthClient, messengerID string) (*TwilioMessengerResponse, *fusionauth.Errors, error) {
	return makeTwilioMessengerRequest(ctx, client, messengerID, TwilioMessengerRequest{}, http.MethodGet)
}

func UpdateTwilioMessenger(ctx context.Context, client fusionauth.FusionAuthClient, messengerID string, request TwilioMessengerRequest) (*TwilioMessengerResponse, *fusionauth.Errors, error) {
	return makeTwilioMessengerRequest(ctx, client, messengerID, request, http.MethodPut)
}

func makeTwilioMessengerRequest(ctx context.Context, client fusionauth.FusionAuthClient, messengerID string, request TwilioMessengerRequest, method string) (*TwilioMessengerResponse, *fusionauth.Errors, error) {
	var resp TwilioMessengerResponse
	var errors fusionauth.Errors

	restClient := client.Start(&resp, &errors)
	err := restClient.WithUri("/api/messenger").
		WithUriSegment(messengerID).
		WithJSONBody(request).
		WithMethod(method).
		Do(ctx)
	if restClient.ErrorRef == nil {
		return &resp, nil, err
	}
	return &resp, &errors, err
}
