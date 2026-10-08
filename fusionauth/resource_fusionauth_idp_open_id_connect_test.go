package fusionauth

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func Test_buildResourceFromOpenIDConnect_clientSecret(t *testing.T) {
	tests := []struct {
		name string
		raw  map[string]interface{}
		want string
	}{
		{
			name: "oauth2_client_secret",
			raw:  map[string]interface{}{},
			want: "secret",
		},
		{
			// The secret is managed through oauth2_client_secret_wo, so it must stay out of state.
			name: "oauth2_client_secret_wo",
			raw:  map[string]interface{}{"oauth2_client_secret_wo_version": 1},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := schema.TestResourceDataRaw(t, newIDPOpenIDConnect().Schema, tt.raw)
			o := fusionauth.OpenIdConnectIdentityProvider{
				Oauth2: fusionauth.IdentityProviderOauth2Configuration{ClientSecret: "secret"},
			}
			if diags := buildResourceFromOpenIDConnect(o, data); diags.HasError() {
				t.Fatalf("buildResourceFromOpenIDConnect() = %v", diags)
			}
			if got := data.Get("oauth2_client_secret").(string); got != tt.want {
				t.Fatalf("oauth2_client_secret = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAccFusionauthIdpOpenIDConnect_clientSecretWriteOnly(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_idp_open_id_connect.test_%s", resourceName)
	startSecret := randString20()
	endSecret := randString20()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				// Test resource create
				Config: testAccIdpOpenIDConnectClientSecretWriteOnlyConfig(resourceName, "Login", startSecret, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthIdentityProviderExists(tfResourcePath),
					testAccCheckFusionauthIdpOpenIDConnectClientSecret(tfResourcePath, startSecret),
					resource.TestCheckNoResourceAttr(tfResourcePath, "oauth2_client_secret_wo"),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_client_secret_wo_version", "1"),
					resource.TestCheckNoResourceAttr(tfResourcePath, "oauth2_client_secret"),
				),
			},
			{
				// Identity provider writes don't write the response back to state, so refresh
				// to check the secret read from FusionAuth is kept out of state.
				RefreshState: true,
				Check:        resource.TestCheckResourceAttr(tfResourcePath, "oauth2_client_secret", ""),
			},
			{
				// Test resource update without a version change keeps the secret
				Config: testAccIdpOpenIDConnectClientSecretWriteOnlyConfig(resourceName, "Login with OIDC", startSecret, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthIdpOpenIDConnectClientSecret(tfResourcePath, startSecret),
					resource.TestCheckResourceAttr(tfResourcePath, "button_text", "Login with OIDC"),
				),
			},
			{
				// Test secret rotation
				Config: testAccIdpOpenIDConnectClientSecretWriteOnlyConfig(resourceName, "Login with OIDC", endSecret, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthIdpOpenIDConnectClientSecret(tfResourcePath, endSecret),
					resource.TestCheckResourceAttr(tfResourcePath, "oauth2_client_secret_wo_version", "2"),
				),
			},
			{
				// Test importing resource into state
				ResourceName:            tfResourcePath,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth2_client_secret", "oauth2_client_secret_wo_version"},
			},
		},
	})
}

// testAccCheckFusionauthIdpOpenIDConnectClientSecret compares the client secret held by
// FusionAuth, for when it is kept out of state.
func testAccCheckFusionauthIdpOpenIDConnectClientSecret(resourceName, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		b, err := readIdentityProvider(rs.Primary.ID, fusionauthProviderClient())
		if err != nil {
			return err
		}

		var ipb OpenIDConnectIdentityProviderBody
		if err := json.Unmarshal(b, &ipb); err != nil {
			return err
		}

		if got := ipb.IdentityProvider.Oauth2.ClientSecret; got != want {
			return fmt.Errorf("client secret = %q, want %q", got, want)
		}

		return nil
	}
}

func testAccIdpOpenIDConnectClientSecretWriteOnlyConfig(resourceName, buttonText, clientSecret string, version int) string {
	return fmt.Sprintf(`
resource "fusionauth_idp_open_id_connect" "test_%[1]s" {
  name                            = "test-acc-oidc %[1]s"
  button_text                     = "%[2]s"
  oauth2_authorization_endpoint   = "https://example.com/oauth2/authorize"
  oauth2_token_endpoint           = "https://example.com/oauth2/token"
  oauth2_user_info_endpoint       = "https://example.com/oauth2/userinfo"
  oauth2_client_id                = "%[1]s"
  oauth2_client_secret_wo         = "%[3]s"
  oauth2_client_secret_wo_version = %[4]d
  oauth2_scope                    = "openid"
}
`, resourceName, buttonText, clientSecret, version)
}
