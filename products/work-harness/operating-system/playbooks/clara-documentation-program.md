# Clara Documentation Program

## Summary
Execution brief for building and maintaining comprehensive user-facing documentation with Technical Writer Clara.

## Intended Audience
- Product Manager and operator coordinating release readiness.
- Technical Writer Clara executing documentation cycles.

## Preconditions
- Source-of-truth artifacts are available:
  - accepted ADRs and platform decisions in `docs/decisions/`
  - implementation stories in `products/work-harness/delivery-backlog/engineering/done/`
  - QA evidence in `products/work-harness/delivery-backlog/engineering/done/`
- Documentation backlog stories exist in `products/work-harness/delivery-backlog/engineering/intake/` or `products/work-harness/delivery-backlog/engineering/active/`.

## Main Flow
1. PM promotes documentation stories into active queue based on value/risk.
2. Clara executes doc story scope and updates current user-facing docs under `docs/`, `workspace/docs/`, and product readmes as appropriate.
3. Clara removes stale/conflicting docs and captures remaining gaps.
4. Clara validates docs with `products/work-harness/tools/run_doc_tests.sh` and ensures publish workflow eligibility.
5. QA validates documentation acceptance criteria and references.
6. Observer runs, then one cycle commit is made using `cycle-<cycle-id>`.

## Operating Modes
### Continuous Mode (Every Relevant Cycle)
- Trigger when any done story changes user-visible behavior.
- Required outputs:
  - updated docs pages
  - linked evidence references
  - gaps captured as intake items
  - publish pipeline-ready docs changes (repo source updated; no direct host edits)

### Release Mode (On-Demand Before Release Decision)
- Trigger during release checkpoint preparation.
- Required outputs:
  - updated release-facing documentation artifacts in the current platform docs surface
  - known limitations section aligned with QA evidence
  - cross-check with release bundle in `products/work-harness/operating-system/handoff/`

## Failure Modes
- Documentation claims drift from accepted ADR or QA-backed behavior.
- Core user journeys lack complete, executable guidance.
- Release notes diverge from release bundle ship/hold decision.

## References
- `products/work-harness/staff-personas/Technical Writer - Clara.md`
- `products/work-harness/operating-system-vault/STAGE_EXIT_GATES.md`
- `products/work-harness/tools/run_observer_cycle.sh`
- `products/work-harness/tools/build_docs_site.sh`
- `products/work-harness/operating-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`
