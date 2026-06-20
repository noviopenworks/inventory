# Architecture Notes

This document captures the resolved architecture for gover (the Wails + Go +
Vue + TypeScript rewrite of IT Asset Inventory). Decisions were locked in
`openspec/changes/archive/2026-06-17-gover-phase1-architecture` and
`2026-06-17-gover-ui-spec`; this file is the canonical summary.

## Runtime Boundaries

- **Go owns** database access, migrations, validation that protects data integrity, file system paths, backup, restore, import, export, and platform-specific behavior.
- **Vue owns** layout, navigation, presentation state, forms, table interactions, filters, and user-facing feedback.
- **TypeScript API clients** in `frontend/src/lib/api/` wrap Wails-generated bindings so Vue components never depend on low-level bridge details.

## Go Packages

Six packages under `internal/`. Dependency order is locked; `bridge` is the
sole package that imports from all others.

| Package | Responsibility |
|---------|----------------|
| `internal/database` | Connection management, SQLite pragmas, transactions, migrations (7-table schema reused verbatim from `db/schema.py`). |
| `internal/models` | Typed records for users, devices, licenses, software, alerts, and settings. |
| `internal/services` | Application operations: list assets, create/update/delete records, CSV export, alert checks, dropdown sources. |
| `internal/bridge` | Wails-exposed methods called by the frontend. Sole package that imports from all others. |
| `internal/config` | App paths, database location, preferences, version metadata. |
| `internal/backup` | Database backup and restore utilities. |

SQLite driver: `modernc.org/sqlite` (pure Go, no CGo).

## Frontend Layout

Feature-folder structure adopted.

| Folder | Responsibility |
|--------|----------------|
| `frontend/src/features/assets` | Asset tables, forms, filters, category views. |
| `frontend/src/features/users` | User list and user forms. |
| `frontend/src/features/alerts` | Expiry and warranty alert views. |
| `frontend/src/features/settings` | Appearance, database path, export, backup, restore. |
| `frontend/src/layouts` | Main shell, navigation, top bar, responsive structure. |
| `frontend/src/components` | Reusable buttons, fields, modals, panels, tables, badges, empty states. |
| `frontend/src/lib/api` | TypeScript wrappers around Wails bindings. |
| `frontend/src/lib/types` | Frontend-facing TypeScript types. |
| `frontend/src/stores` | Pinia stores (`assets.ts`, `ui.ts`). |

Routing: hash-mode Vue Router with 9 routes (one per category plus settings).

## Data Strategy

**Reuse current schema directly** (Option 1 from the original evaluation).

- Pros: fastest path, less migration risk, both apps can read the same `.db` file.
- Cons: carries current modeling limitations forward.
- Mitigation: Phase 5 will document the upgrade path; no schema migration is required for users switching from PyQt6 → gover.

The Go DDL in `internal/database/database.go` carries the 7-table schema
verbatim from `db/schema.py`.

## UI Strategy

**Dashboard plus sidebar** (Option 2 from the original evaluation), with
modifications.

- Persistent left sidebar, 200 px wide, theme-aware (`--c-sidebar*` tokens).
- Computers is the first screen (dashboard deferred — not yet built).
- Slide-in right panel for add/edit, 300 px wide.
- Comfortable / Compact density toggle (33 px / 24 px row height).
- Manual dark mode toggle (Light default). System-follow is deferred.
- 7 fixed categories (matches PyQt6 app).
- Native OS menu bar via `wails/v2/pkg/menu`.

Visual reference: [`UI_DESIGN_REFERENCE.html`](UI_DESIGN_REFERENCE.html).

## Risk Areas

- Data migration from the current SQLite database — **mitigated** by reusing the schema verbatim; risk is now "does gover read the same `.db` correctly", validated by Phase 2.
- Recreating current table and form behavior without losing small usability details — ongoing through Phase 6 regression testing.
- Keeping Wails bindings clean instead of exposing database details directly to Vue — enforced by the `lib/api/` wrapper layer.
- Avoiding a frontend that looks modern but makes bulk data work slower — open concern for Phase 6.
- Packaging differences between PyInstaller and Wails — open concern for Phase 5.
