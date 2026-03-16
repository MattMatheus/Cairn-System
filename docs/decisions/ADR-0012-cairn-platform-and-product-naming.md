# ADR-0012: Cairn Is The Platform Umbrella And Product Naming Surface

Status: Accepted

## Context

The platform now has a clear product direction:

- personal PKM and disciplined human-agent work now
- future offline-friendly ministry utility later

That direction benefits from a stable outward-facing brand that does not overclaim, offend, or depend on explicitly divine naming.

`Cairn` is a better external anchor than direct theological naming:

- less likely to offend
- more durable across personal and ministry contexts
- symbolically appropriate as a marker, guide, and pathfinding reference

At the same time, the codebase is still actively evolving. Broad renames across repos, binaries, docs, and notes would create unnecessary churn while product boundaries are still settling.

## Decision

Use `Cairn` as the platform umbrella and outward-facing brand anchor.

Name the active product surfaces as:

- `memory-cli`
- `tool-cli`
- `work harness`

Keep that naming stable unless a later product-surface change justifies another coordinated rename.

## Consequences

Positive:

- gives the platform a coherent umbrella identity now
- avoids theological overclaim in software branding
- preserves engineering focus by preventing rename churn during active design work
- keeps open the option for a later well-coordinated rename

Tradeoffs:

- some older historical material may still require path and wording cleanup
- product naming is now clearer, but still mixes CLI-style and harness-style surfaces by design

## Required Guidance

Near-term:

- use `Cairn` in north-star, mission, and outward-facing framing
- use `Cairn` for platform-level docs, repo framing, and operator language
- use `memory-cli`, `tool-cli`, and `work harness` consistently in current code, paths, and active product docs

Later:

- only revisit naming if the product boundary itself changes again
