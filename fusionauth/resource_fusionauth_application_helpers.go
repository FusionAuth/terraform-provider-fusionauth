package fusionauth

import (
	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func buildApplication(data *schema.ResourceData) fusionauth.Application {
	a := fusionauth.Application{
		TenantId: data.Get("tenant_id").(string),
		AuthenticationTokenConfiguration: fusionauth.AuthenticationTokenConfiguration{
			Enableable: buildEnableable("authentication_token_configuration_enabled", data),
		},
		CleanSpeakConfiguration: fusionauth.CleanSpeakConfiguration{
			ApplicationIds: handleStringSlice("clean_speak_configuration.0.application_ids", data),
			UsernameModeration: fusionauth.UsernameModeration{
				ApplicationId: data.Get("clean_speak_configuration.0.username_moderation.0.application_id").(string),
				Enableable:    buildEnableable("clean_speak_configuration.0.username_moderation.0.enabled", data),
			},
		},
		Data: data.Get("data").(map[string]interface{}),
		JwtConfiguration: fusionauth.JWTConfiguration{
			Enableable:                      buildEnableable("jwt_configuration.0.enabled", data),
			AccessTokenKeyId:                data.Get("jwt_configuration.0.access_token_id").(string),
			IdTokenKeyId:                    data.Get("jwt_configuration.0.id_token_key_id").(string),
			RefreshTokenTimeToLiveInMinutes: data.Get("jwt_configuration.0.refresh_token_ttl_minutes").(int),
			TimeToLiveInSeconds:             data.Get("jwt_configuration.0.ttl_seconds").(int),
		},
		LambdaConfiguration: fusionauth.LambdaConfiguration{
			AccessTokenPopulateId: data.Get("lambda_configuration.0.access_token_populate_id").(string),
			IdTokenPopulateId:     data.Get("lambda_configuration.0.id_token_populate_id").(string),
			Samlv2PopulateId:      data.Get("lambda_configuration.0.samlv2_populate_id").(string),
		},
		LoginConfiguration: fusionauth.LoginConfiguration{
			AllowTokenRefresh:     data.Get("login_configuration.0.allow_token_refresh").(bool),
			GenerateRefreshTokens: data.Get("login_configuration.0.generate_refresh_tokens").(bool),
			RequireAuthentication: data.Get("login_configuration.0.require_authentication").(bool),
		},
		Name: data.Get("name").(string),
		OauthConfiguration: fusionauth.OAuth2Configuration{
			AuthorizedOriginURLs:        handleStringSlice("oauth_configuration.0.authorized_origin_urls", data),
			AuthorizedRedirectURLs:      handleStringSlice("oauth_configuration.0.authorized_redirect_urls", data),
			ClientSecret:                data.Get("oauth_configuration.0.client_secret").(string),
			DeviceVerificationURL:       data.Get("oauth_configuration.0.device_verification_url").(string),
			GenerateRefreshTokens:       data.Get("oauth_configuration.0.generate_refresh_tokens").(bool),
			LogoutURL:                   data.Get("oauth_configuration.0.logout_url").(string),
			RequireClientAuthentication: data.Get("oauth_configuration.0.require_client_authentication").(bool),
			LogoutBehavior:              fusionauth.LogoutBehavior(data.Get("oauth_configuration.0.logout_behavior").(string)),
			EnabledGrants:               buildGrants("oauth_configuration.0.enabled_grants", data),
		},
		PasswordlessConfiguration: fusionauth.PasswordlessConfiguration{
			Enableable: buildEnableable("passwordless_configuration_enabled", data),
		},
		RegistrationConfiguration: fusionauth.RegistrationConfiguration{
			Enableable:      buildEnableable("registration_configuration.0.enabled", data),
			BirthDate:       buildRequireable("registration_configuration.0.birth_date", data),
			ConfirmPassword: data.Get("registration_configuration.0.confirm_password").(bool),
			FormId:          data.Get("registration_configuration.0.form_id").(string),
			FirstName:       buildRequireable("registration_configuration.0.first_name", data),
			FullName:        buildRequireable("registration_configuration.0.full_name", data),
			LastName:        buildRequireable("registration_configuration.0.last_name", data),
			MiddleName:      buildRequireable("registration_configuration.0.middle_name", data),
			MobilePhone:     buildRequireable("registration_configuration.0.mobile_phone", data),
			LoginIdType:     fusionauth.LoginIdType(data.Get("registration_configuration.0.login_id_type").(string)),
			Type:            fusionauth.RegistrationType(data.Get("registration_configuration.0.type").(string)),
		},
		RegistrationDeletePolicy: fusionauth.ApplicationRegistrationDeletePolicy{
			Unverified: fusionauth.TimeBasedDeletePolicy{
				Enableable:           buildEnableable("registration_delete_policy.0.unverified_enabled", data),
				NumberOfDaysToRetain: data.Get("registration_delete_policy.0.unverified_number_of_days_to_retain").(int),
			},
		},
		Samlv2Configuration: fusionauth.SAMLv2Configuration{
			Enableable:  buildEnableable("samlv2_configuration.0.enabled", data),
			Audience:    data.Get("samlv2_configuration.0.audience").(string),
			CallbackURL: data.Get("samlv2_configuration.0.callback_url").(string),
			Debug:       data.Get("samlv2_configuration.0.debug").(bool),
			Issuer:      data.Get("samlv2_configuration.0.issuer").(string),
			KeyId:       data.Get("samlv2_configuration.0.key_id").(string),
			LogoutURL:   data.Get("samlv2_configuration.0.logout_url").(string),
			XmlSignatureC14nMethod: fusionauth.CanonicalizationMethod(
				data.Get("samlv2_configuration.0.xml_signature_canonicalization_method").(string),
			),
		},
		VerificationEmailTemplateId: data.Get("verification_email_template_id").(string),
		VerifyRegistration:          data.Get("verify_registration").(bool),
		EmailConfiguration: fusionauth.ApplicationEmailConfiguration{
			EmailVerificationEmailTemplateId: data.Get("email_configuration.0.email_verification_template_id").(string),
			ForgotPasswordEmailTemplateId:    data.Get("email_configuration.0.forgot_password_template_id").(string),
			PasswordlessEmailTemplateId:      data.Get("email_configuration.0.passwordless_email_template_id").(string),
			SetPasswordEmailTemplateId:       data.Get("email_configuration.0.set_password_email_template_id").(string),
		},
	}

	return a
}

func buildGrants(key string, data *schema.ResourceData) []fusionauth.GrantType {
	grants := handleStringSlice(key, data)
	gs := make([]fusionauth.GrantType, 0, len(grants))
	for _, g := range grants {
		gs = append(gs, fusionauth.GrantType(g))
	}
	return gs
}

func buildEnableable(key string, data *schema.ResourceData) fusionauth.Enableable {
	return fusionauth.Enableable{
		Enabled: data.Get(key).(bool),
	}
}

func buildRequireable(key string, data *schema.ResourceData) fusionauth.Requirable {
	return fusionauth.Requirable{
		Enableable: buildEnableable(key+".0.enabled", data),
		Required:   data.Get(key + ".0.required").(bool),
	}
}

func buildResourceDataFromApplication(a fusionauth.Application, data *schema.ResourceData) {
	_ = data.Set("tenant_id", a.TenantId)
	_ = data.Set("authentication_token_configuration_enabled", a.AuthenticationTokenConfiguration.Enabled)
	_ = data.Set("clean_speak_configuration.0.application_ids", a.CleanSpeakConfiguration.ApplicationIds)
	_ = data.Set(
		"clean_speak_configuration.0.username_moderation.0.application_id",
		a.CleanSpeakConfiguration.UsernameModeration.ApplicationId,
	)
	_ = data.Set(
		"clean_speak_configuration.0.username_moderation.0.enabled",
		a.CleanSpeakConfiguration.UsernameModeration.Enabled,
	)
	_ = data.Set("data", a.Data)
	_ = data.Set("jwt_configuration.0.enabled", a.JwtConfiguration.Enabled)
	_ = data.Set("jwt_configuration.0.access_token_id", a.JwtConfiguration.AccessTokenKeyId)
	_ = data.Set("jwt_configuration.0.id_token_key_id", a.JwtConfiguration.IdTokenKeyId)
	_ = data.Set("jwt_configuration.0.refresh_token_ttl_minutes", a.JwtConfiguration.RefreshTokenTimeToLiveInMinutes)
	_ = data.Set("jwt_configuration.0.ttl_seconds", a.JwtConfiguration.TimeToLiveInSeconds)
	_ = data.Set("lambda_configuration.0.access_token_populate_id", a.LambdaConfiguration.AccessTokenPopulateId)
	_ = data.Set("lambda_configuration.0.id_token_populate_id", a.LambdaConfiguration.IdTokenPopulateId)
	_ = data.Set("lambda_configuration.0.samlv2_populate_id", a.LambdaConfiguration.Samlv2PopulateId)
	_ = data.Set("login_configuration.0.allow_token_refresh", a.LoginConfiguration.AllowTokenRefresh)
	_ = data.Set("login_configuration.0.generate_refresh_tokens", a.LoginConfiguration.GenerateRefreshTokens)
	_ = data.Set("login_configuration.0.require_authentication", a.LoginConfiguration.RequireAuthentication)
	_ = data.Set("name", a.Name)
	_ = data.Set("oauth_configuration.0.authorized_origin_urls", a.OauthConfiguration.AuthorizedOriginURLs)
	_ = data.Set("oauth_configuration.0.authorized_redirect_urls", a.OauthConfiguration.AuthorizedRedirectURLs)
	_ = data.Set("oauth_configuration.0.client_secret", a.OauthConfiguration.ClientSecret)
	_ = data.Set("oauth_configuration.0.device_verification_url", a.OauthConfiguration.DeviceVerificationURL)
	_ = data.Set("oauth_configuration.0.generate_refresh_tokens", a.OauthConfiguration.GenerateRefreshTokens)
	_ = data.Set("oauth_configuration.0.logout_url", a.OauthConfiguration.LogoutURL)
	_ = data.Set("oauth_configuration.0.require_client_authentication", a.OauthConfiguration.RequireClientAuthentication)
	_ = data.Set("oauth_configuration.0.logout_behavior", a.OauthConfiguration.LogoutBehavior)
	_ = data.Set("oauth_configuration.0.authorized_origin_urls", a.OauthConfiguration.AuthorizedOriginURLs)
	_ = data.Set("oauth_configuration.0.authorized_origienabled_grants_urls", a.OauthConfiguration.EnabledGrants)
	_ = data.Set("passwordless_configuration_enabled", a.PasswordlessConfiguration.Enabled)
	_ = data.Set("registration_configuration.0.enabled", a.RegistrationConfiguration.Enabled)
	_ = data.Set("registration_configuration.0.birth_date.0.enabled", a.RegistrationConfiguration.BirthDate.Enabled)
	_ = data.Set("registration_configuration.0.birth_date.0.required", a.RegistrationConfiguration.BirthDate.Required)
	_ = data.Set("registration_configuration.0.confirm_password", a.RegistrationConfiguration.ConfirmPassword)
	_ = data.Set("registration_configuration.0.first_name.0.enabled", a.RegistrationConfiguration.FirstName.Enabled)
	_ = data.Set("registration_configuration.0.first_name.0.required", a.RegistrationConfiguration.FirstName.Required)
	_ = data.Set("registration_configuration.0.full_name.0.enabled", a.RegistrationConfiguration.FullName.Enabled)
	_ = data.Set("registration_configuration.0.full_name.0.required", a.RegistrationConfiguration.FullName.Required)
	_ = data.Set("registration_configuration.0.last_name.0.enabled", a.RegistrationConfiguration.LastName.Enabled)
	_ = data.Set("registration_configuration.0.last_name.0.required", a.RegistrationConfiguration.LastName.Required)
	_ = data.Set("registration_configuration.0.middle_name.0.enabled", a.RegistrationConfiguration.MiddleName.Enabled)
	_ = data.Set("registration_configuration.0.middle_name.0.required", a.RegistrationConfiguration.MiddleName.Required)
	_ = data.Set("registration_configuration.0.mobile_phone.0.enabled", a.RegistrationConfiguration.MobilePhone.Enabled)
	_ = data.Set("registration_configuration.0.mobile_phone.0.required", a.RegistrationConfiguration.MobilePhone.Required)
	_ = data.Set("registration_configuration.0.login_id_type", a.RegistrationConfiguration.LoginIdType)
	_ = data.Set("registration_configuration.0.type", a.RegistrationConfiguration.Type)
	_ = data.Set("registration_delete_policy.0.unverified_enabled", a.RegistrationDeletePolicy.Unverified.Enabled)
	_ = data.Set(
		"registration_delete_policy.0.unverified_number_of_days_to_retain",
		a.RegistrationDeletePolicy.Unverified.NumberOfDaysToRetain,
	)
	_ = data.Set("registration_configuration.0.form_id", a.RegistrationConfiguration.FormId)
	_ = data.Set("samlv2_configuration.0.enabled", a.Samlv2Configuration.Enabled)
	_ = data.Set("samlv2_configuration.0.audience", a.Samlv2Configuration.Audience)
	_ = data.Set("samlv2_configuration.0.callback_url", a.Samlv2Configuration.CallbackURL)
	_ = data.Set("samlv2_configuration.0.debug", a.Samlv2Configuration.Debug)
	_ = data.Set("samlv2_configuration.0.issuer", a.Samlv2Configuration.Issuer)
	_ = data.Set("samlv2_configuration.0.key_id", a.Samlv2Configuration.KeyId)
	_ = data.Set("samlv2_configuration.0.logout_url", a.Samlv2Configuration.LogoutURL)
	_ = data.Set(
		"samlv2_configuration.0.xml_signature_canonicalization_method",
		a.Samlv2Configuration.XmlSignatureC14nMethod,
	)
	_ = data.Set("verification_email_template_id", a.VerificationEmailTemplateId)
	_ = data.Set("verify_registration", a.VerifyRegistration)

	_ = data.Set("email_configuration.0.email_verification_template_id", a.EmailConfiguration.EmailVerificationEmailTemplateId)
	_ = data.Set("email_configuration.0.forgot_password_template_id", a.EmailConfiguration.ForgotPasswordEmailTemplateId)
	_ = data.Set("email_configuration.0.passwordless_email_template_id", a.EmailConfiguration.PasswordlessEmailTemplateId)
	_ = data.Set("email_configuration.0.set_password_email_template_id", a.EmailConfiguration.SetPasswordEmailTemplateId)
}
