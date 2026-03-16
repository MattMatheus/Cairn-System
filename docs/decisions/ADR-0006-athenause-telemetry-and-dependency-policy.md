# ADR-0006: AthenaUse Telemetry And Dependency Policy

Status: Accepted

## Context

AthenaUse is intended to become the governed tool-context layer for Cairn.

Before implementation starts, the platform needs a clear rule for:

- how tool use is instrumented
- what external dependencies AthenaUse is allowed to introduce

AthenaMind already uses OpenTelemetry for tracing and local telemetry for command behavior.

The platform should not create a second observability pattern for AthenaUse.

## Decision

AthenaUse will use OpenTelemetry as its observability standard.

Instrumentation applies to:

- registry loading
- tool discovery
- tool-context emission
- validation
- any future operator-safe execution path

For AthenaUse-specific implementation, OpenTelemetry is the only external dependency family that should be introduced beyond the Go standard library.

## Consequences

Positive:

- consistent observability posture across AthenaMind and AthenaUse
- simpler operator understanding
- better traceability for tool selection and use
- tighter dependency discipline in secure environments

Tradeoffs:

- AthenaUse must avoid convenience libraries that add non-essential dependency weight
- some parsing and validation helpers may need to stay intentionally simple in v1

## Implementation Direction

- align AthenaUse tracing patterns with `products/athena-mind/internal/telemetry/`
- emit spans for command entry and major operations
- include traceability fields in AthenaUse outputs where appropriate
- keep non-stdlib dependencies limited to OpenTelemetry unless a later ADR explicitly expands the policy
