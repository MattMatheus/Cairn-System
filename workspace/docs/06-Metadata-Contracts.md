# Metadata Contracts

## Required Fields (Operational Notes)
```yaml
id: ath-<id>
type: concept|claim|artifact|decision|run|memory|report|map|task|story|bug
status: idea|candidate|active-test|supported|falsified|parked|intake|active|qa|done|blocked|archived
domain: work|research
updated: YYYY-MM-DD
owner: matt
source_of_truth: human|agent
sensitivity: public|internal|confidential|private_personal
```

## Required Fields (Agent Writes)
```yaml
agent_last_touch: YYYY-MM-DDTHH:MM:SSZ
agent_intent: ingest|summarize|extract|classify|redact|draft|implement|verify
review_state: pending|approved|rejected
```

## Required Fields (Delivery Traceability)
```yaml
idea_id: ""
phase: v0.1|v0.2|v0.3
adr_refs: []
success_metric: ""
```

## Required Fields (Claims)
```yaml
hypothesis: ""
confidence: 0.00
falsifiers: []
related_claims: []
evidence_links: []
```

## Required Fields (Tasks)
```yaml
task_type: implementation|investigation|documentation|quality|operations|decision|maintenance
output_type: code_change|design_note|qa_evidence|runbook_update|decision_record|docs_update|report
acceptance_evidence: []
```
