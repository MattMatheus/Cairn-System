# ADR-0012: Cairn Is The External Brand Anchor While Athena Names Remain Internal Working Names

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

Use `Cairn` as the external umbrella brand and future outward-facing product anchor.

Keep the existing Athena family names as internal working names for now:

- `AthenaPlatform`
- `AthenaMind`
- `AthenaUse`
- `AthenaWork`

Do not attempt a broad rename until the relevant command, repo, and documentation surfaces are stable enough for a coordinated pass.

## Consequences

Positive:

- gives the platform a cleaner outward identity now
- avoids theological overclaim in software branding
- preserves engineering focus by preventing rename churn during active design work
- keeps open the option for a later well-coordinated rename

Tradeoffs:

- internal and external naming will diverge for a while
- some docs must explain the distinction explicitly
- eventual rename work is deferred, not eliminated

## Required Guidance

Near-term:

- use `Cairn` in north-star, mission, and outward-facing framing
- continue using Athena names in current code, paths, and active product docs unless a specific rename is justified

Later:

- perform one deliberate rename pass when stability is high enough to change repo names, binaries, docs, and references together
