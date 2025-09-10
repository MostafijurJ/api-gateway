# High-Performance API Gateway (Go)

This document outlines the **requirements** and **step-by-step implementation plan** for building a production-grade API Gateway using Go.

---

## 📌 Requirements for an API Gateway

### 1. Core Routing & Proxy
- Reverse proxy: route client requests to backend services
- URL/path-based routing (e.g., `/users` → User Service)
- Header/query parameter-based routing
- Protocol support: HTTP/HTTPS (optionally gRPC, WebSockets)

### 2. Load Balancing
- Round-robin, random, least-connections strategies
- Health checks (remove unhealthy instances)
- Service discovery (Consul, Etcd, or Kubernetes)

### 3. Authentication & Authorization
- Support for JWT and OAuth2
- API key validation
- Role-based access control (RBAC)
- Integration with external Identity Providers (Keycloak, Auth0)

### 4. Security
- TLS termination (HTTPS)
- IP whitelisting/blacklisting
- CORS handling
- Protection against common attacks:
    - Rate limiting (DDoS, brute-force prevention)
    - Request validation (size limits, schema validation)
    - CSRF/XSS protection (headers, sanitization)

### 5. Traffic Management
- Rate limiting & request throttling
- Quotas per API key/user/service
- Request retries & timeouts
- Circuit breaker pattern (fail fast if a service is down)

### 6. Observability
- Structured logging (per request/response)
- Metrics (latency, throughput, error rates)
- Distributed tracing (OpenTelemetry, Jaeger)
- Monitoring integrations (Prometheus, Grafana)

### 7. Transformation & Mediation
- Request/response transformation (headers, payload)
- Protocol translation (REST ↔ gRPC, XML ↔ JSON)
- Versioning support (e.g., `/v1/`, `/v2/`)

### 8. Developer Experience
- API documentation auto-generation (OpenAPI/Swagger)
- Sandbox or mock responses
- Developer portal (optional, for external API consumers)

### 9. Extensibility
- Plugin/middleware system (e.g., logging, custom auth)
- Configurable routing without restarting (hot reload)
- Support for custom policies (Lua, Go plugins, etc.)

### 10. Scalability & Deployment
- Horizontal scaling across multiple nodes
- Stateless design (store session/state in Redis or DB)
- Works with Kubernetes ingress or standalone
- Config management (YAML/JSON, or via service registry)

---

## 🚀 Step-by-Step Implementation Plan

### Phase 1: Preparation & Setup
- Initialize repo and Go module
- Define project structure (config, router, proxy, middleware, registry, observability)
- Setup CI/CD (linting, testing, build pipeline)
- Define configuration format (YAML/JSON)

**Deliverable:** Project skeleton with CI pipeline

---

### Phase 2: Basic HTTP Server
- Setup HTTP server with health endpoints (`/healthz`, `/readyz`)
- Implement config loader
- Add structured logging
- Graceful shutdown (handle SIGTERM)

**Deliverable:** Running server with config support

---

### Phase 3: Minimal Reverse Proxy
- Implement reverse proxy forwarding requests to a single upstream
- Pass headers/body correctly
- Handle upstream timeouts and errors
- Add request/response logging

**Deliverable:** Requests forwarded to upstream with logs

---

### Phase 4: Routing Engine
- Implement route matching:
    - Path prefix, exact path
    - HTTP method
    - Header/query param rules
- Support per-route middleware config
- Add route priority and default fallback

**Deliverable:** Configurable routing with multiple upstreams

---

### Phase 5: Load Balancing & Health Checks
- Define backend pools per upstream
- Implement strategies: round-robin, least-connections
- Active & passive health checks
- Auto-remove unhealthy backends

**Deliverable:** Load-balanced requests with failover

---

### Phase 6: Middleware Framework
- Create middleware chain system
- Add:
    - Authentication (JWT, API Key, OAuth2)
    - Logging middleware
    - Tracing middleware (OpenTelemetry)

**Deliverable:** Middleware applied per-route

---

### Phase 7: Rate Limiting & Quotas
- Implement token bucket/leaky bucket algorithm
- Redis-based distributed counters
- Configurable per-IP, per-user, or per-route limits
- Return `429` with `X-RateLimit-*` headers

**Deliverable:** Rate-limited API traffic

---

### Phase 8: Service Discovery & Dynamic Config
- Integrate with Consul/Etcd/Kubernetes for backend discovery
- Implement admin API to:
    - Add/remove routes
    - Reload config
    - Query current config/state

**Deliverable:** Dynamic configuration & discovery support

---

### Phase 9: Observability
- Structured logs with request IDs
- Prometheus metrics:
    - Requests count
    - Latency histograms
    - Backend health states
- Distributed tracing with Jaeger/OpenTelemetry

**Deliverable:** `/metrics` endpoint + tracing support

---

### Phase 10: Security Hardening
- TLS termination & cert rotation
- CORS policies
- IP allow/deny rules
- Request validation & body size limits
- mTLS between gateway & upstreams (optional)

**Deliverable:** Hardened gateway with TLS & secure headers

---

