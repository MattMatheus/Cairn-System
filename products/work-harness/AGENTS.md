# AGENTS

Navigation and operating guide for work harness state-harness agents.

## Mission Context
work harness is a state-harness product with a strict staged delivery workflow.

## Repository Context
Primary product boundary: `products/work-harness/`.
Canonical workspace boundary: `workspace/`.
Local runtime state belongs under `.cairn/`.

## First 5 Minutes
1. Read `products/work-harness/HUMANS.md`.
2. Read `products/work-harness/DEVELOPMENT_CYCLE.md`.
3. Read `workspace/AGENTS.md`.
4. Read `products/work-harness/delivery-backlog/engineering/active/README.md`.
5. Launch requested stage with `products/work-harness/tools/launch_stage.sh <stage>`.

## Canonical Stage Prompts
- Planning: `products/work-harness/stage-prompts/active/planning-seed-prompt.md`
- Engineering: `products/work-harness/stage-prompts/active/next-agent-seed-prompt.md`
- Architect: `products/work-harness/stage-prompts/active/architect-agent-seed-prompt.md`
- QA: `products/work-harness/stage-prompts/active/qa-agent-seed-prompt.md`
- PM: `products/work-harness/stage-prompts/active/pm-refinement-seed-prompt.md`
- Cycle: `products/work-harness/stage-prompts/active/cycle-seed-prompt.md`

## Mandatory Behavioral Rules
- Branch must match `CAIRN_REQUIRED_BRANCH` (default `dev`).
- Respect backlog state model and stage gates.
- Do not fabricate work when engineering reports `no stories`.
- Do not commit during intermediate stage transitions.
- Run observer after each completed cycle: `products/work-harness/tools/run_observer_cycle.sh --cycle-id <cycle-id>`.
- Commit once per cycle with format `cycle-<cycle-id>`.
- Keep the active backlog readmes and queue ordering in sync during PM refinement.
- `dev -> main` requires human-approved PR.
- While a `dev -> main` PR is open, freeze `dev`.
- After merge, fast-forward `dev` to `main` and verify SHAs match.

## Branch Sync Procedure
1. `git fetch origin main dev`
2. `git checkout dev`
3. `git merge --ff-only origin/main`
4. `git push origin dev`
5. Verify `git ls-remote --heads origin main dev` shows equal SHAs.

## Documentation Sync Rule
When workflow behavior changes, update:
1. `HUMANS.md`
2. `DEVELOPMENT_CYCLE.md`
3. affected stage prompt(s)
4. affected `workspace/docs/*` guidance
5. related `products/work-harness/operating-system-vault/*` contracts when behavior actually changes there

## Notes
- Prefer current platform paths over any historical repo-era references that still appear in imported docs.
- The active operator surface is `products/work-harness/`, `workspace/`, and `.cairn/`.
