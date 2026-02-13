package fusionauth

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"text/template"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testAccIdpSamlV2ApplicationConfigurationIdpTFResourcePath  = "fusionauth_idp_saml_v2.test"
	testAccIdpSamlV2ApplicationConfigurationApp0TFResourcePath = "fusionauth_application.test_0"
	testAccIdpSamlV2ApplicationConfigurationApp1TFResourcePath = "fusionauth_application.test_1"
)

func TestAccIdpSamlV2ApplicationConfiguration(t *testing.T) {
	resourceName := randString10()

	// We do not include a CheckDestroy in our tests because destruction of the
	// fusionauth_idp_saml_v2_application_configuration is tested in the happy
	// path test case and destruction of the other resources set up for the
	// test case are tested in acceptance tests for those resources.

	// Happy path.
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			// Create an IdP with no application configurations.
			{
				Config: testAccIDPSAMLv2ApplicationConfigurationConfig(t,
					&testAccIdpSamlV2ApplicationConfigurationResourceTemplateConfig{
						ResourceName: resourceName,
					},
				),
				Check: testAccIDPSAMLv2ApplicationConfigurationCompareAppConfigs(map[string]SAMLAppConfig{}),
			},
			// Add an application configuration.
			{
				Config: testAccIDPSAMLv2ApplicationConfigurationConfig(t,
					&testAccIdpSamlV2ApplicationConfigurationResourceTemplateConfig{
						ResourceName:         resourceName,
						StandaloneApp0Config: true,
					},
				),
				Check: testAccIDPSAMLv2ApplicationConfigurationCompareAppConfigs(map[string]SAMLAppConfig{
					testAccIdpSamlV2ApplicationConfigurationApp0TFResourcePath: {
						ButtonText:         "Login with SAML (test_" + resourceName + "_0)",
						CreateRegistration: true,
						Enabled:            true,
					},
				}),
			},
			// Add an additional application configuration.
			{
				Config: testAccIDPSAMLv2ApplicationConfigurationConfig(t,
					&testAccIdpSamlV2ApplicationConfigurationResourceTemplateConfig{
						ResourceName:         resourceName,
						StandaloneApp0Config: true,
						StandaloneApp1Config: true,
					},
				),
				Check: testAccIDPSAMLv2ApplicationConfigurationCompareAppConfigs(map[string]SAMLAppConfig{
					testAccIdpSamlV2ApplicationConfigurationApp0TFResourcePath: {
						ButtonText:         "Login with SAML (test_" + resourceName + "_0)",
						CreateRegistration: true,
						Enabled:            true,
					},
					testAccIdpSamlV2ApplicationConfigurationApp1TFResourcePath: {
						ButtonText:         "Login with SAML (test_" + resourceName + "_1)",
						CreateRegistration: true,
						Enabled:            true,
					},
				}),
			},
			// Remove the first application configuration.
			{
				Config: testAccIDPSAMLv2ApplicationConfigurationConfig(t,
					&testAccIdpSamlV2ApplicationConfigurationResourceTemplateConfig{
						ResourceName:         resourceName,
						StandaloneApp1Config: true,
					},
				),
				Check: testAccIDPSAMLv2ApplicationConfigurationCompareAppConfigs(map[string]SAMLAppConfig{
					testAccIdpSamlV2ApplicationConfigurationApp1TFResourcePath: {
						ButtonText:         "Login with SAML (test_" + resourceName + "_1)",
						CreateRegistration: true,
						Enabled:            true,
					},
				}),
			},
		},
	})

	// Verify that import is required.
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			// Create an IdP with an inline application configuration.
			{
				Config: testAccIDPSAMLv2ApplicationConfigurationConfig(t,
					&testAccIdpSamlV2ApplicationConfigurationResourceTemplateConfig{
						ResourceName:     resourceName,
						InlineApp0Config: true,
					},
				),
				Check: testAccIDPSAMLv2ApplicationConfigurationCompareAppConfigs(map[string]SAMLAppConfig{
					testAccIdpSamlV2ApplicationConfigurationApp0TFResourcePath: {
						ButtonText:         "Login with SAML (test_" + resourceName + "_0)",
						CreateRegistration: true,
						Enabled:            true,
					},
				}),
			},
			// Try to add a standalone application configuration for the same app.
			// We expect an error that the configuration already exists.
			{
				Config: testAccIDPSAMLv2ApplicationConfigurationConfig(t,
					&testAccIdpSamlV2ApplicationConfigurationResourceTemplateConfig{
						ResourceName:         resourceName,
						InlineApp0Config:     true,
						StandaloneApp0Config: true,
					},
				),
				ExpectError: regexp.MustCompile(fmt.Sprintf(appConfigAlreadyExistsErrorFmt, `[\w-]+`, `[\w-]+`)),
			},
		},
	})

	// Verify that updates to IdP don't overwrite application configuration.
	//
	// XXX: This test case is currently flaky with multiple standalone configurations due to a race
	// condition in the PATCH API for IdPs. Confirmed with both JSON patch and JSON merge patch.
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			// Create an IdP and standalone application configurations.
			{
				Config: testAccIDPSAMLv2ApplicationConfigurationConfig(t,
					&testAccIdpSamlV2ApplicationConfigurationResourceTemplateConfig{
						ResourceName:         resourceName,
						StandaloneApp0Config: true,
						// StandaloneApp1Config: true,
					},
				),
				Check: testAccIDPSAMLv2ApplicationConfigurationCompareAppConfigs(map[string]SAMLAppConfig{
					testAccIdpSamlV2ApplicationConfigurationApp0TFResourcePath: {
						ButtonText:         "Login with SAML (test_" + resourceName + "_0)",
						CreateRegistration: true,
						Enabled:            true,
					},
					// testAccIdpSamlV2ApplicationConfigurationApp1TFResourcePath: {
					// 	ButtonText:         "Login with SAML (test_" + resourceName + "_1)",
					// 	CreateRegistration: true,
					// 	Enabled:            true,
					// },
				}),
			},
			// Make a change to the IdP.
			// We expect the application configurations to be unchanged.
			{
				Config: testAccIDPSAMLv2ApplicationConfigurationConfig(t,
					&testAccIdpSamlV2ApplicationConfigurationResourceTemplateConfig{
						ResourceName:         resourceName,
						DisableIdp:           true,
						StandaloneApp0Config: true,
						// StandaloneApp1Config: true,
					},
				),
				Check: testAccIDPSAMLv2ApplicationConfigurationCompareAppConfigs(map[string]SAMLAppConfig{
					testAccIdpSamlV2ApplicationConfigurationApp0TFResourcePath: {
						ButtonText:         "Login with SAML (test_" + resourceName + "_0)",
						CreateRegistration: true,
						Enabled:            true,
					},
					// testAccIdpSamlV2ApplicationConfigurationApp1TFResourcePath: {
					// 	ButtonText:         "Login with SAML (test_" + resourceName + "_1)",
					// 	CreateRegistration: true,
					// 	Enabled:            true,
					// },
				}),
			},
		},
	})

	// Verify that updates to IdP don't overwrite application configuration after application
	// configuration is moved from inline management to standalone management. This slightly
	// changes the behavior of the Terraform SDK.
	//
	// XXX: This test case is currently flaky with multiple standalone configurations due to a race
	// condition in the PATCH API for IdPs. Confirmed with both JSON patch and JSON merge patch.
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			// Create an IdP with inline application configuration.
			{
				Config: testAccIDPSAMLv2ApplicationConfigurationConfig(t,
					&testAccIdpSamlV2ApplicationConfigurationResourceTemplateConfig{
						ResourceName:     resourceName,
						InlineApp0Config: true,
					},
				),
				Check: testAccIDPSAMLv2ApplicationConfigurationCompareAppConfigs(map[string]SAMLAppConfig{
					testAccIdpSamlV2ApplicationConfigurationApp0TFResourcePath: {
						ButtonText:         "Login with SAML (test_" + resourceName + "_0)",
						CreateRegistration: true,
						Enabled:            true,
					},
				}),
			},
			// Remove the inline application configuration.
			{
				Config: testAccIDPSAMLv2ApplicationConfigurationConfig(t,
					&testAccIdpSamlV2ApplicationConfigurationResourceTemplateConfig{
						ResourceName: resourceName,
					},
				),
				Check: testAccIDPSAMLv2ApplicationConfigurationCompareAppConfigs(map[string]SAMLAppConfig{}),
			},
			// Add standlone application configurations.
			{
				Config: testAccIDPSAMLv2ApplicationConfigurationConfig(t,
					&testAccIdpSamlV2ApplicationConfigurationResourceTemplateConfig{
						ResourceName:         resourceName,
						StandaloneApp0Config: true,
						// StandaloneApp1Config: true,
					},
				),
				Check: testAccIDPSAMLv2ApplicationConfigurationCompareAppConfigs(map[string]SAMLAppConfig{
					testAccIdpSamlV2ApplicationConfigurationApp0TFResourcePath: {
						ButtonText:         "Login with SAML (test_" + resourceName + "_0)",
						CreateRegistration: true,
						Enabled:            true,
					},
					// testAccIdpSamlV2ApplicationConfigurationApp1TFResourcePath: {
					// 	ButtonText:         "Login with SAML (test_" + resourceName + "_1)",
					// 	CreateRegistration: true,
					// 	Enabled:            true,
					// },
				}),
			},
			// Make a change to the IdP.
			// We expect the application configurations to be unchanged.
			{
				Config: testAccIDPSAMLv2ApplicationConfigurationConfig(t,
					&testAccIdpSamlV2ApplicationConfigurationResourceTemplateConfig{
						ResourceName:         resourceName,
						DisableIdp:           true,
						StandaloneApp0Config: true,
						// StandaloneApp1Config: true,
					},
				),
				Check: testAccIDPSAMLv2ApplicationConfigurationCompareAppConfigs(map[string]SAMLAppConfig{
					testAccIdpSamlV2ApplicationConfigurationApp0TFResourcePath: {
						ButtonText:         "Login with SAML (test_" + resourceName + "_0)",
						CreateRegistration: true,
						Enabled:            true,
					},
					// testAccIdpSamlV2ApplicationConfigurationApp1TFResourcePath: {
					// 	ButtonText:         "Login with SAML (test_" + resourceName + "_1)",
					// 	CreateRegistration: true,
					// 	Enabled:            true,
					// },
				}),
			},
		},
	})
}

