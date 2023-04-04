package fusionauth

import (
	"context"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func newReactor() *schema.Resource {
	return &schema.Resource{
		CreateContext: createReactor,
		ReadContext:   readReactor,
		UpdateContext: updateReactor,
		DeleteContext: deleteReactor,
		Schema: map[string]*schema.Schema{
			"license_id": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The license Id to activate",
			},
			"license": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The Base64 encoded license value. This value is necessary in an air gapped configuration where outbound network access is not available.",
			},
		},
	}
}

func buildReactor(data *schema.ResourceData) fusionauth.ReactorRequest {
	reactor := fusionauth.ReactorRequest{
		LicenseId: data.Get("license_id").(string),
		License:   data.Get("license").(string),
	}
	return reactor
}

func createReactor(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	reactor := buildReactor(data)

	resp, faErrs, err := client.FAClient.ActivateReactor(reactor)
	if err != nil {
		return diag.Errorf("CreateReactor err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	// Every terraform resource needs an id. Since none of the fusion auth reactor apis return
	// the id of the reactor, we will need to make one up. There can be only one reactor
	// so this should be safe.
	data.SetId("fusion-auth-reactor-id")

	return nil
}

func updateReactor(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	reactor := buildReactor(data)

	resp, faErrs, err := client.FAClient.ActivateReactor(fusionauth.ReactorRequest{
		LicenseId: reactor.LicenseId,
		License:   reactor.License,
	})
	if err != nil {
		return diag.Errorf("UpdateReactor err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		data.Partial(true)
		return diag.FromErr(err)
	}
	return nil
}

func deleteReactor(_ context.Context, _ *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	resp, err := client.FAClient.DeactivateReactor()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func readReactor(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	resp, err := client.FAClient.RetrieveReactorStatus()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	// The reactor read api does not return the license_id or license values. So we
	// cannot update the state using this response like other resources. However, the
	// reactor response does include a boolean attribute telling us if it is licensed.
	if resp.Status.Licensed {
		// If the reactor is licensed, not setting the attributes here should leave
		// the state unchanged. Adding a warning here so that users know this doesn't
		// work like a normal resource.
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Unable to set reactor license attributes from FusionAuth response",
			},
		}
	}

	// If the reactor is not licensed, we know these attributes must be blank.
	if err := data.Set("license_id", nil); err != nil {
		return diag.Errorf("reactor.license_id: %s", err.Error())
	}
	if err := data.Set("license", nil); err != nil {
		return diag.Errorf("reactor.license: %s", err.Error())
	}

	return nil
}
