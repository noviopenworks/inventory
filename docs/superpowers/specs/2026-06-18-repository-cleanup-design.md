---
comet_change: repository-cleanup
role: technical-design
canonical_spec: openspec
---

# Repository Cleanup Design

> Canonical requirements live in the OpenSpec change `repository-cleanup`
> (`proposal.md`, `specs/repository-hygiene/spec.md`, `tasks.md`). This document
> captures the technical design only.

## Goal

Clean up repository hygiene, code quality, and docs/tooling drift without changing
product behavior or touching unrelated in-progress frontend work.

## Scope

- Remove local generated artifacts already covered by ignore rules: caches, build
  outputs, coverage files, and the runtime database.
- Audit Python application quality with the existing project commands: formatting
  check, lint, type-check, and tests.
- Apply small, verified code cleanup only where issues are concrete and low-risk.
- Update README, Taskfile, or `.gitignore` only when the audit finds concrete drift
  or missing ignore coverage.

## Explicit Non-Goals

- Do not modify the existing dirty `gover2/frontend` worktree changes (or anything
  under `gover2/**`) unless the user explicitly asks.
- Do not remove `.venv/`.
- Do not perform broad refactors, module moves, or behavior changes.
- Do not add compatibility layers or new abstractions unless required by a verified
  issue.
- Do not commit changes unless explicitly requested.

## Approach

Use a safe incremental cleanup pass:

1. Inspect repository status and ignored/generated artifacts.
2. Remove generated local artifacts that are safe to recreate.
3. Run existing verification commands to surface actionable issues.
4. Make minimal edits for clear code/docs/tooling cleanup findings.
5. Re-run verification and report any blockers with exact command output.

## Key Decisions

- **Audit before edit.** Run the four gates first; their output defines the edit
  set. Avoids speculative churn and keeps changes traceable (one edit ↔ one
  diagnostic).
- **Cleanup command set.** `task clean` only removes `dist/`, `build/`, `.task/`,
  and `__pycache__` trees, and it carries an interactive `prompt:`. Therefore:
  - Run `task clean --yes` for the artifacts it covers.
  - Run explicit `rm -rf .coverage .pytest_cache .mypy_cache .ruff_cache inventory.db`
    for the artifacts it misses.
  - All paths are root-level; `gover2/**` and `.venv/` are never named.
- **Tooling-drift fix (verified).** `task clean` is described as "Remove all build
  artefacts" but leaves the lint/type/test caches and `.coverage` behind. Extend
  `task clean` to also remove `.coverage .pytest_cache .mypy_cache .ruff_cache` so
  the description is accurate. Deliberately do **not** add `inventory.db` — it is
  runtime data, not a build artefact, and a user would not expect `clean` to delete
  their database.
- **No `.gitignore` change.** Ignore coverage is already complete for every removal
  target; no rule is added.
- **Strict gover2 isolation.** Removal and edit commands exclude `gover2/**` and
  `.venv/`; a final `git diff --stat` confirms `gover2/frontend` is untouched.
- **Blocker honesty.** If a gate cannot run, record the exact blocker and report the
  pass as partial rather than "verified."

## Verification

Target verification commands:

- `task fmt:check`
- `task lint`
- `task type:check`
- `task test`

If a command cannot run because of environment or dependency constraints, record the
exact blocker instead of treating the cleanup as fully verified. A final
`git diff --stat` must show only intended cleanup edits, with `gover2/frontend`
unmodified by this pass.

## Risks / Trade-offs

- [Removing `inventory.db` disrupts a local dev session] → ignored and recreated by
  the app on launch; low risk, noted in the report.
- [Gates surface large pre-existing debt] → do not mass-fix; report scope and let the
  user decide. Escalate to redesign/split if fixes balloon (comet upgrade rule).
- [Accidentally touching gover2 files] → mitigated by explicit path exclusions and a
  final `git diff --stat` review.
