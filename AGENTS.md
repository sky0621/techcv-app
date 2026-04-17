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

## Working Style
- Keep initial implementations minimal and easy to evolve.
- Favor explicit mappings between layers over implicit framework magic.
