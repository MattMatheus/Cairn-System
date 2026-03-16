# HUMANS

Operator guide for work harness inside Cairn.

## System Concept

Keep the system boundary clear:

- The notebook or vault is the human knowledge surface.
- `memory-cli` is the durable retrieval and promotion surface.
- `tool-cli` is the tool mediation surface.
- `work-harness` is the execution surface for stage flow, prompts, scripts, and observer behavior.

The notebook governs meaning and review.

The harness governs execution.

They should connect through explicit process, not by collapsing into one product.

## Summary
Workflow: Planning (as needed) -> Architect (as needed) -> Engineering -> QA -> PM.
Commit policy is cycle-based:
- no intermediate commits
- run observer each cycle
- one commit per cycle: `cycle-<cycle-id>`

## 60-Second Start
1. Read `workspace/docs/HUMANS.md` for the broader workspace/operator surface.
2. Ensure branch matches `CAIRN_REQUIRED_BRANCH` (default `dev`).
3. Run `./products/work-harness/tools/launch_stage.sh engineering`.
4. Execute the top story from the active queue.
5. Move to QA with `./products/work-harness/tools/launch_stage.sh qa`.
6. Run `./products/work-harness/tools/run_observer_cycle.sh --cycle-id <story-id>`.
7. Commit once for the cycle.

## Active References
- `products/work-harness/DEVELOPMENT_CYCLE.md`
- `products/work-harness/delivery-backlog/README.md`
- `products/work-harness/operating-system/observer/README.md`
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
- Planning: `./products/work-harness/tools/launch_stage.sh planning`
- Engineering: `./products/work-harness/tools/launch_stage.sh engineering`
- Architect: `./products/work-harness/tools/launch_stage.sh architect`
- QA: `./products/work-harness/tools/launch_stage.sh qa`
- PM: `./products/work-harness/tools/launch_stage.sh pm`
- Cycle: `./products/work-harness/tools/launch_stage.sh cycle`
- Observer: `./products/work-harness/tools/run_observer_cycle.sh --cycle-id <cycle-id>`
- Docs/site build: `./products/work-harness/tools/build_docs_site.sh`
- GitNexus preflight when codegraph work is needed: `./products/work-harness/tools/check_gitnexus_readiness.sh`

## Control Plane Docs
- `products/work-harness/operating-system-vault/PROGRAM_CONTROL_PLANE.md`
- `products/work-harness/operating-system-vault/STAGE_EXIT_GATES.md`
- `products/work-harness/operating-system-vault/DELIVERY_STATE_MODEL.md`
- `products/work-harness/operating-system-vault/UNIFIED_METADATA_CONTRACT.md`
- `docs/operator/work-harness/operations/ALPHA_RELEASE_CHECKPOINT.md`

## Backlog State Model
- `products/work-harness/delivery-backlog/engineering/intake/`
- `products/work-harness/delivery-backlog/engineering/active/`
- `products/work-harness/delivery-backlog/engineering/qa/`
- `products/work-harness/delivery-backlog/engineering/done/`
- `products/work-harness/delivery-backlog/engineering/blocked/`
- `products/work-harness/delivery-backlog/architecture/`

## Notes
- Historical and migration material exists, but this file is part of the active operator path.
- If a referenced path no longer exists, prefer the canonical `products/`, `workspace/`, and `.cairn/` surfaces rather than reconstructing old repo-era layouts.
