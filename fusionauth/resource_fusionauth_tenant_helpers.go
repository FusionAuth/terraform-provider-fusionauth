package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func buildTenant(data *schema.ResourceData) (fusionauth.Tenant, diag.Diagnostics) {
	tenant := fusionauth.Tenant{
		Data: data.Get("data").(map[string]interface{}),
		EmailConfiguration: fusionauth.EmailConfiguration{
			EmailUpdateEmailTemplateId:           data.Get("email_configuration.0.email_update_email_template_id").(string),
			EmailVerifiedEmailTemplateId:         data.Get("email_configuration.0.email_verified_email_template_id").(string),
			ImplicitEmailVerificationAllowed:     data.Get("email_configuration.0.implicit_email_verification_allowed").(bool),
			LoginIdInUseOnCreateEmailTemplateId:  data.Get("email_configuration.0.login_id_in_use_on_create_email_template_id").(string),
			LoginIdInUseOnUpdateEmailTemplateId:  data.Get("email_configuration.0.login_id_in_use_on_update_email_template_id").(string),
			LoginNewDeviceEmailTemplateId:        data.Get("email_configuration.0.login_new_device_email_template_id").(string),
			LoginSuspiciousEmailTemplateId:       data.Get("email_configuration.0.login_suspicious_email_template_id").(string),
			PasswordResetSuccessEmailTemplateId:  data.Get("email_configuration.0.password_reset_success_email_template_id").(string),
			PasswordUpdateEmailTemplateId:        data.Get("email_configuration.0.password_update_email_template_id").(string),
			TwoFactorMethodAddEmailTemplateId:    data.Get("email_configuration.0.two_factor_method_add_email_template_id").(string),
			TwoFactorMethodRemoveEmailTemplateId: data.Get("email_configuration.0.two_factor_method_remove_email_template_id").(string),
			ForgotPasswordEmailTemplateId:        data.Get("email_configuration.0.forgot_password_email_template_id").(string),
			Host:                                 data.Get("email_configuration.0.host").(string),
			Password:                             data.Get("email_configuration.0.password").(string),
			PasswordlessEmailTemplateId:          data.Get("email_configuration.0.passwordless_email_template_id").(string),
			Port:                                 data.Get("email_configuration.0.port").(int),
			Properties:                           data.Get("email_configuration.0.properties").(string),
			Security:                             fusionauth.EmailSecurityType(data.Get("email_configuration.0.security").(string)),
			SetPasswordEmailTemplateId:           data.Get("email_configuration.0.set_password_email_template_id").(string),
			Username:                             data.Get("email_configuration.0.username").(string),
			VerificationEmailTemplateId:          data.Get("email_configuration.0.verification_email_template_id").(string),
			VerifyEmail:                          data.Get("email_configuration.0.verify_email").(bool),
			VerifyEmailWhenChanged:               data.Get("email_configuration.0.verify_email_when_changed").(bool),
			VerificationStrategy:                 fusionauth.VerificationStrategy(data.Get("email_configuration.0.verification_strategy").(string)),
			DefaultFromEmail:                     data.Get("email_configuration.0.default_from_email").(string),
			DefaultFromName:                      data.Get("email_configuration.0.default_from_name").(string),
			Unverified: fusionauth.EmailUnverifiedOptions{
				AllowEmailChangeWhenGated: data.Get("email_configuration.0.unverified.0.allow_email_change_when_gated").(bool),
				Behavior:                  fusionauth.UnverifiedBehavior(data.Get("email_configuration.0.unverified.0.behavior").(string)),
			},
		},
		EventConfiguration: buildEventConfiguration("event_configuration", data),
		ExternalIdentifierConfiguration: fusionauth.ExternalIdentifierConfiguration{
			AuthorizationGrantIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.authorization_grant_id_time_to_live_in_seconds",
			).(int),
			ChangePasswordIdGenerator: fusionauth.SecureGeneratorConfiguration{
				Length: data.Get("external_identifier_configuration.0.change_password_id_generator.0.length").(int),
				Type: fusionauth.SecureGeneratorType(
					data.Get("external_identifier_configuration.0.change_password_id_generator.0.type").(string),
				),
			},
			ChangePasswordIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.change_password_id_time_to_live_in_seconds",
			).(int),
			DeviceCodeTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.device_code_time_to_live_in_seconds",
			).(int),
			DeviceUserCodeIdGenerator: fusionauth.SecureGeneratorConfiguration{
				Length: data.Get("external_identifier_configuration.0.device_user_code_id_generator.0.length").(int),
				Type: fusionauth.SecureGeneratorType(
					data.Get("external_identifier_configuration.0.device_user_code_id_generator.0.type").(string),
				),
			},
			EmailVerificationIdGenerator: fusionauth.SecureGeneratorConfiguration{
				Length: data.Get("external_identifier_configuration.0.email_verification_id_generator.0.length").(int),
				Type: fusionauth.SecureGeneratorType(
					data.Get("external_identifier_configuration.0.email_verification_id_generator.0.type").(string),
				),
			},
			EmailVerificationIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.email_verification_id_time_to_live_in_seconds",
			).(int),
			ExternalAuthenticationIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.external_authentication_id_time_to_live_in_seconds",
			).(int),
			LoginIntentTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.login_intent_time_to_live_in_seconds",
			).(int),
			OneTimePasswordTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.one_time_password_time_to_live_in_seconds",
			).(int),
			PasswordlessLoginGenerator: fusionauth.SecureGeneratorConfiguration{
				Length: data.Get("external_identifier_configuration.0.passwordless_login_generator.0.length").(int),
				Type: fusionauth.SecureGeneratorType(
					data.Get("external_identifier_configuration.0.passwordless_login_generator.0.type").(string),
				),
			},
			PasswordlessLoginTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.passwordless_login_time_to_live_in_seconds").(int),
			RegistrationVerificationIdGenerator: fusionauth.SecureGeneratorConfiguration{
				Length: data.Get("external_identifier_configuration.0.registration_verification_id_generator.0.length").(int),
				Type: fusionauth.SecureGeneratorType(
					data.Get("external_identifier_configuration.0.registration_verification_id_generator.0.type").(string),
				),
			},
			RegistrationVerificationIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.registration_verification_id_time_to_live_in_seconds",
			).(int),
			Samlv2AuthNRequestIdTimeToLiveInSeconds: data.Get("external_identifier_configuration.0.saml_v2_authn_request_id_ttl_seconds").(int),
			SetupPasswordIdGenerator: fusionauth.SecureGeneratorConfiguration{
				Length: data.Get("external_identifier_configuration.0.setup_password_id_generator.0.length").(int),
				Type: fusionauth.SecureGeneratorType(
					data.Get("external_identifier_configuration.0.setup_password_id_generator.0.type").(string),
				),
			},
			SetupPasswordIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.setup_password_id_time_to_live_in_seconds",
			).(int),
			TrustTokenTimeToLiveInSeconds: data.Get("external_identifier_configuration.0.trust_token_time_to_live_in_seconds").(int),
			PendingAccountLinkTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.pending_account_link_time_to_live_in_seconds",
			).(int),
			TwoFactorIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.two_factor_id_time_to_live_in_seconds",
			).(int),
			TwoFactorOneTimeCodeIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.two_factor_one_time_code_id_time_to_live_in_seconds",
			).(int),
			TwoFactorTrustIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.two_factor_trust_id_time_to_live_in_seconds",
			).(int),
			EmailVerificationOneTimeCodeGenerator: fusionauth.SecureGeneratorConfiguration{
				Length: data.Get("external_identifier_configuration.0.email_verification_one_time_code_generator.0.length").(int),
				Type: fusionauth.SecureGeneratorType(
					data.Get("external_identifier_configuration.0.email_verification_one_time_code_generator.0.type").(string),
				),
			},
			RegistrationVerificationOneTimeCodeGenerator: fusionauth.SecureGeneratorConfiguration{
				Length: data.Get("external_identifier_configuration.0.registration_verification_one_time_code_generator.0.length").(int),
				Type: fusionauth.SecureGeneratorType(
					data.Get("external_identifier_configuration.0.registration_verification_one_time_code_generator.0.type").(string),
				),
			},
			TwoFactorOneTimeCodeIdGenerator: fusionauth.SecureGeneratorConfiguration{
				Length: data.Get("external_identifier_configuration.0.two_factor_one_time_code_id_generator.0.length").(int),
				Type: fusionauth.SecureGeneratorType(
					data.Get("external_identifier_configuration.0.two_factor_one_time_code_id_generator.0.type").(string),
				),
			},
		},
		FailedAuthenticationConfiguration: fusionauth.FailedAuthenticationConfiguration{
			ActionDuration: int64(data.Get("failed_authentication_configuration.0.action_duration").(int)),
			ActionDurationUnit: fusionauth.ExpiryUnit(
				data.Get("failed_authentication_configuration.0.action_duration_unit").(string),
			),
			ActionCancelPolicy: fusionauth.FailedAuthenticationActionCancelPolicy{
				OnPasswordReset: data.Get("failed_authentication_configuration.0.action_cancel_policy_on_password_reset").(bool),
			},
			EmailUser:           data.Get("failed_authentication_configuration.0.email_user").(bool),
			ResetCountInSeconds: data.Get("failed_authentication_configuration.0.reset_count_in_seconds").(int),
			TooManyAttempts:     data.Get("failed_authentication_configuration.0.too_many_attempts").(int),
			UserActionId:        data.Get("failed_authentication_configuration.0.user_action_id").(string),
		},
		FamilyConfiguration: fusionauth.FamilyConfiguration{
			AllowChildRegistrations:           data.Get("family_configuration.0.allow_child_registrations").(bool),
			ConfirmChildEmailTemplateId:       data.Get("family_configuration.0.confirm_child_email_template_id").(string),
			DeleteOrphanedAccounts:            data.Get("family_configuration.0.delete_orphaned_accounts").(bool),
			DeleteOrphanedAccountsDays:        data.Get("family_configuration.0.delete_orphaned_accounts_days").(int),
			Enableable:                        buildEnableable("family_configuration.0.enabled", data),
			FamilyRequestEmailTemplateId:      data.Get("family_configuration.0.family_request_email_template_id").(string),
			MaximumChildAge:                   data.Get("family_configuration.0.maximum_child_age").(int),
			MinimumOwnerAge:                   data.Get("family_configuration.0.minimum_owner_age").(int),
			ParentEmailRequired:               data.Get("family_configuration.0.parent_email_required").(bool),
			ParentRegistrationEmailTemplateId: data.Get("family_configuration.0.parent_registration_email_template_id").(string),
		},
		FormConfiguration: fusionauth.TenantFormConfiguration{
			AdminUserFormId: data.Get("form_configuration.0.admin_user_form_id").(string),
		},
		HttpSessionMaxInactiveInterval: data.Get("http_session_max_inactive_interval").(int),
		Issuer:                         data.Get("issuer").(string),
		JwtConfiguration: fusionauth.JWTConfiguration{
			AccessTokenKeyId:             data.Get("jwt_configuration.0.access_token_key_id").(string),
			IdTokenKeyId:                 data.Get("jwt_configuration.0.id_token_key_id").(string),
			RefreshTokenExpirationPolicy: fusionauth.RefreshTokenExpirationPolicy(data.Get("jwt_configuration.0.refresh_token_expiration_policy").(string)),
			RefreshTokenRevocationPolicy: fusionauth.RefreshTokenRevocationPolicy{
				OnLoginPrevented:  data.Get("jwt_configuration.0.refresh_token_revocation_policy_on_login_prevented").(bool),
				OnPasswordChanged: data.Get("jwt_configuration.0.refresh_token_revocation_policy_on_password_change").(bool),
			},
			RefreshTokenTimeToLiveInMinutes: data.Get("jwt_configuration.0.refresh_token_time_to_live_in_minutes").(int),
			RefreshTokenUsagePolicy:         fusionauth.RefreshTokenUsagePolicy(data.Get("jwt_configuration.0.refresh_token_usage_policy").(string)),
			TimeToLiveInSeconds:             data.Get("jwt_configuration.0.time_to_live_in_seconds").(int),
			RefreshTokenSlidingWindowConfiguration: fusionauth.RefreshTokenSlidingWindowConfiguration{
				MaximumTimeToLiveInMinutes: data.Get("jwt_configuration.0.refresh_token_sliding_window_maximum_time_to_live_in_minutes").(int),
			},
		},
		LoginConfiguration: fusionauth.TenantLoginConfiguration{
			RequireAuthentication: data.Get("login_configuration.0.require_authentication").(bool),
		},
		LogoutURL: data.Get("logout_url").(string),
		MaximumPasswordAge: fusionauth.MaximumPasswordAge{
			Enableable: buildEnableable("maximum_password_age.0.enabled", data),
			Days:       data.Get("maximum_password_age.0.days").(int),
		},
		MinimumPasswordAge: fusionauth.MinimumPasswordAge{
			Enableable: buildEnableable("minimum_password_age.0.enabled", data),
			Seconds:    data.Get("minimum_password_age.0.seconds").(int),
		},
		MultiFactorConfiguration: fusionauth.TenantMultiFactorConfiguration{
			LoginPolicy: fusionauth.MultiFactorLoginPolicy(data.Get("multi_factor_configuration.0.login_policy").(string)),
			Authenticator: fusionauth.MultiFactorAuthenticatorMethod{
				Enableable: buildEnableable("multi_factor_configuration.0.authenticator.0.enabled", data),
			},
			Email: fusionauth.MultiFactorEmailMethod{
				Enableable: buildEnableable("multi_factor_configuration.0.email.0.enabled", data),
				TemplateId: data.Get("multi_factor_configuration.0.email.0.template_id").(string),
			},
			Sms: fusionauth.MultiFactorSMSMethod{
				Enableable:  buildEnableable("multi_factor_configuration.0.sms.0.enabled", data),
				MessengerId: data.Get("multi_factor_configuration.0.sms.0.messenger_id").(string),
				TemplateId:  data.Get("multi_factor_configuration.0.sms.0.template_id").(string),
			},
		},
		Name: data.Get("name").(string),
		OauthConfiguration: fusionauth.TenantOAuth2Configuration{
			ClientCredentialsAccessTokenPopulateLambdaId: data.Get("oauth_configuration.0.client_credentials_access_token_populate_lambda_id").(string),
		},
		PasswordEncryptionConfiguration: fusionauth.PasswordEncryptionConfiguration{
			EncryptionScheme:       data.Get("password_encryption_configuration.0.encryption_scheme").(string),
			EncryptionSchemeFactor: data.Get("password_encryption_configuration.0.encryption_scheme_factor").(int),
			ModifyEncryptionSchemeOnLogin: data.Get(
				"password_encryption_configuration.0.modify_encryption_scheme_on_login",
			).(bool),
		},
		PasswordValidationRules: fusionauth.PasswordValidationRules{
			BreachDetection: fusionauth.PasswordBreachDetection{
				Enableable: buildEnableable("password_validation_rules.0.breach_detection.0.enabled", data),
				MatchMode: fusionauth.BreachMatchMode(
					data.Get("password_validation_rules.0.breach_detection.0.match_mode").(string),
				),
				NotifyUserEmailTemplateId: data.Get(
					"password_validation_rules.0.breach_detection.0.notify_user_email_template_id",
				).(string),
				OnLogin: fusionauth.BreachAction(
					data.Get("password_validation_rules.0.breach_detection.0.on_login").(string),
				),
			},
			MaxLength: data.Get("password_validation_rules.0.max_length").(int),
			MinLength: data.Get("password_validation_rules.0.min_length").(int),
			RememberPreviousPasswords: fusionauth.RememberPreviousPasswords{
				Enableable: buildEnableable("password_validation_rules.0.remember_previous_passwords.0.enabled", data),
				Count:      data.Get("password_validation_rules.0.remember_previous_passwords.0.count").(int),
			},
			RequireMixedCase: data.Get("password_validation_rules.0.required_mixed_case").(bool),
			RequireNonAlpha:  data.Get("password_validation_rules.0.require_non_alpha").(bool),
			RequireNumber:    data.Get("password_validation_rules.0.require_number").(bool),
			ValidateOnLogin:  data.Get("password_validation_rules.0.validate_on_login").(bool),
		},
		RateLimitConfiguration: fusionauth.TenantRateLimitConfiguration{
			FailedLogin: fusionauth.RateLimitedRequestConfiguration{
				Enableable:          buildEnableable("rate_limit_configuration.0.failed_login.0.enabled", data),
				Limit:               data.Get("rate_limit_configuration.0.failed_login.0.limit").(int),
				TimePeriodInSeconds: data.Get("rate_limit_configuration.0.failed_login.0.time_period_in_seconds").(int),
			},
			ForgotPassword: fusionauth.RateLimitedRequestConfiguration{
				Enableable:          buildEnableable("rate_limit_configuration.0.forgot_password.0.enabled", data),
				Limit:               data.Get("rate_limit_configuration.0.forgot_password.0.limit").(int),
				TimePeriodInSeconds: data.Get("rate_limit_configuration.0.forgot_password.0.time_period_in_seconds").(int),
			},
			SendEmailVerification: fusionauth.RateLimitedRequestConfiguration{
				Enableable:          buildEnableable("rate_limit_configuration.0.send_email_verification.0.enabled", data),
				Limit:               data.Get("rate_limit_configuration.0.send_email_verification.0.limit").(int),
				TimePeriodInSeconds: data.Get("rate_limit_configuration.0.send_email_verification.0.time_period_in_seconds").(int),
			},
			SendPasswordless: fusionauth.RateLimitedRequestConfiguration{
				Enableable:          buildEnableable("rate_limit_configuration.0.send_passwordless.0.enabled", data),
				Limit:               data.Get("rate_limit_configuration.0.send_passwordless.0.limit").(int),
				TimePeriodInSeconds: data.Get("rate_limit_configuration.0.send_passwordless.0.time_period_in_seconds").(int),
			},
			SendRegistrationVerification: fusionauth.RateLimitedRequestConfiguration{
				Enableable:          buildEnableable("rate_limit_configuration.0.send_registration_verification.0.enabled", data),
				Limit:               data.Get("rate_limit_configuration.0.send_registration_verification.0.limit").(int),
				TimePeriodInSeconds: data.Get("rate_limit_configuration.0.send_registration_verification.0.time_period_in_seconds").(int),
			},
			SendTwoFactor: fusionauth.RateLimitedRequestConfiguration{
				Enableable:          buildEnableable("rate_limit_configuration.0.send_two_factor.0.enabled", data),
				Limit:               data.Get("rate_limit_configuration.0.send_two_factor.0.limit").(int),
				TimePeriodInSeconds: data.Get("rate_limit_configuration.0.send_two_factor.0.time_period_in_seconds").(int),
			},
		},
		RegistrationConfiguration: fusionauth.TenantRegistrationConfiguration{
			BlockedDomains: handleStringSlice("registration_configuration.0.blocked_domains", data),
		},
		CaptchaConfiguration: fusionauth.TenantCaptchaConfiguration{
			Enableable:    buildEnableable("captcha_configuration.0.enabled", data),
			CaptchaMethod: fusionauth.CaptchaMethod(data.Get("captcha_configuration.0.captcha_method").(string)),
			SecretKey:     data.Get("captcha_configuration.0.secret_key").(string),
			SiteKey:       data.Get("captcha_configuration.0.site_key").(string),
			Threshold:     data.Get("captcha_configuration.0.threshold").(float64),
		},
		ThemeId: data.Get("theme_id").(string),
		UserDeletePolicy: fusionauth.TenantUserDeletePolicy{
			Unverified: fusionauth.TimeBasedDeletePolicy{
				Enableable:           buildEnableable("user_delete_policy.0.unverified_enabled", data),
				NumberOfDaysToRetain: data.Get("user_delete_policy.0.unverified_number_of_days_to_retain").(int),
			},
		},
		UsernameConfiguration: fusionauth.TenantUsernameConfiguration{
			Unique: fusionauth.UniqueUsernameConfiguration{
				Enableable:     buildEnableable("username_configuration.0.unique.0.enabled", data),
				NumberOfDigits: data.Get("username_configuration.0.unique.0.number_of_digits").(int),
				Separator:      data.Get("username_configuration.0.unique.0.separator").(string),
				Strategy:       fusionauth.UniqueUsernameStrategy(data.Get("username_configuration.0.unique.0.strategy").(string)),
			},
		},
	}

	connectorPolicies, connectorDiags := buildConnectorPolicies(data)
	if connectorDiags == nil {
		tenant.ConnectorPolicies = connectorPolicies
	}

	additionalheaders, emailDiags := buildAdditionalHeaders(data)
	if emailDiags == nil {
		tenant.EmailConfiguration.AdditionalHeaders = additionalheaders
	}

	return tenant, append(connectorDiags, emailDiags...)
}

