# API Gateway

REST API gateway for second-brain. Translates HTTP requests to gRPC calls (auth, content) and HTTP calls (query). Handles JWT authentication middleware for protected routes.

## Setup

1. Install dependencies:
```bash
cd services/gateway
go mod tidy
```

2. Create a `.env` file in `services/gateway/`:
```
AUTH_SERVICE_ADDR=localhost:50051
CONTENT_SERVICE_ADDR=localhost:50052
QUERY_SERVICE_URL=http://localhost:8000
HTTP_PORT=8080
```

## Run

```bash
cd services/gateway
go run cmd/main.go
```

Requires auth, content, and query services to be running.

## API Endpoints

### Public
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/auth/signup | Register new user |
| POST | /api/auth/login | Login, returns tokens |
| POST | /api/auth/refresh | Refresh access token |

### Protected (requires `Authorization: Bearer <token>`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/content | Add content (link or note) |
| GET | /api/content | Get all user content |
| DELETE | /api/content/:id | Delete content by ID |
| POST | /api/query | Query content using RAG |
