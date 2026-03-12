# MongoDB Local Note

## Purpose

This document defines the local MongoDB contract for AthenaPlatform development.

MongoDB is optional. The default AthenaMind path remains `sqlite`.

## Standard Local Contract

- container name: `mongodb-local`
- host port: `27017`
- URI: `mongodb://127.0.0.1:27017`
- database: `athenamind`

Recommended environment:

```bash
export ATHENA_HOME="$PWD/.athena"
export ATHENA_INDEX_BACKEND="mongodb"
export ATHENA_MONGODB_URI="mongodb://127.0.0.1:27017"
export ATHENA_MONGODB_DATABASE="athenamind"
export ATHENA_MEMORY_ROOT="$ATHENA_HOME/memory/default"
mkdir -p "$ATHENA_MEMORY_ROOT"
```

If your local container enforces auth, include credentials in the URI:

```bash
export ATHENA_MONGODB_URI="mongodb://admin:changeme@127.0.0.1:27017/?authSource=admin"
```

## Intended Role

Current intended collections:

- `memory_entries`
- `memory_embeddings`
- `memory_audits`

This contract now supports an optional AthenaMind adapter path for index and embedding persistence. The default v1 runtime posture still remains `sqlite` first.

## Local Validation

If the standard container is already running in Podman:

```bash
./tools/dev/check_mongodb_local.sh
```

Example AthenaMind write using Mongo-backed persistence:

```bash
export ATHENA_HOME="$PWD/.athena"
export ATHENA_MEMORY_ROOT="${ATHENA_MEMORY_ROOT:-$ATHENA_HOME/memory/default}"
mkdir -p "$ATHENA_MEMORY_ROOT"

(cd products/athena-mind && \
  ATHENA_INDEX_BACKEND="mongodb" \
  ATHENA_MONGODB_URI="mongodb://admin:changeme@127.0.0.1:27017/?authSource=admin" \
  ATHENA_MONGODB_DATABASE="athenamind" \
  go run ./cmd/memory-cli write \
    --root "$ATHENA_MEMORY_ROOT" \
    --id mongodb-example \
    --title "MongoDB Example" \
    --type prompt \
    --domain platform \
    --body "Persist index records in MongoDB while keeping markdown files local." \
    --stage planning \
    --reviewer operator \
    --decision approved \
    --reason "mongodb example" \
    --risk "low" \
    --notes "example")
```

## Notes

- Use Mongo when you want a stronger local document-store workflow.
- Do not make Mongo a requirement for the common developer path.
- AthenaWork should continue to integrate with AthenaMind through the CLI contract, not direct DB assumptions.
