# Tooling Integration

Landing area for cross-product tooling contracts in Cairn.

## Purpose

Use this section for the interfaces that connect:

- agents and models
- work harness and memory-cli
- local runtime tools and future fetched artifacts

## Current Boundary

- active platform tooling lives in `tools/`
- work harness product-specific tools live in `products/work-harness/tools/`
- memory-cli command surface lives in `products/memory-cli/cmd/memory-cli/`
- Azure/bootstrap artifact retrieval is intentionally deferred until after the tool contract is finalized

## Next Intended Inputs

- model tool interface spec
- command and payload contract definitions
- runtime configuration contract
- artifact hydration rules after bootstrap scope is approved

## Prep Artifact

- `docs/product/tooling/MODEL_TOOL_INTERFACE_PREP.md`
- `docs/product/tooling/TOOLCLI_V1.md`
- `docs/product/tooling/CAIRN_INTAKE_V0.md`
- `docs/product/tooling/CAIRN_PROMOTION_V0.md`
- `docs/product/tooling/CAIRN_CODE_GRAPH_V0.md`
- `products/tool-cli/registry/approved-tools.yaml`
