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

func TestAccGenericConnector(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_generic_connector.test_%s", resourceName)

	startAuthenticationURL, endAuthenticationURL := "http://example-start.com", "http://example-end.com"
	startConnectTimeout, endConnectTimeout := "2048", "4096"
	startAuthenticationPassword, endAuthenticationPassword := "super-secret-start", "super-secret-end"
	startAuthenicationUsername, endAuthenticationUsername := "me", "me-too"
	startName, endName := "my-test-connector", "my-new-test-connector"
	startReadTimeout, endReadTimeout := "1111", "2222"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckGenericConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGenericConnectorBasicConfig(
					resourceName,
					startAuthenticationURL,
					startConnectTimeout,
					startAuthenticationPassword,
					startAuthenicationUsername,
					startName,
					startReadTimeout,
				),
				Check: testGenericConnectorAccTestCheckFuncs(
					tfResourcePath,
					startAuthenticationURL,
					startConnectTimeout,
					startAuthenticationPassword,
					startAuthenicationUsername,
					startName,
					startReadTimeout,
				),
			},
			{
				Config: testAccGenericConnectorBasicConfig(
					resourceName,
					endAuthenticationURL,
					endConnectTimeout,
					endAuthenticationPassword,
					endAuthenticationUsername,
					endName,
					endReadTimeout,
				),
				Check: testGenericConnectorAccTestCheckFuncs(
					tfResourcePath,
					endAuthenticationURL,
					endConnectTimeout,
					endAuthenticationPassword,
					endAuthenticationUsername,
					endName,
					endReadTimeout,
				),
			},
		},
	})
}

func testAccGenericConnectorBasicConfig(resourceName, authenticationURL, connectTimeout, authenticationPassword, authenticationUsername, name, readTimeout string) string {
	return fmt.Sprintf(`
	# Generic connector setup
	resource "fusionauth_generic_connector" "test_%[1]s" {
		authentication_url           = "%[2]s"
		connect_timeout              = %[3]s
		data                         = { "important-key" : "important-value" }
		debug                        = false
		headers                      = { "header-key-1" : "one", "header-key-2" : "two" }
		http_authentication_password = "%[4]s"
		http_authentication_username = "%[5]s"
		name                         = "%[6]s"
		read_timeout                 = %[7]s
		}
	`, resourceName, authenticationURL, connectTimeout, authenticationPassword, authenticationUsername, name, readTimeout)
}

func testGenericConnectorAccTestCheckFuncs(tfResourcePath, authenticationURL, connectTimeout, authenticationPassword, authenticationUsername, name, readTimeout string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccCheckGenericConnectorExists(tfResourcePath),
		resource.TestCheckResourceAttr(tfResourcePath, "authentication_url", authenticationURL),
		resource.TestCheckResourceAttr(tfResourcePath, "connect_timeout", connectTimeout),
		resource.TestCheckResourceAttr(tfResourcePath, "data.important-key", "important-value"),
		resource.TestCheckResourceAttr(tfResourcePath, "debug", "false"),
		resource.TestCheckResourceAttr(tfResourcePath, "headers.header-key-1", "one"),
		resource.TestCheckResourceAttr(tfResourcePath, "headers.header-key-2", "two"),
		resource.TestCheckResourceAttr(tfResourcePath, "http_authentication_password", authenticationPassword),
		resource.TestCheckResourceAttr(tfResourcePath, "http_authentication_username", authenticationUsername),
		resource.TestCheckResourceAttr(tfResourcePath, "name", name),
		resource.TestCheckResourceAttr(tfResourcePath, "read_timeout", readTimeout),
	)
}

func testAccCheckGenericConnectorExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		connector, faErrs, err := RetrieveConnector(context.Background(), fusionauthClient(), rs.Primary.ID)
		if errs := checkFusionauthErrors(faErrs, err); errs != nil {
			return err
		}

		if connector == nil || connector.StatusCode != http.StatusOK {
			// This is a weird edge case...
			return fmt.Errorf("failed to get resource: %#+v", connector)
		}

		return nil
	}
}

func testAccCheckGenericConnectorDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_generic_connector" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			connector, faErrs, err := RetrieveConnector(context.Background(), fusionauthClient(), rs.Primary.ID)
			if errs := checkFusionauthRetryErrors(faErrs, err); errs != nil {
				return errs
			}

			if connector != nil && connector.StatusCode == http.StatusNotFound {
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
