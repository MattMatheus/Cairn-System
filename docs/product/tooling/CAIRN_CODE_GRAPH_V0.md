# Cairn Code Graph V0

Status: active foundation

## Purpose

Cairn needs a scoped way to analyze old, complex, and poorly documented repositories without adopting a full ambient MCP/editor-integration lifestyle.

GitNexus is the current reference point because its core value is real:

- precomputed repository structure
- symbol context
- impact analysis
- change-impact detection

The goal is to preserve that value while stripping away the surrounding runtime posture that does not fit the personal Cairn system.

## Product Boundary

Cairn Code Graph V0 should expose a thin local contract around a repository-analysis backend.

It should do four things well right now:

1. index a repository on demand
2. report index freshness and status
3. answer symbol-context questions
4. answer blast-radius questions
It should not try to become:

- a default MCP dependency for every thread
- an automatic editor-configuration system
- a hook installer
- an owner of agent context files
- a replacement for memory-cli

### Relationship To Existing Products

- `tool-cli`: advertises code-graph capability and explains when it fits
- `work-harness`: decides when repo analysis is appropriate in a workflow
- `memory-cli`: stores distilled conclusions, not raw graph state
- GitNexus or a future replacement: backend engine for repo structure

## Minimal Command Set

Current v0 command surface:

- `codegraph analyze [path]`
- `codegraph status [path]`
- `codegraph context <symbol>`
- `codegraph impact <symbol>`

Behavior:

- `analyze`: build or refresh the repo graph explicitly
- `status`: report whether the current index is present and stale
- `context`: return a focused structural view for one symbol or target
- `impact`: return upstream or downstream blast radius for a target

Deferred:

- `codegraph detect-changes`

Current flags:

- `--repo`
- `--force`
- `--scope`
- `--direction`
- `--json`
- `--dry-run` where supported later

Explicitly not v0:

- automatic MCP server startup
- automatic hook registration
- global editor setup mutation
- repo-specific skill generation as a required path
- autonomous refactoring or rename application

## Placement In Cairn

Recommended placement:

- registry and guidance home: `products/tool-cli`
- implementation posture: a thin Cairn-facing wrapper or contract around GitNexus, not a new graph engine

Reasoning:

- the problem is tool selection and scoped availability, which belongs to tool-cli
- Cairn does not need to duplicate GitNexus's graph engine unless the retained surface proves too heavy
- keeping the contract thin makes a future backend swap possible

If later experience shows GitNexus remains too heavy even behind a thin contract, then a smaller local graph/index tool can be considered. That is not the right first move.

## Borrow From GitNexus

Useful parts to keep:

- explicit local indexing
- index freshness and status checks
- structural symbol context
- impact analysis
- current-change impact detection
- repo-local persistent graph state

What to preserve conceptually:

- structure should be precomputed, not rediscovered expensively at query time
- repo analysis should be explicit and on-demand
- results should be compact enough to complement memory-cli rather than overwhelm it

## Ignore From GitNexus

Ignore for v0:

- default MCP/editor integration posture
- automatic global setup commands
- automatic hook installation
- auto-generated agent context files as part of the base contract
- product messaging that assumes ambient agent coupling
- any environment mutation that is not explicitly requested

## PM Recommendation

PM recommendation is now implemented as a bounded Cairn-facing wrapper around GitNexus for explicit analyze, status, context, and impact flows.

Why:

- replacing the graph engine immediately would be wasteful
- adopting the full GitNexus ecosystem would add unnecessary ambient complexity
- a thin contract keeps optionality high and context windows small

Success criteria:

- repo analysis is available when needed and absent when not needed
- the retained command set helps with code archaeology and risk analysis in messy repos
- agents can discover and inspect the capability through tool-cli without inheriting full MCP noise
- important findings still end up in memory-cli as durable conclusions

## Follow-On Note

Firecrawl already triggered the stronger conclusion:

- build a smaller native Cairn intake tool instead of integrating the full external product

GitNexus currently looks more suitable as a retained backend than Firecrawl does, but the same discipline applies:

- preserve the useful core
- refuse the ambient complexity
