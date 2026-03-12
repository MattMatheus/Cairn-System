# Development Cycle System

## Human Operator Entry Point
- `products/athena-work/HUMANS.md` is the canonical product operator navigation page.

## Stage Launchers
- Planning: `products/athena-work/stage-prompts/active/planning-seed-prompt.md`
- Engineering: `products/athena-work/stage-prompts/active/next-agent-seed-prompt.md`
- Architect: `products/athena-work/stage-prompts/active/architect-agent-seed-prompt.md`
- QA: `products/athena-work/stage-prompts/active/qa-agent-seed-prompt.md`
- PM: `products/athena-work/stage-prompts/active/pm-refinement-seed-prompt.md`
- Cycle: `products/athena-work/stage-prompts/active/cycle-seed-prompt.md`

Quick commands:
- `products/athena-work/tools/launch_stage.sh <planning|engineering|architect|qa|pm|cycle>`
- `products/athena-work/tools/run_observer_cycle.sh --cycle-id <cycle-id>`

Research Council policy:
- handled as a bounded planning/architecture exploration pattern; do not treat historical repo-era council docs as required active dependencies

## Branch Safety Rule
- Stage launches require branch `ATHENA_REQUIRED_BRANCH` (default `dev`).
- Mismatch aborts launch.

## Main Promotion Hygiene
- `dev -> main` is PR-only with human approval.
- While PR is open, freeze `dev`.
- After merge: fetch, ff-merge `origin/main` into `dev`, push `dev`, verify equal SHAs.

## Canonical Flow
1. Planning may run a timeboxed Research Council when uncertainty is high.
2. PM keeps ranked stories in `products/athena-work/delivery-backlog/engineering/active/`.
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
- `products/athena-work/HUMANS.md`
- `products/athena-work/AGENTS.md`
- relevant `products/athena-work/stage-prompts/active/*`
- relevant `workspace/docs/*` operator guidance
- `products/athena-work/operating-system-vault/*` only when the underlying contract truly changes

## Empty Backlog Rule
If engineering active is empty, output `no stories` and run PM refinement.
