# Program Operating System

Control-plane contract for strategic alignment and execution traceability.

## Required Artifacts
- `docs/operator/work-harness/process/PM-TODO.md`
- `docs/operator/work-harness/process/STAGE_EXIT_GATES.md`
- `products/work-harness/delivery-backlog/engineering/active/README.md`
- `products/work-harness/operating-system/observer/`
- any program-state or roadmap artifacts retained for the current operating cycle

## Control-Plane Rules
1. Every new story/bug includes phase + traceability metadata.
2. PM refinement updates program board and active queue in the same cycle.
3. `done` is not `shipped` until explicit release checkpoint approval.
4. Stage-level commits are disallowed; commit once per cycle after observer report generation.

## Observer Rule
- Run `./products/work-harness/tools/run_observer_cycle.sh --cycle-id <cycle-id>` at cycle boundary.
- Observer reports are committed with cycle commit (`cycle-<cycle-id>`).
- Observer records workflow sync checks, state promotions, and release-impact notes.
