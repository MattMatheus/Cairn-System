---
title: Development Cognition Graph
description: AthenaWork extension that organizes durable development knowledge as a traversable graph with progressive disclosure and canonical source references.
doc_type: architecture
status: draft
---

# Development Cognition Graph

## Summary

AthenaWork should treat development knowledge as a graph, not a flat pile of prompts and docs. The graph gives agents a bounded way to move from repo orientation to action, while [[CANONICAL_SOURCE_REGISTRY]] records pin external tool semantics to official vendor docs.

## Intended Audience

- humans shaping AthenaWork's documentation model
- stage-agent authors
- maintainers of process, architecture, and operational docs

## Preconditions

- AthenaWork remains primarily a development-work system
- durable knowledge stays separate from active backlog state
- new graph nodes follow progressive disclosure and concise summaries

## Main Flow

AthenaWork should split development context into three layers:

1. Durable internal knowledge
2. Active work state
3. External canonical source docs

Durable internal knowledge belongs in a graph of short, linkable markdown nodes. Active work state remains in backlog files, stage prompts, and cycle artifacts. External tool semantics live in [[CANONICAL_SOURCE_REGISTRY]] records that point to official docs and version assumptions.

Agents should traverse the graph in this order:

1. Read an index or MOC summary
2. Read node descriptions before full sections
3. Follow only the most relevant links
4. Pull canonical vendor docs only when the task depends on external tool behavior
5. Stop when the agent has enough context to act safely

## How It Works

The graph should optimize for development cognition, not general knowledge storage. The useful node families are:

- `index`: top-level orientation and attention routing
- `moc`: curated subdomain map such as process, architecture, workflow, or product area
- `claim`: durable engineering claim or invariant
- `procedure`: runbook or action flow
- `heuristic`: decision shortcut for routine choices
- `decision`: design rationale and consequences
- `example`: concrete worked case or code-path illustration
- `canonical_source`: external vendor truth record

The graph should also be layered:

- entry layer: `HUMANS.md`, `AGENTS.md`, top indexes, stage entry docs
- routing layer: MOCs for process, architecture, workflow, codebase, product
- reasoning layer: claims, heuristics, invariants, decisions
- action layer: procedures, checklists, runbooks, verification recipes
- evidence layer: examples, postmortems, code references, experiments

For AthenaWork, the first graph domains should be:

- `process`: stage contracts, exit gates, handoffs, branch rules
- `architecture`: system boundaries, invariants, module responsibilities
- `workflow`: debugging, verification, local environment setup, release flow
- `codebase`: edit maps, entry points, common dependency paths
- `product`: user-facing goals and feature boundaries
- `references`: canonical source records for external tools

## Node Contract

Every graph node should support progressive disclosure:

1. short summary
2. main flow or main claim
3. details and edge cases
4. references

Every node should also declare enough frontmatter for scan-first traversal:

```yaml
title: Verification Patterns
description: How AthenaWork agents decide what to run to prove a code change is correct.
doc_type: heuristic
status: active
domain: workflow
canonical_sources:
  - pytest
  - podman
related:
  - [[COMMON_EDIT_PATHS]]
  - [[STAGE_EXIT_GATES]]
```

## Traversal Policy

Agents should not read the whole graph by default.

Preferred traversal policy:

1. start at the most local entry point for the task
2. read frontmatter and summary only
3. follow one to three high-signal links
4. escalate to deeper nodes only when ambiguity remains
5. consult canonical vendor docs only for external semantics, APIs, syntax, flags, version behavior, or safety constraints

This keeps AthenaWork aligned with bounded attention instead of "load everything."

## Failure Modes

- oversized nodes become mini-handbooks and defeat traversal
- weak descriptions force full-file reads
- untyped links collapse routing into noisy adjacency
- active project state leaks into durable graph nodes and creates staleness
- internal guidance contradicts vendor docs without surfacing the conflict

## References

- [[CANONICAL_SOURCE_REGISTRY]]
- [[AUTHORING_GRAPH_NODES]]
- [AGENTS.md](/home/matt/Workspace/repos/trusted/Cairn/products/athena-work/AGENTS.md)
- [HUMANS.md](/home/matt/Workspace/repos/trusted/Cairn/products/athena-work/HUMANS.md)
- [DEVELOPMENT_CYCLE.md](/home/matt/Workspace/repos/trusted/Cairn/products/athena-work/DEVELOPMENT_CYCLE.md)
