package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceForm() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFormRead,
		Schema: map[string]*schema.Schema{
			"form_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"form_id", "name"},
				Description:  "The unique Id of the Form.",
				ValidateFunc: validation.IsUUID,
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that can hold any information about the Form that should be persisted.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"form_id", "name"},
				Description:  "The unique name of the Form.",
			},
			"steps": {
				Type:        schema.TypeList,
				Description: "An ordered list of objects containing one or more Form Fields. A Form must have at least one step defined.",
				MinItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fields": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							MinItems:    1,
							Required:    true,
							Description: "An ordered list of Form Field Ids assigned to this step.",
						},
					},
				},
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of form being created, a form type cannot be changed after the form has been created.",
				Default:     "registration",
				ForceNew:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"registration",
					"adminRegistration",
					"adminUser",
					"selfServiceUser",
				}, false),
			},
		},
	}
}

func dataSourceFormRead(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	var searchTerm string
	var res *fusionauth.FormResponse
	var err error

	// Either `form_id` or `name` are guaranteed to be set
	if entityID, ok := data.GetOk("form_id"); ok {
		searchTerm = entityID.(string)
		res, err = client.FAClient.RetrieveForm(searchTerm)
	} else {
		searchTerm = data.Get("name").(string)
		res, err = client.FAClient.RetrieveForms()
	}
	if err != nil {
		return diag.FromErr(err)
	}
	if res.StatusCode == http.StatusNotFound {
		return diag.Errorf("couldn't find form '%s'", searchTerm)
	}
	if err := checkResponse(res.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	foundEntity := res.Form
	if len(res.Forms) > 0 {
		// search based on name
		var found = false
		for _, entity := range res.Forms {
			if entity.Name == searchTerm {
				found = true
				foundEntity = entity
				break
			}
		}
		if !found {
			return diag.Errorf("couldn't find form with name '%s'", searchTerm)
		}
	}

	data.SetId(foundEntity.Id)
	return buildResourceDataFromForm(data, foundEntity)
}
