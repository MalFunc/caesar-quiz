# ---- Stage 1: build binary Go ----
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o caesar-quiz ./cmd/server

# ---- Stage 2: image runtime minimal ----
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache postgresql-client ca-certificates

COPY --from=builder /app/caesar-quiz ./caesar-quiz
COPY wait-for-db.sh /wait-for-db.sh
RUN chmod +x /wait-for-db.sh

EXPOSE 8080

ENTRYPOINT ["/wait-for-db.sh", "./caesar-quiz"]
