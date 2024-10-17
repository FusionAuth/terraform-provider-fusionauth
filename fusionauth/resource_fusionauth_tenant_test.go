package fusionauth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gpsinsight/terraform-provider-fusionauth/fusionauth/testdata"
)

func TestAccFusionauthTenant_basic(t *testing.T) {
	resourceName := randString10()
	themeKey := randString10()
	accessTokenKey := randString10()
	idTokenKey := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_tenant.test_%s", resourceName)

	// TODO(tenant_test): test property mutation across all fields
	startFromEmail, endFromEmail := "noreply@example.com", "no-reply@example.com"
	startMinimumPasswordAgeSeconds, endMinimumPasswordAgeSeconds := 10, 5
	startMinimumPasswordAgeEnabled, endMinimumPasswordAgeEnabled := true, false

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthTenantDestroy,
		Steps: []resource.TestStep{
			{
				// Test resource create
				Config: testAccTenantResourceBasicConfig(
					resourceName,
					themeKey,
					accessTokenKey,
					idTokenKey,
					startFromEmail,
					startMinimumPasswordAgeSeconds,
					startMinimumPasswordAgeEnabled,
					false,
				),
				Check: testTenantAccTestCheckFuncs(
					tfResourcePath,
					resourceName,
					startFromEmail,
					startMinimumPasswordAgeSeconds,
					startMinimumPasswordAgeEnabled,
					false,
				),
			},
			{
				// Test resource update/state mutate
				Config: testAccTenantResourceBasicConfig(
					resourceName,
					themeKey,
					accessTokenKey,
					idTokenKey,
					endFromEmail,
					endMinimumPasswordAgeSeconds,
					endMinimumPasswordAgeEnabled,
					true,
				),
				Check: testTenantAccTestCheckFuncs(
					tfResourcePath,
					resourceName,
					endFromEmail,
					endMinimumPasswordAgeSeconds,
					endMinimumPasswordAgeEnabled,
					true,
				),
			},
			{
				// Test importing resource into state
				ResourceName:            tfResourcePath,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// testTenantAccTestCheckFuncs abstracts the test case setup required between
// create and update testing.
func testTenantAccTestCheckFuncs(
	tfResourcePath string,
	resourceName string,
	fromEmail string,
	minimumPasswordAgeSeconds int,
	minimumPasswordAgeEnabled bool,
	genericConnectorIncluded bool,
) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccCheckFusionauthTenantExists(tfResourcePath),
		resource.TestCheckResourceAttrSet(tfResourcePath, "tenant_id"),

		// connector policies
		testAccCheckConnectorPolicies(tfResourcePath, genericConnectorIncluded),

		// user data
		resource.TestCheckResourceAttr(tfResourcePath, "data.user", "data"),
		resource.TestCheckResourceAttr(tfResourcePath, "data.lives", "here"),

		// email_configuration
		testAccCheckEmailConfigurationAdditionalHeaders(tfResourcePath),
		resource.TestCheckResourceAttr(tfResourcePath, "email_configuration.0.default_from_name", "noreply"),
		resource.TestCheckResourceAttr(tfResourcePath, "email_configuration.0.default_from_email", fromEmail),
		// resource.TestCheckResourceAttr(tfResourcePath, "email_configuration.0.forgot_password_email_template_id", "UUID"),
		resource.TestCheckResourceAttr(tfResourcePath, "email_configuration.0.host", "smtp.example.com"),
		resource.TestCheckResourceAttr(tfResourcePath, "email_configuration.0.password", "s3cureP@ssw0rd"),
		// resource.TestCheckResourceAttr(tfResourcePath, "email_configuration.0.passwordless_email_template_id", "UUID"),
		resource.TestCheckResourceAttr(tfResourcePath, "email_configuration.0.port", "587"),
		resource.TestCheckResourceAttr(tfResourcePath, "email_configuration.0.properties", "property=sold"),
		resource.TestCheckResourceAttr(tfResourcePath, "email_configuration.0.security", "TLS"),
		// resource.TestCheckResourceAttr(tfResourcePath, "email_configuration.0.set_password_email_template_id", "UUID"),
		resource.TestCheckResourceAttr(tfResourcePath, "email_configuration.0.unverified.0.allow_email_change_when_gated", "true"),
		resource.TestCheckResourceAttr(tfResourcePath, "email_configuration.0.unverified.0.behavior", "Gated"),

		// event_configuration
		resource.TestCheckResourceAttr(tfResourcePath, "event_configuration.#", "2"),
		resource.TestCheckTypeSetElemNestedAttrs(tfResourcePath, "event_configuration.*", map[string]string{
			"event":            "user.delete",
			"enabled":          "true",
			"transaction_type": "None",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(tfResourcePath, "event_configuration.*", map[string]string{
			"event":            "user.create",
			"enabled":          "true",
			"transaction_type": "SuperMajority",
		}),

		// external_identifier_configuration
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.authorization_grant_id_time_to_live_in_seconds", "30"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.change_password_id_generator.0.length", "32"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.change_password_id_generator.0.type", "randomBytes"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.change_password_id_time_to_live_in_seconds", "600"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.device_code_time_to_live_in_seconds", "1800"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.device_user_code_id_generator.0.length", "6"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.device_user_code_id_generator.0.type", "randomAlphaNumeric"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.email_verification_id_generator.0.length", "32"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.email_verification_id_generator.0.type", "randomBytes"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.email_verification_one_time_code_generator.0.length", "6"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.email_verification_one_time_code_generator.0.type", "randomAlphaNumeric"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.email_verification_id_time_to_live_in_seconds", "86400"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.external_authentication_id_time_to_live_in_seconds", "300"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.login_intent_time_to_live_in_seconds", "3600"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.one_time_password_time_to_live_in_seconds", "60"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.passwordless_login_generator.0.length", "32"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.passwordless_login_generator.0.type", "randomBytes"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.passwordless_login_time_to_live_in_seconds", "600"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.registration_verification_id_generator.0.length", "32"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.registration_verification_id_generator.0.type", "randomBytes"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.registration_verification_one_time_code_generator.0.length", "6"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.registration_verification_one_time_code_generator.0.type", "randomAlphaNumeric"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.registration_verification_id_time_to_live_in_seconds", "86400"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.saml_v2_authn_request_id_ttl_seconds", "300"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.setup_password_id_generator.0.length", "32"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.setup_password_id_generator.0.type", "randomBytes"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.setup_password_id_time_to_live_in_seconds", "86400"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.two_factor_one_time_code_id_generator.0.length", "8"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.two_factor_one_time_code_id_generator.0.type", "randomDigits"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.two_factor_id_time_to_live_in_seconds", "300"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.two_factor_one_time_code_id_time_to_live_in_seconds", "60"),
		resource.TestCheckResourceAttr(tfResourcePath, "external_identifier_configuration.0.two_factor_trust_id_time_to_live_in_seconds", "2592000"),

		// failed_authentication_configuration
		resource.TestCheckResourceAttr(tfResourcePath, "failed_authentication_configuration.0.action_duration", "1"),
		resource.TestCheckResourceAttr(tfResourcePath, "failed_authentication_configuration.0.action_duration_unit", "DAYS"),
		resource.TestCheckResourceAttr(tfResourcePath, "failed_authentication_configuration.0.reset_count_in_seconds", "600"),
		resource.TestCheckResourceAttr(tfResourcePath, "failed_authentication_configuration.0.too_many_attempts", "3"),
		resource.TestCheckResourceAttr(tfResourcePath, "failed_authentication_configuration.0.action_cancel_policy_on_password_reset", "true"),
		resource.TestCheckResourceAttr(tfResourcePath, "failed_authentication_configuration.0.email_user", "true"),
		// resource.TestCheckResourceAttr(tfResourcePath, "failed_authentication_configuration.0.user_action_id", "UUID"),

		// family_configuration
		resource.TestCheckResourceAttr(tfResourcePath, "family_configuration.0.allow_child_registrations", "false"),
		// resource.TestCheckResourceAttr(tfResourcePath, "family_configuration.0.#confirm_child_email_template_id", "UUID"),
		resource.TestCheckResourceAttr(tfResourcePath, "family_configuration.0.delete_orphaned_accounts", "true"),
		resource.TestCheckResourceAttr(tfResourcePath, "family_configuration.0.delete_orphaned_accounts_days", "60"),
		resource.TestCheckResourceAttr(tfResourcePath, "family_configuration.0.enabled", "true"),
		// resource.TestCheckResourceAttr(tfResourcePath, "family_configuration.0.#family_request_email_template_id", "UUID"),
		resource.TestCheckResourceAttr(tfResourcePath, "family_configuration.0.maximum_child_age", "14"),
		resource.TestCheckResourceAttr(tfResourcePath, "family_configuration.0.minimum_owner_age", "18"),
		resource.TestCheckResourceAttr(tfResourcePath, "family_configuration.0.parent_email_required", "false"),
		// resource.TestCheckResourceAttr(tfResourcePath, "family_configuration.0.#parent_registration_email_template_id", "UUID"),

		// form_configuration
		// resource.TestCheckResourceAttr(tfResourcePath, "form_configuration.0.admin_user_form_id", "UUID"),

		resource.TestCheckResourceAttr(tfResourcePath, "http_session_max_inactive_interval", "3400"),
		resource.TestCheckResourceAttr(tfResourcePath, "issuer", "https://example.com"),

		// jwt_configuration
		resource.TestCheckResourceAttrSet(tfResourcePath, "jwt_configuration.0.access_token_key_id"),
		resource.TestCheckResourceAttrSet(tfResourcePath, "jwt_configuration.0.id_token_key_id"),
		resource.TestCheckResourceAttr(tfResourcePath, "jwt_configuration.0.time_to_live_in_seconds", "3600"),
		resource.TestCheckResourceAttr(tfResourcePath, "jwt_configuration.0.refresh_token_sliding_window_maximum_time_to_live_in_minutes", "43200"),
		resource.TestCheckResourceAttr(tfResourcePath, "jwt_configuration.0.refresh_token_time_to_live_in_minutes", "43200"),
		resource.TestCheckResourceAttr(tfResourcePath, "jwt_configuration.0.refresh_token_expiration_policy", "SlidingWindowWithMaximumLifetime"),

		// login_configuration
		resource.TestCheckResourceAttr(tfResourcePath, "login_configuration.0.require_authentication", "true"),

		resource.TestCheckResourceAttr(tfResourcePath, "logout_url", "https://example.com/signed-out"),

		// maximum_password_age
		resource.TestCheckResourceAttr(tfResourcePath, "maximum_password_age.0.days", "90"),
		resource.TestCheckResourceAttr(tfResourcePath, "maximum_password_age.0.enabled", "true"),

		// minimum_password_age
		resource.TestCheckResourceAttr(tfResourcePath, "minimum_password_age.0.seconds", strconv.Itoa(minimumPasswordAgeSeconds)),
		resource.TestCheckResourceAttr(tfResourcePath, "minimum_password_age.0.enabled", strconv.FormatBool(minimumPasswordAgeEnabled)),

		// multi_factor_configuration
		resource.TestCheckResourceAttr(tfResourcePath, "multi_factor_configuration.0.login_policy", "Enabled"),
		resource.TestCheckResourceAttr(tfResourcePath, "multi_factor_configuration.0.authenticator.0.enabled", "true"),
		// resource.TestCheckResourceAttr(tfResourcePath, "multi_factor_configuration.0.authenticator.0.template_id", "UUID"),
		// requires paid edition of FusionAuth
		resource.TestCheckResourceAttr(tfResourcePath, "multi_factor_configuration.0.email.0.enabled", "false"),
		// resource.TestCheckResourceAttr(tfResourcePath, "multi_factor_configuration.0.email.0.template_id", "UUID"),
		// requires paid edition of FusionAuth
		resource.TestCheckResourceAttr(tfResourcePath, "multi_factor_configuration.0.sms.0.enabled", "false"),
		// resource.TestCheckResourceAttr(tfResourcePath, "multi_factor_configuration.0.sms.0.messenger_id", "UUID"),
		// resource.TestCheckResourceAttr(tfResourcePath, "multi_factor_configuration.0.sms.0.template_id", "UUID"),

		resource.TestCheckResourceAttr(tfResourcePath, "name", fmt.Sprintf("test-acc %s", resourceName)),

		// password_encryption_configuration
		resource.TestCheckResourceAttr(tfResourcePath, "password_encryption_configuration.0.encryption_scheme", "bcrypt"),
		resource.TestCheckResourceAttr(tfResourcePath, "password_encryption_configuration.0.encryption_scheme_factor", "14"),
		resource.TestCheckResourceAttr(tfResourcePath, "password_encryption_configuration.0.modify_encryption_scheme_on_login", "true"),

		// password_validation_rules
		resource.TestCheckResourceAttr(tfResourcePath, "password_validation_rules.0.breach_detection.0.enabled", "false"),
		resource.TestCheckResourceAttr(tfResourcePath, "password_validation_rules.0.breach_detection.0.match_mode", "Medium"),
		// resource.TestCheckResourceAttr(tfResourcePath, "password_validation_rules.0.breach_detection.0.notify_user_email_template_id", "UUID"),
		resource.TestCheckResourceAttr(tfResourcePath, "password_validation_rules.0.breach_detection.0.on_login", "NotifyUser"),
		resource.TestCheckResourceAttr(tfResourcePath, "password_validation_rules.0.max_length", "50"),
		resource.TestCheckResourceAttr(tfResourcePath, "password_validation_rules.0.min_length", "6"),
		resource.TestCheckResourceAttr(tfResourcePath, "password_validation_rules.0.remember_previous_passwords.0.count", "3"),
		resource.TestCheckResourceAttr(tfResourcePath, "password_validation_rules.0.remember_previous_passwords.0.enabled", "true"),
		resource.TestCheckResourceAttr(tfResourcePath, "password_validation_rules.0.required_mixed_case", "true"),
		resource.TestCheckResourceAttr(tfResourcePath, "password_validation_rules.0.require_non_alpha", "true"),
		resource.TestCheckResourceAttr(tfResourcePath, "password_validation_rules.0.require_number", "true"),
		resource.TestCheckResourceAttr(tfResourcePath, "password_validation_rules.0.validate_on_login", "true"),

		// rate_limit_configuration
		resource.TestCheckResourceAttrSet(tfResourcePath, "rate_limit_configuration.#"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.failed_login.0.enabled", "true"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.failed_login.0.limit", "6"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.failed_login.0.time_period_in_seconds", "60"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.forgot_password.0.enabled", "false"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.forgot_password.0.limit", "5"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.forgot_password.0.time_period_in_seconds", "59"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.send_email_verification.0.enabled", "true"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.send_email_verification.0.limit", "4"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.send_email_verification.0.time_period_in_seconds", "58"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.send_passwordless.0.enabled", "false"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.send_passwordless.0.limit", "3"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.send_passwordless.0.time_period_in_seconds", "57"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.send_registration_verification.0.enabled", "true"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.send_registration_verification.0.limit", "2"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.send_registration_verification.0.time_period_in_seconds", "56"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.send_two_factor.0.enabled", "false"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.send_two_factor.0.limit", "1"),
		resource.TestCheckResourceAttr(tfResourcePath, "rate_limit_configuration.0.send_two_factor.0.time_period_in_seconds", "55"),

		// registration_configuration
		resource.TestCheckResourceAttrSet(tfResourcePath, "registration_configuration.#"),
		resource.TestCheckResourceAttr(tfResourcePath, "registration_configuration.0.blocked_domains.0", "blocked-domain.com"),

		// captcha_configuration
		resource.TestCheckResourceAttrSet(tfResourcePath, "captcha_configuration.#"),
		resource.TestCheckResourceAttr(tfResourcePath, "captcha_configuration.0.enabled", "true"),
		resource.TestCheckResourceAttr(tfResourcePath, "captcha_configuration.0.captcha_method", "GoogleRecaptchaV3"),
		resource.TestCheckResourceAttr(tfResourcePath, "captcha_configuration.0.secret_key", "captcha_secret_key"),
		resource.TestCheckResourceAttr(tfResourcePath, "captcha_configuration.0.site_key", "captcha_site_key"),
		resource.TestCheckResourceAttr(tfResourcePath, "captcha_configuration.0.threshold", "0.5"),

		resource.TestCheckResourceAttrSet(tfResourcePath, "theme_id"),

		// user_delete_policy
		resource.TestCheckResourceAttr(tfResourcePath, "user_delete_policy.0.unverified_enabled", "true"),
		resource.TestCheckResourceAttr(tfResourcePath, "user_delete_policy.0.unverified_number_of_days_to_retain", "30"),

		// username_configuration
		resource.TestCheckResourceAttr(tfResourcePath, "username_configuration.0.unique.0.enabled", "false"),
		resource.TestCheckResourceAttr(tfResourcePath, "username_configuration.0.unique.0.number_of_digits", "8"),
		resource.TestCheckResourceAttr(tfResourcePath, "username_configuration.0.unique.0.separator", "_"),
		resource.TestCheckResourceAttr(tfResourcePath, "username_configuration.0.unique.0.strategy", "Always"),
	)
}

func testAccCheckConnectorPolicies(tfResourcePath string, genericConnectorIncluded bool) resource.TestCheckFunc {
	if genericConnectorIncluded {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(tfResourcePath, "connector_policy.#", "2"),
			resource.TestCheckResourceAttr(tfResourcePath, "connector_policy.1.connector_id", "e3306678-a53a-4964-9040-1c96f36dda72"),
		)
	}
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(tfResourcePath, "connector_policy.#", "1"),
		resource.TestCheckResourceAttr(tfResourcePath, "connector_policy.0.connector_id", "e3306678-a53a-4964-9040-1c96f36dda72"),
	)
}

func testAccCheckEmailConfigurationAdditionalHeaders(tfResourcePath string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(tfResourcePath, "email_configuration.0.additional_headers.%", "2"),
		resource.TestCheckResourceAttr(tfResourcePath, "email_configuration.0.additional_headers.HeaderName1", "HeaderValue1"),
		resource.TestCheckResourceAttr(tfResourcePath, "email_configuration.0.additional_headers.HeaderName2", "HeaderValue2"),
	)
}

func testAccCheckFusionauthTenantExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		faClient := fusionauthClient()
		tenant, faErrs, err := faClient.RetrieveTenant(rs.Primary.ID)
		if errs := checkFusionauthErrors(faErrs, err); errs != nil {
			return err
		}

		if tenant == nil || tenant.StatusCode != http.StatusOK {
			// This is a weird edge case...
			return fmt.Errorf("failed to get resource: %#+v", tenant)
		}

		return nil
	}
}

