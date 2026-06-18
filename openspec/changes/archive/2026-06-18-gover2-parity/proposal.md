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
