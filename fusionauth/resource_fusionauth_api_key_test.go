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

func TestAccAPIKey_basic(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_api_key.test_%s", resourceName)

	startDescription := "Test API Key"
	endDescription := "Updated Test API Key"
	startName := "test-api-key"
	endName := "updated-test-api-key"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				// Test resource create
				Config: testAccAPIKeyBasicConfig(
					resourceName,
					startDescription,
					startName,
				),
				Check: testAPIKeyAccTestCheckFuncs(
					tfResourcePath,
					startDescription,
					startName,
				),
			},
			{
				// Test resource update
				Config: testAccAPIKeyBasicConfig(
					resourceName,
					endDescription,
					endName,
				),
				Check: testAPIKeyAccTestCheckFuncs(
					tfResourcePath,
					endDescription,
					endName,
				),
			},
			{
				// Test importing resource into state
				ResourceName:            tfResourcePath,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key"},
			},
		},
	})
}

func TestAccAPIKey_withPermissions(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_api_key.test_%s", resourceName)

	description := "Test API Key with Permissions"
	name := "test-api-key-permissions"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIKeyWithPermissionsConfig(
					resourceName,
					description,
					name,
				),
				Check: testAPIKeyWithPermissionsCheckFuncs(
					tfResourcePath,
					description,
					name,
				),
			},
			{
				// Test importing resource into state
				ResourceName:            tfResourcePath,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key"},
			},
		},
	})
}

func TestAccAPIKey_tenantScoped(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_api_key.test_%s", resourceName)

	description := "Tenant Scoped API Key"
	name := "test-tenant-api-key"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIKeyTenantScopedConfig(
					resourceName,
					description,
					name,
				),
				Check: testAPIKeyTenantScopedCheckFuncs(
					tfResourcePath,
					description,
					name,
				),
			},
		},
	})
}

func testAPIKeyAccTestCheckFuncs(
	tfResourcePath string,
	description string,
	name string,
) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccCheckAPIKeyExists(tfResourcePath),
		resource.TestCheckResourceAttrSet(tfResourcePath, "key_id"),
		resource.TestCheckResourceAttrSet(tfResourcePath, "key"),
		resource.TestCheckResourceAttr(tfResourcePath, "description", description),
		resource.TestCheckResourceAttr(tfResourcePath, "name", name),
		resource.TestCheckResourceAttr(tfResourcePath, "retrievable", "true"),
	)
}

func testAPIKeyWithPermissionsCheckFuncs(
	tfResourcePath string,
	description string,
	name string,
) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccCheckAPIKeyExists(tfResourcePath),
		resource.TestCheckResourceAttr(tfResourcePath, "description", description),
		resource.TestCheckResourceAttr(tfResourcePath, "name", name),
		resource.TestCheckResourceAttr(tfResourcePath, "permissions_endpoints.#", "2"),
		resource.TestCheckTypeSetElemNestedAttrs(tfResourcePath, "permissions_endpoints.*", map[string]string{
			"endpoint": "/api/user",
			"get":      "true",
			"post":     "true",
			"put":      "true",
			"delete":   "false",
			"patch":    "false",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(tfResourcePath, "permissions_endpoints.*", map[string]string{
			"endpoint": "/api/application",
			"get":      "true",
			"post":     "false",
			"put":      "false",
			"delete":   "false",
			"patch":    "false",
		}),
	)
}

func testAPIKeyTenantScopedCheckFuncs(
	tfResourcePath string,
	description string,
	name string,
) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccCheckAPIKeyExists(tfResourcePath),
		resource.TestCheckResourceAttr(tfResourcePath, "description", description),
		resource.TestCheckResourceAttr(tfResourcePath, "name", name),
		resource.TestCheckResourceAttrSet(tfResourcePath, "tenant_id"),
	)
}

func testAccCheckAPIKeyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		client := fusionauthClient()
		resp, faErrs, err := client.RetrieveAPIKey(rs.Primary.ID)
		if errs := checkFusionauthErrors(faErrs, err); errs != nil {
			return err
		}

		if resp == nil || resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get resource: %#+v", resp)
		}

		return nil
	}
}

func testAccCheckAPIKeyDestroy(s *terraform.State) error {
	client := fusionauthClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_api_key" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			resp, faErrs, err := client.RetrieveAPIKey(rs.Primary.ID)
			if errs := checkFusionauthRetryErrors(faErrs, err); errs != nil {
				return errs
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

func testAccAPIKeyBasicConfig(resourceName, description, name string) string {
	return fmt.Sprintf(`
resource "fusionauth_api_key" "test_%[1]s" {
  description = "%[2]s"
  name        = "%[3]s"
  key_id      = "bd4efbcc-0667-4946-8b99-fd1474f25540"
  retrievable = true
}
`, resourceName, description, name)
}

func testAccAPIKeyWithPermissionsConfig(resourceName, description, name string) string {
	return fmt.Sprintf(`
resource "fusionauth_api_key" "test_%[1]s" {
  description = "%[2]s"
  name        = "%[3]s"
  retrievable = true

  permissions_endpoints {
    endpoint = "/api/user"
    get      = true
    post     = true
    put      = true
    delete   = false
    patch    = false
  }

  permissions_endpoints {
    endpoint = "/api/application"
    get      = true
    post     = false
    put      = false
    delete   = false
    patch    = false
  }
}
`, resourceName, description, name)
}

func testAccAPIKeyTenantScopedConfig(resourceName, description, name string) string {
	return testAccTenantResourceConfig(
		resourceName,
		"",
		"no-reply@example.com",
		30,
		false,
		false,
	) + fmt.Sprintf(`
resource "fusionauth_api_key" "test_%[1]s" {
  description = "%[2]s"
  name        = "%[3]s"
  tenant_id   = fusionauth_tenant.test_%[1]s.id
  retrievable = true
}
`, resourceName, description, name)
}
