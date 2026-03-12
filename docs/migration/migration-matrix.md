# AthenaPlatform Migration Matrix

## Purpose

This matrix maps the current source trees into the canonical AthenaPlatform scaffold.

Decision states:

- `keep`: preserve in source for now; do not import yet
- `merge`: union into the canonical destination
- `trim`: import selectively after removing non-v1 or non-product content
- `archive`: preserve for history but keep out of the primary platform path
- `exclude`: do not bring into the platform tree

## Canonical Destinations

- Product code: `products/`
- Shared tools: `tools/`
- Workspace content: `workspace/`
- Platform docs: `docs/`

## AthenaMind

| Source | Decision | Destination | Notes |
|---|---|---|---|
| `AthenaMind/.git` | exclude | n/a | Embedded repo metadata stays local only. |
| `AthenaMind/.github/` | archive | `docs/migration/archive/athenamind-github/` | Preserve CI history; evaluate later for shared CI. |
| `AthenaMind/cmd/memory-cli/` | merge | `products/athena-mind/cmd/memory-cli/` | Canonical CLI surface for AthenaMind v1. |
| `AthenaMind/internal/index/` | merge | `products/athena-mind/internal/index/` | Core local indexing path; keep. |
| `AthenaMind/internal/retrieval/` | trim | `products/athena-mind/internal/retrieval/` | Keep sqlite-first retrieval; defer advanced backend complexity not needed for v1. |
| `AthenaMind/internal/snapshot/` | merge | `products/athena-mind/internal/snapshot/` | Keep if lightweight and useful for local recovery. |
| `AthenaMind/internal/gateway/` | trim | `products/athena-mind/internal/gateway/` | Optional for v1; keep only if it supports local CLI workflows cleanly. |
| `AthenaMind/internal/telemetry/` | trim | `products/athena-mind/internal/telemetry/` | Keep minimal health/telemetry only. |
| `AthenaMind/internal/episode/` | archive | `docs/migration/archive/athenamind-episode/` | Not needed for initial slim RAG target. |
| `AthenaMind/internal/governance/` | archive | `docs/migration/archive/athenamind-governance/` | Advanced governance is out of v1 scope. |
| `AthenaMind/internal/types/` | merge | `products/athena-mind/internal/types/` | Shared core types. |
| `AthenaMind/memory/core/` | trim | `products/athena-mind/examples/memory-core/` | Keep only sanitized sample data/layout, not current runtime state. |
| `AthenaMind/knowledge-base/cli/` | merge | `docs/product/athena-mind/cli/` | Product docs for the slim CLI. |
| `AthenaMind/knowledge-base/getting-started/` | merge | `docs/operator/athena-mind/` | Good base for operator setup. |
| `AthenaMind/knowledge-base/how-to/` | trim | `docs/operator/athena-mind/how-to/` | Keep only v1-relevant docs. |
| `AthenaMind/knowledge-base/product/` | merge | `docs/product/athena-mind/` | Product docs and intent. |
| `AthenaMind/knowledge-base/workflows/` | trim | `docs/operator/athena-mind/workflows/` | Keep workflows that match the slim product. |
| `AthenaMind/knowledge-base/process/` | archive | `docs/migration/archive/athenamind-process/` | Process docs overlap with AthenaWork ownership. |
| `AthenaMind/knowledge-base/references/` | trim | `docs/product/athena-mind/references/` | Keep backend/config references relevant to v1. |
| `AthenaMind/knowledge-base/faq/` | trim | `docs/operator/athena-mind/faq/` | Keep only still-valid answers. |
| `AthenaMind/knowledge-base/troubleshooting/` | merge | `docs/operator/athena-mind/troubleshooting/` | Useful for local developer use. |
| `AthenaMind/knowledge-base/release-notes/` | archive | `docs/migration/archive/athenamind-release-notes/` | Historical. |
| `AthenaMind/knowledge-base/public-testing/` | archive | `docs/migration/archive/athenamind-public-testing/` | Not needed for current platform bootstrap. |
| `AthenaMind/products/athena-mind/` | merge | `products/athena-mind/` | Use as product import source if cleaner than root paths. |
| `AthenaMind/products/athena-work/` | archive | `docs/migration/archive/athenamind-athenawork-copy/` | Avoid duplicate AthenaWork import path. |
| `AthenaMind/skills/athena-mind/` | merge | `products/athena-mind/skills/athena-mind/` | Keep product skill. |
| `AthenaMind/skills/athena-work/` | archive | `docs/migration/archive/athenamind-athenawork-skill/` | Use canonical AthenaWork skill from AthenaWork source. |
| `AthenaMind/website/` | archive | `docs/migration/archive/athenamind-website/` | Defer marketing/site concerns. |
| `AthenaMind/*.md` | trim | `docs/product/athena-mind/` | Keep only relevant root docs such as README and selected support docs. |
| `AthenaMind/go.mod` and `AthenaMind/go.sum` | merge | `products/athena-mind/` | Needed for product module import. |
| `AthenaMind/docker-compose*.yml` | trim | `tools/dev/athena-mind/` | Keep only files that support local optional services. |

