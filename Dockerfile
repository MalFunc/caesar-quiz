
# Stage 1: Build static Go binary
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod dan go.sum dulu (biar cache Docker efisien)
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .
RUN go get github.com/gin-contrib/cors@latest && go mod tidy
# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o caesar-quiz ./cmd/server

# Stage 2: Run minimal image
FROM alpine:3.19

WORKDIR /app

# Install postgresql-client for pg_isready
RUN apk add --no-cache postgresql-client

# Copy binary dari builder
COPY --from=builder /app/caesar-quiz ./caesar-quiz

# Copy wait-for-db script
COPY wait-for-db.sh /wait-for-db.sh
RUN chmod +x /wait-for-db.sh

# Expose port (sesuai yang dipakai di main.go)
EXPOSE 8080

# Start app with wait-for-db
ENTRYPOINT ["/wait-for-db.sh", "db", "./caesar-quiz"]