func buildAdditionalHeaders(data *schema.ResourceData) (emailHeaders []fusionauth.EmailHeader, diags diag.Diagnostics) {
	emailHeadersData, ok := data.Get("email_configuration.0.additional_headers").(map[string]interface{})
	if emailHeadersData == nil || !ok {
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Unable to convert additional headers data",
				Detail:   "additional_headers unable to be typecast to map[string]interface{}",
			})
		}

		// Nothing to do here!
		return emailHeaders, diags
	}

	emailHeaders = make([]fusionauth.EmailHeader, len(emailHeadersData))
	i := 0
	for headerName, headerValue := range emailHeadersData {
		emailHeaders[i] = fusionauth.EmailHeader{
			Name:  headerName,
			Value: headerValue.(string),
		}
		i++
	}
	return emailHeaders, diags
}

func buildConnectorPolicies(data *schema.ResourceData) (connectorPolicies []fusionauth.ConnectorPolicy, diags diag.Diagnostics) {
	policiesData, ok := data.Get("connector_policy").([]interface{})
	if policiesData == nil || !ok {
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Unable to convert connector policy data",
				Detail:   "connector_policy unable to be typecast to []interface{}",
			})
		}

		// Nothing to do here!
		return connectorPolicies, diags
	}

	connectorPolicies = make([]fusionauth.ConnectorPolicy, len(policiesData))

	for i, policiesDatum := range policiesData {
		if connectorPolicy, ok := policiesDatum.(map[string]interface{}); !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Unable to convert a connector policy",
				Detail:   fmt.Sprintf("connector_policy.%d: %#+v unable to be typecast to map[string]interface{}", i, policiesDatum),
			})
		} else {
			connectorPolicies[i] = fusionauth.ConnectorPolicy{
				ConnectorId: connectorPolicy["connector_id"].(string),
				Domains:     handleStringSliceFromSet(connectorPolicy["domains"].(*schema.Set)),
				Migrate:     connectorPolicy["migrate"].(bool),
			}
		}
	}

	return connectorPolicies, diags
}

