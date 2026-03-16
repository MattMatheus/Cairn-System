# tool-cli

tool-cli is the tool-context companion to memory-cli.

It gives work harness a governed, scoped tool surface at stage launch so agents receive:

1. seed prompt
2. memory bootstrap
3. tool context

## V1 Role

tool-cli v1 is a discovery and context-emission product.

Its job is to let a small prompt or `SKILL.md` explain the available toolset without loading unnecessary tools into every session.

Initial command surface:

- `tool-cli discover`
- `tool-cli context`
- `tool-cli inspect`
- `tool-cli list`
- `tool-cli validate`
- `intake-cli inspect|url|file|folder|stage`
- `promote-cli inspect|note`

Deferred:

- `tool-cli call`
- memory-backed registry mode
- bootstrap and artifact retrieval concerns

Operating rule:

- tools are called only when needed
- tool availability should remain understandable from a compact context payload
- tool-cli should reduce ambient tool noise, not increase it

Approved tools can now also declare an availability posture:

- `required`
- `default`
- `scoped`

That posture determines whether a tool belongs in normal context or should remain discoverable and inspectable until a session explicitly needs it.

Approved tools can also declare status:

- `active`
- `planned`

Planned tools are part of the planning surface and `inspect` output, but they stay out of active emitted context unless explicitly included.

Current active companion binary:

- `intake-cli` for local-first URL, file, folder, and PDF normalization into Cairn inbox artifacts
- `promote-cli` for deliberate promotion of reviewed Cairn notes into memory-cli

Current intake ergonomics:

- `--title` can override the staged artifact title
- `CAIRN_VAULT` can provide the default vault path
- HTML normalization prefers `main`/`article` content and resolves discovered links against the source URL when possible
- staged intake artifacts include a compact review scaffold for triage and later memory-cli promotion notes

Current active tool-cli intake entries:

- `cairn.intake.inspect_source`
- `cairn.intake.normalize_source`
- `cairn.memory.promote_note`

## Registry Contract

V1 uses a config-backed registry with two trust tiers:

- approved: repo-backed and supported
- local: operator-managed and opt-in

Approved tools should live under:

- `products/tool-cli/registry/approved-tools.yaml`

Local overlays are expected under repo-local runtime state:

- `.cairn/tools/registry.yaml`

## Observability

tool-cli should follow the memory-cli telemetry posture:

- OpenTelemetry is the tracing and metrics standard
- tool discovery, context emission, validation, and later execution paths must emit traceable spans
- no separate observability framework should be introduced

## References

- `docs/product/tooling/TOOLCLI_V1.md`
- `docs/product/tooling/MODEL_TOOL_INTERFACE_PREP.md`
- `docs/decisions/ADR-0005-toolcli-v1-trust-and-registry-policy.md`
- `docs/decisions/ADR-0006-toolcli-telemetry-and-dependency-policy.md`
