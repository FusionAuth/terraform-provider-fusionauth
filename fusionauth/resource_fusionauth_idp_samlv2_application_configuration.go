package fusionauth

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Global mutex map for IDP-level locking to prevent race conditions
// when multiple association resources modify the same IDP concurrently
var (
	idpLocksMutex = sync.RWMutex{}
	idpLocks      = make(map[string]*sync.Mutex)
)

// getIDPMutex returns a mutex for the given IDP ID, creating one if it doesn't exist
func getIDPMutex(idpId string) *sync.Mutex {
	idpLocksMutex.RLock()
	if mutex, exists := idpLocks[idpId]; exists {
		idpLocksMutex.RUnlock()
		return mutex
	}
	idpLocksMutex.RUnlock()

	idpLocksMutex.Lock()
	defer idpLocksMutex.Unlock()

	// Double-check pattern in case another goroutine created it
	if mutex, exists := idpLocks[idpId]; exists {
		return mutex
	}

	// Create new mutex for this IDP
	mutex := &sync.Mutex{}
	idpLocks[idpId] = mutex
	return mutex
}

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
				Default:     true,
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

func createIDPSAMLv2ApplicationConfiguration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	idpId := data.Get("idp_id").(string)
	applicationId := data.Get("application_id").(string)

	// Acquire IDP-specific mutex to prevent race conditions
	mutex := getIDPMutex(idpId)
	mutex.Lock()
	defer func() {
		mutex.Unlock()
	}()

	client := i.(Client)

	// Read current IDP configuration
	b, err := readIdentityProvider(idpId, client)
	if err != nil {
		return diag.FromErr(err)
	}

	var idpBody SAMLIdentityProviderBody
	err = json.Unmarshal(b, &idpBody)
	if err != nil {
		return diag.FromErr(err)
	}

	// Do not check if association already exists
	// When we pull from the API we will get whatever is already there,
	// Which was set by us
	//if idpBody.IdentityProvider.ApplicationConfiguration == nil {
	//	idpBody.IdentityProvider.ApplicationConfiguration = make(map[string]interface{})
	//}

	//if _, exists := idpBody.IdentityProvider.ApplicationConfiguration[applicationId]; exists {
	//	return diag.Errorf("association between IDP %s and application %s already exists", idpId, applicationId)
	//}

	// Create PATCH payload with only the application configuration
	appConfig := SAMLAppConfig{
		Enabled:            data.Get("enabled").(bool),
		CreateRegistration: data.Get("create_registration").(bool),
		ButtonText:         data.Get("button_text").(string),
		ButtonImageURL:     data.Get("button_image_url").(string),
	}

	// Use PATCH to update only the specific application configuration
	patchData := map[string]interface{}{
		"identityProvider": map[string]interface{}{
			"applicationConfiguration": map[string]interface{}{
				applicationId: appConfig,
			},
		},
	}

	patchB, err := json.Marshal(patchData)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = patchIdentityProvider(patchB, idpId, client)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set resource ID as combination of IDP ID and Application ID
	resourceId := fmt.Sprintf("%s:%s", idpId, applicationId)
	data.SetId(resourceId)

	return nil
}

