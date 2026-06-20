## Why

`gover2-scaffold` delivers a compiling Wails skeleton with empty data. This change wires real data: implements the SQLite database layer, models, and services for all 7 asset categories, exposes them through the Wails bridge, and connects the Vue frontend data tables to display actual rows from the existing `inventory.db` database.

After this change, the Gover v2 app replaces the read-only behavior of the current PyQt6 app for all 7 categories. Users can open the app, navigate the sidebar, and see their real asset data.

## What Changes

- `internal/database`: real SQLite connection via `modernc.org/sqlite`, WAL pragma, schema init DDL verbatim from `db/schema.py`
- `internal/models`: complete field mapping (snake_case column → CamelCase Go field, with JSON tags for Wails serialization)
- `internal/services`: `List*` functions for all 7 categories with real SQL SELECT queries; alert calculation (assets expiring within 30 days)
- `internal/bridge`: `List*` bridge methods return real data; `GetDatabasePath`, `GetConfig`, `SetConfig` wired
- `internal/config`: real JSON persistence at `$XDG_CONFIG_HOME/inventory/config.json`; loads existing `inventory.db` path from config or defaults to the Python app's DB location
- `frontend/src/lib/api/index.ts`: updated to call real Wails bridge methods instead of returning empty arrays
- `frontend/src/features/*/`: DataTable component wired with real columns, loading states, and error handling per category
- `frontend/src/components/StatusBadge.vue`: real status color mapping using design tokens
- Tests: Go table-driven tests for `internal/services` List functions ≥70% coverage; Vitest component tests for DataTable

## Capabilities

### New Capabilities

- `gover2-readonly`: The Gover v2 app opens the existing SQLite inventory database and displays all 7 asset categories (computers, smartphones, tablets, windows keys, antivirus, other software, users) in read-only data tables with status badges and correct column layout.

## Impact

- Unblocks: Phase 3 (CRUD add/edit/delete)
- Depends on: `gover2-scaffold` complete
- Does not affect: current Python app, the Python app's database file is read (not written)

## Non-Goals

- No add/edit/delete (Phase 3)
- No slide-in edit panel (Phase 3)
- No alerts modal (Phase 4)
- No CSV export or backup (Phase 5)
- No packaging
- No native menu bar (Phase 3)
