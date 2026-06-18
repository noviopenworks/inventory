# Comet Design Handoff

- Change: gover2-crud
- Phase: design
- Mode: compact
- Context hash: a8dc9194dda259bf71054cf4f66dd3c5f37f053bf046b4334f13dc3e13368b94

Generated-by: comet-handoff.sh

OpenSpec remains the canonical capability spec. This handoff is a deterministic, source-traceable context pack, not an agent-authored summary.

## openspec/changes/gover2-crud/proposal.md

- Source: openspec/changes/gover2-crud/proposal.md
- Lines: 1-34
- SHA256: 9ab3207d77a675caf3c1e29a07d0eb1c1bdbb860666f7def8ba7ac03cc313a86

```md
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
```

## openspec/changes/gover2-crud/design.md

- Source: openspec/changes/gover2-crud/design.md
- Lines: 1-163
- SHA256: ebdc3a99c5c65dd5841bdcda71fd210639f788497470de538b88d404cd93878c

[TRUNCATED]

```md
## Approach

Add write operations layer-by-layer following the same dependency chain as Phase 2:
`services → bridge → api → components → views`.

The Go dependency rule is unchanged: `services` knows only `models` and `database/sql`; `bridge` knows `services` and `config`; nothing imports upward.

## Data Flow — Edit (Update)

```
User clicks row
      │
      ▼
ComputersView.vue sets selectedRow → EditPanel opens (edit mode)
      │  (fields pre-filled from row data)
      ▼
User edits → clicks Save
      │
      ▼
api.updateComputer(id, payload)
      │  window.go.bridge.App.UpdateComputer(id, payload)
      ▼
bridge.UpdateComputer(id, data models.ComputerInput) error
      │  nil-db guard
      ▼
services.UpdateComputer(db, id, data)
      │  UPDATE computers SET … WHERE id = ?
      ▼
bridge returns nil → api resolves → view refetches list → panel closes
```

Create flow is identical with `INSERT` and panel closes after the new row is appended.
Delete flow: `ConfirmDialog` → confirmed → `api.deleteComputer(id)` → bridge → `DELETE` → view refetches.

## Services Layer

New functions per category (Computer shown; pattern repeats for Smartphone, Tablet, WindowsKey, Antivirus, OtherSoftware, User):

```go
func InsertComputer(db *sql.DB, d models.ComputerInput) (int64, error)
func UpdateComputer(db *sql.DB, id int, d models.ComputerInput) error
func DeleteComputer(db *sql.DB, id int) error
```

Input structs live in `internal/models/models.go` alongside existing read structs:

```go
type ComputerInput struct {
    Name           string  `json:"name"`
    Model          string  `json:"model"`
    UserID         *int    `json:"userId"`
    Status         string  `json:"status"`
    PurchaseDate   *string `json:"purchaseDate"`
    WarrantyExpiry *string `json:"warrantyExpiry"`
    Notes          *string `json:"notes"`
}
```

Dropdown helpers:
```go
func ListUsersForDropdown(db *sql.DB) ([]models.DropdownItem, error)
    // SELECT id, name FROM users ORDER BY name

func ListDevicesForDropdown(db *sql.DB) ([]models.DeviceDropdownItem, error)
    // SELECT id, name, 'computer' as kind FROM computers
    // UNION ALL ... smartphones ... tablets ORDER BY kind, name
```

`models.DropdownItem` = `{ ID int; Name string }`.
`models.DeviceDropdownItem` = `{ ID int; Name string; Kind string }` where Kind ∈ `"computer" | "smartphone" | "tablet"`.

`updated_at` is set via `datetime('now','localtime')` in the SQL, never from the client.

## Bridge Layer

```go
// Per-category (Computer shown):
func (a *App) AddComputer(data models.ComputerInput) error
func (a *App) UpdateComputer(id int, data models.ComputerInput) error
func (a *App) DeleteComputer(id int) error
```

Full source: openspec/changes/gover2-crud/design.md

## openspec/changes/gover2-crud/tasks.md

- Source: openspec/changes/gover2-crud/tasks.md
- Lines: 1-81
- SHA256: 3a54f72fc7407be3e374062bce6589444a131235ada1a1ebb9db9f803e09d369

[TRUNCATED]

```md
## 1. Models — Input Structs + Dropdown Types

- [ ] 1.1 Add `ComputerInput`, `SmartphoneInput`, `TabletInput` structs to `internal/models/models.go`
- [ ] 1.2 Add `WindowsKeyInput`, `AntivirusInput`, `OtherSoftwareInput`, `UserInput` structs
- [ ] 1.3 Add `DropdownItem` and `DeviceDropdownItem` structs

## 2. Services — Write Operations

