# CLI Commands

## Summary
Command reference for the memory-cli commands that are intentionally exposed in Cairn.

The research repo may keep a broader command surface, but the platform product exposes only the small sqlite-first workflow.

Recommended runtime pattern for Cairn:

- `CAIRN_HOME="$PWD/.cairn"`
- `CAIRN_MEMORY_ROOT="$CAIRN_HOME/memory/default"`

The CLI still defaults `--root` to `memory` if you do not provide an explicit path, but the platform-recommended operator path is repo-local `.cairn/`.

## Root Command
```bash
memory-cli <write|retrieve|bootstrap|verify|snapshot> [flags]
```

## Core Commands

These are the commands that define the practical memory-cli product:

- `write`
- `retrieve`
- `bootstrap`
- `verify`
- `snapshot`

## `write`
Creates or updates an entry.

Required:
- `--id`
- `--title`
- `--type prompt|instruction|note`
- `--domain`
- `--body` or `--body-file`
- `--stage planning|architect|pm`
- `--reviewer`
- `--decision approved|rejected`
- `--reason`
- `--risk`
- `--notes`

Optional:
- `--root` (default `memory`)
- `--source-ref`
- `--source-kind`
- `--source-type`
- `--session-id`
- `--scenario-id`
- `--memory-type`
- `--operator-verdict`
- `--telemetry-file`
- `--embedding-endpoint` (default `http://localhost:11434`)
- `--approved` (legacy alias for approved decision)
- `--rework-notes` and `--re-reviewed-by` (required for rejected decision)

Note:
- `note` is the intended write type for deliberate promotion of curated Cairn vault material into memory-cli

## `retrieve`
Runs semantic retrieval with deterministic fallback.

Required:
- `--query`

Optional:
- `--root` (default `memory`)
- `--domain`
- `--embedding-endpoint` (default `http://localhost:11434`)
- `--mode classic|hybrid` (default `classic`)
- `--top-k` (default `5`)
- `--session-id`
- `--scenario-id`
- `--memory-type`
- `--operator-verdict`
- `--telemetry-file`

Note:
- Cairn retrieval is local-first and sqlite-backed in the stripped personal product
- broader backend experimentation belongs in the research repo, not the platform product contract

## `snapshot`
Snapshot subcommands:
- `snapshot create`
  - required: `--created-by`, `--reason`
  - optional: `--scope` (`full` only), `--root`, `--session-id`
- `snapshot list`
  - optional: `--root`
- `snapshot restore`
  - required: `--snapshot-id`, restore review evidence (`--reviewer --decision --reason --risk --notes`)
  - optional: `--stage` (default `pm`), `--root`, `--session-id`, rejection fields

## `bootstrap`
Builds a memory bootstrap payload for agent startup.

Required:
- `--repo`
- `--session-id`
- `--scenario`

Optional:
- `--root` (default `memory`)
- telemetry flags (`--memory-type`, `--operator-verdict`, `--telemetry-file`)

Bootstrap payload schema:
- top-level: `repo`, `session_id`, `scenario`, `generated_at`, `memory_entries`
- `memory_entries[]`: `id`, `selection_mode`, `source_path`, `confidence`, `reason`, `type`, `domain`, `title`

## `verify`
Verification subcommands:
- `verify embeddings`
  - optional: `--root`, `--show-missing`
  - reports embedding coverage for indexed entries
- `verify health`
  - optional: `--root`, `--query`, `--domain`, `--session-id`, `--embedding-endpoint`
  - runs semantic retrieval health check and reports pass/fail

Note:
- `verify embeddings` and `verify health` are the only verification commands intentionally exposed in the platform product

## OpenTelemetry

OpenTelemetry remains required for memory-cli even though the public command surface is smaller.

- command execution should continue to emit spans
- telemetry wiring should remain part of the implementation
- CLI simplification should not remove OTel hooks

## References
- `cmd/memory-cli/main.go`
- `cmd/memory-cli/commands.go`
