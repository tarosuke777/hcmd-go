package hvapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchMaxTimestampWithUntrustedCertificate(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"max_created_at":"2026-10-11 03:00:00"}`)
	}))
	defer server.Close()

	got, err := FetchMaxTimestamp(server.URL)
	if err != nil {
		t.Fatalf("FetchMaxTimestamp() error = %v", err)
	}
	if formatted := got.Format(timeLayout); formatted != "2026-10-11 03:00:00" {
		t.Errorf("FetchMaxTimestamp() = %q, want %q", formatted, "2026-10-11 03:00:00")
	}
}
