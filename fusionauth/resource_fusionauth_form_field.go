package fusionauth

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceFormField() *schema.Resource {
	return &schema.Resource{
		CreateContext: createFormField,
		ReadContext:   readFormField,
		UpdateContext: updateFormField,
		DeleteContext: deleteFormField,
		Schema: map[string]*schema.Schema{
			"form_field_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id to use for the new Form Field. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"confirm": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if the user input should be confirmed by requiring the value to be entered twice. If true, a confirmation field is included.",
			},
			"consent_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of an existing Consent. This field will be required when the type is set to consent.",
				ValidateFunc: validation.IsUUID,
			},
			"control": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The Form Field control",
				ValidateFunc: validation.StringInSlice([]string{
					"checkbox",
					"number",
					"password",
					"radio",
					"select",
					"textarea",
					"text",
				}, false),
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that can hold any information about the Form Field that should be persisted.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the Form Field.",
			},
			"key": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The key is the path to the value in the user or registration object.",
				ValidateFunc: validateKey,
				ForceNew:     true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Form Field.",
			},
			"options": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A list of options that are applied to checkbox, radio, or select controls.",
			},
			"required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if a value is required to complete the form.",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "string",
				ValidateFunc: validation.StringInSlice([]string{
					"bool",
					"consent",
					"date",
					"email",
					"number",
					"string",
				}, false),
				Description: "The data type used to store the value in FusionAuth.",
			},
			"validator": {
				Type:       schema.TypeList,
				MaxItems:   1,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines if user input should be validated.",
						},
						"expression": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "A regular expression used to validate user input. Must be a valid regular expression pattern.",
							ValidateFunc: validateRegex,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func createFormField(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	f := buildFormField(data)
	var fid string
	if fi, ok := data.GetOk("form_field_id"); ok {
		fid = fi.(string)
	}
	resp, faErrs, err := client.FAClient.CreateFormField(fid, fusionauth.FormFieldRequest{Field: f})
	if err != nil {
		return diag.Errorf("createFormField err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(resp.Field.Id)
	return buildResourceDataFromFormField(data, resp.Field)
}

func readFormField(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, err := client.FAClient.RetrieveFormField(id)
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

	return buildResourceDataFromFormField(data, resp.Field)
}

func updateFormField(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	f := buildFormField(data)

	resp, faErrs, err := client.FAClient.UpdateFormField(data.Id(), fusionauth.FormFieldRequest{Field: f})
	if err != nil {
		return diag.Errorf("UpdateFormField err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(resp.Field.Id)
	return buildResourceDataFromFormField(data, resp.Field)
}

func deleteFormField(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteFormField(id)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildFormField(data *schema.ResourceData) fusionauth.FormField {
	return fusionauth.FormField{
		Confirm:     data.Get("confirm").(bool),
		ConsentId:   data.Get("consent_id").(string),
		Control:     fusionauth.FormControl(data.Get("control").(string)),
		Data:        data.Get("data").(map[string]interface{}),
		Description: data.Get("description").(string),
		Key:         data.Get("key").(string),
		Name:        data.Get("name").(string),
		Options:     handleStringSliceFromList(data.Get("options").([]interface{})),
		Required:    data.Get("required").(bool),
		Type:        fusionauth.FormDataType(data.Get("type").(string)),
		Validator: fusionauth.FormFieldValidator{
			Enableable: buildEnableable("validator.0.enabled", data),
			Expression: data.Get("validator.0.expression").(string),
		},
	}
}

func buildResourceDataFromFormField(data *schema.ResourceData, f fusionauth.FormField) diag.Diagnostics {
	if err := data.Set("confirm", f.Confirm); err != nil {
		return diag.Errorf("webhook.confirm: %s", err.Error())
	}
	if err := data.Set("consent_id", f.ConsentId); err != nil {
		return diag.Errorf("webhook.consent_id: %s", err.Error())
	}
	if err := data.Set("control", f.Control); err != nil {
		return diag.Errorf("webhook.control: %s", err.Error())
	}
	if err := data.Set("data", f.Data); err != nil {
		return diag.Errorf("webhook.data: %s", err.Error())
	}
	if err := data.Set("description", f.Description); err != nil {
		return diag.Errorf("webhook.description: %s", err.Error())
	}
	if err := data.Set("key", f.Key); err != nil {
		return diag.Errorf("webhook.key: %s", err.Error())
	}
	if err := data.Set("name", f.Name); err != nil {
		return diag.Errorf("webhook.name: %s", err.Error())
	}
	if err := data.Set("options", f.Options); err != nil {
		return diag.Errorf("webhook.options: %s", err.Error())
	}
	if err := data.Set("required", f.Required); err != nil {
		return diag.Errorf("webhook.required: %s", err.Error())
	}
	if err := data.Set("type", f.Type); err != nil {
		return diag.Errorf("webhook.type: %s", err.Error())
	}

	err := data.Set("validator", []map[string]interface{}{
		{
			"enabled":    f.Validator.Enabled,
			"expression": f.Validator.Expression,
		},
	})
	if err != nil {
		return diag.Errorf("form_field.validator: %s", err.Error())
	}
	return nil
}

func validateKey(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	switch v {
	case "registration.preferredLanguages",
		"registration.roles",
		"registration.timezone",
		"registration.username",
		"user.birthDate",
		"user.email",
		"user.firstName",
		"user.fullName",
		"user.imageUrl",
		"user.lastName",
		"user.middleName",
		"user.mobilePhone",
		"user.password",
		"user.preferredLanguages",
		"user.timezone",
		"user.twoFactorEnabled",
		"user.username":
		return warnings, errors
	default:
		if !strings.HasPrefix(v, "user.data.") && !strings.HasPrefix(v, "registration.data.") {
			errors = append(
				errors,
				fmt.Errorf(
					"valid options for %q are: %q or start with %q",
					k,
					[]string{
						"registration.preferredLanguages",
						"registration.roles",
						"registration.timezone",
						"registration.username",
						"user.birthDate",
						"user.email",
						"user.firstName",
						"user.fullName",
						"user.imageUrl",
						"user.lastName",
						"user.middleName",
						"user.mobilePhone",
						"user.password",
						"user.preferredLanguages",
						"user.timezone",
						"user.twoFactorEnabled",
						"user.username",
					},
					[]string{
						"user.data.",
						"registration.data.",
					},
				),
			)
		}
	}
	return warnings, errors
}

func validateRegex(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := regexp.Compile(v); err != nil {
		return warnings, append(errors, err)
	}
	return warnings, errors
}
