package fusionauth

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestImportedKeyTypeValidationAcceptsSecret is a unit test that asserts the
// imported key "type" attribute accepts "Secret" - the key type used to store
// lambda secrets in Key Master - alongside the other supported key types, while
// still rejecting unknown values. It exercises the schema's ValidateFunc
// directly so it requires no running FusionAuth instance.
func TestImportedKeyTypeValidationAcceptsSecret(t *testing.T) {
	typeSchema, ok := resourceImportedKey().Schema["type"]
	if !ok {
		t.Fatal(`imported key resource is missing the "type" attribute`)
	}
	if typeSchema.ValidateFunc == nil {
		t.Fatal(`the "type" attribute does not define a ValidateFunc`)
	}

	// "Secret" must be accepted so lambda secrets can be imported into Key Master.
	for _, validType := range []string{"EC", "RSA", "HMAC", "Secret"} {
		warnings, errs := typeSchema.ValidateFunc(validType, "type")
		if len(warnings) != 0 || len(errs) != 0 {
			t.Errorf("type %q: expected it to be valid, got warnings=%v errors=%v", validType, warnings, errs)
		}
	}

	// An unsupported type must still be rejected.
	if _, errs := typeSchema.ValidateFunc("NotAKeyType", "type"); len(errs) == 0 {
		t.Error(`expected "NotAKeyType" to be rejected by the type ValidateFunc, but it was accepted`)
	}
}

// TestAccFusionauthImportedKey_Secret is an acceptance test that creates an
// imported key of type "Secret" against a live FusionAuth instance and confirms
// the key is created with the expected type. It is skipped automatically unless
// TF_ACC (and the FA_DOMAIN / FA_API_KEY credentials) are set.
func TestAccFusionauthImportedKey_Secret(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_imported_key.test_%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthImportedKeyDestroy,
		Steps: []resource.TestStep{
			{
				// Create a Secret key and verify it lands in FusionAuth as a Secret.
				Config: testAccImportedKeySecretConfig(resourceName, "super-secret-login-validation-value"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthImportedKeyExists(tfResourcePath),
					resource.TestCheckResourceAttr(tfResourcePath, "type", "Secret"),
					resource.TestCheckResourceAttr(tfResourcePath, "name", fmt.Sprintf("test-acc %s", resourceName)),
					resource.TestCheckResourceAttrSet(tfResourcePath, "id"),
				),
			},
		},
	})
}

func testAccCheckFusionauthImportedKeyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		faClient := fusionauthClient()
		key, faErrs, err := faClient.RetrieveKey(rs.Primary.ID)
		if errs := checkFusionauthErrors(faErrs, err); errs != nil {
			return err
		}

		if key == nil || key.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get resource: %#+v", key)
		}

		if key.Key.Type != fusionauth.KeyType_Secret {
			return fmt.Errorf("expected key type %q, got %q", fusionauth.KeyType_Secret, key.Key.Type)
		}

		return nil
	}
}

func testAccCheckFusionauthImportedKeyDestroy(s *terraform.State) error {
	faClient := fusionauthClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_imported_key" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			key, faErrs, err := faClient.RetrieveKey(rs.Primary.ID)
			if errs := checkFusionauthRetryErrors(faErrs, err); errs != nil {
				return errs
			}

			if key != nil && key.StatusCode == http.StatusNotFound {
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

func testAccImportedKeySecretConfig(name, secret string) string {
	return fmt.Sprintf(`
# Imported Key (Secret) Setup
resource "fusionauth_imported_key" "test_%[1]s" {
  name   = "test-acc %[1]s"
  type   = "Secret"
  secret = "%[2]s"
}
`, name, secret)
}
