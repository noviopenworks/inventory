---
comet_change: gover2-crud
role: technical-design
canonical_spec: openspec
archived-with: 2026-06-18-gover2-crud
status: final
---

# gover2-crud — Technical Design

## Overview

Adds full Create / Read / Update / Delete for all 7 asset categories to the `gover2` Wails app. The existing read-only layer (services, bridge, api, views) is extended in-place; no new packages or architectural patterns are introduced.

## Architecture

Dependency chain is unchanged from Phase 2:

```
models  ←  database
  ↑            ↑
services    (sql.DB)
  ↑
bridge  ←  config
  ↑
api/index.ts
  ↑
EditPanel.vue / ConfirmDialog.vue
  ↑
*View.vue  (7 views)
```

Nothing imports upward. `services` knows only `models` and `database/sql`; `bridge` knows `services` and `config`.

## Data Flow

### Update

```
User clicks row
      │
      ▼
*View.vue: selectedRow → panelOpen=true, panelMode='edit'
      │  (EditPanel pre-fills fields from row data)
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
      │  UPDATE computers SET …, updated_at=datetime('now','localtime') WHERE id=?
      ▼
bridge returns nil → api resolves → view calls listComputers() → rows updated → panel closes
```

### Create

Identical with `INSERT`; panel closes after row appears in refreshed list.

### Delete

```
User clicks Delete button
      │
      ▼
confirmOpen = true → ConfirmDialog rendered
      │
User clicks Confirm
      │
      ▼
api.deleteComputer(id) → bridge → DELETE FROM computers WHERE id=? → view refetches
```

## Go — Models

New types added to `internal/models/models.go` alongside existing read structs:

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
// SmartphoneInput, TabletInput, WindowsKeyInput, AntivirusInput,
// OtherSoftwareInput, UserInput — same pattern, fields match DB columns

type DropdownItem struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type DeviceDropdownItem struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Kind string `json:"kind"` // "computer" | "smartphone" | "tablet"
}
```

`updated_at` is never in Input structs — always set by SQL.

## Go — Services

Three write functions per category (Computer shown; pattern × 7):

```go
func InsertComputer(db *sql.DB, d models.ComputerInput) (int64, error)
func UpdateComputer(db *sql.DB, id int, d models.ComputerInput) error
func DeleteComputer(db *sql.DB, id int) error
```

New file `internal/services/dropdowns.go`:

```go
func ListUsersForDropdown(db *sql.DB) ([]models.DropdownItem, error)
    // SELECT id, name FROM users ORDER BY name

func ListDevicesForDropdown(db *sql.DB) ([]models.DeviceDropdownItem, error)
    // SELECT id, name, 'computer' AS kind FROM computers
    // UNION ALL SELECT id, name, 'smartphone' FROM smartphones
    // UNION ALL SELECT id, name, 'tablet'     FROM tablets
    // ORDER BY kind, name
```

## Go — Bridge

23 new methods on `*App` in `internal/bridge/bridge.go` (Computer shown; × 7):

```go
func (a *App) AddComputer(data models.ComputerInput) error
func (a *App) UpdateComputer(id int, data models.ComputerInput) error
func (a *App) DeleteComputer(id int) error

func (a *App) ListUsersForDropdown() ([]models.DropdownItem, error)
func (a *App) ListDevicesForDropdown() ([]models.DeviceDropdownItem, error)
```

All methods open with `if a.db == nil { return errNoDB }`. After a successful write the bridge returns and the frontend calls the relevant `List*` method to refresh.

## Frontend — api/index.ts

23 new typed wrappers following the existing `call<T>` pattern:

```ts
export const addComputer = (data: ComputerInput) =>
  call<void>('AddComputer', data)
export const updateComputer = (id: number, data: ComputerInput) =>
  call<void>('UpdateComputer', id, data)
export const deleteComputer = (id: number) =>
  call<void>('DeleteComputer', id)

export const listUsersForDropdown = () =>
  call<DropdownItem[]>('ListUsersForDropdown')
export const listDevicesForDropdown = () =>
  call<DeviceDropdownItem[]>('ListDevicesForDropdown')
