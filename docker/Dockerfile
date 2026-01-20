# Dev stage (hot reload)
FROM golang:1.24-alpine AS dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["go", "run", "./cmd/{SERVICE_NAME}"]  # Replaced at runtime

# Prod stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG SERVICE_NAME
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/${SERVICE_NAME} ./cmd/${SERVICE_NAME}

FROM alpine:latest AS prod
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/bin .
CMD ["./bin"]
