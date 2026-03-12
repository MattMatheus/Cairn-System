# Stage Exit Gates (Merged)

## Planning Gate
- Session captured.
- Intake items created with required metadata.
- Next stage recommendation is explicit.

## Architect Gate
- Decision scope is explicit.
- Architecture outputs and downstream implementation mapping are explicit.
- Story transitions to architecture QA.

## Engineering Gate
- Acceptance criteria implemented.
- Tests updated and passing.
- Handoff package produced.
- Story transitions to engineering QA.

## QA Gate
- Acceptance evidence explicit.
- Regression evaluation explicit.
- Result artifact explicit.
- Transition is explicit: `qa -> done` or `qa -> active` with defects.

## PM Gate
- Intake validation complete.
- Active queue ranked.
- Control plane board synced.

## Knowledge Promotion Gate
- Claim state transition is explicit.
- For `active-test -> supported|falsified`, at least one artifact link is present.
- Safety decision is attached when ingesting sensitive sources.

## Cycle Closure Gate
- Observer/report artifact exists.
- Workflow sync checklist complete.
- Exactly one cycle commit.
