FROM golang:1.22-alpine AS builder
WORKDIR /app

# GONOSUMDB=* разрешает загрузку модулей без проверки checksum DB
# (нужно из-за pseudo-version pgservicefile, убранного с прокси)
ENV GONOSUMDB=*

COPY backend/go.mod ./
RUN go mod download

COPY backend/ .
RUN go build -o /api ./cmd/api/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /api .
COPY migrations/ ./migrations/
ENV MIGRATIONS_DIR=/app/migrations
EXPOSE 8080
CMD ["/app/api"]
