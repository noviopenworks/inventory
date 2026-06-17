# Brainstorm Summary

- Change: gover-phase1-architecture
- Date: 2026-06-17

## Confirmed Technical Approach

Approach C — Decision narrative + embedded reference tables. Single `specs/architecture/spec.md` with each section opening with a rationale sentence followed by a precise reference table. Mirrors the Q&A style of `gover-ui-spec`.

Project scaffolds at `gover2/` subdirectory inside the current `inventory/` repo. Module path: `module gover2`, Go 1.22+. Root `Taskfile.yml` gets `gover:` namespace with `dir: gover2` on every command.

## Key Trade-offs and Risks

- `modernc.org/sqlite` (pure Go) chosen over `mattn/go-sqlite3` (CGo) — avoids C toolchain requirement, cross-compiles cleanly; slightly slower at high query volumes (not a concern for local asset inventory)
- `internal/bridge` is the only package that imports from all others; keeps coupling explicit and testable
- Wails-generated TypeScript bindings are wrapped in `src/lib/api/` — Vue components never call `window.go.*` directly, so binding changes don't ripple through components
- Vue Router hash mode: no server needed, works in Wails WebView
- `gover2/` co-located with Python app in same repo: easier transition reference, shared git history; small risk of toolchain confusion (mitigated by Taskfile namespace)

## Testing Strategy

- Go: `testing` + `testify/assert`, table-driven, co-located test files, 70% coverage target for `services`, 50% for `bridge`
- Frontend: Vitest + Vue Test Utils, co-located with components
- No E2E in Phase 1 or 2; Playwright deferred to Phase 4
- Python tests with Go equivalents in Phase 2: `test_alerts.py` → `internal/services`, `test_models.py` → `internal/models`

## Spec Patches

None — `specs/architecture/spec.md` is a new capability, no prior delta to reconcile.
