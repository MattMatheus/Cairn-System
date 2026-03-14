# Athena Promotion V0

Status: initial implementation landed

## Purpose

Athena intake produces review artifacts in the vault. Promotion is the deliberate next step that moves selected material into AthenaMind.

This path is intentionally hand-curated.

## Boundary

Promotion V0 should do one thing well:

- take a reviewed vault note and write it into AthenaMind as a durable local note entry with provenance

It should not:

- auto-promote inbox artifacts
- infer trust automatically
- replace human judgment

## Current Command Surface

- `promote-cli inspect <note-path>`
- `promote-cli note <note-path> --reviewer <name> --reason <text> --risk <text> --notes <text>`

Current behavior:

- reads Athena vault frontmatter and markdown body
- trims the vault top-level heading before AthenaMind write to avoid duplicate titles
- preserves provenance into AthenaMind metadata:
  - `source_ref`
  - `source_kind`
  - `source_type`
- serializes promotion writes per AthenaMind root to avoid index races during iteration
- writes into AthenaMind as `type: note`

## Relationship To Existing Products

- Athena vault: human review and curation surface
- `promote-cli`: thin wrapper for deliberate promotion
- `athena-mind`: canonical local memory store

## Current Rule

Review in Athena first. Promote deliberately. Retrieval happens from AthenaMind after promotion, not directly from raw inbox artifacts.

## Follow-On

Future improvements should be driven by real curation use, not guessed automation:

- promotion checklists
- note-to-claim extraction
- claim/artifact linking
- safer bulk promotion only if truly needed
