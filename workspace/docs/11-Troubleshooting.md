# Troubleshooting

## Symptom: Work item moved to wrong lane
- Action: move it to the correct queue lane and correct status metadata.

## Symptom: Claim marked supported without evidence
- Action: revert claim to active-test or candidate; add required artifact links.

## Symptom: Agent proposed canonical change without approval
- Action: move proposal to staging and require review decision.

## Symptom: Sensitive content appears in canonical zones
- Action: apply rejection policy immediately; redact, quarantine, or block.

## Symptom: Contract drift between docs and operations
- Action: reconcile `products/athena-work/operating-system-vault/*`, `workspace/AGENTS.md`, `workspace/docs/HUMANS.md`, and `workspace/docs/*` in one change set.

## Symptom: Task note missing required task metadata
- Action: run `tools/platform/validate_task_metadata.sh` from repo root (checks changed files only) and fix reported fields (`task_type`, `output_type`, `acceptance_evidence`).
- Action: run `tools/platform/validate_task_metadata.sh --all .` for full-repo sweep when doing metadata cleanup.

## Symptom: User requests personal content removal from repository
- Action: run `tools/migration/prune_user_content.sh --export-zip /writable/path/athenawork-user-content-backup.zip`.
- Action: if export path write fails due sandbox/permission constraints, rerun with escalated permissions.
