package fusionauth

import (
	"testing"

	fa "github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestBuildTenantManagerIdentityProviderTypeConfigurationRequest_DefaultDisabledLinkingStrategy(t *testing.T) {
	resourceSchema := resourceTenantManagerConfiguration().Schema
	data := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		"identity_provider_type_configurations": []interface{}{
			map[string]interface{}{
				"type":    "SAMLv2",
				"enabled": true,
			},
		},
	})

	idpConfs := data.Get("identity_provider_type_configurations").(*schema.Set).List()
	req := buildTenantManagerIdentityProviderTypeConfigurationRequest(idpConfs[0].(map[string]interface{}))

	if got, want := req.TypeConfiguration.Type, fa.IdentityProviderType("SAMLv2"); got != want {
		t.Fatalf("type = %q, want %q", got, want)
	}
	if got, want := req.TypeConfiguration.LinkingStrategy, fa.IdentityProviderLinkingStrategy("LinkByEmail"); got != want {
		t.Fatalf("linking strategy = %q, want %q", got, want)
	}
	if got := req.TypeConfiguration.Enabled; !got {
		t.Fatalf("enabled = %v, want true", got)
	}
}

func TestBuildTenantManagerIdentityProviderTypeConfigurationRequest_IgnoresEmptySetItemStrings(t *testing.T) {
	req := buildTenantManagerIdentityProviderTypeConfigurationRequest(map[string]interface{}{
		"type":                       "SAMLv2",
		"enabled":                    true,
		"linking_strategy":           "LinkByEmail",
		"default_attribute_mappings": map[string]interface{}{"user.email": "emailAddress"},
	})

	if got, want := req.TypeConfiguration.Type, fa.IdentityProviderType("SAMLv2"); got != want {
		t.Fatalf("type = %q, want %q", got, want)
	}
	if got, want := req.TypeConfiguration.LinkingStrategy, fa.IdentityProviderLinkingStrategy("LinkByEmail"); got != want {
		t.Fatalf("linking strategy = %q, want %q", got, want)
	}
	if got, want := req.TypeConfiguration.DefaultAttributeMappings["user.email"], "emailAddress"; got != want {
		t.Fatalf("default attribute mapping = %q, want %q", got, want)
	}
	if got := req.TypeConfiguration.Enabled; !got {
		t.Fatalf("enabled = %v, want true", got)
	}
}

func TestGetTenantManagerIdentityProviderTypeConfigurationItems_EmptyWhenUnset(t *testing.T) {
	data := schema.TestResourceDataRaw(t, resourceTenantManagerConfiguration().Schema, map[string]interface{}{})

	items := getTenantManagerIdentityProviderTypeConfigurationItems(data)

	if len(items) != 0 {
		t.Fatalf("len(items) = %d, want 0", len(items))
	}
}

func TestGetTenantManagerIdentityProviderTypeConfigurationItems_SkipsEmptyType(t *testing.T) {
	items := getTenantManagerIdentityProviderTypeConfigurationItems(schema.TestResourceDataRaw(
		t,
		resourceTenantManagerConfiguration().Schema,
		map[string]interface{}{
			"identity_provider_type_configurations": []interface{}{
				map[string]interface{}{
					"type":    "SAMLv2",
					"enabled": true,
				},
				map[string]interface{}{
					"type":    "",
					"enabled": true,
				},
			},
		},
	))

	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if got, want := items[0]["type"], "SAMLv2"; got != want {
		t.Fatalf("type = %v, want %q", got, want)
	}
}

func TestGetTenantManagerIdentityProviderTypesToDelete_UsesCurrentConfiguration(t *testing.T) {
	current := map[string]fa.TenantManagerIdentityProviderTypeConfiguration{
		"SAMLv2": {
			Type: fa.IdentityProviderType("SAMLv2"),
		},
		"Google": {
			Type: fa.IdentityProviderType("Google"),
		},
	}
	desired := map[string]struct{}{
		"Google": {},
	}

	toDelete := getTenantManagerIdentityProviderTypesToDelete(current, desired)

	if len(toDelete) != 1 {
		t.Fatalf("len(toDelete) = %d, want 1", len(toDelete))
	}
	if got, want := toDelete[0], fa.IdentityProviderType("SAMLv2"); got != want {
		t.Fatalf("toDelete[0] = %q, want %q", got, want)
	}
}

func TestGetTenantManagerIdentityProviderTypesToDelete_UsesConfigTypeWhenMapKeyEmpty(t *testing.T) {
	current := map[string]fa.TenantManagerIdentityProviderTypeConfiguration{
		"": {
			Type: fa.IdentityProviderType("SAMLv2"),
		},
	}

	toDelete := getTenantManagerIdentityProviderTypesToDelete(current, map[string]struct{}{})

	if len(toDelete) != 1 {
		t.Fatalf("len(toDelete) = %d, want 1", len(toDelete))
	}
	if got, want := toDelete[0], fa.IdentityProviderType("SAMLv2"); got != want {
		t.Fatalf("toDelete[0] = %q, want %q", got, want)
	}
}

func TestGetTenantManagerIdentityProviderTypesToDelete_DeletesAllWhenDesiredEmpty(t *testing.T) {
	current := map[string]fa.TenantManagerIdentityProviderTypeConfiguration{
		"SAMLv2": {
			Type: fa.IdentityProviderType("SAMLv2"),
		},
		"OpenIDConnect": {
			Type: fa.IdentityProviderType("OpenIDConnect"),
		},
	}

	toDelete := getTenantManagerIdentityProviderTypesToDelete(current, map[string]struct{}{})

	if len(toDelete) != 2 {
		t.Fatalf("len(toDelete) = %d, want 2", len(toDelete))
	}
	got := map[fa.IdentityProviderType]struct{}{}
	for _, idpType := range toDelete {
		got[idpType] = struct{}{}
	}
	for _, want := range []fa.IdentityProviderType{"SAMLv2", "OpenIDConnect"} {
		if _, ok := got[want]; !ok {
			t.Fatalf("missing %q in toDelete", want)
		}
	}
}

func TestBuildResourceFromTenantManagerConfiguration_UsesConfigTypeWhenMapKeyEmpty(t *testing.T) {
	data := schema.TestResourceDataRaw(t, resourceTenantManagerConfiguration().Schema, map[string]interface{}{})
	tmc := fa.TenantManagerConfiguration{
		IdentityProviderTypeConfigurations: map[string]fa.TenantManagerIdentityProviderTypeConfiguration{
			"": {
				Enableable: fa.Enableable{Enabled: true},
				Type:       fa.IdentityProviderType("SAMLv2"),
			},
		},
	}

	if diags := buildResourceFromTenantManagerConfiguration(tmc, data); diags.HasError() {
		t.Fatalf("buildResourceFromTenantManagerConfiguration returned errors: %v", diags)
	}

	items := getTenantManagerIdentityProviderTypeConfigurationItems(data)
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if got, want := items[0]["type"], "SAMLv2"; got != want {
		t.Fatalf("type = %v, want %q", got, want)
	}
}
