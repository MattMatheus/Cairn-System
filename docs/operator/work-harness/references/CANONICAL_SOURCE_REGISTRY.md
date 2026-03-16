---
title: Canonical Source Registry
description: Specification for recording official vendor docs, allowed domains, and version assumptions so work harness agents use the right external sources.
doc_type: reference
status: draft
---

# Canonical Source Registry

## Summary

work harness needs explicit records for official external docs. A canonical source record tells an agent which vendor docs to trust, which version to assume, and when to consult those docs instead of guessing or using stale internal notes.

## Intended Audience

- maintainers authoring work harness operator/reference content
- agent and prompt authors
- humans pinning toolchains and framework versions

## Preconditions

- the task touches an external tool, framework, database, API, or standard
- the team can name an official vendor documentation domain
- version expectations are known or can be declared as intentionally floating

## Main Flow

For any external tool that matters to development work:

1. create a canonical source record
2. pin a version or version policy
3. list official domains and base docs URLs
4. define when agents must consult that source
5. attach the record to relevant graph nodes

This lets work harness combine internal usage guidance with vendor truth.

## Record Schema

Each canonical source record should be one markdown file with frontmatter like this:

```yaml
title: Terraform 1.9
description: Official Terraform docs for work harness tasks that touch Terraform CLI or HCL semantics.
doc_type: canonical_source
status: active
tool: terraform
vendor: HashiCorp
version_policy: pinned_minor
version: "1.9"
docs_base: https://developer.hashicorp.com/terraform
allowed_domains:
  - developer.hashicorp.com
entrypoints:
  - https://developer.hashicorp.com/terraform/docs
  - https://developer.hashicorp.com/terraform/language
  - https://developer.hashicorp.com/terraform/cli
consult_when:
  - interpreting Terraform language behavior
  - using Terraform CLI flags or state commands
  - checking version-specific compatibility
avoid_when:
  - the task is only about work harness process or local conventions
notes:
  - Prefer official docs over blog posts or forum answers.
```

## Version Policy

Use one of these version policies:

- `pinned_exact`: use an exact version such as `1.9.8`
- `pinned_minor`: use a minor line such as `1.9`
- `pinned_major`: use a major line such as `19.x`
- `floating_current`: use current official docs because versioned docs are not stable or not separately published
- `repo_pinned`: infer the version from local project files first, then use the matching official docs

Default rule:

- use `repo_pinned` when the repository declares a version
- otherwise use the narrowest version pin the team can actually maintain
- use `floating_current` only when the tool does not offer clean versioned docs or the team explicitly accepts drift

## Retrieval Policy

Agents should use canonical sources when a task depends on:

- CLI flags or command behavior
- configuration syntax
- framework APIs
- database semantics
- version compatibility
- deployment or runtime requirements

Agents should prefer internal graph nodes first when a task is about:

- work harness process
- local architecture decisions
- repo conventions
- workflow heuristics

If internal docs conflict with vendor docs, the agent should report the conflict and prefer the official source for external semantics.

## Folder Layout

Canonical source records should live under:

- `docs/operator/work-harness/references/canonical-sources/`

Graph nodes should reference them in frontmatter:

```yaml
canonical_sources:
  - terraform-1-9
  - postgresql-16
```

Or in prose when the relationship matters:

- `work harness infra changes depend on [[terraform-1-9]].`

## Failure Modes

- unversioned links cause drift across machines and over time
- unofficial domains creep in and dilute trust
- source records duplicate internal conventions instead of vendor truth
- "current docs" are used where the repo actually pins an older version
- graph nodes quote vendor docs instead of linking to the canonical source record

## References

- [[DEVELOPMENT_COGNITION_GRAPH]]
- [[AUTHORING_GRAPH_NODES]]
- [canonical-sources/README.md](/home/matt/Workspace/repos/trusted/Cairn/docs/operator/work-harness/references/canonical-sources/README.md)
