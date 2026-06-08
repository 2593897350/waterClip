package httpx

import (
	"net/http"

	"waterclip/api/internal/http/handlers"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/health", handlers.Health())
	mux.Handle("/api/jobs/detect", handlers.NewJobHandler())
	return mux
}
