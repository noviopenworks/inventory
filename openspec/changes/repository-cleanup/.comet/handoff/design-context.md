# Comet Design Handoff

- Change: repository-cleanup
- Phase: design
- Mode: compact
- Context hash: 0edcd91e39ff1cbf97cc6751a7c525b7782ff55102ea3dbdee9835c66db4e86d

Generated-by: comet-handoff.sh

OpenSpec remains the canonical capability spec. This handoff is a deterministic, source-traceable context pack, not an agent-authored summary.

## openspec/changes/repository-cleanup/proposal.md

- Source: openspec/changes/repository-cleanup/proposal.md
- Lines: 1-28
- SHA256: 5d263a78760ef0563db5b4d37715a5760e603069da4f03633bcb42b60dd9f833

```md
## Why

The repository has accumulated generated, git-ignored artifacts in the working tree (`.coverage`, the various `*_cache/` dirs, `.task/`, `__pycache__/`, `build/`, `dist/`, `inventory.db`), and the Python code-quality gates have not been audited recently. A focused hygiene pass keeps the tree clean and confirms the project still passes its own formatting, lint, type, and test gates — without disturbing the in-progress `gover2/frontend` work.

## What Changes

- Remove generated, git-ignored artifacts that are safe to recreate from the working tree.
- Audit Python code quality using the existing Taskfile gates: `task fmt:check`, `task lint`, `task type:check`, `task test`.
- Apply only minimal, command-backed fixes for concrete findings (one edit ↔ one diagnostic).
- Fix concrete docs/tooling drift in `README.md`, `Taskfile.yml`, or `.gitignore` where the audit finds it (e.g., the cleanup workflow reinvents `rm -rf` even though a `task clean` target already exists).
- No product behavior changes. No commits unless explicitly requested.

## Capabilities

### New Capabilities

- `repository-hygiene`: defines the clean-working-tree, quality-gate, and docs/tooling-consistency expectations that a repository cleanup pass must satisfy, plus the blocker-reporting rule when a gate cannot run.

### Modified Capabilities

<!-- None — this change does not alter any product/spec-level behavior. -->

## Impact

- **Working tree**: removal of root-level generated artifacts (`.coverage`, `.pytest_cache/`, `.mypy_cache/`, `.ruff_cache/`, `.task/`, `__pycache__/`, `build/`, `dist/`, `inventory.db`).
- **Source (only if audit surfaces it)**: minimal edits under `app/`, `db/`, `tests/`, or `main.py` tied to a specific lint/type/test diagnostic.
- **Docs/tooling (only if drift is verified)**: minimal edits to `README.md`, `Taskfile.yml`, or `.gitignore`.
- **Explicitly NOT touched**: `gover2/**` (including the dirty `gover2/frontend` working-tree changes) and `.venv/`.
```

## openspec/changes/repository-cleanup/design.md

- Source: openspec/changes/repository-cleanup/design.md
- Lines: 1-47
- SHA256: f16d59165044ce847b3d6a1b371d96bcee8e10e059a911bea7e66b6370b3ca8c

```md
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
```

## openspec/changes/repository-cleanup/tasks.md

- Source: openspec/changes/repository-cleanup/tasks.md
- Lines: 1-33
- SHA256: 3448a06eeaca3ba029f1a055dbde76657037abeae887e4ae5d5b9cf3e313c0f8

```md
## 1. Repository hygiene audit and cleanup

- [ ] 1.1 Capture baseline: `git status --short --ignored` (confirm dirty `gover2/frontend` files and `!!` generated artifacts)
- [ ] 1.2 Verify removal candidates are ignored: `git check-ignore -v .coverage .pytest_cache .mypy_cache .ruff_cache .task __pycache__ build dist inventory.db`
- [ ] 1.3 Remove generated artifacts only (root-level; exclude `.venv/` and `gover2/**`); prefer `task clean` where it covers them, fall back to explicit `rm -rf` for the rest
- [ ] 1.4 Confirm no unrelated changes: `git status --short` shows `gover2/frontend` untouched and no source diffs from removal

## 2. Verification audit

- [ ] 2.1 `task fmt:check` — record PASS or the list of files needing formatting
- [ ] 2.2 `task lint` — record PASS or concrete Ruff diagnostics
- [ ] 2.3 `task type:check` — record PASS or concrete mypy diagnostics
- [ ] 2.4 `task test` — record PASS or concrete pytest failures

## 3. Minimal code cleanup fixes (only for Task 2 findings)

- [ ] 3.1 For formatting findings, apply `task fmt`
- [ ] 3.2 For each Ruff finding, make the smallest source edit tied to that diagnostic
- [ ] 3.3 For each mypy finding, make the smallest type-safe edit tied to that diagnostic
- [ ] 3.4 For each pytest failure, diagnose with `task test:verbose` before changing code (environmental vs. real)

## 4. Docs and tooling drift cleanup (only for verified drift)

- [ ] 4.1 Compare README development commands against `task --list`
- [ ] 4.2 Update `README.md` only for verified command drift
- [ ] 4.3 Update `.gitignore` only for verified, uncovered generated artifacts
- [ ] 4.4 Update `Taskfile.yml` only for concrete cleanup gaps (minimal, shell-safe, re-run the affected task)

## 5. Final verification and report

- [ ] 5.1 Re-run `task fmt:check`, `task lint`, `task type:check`, `task test` — record results or exact blockers
- [ ] 5.2 Review `git diff --stat`: only intended cleanup edits; `gover2/frontend` not modified by this pass
- [ ] 5.3 Summarize: removed artifacts, files changed, gate results, remaining risks/blockers
```

## openspec/changes/repository-cleanup/specs/repository-hygiene/spec.md

- Source: openspec/changes/repository-cleanup/specs/repository-hygiene/spec.md
- Lines: 1-57
- SHA256: 03f9479a4679b2ef35e8238497233e78342bc1660730e5959569d25496c00bd2

```md
## ADDED Requirements

### Requirement: Clean working tree of generated artifacts

A repository cleanup pass SHALL remove generated, git-ignored artifacts that are safe to recreate from the working tree, and SHALL NOT remove source files, lockfiles, virtual environments, or another effort's in-progress work.

#### Scenario: Generated artifacts removed

- **WHEN** a cleanup pass runs against a tree containing git-ignored generated artifacts (e.g. `.coverage`, `.pytest_cache/`, `.mypy_cache/`, `.ruff_cache/`, `.task/`, `__pycache__/`, `build/`, `dist/`, `inventory.db`)
- **THEN** those artifacts are removed from the working tree
- **AND** `.venv/` and everything under `gover2/**` remain untouched

#### Scenario: Non-ignored path is not removed blindly

- **WHEN** a candidate removal path is not matched by `.gitignore`
- **THEN** the pass inspects it before removing rather than deleting it unconditionally

### Requirement: Python quality gates pass or are reported

A repository cleanup pass SHALL run the project's existing quality gates — `task fmt:check`, `task lint`, `task type:check`, and `task test` — and SHALL either show them passing or report concrete, actionable findings.

#### Scenario: All gates pass

- **WHEN** the four gates are run after cleanup
- **THEN** each reports success
- **AND** the final report states the gates passed

#### Scenario: A gate reports findings

- **WHEN** a gate reports concrete diagnostics (a Ruff rule, an mypy error, a failing test)
- **THEN** only the smallest edits that resolve those specific diagnostics are applied
- **AND** each edit traces to a named diagnostic from the audit

#### Scenario: A gate cannot run

- **WHEN** a gate cannot run because of an environment or dependency blocker
- **THEN** the exact blocker is recorded
- **AND** the cleanup is reported as partial rather than fully verified

### Requirement: Docs and tooling consistency

A repository cleanup pass SHALL fix concrete drift between `README.md`, `Taskfile.yml`, and `.gitignore` only where the audit proves it, preferring existing tooling over ad-hoc commands.

#### Scenario: README command matches Taskfile

- **WHEN** the audit finds a development command documented in `README.md` that does not exist in `Taskfile.yml` (or vice versa)
- **THEN** the documentation is corrected to match the actual task

#### Scenario: Missing ignore coverage

- **WHEN** the audit finds a generated artifact that is not covered by `.gitignore`
- **THEN** an ignore rule is added without ignoring source files or lockfiles

#### Scenario: Prefer existing task over ad-hoc shell

- **WHEN** an existing `task` target (e.g. `task clean`) already performs a cleanup action
- **THEN** that target is used or its gap is fixed, rather than duplicating the behavior with hand-rolled shell
```
