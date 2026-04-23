# Second Brain

A microservices-based personal knowledge management system. Store links, notes, and YouTube videos . then query your saved content using RAG (Retrieval-Augmented Generation) powered by local LLM.

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
                 ┌──────▼──┐ ┌──▼────┐ ┌▼──────────┐
                 │  Auth   │ │Content│ │  Query     │
                 │  :50051 │ │:50052 │ │  :8000     │
                 │ (Go)    │ │(Go)   │ │ (Python)   │
                 └────┬────┘ └┬────┬─┘ └──┬─────┬──┘
                      │       │    │      │     │
                 ┌────▼───────▼┐ ┌─▼──────▼┐ ┌─▼──────┐
                 │  MongoDB    │ │ Qdrant  │ │ Ollama │
                 │             │ │ :6333   │ │ :11434 │
                 └─────────────┘ └─────────┘ └────────┘
```

## Tech Stack

- **Go** — Auth, Content, Gateway services (gRPC + Gin)
- **Python** — Query service (FastAPI)
- **MongoDB** — User and content storage
- **Qdrant** — Vector database for semantic search
- **Ollama** — Local LLM (llama3) and embeddings (nomic-embed-text)
- **gRPC** — Inter-service communication
- **JWT** — Authentication

## Prerequisites

- Go 1.22+
- Python 3.10+
- Docker (for Qdrant)
- Ollama
- MongoDB Atlas account (or local MongoDB)
- protoc (Protocol Buffers compiler)

## Quick Start

### 1. Start infrastructure

```bash
# Qdrant
docker run -d -p 6333:6333 -p 6334:6334 qdrant/qdrant

# Ollama
ollama serve
ollama pull nomic-embed-text
ollama pull llama3
```

### 2. Start services

```bash
# Terminal 1 — Auth
cd services/auth
go run cmd/main.go

# Terminal 2 — Content
cd services/content
go run cmd/main.go

# Terminal 3 — Query
cd services/query
source venv/bin/activate
python -m uvicorn main:app --port 8000

# Terminal 4 — Gateway
cd services/gateway
go run cmd/main.go
```

### 3. Test

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

# Query your content
curl -X POST http://localhost:8080/api/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"query":"what videos do I have saved?"}'
```

## Services

| Service | Port | Language | Description |
|---------|------|----------|-------------|
| [Auth](services/auth/) | 50051 | Go | User authentication (signup, login, JWT) |
| [Content](services/content/) | 50052 | Go | Content CRUD, scraping, vector embedding |
| [Query](services/query/) | 8000 | Python | RAG-based querying with Ollama LLM |
| [Gateway](services/gateway/) | 8080 | Go | REST API gateway with auth middleware |
