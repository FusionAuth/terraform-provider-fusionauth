package fusionauth

import (
	"fmt"
	"testing"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFusionauthIdpSAMLv2_basic(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_idp_saml_v2.test_%s", resourceName)
	keyResourcePath := fmt.Sprintf("fusionauth_key.test_%s", resourceName)

	startButtonText, endButtonText := "test-acc start button text", "test-acc end button text"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				// Test resource create
				Config: testAccIdpSAMLv2ResourceConfig(resourceName, startButtonText, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthIdentityProviderExists(tfResourcePath),
					resource.TestCheckResourceAttr(tfResourcePath, "button_text", startButtonText),
					resource.TestCheckResourceAttr(tfResourcePath, "name", fmt.Sprintf("test-acc %s", resourceName)),
					resource.TestCheckResourceAttr(tfResourcePath, "idp_endpoint", "https://example.com/fake"),
					resource.TestCheckResourceAttr(tfResourcePath, "enabled", "true"),
					resource.TestCheckResourceAttrPair(tfResourcePath, "key_id", keyResourcePath, "id"),
				),
			},
			{
				// Test resource update/state mutate
				Config: testAccIdpSAMLv2ResourceConfig(resourceName, endButtonText, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthIdentityProviderExists(tfResourcePath),
					resource.TestCheckResourceAttr(tfResourcePath, "button_text", endButtonText),
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

// TestAccFusionauthIdpSAMLv2_issuer covers the service provider EntityId override (ENG-3491),
// FusionAuth >= 1.69.0 only. FusionAuth does not persist the generated default, so an unset issuer
// must stay absent from state rather than reading back as the generated value.
func TestAccFusionauthIdpSAMLv2_issuer(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_idp_saml_v2.test_%s", resourceName)
	buttonText := "test-acc issuer button text"
	// Deliberately not a URI: FusionAuth only requires the value be non-blank when provided.
	customIssuer := "legacy-sp-audience.example"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfFusionAuthBelow(t, "1.69.0") },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				// Empty state attributes get elided from the flatmap, so assert no attribute rather
				// than an empty TestCheckResourceAttr.
				Config: testAccIdpSAMLv2ResourceConfig(resourceName, buttonText, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthIdentityProviderExists(tfResourcePath),
					resource.TestCheckNoResourceAttr(tfResourcePath, "issuer"),
				),
			},
			{
				Config: testAccIdpSAMLv2ResourceConfig(resourceName, buttonText, customIssuer),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthIdentityProviderExists(tfResourcePath),
					resource.TestCheckResourceAttr(tfResourcePath, "issuer", customIssuer),
				),
			},
		},
	})
}

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

func testAccIdpSAMLv2ResourceConfig(resourceName, buttonText, issuer string) string {
	var issuerHCL string
	if issuer != "" {
		issuerHCL = fmt.Sprintf("\n  issuer       = %q", issuer)
	}

	return testAccKeyResourceConfig("", resourceName, fusionauth.Algorithm_RS256, 2048, "FusionAuth") + fmt.Sprintf(`
resource "fusionauth_idp_saml_v2" "test_%[1]s" {
  button_text  = %[2]q
  idp_endpoint = "https://example.com/fake"
  key_id       = fusionauth_key.test_%[1]s.id
  name         = "test-acc %[1]s"
  enabled      = true%[3]s
}
`, resourceName, buttonText, issuerHCL)
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
