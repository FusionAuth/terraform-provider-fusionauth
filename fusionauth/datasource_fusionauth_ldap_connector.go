package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceLDAPConnector() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLDAPConnectorRead,
		Schema: map[string]*schema.Schema{
			// Data Source Parameters
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The unique Id of the LDAP Connector to retrieve.",
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The case-insensitive string to search for in the LDAP Connector name.",
			},
			// Data Source Attributes
			"authentication_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fully qualified LDAP URL to authenticate.",
			},
			"base_structure": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The top of the LDAP directory hierarchy. Typically this contains the dc (domain component) element.",
			},
			"connect_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The connect timeout for the HTTP connection, in milliseconds. Value must be greater than 0.",
			},
			"data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An object that can hold any information about the Connector that should be persisted. Must be a JSON string.",
			},
			"debug": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines if debug should be enabled to create an event log to assist in debugging integration errors.",
			},
			"identifying_attribute": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The entry attribute name which is the first component of the distinguished name of entries in the directory.",
			},
			"lambda_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reconcile_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Id of an existing Lambda. The lambda is executed after the user authenticates with the connector.",
						},
					},
				},
			},
			"login_id_attribute": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The entity attribute name which stores the identifier that is used for logging the user in.",
			},
			"read_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The read timeout for the HTTP connection, in milliseconds. Value must be greater than 0.",
			},
			"requested_attributes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The list of attributes to request from the LDAP server. This is a list of strings.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"security_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The LDAP security method. Possible values are: None, LDAPS, or StartTLS.",
			},
			"system_account_dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The distinguished name of the system account used to authenticate to the LDAP server.",
			},
			"system_account_password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The password of an entry that has read access to the directory.",
			},
		},
	}
}

func dataSourceLDAPConnectorRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	var searchTerm string
	var resp *LDAPConnectorResponse
	var err error

	if entityID, ok := data.GetOk("id"); ok {
		searchTerm = entityID.(string)
		resp, _, err = RetrieveLDAPConnector(ctx, client.FAClient, searchTerm)
	} else if name, ok := data.GetOk("name"); ok {
		searchTerm = name.(string)
		resp, _, err = RetrieveLDAPConnectors(ctx, client.FAClient)
	} else {
		return diag.Errorf("Either 'id' or 'name' must be specified")
	}

	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	var connector fusionauth.LDAPConnectorConfiguration
	if resp.Connector.Id != "" {
		connector = resp.Connector
	} else if len(resp.Connectors) > 0 {
		found := false
		for _, c := range resp.Connectors {
			if c.Name == searchTerm && c.Type == "LDAP" {
				connector = c
				found = true
				break
			}
		}
		if !found {
			return diag.Errorf("Couldn't find LDAP Connector with name '%s'", searchTerm)
		}
	} else {
		return diag.Errorf("No LDAP Connectors found matching '%s'", searchTerm)
	}

	data.SetId(connector.Id)

	if diags := setLDAPConnectorFields(data, connector); diags != nil {
		return diags
	}

	return nil
}

// Helper function to set all connector fields in the schema data
func setLDAPConnectorFields(data *schema.ResourceData, connector fusionauth.LDAPConnectorConfiguration) diag.Diagnostics {
	fields := map[string]interface{}{
		"id":                      connector.Id,
		"name":                    connector.Name,
		"authentication_url":      connector.AuthenticationURL,
		"base_structure":          connector.BaseStructure,
		"connect_timeout":         connector.ConnectTimeout,
		"debug":                   connector.Debug,
		"identifying_attribute":   connector.IdentifyingAttribute,
		"login_id_attribute":      connector.LoginIdAttribute,
		"read_timeout":            connector.ReadTimeout,
		"requested_attributes":    connector.RequestedAttributes,
		"security_method":         connector.SecurityMethod,
		"system_account_dn":       connector.SystemAccountDN,
		"system_account_password": connector.SystemAccountPassword,
	}

	dataJSON, diags := mapStringInterfaceToJSONString(connector.Data)
	if diags != nil {
		return diags
	}
	fields["data"] = dataJSON

	if err := data.Set("lambda_configuration", []interface{}{
		map[string]interface{}{
			"reconcile_id": connector.LambdaConfiguration.ReconcileId,
		},
	}); err != nil {
		return diag.Errorf("error setting connector.lambda_configuration: %s", err.Error())
	}

	for key, value := range fields {
		if err := data.Set(key, value); err != nil {
			return diag.Errorf("error setting connector.%s: %s", key, err.Error())
		}
	}

	return nil
}

func RetrieveLDAPConnectors(ctx context.Context, client fusionauth.FusionAuthClient) (*LDAPConnectorResponse, *fusionauth.Errors, error) {
	return makeLDAPConnectorRequest(ctx, client, "", LDAPConnectorRequest{}, http.MethodGet)
}
