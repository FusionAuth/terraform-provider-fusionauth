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
	return handleStringSliceFromSet(data.Get(key).(*schema.Set))
}

func handleStringSliceFromList(list []interface{}) []string {
	s := make([]string, 0, len(list))

	for _, x := range list {
		s = append(s, x.(string))
	}

	return s
}

func handleStringSliceFromSet(set *schema.Set) []string {
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

// jsonStringToMapStringInterface reads data for a "data" key, which it expects
// to be a json encoded string and transforms the json data to a map[string]interface{}
// to comply to the expected type for the fusionauth client.
func jsonStringToMapStringInterface(in string) (out map[string]interface{}, diags diag.Diagnostics) {
	out = map[string]interface{}{}
	if strings.TrimSpace(in) == "" {
		return out, nil
	}

	if err := json.Unmarshal([]byte(in), &out); err != nil {
		fmt.Printf("jsonStringToMapStringInterface: %s", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to transform data to expected type",
			Detail: fmt.Sprintf(
				"Error unmarshalling data from an expected JSON encoded string to map[string]interface{}.\n"+
					"Please make sure you have wrapped your HCL with jsonencode."+
					"For example, 'jsonencode({ hello = \"world\" })'.\n\n"+
					"error: %s\n",
				err,
			),
		})
	}

	return
}

// mapStringInterfaceToJSONString transforms a map[string]interface{} to a JSON
// string.
func mapStringInterfaceToJSONString(in map[string]interface{}) (out string, diags diag.Diagnostics) {
	if len(in) == 0 {
		return "", nil
	}

	outBytes, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to transform data to expected type",
			Detail: fmt.Sprintf(
				"Error marshalling data from a map[string]interface{} to a JSON string.\n"+
					"error: %s\n",
				err,
			),
		})
		return
	}

	return string(outBytes), nil
}

// setResourceData given a data mapping, will load the provided data into the
// terraform state.
//
// Note:
//   - This performs simple top level loading and returning a build up of errors.
//   - Any sub-objects/maps/lists requiring specific ordering will need to be
//     handled manually.
func setResourceData(resource string, data *schema.ResourceData, dataMapping map[string]interface{}) (diags diag.Diagnostics) {
	for k, v := range dataMapping {
		switch k {
		case "data":
			if resourceData, dataDiags := mapStringInterfaceToJSONString(v.(map[string]interface{})); dataDiags != nil {
				diags = append(diags, dataDiags...)
			} else if err := data.Set("data", resourceData); err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("entity.%s: %s", k, err),
				})
			}

		default:
			if err := data.Set(k, v); err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("%s.%s: %s", resource, k, err),
				})
			}
		}
	}

	return diags
}

func intMapToStringMap(intMap map[string]interface{}) map[string]string {
	m := map[string]string{}
	for k, v := range intMap {
		if s, ok := v.(string); ok {
			m[k] = s
		}
	}

	return m
}

// clientTenantIDOverride takes in the client and the data. As long as the
// resource has a tenant_id attribute set, it will override the client's tenant
// id until the revert function is called.
func clientTenantIDOverride(client *Client, data *schema.ResourceData) (revert func()) {
	if tid := data.Get("tenant_id").(string); tid != "" {
		oldTenantID := client.FAClient.TenantId

		client.FAClient.TenantId = tid
		return func() {
			client.FAClient.TenantId = oldTenantID
		}
	}

	return func() {
		// nothing to revert here...
	}
}
