# Escape Queue

This folder holds escape records — structured signals from agents that cannot safely proceed on a work item.

## Purpose

An escape record is not a failure. It is a deliberate, safe stop. Agents use it when:
- A hard dependency is missing
- Requirements are too ambiguous to determine done
- Two policies conflict and the agent cannot self-resolve
- Insufficient context exists to proceed safely
- Confidence is below the threshold for the current transition
- The same transition has been attempted multiple times without progress
- The task scope exceeds agent authorization

## Zone

Zone A — agent may create escape records freely without human pre-approval.

## Lifecycle

```
[any delivery state] -> blocked (escape record created here)
                            |
                    human reviews escape record
                            |
              ┌─────────────┼─────────────┐
           retry        redirect        cancel
        (return to      (new path)    (close item)
         queue)
```

## Escape Record Naming

```
ESC-<source-id>-<YYYY-MM-DD>.md
```

Example: `ESC-ath-task-042-2026-03-09.md`

## Resolution Contract

When resolving an escape record, the human must:
1. Set `resolution_action` to one of: `retry | redirect | cancel`
2. Set `review_state: approved`
3. Set `resolved_by` and `resolved_at`
4. Update the source item's status and return it to the appropriate queue lane if retrying or redirecting

## Template

Use `workspace/templates/TemplateEscape.md` for all new escape records.
