# Comet Design Handoff

- Change: gover2-readonly
- Phase: design
- Mode: compact
- Context hash: 7469218f6ac88e6ec991822c8aeab41f21eab0695b79b8e8757fb7ec09a829ab

Generated-by: comet-handoff.sh

OpenSpec remains the canonical capability spec. This handoff is a deterministic, source-traceable context pack, not an agent-authored summary.

## openspec/changes/gover2-readonly/proposal.md

- Source: openspec/changes/gover2-readonly/proposal.md
- Lines: 1-38
- SHA256: 903fa07256c16fc51dc3ba760d8ded79cf339955e0776761792d59b9f1a93069

```md
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
```

## openspec/changes/gover2-readonly/design.md

- Source: openspec/changes/gover2-readonly/design.md
- Lines: 1-123
- SHA256: 4f4490c8d6684e23cd3457c299d30b5ab3a5767b2c1c340038d35eb2274e783f

[TRUNCATED]

```md
## Approach

Replace each stub from `gover2-scaffold` with real implementations, working layer by layer: database → models → services → bridge → frontend. Follow the Go dependency rule: no layer imports above itself in the chain.

## Data Flow

```
SQLite file (existing inventory.db)
        │
        ▼
internal/database.Open()  ←── config.Load() resolves path
        │
        ▼
internal/services.ListComputers(db)  ← SQL: SELECT * FROM computers ORDER BY name
        │
        ▼
internal/bridge.ListComputers()  ← delegates to services, returns []models.Computer
        │   (Wails serializes to JSON via struct tags)
        ▼
src/lib/api/index.ts: listComputers()  ← window.go.bridge.App.ListComputers()
        │
        ▼
stores/assets.ts: items = await api.listComputers()
        │
        ▼
features/assets/ComputersView.vue
        │
        ▼
components/DataTable.vue  ← rows prop + columns config
```

## Database Layer

`internal/database/database.go`:
```go
func Open(path string) (*sql.DB, error) {
    db, err := sql.Open("sqlite", path)
    if err != nil { return nil, err }
    _, err = db.Exec("PRAGMA journal_mode=WAL")
    return db, err
}

