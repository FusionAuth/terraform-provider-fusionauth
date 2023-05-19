package fusionauth

import (
	"context"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceEmail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEmailRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the Email Template.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique Id of the Email Template",
			},
			"default_from_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default From Name used when sending emails.",
			},
			"default_html_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default HTML Email Template.",
			},
			"default_subject": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default Subject used when sending emails.",
			},
			"default_text_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default Text Email Template.",
			},
			"from_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email address that this email will be sent from.",
			},
			"localized_from_names": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The From Name used when sending emails to users who speak other languages.",
			},
			"localized_html_templates": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The HTML Email Template used when sending emails to users who speak other languages.",
			},
			"localized_subjects": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The Subject used when sending emails to users who speak other languages.",
			},
			"localized_text_templates": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The Text Email Template used when sending emails to users who speak other languages.",
			},
		},
	}
}

func dataSourceEmailRead(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	resp, err := client.FAClient.RetrieveEmailTemplates()
	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}
	name := data.Get("name").(string)
	var t *fusionauth.EmailTemplate

	if len(resp.EmailTemplates) > 0 {
		for i := range resp.EmailTemplates {
			if resp.EmailTemplates[i].Name == name {
				t = &resp.EmailTemplates[i]
				break
			}
		}
	}
	if t == nil {
		return diag.Errorf("couldn't find email template %s", name)
	}
	data.SetId(t.Id)
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
