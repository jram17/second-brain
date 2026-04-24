# Content Service

gRPC content service for second-brain. Handles storing, retrieving, and deleting user content (links, notes) with automatic metadata scraping via Microlink. Embeddings are processed asynchronously via RabbitMQ worker.

## Setup

1. Install dependencies:
```bash
cd services/content
go mod tidy
```

2. Create a `.env` file in `services/content/`:
```
MONGO_URI=<your_mongodb_uri>
DB_NAME=<your_db_name>
GRPC_PORT=50052
QDRANT_ADDR=localhost:6334
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

3. Ensure infrastructure is running:
```bash
ollama serve
ollama pull nomic-embed-text
docker run -d -p 6333:6333 -p 6334:6334 qdrant/qdrant
docker run -d -p 5672:5672 -p 15672:15672 rabbitmq:management
```

## Run

```bash
# Content service
cd services/content
go run cmd/main.go

# Embedding worker (separate terminal)
cd services/content
go run cmd/worker/main.go
```

RabbitMQ dashboard: `http://localhost:15672` (guest/guest)

## How it works

1. Client adds content → scrape metadata → store in MongoDB → publish to RabbitMQ → return immediately
2. Worker picks up message → embed text via Ollama → store vector in Qdrant
3. Circuit breakers protect all external calls (Microlink, Ollama, Qdrant)

## Test with grpcurl

**Add a link:**
```bash
grpcurl -plaintext -d '{"userId":"testuser123","contentType":"link","content":"https://www.youtube.com/watch?v=dQw4w9WgXcQ"}' localhost:50052 content.ContentService/AddContent
```

**Add a note:**
```bash
grpcurl -plaintext -d '{"userId":"testuser123","contentType":"note","title":"My note","content":"Some text here"}' localhost:50052 content.ContentService/AddContent
```

**Get all contents:**
```bash
grpcurl -plaintext -d '{"userId":"testuser123"}' localhost:50052 content.ContentService/GetContents
```

**Delete content:**
```bash
grpcurl -plaintext -d '{"id":"<content_id>","userId":"testuser123"}' localhost:50052 content.ContentService/DeleteContent
```
