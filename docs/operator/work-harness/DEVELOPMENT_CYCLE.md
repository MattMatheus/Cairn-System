# Development Cycle System

## Human Operator Entry Point
- `HUMANS.md` is the canonical operator navigation page.

## Stage Launchers
- Planning: `products/work-harness/stage-prompts/active/planning-seed-prompt.md`
- Engineering: `products/work-harness/stage-prompts/active/next-agent-seed-prompt.md`
- Architect: `products/work-harness/stage-prompts/active/architect-agent-seed-prompt.md`
- QA: `products/work-harness/stage-prompts/active/qa-agent-seed-prompt.md`
- PM: `products/work-harness/stage-prompts/active/pm-refinement-seed-prompt.md`
- Cycle: `products/work-harness/stage-prompts/active/cycle-seed-prompt.md`

Quick commands:
- `products/work-harness/tools/launch_stage.sh <planning|engineering|architect|qa|pm|cycle>`
- `products/work-harness/tools/run_observer_cycle.sh --cycle-id <cycle-id>`

Research Council policy:
- `docs/operator/work-harness/process/RESEARCH_COUNCIL_BASELINE.md`

## Branch Safety Rule
- Stage launches require branch `CAIRN_REQUIRED_BRANCH` (default `dev`).
- Mismatch aborts launch.

## Main Promotion Hygiene
- `dev -> main` is PR-only with human approval.
- While PR is open, freeze `dev`.
- After merge: fetch, ff-merge `origin/main` into `dev`, push `dev`, verify equal SHAs.

## Canonical Flow
1. Planning may run a timeboxed Research Council when uncertainty is high.
2. PM keeps ranked stories in `products/work-harness/delivery-backlog/engineering/active/`.
3. Engineering executes top story and prepares handoff.
4. QA validates and returns to `active` or promotes to `done`.
5. Observer runs and writes cycle report.
6. Single cycle commit is created.
7. PM refines intake and updates control plane.

## Commit Convention
- Exactly one commit per completed cycle.
- Format: `cycle-<cycle-id>`.
- Include observer report and cycle artifacts in that commit.

## Work-System Doc Sync Rule
If workflow behavior changes, update:
- `HUMANS.md`
- `products/work-harness/AGENTS.md`
- relevant `products/work-harness/stage-prompts/active/*`
- `docs/operator/work-harness/process/RESEARCH_COUNCIL_BASELINE.md` for council-policy changes

## Empty Backlog Rule
If engineering active is empty, output `no stories` and run PM refinement.
