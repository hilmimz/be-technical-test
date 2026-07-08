# ---- Build stage ----
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd

# ---- Final stage ----
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/migration ./migration

EXPOSE 8080

CMD ["./api"]
