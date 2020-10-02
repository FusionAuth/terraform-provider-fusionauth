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

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
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

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	l := resp.Lambda
	if err := data.Set("body", l.Body); err != nil {
		return fmt.Errorf("lambda.body: %s", err.Error())
	}
	if err := data.Set("debug", l.Debug); err != nil {
		return fmt.Errorf("lambda.debug: %s", err.Error())
	}
	if err := data.Set("enabled", l.Enabled); err != nil {
		return fmt.Errorf("lambda.enabled: %s", err.Error())
	}
	if err := data.Set("name", l.Name); err != nil {
		return fmt.Errorf("lambda.name: %s", err.Error())
	}
	if err := data.Set("type", l.Type); err != nil {
		return fmt.Errorf("lambda.type: %s", err.Error())
	}

	return nil
}

func updateLambda(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	l := buildLambda(data)

	resp, faErrs, err := client.FAClient.UpdateLambda(data.Id(), fusionauth.LambdaRequest{
		Lambda: l,
	})
	if err != nil {
		return err
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	return nil
}

func deleteLambda(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteLambda(id)
	if err != nil {
		return err
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	return nil
}
