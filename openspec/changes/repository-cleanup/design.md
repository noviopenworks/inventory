## Context

This change formalizes an existing repository-cleanup effort. A design doc and an implementation plan were already drafted outside the Comet workflow:

- `docs/superpowers/specs/2026-06-18-repository-cleanup-design.md`
- `docs/superpowers/plans/2026-06-18-repository-cleanup.md`

The project is a Python 3.14 / PyQt6 desktop app managed with `uv` and `go-task`, using Ruff (format + lint), mypy, and pytest. A second, in-progress effort (`gover2/`, a Wails/Vue rewrite) has **uncommitted** changes in `gover2/frontend/` that predate this cleanup and must not be disturbed.

Repository reality has been verified: every artifact targeted for removal is present and git-ignored, and every quality gate the plan references (`task fmt:check`, `task lint`, `task type:check`, `task test`) exists in `Taskfile.yml`. There is also an existing `task clean` ("Remove all build artefacts") and `task gover:clean`, which the original plan ignored in favor of a manual `rm -rf`.

## Goals / Non-Goals

**Goals:**

- Remove safe-to-recreate generated artifacts from the working tree.
- Establish (via the project's own gates) whether the Python code currently passes formatting, lint, type-check, and tests.
- Apply only minimal, diagnostic-backed fixes — each edit traceable to a specific gate finding.
- Correct concrete docs/tooling drift (README ↔ Taskfile ↔ .gitignore) where the audit proves it.

**Non-Goals:**

- Modifying the dirty `gover2/frontend` working-tree changes (or anything under `gover2/**`).
- Removing `.venv/`.
- Broad refactors, module moves, behavior changes, or new abstractions.
- Committing anything unless the user explicitly requests it.

## Decisions

- **Audit before edit.** Run the four gates first and let their output define the edit set, rather than guessing fixes up front. Rationale: keeps changes minimal and verifiable; avoids speculative churn. Alternative considered: proactively reformat/retype the whole tree — rejected as over-scoped and risky.
- **Prefer existing tooling over ad-hoc shell.** Where a `task` target already does the job (e.g., `task clean`, `task fmt`), use it and treat any gap as a tooling-drift finding to fix, rather than hand-rolling `rm -rf`/formatting. Rationale: the cleanup should leave the tooling more consistent, not parallel to it.
- **Strict isolation of gover2.** All removal and edit commands target root-level Python artifacts and explicitly exclude `gover2/**` and `.venv/`. Rationale: protects in-progress, uncommitted frontend work.
- **Blocker honesty.** If a gate cannot run (environment/dependency), record the exact blocker and report the pass as partial rather than claiming "verified." Rationale: matches the project's verification-before-completion discipline.

## Risks / Trade-offs

- [Removing `inventory.db` disrupts a local dev session] → It is git-ignored and recreated by the app on launch; low risk, but note it in the report.
- [A gate surfaces large, pre-existing debt (many lint/type errors)] → Do not auto-fix en masse; report the scope and let the user decide what is in-scope for a "minimal cleanup." Escalate to a redesign/split if fixes balloon (Comet upgrade rule).
- [Accidentally touching gover2 files] → Mitigated by explicit path exclusions and a final `git diff --stat` review confirming `gover2/frontend` is untouched.

## Migration Plan

Not applicable — no deployable behavior change, no schema change. Rollback is trivial (artifacts regenerate; edits are minimal and reviewable in `git diff`).

## Open Questions

- None blocking. The exact edit set is intentionally deferred to the build-phase audit output.
