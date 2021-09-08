package fusionauth

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gpsinsight/terraform-provider-fusionauth/fusionauth/testdata"
)

func TestAccFusionauthTheme_basic(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthThemeExists(tfResourcePath),
					resource.TestCheckResourceAttr(tfResourcePath, "default_messages", startMessages),
					resource.TestCheckResourceAttr(tfResourcePath, "name", fmt.Sprintf("test %s", resourceName)),
					resource.TestCheckResourceAttr(tfResourcePath, "stylesheet", startStyles),
					resource.TestCheckResourceAttr(tfResourcePath, "account_edit", startTemplates.AccountEdit),
					resource.TestCheckResourceAttr(tfResourcePath, "account_index", startTemplates.AccountIndex),
					resource.TestCheckResourceAttr(tfResourcePath, "account_two_factor_disable", startTemplates.AccountTwoFactorDisable),
					resource.TestCheckResourceAttr(tfResourcePath, "account_two_factor_enable", startTemplates.AccountTwoFactorEnable),
					resource.TestCheckResourceAttr(tfResourcePath, "account_two_factor_index", startTemplates.AccountTwoFactorIndex),
					resource.TestCheckResourceAttr(tfResourcePath, "email_complete", startTemplates.EmailComplete),
					//resource.TestCheckResourceAttr(tfResourcePath, "email_send", startTemplates.EmailSend), // DEPRECATED
					resource.TestCheckResourceAttr(tfResourcePath, "email_sent", startTemplates.EmailSent),
					resource.TestCheckResourceAttr(tfResourcePath, "email_verification_required", startTemplates.EmailVerificationRequired),
					resource.TestCheckResourceAttr(tfResourcePath, "email_verify", startTemplates.EmailVerify),
					resource.TestCheckResourceAttr(tfResourcePath, "helpers", startTemplates.Helpers),
					resource.TestCheckResourceAttr(tfResourcePath, "index", startTemplates.Index),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_authorize", startTemplates.Oauth2Authorize),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_authorized_not_registered", startTemplates.Oauth2AuthorizedNotRegistered),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_child_registration_not_allowed", startTemplates.Oauth2ChildRegistrationNotAllowed),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_child_registration_not_allowed_complete", startTemplates.Oauth2ChildRegistrationNotAllowedComplete),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_complete_registration", startTemplates.Oauth2CompleteRegistration),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_device", startTemplates.Oauth2Device),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_device_complete", startTemplates.Oauth2DeviceComplete),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_error", startTemplates.Oauth2Error),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_logout", startTemplates.Oauth2Logout),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_passwordless", startTemplates.Oauth2Passwordless),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_register", startTemplates.Oauth2Register),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_start_idp_link", startTemplates.Oauth2StartIdPLink),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_two_factor", startTemplates.Oauth2TwoFactor),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_two_factor_methods", startTemplates.Oauth2TwoFactorMethods),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_wait", startTemplates.Oauth2Wait),
					resource.TestCheckResourceAttr(tfResourcePath, "password_change", startTemplates.PasswordChange),
					resource.TestCheckResourceAttr(tfResourcePath, "password_complete", startTemplates.PasswordComplete),
					resource.TestCheckResourceAttr(tfResourcePath, "password_forgot", startTemplates.PasswordForgot),
					resource.TestCheckResourceAttr(tfResourcePath, "password_sent", startTemplates.PasswordSent),
					resource.TestCheckResourceAttr(tfResourcePath, "registration_complete", startTemplates.RegistrationComplete),
					//resource.TestCheckResourceAttr(tfResourcePath, "registration_send", startTemplates.RegistrationSend), // DEPRECATED
					resource.TestCheckResourceAttr(tfResourcePath, "registration_sent", startTemplates.RegistrationSent),
					resource.TestCheckResourceAttr(tfResourcePath, "registration_verification_required", startTemplates.RegistrationVerificationRequired),
					resource.TestCheckResourceAttr(tfResourcePath, "registration_verify", startTemplates.RegistrationVerify),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_logout", startTemplates.Samlv2Logout),
					resource.TestCheckResourceAttr(tfResourcePath, "unauthorized", startTemplates.Unauthorized),
				),
			},
			{
				// Test resource update/state mutate
				Config: testAccThemeResourceConfig(resourceName, endMessages, endStyles, endTemplates),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthThemeExists(tfResourcePath),
					resource.TestCheckResourceAttr(tfResourcePath, "default_messages", endMessages),
					resource.TestCheckResourceAttr(tfResourcePath, "name", fmt.Sprintf("test %s", resourceName)),
					resource.TestCheckResourceAttr(tfResourcePath, "stylesheet", endStyles),
					resource.TestCheckResourceAttr(tfResourcePath, "account_edit", endTemplates.AccountEdit),
					resource.TestCheckResourceAttr(tfResourcePath, "account_index", endTemplates.AccountIndex),
					resource.TestCheckResourceAttr(tfResourcePath, "account_two_factor_disable", endTemplates.AccountTwoFactorDisable),
					resource.TestCheckResourceAttr(tfResourcePath, "account_two_factor_enable", endTemplates.AccountTwoFactorEnable),
					resource.TestCheckResourceAttr(tfResourcePath, "account_two_factor_index", endTemplates.AccountTwoFactorIndex),
					resource.TestCheckResourceAttr(tfResourcePath, "email_complete", endTemplates.EmailComplete),
					//resource.TestCheckResourceAttr(tfResourcePath, "email_send", startTemplates.EmailSend), // DEPRECATED
					resource.TestCheckResourceAttr(tfResourcePath, "email_sent", endTemplates.EmailSent),
					resource.TestCheckResourceAttr(tfResourcePath, "email_verification_required", endTemplates.EmailVerificationRequired),
					resource.TestCheckResourceAttr(tfResourcePath, "email_verify", endTemplates.EmailVerify),
					resource.TestCheckResourceAttr(tfResourcePath, "helpers", endTemplates.Helpers),
					resource.TestCheckResourceAttr(tfResourcePath, "index", endTemplates.Index),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_authorize", endTemplates.Oauth2Authorize),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_authorized_not_registered", endTemplates.Oauth2AuthorizedNotRegistered),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_child_registration_not_allowed", endTemplates.Oauth2ChildRegistrationNotAllowed),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_child_registration_not_allowed_complete", endTemplates.Oauth2ChildRegistrationNotAllowedComplete),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_complete_registration", endTemplates.Oauth2CompleteRegistration),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_device", endTemplates.Oauth2Device),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_device_complete", endTemplates.Oauth2DeviceComplete),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_error", endTemplates.Oauth2Error),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_logout", endTemplates.Oauth2Logout),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_passwordless", endTemplates.Oauth2Passwordless),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_register", endTemplates.Oauth2Register),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_start_idp_link", endTemplates.Oauth2StartIdPLink),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_two_factor", endTemplates.Oauth2TwoFactor),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_two_factor_methods", endTemplates.Oauth2TwoFactorMethods),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_wait", endTemplates.Oauth2Wait),
					resource.TestCheckResourceAttr(tfResourcePath, "password_change", endTemplates.PasswordChange),
					resource.TestCheckResourceAttr(tfResourcePath, "password_complete", endTemplates.PasswordComplete),
					resource.TestCheckResourceAttr(tfResourcePath, "password_forgot", endTemplates.PasswordForgot),
					resource.TestCheckResourceAttr(tfResourcePath, "password_sent", endTemplates.PasswordSent),
					resource.TestCheckResourceAttr(tfResourcePath, "registration_complete", endTemplates.RegistrationComplete),
					//resource.TestCheckResourceAttr(tfResourcePath, "registration_send", endTemplates.RegistrationSend), // DEPRECATED
					resource.TestCheckResourceAttr(tfResourcePath, "registration_sent", endTemplates.RegistrationSent),
					resource.TestCheckResourceAttr(tfResourcePath, "registration_verification_required", endTemplates.RegistrationVerificationRequired),
					resource.TestCheckResourceAttr(tfResourcePath, "registration_verify", endTemplates.RegistrationVerify),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_logout", endTemplates.Samlv2Logout),
					resource.TestCheckResourceAttr(tfResourcePath, "unauthorized", endTemplates.Unauthorized),
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

func testAccCheckFusionauthThemeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		theme, faErrs, err := fusionauthClient().RetrieveTheme(rs.Primary.ID)
		if errs := checkFusionauthErrors(faErrs, err); errs != nil {
			return err
		}

		if theme == nil || theme.StatusCode != http.StatusOK {
			// This is a weird edge case...
			return fmt.Errorf("failed to get resource: %#+v\n", theme)
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
		err := resource.RetryContext(context.Background(), retryTimeout, func() *resource.RetryError {
			theme, faErrs, err := faClient.RetrieveTheme(rs.Primary.ID)
			if errs := checkFusionauthRetryErrors(faErrs, err); errs != nil {
				return errs
			}

			if theme != nil && theme.StatusCode == http.StatusNotFound {
				// user destroyed!
				return nil
			}

			return resource.RetryableError(fmt.Errorf("fusionauth resource still exists: %s", rs.Primary.ID))
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
		AccountEdit:                       acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		AccountIndex:                      acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		AccountTwoFactorDisable:           acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		AccountTwoFactorEnable:            acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		AccountTwoFactorIndex:             acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		EmailComplete:                     acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		EmailSent:                         acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		EmailVerificationRequired:         acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		EmailVerify:                       acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Helpers:                           acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Index:                             acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Oauth2Authorize:                   acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Oauth2AuthorizedNotRegistered:     acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Oauth2ChildRegistrationNotAllowed: acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Oauth2ChildRegistrationNotAllowedComplete: acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Oauth2CompleteRegistration:                acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Oauth2Device:                              acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Oauth2DeviceComplete:                      acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Oauth2Error:                               acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Oauth2Logout:                              acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Oauth2Passwordless:                        acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Oauth2Register:                            acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Oauth2StartIdPLink:                        acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Oauth2TwoFactor:                           acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Oauth2TwoFactorMethods:                    acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Oauth2Wait:                                acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		PasswordChange:                            acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		PasswordComplete:                          acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		PasswordForgot:                            acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		PasswordSent:                              acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		RegistrationComplete:                      acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		RegistrationSent:                          acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		RegistrationVerificationRequired:          acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		RegistrationVerify:                        acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Samlv2Logout:                              acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		Unauthorized:                              acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),

		// TODO(themes): test for deprecated properties
		//EmailSend:                         acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
		//RegistrationSend:                          acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
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
  name                                           = "test %[1]s"
  default_messages                               = <<EOF
%[2]sEOF
  stylesheet                                     = "%[3]s"
  account_edit                                   = "%[4]s"
  account_index                                  = "%[5]s"
  account_two_factor_disable                     = "%[6]s"
  account_two_factor_enable                      = "%[7]s"
  account_two_factor_index                       = "%[8]s"
  email_complete                                 = "%[9]s"
  email_sent                                     = "%[10]s"
  email_verification_required                    = "%[11]s"
  email_verify                                   = "%[12]s"
  helpers                                        = "%[13]s"
  index                                          = "%[14]s"
  oauth2_authorize                               = "%[15]s"
  oauth2_authorized_not_registered               = "%[16]s"
  oauth2_child_registration_not_allowed          = "%[17]s"
  oauth2_child_registration_not_allowed_complete = "%[18]s"
  oauth2_complete_registration                   = "%[19]s"
  oauth2_device                                  = "%[20]s"
  oauth2_device_complete                         = "%[21]s"
  oauth2_error                                   = "%[22]s"
  oauth2_logout                                  = "%[23]s"
  oauth2_passwordless                            = "%[24]s"
  oauth2_register                                = "%[25]s"
  oauth2_start_idp_link                          = "%[26]s"
  oauth2_two_factor                              = "%[27]s"
  oauth2_two_factor_methods                      = "%[28]s"
  oauth2_wait                                    = "%[29]s"
  password_change                                = "%[30]s"
  password_complete                              = "%[31]s"
  password_forgot                                = "%[32]s"
  password_sent                                  = "%[33]s"
  registration_complete                          = "%[34]s"
  registration_sent                              = "%[35]s"
  registration_verification_required             = "%[36]s"
  registration_verify                            = "%[37]s"
  samlv2_logout                                  = "%[38]s"
  unauthorized                                   = "%[39]s"

  # Deprecated Properties
  email_send                                     = "%[40]s"
  registration_send                              = "%[41]s"
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
		templates.Oauth2Device,
		templates.Oauth2DeviceComplete,
		templates.Oauth2Error,
		templates.Oauth2Logout,
		templates.Oauth2Passwordless,
		templates.Oauth2Register,
		templates.Oauth2StartIdPLink,
		templates.Oauth2TwoFactor,
		templates.Oauth2TwoFactorMethods,
		templates.Oauth2Wait,
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