func InitSchema(db *sql.DB) error {
    // DDL verbatim from db/schema.py — CREATE TABLE IF NOT EXISTS for all 7 tables
}
```

The schema init uses `IF NOT EXISTS` so opening an existing database is safe.

## Models Layer

`internal/models/models.go` — one struct per table with `json:` tags matching camelCase field names (Wails serializes Go → TS using these):

```go
type Computer struct {
    ID             int     `json:"id"`
    Name           string  `json:"name"`
    Model          string  `json:"model"`
    UserID         *int    `json:"userId"`
    Status         string  `json:"status"`
    PurchaseDate   *string `json:"purchaseDate"`
    WarrantyExpiry *string `json:"warrantyExpiry"`
    Notes          *string `json:"notes"`
    CreatedAt      string  `json:"createdAt"`
    UpdatedAt      string  `json:"updatedAt"`
}
// Smartphone, Tablet: same shape
// WindowsKey: id, licenseKey, computerId, status, notes, createdAt, updatedAt
// Antivirus, OtherSoftware: id, name, licenseKey, computerId, smartphoneId, tabletId, status, expiryDate, notes, createdAt, updatedAt
// User: id, name, surname, status, notes, createdAt, updatedAt
```

## Services Layer

Each List function runs a simple SELECT and scans into the model slice:

```go
func ListComputers(db *sql.DB) ([]models.Computer, error) {
    rows, err := db.Query("SELECT id, name, model, user_id, status, purchase_date, warranty_expiry, notes, created_at, updated_at FROM computers ORDER BY name")
    // ... scan loop
```

Full source: openspec/changes/gover2-readonly/design.md

## openspec/changes/gover2-readonly/tasks.md

- Source: openspec/changes/gover2-readonly/tasks.md
- Lines: 1-69
- SHA256: e8627388bf1b24c07a1ec61a839e9e999951e40f769019b7f084baf1260564cf

```md
## 1. Database Package

- [ ] 1.1 Implement `internal/database/database.go`: `Open(path string) (*sql.DB, error)` with WAL pragma
- [ ] 1.2 Implement `InitSchema(db *sql.DB) error` — DDL verbatim from `db/schema.py` (7 tables + _meta, all IF NOT EXISTS)
- [ ] 1.3 Write `internal/database/database_test.go`: open in-memory DB, verify schema initialises without error
- [ ] 1.4 Run `go test ./internal/database/...` — must pass

## 2. Models Package

- [ ] 2.1 Implement all 7 model structs with correct JSON tags in `internal/models/models.go`
- [ ] 2.2 Implement `AppConfig` in `internal/config/config.go` with `Load() / Save()` using `$XDG_CONFIG_HOME`
- [ ] 2.3 Write `internal/models/models_test.go`: JSON marshal/unmarshal round-trip for each struct
- [ ] 2.4 Run `go test ./internal/models/...` — must pass

## 3. Services Package — List Operations

- [ ] 3.1 Implement `ListComputers(db *sql.DB) ([]models.Computer, error)` in `internal/services/`
- [ ] 3.2 Implement `ListSmartphones`, `ListTablets` (same pattern as computers)
- [ ] 3.3 Implement `ListWindowsKeys(db *sql.DB) ([]models.WindowsKey, error)`
- [ ] 3.4 Implement `ListAntivirus`, `ListOtherSoftware` (same pattern)
- [ ] 3.5 Implement `ListUsers(db *sql.DB) ([]models.User, error)`
- [ ] 3.6 Implement `GetAlerts(db *sql.DB) ([]models.Alert, error)` — expiry within 30 days from antivirus + other_software
- [ ] 3.7 Write table-driven tests for all 7 List functions using in-memory SQLite with fixture data
- [ ] 3.8 Run `go test -cover ./internal/services/...` — must hit ≥70% line coverage

## 4. Bridge Package — Real Methods

- [ ] 4.1 Update `internal/bridge/bridge.go`: App struct holds `*sql.DB` and `*config.AppConfig`
- [ ] 4.2 Implement real `ListComputers()` (and all 6 other List methods) delegating to services
- [ ] 4.3 Implement `GetDatabasePath()` returning config.DBPath
- [ ] 4.4 Implement `GetConfig()` and `SetConfig()` using config package
- [ ] 4.5 Implement `NewDatabase(path string)` — updates config, re-opens DB
- [ ] 4.6 Implement `OpenDatabase(path string)` — same as NewDatabase but validates file exists first
- [ ] 4.7 Implement `ExportCSV(category, destPath string)` — stub returning nil for now (full implementation Phase 5)
- [ ] 4.8 Update `app.go`/`main.go` to wire DB open → bridge init → wails.Run
- [ ] 4.9 Write `internal/bridge/bridge_test.go`: bridge methods return correct data via test DB
- [ ] 4.10 Run `go test -cover ./internal/bridge/...` — must hit ≥50% line coverage

## 5. Vue Data Tables

- [ ] 5.1 Update `src/lib/api/index.ts` to call real Wails-generated bridge bindings
- [ ] 5.2 Implement `ComputersView.vue`: columns = [Name, Model, User, Status, Purchase Date, Warranty Expiry], loads via `api.listComputers()`
- [ ] 5.3 Implement `SmartphonesView.vue`, `TabletsView.vue` (same columns as computers)
- [ ] 5.4 Implement `WindowsKeysView.vue`: columns = [License Key, Computer, Status, Notes]
- [ ] 5.5 Implement `AntivirusView.vue`: columns = [Name, License Key, Computer, Status, Expiry Date]
- [ ] 5.6 Implement `OtherSoftwareView.vue`: same columns as Antivirus
- [ ] 5.7 Implement `UsersView.vue`: columns = [Name, Surname, Status, Notes]
- [ ] 5.8 Implement `AllAssetsView.vue`: union of computers + smartphones + tablets with Category column
- [ ] 5.9 Implement `DataTable.vue`: props = { columns: ColDef[], rows: Row[], loading: boolean }; renders table with header, rows, empty-state
- [ ] 5.10 Implement `StatusBadge.vue`: applies correct `text-s-*` + `bg-*` classes from Tailwind tokens per status value

## 6. Frontend Tests

- [ ] 6.1 Write `DataTable.test.ts`: renders columns, renders 3 fixture rows, shows empty-state div when rows=[]
- [ ] 6.2 Write `StatusBadge.test.ts`: each status value → correct CSS class applied
- [ ] 6.3 Run `pnpm --prefix frontend run test` — all tests pass

## 7. TypeScript & Build Verification

- [ ] 7.1 Run `pnpm --prefix frontend run typecheck` — exit 0, no errors
- [ ] 7.2 Run `go build ./...` — exit 0, no import cycles
- [ ] 7.3 Run `wails build` — exit 0, binary produced

## 8. Integration Smoke Test

- [ ] 8.1 Run `wails dev` — app opens, sidebar shows 7 categories
- [ ] 8.2 Navigate to Computers — real rows from `inventory.db` appear in data table with correct columns
- [ ] 8.3 Navigate to each of the other 6 categories — rows appear (or empty table if no data, no crash)
- [ ] 8.4 StatusBadge shows correct color for each status value on real data
```
