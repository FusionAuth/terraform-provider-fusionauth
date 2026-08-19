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

// TestAccFusionauthEntityType_verificationKeyIds covers the FusionAuth 1.69.0 entity JWT
// verification key list; skips on older servers.
func TestAccFusionauthEntityType_verificationKeyIds(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_entity_type.test_%s", resourceName)
	keyIDs := sortedUUIDs(t, 2)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfFusionAuthBelow(t, "1.69.0") },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthEntityTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityTypeVerificationKeyIDsConfig(resourceName, keyIDs, idpKeyRefs(resourceName, 0, 1)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthEntityTypeExists(tfResourcePath),
					resource.TestCheckResourceAttr(tfResourcePath, "jwt_configuration.0.access_token_verification_key_ids.#", "2"),
					resource.TestCheckTypeSetElemAttr(tfResourcePath, "jwt_configuration.0.access_token_verification_key_ids.*", keyIDs[0]),
					resource.TestCheckTypeSetElemAttr(tfResourcePath, "jwt_configuration.0.access_token_verification_key_ids.*", keyIDs[1]),
				),
			},
			{
				// Rotate off the first key.
				Config: testAccEntityTypeVerificationKeyIDsConfig(resourceName, keyIDs, idpKeyRefs(resourceName, 1)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthEntityTypeExists(tfResourcePath),
					resource.TestCheckResourceAttr(tfResourcePath, "jwt_configuration.0.access_token_verification_key_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(tfResourcePath, "jwt_configuration.0.access_token_verification_key_ids.*", keyIDs[1]),
				),
			},
		},
	})
}

// testAccEntityTypeVerificationKeyIDsConfig returns terraform configuration for an entity type that
// trusts additional access token verification keys. The signing key is a separate key, so the
// trusted keys are always keys the entity type doesn't already sign with.
func testAccEntityTypeVerificationKeyIDsConfig(resourceName string, keyIDs []string, keyRefs string) string {
	return testAccIdpVerificationKeysConfig(resourceName, keyIDs) +
		testAccAccessTokenKeyResourceConfig(resourceName) +
		fmt.Sprintf(`
resource "fusionauth_entity_type" "test_%[1]s" {
  name = "test-acc-verifykeys %[1]s"

  jwt_configuration {
    enabled                           = true
    access_token_key_id               = fusionauth_key.test_%[2]s.id
    time_to_live_in_seconds           = 3600
    access_token_verification_key_ids = %[3]s
  }
}
`, resourceName, testKeyName(testAccessTokenKey, resourceName), keyRefs)
}

func testAccCheckFusionauthEntityTypeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		faClient := fusionauthClient()
		res, faErrs, err := faClient.RetrieveEntityType(rs.Primary.ID)
		if errs := checkFusionauthErrors(faErrs, err); errs != nil {
			return errs
		}

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get resource: %#+v", res)
		}

		return nil
	}
}

func testAccCheckFusionauthEntityTypeDestroy(s *terraform.State) error {
	faClient := fusionauthClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_entity_type" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			res, faErrs, err := faClient.RetrieveEntityType(rs.Primary.ID)
			if errs := checkFusionauthRetryErrors(faErrs, err); errs != nil {
				return errs
			}

			if res.StatusCode == http.StatusNotFound {
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
