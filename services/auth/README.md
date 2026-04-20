# Auth Service

gRPC authentication service for second-brain. Handles user signup, login, token validation, and token refresh using JWT and MongoDB.

## Setup

1. Install dependencies:
```bash
cd services/auth
go mod tidy
```

2. Create a `.env` file in `services/auth/`:
```
MONGO_URI=<your_mongodb_uri>
DB_NAME=<your_db_name>
JWT_SECRET=<min_32_char_secret>
GRPC_PORT=50051
```

3. Generate proto files (from repo root):
```bash
protoc --go_out=services/auth/pkg/pb --go_opt=module=github.com/jram17/second-brain/services/auth/pkg/pb --go-grpc_out=services/auth/pkg/pb --go-grpc_opt=module=github.com/jram17/second-brain/services/auth/pkg/pb proto/auth/auth.proto
```

## Run

```bash
cd services/auth
go run cmd/main.go
```

## Test with grpcurl

Install grpcurl:
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

**Signup:**
```bash
grpcurl -plaintext -d '{"username":"testuser","email":"test@test.com","password":"password123"}' localhost:50051 auth.AuthService/Signup
```

**Login:**
```bash
grpcurl -plaintext -d '{"email":"test@test.com","password":"password123"}' localhost:50051 auth.AuthService/Login
```

**Validate Token:**
```bash
grpcurl -plaintext -d '{"accessToken":"<token>"}' localhost:50051 auth.AuthService/ValidateToken
```

**Refresh Token:**
```bash
grpcurl -plaintext -d '{"refreshToken":"<token>"}' localhost:50051 auth.AuthService/RefreshToken
```
