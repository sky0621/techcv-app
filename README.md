# techcv-app

Implementation repository for the TechCV product.

## Stack

- Frontend: Next.js + TypeScript + Tailwind CSS + shadcn/ui
- Backend: Go + Chi + sqlc
- Database: MySQL
- Auth: SessionAuth

## Structure

```text
frontend/
backend/
```

## Frontend

```sh
cd frontend
pnpm install
pnpm dev
```

`/api/*` requests are proxied to the backend with `BACKEND_ORIGIN` (default: `http://127.0.0.1:8080`).

## Backend

```sh
cd backend
export MYSQL_DSN="root:password@tcp(127.0.0.1:3306)/techcv_app?parseTime=true"
go run ./cmd/server
```

By default, the backend serves `GET /healthz` on port `8080`.

## Local development

Trust the project config once before running tasks:

```sh
mise trust
```

Start MySQL with Docker Compose. The schema is loaded from `backend/migrations/schema.sql` on first boot.

```sh
mise run mysql-up
```

Then start the backend and frontend in separate terminals:

```sh
mise run backend-up
```

```sh
mise run frontend-up
```

Open `http://127.0.0.1:3000`. The profile screen now:

- loads profile data with `GET /api/profile`
- saves edits with `PUT /api/profile`
- reads the saved values back from MySQL through the backend

Stop MySQL when you are done:

```sh
mise run frontend-down
mise run backend-down
mise run mysql-down
```

This stops only the `mysql` service defined in `docker-compose.yml`.

### Schema management

The backend schema is managed with `mysqldef` using `backend/migrations/schema.sql`.

```sh
export MYSQL_HOST=127.0.0.1
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASSWORD=password
export MYSQL_DATABASE=techcv_app
mise run schema-dry-run
mise run schema-apply
```

### SQL management

Repository SQL is managed with `sqlc`.

```sh
mise run sqlc-generate
```

### OpenAPI management

Profile API schema is extracted from `techcv-design/openapi.yaml` and maintained under `backend/openapi/openapi.yaml`.
Backend server/models are generated from that file with `oapi-codegen`.

```sh
mise run oapi-generate
```

## Dependency updates

This repository includes [renovate.json](air-file://74mpjbg0chpcohbk3d4o/Users/sky0621/work/github.com/sky0621/techcv-products/techcv-app/renovate.json?type=file&root=%252F) for automated dependency updates.
Enable the Renovate GitHub App for the repository to start receiving update PRs for Go modules, npm packages, Docker Compose images, and mise-managed tools.
