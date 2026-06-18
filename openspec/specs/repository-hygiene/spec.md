# repository-hygiene Specification

## Purpose
TBD - created by archiving change repository-cleanup. Update Purpose after archive.
## Requirements
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
