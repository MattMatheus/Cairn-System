# Local Control Plane Quickstart

Historical or future-facing reference for a larger AthenaWork control-plane stack.
This is not part of the current internal-beta default path.
Use repo-local `.athena/`, `INTERNAL_BETA.md`, and the platform smoke scripts first.

## Purpose
Bootstrap AthenaWork local control-plane services (API + DB + UI) using Docker Compose.

## Prerequisites
- Docker Desktop or Docker Engine with Compose plugin
- Repo root as current working directory

## Environment Setup
1. Copy environment defaults:
```bash
cp .env.example .env
```
2. Configure one model provider in `.env`:
- Azure OpenAI (preferred):
  - `AZURE_OPENAI_ENDPOINT`
  - `AZURE_OPENAI_API_KEY`
  - `AZURE_OPENAI_DEPLOYMENT_NAME`
  - optional: `AZURE_OPENAI_API_VERSION`
- OR OpenAI:
  - `OPENAI_API_KEY`
  - optional: `OPENAI_MODEL`

## Start
1. Build and start stack:
```bash
docker compose -f docker-compose.local.yml up --build -d
```
2. Check health:
```bash
docker compose -f docker-compose.local.yml ps
curl -fsS http://127.0.0.1:${ATHENAWORK_API_PORT:-8787}/health
curl -fsS http://127.0.0.1:${ATHENAWORK_API_PORT:-8787}/api/v1/read-model/board
curl -fsS http://127.0.0.1:${ATHENAWORK_API_PORT:-8787}/api/v1/read-model/timeline
curl -fsS http://127.0.0.1:${ATHENAWORK_UI_PORT:-8080}
curl -fsS -X POST http://127.0.0.1:${ATHENAWORK_API_PORT:-8787}/api/v1/model/respond \
  -H 'content-type: application/json' \
  -d '{"prompt":"Return one sentence confirming model path is online."}'
```

## Stop (Preserve Data)
```bash
docker compose -f docker-compose.local.yml down
```

## Reset (Delete Data)
```bash
./products/athena-work/tools/workspace_reset.sh
```

## Optional Worker Profile
```bash
docker compose -f docker-compose.local.yml --profile worker up -d
```

## Notes
- DB data persists across restarts in volume `athenawork-db-data`.
- Use reset only when you explicitly want a clean local state.