func testAccCheckFusionauthTenantDestroy(s *terraform.State) error {
	faClient := fusionauthClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_tenant" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			tenant, faErrs, err := faClient.RetrieveTenant(rs.Primary.ID)
			if errs := checkFusionauthRetryErrors(faErrs, err); errs != nil {
				return errs
			}

			if tenant != nil && tenant.StatusCode == http.StatusNotFound {
				// resource destroyed!
				return nil
			}

			return retry.RetryableError(fmt.Errorf("fusionauth resource still exists: %s", rs.Primary.ID))
		})

		if err != nil {
			// We failed destroying the resource...
			return err
		}
	}

	return nil
}

// testAccTenantResourceBasicConfig generates the terraform config for the
// resources a tenant requires.
func testAccTenantResourceBasicConfig(
	resourceName string,
	themeKey string,
	accessTokenKey string,
	idTokenKey string,
	fromEmail string,
	minimumPasswordAgeSeconds int,
	minimumPasswordAgeEnabled bool,
	genericConnectorIncluded bool,
) string {
	return testAccKeyResourceConfig(
		"",
		accessTokenKey,
		fusionauth.Algorithm_RS256,
		2048,
	) +
		testAccKeyResourceConfig(
			"",
			idTokenKey,
			fusionauth.Algorithm_RS256,
			2048,
		) +
		testAccThemeResourceConfig(
			themeKey,
			testdata.MessageProperties(""),
			"/* stylez */",
			generateFusionAuthTemplate(),
		) +
		testAccGenericConnectorBasicConfig(
			resourceName,
			"http://example.com/connector",
			"1000",
			"password1",
			"username",
			"test-generic-connector"+resourceName,
			"1000",
		) +
		testAccTenantResourceConfig(
			resourceName,
			themeKey,
			fromEmail,
			minimumPasswordAgeSeconds,
			minimumPasswordAgeEnabled,
			genericConnectorIncluded,
		)
}

