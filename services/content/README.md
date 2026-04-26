# Content Service

gRPC content service. Stores links and notes with automatic metadata scraping (Microlink). Embeddings are processed asynchronously via RabbitMQ → Worker → Qdrant.

## Setup

```bash
cd services/content
go mod tidy
```

Create `.env`:
```
MONGO_URI=<your_mongodb_uri>
DB_NAME=<your_db_name>
GRPC_PORT=50052
QDRANT_ADDR=localhost:6334
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
OLLAMA_URL=http://localhost:11434
```

Requires Qdrant, RabbitMQ, and Ollama running.

## Run

```bash
# Content service
go run cmd/main.go

# Embedding worker (separate terminal)
go run cmd/worker/main.go
```

## How it works

```
AddContent → scrape metadata → store in MongoDB → publish to RabbitMQ → return
                                                          ↓
                                                   Worker (background)
                                                   embed via Ollama → store in Qdrant
```

Circuit breakers protect Microlink, Ollama, and Qdrant calls. Only network-level failures trip the breaker — application errors (e.g., Microlink rejecting a URL) do not.

## gRPC API

| RPC | Description |
|-----|-------------|
| AddContent | Store a link (with scraping) or note |
| GetContents | Get all content for a user |
| DeleteContent | Delete content by ID (also removes from Qdrant) |

## Test

```bash
# Add link
grpcurl -plaintext -d '{"userId":"user1","contentType":"link","content":"https://www.youtube.com/watch?v=dQw4w9WgXcQ"}' localhost:50052 content.ContentService/AddContent

# Add note
grpcurl -plaintext -d '{"userId":"user1","contentType":"note","title":"My note","content":"Some text"}' localhost:50052 content.ContentService/AddContent

# Get contents
grpcurl -plaintext -d '{"userId":"user1"}' localhost:50052 content.ContentService/GetContents

# Delete
grpcurl -plaintext -d '{"id":"<id>","userId":"user1"}' localhost:50052 content.ContentService/DeleteContent
```
