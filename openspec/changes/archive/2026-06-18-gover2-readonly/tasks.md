## 1. Database Package

- [x] 1.1 Implement `internal/database/database.go`: `Open(path string) (*sql.DB, error)` with WAL pragma
- [x] 1.2 Implement `InitSchema(db *sql.DB) error` — DDL verbatim from `db/schema.py` (7 tables + _meta, all IF NOT EXISTS)
- [x] 1.3 Write `internal/database/database_test.go`: open in-memory DB, verify schema initialises without error
- [x] 1.4 Run `go test ./internal/database/...` — must pass

## 2. Models Package

- [x] 2.1 Implement all 7 model structs with correct JSON tags in `internal/models/models.go`
- [x] 2.2 Implement `AppConfig` in `internal/config/config.go` with `Load() / Save()` using `$XDG_CONFIG_HOME`
- [x] 2.3 Write `internal/models/models_test.go`: JSON marshal/unmarshal round-trip for each struct
- [x] 2.4 Run `go test ./internal/models/...` — must pass

## 3. Services Package — List Operations

- [x] 3.1 Implement `ListComputers(db *sql.DB) ([]models.Computer, error)` in `internal/services/`
- [x] 3.2 Implement `ListSmartphones`, `ListTablets` (same pattern as computers)
- [x] 3.3 Implement `ListWindowsKeys(db *sql.DB) ([]models.WindowsKey, error)`
- [x] 3.4 Implement `ListAntivirus`, `ListOtherSoftware` (same pattern)
- [x] 3.5 Implement `ListUsers(db *sql.DB) ([]models.User, error)`
- [x] 3.6 Implement `GetAlerts(db *sql.DB, warningDays int) ([]models.Alert, error)` — expiry within warningDays from antivirus + other_software
- [x] 3.7 Write table-driven tests for all 7 List functions using in-memory SQLite with fixture data
- [x] 3.8 Run `go test -cover ./internal/services/...` — must hit ≥70% line coverage

## 4. Bridge Package — Real Methods

- [x] 4.1 Update `internal/bridge/bridge.go`: App struct holds `*sql.DB` and `*config.AppConfig`
- [x] 4.2 Implement real `ListComputers()` (and all 6 other List methods) delegating to services
- [x] 4.3 Implement `GetDatabasePath()` returning config.DBPath
- [x] 4.4 Implement `GetConfig()` and `SetConfig()` using config package
- [x] 4.5 Implement `NewDatabase(path string)` — updates config, re-opens DB
- [x] 4.6 Implement `OpenDatabase(path string)` — same as NewDatabase but validates file exists first
- [x] 4.7 Implement `ExportCSV(category, destPath string)` — stub returning nil for now (full implementation Phase 5)
- [x] 4.8 Update `app.go`/`main.go` to wire DB open → bridge init → wails.Run
- [x] 4.9 Write `internal/bridge/bridge_test.go`: bridge methods return correct data via test DB
- [x] 4.10 Run `go test -cover ./internal/bridge/...` — must hit ≥50% line coverage

## 5. Vue Data Tables

- [x] 5.1 Update `src/lib/api/index.ts` to call real Wails-generated bridge bindings
- [x] 5.2 Implement `ComputersView.vue`: columns = [Name, Model, User, Status, Purchase Date, Warranty Expiry], loads via `api.listComputers()`
- [x] 5.3 Implement `SmartphonesView.vue`, `TabletsView.vue` (same columns as computers)
- [x] 5.4 Implement `WindowsKeysView.vue`: columns = [License Key, Computer, Status, Notes]
- [x] 5.5 Implement `AntivirusView.vue`: columns = [Name, License Key, Computer, Status, Expiry Date]
- [x] 5.6 Implement `OtherSoftwareView.vue`: same columns as Antivirus
- [x] 5.7 Implement `UsersView.vue`: columns = [Name, Surname, Status, Notes]
- [x] 5.8 Implement `AllAssetsView.vue`: union of computers + smartphones + tablets with Category column
- [x] 5.9 Implement `DataTable.vue`: props = { columns: ColDef[], rows: Row[], loading: boolean }; renders table with header, rows, empty-state
- [x] 5.10 Implement `StatusBadge.vue`: applies correct `text-s-*` + `bg-*` classes from Tailwind tokens per status value

## 6. Frontend Tests

- [x] 6.1 Write `DataTable.test.ts`: renders columns, renders 3 fixture rows, shows empty-state div when rows=[]
- [x] 6.2 Write `StatusBadge.test.ts`: each status value → correct CSS class applied
- [x] 6.3 Run `pnpm --prefix frontend run test` — all tests pass

## 7. TypeScript & Build Verification

- [x] 7.1 Run `pnpm --prefix frontend run typecheck` — exit 0, no errors
- [x] 7.2 Run `go build ./...` — exit 0, no import cycles
- [x] 7.3 Run `wails build` — exit 0, binary produced

## 8. Integration Smoke Test

- [x] 8.1 Run `wails dev` — app opens, sidebar shows 7 categories
- [x] 8.2 Navigate to Computers — real rows from `inventory.db` appear in data table with correct columns
- [x] 8.3 Navigate to each of the other 6 categories — rows appear (or empty table if no data, no crash)
- [x] 8.4 StatusBadge shows correct color for each status value on real data
