package fusionauth

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testAccCheckFusionauthIdentityProviderExists verifies the identity provider can be read back.
func testAccCheckFusionauthIdentityProviderExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		if _, err := readIdentityProvider(rs.Primary.ID, fusionauthProviderClient()); err != nil {
			return err
		}

		return nil
	}
}

// testAccCheckFusionauthIdentityProviderDestroy verifies every identity provider in state is gone.
func testAccCheckFusionauthIdentityProviderDestroy(s *terraform.State) error {
	client := fusionauthProviderClient()

	for _, rs := range s.RootModule().Resources {
		if !strings.HasPrefix(rs.Type, "fusionauth_idp_") {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			if _, err := readIdentityProvider(rs.Primary.ID, client); err != nil {
				if err.Error() == NotFoundError {
					// resource destroyed!
					return nil
				}

				return retry.NonRetryableError(err)
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

// testAccIdpVerificationKeysConfig returns terraform configuration for keys pinned to the supplied
// Ids, named after the identity provider under test.
func testAccIdpVerificationKeysConfig(resourceName string, keyIDs []string) string {
	config := make([]string, 0, len(keyIDs))
	for i, keyID := range keyIDs {
		config = append(config, testAccPinnedKeyResourceConfig(keyID, fmt.Sprintf("%s_%d", resourceName, i)))
	}

	return strings.Join(config, "")
}

// idpKeyRefs renders an HCL list of key resource references for the supplied key indexes.
func idpKeyRefs(resourceName string, indexes ...int) string {
	refs := make([]string, 0, len(indexes))
	for _, i := range indexes {
		refs = append(refs, fmt.Sprintf("fusionauth_key.test_%s_%d.id", resourceName, i))
	}

	return "[" + strings.Join(refs, ", ") + "]"
}
