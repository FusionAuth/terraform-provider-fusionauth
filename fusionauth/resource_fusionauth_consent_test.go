package fusionauth

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccConsent(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_consent.test_%s", resourceName)

	startName, endName := "Test Consent", "Updated Test Consent"
	startDefaultAge, endDefaultAge := "13", "16"
	startMultipleValues, endMultipleValues := "false", "true"
	startData, endData := `{"purpose": "marketing"}`, `{"purpose": "marketing", "detail": "email campaigns"}`
	startValues := []string{"email", "sms"}
	endValues := []string{"email", "sms", "phone"}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckConsentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConsentBasicConfig(
					resourceName,
					startName,
					startDefaultAge,
					startMultipleValues,
					startData,
					startValues,
					false,
				),
				Check: testConsentAccTestCheckFuncs(
					tfResourcePath,
					startName,
					startDefaultAge,
					startMultipleValues,
					startData,
					len(startValues),
				),
			},
			{
				Config: testAccConsentBasicConfig(
					resourceName,
					endName,
					endDefaultAge,
					endMultipleValues,
					endData,
					endValues,
					true,
				),
				Check: testConsentAccTestCheckFuncs(
					tfResourcePath,
					endName,
					endDefaultAge,
					endMultipleValues,
					endData,
					len(endValues),
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

func testAccConsentBasicConfig(
	resourceName,
	name,
	defaultAge,
	multipleValues,
	data string,
	values []string,
	withEmailPlus bool) string {
	valuesStr := ""
	for _, v := range values {
		valuesStr += fmt.Sprintf(`"%s", `, v)
	}

	emailPlusConfig := ""
	if withEmailPlus {
		emailPlusConfig = `
        email_plus {
            email_template_id = fusionauth_email.test_email_%[1]s.id
            enabled = true
            minimum_time_to_send_email_in_hours = 24
            maximum_time_to_send_email_in_hours = 48
        }
        `
		emailPlusConfig = fmt.Sprintf(emailPlusConfig, resourceName)
	}

	countryMinAgeConfig := `
    country_minimum_age_for_self_consent = {
        "us" = 13
        "ca" = 16
    }
    `

	return fmt.Sprintf(`
    resource "fusionauth_email" "test_email_%[1]s" {
        name = "consent-reminder-%[1]s"
        default_from_name = "No Reply"
        default_subject = "Consent Reminder"
        default_html_template = "<p>This is a reminder that you provided consent.</p>"
        default_text_template = "This is a reminder that you provided consent."
    }

    resource "fusionauth_consent" "test_%[1]s" {
        name = "%[2]s"
		consent_email_template_id = fusionauth_email.test_email_%[1]s.id
        default_minimum_age_for_self_consent = %[3]s
        multiple_values_allowed = %[4]s
        data = jsonencode(%[5]s)
        values = [%[6]s]
        %[7]s
        %[8]s
    }
    `, resourceName, name, defaultAge, multipleValues, data, valuesStr, countryMinAgeConfig, emailPlusConfig)
}

func testConsentAccTestCheckFuncs(
	tfResourcePath,
	name,
	defaultAge,
	multipleValues,
	data string,
	valuesCount int) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccCheckConsentExists(tfResourcePath),
		resource.TestCheckResourceAttr(tfResourcePath, "name", name),
		resource.TestCheckResourceAttr(tfResourcePath, "default_minimum_age_for_self_consent", defaultAge),
		resource.TestCheckResourceAttr(tfResourcePath, "multiple_values_allowed", multipleValues),
		testCheckResourceAttrJSON(tfResourcePath, "data", data),
		resource.TestCheckResourceAttr(tfResourcePath, "values.#", fmt.Sprintf("%d", valuesCount)),
		resource.TestCheckResourceAttr(tfResourcePath, "country_minimum_age_for_self_consent.us", "13"),
		resource.TestCheckResourceAttr(tfResourcePath, "country_minimum_age_for_self_consent.ca", "16"),
	)
}

func testAccCheckConsentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		client := fusionauthClient()
		resp, err := client.RetrieveConsent(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp == nil || resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get resource: %#+v", resp)
		}

		return nil
	}
}

func testAccCheckConsentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_consent" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			client := fusionauthClient()
			resp, err := client.RetrieveConsent(rs.Primary.ID)
			if err != nil {
				return retry.NonRetryableError(err)
			}

			if resp != nil && resp.StatusCode == http.StatusNotFound {
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
