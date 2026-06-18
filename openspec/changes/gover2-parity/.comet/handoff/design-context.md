# Comet Design Handoff

- Change: gover2-parity
- Phase: design
- Mode: compact
- Context hash: 9b7a3a210598c021fd022ffeab241f6cb878756c79d3e2def28bf05cc49f339d

Generated-by: comet-handoff.sh

OpenSpec remains the canonical capability spec. This handoff is a deterministic, source-traceable context pack, not an agent-authored summary.

## openspec/changes/gover2-parity/proposal.md

- Source: openspec/changes/gover2-parity/proposal.md
- Lines: 1-33
- SHA256: 28a616df8261c429f04e0b34bdcf1b7c81b4f92e0ae19b39c7341b9ad9571100

```md
## Why

The gover2 Wails app has complete CRUD for all 7 asset categories but is missing the supporting features that make the original PyQt inventory app usable day-to-day: a startup expiry-alerts popup, CSV export, working dark mode, search, and database switching from the UI. This change closes the parity gap while polishing (not redesigning) the existing UI.

## What Changes

- **Alerts popup**: wire `GetAlerts` into the existing `AlertsModal`, show it automatically on startup when any asset expires within the warning window, and re-openable from a ⚠ Alerts topbar action.
- **CSV export**: implement the stubbed `ExportCSV` bridge method using Wails' native save dialog, writing the currently active view's columns and rows to a CSV file.
- **Dark mode**: complete the half-wired theme toggle — a 🌙 topbar button flips the whole UI between light/dark, applied via a root theme class, and the choice persists through `AppConfig.DarkMode`.
- **Search/filter**: a topbar search box that filters the current view's already-loaded rows across all visible fields in real time (client-side, no new bridge methods).
- **Database switching**: surface the existing `NewDatabase`/`OpenDatabase` bridge methods through a topbar overflow menu (New / Open) using native file dialogs.
- **About / License**: simple modal dialogs reachable from the topbar overflow menu.
- **Density**: apply the existing `AppConfig.Density` (comfortable/compact) plus alternating row colors to the tables.
- **Topbar action row**: consolidate search, ⚠ Alerts, ↓ CSV, 🌙 dark-mode, and the ⋯ overflow menu into a single polished topbar row.

## Capabilities

### New Capabilities
- `asset-alerts`: startup + on-demand expiry/warranty alerts surfaced through a modal.
- `data-export`: export the active view to a CSV file via a native save dialog.
- `app-preferences`: dark mode + density preferences applied to the UI and persisted in app config.
- `asset-search`: client-side search/filter of the active view across visible fields.
- `database-management`: create/open a database file from the UI via native dialogs.

### Modified Capabilities
<!-- No existing main specs in openspec/specs/ to modify; all capabilities are new. -->

## Impact

- **Go**: implement `ExportCSV` in `internal/bridge`; add CSV-building + native-dialog logic (likely a new `internal/services/export.go` and use of `wails/v2/pkg/runtime`). `GetAlerts`, `GetConfig`/`SetConfig`, `NewDatabase`/`OpenDatabase` already exist and are reused.
- **Frontend**: rework `Topbar.vue` into an action row; finish `stores/ui.ts` dark-mode wiring + persistence; populate `AlertsModal.vue`; add search box + filtering to the 7 category views and All Assets; add About/License modals; apply density/zebra styling to `DataTable.vue`.
- **Config**: `AppConfig` (Density, DarkMode, ExpiryWarningDays, DBPath) already exists; dark mode + density now actually read/written.
- **No** schema changes, no new asset categories, no inline cell editing, no visual redesign.
```

## openspec/changes/gover2-parity/design.md

- Source: openspec/changes/gover2-parity/design.md
- Lines: 1-39
- SHA256: ea14ece95dffc7e61dd02d90b38d5f20292bda4b6bd57a0f4660f2c6bfa9b8e0

