# Platform Status

The work harness tools folder currently contains a mix of:

- active platform scripts
- historical validation scripts imported from the standalone repository

## Active Platform Scripts

- `launch_stage.sh`
- `run_observer_cycle.sh`
- `run_doc_tests.sh`
- `check_markdown_drift.sh`
- `check_gitnexus_readiness.sh`
- `validate_intake_items.sh`
- `lib/`

These are the work harness-owned scripts that still matter for the internal-beta path.
When in doubt, prefer the root platform checks in `tools/platform/` first, then use work harness product tools for stage flow and observer behavior.

## Historical Or Transitional Scripts

Most `test_*` scripts should currently be treated as imported validation references rather than guaranteed active platform checks.

They preserve validation intent, but many still assume the old standalone work harness repository shape.
