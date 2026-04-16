# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app/backend

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates bash tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/backend/migrations ./migrations
COPY --from=builder /app/backend/.env .

ENV PORT=8080
EXPOSE 8080

CMD ["./main"]
