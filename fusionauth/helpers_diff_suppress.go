package fusionauth

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// diffSuppressCertKey suppresses terraform reporting differences in certificate
// keys if the returned content is equivalent.
func diffSuppressCertKey(_, oldStr, newStr string, _ *schema.ResourceData) bool {
	clean := func(s string) string {
		s = strings.ReplaceAll(s, "\r\n", "\n")
		s = strings.ReplaceAll(s, "\n", "")
		s = strings.ReplaceAll(s, "-----BEGIN CERTIFICATE-----", "")
		s = strings.ReplaceAll(s, "-----END CERTIFICATE-----", "")
		s = strings.ReplaceAll(s, "-----BEGIN PUBLIC KEY-----", "")
		s = strings.ReplaceAll(s, "-----END PUBLIC KEY-----", "")
		s = strings.ReplaceAll(s, "-----BEGIN PRIVATE KEY-----", "")
		s = strings.ReplaceAll(s, "-----END PRIVATE KEY-----", "")
		return s
	}
	oldStr = clean(oldStr)
	newStr = clean(newStr)
	return oldStr == newStr
}

// diffSuppressJSON suppresses terraform reporting differences in schema if the
// returned JSON is equivalent.
func diffSuppressJSON(_, oldJSON, newJSON string, _ *schema.ResourceData) bool {
	equal, err := isEqualJSON(oldJSON, newJSON)
	if err != nil {
		return false
	}

	return equal
}

// diffSuppressTemplate suppresses terraform reporting differences in schema if
// the returned template content is equivalent.
func diffSuppressTemplate(_, oldStr, newStr string, _ *schema.ResourceData) bool {
	clean := func(s string) string {
		s = strings.ReplaceAll(s, " ", "")
		s = strings.ReplaceAll(s, "\t", "")
		s = strings.ReplaceAll(s, "\r\n", "\n")
		s = strings.ReplaceAll(s, "\n", "")
		return s
	}
	oldStr = clean(oldStr)
	newStr = clean(newStr)
	return oldStr == newStr
}
