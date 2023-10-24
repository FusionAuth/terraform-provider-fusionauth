package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newEmail() *schema.Resource {
	return &schema.Resource{
		CreateContext: createEmail,
		ReadContext:   readEmail,
		UpdateContext: updateEmail,
		DeleteContext: deleteEmail,
		Schema: map[string]*schema.Schema{
			"email_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id to use for the new Email Template. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"default_from_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The default From Name used when sending emails. If not provided, and a localized value cannot be determined, the default value for the tenant will be used. This is the display name part of the email address ( i.e. Jared Dunn <jared@piedpiper.com>).",
			},
			"default_html_template": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The default HTML Email Template.",
				DiffSuppressFunc: diffSuppressTemplate,
				ValidateFunc:     validation.StringIsNotWhiteSpace,
			},
			"default_subject": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The default Subject used when sending emails.",
			},
			"default_text_template": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The default Text Email Template.",
				DiffSuppressFunc: diffSuppressTemplate,
				ValidateFunc:     validation.StringIsNotWhiteSpace,
			},
			"from_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email address that this email will be sent from. If not provided, the default value for the tenant will be used. This is the address part email address (i.e. Jared Dunn <jared@piedpiper.com>).",
			},
			"localized_from_names": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The From Name used when sending emails to users who speak other languages. This overrides the default From Name based on the user’s list of preferred languages.",
			},
			"localized_html_templates": {
				Type:             schema.TypeMap,
				Optional:         true,
				Description:      "The HTML Email Template used when sending emails to users who speak other languages. This overrides the default HTML Email Template based on the user’s list of preferred languages.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"localized_subjects": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The Subject used when sending emails to users who speak other languages. This overrides the default Subject based on the user’s list of preferred languages.",
			},
			"localized_text_templates": {
				Type:             schema.TypeMap,
				Optional:         true,
				Description:      "The Text Email Template used when sending emails to users who speak other languages. This overrides the default Text Email Template based on the user’s list of preferred languages.",
				DiffSuppressFunc: diffSuppressTemplate,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  `A descriptive name for the email template (i.e. "April 2016 Coupon Email")`,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func buildEmail(data *schema.ResourceData) fusionauth.EmailTemplate {
	e := fusionauth.EmailTemplate{
		DefaultFromName:     data.Get("default_from_name").(string),
		DefaultHtmlTemplate: data.Get("default_html_template").(string),
		DefaultSubject:      data.Get("default_subject").(string),
		DefaultTextTemplate: data.Get("default_text_template").(string),
		FromEmail:           data.Get("from_email").(string),
		Name:                data.Get("name").(string),
	}

	if i, ok := data.GetOk("localized_from_names"); ok {
		e.LocalizedFromNames = intMapToStringMap(i.(map[string]interface{}))
	}

	if i, ok := data.GetOk("localized_html_templates"); ok {
		e.LocalizedHtmlTemplates = intMapToStringMap(i.(map[string]interface{}))
	}

	if i, ok := data.GetOk("localized_subjects"); ok {
		e.LocalizedSubjects = intMapToStringMap(i.(map[string]interface{}))
	}

	if i, ok := data.GetOk("localized_text_templates"); ok {
		e.LocalizedTextTemplates = intMapToStringMap(i.(map[string]interface{}))
	}
	return e
}

func createEmail(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	e := buildEmail(data)

	var eid string
	if ei, ok := data.GetOk("email_id"); ok {
		eid = ei.(string)
	}

	resp, faErrs, err := client.FAClient.CreateEmailTemplate(eid, fusionauth.EmailTemplateRequest{
		EmailTemplate: e,
	})

	if err != nil {
		return diag.Errorf("CreateEmailTemplate err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.EmailTemplate.Id)
	return nil
}

func readEmail(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, err := client.FAClient.RetrieveEmailTemplate(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	t := resp.EmailTemplate
	if err := data.Set("default_from_name", t.DefaultFromName); err != nil {
		return diag.Errorf("email.default_from_name: %s", err.Error())
	}
	if err := data.Set("default_html_template", t.DefaultHtmlTemplate); err != nil {
		return diag.Errorf("email.default_html_template: %s", err.Error())
	}
	if err := data.Set("default_subject", t.DefaultSubject); err != nil {
		return diag.Errorf("email.default_subject: %s", err.Error())
	}
	if err := data.Set("default_text_template", t.DefaultTextTemplate); err != nil {
		return diag.Errorf("email.default_text_template: %s", err.Error())
	}
	if err := data.Set("from_email", t.FromEmail); err != nil {
		return diag.Errorf("email.from_email: %s", err.Error())
	}
	if err := data.Set("localized_from_names", t.LocalizedFromNames); err != nil {
		return diag.Errorf("email.localized_from_names: %s", err.Error())
	}
	if err := data.Set("localized_html_templates", t.LocalizedHtmlTemplates); err != nil {
		return diag.Errorf("email.localized_html_templates: %s", err.Error())
	}
	if err := data.Set("localized_subjects", t.LocalizedSubjects); err != nil {
		return diag.Errorf("email.localized_subjects: %s", err.Error())
	}
	if err := data.Set("localized_text_templates", t.LocalizedTextTemplates); err != nil {
		return diag.Errorf("email.localized_text_templates: %s", err.Error())
	}
	if err := data.Set("name", t.Name); err != nil {
		return diag.Errorf("email.name: %s", err.Error())
	}

	return nil
}

func updateEmail(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	e := buildEmail(data)

	resp, faErrs, err := client.FAClient.UpdateEmailTemplate(data.Id(), fusionauth.EmailTemplateRequest{
		EmailTemplate: e,
	})

	if err != nil {
		return diag.Errorf("UpdateEmailTemplate err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteEmail(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteEmailTemplate(id)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
