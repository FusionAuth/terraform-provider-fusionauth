package fusionauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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

// jsonStringToMapStringInterface takes in a json encoded string and transforms
// the data to a map[string]interface{} to comply to the expected type for the
// fusionauth client.
func jsonStringToMapStringInterface(fieldName string, data *schema.ResourceData) (out map[string]interface{}, diags diag.Diagnostics) {
	in := data.Get(fieldName).(string)
	out = map[string]interface{}{}
	if strings.TrimSpace(in) == "" {
		return out, nil
	}

	if err := json.Unmarshal([]byte(in), &out); err != nil {
		fmt.Printf("jsonStringToMapStringInterface: %s", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to transform %s to expected type", fieldName),
			Detail: fmt.Sprintf(
				"Error unmarshalling %s from an expected JSON encoded string to map[string]interface{}.\n"+
					"Please make sure you have wrapped your HCL with jsonencode."+
					"For example, 'jsonencode({ hello = \"world\" })'.\n\n"+
					"error: %s\n",
				fieldName, err,
			),
		})
	}

	return
}

// mapStringInterfaceToJSONString transforms a map[string]interface{} to a JSON
// string.
func mapStringInterfaceToJSONString(fieldName string, in map[string]interface{}) (out string, diags diag.Diagnostics) {
	if len(in) == 0 {
		return "", nil
	}

	outBytes, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to transform %s to expected type", fieldName),
			Detail: fmt.Sprintf(
				"Error marshalling %s from a map[string]interface{} to a JSON string.\n"+
					"error: %s\n",
				fieldName, err,
			),
		})
		return
	}

	return string(outBytes), nil
}