// … × 7 categories
```

New local types added alongside existing ones: `ComputerInput`, `SmartphoneInput`, etc., `DropdownItem`, `DeviceDropdownItem`.

## Frontend — EditPanel.vue

Replaces the existing stub. Props and emits:

```ts
defineProps<{
  category: string          // 'computers' | 'smartphones' | 'tablets' | 'windowskeys'
                            // | 'antivirus' | 'othersoftware' | 'users'
  mode: 'create' | 'edit'
  row: Record<string, unknown> | null
  users: DropdownItem[]
  devices: DeviceDropdownItem[]
}>()
defineEmits<{ saved: []; cancelled: [] }>()
```

**Layout**: 300 px right-edge panel (`w-80`), slides in via `transform: translateX(0)` (200 ms ease-in-out). Overlay covers the rest of the viewport to capture Escape.

**Field rendering**: A `FIELD_CONFIGS` constant maps each category to its ordered field list. Each entry has `{ key, label, type, required? }` where `type` is one of `text | textarea | select-status | select-user | select-device | date`. The panel iterates and renders the appropriate control using `FormField.vue` as the wrapper.

**Status options**:
- Devices (computers, smartphones, tablets): `active | inactive | repair | decommissioned`
- Antivirus / Other Software: `active | inactive | expired`
- Windows Keys: `active | inactive | expired`
- Users: `active | inactive`

**Device FK**: A single `<select>` built from `devices` prop. Each `<option>` value is `"kind:id"` (e.g. `"computer:3"`). On Save the panel splits the value and writes the correct FK key in the payload (`computerId`, `smartphoneId`, or `tabletId`); the other two are set to `null`.

**Validation (on Save)**:
1. Required field check: `name` for all categories except Windows Keys (required: `license_key`). Inline error below the field.
2. Date format check: if a date field is non-empty, must match `/^\d{4}-\d{2}-\d{2}$/`. Inline error below the field.
3. If any error: do not call api, keep panel open.

**Escape key**: `keydown` listener on the panel overlay emits `cancelled`.

## Frontend — ConfirmDialog.vue

New component. Props: `{ message: string }`. Emits: `confirmed`, `cancelled`.

Renders as a full-viewport overlay (`fixed inset-0 bg-black/50`) with a centered card containing the message, a "Delete" button (emits `confirmed`) and a "Cancel" button (emits `cancelled`).

## Frontend — View Wiring

Pattern shown for `ComputersView.vue`; identical for the other 6 views:

```vue
<script setup lang="ts">
const rows = ref<Computer[]>([])
const loading = ref(true)
const panelOpen = ref(false)
const panelMode = ref<'create' | 'edit'>('create')
const selectedRow = ref<Record<string, unknown> | null>(null)
const confirmOpen = ref(false)
const users = ref<DropdownItem[]>([])
const devices = ref<DeviceDropdownItem[]>([])

onMounted(async () => {
  ;[rows.value, users.value, devices.value] = await Promise.all([
    listComputers(), listUsersForDropdown(), listDevicesForDropdown(),
  ])
  loading.value = false
})

function onRowClick(row: Computer) {
  selectedRow.value = row as Record<string, unknown>
  panelMode.value = 'edit'
  panelOpen.value = true
}
function onAdd() { selectedRow.value = null; panelMode.value = 'create'; panelOpen.value = true }
async function onSaved() { panelOpen.value = false; rows.value = await listComputers() }
function onDeleteClick(row: Computer) { selectedRow.value = row as Record<string, unknown>; confirmOpen.value = true }
async function onDeleteConfirmed() {
  confirmOpen.value = false
  await deleteComputer(selectedRow.value!.id as number)
  rows.value = await listComputers()
}
</script>

<template>
  <div class="p-4 text-text-primary">
    <div class="flex justify-between items-center mb-4">
      <h1 class="text-lg font-semibold">Computers</h1>
      <button @click="onAdd">Add</button>
    </div>
    <DataTable :columns="columns" :rows="rows" :loading="loading"
               @row-click="onRowClick" @delete="onDeleteClick" />
    <EditPanel v-if="panelOpen" category="computers" :mode="panelMode"
               :row="selectedRow" :users="users" :devices="devices"
               @saved="onSaved" @cancelled="panelOpen = false" />
    <ConfirmDialog v-if="confirmOpen" message="Delete this computer?"
                   @confirmed="onDeleteConfirmed" @cancelled="confirmOpen = false" />
  </div>
</template>
```

`DataTable.vue` gains two new emits: `row-click` (when a row body is clicked) and `delete` (when a per-row Delete button is clicked). The existing `columns` + `rows` props are unchanged.

## Tests

### Go

**`internal/services/services_test.go`** (table-driven, in-memory SQLite via `database.Open(":memory:")`):

For each category: insert a record → verify it appears in `List*`; update it → verify field changed; delete it → verify `List*` is empty. Dropdown helpers: `ListDevicesForDropdown` returns merged results from computers/smartphones/tablets.

Coverage target: ≥ 70% for `internal/services`.

**`internal/bridge/bridge_test.go`**:

`AddComputer` round-trip (uses `NewAppWithDB`); `DeleteComputer` removes record; `ListUsersForDropdown` after insert; nil-db guard returns `errNoDB`.

Coverage target: ≥ 50% for `internal/bridge`.

### Frontend

**`EditPanel.test.ts`** (Vitest + @vue/test-utils):

- Edit mode with fixture Computer row → `name` input value equals fixture name
- Create mode → `name` input is empty
- Click Cancel → `cancelled` emitted
- Dispatch Escape keydown → `cancelled` emitted

**`ConfirmDialog.test.ts`**:

- Click "Delete" button → `confirmed` emitted
- Click "Cancel" button → `cancelled` emitted

## Error Handling

- Bridge nil-db guard: all methods return `errNoDB` if `a.db == nil`; api.ts propagates the rejected promise; views catch and can show an error toast (out of scope for Phase 3)
- SQL errors propagate from services → bridge → api → view; unhandled in Phase 3 (console.error)
- Validation errors are inline, client-side only; no server-side re-validation

## Non-Goals (Phase 3)

- Bulk operations
- Inline cell editing
- Alerts modal (Phase 4)
- CSV export (Phase 5)
- Undo / redo
- Error toast UI
