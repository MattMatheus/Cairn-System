# Unified Metadata Contract

## Required (All Operational Notes)
```yaml
---
id: ath-<id>
type: concept|claim|artifact|decision|run|memory|report|map|task|story|bug|escape
status: idea|candidate|active-test|supported|falsified|parked|intake|active|qa|done|blocked|archived
domain: work|research
updated: YYYY-MM-DD
owner: matt
source_of_truth: human|agent
sensitivity: public|internal|confidential|private_personal
---
```

## Required (Agent Writes)
```yaml
agent_last_touch: YYYY-MM-DDTHH:MM:SSZ
agent_intent: ingest|summarize|extract|classify|redact|draft|implement|verify|escalate
review_state: pending|approved|rejected
```

## Required (Delivery Traceability)
```yaml
idea_id: ""
phase: v0.1|v0.2|v0.3
adr_refs: []
success_metric: ""
```

## Required (Claim Notes)
```yaml
hypothesis: ""
confidence: 0.00
falsifiers: []
related_claims: []
evidence_links: []
```

## Required (Task Notes)
```yaml
task_type: implementation|investigation|documentation|quality|operations|decision|maintenance
output_type: code_change|design_note|qa_evidence|runbook_update|decision_record|docs_update|report
acceptance_evidence: []
```

## Required (Escape Records)
```yaml
escape_class: blocked|ambiguous_requirements|policy_conflict|missing_context|low_confidence|loop_detected|scope_exceeded
escape_summary: ""
attempt_count: 1
source_item_id: ath-
source_item_path: ""
state_at_escape: intake|active|qa
resolution_action: retry|redirect|cancel
resolved_by: ""
resolved_at: ""
resolution_notes: ""
```