## AthenaWork Productized Variant

| Source | Decision | Destination | Notes |
|---|---|---|---|
| `AthenaWork/.git` | exclude | n/a | Embedded repo metadata stays local only. |
| `AthenaWork/.github/` | archive | `docs/migration/archive/athenawork-github/` | Preserve CI history; adapt later if useful. |
| `AthenaWork/products/athena-work/` | merge | `products/athena-work/` | Canonical executable AthenaWork product surface. |
| `AthenaWork/skills/athena-work/` | merge | `products/athena-work/skills/athena-work/` | Canonical AthenaWork skill. |
| `AthenaWork/knowledge-base/process/` | merge | `docs/operator/athena-work/process/` | Good source for workflow/operator docs. |
| `AthenaWork/knowledge-base/product/` | merge | `docs/product/athena-work/` | Product docs. |
| `AthenaWork/knowledge-base/operations/` | merge | `docs/operator/athena-work/operations/` | Operational runbooks. |
| `AthenaWork/knowledge-base/how-to/` | merge | `docs/operator/athena-work/how-to/` | Practical guidance. |
| `AthenaWork/knowledge-base/architecture/` | trim | `docs/product/athena-work/architecture/` | Keep canonical architecture docs. |
| `AthenaWork/knowledge-base/references/` | trim | `docs/product/athena-work/references/` | Keep pinned references that still matter. |
| `AthenaWork/knowledge-base/decisions/` | archive | `docs/migration/archive/athenawork-decisions/` | Review before promoting to platform decisions. |
| `AthenaWork/product-research/roadmap/` | trim | `docs/migration/athenawork-roadmap/` | Keep decision-driving items; archive project-specific items. |
| `AthenaWork/product-research/planning/` | archive | `docs/migration/archive/athenawork-planning/` | Historical planning context. |
| `AthenaWork/website/` | archive | `docs/migration/archive/athenawork-website/` | Defer site concerns. |
| `AthenaWork/.env.example` | trim | `tools/dev/examples/athenawork.env.example` | Keep only if still useful. |
| `AthenaWork/docker-compose*.yml` | trim | `tools/dev/athena-work/` | Keep local dev support only if still relevant. |
| `AthenaWork/azure-pipelines.yml` | archive | `docs/migration/archive/athenawork-ci/` | Historical CI config. |
| `AthenaWork/go.mod` | trim | `products/athena-work/` | Keep only if AthenaWork product module remains Go-backed. |
| `AthenaWork/HUMANS.md` | merge | `docs/operator/athena-work/HUMANS.md` | Productized operator instructions. |
| `AthenaWork/AGENTS.md` | merge | `products/athena-work/AGENTS.md` | Product agent contract. |
| `AthenaWork/README.md` | merge | `docs/product/athena-work/README.md` | Product overview. |
| `AthenaWork/DEVELOPMENT_CYCLE.md` | merge | `docs/operator/athena-work/DEVELOPMENT_CYCLE.md` | Core stage workflow reference. |

## AthenaWork Vault Variant

