# Precompiled Binaries

## Summary

Cairn now supports a canonical build path for precompiled CLI artifacts through Azure DevOps and the shared release build script.

## Artifact Naming

Current release artifacts use this deterministic format:

- `memory-cli_<version>_darwin_arm64.zip`
- `memory-cli_<version>_windows_amd64.zip`
- `use-cli_<version>_darwin_arm64.zip`
- `use-cli_<version>_windows_amd64.zip`
- `SHA256SUMS`

`<version>` matches the release label used by the build, for example `alpha-2026-03-12` or `azure-1234`.

## Supported Build Paths

### Local build

```bash
./tools/dev/build_release_artifacts.sh --version alpha-local
```

Artifacts are written under `.athena/artifacts/releases/<version>/`.

### Azure DevOps

`azure-pipelines.yml` runs:

1. `./tools/dev/bootstrap_platform.sh`
2. AthenaWork doc tests
3. `go test ./...` in `products/athena-use`
4. `go test ./...` in `products/athena-mind`
5. `./tools/dev/build_release_artifacts.sh --version "azure-<build-id>"`

The resulting zipped binaries and `SHA256SUMS` file are published as the `athena-release-artifacts` build artifact.

## Verify Checksums

macOS/Linux:

```bash
shasum -a 256 -c SHA256SUMS
```

Windows PowerShell:

```powershell
Get-FileHash .\memory-cli_alpha-local_windows_amd64.zip -Algorithm SHA256
```

## Platform Note

- Source builds remain the default for development.
- Precompiled artifacts are now intended for alpha distribution and operator convenience, not as a replacement for source-based debugging.
