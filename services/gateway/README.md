# API Gateway

REST API gateway. Translates HTTP → gRPC (auth, content) and HTTP → HTTP (query). JWT middleware on protected routes. Circuit breakers on all downstream calls.

## Architecture

```
                        Client (HTTP)
                                │
                                ▼
                         ┌───────────────┐
                         │    Gateway    │
                         │    :8080      │
                         └──┬──┬──┬──┬───┘
                            │  │  │  │
              ┌─────────────┘  │  │  └──────────────┐
              │                │  │                 │
              ▼                │  │                 ▼
        ┌──────────┐           │  │           ┌────────────┐
        │  Redis   │ cache     │  │           │   Query    │
        │  :6379   │◄─────┐    │  │    HTTP   │   :8000    │
        └──────────┘      │    │  │   ──────► │  (Python)  │
                          │    │  │           └────────────┘
                          │    │  │
              JWT validate│    │  │
              (cache miss)│    │  │
                          │    │  │
                          ▼    ▼  ▼
                     ┌─────┐  ┌───────┐
                     │Auth:│  │Content│
                     │50051│  │:50052 │
                     └─────┘  └───────┘
                      gRPC      gRPC

  Request Flow (protected routes):
  ─────────────────────────────────
  1. Client sends request with Bearer token
  2. Auth middleware checks Redis for cached token
  3. Cache hit  → extract userId, skip gRPC call
     Cache miss → validate via Auth gRPC, cache result (5 min TTL)
  4. Route handler calls downstream service via circuit breaker
  5. Response returned to client
```

## Setup

```bash
cd services/gateway
go mod tidy
```

Create `.env`:
```
AUTH_SERVICE_ADDR=localhost:50051
CONTENT_SERVICE_ADDR=localhost:50052
QUERY_SERVICE_URL=http://localhost:8000
HTTP_PORT=8080
REDIS_ADDR=localhost:6379
```

## Run

```bash
go run cmd/main.go
```

Requires auth, content, query services, and worker to be running. Redis is optional — if unavailable, JWT validation falls back to gRPC.

## Endpoints

### Public

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/auth/signup | Register new user |
| POST | /api/auth/login | Login, returns tokens |
| POST | /api/auth/refresh | Refresh access token |

### Protected (`Authorization: Bearer <token>`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/content | Add content (link or note) |
| GET | /api/content | Get all user content |
| DELETE | /api/content/:id | Delete content by ID |
| POST | /api/query | Query content using RAG |

## Resilience

- Circuit breakers on auth, content, and query service calls. Only infrastructure failures (Unavailable, DeadlineExceeded) trip the breaker — business errors (e.g., "user already exists") do not.
- **Redis Caching** — JWT validation results cached for 5 min (~25x faster). Falls back to gRPC if Redis is down.
- Returns `503` only when a downstream service is actually unavailable (circuit open)
- Proper gRPC error forwarding for business logic errors (400 with real message)
