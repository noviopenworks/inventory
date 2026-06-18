# Brainstorm Summary

- Change: repository-cleanup
- Date: 2026-06-18

## Confirmed Technical Approach

Low-risk repository hygiene + Python quality-gate audit, formalizing pre-existing
design+plan into comet. Audit-before-edit: run the project's own gates and let
their output define the (minimal) edit set.

- **Cleanup:** `task clean --yes` (covers dist/, build/, .task/, __pycache__) then
  explicit `rm -rf .coverage .pytest_cache .mypy_cache .ruff_cache inventory.db`.
  All root-level; `gover2/**` and `.venv/` never named.
- **Audit:** `task fmt:check → lint → type:check → test`, capture verbatim output.
- **Fixes:** only diagnostic-backed edits, one edit ↔ one finding; diagnose test
  failures with `task test:verbose` before editing.
- **Tooling drift (the one verified finding):** extend `task clean` to also remove
  the four caches + `.coverage` so its "Remove all build artefacts" description is
  accurate. Deliberately NOT adding `inventory.db` (runtime data, not a build
  artefact). No `.gitignore` gap — coverage is already complete.

## Key Trade-offs and Risks

- `task clean` had a real gap (only dist/build/.task/pycache; missed caches +
  coverage) and an interactive `prompt:` → automated use needs `--yes`.
- Pre-existing lint/type/test debt could balloon scope → cap at minimal fixes,
  escalate to user if large (comet upgrade rule).
- `inventory.db` removal low-risk (ignored, app-recreated) but noted in report.

## Testing Strategy

The project's own gates are the test: `task clean` (verify correct set removed,
gover2 untouched) + the four quality gates re-run green at the end; final
`git diff --stat` confirms `gover2/frontend` is unmodified.

## Spec Patches

None. The delta spec already covers this — its "prefer existing task over ad-hoc
shell" scenario uses `task clean` as the example, now validated as accurate.
