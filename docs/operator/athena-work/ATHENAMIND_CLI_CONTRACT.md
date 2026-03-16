# AthenaMind CLI Contract For AthenaWork

## Purpose

This document defines the initial integration contract between AthenaWork and AthenaMind inside Cairn.

The contract is intentionally simple:

- AthenaWork owns workflow and workspace structure
- AthenaMind owns local memory and retrieval
- integration happens through the AthenaMind CLI

AthenaWork now also composes tool context through AthenaUse at stage launch.

Recommended repo-local runtime defaults:

- `ATHENA_HOME="$PWD/.athena"`
- `ATHENA_MEMORY_ROOT="$ATHENA_HOME/memory/default"`

## Default Contract

AthenaWork should treat AthenaMind as a local sidecar tool invoked through:

```bash
(cd products/athena-mind && go run ./cmd/memory-cli <command> ...)
```

For tool context, AthenaWork should treat AthenaUse as a local companion invoked through:

```bash
(cd products/athena-use && go run ./cmd/use-cli context --stage <stage>)
```

## Preferred V1 Commands

- `write`
- `retrieve`
- `bootstrap`
- `verify`
- `snapshot`

## Expected Inputs

Typical AthenaWork-originating inputs:

- queue item context
- workflow stage
- domain or project label
- markdown content or note excerpts
- cycle identifier when relevant

## Expected Outputs

AthenaWork should expect AthenaMind to return:

- selected memory ids
- selection mode
- source path
- confidence or warning data
- verification/health output for local memory readiness

AthenaWork should expect AthenaUse to return:

- approved tool ids
- short descriptions
- parameter summaries
- credential references by name only
- support tier markers

## Example Use Cases

### 1. Queue Bootstrap

Before implementing a task, retrieve relevant memory:

```bash
(cd products/athena-mind && go run ./cmd/memory-cli retrieve \
  --root "$ATHENA_MEMORY_ROOT" \
  --query "queue metadata observer cycle handoff" \
  --domain platform)
```

### 2. Planning Bootstrap

Create reusable planning guidance:

```bash
(cd products/athena-mind && go run ./cmd/memory-cli write \
  --root "$ATHENA_MEMORY_ROOT" \
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
(cd products/athena-mind && go run ./cmd/memory-cli verify health \
  --root "$ATHENA_MEMORY_ROOT" \
  --query "planning bootstrap")
```

## V1 Constraints

1. Default storage is local and sqlite-first.
2. AthenaWork should assume AthenaMind is sqlite-only in Cairn unless the product boundary is intentionally expanded again.
3. AthenaWork should not depend on AthenaMind-specific internal packages.
4. AthenaWork should not depend on AthenaUse-specific internal packages.
5. The CLI boundary is the source of truth for product integration until a stronger service boundary is intentionally introduced.
