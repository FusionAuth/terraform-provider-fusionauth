package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceConsent() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConsentRead,
		Schema: map[string]*schema.Schema{
			// Data Source Parameters
			"consent_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"consent_id", "name"},
				Description:  "The unique Id of the Consent to retrieve.",
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"consent_id", "name"},
				Description:  "The case-insensitive string to search for in the Consent name.",
			},
			// Data Source Attributes
			"consent_email_template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Id of the Email Template that is used to send confirmation to the end user.",
			},
			"country_minimum_age_for_self_consent": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "This property optionally overrides the value provided in defaultMinimumAgeForSelfConsent if a more specific value is defined.",
			},
			"data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An object that can hold any information about the Consent that should be persisted. Must be a JSON string.",
			},
			"default_minimum_age_for_self_consent": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The default age of self consent used when granting this consent to a user unless a more specific one is provided by the countryMinimumAgeForSelfConsent.",
			},
			"email_plus": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Email Plus provides and additional opportunity to notify the giver that consent was provided.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email_template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Id of the Email Template that is used to send the reminder notice to the consent giver.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Set to true when the Email Plus workflow is enabled.",
						},
						"maximum_time_to_send_email_in_hours": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum number of hours to wait until sending the reminder notice the consent giver.",
						},
						"minimum_time_to_send_email_in_hours": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum number of hours to wait until sending the reminder notice the consent giver.",
						},
					},
				},
			},
			"multiple_values_allowed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Set to true if more than one value may be used when granting this consent to a User.",
			},
			"values": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of values that may be assigned to this consent.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceConsentRead(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	// Retrieve consent based on provided identifiers
	consent, diags := retrieveConsent(client, data)
	if diags != nil {
		return diags
	}

	// Set the resource ID
	data.SetId(consent.Id)

	// Set all the consent properties
	return setConsentProperties(data, consent)
}

// retrieveConsent gets a consent using either ID or name
func retrieveConsent(client Client, data *schema.ResourceData) (*fusionauth.Consent, diag.Diagnostics) {
	if id, ok := data.GetOk("consent_id"); ok {
		return retrieveConsentByID(client, id.(string))
	}

	if name, ok := data.GetOk("name"); ok {
		return retrieveConsentByName(client, name.(string))
	}

	return nil, diag.Errorf("Either 'consent_id' or 'name' must be specified")
}

// retrieveConsentByID retrieves a consent using its ID
func retrieveConsentByID(client Client, id string) (*fusionauth.Consent, diag.Diagnostics) {
	resp, err := client.FAClient.RetrieveConsent(id)
	if err != nil {
		return nil, diag.Errorf("Error retrieving consent with id %s: %s", id, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, diag.Errorf("Consent not found")
	}

	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return nil, diag.FromErr(err)
	}

	return &resp.Consent, nil
}

// retrieveConsentByName retrieves a consent using its name
func retrieveConsentByName(client Client, name string) (*fusionauth.Consent, diag.Diagnostics) {
	consentsResp, err := client.FAClient.RetrieveConsents()
	if err != nil {
		return nil, diag.Errorf("Error retrieving consents: %s", err)
	}

	if err := checkResponse(consentsResp.StatusCode, nil); err != nil {
		return nil, diag.FromErr(err)
	}

	// Find consent with matching name
	for _, c := range consentsResp.Consents {
		if c.Name == name {
			// Once found, get the full consent details
			return retrieveConsentByID(client, c.Id)
		}
	}

	return nil, diag.Errorf("Couldn't find Consent with name '%s'", name)
}

// setConsentProperties sets all properties from the consent to the resource data
func setConsentProperties(data *schema.ResourceData, consent *fusionauth.Consent) diag.Diagnostics {
	// Helper function to reduce repeated error handling code
	setField := func(key string, value interface{}) diag.Diagnostics {
		if err := data.Set(key, value); err != nil {
			return diag.Errorf("error setting consent.%s: %s", key, err.Error())
		}
		return nil
	}

	// Set all fields
	if diags := setField("name", consent.Name); diags != nil {
		return diags
	}

	if diags := setField("consent_email_template_id", consent.ConsentEmailTemplateId); diags != nil {
		return diags
	}

	if diags := setField("country_minimum_age_for_self_consent", consent.CountryMinimumAgeForSelfConsent); diags != nil {
		return diags
	}

	if diags := setField("default_minimum_age_for_self_consent", consent.DefaultMinimumAgeForSelfConsent); diags != nil {
		return diags
	}

	if diags := setField("multiple_values_allowed", consent.MultipleValuesAllowed); diags != nil {
		return diags
	}

	if diags := setField("values", consent.Values); diags != nil {
		return diags
	}

	// Handle the JSON data conversion
	dataJSON, diags := mapStringInterfaceToJSONString(consent.Data)
	if diags != nil {
		return diags
	}

	if diags := setField("data", dataJSON); diags != nil {
		return diags
	}

	// Set the email_plus nested object
	emailPlus := []interface{}{
		map[string]interface{}{
			"email_template_id":                   consent.EmailPlus.EmailTemplateId,
			"enabled":                             consent.EmailPlus.Enabled,
			"maximum_time_to_send_email_in_hours": consent.EmailPlus.MaximumTimeToSendEmailInHours,
			"minimum_time_to_send_email_in_hours": consent.EmailPlus.MinimumTimeToSendEmailInHours,
		},
	}

	if diags := setField("email_plus", emailPlus); diags != nil {
		return diags
	}

	return nil
}
