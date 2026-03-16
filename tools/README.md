# Tools

Shared Cairn tooling from the repo root.

## Use These First

- `tools/platform/` for supported platform checks and smoke runs
- `tools/dev/` for local environment helpers

## Promotion Rule

Product-specific scripts should remain with their owning product unless they are promoted into shared platform use.

## Current Boundary

- active runtime validation lives under `tools/platform/`
- local helper scripts live under `tools/dev/`
- Azure/bootstrap artifact retrieval is intentionally deferred until the platform surface and tool contract are finalized
- local bootstrap and release artifact creation now live under `tools/dev/`
