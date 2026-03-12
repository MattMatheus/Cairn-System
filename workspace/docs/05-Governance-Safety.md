# Governance and Safety

## Zone Permissions
- Zone A (Agent Write): queue, runs, memory candidates, staging.
- Zone B (Agent Suggest): canonical work/research zones.
- Zone C (Human Only): private-personal or explicitly restricted notes.

## Content Safety Decisions
- ALLOW
- REDACT_AND_ALLOW
- QUARANTINE
- BLOCK

## Non-Negotiable Rules
1. Never store raw secrets in canonical graph notes.
2. Never auto-downgrade BLOCK to ALLOW.
3. Never ingest without a decision when safety signals are present.
4. Never promote staging output without explicit review.

## Canonical Policies
- `workspace/docs/policies/Content Rejection Policy.md`
- `workspace/agents/policies/Agent Operating Policy.md`
- `workspace/agents/policies/Schema Contract.md`
