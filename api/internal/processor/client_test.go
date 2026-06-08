package processor

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthReturnsOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	client := New(server.URL)
	if err := client.Health(); err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
}
