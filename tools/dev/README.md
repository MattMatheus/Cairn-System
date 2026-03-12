# Developer Tools

Local developer helpers for running, validating, and operating AthenaPlatform.

## Current Helpers

- `check_mongodb_local.sh`: validates the expected local Podman Mongo container and prints the standard AthenaPlatform connection contract
- `bootstrap_platform.sh`: prepares repo-local `.athena/` runtime folders and validates the shared Go toolchain
- `build_release_artifacts.sh`: builds zipped AthenaMind and AthenaUse CLI artifacts for supported release targets
