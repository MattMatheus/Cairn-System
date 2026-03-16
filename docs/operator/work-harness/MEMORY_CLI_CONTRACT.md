# memory-cli CLI Contract For work harness

## Purpose

This document defines the initial integration contract between work harness and memory-cli inside Cairn.

The contract is intentionally simple:

- work harness owns workflow and workspace structure
- memory-cli owns local memory and retrieval
- integration happens through the memory-cli CLI

work harness now also composes tool context through tool-cli at stage launch.

Recommended repo-local runtime defaults:

- `CAIRN_HOME="$PWD/.cairn"`
- `CAIRN_MEMORY_ROOT="$CAIRN_HOME/memory/default"`

## Default Contract

work harness should treat memory-cli as a local sidecar tool invoked through:

```bash
(cd products/memory-cli && go run ./cmd/memory-cli <command> ...)
```

For tool context, work harness should treat tool-cli as a local companion invoked through:

```bash
(cd products/tool-cli && go run ./cmd/tool-cli context --stage <stage>)
```

## Preferred V1 Commands

- `write`
- `retrieve`
- `bootstrap`
- `verify`
- `snapshot`

## Expected Inputs

Typical work harness-originating inputs:

- queue item context
- workflow stage
- domain or project label
- markdown content or note excerpts
- cycle identifier when relevant

## Expected Outputs

work harness should expect memory-cli to return:

- selected memory ids
- selection mode
- source path
- confidence or warning data
- verification/health output for local memory readiness

work harness should expect tool-cli to return:

- approved tool ids
- short descriptions
- parameter summaries
- credential references by name only
- support tier markers

## Example Use Cases

### 1. Queue Bootstrap

Before implementing a task, retrieve relevant memory:

```bash
(cd products/memory-cli && go run ./cmd/memory-cli retrieve \
  --root "$CAIRN_MEMORY_ROOT" \
  --query "queue metadata observer cycle handoff" \
  --domain platform)
```

### 2. Planning Bootstrap

Create reusable planning guidance:

```bash
(cd products/memory-cli && go run ./cmd/memory-cli write \
  --root "$CAIRN_MEMORY_ROOT" \
  --id planning-bootstrap \
  --title "Planning Bootstrap" \
  --type prompt \
  --domain platform \
  --body "Create explicit next-stage recommendations and preserve rationale in workspace notes." \
  --stage planning \
  --reviewer matt \
  --decision approved \
  --reason "bootstrap planning memory" \
  --risk "low" \
  --notes "example")
```

### 3. Health Check Before Use

```bash
(cd products/memory-cli && go run ./cmd/memory-cli verify health \
  --root "$CAIRN_MEMORY_ROOT" \
  --query "planning bootstrap")
```

## V1 Constraints

1. Default storage is local and sqlite-first.
2. work harness should assume memory-cli is sqlite-only in Cairn unless the product boundary is intentionally expanded again.
3. work harness should not depend on memory-cli-specific internal packages.
4. work harness should not depend on tool-cli-specific internal packages.
5. The CLI boundary is the source of truth for product integration until a stronger service boundary is intentionally introduced.