- [ ] 2.1 Implement `InsertComputer`, `UpdateComputer`, `DeleteComputer` in `internal/services/`
- [ ] 2.2 Implement `InsertSmartphone`, `UpdateSmartphone`, `DeleteSmartphone`
- [ ] 2.3 Implement `InsertTablet`, `UpdateTablet`, `DeleteTablet`
- [ ] 2.4 Implement `InsertWindowsKey`, `UpdateWindowsKey`, `DeleteWindowsKey`
- [ ] 2.5 Implement `InsertAntivirus`, `UpdateAntivirus`, `DeleteAntivirus`
- [ ] 2.6 Implement `InsertOtherSoftware`, `UpdateOtherSoftware`, `DeleteOtherSoftware`
- [ ] 2.7 Implement `InsertUser`, `UpdateUser`, `DeleteUser`
- [ ] 2.8 Implement `ListUsersForDropdown` and `ListDevicesForDropdown`
- [ ] 2.9 Write table-driven tests: insert→list, update→verify field, delete→list empty (all 7 categories)
- [ ] 2.10 Run `go test -cover ./internal/services/...` — must stay ≥ 70%

## 3. Bridge — Write Bridge Methods

- [ ] 3.1 Add `AddComputer`, `UpdateComputer`, `DeleteComputer` to `internal/bridge/bridge.go`
- [ ] 3.2 Add `AddSmartphone`, `UpdateSmartphone`, `DeleteSmartphone`
- [ ] 3.3 Add `AddTablet`, `UpdateTablet`, `DeleteTablet`
- [ ] 3.4 Add `AddWindowsKey`, `UpdateWindowsKey`, `DeleteWindowsKey`
- [ ] 3.5 Add `AddAntivirus`, `UpdateAntivirus`, `DeleteAntivirus`
- [ ] 3.6 Add `AddOtherSoftware`, `UpdateOtherSoftware`, `DeleteOtherSoftware`
- [ ] 3.7 Add `AddUser`, `UpdateUser`, `DeleteUser`
- [ ] 3.8 Add `ListUsersForDropdown` and `ListDevicesForDropdown` bridge methods
- [ ] 3.9 Write bridge tests: round-trip Add→List, Delete removes record, nil-db guard
- [ ] 3.10 Run `go test -cover ./internal/bridge/...` — must stay ≥ 50%

## 4. Go Build Verification

- [ ] 4.1 Run `go build ./...` — exit 0
- [ ] 4.2 Run `wails build` — exit 0, binary produced, wailsjs regenerated

## 5. Frontend — api/index.ts

- [ ] 5.1 Add `addComputer`, `updateComputer`, `deleteComputer` to `src/lib/api/index.ts`
- [ ] 5.2 Add same for Smartphone, Tablet, WindowsKey, Antivirus, OtherSoftware, User (6 × 3 = 18 functions)
- [ ] 5.3 Add `listUsersForDropdown` and `listDevicesForDropdown`

## 6. Frontend — EditPanel.vue

- [ ] 6.1 Create `src/components/EditPanel.vue` with slide-in animation (300 px, right edge)
- [ ] 6.2 Implement `FIELD_CONFIGS` map: fields per category with labels, types, required flags
- [ ] 6.3 Implement user dropdown (blank + all users from `listUsersForDropdown`)
- [ ] 6.4 Implement device dropdown (mixed, `"Name (Type)"` labels from `listDevicesForDropdown`)
- [ ] 6.5 Implement status dropdown (device statuses vs user statuses)
- [ ] 6.6 Implement date field with `YYYY-MM-DD` placeholder + save-time format validation
- [ ] 6.7 Implement required-field validation (inline error message on Save)
- [ ] 6.8 Implement Escape key → close panel
- [ ] 6.9 Write `EditPanel.test.ts`: edit mode pre-fills fields; create mode fields empty; Cancel emits cancelled; Escape emits cancelled

## 7. Frontend — ConfirmDialog.vue

- [ ] 7.1 Create `src/components/ConfirmDialog.vue` (centered overlay modal)
- [ ] 7.2 Write `ConfirmDialog.test.ts`: confirm emits confirmed; cancel emits cancelled

## 8. Frontend — Wire All 7 Views

- [ ] 8.1 Update `ComputersView.vue`: row click → edit panel; Add button → create panel; Delete → confirm → delete
- [ ] 8.2 Update `SmartphonesView.vue` (same pattern)
- [ ] 8.3 Update `TabletsView.vue` (same pattern)
- [ ] 8.4 Update `WindowsKeysView.vue` (same pattern)
- [ ] 8.5 Update `AntivirusView.vue` (same pattern)
- [ ] 8.6 Update `OtherSoftwareView.vue` (same pattern)
- [ ] 8.7 Update `UsersView.vue` (same pattern, user status options)

## 9. Frontend Verification

- [ ] 9.1 Run `pnpm run typecheck` — exit 0
- [ ] 9.2 Run `pnpm run test` — all tests pass

## 10. Integration Smoke Test

- [ ] 10.1 Add a new Computer record — appears in table
- [ ] 10.2 Edit an existing record — changes reflected in table
- [ ] 10.3 Delete a record — row removed
```

Full source: openspec/changes/gover2-crud/tasks.md
