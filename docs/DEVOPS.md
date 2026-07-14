# DevOps Guide

## Local development

1. Copy `.env.example` or `setup.txt` values into `.env` and set:
   - `POSTGRES_USER`
   - `POSTGRES_PASSWORD`
   - `POSTGRES_DB`
   - `POSTGRES_HOST`
   - `POSTGRES_PORT`
   - `DIST_DIR` (default: `dist`)
   - `UPLOAD_DIR` (optional, default: `./uploads`)

2. Install Go dependencies:

```powershell
cd <repo-root>
go mod tidy
```

3. Generate Swagger docs if needed:

```powershell
swag init
```

4. Run the server locally:

```powershell
go run main.go
```

## Build and verify

```powershell
go test ./...
gofmt -w .
go build -o bin/industrial-automation .
```

## Docker and containerized development

Use the existing `docker-compose.yml` to start the backend with mounted static UI content.

```powershell
docker compose up --build
```

The service listens on `8022` by default and forwards to container port `8080`.

## Key operational notes

- The backend exposes a health route at `/api/health` and readiness at `/api/health/ready`.
- Static front-end assets are served from `DIST_DIR` or the default `dist` directory.
- Database connection pooling is configured in `db/config.go` with sensible defaults.
- Keep environment-specific settings separate using `.env.dev`, `.env.uat`, and `.env.prod`.
