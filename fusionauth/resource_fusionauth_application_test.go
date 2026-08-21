package fusionauth

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"testing"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testSAMLLogoutKey names the key used to sign SAML logout responses in the application tests.
const testSAMLLogoutKey = "samllogout"

func TestAccFusionauthApplication_verificationKeyIds(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_application.test_%s", resourceName)
	keyIDs := sortedUUIDs(t, 3)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfFusionAuthBelow(t, "1.69.0") },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationVerificationKeyIDsConfig(resourceName, keyIDs, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthApplicationExists(tfResourcePath),

					// jwt_configuration: unordered.
					resource.TestCheckResourceAttr(tfResourcePath, "jwt_configuration.0.access_token_verification_key_ids.#", "2"),
					resource.TestCheckTypeSetElemAttr(tfResourcePath, "jwt_configuration.0.access_token_verification_key_ids.*", keyIDs[0]),
					resource.TestCheckTypeSetElemAttr(tfResourcePath, "jwt_configuration.0.access_token_verification_key_ids.*", keyIDs[1]),
					resource.TestCheckResourceAttr(tfResourcePath, "jwt_configuration.0.id_token_verification_key_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(tfResourcePath, "jwt_configuration.0.id_token_verification_key_ids.*", keyIDs[2]),

					// samlv2_configuration: ordered, configured order must survive the read.
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.verification_key_ids.#", "3"),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.verification_key_ids.0", keyIDs[0]),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.verification_key_ids.1", keyIDs[2]),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.verification_key_ids.2", keyIDs[1]),
					// The deprecated scalar is derived from the first entry.
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.default_verification_key_id", keyIDs[0]),

					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.logout.0.verification_key_ids.#", "3"),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.logout.0.verification_key_ids.0", keyIDs[1]),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.logout.0.verification_key_ids.1", keyIDs[2]),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.logout.0.verification_key_ids.2", keyIDs[0]),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.logout.0.default_verification_key_id", keyIDs[1]),
				),
			},
			{
				// Rotate the trusted keys, including moving the SAML v2 default to another key.
				Config: testAccApplicationVerificationKeyIDsConfig(resourceName, keyIDs, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthApplicationExists(tfResourcePath),

					resource.TestCheckResourceAttr(tfResourcePath, "jwt_configuration.0.access_token_verification_key_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(tfResourcePath, "jwt_configuration.0.access_token_verification_key_ids.*", keyIDs[1]),
					resource.TestCheckResourceAttr(tfResourcePath, "jwt_configuration.0.id_token_verification_key_ids.#", "2"),
					resource.TestCheckTypeSetElemAttr(tfResourcePath, "jwt_configuration.0.id_token_verification_key_ids.*", keyIDs[0]),
					resource.TestCheckTypeSetElemAttr(tfResourcePath, "jwt_configuration.0.id_token_verification_key_ids.*", keyIDs[2]),

					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.verification_key_ids.#", "2"),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.verification_key_ids.0", keyIDs[2]),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.verification_key_ids.1", keyIDs[0]),

					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.logout.0.verification_key_ids.#", "1"),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.logout.0.verification_key_ids.0", keyIDs[0]),
				),
			},
			{
				// updateApplication doesn't write the response back to state, so refresh before
				// asserting the scalars FusionAuth re-derived from the rotated lists.
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.verification_key_ids.0", keyIDs[2]),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.verification_key_ids.1", keyIDs[0]),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.default_verification_key_id", keyIDs[2]),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.logout.0.verification_key_ids.0", keyIDs[0]),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.logout.0.default_verification_key_id", keyIDs[0]),
				),
			},
		},
	})
}

// TestAccFusionauthApplication_deprecatedVerificationKeyId proves the deprecated scalars still
// round-trip on 1.69.0: FusionAuth seeds the list from the scalar, and the computed list must not
// then diff against the config that omits it.
func TestAccFusionauthApplication_deprecatedVerificationKeyId(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_application.test_%s", resourceName)
	keyIDs := sortedUUIDs(t, 3)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfFusionAuthBelow(t, "1.69.0") },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDeprecatedVerificationKeyIDConfig(resourceName, keyIDs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthApplicationExists(tfResourcePath),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.default_verification_key_id", keyIDs[0]),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.verification_key_ids.#", "1"),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.verification_key_ids.0", keyIDs[0]),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.logout.0.default_verification_key_id", keyIDs[1]),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.logout.0.verification_key_ids.#", "1"),
					resource.TestCheckResourceAttr(tfResourcePath, "samlv2_configuration.0.logout.0.verification_key_ids.0", keyIDs[1]),
				),
			},
		},
	})
}

// sortedUUIDs generates n UUIDs sorted ascending, matching how FusionAuth orders verification keys
// that aren't the default.
func sortedUUIDs(t *testing.T, n int) []string {
	t.Helper()

	ids := make([]string, n)
	for i := range ids {
		id, err := uuid.GenerateUUID()
		if err != nil {
			t.Fatalf("error generating uuid: %s", err)
		}
		ids[i] = id
	}
	sort.Strings(ids)

	return ids
}

