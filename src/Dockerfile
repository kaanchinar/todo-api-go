# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /todo-app \
    .

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

RUN adduser --disabled-password --no-create-home appuser

WORKDIR /app

COPY --from=builder /todo-app .

USER appuser

EXPOSE 8080

ENTRYPOINT ["./todo-app"]
