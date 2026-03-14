# AthenaUse V1

Status: proposed product design for implementation planning.

## Purpose

AthenaUse is the tool-surface companion to AthenaMind.

AthenaWork already launches agents with:

1. a seed prompt
2. a memory bootstrap

AthenaUse adds the third governed input:

3. a scoped tool context

The goal is to give agents a small, relevant, machine-readable set of approved tools instead of either no tool guidance or a bloated manifest dump.

AthenaUse should let a simple prompt or `SKILL.md` explain available tool options without forcing those tools into every startup context.

## Product Role

- AthenaMind retrieves knowledge context
- AthenaUse retrieves tool context
- AthenaWork composes stage prompt, memory context, and tool context

AthenaUse is not a daemon, sandbox, MCP broker, OAuth manager, or secret store.

It is also not an excuse to ambiently load every possible tool into every thread.

## V1 Scope

V1 is discovery and context emission, not a full execution runtime.

Supported commands:

- `use-cli discover <query>`
- `use-cli context [--stage <stage>] [--query <query>]`
- `use-cli inspect <tool-id>`
- `use-cli list [--stage <stage>] [--tag <tag>]`
- `use-cli validate`

Deferred from v1:

- `use-cli call`
- memory-backed tool registry
- Azure/bootstrap artifact retrieval
- full JSON Schema support

Operational rule:

- tools should be discoverable and understandable from a compact context payload
- tools should only be injected into a session when the work requires them

## Observability And Dependency Policy

AthenaUse should align with AthenaMind's tracing posture.

Requirements:

- OpenTelemetry is the observability standard
- discovery, context emission, validation, and future execution paths should emit traceable spans
- tool-selection behavior should remain auditable
- no separate observability framework should be introduced

Dependency rule:

- for AthenaUse-specific implementation, OpenTelemetry should be the only non-stdlib dependency family introduced in v1

## Naming

- product concept: `AthenaUse`
- binary name: `use-cli`

This matches `memory-cli` and keeps the executable name short and functional.

## Registry Model

AthenaUse uses a curated registry of callable tools.

Each tool entry describes:

- stable identifier
- human-readable name
- description for intent matching
- implementation status
- availability posture
- operator guidance
- complementary tools or systems when relevant
- tags
- stage affinity
- credential reference
- call contract
- parameter schema

## Trust Tiers

AthenaUse supports two trust tiers.

### Athena Approved

Repo-backed, curated, and supported.

Properties:

- committed to the repository
- reviewed and versioned with platform changes
- validated by `use-cli validate`
- eligible for AthenaWork stage-context injection according to tool availability posture
- documented as supported

### Caveat Emptor

Local or operator-defined tools outside the supported platform contract.

Properties:

- user-managed
- optional and additive
- never assumed to be stable
- never injected by default into AthenaWork stage context
- discoverable only when explicitly requested

## Registry Locations

V1 should support layered sources with explicit precedence.

1. env override: `$ATHENA_USE_REGISTRY`
2. repo-approved registry: committed platform registry
3. local registry: user-local or repo-local runtime overlay

Recommended defaults:

- approved registry: repo-backed under the committed platform tree
- recommended file: `products/athena-use/registry/approved-tools.yaml`
- local registry: `.athena/tools/registry.yaml`

Operational rule:

- default command behavior uses approved tools only
- local tools require explicit inclusion

## Default Behavior

`use-cli context` should emit approved tools only unless the operator explicitly opts in to local overlays.

Suggested flag model:

- default: approved only
- `--include-local`: merge local tools into results
- `--include-scoped`: include approved scoped tools even when no query is provided

All output should mark support tier:

- `support_tier: approved`
- `support_tier: local`

Approved tools should also carry availability posture:

- `required`: part of the normal baseline when stage-matched
- `default`: normal stage-matched approved tool
- `scoped`: approved but excluded from default context unless explicitly relevant

Approved tools may also carry implementation status:

- `active`: implemented and eligible for active use
- `planned`: part of the planning surface, inspectable and discoverable, but excluded from active context unless explicitly included for planning work

## Registry Shape

