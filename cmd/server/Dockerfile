# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/server ./cmd/server/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/server /app/server
COPY --from=builder /app/data.db /app/data.db


RUN mkdir -p /app/data && chmod -R 777 /app

EXPOSE 8080
CMD ["/app/server", "--db", "/app/data.db", "--listen", "8080"]