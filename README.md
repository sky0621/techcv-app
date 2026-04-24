# techcv-app

Implementation repository for the TechCV product.

## Stack

- Frontend: Next.js + TypeScript + Tailwind CSS + shadcn/ui
- Backend: Go + Chi + sqlc
- Database: SQLite
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
export SQLITE_DSN="data/techcv.db"
go run ./cmd/server
```

By default, the backend stores data in `backend/data/techcv.db` and serves `GET /healthz` on port `8080`.

## Local development

Trust the project config once before running tasks:

```sh
mise trust
```

Start the backend and frontend in separate terminals:

```sh
mise run backend-up
```

```sh
mise run frontend-up
```

Open `http://127.0.0.1:3000`. The profile screen now:

- loads profile data with `GET /api/profile`
- saves edits with `PUT /api/profile`
- reads the saved values back from SQLite through the backend

Stop the local servers when you are done:

```sh
mise run frontend-down
mise run backend-down
```

### Schema management

The backend schema is managed in `backend/migrations/schema.sql` and applied to SQLite on backend startup.

### SQL management

Repository SQL is managed with `sqlc`.

```sh
mise run sqlc-generate
```

### Backend tests

Run backend Go tests with:

```sh
mise run backend-test
```

## Backend internal layout

`backend/internal` is organized by technical layer, not by assumed domain boundaries.

- `internal/domain`
- `internal/handler`
- `internal/repository`
- `internal/usecase`
- `internal/shared`

## Dependency updates

This repository includes [renovate.json](air-file://74mpjbg0chpcohbk3d4o/Users/sky0621/work/github.com/sky0621/techcv-products/techcv-app/renovate.json?type=file&root=%252F) for automated dependency updates.
Enable the Renovate GitHub App for the repository to start receiving update PRs for Go modules, npm packages, Docker Compose images, and mise-managed tools.
