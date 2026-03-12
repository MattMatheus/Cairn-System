# Cycle Index

Single navigation entrypoint for operators and agents running the AthenaWork delivery cycle.

## First 5 Minutes
1. Confirm branch safety:
   - run `git branch --show-current`
   - ensure it matches `ATHENA_REQUIRED_BRANCH` (default: `dev`)
2. Launch the stage you need:
   - `./products/athena-work/tools/launch_stage.sh planning`
   - `./products/athena-work/tools/launch_stage.sh architect`
   - `./products/athena-work/tools/launch_stage.sh engineering`
   - `./products/athena-work/tools/launch_stage.sh qa`
   - `./products/athena-work/tools/launch_stage.sh pm`
3. Open the returned seed prompt in `products/athena-work/stage-prompts/active/` and follow it as directive.
4. Run docs validation command before handoff/decision points:
   - `./products/athena-work/tools/run_doc_tests.sh`
5. Apply backlog state movement only through the canonical flow (`engineering/active -> engineering/qa -> engineering/done`, with intake/active loop for defects).
6. Close each cycle with Observer + single commit:
   - `./products/athena-work/tools/run_observer_cycle.sh --cycle-id <cycle-id>`
   - `git commit -m "cycle-<cycle-id>"`
7. Apply stage and cycle gates:
   - `docs/operator/athena-work/process/STAGE_EXIT_GATES.md`

## Branch Rule and Empty Active Behavior
- Branch safety rule: launcher requires branch `ATHENA_REQUIRED_BRANCH` (default `dev`).
- If branch differs, launcher aborts:
  - `abort: active branch is '<branch>'; expected '<required-branch>'`
- Agent branch policy:
  - default execution on `dev`; optional isolation on `agent/<agent-id>/<cycle-id>`
  - human role is management/approval; implementation is agent-owned
- If engineering is launched with no active stories, the expected output is:
  - `no stories`

## Canonical References
- Development cycle overview:
  - `docs/operator/athena-work/DEVELOPMENT_CYCLE.md`
- Stage launch script:
  - `products/athena-work/tools/launch_stage.sh`
- Observer script:
  - `products/athena-work/tools/run_observer_cycle.sh`
- Stage seed prompts:
  - `products/athena-work/stage-prompts/active/planning-seed-prompt.md`
  - `products/athena-work/stage-prompts/active/architect-agent-seed-prompt.md`
  - `products/athena-work/stage-prompts/active/next-agent-seed-prompt.md`
  - `products/athena-work/stage-prompts/active/qa-agent-seed-prompt.md`
  - `products/athena-work/stage-prompts/active/pm-refinement-seed-prompt.md`
  - `docs/operator/athena-work/process/RESEARCH_COUNCIL_BASELINE.md`
- Backlog state directories:
  - `products/athena-work/delivery-backlog/architecture/intake/`
  - `products/athena-work/delivery-backlog/architecture/ready/`
  - `products/athena-work/delivery-backlog/architecture/active/`
  - `products/athena-work/delivery-backlog/architecture/qa/`
  - `products/athena-work/delivery-backlog/architecture/done/`
  - `products/athena-work/delivery-backlog/engineering/intake/`
  - `products/athena-work/delivery-backlog/engineering/ready/`
  - `products/athena-work/delivery-backlog/engineering/active/`
  - `products/athena-work/delivery-backlog/engineering/qa/`
  - `products/athena-work/delivery-backlog/engineering/done/`
  - `products/athena-work/delivery-backlog/engineering/blocked/`
  - `products/athena-work/delivery-backlog/engineering/archive/`
- Active queue ordering:
  - `products/athena-work/delivery-backlog/engineering/active/README.md`
- Program control plane:
  - `docs/operator/athena-work/process/PROGRAM_OPERATING_SYSTEM.md`
- Observer artifacts:
  - `products/athena-work/operating-system/observer/README.md`
  - `products/athena-work/operating-system/observer/OBSERVER_REPORT_TEMPLATE.md`
- Release checkpoint template:
  - `products/athena-work/operating-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`
- Personas directory and role index:
  - `products/athena-work/staff-personas/`
  - `products/athena-work/staff-personas/STAFF_DIRECTORY.md`
- Handoff docs:
  - `products/athena-work/delivery-backlog/engineering/qa/HANDOFF-*.md`
  - `products/athena-work/delivery-backlog/engineering/done/QA-RESULT-*.md`
- Doc test harness standard:
  - `products/athena-work/operating-system/observer/LATEST_DOCS_INDEX_READ_MODEL.md`
- Agent branch strategy playbook:
  - `products/athena-work/operating-system/playbooks/agent-branch-strategy.md`
