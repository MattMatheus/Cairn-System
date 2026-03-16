# ADR-0002: memory-cli V1 Command Surface

## Status

Accepted for platform bootstrap.

## Date

2026-03-12

## Context

The imported memory-cli module includes a broader historical CLI surface than the unified platform should emphasize initially.

Cairn needs a practical command posture for developer use that stays simple and local-first.

## Decision

The preferred memory-cli v1 command surface is:

- `write`
- `retrieve`
- `bootstrap`
- `verify`
- `snapshot`

These commands define the primary developer workflow for local markdown-backed memory.

Other imported commands may remain available for compatibility, testing, or later evaluation, but they are not the primary v1 posture.

## Consequences

### Positive

- The product is easier to explain and operate.
- The documented path aligns with the sqlite-first local workflow.
- We avoid overcommitting to advanced flows before platform integration is stable.

### Negative

- The codebase still contains a wider command surface than the preferred documented path.
- Future cleanup may still be needed if we want the binary surface itself to be narrower.

