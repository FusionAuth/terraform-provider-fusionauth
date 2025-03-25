package fusionauth

import (
	"context"
	"log"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newSMSMessageTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: createSMSMessageTemplate,
		ReadContext:   readSMSMessageTemplate,
		UpdateContext: updateSMSMessageTemplate,
		DeleteContext: deleteSMSMessageTemplate,
		Schema: map[string]*schema.Schema{
			"message_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The Id to use for the new Message Template. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"data": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "An object that can hold any information about the Message Template that should be persisted. Must be a JSON string.",
				DiffSuppressFunc: diffSuppressJSON,
				ValidateFunc:     validation.StringIsJSON,
			},
			"default_template": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The default Message Template.",
			},
			"localized_templates": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The Message Template used when sending messages to users who speak other languages. This overrides the default Message Template based on the user's list of preferred languages.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A descriptive name for the Message Template (i.e. \"Two Factor Code Message\")",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of Message Template. This is always 'SMS'.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func buildSMSMessageTemplate(data *schema.ResourceData) fusionauth.SMSMessageTemplate {
	template := fusionauth.SMSMessageTemplate{
		MessageTemplate: fusionauth.MessageTemplate{
			Id:   data.Get("message_template_id").(string),
			Name: data.Get("name").(string),
			Type: fusionauth.MessageType_SMS,
		},
		DefaultTemplate: data.Get("default_template").(string),
	}

	if i, ok := data.GetOk("data"); ok {
		resourceData, _ := jsonStringToMapStringInterface(i.(string))
		template.Data = resourceData
	}

	if i, ok := data.GetOk("localized_templates"); ok {
		template.LocalizedTemplates = intMapToStringMap(i.(map[string]interface{}))
	}

	return template
}

func createSMSMessageTemplate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	template := buildSMSMessageTemplate(data)
	resp, faErrs, err := createMessageTemplate(ctx, client.FAClient, template.Id, SMSMessageTemplateRequest{MessageTemplate: template})
	if err != nil {
		return diag.Errorf("CreateSMSMessageTemplate err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(resp.MessageTemplate.Id)
	return nil
}

func readSMSMessageTemplate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := retrieveMessageTemplate(ctx, client.FAClient, id)
	if err != nil {
		return diag.FromErr(err)
	}

	// Log the response for debugging
	log.Printf("[DEBUG] readSMSMessageTemplate response: %+v", resp)

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	template := resp.MessageTemplate
	dataJSON, diags := mapStringInterfaceToJSONString(template.Data)
	if diags != nil {
		return diags
	}
	if err := data.Set("data", dataJSON); err != nil {
		return diag.Errorf("message_template.data: %s", err.Error())
	}
	if err := data.Set("default_template", template.DefaultTemplate); err != nil {
		return diag.Errorf("message_template.default_template: %s", err.Error())
	}
	if err := data.Set("localized_templates", template.LocalizedTemplates); err != nil {
		return diag.Errorf("message_template.localized_templates: %s", err.Error())
	}
	if err := data.Set("message_template_id", template.Id); err != nil {
		return diag.Errorf("message_template.message_template_id: %s", err.Error())
	}
	if err := data.Set("name", template.Name); err != nil {
		return diag.Errorf("message_template.name: %s", err.Error())
	}
	if err := data.Set("type", template.Type); err != nil {
		return diag.Errorf("message_template.type: %s", err.Error())
	}

	return nil
}

func updateSMSMessageTemplate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	template := buildSMSMessageTemplate(data)

	resp, faErrs, err := updateMessageTemplate(ctx, client.FAClient, data.Id(), SMSMessageTemplateRequest{MessageTemplate: template})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteSMSMessageTemplate(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteMessageTemplate(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

type SMSMessageTemplateRequest struct {
	MessageTemplate fusionauth.SMSMessageTemplate `json:"messageTemplate,omitempty"`
}

type SMSMessageTemplateResponse struct {
	fusionauth.BaseHTTPResponse
	MessageTemplate  fusionauth.SMSMessageTemplate   `json:"messageTemplate,omitempty"`
	MessageTemplates []fusionauth.SMSMessageTemplate `json:"messageTemplates,omitempty"`
}

func (b *SMSMessageTemplateResponse) SetStatus(status int) {
	b.StatusCode = status
}

func createMessageTemplate(ctx context.Context, client fusionauth.FusionAuthClient, templateID string, request SMSMessageTemplateRequest) (*SMSMessageTemplateResponse, *fusionauth.Errors, error) {
	return makeMessageTemplateRequest(ctx, client, templateID, request, http.MethodPost)
}

func retrieveMessageTemplate(ctx context.Context, client fusionauth.FusionAuthClient, templateID string) (*SMSMessageTemplateResponse, *fusionauth.Errors, error) {
	return makeMessageTemplateRequest(ctx, client, templateID, SMSMessageTemplateRequest{}, http.MethodGet)
}

func updateMessageTemplate(ctx context.Context, client fusionauth.FusionAuthClient, templateID string, request SMSMessageTemplateRequest) (*SMSMessageTemplateResponse, *fusionauth.Errors, error) {
	return makeMessageTemplateRequest(ctx, client, templateID, request, http.MethodPut)
}

func makeMessageTemplateRequest(ctx context.Context, client fusionauth.FusionAuthClient, templateID string, request SMSMessageTemplateRequest, method string) (*SMSMessageTemplateResponse, *fusionauth.Errors, error) {
	var resp SMSMessageTemplateResponse
	var errors fusionauth.Errors

	restClient := client.Start(&resp, &errors)

	/* Pattern elsewhere is to include an empty Request object for GET requests:
	   With the /api/message/template/{id} endpoint, this causes a 400 error.
	*/
	if method != http.MethodGet {
		restClient.WithJSONBody(request)
	}

	err := restClient.WithUri("/api/message/template").
		WithUriSegment(templateID).
		WithMethod(method).
		Do(ctx)
	if restClient.ErrorRef == nil {
		return &resp, nil, err
	}
	return &resp, &errors, err
}
