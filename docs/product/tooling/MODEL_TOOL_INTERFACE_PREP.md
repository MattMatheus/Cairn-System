# Model Tool Interface Prep

Preparation document for the upcoming Cairn model tool interface spec.

Use this file as the first landing point for the formal spec once it is ready.

## Scope

The planned interface should clarify how Cairn connects:

- model providers
- agent tools
- work harness workflow actions
- memory-cli retrieval and write operations

## Current Status

This document remains a preparation artifact, not the next implementation contract.

Per `ADR-0007`, Cairn should complete the bounded tool-cli context/schema slice before promoting this prep document into a formal model/tool interface specification.

## Decisions To Capture

- command surface or RPC shape
- request and response schema
- tool discovery model
- auth and local trust boundary
- error model
- logging and traceability requirements
- repo-local runtime expectations under `.cairn/`
- OpenTelemetry instrumentation requirements
- dependency-boundary rules for tool-cli

## Deferred Until Later

- Azure artifact/bootstrap retrieval
- environment-specific provisioning details
- distribution packaging decisions

## Existing Surfaces To Reconcile

- `products/memory-cli/cmd/memory-cli/`
- `products/memory-cli/internal/telemetry/`
- `products/work-harness/tools/`
- `docs/operator/work-harness/MEMORYCLI_CLI_CONTRACT.md`
- `docs/runtime-layout.md`
- `docs/product/tooling/README.md`