func buildEventConfiguration(key string, data *schema.ResourceData) fusionauth.EventConfiguration {
	s := data.Get(key)
	set, ok := s.(*schema.Set)
	if !ok {
		return fusionauth.EventConfiguration{}
	}
	l := set.List()

	ev := make(map[fusionauth.EventType]fusionauth.EventConfigurationData)

	for _, x := range l {
		r := x.(map[string]interface{})
		ev[fusionauth.EventType(r["event"].(string))] = fusionauth.EventConfigurationData{
			TransactionType: fusionauth.TransactionType(r["transaction_type"].(string)),
			Enableable: fusionauth.Enableable{
				Enabled: r["enabled"].(bool),
			},
		}
	}

	return fusionauth.EventConfiguration{Events: ev}
}

func buildResourceDataFromTenant(t fusionauth.Tenant, data *schema.ResourceData) diag.Diagnostics {
	if err := data.Set("tenant_id", t.Id); err != nil {
		return diag.Errorf("tenant.tenant_id: %s", err.Error())
	}

	if err := data.Set("data", t.Data); err != nil {
		return diag.Errorf("tenant.data: %s", err.Error())
	}

	connectorPolicies := []map[string]interface{}{}
	for _, policy := range t.ConnectorPolicies {
		connectorPolicies = append(connectorPolicies, map[string]interface{}{
			"connector_id": policy.ConnectorId,
			"domains":      policy.Domains,
			"migrate":      policy.Migrate,
		})
	}
	if err := data.Set("connector_policy", connectorPolicies); err != nil {
		return diag.Errorf("tenant.connector_policy: %s", err.Error())
	}

	additionalHeaders := make(map[string]string, len(t.EmailConfiguration.AdditionalHeaders))
	for _, additionalHeader := range t.EmailConfiguration.AdditionalHeaders {
		additionalHeaders[additionalHeader.Name] = additionalHeader.Value
	}

	err := data.Set("email_configuration", []map[string]interface{}{
		{
			"additional_headers":                          additionalHeaders,
			"email_update_email_template_id":              t.EmailConfiguration.EmailUpdateEmailTemplateId,
			"email_verified_email_template_id":            t.EmailConfiguration.EmailVerifiedEmailTemplateId,
			"forgot_password_email_template_id":           t.EmailConfiguration.ForgotPasswordEmailTemplateId,
			"host":                                        t.EmailConfiguration.Host,
			"implicit_email_verification_allowed":         t.EmailConfiguration.ImplicitEmailVerificationAllowed,
			"login_id_in_use_on_create_email_template_id": t.EmailConfiguration.LoginIdInUseOnCreateEmailTemplateId,
			"login_id_in_use_on_update_email_template_id": t.EmailConfiguration.LoginIdInUseOnUpdateEmailTemplateId,
			"login_new_device_email_template_id":          t.EmailConfiguration.LoginNewDeviceEmailTemplateId,
			"login_suspicious_email_template_id":          t.EmailConfiguration.LoginSuspiciousEmailTemplateId,
			"password":                                    t.EmailConfiguration.Password,
			"passwordless_email_template_id":              t.EmailConfiguration.PasswordlessEmailTemplateId,
			"password_reset_success_email_template_id":    t.EmailConfiguration.PasswordResetSuccessEmailTemplateId,
			"password_update_email_template_id":           t.EmailConfiguration.PasswordUpdateEmailTemplateId,
			"port":                                        t.EmailConfiguration.Port,
			"properties":                                  t.EmailConfiguration.Properties,
			"security":                                    t.EmailConfiguration.Security,
			"set_password_email_template_id":              t.EmailConfiguration.SetPasswordEmailTemplateId,
			"two_factor_method_add_email_template_id":     t.EmailConfiguration.TwoFactorMethodAddEmailTemplateId,
			"two_factor_method_remove_email_template_id":  t.EmailConfiguration.TwoFactorMethodRemoveEmailTemplateId,
			"username":                                    t.EmailConfiguration.Username,
			"verification_email_template_id":              t.EmailConfiguration.VerificationEmailTemplateId,
			"verification_strategy":                       t.EmailConfiguration.VerificationStrategy,
			"verify_email":                                t.EmailConfiguration.VerifyEmail,
			"verify_email_when_changed":                   t.EmailConfiguration.VerifyEmailWhenChanged,
			"default_from_email":                          t.EmailConfiguration.DefaultFromEmail,
			"default_from_name":                           t.EmailConfiguration.DefaultFromName,
			"unverified": []map[string]interface{}{{
				"allow_email_change_when_gated": t.EmailConfiguration.Unverified.AllowEmailChangeWhenGated,
				"behavior":                      t.EmailConfiguration.Unverified.Behavior,
			}},
		},
	})
	if err != nil {
		return diag.Errorf("tenant.email_configuration: %s", err.Error())
	}

	err = data.Set("external_identifier_configuration", []map[string]interface{}{
		{
			"authorization_grant_id_time_to_live_in_seconds": t.ExternalIdentifierConfiguration.AuthorizationGrantIdTimeToLiveInSeconds,
			"change_password_id_generator": []map[string]interface{}{{
				"length": t.ExternalIdentifierConfiguration.ChangePasswordIdGenerator.Length,
				"type":   t.ExternalIdentifierConfiguration.ChangePasswordIdGenerator.Type,
			}},
			"change_password_id_time_to_live_in_seconds": t.ExternalIdentifierConfiguration.ChangePasswordIdTimeToLiveInSeconds,
			"device_code_time_to_live_in_seconds":        t.ExternalIdentifierConfiguration.DeviceCodeTimeToLiveInSeconds,
			"device_user_code_id_generator": []map[string]interface{}{{
				"length": t.ExternalIdentifierConfiguration.DeviceUserCodeIdGenerator.Length,
				"type":   t.ExternalIdentifierConfiguration.DeviceUserCodeIdGenerator.Type,
			}},
			"email_verification_id_generator": []map[string]interface{}{{
				"length": t.ExternalIdentifierConfiguration.EmailVerificationIdGenerator.Length,
				"type":   t.ExternalIdentifierConfiguration.EmailVerificationIdGenerator.Type,
			}},
			"email_verification_id_time_to_live_in_seconds":      t.ExternalIdentifierConfiguration.EmailVerificationIdTimeToLiveInSeconds,
			"external_authentication_id_time_to_live_in_seconds": t.ExternalIdentifierConfiguration.ExternalAuthenticationIdTimeToLiveInSeconds,
			"login_intent_time_to_live_in_seconds":               t.ExternalIdentifierConfiguration.LoginIntentTimeToLiveInSeconds,
			"one_time_password_time_to_live_in_seconds":          t.ExternalIdentifierConfiguration.OneTimePasswordTimeToLiveInSeconds,
			"passwordless_login_generator": []map[string]interface{}{{
				"length": t.ExternalIdentifierConfiguration.PasswordlessLoginGenerator.Length,
				"type":   t.ExternalIdentifierConfiguration.PasswordlessLoginGenerator.Type,
			}},
			"passwordless_login_time_to_live_in_seconds": t.ExternalIdentifierConfiguration.PasswordlessLoginTimeToLiveInSeconds,
			"registration_verification_id_generator": []map[string]interface{}{{
				"length": t.ExternalIdentifierConfiguration.RegistrationVerificationIdGenerator.Length,
				"type":   t.ExternalIdentifierConfiguration.RegistrationVerificationIdGenerator.Type,
			}},
			"registration_verification_id_time_to_live_in_seconds": t.ExternalIdentifierConfiguration.RegistrationVerificationIdTimeToLiveInSeconds,
			"saml_v2_authn_request_id_ttl_seconds":                 t.ExternalIdentifierConfiguration.Samlv2AuthNRequestIdTimeToLiveInSeconds,
			"setup_password_id_generator": []map[string]interface{}{{
				"length": t.ExternalIdentifierConfiguration.SetupPasswordIdGenerator.Length,
				"type":   t.ExternalIdentifierConfiguration.SetupPasswordIdGenerator.Type,
			}},
			"setup_password_id_time_to_live_in_seconds":           t.ExternalIdentifierConfiguration.SetupPasswordIdTimeToLiveInSeconds,
			"trust_token_time_to_live_in_seconds":                 t.ExternalIdentifierConfiguration.TrustTokenTimeToLiveInSeconds,
			"pending_account_link_time_to_live_in_seconds":        t.ExternalIdentifierConfiguration.PendingAccountLinkTimeToLiveInSeconds,
			"two_factor_id_time_to_live_in_seconds":               t.ExternalIdentifierConfiguration.TwoFactorIdTimeToLiveInSeconds,
			"two_factor_one_time_code_id_time_to_live_in_seconds": t.ExternalIdentifierConfiguration.TwoFactorOneTimeCodeIdTimeToLiveInSeconds,
			"two_factor_trust_id_time_to_live_in_seconds":         t.ExternalIdentifierConfiguration.TwoFactorTrustIdTimeToLiveInSeconds,
			"email_verification_one_time_code_generator": []map[string]interface{}{{
				"length": t.ExternalIdentifierConfiguration.EmailVerificationOneTimeCodeGenerator.Length,
				"type":   t.ExternalIdentifierConfiguration.EmailVerificationOneTimeCodeGenerator.Type,
			}},
			"registration_verification_one_time_code_generator": []map[string]interface{}{{
				"length": t.ExternalIdentifierConfiguration.RegistrationVerificationOneTimeCodeGenerator.Length,
				"type":   t.ExternalIdentifierConfiguration.RegistrationVerificationOneTimeCodeGenerator.Type,
			}},
			"two_factor_one_time_code_id_generator": []map[string]interface{}{{
				"length": t.ExternalIdentifierConfiguration.TwoFactorOneTimeCodeIdGenerator.Length,
				"type":   t.ExternalIdentifierConfiguration.TwoFactorOneTimeCodeIdGenerator.Type,
			}},
		},
	})
	if err != nil {
		return diag.Errorf("tenant.external_identifier_configuration: %s", err.Error())
	}

	err = data.Set("failed_authentication_configuration", []map[string]interface{}{
		{
			"action_duration":                        t.FailedAuthenticationConfiguration.ActionDuration,
			"action_duration_unit":                   t.FailedAuthenticationConfiguration.ActionDurationUnit,
			"action_cancel_policy_on_password_reset": t.FailedAuthenticationConfiguration.ActionCancelPolicy.OnPasswordReset,
			"email_user":                             t.FailedAuthenticationConfiguration.EmailUser,
			"reset_count_in_seconds":                 t.FailedAuthenticationConfiguration.ResetCountInSeconds,
			"too_many_attempts":                      t.FailedAuthenticationConfiguration.TooManyAttempts,
			"user_action_id":                         t.FailedAuthenticationConfiguration.UserActionId,
		},
	})
	if err != nil {
		return diag.Errorf("tenant.failed_authentication_configuration: %s", err.Error())
	}

	err = data.Set("family_configuration", []map[string]interface{}{
		{
			"allow_child_registrations":             t.FamilyConfiguration.AllowChildRegistrations,
			"confirm_child_email_template_id":       t.FamilyConfiguration.ConfirmChildEmailTemplateId,
			"delete_orphaned_accounts":              t.FamilyConfiguration.DeleteOrphanedAccounts,
			"delete_orphaned_accounts_days":         t.FamilyConfiguration.DeleteOrphanedAccountsDays,
			"enabled":                               t.FamilyConfiguration.Enabled,
			"family_request_email_template_id":      t.FamilyConfiguration.FamilyRequestEmailTemplateId,
			"maximum_child_age":                     t.FamilyConfiguration.MaximumChildAge,
			"minimum_owner_age":                     t.FamilyConfiguration.MinimumOwnerAge,
			"parent_email_required":                 t.FamilyConfiguration.ParentEmailRequired,
			"parent_registration_email_template_id": t.FamilyConfiguration.ParentRegistrationEmailTemplateId,
		},
	})
	if err != nil {
		return diag.Errorf("tenant.family_configuration: %s", err.Error())
	}

	err = data.Set("form_configuration", []map[string]interface{}{
		{
			"admin_user_form_id": t.FormConfiguration.AdminUserFormId,
		},
	})
	if err != nil {
		return diag.Errorf("tenant.form_configuration: %s", err.Error())
	}

	if err := data.Set("http_session_max_inactive_interval", t.HttpSessionMaxInactiveInterval); err != nil {
		return diag.Errorf("tenant.http_session_max_inactive_interval: %s", err.Error())
	}

	if err := data.Set("issuer", t.Issuer); err != nil {
		return diag.Errorf("tenant.issuer: %s", err.Error())
	}

	err = data.Set("jwt_configuration", []map[string]interface{}{
		{
			"access_token_key_id":                                          t.JwtConfiguration.AccessTokenKeyId,
			"id_token_key_id":                                              t.JwtConfiguration.IdTokenKeyId,
			"refresh_token_expiration_policy":                              t.JwtConfiguration.RefreshTokenExpirationPolicy,
			"refresh_token_revocation_policy_on_login_prevented":           t.JwtConfiguration.RefreshTokenRevocationPolicy.OnLoginPrevented,
			"refresh_token_revocation_policy_on_password_change":           t.JwtConfiguration.RefreshTokenRevocationPolicy.OnPasswordChanged,
			"refresh_token_usage_policy":                                   t.JwtConfiguration.RefreshTokenUsagePolicy,
			"refresh_token_time_to_live_in_minutes":                        t.JwtConfiguration.RefreshTokenTimeToLiveInMinutes,
			"time_to_live_in_seconds":                                      t.JwtConfiguration.TimeToLiveInSeconds,
			"refresh_token_sliding_window_maximum_time_to_live_in_minutes": t.JwtConfiguration.RefreshTokenSlidingWindowConfiguration.MaximumTimeToLiveInMinutes,
		},
	})
	if err != nil {
		return diag.Errorf("tenant.jwt_configuration: %s", err.Error())
	}

	err = data.Set("login_configuration", []map[string]interface{}{
		{
			"require_authentication": t.LoginConfiguration.RequireAuthentication,
		},
	})
	if err != nil {
		return diag.Errorf("tenant.login_configuration: %s", err.Error())
	}

	if err := data.Set("logout_url", t.LogoutURL); err != nil {
		return diag.Errorf("tenant.logout_url: %s", err.Error())
	}

	err = data.Set("maximum_password_age", []map[string]interface{}{
		{
			"enabled": t.MaximumPasswordAge.Enabled,
			"days":    t.MaximumPasswordAge.Days,
		},
	})
	if err != nil {
		return diag.Errorf("tenant.maximum_password_age: %s", err.Error())
	}
	err = data.Set("minimum_password_age", []map[string]interface{}{
		{
			"enabled": t.MinimumPasswordAge.Enabled,
			"seconds": t.MinimumPasswordAge.Seconds,
		},
	})
	if err != nil {
		return diag.Errorf("tenant.minimum_password_age: %s", err.Error())
	}

	err = data.Set("multi_factor_configuration", []map[string]interface{}{
		{
			"login_policy": t.MultiFactorConfiguration.LoginPolicy,
			"authenticator": []map[string]interface{}{{
				"enabled": t.MultiFactorConfiguration.Authenticator.Enabled,
			}},
			"email": []map[string]interface{}{{
				"enabled":     t.MultiFactorConfiguration.Email.Enabled,
				"template_id": t.MultiFactorConfiguration.Email.TemplateId,
			}},
			"sms": []map[string]interface{}{{
				"enabled":      t.MultiFactorConfiguration.Sms.Enabled,
				"messenger_id": t.MultiFactorConfiguration.Sms.MessengerId,
				"template_id":  t.MultiFactorConfiguration.Sms.TemplateId,
			}},
		},
	})
	if err != nil {
		return diag.Errorf("tenant.multi_factor_configuration: %s", err.Error())
	}

	if err := data.Set("name", t.Name); err != nil {
		return diag.Errorf("tenant.name: %s", err.Error())
	}

	if lambdaID := t.OauthConfiguration.ClientCredentialsAccessTokenPopulateLambdaId; lambdaID != "" {
		err = data.Set("oauth_configuration", []map[string]interface{}{
			{
				"client_credentials_access_token_populate_lambda_id": lambdaID,
			},
		})
		if err != nil {
			return diag.Errorf("tenant.oauth_configuration: %s", err.Error())
		}
	}

	err = data.Set("password_encryption_configuration", []map[string]interface{}{
		{
			"encryption_scheme":                 t.PasswordEncryptionConfiguration.EncryptionScheme,
			"encryption_scheme_factor":          t.PasswordEncryptionConfiguration.EncryptionSchemeFactor,
			"modify_encryption_scheme_on_login": t.PasswordEncryptionConfiguration.ModifyEncryptionSchemeOnLogin,
		},
	})
	if err != nil {
		return diag.Errorf("tenant.password_encryption_configuration: %s", err.Error())
	}

	err = data.Set("captcha_configuration", []map[string]interface{}{
		{
			"enabled":        t.CaptchaConfiguration.Enabled,
			"captcha_method": t.CaptchaConfiguration.CaptchaMethod,
			"secret_key":     t.CaptchaConfiguration.SecretKey,
			"site_key":       t.CaptchaConfiguration.SiteKey,
			"threshold":      t.CaptchaConfiguration.Threshold,
		},
	})
	if err != nil {
		return diag.Errorf("tenant.captcha_configuration: %s", err.Error())
	}

	err = data.Set("password_validation_rules", []map[string]interface{}{
		{
			"breach_detection": []map[string]interface{}{{
				"enabled":                       t.PasswordValidationRules.BreachDetection.Enabled,
				"match_mode":                    t.PasswordValidationRules.BreachDetection.MatchMode,
				"notify_user_email_template_id": t.PasswordValidationRules.BreachDetection.NotifyUserEmailTemplateId,
				"on_login":                      t.PasswordValidationRules.BreachDetection.OnLogin,
			}},
			"max_length": t.PasswordValidationRules.MaxLength,
			"min_length": t.PasswordValidationRules.MinLength,
			"remember_previous_passwords": []map[string]interface{}{{
				"enabled": t.PasswordValidationRules.RememberPreviousPasswords.Enabled,
				"count":   t.PasswordValidationRules.RememberPreviousPasswords.Count,
			}},
			"required_mixed_case": t.PasswordValidationRules.RequireMixedCase,
			"require_non_alpha":   t.PasswordValidationRules.RequireNonAlpha,
			"require_number":      t.PasswordValidationRules.RequireNumber,
			"validate_on_login":   t.PasswordValidationRules.ValidateOnLogin,
		},
	})
	if err != nil {
		return diag.Errorf("tenant.password_validation_rules: %s", err.Error())
	}

	err = data.Set("rate_limit_configuration", []map[string]interface{}{
		{
			"failed_login": []map[string]interface{}{{
				"enabled":                t.RateLimitConfiguration.FailedLogin.Enabled,
				"limit":                  t.RateLimitConfiguration.FailedLogin.Limit,
				"time_period_in_seconds": t.RateLimitConfiguration.FailedLogin.TimePeriodInSeconds,
			}},
			"forgot_password": []map[string]interface{}{{
				"enabled":                t.RateLimitConfiguration.ForgotPassword.Enabled,
				"limit":                  t.RateLimitConfiguration.ForgotPassword.Limit,
				"time_period_in_seconds": t.RateLimitConfiguration.ForgotPassword.TimePeriodInSeconds,
			}},
			"send_email_verification": []map[string]interface{}{{
				"enabled":                t.RateLimitConfiguration.SendEmailVerification.Enabled,
				"limit":                  t.RateLimitConfiguration.SendEmailVerification.Limit,
				"time_period_in_seconds": t.RateLimitConfiguration.SendEmailVerification.TimePeriodInSeconds,
			}},
			"send_passwordless": []map[string]interface{}{{
				"enabled":                t.RateLimitConfiguration.SendPasswordless.Enabled,
				"limit":                  t.RateLimitConfiguration.SendPasswordless.Limit,
				"time_period_in_seconds": t.RateLimitConfiguration.SendPasswordless.TimePeriodInSeconds,
			}},
			"send_registration_verification": []map[string]interface{}{{
				"enabled":                t.RateLimitConfiguration.SendRegistrationVerification.Enabled,
				"limit":                  t.RateLimitConfiguration.SendRegistrationVerification.Limit,
				"time_period_in_seconds": t.RateLimitConfiguration.SendRegistrationVerification.TimePeriodInSeconds,
			}},
			"send_two_factor": []map[string]interface{}{{
				"enabled":                t.RateLimitConfiguration.SendTwoFactor.Enabled,
				"limit":                  t.RateLimitConfiguration.SendTwoFactor.Limit,
				"time_period_in_seconds": t.RateLimitConfiguration.SendTwoFactor.TimePeriodInSeconds,
			}},
		},
	})
	if err != nil {
		return diag.Errorf("tenant.rate_limit_configuration: %s", err.Error())
	}

	regDiag := buildAdditionalResourceDataFromTenant(t, data)
	if regDiag != nil {
		return regDiag
	}

	return nil
}

