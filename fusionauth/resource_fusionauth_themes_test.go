package fusionauth

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gpsinsight/terraform-provider-fusionauth/fusionauth/testdata"
)

func TestAccFusionauthTheme_basic(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_theme.test_%s", resourceName)

	startMessages, endMessages := testdata.MessageProperties(""), testdata.MessageProperties("Terraform")
	startStyles, endStyles := "/* styles */", "/* changed styles */"
	startTemplates, endTemplates := generateFusionAuthTemplate(), generateFusionAuthTemplate()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthThemeDestroy,
		Steps: []resource.TestStep{
			{
				// Test resource create
				Config: testAccThemeResourceConfig(resourceName, startMessages, startStyles, startTemplates),
				Check: testThemeAccTestCheckFuncs(
					tfResourcePath,
					resourceName,
					startMessages,
					startStyles,
					startTemplates,
				),
			},
			{
				// Test resource update/state mutate
				Config: testAccThemeResourceConfig(resourceName, endMessages, endStyles, endTemplates),
				Check: testThemeAccTestCheckFuncs(
					tfResourcePath,
					resourceName,
					endMessages,
					endStyles,
					endTemplates,
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

// testThemeAccTestCheckFuncs abstracts the test case setup required between
// create and update testing.
func testThemeAccTestCheckFuncs(
	tfResourcePath string,
	resourceName string,
	defaultMessages string,
	stylesheet string,
	templates fusionauth.Templates,
) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccCheckFusionauthThemeExists(tfResourcePath),
		resource.TestCheckResourceAttr(tfResourcePath, "default_messages", defaultMessages),
		resource.TestCheckResourceAttr(tfResourcePath, "name", fmt.Sprintf("test-acc %s", resourceName)),
		resource.TestCheckResourceAttr(tfResourcePath, "stylesheet", stylesheet),
		resource.TestCheckResourceAttr(tfResourcePath, "account_edit", templates.AccountEdit),
		resource.TestCheckResourceAttr(tfResourcePath, "account_index", templates.AccountIndex),
		resource.TestCheckResourceAttr(tfResourcePath, "account_two_factor_disable", templates.AccountTwoFactorDisable),
		resource.TestCheckResourceAttr(tfResourcePath, "account_two_factor_enable", templates.AccountTwoFactorEnable),
		resource.TestCheckResourceAttr(tfResourcePath, "account_two_factor_index", templates.AccountTwoFactorIndex),
		resource.TestCheckResourceAttr(tfResourcePath, "account_webauthn_add", templates.AccountWebAuthnAdd),
		resource.TestCheckResourceAttr(tfResourcePath, "account_webauthn_delete", templates.AccountWebAuthnDelete),
		resource.TestCheckResourceAttr(tfResourcePath, "account_webauthn_index", templates.AccountWebAuthnIndex),
		resource.TestCheckResourceAttr(tfResourcePath, "confirmation_required", templates.ConfirmationRequired),
		resource.TestCheckResourceAttr(tfResourcePath, "email_complete", templates.EmailComplete),
		// resource.TestCheckResourceAttr(tfResourcePath, "email_send", startTemplates.EmailSend), // DEPRECATED
		resource.TestCheckResourceAttr(tfResourcePath, "email_sent", templates.EmailSent),
		resource.TestCheckResourceAttr(tfResourcePath, "email_verification_required", templates.EmailVerificationRequired),
		resource.TestCheckResourceAttr(tfResourcePath, "email_verify", templates.EmailVerify),
		resource.TestCheckResourceAttr(tfResourcePath, "helpers", templates.Helpers),
		resource.TestCheckResourceAttr(tfResourcePath, "index", templates.Index),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_authorize", templates.Oauth2Authorize),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_authorized_not_registered", templates.Oauth2AuthorizedNotRegistered),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_child_registration_not_allowed", templates.Oauth2ChildRegistrationNotAllowed),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_child_registration_not_allowed_complete", templates.Oauth2ChildRegistrationNotAllowedComplete),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_complete_registration", templates.Oauth2CompleteRegistration),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_consent", templates.Oauth2Consent),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_device", templates.Oauth2Device),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_device_complete", templates.Oauth2DeviceComplete),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_error", templates.Oauth2Error),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_logout", templates.Oauth2Logout),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_passwordless", templates.Oauth2Passwordless),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_register", templates.Oauth2Register),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_start_idp_link", templates.Oauth2StartIdPLink),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_two_factor", templates.Oauth2TwoFactor),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_two_factor_methods", templates.Oauth2TwoFactorMethods),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_two_factor_enable", templates.Oauth2TwoFactorEnable),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_two_factor_enable_complete", templates.Oauth2TwoFactorEnableComplete),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_wait", templates.Oauth2Wait),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_webauthn", templates.Oauth2WebAuthn),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_webauthn_reauth", templates.Oauth2WebAuthnReauth),
		resource.TestCheckResourceAttr(tfResourcePath, "oauth2_webauthn_reauth_enable", templates.Oauth2WebAuthnReauthEnable),
		resource.TestCheckResourceAttr(tfResourcePath, "password_change", templates.PasswordChange),
		resource.TestCheckResourceAttr(tfResourcePath, "password_complete", templates.PasswordComplete),
		resource.TestCheckResourceAttr(tfResourcePath, "password_forgot", templates.PasswordForgot),
		resource.TestCheckResourceAttr(tfResourcePath, "password_sent", templates.PasswordSent),
		resource.TestCheckResourceAttr(tfResourcePath, "registration_complete", templates.RegistrationComplete),
		// resource.TestCheckResourceAttr(tfResourcePath, "registration_send", templates.RegistrationSend), // DEPRECATED
		resource.TestCheckResourceAttr(tfResourcePath, "registration_sent", templates.RegistrationSent),
		resource.TestCheckResourceAttr(tfResourcePath, "registration_verification_required", templates.RegistrationVerificationRequired),
		resource.TestCheckResourceAttr(tfResourcePath, "registration_verify", templates.RegistrationVerify),
		resource.TestCheckResourceAttr(tfResourcePath, "samlv2_logout", templates.Samlv2Logout),
		resource.TestCheckResourceAttr(tfResourcePath, "unauthorized", templates.Unauthorized),
	)
}

func testAccCheckFusionauthThemeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		faClient := fusionauthClient()
		theme, faErrs, err := faClient.RetrieveTheme(rs.Primary.ID)
		if errs := checkFusionauthErrors(faErrs, err); errs != nil {
			return err
		}

		if theme == nil || theme.StatusCode != http.StatusOK {
			// This is a weird edge case...
			return fmt.Errorf("failed to get resource: %#+v", theme)
		}

		return nil
	}
}

func testAccCheckFusionauthThemeDestroy(s *terraform.State) error {
	faClient := fusionauthClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_theme" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			theme, faErrs, err := faClient.RetrieveTheme(rs.Primary.ID)
			if errs := checkFusionauthRetryErrors(faErrs, err); errs != nil {
				return errs
			}

			if theme != nil && theme.StatusCode == http.StatusNotFound {
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

// generateFusionAuthTemplate generates random template data to ensure each property is being set correctly.
func generateFusionAuthTemplate() fusionauth.Templates {
	return fusionauth.Templates{
		AccountEdit:                       randString20(),
		AccountIndex:                      randString20(),
		AccountTwoFactorDisable:           randString20(),
		AccountTwoFactorEnable:            randString20(),
		AccountTwoFactorIndex:             randString20(),
		AccountWebAuthnAdd:                randString20(),
		AccountWebAuthnDelete:             randString20(),
		AccountWebAuthnIndex:              randString20(),
		ConfirmationRequired:              randString20(),
		EmailComplete:                     randString20(),
		EmailSent:                         randString20(),
		EmailVerificationRequired:         randString20(),
		EmailVerify:                       randString20(),
		Helpers:                           randString20(),
		Index:                             randString20(),
		Oauth2Authorize:                   randString20(),
		Oauth2AuthorizedNotRegistered:     randString20(),
		Oauth2ChildRegistrationNotAllowed: randString20(),
		Oauth2ChildRegistrationNotAllowedComplete: randString20(),
		Oauth2CompleteRegistration:                randString20(),
		Oauth2Consent:                             randString20(),
		Oauth2Device:                              randString20(),
		Oauth2DeviceComplete:                      randString20(),
		Oauth2Error:                               randString20(),
		Oauth2Logout:                              randString20(),
		Oauth2Passwordless:                        randString20(),
		Oauth2Register:                            randString20(),
		Oauth2StartIdPLink:                        randString20(),
		Oauth2TwoFactor:                           randString20(),
		Oauth2TwoFactorEnable:                     randString20(),
		Oauth2TwoFactorEnableComplete:             randString20(),
		Oauth2TwoFactorMethods:                    randString20(),
		Oauth2Wait:                                randString20(),
		Oauth2WebAuthn:                            randString20(),
		Oauth2WebAuthnReauth:                      randString20(),
		Oauth2WebAuthnReauthEnable:                randString20(),
		PasswordChange:                            randString20(),
		PasswordComplete:                          randString20(),
		PasswordForgot:                            randString20(),
		PasswordSent:                              randString20(),
		RegistrationComplete:                      randString20(),
		RegistrationSent:                          randString20(),
		RegistrationVerificationRequired:          randString20(),
		RegistrationVerify:                        randString20(),
		Samlv2Logout:                              randString20(),
		Unauthorized:                              randString20(),

		// TODO(themes): test for deprecated properties
		// EmailSend:        randString20(),
		// RegistrationSend: randString20(),
	}
}

// testAccThemeResourceConfig returns terraform configuration to generate a test
// theme.
func testAccThemeResourceConfig(
	resourceName string,
	defaultMessages string,
	stylesheet string,
	templates fusionauth.Templates,
) string {
	return fmt.Sprintf(`
# Theme Setup
resource "fusionauth_theme" "test_%[1]s" {
  name                                           = "test-acc %[1]s"
  default_messages                               = <<EOF
%[2]sEOF
  stylesheet                                     = "%[3]s"
  account_edit                                   = "%[4]s"
  account_index                                  = "%[5]s"
  account_two_factor_disable                     = "%[6]s"
  account_two_factor_enable                      = "%[7]s"
  account_two_factor_index                       = "%[8]s"
  account_webauthn_add                           = "%[9]s"
  account_webauthn_delete                        = "%[10]s"
  account_webauthn_index                         = "%[11]s"
  confirmation_required                          = "%[12]s"
  email_complete                                 = "%[13]s"
  email_sent                                     = "%[14]s"
  email_verification_required                    = "%[15]s"
  email_verify                                   = "%[16]s"
  helpers                                        = "%[17]s"
  index                                          = "%[18]s"
  oauth2_authorize                               = "%[19]s"
  oauth2_authorized_not_registered               = "%[20]s"
  oauth2_child_registration_not_allowed          = "%[21]s"
  oauth2_child_registration_not_allowed_complete = "%[22]s"
  oauth2_complete_registration                   = "%[23]s"
  oauth2_consent                                 = "%[24]s"
  oauth2_device                                  = "%[25]s"
  oauth2_device_complete                         = "%[26]s"
  oauth2_error                                   = "%[27]s"
  oauth2_logout                                  = "%[28]s"
  oauth2_passwordless                            = "%[29]s"
  oauth2_register                                = "%[30]s"
  oauth2_start_idp_link                          = "%[31]s"
  oauth2_two_factor                              = "%[32]s"
  oauth2_two_factor_methods                      = "%[33]s"
  oauth2_two_factor_enable                       = "%[34]s"
  oauth2_two_factor_enable_complete              = "%[35]s"
  oauth2_wait                                    = "%[36]s"
  oauth2_webauthn                                = "%[37]s"
  oauth2_webauthn_reauth                         = "%[38]s"
  oauth2_webauthn_reauth_enable                  = "%[39]s"
  password_change                                = "%[40]s"
  password_complete                              = "%[41]s"
  password_forgot                                = "%[42]s"
  password_sent                                  = "%[43]s"
  registration_complete                          = "%[44]s"
  registration_sent                              = "%[45]s"
  registration_verification_required             = "%[46]s"
  registration_verify                            = "%[47]s"
  samlv2_logout                                  = "%[48]s"
  unauthorized                                   = "%[49]s"

  # Deprecated Properties
  email_send                                     = "%[50]s"
  registration_send                              = "%[51]s"
}
`,
		resourceName,
		defaultMessages,
		stylesheet,
		templates.AccountEdit,
		templates.AccountIndex,
		templates.AccountTwoFactorDisable,
		templates.AccountTwoFactorEnable,
		templates.AccountTwoFactorIndex,
		templates.AccountWebAuthnAdd,
		templates.AccountWebAuthnDelete,
		templates.AccountWebAuthnIndex,
		templates.ConfirmationRequired,
		templates.EmailComplete,
		templates.EmailSent,
		templates.EmailVerificationRequired,
		templates.EmailVerify,
		templates.Helpers,
		templates.Index,
		templates.Oauth2Authorize,
		templates.Oauth2AuthorizedNotRegistered,
		templates.Oauth2ChildRegistrationNotAllowed,
		templates.Oauth2ChildRegistrationNotAllowedComplete,
		templates.Oauth2CompleteRegistration,
		templates.Oauth2Consent,
		templates.Oauth2Device,
		templates.Oauth2DeviceComplete,
		templates.Oauth2Error,
		templates.Oauth2Logout,
		templates.Oauth2Passwordless,
		templates.Oauth2Register,
		templates.Oauth2StartIdPLink,
		templates.Oauth2TwoFactor,
		templates.Oauth2TwoFactorMethods,
		templates.Oauth2TwoFactorEnable,
		templates.Oauth2TwoFactorEnableComplete,
		templates.Oauth2Wait,
		templates.Oauth2WebAuthn,
		templates.Oauth2WebAuthnReauth,
		templates.Oauth2WebAuthnReauthEnable,
		templates.PasswordChange,
		templates.PasswordComplete,
		templates.PasswordForgot,
		templates.PasswordSent,
		templates.RegistrationComplete,
		templates.RegistrationSent,
		templates.RegistrationVerificationRequired,
		templates.RegistrationVerify,
		templates.Samlv2Logout,
		templates.Unauthorized,
		templates.EmailSend,
		templates.RegistrationSend,
	)
}
