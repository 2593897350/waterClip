package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"waterclip/api/internal/jobs"
	processorclient "waterclip/api/internal/processor"
)

type processorAPI interface {
	Detect(sourcePath string) (processorclient.DetectResult, error)
	Inpaint(sourcePath, maskPath, mode string) (processorclient.InpaintResult, error)
}

type detectRequest struct {
	SourcePath string `json:"source_path"`
}

type processRequest struct {
	SourcePath string `json:"source_path"`
	MaskPath   string `json:"mask_path"`
	Mode       string `json:"mode"`
}

func NewJobHandler() http.Handler {
	return NewDetectHandler(jobs.NewMemoryStore(), nil)
}

func NewDetectHandler(store *jobs.MemoryStore, processor processorAPI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload detectRequest
		var job jobs.Job
		jobCreated := false
		if file, header, err := r.FormFile("file"); err == nil {
			defer file.Close()
			job = store.Create("detect", "")
			jobCreated = true
			payload.SourcePath = jobs.UploadPath(job.ID, filepath.Base(header.Filename))
			if err := jobs.SaveFile(payload.SourcePath, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			job.SourcePath = payload.SourcePath
			store.Update(job)
		} else {
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		if !jobCreated {
			job = store.Create("detect", payload.SourcePath)
		}
		if processor == nil {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"job_id": job.ID, "status": job.Status})
			return
		}

		result, err := processor.Detect(payload.SourcePath)
		if err != nil {
			job.Status = "failed"
			job.Error = err.Error()
			store.Update(job)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		job.MaskPath = result.MaskPath
		job.Status = "done"
		store.Update(job)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"job_id":    job.ID,
			"status":    job.Status,
			"mask_path": job.MaskPath,
		})
	})
}

func NewProcessHandler(store *jobs.MemoryStore, processor processorAPI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload processRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		job := store.Create("process", payload.SourcePath)
		job.Mode = payload.Mode

		if processor == nil {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"job_id": job.ID, "status": job.Status})
			return
		}

		result, err := processor.Inpaint(payload.SourcePath, payload.MaskPath, payload.Mode)
		if err != nil {
			job.Status = "failed"
			job.Error = err.Error()
			store.Update(job)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		job.MaskPath = payload.MaskPath
		job.ResultPath = result.OutputPath
		job.Status = "done"
		store.Update(job)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"job_id":      job.ID,
			"status":      job.Status,
			"result_path": job.ResultPath,
		})
	})
}

func NewStatusHandler(store *jobs.MemoryStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobID := r.PathValue("jobID")
		job, ok := store.Get(jobID)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"job_id":      job.ID,
			"status":      job.Status,
			"source_path": job.SourcePath,
			"mask_path":   job.MaskPath,
			"result_path": job.ResultPath,
			"mode":        job.Mode,
		})
	})
}
