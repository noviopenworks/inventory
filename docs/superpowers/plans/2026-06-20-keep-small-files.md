# Keep-Small-Files Architecture Pass — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Split the codebase's >200-line files into small, per-entity files within the existing layered packages, with zero behavior change.

**Architecture:** Six mechanical tasks on `feat/go-default-cutover`: split `models`, `services` (source + tests), `bridge` (source + tests), and `EditPanel.vue`. Tasks 1–5 move symbols verbatim between files in the same package; only Task 6 edits code (parametrizing extracted functions). The existing test suite is the regression net — it must pass unchanged after every task.

**Tech Stack:** Go 1.25 (local 1.26), Wails v2.12, golangci-lint v2.12.2, Vue 3 + TypeScript, Vitest, Task.

## Global Constraints

- **Pure refactor — no behavior change, no logic edits** (except Task 6's mechanical parametrization of moved functions). Moved Go funcs/types and test funcs are copied **verbatim**.
- Module `inventory`; build tag `webkit2_41`; CGO for the bridge package. `golangci-lint` is at `$(go env GOPATH)/bin` — prefix lint with `export PATH="$(go env GOPATH)/bin:$PATH"`.
- **Public API unchanged:** exported Go identifiers keep their names/signatures; bridge public method signatures unchanged so `frontend/wailsjs/` bindings do NOT change; `EditPanel.vue` props/emits and rendered `data-testid`s unchanged.
- Per-task gates: `task lint` → `0 issues.`; `CGO_ENABLED=1 go test -tags webkit2_41 ./...` all pass; for Task 6 also `pnpm --prefix frontend run typecheck` + Vitest. The full `task check` and `task test:cov` (≥75%) must still pass.
- Target: no in-scope source file over ~200 lines after its split.
- Commit with `git commit --no-verify`.
- Each new file declares the same `package` as the file it came from and imports only what it uses (`go build` errors guide imports; run `gofmt -w` on touched dirs).

---

### Task 1: Split `internal/models/models.go` (200) into per-entity files

**Files:**
- Create: `internal/models/{computer,smartphone,tablet,windowskey,antivirus,othersoftware,user,shared}.go`
- Delete: `internal/models/models.go`
- Unchanged: `internal/models/models_test.go` (101 lines, stays)

**Interfaces:**
- Produces: all the same exported types in package `models`, just relocated. No signature changes.

- [ ] **Step 1: Create per-entity model files (move type declarations verbatim)**

Each file starts with `package models` and the file's structs copied verbatim from `models.go`:

| New file | Types to move (verbatim) |
|----------|--------------------------|
| `computer.go` | `Computer`, `ComputerInput` |
| `smartphone.go` | `Smartphone`, `SmartphoneInput` |
| `tablet.go` | `Tablet`, `TabletInput` |
| `windowskey.go` | `WindowsKey`, `WindowsKeyInput` |
| `antivirus.go` | `Antivirus`, `AntivirusInput` |
| `othersoftware.go` | `OtherSoftware`, `OtherSoftwareInput` |
| `user.go` | `User`, `UserInput` |
| `shared.go` | `AppConfig`, `Alert`, `DropdownItem`, `DeviceDropdownItem` |

If `models.go` has a package doc comment (`// Package models ...`), keep it on exactly one file — put it on `shared.go`.

- [ ] **Step 2: Delete the emptied source file**

```bash
git rm internal/models/models.go
```

- [ ] **Step 3: Build, vet, lint**

Run: `export PATH="$(go env GOPATH)/bin:$PATH" && CGO_ENABLED=1 go build -tags webkit2_41 ./... && task lint`
Expected: build succeeds, lint `0 issues.` (structs are pure data — no import needed beyond what each had).

- [ ] **Step 4: Run the model tests**

Run: `CGO_ENABLED=1 go test -tags webkit2_41 ./internal/models/...`
Expected: `ok inventory/internal/models`.

- [ ] **Step 5: Commit**

```bash
git add -A && git commit --no-verify -m "refactor(models): split models.go into per-entity files"
```

---

### Task 2: Split `internal/services` source (`services.go` 221 + `write.go` 238) by entity

**Files:**
- Create: `internal/services/{computer,smartphone,tablet,windowskey,antivirus,othersoftware,user,alerts}.go`
- Delete: `internal/services/services.go`, `internal/services/write.go`
- Unchanged: `internal/services/{export.go, dropdowns.go}`

**Interfaces:**
- Consumes: package `models` (Task 1).
- Produces: all the same exported funcs in package `services`, relocated. Signatures unchanged: `ListComputers(db) ([]models.Computer, error)`, `InsertComputer(db, models.ComputerInput) (int64, error)`, `UpdateComputer(db, int, models.ComputerInput) error`, `DeleteComputer(db, int) error`, … per entity; `GetAlerts(db, int) ([]models.Alert, error)`.

- [ ] **Step 1: Create per-entity service files (move funcs verbatim, read+write together)**

Each file: `package services` + the entity's read func (from `services.go`) and its Insert/Update/Delete funcs (from `write.go`), copied verbatim. Add imports `"database/sql"` and `"inventory/internal/models"` (and anything else the moved funcs reference — `go build` will flag exact needs).

| New file | Funcs to move (verbatim) |
|----------|--------------------------|
| `computer.go` | `ListComputers`, `InsertComputer`, `UpdateComputer`, `DeleteComputer` |
| `smartphone.go` | `ListSmartphones`, `InsertSmartphone`, `UpdateSmartphone`, `DeleteSmartphone` |
| `tablet.go` | `ListTablets`, `InsertTablet`, `UpdateTablet`, `DeleteTablet` |
| `windowskey.go` | `ListWindowsKeys`, `InsertWindowsKey`, `UpdateWindowsKey`, `DeleteWindowsKey` |
| `antivirus.go` | `ListAntivirus`, `InsertAntivirus`, `UpdateAntivirus`, `DeleteAntivirus` |
| `othersoftware.go` | `ListOtherSoftware`, `InsertOtherSoftware`, `UpdateOtherSoftware`, `DeleteOtherSoftware` |
| `user.go` | `ListUsers`, `InsertUser`, `UpdateUser`, `DeleteUser` |
| `alerts.go` | `GetAlerts` |

Preserve any package doc comment on one file (`alerts.go`). `services.go`/`write.go` contain no shared helpers (verified — only the listed funcs), so nothing else needs relocating.

- [ ] **Step 2: Delete the emptied source files**

```bash
git rm internal/services/services.go internal/services/write.go
```

- [ ] **Step 3: Format, build, lint**

Run: `export PATH="$(go env GOPATH)/bin:$PATH" && gofmt -w internal/services && CGO_ENABLED=1 go build -tags webkit2_41 ./... && task lint`
Expected: build + lint clean. Fix any "imported and not used" / "undefined" by adjusting each new file's import block.

- [ ] **Step 4: Run services + bridge tests (bridge consumes services)**

Run: `CGO_ENABLED=1 go test -tags webkit2_41 ./internal/services/... ./internal/bridge/...`
Expected: both `ok`.

- [ ] **Step 5: Commit**

```bash
git add -A && git commit --no-verify -m "refactor(services): split services.go/write.go into per-entity files"
```

---

### Task 3: Split `internal/services/services_test.go` (516) by entity

**Files:**
- Create: `internal/services/{helpers_test.go, computer_test.go, smartphone_test.go, tablet_test.go, windowskey_test.go, antivirus_test.go, othersoftware_test.go, user_test.go, alerts_test.go, dropdowns_test.go}`
- Delete: `internal/services/services_test.go`
- Unchanged: `internal/services/export_test.go`

**Interfaces:**
- Consumes: the `services` package funcs. All test funcs move **verbatim**; only their file location changes.

- [ ] **Step 1: Create the shared test helper file**

`helpers_test.go`: `package services_test` + the `openTestDB(t *testing.T) *sql.DB` helper moved verbatim, with its imports (`database/sql`, `testing`, `require`, `inventory/internal/database`).

- [ ] **Step 2: Create per-concern test files (move test funcs verbatim)**

Each file `package services_test`, funcs copied verbatim, imports as each needs:

| New file | Test funcs to move |
|----------|--------------------|
| `computer_test.go` | `TestListComputers_Empty`, `TestListComputers_WithData`, `TestComputer_InsertUpdateDelete` |
| `smartphone_test.go` | `TestListSmartphones_WithData`, `TestSmartphone_InsertUpdateDelete` |
| `tablet_test.go` | `TestListTablets_WithData`, `TestTablet_InsertUpdateDelete` |
| `windowskey_test.go` | `TestListWindowsKeys_WithData`, `TestWindowsKey_InsertUpdateDelete` |
| `antivirus_test.go` | `TestListAntivirus_WithData`, `TestAntivirus_InsertUpdateDelete` |
| `othersoftware_test.go` | `TestListOtherSoftware_WithData`, `TestOtherSoftware_InsertUpdateDelete` |
| `user_test.go` | `TestListUsers_WithData`, `TestUser_InsertUpdateDelete` |
| `alerts_test.go` | the 7 `TestGetAlerts_*` funcs (`_ExpiringWithin30Days`, `_NotExpiringSoon`, `_AlreadyExpired`, `_ComputerWarrantyExpiring`, `_SmartphoneWarrantyExpired`, `_TabletWarrantyCovered`, `_RetiredExcluded`) |
| `dropdowns_test.go` | `TestListUsersForDropdown`, `TestListDevicesForDropdown` |

- [ ] **Step 3: Delete the emptied test file**

```bash
git rm internal/services/services_test.go
```

- [ ] **Step 4: Format, lint, run services tests**

Run: `export PATH="$(go env GOPATH)/bin:$PATH" && gofmt -w internal/services && task lint && CGO_ENABLED=1 go test -tags webkit2_41 ./internal/services/...`
Expected: lint `0 issues.`; services tests `ok` with the same total test count as before.

- [ ] **Step 5: Commit**

```bash
git add -A && git commit --no-verify -m "test(services): split services_test.go into per-entity test files"
```

---

### Task 4: Split `internal/bridge/bridge.go` (386) by entity + concern

**Files:**
- Create: `internal/bridge/{app.go, computer.go, smartphone.go, tablet.go, windowskey.go, antivirus.go, othersoftware.go, user.go, database.go, queries.go, export.go}`
- Delete: `internal/bridge/bridge.go`

**Interfaces:**
- Consumes: `services`, `models`, `config`, `database`, `backup`, Wails `runtime`.
- Produces: identical `bridge.App` with identical public method signatures (so `frontend/wailsjs/` bindings stay byte-identical). `ErrNoDB`, `list[T]`, `withDB` remain in the package.

- [ ] **Step 1: Create `app.go` (package core)**

`package bridge` + move verbatim: the `ErrNoDB` var (with its doc comment), the `App` struct, `NewApp`, `NewAppWithDB`, `Startup`, and the `list[T any](...)` func + `withDB` method. Keep the `// Package bridge ...` doc comment here if present. Imports: `context`, `database/sql`, `errors`, `config`, `database`, `models`.

- [ ] **Step 2: Create per-entity method files (move methods verbatim)**

Each file `package bridge`, the entity's four methods moved verbatim, imports `database/sql`, `inventory/internal/models`, `inventory/internal/services`:

| New file | Methods to move |
|----------|-----------------|
| `computer.go` | `ListComputers`, `AddComputer`, `UpdateComputer`, `DeleteComputer` |
| `smartphone.go` | `ListSmartphones`, `AddSmartphone`, `UpdateSmartphone`, `DeleteSmartphone` |
| `tablet.go` | `ListTablets`, `AddTablet`, `UpdateTablet`, `DeleteTablet` |
| `windowskey.go` | `ListWindowsKeys`, `AddWindowsKey`, `UpdateWindowsKey`, `DeleteWindowsKey` |
| `antivirus.go` | `ListAntivirus`, `AddAntivirus`, `UpdateAntivirus`, `DeleteAntivirus` |
| `othersoftware.go` | `ListOtherSoftware`, `AddOtherSoftware`, `UpdateOtherSoftware`, `DeleteOtherSoftware` |
| `user.go` | `ListUsers`, `AddUser`, `UpdateUser`, `DeleteUser` |

- [ ] **Step 3: Create `database.go`, `queries.go`, `export.go` (concern files)**

`package bridge`, move verbatim:
- `database.go`: `GetDatabasePath`, `GetConfig`, `SetConfig`, `NewDatabase`, `OpenDatabase`, `NewDatabaseDialog`, `OpenDatabaseDialog`, `BackupDatabase`, `BackupDatabaseDialog`. Imports: `database/sql`, `fmt`, `os`, `time`, `runtime`, `config`, `database`, `backup`, `models`.
- `queries.go`: `GetAlerts`, `ListUsersForDropdown`, `ListDevicesForDropdown`. Imports: `models`, `services`.
- `export.go`: `ExportCSV`. Imports: `runtime`, `services`.

- [ ] **Step 4: Delete the emptied source file**

```bash
git rm internal/bridge/bridge.go
```

- [ ] **Step 5: Format, build, lint**

Run: `export PATH="$(go env GOPATH)/bin:$PATH" && gofmt -w internal/bridge && CGO_ENABLED=1 go build -tags webkit2_41 ./... && task lint`
Expected: build + lint clean. Adjust per-file imports until clean.

- [ ] **Step 6: Run bridge tests AND confirm bindings unchanged**

Run: `CGO_ENABLED=1 go test -tags webkit2_41 ./internal/bridge/... && wails build -tags webkit2_41 >/dev/null && git status --porcelain frontend/wailsjs/`
Expected: bridge tests `ok`; `git status --porcelain frontend/wailsjs/` prints nothing (public signatures unchanged → bindings identical).

- [ ] **Step 7: Commit**

```bash
git add -A && git commit --no-verify -m "refactor(bridge): split bridge.go into per-entity + concern files"
```

---

### Task 5: Split `internal/bridge/bridge_test.go` (585) by entity + concern

**Files:**
- Create: `internal/bridge/{helpers_test.go, computer_test.go, smartphone_test.go, tablet_test.go, windowskey_test.go, antivirus_test.go, othersoftware_test.go, user_test.go, dropdowns_test.go, alerts_test.go, config_test.go, database_test.go, export_test.go}`
- Delete: `internal/bridge/bridge_test.go`

**Interfaces:**
- All test funcs move **verbatim**; only file location changes. `newBridgeWithDB` becomes the shared helper.

- [ ] **Step 1: Create `helpers_test.go`**

`package bridge_test` + `newBridgeWithDB(t *testing.T) *bridge.App` moved verbatim, with its imports.

- [ ] **Step 2: Create per-entity test files (move verbatim)**

Each `package bridge_test`. Group each entity's list + write tests:

| New file | Test funcs to move |
|----------|--------------------|
| `computer_test.go` | `TestListComputers_NoDB`, `TestListComputers_WithDB`, `TestAddComputer_NoDB`, `TestAddComputer_RoundTrip`, `TestUpdateComputer`, `TestDeleteComputer` |
| `smartphone_test.go` | `TestListSmartphones_WithDB`, `TestAddSmartphone_NoDB`, `TestAddSmartphone_RoundTrip`, `TestUpdateSmartphone`, `TestDeleteSmartphone` |
| `tablet_test.go` | `TestListTablets_WithDB`, `TestAddTablet_NoDB`, `TestAddTablet_RoundTrip`, `TestUpdateTablet`, `TestDeleteTablet` |
| `windowskey_test.go` | `TestListWindowsKeys_WithDB`, `TestAddWindowsKey_NoDB`, `TestAddWindowsKey_RoundTrip`, `TestUpdateWindowsKey`, `TestDeleteWindowsKey` |
| `antivirus_test.go` | `TestListAntivirus_WithDB`, `TestAddAntivirus_NoDB`, `TestAddAntivirus_RoundTrip`, `TestUpdateAntivirus`, `TestDeleteAntivirus` |
| `othersoftware_test.go` | `TestListOtherSoftware_WithDB`, `TestAddOtherSoftware_NoDB`, `TestAddOtherSoftware_RoundTrip`, `TestUpdateOtherSoftware`, `TestDeleteOtherSoftware` |
| `user_test.go` | `TestListUsers_WithDB`, `TestAddUser_NoDB`, `TestAddUser_RoundTrip`, `TestUpdateUser`, `TestDeleteUser` |

- [ ] **Step 3: Create concern test files (move verbatim)**

| New file | Test funcs to move |
|----------|--------------------|
| `dropdowns_test.go` | `TestListUsersForDropdown_NoDB`, `TestListDevicesForDropdown_NoDB`, `TestListUsersForDropdown`, `TestListDevicesForDropdown` |
| `alerts_test.go` | `TestApp_GetAlerts_NoDB`, `TestApp_GetAlerts_WithDB` |
| `config_test.go` | `TestGetDatabasePath_NoConfig`, `TestGetConfig_Default`, `TestApp_GetSetConfig_RoundTrip` |
| `database_test.go` | `TestApp_NewDatabase`, `TestApp_NewDatabase_BadPath`, `TestApp_OpenDatabase`, `TestApp_OpenDatabase_NotFound`, `TestApp_BackupDatabaseDialog_NoDB`, `TestApp_BackupDatabase_NoDB`, `TestApp_BackupDatabase_WritesCopy` |
| `export_test.go` | `TestExportCSV_NoDB` |

- [ ] **Step 4: Delete the emptied test file**

```bash
git rm internal/bridge/bridge_test.go
```

- [ ] **Step 5: Format, lint, run bridge tests**

Run: `export PATH="$(go env GOPATH)/bin:$PATH" && gofmt -w internal/bridge && task lint && CGO_ENABLED=1 go test -tags webkit2_41 ./internal/bridge/...`
Expected: lint `0 issues.`; bridge tests `ok` with the same total test count as before.

- [ ] **Step 6: Commit**

```bash
git add -A && git commit --no-verify -m "test(bridge): split bridge_test.go into per-entity + concern test files"
```

---

### Task 6: Split `frontend/src/components/EditPanel.vue` (460) into shell + 3 modules

**Files:**
- Create: `frontend/src/components/editpanel/{fieldConfigs.ts, payload.ts, save.ts}`
- Modify: `frontend/src/components/EditPanel.vue` (slim the `<script setup>`; template unchanged)
- Unchanged: `frontend/src/components/EditPanel.test.ts` (gate — must pass untouched)

**Interfaces:**
- Produces: `fieldConfigs.ts` exports `FieldConfig`, `FIELD_CONFIGS`, `STATUS_MAP`, `CATEGORY_LABELS`. `payload.ts` exports `initForm(category, row)`, `buildPayload(category, form)`. `save.ts` exports `callApi(category, isEdit, id, payload)`.

- [ ] **Step 1: Create `frontend/src/components/editpanel/fieldConfigs.ts`**

```ts
export interface FieldConfig {
  key: string
  label: string
  type: 'text' | 'textarea' | 'select-status' | 'select-user' | 'select-device' | 'select-computer' | 'date'
  required?: boolean
}

export const FIELD_CONFIGS: Record<string, FieldConfig[]> = {
  computers: [
    { key: 'name', label: 'Name', type: 'text', required: true },
    { key: 'model', label: 'Model', type: 'text' },
    { key: 'userId', label: 'User', type: 'select-user' },
    { key: 'status', label: 'Status', type: 'select-status' },
    { key: 'purchaseDate', label: 'Purchase Date', type: 'date' },
    { key: 'warrantyExpiry', label: 'Warranty Expiry', type: 'date' },
    { key: 'notes', label: 'Notes', type: 'textarea' },
  ],
  smartphones: [
    { key: 'name', label: 'Name', type: 'text', required: true },
    { key: 'model', label: 'Model', type: 'text' },
    { key: 'userId', label: 'User', type: 'select-user' },
    { key: 'status', label: 'Status', type: 'select-status' },
    { key: 'purchaseDate', label: 'Purchase Date', type: 'date' },
    { key: 'warrantyExpiry', label: 'Warranty Expiry', type: 'date' },
    { key: 'notes', label: 'Notes', type: 'textarea' },
  ],
  tablets: [
    { key: 'name', label: 'Name', type: 'text', required: true },
    { key: 'model', label: 'Model', type: 'text' },
    { key: 'userId', label: 'User', type: 'select-user' },
    { key: 'status', label: 'Status', type: 'select-status' },
    { key: 'purchaseDate', label: 'Purchase Date', type: 'date' },
    { key: 'warrantyExpiry', label: 'Warranty Expiry', type: 'date' },
    { key: 'notes', label: 'Notes', type: 'textarea' },
  ],
  windowskeys: [
    { key: 'licenseKey', label: 'License Key', type: 'text', required: true },
    { key: '_computerSelect', label: 'Computer', type: 'select-computer' },
    { key: 'status', label: 'Status', type: 'select-status' },
    { key: 'notes', label: 'Notes', type: 'textarea' },
  ],
  antivirus: [
    { key: 'name', label: 'Name', type: 'text', required: true },
    { key: 'licenseKey', label: 'License Key', type: 'text' },
    { key: '_deviceSelect', label: 'Device', type: 'select-device' },
    { key: 'status', label: 'Status', type: 'select-status' },
    { key: 'expiryDate', label: 'Expiry Date', type: 'date' },
    { key: 'notes', label: 'Notes', type: 'textarea' },
  ],
  othersoftware: [
    { key: 'name', label: 'Name', type: 'text', required: true },
    { key: 'licenseKey', label: 'License Key', type: 'text' },
    { key: '_deviceSelect', label: 'Device', type: 'select-device' },
    { key: 'status', label: 'Status', type: 'select-status' },
    { key: 'expiryDate', label: 'Expiry Date', type: 'date' },
    { key: 'notes', label: 'Notes', type: 'textarea' },
  ],
  users: [
    { key: 'name', label: 'Name', type: 'text', required: true },
    { key: 'surname', label: 'Surname', type: 'text' },
    { key: 'status', label: 'Status', type: 'select-status' },
    { key: 'notes', label: 'Notes', type: 'textarea' },
  ],
}

const DEVICE_STATUSES = ['active', 'inactive', 'repair', 'decommissioned']
const LICENSE_STATUSES = ['active', 'inactive', 'expired']
const USER_STATUSES = ['active', 'inactive']

export const STATUS_MAP: Record<string, string[]> = {
  computers: DEVICE_STATUSES,
  smartphones: DEVICE_STATUSES,
  tablets: DEVICE_STATUSES,
  windowskeys: LICENSE_STATUSES,
  antivirus: LICENSE_STATUSES,
  othersoftware: LICENSE_STATUSES,
  users: USER_STATUSES,
}

export const CATEGORY_LABELS: Record<string, string> = {
  computers: 'Computer',
  smartphones: 'Smartphone',
  tablets: 'Tablet',
  windowskeys: 'Windows Key',
  antivirus: 'Antivirus',
  othersoftware: 'Other Software',
  users: 'User',
}
```

- [ ] **Step 2: Create `frontend/src/components/editpanel/payload.ts`**

```ts
import { FIELD_CONFIGS, STATUS_MAP } from './fieldConfigs'

function deriveDeviceSelect(row: Record<string, unknown> | null): string {
  if (!row) return ''
  if (row.computerId != null) return `computer:${row.computerId}`
  if (row.smartphoneId != null) return `smartphone:${row.smartphoneId}`
  if (row.tabletId != null) return `tablet:${row.tabletId}`
  return ''
}

function deriveComputerSelect(row: Record<string, unknown> | null): string {
  if (!row || row.computerId == null) return ''
  return `computer:${row.computerId}`
}

function nullOrStr(val: unknown): string | null {
  if (val === '' || val === null || val === undefined) return null
  return String(val)
}

function nullOrNum(val: unknown): number | null {
  if (val === '' || val === null || val === undefined) return null
  const n = Number(val)
  return isNaN(n) ? null : n
}

export function initForm(category: string, row: Record<string, unknown> | null): Record<string, unknown> {
  const base: Record<string, unknown> = {}
  const fields = FIELD_CONFIGS[category] ?? []
  for (const field of fields) {
    if (field.key === '_deviceSelect') {
      base['_deviceSelect'] = deriveDeviceSelect(row)
    } else if (field.key === '_computerSelect') {
      base['_computerSelect'] = deriveComputerSelect(row)
    } else {
      base[field.key] = row != null ? (row[field.key] ?? '') : ''
    }
  }
  if (!base['status']) {
    const statuses = STATUS_MAP[category]
    base['status'] = statuses?.[0] ?? ''
  }
  return base
}

export function buildPayload(category: string, form: Record<string, unknown>): unknown {
  const cat = category

  if (cat === 'computers' || cat === 'smartphones' || cat === 'tablets') {
    return {
      name: String(form['name'] ?? ''),
      model: nullOrStr(form['model']) ?? '',
      userId: nullOrNum(form['userId']),
      status: String(form['status'] ?? ''),
      purchaseDate: nullOrStr(form['purchaseDate']),
      warrantyExpiry: nullOrStr(form['warrantyExpiry']),
      notes: nullOrStr(form['notes']),
    }
  }

  if (cat === 'windowskeys') {
    const computerSel = String(form['_computerSelect'] ?? '')
    const computerId = computerSel ? Number(computerSel.split(':')[1]) : null
    return {
      licenseKey: String(form['licenseKey'] ?? ''),
      computerId: isNaN(computerId as number) ? null : computerId,
      status: String(form['status'] ?? ''),
      notes: nullOrStr(form['notes']),
    }
  }

  if (cat === 'antivirus' || cat === 'othersoftware') {
    const deviceSel = String(form['_deviceSelect'] ?? '')
    let computerId: number | null = null
    let smartphoneId: number | null = null
    let tabletId: number | null = null
    if (deviceSel) {
      const [kind, idStr] = deviceSel.split(':')
      const id = Number(idStr)
      if (kind === 'computer') computerId = id
      else if (kind === 'smartphone') smartphoneId = id
      else if (kind === 'tablet') tabletId = id
    }
    return {
      name: String(form['name'] ?? ''),
      licenseKey: nullOrStr(form['licenseKey']) ?? '',
      computerId,
      smartphoneId,
      tabletId,
      status: String(form['status'] ?? ''),
      expiryDate: nullOrStr(form['expiryDate']),
      notes: nullOrStr(form['notes']),
    }
  }

  if (cat === 'users') {
    return {
      name: String(form['name'] ?? ''),
      surname: nullOrStr(form['surname']),
      status: String(form['status'] ?? ''),
      notes: nullOrStr(form['notes']),
    }
  }

  return {}
}
```

- [ ] **Step 3: Create `frontend/src/components/editpanel/save.ts`**

```ts
import {
  addComputer, updateComputer,
  addSmartphone, updateSmartphone,
  addTablet, updateTablet,
  addWindowsKey, updateWindowsKey,
  addAntivirus, updateAntivirus,
  addOtherSoftware, updateOtherSoftware,
  addUser, updateUser,
} from '@/lib/api'

export async function callApi(category: string, isEdit: boolean, id: number, payload: unknown): Promise<void> {
  const cat = category

  if (cat === 'computers') {
    return isEdit
      ? updateComputer(id, payload as Parameters<typeof updateComputer>[1])
      : addComputer(payload as Parameters<typeof addComputer>[0])
  }
  if (cat === 'smartphones') {
    return isEdit
      ? updateSmartphone(id, payload as Parameters<typeof updateSmartphone>[1])
      : addSmartphone(payload as Parameters<typeof addSmartphone>[0])
  }
  if (cat === 'tablets') {
    return isEdit
      ? updateTablet(id, payload as Parameters<typeof updateTablet>[1])
      : addTablet(payload as Parameters<typeof addTablet>[0])
  }
  if (cat === 'windowskeys') {
    return isEdit
      ? updateWindowsKey(id, payload as Parameters<typeof updateWindowsKey>[1])
      : addWindowsKey(payload as Parameters<typeof addWindowsKey>[0])
  }
  if (cat === 'antivirus') {
    return isEdit
      ? updateAntivirus(id, payload as Parameters<typeof updateAntivirus>[1])
      : addAntivirus(payload as Parameters<typeof addAntivirus>[0])
  }
  if (cat === 'othersoftware') {
    return isEdit
      ? updateOtherSoftware(id, payload as Parameters<typeof updateOtherSoftware>[1])
      : addOtherSoftware(payload as Parameters<typeof addOtherSoftware>[0])
  }
  if (cat === 'users') {
    return isEdit
      ? updateUser(id, payload as Parameters<typeof updateUser>[1])
      : addUser(payload as Parameters<typeof addUser>[0])
  }
}
```

- [ ] **Step 4: Replace `EditPanel.vue`'s `<script setup>` (keep the `<template>` block byte-for-byte)**

Replace everything between `<script setup lang="ts">` and `</script>` with:

```ts
import { reactive, computed, onMounted, onUnmounted } from 'vue'
import FormField from './FormField.vue'
import type { DropdownItem, DeviceDropdownItem } from '@/lib/types'
import { FIELD_CONFIGS, STATUS_MAP, CATEGORY_LABELS, type FieldConfig } from './editpanel/fieldConfigs'
import { initForm, buildPayload } from './editpanel/payload'
import { callApi } from './editpanel/save'

const props = defineProps<{
  category: string
  mode: 'create' | 'edit'
  row: Record<string, unknown> | null
  users: DropdownItem[]
  devices: DeviceDropdownItem[]
}>()

const emit = defineEmits<{ saved: []; cancelled: [] }>()

const currentFields = computed<FieldConfig[]>(() => FIELD_CONFIGS[props.category] ?? [])
const currentStatuses = computed<string[]>(() => STATUS_MAP[props.category] ?? [])
const categoryLabel = computed(() => CATEGORY_LABELS[props.category] ?? props.category)
const computerDevices = computed(() => props.devices.filter(d => d.kind === 'computer'))

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const form = reactive<Record<string, any>>(initForm(props.category, props.row))
const errors = reactive<Record<string, string>>({})

// Global Escape key handler (catches key events even when focus is on inputs)
function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    emit('cancelled')
  }
}

onMounted(() => {
  document.addEventListener('keydown', onKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', onKeyDown)
})

function validateForm(): boolean {
  for (const key of Object.keys(errors)) {
    delete errors[key]
  }

  const fields = currentFields.value
  let valid = true

  for (const field of fields) {
    const value = form[field.key]

    if (field.required) {
      if (value === '' || value === null || value === undefined) {
        errors[field.key] = 'This field is required'
        valid = false
      }
    }

    if (field.type === 'date' && value !== '' && value !== null && value !== undefined) {
      if (!/^\d{4}-\d{2}-\d{2}$/.test(String(value))) {
        errors[field.key] = 'Date must be in YYYY-MM-DD format'
        valid = false
      }
    }
  }

  return valid
}

async function onSave() {
  if (!validateForm()) return

  const payload = buildPayload(props.category, form)
  const id = props.mode === 'edit' ? (props.row!['id'] as number) : 0
  try {
    await callApi(props.category, props.mode === 'edit', id, payload)
    emit('saved')
  } catch (err) {
    console.error('[EditPanel] save failed:', err)
  }
}
```

(`currentStatuses`, `categoryLabel`, `currentFields`, `computerDevices`, `errors`, `form` are all still referenced by the unchanged template — do not remove them.)

- [ ] **Step 5: Typecheck and run the EditPanel tests**

Run: `pnpm --prefix frontend run typecheck && pnpm --prefix frontend run test -- EditPanel`
Expected: `vue-tsc` clean; `EditPanel.test.ts` passes unchanged (same DOM, same `data-testid`s, same save behavior).

- [ ] **Step 6: Run the full frontend suite**

Run: `pnpm --prefix frontend run test`
Expected: all suites pass (69 tests).

- [ ] **Step 7: Commit**

```bash
git add -A && git commit --no-verify -m "refactor(editpanel): extract fieldConfigs/payload/save modules, slim EditPanel.vue"
```

---

### Task 7: Final verification

- [ ] **Step 1: Full gate**

Run: `export PATH="$(go env GOPATH)/bin:$PATH" && task check && task test:cov`
Expected: lint `0 issues.`, typecheck clean, all Go + frontend tests pass, coverage ≥ 75%.

- [ ] **Step 2: Confirm no in-scope file still exceeds ~200 lines**

Run:
```bash
find internal frontend/src -name '*.go' -o -name '*.vue' -o -name '*.ts' | grep -v node_modules \
  | xargs wc -l 2>/dev/null | sort -rn | awk '$1 > 200 && $2 != "total"'
```
Expected: only files explicitly left out of scope (`bridge_test`/`services_test` are now gone; any remaining >200 should be a deliberately-excluded test file or barrel — confirm none of the six split sources reappear).

- [ ] **Step 3: Confirm bindings still unchanged on the final tree**

Run: `wails build -tags webkit2_41 >/dev/null && git status --porcelain frontend/wailsjs/`
Expected: empty.

---

## Self-Review

**Spec coverage** (design §3 → tasks): §3.1 models → Task 1 ✓ · §3.2 services source → Task 2 ✓ · §3.3 bridge → Task 4 ✓ · §3.4 EditPanel → Task 6 ✓ · §3.5 tests → Tasks 3 (services) + 5 (bridge) ✓. §4 not-split files are untouched (no task references them). §6 acceptance → Task 7 verification.

**Placeholder scan:** Tasks 1–5 are verbatim moves specified by exact symbol→file tables (the bodies live in the named source files; reproducing them would risk transcription drift — the instruction is "move verbatim, do not edit"). Task 6 contains the complete new module source. No "TBD"/"handle edge cases".

**Type/name consistency:** Extracted signatures are consistent across files: `initForm(category, row)` and `buildPayload(category, form)` (payload.ts) match their call sites in EditPanel.vue Step 4; `callApi(category, isEdit, id, payload)` matches `callApi(props.category, props.mode === 'edit', id, payload)`. `FieldConfig`/`FIELD_CONFIGS`/`STATUS_MAP`/`CATEGORY_LABELS` exported by fieldConfigs.ts and imported by both payload.ts and EditPanel.vue. Go: all moved identifiers keep their exact names (no renames), so cross-package references in `services`→`models` and `bridge`→`services` remain valid.
