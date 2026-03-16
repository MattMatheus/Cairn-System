# Research Council Baseline

## Purpose
The Research Council is a bounded exploration mechanism used before implementation when uncertainty is high.
It generates options, risks, and a recommendation. Final decision authority remains human-owned.

## Trigger Criteria
Run a council cycle only when at least one applies:
- unclear solution path with multiple viable approaches
- high architectural or product risk
- high expected rework cost if wrong path is chosen
- explicit operator direction to run a council

## Roles
Use three explicit roles in every council run:
- `Explorer`: proposes alternatives and upside opportunities
- `Skeptic`: stress-tests assumptions, failure modes, and hidden costs
- `Synthesizer`: converges to 1-2 viable options and recommendation

## Timebox And Limits
- Timebox each council run (default: 30-90 minutes).
- Maximum options in final output: 3.
- Council output must end with a clear recommendation and next action.
- No open-ended research without an explicit stop condition.

## Required Output Artifact
Create a single artifact under `workspace/research/` or another clearly tracked planning/research location:
- `COUNCIL-<YYYYMMDD>-<slug>.md`

Required sections:
- context and decision scope
- assumptions and constraints
- options considered (max 3)
- tradeoffs and risks
- recommendation
- confidence and unknowns
- execution next step (architect or pm)
- operator confirmation field (`confirmed_by`, `confirmed_at`)

## Integration Rules
- Planning stage may run council cycles.
- PM consumes council output and converts it into ranked intake/active stories.
- Architect consumes council output for ADR or architecture-story updates.

## Governance: User-Directed, Agent-Executed Updates
Research Council specs may be updated by agents only under explicit human direction.
When changing council policy docs/prompts:
1. Use explicit direction confirmation evidence (`CAIRN_DIRECTION_CHANGE=true` + valid `CAIRN_DIRECTION_CONFIRMATION_ID`).
2. Update all affected docs together: `AGENTS.md`, `HUMANS.md`, `DEVELOPMENT_CYCLE.md`, relevant stage prompts, and this file.
3. Record rationale and scope in observer report and cycle commit.
4. Do not merge council-policy changes without human review.
