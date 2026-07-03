package stats

import "sync"

type Stats struct {
	mu      sync.Mutex
	success int
	failure int
}

func New() *Stats {
	return &Stats{}
}

func (s *Stats) RecordSuccess() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.success++
}

func (s *Stats) RecordFailure() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failure++
}

func (s *Stats) Snapshot() (success, failure int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.success, s.failure
}
