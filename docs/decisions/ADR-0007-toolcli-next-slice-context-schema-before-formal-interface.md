# ADR-0007: tool-cli Next Slice Prioritizes Context Schema Before Formal Interface Spec

Status: Accepted

## Context

tool-cli v1 is now live enough to validate three things:

- approved registry curation is working
- stage launch consumes approved tool context
- shared platform checks can enforce registry validity

The next roadmap fork is whether to:

1. deepen the current context-emission surface with stronger schema/context output, or
2. begin a formal cross-platform tool interface specification immediately

The current design set already constrains this decision:

- ADR-0005 limits tool-cli v1 to discovery, listing, context, and validation
- ADR-0006 keeps the implementation intentionally narrow and OpenTelemetry-aligned
- `docs/product/tooling/TOOLCLI_V1.md` explicitly defers full JSON Schema support and execution/runtime expansion
- `docs/product/tooling/MODEL_TOOL_INTERFACE_PREP.md` exists as preparation material, not an approved implementation contract

## Decision

The next approved tool-cli implementation slice is:

- bounded context/schema shaping for `tool-cli context`

The formal tool interface specification remains deferred until after that bounded slice is complete and validated.

## Consequences

Positive:

- keeps work inside the already-approved tool-cli v1 boundary
- improves stage-launch usefulness without authorizing execution/runtime expansion
- gives PM and engineering a concrete follow-on implementation target
- preserves the prep document as a staging area instead of prematurely turning it into a contract

Tradeoffs:

- the broader model/tool interface remains unresolved for now
- some future contract questions will be answered incrementally instead of in one formal pass
- downstream consumers must continue treating the current tooling surface as product-specific rather than platform-final

## Required Boundaries

The bounded next slice may:

- enrich `tool-cli context` output with schema-oriented fields for required parameters
- tighten tests around emitted context shape
- update work harness stage-launch consumption in the same cycle if needed

The bounded next slice may not:

- introduce `tool-cli call`
- define a full platform-wide execution/runtime contract
- add full JSON Schema support
- reopen Azure/bootstrap scope
- weaken approved/local trust-tier rules

## Follow-On Documentation Targets

Before implementation stories are promoted from this decision:

- align `docs/product/tooling/TOOLCLI_V1.md` with the selected next slice
- update `docs/product/tooling/MODEL_TOOL_INTERFACE_PREP.md` to note that the formal interface spec remains deferred
- keep backlog items for formal interface work in architecture lanes until the bounded slice is complete
