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
