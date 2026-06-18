package fusionauth

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

func fusionauthServerVersion(t *testing.T) string {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, os.Getenv("FA_DOMAIN")+"/api/system/version", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", os.Getenv("FA_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("fetching FusionAuth version: %s", err)
	}
	defer resp.Body.Close()

	var body struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("reading FusionAuth version: %s", err)
	}
	return body.Version
}

func skipIfFusionAuthBelow(t *testing.T, minVersion string) {
	t.Helper()
	if v := fusionauthServerVersion(t); versionOlderThan(v, minVersion) {
		t.Skipf("requires FusionAuth >= %s, server is %s", minVersion, v)
	}
}

func versionOlderThan(a, b string) bool {
	x, y := parseVersion(a), parseVersion(b)
	for i := 0; i < 3; i++ {
		if x[i] != y[i] {
			return x[i] < y[i]
		}
	}
	return false
}

func parseVersion(v string) [3]int {
	v = strings.SplitN(v, "-", 2)[0]
	v = strings.SplitN(v, "+", 2)[0]
	var out [3]int
	for i, part := range strings.SplitN(v, ".", 3) {
		out[i], _ = strconv.Atoi(part)
	}
	return out
}

func Test_versionOlderThan(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"1.67.0", "1.68.0", true},
		{"1.70.0", "1.68.0", false},
		{"1.68.0", "1.68.0", false},
		{"1.7.0", "1.68.0", true},
		{"1.68.0-beta.284", "1.68.0", false},
	}
	for _, c := range cases {
		if got := versionOlderThan(c.a, c.b); got != c.want {
			t.Errorf("versionOlderThan(%q, %q) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}
