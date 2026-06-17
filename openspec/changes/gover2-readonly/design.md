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
}
```

Alert calculation:
```go
func GetAlerts(db *sql.DB) ([]models.Alert, error) {
    // SELECT from antivirus + other_software where expiry_date within 30 days
}
```

## Config Layer

`config.Load()` reads `$XDG_CONFIG_HOME/inventory/config.json`. If the file doesn't exist, returns defaults:
- `dbPath`: Python app's default DB path (`~/.local/share/inventory/inventory.db`)
- `density`: `"comfortable"`
- `darkMode`: false
- `expiryWarningDays`: 30

## Vue Frontend

`lib/api/index.ts` replaces stub `Promise.resolve([])` with real Wails calls:
```typescript
import { ListComputers } from '../../wailsjs/go/bridge/App'
export const listComputers = (): Promise<Computer[]> => ListComputers()
```

`features/assets/ComputersView.vue` subscribes to assets store, calls `api.listComputers()` on mount:
```html
<DataTable :columns="computerColumns" :rows="store.items" :loading="loading" />
```

Column configs per category are defined as constants — no dynamic schema detection.

## Tests

Go (table-driven, `testing` + `testify/assert`):
- `internal/services/computers_test.go`: open in-memory SQLite, insert fixtures, verify ListComputers returns correct rows
- Same pattern for all 7 categories
- `internal/services/alerts_test.go`: test expiry calculation logic

Frontend (Vitest + Vue Test Utils):
- `DataTable.test.ts`: renders columns and rows, empty state when rows=[]
- `StatusBadge.test.ts`: correct class for each status value
