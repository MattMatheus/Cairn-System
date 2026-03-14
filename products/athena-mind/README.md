# AthenaMind

Slim AthenaMind product within AthenaPlatform.

## Product Posture

AthenaMind should first succeed as a useful local memory tool for everyday work.

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

Within that small surface, `write` may store prompts, instructions, or curated notes promoted from the Athena vault.

These commands are the default operating story for AthenaMind inside AthenaPlatform.

The platform CLI is intentionally stripped to this small command surface.

## OpenTelemetry

OpenTelemetry remains required for AthenaMind and should stay wired into the CLI and supporting systems.

- command execution should continue to emit spans and telemetry events
- simplification should not remove OTel initialization or tracing hooks
- broader platform systems should keep using the same telemetry posture

## Research Boundary

The broader AthenaMind research repo may continue to carry experimental commands, alternate backends, and deeper evaluation surfaces.

## Integration Notes

- AthenaMind in AthenaPlatform is a productized subset, not the full AthenaMind research repo.
- Research-heavy docs, publication scaffolding, and broader experimentation surfaces may inform future work, but they should not define the default mental model for this product.