// testAccIdpSamlV2GetApplicationConfigurations fetches the Application Configurations
// from a FusionAuth SAML v2 IdP by id.
//
// The FusionAuthClient does not have methods that do this natively in a convenient way.
func testAccIdpSamlV2GetApplicationConfigurations(id string) (map[string]*SAMLAppConfig, error) {
	client := fusionauthClient()
	// This method of using the client to build and execute REST requests is not well documented,
	// but from reading source code it takes care of authentication and error handling and spares
	// us setting up our own HTTP client.
	// https://pkg.go.dev/github.com/FusionAuth/go-client@v1.59.0/pkg/fusionauth#FusionAuthClient.Start
	var errors fusionauth.Errors
	var resp samlV2IDPApplicationConfiguration
	err := client.Start(&resp, &errors).
		WithUri(fmt.Sprintf("/api/identity-provider/%s", id)).
		WithMethod(http.MethodGet).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	if errors.Present() {
		return nil, fmt.Errorf("errors returned from FusionAuth API: %w", errors)
	}
	return resp.IdentityProvider.ApplicationConfiguration, nil
}

// testAccIdpSamlV2ApplicationConfigurationResourceTemplateConfig provides configuration
// to the testdata/fusionauth_idp_samlv2_application_configuration.tf.tmpl template.
type testAccIdpSamlV2ApplicationConfigurationResourceTemplateConfig struct {
	// ResourceName is the name used in generation of all test resources. Required.
	ResourceName string

	// DisableIdp sets the enabled field on the SAML v2 IdP to false.
	DisableIdp bool

	// InlineApp0Config enables inline application_configuration for the first test app.
	InlineApp0Config bool

	// StandaloneApp0Config enables a standalone fusionauth_idp_saml_v2_application_configuration
	// resource for the first test app.
	StandaloneApp0Config bool

	// StandaloneApp1Config enables a standalone fusionauth_idp_saml_v2_application_configuration
	// resource for the second test app.
	StandaloneApp1Config bool
}

