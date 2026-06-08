package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	Health().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestCreateDetectJobReturnsAccepted(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/jobs/detect", strings.NewReader(`{"source_path":"var/uploads/a.ppm"}`))
	rec := httptest.NewRecorder()
	NewJobHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rec.Code)
	}
}
