package fusionauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSMSMessageTemplate(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_sms_message_template.test_%s", resourceName)

	startName, endName := "Test SMS Template", "Updated SMS Template"
	startDefaultTemplate, endDefaultTemplate := "Your verification code is $${code}", "Your new verification code is $${code}"
	startLocalizedFR, endLocalizedFR := "Votre code de vérification est $${code}", "Votre nouveau code de vérification est $${code}"
	startLocalizedES, endLocalizedES := "Su código de verificación es $${code}", "Su nuevo código de verificación es $${code}"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSMSMessageTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSMSMessageTemplateBasicConfig(
					resourceName,
					startName,
					startDefaultTemplate,
					startLocalizedFR,
					startLocalizedES,
				),
				Check: testSMSMessageTemplateAccTestCheckFuncs(
					tfResourcePath,
					startName,
					startDefaultTemplate,
					startLocalizedFR,
					startLocalizedES,
				),
			},
			{
				Config: testAccSMSMessageTemplateBasicConfig(
					resourceName,
					endName,
					endDefaultTemplate,
					endLocalizedFR,
					endLocalizedES,
				),
				Check: testSMSMessageTemplateAccTestCheckFuncs(
					tfResourcePath,
					endName,
					endDefaultTemplate,
					endLocalizedFR,
					endLocalizedES,
				),
			},
			{
				ResourceName:      tfResourcePath,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSMSMessageTemplateBasicConfig(
	resourceName,
	name,
	defaultTemplate,
	localizedFR,
	localizedES string) string {
	return fmt.Sprintf(`
    resource "fusionauth_sms_message_template" "test_%[1]s" {
        name             = "%[2]s"
        default_template = "%[3]s"
        data             = jsonencode({ "usage": "verification", "notes": "Used for 2FA verification" })
        localized_templates = {
            "fr" = "%[4]s",
            "es" = "%[5]s"
        }
    }
    `, resourceName, name, defaultTemplate, localizedFR, localizedES)
}

func testSMSMessageTemplateAccTestCheckFuncs(
	tfResourcePath,
	name,
	defaultTemplate,
	localizedFR,
	localizedES string) resource.TestCheckFunc {
	// Replace $$ with $ to avoid escaping in the test.
	defaultTemplate = strings.ReplaceAll(defaultTemplate, "$$", "$")
	localizedFR = strings.ReplaceAll(localizedFR, "$$", "$")
	localizedES = strings.ReplaceAll(localizedES, "$$", "$")

	return resource.ComposeTestCheckFunc(
		testAccCheckSMSMessageTemplateExists(tfResourcePath),
		resource.TestCheckResourceAttr(tfResourcePath, "name", name),
		resource.TestCheckResourceAttr(tfResourcePath, "default_template", defaultTemplate),
		resource.TestCheckResourceAttr(tfResourcePath, "localized_templates.fr", localizedFR),
		resource.TestCheckResourceAttr(tfResourcePath, "localized_templates.es", localizedES),
		resource.TestCheckResourceAttrSet(tfResourcePath, "data"),
	)
}

func testAccCheckSMSMessageTemplateExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		template, faErrs, err := retrieveMessageTemplate(context.Background(), fusionauthClient(), rs.Primary.ID)
		if errs := checkFusionauthErrors(faErrs, err); errs != nil {
			return err
		}

		if template == nil || template.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get resource: %#+v", template)
		}

		return nil
	}
}

func testAccCheckSMSMessageTemplateDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_sms_message_template" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			template, faErrs, err := retrieveMessageTemplate(context.Background(), fusionauthClient(), rs.Primary.ID)
			if errs := checkFusionauthRetryErrors(faErrs, err); errs != nil {
				return errs
			}

			if template != nil && template.StatusCode == http.StatusNotFound {
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
