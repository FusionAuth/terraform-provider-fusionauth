package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTenantManagerConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: createTenantManagerConfiguration,
		ReadContext:   readTenantManagerConfiguration,
		UpdateContext: updateTenantManagerConfiguration,
		DeleteContext: deleteTenantManagerConfiguration,
		Schema: map[string]*schema.Schema{
			"application_configurations": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The application configurations for the Tenant Manager.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The unique Id of the Application that the Tenant Manager will use.",
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},
			"attribute_form_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The unique Id of the Form to use for collecting additional user attributes during Tenant Manager registration.",
				ValidateFunc: validation.IsUUID,
			},
			"brand_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The brand name for the Tenant Manager.",
			},
			"identity_provider_type_configurations": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The identity provider type configurations allowed in the Tenant Manager. Each entry corresponds to one identity provider type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The identity provider type. Valid values are: OpenIDConnect and SAMLv2.",
							ValidateFunc: validation.StringInSlice([]string{
								"OpenIDConnect",
								"SAMLv2",
							}, false),
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether or not this identity provider type is enabled in the Tenant Manager.",
						},
						"default_attribute_mappings": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "A map of default attribute mappings for this identity provider type.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"linking_strategy": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "LinkByEmail",
							Description: "The linking strategy for this identity provider type. Valid values are: LinkByEmail, LinkByEmailForExistingUser, LinkByUsername, and LinkByUsernameForExistingUser.",
							ValidateFunc: validation.StringInSlice([]string{
								"LinkByEmail", "LinkByEmailForExistingUser",
								"LinkByUsername", "LinkByUsernameForExistingUser",
							}, false),
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func createTenantManagerConfiguration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	data.SetId("tenantmanager_cfg")
	if diags := updateTenantManagerCfg(buildTenantManagerConfigurationRequest(data), client); diags != nil {
		return diags
	}
	return updateTenantManagerIdpTypeConfigurations(data, client)
}

func readTenantManagerConfiguration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	resp, err := client.FAClient.RetrieveTenantManagerConfiguration()
	if err != nil {
		return diag.Errorf("RetrieveTenantManagerConfiguration err: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	return buildResourceFromTenantManagerConfiguration(resp.TenantManagerConfiguration, data)
}

func updateTenantManagerConfiguration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	if diags := updateTenantManagerCfg(buildTenantManagerConfigurationRequest(data), client); diags != nil {
		return diags
	}
	return updateTenantManagerIdpTypeConfigurations(data, client)
}

func deleteTenantManagerConfiguration(_ context.Context, _ *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	if diags := updateTenantManagerCfg(getDefaultTenantManagerConfigurationRequest(), client); diags != nil {
		return diags
	}

	return updateTenantManagerIdpTypeCfgItems(nil, client)
}

