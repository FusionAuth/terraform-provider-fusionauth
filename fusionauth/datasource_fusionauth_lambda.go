package fusionauth

import (
	"errors"
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceLambda() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLambdaRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The ID of the Lambda.",
			},
			"body": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lambda function body, a JavaScript function.",
			},
			"debug": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not debug event logging is enabled for this Lambda.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the Lambda.",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Description: "The Lambda type.",
			},
		},
	}
}

func dataSourceLambdaRead(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)

	lambdaType := data.Get("type").(string)
	resp, err := client.FAClient.RetrieveLambdasByType(fusionauth.LambdaType(lambdaType))
	if err != nil {
		return err
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return err
	}

	name := data.Get("name").(string)
	preFilteredLambdas := resp.Lambdas
	var filteredLambdas []fusionauth.Lambda

	if id, ok := data.GetOk("id"); ok {
		for _, l := range preFilteredLambdas {
			if l.Id == id.(string) {
				filteredLambdas = append(filteredLambdas, l)
			}
		}
		preFilteredLambdas = filteredLambdas
	}

	if name, ok := data.GetOk("name"); ok {
		filteredLambdas = []fusionauth.Lambda{}
		for _, l := range preFilteredLambdas {
			if l.Name == name.(string) {
				filteredLambdas = append(filteredLambdas, l)
			}
		}
	}

	if len(filteredLambdas) < 1 {
		return fmt.Errorf("couldn't find lambda %s of type %s", name, lambdaType)
	}
	if len(filteredLambdas) > 1 {
		return errors.New("query returned more than one lambda. Use a more specific search creteria")
	}

	l := filteredLambdas[0]

	data.SetId(l.Id)
	data.Set("body", l.Body)
	data.Set("debug", l.Debug)
	data.Set("name", l.Name)
	return nil
}
