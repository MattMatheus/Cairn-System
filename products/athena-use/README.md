# AthenaUse

AthenaUse is the tool-context companion to AthenaMind.

It gives AthenaWork a governed, scoped tool surface at stage launch so agents receive:

1. seed prompt
2. memory bootstrap
3. tool context

## V1 Role

AthenaUse v1 is a discovery and context-emission product.

Initial command surface:

- `use-cli discover`
- `use-cli context`
- `use-cli list`
- `use-cli validate`

Deferred:

- `use-cli call`
- memory-backed registry mode
- bootstrap and artifact retrieval concerns

## Registry Contract

V1 uses a config-backed registry with two trust tiers:

- approved: repo-backed and supported
- local: operator-managed and opt-in

Approved tools should live under:

- `products/athena-use/registry/approved-tools.yaml`

Local overlays are expected under repo-local runtime state:

- `.athena/tools/registry.yaml`

## Observability

AthenaUse should follow the AthenaMind telemetry posture:

- OpenTelemetry is the tracing and metrics standard
- tool discovery, context emission, validation, and later execution paths must emit traceable spans
- no separate observability framework should be introduced

## References

- `docs/product/tooling/ATHENAUSE_V1.md`
- `docs/product/tooling/MODEL_TOOL_INTERFACE_PREP.md`
- `docs/decisions/ADR-0005-athenause-v1-trust-and-registry-policy.md`
- `docs/decisions/ADR-0006-athenause-telemetry-and-dependency-policy.md`
