# ADR-0004: Repo-Local Runtime Boundary

## Status

Accepted

## Context

The current power-user workflow relies on a large shared workspace that mixes:

- vault content
- tool references
- source code
- generated artifacts
- local operational state

That pattern works for experienced operators, but it is a poor default for internal beta users because it hides boundaries and creates unnecessary context pollution.

## Decision

Cairn will adopt a repo-local `.cairn/` directory as the default runtime boundary.

The repository itself remains the committed product and documentation surface.

The `.cairn/` directory is the default location for:

- work harness runtime state
- memory-cli memory roots
- fetched bootstrap artifacts
- caches
- run outputs
- machine-local config

The `.cairn/` directory must be excluded from version control by default.

## Consequences

Positive:

- each repository gets an isolated Cairn operating area
- onboarding becomes clearer for internal beta users
- tooling can assume predictable local paths
- generated state stays out of the committed tree

Negative:

- advanced users may need explicit override support for shared setups
- bootstrap tooling must hydrate repo-local runtime folders intentionally

## Follow-Up

1. Use `.cairn/artifacts/` as the default bootstrap destination for fetched Azure-delivered assets.
2. Add environment-variable override support where needed for advanced users.
3. Continue pruning docs and scripts that assume a shared global workspace.
