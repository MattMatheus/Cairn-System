# Content Rejection Policy (Agent-Compatible)

## Purpose

Prevent sensitive or low-safety content from entering the knowledge graph during capture, ingestion, and agent write-back.

Scope: all of `Athena`.

## Policy Outcomes

- `ALLOW`
- `REDACT_AND_ALLOW`
- `QUARANTINE`
- `BLOCK`

## Decision Contract (required agent output)

```yaml
policy_version: "1.0"
decision: "ALLOW|REDACT_AND_ALLOW|QUARANTINE|BLOCK"
severity: "low|medium|high|critical"
categories:
  - "credential"
  - "personal_identifier"
  - "financial"
  - "health"
  - "private_personal"
signals:
  deterministic_hits: 0
  statistical_hits: 0
  classifier_label: "safe|sensitive|unknown"
  classifier_confidence: 0.0
  semantic_confidence: 0.0
spans:
  - start: 0
    end: 0
    category: "credential"
    confidence: 0.0
actions:
  redact_spans: true
  quarantine_path: "Athena/30 AI/Quarantine"
  allow_graph_ingest: false
  require_human_review: true
reason: "Short rationale"
timestamp_utc: "YYYY-MM-DDTHH:MM:SSZ"
agent_id: "agent-name"
```

## Detection Stack

1. Deterministic detectors
- Pattern + checksum/format validators for SSN, API keys, tokens, private keys, connection strings, emails, phone numbers, and account IDs.

2. Statistical detectors
- Entropy and token-shape scoring for secret-like strings and dumps.

3. Context classifier
- Label context as `safe`, `sensitive`, or `unknown` to reduce false positives.

4. Semantic review
- LLM judge returns strict JSON only (no free text decisions).

## Action Rules

- `BLOCK`:
  - Confirmed high-risk personal identifiers with strong context (example: SSN + name + DOB).
  - Raw private keys or active bearer/refresh tokens.
- `QUARANTINE`:
  - High-confidence secrets or identifiers where intent is unclear.
  - Any ambiguous credential blob above confidence threshold.
- `REDACT_AND_ALLOW`:
  - Sensitive spans detected with moderate confidence and non-critical context.
  - Replace spans with stable placeholders: `[REDACTED:<category>:<hash8>]`.
- `ALLOW`:
  - No material sensitive findings.

## Thresholds (initial defaults)

```yaml
thresholds:
  block:
    min_confidence: 0.92
  quarantine:
    min_confidence: 0.80
  redact_and_allow:
    min_confidence: 0.60
  allow:
    max_confidence: 0.59
```

## Non-Negotiable Guardrails

- Never store raw secrets in graph notes.
- Never auto-downgrade a `BLOCK` to `ALLOW`.
- Never silently redact without audit metadata.
- Never permit ingestion if `decision` is missing.

## Override Protocol

- Overrides are interactive only.
- Allowed overrides:
  - `BLOCK -> QUARANTINE`
  - `QUARANTINE -> REDACT_AND_ALLOW`
- Disallowed override:
  - `BLOCK -> ALLOW`

Override record:

```yaml
override:
  approved_by: "matt"
  previous_decision: "BLOCK"
  new_decision: "QUARANTINE"
  reason: "User-confirmed research artifact handling"
  timestamp_utc: "YYYY-MM-DDTHH:MM:SSZ"
```

## Audit Record (required for non-ALLOW)

```yaml
audit:
  note_path: "absolute/or/vault/path.md"
  note_hash_sha256: "..."
  policy_version: "1.0"
  decision: "..."
  reviewer: "agent|human"
  queued_for_review: true
```

## Routing

- `ALLOW`: proceed to normal ingest pipeline.
- `REDACT_AND_ALLOW`: write redacted content, attach `audit`.
- `QUARANTINE`: move to `Athena/30 AI/90 Quarantine/<YYYY-MM-DD>/`.
- `BLOCK`: reject write, emit alert item to agent queue.

## Claim-Graph Compatibility

- Redaction must preserve markdown structure and links.
- Claims may ingest redacted artifacts if provenance is intact.
- Redacted spans must be stable and referenceable by placeholder ID.
