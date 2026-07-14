# Deployment Guide

This guide covers deployment promotion for the industrial automation backend.

## Environment strategy

Use separate environment definitions for:
- Development (`DEV`)
- User acceptance testing (`UAT`)
- Production (`PROD`)

Maintain dedicated environment files such as `.env.dev`, `.env.uat`, and `.env.prod`.

## Recommended deployment flow

1. **Prepare the environment**
   - Ensure `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_HOST`, and `POSTGRES_PORT` are set.
   - Set `DIST_DIR` to the built front-end directory.
   - Optionally set `UPLOAD_DIR` for persisted uploads.
   - Optionally set `BUILD_TIME` and `RELEASE_TAG` for runtime metadata.

2. **Build the service**

```powershell
go mod tidy
go test ./...
go build -o bin/industrial-automation .
```

3. **Package the container**

```powershell
docker build -t industrial-automation:latest .
```

4. **Deploy to DEV**
   - Start with a dev database and smaller pool sizes if necessary.
   - Use `docker compose up --build` or platform-specific deployment tooling.
   - Verify `/api/health`, `/api/health/ready`, and `/ping`.

5. **Promote to UAT**
   - Use UAT environment variables and DB credentials.
   - Run smoke tests and perform business process validations.
   - Validate file upload, maintenance schedule, vendor and feed workflows.

6. **Release to PROD**
   - Use production database credentials and separate `UPLOAD_DIR` storage.
   - Apply DB migrations before starting the service.
   - Monitor health endpoints and logs.

## Deployment checklist

- [ ] Database migrations applied
- [ ] Swagger docs regenerated if API changed (`swag init`)
- [ ] Application built successfully
- [ ] Health and readiness checks pass
- [ ] Static assets available in `DIST_DIR`
- [ ] Service listening on expected port

## Production considerations

- Use a managed Postgres service or a dedicated Postgres cluster.
- Ensure backups and disaster recovery are in place.
- Run the service behind a load balancer or API gateway for gradual rollout.
- Use environment variables for secrets, not baked into images.
