# V1 Checks

## Purpose

This folder contains root-level checks for the current Cairn v1 posture.

## Current Checks

- `validate_task_metadata.sh`: validates workspace task metadata
- `smoke_v1.sh`: exercises the approved sqlite-first memory-cli path plus basic work harness launch and observer flow

## Notes

- These checks are intentionally small and trustworthy.
- They complement, rather than replace, imported work harness historical validation scripts.
- The smoke path validates the approved tool-cli registry before running downstream platform checks.
- `smoke_v1.sh` writes embeddings when `CAIRN_EMBEDDING_ENDPOINT` is set.
- `smoke_v1.sh` disables memory-cli latency fallback during semantic health validation unless you explicitly provide `MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS`.
- `smoke_v1.sh` defaults to repo-local runtime paths under `.cairn/` and can be overridden with `CAIRN_HOME`, `CAIRN_MEMORY_ROOT`, and related env vars.
