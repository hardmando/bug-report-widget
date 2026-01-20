# Bug Report Widget SaaS

A scalable SaaS platform for collecting, managing, and tracking bug reports from embedded widgets.

## Architecture

The system is composed of several microservices orchestrated via Docker Compose:

| Service | Port (Host) | Description |
|---------|-------------|-------------|
| **Dashboard** | `3001` | Next.js admin interface for users to manage reports and API keys. |
| **Auth Service** | `8081` | Handles user registration, login, JWT issuance, and API Key management. |
| **Ingestion Service** | `8082` | High-throughput public endpoint for receiving bug reports from widgets. Pushes to RabbitMQ. |
| **Bug Service** | `8083` | Consumes reports from RabbitMQ and manages bug lifecycle/storage. |
| **Postgres** | `5432` | Primary relational database. |
| **Redis** | `6379` | Caching layer. |
| **RabbitMQ** | `5672` | Message broker for decoupling ingestion from processing. |

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.22+ (for local development without Docker)
- Node.js 20+ (for Dashboard local dev)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd bug-report-widget
   ```

2. **Environment Setup**
   Copy the example environment file:
   ```bash
   cp .env.example .env.local
   ```
   *Note: Adjust variables in `.env.local` if necessary.*

3. **Start the Stack**
   Run the entire system with Docker Compose:
   ```bash
   docker-compose up --build
   ```

   - Dashboard will be available at [http://localhost:3001](http://localhost:3001)
   - API endpoints will be reachable via their respective ports (e.g., Auth at `8081`).

## Development

### Directory Structure

- `/cmd` - Entry points for Go services (`auth-service`, `ingestion-service`, `bug-service`).
- `/internal` - Shared application logic and internal packages.
- `/dashboard` - Next.js frontend application.
- `/migrations` - SQL database migrations.
- `/docker` - Dockerfiles for services.

### API Usage

#### Ingestion (Public)

**POST** `http://localhost:8082/ingest/bugs`
Headers: `X-API-Key: <your_tenant_api_key>`
Body:
```json
{
  "description": "Something went wrong",
  "metadata": {
    "url": "https://example.com",
    "browser": "Chrome"
  }
}
```

#### Auth (Internal/Admin)

- **POST** `/register` - Create a new user/tenant.
- **POST** `/login` - Get JWT.
- **POST** `/api-keys` - Generate new API Keys (Protected).

## Testing

Run Go tests:
```bash
go test ./...
```