func updateTenantManagerCfg(req fusionauth.TenantManagerConfigurationRequest, client Client) diag.Diagnostics {
	resp, faErrs, err := client.FAClient.UpdateTenantManagerConfiguration(req)
	if err != nil {
		return diag.Errorf("UpdateTenantManagerConfiguration err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func buildTenantManagerConfigurationRequest(data *schema.ResourceData) fusionauth.TenantManagerConfigurationRequest {
	tmc := getDefaultTenantManagerConfigurationRequest()

	if appConfs, isSet := data.GetOk("application_configurations"); isSet {
		list := appConfs.([]interface{})
		configs := make([]fusionauth.TenantManagerApplicationConfiguration, 0, len(list))
		for _, item := range list {
			m := item.(map[string]interface{})
			configs = append(configs, fusionauth.TenantManagerApplicationConfiguration{
				ApplicationId: m["application_id"].(string),
			})
		}
		tmc.TenantManagerConfiguration.ApplicationConfigurations = configs
	}

	if val, isSet := getValueAndIsSet[string](data, "attribute_form_id"); isSet {
		tmc.TenantManagerConfiguration.AttributeFormId = val
	}

	if val, isSet := getValueAndIsSet[string](data, "brand_name"); isSet {
		tmc.TenantManagerConfiguration.BrandName = val
	}

	return tmc
}

func updateTenantManagerIdpTypeConfigurations(data *schema.ResourceData, client Client) diag.Diagnostics {
	desiredItems := getTenantManagerIdentityProviderTypeConfigurationItems(data)
	return updateTenantManagerIdpTypeCfgItems(desiredItems, client)
}

func updateTenantManagerIdpTypeCfgItems(desiredItems []map[string]interface{}, client Client) diag.Diagnostics {
	desiredTypes := getTenantManagerIdentityProviderTypes(desiredItems)

	resp, err := client.FAClient.RetrieveTenantManagerConfiguration()
	if err != nil {
		return diag.Errorf("RetrieveTenantManagerConfiguration err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	for _, idpType := range getTenantManagerIdentityProviderTypesToDelete(resp.TenantManagerConfiguration.IdentityProviderTypeConfigurations, desiredTypes) {
		deleteResp, faErrs, deleteErr := client.FAClient.DeleteTenantManagerIdentityProviderTypeConfiguration(idpType)
		if deleteErr != nil {
			return diag.Errorf("DeleteTenantManagerIdentityProviderTypeConfiguration err: %v", deleteErr)
		}
		if err := checkResponse(deleteResp.StatusCode, faErrs); err != nil {
			return diag.FromErr(err)
		}
	}

	for _, m := range desiredItems {
		idpType := m["type"].(string)
		req := buildTenantManagerIdentityProviderTypeConfigurationRequest(m)
		updateResp, faErrs, updateErr := client.FAClient.UpdateTenantManagerIdentityProviderTypeConfiguration(
			fusionauth.IdentityProviderType(idpType), req,
		)
		if updateErr != nil {
			return diag.Errorf("UpdateTenantManagerIdentityProviderTypeConfiguration err: %v", updateErr)
		}
		if updateResp.StatusCode == http.StatusNotFound {
			createResp, createErrs, createErr := client.FAClient.CreateTenantManagerIdentityProviderTypeConfiguration(
				fusionauth.IdentityProviderType(idpType), req,
			)
			if createErr != nil {
				return diag.Errorf("CreateTenantManagerIdentityProviderTypeConfiguration err: %v", createErr)
			}
			if err := checkResponse(createResp.StatusCode, createErrs); err != nil {
				return diag.FromErr(err)
			}
			continue
		}
		if err := checkResponse(updateResp.StatusCode, faErrs); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func getTenantManagerIdentityProviderTypeConfigurationItems(data *schema.ResourceData) []map[string]interface{} {
	raw, ok := data.GetOk("identity_provider_type_configurations")
	if !ok {
		return nil
	}

	set, ok := raw.(*schema.Set)
	if !ok {
		return nil
	}

	items := set.List()
	configs := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		config := item.(map[string]interface{})
		if getTenantManagerIdentityProviderTypeString(config, "type") == "" {
			continue
		}
		configs = append(configs, config)
	}

	return configs
}

func getTenantManagerIdentityProviderTypes(items []map[string]interface{}) map[string]struct{} {
	types := make(map[string]struct{}, len(items))
	for _, item := range items {
		idpType, ok := item["type"].(string)
		if ok && idpType != "" {
			types[idpType] = struct{}{}
		}
	}

	return types
}

func getTenantManagerIdentityProviderTypesToDelete(current map[string]fusionauth.TenantManagerIdentityProviderTypeConfiguration, desiredTypes map[string]struct{}) []fusionauth.IdentityProviderType {
	types := make([]fusionauth.IdentityProviderType, 0, len(current))
	for idpType, cfg := range current {
		resolvedType := getTenantManagerIdentityProviderType(cfg, idpType)
		if resolvedType == "" {
			continue
		}
		if _, ok := desiredTypes[string(resolvedType)]; ok {
			continue
		}
		types = append(types, resolvedType)
	}

	return types
}

func getTenantManagerIdentityProviderType(cfg fusionauth.TenantManagerIdentityProviderTypeConfiguration, fallback string) fusionauth.IdentityProviderType {
	if cfg.Type != "" {
		return cfg.Type
	}

	return fusionauth.IdentityProviderType(fallback)
}

func buildTenantManagerIdentityProviderTypeConfigurationRequest(m map[string]interface{}) fusionauth.TenantManagerIdentityProviderTypeConfigurationRequest {
	idpType := getTenantManagerIdentityProviderTypeString(m, "type")
	cfg := fusionauth.TenantManagerIdentityProviderTypeConfiguration{
		Enableable: fusionauth.Enableable{
			Enabled: m["enabled"].(bool),
		},
		LinkingStrategy: fusionauth.IdentityProviderLinkingStrategy(getTenantManagerIdentityProviderTypeString(m, "linking_strategy")),
		Type:            fusionauth.IdentityProviderType(idpType),
	}
	if mappings, ok := m["default_attribute_mappings"].(map[string]interface{}); ok && len(mappings) > 0 {
		cfg.DefaultAttributeMappings = make(map[string]string, len(mappings))
		for k, v := range mappings {
			cfg.DefaultAttributeMappings[k] = v.(string)
		}
	}

	return fusionauth.TenantManagerIdentityProviderTypeConfigurationRequest{
		TypeConfiguration: cfg,
	}
}

func getTenantManagerIdentityProviderTypeString(m map[string]interface{}, key string) string {
	value, ok := m[key].(string)
	if ok && value != "" {
		return value
	}

	field := resourceTenantManagerConfiguration().Schema["identity_provider_type_configurations"].Elem.(*schema.Resource).Schema[key]
	if field == nil || field.Default == nil {
		return ""
	}

	defaultValue, ok := field.Default.(string)
	if !ok {
		return ""
	}

	return defaultValue
}

func buildResourceFromTenantManagerConfiguration(tmc fusionauth.TenantManagerConfiguration, data *schema.ResourceData) diag.Diagnostics {
	appConfigs := make([]map[string]interface{}, 0, len(tmc.ApplicationConfigurations))
	for _, ac := range tmc.ApplicationConfigurations {
		appConfigs = append(appConfigs, map[string]interface{}{
			"application_id": ac.ApplicationId,
		})
	}
	if err := data.Set("application_configurations", appConfigs); err != nil {
		return diag.Errorf("tenant_manager_configuration.application_configurations: %s", err.Error())
	}

	if err := data.Set("attribute_form_id", tmc.AttributeFormId); err != nil {
		return diag.Errorf("tenant_manager_configuration.attribute_form_id: %s", err.Error())
	}

	if err := data.Set("brand_name", tmc.BrandName); err != nil {
		return diag.Errorf("tenant_manager_configuration.brand_name: %s", err.Error())
	}

	idpConfigs := make([]map[string]interface{}, 0, len(tmc.IdentityProviderTypeConfigurations))
	for idpType, cfg := range tmc.IdentityProviderTypeConfigurations {
		resolvedType := getTenantManagerIdentityProviderType(cfg, idpType)
		if resolvedType == "" {
			continue
		}
		mappings := make(map[string]interface{}, len(cfg.DefaultAttributeMappings))
		for k, v := range cfg.DefaultAttributeMappings {
			mappings[k] = v
		}
		idpConfigs = append(idpConfigs, map[string]interface{}{
			"type":                       string(resolvedType),
			"enabled":                    cfg.Enabled,
			"linking_strategy":           string(cfg.LinkingStrategy),
			"default_attribute_mappings": mappings,
		})
	}
	if err := data.Set("identity_provider_type_configurations", idpConfigs); err != nil {
		return diag.Errorf("tenant_manager_configuration.identity_provider_type_configurations: %s", err.Error())
	}

	return nil
}

func getDefaultTenantManagerConfigurationRequest() fusionauth.TenantManagerConfigurationRequest {
	return fusionauth.TenantManagerConfigurationRequest{
		TenantManagerConfiguration: fusionauth.TenantManagerConfiguration{},
	}
}
