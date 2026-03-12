# Founder-Operator Daily Workflow (Codex App)

Single-operator daily script for AthenaWork state harness cycles.

## Startup Routine
1. Confirm branch discipline (`dev`).
2. Review `docs/operator/athena-work/process/CYCLE_INDEX.md` and `docs/operator/athena-work/DEVELOPMENT_CYCLE.md`.
3. Run `./products/athena-work/tools/run_doc_tests.sh`.
4. Confirm active queue in `products/athena-work/delivery-backlog/engineering/active/README.md`.
5. Review `docs/operator/athena-work/process/PROGRAM_OPERATING_SYSTEM.md`.

## Engineering + QA Loop
1. `./products/athena-work/tools/launch_stage.sh engineering`
2. Execute top active story and move it to QA with handoff.
3. `./products/athena-work/tools/launch_stage.sh qa`
4. QA returns to active with defects or promotes to done.
5. `./products/athena-work/tools/run_observer_cycle.sh --cycle-id <story-id>`
6. Commit once: `cycle-<cycle-id>`

## PM Loop
1. `./products/athena-work/tools/launch_stage.sh pm`
2. Refine intake and rank active queue.
3. Update program-state and queue references described in `docs/operator/athena-work/process/PROGRAM_OPERATING_SYSTEM.md`.
4. Run observer and commit cycle.

## Rules
- If engineering returns `no stories`, run PM refinement.
- Do not advance state on failing tests.
- Do not bypass stage gates.
