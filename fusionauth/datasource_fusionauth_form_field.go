package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceFormField() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFormFieldRead,
		Schema: map[string]*schema.Schema{
			"form_field_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"form_field_id", "name"},
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
				Optional:     true,
				Description:  "The key is the path to the value in the user or registration object.",
				ValidateFunc: validateKey,
				ForceNew:     true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"form_field_id", "name"},
				Description:  "The unique name of the Form Field.",
			},
			"options": {
				Type:        schema.TypeSet,
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
	}
}

func dataSourceFormFieldRead(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	var searchTerm string
	var res *fusionauth.FormFieldResponse
	var err error

	// Either `form_field_id` or `name` are guaranteed to be set
	if entityID, ok := data.GetOk("form_field_id"); ok {
		searchTerm = entityID.(string)
		res, err = client.FAClient.RetrieveFormField(searchTerm)
	} else {
		searchTerm = data.Get("name").(string)
		res, err = client.FAClient.RetrieveFormFields()
	}
	if err != nil {
		return diag.FromErr(err)
	}
	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("couldn't find form field '%s'", searchTerm)
	}
	if err := checkResponse(res.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	foundEntity := res.Field
	if len(res.Fields) > 0 {
		// search based on name
		var found = false
		for _, entity := range res.Fields {
			if entity.Name == searchTerm {
				found = true
				foundEntity = entity
				break
			}
		}
		if !found {
			return diag.Errorf("couldn't find form field with name '%s'", searchTerm)
		}
	}

	data.SetId(foundEntity.Id)
	return buildResourceDataFromFormField(data, foundEntity)
}
