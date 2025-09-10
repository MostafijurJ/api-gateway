# API Gateway (Go)

Production-ready API Gateway in Go featuring routing, reverse proxying, load balancing, middleware, and observability.

## Features
- Reverse proxy with configurable routes (path/method/headers)
- Upstream pools with round-robin and least-connections
- CORS, API key, optional JWT checks, request ID
- In-memory IP rate limiting
- Prometheus metrics at `/metrics`

## Quick Start
1) Start example backends (in separate terminals):
```bash
go run ./cmd/usersvc
go run ./cmd/inventorysvc
```

2) Run the gateway with the provided config:
```bash
GATEWAY_CONFIG=$(pwd)/config.yml go run .
```

3) Test through the gateway:
```bash
curl -H "X-API-Key: DEV_KEY_12345" http://localhost:8081/users
curl -H "X-API-Key: DEV_KEY_12345" http://localhost:8081/inventory
curl http://localhost:8081/metrics
```

## Configuration
Edit `config.yml` to change:
- `http`: address/timeouts
- `routes`: matching rules and upstream/pool
- `pools`: backend lists and strategy
- `cors`, `auth` (API keys/JWT), `rate`

## Project Structure
- `config/` config loader
- `internal/` router, proxy, middleware, observability, server
- `cmd/usersvc`, `cmd/inventorysvc` example services


