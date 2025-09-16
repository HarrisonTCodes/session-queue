# Session Queue

[![Go Report Card](https://goreportcard.com/badge/github.com/HarrisonTCodes/session-queue)](https://goreportcard.com/report/github.com/HarrisonTCodes/session-queue)
[![Go Version](https://img.shields.io/badge/go-1.24-blue?logo=go)](https://golang.org)
![CI](https://github.com/HarrisonTCodes/session-queue/actions/workflows/ci.yaml/badge.svg)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A solution to the system design problem of a fair, scalable queue and waiting list service, similar those seen with popular events retailers such as Ticketmaster. Written in Go.

## Design
The design of the system is documented in the [system doc](./docs/system.md).
The server(s) issue JWTs to users designating their queue position, tracking the number of issued tokens in Redis. A sliding window algorithm is used to determine, in batches, which tickets are active, and which are still waiting (or expired), sliding at a regular interval.

## Run Locally
The system requires a Redis instance to be accessible. You can spin one up with the `docker-compose` config:
```bash
docker compose up
```

You can then run the Go app on bare metal with:
```bash
go run ./cmd/server
```

You will need various environment variables set, listed below:
```bash
REDIS_ADDR          # Full address of Redis instance (default: localhost:6379)
PORT                # Port for the HTTP server to run on
JWT_SECRET          # Secret used to sign JWTs
WINDOW_SIZE         # Size of active ticket window(s)
WINDOW_INTERVAL     # How often to slide window to next batch (in seconds)
ACTIVE_WINDOW_COUNT # How many windows (counted backwards from current) considered active
```

You should then be able to hit the server on the port specified with a HTTP client. For example, with HTTPie:
```bash
http POST localhost:{port}/join
http GET localhost:{port}/status "Authorization: Bearer {token}"
```

The repo also has some default Kubernetes configuration, that can be used for running the system, by default with 3 server replicas.
You can find this config in the [k8s folder](./k8s/), and documentation on it and how to run locally (with a minikube cluster) in the [k8s docs](./docs/k8s.md).
You can apply the [k8s folder](./k8s/) with the following:
```bash
kubectl apply -f k8s
```
