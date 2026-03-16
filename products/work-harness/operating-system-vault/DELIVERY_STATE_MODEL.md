# Delivery State Model

## Lanes
- Architecture: `intake -> active -> qa -> done`
- Engineering: `intake -> active -> qa -> done`

## Queue Paths
- `70 Agents/05 Delivery Queue/Architecture Intake/`
- `70 Agents/05 Delivery Queue/Architecture Active/`
- `70 Agents/05 Delivery Queue/Architecture QA/`
- `70 Agents/05 Delivery Queue/Architecture Done/`
- `70 Agents/05 Delivery Queue/Engineering Intake/`
- `70 Agents/05 Delivery Queue/Engineering Active/`
- `70 Agents/05 Delivery Queue/Engineering QA/`
- `70 Agents/05 Delivery Queue/Engineering Done/`

## Rules
- No direct `intake -> done`.
- No direct `active -> done`.
- Defects discovered in QA return work to `active` with linked bug/task.

## Escape Path

Any item in any delivery state may transition to `blocked` via an escape record.

```
[intake | active | qa] -> blocked (escape)
                               |
                       human resolves
                               |
               retry (return to lane) | redirect | cancel
```

- Escape record created in `70 Agents/70 Escape/` using `TemplateEscape.md`.
- Source item status set to `blocked`; item remains in its current lane folder.
- Human sets `resolution_action` and returns item to queue or closes it.
- Escape records are permanent audit artifacts; do not delete after resolution.