### Phase 11: Reliability Patterns
- Request retries with backoff
- Per-route timeouts
- Circuit breaker per upstream
- Bulkheading (limit connections per backend)

**Deliverable:** Resilient gateway against upstream failures

---

### Phase 12: Caching & Transformations (Optional)
- Response caching with Redis
- Cache invalidation API
- Request/response transformations
- Protocol translation (REST ↔ gRPC)

**Deliverable:** API mediation features

---

### Phase 13: Testing & Validation
- Unit tests for router, proxy, middleware
- Integration tests with mock upstreams
- Load testing (k6, Vegeta)
- Fault injection (simulate failures)

**Deliverable:** Automated test suite with performance reports

---

### Phase 14: Deployment & CI/CD
- Containerize with Docker
- Provide Helm chart / K8s manifests
- Add health/readiness probes
- Blue/Green or rolling deployment strategy
- Security scanning in pipeline

**Deliverable:** Production-ready deployment

---

### Phase 15: Monitoring & Runbooks
- Create Grafana dashboards
- Setup alerting (Prometheus Alertmanager)
- Write SRE runbooks for:
    - Backend failures
    - High latency
    - Certificate renewal

**Deliverable:** Operational playbooks & dashboards

---

## ✅ MVP Checklist
Before production, ensure you have:
- Reverse proxy + routing
- Load balancing + health checks
- JWT + API key auth
- Rate limiting
- Logging + metrics + tracing
- Admin API for config reload
- Docker/K8s deployment manifests

---

## 🎯 Key Learning Outcomes
- Mastery of Go networking (`net/http`, `httputil`)
- Middleware & observability design
- Resilient distributed system patterns
- Secure, scalable API management
- Production-ready deployment with CI/CD


## DB and other Connections Configurations
REDIS_URL=redis://localhost:6379
ACTIVE_PROFILES=dev


## 📌 Project structure?

- **`cmd/gateway`** → keeps the entry point clean, only bootstraps the app.
- **`config`** → central place for configuration (env, YAML, etc.).
- **`internal/router`** → defines routes and middleware setup.
- **`internal/proxy`** → core proxying logic, load balancer, service discovery.
- **`internal/auth`** → authentication/authorization (JWT, OAuth2, API key).
- **`internal/security`** → request safety (rate limiting, CORS, circuit breakers).
- **`internal/observability`** → logging, metrics, tracing.
- **`internal/transform`** → request/response modifications.
- **`pkg/plugins`** → for future extensibility (e.g., custom policies).
- **`test`** → separation of unit and integration tests for maintainability.  






# Sample Applications for Testing API Gateway

To test the API Gateway, we will build two simple Go microservices with dummy endpoints.  
Each service will expose **5 REST endpoints** and return JSON responses.

---

## 1. User Service

### Purpose
Provides user-related data for testing routing, authentication, and transformations.

### Requirements
- **Endpoints**
    - `GET /users` → Returns a list of users
    - `GET /users/{id}` → Returns details of a specific user
    - `POST /users` → Creates a new user (dummy response)
    - `PUT /users/{id}` → Updates a user (dummy response)
    - `DELETE /users/{id}` → Deletes a user (dummy response)

- **Implementation Notes**
    - Use `net/http` with a simple router (e.g., Chi, Gorilla Mux, or just default mux).
    - Respond with JSON (`Content-Type: application/json`).
    - Hardcode/dummy data, no DB needed.

---

## 2. Product Service

### Purpose
Provides product-related data for testing API Gateway load balancing, rate limiting, and response transformations.

### Requirements
- **Endpoints**
    - `GET /products` → Returns a list of products
    - `GET /products/{id}` → Returns details of a specific product
    - `POST /products` → Creates a new product (dummy response)
    - `PUT /products/{id}` → Updates a product (dummy response)
    - `DELETE /products/{id}` → Deletes a product (dummy response)

- **Implementation Notes**
    - Use `net/http` with the same router style as User Service.
    - Respond with JSON (`Content-Type: application/json`).
    - Dummy/hardcoded product data, no DB needed.

---

## 3. General Requirements
- Both applications should:
    - Run on **different ports** (e.g., `:8081` for User Service, `:8082` for Product Service).
    - Return clear, structured JSON for easy Gateway logging and transformation.
    - Be simple enough for unit & integration tests with the Gateway.
    - Include minimal logging (`fmt.Println` or structured logging if desired).

---

## 4. Usage
- Start **User Service** on port `8081`.
- Start **Product Service** on port `8082`.
- Configure the **API Gateway** to route:
    - `/api/users/*` → User Service
    - `/api/products/*` → Product Service


---
## Testing the API Gateway Applications
### 🟢 Start the Services

```sh
# Start User Service (port 8081)
go run ./cmd/usersvc

# Start Inventory/Product Service (port 8082)
go run ./cmd/inventorysvc

GATEWAY_CONFIG=/home/mr/GolandProjects/api-gateway/config.yml go run .

curl -H "X-API-Key: DEV_KEY_12345" http://localhost:8081/users
curl -H "X-API-Key: DEV_KEY_12345" http://localhost:8081/inventory
curl http://localhost:8081/metrics