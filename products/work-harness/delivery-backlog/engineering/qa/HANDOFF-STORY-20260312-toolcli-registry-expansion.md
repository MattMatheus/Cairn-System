# Engineering Handoff: STORY-20260312-toolcli-registry-expansion

## What Changed
- Added four approved tool-cli registry entries in `products/tool-cli/registry/approved-tools.yaml`.
- Added coverage in `products/tool-cli/internal/registry/registry_test.go` to validate the shipped approved registry file from the repo.
- Added `products/work-harness/operating-system/handoff/CURRENT.md` to `.gitignore` so the local handoff file does not show up as normal tracked work.
- Created the missing `products/work-harness/delivery-backlog/engineering/qa/` lane and moved this story into it with `status: qa`.

## Why It Changed
- tool-cli stage launch integration was live, but the approved registry only exposed two tools and did not yet provide a credible shared tool-context surface.
- The added entries expand the approved catalog using existing repo-supported commands without widening the trust model or execution scope.

## Test Updates Made
- Added a registry test that loads the repo's approved registry file through the normal path resolver and validates the resulting contract.

## Test Run Results
- `gofmt -w products/tool-cli/internal/registry/registry_test.go`
- `go test ./...` in `products/tool-cli`: pass
- `go run ./cmd/tool-cli validate` in `products/tool-cli`: pass (`registry valid: 6 tools`)
- `go run ./cmd/tool-cli list --stage engineering --format text` in `products/tool-cli`: pass
- `go run ./cmd/tool-cli discover validation --format text` in `products/tool-cli`: pass

## Open Risks/Questions
- The new tool entries are approved and discoverable, but the shared platform validation path still does not enforce tool-cli validation; that is the next active story.
- The registry remains command-string based and intentionally narrow; richer schema/context shaping is still deferred pending the architecture decision item.

## Recommended QA Focus Areas
- Verify each added command path exists and matches its description and stage affinity.
- Confirm the registry test is validating the shipped file rather than only a fixture.
- Confirm ignoring `CURRENT.md` does not mask any other handoff artifacts unintentionally.

## New Gaps Discovered
- `products/work-harness/delivery-backlog/engineering/qa/` was missing even though the documented state model requires it; this cycle created the lane and its queue README so the workflow can proceed as documented.
