# Precompiled Binaries

## Summary

Precompiled binaries are not yet a canonical AthenaPlatform distribution path.

This page is retained as a historical note from AthenaMind’s prior standalone distribution model.

## Historical Artifact Naming

Release artifacts use this deterministic format:

- `memory-cli_<version>_darwin_amd64.zip`
- `memory-cli_<version>_darwin_arm64.zip`
- `memory-cli_<version>_windows_amd64.zip`
- `memory-cli_<version>_windows_arm64.zip`
- `SHA256SUMS`

`<version>` matches the git tag, for example `v0.2.0`.

## Historical Download Flow

1. Open releases: `https://github.com/MattMatheus/AthenaMind/releases`
2. Select the desired version tag.
3. Download the artifact for your OS/architecture.
4. Verify checksum using `SHA256SUMS`.
5. Unzip and place `memory-cli` (or `memory-cli.exe`) in your PATH.

## Historical Automation Flow

Latest release metadata:

```bash
curl -fsSL https://api.github.com/repos/MattMatheus/AthenaMind/releases/latest
```

Example: download latest Apple Silicon binary:

```bash
curl -fsSL -o memory-cli_darwin_arm64.zip \
  https://github.com/MattMatheus/AthenaMind/releases/latest/download/memory-cli_$(curl -fsSL https://api.github.com/repos/MattMatheus/AthenaMind/releases/latest | jq -r .tag_name)_darwin_arm64.zip
```

Example: download a pinned version (recommended for reproducibility):

```bash
VERSION=v0.2.0
curl -fsSL -o memory-cli_${VERSION}_windows_amd64.zip \
  https://github.com/MattMatheus/AthenaMind/releases/download/${VERSION}/memory-cli_${VERSION}_windows_amd64.zip
```

## Verify Checksums

macOS/Linux:

```bash
shasum -a 256 -c SHA256SUMS
```

Windows PowerShell:

```powershell
Get-FileHash .\memory-cli_v0.2.0_windows_amd64.zip -Algorithm SHA256
```

## Platform Note

- AthenaPlatform currently treats source builds as the primary supported path.
- If a platform release pipeline is added later, this document should be rewritten around the canonical AthenaPlatform release process.
