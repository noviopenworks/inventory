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

// Dropdown helpers:
func (a *App) ListUsersForDropdown() ([]models.DropdownItem, error)
func (a *App) ListDevicesForDropdown() ([]models.DeviceDropdownItem, error)
```

All methods guard with `if a.db == nil { return errNoDB }`.

After a successful write, the bridge does **not** refetch — the frontend calls the relevant `List*` method itself to refresh the table.

## Frontend Components

### EditPanel.vue

Props:
```ts
defineProps<{
  category: string          // 'computers' | 'smartphones' | ...
  mode: 'create' | 'edit'
  row: Record<string, unknown> | null   // null in create mode
  users: DropdownItem[]
  devices: DeviceDropdownItem[]
}>()

defineEmits<{
  saved: []
  cancelled: []
}>()
```

Behavior:
- 300 px right-edge slide-in (CSS `transform: translateX`)
- Transition: 200 ms ease-in-out
- Escape key closes (calls `cancelled`)
- Fields per category are driven by a `FIELD_CONFIGS` constant map (mirrors Python's TABLE_CONFIG)
- Status dropdown: device statuses vs user statuses depending on category
- Device FK: single "Device" dropdown combining computers + smartphones + tablets with `"Name (Type)"` labels; selection resolves to the correct FK field on save
- License key: text input with "N/A" preset button
- Date fields: text input with `YYYY-MM-DD` placeholder, validated on save
- Required field: `name` for all categories except Windows Keys (required: `license_key`)
- Save calls the appropriate api function, emits `saved` on success

### ConfirmDialog.vue

Props: `{ message: string }`. Emits: `confirmed`, `cancelled`.
Rendered as a centered modal overlay.

## View Wiring (pattern, shown for ComputersView)

```vue
<EditPanel
  v-if="panelOpen"
  :category="'computers'"
  :mode="panelMode"
  :row="selectedRow"
  :users="users"
  :devices="devices"
  @saved="onSaved"
  @cancelled="panelOpen = false"
/>
<ConfirmDialog
  v-if="confirmOpen"
  message="Delete this computer?"
  @confirmed="onDeleteConfirmed"
  @cancelled="confirmOpen = false"
/>
```

`onSaved`: closes panel, calls `listComputers()`, updates `rows`.
Delete button: sets `confirmOpen = true`. `onDeleteConfirmed`: calls `deleteComputer(selectedRow.id)`, refetches.

## Tests

### Go (table-driven, in-memory SQLite)

- `services_test.go`: for each category — insert a record, verify it appears in List; update it, verify field changed; delete it, verify List is empty
- Coverage target: ≥ 70% for services package (same as Phase 2)
- Bridge tests: `AddComputer` round-trip, `DeleteComputer` removes record, nil-db guard returns error

### Frontend (Vitest + Vue Test Utils)

- `EditPanel.test.ts`: mounts in edit mode with fixture row → fields pre-filled; mounts in create mode → fields empty; clicking Cancel emits `cancelled`; Escape key emits `cancelled`
- `ConfirmDialog.test.ts`: clicking confirm emits `confirmed`; clicking cancel emits `cancelled`
