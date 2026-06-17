---
comet_change: gover2-readonly
role: technical-design
canonical_spec: openspec
---

# Technical Design: gover2-readonly

## Context

Builds directly on `gover2-scaffold`. All Go stubs return nil/empty. This change replaces every stub with a real implementation, layer by layer. The goal: open `inventory.db` (the existing Python app's SQLite file) and display all 7 asset categories in read-only data tables.

## Implementation Order

```
database → config → models → services → bridge → frontend
```

Each layer imports only from layers below it (Go dependency rule unchanged).

## Bridge Startup Pattern

`bridge.App` becomes stateful using the Wails `OnStartup` lifecycle callback:

```go
type App struct {
    db  *sql.DB
    cfg models.AppConfig
    ctx context.Context
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    cfg, err := config.Load()
    if err == nil {
        a.cfg = cfg
    }
    if a.cfg.DBPath != "" {
        if db, err := database.Open(a.cfg.DBPath); err == nil {
            database.InitSchema(db)
            a.db = db
        }
    }
}
```

`main.go` wires `OnStartup: app.startup`. If `a.db == nil` (path missing, file not found, or open failed), every List bridge method returns `(nil, errors.New("no database open"))`. Vue shows the DataTable empty state — no crash.

## Database Layer

`internal/database/database.go`:

```go
func Open(path string) (*sql.DB, error) {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        return nil, err
    }
    if _, err = db.Exec("PRAGMA journal_mode=WAL"); err != nil {
        db.Close()
        return nil, err
    }
    return db, nil
}

func InitSchema(db *sql.DB) error {
    _, err := db.Exec(`
        PRAGMA foreign_keys = ON;
        CREATE TABLE IF NOT EXISTS _meta (...);
        CREATE TABLE IF NOT EXISTS users (...);
        CREATE TABLE IF NOT EXISTS computers (...);
        -- ... all 7 tables verbatim from db/schema.py
    `)
    return err
}
```

**WAL mode note**: Opening an existing Python-app DB with WAL changes the journal mode permanently. This is safe because both apps never run simultaneously; WAL is strictly better for concurrent reads anyway.

## Models Layer

`internal/models/models.go` — complete struct definitions with `json:` tags for Wails serialization:

| Table | Key fields |
|---|---|
| `computers`, `smartphones`, `tablets` | id, name, model, userId\*, status, purchaseDate\*, warrantyExpiry\*, notes\*, createdAt, updatedAt |
| `windows_keys` | id, licenseKey\*, computerId\*, status, notes\*, createdAt, updatedAt |
| `antivirus`, `other_software` | id, name\*, licenseKey\*, computerId\*, smartphoneId\*, tabletId\*, status, expiryDate\*, notes\*, createdAt, updatedAt |
| `users` | id, name, surname\*, status, notes\*, createdAt, updatedAt |
| `Alert` | id, category, name, expiryDate, status, daysUntilExpiry |

`*` = nullable (`*string` or `*int` in Go, `string | null` or `number | null` in generated TS).

Wails regenerates `wailsjs/go/models.ts` on `wails build` — api wrapper imports from there.

## Config Layer

`internal/config/config.go`:

```go
func Load() (models.AppConfig, error) {
    path := configPath() // $XDG_CONFIG_HOME/inventory/config.json
    // if file missing: return defaults (DBPath="~/.local/share/inventory/inventory.db", density="comfortable", darkMode=false, expiryWarningDays=30)
}

func Save(cfg models.AppConfig) error {
    // write JSON atomically (write to .tmp, rename)
}
```

## Services Layer

`internal/services/services.go` — 7 `List*` functions + `GetAlerts`:

```go
func ListComputers(db *sql.DB) ([]models.Computer, error) {
    rows, err := db.Query(
        "SELECT id, name, model, user_id, status, purchase_date, warranty_expiry, notes, created_at, updated_at FROM computers ORDER BY name")
    if err != nil { return nil, err }
    defer rows.Close()
    var out []models.Computer
    for rows.Next() {
        var c models.Computer
        rows.Scan(&c.ID, &c.Name, &c.Model, &c.UserID, &c.Status, &c.PurchaseDate, &c.WarrantyExpiry, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
        out = append(out, c)
    }
    return out, rows.Err()
}

func GetAlerts(db *sql.DB) ([]models.Alert, error) {
    // SELECT from antivirus UNION SELECT from other_software
    // WHERE expiry_date IS NOT NULL AND date(expiry_date) <= date('now', '+30 days')
    // ORDER BY expiry_date
}
```

`ListWindowsKeys` selects from `windows_keys`. `ListAntivirus`/`ListOtherSoftware` selects from their respective tables. `ListUsers` selects from `users`.

## Bridge Layer

```go
func (a *App) ListComputers() ([]models.Computer, error) {
    if a.db == nil { return nil, errors.New("no database open") }
    return services.ListComputers(a.db)
}
// same pattern for all 7 List methods + GetAlerts

func (a *App) GetDatabasePath() (string, error) { return a.cfg.DBPath, nil }
func (a *App) GetConfig() (models.AppConfig, error) { return a.cfg, nil }
func (a *App) SetConfig(cfg models.AppConfig) error {
    a.cfg = cfg
    return config.Save(cfg)
}
func (a *App) NewDatabase(path string) error {
    db, err := database.Open(path)
    if err != nil { return err }
    database.InitSchema(db)
    if a.db != nil { a.db.Close() }
    a.db = db
    a.cfg.DBPath = path
    return config.Save(a.cfg)
}
func (a *App) OpenDatabase(path string) error {
    if _, err := os.Stat(path); err != nil { return fmt.Errorf("file not found: %s", path) }
    return a.NewDatabase(path)
}
func (a *App) ExportCSV(category, destPath string) error { return nil } // Phase 5
```

## Frontend

### api/index.ts

Replaces all `Promise.resolve([])` stubs with real Wails bridge calls:

```ts
import { ListComputers, ListSmartphones, ListTablets, ListWindowsKeys,
         ListAntivirus, ListOtherSoftware, ListUsers, GetAlerts,
         GetConfig, SetConfig, GetDatabasePath, OpenDatabase, NewDatabase } from '../../wailsjs/go/bridge/App'
import type { models } from '../../wailsjs/go/models'

export const listComputers = (): Promise<models.Computer[]> => ListComputers().catch(() => [])
// same pattern for all 7 list functions + getAlerts
export const getConfig = (): Promise<models.AppConfig> => GetConfig().catch(() => defaultConfig)
export const setConfig = (cfg: models.AppConfig): Promise<void> => SetConfig(cfg)
export const getDatabasePath = (): Promise<string> => GetDatabasePath().catch(() => '')
export const openDatabase = (path: string): Promise<void> => OpenDatabase(path)
export const newDatabase = (path: string): Promise<void> => NewDatabase(path)
```

### DataTable.vue

```ts
defineProps<{
  columns: { key: string; label: string }[]
  rows: Record<string, unknown>[]
  loading: boolean
}>()
```

Renders `<thead>` from `columns`, `<tbody>` from `rows` (access `row[col.key]`). Shows empty-state `<td colspan>` when `rows.length === 0 && !loading`.

### Feature Views

Each view calls its api function on `onMounted`, stores result in `ref<T[]>`, passes to `<DataTable>`. Column definitions are typed constants per view file:

```ts
// ComputersView.vue
const columns = [
  { key: 'name', label: 'Name' },
  { key: 'model', label: 'Model' },
  { key: 'userId', label: 'User' },
  { key: 'status', label: 'Status' },
  { key: 'purchaseDate', label: 'Purchase Date' },
  { key: 'warrantyExpiry', label: 'Warranty Expiry' },
]
```

`AllAssetsView.vue` uses `Promise.all([listComputers(), listSmartphones(), listTablets()])`, tags each entry with `{ category: 'Computer' | 'Smartphone' | 'Tablet' }`, passes the union to DataTable with Category column first.

`StatusBadge.vue` uses the computed class map already implemented in scaffold — no changes needed.

### Vitest Setup

`package.json` gains: `"vitest": "^2.0"`, `"@vue/test-utils": "^2"`, `"jsdom": "^25"`
`package.json` scripts: `"test": "vitest run"`
`vite.config.ts` gains:
```ts
test: {
  environment: 'jsdom',
  globals: true,
}
```

## Test Strategy

### Go

**`internal/database/database_test.go`**:
- Open `:memory:` DB, call `InitSchema`, verify 7 tables exist (query `sqlite_master`)

**`internal/services/*_test.go`** (table-driven):
- Open `:memory:` + `InitSchema`, insert 2–3 fixture rows via `db.Exec`, call `List*`, assert length and field values
- `GetAlerts`: expiry = today+15 days → in result; today+45 days → not in result
- Target: ≥70% line coverage

**`internal/bridge/bridge_test.go`**:
- Create `App` with in-memory DB, call bridge methods, verify correct data returned
- Call bridge methods with `nil` db → error returned
- Target: ≥50% line coverage

### Frontend

**`DataTable.test.ts`**:
- Mount with 3 columns + 3 rows → assert 3 `<th>`, 3 `<tr>` in tbody
- Mount with `rows=[]` → assert empty-state `<td>` present
- Mount with `loading=true` → assert loading indicator visible

**`StatusBadge.test.ts`**:
- For each of 7 status values → assert correct `bg-s-*` class in rendered HTML

## Non-Goals

- No add/edit/delete (Phase 3)
- No edit panel slide-in (Phase 3)
- No alerts modal (Phase 4)
- No CSV export (Phase 5)
- No packaging
- ExportCSV stub returns nil — real implementation in Phase 5
