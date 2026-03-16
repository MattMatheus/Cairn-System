# ADR-0008: tool-cli Uses GitNexus As A Scoped Analysis Tool, Not An Ambient Default

Status: Accepted

## Context

tool-cli exists to keep the tool surface small, relevant, and stage-scoped.

The personal Cairn direction is intentionally resistant to loading broad MCP server surfaces into every agent session. The problem is not MCP in principle; the problem is ambient tool noise, startup overhead, and focus drift when agents inherit large generic manifests by default.

The target use case for GitNexus is narrow and practical:

- old, complex, poorly documented codebases
- structural code archaeology
- impact analysis
- repo-local graph queries that help agents orient quickly

At the same time, memory-cli already owns durable memory and working conclusions. GitNexus should not become the long-term system of record for operator knowledge.

## Decision

tool-cli will treat GitNexus as an approved but scoped code-analysis tool.

Approved posture:

- useful for engineering and architecture work on disorganized repositories
- suitable for repo-scoped structural graph queries
- complementary to memory-cli rather than a replacement for it

Boundary posture:

- not loaded by default into every thread
- not part of ambient stage startup unless explicitly requested or strongly justified by the stage
- not treated as the durable memory layer

Operating model:

- GitNexus finds structural relationships in a repository
- memory-cli stores distilled conclusions, decisions, and reusable notes
- work harness decides when GitNexus should be present in a given workflow
- tool-cli advertises GitNexus as an optional tool surface, not a universal baseline

## Consequences

Positive:

- preserves agent focus by avoiding broad default MCP/tool manifest injection
- gives engineering work a high-leverage structural analysis tool when needed
- keeps memory-cli cleanly positioned as durable memory instead of repo graph middleware
- fits the personal-platform preference for on-demand capability over ambient complexity

Tradeoffs:

- some sessions will require explicit opt-in to gain GitNexus context
- cross-repo structural state remains tool-local and transient unless conclusions are written into memory-cli
- tool-cli must carry enough metadata to explain why GitNexus is scoped and not default

## Required Boundaries

tool-cli may:

- register GitNexus as an approved tool candidate
- mark it as engineering/architecture oriented
- surface it through discovery, listing, and context when explicitly requested
- describe it as complementary to memory-cli for messy codebase work

tool-cli may not:

- make GitNexus a default ambient dependency for all stages
- assume MCP-backed tools belong in every agent startup
- blur the product boundary between structural code analysis and durable memory

## Follow-On Work

- add registry metadata that distinguishes required tools from scoped optional tools
- add a per-tool inspection/evaluation command so operators can see why a tool is or is not a fit for a session
- evaluate Firecrawl separately under the same focus-preserving policy
