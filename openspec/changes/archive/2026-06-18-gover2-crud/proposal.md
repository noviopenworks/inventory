## Why

`gover2-readonly` displays all 7 asset categories but is fully read-only. Users cannot add new records, edit existing ones, or delete stale entries. This change delivers full CRUD for every category, making the Go app functionally equivalent to the PyQt6 app for day-to-day data management.

## What Changes

- `internal/services`: `Insert*`, `Update*`, `Delete*` functions for all 7 categories; `ListUsersForDropdown`, `ListDevicesForDropdown` helper queries
- `internal/bridge`: 21 new bridge methods (`AddComputer`, `UpdateComputer`, `DeleteComputer` × 7 categories) + 2 dropdown helpers
- `src/lib/api/index.ts`: 23 new api wrappers calling the bridge
- `src/components/EditPanel.vue`: 300 px slide-in right panel (add / edit modes), closes on Save / Cancel / Escape
- `src/components/ConfirmDialog.vue`: confirmation modal used before hard-delete
- `src/features/*/`: each view wires row-click → edit panel, "Add" button → create panel, delete button → confirm → delete
- Tests: Go table-driven tests for all Insert/Update/Delete functions; Vitest tests for EditPanel and ConfirmDialog

## Capabilities

### New Capabilities

- `gover2-crud`: Full add / edit / delete for all 7 asset categories via a slide-in 300 px panel. Row selection opens the panel in edit mode; "Add" button opens it in create mode. Hard delete is guarded by a confirmation dialog.

## Impact

- Unblocks: Phase 4 (alerts modal), Phase 5 (CSV export)
- Depends on: `gover2-readonly` complete (bridge + api layer already in place)
- Does not change the SQLite schema — same tables, same columns

## Non-Goals

- No bulk operations (multi-select delete)
- No inline cell editing
- No alerts modal (Phase 4)
- No CSV export (Phase 5)
- No packaging
- No undo / redo
