# Cycle Index

Single navigation entrypoint for operators and agents running the work harness delivery cycle.

## First 5 Minutes
1. Confirm branch safety:
   - run `git branch --show-current`
   - ensure it matches `CAIRN_REQUIRED_BRANCH` (default: `dev`)
2. Launch the stage you need:
   - `./products/work-harness/tools/launch_stage.sh planning`
   - `./products/work-harness/tools/launch_stage.sh architect`
   - `./products/work-harness/tools/launch_stage.sh engineering`
   - `./products/work-harness/tools/launch_stage.sh qa`
   - `./products/work-harness/tools/launch_stage.sh pm`
3. Open the returned seed prompt in `products/work-harness/stage-prompts/active/` and follow it as directive.
4. Run docs validation command before handoff/decision points:
   - `./products/work-harness/tools/run_doc_tests.sh`
5. Apply backlog state movement only through the canonical flow (`engineering/active -> engineering/qa -> engineering/done`, with intake/active loop for defects).
6. Close each cycle with Observer + single commit:
   - `./products/work-harness/tools/run_observer_cycle.sh --cycle-id <cycle-id>`
   - `git commit -m "cycle-<cycle-id>"`
7. Apply stage and cycle gates:
   - `docs/operator/work-harness/process/STAGE_EXIT_GATES.md`

## Branch Rule and Empty Active Behavior
- Branch safety rule: launcher requires branch `CAIRN_REQUIRED_BRANCH` (default `dev`).
- If branch differs, launcher aborts:
  - `abort: active branch is '<branch>'; expected '<required-branch>'`
- Agent branch policy:
  - default execution on `dev`; optional isolation on `agent/<agent-id>/<cycle-id>`
  - human role is management/approval; implementation is agent-owned
- If engineering is launched with no active stories, the expected output is:
  - `no stories`

## Canonical References
- Development cycle overview:
  - `docs/operator/work-harness/DEVELOPMENT_CYCLE.md`
- Stage launch script:
  - `products/work-harness/tools/launch_stage.sh`
- Observer script:
  - `products/work-harness/tools/run_observer_cycle.sh`
- Stage seed prompts:
  - `products/work-harness/stage-prompts/active/planning-seed-prompt.md`
  - `products/work-harness/stage-prompts/active/architect-agent-seed-prompt.md`
  - `products/work-harness/stage-prompts/active/next-agent-seed-prompt.md`
  - `products/work-harness/stage-prompts/active/qa-agent-seed-prompt.md`
  - `products/work-harness/stage-prompts/active/pm-refinement-seed-prompt.md`
  - `docs/operator/work-harness/process/RESEARCH_COUNCIL_BASELINE.md`
- Backlog state directories:
  - `products/work-harness/delivery-backlog/architecture/intake/`
  - `products/work-harness/delivery-backlog/architecture/ready/`
  - `products/work-harness/delivery-backlog/architecture/active/`
  - `products/work-harness/delivery-backlog/architecture/qa/`
  - `products/work-harness/delivery-backlog/architecture/done/`
  - `products/work-harness/delivery-backlog/engineering/intake/`
  - `products/work-harness/delivery-backlog/engineering/ready/`
  - `products/work-harness/delivery-backlog/engineering/active/`
  - `products/work-harness/delivery-backlog/engineering/qa/`
  - `products/work-harness/delivery-backlog/engineering/done/`
  - `products/work-harness/delivery-backlog/engineering/blocked/`
  - `products/work-harness/delivery-backlog/engineering/archive/`
- Active queue ordering:
  - `products/work-harness/delivery-backlog/engineering/active/README.md`
- Program control plane:
  - `docs/operator/work-harness/process/PROGRAM_OPERATING_SYSTEM.md`
- Observer artifacts:
  - `products/work-harness/operating-system/observer/README.md`
  - `products/work-harness/operating-system/observer/OBSERVER_REPORT_TEMPLATE.md`
- Release checkpoint template:
  - `products/work-harness/operating-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`
- Personas directory and role index:
  - `products/work-harness/staff-personas/`
  - `products/work-harness/staff-personas/STAFF_DIRECTORY.md`
- Handoff docs:
  - `products/work-harness/delivery-backlog/engineering/qa/HANDOFF-*.md`
  - `products/work-harness/delivery-backlog/engineering/done/QA-RESULT-*.md`
- Doc test harness standard:
  - `products/work-harness/operating-system/observer/LATEST_DOCS_INDEX_READ_MODEL.md`
- Agent branch strategy playbook:
  - `products/work-harness/operating-system/playbooks/agent-branch-strategy.md`
