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

func TestAccLDAPConnector(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_ldap_connector.test_%s", resourceName)

	startAuthenticationURL, endAuthenticationURL := "ldap://example-start.com", "ldap://example-end.com"
	startBaseStructure, endBaseStructure := "dc=example,dc=start,dc=com", "dc=example,dc=end,dc=com"
	startConnectTimeout, endConnectTimeout := "2048", "4096"
	startIdentifyingAttribute, endIdentifyingAttribute := "uid", "cn"
	startLoginIdAttribute, endLoginIdAttribute := "mail", "emailAddress"
	startName, endName := "my-test-ldap-connector", "my-new-test-ldap-connector"
	startReadTimeout, endReadTimeout := "1111", "2222"
	startSystemAccountDN, endSystemAccountDN := "cn=admin,dc=example,dc=start,dc=com", "cn=admin,dc=example,dc=end,dc=com"
	startSystemAccountPassword, endSystemAccountPassword := "ldap-secret-start", "ldap-secret-end"
	startSecurityMethod, endSecurityMethod := "None", "LDAPS"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckLDAPConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLDAPConnectorBasicConfig(
					resourceName,
					startAuthenticationURL,
					startBaseStructure,
					startConnectTimeout,
					startIdentifyingAttribute,
					startLoginIdAttribute,
					startName,
					startReadTimeout,
					startSecurityMethod,
					startSystemAccountDN,
					startSystemAccountPassword,
				),
				Check: testLDAPConnectorAccTestCheckFuncs(
					tfResourcePath,
					startAuthenticationURL,
					startBaseStructure,
					startConnectTimeout,
					startIdentifyingAttribute,
					startLoginIdAttribute,
					startName,
					startReadTimeout,
					startSecurityMethod,
					startSystemAccountDN,
					startSystemAccountPassword,
				),
			},
			{
				Config: testAccLDAPConnectorBasicConfig(
					resourceName,
					endAuthenticationURL,
					endBaseStructure,
					endConnectTimeout,
					endIdentifyingAttribute,
					endLoginIdAttribute,
					endName,
					endReadTimeout,
					endSecurityMethod,
					endSystemAccountDN,
					endSystemAccountPassword,
				),
				Check: testLDAPConnectorAccTestCheckFuncs(
					tfResourcePath,
					endAuthenticationURL,
					endBaseStructure,
					endConnectTimeout,
					endIdentifyingAttribute,
					endLoginIdAttribute,
					endName,
					endReadTimeout,
					endSecurityMethod,
					endSystemAccountDN,
					endSystemAccountPassword,
				),
			},
		},
	})
}

func testAccLDAPConnectorBasicConfig(
	resourceName,
	authenticationURL,
	baseStructure,
	connectTimeout,
	identifyingAttribute,
	loginIdAttribute,
	name,
	readTimeout,
	securityMethod,
	systemAccountDN,
	systemAccountPassword string) string {

	return fmt.Sprintf(`
	resource "fusionauth_lambda" "test_%[1]s" {
		name    = "test_%[1]s"
		type    = "LDAPConnectorReconcile"
		enabled = true
		body    = <<EOT
		// Using the response from an LDAP connector, reconcile the User.
		function reconcile(user, userAttributes) {
		console.info('Hello World!');
		}
		EOT
	}

    # LDAP connector setup
    resource "fusionauth_ldap_connector" "test_%[1]s" {
        authentication_url     = "%[2]s"
        base_structure         = "%[3]s"
        connect_timeout        = %[4]s
        data                   = jsonencode({ "important-key" : "important-value" })
        debug                  = false
        identifying_attribute  = "%[5]s"
        lambda_configuration {
            reconcile_id       = fusionauth_lambda.test_%[1]s.id
        }
        login_id_attribute     = "%[6]s"
        name                   = "%[7]s"
        read_timeout          = %[8]s
        requested_attributes   = ["givenName", "sn", "mail"]
        security_method        = "%[9]s"
        system_account_dn      = "%[10]s"
        system_account_password = "%[11]s"
    }
    `, resourceName, authenticationURL, baseStructure, connectTimeout, identifyingAttribute, loginIdAttribute, name, readTimeout, securityMethod, systemAccountDN, systemAccountPassword)
}

func testLDAPConnectorAccTestCheckFuncs(
	tfResourcePath,
	authenticationURL,
	baseStructure,
	connectTimeout,
	identifyingAttribute,
	loginIdAttribute,
	name,
	readTimeout,
	securityMethod,
	systemAccountDN,
	systemAccountPassword string) resource.TestCheckFunc {

	return resource.ComposeTestCheckFunc(
		testAccCheckLDAPConnectorExists(tfResourcePath),
		resource.TestCheckResourceAttr(tfResourcePath, "authentication_url", authenticationURL),
		resource.TestCheckResourceAttr(tfResourcePath, "base_structure", baseStructure),
		resource.TestCheckResourceAttr(tfResourcePath, "connect_timeout", connectTimeout),
		resource.TestCheckResourceAttr(tfResourcePath, "debug", "false"),
		resource.TestCheckResourceAttr(tfResourcePath, "identifying_attribute", identifyingAttribute),
		resource.TestCheckResourceAttr(tfResourcePath, "login_id_attribute", loginIdAttribute),
		resource.TestCheckResourceAttr(tfResourcePath, "name", name),
		resource.TestCheckResourceAttr(tfResourcePath, "read_timeout", readTimeout),
		resource.TestCheckResourceAttr(tfResourcePath, "security_method", securityMethod),
		resource.TestCheckResourceAttr(tfResourcePath, "system_account_dn", systemAccountDN),
		resource.TestCheckResourceAttr(tfResourcePath, "system_account_password", systemAccountPassword),
		resource.TestCheckResourceAttr(tfResourcePath, "requested_attributes.#", "3"),
		resource.TestCheckTypeSetElemAttr(tfResourcePath, "requested_attributes.*", "givenName"),
		resource.TestCheckTypeSetElemAttr(tfResourcePath, "requested_attributes.*", "sn"),
		resource.TestCheckTypeSetElemAttr(tfResourcePath, "requested_attributes.*", "mail"),
	)
}

func testAccCheckLDAPConnectorExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		connector, faErrs, err := RetrieveLDAPConnector(context.Background(), fusionauthClient(), rs.Primary.ID)
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

func testAccCheckLDAPConnectorDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_ldap_connector" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			connector, faErrs, err := RetrieveLDAPConnector(context.Background(), fusionauthClient(), rs.Primary.ID)
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
