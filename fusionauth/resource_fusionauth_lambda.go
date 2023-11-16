package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newLambda() *schema.Resource {
	return &schema.Resource{
		CreateContext: createLambda,
		ReadContext:   readLambda,
		UpdateContext: updateLambda,
		DeleteContext: deleteLambda,
		Schema: map[string]*schema.Schema{
			"lambda_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The Id to use for the new lambda. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
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
				Deprecated:  "Not currently used and may be removed in a future version.",
			},
			"engine_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(fusionauth.LambdaEngineType_GraalJS),
				Description:  "The JavaScript execution engine for the lambda.",
				ValidateFunc: validation.StringInSlice([]string{string(fusionauth.LambdaEngineType_GraalJS), string(fusionauth.LambdaEngineType_Nashorn)}, false),
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
				Description: "The lambda type.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func buildLambda(data *schema.ResourceData) fusionauth.Lambda {
	l := fusionauth.Lambda{
		Id:         data.Get("lambda_id").(string),
		Body:       data.Get("body").(string),
		Debug:      data.Get("debug").(bool),
		Name:       data.Get("name").(string),
		Type:       fusionauth.LambdaType(data.Get("type").(string)),
		EngineType: fusionauth.LambdaEngineType(data.Get("engine_type").(string)),
	}
	return l
}

func createLambda(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	l := buildLambda(data)
	resp, faErrs, err := client.FAClient.CreateLambda(l.Id, fusionauth.LambdaRequest{
		Lambda: l,
	})
	if err != nil {
		return diag.Errorf("CreateLambda err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(resp.Lambda.Id)
	return nil
}

func readLambda(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveLambda(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	l := resp.Lambda
	if err := data.Set("body", l.Body); err != nil {
		return diag.Errorf("lambda.body: %s", err.Error())
	}
	if err := data.Set("debug", l.Debug); err != nil {
		return diag.Errorf("lambda.debug: %s", err.Error())
	}
	if err := data.Set("engine_type", l.EngineType); err != nil {
		return diag.Errorf("lambda.engine_type: %s", err.Error())
	}
	if err := data.Set("name", l.Name); err != nil {
		return diag.Errorf("lambda.name: %s", err.Error())
	}
	if err := data.Set("type", l.Type); err != nil {
		return diag.Errorf("lambda.type: %s", err.Error())
	}

	return nil
}

func updateLambda(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	l := buildLambda(data)

	resp, faErrs, err := client.FAClient.UpdateLambda(data.Id(), fusionauth.LambdaRequest{
		Lambda: l,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteLambda(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteLambda(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
