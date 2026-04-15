FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN go build -o wishlist-api ./cmd/api


FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/wishlist-api ./wishlist-api
COPY --from=builder /app/.env.example ./.env.example
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./wishlist-api"]
