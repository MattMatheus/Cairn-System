# HUMANS

Operator guide for AthenaWork inside AthenaPlatform.

## Summary
Workflow: Planning (as needed) -> Architect (as needed) -> Engineering -> QA -> PM.
Commit policy is cycle-based:
- no intermediate commits
- run observer each cycle
- one commit per cycle: `cycle-<cycle-id>`

## 60-Second Start
1. Read `workspace/docs/HUMANS.md` for the broader workspace/operator surface.
2. Ensure branch matches `ATHENA_REQUIRED_BRANCH` (default `dev`).
3. Run `./products/athena-work/tools/launch_stage.sh engineering`.
4. Execute the top story from the active queue.
5. Move to QA with `./products/athena-work/tools/launch_stage.sh qa`.
6. Run `./products/athena-work/tools/run_observer_cycle.sh --cycle-id <story-id>`.
7. Commit once for the cycle.

## Active References
- `products/athena-work/DEVELOPMENT_CYCLE.md`
- `products/athena-work/delivery-backlog/README.md`
- `products/athena-work/operating-system/observer/README.md`
- `workspace/docs/HUMANS.md`
- `workspace/docs/08-Human-Runbook.md`

## Main Promotion Policy
- Use PRs only for `dev -> main`.
- While PR is open, freeze `dev`.
- After merge, sync `dev` to `main`:
  1. `git fetch origin main dev`
  2. `git checkout dev`
  3. `git merge --ff-only origin/main`
  4. `git push origin dev`
  5. `git ls-remote --heads origin main dev` (SHAs must match)

## Stage Commands
- Planning: `./products/athena-work/tools/launch_stage.sh planning`
- Engineering: `./products/athena-work/tools/launch_stage.sh engineering`
- Architect: `./products/athena-work/tools/launch_stage.sh architect`
- QA: `./products/athena-work/tools/launch_stage.sh qa`
- PM: `./products/athena-work/tools/launch_stage.sh pm`
- Cycle: `./products/athena-work/tools/launch_stage.sh cycle`
- Observer: `./products/athena-work/tools/run_observer_cycle.sh --cycle-id <cycle-id>`
- Docs/site build: `./products/athena-work/tools/build_docs_site.sh`

## Control Plane Docs
- `products/athena-work/operating-system-vault/PROGRAM_CONTROL_PLANE.md`
- `products/athena-work/operating-system-vault/STAGE_EXIT_GATES.md`
- `products/athena-work/operating-system-vault/DELIVERY_STATE_MODEL.md`
- `products/athena-work/operating-system-vault/UNIFIED_METADATA_CONTRACT.md`

## Backlog State Model
- `products/athena-work/delivery-backlog/engineering/intake/`
- `products/athena-work/delivery-backlog/engineering/active/`
- `products/athena-work/delivery-backlog/engineering/qa/`
- `products/athena-work/delivery-backlog/engineering/done/`
- `products/athena-work/delivery-backlog/engineering/blocked/`
- `products/athena-work/delivery-backlog/architecture/`

## Notes
- AthenaWork historical/migration material exists, but this file is part of the active operator path.
- If a referenced path no longer exists, prefer the canonical `products/`, `workspace/`, and `.athena/` surfaces rather than reconstructing old repo-era layouts.
