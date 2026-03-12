# Model Tool Interface Prep

Preparation document for the upcoming AthenaPlatform model tool interface spec.

Use this file as the first landing point for the formal spec once it is ready.

## Scope

The planned interface should clarify how AthenaPlatform connects:

- model providers
- agent tools
- AthenaWork workflow actions
- AthenaMind retrieval and write operations

## Current Status

This document remains a preparation artifact, not the next implementation contract.

Per `ADR-0007`, AthenaPlatform should complete the bounded AthenaUse context/schema slice before promoting this prep document into a formal model/tool interface specification.

## Decisions To Capture

- command surface or RPC shape
- request and response schema
- tool discovery model
- auth and local trust boundary
- error model
- logging and traceability requirements
- repo-local runtime expectations under `.athena/`
- OpenTelemetry instrumentation requirements
- dependency-boundary rules for AthenaUse

## Deferred Until Later

- Azure artifact/bootstrap retrieval
- environment-specific provisioning details
- distribution packaging decisions

## Existing Surfaces To Reconcile

- `products/athena-mind/cmd/memory-cli/`
- `products/athena-mind/internal/telemetry/`
- `products/athena-work/tools/`
- `docs/operator/athena-work/ATHENAMIND_CLI_CONTRACT.md`
- `docs/runtime-layout.md`
- `docs/product/tooling/README.md`
