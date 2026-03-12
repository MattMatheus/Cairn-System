# AthenaWork Active Surface

## Purpose

This document separates AthenaWork’s active platform surface from imported historical material.

## Active Surface

These are part of the intended active AthenaPlatform path:

### Product entrypoints

- `products/athena-work/README.md`
- `products/athena-work/HUMANS.md`
- `products/athena-work/AGENTS.md`
- `products/athena-work/DEVELOPMENT_CYCLE.md`

### Workflow machinery

- `products/athena-work/tools/launch_stage.sh`
- `products/athena-work/tools/run_observer_cycle.sh`
- `products/athena-work/tools/run_doc_tests.sh`
- `products/athena-work/tools/check_markdown_drift.sh`
- `products/athena-work/tools/validate_intake_items.sh`
- `products/athena-work/tools/lib/`

### Contracts and prompts

- `products/athena-work/operating-system/`
- `products/athena-work/operating-system-vault/`
- `products/athena-work/stage-prompts/`
- `products/athena-work/staff-personas/`
- `products/athena-work/delivery-backlog/engineering/`
- `products/athena-work/delivery-backlog/architecture/`

## Historical Or Transitional Surface

These are preserved for migration context, but should not be treated as the default active platform surface:

### Historical handoff artifacts

- dated files in `products/athena-work/operating-system/handoff/`

### Historical archived backlog items

- dated files in `products/athena-work/delivery-backlog/architecture/archive/`

### Repo-shape-dependent tests

Many `products/athena-work/tools/test_*` files were imported from the standalone repo and still assume the older repo layout.

They should be treated as:

- useful reference for validation intent
- candidates for selective promotion or rewrite
- not automatically trusted as active platform checks yet

## Recommended Next Classification Work

1. keep and adapt a small set of active validation scripts
2. mark repo-shape-dependent tests as historical until rewritten
3. move clearly stale tests into a dedicated historical test area if they do not serve the unified platform

