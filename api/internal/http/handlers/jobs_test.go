package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"waterclip/api/internal/jobs"
	"waterclip/api/internal/processor"
)

type stubProcessor struct{}

func (stubProcessor) Detect(sourcePath string) (processor.DetectResult, error) {
	return processor.DetectResult{MaskPath: sourcePath + ".mask.pgm"}, nil
}

func (stubProcessor) Inpaint(sourcePath, maskPath, mode string) (processor.InpaintResult, error) {
	return processor.InpaintResult{OutputPath: sourcePath + "." + mode + ".ppm"}, nil
}

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
	NewDetectHandler(jobs.NewMemoryStore(), stubProcessor{}).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["status"] != "done" {
		t.Fatalf("expected done status, got %q", body["status"])
	}
	if body["mask_path"] == "" {
		t.Fatalf("expected mask_path in response")
	}
}

func TestCreateProcessJobReturnsResultPath(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/jobs/process",
		strings.NewReader(`{"source_path":"var/uploads/a.ppm","mask_path":"var/masks/a.pgm","mode":"fast"}`),
	)
	rec := httptest.NewRecorder()
	NewProcessHandler(jobs.NewMemoryStore(), stubProcessor{}).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["result_path"] == "" {
		t.Fatalf("expected result_path in response")
	}
}

func TestStatusHandlerReturnsStoredJob(t *testing.T) {
	store := jobs.NewMemoryStore()
	job := store.Create("detect", "var/uploads/a.ppm")
	job.Status = "done"
	job.MaskPath = "var/masks/a.pgm"
	store.Update(job)

	req := httptest.NewRequest(http.MethodGet, "/api/jobs/"+job.ID, nil)
	req.SetPathValue("jobID", job.ID)
	rec := httptest.NewRecorder()

	NewStatusHandler(store).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["job_id"] != job.ID {
		t.Fatalf("expected job_id %q, got %q", job.ID, body["job_id"])
	}
	if body["status"] != "done" {
		t.Fatalf("expected done status, got %q", body["status"])
	}
}

func TestCreateDetectJobFromMultipartUpload(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	fileWriter, err := writer.CreateFormFile("file", "sample.ppm")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	if _, err := fileWriter.Write([]byte("P3\n1 1\n255\n255 255 255\n")); err != nil {
		t.Fatalf("failed to write file body: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/jobs/detect", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	NewDetectHandler(jobs.NewMemoryStore(), stubProcessor{}).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response["job_id"] == "" {
		t.Fatalf("expected job_id in response")
	}
	if response["mask_path"] == "" {
		t.Fatalf("expected mask_path in response")
	}
}
