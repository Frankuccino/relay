package main

import (
	"bytes"
	"encoding/json"
	// "time"

	"log"
	"net/http"

	"github.com/Frankuccino/relay/job"
	"github.com/Frankuccino/relay/queue"
	"github.com/Frankuccino/relay/stats"
	"github.com/google/uuid"
)

func main() {
	q := queue.New(100)
	st := stats.New()

	const numWorkers = 5

	for range numWorkers {
		go worker(q, st)
	}

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		success, failure := st.Snapshot()
		json.NewEncoder(w).Encode(map[string]int{
			"success":     success,
			"failure":     failure,
			"queue_depth": q.Len(),
		})
	})

	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			TargetURL string `json:"target_url"`
			Payload   string `json:"payload"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		j := job.New(generateID(), req.TargetURL, []byte(req.Payload))
		q.Enqueue(j)

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"job_id": j.ID})
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Client!"))
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func worker(q *queue.Queue, st *stats.Stats) {
	for j := range q.Jobs() {
		deliver(j, st)
	}
}

func deliver(j job.Job, st *stats.Stats) {
	// time.Sleep(2 * time.Second)
	resp, err := http.Post(j.TargetURL, "application/json", bytes.NewReader(j.Payload))
	if err != nil {
		log.Printf("job %s failed: %v", j.ID, err)
		st.RecordFailure()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("job %s delivered successfully (status %d)", j.ID, resp.StatusCode)
		st.RecordSuccess()
	} else {
		log.Printf("job %s failed with status %d", j.ID, resp.StatusCode)
		st.RecordFailure()
	}
}

func generateID() string {
	return uuid.New().String()
}
