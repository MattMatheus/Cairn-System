# tool-cli V1

Status: active foundation design.

## Purpose

tool-cli is the tool-surface companion to memory-cli.

work harness already launches agents with:

1. a seed prompt
2. a memory bootstrap

tool-cli adds the third governed input:

3. a scoped tool-system context

The goal is to expose a very small set of approved tool systems and only the bounded capabilities Cairn actually supports from each one.

## Product Role

- memory-cli retrieves knowledge context
- tool-cli retrieves tool-system context
- work harness composes stage prompt, memory context, and tool context

tool-cli is not a daemon, sandbox, MCP broker, or generic tool marketplace.

## V1 Scope

V1 is discovery and context emission, not a full execution runtime.

Supported commands:

- `tool-cli discover <query>`
- `tool-cli context [--stage <stage>] [--query <query>]`
- `tool-cli inspect <tool-system-id-or-query>`
- `tool-cli list [--stage <stage>] [--tag <tag>]`
- `tool-cli validate`
- `codegraph-cli analyze|status|context|impact`

Deferred from v1:

- `tool-cli call`
- memory-backed tool registry
- wide external integration surfaces

## Operating Model

V1 models tool systems, not individual one-off operations.

Current approved systems:

- `cairn`
- `obsidian`
- `firecrawl`
- `gitnexus`

Each system can expose bounded capabilities. Cairn-owned capabilities stay under the `cairn` system. External systems should expose only the specific supported integration surface, not their entire native product area.

Current external-integration posture:

- `gitnexus` is active only through the Cairn-owned `codegraph-cli` wrapper
- `firecrawl` remains planned until a similarly narrow Cairn-facing wrapper exists

## Selection Rules

- prefer exact tool-system IDs
- allow unique query resolution for operator convenience
- reject ambiguous queries with candidate IDs
- keep default context small
- include scoped capabilities only when the work explicitly calls for them
- keep planned systems and capabilities visible for planning, but out of active context by default

## Registry Model

The approved registry is:

- `products/tool-cli/registry/approved-tools.yaml`

The local overlay is:

- `.cairn/tools/registry.yaml`

Each tool-system entry describes:

- stable system identifier
- human-readable name and description
- implementation status
- operator guidance
- tags
- complementary systems
- credential reference
- bounded capabilities

Each capability describes:

- stable capability identifier
- name and description
- status
- availability posture
- stage affinity
- guidance
- call contract
- parameter schema
