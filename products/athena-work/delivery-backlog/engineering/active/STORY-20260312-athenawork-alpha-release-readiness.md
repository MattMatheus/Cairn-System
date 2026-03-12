# Story: Prepare AthenaWork alpha release checkpoint

## Metadata
- `id`: STORY-20260312-athenawork-alpha-release-readiness
- `owner_persona`: Product Manager - AthenaWork
- `status`: active
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0004]
- `success_metric`: AthenaWork can produce a truthful alpha release bundle and launch authorization package where the only remaining ship blocker is explicit human approval
- `release_checkpoint`: required

## Problem Statement
- AthenaWork’s current implementation surface is largely runnable, but it cannot be credibly shipped as an alpha because the release evidence path is incomplete. We need a bounded slice that turns current validation coverage into a usable release checkpoint for feedback collection.

## Scope
- In:
  - generate the first real AthenaWork release bundle from the current template
  - collect and record current validation evidence for the alpha checkpoint
  - align launch authorization expectations with the current queue and release-bundle flow
  - document the operator path for preparing an alpha feedback build
- Out:
  - broader cleanup of historical scripts outside the active release path
  - new product features beyond release readiness
  - production deployment automation beyond current local/operator scope

## Acceptance Criteria
1. A real AthenaWork release bundle exists under `products/athena-work/operating-system/handoff/` with explicit `ship|hold` decision, scope, evidence, and rollback notes.
2. Launch authorization generation reports only true remaining blockers after the release bundle is created.
3. The alpha preparation path is documented clearly enough that another operator can reproduce the checkpoint and gather feedback.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `products/athena-work/operating-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`
- `products/athena-work/tools/generate_launch_authorization_package.sh`
- `docs/operator/athena-work/process/STAGE_EXIT_GATES.md`

## Notes
- Prioritize truthful alpha shipment over additional cleanup. Historical validation cleanup can continue after feedback starts flowing.