// testAccTenantResource returns terraform configuration to create a test
// tenant.
//
// Note:
//   - A bug in the terraform SDK means defaults configured for TypeList/TypeSet
//     schemas aren't applied unless the top level object is defined in the
//     config, for example, you have to explicitly add `minimum_password_age {}`
//     to get the defaults to propagate down into the object's properties.
//     Refer: https://github.com/hashicorp/terraform-plugin-sdk/issues/142
//   - `form_configuration.admin_user_form_id` is commented out as it requires a
//     paid edition of fusionauth.
//   - `multi_factor_configuration.email.enabled` is set to false, as it requires
//     a paid edition of fusionauth.
//   - `multi_factor_configuration.sms.enabled` is set to false, as it requires
//     a paid edition of fusionauth.
//   - `password_validation_rules.breach_detection.enabled` is set to false, as
//     it requires a paid edition of fusionauth.
//   - `username_configuration.unique.enabled` is set to false, as it requires a
//     paid edition of fusionauth.
func testAccTenantResourceConfig(
	resourceName string,
	themeKey string,
	fromEmail string,
	minimumPasswordAgeSeconds int,
	minimumPasswordAgeEnabled bool,
	genericConnectorIncluded bool,
) string {
	if themeKey != "" {
		themeKey = fmt.Sprintf(
			"\n  theme_id = fusionauth_theme.test_%s.id\n",
			themeKey,
		)
	}
	connectorPolicies := ""
	if genericConnectorIncluded {
		connectorPolicies = fmt.Sprintf(`
		  connector_policy {
				connector_id = fusionauth_generic_connector.test_%s.id
				domains      = ["*"]
				migrate      = true
			}
		  connector_policy {
				connector_id = "e3306678-a53a-4964-9040-1c96f36dda72"
				domains      = ["*"]
				migrate      = false
			}
		`, resourceName)
	}

	return fmt.Sprintf(`
# Tenant Setup
resource "fusionauth_tenant" "test_%[1]s" {
  #source_tenant_id = "UUID"
  #tenant_id        = "UUID"
  # connector policies %[6]s
  data = {
    user  = "data"
    lives = "here"
  }
  email_configuration {
    default_from_name  = "noreply"
    additional_headers = {
      "HeaderName1" = "HeaderValue1"
      "HeaderName2" = "HeaderValue2"
    }
    default_from_email = "%[3]s"
    #forgot_password_email_template_id = ""
    host               = "smtp.example.com"
    password           = "s3cureP@ssw0rd"
    #passwordless_email_template_id = ""
    port               = 587
    properties         = "property=sold"
    security           = "TLS"
    # set_password_email_template_id = ""
    unverified {
      allow_email_change_when_gated = true
      behavior                      = "Gated"
    }
  }
  event_configuration {
    event            = "user.delete"
    enabled          = true
    transaction_type = "None"
  }
  event_configuration {
    event            = "user.create"
    enabled          = true
    transaction_type = "SuperMajority"
  }
  external_identifier_configuration {
    authorization_grant_id_time_to_live_in_seconds = 30
    change_password_id_generator {
      length = 32
      type   = "randomBytes"
    }
    change_password_id_time_to_live_in_seconds = 600
    device_code_time_to_live_in_seconds        = 1800
    device_user_code_id_generator {
      length = 6
      type   = "randomAlphaNumeric"
    }
    email_verification_id_generator {
      length = 32
      type   = "randomBytes"
    }
    email_verification_one_time_code_generator {   # Technically Optional
      length = 6
      type   = "randomAlphaNumeric"
    }
    email_verification_id_time_to_live_in_seconds      = 86400
    external_authentication_id_time_to_live_in_seconds = 300
    login_intent_time_to_live_in_seconds               = 3600
    one_time_password_time_to_live_in_seconds          = 60
    passwordless_login_generator {
      length = 32
      type   = "randomBytes"
    }
    passwordless_login_time_to_live_in_seconds = 600
    registration_verification_id_generator {
      length = 32
      type   = "randomBytes"
    }
    registration_verification_one_time_code_generator { # Optional since v1.28.0
      length = 6
      type   = "randomAlphaNumeric"
    }
    registration_verification_id_time_to_live_in_seconds = 86400
    saml_v2_authn_request_id_ttl_seconds = 300
    setup_password_id_generator {
      length = 32
      type   = "randomBytes"
    }
    setup_password_id_time_to_live_in_seconds   = 86400
    two_factor_one_time_code_id_generator { # Optional since v1.28.0
      length = 8
      type   = "randomDigits"
    }
    two_factor_id_time_to_live_in_seconds               = 300
    two_factor_one_time_code_id_time_to_live_in_seconds = 60
    two_factor_trust_id_time_to_live_in_seconds         = 2592000
  }
  failed_authentication_configuration {
    action_duration        = 1
    action_duration_unit   = "DAYS"
    reset_count_in_seconds = 600
    too_many_attempts      = 3
    action_cancel_policy_on_password_reset = true
    email_user             = true
    #user_action_id         = "UUID"
  }
  family_configuration {
    allow_child_registrations             = false
    #confirm_child_email_template_id       = "UUID"
    delete_orphaned_accounts              = true
    delete_orphaned_accounts_days         = 60
    enabled                               = true
    #family_request_email_template_id      = "UUID"
    maximum_child_age                     = 14
    minimum_owner_age                     = 18
    parent_email_required                 = false
    #parent_registration_email_template_id = "UUID"
  }
  form_configuration {
    # requires paid edition of FusionAuth
    #admin_user_form_id = "UUID"
  }
  http_session_max_inactive_interval = 3400
  issuer   = "https://example.com"
  jwt_configuration {
    refresh_token_expiration_policy = "SlidingWindowWithMaximumLifetime"
    refresh_token_time_to_live_in_minutes = 43200
    time_to_live_in_seconds               = 3600
    refresh_token_sliding_window_maximum_time_to_live_in_minutes = 43200
  }
  login_configuration {
    require_authentication = true
  }
  logout_url = "https://example.com/signed-out"
  maximum_password_age {
    days    = 90
    enabled = true
  }
  minimum_password_age {
    seconds = %[4]d
    enabled = %[5]t
  }
  multi_factor_configuration {
    login_policy = "Enabled"
    authenticator {
      enabled     = true
      #template_id = "UUID"
    }
    email {
      # requires paid edition of FusionAuth
      enabled     = false
      #template_id = "UUID"
    }
    sms {
      # requires paid edition of FusionAuth
      enabled      = false
      #messenger_id = "UUID"
      #template_id  = "UUID"
    }
  }
  name = "test-acc %[1]s"
  password_encryption_configuration {
    encryption_scheme                 = "bcrypt"
    encryption_scheme_factor          = 14
    modify_encryption_scheme_on_login = true
  }
  password_validation_rules {
    breach_detection {
      # requires paid edition of FusionAuth
      enabled                       = false
      match_mode                    = "Medium"
      #notify_user_email_template_id = "UUID"
      on_login                      = "NotifyUser"
    }
    max_length = 50
    min_length = 6
    remember_previous_passwords {
      count   = 3
      enabled = true
    }
    required_mixed_case = true
    require_non_alpha   = true
    require_number      = true
    validate_on_login   = true
  }
  rate_limit_configuration {
    failed_login {
      enabled                = true
      limit                  = 6
      time_period_in_seconds = 60
    }
    forgot_password {
      enabled                = false
      limit                  = 5
      time_period_in_seconds = 59
    }
    send_email_verification {
      enabled                = true
      limit                  = 4
      time_period_in_seconds = 58
    }
    send_passwordless {
      enabled                = false
      limit                  = 3
      time_period_in_seconds = 57
    }
    send_registration_verification {
      enabled                = true
      limit                  = 2
      time_period_in_seconds = 56
    }
    send_two_factor {
      enabled                = false
      limit                  = 1
      time_period_in_seconds = 55
    }
  }
  captcha_configuration {
    enabled    		= true
    captcha_method  = "GoogleRecaptchaV3"
    site_key   		= "captcha_site_key"
    secret_key 		= "captcha_secret_key"
    threshold  		= 0.5
  }
  registration_configuration {
	blocked_domains = ["blocked-domain.com"]
  }
  # theme_id%[2]s
  user_delete_policy {
    unverified_enabled                  = true
    unverified_number_of_days_to_retain = 30
  }
  username_configuration {
    unique {
      # requires paid edition of FusionAuth
      enabled          = false
      number_of_digits = 8
      separator        = "_"
      strategy         = "Always"
    }
  }
}
`,
		resourceName,
		themeKey,
		fromEmail,
		minimumPasswordAgeSeconds,
		minimumPasswordAgeEnabled,
		connectorPolicies,
	)
}