// testAccIDPSAMLv2ApplicationConfigurationConfig invokes the
// testdata/fusionauth_idp_samlv2_application_configuration.tf.tmpl with the provided configuration
// to generate test config. If errors are encountered the test is failed.
func testAccIDPSAMLv2ApplicationConfigurationConfig(t *testing.T,
	config *testAccIdpSamlV2ApplicationConfigurationResourceTemplateConfig) string {

	t.Helper()

	tmpl, err := template.ParseFiles("testdata/fusionauth_idp_saml_v2_application_configuration.tf.tmpl")
	if err != nil {
		t.Fatalf("unable to parse template testdata/fusionauth_idp_saml_v2_application_configuration.tf.tmpl: %s", err)
	}

	b := strings.Builder{}
	err = tmpl.Execute(&b, config)
	if err != nil {
		t.Fatalf("unable to execute template testdata/fusionauth_idp_saml_v2_application_configuration.tf.tmpl: %s", err)
	}

	return b.String()
}

// testAccIDPSAMLv2ApplicationConfigurationCompareAppConfigs returns a resource.TestCheckFunc that
// compares the application configuration of the SAML v2 IdP to the provided configurations.
//
// The keys in expected are Terraform fusionauth_application resource paths, not application IDs.
// This function uses those paths to look up the application IDs.
func testAccIDPSAMLv2ApplicationConfigurationCompareAppConfigs(expected map[string]SAMLAppConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		actual, err := testAccIdpSamlV2GetApplicationConfigurations(
			s.RootModule().Resources[testAccIdpSamlV2ApplicationConfigurationIdpTFResourcePath].Primary.ID)
		if err != nil {
			return err
		}

		if len(actual) != len(expected) {
			return fmt.Errorf("expected %d application configurations, got %d", len(expected), len(actual))
		}

		for tfPath, expectedConfig := range expected {
			appID := s.RootModule().Resources[tfPath].Primary.ID
			config, ok := actual[appID]
			if !ok {
				return fmt.Errorf("expected application configuration for app %s to exist", appID)
			}
			if *config != expectedConfig {
				return fmt.Errorf("expected application configuration for app %s to be %#v, got %#v", appID, expectedConfig, config)
			}
		}

		return nil
	}
}
