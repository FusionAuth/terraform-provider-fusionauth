package fusionauth

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Test_buildApplication_baseURL(t *testing.T) {
	tests := []struct {
		name string
		raw  map[string]interface{}
		want string
	}{
		{
			name: "set",
			raw:  map[string]interface{}{"base_url": "https://auth.example.com"},
			want: "https://auth.example.com",
		},
		{
			// Unset maps to "" so servers that omit baseURL don't perpetually diff.
			name: "unset",
			raw:  map[string]interface{}{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := schema.TestResourceDataRaw(t, newApplication().Schema, tt.raw)
			if got := buildApplication(data).BaseURL; got != tt.want {
				t.Fatalf("BaseURL = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_resolvePasswordlessEnabledState(t *testing.T) {
	tests := []struct {
		name               string
		apiEnabled         bool
		priorLegacyEnabled bool
		priorBlockEnabled  bool
		wantLegacyEnabled  bool
		wantBlockEnabled   bool
	}{
		{
			name:               "legacy alias preserves disabled block when api is enabled",
			apiEnabled:         true,
			priorLegacyEnabled: true,
			priorBlockEnabled:  false,
			wantLegacyEnabled:  true,
			wantBlockEnabled:   false,
		},
		{
			name:               "nested block preserves disabled legacy alias when api is enabled",
			apiEnabled:         true,
			priorLegacyEnabled: false,
			priorBlockEnabled:  true,
			wantLegacyEnabled:  false,
			wantBlockEnabled:   true,
		},
		{
			name:               "both remain enabled when both were configured",
			apiEnabled:         true,
			priorLegacyEnabled: true,
			priorBlockEnabled:  true,
			wantLegacyEnabled:  true,
			wantBlockEnabled:   true,
		},
		{
			name:               "external enablement surfaces drift on both aliases",
			apiEnabled:         true,
			priorLegacyEnabled: false,
			priorBlockEnabled:  false,
			wantLegacyEnabled:  true,
			wantBlockEnabled:   true,
		},
		{
			name:               "api disabled clears both aliases",
			apiEnabled:         false,
			priorLegacyEnabled: true,
			priorBlockEnabled:  true,
			wantLegacyEnabled:  false,
			wantBlockEnabled:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLegacyEnabled, gotBlockEnabled := resolvePasswordlessEnabledState(
				tt.apiEnabled,
				tt.priorLegacyEnabled,
				tt.priorBlockEnabled,
			)

			if gotLegacyEnabled != tt.wantLegacyEnabled {
				t.Fatalf("legacy enabled = %v, want %v", gotLegacyEnabled, tt.wantLegacyEnabled)
			}
			if gotBlockEnabled != tt.wantBlockEnabled {
				t.Fatalf("block enabled = %v, want %v", gotBlockEnabled, tt.wantBlockEnabled)
			}
		})
	}
}
