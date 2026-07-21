package fusionauth

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	appConfigAlreadyExistsErrorFmt = "configuration for application %s on IDP %s already exists"
)

func resourceIDPSAMLv2ApplicationConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: createIDPSAMLv2ApplicationConfiguration,
		ReadContext:   readIDPSAMLv2ApplicationConfiguration,
		UpdateContext: updateIDPSAMLv2ApplicationConfiguration,
		DeleteContext: deleteIDPSAMLv2ApplicationConfiguration,
		Schema: map[string]*schema.Schema{
			"idp_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The ID of the SAML v2 identity provider.",
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
			},
			"application_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The ID of the FusionAuth application to associate with the identity provider.",
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether this identity provider is enabled for the specified application.",
			},
			"create_registration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether FusionAuth will create a UserRegistration for the User automatically when they log in through this identity provider.",
			},
			"button_text": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Application-specific button text override for the identity provider login button.",
			},
			"button_image_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Application-specific button image URL override for the identity provider login button.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: importIDPSAMLv2ApplicationConfiguration,
		},
	}
}

type samlV2IDPApplicationConfigurationIdentityProvider struct {
	ApplicationConfiguration map[string]*SAMLAppConfig `json:"applicationConfiguration"`
}

// samlV2IDPApplicationConfiguration retrieves only the applicationConfiguration
// from a FusionAuth identity provider API response.
type samlV2IDPApplicationConfiguration struct {
	fusionauth.StatusResponse

	IdentityProvider samlV2IDPApplicationConfigurationIdentityProvider `json:"identityProvider"`
}

// samlV2AppConfigPatchBody builds a merge-patch body for a single application configuration
// entry. A nil cfg deletes the entry.
func samlV2AppConfigPatchBody(applicationID string, cfg interface{}) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"identityProvider": map[string]interface{}{
			"applicationConfiguration": map[string]interface{}{applicationID: cfg},
		},
	})
}

// samlV2AppConfigPatchEntry maps resource data to a merge-patch entry. Cleared string overrides
// are sent as explicit nulls — omitting them would leave the previous value in place.
func samlV2AppConfigPatchEntry(data *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"enabled":            data.Get("enabled").(bool),
		"createRegistration": data.Get("create_registration").(bool),
		"buttonText":         stringOrNil(data.Get("button_text").(string)),
		"buttonImageURL":     stringOrNil(data.Get("button_image_url").(string)),
	}
}

// stringOrNil returns nil for an empty string so a JSON merge patch clears the field.
func stringOrNil(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func createIDPSAMLv2ApplicationConfiguration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	idpID := data.Get("idp_id").(string)
	applicationID := data.Get("application_id").(string)

	client := i.(Client)

	b, err := readIdentityProvider(idpID, client)
	if err != nil {
		return diag.FromErr(err)
	}

	var idpBody samlV2IDPApplicationConfiguration
	err = json.Unmarshal(b, &idpBody)
	if err != nil {
		return diag.FromErr(err)
	}

	// Refuse to adopt an existing configuration that hasn't been explicitly imported.
	if _, ok := idpBody.IdentityProvider.ApplicationConfiguration[applicationID]; ok {
		return diag.Errorf(appConfigAlreadyExistsErrorFmt, applicationID, idpID)
	}

	b, err = samlV2AppConfigPatchBody(applicationID, samlV2AppConfigPatchEntry(data))
	if err != nil {
		return diag.FromErr(err)
	}

	err = patchIdentityProvider(b, idpID, client)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceID := fmt.Sprintf("%s:%s", idpID, applicationID)
	data.SetId(resourceID)

	return nil
}

func readIDPSAMLv2ApplicationConfiguration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	idpID := data.Get("idp_id").(string)
	applicationID := data.Get("application_id").(string)

	client := i.(Client)

	b, err := readIdentityProvider(idpID, client)
	if err != nil {
		if err.Error() == NotFoundError {
			data.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	var idpBody samlV2IDPApplicationConfiguration
	err = json.Unmarshal(b, &idpBody)
	if err != nil {
		return diag.FromErr(err)
	}

	appConfig, ok := idpBody.IdentityProvider.ApplicationConfiguration[applicationID]
	if !ok {
		data.SetId("")
		return nil
	}

	if err := data.Set("enabled", appConfig.Enabled); err != nil {
		return diag.Errorf("idpSAMLv2ApplicationConfiguration.enabled: %s", err.Error())
	}
	if err := data.Set("create_registration", appConfig.CreateRegistration); err != nil {
		return diag.Errorf("idpSAMLv2ApplicationConfiguration.create_registration: %s", err.Error())
	}
	if err := data.Set("button_text", appConfig.ButtonText); err != nil {
		return diag.Errorf("idpSAMLv2ApplicationConfiguration.button_text: %s", err.Error())
	}
	if err := data.Set("button_image_url", appConfig.ButtonImageURL); err != nil {
		return diag.Errorf("idpSAMLv2ApplicationConfiguration.button_image_url: %s", err.Error())
	}

	return nil
}

func updateIDPSAMLv2ApplicationConfiguration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	idpID := data.Get("idp_id").(string)
	applicationID := data.Get("application_id").(string)

	client := i.(Client)

	b, err := samlV2AppConfigPatchBody(applicationID, samlV2AppConfigPatchEntry(data))
	if err != nil {
		return diag.FromErr(err)
	}

	err = patchIdentityProvider(b, idpID, client)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteIDPSAMLv2ApplicationConfiguration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	idpID := data.Get("idp_id").(string)
	applicationID := data.Get("application_id").(string)

	client := i.(Client)

	// A null entry in the merge patch deletes the application configuration
	b, err := samlV2AppConfigPatchBody(applicationID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	err = patchIdentityProvider(b, idpID, client)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func importIDPSAMLv2ApplicationConfiguration(_ context.Context, data *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	id := data.Id()

	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import ID format, expected 'idp_id:application_id', got: %s", id)
	}

	idpID := parts[0]
	applicationID := parts[1]

	if err := data.Set("idp_id", idpID); err != nil {
		return nil, err
	}
	if err := data.Set("application_id", applicationID); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{data}, nil
}
