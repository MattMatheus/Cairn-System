# tool-cli

tool-cli is the tool-context companion to memory-cli.

It does not model every small operation as a separate tool. It models a small set of approved tool systems and the bounded capabilities Cairn is willing to expose from each one.

## Current Tool Systems

- `cairn`
- `obsidian`
- `firecrawl`
- `gitnexus`

Current rule:

- keep the top-level tool-system list small
- keep Cairn-owned capabilities under `cairn`
- expose only the specific external capabilities Cairn actually supports
- fail explicitly when selection is ambiguous
- wrap GitNexus behind `codegraph-cli` instead of exposing its full native surface

## Command Surface

- `tool-cli discover`
- `tool-cli context`
- `tool-cli inspect`
- `tool-cli list`
- `tool-cli validate`
- `codegraph-cli analyze|status|context|impact`
- `intake-cli inspect|url|file|folder|stage`
- `promote-cli inspect|note`

`inspect` prefers an exact tool-system ID, but may resolve a unique query. Ambiguous queries fail with candidate IDs rather than guessing.

## Registry Contract

V1 uses a config-backed registry with two trust tiers:

- approved: repo-backed and supported
- local: operator-managed and opt-in

Approved tool systems live under:

- `products/tool-cli/registry/approved-tools.yaml`

Local overlays are expected under:

- `.cairn/tools/registry.yaml`

The registry stores:

- tool-system identity and guidance
- bounded capabilities for each system
- capability status and availability posture
- stage affinity
- call contract
- minimal parameter schema

## Observability

tool-cli follows the memory-cli telemetry posture:

- OpenTelemetry is the tracing and metrics standard
- discovery, context emission, validation, and future execution paths emit traceable spans
- no separate observability framework should be introduced
