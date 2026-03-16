# Engineering Handoff: STORY-20260312-toolcli-registry-parameter-metadata

## What Changed
- Enriched approved tool-cli registry metadata in `products/tool-cli/registry/approved-tools.yaml`.
- Added meaningful descriptions for `cairn.memory.verify_health` parameters.
- Expanded `cairn.workspace.validate_task_metadata` to describe its existing underlying mode/target-dir contract, including a stable `changed|all` enum for `mode`.
- Added repo-backed tests that assert the shipped approved registry and emitted context output include the richer metadata.

## Why It Changed
- The bounded context/schema slice was functionally complete, but the shipped registry was still too sparse for the richer payload to help operators or downstream consumers.
- This keeps improvements inside the approved tool-cli v1 surface and avoids reopening interface-scope questions.

## Test Updates Made
- Added approved-registry metadata assertions in `products/tool-cli/internal/registry/registry_test.go`.
- Added approved-registry context-output assertions in `products/tool-cli/cmd/tool-cli/main_test.go`.

## Test Run Results
- `gofmt -w products/tool-cli/internal/registry/registry_test.go products/tool-cli/cmd/tool-cli/main_test.go`
- `GOCACHE=/tmp/cairn-gocache go test ./...` in `products/tool-cli`: pass
- `GOCACHE=/tmp/cairn-gocache go run ./cmd/tool-cli validate` in `products/tool-cli`: pass
- `GOCACHE=/tmp/cairn-gocache go run ./cmd/tool-cli context --stage engineering --format json` in `products/tool-cli`: pass
- `env CAIRN_REQUIRED_BRANCH=main GOCACHE=/tmp/cairn-gocache ./products/work-harness/tools/launch_stage.sh engineering`: pass

## Open Risks/Questions
- The registry still only has two parameterized tools; future registry growth may surface additional metadata gaps.
- `cairn.workspace.validate_task_metadata` now exposes the script’s supported arguments in the registry, but tool-cli still remains a context product, not an execution runtime.

## Recommended QA Focus Areas
- Verify descriptions are present for all shipped parameter fields in the approved registry and context output.
- Verify the `mode` enum for `cairn.workspace.validate_task_metadata` appears as `[changed, all]`.
- Verify launch-stage output remains readable with the richer metadata.

## New Gaps Discovered
- None
