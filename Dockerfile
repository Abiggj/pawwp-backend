# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build both binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o seed ./cmd/seed/main.go

# Stage 2: Run the binary
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/seed .
COPY --from=builder /app/migrations ./migrations
COPY .env .

EXPOSE 8080

CMD ["./main"]
