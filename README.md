# 🔗 URL Shortener - Microservices Project

A modern URL shortener built with Go microservices, Redis, and PostgreSQL.

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Frontend  │────▶│   Nginx     │────▶│   Services  │
│   (React)   │     │   Gateway   │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
                                              │
                    ┌─────────────────────────┼─────────────────────────┐
                    │                         │                         │
              ┌─────▼─────┐           ┌───────▼───────┐         ┌───────▼───────┐
              │   User    │           │     URL       │         │    Stats      │
              │  Service  │           │   Service     │         │   Service     │
              │  :8081    │           │   :8082       │         │   :8083       │
              └─────┬─────┘           └───────┬───────┘         └───────┬───────┘
                    │                         │                         │
                    │                         │      ┌──────────┐       │
                    │                         └─────▶│  Redis   │◀──────┘
                    │                                │  Pub/Sub │
                    │                                └──────────┘
                    │                                      │
              ┌─────▼──────────────────────────────────────▼─────┐
              │                   PostgreSQL                      │
              └───────────────────────────────────────────────────┘
```

## Services

| Service | Port | Description |
|---------|------|-------------|
| User Service | 8081 | User registration & authentication |
| URL Service | 8082 | URL shortening & redirection |
| Stats Service | 8083 | Click tracking & analytics |
| Frontend | 3000 | Web UI |

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)

### Run with Docker Compose

```bash
# Start all services
docker-compose up --build

# Access the application
open http://localhost:3000
```

### API Endpoints

#### User Service (http://localhost:8081)
```bash
# Register
curl -X POST http://localhost:8081/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123", "name": "Test User"}'

# Login
curl -X POST http://localhost:8081/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}'
```

#### URL Service (http://localhost:8082)
```bash
# Shorten URL
curl -X POST http://localhost:8082/api/urls \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"original_url": "https://github.com/very/long/url"}'

# Redirect (use in browser)
curl -L http://localhost:8082/abc123
```

#### Stats Service (http://localhost:8083)
```bash
# Get URL stats
curl http://localhost:8083/api/stats/abc123 \
  -H "Authorization: Bearer <token>"
```

## Tech Stack

- **Language**: Go 1.21
- **Framework**: Gin
- **Database**: PostgreSQL
- **Cache/Pub-Sub**: Redis
- **Container**: Docker
- **Cloud**: Railway
