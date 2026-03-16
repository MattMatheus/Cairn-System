# memory-cli

Slim memory-cli product within Cairn.

## Product Posture

memory-cli should first succeed as a useful local memory tool for everyday work.

The default product posture is:

- markdown-first ingestion
- local-first runtime
- `sqlite` as the default storage path
- Go CLI as the main interface
- small, dependable command surface before advanced retrieval infrastructure

## Preferred V1 Path

The practical v1 path is the local sqlite-first workflow:

- `write`
- `retrieve`
- `bootstrap`
- `verify`
- `snapshot`

Within that small surface, `write` may store prompts, instructions, or curated notes promoted from the Cairn vault.

These commands are the default operating story for memory-cli inside Cairn.

The platform CLI is intentionally stripped to this small command surface.

## OpenTelemetry

OpenTelemetry remains required for memory-cli and should stay wired into the CLI and supporting systems.

- command execution should continue to emit spans and telemetry events
- simplification should not remove OTel initialization or tracing hooks
- broader platform systems should keep using the same telemetry posture

## Research Boundary

The broader memory-cli research repo may continue to carry experimental commands, alternate backends, and deeper evaluation surfaces.

## Integration Notes

- memory-cli in Cairn is a productized subset, not the full memory-cli research repo.
- Research-heavy docs, publication scaffolding, and broader experimentation surfaces may inform future work, but they should not define the default mental model for this product.
