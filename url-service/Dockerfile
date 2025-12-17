FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copy all source code first
COPY . .

# Download and tidy modules
RUN go mod tidy
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /url-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /url-service .

EXPOSE 8082

CMD ["./url-service"]
