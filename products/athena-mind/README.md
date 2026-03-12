# AthenaMind

Slim AthenaMind product within AthenaPlatform.

## Current V1 Baseline

- markdown-first ingestion from AthenaWork content
- `sqlite` as the default backend
- optional MongoDB-backed index and embedding persistence is available for advanced local developer use
- Go CLI as the primary integration surface

## Current Status

The imported AthenaMind module is buildable and testable inside `AthenaPlatform`.

Preferred v1 usage is the local sqlite-first path through:

- `write`
- `retrieve`
- `bootstrap`
- `verify`
- `snapshot`

Current MongoDB readiness can be checked through:

- `verify mongodb`

Optional MongoDB-backed runtime can be enabled through environment:

- `ATHENA_INDEX_BACKEND=mongodb`
- `ATHENA_MONGODB_URI=...`
- `ATHENA_MONGODB_DATABASE=athenamind`

Additional historical commands remain present in the imported module, but they are not the primary product posture for the unified platform.
