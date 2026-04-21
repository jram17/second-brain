# Content Service

gRPC content service for second-brain. Handles storing, retrieving, and deleting user content (links, notes) with automatic metadata scraping via Microlink.

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
```

3. Generate proto files (from repo root):
```bash
protoc --go_out=services/content/pkg/pb --go_opt=module=github.com/jram17/second-brain/services/content/pkg/pb --go-grpc_out=services/content/pkg/pb --go-grpc_opt=module=github.com/jram17/second-brain/services/content/pkg/pb proto/content/content.proto
```

## Run

```bash
cd services/content
go run cmd/main.go
```

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

## Via API Gateway

Requires auth service and gateway running. All content routes are protected.

```bash
# Add content
curl -X POST http://localhost:8080/api/content -H "Content-Type: application/json" -H "Authorization: Bearer <token>" -d '{"contentType":"link","content":"https://github.com"}'

# Get contents
curl -X GET http://localhost:8080/api/content -H "Authorization: Bearer <token>"

# Delete content
curl -X DELETE http://localhost:8080/api/content/<id> -H "Authorization: Bearer <token>"
```