```md
## Context

gover2 is a Wails v2 + Vue 3 (Pinia, vue-router, Tailwind) desktop app. The Go layer follows a strict dependency chain: `models` ← `database` ← `services` ← `bridge` (the `App` struct holds `*sql.DB`, `models.AppConfig`, and `context.Context`). The frontend calls bridge methods through `lib/api/index.ts` (`call<T>` proxy over `window.go.bridge.App`).

Existing, reusable building blocks:
- Bridge: `GetAlerts`, `GetConfig`/`SetConfig`, `NewDatabase`/`OpenDatabase`, and a stubbed `ExportCSV(category, destPath) error`.
- Frontend: `AlertsModal.vue` (placeholder), `stores/ui.ts` (`darkMode`/`density` refs, not applied), `Topbar.vue` (title only), `DataTable.vue`, the 7 category views + `AllAssetsView` + `AlertsView`.
- `AppConfig` already carries `Density`, `DarkMode`, `ExpiryWarningDays`, `DBPath`.

This change wires these together and polishes styling; it does not introduce new architectural patterns or asset categories.

## Goals / Non-Goals

**Goals:**
- Feature parity with the PyQt app: startup alerts popup, CSV export, dark mode, search, DB new/open, About/License, density.
- Keep the current layout; refine spacing/readability and make dark mode work.
- Consolidate global actions into a single Topbar action row.
- Reuse existing bridge methods wherever they already exist; add Go code only for CSV export.

**Non-Goals:**
- No visual redesign / new design language (polish only).
- No new asset categories, no inline cell editing, no undo/redo, no multi-user/cloud/server.
- No new persistence schema; `AppConfig` is unchanged.

## Decisions

- **Theme application**: apply dark mode by toggling a root class (e.g. `dark` on `<html>`/AppShell) that Tailwind tokens switch on. `stores/ui.ts` becomes the single source of truth; on startup it reads `AppConfig.DarkMode` via `GetConfig`, and `toggleDarkMode` persists via `SetConfig`. Avoids per-component theme props.
- **CSV export from Go**: native file dialogs require the Wails runtime + app context, so export is driven from the bridge. Reshape to `ExportCSV(category string) error`: the bridge opens `runtime.SaveFileDialog(ctx, ...)`, and on a chosen path delegates to a new `services` export function that queries the category and writes CSV. Empty path (user cancelled) returns nil with no file written. Column order matches the view.
- **Search is client-side**: each view keeps its loaded rows and derives a filtered list from a search string; the Topbar emits/owns the search term for the active route. No bridge changes, no DB round-trips per keystroke.
- **Alerts on startup**: the alerts modal is owned at the app-shell level so it can auto-open after the first `GetAlerts` returns a non-empty list, and the ⚠ Topbar action re-opens it on demand.
- **Topbar as action host**: `Topbar.vue` becomes an action row (search, ⚠ Alerts, ↓ CSV, 🌙 dark-mode, ⋯ overflow with New/Open DB + About + License). Per-view headers keep their own Add button.
- **Density + zebra**: `DataTable.vue` reads density from `stores/ui.ts` to switch row padding and applies alternating row backgrounds, both theme-aware.

## Risks / Trade-offs

- **CSV export signature change**: changing `ExportCSV(category, destPath)` to `ExportCSV(category)` updates the bridge API and regenerates wailsjs bindings; the existing stub has no real callers, so blast radius is low.
- **Native dialog testing**: `runtime.SaveFileDialog` cannot run in headless Go tests; the CSV-writing logic will be factored into a pure `services` function that is unit-tested with an explicit path, while the dialog wrapper stays thin and is covered by manual/smoke testing.
- **Dark mode coverage**: existing components use Tailwind semantic tokens (`bg-surface`, `text-text-primary`, etc.); any hardcoded colors must be migrated to tokens or they won't theme. Audit during build.
- **Search across "all visible fields"**: filtering stringifies/normalizes row values; nested or formatted fields (dates, FK names) must match what the user sees — handled per-view using the same display values as the table.
```

## openspec/changes/gover2-parity/tasks.md

- Source: openspec/changes/gover2-parity/tasks.md
- Lines: 1-46
- SHA256: 1a02a2bf38b6fdcb36bc3db2b7db6290f6ef955f5b6f781bbc26c062cc020dbd

```md
# Tasks — gover2-parity

## 1. Go: CSV export (data-export)

- [ ] 1.1 Add `services/export.go` with a pure function that builds CSV bytes/rows for a given category (header = view columns, rows in column order), unit-tested with in-memory SQLite
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
```

## openspec/changes/gover2-parity/specs/app-preferences/spec.md

- Source: openspec/changes/gover2-parity/specs/app-preferences/spec.md
- Lines: 1-24
- SHA256: 6ae5b181e29003ab7bb7405e43d57a1ac671f802a81af0c74a9de0eae67a045b

