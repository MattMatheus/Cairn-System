# Engineering Handoff: STORY-20260312-toolcli-context-schema-output

## What Changed
- Enriched `tool-cli context` output in `products/tool-cli/cmd/tool-cli/commands.go` with structured `parameters` and `schema_summary` fields while preserving the existing top-level summary fields.
- Kept the legacy `schema` string in the context payload as a compatibility alias for the richer `schema_summary`.
- Added focused tests in `products/tool-cli/cmd/tool-cli/main_test.go` for both JSON and YAML context output.

## Why It Changed
- The previous context output only exposed a flattened schema string, which was too shallow for downstream consumers as the approved tool surface grew.
- This keeps tool-cli inside the bounded v1 context/schema slice from ADR-0007 without changing tool execution/runtime scope.

## Test Updates Made
- Added JSON context-output coverage for structured parameter fields.
- Added YAML context-output coverage for the `parameters` section and `schema_summary`.

## Test Run Results
- `gofmt -w products/tool-cli/cmd/tool-cli/commands.go products/tool-cli/cmd/tool-cli/main_test.go`
- `go test ./...` in `products/tool-cli`: pass
- `go run ./cmd/tool-cli context --stage engineering --format json` in `products/tool-cli`: pass
- `go run ./cmd/tool-cli context --stage engineering --format yaml` in `products/tool-cli`: pass
- `./products/work-harness/tools/test_launch_stage_workspace_api_adapter.sh`: pass

## Open Risks/Questions
- `launch_stage.sh` still prints tool-cli output verbatim and does not parse the richer fields yet; that is acceptable for this slice but means downstream structural consumption is still shallow.
- The registry currently does not populate schema descriptions or enum values for most approved tools, so structured fields are present but sparse on today’s catalog.

## Recommended QA Focus Areas
- Verify `tool-cli context` includes `parameters` and `schema_summary` in both JSON and YAML outputs.
- Verify tools with no schema emit stable empty shapes (`parameters: []` and empty JSON arrays).
- Verify launch-stage output remains readable and unchanged enough for current operators.

## New Gaps Discovered
- None