// testAccApplicationVerificationKeysConfig returns terraform configuration for the keys the
// verification key tests share: three keys pinned to sorted Ids, the two JWT signing keys, the SAML
// logout signing key, and the default tenant lookup.
//
// Note:
//   - `tenant_id` and `samlv2_configuration.logout.key_id` are set explicitly because both are
//     server assigned but not Computed, so omitting them leaves a pre-existing perpetual diff that
//     would mask the empty plan assertion these tests rely on.
func testAccApplicationVerificationKeysConfig(resourceName string, keyIDs []string) string {
	return testAccPinnedKeyResourceConfig(keyIDs[0], resourceName+"_a") +
		testAccPinnedKeyResourceConfig(keyIDs[1], resourceName+"_b") +
		testAccPinnedKeyResourceConfig(keyIDs[2], resourceName+"_c") +
		testAccAccessTokenKeyResourceConfig(resourceName) +
		testAccIDTokenKeyResourceConfig(resourceName) +
		testKeyConfig(testSAMLLogoutKey, resourceName) + `
data "fusionauth_tenant" "default" {
  name = "Default"
}
`
}

// testAccPinnedKeyResourceConfig returns terraform configuration for a key with a known Id.
func testAccPinnedKeyResourceConfig(keyID, name string) string {
	return testAccKeyResourceConfig(keyID, name, fusionauth.Algorithm_RS256, 2048, "FusionAuth")
}

func testAccApplicationVerificationKeyIDsConfig(resourceName string, keyIDs []string, rotated bool) string {
	accessTokenKeys := fmt.Sprintf("[fusionauth_key.test_%[1]s_a.id, fusionauth_key.test_%[1]s_b.id]", resourceName)
	idTokenKeys := fmt.Sprintf("[fusionauth_key.test_%[1]s_c.id]", resourceName)
	// The tail is deliberately not key Id sorted, which is the order FusionAuth reads back.
	samlv2Keys := fmt.Sprintf("[fusionauth_key.test_%[1]s_a.id, fusionauth_key.test_%[1]s_c.id, fusionauth_key.test_%[1]s_b.id]", resourceName)
	logoutKeys := fmt.Sprintf("[fusionauth_key.test_%[1]s_b.id, fusionauth_key.test_%[1]s_c.id, fusionauth_key.test_%[1]s_a.id]", resourceName)
	if rotated {
		accessTokenKeys = fmt.Sprintf("[fusionauth_key.test_%[1]s_b.id]", resourceName)
		idTokenKeys = fmt.Sprintf("[fusionauth_key.test_%[1]s_c.id, fusionauth_key.test_%[1]s_a.id]", resourceName)
		samlv2Keys = fmt.Sprintf("[fusionauth_key.test_%[1]s_c.id, fusionauth_key.test_%[1]s_a.id]", resourceName)
		logoutKeys = fmt.Sprintf("[fusionauth_key.test_%[1]s_a.id]", resourceName)
	}

	return testAccApplicationVerificationKeysConfig(resourceName, keyIDs) + fmt.Sprintf(`
resource "fusionauth_application" "test_%[1]s" {
  name      = "test-acc-verifykeys %[1]s"
  tenant_id = data.fusionauth_tenant.default.id

  jwt_configuration {
    enabled                           = true
    access_token_id                   = fusionauth_key.test_%[2]s.id
    id_token_key_id                   = fusionauth_key.test_%[3]s.id
    access_token_verification_key_ids = %[5]s
    id_token_verification_key_ids     = %[6]s
  }

  samlv2_configuration {
    enabled                  = true
    issuer                   = "https://example.com/%[1]s"
    authorized_redirect_urls = ["https://example.com/acs"]
    verification_key_ids     = %[7]s

    logout {
      key_id               = fusionauth_key.test_%[4]s.id
      verification_key_ids = %[8]s
    }
  }
}
`,
		resourceName,
		testKeyName(testAccessTokenKey, resourceName),
		testKeyName(testIDTokenKey, resourceName),
		testKeyName(testSAMLLogoutKey, resourceName),
		accessTokenKeys,
		idTokenKeys,
		samlv2Keys,
		logoutKeys,
	)
}

func testAccApplicationDeprecatedVerificationKeyIDConfig(resourceName string, keyIDs []string) string {
	return testAccApplicationVerificationKeysConfig(resourceName, keyIDs) + fmt.Sprintf(`
resource "fusionauth_application" "test_%[1]s" {
  name      = "test-acc-deprecated-verifykey %[1]s"
  tenant_id = data.fusionauth_tenant.default.id

  samlv2_configuration {
    enabled                     = true
    issuer                      = "https://example.com/%[1]s"
    authorized_redirect_urls    = ["https://example.com/acs"]
    default_verification_key_id = fusionauth_key.test_%[1]s_a.id

    logout {
      key_id                      = fusionauth_key.test_%[2]s.id
      default_verification_key_id = fusionauth_key.test_%[1]s_b.id
    }
  }
}
`, resourceName, testKeyName(testSAMLLogoutKey, resourceName))
}

func testAccCheckFusionauthApplicationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		faClient := fusionauthClient()
		app, err := faClient.RetrieveApplication(rs.Primary.ID)
		if err != nil {
			return err
		}

		if app == nil || app.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get resource: %#+v", app)
		}

		return nil
	}
}

// testAccCheckFusionauthApplicationDestroy verifies the application is gone. Deleting an
// application deactivates it rather than removing it, so an inactive application counts as
// destroyed.
func testAccCheckFusionauthApplicationDestroy(s *terraform.State) error {
	faClient := fusionauthClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_application" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			app, err := faClient.RetrieveApplication(rs.Primary.ID)
			if err != nil {
				return retry.NonRetryableError(err)
			}

			if app.StatusCode == http.StatusNotFound || app.Application.State != fusionauth.ObjectState_Active {
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
