# Human Runbook

Use this as the daily operator checklist after reading `INTERNAL_BETA.md` and `PLATFORM_QUICKSTART.md`.

## Daily Startup
1. Open `workspace/docs/HUMANS.md`.
2. Open `workspace/docs/README.md`.
3. Confirm repo-local runtime state under `.athena/`.
4. Check active delivery queues.
5. Check claim and artifact queues.
6. Review staging proposals needing approval.
7. Confirm template set in `workspace/templates/` for today’s new items.

## Delivery Operations
1. Confirm top item in `Engineering Active` or `Architecture Active`.
2. Execute stage-specific work.
3. Validate exit gates before moving state.
4. Ensure QA outcome is explicit.
5. Close cycle and record observer/report artifact.

Main supporting commands:

- `./products/athena-work/tools/launch_stage.sh <stage>`
- `./products/athena-work/tools/run_observer_cycle.sh --cycle-id <cycle-id>`
- `./tools/platform/validate_task_metadata.sh`

## Knowledge Operations
1. Promote ideas to candidate claims only when statements are testable.
2. Require explicit test plan before active-test.
3. Require artifact evidence before supported/falsified.
4. Park claims that are not yet testable or evidence-backed.

## Human Review Responsibilities
- Approve or reject staging promotions.
- Resolve safety/quarantine decisions.
- Maintain consistency between operating contracts and actual practice.
- Enforce template-first starts for new tasks/claims/artifacts.
