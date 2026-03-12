# First Pass Sequence

## Goal

Execute the first AthenaPlatform import pass without rewriting the embedded source repositories or forcing premature content decisions.

## Scope

This pass is limited to:

1. canonical platform scaffolding
2. migration rules and decision records
3. product skills
4. shared validation and migration utilities
5. safe, low-ambiguity documentation promotion

## Not In Scope

- importing live project notes into canonical workspace locations
- trimming AthenaMind internals in place
- reconciling duplicate workflow contracts
- deleting any source content

## Sequence

1. Keep the embedded repositories intact as source references:
   - `AthenaMind/`
   - `AthenaWork/`
   - `work/AthenaWork/`
2. Establish canonical root destinations:
   - `products/`
   - `tools/`
   - `workspace/`
   - `docs/`
3. Promote product skills into canonical product-owned paths.
4. Promote reusable utilities into `tools/` with path corrections for AthenaPlatform.
5. Add migration and decision docs that define import rules before content moves.
6. Prepare the second pass:
   - AthenaMind slim import
   - AthenaWork product import
   - workspace structure import

## Exit Criteria

The first pass is complete when:

- canonical directories exist
- migration rules are written
- backend policy is recorded
- product skills exist in canonical destinations
- at least one shared migration tool and one shared validation tool are ready under `tools/`

## Second Pass Inputs

Before the next pass starts, confirm:

- which AthenaMind backends are in v1 import scope beyond `sqlite` and optional Mongo
- whether the workspace should ship example content or structure-only defaults
- whether AthenaWork prompts should live only under `products/athena-work/` or also be projected into workspace docs

