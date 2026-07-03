# relay

A rate-limited job queue with a concurrent worker pool, HTTP delivery, retries, and SQLite persistence — built as a hands-on deep dive into Go's concurrency primitives (goroutines, channels, `sync`, `context`) rather than a toy exercise.

Modeled loosely after webhook delivery systems like Svix: accept a job over HTTP, attempt delivery to a target URL, retry on failure with backoff, and persist anything that ultimately fails for good.

## Why this exists

This project was built to move past reading Go's spec and idioms and into actually feeling the failure modes — data races, resource leaks, goroutine leaks, graceful shutdown — that only show up when you write concurrent code under real (if simulated) load. Every core concurrency claim in this codebase was verified by hand before being trusted: worker pool concurrency was proven with timestamped delayed jobs, and the need for the `Stats` mutex was proven by first reproducing an actual data race with `go run -race`, then fixing it.

## Features

- **HTTP ingestion** — `POST /jobs` accepts a job (target URL + payload) and enqueues it
- **Concurrent worker pool** — a fixed number of goroutines process jobs off a shared buffered channel
- **Real HTTP delivery** — each job is delivered via an actual HTTP POST to its target URL, with explicit success/failure handling
- **Thread-safe stats tracking** — success/failure counts are tracked safely across concurrent workers via a mutex-guarded `Stats` type
- **Status endpoint** — `GET /status` reports current queue depth and delivery counts
- *(in progress)* Retry with exponential backoff
- *(planned)* Dead-letter persistence via SQLite
- *(planned)* Graceful shutdown on SIGTERM/SIGINT

## Architecture

```
relay/
├── main.go              # wires everything together, HTTP handlers, worker startup
├── job/
│   └── job.go             # Job struct and status constants
├── queue/
│   └── queue.go             # Queue type wrapping a buffered channel
├── stats/
│   └── stats.go               # Thread-safe success/failure counters
├── worker/                     # (planned) worker pool logic, extracted from main.go
├── delivery/                    # (planned) HTTP delivery + retry/backoff logic
└── store/                        # (planned) SQLite dead-letter persistence
```

## Running locally

```bash
go run main.go
```

Server listens on `:8080`.

**Submit a job:**
```bash
curl -X POST localhost:8080/jobs \
  -d '{"target_url":"https://example.com/webhook","payload":"hello world"}'
```

**Check status:**
```bash
curl localhost:8080/status
```

## Design notes

- **Worker pool** — N goroutines all read from the same buffered channel (`chan job.Job`); Go's channel semantics handle fan-out automatically, no manual load balancing needed.
- **Stats safety** — the `Stats` struct holds its own `sync.Mutex` internally, guarding `success`/`failure` counters. This was deliberately verified against a real data race (reproduced in isolation, then fixed) before being trusted in the main codebase — see commit history.
- **Dependency injection over globals** — shared state (`*queue.Queue`, `*stats.Stats`) is passed explicitly into `worker`/`deliver` rather than stored as package-level globals, to keep the code testable and its dependencies visible from function signatures alone.

## Status

Actively under development. Core ingestion, delivery, and concurrent worker pool are working and verified. Retry/backoff, persistence, and graceful shutdown are next.

## Author

Built by [Frank](https://github.com/frankuccino) (Franyx Studios) as part of a deliberate, hands-on push into backend and distributed systems engineering with Go.
