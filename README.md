# techcv-app

Implementation repository for the TechCV product.

## Stack

- Frontend: Next.js + TypeScript + Tailwind CSS + shadcn/ui
- Backend: Go + Chi + sqlc
- Database: PostgreSQL
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

## Backend

```sh
cd backend
go run ./cmd/server
```

By default, the backend serves `GET /healthz` on port `8080`.
