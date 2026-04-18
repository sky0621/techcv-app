# AGENTS.md

## Scope
- These instructions apply to the entire `techcv-app` repository unless a deeper `AGENTS.md` overrides them.

## Layering Rules
- Keep transport concerns in the transport layer.
- JSON, HTTP, and serialization tags belong in handler-layer DTOs, not in domain models.
- Do not add `json`, `form`, or similar transport tags to domain or usecase-layer structs.
- Use handler-layer request/response structs to translate between HTTP payloads and usecase inputs/outputs.

## Backend Structure
- Prefer separating `handler`, `usecase`, and `repository` once an API is expected to grow.
- Keep domain models focused on business meaning, not framework or protocol details.
- Add shared helpers only for clearly cross-cutting concerns.
- Keep database schema changes in `backend/migrations/schema.sql`.
- Manage MySQL schema changes with `mysqldef`-compatible schema files and commands.
- Manage repository SQL in `sqlc` query files rather than inline SQL in repository implementations.
- Treat `backend/sqlc.yaml` and `backend/internal/**/repository/queries/*.sql` as the source of truth for SQL managed by `sqlc`.

## Working Style
- Keep initial implementations minimal and easy to evolve.
- Favor explicit mappings between layers over implicit framework magic.
