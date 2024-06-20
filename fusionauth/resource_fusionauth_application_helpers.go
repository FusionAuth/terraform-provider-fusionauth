package fusionauth

import (
	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func buildApplication(data *schema.ResourceData) fusionauth.Application {
	a := fusionauth.Application{
		TenantId: data.Get("tenant_id").(string),
		AuthenticationTokenConfiguration: fusionauth.AuthenticationTokenConfiguration{
			Enableable: buildEnableable("authentication_token_configuration_enabled", data),
		},
		AccessControlConfiguration: fusionauth.ApplicationAccessControlConfiguration{
			UiIPAccessControlListId: data.Get("access_control_configuration.0.ui_ip_access_control_list_id").(string),
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
			SelfServiceFormId:       data.Get("form_configuration.0.self_service_form_id").(string),
		},
		JwtConfiguration: fusionauth.JWTConfiguration{
			Enableable:                      buildEnableable("jwt_configuration.0.enabled", data),
			AccessTokenKeyId:                data.Get("jwt_configuration.0.access_token_id").(string),
			IdTokenKeyId:                    data.Get("jwt_configuration.0.id_token_key_id").(string),
			RefreshTokenTimeToLiveInMinutes: data.Get("jwt_configuration.0.refresh_token_ttl_minutes").(int),
			TimeToLiveInSeconds:             data.Get("jwt_configuration.0.ttl_seconds").(int),
			RefreshTokenExpirationPolicy:    fusionauth.RefreshTokenExpirationPolicy(data.Get("jwt_configuration.0.refresh_token_expiration_policy").(string)),
			RefreshTokenUsagePolicy:         fusionauth.RefreshTokenUsagePolicy(data.Get("jwt_configuration.0.refresh_token_usage_policy").(string)),
			RefreshTokenSlidingWindowConfiguration: fusionauth.RefreshTokenSlidingWindowConfiguration{
				MaximumTimeToLiveInMinutes: data.Get("jwt_configuration.0.refresh_token_sliding_window_maximum_time_to_live_in_minutes").(int),
			},
		},
		LambdaConfiguration: fusionauth.LambdaConfiguration{
			AccessTokenPopulateId:               data.Get("lambda_configuration.0.access_token_populate_id").(string),
			IdTokenPopulateId:                   data.Get("lambda_configuration.0.id_token_populate_id").(string),
			Samlv2PopulateId:                    data.Get("lambda_configuration.0.samlv2_populate_id").(string),
			SelfServiceRegistrationValidationId: data.Get("lambda_configuration.0.self_service_registration_validation_id").(string),
			UserinfoPopulateId:                  data.Get("lambda_configuration.0.userinfo_populate_id").(string),
		},
		LoginConfiguration: fusionauth.LoginConfiguration{
			AllowTokenRefresh:     data.Get("login_configuration.0.allow_token_refresh").(bool),
			GenerateRefreshTokens: data.Get("login_configuration.0.generate_refresh_tokens").(bool),
			RequireAuthentication: data.Get("login_configuration.0.require_authentication").(bool),
		},
		MultiFactorConfiguration: fusionauth.ApplicationMultiFactorConfiguration{
			Email: fusionauth.MultiFactorEmailTemplate{
				TemplateId: data.Get("multi_factor_configuration.0.email_template_id").(string),
			},
			LoginPolicy: fusionauth.MultiFactorLoginPolicy(data.Get("multi_factor_configuration.0.login_policy").(string)),
			Sms: fusionauth.MultiFactorSMSTemplate{
				TemplateId: data.Get("multi_factor_configuration.0.sms_template_id").(string),
			},
			TrustPolicy: fusionauth.ApplicationMultiFactorTrustPolicy(data.Get("multi_factor_configuration.0.trust_policy").(string)),
		},
		Name: data.Get("name").(string),
		OauthConfiguration: fusionauth.OAuth2Configuration{
			AuthorizedOriginURLs:          handleStringSlice("oauth_configuration.0.authorized_origin_urls", data),
			AuthorizedRedirectURLs:        handleStringSliceFromList(data.Get("oauth_configuration.0.authorized_redirect_urls").([]interface{})),
			AuthorizedURLValidationPolicy: fusionauth.Oauth2AuthorizedURLValidationPolicy(data.Get("oauth_configuration.0.authorized_url_validation_policy").(string)),
			ClientAuthenticationPolicy:    fusionauth.ClientAuthenticationPolicy(data.Get("oauth_configuration.0.client_authentication_policy").(string)),
			ClientSecret:                  data.Get("oauth_configuration.0.client_secret").(string),
			ConsentMode:                   fusionauth.OAuthScopeConsentMode(data.Get("oauth_configuration.0.consent_mode").(string)),
			Debug:                         data.Get("oauth_configuration.0.debug").(bool),
			DeviceVerificationURL:         data.Get("oauth_configuration.0.device_verification_url").(string),
			EnabledGrants:                 buildGrants("oauth_configuration.0.enabled_grants", data),
			GenerateRefreshTokens:         data.Get("oauth_configuration.0.generate_refresh_tokens").(bool),
			LogoutBehavior:                fusionauth.LogoutBehavior(data.Get("oauth_configuration.0.logout_behavior").(string)),
			LogoutURL:                     data.Get("oauth_configuration.0.logout_url").(string),
			ProofKeyForCodeExchangePolicy: fusionauth.ProofKeyForCodeExchangePolicy(data.Get("oauth_configuration.0.proof_key_for_code_exchange_policy").(string)),
			ProvidedScopePolicy: fusionauth.ProvidedScopePolicy{
				Address: buildRequireable("oauth_configuration.0.provided_scope_policy.0.address", data),
				Email:   buildRequireable("oauth_configuration.0.provided_scope_policy.0.email", data),
				Phone:   buildRequireable("oauth_configuration.0.provided_scope_policy.0.phone", data),
				Profile: buildRequireable("oauth_configuration.0.provided_scope_policy.0.profile", data),
			},
			Relationship:                fusionauth.OAuthApplicationRelationship(data.Get("oauth_configuration.0.relationship").(string)),
			RequireClientAuthentication: data.Get("oauth_configuration.0.require_client_authentication").(bool),
			RequireRegistration:         data.Get("oauth_configuration.0.require_registration").(bool),
			ScopeHandlingPolicy:         fusionauth.OAuthScopeHandlingPolicy(data.Get("oauth_configuration.0.scope_handling_policy").(string)),
			UnknownScopePolicy:          fusionauth.UnknownScopePolicy(data.Get("oauth_configuration.0.unknown_scope_policy").(string)),
		},
		PasswordlessConfiguration: fusionauth.PasswordlessConfiguration{
			Enableable: buildEnableable("passwordless_configuration_enabled", data),
		},
		RegistrationConfiguration: fusionauth.RegistrationConfiguration{
			Enableable:         buildEnableable("registration_configuration.0.enabled", data),
			BirthDate:          buildRequireable("registration_configuration.0.birth_date", data),
			ConfirmPassword:    data.Get("registration_configuration.0.confirm_password").(bool),
			FormId:             data.Get("registration_configuration.0.form_id").(string),
			FirstName:          buildRequireable("registration_configuration.0.first_name", data),
			FullName:           buildRequireable("registration_configuration.0.full_name", data),
			LastName:           buildRequireable("registration_configuration.0.last_name", data),
			MiddleName:         buildRequireable("registration_configuration.0.middle_name", data),
			MobilePhone:        buildRequireable("registration_configuration.0.mobile_phone", data),
			PreferredLanguages: buildRequireable("registration_configuration.0.preferred_languages", data),
			LoginIdType:        fusionauth.LoginIdType(data.Get("registration_configuration.0.login_id_type").(string)),
			Type:               fusionauth.RegistrationType(data.Get("registration_configuration.0.type").(string)),
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
			AuthorizedRedirectURLs:   handleStringSliceFromList(data.Get("samlv2_configuration.0.authorized_redirect_urls").([]interface{})),
			CallbackURL:              data.Get("samlv2_configuration.0.callback_url").(string),
			Debug:                    data.Get("samlv2_configuration.0.debug").(bool),
			DefaultVerificationKeyId: data.Get("samlv2_configuration.0.default_verification_key_id").(string),
			Issuer:                   data.Get("samlv2_configuration.0.issuer").(string),
			KeyId:                    data.Get("samlv2_configuration.0.key_id").(string),
			LogoutURL:                data.Get("samlv2_configuration.0.logout_url").(string),
			Logout: fusionauth.SAMLv2Logout{
				Behavior:                 fusionauth.SAMLLogoutBehavior(data.Get("samlv2_configuration.0.logout.0.behavior").(string)),
				DefaultVerificationKeyId: data.Get("samlv2_configuration.0.logout.0.default_verification_key_id").(string),
				KeyId:                    data.Get("samlv2_configuration.0.logout.0.key_id").(string),
				RequireSignedRequests:    data.Get("samlv2_configuration.0.logout.0.require_signed_requests").(bool),
				SingleLogout: fusionauth.SAMLv2SingleLogout{
					Enableable: buildEnableable("samlv2_configuration.0.logout.0.single_logout.0.enabled", data),
					KeyId:      data.Get("samlv2_configuration.0.logout.0.single_logout.0.key_id").(string),
					Url:        data.Get("samlv2_configuration.0.logout.0.single_logout.0.url").(string),
					XmlSignatureC14nMethod: fusionauth.CanonicalizationMethod(
						data.Get("samlv2_configuration.0.logout.0.single_logout.0.xml_signature_canonicalization_method").(string),
					),
				},
				XmlSignatureC14nMethod: fusionauth.CanonicalizationMethod(
					data.Get("samlv2_configuration.0.logout.0.xml_signature_canonicalization_method").(string),
				),
			},
			RequireSignedRequests: data.Get("samlv2_configuration.0.required_signed_requests").(bool),
			XmlSignatureC14nMethod: fusionauth.CanonicalizationMethod(
				data.Get("samlv2_configuration.0.xml_signature_canonicalization_method").(string),
			),
			XmlSignatureLocation: fusionauth.XMLSignatureLocation(
				data.Get("samlv2_configuration.0.xml_signature_location").(string),
			),
		},
		ThemeId:                     data.Get("theme_id").(string),
		VerificationEmailTemplateId: data.Get("verification_email_template_id").(string),
		VerificationStrategy:        fusionauth.VerificationStrategy(data.Get("verification_strategy").(string)),
		VerifyRegistration:          data.Get("verify_registration").(bool),
		EmailConfiguration: fusionauth.ApplicationEmailConfiguration{
			EmailVerificationEmailTemplateId:     data.Get("email_configuration.0.email_verification_template_id").(string),
			EmailUpdateEmailTemplateId:           data.Get("email_configuration.0.email_update_template_id").(string),
			EmailVerifiedEmailTemplateId:         data.Get("email_configuration.0.email_verified_template_id").(string),
			ForgotPasswordEmailTemplateId:        data.Get("email_configuration.0.forgot_password_template_id").(string),
			LoginIdInUseOnCreateEmailTemplateId:  data.Get("email_configuration.0.login_id_in_use_on_create_template_id").(string),
			LoginIdInUseOnUpdateEmailTemplateId:  data.Get("email_configuration.0.login_id_in_use_on_update_template_id").(string),
			LoginNewDeviceEmailTemplateId:        data.Get("email_configuration.0.login_new_device_template_id").(string),
			LoginSuspiciousEmailTemplateId:       data.Get("email_configuration.0.login_suspicious_template_id").(string),
			PasswordlessEmailTemplateId:          data.Get("email_configuration.0.passwordless_email_template_id").(string),
			PasswordResetSuccessEmailTemplateId:  data.Get("email_configuration.0.password_reset_success_template_id").(string),
			PasswordUpdateEmailTemplateId:        data.Get("email_configuration.0.password_update_template_id").(string),
			SetPasswordEmailTemplateId:           data.Get("email_configuration.0.set_password_email_template_id").(string),
			TwoFactorMethodAddEmailTemplateId:    data.Get("email_configuration.0.two_factor_method_add_template_id").(string),
			TwoFactorMethodRemoveEmailTemplateId: data.Get("email_configuration.0.two_factor_method_remove_template_id").(string),
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

func buildResourceDataFromApplication(a fusionauth.Application, data *schema.ResourceData) diag.Diagnostics {
	if err := data.Set("tenant_id", a.TenantId); err != nil {
		return diag.Errorf("application.tenant_id: %s", err.Error())
	}
	if err := data.Set("authentication_token_configuration_enabled", a.AuthenticationTokenConfiguration.Enabled); err != nil {
		return diag.Errorf("application.authentication_token_configuration_enabled: %s", err.Error())
	}

	err := data.Set("access_control_configuration", []map[string]interface{}{
		{
			"ui_ip_access_control_list_id": a.AccessControlConfiguration.UiIPAccessControlListId,
		},
	})
	if err != nil {
		return diag.Errorf("application.access_control_configuration: %s", err.Error())
	}

	err = data.Set("clean_speak_configuration", []map[string]interface{}{
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
		return diag.Errorf("application.clean_speak_configuration: %s", err.Error())
	}

	if err := data.Set("data", a.Data); err != nil {
		return diag.Errorf("application.data: %s", err.Error())
	}

	err = data.Set("form_configuration", []map[string]interface{}{
		{
			"admin_registration_form_id": a.FormConfiguration.AdminRegistrationFormId,
			"self_service_form_id":       a.FormConfiguration.SelfServiceFormId,
		},
	})
	if err != nil {
		return diag.Errorf("application.form_configuration: %s", err.Error())
	}

	err = data.Set("jwt_configuration", []map[string]interface{}{
		{
			"enabled":                         a.JwtConfiguration.Enabled,
			"access_token_id":                 a.JwtConfiguration.AccessTokenKeyId,
			"id_token_key_id":                 a.JwtConfiguration.IdTokenKeyId,
			"refresh_token_expiration_policy": a.JwtConfiguration.RefreshTokenExpirationPolicy,
			"refresh_token_sliding_window_maximum_time_to_live_in_minutes": a.JwtConfiguration.RefreshTokenSlidingWindowConfiguration.MaximumTimeToLiveInMinutes,
			"refresh_token_ttl_minutes":                                    a.JwtConfiguration.RefreshTokenTimeToLiveInMinutes,
			"refresh_token_usage_policy":                                   a.JwtConfiguration.RefreshTokenUsagePolicy,
			"ttl_seconds":                                                  a.JwtConfiguration.TimeToLiveInSeconds,
		},
	})
	if err != nil {
		return diag.Errorf("application.jwt_configuration: %s", err.Error())
	}

	err = data.Set("lambda_configuration", []map[string]interface{}{
		{
			"access_token_populate_id":                a.LambdaConfiguration.AccessTokenPopulateId,
			"id_token_populate_id":                    a.LambdaConfiguration.IdTokenPopulateId,
			"samlv2_populate_id":                      a.LambdaConfiguration.Samlv2PopulateId,
			"self_service_registration_validation_id": a.LambdaConfiguration.SelfServiceRegistrationValidationId,
			"userinfo_populate_id":                    a.LambdaConfiguration.UserinfoPopulateId,
		},
	})
	if err != nil {
		return diag.Errorf("application.lambda_configuration: %s", err.Error())
	}

	err = data.Set("login_configuration", []map[string]interface{}{
		{
			"allow_token_refresh":     a.LoginConfiguration.AllowTokenRefresh,
			"generate_refresh_tokens": a.LoginConfiguration.GenerateRefreshTokens,
			"require_authentication":  a.LoginConfiguration.RequireAuthentication,
		},
	})
	if err != nil {
		return diag.Errorf("application.login_configuration: %s", err.Error())
	}

	err = data.Set("multi_factor_configuration", []map[string]interface{}{
		{
			"email_template_id": a.MultiFactorConfiguration.Email.TemplateId,
			"sms_template_id":   a.MultiFactorConfiguration.Sms.TemplateId,
			"login_policy":      a.MultiFactorConfiguration.LoginPolicy,
			"trust_policy":      a.MultiFactorConfiguration.TrustPolicy,
		},
	})
	if err != nil {
		return diag.Errorf("application.multi_factor_configuration: %s", err.Error())
	}

	if err := data.Set("name", a.Name); err != nil {
		return diag.Errorf("application.name: %s", err.Error())
	}

	err = data.Set("oauth_configuration", []map[string]interface{}{
		{
			"authorized_origin_urls":             a.OauthConfiguration.AuthorizedOriginURLs,
			"authorized_redirect_urls":           a.OauthConfiguration.AuthorizedRedirectURLs,
			"authorized_url_validation_policy":   a.OauthConfiguration.AuthorizedURLValidationPolicy,
			"client_authentication_policy":       a.OauthConfiguration.ClientAuthenticationPolicy,
			"client_secret":                      a.OauthConfiguration.ClientSecret,
			"client_id":                          a.OauthConfiguration.ClientId,
			"consent_mode":                       a.OauthConfiguration.ConsentMode,
			"debug":                              a.OauthConfiguration.Debug,
			"device_verification_url":            a.OauthConfiguration.DeviceVerificationURL,
			"enabled_grants":                     a.OauthConfiguration.EnabledGrants,
			"generate_refresh_tokens":            a.OauthConfiguration.GenerateRefreshTokens,
			"logout_behavior":                    a.OauthConfiguration.LogoutBehavior,
			"logout_url":                         a.OauthConfiguration.LogoutURL,
			"proof_key_for_code_exchange_policy": a.OauthConfiguration.ProofKeyForCodeExchangePolicy,
			"provided_scope_policy": []map[string]interface{}{
				{
					"address": []map[string]interface{}{
						{
							"enabled":  a.OauthConfiguration.ProvidedScopePolicy.Address.Enabled,
							"required": a.OauthConfiguration.ProvidedScopePolicy.Address.Required,
						},
					},
					"email": []map[string]interface{}{
						{
							"enabled":  a.OauthConfiguration.ProvidedScopePolicy.Email.Enabled,
							"required": a.OauthConfiguration.ProvidedScopePolicy.Email.Required,
						},
					},
					"phone": []map[string]interface{}{
						{
							"enabled":  a.OauthConfiguration.ProvidedScopePolicy.Phone.Enabled,
							"required": a.OauthConfiguration.ProvidedScopePolicy.Phone.Required,
						},
					},
					"profile": []map[string]interface{}{
						{
							"enabled":  a.OauthConfiguration.ProvidedScopePolicy.Profile.Enabled,
							"required": a.OauthConfiguration.ProvidedScopePolicy.Profile.Required,
						},
					},
				},
			},
			"relationship":                  a.OauthConfiguration.Relationship,
			"require_client_authentication": a.OauthConfiguration.RequireClientAuthentication,
			"require_registration":          a.OauthConfiguration.RequireRegistration,
			"scope_handling_policy":         a.OauthConfiguration.ScopeHandlingPolicy,
			"unknown_scope_policy":          a.OauthConfiguration.UnknownScopePolicy,
		},
	})
	if err != nil {
		return diag.Errorf("application.oauth_configuration: %s", err.Error())
	}

	if err := data.Set("passwordless_configuration_enabled", a.PasswordlessConfiguration.Enabled); err != nil {
		return diag.Errorf("application.passwordless_configuration_enabled: %s", err.Error())
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
			"preferred_languages": []map[string]interface{}{
				{
					"enabled":  a.RegistrationConfiguration.PreferredLanguages.Enabled,
					"required": a.RegistrationConfiguration.PreferredLanguages.Required,
				},
			},
			"login_id_type": a.RegistrationConfiguration.LoginIdType,
			"type":          a.RegistrationConfiguration.Type,
			"form_id":       a.RegistrationConfiguration.FormId,
		},
	})
	if err != nil {
		return diag.Errorf("application.registration_configuration: %s", err.Error())
	}

	err = data.Set("registration_delete_policy", []map[string]interface{}{
		{
			"unverified_enabled":                  a.RegistrationDeletePolicy.Unverified.Enabled,
			"unverified_number_of_days_to_retain": a.RegistrationDeletePolicy.Unverified.NumberOfDaysToRetain,
		},
	})
	if err != nil {
		return diag.Errorf("application.registration_delete_policy: %s", err.Error())
	}

	err = data.Set("samlv2_configuration", []map[string]interface{}{
		{
			"enabled":                     a.Samlv2Configuration.Enabled,
			"audience":                    a.Samlv2Configuration.Audience,
			"authorized_redirect_urls":    a.Samlv2Configuration.AuthorizedRedirectURLs,
			"callback_url":                a.Samlv2Configuration.CallbackURL,
			"debug":                       a.Samlv2Configuration.Debug,
			"default_verification_key_id": a.Samlv2Configuration.DefaultVerificationKeyId,
			"issuer":                      a.Samlv2Configuration.Issuer,
			"key_id":                      a.Samlv2Configuration.KeyId,
			"logout": []map[string]interface{}{
				{
					"behavior":                    a.Samlv2Configuration.Logout.Behavior,
					"default_verification_key_id": a.Samlv2Configuration.Logout.DefaultVerificationKeyId,
					"key_id":                      a.Samlv2Configuration.Logout.KeyId,
					"require_signed_requests":     a.Samlv2Configuration.Logout.RequireSignedRequests,
					"single_logout": []map[string]interface{}{
						{
							"enabled":                               a.Samlv2Configuration.Logout.SingleLogout.Enabled,
							"key_id":                                a.Samlv2Configuration.Logout.SingleLogout.KeyId,
							"url":                                   a.Samlv2Configuration.Logout.SingleLogout.Url,
							"xml_signature_canonicalization_method": a.Samlv2Configuration.Logout.SingleLogout.XmlSignatureC14nMethod,
						},
					},
					"xml_signature_canonicalization_method": a.Samlv2Configuration.Logout.XmlSignatureC14nMethod,
				},
			},
			"logout_url":                            a.Samlv2Configuration.LogoutURL,
			"required_signed_requests":              a.Samlv2Configuration.RequireSignedRequests,
			"xml_signature_canonicalization_method": a.Samlv2Configuration.XmlSignatureC14nMethod,
			"xml_signature_location":                a.Samlv2Configuration.XmlSignatureLocation,
		},
	})
	if err != nil {
		return diag.Errorf("application.samlv2_configuration: %s", err.Error())
	}

	if err := data.Set("verification_email_template_id", a.VerificationEmailTemplateId); err != nil {
		return diag.Errorf("application.verification_email_template_id: %s", err.Error())
	}
	if err := data.Set("verification_strategy", a.VerificationStrategy); err != nil {
		return diag.Errorf("application.verification_strategy: %s", err.Error())
	}
	if err := data.Set("theme_id", a.ThemeId); err != nil {
		return diag.Errorf("application.theme_id: %s", err.Error())
	}

	if err := data.Set("verify_registration", a.VerifyRegistration); err != nil {
		return diag.Errorf("application.verify_registration: %s", err.Error())
	}

	err = data.Set("email_configuration", []map[string]interface{}{
		{
			"email_verification_template_id":        a.EmailConfiguration.EmailVerificationEmailTemplateId,
			"email_update_template_id":              a.EmailConfiguration.EmailUpdateEmailTemplateId,
			"email_verified_template_id":            a.EmailConfiguration.EmailVerifiedEmailTemplateId,
			"forgot_password_template_id":           a.EmailConfiguration.ForgotPasswordEmailTemplateId,
			"login_id_in_use_on_create_template_id": a.EmailConfiguration.LoginIdInUseOnCreateEmailTemplateId,
			"login_id_in_use_on_update_template_id": a.EmailConfiguration.LoginIdInUseOnUpdateEmailTemplateId,
			"login_new_device_template_id":          a.EmailConfiguration.LoginNewDeviceEmailTemplateId,
			"login_suspicious_template_id":          a.EmailConfiguration.LoginSuspiciousEmailTemplateId,
			"passwordless_email_template_id":        a.EmailConfiguration.PasswordlessEmailTemplateId,
			"password_reset_success_template_id":    a.EmailConfiguration.PasswordResetSuccessEmailTemplateId,
			"password_update_template_id":           a.EmailConfiguration.PasswordUpdateEmailTemplateId,
			"set_password_email_template_id":        a.EmailConfiguration.SetPasswordEmailTemplateId,
			"two_factor_method_add_template_id":     a.EmailConfiguration.TwoFactorMethodAddEmailTemplateId,
			"two_factor_method_remove_template_id":  a.EmailConfiguration.TwoFactorMethodRemoveEmailTemplateId,
		},
	})
	if err != nil {
		return diag.Errorf("application.email_configuration: %s", err.Error())
	}

	return nil
}
