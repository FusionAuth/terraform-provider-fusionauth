package fusionauth

import (
	"context"
	"errors"
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceLambda() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLambdaRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				AtLeastOneOf: []string{"id", "name"},
				Description:  "The ID of the Lambda.",
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
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"id", "name"},
				Description:  "The name of the Lambda.",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(fusionauth.LambdaType_JWTPopulate),
					string(fusionauth.LambdaType_OpenIDReconcile),
					string(fusionauth.LambdaType_SAMLv2Reconcile),
					string(fusionauth.LambdaType_SAMLv2Populate),
					string(fusionauth.LambdaType_AppleReconcile),
					string(fusionauth.LambdaType_ExternalJWTReconcile),
					string(fusionauth.LambdaType_FacebookReconcile),
					string(fusionauth.LambdaType_GoogleReconcile),
					string(fusionauth.LambdaType_HYPRReconcile),
					string(fusionauth.LambdaType_TwitterReconcile),
					string(fusionauth.LambdaType_LDAPConnectorReconcile),
					string(fusionauth.LambdaType_LinkedInReconcile),
					string(fusionauth.LambdaType_EpicGamesReconcile),
					string(fusionauth.LambdaType_NintendoReconcile),
					string(fusionauth.LambdaType_SonyPSNReconcile),
					string(fusionauth.LambdaType_SteamReconcile),
					string(fusionauth.LambdaType_TwitchReconcile),
					string(fusionauth.LambdaType_XboxReconcile),
					string(fusionauth.LambdaType_SelfServiceRegistrationValidation),
					string(fusionauth.LambdaType_ClientCredentialsJWTPopulate),
				}, false),
				Description: "The Lambda type.",
			},
		},
	}
}

func dataSourceLambdaRead(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	lambdaType := data.Get("type").(string)
	resp, err := client.FAClient.RetrieveLambdasByType(fusionauth.LambdaType(lambdaType))
	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(fmt.Errorf("couldn't find lambda %s of type %s", name, lambdaType))
	}
	if len(filteredLambdas) > 1 {
		return diag.FromErr(errors.New("query returned more than one lambda. Use a more specific search creteria"))
	}

	l := filteredLambdas[0]

	data.SetId(l.Id)
	err = data.Set("body", l.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("debug", l.Debug)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("name", l.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
