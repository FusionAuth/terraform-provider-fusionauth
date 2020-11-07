package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceForm() *schema.Resource {
	return &schema.Resource{
		Create: createForm,
		Read:   readForm,
		Update: updateForm,
		Delete: deleteForm,
		Schema: map[string]*schema.Schema{
			"form_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id to use for the new Form. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that can hold any information about the Form that should be persisted.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Form.",
			},
			"steps": {
				Type:        schema.TypeList,
				Description: "An ordered list of objects containing one or more Form Fields. A Form must have at least one step defined.",
				MinItems:    1,
				Required:    true,
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
				}, false),
			},
		},
	}
}

func createForm(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	f := buildForm(data)
	var fid string
	if fi, ok := data.GetOk("form_id"); ok {
		fid = fi.(string)
	}
	resp, faErrs, err := client.FAClient.CreateForm(fid, fusionauth.FormRequest{Form: f})
	if err != nil {
		return fmt.Errorf("createForm err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}
	data.SetId(resp.Form.Id)
	return buildResourceDataFromForm(data, resp.Form)
}

func readForm(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, err := client.FAClient.RetrieveForm(id)
	if err != nil {
		return err
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return err
	}

	return buildResourceDataFromForm(data, resp.Form)
}

func updateForm(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	f := buildForm(data)

	resp, faErrs, err := client.FAClient.UpdateForm(data.Id(), fusionauth.FormRequest{Form: f})
	if err != nil {
		return fmt.Errorf("updateForm err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}
	data.SetId(resp.Form.Id)
	return buildResourceDataFromForm(data, resp.Form)
}

func deleteForm(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteForm(id)
	if err != nil {
		return err
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	return nil
}

func buildForm(data *schema.ResourceData) fusionauth.Form {
	s := data.Get("steps")
	slice, _ := s.([]interface{})
	steps := make([]fusionauth.FormStep, 0, len(slice))
	for _, i := range slice {
		m := i.(map[string]interface{})
		fields := m["fields"].([]interface{})
		var f []string
		for _, x := range fields {
			f = append(f, x.(string))
		}
		steps = append(steps, fusionauth.FormStep{
			Fields: f,
		})
	}

	return fusionauth.Form{
		Data:  data.Get("data").(map[string]interface{}),
		Name:  data.Get("name").(string),
		Steps: steps,
		Type:  fusionauth.FormType(data.Get("type").(string)),
	}
}

func buildResourceDataFromForm(data *schema.ResourceData, f fusionauth.Form) error {
	if err := data.Set("data", f.Data); err != nil {
		return fmt.Errorf("form.data: %s", err.Error())
	}
	if err := data.Set("name", f.Name); err != nil {
		return fmt.Errorf("form.name: %s", err.Error())
	}

	fs := make([]map[string]interface{}, 0, len(f.Steps))
	for _, step := range f.Steps {
		fs = append(fs, map[string]interface{}{
			"fields": step.Fields,
		},
		)
	}
	if err := data.Set("steps", fs); err != nil {
		return fmt.Errorf("form.steps: %s", err.Error())
	}
	if err := data.Set("type", f.Type); err != nil {
		return fmt.Errorf("form.type: %s", err.Error())
	}

	return nil
}
