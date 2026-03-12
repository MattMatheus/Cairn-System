<!-- AUDIENCE: Internal/Technical -->

# Planning Session Directive

Your task is to run an interactive idea-generation session with the human operator before architecture or PM execution.
When uncertainty is high, run a bounded Research Council cycle per `knowledge-base/process/RESEARCH_COUNCIL_BASELINE.md`.

## Accessibility Profile (Mandatory)
Before step 1, select and record active profile(s) for this cycle:
- `low-vision-default`
- `high-variability-attention`
- both (if needed)
Keep planning views high-clarity and keep active checklist sections to 5 items max.

## Planning Cycle (Mandatory)
1. Start a structured conversation to capture goals, users, constraints, risks, and success metrics.
2. Record notes in a new file under `product-research/planning/sessions/` using `product-research/planning/PLANNING_SESSION_TEMPLATE.md`.
3. Classify each proposed idea:
   - Implementation work -> create story in `delivery-backlog/engineering/intake/` using `delivery-backlog/engineering/intake/STORY_TEMPLATE.md`.
   - Architecture/ADR decision work -> create story in `delivery-backlog/architecture/intake/` using `delivery-backlog/architecture/intake/ARCH_STORY_TEMPLATE.md`.
4. Ensure all created intake items include traceability metadata (`idea_id`, `phase`, `adr_refs`, metric fields).
5. Provide a next-stage recommendation:
   - `architect` when architecture decisions are required first.
   - `pm` when intake is ready for refinement and ranking.
6. Set planning session status to `finalized` once intake artifacts are created and linked.
7. If Research Council is used:
   - run explicit `explorer`, `skeptic`, `synthesizer` roles;
   - keep options <= 3 and timebox <= 90 minutes;
   - create `product-research/planning/sessions/COUNCIL-<YYYYMMDD>-<slug>.md`.
8. Council/spec updates are allowed only under explicit human direction confirmation evidence (`ATHENA_DIRECTION_CONFIRMATION_ID`).
9. Run observer and capture cycle delta:
   - `tools/run_observer_cycle.sh --cycle-id <plan-id>`
10. Commit once for this cycle:
   - `cycle-<cycle-id>`

## Session Output Requirements
- Problem framing and target outcomes
- Assumptions and constraints
- Candidate ideas/options considered
- Decision/gap list by owner lane (engineering vs architecture)
- Concrete intake items created (paths + ids)
- Recommended next stage and rationale
- One-screen direction confirmation summary card:
  - `direction`
  - `constraints`
  - `next_stage`
  - `confirmed_by`

## Constraints
- Do not implement production changes in planning mode.
- Do not skip writing session notes.
- Do not place architecture decision work in engineering intake.
- Do not commit before observer report is generated.
