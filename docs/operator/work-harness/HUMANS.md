# HUMANS

Operator guide for work harness state harness.

## Summary
Workflow: Planning (as needed) -> Architect (as needed) -> Engineering -> QA -> PM.
Commit policy is cycle-based:
- no intermediate commits
- run observer each cycle
- one commit per cycle: `cycle-<cycle-id>`

## 60-Second Start
1. Ensure branch is `dev`.
2. Run `./products/work-harness/tools/launch_stage.sh engineering`.
3. Execute the top story from active queue.
4. Move to QA and run `./products/work-harness/tools/launch_stage.sh qa`.
5. Run `./products/work-harness/tools/run_observer_cycle.sh --cycle-id <story-id>`.
6. Commit once for the cycle.

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

## Research Council
- Use Research Council for high-uncertainty exploration before implementation.
- Canonical policy: `docs/operator/work-harness/process/RESEARCH_COUNCIL_BASELINE.md`.
- Keep council runs timeboxed and output one council artifact under `workspace/research/`.
- Council spec updates may be agent-executed only under your explicit direction confirmation.

## Control Plane Docs
- `docs/operator/work-harness/process/PROGRAM_OPERATING_SYSTEM.md`
- `docs/operator/work-harness/process/STAGE_EXIT_GATES.md`
- `docs/operator/work-harness/process/BACKLOG_WEIGHTING_POLICY.md`
- `docs/operator/work-harness/process/OPERATOR_DAILY_WORKFLOW.md`
- `docs/operator/work-harness/operations/ALPHA_RELEASE_CHECKPOINT.md`

## Backlog State Model
- `products/work-harness/delivery-backlog/engineering/intake/`
- `products/work-harness/delivery-backlog/engineering/active/`
- `products/work-harness/delivery-backlog/engineering/qa/`
- `products/work-harness/delivery-backlog/engineering/done/`
- `products/work-harness/delivery-backlog/engineering/blocked/`
- `products/work-harness/delivery-backlog/architecture/`
