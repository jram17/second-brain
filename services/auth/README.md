# Auth Service

gRPC authentication service. Handles signup, login, token validation, and refresh using JWT + MongoDB.

## Setup

```bash
cd services/auth
go mod tidy
```

Create `.env`:
```
MONGO_URI=<your_mongodb_uri>
DB_NAME=<your_db_name>
JWT_SECRET=<min_32_char_secret>
GRPC_PORT=50051
```

## Run

```bash
go run cmd/main.go
```

## gRPC API

| RPC | Description |
|-----|-------------|
| Signup | Register new user, returns access + refresh tokens |
| Login | Authenticate user, returns tokens |
| ValidateToken | Verify access token, returns userId |
| RefreshToken | Issue new access token from refresh token |

## Test

```bash
# Signup
grpcurl -plaintext -d '{"username":"testuser","email":"test@test.com","password":"password123"}' localhost:50051 auth.AuthService/Signup

# Login
grpcurl -plaintext -d '{"email":"test@test.com","password":"password123"}' localhost:50051 auth.AuthService/Login

# Validate
grpcurl -plaintext -d '{"accessToken":"<token>"}' localhost:50051 auth.AuthService/ValidateToken

# Refresh
grpcurl -plaintext -d '{"refreshToken":"<token>"}' localhost:50051 auth.AuthService/RefreshToken
```
