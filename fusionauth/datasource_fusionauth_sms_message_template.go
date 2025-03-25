package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceSMSMessageTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSMSMessageTemplateRead,
		Schema: map[string]*schema.Schema{
			// Data Source Parameters
			"message_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"message_template_id", "name"},
				Description:  "The unique Id of the SMS Message Template to retrieve.",
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"message_template_id", "name"},
				Description:  "The case-insensitive string to search for in the SMS Message Template name.",
			},
			// Data Source Attributes
			"data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An object that can hold any information about the Message Template that should be persisted. Must be a JSON string.",
			},
			"default_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default Message Template.",
			},
			"localized_templates": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The Message Template used when sending messages to users who speak other languages. This overrides the default Message Template based on the user's list of preferred languages.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of Message Template. This is always 'SMS'.",
			},
		},
	}
}

func dataSourceSMSMessageTemplateRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	var searchTerm string
	var resp *SMSMessageTemplateResponse
	var err error

	if entityID, ok := data.GetOk("message_template_id"); ok {
		searchTerm = entityID.(string)
		resp, _, err = retrieveMessageTemplate(ctx, client.FAClient, searchTerm)
	} else if name, ok := data.GetOk("name"); ok {
		searchTerm = name.(string)
		resp, _, err = retrieveAllMessageTemplates(ctx, client.FAClient)
	} else {
		return diag.Errorf("Either 'message_template_id' or 'name' must be specified")
	}

	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	var template fusionauth.SMSMessageTemplate
	if resp.MessageTemplate.Id != "" {
		template = resp.MessageTemplate
	} else if len(resp.MessageTemplates) > 0 {
		found := false
		for _, t := range resp.MessageTemplates {
			if t.Name == searchTerm && t.Type == fusionauth.MessageType_SMS {
				template = t
				found = true
				break
			}
		}
		if !found {
			return diag.Errorf("Couldn't find SMS Message Template with name '%s'", searchTerm)
		}
	} else {
		return diag.Errorf("No SMS Message Templates found matching '%s'", searchTerm)
	}

	data.SetId(template.Id)

	if diags := setSMSMessageTemplateFields(data, template); diags != nil {
		return diags
	}

	return nil
}

// Helper function to set all SMS message template fields in the schema data
func setSMSMessageTemplateFields(data *schema.ResourceData, template fusionauth.SMSMessageTemplate) diag.Diagnostics {
	dataJSON, diags := mapStringInterfaceToJSONString(template.Data)
	if diags != nil {
		return diags
	}

	fields := map[string]interface{}{
		"message_template_id": template.Id,
		"name":                template.Name,
		"default_template":    template.DefaultTemplate,
		"localized_templates": template.LocalizedTemplates,
		"type":                template.Type,
		"data":                dataJSON,
	}

	for key, value := range fields {
		if err := data.Set(key, value); err != nil {
			return diag.Errorf("error setting message_template.%s: %s", key, err.Error())
		}
	}

	return nil
}

// Function to retrieve all message templates
func retrieveAllMessageTemplates(ctx context.Context, client fusionauth.FusionAuthClient) (*SMSMessageTemplateResponse, *fusionauth.Errors, error) {
	return makeMessageTemplateRequest(ctx, client, "", SMSMessageTemplateRequest{}, http.MethodGet)
}
