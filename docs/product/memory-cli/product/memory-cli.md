# memory-cli Product Guide

## What memory-cli Is

memory-cli is a local-first memory layer for agentic coding workflows. It focuses on durable, governable memory operations rather than runtime execution orchestration.

## Product Scope

In scope:
- Procedural memory (instructions/prompts).
- Semantic retrieval and startup bootstrap context.
- Governance and auditability (review evidence, deterministic behavior).
- Observability (telemetry events and OTel traces).
- A small sqlite-first local memory workflow that is useful in daily work.

Out of scope for v0.1:
- Owning code execution runtime orchestration.
- Pod/container lifecycle management as a product responsibility.
- Tool proxy/wrapper identity. That belongs in `tool-cli`.

## User Personas

- `Engineer`: wants reliable recall and deterministic retrieval behavior.
- `Data Scientist`: wants tunable retrieval/backends and measurable quality gates.
- `Operator`: wants policy controls, audit evidence, and telemetry export.

## Architecture (Practical View)

- CLI surface: `cmd/memory-cli`
- Core modules: `internal/index`, `internal/retrieval`, `internal/governance`, `internal/snapshot`, `internal/telemetry`
- Storage: local memory root (index, metadata, entries, telemetry)

## Default Versus Experimental

Default path:

- `write`
- `retrieve`
- `bootstrap`
- `verify`
- `snapshot`

The Cairn product intentionally exposes only this smaller path.

The broader memory-cli research repo may continue to carry experimental commands and backend work, but those are not part of the default product contract here.

## Quality Model

- Retrieval modes: `classic`, `hybrid`
- Backend: `sqlite`
- Deterministic fallback when semantic confidence is insufficient
- OpenTelemetry is required even in the stripped product

## Governance Model

- Mutation writes require reviewer evidence fields.
- Constraint checks enforce cost/traceability/reliability rules.
- Latency degradation forces deterministic fallback behavior.

## What To Tune

- Retrieval mode and top-k
- Embedding endpoint choice
- Constraint env vars
- OTel/OTLP export configuration

## Start Here

- `docs/operator/memory-cli/getting-started/installation.md`
- `docs/operator/memory-cli/getting-started/quickstart.md`
- `docs/product/memory-cli/cli/commands.md`
- `products/memory-cli/README.md`
