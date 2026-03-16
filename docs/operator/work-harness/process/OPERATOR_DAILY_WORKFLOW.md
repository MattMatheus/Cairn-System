# Founder-Operator Daily Workflow (Codex App)

Single-operator daily script for work harness state harness cycles.

## Startup Routine
1. Confirm branch discipline (`dev`).
2. Review `docs/operator/work-harness/process/CYCLE_INDEX.md` and `docs/operator/work-harness/DEVELOPMENT_CYCLE.md`.
3. Run `./products/work-harness/tools/check_go_toolchain.sh`.
4. If codegraph work is likely, run `./products/work-harness/tools/check_gitnexus_readiness.sh`.
5. Run `./products/work-harness/tools/run_doc_tests.sh`.
6. Confirm active queue in `products/work-harness/delivery-backlog/engineering/active/README.md`.
7. Review `docs/operator/work-harness/process/PROGRAM_OPERATING_SYSTEM.md`.

## Engineering + QA Loop
1. `./products/work-harness/tools/launch_stage.sh engineering`
2. Execute top active story and move it to QA with handoff.
3. `./products/work-harness/tools/launch_stage.sh qa`
4. QA returns to active with defects or promotes to done.
5. `./products/work-harness/tools/run_observer_cycle.sh --cycle-id <story-id>`
6. Commit once: `cycle-<cycle-id>`

## PM Loop
1. `./products/work-harness/tools/launch_stage.sh pm`
2. Refine intake and rank active queue.
3. Update program-state and queue references described in `docs/operator/work-harness/process/PROGRAM_OPERATING_SYSTEM.md`.
4. Run observer and commit cycle.

## Rules
- If engineering returns `no stories`, run PM refinement.
- Do not advance state on failing tests.
- Do not bypass stage gates.
