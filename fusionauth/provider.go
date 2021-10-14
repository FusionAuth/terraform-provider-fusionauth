package fusionauth

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider configures and returns a fusionauth terraform provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FA_DOMAIN", nil),
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FA_API_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"fusionauth_lambda":                   newLambda(),
			"fusionauth_application":              newApplication(),
			"fusionauth_user":                     newUser(),
			"fusionauth_theme":                    newTheme(),
			"fusionauth_tenant":                   newTenant(),
			"fusionauth_email":                    newEmail(),
			"fusionauth_key":                      newKey(),
			"fusionauth_webhook":                  newWebhook(),
			"fusionauth_group":                    newGroup(),
			"fusionauth_application_role":         newApplicationRole(),
			"fusionauth_idp_open_id_connect":      newIDPOpenIDConnect(),
			"fusionauth_idp_google":               newIDPGoogle(),
			"fusionauth_registration":             newRegistration(),
			"fusionauth_system_configuration":     resourceSystemConfiguration(),
			"fusionauth_form_field":               resourceFormField(),
			"fusionauth_form":                     resourceForm(),
			"fusionauth_idp_apple":                resourceIDPApple(),
			"fusionauth_imported_key":             resourceImportedKey(),
			"fusionauth_user_action":              resourceUserAction(),
			"fusionauth_idp_external_jwt":         resourceIDPExternalJWT(),
			"fusionauth_idp_saml_v2":              resourceIDPSAMLv2(),
			"fusionauth_api_key":                  resourceAPIKey(),
			"fusionauth_idp_saml_v2_idp_initated": resourceIDPSAMLv2IdPInitiated(),
			"fusionauth_idp_xbox":                 resourceIDPXbox(),
			"fusionauth_idp_sony_psn":             resourceIDPSonyPSN(),
			"fusionauth_idp_steam":                resourceIDPSteam(),
			"fusionauth_idp_twitch":               resourceIDPTwitch(),
			"fusionauth_idp_facebook":             resourceIDPFacebook(),
			"fusionauth_entity_type":              newEntityType(),
			"fusionauth_entity":                   newEntity(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"fusionauth_lambda":           dataSourceLambda(),
			"fusionauth_application":      dataSourceApplication(),
			"fusionauth_tenant":           dataSourceTenant(),
			"fusionauth_application_role": dataSourceApplicationRole(),
			"fusionauth_idp":              dataSourceIDP(),
		},
		ConfigureContextFunc: configureClient,
	}
}
