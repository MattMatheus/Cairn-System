---
id: ath-task-firefox-golden-pages-ingestion-2026-02-24
type: task
status: active
domain: research
updated: 2026-02-24
owner: matt
source_of_truth: human
sensitivity: internal
agent_last_touch: 2026-02-24T00:00:00Z
agent_intent: ingest
review_state: pending
---

# Idea: Firefox History Golden Pages Ingestion

Yes, that is a strong pattern.

Treat Firefox history as an input stream, score for golden pages, then write curated notes into an Obsidian inbox as markdown files.

## Minimal Pipeline

1. Read history from `places.sqlite` (copy or read-only).
2. Score URLs by signals:
- revisit count
- dwell/proximity (multiple related visits in session)
- domain priority (work docs, RFCs, tickets, deep references)
- recency plus repeat over days
- URL/path quality (exclude login/callback/search noise)
3. Deduplicate by canonical URL/title.
4. Generate markdown note(s) in an Obsidian inbox folder with:
- title
- URL
- why it was selected
- tags (`#golden`, `#topic/...`)
- optional short summary
5. Send to review queue, then promote to permanent notes/bookmarks.

## Safety Defaults

- Never touch live Firefox DB directly.
- Keep read-only mode default.
- Add allowlist/blocklist domains.
- Strip sensitive query params before writing notes.

## Next Step

Sketch exact CLI command set and note template.
