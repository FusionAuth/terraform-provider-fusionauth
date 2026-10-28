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

// TestAccFusionauthEntity_actAsTokenExchangeClient covers the FusionAuth 1.70.0 token exchange
// client flag; skips on older servers.
func TestAccFusionauthEntity_actAsTokenExchangeClient(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_entity.test_%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfFusionAuthBelow(t, "1.70.0") },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthEntityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityActAsTokenExchangeClientConfig(resourceName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthEntityExists(tfResourcePath, true),
					resource.TestCheckResourceAttr(tfResourcePath, "act_as_token_exchange_client", "true"),
				),
			},
			{
				Config: testAccEntityActAsTokenExchangeClientConfig(resourceName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthEntityExists(tfResourcePath, false),
					resource.TestCheckResourceAttr(tfResourcePath, "act_as_token_exchange_client", "false"),
				),
			},
		},
	})
}

func testAccEntityActAsTokenExchangeClientConfig(resourceName string, actAsTokenExchangeClient bool) string {
	return fmt.Sprintf(`
resource "fusionauth_tenant" "test_%[1]s" {
  name = "test-acc-token-exchange %[1]s"
}

resource "fusionauth_entity_type" "test_%[1]s" {
  name = "test-acc-token-exchange %[1]s"
}

resource "fusionauth_entity" "test_%[1]s" {
  name                         = "test-acc-token-exchange %[1]s"
  entity_type_id               = fusionauth_entity_type.test_%[1]s.id
  tenant_id                    = fusionauth_tenant.test_%[1]s.id
  act_as_token_exchange_client = %[2]t
}
`, resourceName, actAsTokenExchangeClient)
}

func testAccCheckFusionauthEntityExists(resourceName string, actAsTokenExchangeClient bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		faClient := fusionauthClient()
		res, faErrs, err := faClient.RetrieveEntity(rs.Primary.ID)
		if errs := checkFusionauthErrors(faErrs, err); errs != nil {
			return errs
		}

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get resource: %#+v", res)
		}

		if res.Entity.ActAsTokenExchangeClient != actAsTokenExchangeClient {
			return fmt.Errorf("expected actAsTokenExchangeClient %t, got %t", actAsTokenExchangeClient, res.Entity.ActAsTokenExchangeClient)
		}

		return nil
	}
}

func testAccCheckFusionauthEntityDestroy(s *terraform.State) error {
	faClient := fusionauthClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_entity" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			res, faErrs, err := faClient.RetrieveEntity(rs.Primary.ID)
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
