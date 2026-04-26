# Second Brain

A microservices-based personal knowledge management system. Store links, notes, and YouTube videos — then query your saved content using RAG (Retrieval-Augmented Generation) powered by local LLM.

## Architecture

```
                         ┌──────────────┐
                         │   Client     │
                         └──────┬───────┘
                                │ HTTP
                         ┌──────▼───────┐
                         │   Gateway    │ (Go + Gin)
                         │   :8080      │
                         └──┬───┬───┬───┘
                   gRPC ┌───┘   │   └───┐ HTTP
                        │       │       │
                 ┌──────▼──┐ ┌──▼────┐ ┌▼───────────┐
                 │  Auth   │ │Content│ │  Query     │
                 │  :50051 │ │:50052 │ │  :8000     │
                 │ (Go)    │ │(Go)   │ │ (Python)   │
                 └────┬────┘ └┬──┬─┬─┘ └──┬─────┬───┘
                      │       │  │ │      │     │
                 ┌────▼───────▼┐ │ │  ┌───▼┐ ┌──▼─────┐
                 │  MongoDB    │ │ │  │Qdr.│ │ Ollama │
                 │             │ │ │  │6333│ │ :11434 │
                 └─────────────┘ │ │  └────┘ └────────┘
                          ┌──────┘ │
                          │  ┌─────┘
                    ┌─────▼──▼──┐    ┌───────────┐
                    │ RabbitMQ  │───►│  Worker   │──► Qdrant
                    │  :5672    │    │(embedding)│
                    └───────────┘    └───────────┘
```

## Tech Stack

- **Go** — Auth, Content, Gateway services (gRPC + Gin)
- **Python** — Query service (FastAPI)
- **MongoDB** — User and content storage
- **Qdrant** — Vector database for semantic search
- **Ollama** — Local LLM (llama3) and embeddings (nomic-embed-text)
- **RabbitMQ** — Async message queue for embedding pipeline
- **Redis** — JWT validation caching in gateway
- **gRPC** — Inter-service communication
- **JWT** — Authentication
- **Circuit Breaker** — Resilience pattern (gobreaker + pybreaker)
- **Docker Compose** — Container orchestration

## Quick Start (Docker)

```bash
git clone https://github.com/jram17/second-brain.git
cd second-brain
cp .env.example .env   # edit with your values
docker compose up -d --build
```

> **Important:** Ollama must be running on the host and listening on all interfaces:
> ```bash
> OLLAMA_HOST=0.0.0.0 ollama serve
> ```
> Or for systemd: add `Environment="OLLAMA_HOST=0.0.0.0"` to the service override.

Pull Ollama models (one-time):
```bash
ollama pull nomic-embed-text
ollama pull llama3
```

Dashboards:
- RabbitMQ: `http://localhost:15672` (guest/guest)
- Qdrant: `http://localhost:6333/dashboard`
- Dozzle (logs): `http://localhost:9999`

## Quick Start (Local)

### 1. Prerequisites

- Go 1.22+
- Python 3.10+
- Docker (for Qdrant, RabbitMQ)
- Ollama
- MongoDB (Atlas or local)
- protoc (Protocol Buffers compiler)

### 2. Start infrastructure

```bash
docker run -d -p 6333:6333 -p 6334:6334 qdrant/qdrant
docker run -d -p 5672:5672 -p 15672:15672 rabbitmq:management
ollama serve
ollama pull nomic-embed-text
ollama pull llama3
```

### 3. Start services

```bash
# Terminal 1 — Auth
cd services/auth && go run cmd/main.go

# Terminal 2 — Content
cd services/content && go run cmd/main.go

# Terminal 3 — Embedding Worker
cd services/content && go run cmd/worker/main.go

# Terminal 4 — Query
cd services/query && source venv/bin/activate && python -m uvicorn main:app --port 8000

# Terminal 5 — Gateway
cd services/gateway && go run cmd/main.go
```

### 4. Test

```bash
# Signup
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@test.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"password123"}'

# Add content (use token from login)
curl -X POST http://localhost:8080/api/content \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"contentType":"link","content":"https://www.youtube.com/watch?v=dQw4w9WgXcQ"}'

# Query (wait ~10s for embedding)
curl -X POST http://localhost:8080/api/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"query":"what videos do I have saved?"}'
```

## Services

| Service | Port | Language | Description |
|---------|------|----------|-------------|
| [Auth](services/auth/) | 50051 | Go | User authentication (signup, login, JWT) |
| [Content](services/content/) | 50052 | Go | Content CRUD, scraping, async embedding via RabbitMQ |
| [Worker](services/content/) | — | Go | Consumes embedding jobs, stores vectors in Qdrant |
| [Query](services/query/) | 8000 | Python | RAG-based querying with Qdrant + Ollama LLM |
| [Gateway](services/gateway/) | 8080 | Go | REST API gateway with auth middleware + circuit breakers |

## Resilience

- **Circuit Breakers** on all external calls (gobreaker for Go, pybreaker for Python). Only infrastructure failures (network errors, timeouts) trip the breaker — application errors (e.g., "user already exists") do not.
- **Redis Caching** — JWT validation results cached in Redis (~25x faster on subsequent requests). Falls back to gRPC if Redis is unavailable.
- **Async Embedding** via RabbitMQ — content creation returns immediately, embedding happens in background
- **Timeouts** on all HTTP and gRPC calls
- **Health Checks** in Docker Compose for dependency ordering
