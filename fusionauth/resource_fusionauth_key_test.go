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
)

func TestAccFusionauthKey_basic(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	tfResourcePath := fmt.Sprintf("fusionauth_key.test_%s", resourceName)

	startAlgorithm, endAlgorithm := fusionauth.Algorithm_RS256, fusionauth.Algorithm_RS512
	startLength, endLength := 2048, 4096

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthKeyDestroy,
		Steps: []resource.TestStep{
			{
				// Test resource create
				Config: testAccKeyResourceConfig(resourceName, startAlgorithm, startLength),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthKeyExists(tfResourcePath),
					resource.TestCheckResourceAttrSet(tfResourcePath, "key_id"),
					resource.TestCheckResourceAttr(tfResourcePath, "name", fmt.Sprintf("test %s", resourceName)),
					resource.TestCheckResourceAttr(tfResourcePath, "algorithm", string(startAlgorithm)),
					resource.TestCheckResourceAttr(tfResourcePath, "length", fmt.Sprintf("%d", startLength)),
				),
			},
			{
				// Test resource update/state mutate
				Config: testAccKeyResourceConfig(resourceName, endAlgorithm, endLength),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthKeyExists(tfResourcePath),
					resource.TestCheckResourceAttrSet(tfResourcePath, "key_id"),
					resource.TestCheckResourceAttr(tfResourcePath, "name", fmt.Sprintf("test %s", resourceName)),
					resource.TestCheckResourceAttr(tfResourcePath, "algorithm", string(endAlgorithm)),
					resource.TestCheckResourceAttr(tfResourcePath, "length", fmt.Sprintf("%d", endLength)),
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

func testAccCheckFusionauthKeyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		key, faErrs, err := fusionauthClient().RetrieveKey(rs.Primary.ID)
		if errs := checkFusionauthErrors(faErrs, err); errs != nil {
			return err
		}

		if key == nil || key.StatusCode != http.StatusOK {
			// This is a weird edge case...
			return fmt.Errorf("failed to get resource: %#+v\n", key)
		}

		return nil
	}
}

func testAccCheckFusionauthKeyDestroy(s *terraform.State) error {
	faClient := fusionauthClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_key" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := resource.RetryContext(context.Background(), retryTimeout, func() *resource.RetryError {
			key, faErrs, err := faClient.RetrieveKey(rs.Primary.ID)
			if errs := checkFusionauthRetryErrors(faErrs, err); errs != nil {
				return errs
			}

			if key != nil && key.StatusCode == http.StatusNotFound {
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

func testAccKeyResourceConfig(name string, algorithm fusionauth.Algorithm, length int) string {
	return fmt.Sprintf(`
# Key Setup
resource "fusionauth_key" "test_%[1]s" {
  name      = "test %[1]s"
  algorithm = "%[2]s"
  length    = %[3]d
}
`, name, algorithm, length)
}

// testAccAccessTokenKeyResourceConfig returns terraform configuration to
// generate a standalone test Access Token key.
func testAccAccessTokenKeyResourceConfig() string {
	return testAccKeyResourceConfig(
		"accesstoken",
		fusionauth.Algorithm_RS256,
		2048,
	)
}

// testAccIdTokenKeyResourceConfig returns terraform configuration to generate a
// standalone test ID Token key.
func testAccIdTokenKeyResourceConfig() string {
	return testAccKeyResourceConfig(
		"idtoken",
		fusionauth.Algorithm_RS256,
		2048,
	)
}
