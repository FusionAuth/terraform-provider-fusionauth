package fusionauth

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFusionauthIdpExternalJWT_verificationKeyIds(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_idp_external_jwt.test_%s", resourceName)
	keyIDs := sortedUUIDs(t, 3)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfFusionAuthBelow(t, "1.69.0") },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdpExternalJWTVerificationKeyIDsConfig(resourceName, keyIDs, idpKeyRefs(resourceName, 0, 2, 1)),
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.0", keyIDs[0]),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.1", keyIDs[2]),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.2", keyIDs[1]),
					// The deprecated scalar is derived from the first entry.
					resource.TestCheckResourceAttr(tfResourcePath, "default_key_id", keyIDs[0]),
				),
			},
			{
				// Rotate the fallback key onto another entry and drop one of the trusted keys.
				Config: testAccIdpExternalJWTVerificationKeyIDsConfig(resourceName, keyIDs, idpKeyRefs(resourceName, 2, 0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthIdentityProviderExists(tfResourcePath),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.#", "2"),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.0", keyIDs[2]),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.1", keyIDs[0]),
				),
			},
			{
				// Identity provider updates don't write the response back to state, so refresh
				// before asserting the scalar FusionAuth re-derived from the rotated list.
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.0", keyIDs[2]),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.1", keyIDs[0]),
					resource.TestCheckResourceAttr(tfResourcePath, "default_key_id", keyIDs[2]),
				),
			},
		},
	})
}

// TestAccFusionauthIdpExternalJWT_deprecatedDefaultKeyId proves the deprecated default_key_id still
// round-trips on 1.69.0: FusionAuth seeds verification_key_ids from it, and the computed list must
// not then diff against the config that omits it.
func TestAccFusionauthIdpExternalJWT_deprecatedDefaultKeyId(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_idp_external_jwt.test_%s", resourceName)
	keyIDs := sortedUUIDs(t, 1)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfFusionAuthBelow(t, "1.69.0") },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdpExternalJWTDeprecatedDefaultKeyIDConfig(resourceName, keyIDs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthIdentityProviderExists(tfResourcePath),
					resource.TestCheckResourceAttr(tfResourcePath, "default_key_id", keyIDs[0]),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.#", "1"),
					resource.TestCheckResourceAttr(tfResourcePath, "verification_key_ids.0", keyIDs[0]),
				),
			},
		},
	})
}

func testAccIdpExternalJWTVerificationKeyIDsConfig(resourceName string, keyIDs []string, keyRefs string) string {
	return testAccIdpVerificationKeysConfig(resourceName, keyIDs) + fmt.Sprintf(`
resource "fusionauth_idp_external_jwt" "test_%[1]s" {
  name                 = "test-acc-extjwt %[1]s"
  header_key_parameter = "kid"
  enabled              = true
  verification_key_ids = %[2]s
}
`, resourceName, keyRefs)
}

func testAccIdpExternalJWTDeprecatedDefaultKeyIDConfig(resourceName string, keyIDs []string) string {
	return testAccIdpVerificationKeysConfig(resourceName, keyIDs) + fmt.Sprintf(`
resource "fusionauth_idp_external_jwt" "test_%[1]s" {
  name                 = "test-acc-extjwt-deprecated %[1]s"
  header_key_parameter = "kid"
  enabled              = true
  default_key_id       = fusionauth_key.test_%[1]s_0.id
}
`, resourceName)
}
