# ADR-0011: Platform North Star Is Personal PKM First And Remote Ministry Utility Later

Status: Accepted

## Context

AthenaPlatform started as a personal system, but the intended destination is broader.

The near-term operator is Matt:

- personal knowledge management
- research capture
- working memory
- staged human-agent collaboration

The longer-term destination is service to a real ministry context with constrained connectivity and uneven access to curated knowledge, including church and ministry partners such as MOHI in Africa.

That downstream context changes the design standard.

The platform should not optimize for feature breadth, ambient cloud dependence, or heavyweight agent runtime patterns. It should optimize for:

- durable local operation
- reviewable artifacts
- explicit provenance
- low operational weight
- graceful use under weak or intermittent connectivity

## Decision

AthenaPlatform's north star is:

- personal PKM and disciplined human-agent work for Matt now
- eventual utility as an offline-friendly, curated knowledge tool for remote ministry contexts later

This is one continuous direction, not two separate products.

The personal system is the proving ground for the same qualities the ministry context will later require:

- trustworthiness
- curation discipline
- portability
- low dependency burden
- understandable workflows

## Consequences

Positive:

- justifies the current preference for small, dependable tools over expansive platforms
- keeps markdown, sqlite, local files, and explicit review gates central
- discourages cloud-first assumptions and ambient tool/runtime bloat
- creates a clear filter for future product decisions

Tradeoffs:

- some advanced integrations will be deferred or rejected
- convenience features that depend on constant connectivity should face a higher bar
- product evolution will favor portability and clarity over maximal automation

## Required Product Implications

AthenaPlatform should prefer:

- local-first storage and retrieval
- explicit staging and review before canonical promotion
- compact tool surfaces that are injected only when needed
- outputs that remain useful outside a live service environment
- provenance and authority metadata that survive export and reuse

AthenaPlatform should avoid by default:

- mandatory hosted-service dependencies
- always-on MCP or agent-runtime assumptions
- operationally heavy subsystems without a proven need
- opaque pipelines that are hard to audit or transplant

## Immediate Guidance

This north star supports and reinforces existing product choices:

- `athena-mind` remains small and dependable
- `athena-use` remains a thin selector and explainer
- `intake-cli` remains markdown-first and review-first
- Firecrawl-style heavy runtime posture remains non-default
- GitNexus remains optional and scoped rather than ambient

## Follow-On Use

Future roadmap, architecture, and tooling decisions should be checked against this question:

"Does this help the personal system now while also moving toward a portable, curated, offline-friendly ministry tool later?"
