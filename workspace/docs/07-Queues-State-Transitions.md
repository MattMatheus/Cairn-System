# Queues and State Transitions

## Delivery Queue Paths
- `workspace/agents/delivery/Architecture Intake/`
- `workspace/agents/delivery/Architecture Active/`
- `workspace/agents/delivery/Architecture QA/`
- `workspace/agents/delivery/Architecture Done/`
- `workspace/agents/delivery/Engineering Intake/`
- `workspace/agents/delivery/Engineering Active/`
- `workspace/agents/delivery/Engineering QA/`
- `workspace/agents/delivery/Engineering Done/`

## Allowed Delivery Transitions
- `intake -> active`
- `active -> qa`
- `qa -> done`
- `qa -> active` (with linked defect)

## Disallowed Delivery Transitions
- `intake -> done`
- `active -> done`

## Knowledge Transitions
- `idea -> candidate`
- `candidate -> active-test`
- `active-test -> supported|falsified|parked`

## Transition Discipline
- Every transition requires explicit status update and rationale in note history.
