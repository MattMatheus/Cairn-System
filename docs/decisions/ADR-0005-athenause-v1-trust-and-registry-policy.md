# ADR-0005: AthenaUse V1 Trust And Registry Policy

Status: Accepted

## Context

Cairn needs a governed way to expose tools to agents without either:

- giving them no tool context
- dumping a full unscoped manifest into every stage

The desired product is AthenaUse, a Go-native tool-context layer that complements AthenaMind and AthenaWork.

The key unresolved design issue was trust and source policy:

- should tools be repo-backed, local, or both?
- should all tools be treated equally?
- should stage launches inject every available tool by default?

## Decision

AthenaUse v1 will use a two-tier registry model.

### Approved Tier

Approved tools are:

- repo-backed
- curated
- reviewed
- validated
- injected by default into AthenaWork stage context

### Local Tier

Local tools are:

- additive
- user-managed
- not platform-guaranteed
- excluded from default stage-context injection
- available only when explicitly requested

V1 will use a config-backed registry as the active runtime path.

Memory-backed registry support is deferred.

V1 command scope is limited to:

- `discover`
- `context`
- `list`
- `validate`

`call` is deferred.

## Consequences

Positive:

- keeps the supported tool surface small and trustworthy
- reduces agent context pollution
- makes beta-user behavior more predictable
- allows local experimentation without weakening the approved contract

Tradeoffs:

- local tools become intentionally second-class
- execution support is deferred even when some operators may want it
- memory integration for tools remains a later enhancement

## Implementation Direction

- use `use-cli` as the binary name
- store approved registry data in committed repo-backed files
- use `products/athena-use/registry/approved-tools.yaml` as the default approved registry location
- allow local overlays through `.athena/` or explicit env override
- include support-tier markers in AthenaUse output
- validate approved registry in shared platform checks
