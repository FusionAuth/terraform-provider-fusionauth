package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newLDAPConnector() *schema.Resource {
	return &schema.Resource{
		CreateContext: createLDAPConnector,
		ReadContext:   readLDAPConnector,
		UpdateContext: updateLDAPConnector,
		DeleteContext: deleteLDAPConnector,
		Schema: map[string]*schema.Schema{
			"authentication_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The fully qualified LDAP URL to authenticate.",
			},
			"base_structure": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The top of the LDAP directory hierarchy. Typically this contains the dc (domain component) element.",
			},
			"connect_timeout": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The connect timeout for the HTTP connection, in milliseconds. Value must be greater than 0.",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"data": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "An object that can hold any information about the Connector that should be persisted.. Must be a JSON encoded string.",
				DiffSuppressFunc: diffSuppressJSON,
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable debug logging.",
			},
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The Id to use for the new Connector. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"identifying_attribute": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The entry attribute name which is the first component of the distinguished name of entries in the directory.",
			},
			"lambda_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reconcile_id": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The Id of an existing Lambda. The lambda is executed after the user authenticates with the connector. This lambda can create a user, registrations, and group memberships in FusionAuth based on attributes returned from the connector.",
							ValidateFunc: validation.IsUUID,
						},
					},
				},
				DiffSuppressFunc: suppressBlockDiff,
				Required:         true,
			},
			"login_id_attribute": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The entity attribute name which stores the identifier that is used for logging the user in.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique LDAP Connector name.",
			},
			"read_timeout": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The read timeout for the HTTP connection, in milliseconds. Value must be greater than 0.",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"requested_attributes": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The list of attributes to request from the LDAP server. This is a list of strings.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"security_method": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The LDAP security method. Possible values are: None (Requests will be made without encryption), LDAPS (A secure connection will be made to a secure port over using the LDAPS protocol) or StartTLS (An un-secured connection will initially be established, followed by secure connection established using the StartTLS extension).",
				ValidateFunc: validation.StringInSlice([]string{
					"None",
					"LDAPS",
					"StartTLS",
				}, false),
			},
			"system_account_dn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The distinguished name of the system account used to authenticate to the LDAP server. This account must have permission to search the directory.",
			},
			"system_account_password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The password of an entry that has read access to the directory.",
				Sensitive:   true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func buildLDAPConnector(data *schema.ResourceData) fusionauth.LDAPConnectorConfiguration {
	resourceData, _ := jsonStringToMapStringInterface(data.Get("data").(string))
	connector := fusionauth.LDAPConnectorConfiguration{
		AuthenticationURL: data.Get("authentication_url").(string),
		BaseConnectorConfiguration: fusionauth.BaseConnectorConfiguration{
			Data:  resourceData,
			Debug: data.Get("debug").(bool),
			Id:    data.Get("id").(string),
			Name:  data.Get("name").(string),
			Type:  fusionauth.ConnectorType_LDAP,
		},
		BaseStructure:        data.Get("base_structure").(string),
		ConnectTimeout:       data.Get("connect_timeout").(int),
		IdentifyingAttribute: data.Get("identifying_attribute").(string),
		LambdaConfiguration: fusionauth.ConnectorLambdaConfiguration{
			ReconcileId: data.Get("lambda_configuration.0.reconcile_id").(string),
		},
		LoginIdAttribute:      data.Get("login_id_attribute").(string),
		ReadTimeout:           data.Get("read_timeout").(int),
		RequestedAttributes:   handleStringSlice("requested_attributes", data),
		SecurityMethod:        fusionauth.LDAPSecurityMethod(data.Get("security_method").(string)),
		SystemAccountDN:       data.Get("system_account_dn").(string),
		SystemAccountPassword: data.Get("system_account_password").(string),
	}

	return connector
}