func readIDPSAMLv2ApplicationConfiguration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	idpId := data.Get("idp_id").(string)
	applicationId := data.Get("application_id").(string)

	// Acquire IDP-specific mutex to ensure consistent reads during concurrent modifications
	mutex := getIDPMutex(idpId)
	mutex.Lock()
	defer func() {
		mutex.Unlock()
	}()

	client := i.(Client)

	// Read current IDP configuration
	b, err := readIdentityProvider(idpId, client)
	if err != nil {
		return diag.FromErr(err)
	}

	var idpBody SAMLIdentityProviderBody
	err = json.Unmarshal(b, &idpBody)
	if err != nil {
		return diag.FromErr(err)
	}

	// Check if association exists
	appConfigInterface, exists := idpBody.IdentityProvider.ApplicationConfiguration[applicationId]
	if !exists {
		data.SetId("")
		return nil
	}

	// Convert to our struct format
	appConfigJson, err := json.Marshal(appConfigInterface)
	if err != nil {
		return diag.FromErr(err)
	}

	var appConfig SAMLAppConfig
	err = json.Unmarshal(appConfigJson, &appConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set values in Terraform state
	if err := data.Set("enabled", appConfig.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("create_registration", appConfig.CreateRegistration); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("button_text", appConfig.ButtonText); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("button_image_url", appConfig.ButtonImageURL); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func updateIDPSAMLv2ApplicationConfiguration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	idpId := data.Get("idp_id").(string)
	applicationId := data.Get("application_id").(string)

	// Acquire IDP-specific mutex to prevent race conditions
	mutex := getIDPMutex(idpId)
	mutex.Lock()
	defer func() {
		mutex.Unlock()
	}()

	client := i.(Client)

	// Read current IDP configuration
	b, err := readIdentityProvider(idpId, client)
	if err != nil {
		return diag.FromErr(err)
	}

	var idpBody SAMLIdentityProviderBody
	err = json.Unmarshal(b, &idpBody)
	if err != nil {
		return diag.FromErr(err)
	}

	// Check if association exists
	if idpBody.IdentityProvider.ApplicationConfiguration == nil {
		idpBody.IdentityProvider.ApplicationConfiguration = make(map[string]interface{})
	}

	_, exists := idpBody.IdentityProvider.ApplicationConfiguration[applicationId]
	if !exists {
		return diag.Errorf("association between IDP %s and application %s does not exist", idpId, applicationId)
	}

	// Create PATCH payload with only the application configuration
	appConfig := SAMLAppConfig{
		Enabled:            data.Get("enabled").(bool),
		CreateRegistration: data.Get("create_registration").(bool),
		ButtonText:         data.Get("button_text").(string),
		ButtonImageURL:     data.Get("button_image_url").(string),
	}

	// Use PATCH to update only the specific application configuration
	patchData := map[string]interface{}{
		"identityProvider": map[string]interface{}{
			"applicationConfiguration": map[string]interface{}{
				applicationId: appConfig,
			},
		},
	}

	patchB, err := json.Marshal(patchData)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = patchIdentityProvider(patchB, idpId, client)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteIDPSAMLv2ApplicationConfiguration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	idpId := data.Get("idp_id").(string)
	applicationId := data.Get("application_id").(string)

	// Acquire IDP-specific mutex to prevent race conditions
	mutex := getIDPMutex(idpId)
	mutex.Lock()
	defer func() {
		mutex.Unlock()
	}()

	client := i.(Client)

	// Read current IDP configuration
	b, err := readIdentityProvider(idpId, client)
	if err != nil {
		return diag.FromErr(err)
	}

	var idpBody SAMLIdentityProviderBody
	err = json.Unmarshal(b, &idpBody)
	if err != nil {
		return diag.FromErr(err)
	}

	// Remove association using JSON Patch (RFC 6902)
	if idpBody.IdentityProvider.ApplicationConfiguration != nil {
		if _, exists := idpBody.IdentityProvider.ApplicationConfiguration[applicationId]; exists {
			p := patch{
				{
					Op:   removeOp,
					Path: fmt.Sprintf("/identityProvider/applicationConfiguration/%s", applicationId),
				},
			}

			_, err = patchIdentityProviderJsonPatch(p, idpId, client)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func importIDPSAMLv2ApplicationConfiguration(_ context.Context, data *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	id := data.Id()

	// Parse composite ID: "idp_id:application_id"
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import ID format, expected 'idp_id:application_id', got: %s", id)
	}

	idpId := parts[0]
	applicationId := parts[1]

	// Set the individual fields
	if err := data.Set("idp_id", idpId); err != nil {
		return nil, err
	}
	if err := data.Set("application_id", applicationId); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{data}, nil
}
