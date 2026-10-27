package fusionauth

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	fa "github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// eventTypesNotConfigurable are go-client event types the provider deliberately doesn't expose.
var eventTypesNotConfigurable = map[fa.EventType]string{
	fa.EventType_Test: "only sent by the admin UI's webhook test button (TestEvent)",
}

// sdkEventTypes returns every EventType constant declared by the go-client in use.
func sdkEventTypes(t *testing.T) []fa.EventType {
	t.Helper()

	out, err := exec.Command("go", "list", "-f", "{{.Dir}}", "github.com/FusionAuth/go-client/pkg/fusionauth").Output()
	if err != nil {
		t.Fatalf("locating go-client: %v", err)
	}

	file, err := parser.ParseFile(token.NewFileSet(), filepath.Join(strings.TrimSpace(string(out)), "Domain.go"), nil, 0)
	if err != nil {
		t.Fatalf("parsing go-client Domain.go: %v", err)
	}

	var events []fa.EventType
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			vs := spec.(*ast.ValueSpec)
			if ident, ok := vs.Type.(*ast.Ident); !ok || ident.Name != "EventType" {
				continue
			}
			for _, v := range vs.Values {
				if lit, ok := v.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					s, _ := strconv.Unquote(lit.Value)
					events = append(events, fa.EventType(s))
				}
			}
		}
	}

	if len(events) == 0 {
		t.Fatal("found no EventType constants in go-client Domain.go")
	}

	return events
}

// TestWebhookEventsEnabled_CoversEveryEventType fails when a go-client upgrade brings an event the webhook
// resource doesn't expose yet.
func TestWebhookEventsEnabled_CoversEveryEventType(t *testing.T) {
	eventsSchema := newWebhook().Schema["events_enabled"].Elem.(*schema.Resource).Schema
	allEnabled := map[string]interface{}{}
	for k := range eventsSchema {
		allEnabled[k] = true
	}
	data := schema.TestResourceDataRaw(t, newWebhook().Schema, map[string]interface{}{
		"events_enabled": []interface{}{allEnabled},
	})

	events := buildWebhook(data).EventsEnabled

	for _, e := range sdkEventTypes(t) {
		if _, skip := eventTypesNotConfigurable[e]; skip {
			continue
		}
		if !events[e] {
			t.Errorf("fusionauth_webhook events_enabled has no field for %q", e)
		}
	}
	if len(events) != len(eventsSchema) {
		t.Errorf("events_enabled has %d fields but buildWebhook maps %d events", len(eventsSchema), len(events))
	}
}

// TestTenantEventConfiguration_CoversEveryEventType fails when a go-client upgrade brings an event the tenant
// event_configuration doesn't accept yet.
func TestTenantEventConfiguration_CoversEveryEventType(t *testing.T) {
	validate := newTenant().Schema["event_configuration"].Elem.(*schema.Resource).Schema["event"].ValidateFunc

	for _, e := range sdkEventTypes(t) {
		if _, skip := eventTypesNotConfigurable[e]; skip {
			continue
		}
		if _, errs := validate(string(e), "event"); len(errs) > 0 {
			t.Errorf("fusionauth_tenant event_configuration rejects %q", e)
		}
	}
}

func TestBuildEventsEnabled_OmitsDisabledNewerEvents(t *testing.T) {
	data := schema.TestResourceDataRaw(t, newWebhook().Schema, map[string]interface{}{
		"events_enabled": []interface{}{
			map[string]interface{}{
				"user_create":   true,
				"entity_create": false,
			},
		},
	})

	events := buildWebhook(data).EventsEnabled

	for _, e := range webhookEventsOmittedWhenDisabled() {
		if _, ok := events[e]; ok {
			t.Errorf("events contains %q, want it omitted so older servers accept the request", e)
		}
	}
	if got := events[fa.EventType_UserCreate]; !got {
		t.Errorf("user.create = %v, want true", got)
	}
	if _, ok := events[fa.EventType_UserDelete]; !ok {
		t.Errorf("events is missing %q, want disabled events older than 1.68.0 still sent", fa.EventType_UserDelete)
	}
}

func TestBuildEventsEnabled_SendsEnabledNewerEvents(t *testing.T) {
	allEnabled := map[string]interface{}{}
	for k := range newWebhook().Schema["events_enabled"].Elem.(*schema.Resource).Schema {
		allEnabled[k] = true
	}
	data := schema.TestResourceDataRaw(t, newWebhook().Schema, map[string]interface{}{
		"events_enabled": []interface{}{allEnabled},
	})

	events := buildWebhook(data).EventsEnabled

	for _, e := range webhookEventsOmittedWhenDisabled() {
		if got := events[e]; !got {
			t.Errorf("%q = %v, want true", e, got)
		}
	}
}

