# Getting Started

## Summary

This section gets you from zero to a working AthenaMind memory workflow quickly, then branches into AthenaWork operating flows.

The intended platform command surface is small:

- `write`
- `retrieve`
- `bootstrap`
- `verify`
- `snapshot`

Recommended platform posture:

- repo-local runtime state under `.athena/`
- `ATHENA_HOME="$PWD/.athena"`
- `ATHENA_MEMORY_ROOT="$ATHENA_HOME/memory/default"`

## Audience

- Engineers and data scientists who want to run or tune memory-assisted coding workflows.
- Technical operators preparing telemetry and governance gates for real usage.

## Paths

1. `docs/operator/athena-mind/getting-started/installation.md`
2. `docs/operator/athena-mind/getting-started/binaries.md`
3. `docs/operator/athena-mind/getting-started/quickstart.md`
4. `docs/operator/athena-mind/getting-started/athenawork-quickstart.md`

## Next Steps

- Command deep dive: `docs/product/athena-mind/cli/commands.md`
- Product scope: `products/athena-mind/README.md`
- Repo-local runtime contract: `docs/runtime-layout.md`
- OpenTelemetry remains required across Cairn systems even while AthenaMind's public command surface stays small
