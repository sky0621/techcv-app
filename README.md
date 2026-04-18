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

Start MySQL with Docker Compose. The schema is loaded from `backend/migrations/schema.sql` on first boot.

```sh
cd backend
make mysql-up
```

Then start the backend and frontend in separate terminals:

```sh
cd backend
export MYSQL_DSN="root:password@tcp(127.0.0.1:3306)/techcv_app?parseTime=true"
go run ./cmd/server
```

```sh
cd frontend
pnpm install
pnpm dev
```

Open `http://127.0.0.1:3000`. The profile screen now:

- loads profile data with `GET /api/profile`
- saves edits with `PUT /api/profile`
- reads the saved values back from MySQL through the backend

Stop MySQL when you are done:

```sh
cd backend
make mysql-down
```

This stops only the `mysql` service defined in `docker-compose.yml`.

### Schema management

The backend schema is managed with `mysqldef` using `backend/migrations/schema.sql`.

```sh
cd backend
export MYSQL_HOST=127.0.0.1
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASSWORD=password
export MYSQL_DATABASE=techcv_app
make schema-dry-run
make schema-apply
```

### SQL management

Repository SQL is managed with `sqlc`.

```sh
cd backend
make sqlc-generate
```
