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

func dataSourceConsentRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	var resp *fusionauth.ConsentResponse
	var retrieveError error

	if id, ok := data.GetOk("consent_id"); ok {
		resp, retrieveError = client.FAClient.RetrieveConsent(id.(string))
		if retrieveError != nil {
			return diag.Errorf("Error retrieving consent with id %s: %s", id, retrieveError)
		}
	} else if name, ok := data.GetOk("name"); ok {
		// Retrieve all consents and find the one with matching name
		searchName := name.(string)
		consentsResp, retrieveError := client.FAClient.RetrieveConsents()
		if retrieveError != nil {
			return diag.Errorf("Error retrieving consents: %s", retrieveError)
		}

		if err := checkResponse(consentsResp.StatusCode, nil); err != nil {
			return diag.FromErr(err)
		}

		// Find consent with matching name
		found := false
		for _, c := range consentsResp.Consents {
			if c.Name == searchName {
				// Once found, get the full consent details
				resp, retrieveError = client.FAClient.RetrieveConsent(c.Id)
				if retrieveError != nil {
					return diag.Errorf("Error retrieving consent with id %s: %s", c.Id, retrieveError)
				}
				found = true
				break
			}
		}

		if !found {
			return diag.Errorf("Couldn't find Consent with name '%s'", searchName)
		}
	} else {
		return diag.Errorf("Either 'consent_id' or 'name' must be specified")
	}

	if resp.StatusCode == http.StatusNotFound {
		return diag.Errorf("Consent not found")
	}

	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	consent := resp.Consent
	data.SetId(consent.Id)

	if err := data.Set("name", consent.Name); err != nil {
		return diag.Errorf("error setting consent.name: %s", err.Error())
	}

	if err := data.Set("consent_email_template_id", consent.ConsentEmailTemplateId); err != nil {
		return diag.Errorf("error setting consent.consent_email_template_id: %s", err.Error())
	}

	if err := data.Set("country_minimum_age_for_self_consent", consent.CountryMinimumAgeForSelfConsent); err != nil {
		return diag.Errorf("error setting consent.country_minimum_age_for_self_consent: %s", err.Error())
	}

	if err := data.Set("default_minimum_age_for_self_consent", consent.DefaultMinimumAgeForSelfConsent); err != nil {
		return diag.Errorf("error setting consent.default_minimum_age_for_self_consent: %s", err.Error())
	}

	if err := data.Set("multiple_values_allowed", consent.MultipleValuesAllowed); err != nil {
		return diag.Errorf("error setting consent.multiple_values_allowed: %s", err.Error())
	}

	if err := data.Set("values", consent.Values); err != nil {
		return diag.Errorf("error setting consent.values: %s", err.Error())
	}

	dataJSON, diags := mapStringInterfaceToJSONString(consent.Data)
	if diags != nil {
		return diags
	}
	if err := data.Set("data", dataJSON); err != nil {
		return diag.Errorf("error setting consent.data: %s", err.Error())
	}

	if err := data.Set("email_plus", []interface{}{
		map[string]interface{}{
			"email_template_id":                   consent.EmailPlus.EmailTemplateId,
			"enabled":                             consent.EmailPlus.Enabled,
			"maximum_time_to_send_email_in_hours": consent.EmailPlus.MaximumTimeToSendEmailInHours,
			"minimum_time_to_send_email_in_hours": consent.EmailPlus.MinimumTimeToSendEmailInHours,
		},
	}); err != nil {
		return diag.Errorf("error setting consent.email_plus: %s", err.Error())
	}

	return nil
}
