package handlers

import "net/http"

func NewJobHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"job_id":"demo-job","status":"queued"}`))
	})
}
