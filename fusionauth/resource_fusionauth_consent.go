package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newConsent() *schema.Resource {
	return &schema.Resource{
		CreateContext: createConsent,
		ReadContext:   readConsent,
		UpdateContext: updateConsent,
		DeleteContext: deleteConsent,
		Schema: map[string]*schema.Schema{
			"consent_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template that is used to send confirmation to the end user. If this value is omitted an email will not be sent to the user.",
				ValidateFunc: validation.IsUUID,
			},
			"consent_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The Id to use for the new Consent. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
			},
			"country_minimum_age_for_self_consent": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "This property optionally overrides the value provided in defaultMinimumAgeForSelfConsent if a more specific value is defined. This can be useful when the age of self consent varies by country.",
			},
			"data": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "An object that can hold any information about the Consent that should be persisted. Please review the limits on data field types as you plan for and build your custom data schema. Must be a JSON string.",
				DiffSuppressFunc: diffSuppressJSON,
				ValidateFunc:     validation.StringIsJSON,
			},
			"default_minimum_age_for_self_consent": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The default age of self consent used when granting this consent to a user unless a more specific one is provided by the countryMinimumAgeForSelfConsent. A user that meets the minimum age of self consent may self-consent, this means the recipient may also be the giver.",
			},
			"email_plus": {
				Type:             schema.TypeList,
				MaxItems:         1,
				Optional:         true,
				ConfigMode:       schema.SchemaConfigModeAttr,
				DiffSuppressFunc: suppressBlockDiff,
				Description:      "Email Plus provides and additional opportunity to notify the giver that consent was provided. For example, if consentEmailTemplateId is provided then when the consent is granted an email will be sent to notify the giver that consent was granted to the user. When using Email Plus a follow up email will be sent to the giver at a randomly selected time within the configured minimum and maximum range of hours.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email_template_id": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The Id of the Email Template that is used to send the reminder notice to the consent giver. This value is required when emailPlus.enabled is set to true.",
							ValidateFunc: validation.IsUUID,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Set this value to true to enable the Email Plus workflow.",
						},
						"maximum_time_to_send_email_in_hours": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     48,
							Description: "The maximum number of hours to wait until sending the reminder notice the consent giver.",
						},
						"minimum_time_to_send_email_in_hours": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     24,
							Description: "The minimum number of hours to wait until sending the reminder notice the consent giver.",
						},
					},
				},
			},
			"multiple_values_allowed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set this value to true if more than one value may be used when granting this consent to a User. When this value is false a maximum of one value may be assigned. This value is not used when no values have been defined for this consent.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the consent.",
			},
			"values": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of values that may be assigned to this consent. This is a required field when multipleValuesAllowed is set to true. If this value is not provided then the consent will be created with no values.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func buildConsent(data *schema.ResourceData) fusionauth.Consent {
	countryMinAgeMap := make(map[string]int)
	for country, ageInterface := range data.Get("country_minimum_age_for_self_consent").(map[string]interface{}) {
		if age, ok := ageInterface.(int); ok {
			countryMinAgeMap[country] = age
		}
	}
	resourceData, _ := jsonStringToMapStringInterface(data.Get("data").(string))
	consent := fusionauth.Consent{
		ConsentEmailTemplateId:          data.Get("consent_email_template_id").(string),
		CountryMinimumAgeForSelfConsent: countryMinAgeMap,
		Data:                            resourceData,
		DefaultMinimumAgeForSelfConsent: data.Get("default_minimum_age_for_self_consent").(int),
		EmailPlus: fusionauth.EmailPlus{
			Enableable:                    buildEnableable("email_plus.0.enabled", data),
			EmailTemplateId:               data.Get("email_plus.0.email_template_id").(string),
			MaximumTimeToSendEmailInHours: data.Get("email_plus.0.maximum_time_to_send_email_in_hours").(int),
			MinimumTimeToSendEmailInHours: data.Get("email_plus.0.minimum_time_to_send_email_in_hours").(int),
		},
		Id:                    data.Get("consent_id").(string),
		MultipleValuesAllowed: data.Get("multiple_values_allowed").(bool),
		Name:                  data.Get("name").(string),
		Values:                handleStringSliceFromList(data.Get("values").([]interface{})),
	}

	return consent
}

func createConsent(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	consent := buildConsent(data)

	resp, faErrs, err := client.FAClient.CreateConsent(consent.Id, fusionauth.ConsentRequest{
		Consent: consent,
	})
	if err != nil {
		return diag.Errorf("CreateConsent err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(resp.Consent.Id)
	return nil
}

func readConsent(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, err := client.FAClient.RetrieveConsent(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}

	consent := resp.Consent

	if err := data.Set("consent_email_template_id", consent.ConsentEmailTemplateId); err != nil {
		return diag.Errorf("consent.consent_email_template_id: %s", err.Error())
	}
	if err := data.Set("country_minimum_age_for_self_consent", consent.CountryMinimumAgeForSelfConsent); err != nil {
		return diag.Errorf("consent.country_minimum_age_for_self_consent: %s", err.Error())
	}
	dataJSON, diags := mapStringInterfaceToJSONString(consent.Data)
	if diags != nil {
		return diags
	}
	err = data.Set("data", dataJSON)
	if err != nil {
		return diag.Errorf("consent.data: %s", err.Error())
	}
	if err := data.Set("default_minimum_age_for_self_consent", consent.DefaultMinimumAgeForSelfConsent); err != nil {
		return diag.Errorf("consent.default_minimum_age_for_self_consent: %s", err.Error())
	}
	if err := data.Set("email_plus", []interface{}{
		map[string]interface{}{
			"email_template_id":                   consent.EmailPlus.EmailTemplateId,
			"enabled":                             consent.EmailPlus.Enabled,
			"maximum_time_to_send_email_in_hours": consent.EmailPlus.MaximumTimeToSendEmailInHours,
			"minimum_time_to_send_email_in_hours": consent.EmailPlus.MinimumTimeToSendEmailInHours,
		},
	}); err != nil {
		return diag.Errorf("consent.email_plus: %s", err.Error())
	}
	if err := data.Set("multiple_values_allowed", consent.MultipleValuesAllowed); err != nil {
		return diag.Errorf("consent.multiple_values_allowed: %s", err.Error())
	}
	if err := data.Set("name", consent.Name); err != nil {
		return diag.Errorf("consent.name: %s", err.Error())
	}
	if err := data.Set("values", consent.Values); err != nil {
		return diag.Errorf("consent.values: %s", err.Error())
	}

	return nil
}

func updateConsent(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	consent := buildConsent(data)

	resp, faErrs, err := client.FAClient.UpdateConsent(data.Id(), fusionauth.ConsentRequest{
		Consent: consent,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteConsent(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteConsent(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
