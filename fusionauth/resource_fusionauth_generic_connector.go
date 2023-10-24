package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newGenericConnector() *schema.Resource {
	return &schema.Resource{
		CreateContext: createGenericConnector,
		ReadContext:   readGenericConnector,
		UpdateContext: updateGenericConnector,
		DeleteContext: deleteGenericConnector,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The Id to use for the new Connector. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"authentication_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The fully qualified URL used to send an HTTP request to authenticate the user.",
			},
			"connect_timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The connect timeout for the HTTP connection, in milliseconds. Value must be greater than 0.",
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that can hold any information about the Connector that should be persisted.",
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
				Description: "The basic authentication password to use for requests to the Connector.",
			},
			"http_authentication_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The basic authentication username to use for requests to the Connector.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique Connector name.",
			},
			"read_timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The read timeout for the HTTP connection, in milliseconds. Value must be greater than 0.",
			},
			"ssl_certificate_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Id of an existing Key. The X509 certificate is used for client certificate authentication in requests to the Connector.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func buildGenericConnector(data *schema.ResourceData) fusionauth.GenericConnectorConfiguration {
	connector := fusionauth.GenericConnectorConfiguration{
		AuthenticationURL: data.Get("authentication_url").(string),
		BaseConnectorConfiguration: fusionauth.BaseConnectorConfiguration{
			Id:    data.Get("id").(string),
			Data:  data.Get("data").(map[string]interface{}),
			Debug: data.Get("debug").(bool),
			Name:  data.Get("name").(string),
			Type:  fusionauth.ConnectorType_Generic,
		},
		ConnectTimeout:             data.Get("connect_timeout").(int),
		HttpAuthenticationPassword: data.Get("http_authentication_password").(string),
		HttpAuthenticationUsername: data.Get("http_authentication_username").(string),
		ReadTimeout:                data.Get("read_timeout").(int),
		SslCertificateKeyId:        data.Get("ssl_certificate_key_id").(string),
	}

	if i, ok := data.GetOk("headers"); ok {
		connector.Headers = intMapToStringMap(i.(map[string]interface{}))
	}

	return connector
}

func createGenericConnector(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	connector := buildGenericConnector(data)
	resp, faErrs, err := CreateConnector(ctx, client.FAClient, connector.Id, GenericConnectorRequest{Connector: connector})
	if err != nil {
		return diag.Errorf("CreateGenericConnector err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(resp.Connector.Id)
	return nil
}

func readGenericConnector(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := RetrieveConnector(ctx, client.FAClient, id)
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

	connector := resp.Connector
	if err := data.Set("authentication_url", connector.AuthenticationURL); err != nil {
		return diag.Errorf("connector.authentication_url: %s", err.Error())
	}
	if err := data.Set("connect_timeout", connector.ConnectTimeout); err != nil {
		return diag.Errorf("connector.connect_timeout: %s", err.Error())
	}
	if err := data.Set("data", connector.Data); err != nil {
		return diag.Errorf("connector.data: %s", err.Error())
	}
	if err := data.Set("debug", connector.Debug); err != nil {
		return diag.Errorf("connector.debug: %s", err.Error())
	}
	if err := data.Set("headers", connector.Headers); err != nil {
		return diag.Errorf("connector.headers: %s", err.Error())
	}
	if err := data.Set("http_authentication_password", connector.HttpAuthenticationPassword); err != nil {
		return diag.Errorf("connector.http_authentication_password: %s", err.Error())
	}
	if err := data.Set("http_authentication_username", connector.HttpAuthenticationUsername); err != nil {
		return diag.Errorf("connector.http_authentication_username: %s", err.Error())
	}
	if err := data.Set("name", connector.Name); err != nil {
		return diag.Errorf("connector.name: %s", err.Error())
	}
	if err := data.Set("read_timeout", connector.ReadTimeout); err != nil {
		return diag.Errorf("connector.read_timeout: %s", err.Error())
	}
	if err := data.Set("ssl_certificate_key_id", connector.SslCertificateKeyId); err != nil {
		return diag.Errorf("connector.ssl_certificate_key_id: %s", err.Error())
	}

	return nil
}

func updateGenericConnector(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	connector := buildGenericConnector(data)

	resp, faErrs, err := UpdateConnector(ctx, client.FAClient, data.Id(), GenericConnectorRequest{Connector: connector})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteGenericConnector(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteConnector(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

type GenericConnectorRequest struct {
	Connector fusionauth.GenericConnectorConfiguration `json:"connector,omitempty"`
}

type GenericConnectorResponse struct {
	fusionauth.BaseHTTPResponse
	Connector  fusionauth.GenericConnectorConfiguration   `json:"connector,omitempty"`
	Connectors []fusionauth.GenericConnectorConfiguration `json:"connectors,omitempty"`
}

func (b *GenericConnectorResponse) SetStatus(status int) {
	b.StatusCode = status
}

func CreateConnector(ctx context.Context, client fusionauth.FusionAuthClient, connectorID string, request GenericConnectorRequest) (*GenericConnectorResponse, *fusionauth.Errors, error) {
	return makeConnectorRequest(ctx, client, connectorID, request, http.MethodPost)
}

func RetrieveConnector(ctx context.Context, client fusionauth.FusionAuthClient, connectorID string) (*GenericConnectorResponse, *fusionauth.Errors, error) {
	return makeConnectorRequest(ctx, client, connectorID, GenericConnectorRequest{}, http.MethodGet)
}

func UpdateConnector(ctx context.Context, client fusionauth.FusionAuthClient, connectorID string, request GenericConnectorRequest) (*GenericConnectorResponse, *fusionauth.Errors, error) {
	return makeConnectorRequest(ctx, client, connectorID, request, http.MethodPut)
}

func makeConnectorRequest(ctx context.Context, client fusionauth.FusionAuthClient, connectorID string, request GenericConnectorRequest, method string) (*GenericConnectorResponse, *fusionauth.Errors, error) {
	var resp GenericConnectorResponse
	var errors fusionauth.Errors

	restClient := client.Start(&resp, &errors)
	err := restClient.WithUri("/api/connector").
		WithUriSegment(connectorID).
		WithJSONBody(request).
		WithMethod(method).
		Do(ctx)
	if restClient.ErrorRef == nil {
		return &resp, nil, err
	}
	return &resp, &errors, err
}