// TestAccFusionauthWebhook_basic runs on every server version, so CI's older-version runs fail when a newer
// event is sent to a server that doesn't know it.
func TestAccFusionauthWebhook_basic(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_webhook.test_%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookResourceConfig(resourceName, "user_create = true\n    user_delete = false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthWebhookEvents(tfResourcePath, map[fa.EventType]bool{
						fa.EventType_UserCreate: true,
						fa.EventType_UserDelete: false,
					}),
					resource.TestCheckResourceAttr(tfResourcePath, "events_enabled.0.user_create", "true"),
					resource.TestCheckResourceAttr(tfResourcePath, "events_enabled.0.user_delete", "false"),
				),
			},
			{
				Config: testAccWebhookResourceConfig(resourceName, "user_create = false\n    user_delete = true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthWebhookEvents(tfResourcePath, map[fa.EventType]bool{
						fa.EventType_UserCreate: false,
						fa.EventType_UserDelete: true,
					}),
					resource.TestCheckResourceAttr(tfResourcePath, "events_enabled.0.user_create", "false"),
					resource.TestCheckResourceAttr(tfResourcePath, "events_enabled.0.user_delete", "true"),
				),
			},
		},
	})
}

// TestAccFusionauthWebhook_newerEvents covers the 1.68.0 two-factor and 1.69.0 entity events; skips on older
// servers.
func TestAccFusionauthWebhook_newerEvents(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_webhook.test_%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfFusionAuthBelow(t, "1.69.0") },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookResourceConfig(resourceName, "entity_create = true\n    entity_update_complete = true\n    user_identity_update = true\n    user_two_factor_challenge = true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthWebhookEvents(tfResourcePath, map[fa.EventType]bool{
						fa.EventType_EntityCreate:           true,
						fa.EventType_EntityUpdateComplete:   true,
						fa.EventType_UserIdentityUpdate:     true,
						fa.EventType_UserTwoFactorChallenge: true,
						fa.EventType_EntityDelete:           false,
					}),
					resource.TestCheckResourceAttr(tfResourcePath, "events_enabled.0.entity_create", "true"),
					resource.TestCheckResourceAttr(tfResourcePath, "events_enabled.0.entity_delete", "false"),
					resource.TestCheckResourceAttr(tfResourcePath, "events_enabled.0.user_identity_update", "true"),
				),
			},
			{
				// Disabling an omitted event must turn it off on the server, not leave the previous value.
				Config: testAccWebhookResourceConfig(resourceName, "entity_create = false\n    entity_delete = true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFusionauthWebhookEvents(tfResourcePath, map[fa.EventType]bool{
						fa.EventType_EntityCreate:           false,
						fa.EventType_EntityDelete:           true,
						fa.EventType_EntityUpdateComplete:   false,
						fa.EventType_UserIdentityUpdate:     false,
						fa.EventType_UserTwoFactorChallenge: false,
					}),
					resource.TestCheckResourceAttr(tfResourcePath, "events_enabled.0.entity_create", "false"),
					resource.TestCheckResourceAttr(tfResourcePath, "events_enabled.0.entity_delete", "true"),
				),
			},
		},
	})
}

func testAccWebhookResourceConfig(resourceName, events string) string {
	return fmt.Sprintf(`
resource "fusionauth_webhook" "test_%[1]s" {
  connect_timeout = 1000
  read_timeout    = 2000
  url             = "http://localhost:12345/test-acc-webhook/%[1]s"
  events_enabled {
    %[2]s
  }
}
`, resourceName, events)
}

// testAccCheckFusionauthWebhookEvents checks the events FusionAuth stored; a missing key counts as disabled.
func testAccCheckFusionauthWebhookEvents(resourceName string, want map[fa.EventType]bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		faClient := fusionauthClient()
		res, err := faClient.RetrieveWebhook(rs.Primary.ID)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get resource: %#+v", res)
		}

		for e, enabled := range want {
			if got := res.Webhook.EventsEnabled[e]; got != enabled {
				return fmt.Errorf("event %q = %t, want %t", e, got, enabled)
			}
		}

		return nil
	}
}

func testAccCheckFusionauthWebhookDestroy(s *terraform.State) error {
	faClient := fusionauthClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_webhook" {
			continue
		}

		// Ensure we retry for eventual consistency in HA setups.
		err := retry.RetryContext(context.Background(), retryTimeout, func() *retry.RetryError {
			res, err := faClient.RetrieveWebhook(rs.Primary.ID)
			if err != nil {
				return retry.NonRetryableError(err)
			}

			if res.StatusCode == http.StatusNotFound {
				// resource destroyed!
				return nil
			}

			return retry.RetryableError(fmt.Errorf("fusionauth resource still exists: %s", rs.Primary.ID))
		})

		if err != nil {
			// We failed destroying the resource...
			return err
		}
	}

	return nil
}
