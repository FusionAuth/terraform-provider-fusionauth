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

func TestAccTwilioMessenger(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_twilio_messenger.test_%s", resourceName)

	startAccountSID, endAccountSID := "SID12345678", "SID87654321"
	startAuthToken, endAuthToken := "token-secret-start", "token-secret-end"
	startFromPhoneNumber, endFromPhoneNumber := "+15551234567", "+15559876543"
	startName, endName := "my-test-messenger", "my-new-test-messenger"
	startURL, endURL := "https://api.twilio.com", "https://alt-api.twilio.com"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTwilioMessengerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTwilioMessengerBasicConfig(
					resourceName,
					startAccountSID,
					startAuthToken,
					startFromPhoneNumber,
					startName,
					startURL,
				),
				Check: testTwilioMessengerAccTestCheckFuncs(
					tfResourcePath,
					startAccountSID,
					startAuthToken,
					startFromPhoneNumber,
					startName,
					startURL,
				),
			},
			{
				Config: testAccTwilioMessengerBasicConfig(
					resourceName,
					endAccountSID,
					endAuthToken,
					endFromPhoneNumber,
					endName,
					endURL,
				),
				Check: testTwilioMessengerAccTestCheckFuncs(
					tfResourcePath,
					endAccountSID,
					endAuthToken,
					endFromPhoneNumber,
					endName,
					endURL,
				),
			},
		},
	})
}

func testAccTwilioMessengerBasicConfig(resourceName, accountSID, authToken, fromPhoneNumber, name, url string) string {
	return fmt.Sprintf(`
    # Twilio messenger setup
    resource "fusionauth_twilio_messenger" "test_%[1]s" {
        account_sid       = "%[2]s"
        auth_token        = "%[3]s"
        data              = jsonencode({ "important-key" : "important-value" })
        debug             = false
        from_phone_number = "%[4]s"
        name              = "%[5]s"
        url               = "%[6]s"
    }
    `, resourceName, accountSID, authToken, fromPhoneNumber, name, url)
}

func testTwilioMessengerAccTestCheckFuncs(tfResourcePath, accountSID, authToken, fromPhoneNumber, name, url string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccCheckTwilioMessengerExists(tfResourcePath),
		resource.TestCheckResourceAttr(tfResourcePath, "account_sid", accountSID),
		resource.TestCheckResourceAttr(tfResourcePath, "auth_token", authToken),
		testCheckResourceAttrJSON(tfResourcePath, "data", `{"important-key":"important-value"}`),
		resource.TestCheckResourceAttr(tfResourcePath, "debug", "false"),
		resource.TestCheckResourceAttr(tfResourcePath, "from_phone_number", fromPhoneNumber),
		resource.TestCheckResourceAttr(tfResourcePath, "name", name),
		resource.TestCheckResourceAttr(tfResourcePath, "url", url),
	)
}

func testAccCheckTwilioMessengerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		messenger, faErrs, err := RetrieveTwilioMessenger(context.Background(), fusionauthClient(), rs.Primary.ID)
		if errs := checkFusionauthErrors(faErrs, err); errs != nil {
			return err
		}

		if messenger == nil || messenger.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get resource: %#+v", messenger)
		}

		return nil
	}
}

func testAccCheckTwilioMessengerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_twilio_messenger" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			messenger, faErrs, err := RetrieveTwilioMessenger(context.Background(), fusionauthClient(), rs.Primary.ID)
			if errs := checkFusionauthRetryErrors(faErrs, err); errs != nil {
				return errs
			}

			if messenger != nil && messenger.StatusCode == http.StatusNotFound {
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