| Source | Decision | Destination | Notes |
|---|---|---|---|
| `work/AthenaWork/.git` | exclude | n/a | Embedded repo metadata stays local only. |
| `work/AthenaWork/.obsidian/` | archive | `docs/migration/archive/athenawork-obsidian/` | Preserve only if settings/templates are useful. |
| `work/AthenaWork/00-Docs/` | merge | `workspace/docs/` | Strong source for workspace-native docs. |
| `work/AthenaWork/10 Work/` | trim | `workspace/work/` | Merge structure first; classify and archive project-specific/private content later. |
| `work/AthenaWork/30 AI/10 Policies/` | merge | `workspace/docs/policies/` | Good workspace policy source. |
| `work/AthenaWork/30 AI/20 Prompts/` | merge | `products/athena-work/prompts/` | Product prompts belong with AthenaWork. |
| `work/AthenaWork/30 AI/30 Workflows/` | merge | `workspace/docs/workflows/` | Workspace execution patterns. |
| `work/AthenaWork/30 AI/40 Operating System/` | merge | `products/athena-work/operating-system-vault/` | Merge against productized AthenaWork contracts. |
| `work/AthenaWork/30 AI/50 References/` | trim | `docs/product/athena-work/references-vault/` | Keep selected references. |
| `work/AthenaWork/40 Assets/Templates/` | merge | `workspace/templates/` | Canonical source for markdown templates. |
| `work/AthenaWork/40 Assets/*` | trim | `workspace/assets/` | Keep reusable assets only. |
| `work/AthenaWork/70 Agents/00 Queue/` | merge | `workspace/agents/queue/` | Canonical queue structure source. |
| `work/AthenaWork/70 Agents/05 Delivery Queue/` | merge | `workspace/agents/delivery/` | Unify with productized stage flow. |
| `work/AthenaWork/70 Agents/20 Memory/` | merge | `workspace/agents/memory/` | Important integration point with AthenaMind. |
| `work/AthenaWork/70 Agents/40 Policies/` | merge | `workspace/agents/policies/` | Useful workspace governance. |
| `work/AthenaWork/70 Agents/60 Maps/` | merge | `workspace/agents/maps/` | Navigation/control surfaces. |
| `work/AthenaWork/70 Agents/70 Escape/` | merge | `workspace/agents/escape/` | Good blocked-work handling model. |
| `work/AthenaWork/70 Agents/90 Staging/` | merge | `workspace/agents/staging/` | Useful staging area. |
| `work/AthenaWork/70 Agents/system.json` | merge | `workspace/agents/system.json` | Workspace control contract. |
| `work/AthenaWork/80 Research/` | trim | `workspace/research/` | Keep structure and product-safe exemplars; archive live/project-specific research. |
| `work/AthenaWork/tools/prune_user_content.sh` | merge | `tools/migration/prune_user_content.sh` | Useful cleanup utility. |
| `work/AthenaWork/tools/validate_task_metadata.sh` | merge | `tools/platform/validate_task_metadata.sh` | Useful shared validation utility. |
| `work/AthenaWork/HUMANS.md` | merge | `workspace/docs/HUMANS.md` | Workspace-native operator landing page. |
| `work/AthenaWork/AGENTS.md` | merge | `workspace/AGENTS.md` | Workspace-native agent contract. |
| `work/AthenaWork/README.md` | merge | `workspace/README.md` | Workspace overview. |
| `work/AthenaWork/EXTRACTION.md` | archive | `docs/migration/archive/athenawork-extraction/` | Historical/specialized. |
| `work/AthenaWork/Index.md` | trim | `workspace/docs/Index.md` | Keep if it remains the best workspace map. |

## Root-Level Merge Rules

1. Do not import `.git/` directories.
2. Do not import generated runtime data, WAL files, or local indexes as product defaults.
3. Prefer union merges over replacement for docs, prompts, templates, and contracts.
4. When two files define the same concept differently, import both into adjacent destinations and reconcile later.
5. Preserve product instructions and skills even if duplicated elsewhere.
6. Treat live work notes, task output, and project-specific research as `trim` until explicitly classified.

## First Execution Pass

The first migration pass should focus on:

1. `products/athena-mind/` slim import
2. `products/athena-work/` productized workflow import
3. `workspace/` structure and templates import
4. shared tool promotion into `tools/`
5. documentation promotion into `docs/`

