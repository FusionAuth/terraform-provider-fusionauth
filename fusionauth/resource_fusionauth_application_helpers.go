package fusionauth

import (
	"fmt"

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
		FormConfiguration: fusionauth.ApplicationFormConfiguration{
			AdminRegistrationFormId: data.Get("form_configuration.0.admin_registration_form_id").(string),
		},
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
			Enableable:               buildEnableable("samlv2_configuration.0.enabled", data),
			Audience:                 data.Get("samlv2_configuration.0.audience").(string),
			CallbackURL:              data.Get("samlv2_configuration.0.callback_url").(string),
			Debug:                    data.Get("samlv2_configuration.0.debug").(bool),
			DefaultVerificationKeyId: data.Get("samlv2_configuration.0.default_verification_key_id").(string),
			Issuer:                   data.Get("samlv2_configuration.0.issuer").(string),
			KeyId:                    data.Get("samlv2_configuration.0.key_id").(string),
			LogoutURL:                data.Get("samlv2_configuration.0.logout_url").(string),
			RequireSignedRequests:    data.Get("samlv2_configuration.0.required_signed_requests").(bool),
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

func buildResourceDataFromApplication(a fusionauth.Application, data *schema.ResourceData) error {
	if err := data.Set("tenant_id", a.TenantId); err != nil {
		return fmt.Errorf("application.tenant_id: %s", err.Error())
	}
	if err := data.Set("authentication_token_configuration_enabled", a.AuthenticationTokenConfiguration.Enabled); err != nil {
		return fmt.Errorf("application.authentication_token_configuration_enabled: %s", err.Error())
	}

	err := data.Set("clean_speak_configuration", []map[string]interface{}{
		{
			"application_ids": a.CleanSpeakConfiguration.ApplicationIds,
			"username_moderation": []map[string]interface{}{
				{
					"application_id": a.CleanSpeakConfiguration.UsernameModeration.ApplicationId,
					"enabled":        a.CleanSpeakConfiguration.UsernameModeration.Enabled,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("application.clean_speak_configuration: %s", err.Error())
	}

	if err := data.Set("data", a.Data); err != nil {
		return fmt.Errorf("application.data: %s", err.Error())
	}

	err = data.Set("form_configuration", []map[string]interface{}{
		{
			"admin_registration_form_id": a.FormConfiguration.AdminRegistrationFormId,
		},
	})
	if err != nil {
		return fmt.Errorf("application.form_configuration: %s", err.Error())
	}

	err = data.Set("jwt_configuration", []map[string]interface{}{
		{
			"enabled":                   a.JwtConfiguration.Enabled,
			"access_token_id":           a.JwtConfiguration.AccessTokenKeyId,
			"id_token_key_id":           a.JwtConfiguration.IdTokenKeyId,
			"refresh_token_ttl_minutes": a.JwtConfiguration.RefreshTokenTimeToLiveInMinutes,
			"ttl_seconds":               a.JwtConfiguration.TimeToLiveInSeconds,
		},
	})
	if err != nil {
		return fmt.Errorf("application.jwt_configuration: %s", err.Error())
	}

	err = data.Set("lambda_configuration", []map[string]interface{}{
		{
			"access_token_populate_id": a.LambdaConfiguration.AccessTokenPopulateId,
			"id_token_populate_id":     a.LambdaConfiguration.IdTokenPopulateId,
			"samlv2_populate_id":       a.LambdaConfiguration.Samlv2PopulateId,
		},
	})
	if err != nil {
		return fmt.Errorf("application.lambda_configuration: %s", err.Error())
	}

	err = data.Set("login_configuration", []map[string]interface{}{
		{
			"allow_token_refresh":     a.LoginConfiguration.AllowTokenRefresh,
			"generate_refresh_tokens": a.LoginConfiguration.GenerateRefreshTokens,
			"require_authentication":  a.LoginConfiguration.RequireAuthentication,
		},
	})
	if err != nil {
		return fmt.Errorf("application.login_configuration: %s", err.Error())
	}

	if err := data.Set("name", a.Name); err != nil {
		return fmt.Errorf("application.name: %s", err.Error())
	}

	err = data.Set("oauth_configuration", []map[string]interface{}{
		{
			"authorized_origin_urls":        a.OauthConfiguration.AuthorizedOriginURLs,
			"authorized_redirect_urls":      a.OauthConfiguration.AuthorizedRedirectURLs,
			"client_secret":                 a.OauthConfiguration.ClientSecret,
			"device_verification_url":       a.OauthConfiguration.DeviceVerificationURL,
			"generate_refresh_tokens":       a.OauthConfiguration.GenerateRefreshTokens,
			"logout_url":                    a.OauthConfiguration.LogoutURL,
			"require_client_authentication": a.OauthConfiguration.RequireClientAuthentication,
			"logout_behavior":               a.OauthConfiguration.LogoutBehavior,
			"enabled_grants":                a.OauthConfiguration.EnabledGrants,
		},
	})
	if err != nil {
		return fmt.Errorf("application.oauth_configuration: %s", err.Error())
	}

	if err := data.Set("passwordless_configuration_enabled", a.PasswordlessConfiguration.Enabled); err != nil {
		return fmt.Errorf("application.passwordless_configuration_enabled: %s", err.Error())
	}

	err = data.Set("registration_configuration", []map[string]interface{}{
		{
			"enabled": a.RegistrationConfiguration.Enabled,
			"birth_date": []map[string]interface{}{
				{
					"enabled":  a.RegistrationConfiguration.BirthDate.Enabled,
					"required": a.RegistrationConfiguration.BirthDate.Required,
				},
			},
			"confirm_password": a.RegistrationConfiguration.ConfirmPassword,
			"first_name": []map[string]interface{}{
				{
					"enabled":  a.RegistrationConfiguration.FirstName.Enabled,
					"required": a.RegistrationConfiguration.FirstName.Required,
				},
			},
			"full_name": []map[string]interface{}{
				{
					"enabled":  a.RegistrationConfiguration.FullName.Enabled,
					"required": a.RegistrationConfiguration.FullName.Required,
				},
			},
			"last_name": []map[string]interface{}{
				{
					"enabled":  a.RegistrationConfiguration.LastName.Enabled,
					"required": a.RegistrationConfiguration.LastName.Required,
				},
			},
			"middle_name": []map[string]interface{}{
				{
					"enabled":  a.RegistrationConfiguration.MiddleName.Enabled,
					"required": a.RegistrationConfiguration.MiddleName.Required,
				},
			},
			"mobile_phone": []map[string]interface{}{
				{
					"enabled":  a.RegistrationConfiguration.MobilePhone.Enabled,
					"required": a.RegistrationConfiguration.MobilePhone.Required,
				},
			},
			"login_id_type": a.RegistrationConfiguration.LoginIdType,
			"type":          a.RegistrationConfiguration.Type,
			"form_id":       a.RegistrationConfiguration.FormId,
		},
	})
	if err != nil {
		return fmt.Errorf("application.registration_configuration: %s", err.Error())
	}

	err = data.Set("registration_delete_policy", []map[string]interface{}{
		{
			"unverified_enabled":                  a.RegistrationDeletePolicy.Unverified.Enabled,
			"unverified_number_of_days_to_retain": a.RegistrationDeletePolicy.Unverified.NumberOfDaysToRetain,
		},
	})
	if err != nil {
		return fmt.Errorf("application.registration_delete_policy: %s", err.Error())
	}

	err = data.Set("samlv2_configuration", []map[string]interface{}{
		{
			"enabled":                               a.Samlv2Configuration.Enabled,
			"audience":                              a.Samlv2Configuration.Audience,
			"callback_url":                          a.Samlv2Configuration.CallbackURL,
			"debug":                                 a.Samlv2Configuration.Debug,
			"default_verification_key_id":           a.Samlv2Configuration.DefaultVerificationKeyId,
			"issuer":                                a.Samlv2Configuration.Issuer,
			"key_id":                                a.Samlv2Configuration.KeyId,
			"logout_url":                            a.Samlv2Configuration.LogoutURL,
			"required_signed_requests":              a.Samlv2Configuration.RequireSignedRequests,
			"xml_signature_canonicalization_method": a.Samlv2Configuration.XmlSignatureC14nMethod,
		},
	})
	if err != nil {
		return fmt.Errorf("application.samlv2_configuration: %s", err.Error())
	}

	if err := data.Set("verification_email_template_id", a.VerificationEmailTemplateId); err != nil {
		return fmt.Errorf("application.verification_email_template_id: %s", err.Error())
	}
	if err := data.Set("verify_registration", a.VerifyRegistration); err != nil {
		return fmt.Errorf("application.verify_registration: %s", err.Error())
	}

	err = data.Set("email_configuration", []map[string]interface{}{
		{
			"email_verification_template_id": a.EmailConfiguration.EmailVerificationEmailTemplateId,
			"forgot_password_template_id":    a.EmailConfiguration.ForgotPasswordEmailTemplateId,
			"passwordless_email_template_id": a.EmailConfiguration.PasswordlessEmailTemplateId,
			"set_password_email_template_id": a.EmailConfiguration.SetPasswordEmailTemplateId,
		},
	})
	if err != nil {
		return fmt.Errorf("application.email_configuration: %s", err.Error())
	}

	return nil
}
