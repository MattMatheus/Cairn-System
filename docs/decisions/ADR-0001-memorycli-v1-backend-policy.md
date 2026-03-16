# ADR-0001: memory-cli V1 Backend Policy

## Status

Accepted for platform bootstrap.

## Date

2026-03-12

## Context

Cairn needs a practical memory subsystem for developers using work harness in agent-driven development.

The current memory-cli implementation supports multiple backends and more advanced governance-oriented capabilities. That is too broad for the first unified platform.

The v1 memory product must be:

- easy to run locally
- markdown-first
- CLI-driven
- compatible with future backend expansion

## Decision

memory-cli v1 will use:

- `sqlite` as the default local backend
- markdown files from work harness as the primary source content
- a slim Go CLI as the integration surface
- optional document database support for advanced users

MongoDB is the preferred optional document database path for now because it is familiar to developers and easy to run in Podman.

## Consequences

### Positive

- Default setup remains simple and local-first.
- The CLI contract stays stable even when optional backends are added.
- Advanced users still have a stronger storage option without burdening the default path.

### Negative

- We need to keep backend abstraction disciplined to avoid leaking optional complexity into the common path.
- Mongo support must be explicitly framed as optional, not a required dependency.

## Follow-Up

1. Confirm whether Mongo is used only for document storage or also for retrieval-serving patterns.
2. Decide whether optional backend startup helpers belong under `tools/dev/` or product-specific tooling.
3. Keep the memory-cli product posture sqlite-first even if broader historical commands remain available in the imported module.
4. Decide when Mongo moves from planned support to implemented support.

