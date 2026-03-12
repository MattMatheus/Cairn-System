# AthenaUse Next Slice Recommendation

## Decision

Implement stronger AthenaUse context/schema output next.

Do not start the formal cross-platform tool interface spec as the next implementation slice.

## Why

- The current AthenaUse v1 design already defines `context` as the stage-launch integration surface.
- Registry and shared-validation work are now in place, so the highest-leverage remaining gap is richer, still-bounded context output.
- A formal interface spec would broaden scope too early and pull execution/runtime questions back into a v1 surface that explicitly defers them.

## Accepted Boundary

The next engineering work should stay inside these limits:

- improve `use-cli context` to emit stronger schema-oriented fields
- preserve approved/local trust-tier behavior
- preserve discovery/context/validate/list as the v1 command surface
- update tests for output shape and any stage-launch consumer changes

## Deferred Boundary

The following remain deferred after this recommendation:

- full model/tool interface specification
- `use-cli call`
- full JSON Schema support
- Azure/bootstrap retrieval concerns
- broader execution/auth/runtime contracts

## Effects On Planning

- The existing engineering intake item for context/schema output can be promoted when PM is ready.
- Future formal interface work should stay in architecture intake until the bounded context/schema slice has shipped and its consumer needs are better understood.

## Required References

- `docs/decisions/ADR-0005-athenause-v1-trust-and-registry-policy.md`
- `docs/decisions/ADR-0006-athenause-telemetry-and-dependency-policy.md`
- `docs/decisions/ADR-0007-athenause-next-slice-context-schema-before-formal-interface.md`
- `docs/product/tooling/ATHENAUSE_V1.md`
- `docs/product/tooling/MODEL_TOOL_INTERFACE_PREP.md`
