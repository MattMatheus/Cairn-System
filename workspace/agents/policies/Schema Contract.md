# Schema Contract

## Required Frontmatter (all operational notes)

```yaml
---
id: ath-<uuid-or-stable-id>
type: concept|claim|artifact|decision|run|memory|report|map|task|daily
status: idea|candidate|active-test|supported|falsified|parked|active|blocked|done|archived
domain: work|personal|home|research
updated: YYYY-MM-DD
owner: matt
source_of_truth: human|agent
sensitivity: public|internal|confidential|private_personal
---
```

## Claim Fields (required when `type: claim`)

```yaml
hypothesis: "Testable statement"
confidence: 0.00
falsifiers: []
related_claims: []
evidence_links: []
```

## Artifact Fields (required when `type: artifact`)

```yaml
artifact_kind: run|experiment|dataset|note|decision
supports_claims: []
refutes_claims: []
```

## Agent Trace Fields (required on agent write)

```yaml
agent_last_touch: YYYY-MM-DDTHH:MM:SSZ
agent_intent: ingest|summarize|extract|classify|redact|draft
review_state: pending|approved|rejected
```
