# Agent Operating Policy

## Scope

This policy governs agent behavior across all of `Athena`.

## Runtime Model

- Human runtime: conceptual work, interpretation, judgment, final decisions.
- Agent runtime: ingestion, triage, metadata maintenance, run artifacts, drafts.

## Zone Permissions

- `Zone A - Agent Write`: `workspace/agents/queue`, `workspace/agents/memory`, `workspace/agents/staging`, `workspace/agents/escape`
  - Agents may create/update/delete freely.
- `Zone B - Agent Suggest`: `workspace/research/*`, `workspace/work/*`
  - Agents may draft and stage updates; human promotes.
- `Zone C - Human Only`: personal/private/reflection zones as tagged by `sensitivity: private_personal`.
  - Agents may read only when explicitly requested; no writes.

## Claim Lifecycle

`idea -> candidate -> active-test -> supported|falsified|parked`

Promotion rule:
- `candidate -> active-test`: requires explicit test plan.
- `active-test -> supported|falsified`: requires at least one linked artifact.

## Ingestion Routing

- Human raw capture: route to `workspace/work/`.
- Agent raw capture: route to `workspace/agents/queue/`.
- Sensitive findings: follow `Content Rejection Policy`.

## Supervision Rule

Until trust is upgraded, all promotions from `Zone A` into canonical graph zones are interactive.
