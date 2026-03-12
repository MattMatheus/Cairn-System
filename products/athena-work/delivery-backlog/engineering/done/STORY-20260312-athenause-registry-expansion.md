# Story: Expand AthenaUse approved registry with additional platform tools

## Metadata
- `id`: STORY-20260312-athenause-registry-expansion
- `owner_persona`: Product Manager - AthenaWork
- `status`: done
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0005, ADR-0006]
- `success_metric`: approved registry grows from 2 to at least 5 validated tools with stage affinity and schema coverage
- `release_checkpoint`: required

## Problem Statement
- AthenaUse is integrated into stage launch, but the approved registry is still too narrow to be a credible shared tool-context surface for the platform.

## Scope
- In:
  - Add several real, repo-supported platform tools to `products/athena-use/registry/approved-tools.yaml`
  - Preserve the existing narrow registry contract and dependency policy
  - Add or update tests that cover the expanded registry shape where needed
- Out:
  - Executing tools through AthenaUse
  - Memory-backed registry mode
  - Azure/bootstrap integrations

## Acceptance Criteria
1. The approved AthenaUse registry contains at least three additional real platform tools beyond the current two entries.
2. Each new tool entry has a clear description, valid stage affinity, and schema fields where parameters are required.
3. `go test ./...` in `products/athena-use` and `go run ./cmd/use-cli validate` both pass after the registry expansion.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- Existing AthenaUse registry contract in `products/athena-use/registry/approved-tools.yaml`
- ADR-0005 trust and registry policy
- ADR-0006 telemetry and dependency policy

## Notes
- Product-first ranking: this is the smallest step that improves tool-context usefulness without reopening architecture scope.
