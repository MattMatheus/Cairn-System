# Architecture Handoff: ARCH-20260312-toolcli-tool-interface-spec

## What Changed
- Added `docs/decisions/ADR-0007-toolcli-next-slice-context-schema-before-formal-interface.md` to record the architecture decision.
- Added `docs/product/tooling/TOOLCLI_NEXT_SLICE.md` as the recommendation note for PM and engineering.
- Updated `docs/product/tooling/TOOLCLI_V1.md` and `docs/product/tooling/MODEL_TOOL_INTERFACE_PREP.md` to reflect that bounded context/schema shaping is the next approved slice and the formal interface spec remains deferred.
- Repaired the canonical doc-test pathing under `products/work-harness/tools/` so architecture validation can run from the current repo layout.

## Why It Changed
- The roadmap needed an explicit decision between a bounded tool-cli context/schema enhancement and a broader formal tool interface spec.
- Existing design docs already constrained the answer toward the bounded context/schema slice; this cycle made that choice explicit and actionable.

## Validation Run Results
- `./products/work-harness/tools/run_doc_tests.sh`: pass
- The successful run required fixing stale work harness test-script paths and an empty-array bug in `test_no_personal_paths.sh`.

## Risks/Tradeoffs
- The broader formal tool interface remains unresolved and deferred.
- Some historical work harness docs/tests still assume older repo-era paths; this cycle fixed the canonical doc-test chain but not every historical script in the tree.

## Recommended QA Focus
- Verify ADR-0007 stays inside tool-cli v1 boundaries and does not authorize execution/runtime scope.
- Verify the recommendation note is specific enough for PM to promote the deferred context/schema engineering story.
- Verify the repaired doc-test pathing is scoped to the canonical validation chain and does not introduce path regressions.
