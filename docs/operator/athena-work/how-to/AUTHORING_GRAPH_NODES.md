---
title: Authoring Graph Nodes
description: How to write AthenaWork knowledge-base nodes so agents can scan, route, and compose them without loading everything.
doc_type: how_to
status: draft
---

# Authoring Graph Nodes

## Summary

Write graph nodes as one complete thought per file. Small, specific nodes are easier for agents to scan and compose than broad handbook pages. Use [[DEVELOPMENT_COGNITION_GRAPH]] for internal structure and [[CANONICAL_SOURCE_REGISTRY]] for external tool truth.

## Intended Audience

- people adding or refactoring AthenaWork docs
- prompt authors designing stage traversal behavior

## Preconditions

- the topic is durable knowledge, not active sprint state
- the file can be expressed as one claim, procedure, heuristic, decision, or example

## Main Flow

When authoring a node:

1. choose one node type
2. write frontmatter with a description
3. keep the summary short enough to scan quickly
4. add only links that materially change what an agent should read next
5. attach canonical source records when external tool semantics matter

## Required Structure

Every new node should include:

- frontmatter with `title`, `description`, `doc_type`, and `status`
- `Summary`
- `Intended Audience`
- `Preconditions`
- `Main Flow`
- `Failure Modes`
- `References`

Recommended frontmatter:

```yaml
title: Common Edit Paths
description: High-signal codebase map for file clusters that usually need to change together.
doc_type: procedure
status: active
domain: codebase
canonical_sources:
  - react-19
related:
  - [[VERIFICATION_PATTERNS]]
  - [[PRODUCT_BOUNDARY]]
```

## Link Rules

Use links when they answer one of these questions:

- what should the agent read next
- what is a prerequisite
- what is the contrasting rule
- what is the authoritative external source

Do not add links only because two pages mention a similar noun.

## Node Size Rules

Aim for:

- one claim, procedure, or decision per file
- summary in under 120 words
- full page short enough that an agent can load it without sacrificing the rest of its context budget

If a page needs multiple unrelated sections, split it into a MOC plus child nodes.

## Canonical Source Rules

If the page tells an agent how AthenaWork uses a tool, attach a canonical source record for the tool. Keep the node focused on AthenaWork-specific usage and let the source record point to the vendor docs.

Good pattern:

- node says how AthenaWork uses Podman locally
- source record says which Podman docs and version policy to trust

Bad pattern:

- node pastes large excerpts from Podman docs into AthenaWork docs

## Failure Modes

- mixed concerns inside one file
- frontmatter without a useful description
- links that do not change routing behavior
- stale version claims with no source record
- process docs that silently include tool semantics from memory

## References

- [[DEVELOPMENT_COGNITION_GRAPH]]
- [[CANONICAL_SOURCE_REGISTRY]]
