package fusionauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// isEqualJSON checks to see if the two JSON encoded strings are equal.
func isEqualJSON(a, b string) (equal bool, err error) {
	var x interface{}
	if err = json.Unmarshal([]byte(a), &x); err != nil {
		return false, fmt.Errorf("error unmarshaling %s to JSON: %s", a, err)
	}

	var y interface{}
	if err = json.Unmarshal([]byte(b), &y); err != nil {
		return false, fmt.Errorf("error unmarshaling %s to JSON: %s", b, err)
	}

	return reflect.DeepEqual(x, y), nil
}

// injectSchemaChanges pushes the provided schema edits into the provided schema.
// Primarily for working with multiple schema versions.
func injectSchemaChanges(schemaToEdit, schemaEdits *schema.Resource) *schema.Resource {
	for attributeName, schemaEdit := range schemaEdits.Schema {
		schemaToEdit.Schema[attributeName] = schemaEdit
	}

	return schemaToEdit
}
