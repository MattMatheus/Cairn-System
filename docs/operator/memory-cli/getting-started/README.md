# Getting Started

## Summary

This section gets you from zero to a working memory-cli memory workflow quickly, then branches into work harness operating flows.

The intended platform command surface is small:

- `write`
- `retrieve`
- `bootstrap`
- `verify`
- `snapshot`

Recommended platform posture:

- repo-local runtime state under `.cairn/`
- `CAIRN_HOME="$PWD/.cairn"`
- `CAIRN_MEMORY_ROOT="$CAIRN_HOME/memory/default"`

## Audience

- Engineers and data scientists who want to run or tune memory-assisted coding workflows.
- Technical operators preparing telemetry and governance gates for real usage.

## Paths

1. `docs/operator/memory-cli/getting-started/installation.md`
2. `docs/operator/memory-cli/getting-started/binaries.md`
3. `docs/operator/memory-cli/getting-started/quickstart.md`
4. `docs/operator/memory-cli/getting-started/work-harness-quickstart.md`

## Next Steps

- Command deep dive: `docs/product/memory-cli/cli/commands.md`
- Product scope: `products/memory-cli/README.md`
- Repo-local runtime contract: `docs/runtime-layout.md`
- OpenTelemetry remains required across Cairn systems even while memory-cli's public command surface stays small
