package fusionauth

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFusionauthKey_basic(t *testing.T) {
	resourceName := randString10()
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
				Config: testAccKeyResourceConfig("", resourceName, startAlgorithm, startLength),
				Check: testKeyAccTestCheckFuncs(
					tfResourcePath,
					resourceName,
					startAlgorithm,
					startLength,
				),
			},
			{
				// Test resource update/state mutate
				Config: testAccKeyResourceConfig("", resourceName, endAlgorithm, endLength),
				Check: testKeyAccTestCheckFuncs(
					tfResourcePath,
					resourceName,
					endAlgorithm,
					endLength,
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

func TestAccFusionauthKey_SetID(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_key.test_%s", resourceName)

	id, err := uuid.GenerateUUID()
	if err != nil {
		t.Errorf("error generating uuid: %s", err)
		return
	}
	startAlgorithm, endAlgorithm := fusionauth.Algorithm_RS256, fusionauth.Algorithm_RS512
	startLength, endLength := 2048, 4096

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthKeyDestroy,
		Steps: []resource.TestStep{
			{
				// Test resource create
				Config: testAccKeyResourceConfig(id, resourceName, startAlgorithm, startLength),
				Check: testKeyAccTestCheckFuncs(
					tfResourcePath,
					resourceName,
					startAlgorithm,
					startLength,
					resource.TestCheckResourceAttr(tfResourcePath, "key_id", id),
				),
			},
			{
				// Test resource update/state mutate
				Config: testAccKeyResourceConfig(id, resourceName, endAlgorithm, endLength),
				Check: testKeyAccTestCheckFuncs(
					tfResourcePath,
					resourceName,
					endAlgorithm,
					endLength,
					resource.TestCheckResourceAttr(tfResourcePath, "key_id", id),
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

// testKeyAccTestCheckFuncs abstracts the test case setup required between
// create and update testing.
func testKeyAccTestCheckFuncs(
	tfResourcePath string,
	resourceName string,
	algorithm fusionauth.Algorithm,
	length int,
	extraFuncs ...resource.TestCheckFunc,
) resource.TestCheckFunc {
	testFuncs := []resource.TestCheckFunc{
		testAccCheckFusionauthKeyExists(tfResourcePath),
		resource.TestCheckResourceAttrSet(tfResourcePath, "key_id"),
		resource.TestCheckResourceAttr(tfResourcePath, "name", fmt.Sprintf("test-acc %s", resourceName)),
		resource.TestCheckResourceAttr(tfResourcePath, "algorithm", string(algorithm)),
		resource.TestCheckResourceAttr(tfResourcePath, "length", fmt.Sprintf("%d", length)),
		resource.TestCheckResourceAttrSet(tfResourcePath, "kid"),
	}

	if len(extraFuncs) > 0 {
		testFuncs = append(testFuncs, extraFuncs...)
	}

	return resource.ComposeTestCheckFunc(testFuncs...)
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

		faClient := fusionauthClient()
		key, faErrs, err := faClient.RetrieveKey(rs.Primary.ID)
		if errs := checkFusionauthErrors(faErrs, err); errs != nil {
			return err
		}

		if key == nil || key.StatusCode != http.StatusOK {
			// This is a weird edge case...
			return fmt.Errorf("failed to get resource: %#+v", key)
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

func testAccKeyResourceConfig(id string, name string, algorithm fusionauth.Algorithm, length int) string {
	var keyID string
	if id != "" {
		keyID = fmt.Sprintf("\n  key_id    = \"%s\"\n", id)
	}

	return fmt.Sprintf(`
# Key Setup
resource "fusionauth_key" "test_%[2]s" {%[1]s
  name      = "test-acc %[2]s"
  algorithm = "%[3]s"
  length    = %[4]d
}
`, keyID, name, algorithm, length)
}

const (
	testAccessTokenKey = "accesstoken"
	testIDTokenKey     = "idtoken"
)

// testAccAccessTokenKeyResourceConfig returns terraform configuration to
// generate a standalone test Access Token key.
func testAccAccessTokenKeyResourceConfig(suffix string) string {
	return testKeyConfig(testAccessTokenKey, suffix)
}

// testAccIDTokenKeyResourceConfig returns terraform configuration to generate a
// standalone test ID Token key.
func testAccIDTokenKeyResourceConfig(suffix string) string {
	return testKeyConfig(testIDTokenKey, suffix)
}

func testKeyConfig(name, suffix string) string {
	return testAccKeyResourceConfig(
		"",
		testKeyName(name, suffix),
		fusionauth.Algorithm_RS256,
		2048,
	)
}

func testKeyName(name, suffix string) string {
	if suffix != "" {
		name = fmt.Sprintf("%s_%s", name, suffix)
	}

	return name
}
