# Engineering Handoff: STORY-20260312-athenause-shared-validation

## What Changed
- Added `tools/platform/validate_athenause_registry.sh` as the shared AthenaUse registry validation helper.
- Wired that helper into `tools/platform/smoke_v1.sh` and `tools/platform/smoke_mongodb.sh` so shared platform smoke checks fail before downstream work when the approved registry is invalid.
- Added `products/athena-work/tools/test_athenause_registry_validation.sh` to verify both valid and invalid registry behavior.
- Updated `tools/platform/README_V1_CHECKS.md` to document that both smoke scripts now validate the approved AthenaUse registry up front.

## Why It Changed
- ADR-0005 already required approved-registry validation in shared platform checks, but the actual smoke entrypoints were not enforcing it.
- The helper keeps the validation contract in one place and gives the platform scripts a stable preflight gate with actionable errors.

## Test Updates Made
- Added `products/athena-work/tools/test_athenause_registry_validation.sh`.

## Test Run Results
- `./tools/platform/validate_athenause_registry.sh`: pass (`registry valid: 6 tools`)
- `./products/athena-work/tools/test_athenause_registry_validation.sh`: pass
- `./tools/platform/smoke_v1.sh`: pass when rerun outside the sandbox
- `shellcheck tools/platform/validate_athenause_registry.sh tools/platform/smoke_v1.sh tools/platform/smoke_mongodb.sh products/athena-work/tools/test_athenause_registry_validation.sh`: not run (`shellcheck` unavailable in this environment)

## Open Risks/Questions
- `smoke_mongodb.sh` now validates the registry as well, but its optional Mongo-backed execution path was not exercised in this cycle.
- `smoke_v1.sh` required an unrestricted rerun because sandboxed `httptest` socket binding inside AthenaMind tests is blocked; the story itself is not dependent on that limitation.

## Recommended QA Focus Areas
- Verify the helper fails with actionable output when `ATHENA_USE_REGISTRY` points to an invalid registry file.
- Confirm `smoke_v1.sh` and `smoke_mongodb.sh` invoke the helper before their downstream checks.
- Confirm documentation matches the enforced smoke-check behavior.

## New Gaps Discovered
- None beyond the already-tracked deferred context/schema story and architecture decision item.