func buildAdditionalResourceDataFromTenant(t fusionauth.Tenant, data *schema.ResourceData) diag.Diagnostics {
	err := data.Set("registration_configuration", []map[string]interface{}{
		{
			"blocked_domains": t.RegistrationConfiguration.BlockedDomains,
		},
	})
	if err != nil {
		return diag.Errorf("tenant.registration_configuration: %s", err.Error())
	}

	if err := data.Set("theme_id", t.ThemeId); err != nil {
		return diag.Errorf("tenant.theme_id: %s", err.Error())
	}

	err = data.Set("user_delete_policy", []map[string]interface{}{
		{
			"unverified_enabled":                  t.UserDeletePolicy.Unverified.Enabled,
			"unverified_number_of_days_to_retain": t.UserDeletePolicy.Unverified.NumberOfDaysToRetain,
		},
	})
	if err != nil {
		return diag.Errorf("tenant.user_delete_policy: %s", err.Error())
	}

	err = data.Set("username_configuration", []map[string]interface{}{
		{
			"unique": []map[string]interface{}{{
				"enabled":          t.UsernameConfiguration.Unique.Enabled,
				"number_of_digits": t.UsernameConfiguration.Unique.NumberOfDigits,
				"separator":        t.UsernameConfiguration.Unique.Separator,
				"strategy":         t.UsernameConfiguration.Unique.Strategy,
			}},
		},
	})
	if err != nil {
		return diag.Errorf("tenant.username_configuration: %s", err.Error())
	}

	e := make([]map[string]interface{}, 0, len(t.EventConfiguration.Events))
	for k, v := range t.EventConfiguration.Events {
		e = append(e, map[string]interface{}{
			"event":            k,
			"transaction_type": v.TransactionType,
			"enabled":          v.Enabled,
		})
	}
	err = data.Set("event_configuration", e)
	if err != nil {
		return diag.Errorf("tenant.event_configuration: %s", err.Error())
	}

	return nil
}
