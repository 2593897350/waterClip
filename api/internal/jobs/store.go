package jobs

import (
	"crypto/rand"
	"encoding/hex"
)

type MemoryStore struct {
	items map[string]Job
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{items: map[string]Job{}}
}

func (s *MemoryStore) Create(kind, sourcePath string) Job {
	job := Job{
		ID:         newJobID(),
		Kind:       kind,
		SourcePath: sourcePath,
		Status:     "queued",
	}
	s.items[job.ID] = job
	return job
}

func (s *MemoryStore) Update(job Job) {
	s.items[job.ID] = job
}

func (s *MemoryStore) Get(id string) (Job, bool) {
	job, ok := s.items[id]
	return job, ok
}

func newJobID() string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
