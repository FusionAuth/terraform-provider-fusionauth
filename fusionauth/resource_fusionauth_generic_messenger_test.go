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

func TestAccGenericMessenger(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_generic_messenger.test_%s", resourceName)

	startURL, endURL := "http://example-start.com", "http://example-end.com"
	startConnectTimeout, endConnectTimeout := "2048", "4096"
	startHttpAuthenticationPassword, endHttpAuthenticationPassword := "super-secret-start", "super-secret-end"
	startHttpAuthenticationUsername, endHttpAuthenticationUsername := "me", "me-too"
	startName, endName := "my-test-messenger", "my-new-test-messenger"
	startReadTimeout, endReadTimeout := "1111", "2222"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckGenericMessengerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGenericMessengerBasicConfig(
					resourceName,
					startURL,
					startConnectTimeout,
					startHttpAuthenticationPassword,
					startHttpAuthenticationUsername,
					startName,
					startReadTimeout,
				),
				Check: testGenericMessengerAccTestCheckFuncs(
					tfResourcePath,
					startURL,
					startConnectTimeout,
					startHttpAuthenticationPassword,
					startHttpAuthenticationUsername,
					startName,
					startReadTimeout,
				),
			},
			{
				Config: testAccGenericMessengerBasicConfig(
					resourceName,
					endURL,
					endConnectTimeout,
					endHttpAuthenticationPassword,
					endHttpAuthenticationUsername,
					endName,
					endReadTimeout,
				),
				Check: testGenericMessengerAccTestCheckFuncs(
					tfResourcePath,
					endURL,
					endConnectTimeout,
					endHttpAuthenticationPassword,
					endHttpAuthenticationUsername,
					endName,
					endReadTimeout,
				),
			},
		},
	})
}

func testAccGenericMessengerBasicConfig(resourceName, url, connectTimeout, httpAuthenticationPassword, httpAuthenticationUsername, name, readTimeout string) string {
	return fmt.Sprintf(`
    # Generic messenger setup
    resource "fusionauth_generic_messenger" "test_%[1]s" {
        url                         = "%[2]s"
        connect_timeout             = %[3]s
        data                        = jsonencode({ "important-key" : "important-value" })
        debug                       = false
        headers                     = { "header-key-1" : "one", "header-key-2" : "two" }
        http_authentication_password = "%[4]s"
        http_authentication_username = "%[5]s"
        name                        = "%[6]s"
        read_timeout                = %[7]s
        }
    `, resourceName, url, connectTimeout, httpAuthenticationPassword, httpAuthenticationUsername, name, readTimeout)
}

func testGenericMessengerAccTestCheckFuncs(tfResourcePath, url, connectTimeout, httpAuthenticationPassword, httpAuthenticationUsername, name, readTimeout string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccCheckGenericMessengerExists(tfResourcePath),
		resource.TestCheckResourceAttr(tfResourcePath, "url", url),
		resource.TestCheckResourceAttr(tfResourcePath, "connect_timeout", connectTimeout),
		testCheckResourceAttrJSON(tfResourcePath, "data", `{"important-key":"important-value"}`),
		resource.TestCheckResourceAttr(tfResourcePath, "debug", "false"),
		resource.TestCheckResourceAttr(tfResourcePath, "headers.header-key-1", "one"),
		resource.TestCheckResourceAttr(tfResourcePath, "headers.header-key-2", "two"),
		resource.TestCheckResourceAttr(tfResourcePath, "http_authentication_password", httpAuthenticationPassword),
		resource.TestCheckResourceAttr(tfResourcePath, "http_authentication_username", httpAuthenticationUsername),
		resource.TestCheckResourceAttr(tfResourcePath, "name", name),
		resource.TestCheckResourceAttr(tfResourcePath, "read_timeout", readTimeout),
	)
}

func testAccCheckGenericMessengerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		messenger, faErrs, err := RetrieveMessenger(context.Background(), fusionauthClient(), rs.Primary.ID)
		if errs := checkFusionauthErrors(faErrs, err); errs != nil {
			return err
		}

		if messenger == nil || messenger.StatusCode != http.StatusOK {
			// This is a weird edge case...
			return fmt.Errorf("failed to get resource: %#+v", messenger)
		}

		return nil
	}
}

func testAccCheckGenericMessengerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_generic_messenger" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			messenger, faErrs, err := RetrieveMessenger(context.Background(), fusionauthClient(), rs.Primary.ID)
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
