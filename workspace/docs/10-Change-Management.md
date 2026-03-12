# Change Management

## When Workflow Behavior Changes
Update all of the following together:
- `workspace/AGENTS.md`
- `workspace/docs/HUMANS.md`
- `products/athena-work/operating-system-vault/*` affected contracts
- `workspace/docs/*` affected explanatory docs

## Required Change Record
Every workflow change must include:
1. Intent
2. Scope
3. Contract updates
4. Migration impact
5. Validation checks

## Safe Change Sequence
1. Update contracts.
2. Update docs.
3. Validate links and required fields.
4. Announce control-plane change summary.
