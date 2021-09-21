package fusionauth

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// WarnStringInSlice returns a SchemaValidateFunc which tests if the provided
// value is of type string and will not error, but rather, return a warning if
// the string provided is not in the valid slice. This enables soft-constraints,
// where a plugin may expose extra values over and above the core values.
// Will test using case-insensitivity if ignoreCase is true.
func WarnStringInSlice(valid []string, ignoreCase bool) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) (diags diag.Diagnostics) {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{{
				Severity:      diag.Error,
				Summary:       "expected type of string",
				Detail:        fmt.Sprintf("expected string, got: %+#v instead", value),
				AttributePath: path,
			}}
		}

		for _, validator := range valid {
			if value == validator || (ignoreCase && strings.EqualFold(value, validator)) {
				return
			}
		}

		return diag.Diagnostics{{
			Severity:      diag.Warning,
			Summary:       "provided value not in expected list (may be expected)",
			Detail:        fmt.Sprintf("expected the provided value to be one of %s, got %s", valid, value),
			AttributePath: path,
		}}
	}
}
