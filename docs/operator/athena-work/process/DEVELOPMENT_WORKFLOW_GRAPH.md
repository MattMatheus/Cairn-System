---
title: Development Workflow Graph
description: Map of the first AthenaWork development-work subgraph covering cycle navigation, stage contracts, handoffs, verification, and commit invariants.
doc_type: moc
status: draft
domain: process
---

# Development Workflow Graph

## Summary

This MOC is the entry point for AthenaWork development execution. It routes agents and operators to the minimum process nodes needed to choose the next action, verify work, hand off safely, and close a cycle without rereading every process document.

## Intended Audience

- engineering-stage agents
- QA-stage agents
- PM and operator roles coordinating cycle flow

## Preconditions

- the task is part of the AthenaWork development cycle
- durable process knowledge is needed, not just active backlog state

## Main Flow

Start with the node that matches the immediate question:

- [[CYCLE_NAVIGATION]]: how to enter and move through a cycle
- [[STAGE_CONTRACTS]]: what each stage must produce before handoff
- [[HANDOFF_EXPECTATIONS]]: what a good handoff package must contain
- [[VERIFICATION_PATTERNS]]: how to decide what to run before claiming done
- [[BRANCH_AND_COMMIT_INVARIANTS]]: branch discipline and one-commit cycle closure

Related high-level sources:

- `docs/operator/athena-work/process/CYCLE_INDEX.md`
- `docs/operator/athena-work/process/STAGE_EXIT_GATES.md`
- `docs/operator/athena-work/process/OPERATOR_DAILY_WORKFLOW.md`

## How It Works

The workflow subgraph keeps stable execution rules separate from active story state:

- stable rules live in these nodes
- current priorities stay in backlog directories and program-state files
- stage prompts remain execution entrypoints, not the sole source of process knowledge

Use this graph when the question is "what should I do next or verify next?" rather than "what story is active?"

## Failure Modes

- reading the full process library for every task
- confusing cycle policy with current backlog state
- skipping handoff or verification details because they are buried in larger docs

## References

- [[DEVELOPMENT_COGNITION_GRAPH]]
- `docs/operator/athena-work/process/README.md`
