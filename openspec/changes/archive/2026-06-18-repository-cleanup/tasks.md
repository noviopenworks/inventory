## 1. Repository hygiene audit and cleanup

- [x] 1.1 Capture baseline: `git status --short --ignored` (confirm dirty `gover2/frontend` files and `!!` generated artifacts)
- [x] 1.2 Verify removal candidates are ignored: `git check-ignore -v .coverage .pytest_cache .mypy_cache .ruff_cache .task __pycache__ build dist inventory.db`
- [x] 1.3 Remove generated artifacts only (root-level; exclude `.venv/` and `gover2/**`); prefer `task clean` where it covers them, fall back to explicit `rm -rf` for the rest
- [x] 1.4 Confirm no unrelated changes: `git status --short` shows `gover2/frontend` untouched and no source diffs from removal

## 2. Verification audit

- [x] 2.1 `task fmt:check` — PASS (24 files already formatted)
- [x] 2.2 `task lint` — PASS (All checks passed)
- [x] 2.3 `task type:check` — PASS (no issues in 13 source files)
- [x] 2.4 `task test` — PASS (175 passed in 2m53s)

## 3. Minimal code cleanup fixes (only for Task 2 findings)

- [x] 3.1 For formatting findings, apply `task fmt` — N/A
- [x] 3.2 For each Ruff finding, make the smallest source edit — N/A
- [x] 3.3 For each mypy finding, make the smallest type-safe edit — N/A
- [x] 3.4 For each pytest failure, diagnose with `task test:verbose` — N/A

## 4. Docs and tooling drift cleanup (only for verified drift)

- [x] 4.1 Compare README development commands against `task --list`
- [x] 4.2 Update `README.md` only for verified command drift — N/A (no drift)
- [x] 4.3 Update `.gitignore` only for verified, uncovered generated artifacts — N/A
- [x] 4.4 Update `Taskfile.yml` only for concrete cleanup gaps — Done: extended `clean` to remove .coverage + 3 caches

## 5. Final verification and report

- [x] 5.1 Re-run `task fmt:check`, `task lint`, `task type:check`, `task test` — all PASS
- [x] 5.2 Review `git diff --stat` — only `Taskfile.yml` modified by this pass; `gover2/frontend` pre-existing and untouched
- [x] 5.3 Summarize: see outcome report
