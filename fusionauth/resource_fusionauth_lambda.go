package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func newLambda() *schema.Resource {
	return &schema.Resource{
		Create: createLambda,
		Read:   readLambda,
		Update: updateLambda,
		Delete: deleteLambda,
		Schema: map[string]*schema.Schema{
			"body": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The lambda function body, a JavaScript function.",
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not debug event logging is enabled for this Lambda.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not this Lambda is enabled.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the lambda.",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"JWTPopulate",
					"OpenIDReconcile",
					"SAMLv2Reconcile",
					"SAMLv2Populate",
					"AppleReconcile",
					"ExternalJWTReconcile",
					"FacebookReconcile",
					"GoogleReconcile",
					"HYPRReconcile",
					"TwitterReconcile",
					"LDAPConnectorReconcile",
				}, false),
				Description: "The lambda type.",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func buildLambda(data *schema.ResourceData) fusionauth.Lambda {
	l := fusionauth.Lambda{
		Body:       data.Get("body").(string),
		Debug:      data.Get("debug").(bool),
		Name:       data.Get("name").(string),
		Type:       fusionauth.LambdaType(data.Get("type").(string)),
		Enableable: buildEnableable("enabled", data),
	}
	return l
}

func createLambda(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	l := buildLambda(data)
	resp, faErrs, err := client.FAClient.CreateLambda("", fusionauth.LambdaRequest{
		Lambda: l,
	})
	if err != nil {
		return fmt.Errorf("CreateLambda err: %v", err)
	}

	if faErrs != nil {
		return fmt.Errorf("CreateLambda errors: %v", faErrs)
	}
	data.SetId(resp.Lambda.Id)
	return nil
}

func readLambda(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveLambda(id)
	if err != nil {
		return err
	}

	if faErrs != nil {
		return fmt.Errorf("UpdateLambda errors: %v", faErrs)
	}

	l := resp.Lambda
	_ = data.Set("body", l.Body)
	_ = data.Set("debug", l.Debug)
	_ = data.Set("enabled", l.Enabled)
	_ = data.Set("name", l.Name)
	_ = data.Set("type", l.Type)

	return nil
}

func updateLambda(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	l := buildLambda(data)

	_, faErrs, err := client.FAClient.UpdateLambda(data.Id(), fusionauth.LambdaRequest{
		Lambda: l,
	})
	if err != nil {
		return err
	}

	if faErrs != nil {
		return fmt.Errorf("UpdateLambda errors: %v", faErrs)
	}

	return nil
}

func deleteLambda(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	_, faErrs, err := client.FAClient.DeleteLambda(id)
	if err != nil {
		return err
	}

	if faErrs != nil {
		return fmt.Errorf("DeleteLambda errors: %v", faErrs)
	}

	return nil
}
