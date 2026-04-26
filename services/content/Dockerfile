FROM golang:1.24-alpine AS builder
ENV GOTOOLCHAIN=auto
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /content ./cmd/main.go

FROM alpine:3.20
COPY --from=builder /content /content
EXPOSE 50052
CMD ["/content"]