func createLDAPConnector(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	connector := buildLDAPConnector(data)
	resp, faErrs, err := CreateLDAPConnector(ctx, client.FAClient, connector.Id, LDAPConnectorRequest{Connector: connector})
	if err != nil {
		return diag.Errorf("CreateLDAPConnector err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(resp.Connector.Id)
	return nil
}

func readLDAPConnector(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := RetrieveLDAPConnector(ctx, client.FAClient, id)
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
	if err := data.Set("base_structure", connector.BaseStructure); err != nil {
		return diag.Errorf("connector.base_structure: %s", err.Error())
	}
	if err := data.Set("connect_timeout", connector.ConnectTimeout); err != nil {
		return diag.Errorf("connector.connect_timeout: %s", err.Error())
	}
	dataJSON, diags := mapStringInterfaceToJSONString(connector.Data)
	if diags != nil {
		return diags
	}
	err = data.Set("data", dataJSON)
	if err != nil {
		return diag.Errorf("connector.data: %s", err.Error())
	}
	if err := data.Set("debug", connector.Debug); err != nil {
		return diag.Errorf("connector.debug: %s", err.Error())
	}
	if err := data.Set("id", connector.Id); err != nil {
		return diag.Errorf("connector.id: %s", err.Error())
	}
	if err := data.Set("identifying_attribute", connector.IdentifyingAttribute); err != nil {
		return diag.Errorf("connector.identifying_attribute: %s", err.Error())
	}
	if err := data.Set("lambda_configuration", []interface{}{
		map[string]interface{}{
			"reconcile_id": connector.LambdaConfiguration.ReconcileId,
		},
	}); err != nil {
		return diag.Errorf("connector.lambda_configuration: %s", err.Error())
	}
	if err := data.Set("login_id_attribute", connector.LoginIdAttribute); err != nil {
		return diag.Errorf("connector.login_id_attribute: %s", err.Error())
	}
	if err := data.Set("name", connector.Name); err != nil {
		return diag.Errorf("connector.name: %s", err.Error())
	}
	if err := data.Set("read_timeout", connector.ReadTimeout); err != nil {
		return diag.Errorf("connector.read_timeout: %s", err.Error())
	}
	if err := data.Set("requested_attributes", connector.RequestedAttributes); err != nil {
		return diag.Errorf("connector.requested_attributes: %s", err.Error())
	}
	if err := data.Set("security_method", connector.SecurityMethod); err != nil {
		return diag.Errorf("connector.security_method: %s", err.Error())
	}
	if err := data.Set("system_account_dn", connector.SystemAccountDN); err != nil {
		return diag.Errorf("connector.system_account_dn: %s", err.Error())
	}
	if err := data.Set("system_account_password", connector.SystemAccountPassword); err != nil {
		return diag.Errorf("connector.system_account_password: %s", err.Error())
	}
	return nil
}

func updateLDAPConnector(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	connector := buildLDAPConnector(data)

	resp, faErrs, err := UpdateLDAPConnector(ctx, client.FAClient, data.Id(), LDAPConnectorRequest{Connector: connector})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteLDAPConnector(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
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

type LDAPConnectorRequest struct {
	Connector fusionauth.LDAPConnectorConfiguration `json:"connector,omitempty"`
}

type LDAPConnectorResponse struct {
	fusionauth.BaseHTTPResponse
	Connector  fusionauth.LDAPConnectorConfiguration   `json:"connector,omitempty"`
	Connectors []fusionauth.LDAPConnectorConfiguration `json:"connectors,omitempty"`
}

func (b *LDAPConnectorResponse) SetStatus(status int) {
	b.StatusCode = status
}

func CreateLDAPConnector(ctx context.Context, client fusionauth.FusionAuthClient, connectorID string, request LDAPConnectorRequest) (*LDAPConnectorResponse, *fusionauth.Errors, error) {
	return makeLDAPConnectorRequest(ctx, client, connectorID, request, http.MethodPost)
}

func RetrieveLDAPConnector(ctx context.Context, client fusionauth.FusionAuthClient, connectorID string) (*LDAPConnectorResponse, *fusionauth.Errors, error) {
	return makeLDAPConnectorRequest(ctx, client, connectorID, LDAPConnectorRequest{}, http.MethodGet)
}

func UpdateLDAPConnector(ctx context.Context, client fusionauth.FusionAuthClient, connectorID string, request LDAPConnectorRequest) (*LDAPConnectorResponse, *fusionauth.Errors, error) {
	return makeLDAPConnectorRequest(ctx, client, connectorID, request, http.MethodPut)
}

func makeLDAPConnectorRequest(ctx context.Context, client fusionauth.FusionAuthClient, connectorID string, request LDAPConnectorRequest, method string) (*LDAPConnectorResponse, *fusionauth.Errors, error) {
	var resp LDAPConnectorResponse
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