V1 should use a minimal structured YAML format.

Example:

```yaml
version: 1
tools:
  - id: github.create_pr
    name: Create Pull Request
    description: Opens a pull request against a target branch in a GitHub repository
    tags: [github, vcs, code-review]
    stage_affinity: [engineering, qa]
    credential: GITHUB_TOKEN
    call:
      type: http
      method: POST
      url: "https://api.github.com/repos/{owner}/{repo}/pulls"
    schema:
      - name: owner
        type: string
        required: true
      - name: repo
        type: string
        required: true
      - name: title
        type: string
        required: true
      - name: head
        type: string
        required: true
      - name: base
        type: string
        required: true
      - name: body
        type: string
        required: false
```

V1 schema fields should support:

- `name`
- `type`
- `required`
- `enum`
- `description`

This keeps authoring simple while remaining machine-readable.

## Discovery Model

Key design principle: intent-based discovery over manifest dumping.

Agents should ask for tools by what they want to do, not by already knowing the tool name.

`discover` should return:

- `id`
- `name`
- `description`
- `status`
- `availability`
- `guidance`
- `credential_ref`
- summarized schema
- support tier

Formats:

- human-readable default
- `--format=json`
- `--format=yaml`

## Context Model

`context` is the stage-launch integration surface.

It should filter by:

- `stage_affinity`
- query intent when provided
- support tier rules

Output should stay intentionally small.

It should help an agent understand what is available, not bury the agent under a random execution manifest.

Recommended v1 fields:

- `id`
- `name`
- `description`
- `status`
- `availability`
- `guidance`
- `schema`
- `credential_ref`
- `call_type`
- `support_tier`

It should not emit full secret material or large transport templates by default.

If a scoped tool is not explicitly relevant, it should stay out of default emitted context.

If a tool is still `planned`, it should also stay out of default emitted context even when it is discoverable.

## Inspect Model

`inspect` is the operator and PM-facing fit-check surface.

It should answer:

- what the tool is for
- whether it is active or planned
- whether it is `required`, `default`, or `scoped`
- whether it matches the current stage
- whether it belongs in default context
- what it complements in the wider Athena system

This keeps tool reasoning compact and auditable without reopening runtime execution scope.

## Next Approved Slice

After approved-registry expansion and shared-platform validation enforcement, the next approved AthenaUse implementation slice is:

- stronger `context` schema shaping within the existing v1 boundary

Deferred beyond that slice:

- formal model/tool interface specification
- execution/runtime contract work
- full JSON Schema support

## Credential Handling

AthenaUse never stores or emits credential values.

Registry entries contain only credential references such as env var names.

At v1 scope:

- context and discovery emit reference names only
- validation checks whether required references are present when appropriate
- secret resolution belongs to execution-time tooling, not context generation

## Validation

`use-cli validate` should enforce the approved registry.

Minimum checks:

- valid YAML structure
- supported version
- duplicate ID detection
- required field presence
- schema field integrity
- stage-affinity value validity
- credential reference presence checks when requested

Because AthenaUse should stay dependency-light, v1 validation and config handling should prefer simple native parsing and strict contract checks over broad framework adoption.

Validation of local tools should be advisory, not equivalent to approved-platform validation.

## AthenaWork Integration

AthenaWork should add `emit_tool_context` alongside memory bootstrap emission.

Behavior:

- call `use-cli context --stage <stage>`
- inject approved tool context after memory bootstrap
- warn and continue if AthenaUse is unavailable
- do not fail stage launch solely because tool context is unavailable in v1
- instrument the context emission path with OpenTelemetry spans

## Non-Goals

- general plugin runtime
- policy engine for arbitrary tool brokering
- long-running tool server
- automatic credential brokering
- replacing AthenaMind retrieval

## Open V2 Questions

- should `use-cli call` become an operator-safe execution path?
- should memory-backed registry mode exist in addition to config-backed mode?
- should workspace policy overlays refine `stage_affinity`?
- should schema evolve toward JSON Schema or OpenAPI fragments?
- should approved registries support signed or checksum-validated bundles later?
