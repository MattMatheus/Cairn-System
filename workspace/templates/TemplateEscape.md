---
id: ath-escape-
type: escape
status: blocked
domain: work
updated: YYYY-MM-DD
owner: matt
source_of_truth: agent
sensitivity: internal
agent_last_touch: YYYY-MM-DDTHH:MM:SSZ
agent_intent: escalate
review_state: pending
---

# Escape: <short description>

## Source Item

| Field | Value |
|-------|-------|
| id | ath- |
| path | relative/path/to/item.md |
| type | task\|story\|bug |
| state_at_escape | intake\|active\|qa |
| lane | Architecture\|Engineering |

## Escape Classification

```yaml
escape_class: blocked|ambiguous_requirements|policy_conflict|missing_context|low_confidence|loop_detected|scope_exceeded
escape_summary: "One sentence description of why the agent cannot proceed."
attempt_count: 1
```

### Escape Class Reference

| Class | When to Use |
|-------|-------------|
| `blocked` | Hard prerequisite not met (dependency, external system, missing artifact) |
| `ambiguous_requirements` | Acceptance criteria too unclear to determine done |
| `policy_conflict` | Two rules contradict and agent cannot self-resolve |
| `missing_context` | Insufficient information to proceed safely |
| `low_confidence` | Agent can attempt but confidence is below threshold for this transition |
| `loop_detected` | Same transition attempted 2+ times without forward progress |
| `scope_exceeded` | Task scope requires authorization beyond agent write boundaries |

## What Was Attempted

-

## What Is Needed to Unblock

-

## Relevant Context

> Paste or link any relevant excerpts, error messages, conflicting references, or prior run notes here.

---

## Resolution

```yaml
resolution_action: retry|redirect|cancel
resolved_by: matt
resolved_at: YYYY-MM-DD
resolution_notes: ""
```

### Resolution Notes

_Human completes this section._
