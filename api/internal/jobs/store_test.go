package jobs

import "testing"

func TestCreateJobDefaultsToQueuedStatus(t *testing.T) {
	store := NewMemoryStore()
	job := store.Create("detect", "source.png")
	if job.Status != "queued" {
		t.Fatalf("expected queued, got %s", job.Status)
	}
}
