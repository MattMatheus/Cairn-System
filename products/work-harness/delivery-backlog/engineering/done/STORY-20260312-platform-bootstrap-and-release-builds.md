# Story: Add platform bootstrap and Azure release build path

## Metadata
- `id`: STORY-20260312-platform-bootstrap-and-release-builds
- `owner_persona`: Product Manager - work harness
- `status`: done
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0004]
- `success_metric`: a new operator can bootstrap the repo locally and Azure DevOps can produce deterministic macOS Apple silicon and Windows CLI artifacts from the unified repo
- `release_checkpoint`: deferred

## Problem Statement
- Cairn had no canonical bootstrap entrypoint and no supported CI path for building distributable binaries from the unified repo. That made setup slower and alpha distribution harder than necessary.

## Scope
- In:
  - add a shared local bootstrap script
  - add a shared release artifact build script for supported targets
  - add Azure DevOps pipeline configuration for tests and artifact publishing
  - update docs and validation to use module-scoped Go commands instead of nonexistent root-module assumptions
- Out:
  - artifact download/install automation
  - packaging beyond zipped binaries and checksums
  - broader cleanup of historical workflow references outside the touched setup/build path

## Acceptance Criteria
1. Cairn has a canonical bootstrap command that prepares repo-local runtime folders and validates the Go toolchain.
2. A canonical build command produces deterministic `memory-cli` and `tool-cli` artifacts for macOS Apple silicon and Windows PCs.
3. Azure DevOps pipeline configuration exists and publishes release artifacts from the unified repo.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `products/memory-cli/go.mod`
- `products/tool-cli/go.mod`
- `products/work-harness/tools/check_go_toolchain.sh`

## Notes
- Keep the distribution path narrow for now: source-first development, zipped binary artifacts for alpha convenience.
