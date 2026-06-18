# Tasks — gover2-parity

## 1. Go: CSV export (data-export)

- [x] 1.1 Add `services/export.go` with a pure function that builds CSV bytes/rows for a given category (header = view columns, rows in column order), unit-tested with in-memory SQLite
- [ ] 1.2 Reshape bridge `ExportCSV(category)` to open `runtime.SaveFileDialog(ctx, ...)`, write the CSV to the chosen path, and return nil on cancel (nil-db guard preserved)
- [ ] 1.3 Add bridge test for the path-based export (write to temp file, assert contents); thin dialog wrapper left to smoke test

## 2. Frontend: theme + preferences (app-preferences)

- [ ] 2.1 Finish `stores/ui.ts`: load `darkMode`/`density` from `GetConfig` on startup; `toggleDarkMode` and density setters persist via `SetConfig`
- [ ] 2.2 Apply dark mode via a root class on AppShell/`<html>`; audit components for hardcoded colors and migrate to Tailwind semantic tokens
- [ ] 2.3 Apply density + theme-aware alternating row colors in `DataTable.vue`
- [ ] 2.4 Tests: ui store persistence/toggle; DataTable density class rendering

## 3. Frontend: topbar action row

- [ ] 3.1 Rework `Topbar.vue` into an action row: search box, ⚠ Alerts, ↓ CSV, 🌙 dark-mode toggle, ⋯ overflow menu
- [ ] 3.2 Wire the 🌙 button to `stores/ui.ts`; reflect current theme state in the icon/label
- [ ] 3.3 Wire the ↓ CSV action to call `exportCsv(activeCategory)` for the active route
- [ ] 3.4 Overflow menu entries: New Database, Open Database, About, License

## 4. Frontend: alerts modal (asset-alerts)

- [ ] 4.1 Populate `AlertsModal.vue` from `GetAlerts` (category, name, expiry date, days remaining, severity), with severity styling
- [ ] 4.2 Own the modal at app-shell level; auto-open on startup when `GetAlerts` returns a non-empty list
- [ ] 4.3 Wire the ⚠ Alerts topbar action to reopen the modal
- [ ] 4.4 Test: modal renders rows from a fixture; auto-open logic for empty vs non-empty

## 5. Frontend: search/filter (asset-search)

- [ ] 5.1 Add a shared filter helper that matches a query against a row's visible display values
- [ ] 5.2 Wire the topbar search term to the active view; filter the 7 category views + AllAssets client-side, updating as the user types
- [ ] 5.3 Test: filter helper matches/clears correctly

## 6. Frontend: database management (database-management) + dialogs

- [ ] 6.1 Wire New Database / Open Database overflow entries to `NewDatabase`/`OpenDatabase` bridge methods (native dialogs); refresh views after switch
- [ ] 6.2 Add About and License modal components reachable from the overflow menu
- [ ] 6.3 Test: About/License modals open/close

## 7. Integration & verification

- [ ] 7.1 Regenerate wailsjs bindings (`wails build`) after the `ExportCSV` signature change; `go build ./...` and `pnpm run typecheck` clean
- [ ] 7.2 All Go tests + frontend tests pass; coverage targets met (services ≥70%, bridge ≥50%)
- [ ] 7.3 Smoke test: startup alerts, CSV export, dark-mode persistence across restart, search, New/Open DB, About/License
