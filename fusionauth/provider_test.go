package fusionauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	providerFusionauth = "fusionauth"
	retryTimeout       = 10 * time.Second
)

// testAccProviderFactories is a static map containing only the main provider instance
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	// Always allocate a new provider instance each invocation, otherwise gRPC
	// ProviderConfigure() can overwrite configuration during concurrent testing.
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		providerFusionauth: func() (*schema.Provider, error) {
			// Create a provider instance...
			p := Provider()

			// Then pump the provider with the required Terraform configuration.
			if diags := p.Configure(context.Background(), terraform.NewResourceConfigRaw(nil)); diags.HasError() {
				return nil, errors.New(fmt.Sprintf("error configuring provider: %#+v\n", diags))
			}

			return p, nil
		},
	}
}

// fusionauthClient extracts the underlying client from a configured provider
func fusionauthClient() *fusionauth.FusionAuthClient {
	provider, err := testAccProviderFactories[providerFusionauth]()
	if err != nil {
		log.Println("[ERROR] error getting Fusionauth Provider")
	}

	return provider.Meta().(Client).FAClient
}

// testAccPreCheck validates the necessary test API keys exist in the testing
// environment
func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("FA_DOMAIN"); v == "" {
		t.Fatal("FA_DOMAIN must be set for acceptance tests")
	}
	if v := os.Getenv("FA_API_KEY"); v == "" {
		t.Fatal("FA_API_KEY must be set for acceptance tests")
	}
}

// checkFusionauthErrors checks for low-level client errors and then any
// reported fusionauth errors.
func checkFusionauthErrors(faErrs *fusionauth.Errors, err error) error {
	if err != nil {
		// low-level error performing api request
		return err
	}

	if faErrs != nil && faErrs.Present() {
		// Fusionauth has errors to report
		return fmt.Errorf("fusionauth errors: %#+v\n", faErrs)
	}

	return nil
}

// checkFusionauthRetryErrors wraps checking for fusionauth or low-level client
// errors and returns a non-retryable error on failure.
func checkFusionauthRetryErrors(faErrs *fusionauth.Errors, err error) *resource.RetryError {
	if anErr := checkFusionauthErrors(faErrs, err); anErr != nil {
		return resource.NonRetryableError(anErr)
	}

	return nil
}

// stringSliceToHCL takes a string slice and marshals it to JSON in order to
// generate an HCL syntactically compatible array.
func stringSliceToHCL(values []string) string {
	output := "[]"
	if len(values) > 0 {
		bytes, _ := json.Marshal(values)
		output = string(bytes)
	}

	return output
}

// randString10 returns a random alpha-numeric string of 10 characters.
func randString10() string {
	return acctest.RandString(10)
}

// randString20 returns a random alpha-numeric string of 20 characters.
func randString20() string {
	return acctest.RandString(20)
}