```md
## ADDED Requirements

### Requirement: Dark mode toggle and persistence

The user SHALL be able to toggle the UI between light and dark themes, and the chosen theme SHALL persist across application restarts.

#### Scenario: Toggle dark mode

- **WHEN** the user clicks the 🌙 dark-mode action
- **THEN** the entire UI switches between light and dark themes

#### Scenario: Theme persists across restart

- **WHEN** the user has selected dark mode and restarts the app
- **THEN** the app launches in dark mode

### Requirement: Table density

The tables SHALL render according to the configured density (comfortable or compact) with alternating row colors.

#### Scenario: Apply density to tables

- **WHEN** the density preference is set to compact or comfortable
- **THEN** all data tables render with the corresponding row spacing and theme-aware alternating row colors
```

## openspec/changes/gover2-parity/specs/asset-alerts/spec.md

- Source: openspec/changes/gover2-parity/specs/asset-alerts/spec.md
- Lines: 1-24
- SHA256: fa0147a779dea992c93816651c007ac4dc66e210a51f475bfa34ea011ed58def

```md
## ADDED Requirements

### Requirement: Startup expiry alerts

On application startup, the app SHALL query for assets whose expiry or warranty date falls within the configured warning window and, if any exist, present them in an alerts modal automatically.

#### Scenario: Assets expiring within the window on launch

- **WHEN** the app starts and at least one asset expires within `ExpiryWarningDays`
- **THEN** the alerts modal opens automatically listing each expiring asset with its category, name, expiry date, days remaining, and severity

#### Scenario: No expiring assets on launch

- **WHEN** the app starts and no asset expires within the warning window
- **THEN** no alerts modal is shown

### Requirement: On-demand alerts

The user SHALL be able to reopen the alerts modal at any time from a topbar action.

#### Scenario: Reopen alerts from the topbar

- **WHEN** the user clicks the ⚠ Alerts action in the topbar
- **THEN** the alerts modal opens showing the current expiring assets
```

## openspec/changes/gover2-parity/specs/asset-search/spec.md

- Source: openspec/changes/gover2-parity/specs/asset-search/spec.md
- Lines: 1-15
- SHA256: f684b98217cca680d4629862156c0e77f8d5587984907c375e21ad79f52ae71f

```md
## ADDED Requirements

### Requirement: Search the active view

The user SHALL be able to filter the currently active view's rows in real time by typing in a topbar search box, matching across the view's visible fields.

#### Scenario: Filter rows by query

- **WHEN** the user types text in the search box on a view
- **THEN** the table shows only rows whose visible field values contain the query, updating as the user types

#### Scenario: Clear the query

- **WHEN** the user clears the search box
- **THEN** the table shows all rows of the active view again
```

## openspec/changes/gover2-parity/specs/database-management/spec.md

- Source: openspec/changes/gover2-parity/specs/database-management/spec.md
- Lines: 1-15
- SHA256: 35366e93b508a929e3bc572b9fdf70f8d92194cfe71e4f5c765b6f55314dbc10

```md
## ADDED Requirements

### Requirement: Create and open databases from the UI

The user SHALL be able to create a new database file or open an existing one through native dialogs reachable from the topbar.

#### Scenario: Create a new database

- **WHEN** the user chooses New Database and selects a destination path
- **THEN** a new database file with the initialized schema is created and becomes the active database

#### Scenario: Open an existing database

- **WHEN** the user chooses Open Database and selects an existing database file
- **THEN** that database becomes the active database and the views reflect its contents
```

## openspec/changes/gover2-parity/specs/data-export/spec.md

- Source: openspec/changes/gover2-parity/specs/data-export/spec.md
- Lines: 1-15
- SHA256: b18485034759bb39a89e856fceb083a8f4ef52a1fcaf490ff25d7ad8ad6b4d60

```md
## ADDED Requirements

### Requirement: Export active view to CSV

The user SHALL be able to export the currently active asset view to a CSV file chosen through a native save dialog.

#### Scenario: Export with a chosen path

- **WHEN** the user triggers Export CSV on a category view and selects a destination path
- **THEN** a CSV file is written at that path containing the view's columns as the header row and one row per record in the view's column order

#### Scenario: Cancel the save dialog

- **WHEN** the user triggers Export CSV and cancels the native save dialog
- **THEN** no file is written and no error is reported
```
