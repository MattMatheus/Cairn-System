# Embedding Endpoint Note

## Verified Endpoint Pattern

AthenaMind currently works with an Ollama-compatible base endpoint.

Validated host pattern:

- base endpoint: `http://192.168.1.35:11434`
- API path used by AthenaMind: `/api/embeddings`

## Verified Model

- `mxbai-embed-large:latest`

## Model Override

AthenaMind now supports an Ollama model override through:

- `ATHENA_OLLAMA_EMBED_MODEL`

Example:

```bash
export ATHENA_OLLAMA_EMBED_MODEL="mxbai-embed-large:latest"
```

## Important Detail

For AthenaMind, pass the base host URL as `--embedding-endpoint`.

Do not pass `/api/embed` directly for the current CLI path.

Example:

```bash
export ATHENA_HOME="$PWD/.athena"
export ATHENA_MEMORY_ROOT="${ATHENA_MEMORY_ROOT:-$ATHENA_HOME/memory/default}"
mkdir -p "$ATHENA_MEMORY_ROOT"

(cd products/athena-mind && ATHENA_OLLAMA_EMBED_MODEL="mxbai-embed-large:latest" go run ./cmd/memory-cli verify health \
  --root "$ATHENA_MEMORY_ROOT" \
  --query "planning bootstrap" \
  --embedding-endpoint http://192.168.1.35:11434)
```

## Validation Note

AthenaMind applies a latency degradation policy during retrieval.

The default threshold is tuned for local iteration rather than strict low-latency production assumptions:

- default fallback threshold: `1500ms`
- set `MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS=0` to disable latency fallback entirely during local testing

The default threshold is `700ms`. When you are validating semantic retrieval against a real embedding service, that policy can force deterministic fallback even when the vector path is working correctly.

For semantic validation runs, disable the fallback gate explicitly:

```bash
export ATHENA_HOME="$PWD/.athena"
export ATHENA_MEMORY_ROOT="${ATHENA_MEMORY_ROOT:-$ATHENA_HOME/memory/default}"
mkdir -p "$ATHENA_MEMORY_ROOT"

(cd products/athena-mind && \
  MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS=0 \
  ATHENA_OLLAMA_EMBED_MODEL="mxbai-embed-large:latest" \
  go run ./cmd/memory-cli verify health \
    --root "$ATHENA_MEMORY_ROOT" \
    --query "planning bootstrap" \
    --embedding-endpoint http://192.168.1.35:11434)
```

## Notes

- The platform smoke path works without embeddings via lexical fallback.
- A reachable embedding endpoint improves retrieval quality validation.
- The platform smoke script disables latency fallback during semantic validation when an embedding endpoint is configured.
