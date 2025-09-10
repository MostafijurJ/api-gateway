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

## Request Flow Architecture

The API Gateway processes requests through a middleware chain and routes them to appropriate backend services.

### Flow Diagram

```
┌─────────────┐    ┌─────────────────────────────────────────────────────────┐
│   Client    │────▶│                API Gateway (:8081)                      │
│ Application │    │                                                         │
└─────────────┘    │  ┌─────────────────────────────────────────────────────┐ │
                   │  │              Middleware Chain                       │ │
                   │  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────────┐ │ │
                   │  │  │  CORS   │▶│ API Key │▶│  JWT    │▶│Rate Limiting│ │ │
                   │  │  │         │ │  Auth   │ │  Auth   │ │             │ │ │
                   │  │  └─────────┘ └─────────┘ └─────────┘ └─────────────┘ │ │
                   │  └─────────────────────────────────────────────────────┘ │
                   │                            │                             │
                   │  ┌─────────────────────────▼─────────────────────────────┐ │
                   │  │                    Router                             │ │
                   │  │         (Path/Method/Header Matching)                 │ │
                   │  └─────────────────┬─────────────────┬─────────────────────┘ │
                   └────────────────────┼─────────────────┼─────────────────────────┘
                                       │                 │
                    ┌──────────────────▼─────────────────▼──────────────────┐
                    │                Load Balancer                         │
                    │         (Round Robin / Least Connections)            │
                    └┬─────────────────────────────────────────────────────┬┘
                     │                                                     │
        ┌────────────▼────────────┐                         ┌─────────────▼────────────┐
        │    Users Pool           │                         │   Inventory Service      │
        │                        │                         │     (:9002)              │
        │ ┌─────────┐ ┌─────────┐ │                         │                          │
        │ │User Svc │ │User Svc │ │                         │ ┌─────────────────────┐  │
        │ │ (:9001) │ │ (:9002) │ │                         │ │  GET /inventory     │  │
        │ └─────────┘ └─────────┘ │                         │ │  Returns: Items     │  │
        │                        │                         │ └─────────────────────┘  │
        │ GET/POST /users         │                         └──────────────────────────┘
        │ Returns: User data      │
        └─────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────┐
│                              Observability                                     │
│                                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐ │
│  │ Structured  │  │ Prometheus  │  │ Request ID  │  │     Health Checks       │ │
│  │   Logging   │  │  Metrics    │  │  Tracking   │  │  /healthz  /readyz      │ │
│  │             │  │ (/metrics)  │  │             │  │                         │ │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### Request Processing Steps

1. **Client Request**: HTTP request arrives at the API Gateway (port 8081)

2. **Middleware Chain Processing**:
   - **CORS**: Handles cross-origin resource sharing policies
   - **API Key Authentication**: Validates `X-API-Key` header against configured keys
   - **JWT Authentication**: Optional JWT token validation
   - **Rate Limiting**: In-memory IP-based request throttling
   - **Request ID**: Assigns unique identifier for request tracing

3. **Routing**: The router matches requests based on:
   - Path prefix (e.g., `/users`, `/inventory`)
   - HTTP methods (GET, POST, PUT, DELETE)
   - Header values (optional)

4. **Load Balancing & Proxying**:
   - **Pool-based**: Routes to backend pools with load balancing strategies
     - Round-robin distribution
     - Least-connections distribution
   - **Direct upstream**: Routes directly to a specific backend service

5. **Backend Services**:
   - **User Service** (port 9001/9002): Handles user-related operations
   - **Inventory Service** (port 9002): Manages inventory data

6. **Observability**:
   - Structured logging with request IDs
   - Prometheus metrics collection
   - Health monitoring endpoints

### Example Request Flows

**User Request with Load Balancing**:
```
GET /users → Middleware → Router → Users Pool → User Service (9001 or 9002)
```

**Inventory Request to Direct Upstream**:
```
GET /inventory → Middleware → Router → Inventory Service (9002)
```

**Metrics Collection**:
```
GET /metrics → Skip Auth Middleware → Prometheus Handler
```

### Configuration Notes

The example `config.yml` demonstrates both load balancing patterns:

- **Users Pool**: Configured with multiple backends for load balancing (though in the demo, one backend runs inventory service for simplicity)
- **Inventory Direct**: Routes directly to a specific upstream service

For production use, ensure all backends in a pool serve the same service endpoints.


