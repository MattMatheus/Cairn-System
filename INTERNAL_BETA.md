# Internal Beta

This is the simplest supported starting path for internal AthenaPlatform testers.

## Who This Is For

- engineers comfortable with local tooling
- AI champions trying the platform in real repo work
- operators validating AthenaMind and AthenaWork together

## Start Here

1. Read `PLATFORM_QUICKSTART.md`.
2. Use repo-local runtime state under `.athena/`.
3. Run the supported smoke path:
   - sqlite-first: `./tools/platform/smoke_v1.sh`
4. For AthenaWork human workflow, read `workspace/docs/HUMANS.md`.
5. For AthenaMind operator setup, read `docs/operator/athena-mind/getting-started/README.md`.

## Recommended Local Setup

```bash
export ATHENA_HOME="$PWD/.athena"
export ATHENA_MEMORY_ROOT="$ATHENA_HOME/memory/default"
mkdir -p "$ATHENA_MEMORY_ROOT"
```

Optional embeddings:

```bash
export ATHENA_EMBEDDING_ENDPOINT="http://192.168.1.35:11434"
export ATHENA_OLLAMA_EMBED_MODEL="mxbai-embed-large:latest"
```

## Active Platform Surface

- product code: `products/`
- shared tooling: `tools/`
- workspace contract: `workspace/`
- operator and product docs: `docs/`
- local runtime state: `.athena/`

## Historical Material

The beta path intentionally omits historical integration material.
