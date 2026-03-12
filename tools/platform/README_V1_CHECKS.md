# V1 Checks

## Purpose

This folder contains root-level checks for the current AthenaPlatform v1 posture.

## Current Checks

- `validate_task_metadata.sh`: validates workspace task metadata
- `smoke_v1.sh`: exercises the approved sqlite-first AthenaMind path plus basic AthenaWork launch and observer flow
- `smoke_mongodb.sh`: exercises the optional Mongo-backed AthenaMind persistence path

## Notes

- These checks are intentionally small and trustworthy.
- They complement, rather than replace, imported AthenaWork historical validation scripts.
- `smoke_v1.sh` writes embeddings when `ATHENA_EMBEDDING_ENDPOINT` is set.
- `smoke_v1.sh` disables AthenaMind latency fallback during semantic health validation unless you explicitly provide `MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS`.
- `smoke_mongodb.sh` keeps Mongo optional and validates only AthenaMind's Mongo-backed path.
- Both smoke scripts default to repo-local runtime paths under `.athena/` and can be overridden with `ATHENA_HOME`, `ATHENA_MEMORY_ROOT`, and related env vars.
