package fusionauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestNormalizeCIMDDomain(t *testing.T) {
	for in, want := range map[string]string{
		"https://app.example.com":          "https://app.example.com",
		"HTTPS://App.Example.com:443":      "https://app.example.com",
		"https://app.example.com:8443":     "https://app.example.com:8443",
		"*":                                "*",
		"https://app.example.com:4430":     "https://app.example.com:4430",
		"https://APP.example.com:443/path": "https://app.example.com:443/path",
	} {
		if got := normalizeCIMDDomain(in); got != want {
			t.Errorf("normalizeCIMDDomain(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestAccFusionauthCIMDDomain_basic covers the FusionAuth 1.70.0 CIMD domain and scope resources; skips on
// older servers.
func TestAccFusionauthCIMDDomain_basic(t *testing.T) {
	resourceName := randString10()
	domainPath := fmt.Sprintf("fusionauth_cimd_domain.test_%s", resourceName)
	scopePath := fmt.Sprintf("fusionauth_cimd_domain_oauth_scope.test_%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfFusionAuthBelow(t, "1.70.0") },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthCIMDDomainDestroy,
		Steps: []resource.TestStep{
			{
				// Minimal: omitted fields take the server defaults, and the domain's case and default port are
				// normalized by FusionAuth without a perpetual diff.
				Config: testAccCIMDDomainMinimalConfig(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthCIMDDomainExists(domainPath),
					resource.TestCheckResourceAttr(domainPath, "domain", fmt.Sprintf("https://cimd-%s.example.com", resourceName)),
					resource.TestCheckResourceAttr(domainPath, "name", fmt.Sprintf("https://cimd-%s.example.com", resourceName)),
					resource.TestCheckResourceAttr(domainPath, "consent_mode", "AlwaysPrompt"),
					resource.TestCheckResourceAttr(domainPath, "max_cache_ttl_in_seconds", "86400"),
					resource.TestCheckResourceAttr(domainPath, "grace_period_in_seconds", "0"),
					resource.TestCheckResourceAttr(domainPath, "provided_scope_policy.0.address.0.enabled", "true"),
					resource.TestCheckResourceAttr(domainPath, "provided_scope_policy.0.profile.0.enabled", "true"),
					resource.TestCheckResourceAttr(domainPath, "provided_scope_policy.0.email.0.required", "false"),
				),
			},
			{
				Config: testAccCIMDDomainFullConfig(resourceName, "read-only access"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthCIMDDomainExists(domainPath),
					resource.TestCheckResourceAttr(domainPath, "name", "test-acc-cimd "+resourceName),
					resource.TestCheckResourceAttr(domainPath, "allow_localhost_redirect_uris", "true"),
					resource.TestCheckResourceAttr(domainPath, "authorized_resource_uris.#", "2"),
					resource.TestCheckResourceAttr(domainPath, "authorized_resource_uris.0", "https://api.example.com/mcp"),
					resource.TestCheckResourceAttr(domainPath, "consent_mode", "RememberDecision"),
					resource.TestCheckResourceAttr(domainPath, "debug", "true"),
					resource.TestCheckResourceAttr(domainPath, "grace_period_in_seconds", "3600"),
					resource.TestCheckResourceAttr(domainPath, "max_cache_ttl_in_seconds", "900"),
					resource.TestCheckResourceAttrPair(domainPath, "registration_lambda_id", fmt.Sprintf("fusionauth_lambda.registration_%s", resourceName), "id"),
					resource.TestCheckResourceAttrPair(domainPath, "jwt_populate_lambda_id", fmt.Sprintf("fusionauth_lambda.populate_%s", resourceName), "id"),
					// Partial policy: the omitted address and phone scopes keep the server default of enabled.
					resource.TestCheckResourceAttr(domainPath, "provided_scope_policy.0.address.0.enabled", "true"),
					resource.TestCheckResourceAttr(domainPath, "provided_scope_policy.0.email.0.required", "true"),
					resource.TestCheckResourceAttr(domainPath, "provided_scope_policy.0.phone.0.enabled", "true"),
					resource.TestCheckResourceAttr(domainPath, "provided_scope_policy.0.profile.0.enabled", "false"),
					resource.TestCheckResourceAttrPair(scopePath, "cimd_domain_id", domainPath, "id"),
					resource.TestCheckResourceAttr(scopePath, "name", "data:read"),
					resource.TestCheckResourceAttrPair(scopePath, "scope_id", scopePath, "id"),
					resource.TestCheckResourceAttr(scopePath, "description", "read-only access"),
					resource.TestCheckResourceAttr(scopePath, "required", "true"),
				),
			},
			{
				Config: testAccCIMDDomainFullConfig(resourceName, "updated access"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(scopePath, "description", "updated access"),
				),
			},
			{
				ResourceName:      domainPath,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      scopePath,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources[scopePath]
					return rs.Primary.Attributes["cimd_domain_id"] + ":" + rs.Primary.ID, nil
				},
			},
		},
	})
}

func testAccCIMDDomainMinimalConfig(resourceName string) string {
	return fmt.Sprintf(`
resource "fusionauth_tenant" "test_%[1]s" {
  name = "test-acc-cimd %[1]s"
}

resource "fusionauth_cimd_domain" "test_%[1]s" {
  tenant_id = fusionauth_tenant.test_%[1]s.id
  domain    = "HTTPS://CIMD-%[2]s.Example.com:443"
}
`, resourceName, strings.ToUpper(resourceName))
}

func testAccCIMDDomainFullConfig(resourceName, scopeDescription string) string {
	return fmt.Sprintf(`
resource "fusionauth_tenant" "test_%[1]s" {
  name = "test-acc-cimd %[1]s"
}

resource "fusionauth_lambda" "registration_%[1]s" {
  name = "test-acc-cimd-registration %[1]s"
  type = "CimdClientRegistration"
  body = <<EOT
function evaluate(result, clientMetadata, context) {
}
EOT
}

resource "fusionauth_lambda" "populate_%[1]s" {
  name = "test-acc-cimd-populate %[1]s"
  type = "JWTPopulate"
  body = <<EOT
function populate(jwt, user, registration) {
}
EOT
}

resource "fusionauth_cimd_domain" "test_%[1]s" {
  tenant_id                     = fusionauth_tenant.test_%[1]s.id
  domain                        = "https://cimd-%[1]s.example.com"
  name                          = "test-acc-cimd %[1]s"
  allow_localhost_redirect_uris = true
  authorized_resource_uris      = ["https://api.example.com/mcp", "https://api.example.com/other"]
  consent_mode                  = "RememberDecision"
  data                          = jsonencode({ owner = "test-acc" })
  debug                         = true
  grace_period_in_seconds       = 3600
  jwt_populate_lambda_id        = fusionauth_lambda.populate_%[1]s.id
  max_cache_ttl_in_seconds      = 900
  registration_lambda_id        = fusionauth_lambda.registration_%[1]s.id

  provided_scope_policy {
    email {
      required = true
    }
    profile {
      enabled = false
    }
  }
}

resource "fusionauth_cimd_domain_oauth_scope" "test_%[1]s" {
  cimd_domain_id          = fusionauth_cimd_domain.test_%[1]s.id
  name                    = "data:read"
  description             = "%[2]s"
  default_consent_message = "View your data"
  default_consent_detail  = "Read-only access to your data"
  data                    = jsonencode({ createdBy = "test-acc" })
  required                = true
}
`, resourceName, scopeDescription)
}

func testAccCheckFusionauthCIMDDomainExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		faClient := fusionauthClient()
		res, faErrs, err := faClient.RetrieveCimdDomain(rs.Primary.ID)
		if errs := checkFusionauthErrors(faErrs, err); errs != nil {
			return errs
		}

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get resource: %#+v", res)
		}

		return nil
	}
}

func testAccCheckFusionauthCIMDDomainDestroy(s *terraform.State) error {
	faClient := fusionauthClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_cimd_domain" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			res, faErrs, err := faClient.RetrieveCimdDomain(rs.Primary.ID)
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
