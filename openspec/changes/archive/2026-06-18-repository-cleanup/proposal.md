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
