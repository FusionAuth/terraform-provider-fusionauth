package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newGenericMessenger() *schema.Resource {
	return &schema.Resource{
		CreateContext: createGenericMessenger,
		ReadContext:   readGenericMessenger,
		UpdateContext: updateGenericMessenger,
		DeleteContext: deleteGenericMessenger,
		Schema: map[string]*schema.Schema{
			"connect_timeout": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The connect timeout for the HTTP connection, in milliseconds. Value must be greater than 0.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"data": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "An object that can hold any information about the Generic Messenger that should be persisted. Please review the limits on data field types as you plan for and build your custom data schema. Must be a JSON string.",
				DiffSuppressFunc: diffSuppressJSON,
				ValidateFunc:     validation.StringIsJSON,
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if debug should be enabled to create an event log to assist in debugging integration errors.",
			},
			"headers": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that can hold HTTPHeader key and value pairs.",
			},
			"http_authentication_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The basic authentication password to use for requests to the Messenger.",
			},
			"http_authentication_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The basic authentication username to use for requests to the Messenger.",
			},
			"messenger_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The Id to use for the new Messenger. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique Messenger name.",
			},
			"read_timeout": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The read timeout for the HTTP connection, in milliseconds. Value must be greater than 0.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"ssl_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "An SSL certificate. The certificate is used for client certificate authentication in requests to the Messenger.",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The fully qualified URL used to send an HTTP request.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func buildGenericMessenger(data *schema.ResourceData) fusionauth.GenericMessengerConfiguration {
	messenger := fusionauth.GenericMessengerConfiguration{
		BaseMessengerConfiguration: fusionauth.BaseMessengerConfiguration{
			Debug: data.Get("debug").(bool),
			Id:    data.Get("messenger_id").(string),
			Name:  data.Get("name").(string),
			Type:  fusionauth.MessengerType_Generic,
		},
		ConnectTimeout:             data.Get("connect_timeout").(int),
		HttpAuthenticationPassword: data.Get("http_authentication_password").(string),
		HttpAuthenticationUsername: data.Get("http_authentication_username").(string),
		ReadTimeout:                data.Get("read_timeout").(int),
		SslCertificate:             data.Get("ssl_certificate").(string),
		Url:                        data.Get("url").(string),
	}

	if i, ok := data.GetOk("data"); ok {
		resourceData, _ := jsonStringToMapStringInterface(i.(string))
		messenger.Data = resourceData
	}

	if i, ok := data.GetOk("headers"); ok {
		messenger.Headers = intMapToStringMap(i.(map[string]interface{}))
	}

	return messenger
}

func createGenericMessenger(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	messenger := buildGenericMessenger(data)
	resp, faErrs, err := CreateMessenger(ctx, client.FAClient, messenger.Id, GenericMessengerRequest{Messenger: messenger})
	if err != nil {
		return diag.Errorf("CreateGenericMessenger err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(resp.Messenger.Id)
	return nil
}

func readGenericMessenger(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := RetrieveMessenger(ctx, client.FAClient, id)
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
	if err := data.Set("connect_timeout", messenger.ConnectTimeout); err != nil {
		return diag.Errorf("messenger.connect_timeout: %s", err.Error())
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
	if err := data.Set("headers", messenger.Headers); err != nil {
		return diag.Errorf("messenger.headers: %s", err.Error())
	}
	if err := data.Set("http_authentication_password", messenger.HttpAuthenticationPassword); err != nil {
		return diag.Errorf("messenger.http_authentication_password: %s", err.Error())
	}
	if err := data.Set("http_authentication_username", messenger.HttpAuthenticationUsername); err != nil {
		return diag.Errorf("messenger.http_authentication_username: %s", err.Error())
	}
	if err := data.Set("name", messenger.Name); err != nil {
		return diag.Errorf("messenger.name: %s", err.Error())
	}
	if err := data.Set("read_timeout", messenger.ReadTimeout); err != nil {
		return diag.Errorf("messenger.read_timeout: %s", err.Error())
	}
	if err := data.Set("ssl_certificate", messenger.SslCertificate); err != nil {
		return diag.Errorf("messenger.ssl_certificate: %s", err.Error())
	}
	if err := data.Set("url", messenger.Url); err != nil {
		return diag.Errorf("messenger.url: %s", err.Error())
	}

	return nil
}

func updateGenericMessenger(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	messenger := buildGenericMessenger(data)

	resp, faErrs, err := UpdateMessenger(ctx, client.FAClient, data.Id(), GenericMessengerRequest{Messenger: messenger})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteGenericMessenger(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
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

type GenericMessengerRequest struct {
	Messenger fusionauth.GenericMessengerConfiguration `json:"messenger,omitempty"`
}

type GenericMessengerResponse struct {
	fusionauth.BaseHTTPResponse
	Messenger  fusionauth.GenericMessengerConfiguration   `json:"messenger,omitempty"`
	Messengers []fusionauth.GenericMessengerConfiguration `json:"messengers,omitempty"`
}

func (b *GenericMessengerResponse) SetStatus(status int) {
	b.StatusCode = status
}

func CreateMessenger(ctx context.Context, client fusionauth.FusionAuthClient, messengerID string, request GenericMessengerRequest) (*GenericMessengerResponse, *fusionauth.Errors, error) {
	return makeMessengerRequest(ctx, client, messengerID, request, http.MethodPost)
}

func RetrieveMessenger(ctx context.Context, client fusionauth.FusionAuthClient, messengerID string) (*GenericMessengerResponse, *fusionauth.Errors, error) {
	return makeMessengerRequest(ctx, client, messengerID, GenericMessengerRequest{}, http.MethodGet)
}

func UpdateMessenger(ctx context.Context, client fusionauth.FusionAuthClient, messengerID string, request GenericMessengerRequest) (*GenericMessengerResponse, *fusionauth.Errors, error) {
	return makeMessengerRequest(ctx, client, messengerID, request, http.MethodPut)
}

func makeMessengerRequest(ctx context.Context, client fusionauth.FusionAuthClient, messengerID string, request GenericMessengerRequest, method string) (*GenericMessengerResponse, *fusionauth.Errors, error) {
	var resp GenericMessengerResponse
	var errors fusionauth.Errors

	restClient := client.Start(&resp, &errors)

	// Set the request body only if the method is not GET
	if method != http.MethodGet {
		restClient.WithJSONBody(request)
	}

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
