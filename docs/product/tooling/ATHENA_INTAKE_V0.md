# Athena Intake V0

Status: initial implementation landed

## Purpose

AthenaPlatform needs a lightweight ingestion path for messy external material without adopting a heavy crawling platform as a default dependency.

The target material is practical and narrow:

- personal blogs
- documentation sites
- random exported folders of files
- markdown, text, HTML, PDF, and similar local artifacts

This is not a general web-scale crawling problem.

## Product Boundary

Athena Intake V0 should be a small local-first ingestion tool that prepares candidate markdown for review in Obsidian before anything is promoted into AthenaMind.

It should do three things well:

1. fetch or read messy source material
2. normalize it into reviewable markdown plus metadata
3. stage the result into the Athena inbox or a designated intake lane

It should not try to become:

- a general-purpose browser automation system
- a social-media crawler
- a distributed crawl service
- an always-on agent runtime
- a replacement for AthenaMind

### Relationship To Existing Products

- `athena-use`: advertises intake capability and explains when it fits
- `athena-work`: decides when intake is appropriate in a stage workflow
- `Athena` Obsidian vault: review and triage surface for candidate artifacts
- `athena-mind`: canonical curated memory after human acceptance

## Minimal Command Set

Implemented v0 command surface:

- `intake url <target>`
- `intake folder <path>`
- `intake file <path>`
- `intake inspect <target>`
- `intake stage <input> [--into <vault-path>]`

Current AthenaUse registry surface:

- `athena.intake.inspect_source`
- `athena.intake.normalize_source`

Behavior:

- `url`: fetch a single page and extract main content to markdown
- `folder`: walk a local folder and normalize supported markdown, text, HTML, and PDF files into one intake artifact
- `file`: normalize one markdown, text, HTML, or PDF file
- `inspect`: show what would be ingested, expected output format, and risk notes without writing anything
- `stage`: write candidate output into Obsidian with frontmatter and provenance

Implemented flags:

- `--out`
- `--into`
- `--vault`
- `--title`
- `--format=json`

Deferred flags and features:

- `--title`
- `--tags`
- `--source-type`
- `--max-depth` for shallow docs/blog traversal
- `--dry-run`
- DOCX normalization

Explicitly not v0:

- browser sessions
- authenticated browser automation
- search-engine orchestration
- multi-worker crawl queues
- agentic extraction loops

## Placement In AthenaPlatform

Recommended placement:

- product home: `products/athena-use`
- implementation shape: a small sibling binary under AthenaUse, such as `cmd/intake-cli/`

Reasoning:

- intake is a tool-surface concern, not a memory-store concern
- it fits the AthenaUse rule of discoverable, scoped capabilities
- it keeps AthenaMind from becoming the raw ingestion layer

The alternative is a separate product, but that is not justified yet. V0 is small enough to live as a tightly scoped AthenaUse-adjacent utility until real complexity proves otherwise.

## Borrow From Firecrawl

Useful ideas to borrow:

- HTML to markdown normalization as the main output path
- consistent source metadata capture
- simple scrape contract for a single URL
- optional shallow map or link discovery for docs sites
- support for mixed source formats like HTML and PDF now, and DOCX later if justified

What to preserve conceptually:

- ingestion output should be LLM-friendly and human-reviewable
- markdown should be first-class output, not a side effect
- provenance should remain attached to every artifact

## Ignore From Firecrawl

Ignore for v0:

- hosted-service assumption
- full crawl/search/browser/agent product surface
- MCP-oriented setup path
- cloud-first pricing or account model
- browser automation and session control
- generalized website operations platform behavior
- monorepo operational weight that exceeds the personal-system use case

## PM Recommendation

PM recommendation is to build Athena Intake V0 instead of integrating Firecrawl directly.

Why:

- lower operational burden
- tighter fit to the personal workflow
- easier to audit and evolve
- avoids importing broad product complexity that does not serve the real use case

Success criteria:

- one command can turn a messy page or local document set into reviewable markdown
- output lands in Obsidian with provenance and stable filenames
- agents can discover the capability through AthenaUse without loading it into every session
- accepted artifacts can be promoted into AthenaMind intentionally, not automatically

## Current Implementation Notes

The current implementation lives in:

- `products/athena-use/cmd/intake-cli/`
- `products/athena-use/internal/intake/`

Current posture:

- local-first
- no background service
- no browser automation
- main-content HTML-to-markdown extraction with absolute discovered links when a base URL is known
- best-effort PDF text extraction with page grouping
- explicit staging into `Athena/01 Inbox` by default
- `ATHENA_VAULT` env support for vault targeting
- stable staged filenames derived from artifact creation time and title
- staged artifacts include a light review scaffold for:
  - highlights
  - keep/discard decisions
  - promotion notes for later AthenaMind import

## Follow-On Note

GitNexus may deserve the same treatment later:

- keep the useful structural-analysis value
- challenge the surrounding runtime surface
- decide whether a thinner local contract would fit AthenaPlatform better than adopting the full external product posture
