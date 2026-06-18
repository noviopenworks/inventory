---
change: repository-cleanup
design-doc: docs/superpowers/specs/2026-06-18-repository-cleanup-design.md
base-ref: 07cfbe8a610601dc417eee65b608402c729166f1
---

# Repository Cleanup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans (or superpowers:subagent-driven-development) to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.
>
> Canonical tasks: `openspec/changes/repository-cleanup/tasks.md`. Keep both in sync.

**Goal:** Clean generated artifacts, fix concrete Python code/docs/tooling cleanup findings, and verify the repository without touching existing `gover2/frontend` changes.

**Architecture:** Safe incremental cleanup pass. Audit-before-edit: existing repository tasks define the edit set; only minimal, diagnostic-backed edits are applied; no product behavior changes.

**Tech Stack:** Python 3.14, PyQt6, uv, go-task, Ruff, mypy, pytest.

---

## File Structure

- Modify: `Taskfile.yml` — extend `clean` to also remove the lint/type/test caches + `.coverage` (verified drift fix).
- Potentially modify: `.gitignore` only if a generated artifact is found uncovered (none found at design time).
- Potentially modify: `README.md` if development commands drift from `Taskfile.yml`.
- Potentially modify: Python files under `app/`, `db/`, `tests/`, or `main.py` only for concrete lint/type/test findings.
- Do not modify: anything under `gover2/**` (including `gover2/frontend/*` dirty changes) or `.venv/`.
- Do not commit changes unless the user explicitly requests it.

---

### Task 1: Repository Hygiene Audit And Cleanup

**Files:** Inspect `.gitignore`, `Taskfile.yml`. Remove root-level generated artifacts only.

- [x] **Step 1: Capture baseline status**

  Run: `git status --short --ignored`

  Expected: dirty `gover2/frontend` files appear; generated ignored artifacts appear with `!!`.

- [x] **Step 2: Verify removal candidates are ignored**

  Run: `git check-ignore -v .coverage .pytest_cache .mypy_cache .ruff_cache .task __pycache__ build dist inventory.db`

  Expected: each existing generated path matches a `.gitignore` rule. If a path is not ignored, inspect before removing.

- [x] **Step 3: Remove generated artifacts only**

  Run: `task clean --yes` (covers `dist/`, `build/`, `.task/`, `__pycache__` trees), then
  `rm -rf .coverage .pytest_cache .mypy_cache .ruff_cache inventory.db` (the artifacts `clean` does not cover).

  Expected: generated artifacts removed; source files, `.venv/`, and `gover2/**` untouched.

- [x] **Step 4: Confirm no unrelated changes**

  Run: `git status --short`

  Expected: `gover2/frontend` changes remain; no source diffs from artifact removal.

---

### Task 2: Verification Audit

**Files:** Inspect command output only.

- [x] **Step 1: Check formatting** — Run: `task fmt:check` — Result: PASS (24 files already formatted)
- [x] **Step 2: Run lint** — Run: `task lint` — Result: PASS (All checks passed)
- [x] **Step 3: Run type-check** — Run: `task type:check` — Result: PASS (no issues in 13 source files)
- [x] **Step 4: Run tests** — Run: `task test` — Result: PASS (175 passed in 2m53s)

---

### Task 3: Minimal Code Cleanup Fixes (only for Task 2 findings)

**Files:** Modify only files named by Task 2 diagnostics (under `app/`, `db/`, `tests/`, or `main.py`).

- [x] **Step 1:** For formatting findings, apply `task fmt`. — N/A (fmt:check passed)
- [x] **Step 2:** For each Ruff finding, make the smallest source edit. — N/A (lint passed)
- [x] **Step 3:** For each mypy finding, make the smallest type-safe edit. — N/A (type:check passed)
- [x] **Step 4:** For each pytest failure, diagnose with `task test:verbose`. — N/A (tests passed)

---

### Task 4: Docs And Tooling Drift Cleanup

**Files:** `Taskfile.yml` (the verified drift fix), and `README.md` / `.gitignore` only if drift is found.

- [x] **Step 1: Extend `task clean` (verified drift fix)**

  Edit `Taskfile.yml` `clean` task to also remove `.coverage .pytest_cache .mypy_cache .ruff_cache` so its "Remove all build artefacts" description is accurate. Do NOT add `inventory.db` (runtime data, not a build artefact).

  Result: added `rm -rf .coverage .pytest_cache .mypy_cache .ruff_cache` to clean cmds; prompt updated.

- [x] **Step 2: Compare README commands to Taskfile** — Run: `task --list`; update `README.md` only for verified command drift.

  Result: No drift. README already says `task clean # remove dist/, build/ and caches`; all task names match.

- [x] **Step 3: Update `.gitignore`** only for a verified, uncovered generated artifact — N/A (all targets already covered).

---

### Task 5: Final Verification And Report

**Files:** Inspect `git status --short`, `git diff --stat`, `git diff`.

- [x] **Step 1:** Re-run `task fmt:check` — PASS (24 files already formatted)
- [x] **Step 2:** Re-run `task lint` — PASS (All checks passed)
- [x] **Step 3:** Re-run `task type:check` — PASS (no issues in 13 source files)
- [x] **Step 4:** Re-run `task test` — PASS (175 passed in 2m42s)
- [x] **Step 5:** `git diff --stat` — only `Taskfile.yml` (+2 lines) changed by this pass; `gover2/frontend` diffs are pre-existing and untouched.
- [x] **Step 6:** See outcome summary below.
