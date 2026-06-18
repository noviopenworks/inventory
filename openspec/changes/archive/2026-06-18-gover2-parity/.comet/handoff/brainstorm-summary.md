# Brainstorm Summary

- Change: gover2-parity
- Date: 2026-06-18

## Confirmed Technical Approach

Feature-parity + polish for gover2. Reuse existing bridge methods; add Go only for CSV export.

- **Dark mode (CSS variables)**: convert tailwind.config color tokens to `var(--c-*)`; define light values under `:root` and dark under `.dark` in style.css; `stores/ui.ts` toggles `.dark` on `document.documentElement`. Same token names → no component class churn. Hardcoded overlay/red/status colors stay theme-agnostic.
- **CSV export (Go-side)**: `services/export.go BuildCSV(db, category) ([][]string, error)` pure + canonical per-category columns; `bridge.ExportCSV(category)` nil-db guard → `runtime.SaveFileDialog(ctx)` → cancel returns nil → `writeCSV(path, rows)`. api wrapper drops destPath → `exportCSV(category)`; regenerate wailsjs.
- **Search**: new `stores/search.ts` holds `query`; topbar binds input; views compute `filteredRows` via shared `matchesQuery(row, columns, query)` on display values; `watch(route)` clears query on navigation.
- **Topbar**: action row (search, ⚠ Alerts, ↓ CSV, 🌙 dark-mode, ⋯ overflow=New/Open DB + About + License). `routePath→category` map drives export; export button hidden on `/all`; per-view Add buttons stay.
- **Alerts**: `AppShell.vue` owns modal state, auto-opens on mount when `getAlerts()` non-empty; ⚠ action re-fetches + opens; `AlertsModal.vue` renders rows with severity badges.
- **Density + zebra**: `DataTable.vue` reads density from ui store (comfortable=py-2 / compact=py-1) + themed odd/even zebra.

## Key Trade-offs and Risks

- CSV export signature change (`ExportCSV(category, destPath)` → `ExportCSV(category)`) regenerates wailsjs; stub has no real callers → low blast radius.
- `runtime.SaveFileDialog` not testable headless → factor pure `BuildCSV`/`writeCSV` for unit tests; dialog wrapper smoke-tested.
- Dark mode coverage depends on tokens being CSS vars; audit hardcoded colors during build.
- Search must match displayed values (dates, FK names) not raw fields.
- `/all` export hidden in v1 (Go export strictly per-category); search still works there.

## Testing Strategy

- Go: `BuildCSV` table-driven per category (in-memory SQLite); `writeCSV` temp-file content assertion; bridge nil-db guard.
- Frontend: `matchesQuery` unit tests; `ui` store toggle/persist; `DataTable` density class; `AlertsModal` fixture render + empty-vs-nonempty auto-open.

## Spec Patches

None — existing delta specs (asset-alerts, data-export, app-preferences, asset-search, database-management) already cover the confirmed scenarios. `/all` export-hidden is an implementation choice, not a spec change.
