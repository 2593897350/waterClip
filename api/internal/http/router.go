package httpx

import (
	"net/http"

	"waterclip/api/internal/jobs"
	"waterclip/api/internal/http/handlers"
	processorclient "waterclip/api/internal/processor"
)

func NewRouter(processorBaseURL string) http.Handler {
	mux := http.NewServeMux()
	store := jobs.NewMemoryStore()
	processor := processorclient.New(processorBaseURL)
	mux.Handle("/health", handlers.Health())
	mux.Handle("/api/jobs/detect", handlers.NewDetectHandler(store, processor))
	mux.Handle("/api/jobs/process", handlers.NewProcessHandler(store, processor))
	mux.Handle("/api/jobs/{jobID}", handlers.NewStatusHandler(store))
	return mux
}
