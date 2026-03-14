# ADR-0009: AthenaUse Uses Firecrawl As A Scoped Ingestion Tool, Not An Ambient Web Runtime

Status: Accepted

## Context

The personal AthenaPlatform system needs help with disorganized external material as much as it needs help with disorganized codebases.

Firecrawl is under consideration because it can turn messy web pages and similar external sources into markdown-oriented output that fits the local Athena workflow:

- ingest unstructured material
- normalize it into markdown
- review it in Obsidian
- selectively promote durable knowledge into AthenaMind

The same focus constraint that applies to GitNexus also applies here:

- tools should be discoverable
- tools should be understandable from a small skill or tool-context payload
- tools should only be brought into a session when the work requires them

AthenaUse must not become a back door for ambient browser or MCP-heavy startup behavior.

## Decision

AthenaUse will treat Firecrawl as an approved but scoped ingestion tool.

Approved posture:

- useful for research, documentation intake, and messy web-content normalization
- suitable for generating markdown that can be reviewed before import
- complementary to Obsidian and AthenaMind rather than a replacement for either

Boundary posture:

- not loaded by default into every thread
- not treated as a general-purpose ambient browsing runtime
- not treated as canonical storage

Operating model:

- Firecrawl gathers and normalizes external source material
- Obsidian holds candidate markdown for human review
- AthenaMind stores curated conclusions or accepted canonical notes
- AthenaUse advertises Firecrawl as an optional ingestion capability, not a universal baseline

## Consequences

Positive:

- gives the personal system a strong intake path for messy external information
- keeps AthenaMind focused on curated memory instead of raw web ingestion
- preserves operator review through Obsidian before canonical import
- remains consistent with the no-bloated-context-window rule

Tradeoffs:

- external ingestion remains an explicit step instead of an always-present background capability
- some research sessions will require deliberate opt-in to Firecrawl
- AthenaUse must provide enough metadata for agents to know Firecrawl exists without forcing it into every startup

## Required Boundaries

AthenaUse may:

- register Firecrawl as an approved tool candidate
- mark it as research, documentation, and intake oriented
- surface it through discovery, listing, and context when explicitly requested
- describe it as a markdown-ingestion complement to Obsidian and AthenaMind

AthenaUse may not:

- make Firecrawl an ambient default dependency for all stages
- imply that web ingestion should automatically become canonical memory
- authorize broad browser or MCP-heavy startup behavior by default

## Follow-On Work

- extend AthenaUse registry metadata so required tools, scoped tools, and default-injection behavior are explicit
- add a per-tool inspection or evaluation command oriented around fit-for-session decisions
- keep Obsidian documented as the required review surface in the personal workflow
