package fusionauth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func handleStringSlice(key string, data *schema.ResourceData) []string {
	set := data.Get(key).(*schema.Set)
	l := set.List()
	s := make([]string, 0, len(l))

	for _, x := range l {
		s = append(s, x.(string))
	}

	return s
}

func checkResponse(statusCode int, faErrors *fusionauth.Errors) error {
	switch {
	case statusCode >= 200 && statusCode <= 299:
		return nil
	case faErrors == nil:
		return fmt.Errorf("unexpected status code: %d(%s)", statusCode, http.StatusText(statusCode))
	default:
		return fmt.Errorf("unexpected status code: %d(%s) Errors: %v", statusCode, http.StatusText(statusCode), faErrors)
	}
}

func templateCompare(k, oldStr, newStr string, d *schema.ResourceData) bool {
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

func certKeyCompare(k, oldStr, newStr string, d *schema.ResourceData) bool {
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
