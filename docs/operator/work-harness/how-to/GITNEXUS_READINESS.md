# GitNexus Readiness

Use GitNexus through Cairn's bounded `codegraph-cli` wrapper, not through raw GitNexus setup surfaces.

Source of truth:

- `products/tool-cli/cmd/codegraph-cli`
- `products/tool-cli/internal/codegraph`
- `products/work-harness/tools/check_gitnexus_readiness.sh`

## Goal

Make the GitNexus backend available without introducing ambient MCP setup, hooks, or global editor mutation.

## Supported readiness modes

1. Preferred: `CAIRN_GITNEXUS_BIN`
2. Acceptable: `gitnexus` on `PATH`
3. Controlled fallback: built local checkout plus `node`

## Preflight

Run:

```bash
./products/work-harness/tools/check_gitnexus_readiness.sh
```

If Cairn should inspect a non-default checkout root, run:

```bash
./products/work-harness/tools/check_gitnexus_readiness.sh --root /home/matt/Workspace
```

## Expected local checkout

Default lookup path:

```text
/home/matt/Workspace/repos/untrusted/GitNexus/gitnexus
```

For the local checkout fallback to work, Cairn expects:

- `node` on `PATH`
- built entry point at `dist/cli/index.js`

## Recommended operator posture

- Prefer a pinned local binary or pinned built checkout.
- Do not use `npx gitnexus@latest` as the normal Cairn runtime path.
- Do not enable GitNexus hooks, MCP setup, or generated context files as part of Cairn bootstrap.

## Runtime path

Discovery and fit:

```bash
cd products/tool-cli && go run ./cmd/tool-cli inspect gitnexus
```

Execution:

```bash
cd products/tool-cli && go run ./cmd/codegraph-cli status --repo /path/to/repo
cd products/tool-cli && go run ./cmd/codegraph-cli analyze --repo /path/to/repo
cd products/tool-cli && go run ./cmd/codegraph-cli context --repo /path/to/repo SymbolName
cd products/tool-cli && go run ./cmd/codegraph-cli impact --repo /path/to/repo --direction upstream SymbolName
```
