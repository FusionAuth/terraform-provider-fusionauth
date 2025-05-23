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
			"fusionauth_api_key":                  resourceAPIKey(),
			"fusionauth_application":              newApplication(),
			"fusionauth_application_role":         newApplicationRole(),
			"fusionauth_application_oauth_scope":  newApplicationOAuthScope(),
			"fusionauth_consent":                  newConsent(),
			"fusionauth_email":                    newEmail(),
			"fusionauth_entity":                   resourceEntity(),
			"fusionauth_entity_grant":             resourceEntityGrant(),
			"fusionauth_entity_type":              resourceEntityType(),
			"fusionauth_entity_type_permission":   resourceEntityTypePermission(),
			"fusionauth_generic_connector":        newGenericConnector(),
			"fusionauth_generic_messenger":        newGenericMessenger(),
			"fusionauth_form":                     resourceForm(),
			"fusionauth_form_field":               resourceFormField(),
			"fusionauth_group":                    newGroup(),
			"fusionauth_idp_apple":                resourceIDPApple(),
			"fusionauth_idp_external_jwt":         resourceIDPExternalJWT(),
			"fusionauth_idp_facebook":             resourceIDPFacebook(),
			"fusionauth_idp_google":               newIDPGoogle(),
			"fusionauth_idp_linkedin":             resourceIDPLinkedIn(),
			"fusionauth_idp_open_id_connect":      newIDPOpenIDConnect(),
			"fusionauth_idp_saml_v2":              resourceIDPSAMLv2(),
			"fusionauth_idp_saml_v2_idp_initated": resourceIDPSAMLv2IdPInitiated(),
			"fusionauth_idp_sony_psn":             resourceIDPSonyPSN(),
			"fusionauth_idp_steam":                resourceIDPSteam(),
			"fusionauth_idp_twitch":               resourceIDPTwitch(),
			"fusionauth_idp_xbox":                 resourceIDPXbox(),
			"fusionauth_imported_key":             resourceImportedKey(),
			"fusionauth_key":                      newKey(),
			"fusionauth_lambda":                   newLambda(),
			"fusionauth_ldap_connector":           newLDAPConnector(),
			"fusionauth_reactor":                  newReactor(),
			"fusionauth_registration":             newRegistration(),
			"fusionauth_sms_message_template":     newSMSMessageTemplate(),
			"fusionauth_system_configuration":     resourceSystemConfiguration(),
			"fusionauth_theme":                    newTheme(),
			"fusionauth_tenant":                   newTenant(),
			"fusionauth_twilio_messenger":         newTwilioMessenger(),
			"fusionauth_user":                     newUser(),
			"fusionauth_user_action":              resourceUserAction(),
			"fusionauth_user_group_membership":    newUserGroupMembership(),
			"fusionauth_webhook":                  newWebhook(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"fusionauth_application":             dataSourceApplication(),
			"fusionauth_application_oauth_scope": dataSourceApplicationOAuthScope(),
			"fusionauth_application_role":        dataSourceApplicationRole(),
			"fusionauth_consent":                 dataSourceConsent(),
			"fusionauth_email":                   dataSourceEmail(),
			"fusionauth_form":                    dataSourceForm(),
			"fusionauth_form_field":              dataSourceFormField(),
			"fusionauth_generic_connector":       dataSourceGenericConnector(),
			"fusionauth_generic_messenger":       dataSourceGenericMessenger(),
			"fusionauth_idp":                     dataSourceIDP(),
			"fusionauth_lambda":                  dataSourceLambda(),
			"fusionauth_ldap_connector":          dataSourceLDAPConnector(),
			"fusionauth_sms_message_template":    dataSourceSMSMessageTemplate(),
			"fusionauth_tenant":                  dataSourceTenant(),
			"fusionauth_theme":                   dataSourceTheme(),
			"fusionauth_twilio_messenger":        dataSourceTwilioMessenger(),
			"fusionauth_user":                    dataSourceUser(),
			"fusionauth_user_group_membership":   dataSourceUserGroupMembership(),
		},
		ConfigureContextFunc: configureClient,
	}
}
