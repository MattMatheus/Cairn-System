# ADR-0003: AthenaMind Optional MongoDB Contract

## Status

Accepted

## Context

AthenaMind v1 is intentionally `sqlite`-first. That default path is already buildable, testable, and easy for developers to run locally.

The platform also needs an optional stronger document-store path for developers who want:

- larger local datasets
- richer document inspection and ad hoc querying
- a backend that can be containerized cleanly with Podman

MongoDB is the preferred document DB for that optional path.

## Decision

AthenaPlatform will standardize an optional local MongoDB contract now, while keeping `sqlite` as the default active backend.

The MongoDB contract for local development is:

- container runtime: Podman
- default container name: `mongodb-local`
- default URI: `mongodb://127.0.0.1:27017`
- default database: `athenamind`
- default collections:
  - `memory_entries`
  - `memory_embeddings`
  - `memory_audits`

Recommended environment variables:

- `ATHENA_MONGODB_URI`
- `ATHENA_MONGODB_DATABASE`

## Consequences

Positive:

- developers get a clear optional backend target without complicating the default path
- local Podman usage remains straightforward
- future Mongo integration work can converge on one contract
- AthenaMind can already use MongoDB for optional index and embedding persistence

Negative:

- docs and tooling must stay explicit that Mongo is optional and not the default backend
- active retrieval behavior still remains optimized for the sqlite-first path

## Follow-Up

1. Decide whether Mongo should become a documented smoke-check path under `tools/platform/`.
2. Decide whether Mongo is used only for index/embedding persistence or also for additional retrieval-side experiments.
3. Add richer migration and inspection tooling for Mongo-backed memory roots.
