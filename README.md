# Users API
An API to interact with users data.

# Features
- Create, update, search, and delete user records.
- View users individually or in bulk.
- OpenAPI (Swagger) documentation available.
- Built-in tracing (via OpenTelemetry).

## Requirements
- Docker & Docker Compose.
- [Go 1.23.9](https://go.dev/doc/install) and [PostgreSQL 16](https://www.postgresql.org/download/) (if running locally without containers).

## Setup Instructions
### 1. Prepare your environment
Copy the `.env.sample` file to `.env` to set up environment variables:
```sh
cp .env.sample .env
```
### 2. Launch everything with Docker Compose
```bash
docker-compose up --build
```
This will spin up:
- The API server.
- A PostgreSQL database (with user table migrations applied).
- Swagger documentation at: [http://localhost:3001/company/docs/index.html](http://localhost:3001/company/docs/index.html).

## Helpful Commands
Build docs manually:
```bash
make generate-openapi
```
Run tests:
```bash
make test
```
Lint your code:
```bash
make lint
```

## Want to Learn More?
This project is a great entry point to:
- Structuring Go apps using Clean Architecture.
- Writing APIs with Gin.
- Integrating Postgres using sqlc.
- Observability with OpenTelemetry.

## Contributions
Thank you for considering contributing to this project. It's awesome of you!
