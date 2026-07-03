package job

import "time"

type Status string

const (
	StatusPending Status = "pending"
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
)

type Job struct {
	ID        string
	TargetURL string
	Payload   []byte
	Status    Status
	CreatedAt time.Time
}

func New(id, targetURL string, payload []byte) Job {
	return Job{
		ID:        id,
		TargetURL: targetURL,
		Payload:   payload,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
}
