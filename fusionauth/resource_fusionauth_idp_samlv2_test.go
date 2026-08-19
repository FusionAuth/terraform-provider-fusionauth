package fusionauth

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFusionauthIdpSAMLv2_verificationKeyIds(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_idp_saml_v2.test_%s", resourceName)
	keyIDs := sortedUUIDs(t, 3)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfFusionAuthBelow(t, "1.69.0") },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				// The tail is deliberately not sorted by key Id, which is the order FusionAuth reads
				// it back in; the empty plan after this step is the check on alignVerificationKeyIDs.
				Config: testAccIdpSAMLv2VerificationKeyIDsConfig(resourceName, keyIDs, idpKeyRefs(resourceName, 0, 2, 1)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthIdentityProviderExists(tfResourcePath),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.#", "3"),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.0", keyIDs[0]),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.1", keyIDs[2]),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.2", keyIDs[1]),
				),
			},
			{
				// Identity provider writes don't write the response back to state, so refresh
				// before asserting the scalar FusionAuth derived from the list.
				RefreshState: true,
				Check:        resource.TestCheckResourceAttr(tfResourcePath, "key_id", keyIDs[0]),
			},
			{
				// Promote a different key to the front of the list.
				Config: testAccIdpSAMLv2VerificationKeyIDsConfig(resourceName, keyIDs, idpKeyRefs(resourceName, 2, 0, 1)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthIdentityProviderExists(tfResourcePath),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.0", keyIDs[2]),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.1", keyIDs[0]),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.2", keyIDs[1]),
				),
			},
			{
				// The reorder must have moved the derived default, not just the state order.
				RefreshState: true,
				Check:        resource.TestCheckResourceAttr(tfResourcePath, "key_id", keyIDs[2]),
			},
		},
	})
}

// TestAccFusionauthIdpSAMLv2_deprecatedKeyId proves the deprecated key_id still round-trips on
// 1.69.0: FusionAuth seeds verification_key_ids from it, and the computed list must not then diff
// against the config that omits it.
func TestAccFusionauthIdpSAMLv2_deprecatedKeyId(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_idp_saml_v2.test_%s", resourceName)
	keyIDs := sortedUUIDs(t, 1)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfFusionAuthBelow(t, "1.69.0") },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdpSAMLv2DeprecatedKeyIDConfig(resourceName, keyIDs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthIdentityProviderExists(tfResourcePath),
					resource.TestCheckResourceAttr(tfResourcePath, "key_id", keyIDs[0]),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.#", "1"),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.0", keyIDs[0]),
				),
			},
		},
	})
}

func testAccIdpSAMLv2VerificationKeyIDsConfig(resourceName string, keyIDs []string, keyRefs string) string {
	return testAccIdpVerificationKeysConfig(resourceName, keyIDs) + fmt.Sprintf(`
resource "fusionauth_idp_saml_v2" "test_%[1]s" {
  name                 = "test-acc-samlv2 %[1]s"
  button_text          = "Login with SAML"
  idp_endpoint         = "https://example.com/%[1]s/login"
  use_name_for_email   = true
  verification_key_ids = %[2]s
}
`, resourceName, keyRefs)
}

func testAccIdpSAMLv2DeprecatedKeyIDConfig(resourceName string, keyIDs []string) string {
	return testAccIdpVerificationKeysConfig(resourceName, keyIDs) + fmt.Sprintf(`
resource "fusionauth_idp_saml_v2" "test_%[1]s" {
  name               = "test-acc-samlv2-deprecated %[1]s"
  button_text        = "Login with SAML"
  idp_endpoint       = "https://example.com/%[1]s/login"
  use_name_for_email = true
  key_id             = fusionauth_key.test_%[1]s_0.id
}
`, resourceName)
}
