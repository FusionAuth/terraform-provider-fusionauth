package fusionauth

import (
	"strings"
	"testing"

	fa "github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const linkingStrategyAttr = "linking_strategy"

// idpResourceCount stops the resource walk below from passing vacuously.
const idpResourceCount = 12

// TestIDPLinkingStrategies_MatchesSDKEnum pins the allow-list to the go-client enum.
func TestIDPLinkingStrategies_MatchesSDKEnum(t *testing.T) {
	want := []string{
		string(fa.IdentityProviderLinkingStrategy_CreatePendingLink),
		string(fa.IdentityProviderLinkingStrategy_Disabled),
		string(fa.IdentityProviderLinkingStrategy_LinkAnonymously),
		string(fa.IdentityProviderLinkingStrategy_LinkByEmail),
		string(fa.IdentityProviderLinkingStrategy_LinkByEmailForExistingUser),
		string(fa.IdentityProviderLinkingStrategy_LinkByUsername),
		string(fa.IdentityProviderLinkingStrategy_LinkByUsernameForExistingUser),
	}

	strategies := idpLinkingStrategies()

	if len(strategies) != len(want) {
		t.Fatalf("idpLinkingStrategies() = %v, want %v", strategies, want)
	}
	for i, w := range want {
		if got := strategies[i]; got != w {
			t.Errorf("idpLinkingStrategies()[%d] = %q, want %q", i, got, w)
		}
	}

	unsupported := string(fa.IdentityProviderLinkingStrategy_Unsupported)
	for _, got := range strategies {
		if got == unsupported {
			t.Errorf("idpLinkingStrategies() contains %q, which the API always rejects", unsupported)
		}
	}
}

// TestIDPResourcesShareLinkingStrategySchema checks the shared linking_strategy attribute on every
// fusionauth_idp_* resource.
func TestIDPResourcesShareLinkingStrategySchema(t *testing.T) {
	seen := 0
	for name, resource := range Provider().ResourcesMap {
		if !strings.HasPrefix(name, "fusionauth_idp_") {
			continue
		}
		seen++

		t.Run(name, func(t *testing.T) {
			attr, ok := resource.Schema[linkingStrategyAttr]
			if !ok {
				t.Fatalf("%s has no %s attribute", name, linkingStrategyAttr)
			}

			assertLinkingStrategyValidation(t, name, attr)
			assertLinkingStrategyDescription(t, name, attr.Description)
		})
	}

	if seen != idpResourceCount {
		t.Errorf("checked %d fusionauth_idp_* resources, want %d", seen, idpResourceCount)
	}
}

func assertLinkingStrategyValidation(t *testing.T, name string, attr *schema.Schema) {
	t.Helper()

	// Optional without Computed on a server defaulted attribute gives a perpetual diff.
	if !attr.Optional || !attr.Computed {
		t.Errorf("%s: Optional = %v, Computed = %v, want both true", name, attr.Optional, attr.Computed)
	}
	if attr.ValidateFunc == nil {
		t.Fatalf("%s: %s has no ValidateFunc", name, linkingStrategyAttr)
	}

	for _, strategy := range idpLinkingStrategies() {
		if _, errs := attr.ValidateFunc(strategy, linkingStrategyAttr); len(errs) > 0 {
			t.Errorf("%s: %q rejected, want accepted: %v", name, strategy, errs)
		}
	}

	rejected := []string{string(fa.IdentityProviderLinkingStrategy_Unsupported), "Bogus", "disabled"}
	for _, strategy := range rejected {
		if _, errs := attr.ValidateFunc(strategy, linkingStrategyAttr); len(errs) == 0 {
			t.Errorf("%s: %q accepted, want rejected", name, strategy)
		}
	}
}

func assertLinkingStrategyDescription(t *testing.T, name, desc string) {
	t.Helper()

	const prefix = "The linking strategy to use when creating the link between the "
	const suffix = " Identity Provider and the user."

	if !strings.HasPrefix(desc, prefix) || !strings.HasSuffix(desc, suffix) {
		t.Fatalf("%s: description = %q, want %q<display name>%q", name, desc, prefix, suffix)
	}

	// An unrendered display name reaches the registry docs as a stray placeholder.
	switch displayName := strings.TrimSuffix(strings.TrimPrefix(desc, prefix), suffix); {
	case displayName == "":
		t.Errorf("%s: description has an empty display name: %q", name, desc)
	case strings.ContainsAny(displayName, "{}"):
		t.Errorf("%s: description display name is unrendered: %q", name, displayName)
	}
}
