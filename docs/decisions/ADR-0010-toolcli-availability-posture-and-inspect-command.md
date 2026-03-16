# ADR-0010: tool-cli Uses Availability Posture And Inspect To Keep Tool Context Small

Status: Accepted

## Context

tool-cli originally split tools into two trust tiers:

- approved
- local

That trust split is still useful, but it is not enough for the personal Cairn workflow.

The missing distinction is inside the approved tier itself:

- some tools are required and should normally be present
- some tools are ordinary defaults for matching stages
- some tools are approved but should stay scoped until the session actually needs them

Without that second layer, tool-cli would still tend toward bloated default context even if the tools were technically curated.

The product also needs a compact way for a PM, operator, or agent to ask:

- what is this tool for
- should it be in this session
- why is it default or scoped

## Decision

tool-cli will add an availability posture inside the approved registry model.

Supported availability values:

- `required`
- `default`
- `scoped`

Meaning:

- `required`: appropriate to include by default when the stage matches and the personal workflow depends on it
- `default`: approved and normally eligible for stage-matched context injection
- `scoped`: approved and discoverable, but excluded from default context unless explicitly requested or strongly narrowed by query

tool-cli will also add:

- `tool-cli inspect <tool-id>`

The inspect command exists to explain tool fit, context posture, and stage relevance without forcing broad tool startup.

## Consequences

Positive:

- keeps approved tooling curated without making all approved tooling ambient
- gives PM and operator workflows a compact way to understand tool fit
- aligns tool-cli with the rule that tools are called only when needed
- supports simple `SKILL.md`-driven guidance instead of large startup manifests

Tradeoffs:

- registry entries now carry slightly more metadata
- context behavior becomes more opinionated
- operators must understand that approved is not the same as default-injected

## Required Boundaries

tool-cli may:

- include required and default tools in normal context when stage filters match
- keep scoped tools discoverable through `discover`, `list`, and `inspect`
- include scoped tools in `context` only when explicitly requested or query-narrowed

tool-cli may not:

- treat all approved tools as default ambient context
- turn inspection into execution
- expand this change into a general-purpose runtime layer

## Follow-On Work

- extend approved tool definitions for Obsidian, GitNexus, and Firecrawl using the new availability posture
- let work harness stage launch consume the smaller context output as the default path
