# Internal Beta

This is the simplest supported starting path for internal Cairn testers.

## Who This Is For

- engineers comfortable with local tooling
- AI champions trying the platform in real repo work
- operators validating memory-cli and work harness together

## Start Here

1. Read `PLATFORM_QUICKSTART.md`.
2. Use repo-local runtime state under `.cairn/`.
3. Run the supported smoke path:
   - sqlite-first: `./tools/platform/smoke_v1.sh`
4. For work harness human workflow, read `workspace/docs/HUMANS.md`.
5. For memory-cli operator setup, read `docs/operator/memory-cli/getting-started/README.md`.

## Recommended Local Setup

```bash
export CAIRN_HOME="$PWD/.cairn"
export CAIRN_MEMORY_ROOT="$CAIRN_HOME/memory/default"
mkdir -p "$CAIRN_MEMORY_ROOT"
```

Optional embeddings:

```bash
export CAIRN_EMBEDDING_ENDPOINT="http://192.168.1.35:11434"
export CAIRN_OLLAMA_EMBED_MODEL="mxbai-embed-large:latest"
```

## Active Platform Surface

- product code: `products/`
- shared tooling: `tools/`
- workspace contract: `workspace/`
- operator and product docs: `docs/`
- local runtime state: `.cairn/`

## Historical Material

The beta path intentionally omits historical integration material.
