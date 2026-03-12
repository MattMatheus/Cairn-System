# Tooling Integration

Landing area for cross-product tooling contracts in AthenaPlatform.

## Purpose

Use this section for the interfaces that connect:

- agents and models
- AthenaWork and AthenaMind
- local runtime tools and future fetched artifacts

## Current Boundary

- active platform tooling lives in `tools/`
- AthenaWork product-specific tools live in `products/athena-work/tools/`
- AthenaMind command surface lives in `products/athena-mind/cmd/memory-cli/`
- Azure/bootstrap artifact retrieval is intentionally deferred until after the tool contract is finalized

## Next Intended Inputs

- model tool interface spec
- command and payload contract definitions
- runtime configuration contract
- artifact hydration rules after bootstrap scope is approved

## Prep Artifact

- `docs/product/tooling/MODEL_TOOL_INTERFACE_PREP.md`
- `docs/product/tooling/ATHENAUSE_V1.md`
- `products/athena-use/registry/approved-tools.yaml`
