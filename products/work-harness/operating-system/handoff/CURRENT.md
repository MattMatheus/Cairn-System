# Current Process Handoff

## Current Stage
- Cairn platform rename is complete and the repo no longer uses `Athena` naming internally.
- Active platform surface is unified under `products/`, `tools/`, `workspace/`, `docs/`, and repo-local runtime `.cairn/`.
- tool-cli now models a small set of tool systems with bounded capabilities instead of many operation-level tools.

## Active Process Stories
- memory-cli:
  - intentionally restrained; no broad expansion work is active right now
- work harness:
  - now owns explicit GitNexus readiness preflight through `products/work-harness/tools/check_gitnexus_readiness.sh`
  - startup/operator docs include GitNexus preflight only when codegraph work is likely
- tool-cli:
  - registry shape is now tool-system based
  - approved systems are `cairn`, `obsidian`, `firecrawl`, and `gitnexus`
  - `gitnexus` is active only through the bounded `codegraph-cli` wrapper
  - `firecrawl` remains planned
  - approved registry lives at `products/tool-cli/registry/approved-tools.yaml`
  - bounded GitNexus wrapper commands exist at `products/tool-cli/cmd/codegraph-cli`

## Risks
- Historical work harness material still exists for traceability; active paths are much cleaner, but not every historical reference has been normalized.
- GitNexus backend is not yet runnable in the current shell environment without additional local readiness:
  - `node` must be available on `PATH`
  - either `CAIRN_GITNEXUS_BIN` must point at a runnable binary, or the local checkout under `repos/untrusted/GitNexus/gitnexus` must be built so `dist/cli/index.js` exists
- Full `go test ./...` in `products/tool-cli` still hits the known sandbox restriction in `internal/intake` because `httptest.NewServer` cannot open a local listener here.

## Next Improvement Target
- Bring the local GitNexus backend to readiness in a controlled way:
  - install/confirm `node`
  - build the local GitNexus checkout or point Cairn at a pinned binary
  - rerun `products/work-harness/tools/check_gitnexus_readiness.sh`
- After GitNexus is runnable, exercise the bounded path end to end:
  - `tool-cli inspect gitnexus`
  - `codegraph-cli status --repo <path>`
  - `codegraph-cli analyze --repo <path>`
  - `codegraph-cli context --repo <path> <symbol>`
  - `codegraph-cli impact --repo <path> --direction upstream <symbol>`
- Keep the external-tool posture narrow. Do not widen GitNexus or Firecrawl beyond Cairn-owned wrappers without explicit review.

## Verification Snapshot
- `go test ./cmd/codegraph-cli ./cmd/tool-cli ./internal/codegraph ./internal/registry` passed in `products/tool-cli`
- `products/work-harness/tools/test_gitnexus_readiness.sh` passed
- `products/work-harness/tools/check_gitnexus_readiness.sh` currently fails with `node is required to run a built local GitNexus checkout`
