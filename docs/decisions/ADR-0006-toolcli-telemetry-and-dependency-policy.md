# ADR-0006: tool-cli Telemetry And Dependency Policy

Status: Accepted

## Context

tool-cli is intended to become the governed tool-context layer for Cairn.

Before implementation starts, the platform needs a clear rule for:

- how tool use is instrumented
- what external dependencies tool-cli is allowed to introduce

memory-cli already uses OpenTelemetry for tracing and local telemetry for command behavior.

The platform should not create a second observability pattern for tool-cli.

## Decision

tool-cli will use OpenTelemetry as its observability standard.

Instrumentation applies to:

- registry loading
- tool discovery
- tool-context emission
- validation
- any future operator-safe execution path

For tool-cli-specific implementation, OpenTelemetry is the only external dependency family that should be introduced beyond the Go standard library.

## Consequences

Positive:

- consistent observability posture across memory-cli and tool-cli
- simpler operator understanding
- better traceability for tool selection and use
- tighter dependency discipline in secure environments

Tradeoffs:

- tool-cli must avoid convenience libraries that add non-essential dependency weight
- some parsing and validation helpers may need to stay intentionally simple in v1

## Implementation Direction

- align tool-cli tracing patterns with `products/memory-cli/internal/telemetry/`
- emit spans for command entry and major operations
- include traceability fields in tool-cli outputs where appropriate
- keep non-stdlib dependencies limited to OpenTelemetry unless a later ADR explicitly expands the policy
