package fusionauth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const retryableConflictBody = `{"fieldErrors":{},"generalErrors":[{"code":"[retryableConflict]","message":"The request could not be completed because of a conflict. Please use an exponential backoff and try again."}]}`

func TestPatchIdentityProviderRetriesRetryableConflict(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		if calls == 1 {
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(retryableConflictBody))
			return
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	err := patchIdentityProvider([]byte(`{}`), "id", Client{Host: srv.URL, APIKey: "key"})
	if err != nil {
		t.Fatalf("expected retried patch to succeed, got: %s", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls (conflict then success), got %d", calls)
	}
}

func TestPatchIdentityProviderDoesNotRetryOtherConflicts(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"generalErrors":[{"code":"[somethingElse]"}]}`))
	}))
	defer srv.Close()

	err := patchIdentityProvider([]byte(`{}`), "id", Client{Host: srv.URL, APIKey: "key"})
	if err == nil {
		t.Fatal("expected non-retryable conflict to fail")
	}
	if calls != 1 {
		t.Fatalf("expected 1 call for a non-retryable conflict, got %d", calls)
	}
}

func TestPatchIdentityProviderRetryExhaustion(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(retryableConflictBody))
	}))
	defer srv.Close()

	err := patchIdentityProvider([]byte(`{}`), "id", Client{Host: srv.URL, APIKey: "key"})
	if err == nil {
		t.Fatal("expected exhausted retries to fail")
	}
	if !strings.Contains(err.Error(), "retryableConflict") {
		t.Fatalf("expected the last conflict error, got: %s", err)
	}
	if calls != 6 {
		t.Fatalf("expected 6 attempts before giving up, got %d", calls)
	}
}
