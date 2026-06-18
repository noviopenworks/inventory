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
