---
change: gover2-readonly
design-doc: docs/superpowers/specs/2026-06-18-gover2-readonly-design.md
base-ref: ddc19f352f8a79cbc4233328594b7535ee82e8a6
---

# gover2-readonly Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace every stub in the gover2 Wails v2 + Vue 3 desktop app with real implementations that open `inventory.db` (the Python app's SQLite file) and display all 7 asset categories in read-only data tables.

**Architecture:** The Go backend is layered as `database → config → models → services → bridge`; each layer imports only from layers below it. `bridge.App` becomes stateful — it holds `*sql.DB` and `models.AppConfig`, opened once in the Wails `OnStartup` callback. The Vue frontend replaces `Promise.resolve([])` stubs in `api/index.ts` with real Wails bridge calls; `DataTable.vue` becomes a generic column-driven component.

**Tech Stack:** Go 1.25, Wails v2, modernc.org/sqlite (pure-Go, already in go.mod), testify/assert (to add), Vue 3, TypeScript, Vite, Vitest 2, @vue/test-utils 2, jsdom 25, pnpm, Tailwind CSS 3.

## Global Constraints

- Module name: `gover2` (all Go imports use this prefix, e.g. `gover2/internal/database`)
- SQLite driver: `modernc.org/sqlite` — import side-effect only as `_ "modernc.org/sqlite"` in `database.go`; driver name is `"sqlite"` (not `"sqlite3"`)
- WAL mode: `PRAGMA journal_mode=WAL` set on every `Open()` call — safe because Python and Go apps never run simultaneously
- Config path: `$XDG_CONFIG_HOME/inventory/config.json` (fallback: `~/.config/inventory/config.json`)
- Default DB path for new installs: `~/.local/share/inventory/inventory.db`
- All bridge `List*` methods return `(nil, errors.New("no database open"))` when `a.db == nil`
- `GetAlerts` threshold: `expiryWarningDays` from config (default 30); use `date('now', '+N days')` SQL
- Frontend tests run via `pnpm --prefix gover2/frontend run test`
- Go tests run via `go test ./...` from inside `gover2/`
- `wails build` must be run from inside `gover2/` to regenerate `wailsjs/` bindings
- Severity logic for alerts: `daysRemaining < 0` → `"expired"`, otherwise `"expiring"`
- ExportCSV stub keeps returning `nil` — Phase 5 implements real export

---

### Task 1: Database Package — Open + InitSchema

**Files:**
- Modify: `gover2/internal/database/database.go` (replace 2-line stub with real implementation)
- Create: `gover2/internal/database/database_test.go`

**Interfaces:**
- Produces:
  - `database.Open(path string) (*sql.DB, error)` — opens SQLite with WAL pragma; returns `(*sql.DB, nil)` on success
  - `database.InitSchema(db *sql.DB) error` — creates 8 tables (\_meta + 7 asset tables) with `IF NOT EXISTS`

- [x] **Step 1: Write the failing test**

```go
// gover2/internal/database/database_test.go
package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestOpenInMemory(t *testing.T) {
	db, err := Open(":memory:")
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()
}

func TestInitSchema_CreatesAllTables(t *testing.T) {
	db, err := Open(":memory:")
	require.NoError(t, err)
	defer db.Close()

	err = InitSchema(db)
	require.NoError(t, err)

	want := []string{"_meta", "users", "computers", "smartphones", "tablets", "windows_keys", "other_software", "antivirus"}
	for _, tbl := range want {
		var name string
		err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", tbl).Scan(&name)
		assert.NoError(t, err, "table %q should exist", tbl)
		assert.Equal(t, tbl, name)
	}
}

func TestInitSchema_Idempotent(t *testing.T) {
	db, _ := Open(":memory:")
	defer db.Close()
	assert.NoError(t, InitSchema(db))
	assert.NoError(t, InitSchema(db), "second call must not error (IF NOT EXISTS)")
}
```

- [x] **Step 2: Add testify to go.mod**

Run from `gover2/`:
```bash
go get github.com/stretchr/testify@v1.10.0
```

Expected: `go.mod` and `go.sum` updated.

- [x] **Step 3: Run test to verify it fails**

Run from `gover2/`:
```bash
go test ./internal/database/...
```

Expected: FAIL — `Open` returns `nil, nil` so `NotNil(t, db)` fails.

- [x] **Step 4: Implement database.go**

Replace the entire contents of `gover2/internal/database/database.go`:

```go
package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// Open opens the SQLite database at path, enables WAL mode, and returns the connection.
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

// InitSchema creates all application tables if they do not already exist.
// Safe to call on an existing database (uses IF NOT EXISTS).
func InitSchema(db *sql.DB) error {
	_, err := db.Exec(`
		PRAGMA foreign_keys = ON;

		CREATE TABLE IF NOT EXISTS _meta (
			key   TEXT PRIMARY KEY,
			value TEXT
		);

		CREATE TABLE IF NOT EXISTS users (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT NOT NULL,
			surname    TEXT,
			status     TEXT DEFAULT 'Active',
			notes      TEXT,
			created_at TEXT DEFAULT (datetime('now','localtime')),
			updated_at TEXT DEFAULT (datetime('now','localtime'))
		);

		CREATE TABLE IF NOT EXISTS computers (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			name            TEXT,
			model           TEXT,
			user_id         INTEGER REFERENCES users(id) ON DELETE SET NULL,
			status          TEXT DEFAULT 'Active',
			purchase_date   TEXT,
			warranty_expiry TEXT,
			notes           TEXT,
			created_at      TEXT DEFAULT (datetime('now','localtime')),
			updated_at      TEXT DEFAULT (datetime('now','localtime'))
		);

		CREATE TABLE IF NOT EXISTS smartphones (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			name            TEXT,
			model           TEXT,
			user_id         INTEGER REFERENCES users(id) ON DELETE SET NULL,
			status          TEXT DEFAULT 'Active',
			purchase_date   TEXT,
			warranty_expiry TEXT,
			notes           TEXT,
			created_at      TEXT DEFAULT (datetime('now','localtime')),
			updated_at      TEXT DEFAULT (datetime('now','localtime'))
		);

		CREATE TABLE IF NOT EXISTS tablets (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			name            TEXT,
			model           TEXT,
			user_id         INTEGER REFERENCES users(id) ON DELETE SET NULL,
			status          TEXT DEFAULT 'Active',
			purchase_date   TEXT,
			warranty_expiry TEXT,
			notes           TEXT,
			created_at      TEXT DEFAULT (datetime('now','localtime')),
			updated_at      TEXT DEFAULT (datetime('now','localtime'))
		);

		CREATE TABLE IF NOT EXISTS windows_keys (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			license_key TEXT,
			computer_id INTEGER REFERENCES computers(id) ON DELETE SET NULL,
			status      TEXT DEFAULT 'Active',
			notes       TEXT,
			created_at  TEXT DEFAULT (datetime('now','localtime')),
			updated_at  TEXT DEFAULT (datetime('now','localtime'))
		);

		CREATE TABLE IF NOT EXISTS other_software (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			name          TEXT,
			license_key   TEXT,
			computer_id   INTEGER REFERENCES computers(id) ON DELETE SET NULL,
			smartphone_id INTEGER REFERENCES smartphones(id) ON DELETE SET NULL,
			tablet_id     INTEGER REFERENCES tablets(id) ON DELETE SET NULL,
			status        TEXT DEFAULT 'Active',
			expiry_date   TEXT,
			notes         TEXT,
			created_at    TEXT DEFAULT (datetime('now','localtime')),
			updated_at    TEXT DEFAULT (datetime('now','localtime'))
		);

		CREATE TABLE IF NOT EXISTS antivirus (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			name          TEXT,
			license_key   TEXT,
			computer_id   INTEGER REFERENCES computers(id) ON DELETE SET NULL,
			smartphone_id INTEGER REFERENCES smartphones(id) ON DELETE SET NULL,
			tablet_id     INTEGER REFERENCES tablets(id) ON DELETE SET NULL,
			status        TEXT DEFAULT 'Active',
			expiry_date   TEXT,
			notes         TEXT,
			created_at    TEXT DEFAULT (datetime('now','localtime')),
			updated_at    TEXT DEFAULT (datetime('now','localtime'))
		);
	`)
	return err
}
```

- [x] **Step 5: Run test to verify it passes**

Run from `gover2/`:
```bash
go test ./internal/database/...
```

Expected: `ok gover2/internal/database`

- [x] **Step 6: Commit**

```bash
git add gover2/internal/database/database.go gover2/internal/database/database_test.go gover2/go.mod gover2/go.sum
git commit -m "feat(gover2): implement database.Open and InitSchema with WAL mode"
```

---

### Task 2: Config Package — Load and Save

**Files:**
- Modify: `gover2/internal/config/config.go` (replace 7-line stub with real JSON-backed implementation)
- Create: `gover2/internal/config/config_test.go`

**Interfaces:**
- Consumes: `models.AppConfig` from `gover2/internal/models`
- Produces:
  - `config.Load() (models.AppConfig, error)` — reads `$XDG_CONFIG_HOME/inventory/config.json`; returns defaults when file missing
  - `config.Save(cfg models.AppConfig) error` — writes atomically via temp file + rename

- [x] **Step 1: Write the failing test**

```go
// gover2/internal/config/config_test.go
package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gover2/internal/models"
)

func TestLoad_ReturnsDefaultsWhenFileMissing(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "comfortable", cfg.Density)
	assert.Equal(t, false, cfg.DarkMode)
	assert.Equal(t, 30, cfg.ExpiryWarningDays)
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	want := models.AppConfig{
		DBPath:            "/tmp/test.db",
		Density:           "compact",
		DarkMode:          true,
		ExpiryWarningDays: 14,
	}
	require.NoError(t, Save(want))

	got, err := Load()
	require.NoError(t, err)
	assert.Equal(t, want, got)

	// Verify file exists at expected path
	_, err = os.Stat(filepath.Join(dir, "inventory", "config.json"))
	assert.NoError(t, err)
}
```

- [x] **Step 2: Run test to verify it fails**

Run from `gover2/`:
```bash
go test ./internal/config/...
```

Expected: FAIL — stub `Load()` returns empty `DBPath`, not the saved value.

- [x] **Step 3: Implement config.go**

Replace the entire contents of `gover2/internal/config/config.go`:

```go
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gover2/internal/models"
)

func configPath() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "inventory", "config.json")
}

// Load reads config from disk. If the file does not exist, returns safe defaults.
func Load() (models.AppConfig, error) {
	defaults := models.AppConfig{
		Density:           "comfortable",
		DarkMode:          false,
		ExpiryWarningDays: 30,
	}
	// Default DB path
	home, _ := os.UserHomeDir()
	defaults.DBPath = filepath.Join(home, ".local", "share", "inventory", "inventory.db")

	data, err := os.ReadFile(configPath())
	if errors.Is(err, os.ErrNotExist) {
		return defaults, nil
	}
	if err != nil {
		return defaults, err
	}

	var cfg models.AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaults, err
	}
	return cfg, nil
}

// Save writes cfg to disk atomically (temp file + rename).
func Save(cfg models.AppConfig) error {
	path := configPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
```

- [x] **Step 4: Run test to verify it passes**

Run from `gover2/`:
```bash
go test ./internal/config/...
```

Expected: `ok gover2/internal/config`

- [x] **Step 5: Commit**

```bash
git add gover2/internal/config/config.go gover2/internal/config/config_test.go
git commit -m "feat(gover2): implement config Load/Save with XDG path and atomic write"
```

---

### Task 3: Services Package — All 7 List Functions + GetAlerts

**Files:**
- Modify: `gover2/internal/services/services.go` (replace 7 stubs with real SQL queries; add `GetAlerts`)
- Create: `gover2/internal/services/services_test.go`

**Interfaces:**
- Consumes: `*sql.DB`, structs from `gover2/internal/models`; `database.Open` and `database.InitSchema` from Task 1
- Produces:
  - `services.ListComputers(db *sql.DB) ([]models.Computer, error)`
  - `services.ListSmartphones(db *sql.DB) ([]models.Smartphone, error)`
  - `services.ListTablets(db *sql.DB) ([]models.Tablet, error)`
  - `services.ListWindowsKeys(db *sql.DB) ([]models.WindowsKey, error)`
  - `services.ListAntivirus(db *sql.DB) ([]models.Antivirus, error)`
  - `services.ListOtherSoftware(db *sql.DB) ([]models.OtherSoftware, error)`
  - `services.ListUsers(db *sql.DB) ([]models.User, error)`
  - `services.GetAlerts(db *sql.DB, warningDays int) ([]models.Alert, error)`

- [x] **Step 1: Write the failing test**

```go
// gover2/internal/services/services_test.go
package services_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gover2/internal/database"
	"gover2/internal/models"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, database.InitSchema(db))
	t.Cleanup(func() { db.Close() })
	return db
}

func TestListComputers(t *testing.T) {
	db := openTestDB(t)
	_, err := db.Exec(`INSERT INTO computers (name, model, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		"PC-001", "Dell OptiPlex", "active", "2024-01-01", "2024-01-01")
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO computers (name, model, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		"PC-002", "HP EliteDesk", "retired", "2024-01-02", "2024-01-02")
	require.NoError(t, err)

	rows, err := ListComputers(db)
	require.NoError(t, err)
	assert.Len(t, rows, 2)
	// Ordered by name
	assert.Equal(t, "PC-001", rows[0].Name)
	assert.Equal(t, "PC-002", rows[1].Name)
	assert.Equal(t, "active", rows[0].Status)
}

func TestListSmartphones(t *testing.T) {
	db := openTestDB(t)
	db.Exec(`INSERT INTO smartphones (name, model, status, created_at, updated_at) VALUES ('iPhone 14','Apple iPhone 14','active','2024-01-01','2024-01-01')`)
	rows, err := ListSmartphones(db)
	require.NoError(t, err)
	assert.Len(t, rows, 1)
	assert.Equal(t, "iPhone 14", rows[0].Name)
}

func TestListTablets(t *testing.T) {
	db := openTestDB(t)
	db.Exec(`INSERT INTO tablets (name, model, status, created_at, updated_at) VALUES ('iPad Pro','Apple iPad Pro','active','2024-01-01','2024-01-01')`)
	rows, err := ListTablets(db)
	require.NoError(t, err)
	assert.Len(t, rows, 1)
}

func TestListWindowsKeys(t *testing.T) {
	db := openTestDB(t)
	db.Exec(`INSERT INTO windows_keys (license_key, status, created_at, updated_at) VALUES ('XXXXX-YYYYY-ZZZZZ-AAAAA-BBBBB','active','2024-01-01','2024-01-01')`)
	rows, err := ListWindowsKeys(db)
	require.NoError(t, err)
	assert.Len(t, rows, 1)
	assert.Equal(t, "XXXXX-YYYYY-ZZZZZ-AAAAA-BBBBB", rows[0].LicenseKey)
}

func TestListAntivirus(t *testing.T) {
	db := openTestDB(t)
	db.Exec(`INSERT INTO antivirus (name, license_key, status, created_at, updated_at) VALUES ('Defender Pro','AV-KEY-001','active','2024-01-01','2024-01-01')`)
	rows, err := ListAntivirus(db)
	require.NoError(t, err)
	assert.Len(t, rows, 1)
	assert.Equal(t, "Defender Pro", rows[0].Name)
}

func TestListOtherSoftware(t *testing.T) {
	db := openTestDB(t)
	db.Exec(`INSERT INTO other_software (name, license_key, status, created_at, updated_at) VALUES ('Adobe CC','ADOBE-KEY-001','active','2024-01-01','2024-01-01')`)
	rows, err := ListOtherSoftware(db)
	require.NoError(t, err)
	assert.Len(t, rows, 1)
}

func TestListUsers(t *testing.T) {
	db := openTestDB(t)
	db.Exec(`INSERT INTO users (name, status, created_at, updated_at) VALUES ('Alice','active','2024-01-01','2024-01-01')`)
	db.Exec(`INSERT INTO users (name, surname, status, created_at, updated_at) VALUES ('Bob','Smith','active','2024-01-01','2024-01-01')`)
	rows, err := ListUsers(db)
	require.NoError(t, err)
	assert.Len(t, rows, 2)
	// surname nullable
	assert.Nil(t, rows[0].Surname) // Alice has no surname
	assert.NotNil(t, rows[1].Surname)
	assert.Equal(t, "Smith", *rows[1].Surname)
}

func TestGetAlerts_WithinWarning(t *testing.T) {
	db := openTestDB(t)
	// Expiry 15 days from now — within 30-day window
	near := time.Now().AddDate(0, 0, 15).Format("2006-01-02")
	far := time.Now().AddDate(0, 0, 45).Format("2006-01-02")
	db.Exec(fmt.Sprintf(`INSERT INTO antivirus (name, license_key, status, expiry_date, created_at, updated_at) VALUES ('AV Near','KEY1','active','%s','2024-01-01','2024-01-01')`, near))
	db.Exec(fmt.Sprintf(`INSERT INTO antivirus (name, license_key, status, expiry_date, created_at, updated_at) VALUES ('AV Far','KEY2','active','%s','2024-01-01','2024-01-01')`, far))

	alerts, err := GetAlerts(db, 30)
	require.NoError(t, err)
	assert.Len(t, alerts, 1)
	assert.Equal(t, "AV Near", alerts[0].Name)
	assert.Equal(t, "antivirus", alerts[0].Category)
	assert.Equal(t, "expiring", alerts[0].Severity)
}

func TestGetAlerts_Expired(t *testing.T) {
	db := openTestDB(t)
	past := time.Now().AddDate(0, 0, -5).Format("2006-01-02")
	db.Exec(fmt.Sprintf(`INSERT INTO other_software (name, license_key, status, expiry_date, created_at, updated_at) VALUES ('Old App','KEY3','active','%s','2024-01-01','2024-01-01')`, past))

	alerts, err := GetAlerts(db, 30)
	require.NoError(t, err)
	assert.Len(t, alerts, 1)
	assert.Equal(t, "expired", alerts[0].Severity)
	assert.True(t, alerts[0].DaysRemaining < 0)
}
```

- [x] **Step 2: Run test to verify it fails**

Run from `gover2/`:
```bash
go test ./internal/services/...
```

Expected: FAIL — stubs return empty slices, `TestListComputers` fails on `Len(t, rows, 2)`.

- [x] **Step 3: Implement services.go**

Replace the entire contents of `gover2/internal/services/services.go`:

```go
package services

import (
	"database/sql"
	"fmt"

	"gover2/internal/models"
)

func ListComputers(db *sql.DB) ([]models.Computer, error) {
	rows, err := db.Query(`SELECT id, name, model, user_id, status, purchase_date, warranty_expiry, notes, created_at, updated_at FROM computers ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Computer
	for rows.Next() {
		var c models.Computer
		if err := rows.Scan(&c.ID, &c.Name, &c.Model, &c.UserID, &c.Status, &c.PurchaseDate, &c.WarrantyExpiry, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func ListSmartphones(db *sql.DB) ([]models.Smartphone, error) {
	rows, err := db.Query(`SELECT id, name, model, user_id, status, purchase_date, warranty_expiry, notes, created_at, updated_at FROM smartphones ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Smartphone
	for rows.Next() {
		var s models.Smartphone
		if err := rows.Scan(&s.ID, &s.Name, &s.Model, &s.UserID, &s.Status, &s.PurchaseDate, &s.WarrantyExpiry, &s.Notes, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func ListTablets(db *sql.DB) ([]models.Tablet, error) {
	rows, err := db.Query(`SELECT id, name, model, user_id, status, purchase_date, warranty_expiry, notes, created_at, updated_at FROM tablets ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Tablet
	for rows.Next() {
		var t models.Tablet
		if err := rows.Scan(&t.ID, &t.Name, &t.Model, &t.UserID, &t.Status, &t.PurchaseDate, &t.WarrantyExpiry, &t.Notes, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func ListWindowsKeys(db *sql.DB) ([]models.WindowsKey, error) {
	rows, err := db.Query(`SELECT id, license_key, computer_id, status, notes, created_at, updated_at FROM windows_keys ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.WindowsKey
	for rows.Next() {
		var w models.WindowsKey
		if err := rows.Scan(&w.ID, &w.LicenseKey, &w.ComputerID, &w.Status, &w.Notes, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

func ListAntivirus(db *sql.DB) ([]models.Antivirus, error) {
	rows, err := db.Query(`SELECT id, name, license_key, computer_id, smartphone_id, tablet_id, status, expiry_date, notes, created_at, updated_at FROM antivirus ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Antivirus
	for rows.Next() {
		var a models.Antivirus
		if err := rows.Scan(&a.ID, &a.Name, &a.LicenseKey, &a.ComputerID, &a.SmartphoneID, &a.TabletID, &a.Status, &a.ExpiryDate, &a.Notes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func ListOtherSoftware(db *sql.DB) ([]models.OtherSoftware, error) {
	rows, err := db.Query(`SELECT id, name, license_key, computer_id, smartphone_id, tablet_id, status, expiry_date, notes, created_at, updated_at FROM other_software ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.OtherSoftware
	for rows.Next() {
		var o models.OtherSoftware
		if err := rows.Scan(&o.ID, &o.Name, &o.LicenseKey, &o.ComputerID, &o.SmartphoneID, &o.TabletID, &o.Status, &o.ExpiryDate, &o.Notes, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

func ListUsers(db *sql.DB) ([]models.User, error) {
	rows, err := db.Query(`SELECT id, name, surname, status, notes, created_at, updated_at FROM users ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Status, &u.Notes, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// GetAlerts returns items from antivirus and other_software whose expiry_date
// falls within the next warningDays days (or is already past).
// Results are ordered by expiry_date ascending.
func GetAlerts(db *sql.DB, warningDays int) ([]models.Alert, error) {
	threshold := fmt.Sprintf("+%d days", warningDays)
	q := `
		SELECT 'antivirus' AS category, id, name, expiry_date,
		       CAST(julianday(expiry_date) - julianday('now') AS INTEGER) AS days_remaining
		FROM antivirus
		WHERE expiry_date IS NOT NULL AND date(expiry_date) <= date('now', ?)
		UNION ALL
		SELECT 'other_software' AS category, id, name, expiry_date,
		       CAST(julianday(expiry_date) - julianday('now') AS INTEGER) AS days_remaining
		FROM other_software
		WHERE expiry_date IS NOT NULL AND date(expiry_date) <= date('now', ?)
		ORDER BY expiry_date`

	rows, err := db.Query(q, threshold, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Alert
	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(&a.Category, &a.ID, &a.Name, &a.ExpiryDate, &a.DaysRemaining); err != nil {
			return nil, err
		}
		if a.DaysRemaining < 0 {
			a.Severity = "expired"
		} else {
			a.Severity = "expiring"
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

```

- [x] **Step 4: Run test to verify it passes**

Run from `gover2/`:
```bash
go test -cover ./internal/services/...
```

Expected: `ok gover2/internal/services` with coverage >= 70%.

- [x] **Step 5: Commit**

```bash
git add gover2/internal/services/services.go gover2/internal/services/services_test.go
git commit -m "feat(gover2): implement all 7 List service functions and GetAlerts"
```

---

### Task 4: Bridge Package — Stateful App + All Methods

**Files:**
- Modify: `gover2/internal/bridge/bridge.go` (replace stateless `App{}` with stateful struct; implement all methods)
- Modify: `gover2/app.go` (add `startup` method wiring)
- Modify: `gover2/main.go` (add `OnStartup` to wails.Run options)
- Create: `gover2/internal/bridge/bridge_test.go`

**Interfaces:**
- Consumes: `database.Open`, `database.InitSchema` (Task 1); `config.Load`, `config.Save` (Task 2); all `services.List*` and `services.GetAlerts` (Task 3)
- Produces: all bridge methods used by the frontend; `App.startup(ctx context.Context)` wired to Wails `OnStartup`

- [ ] **Step 1: Write the failing test**

```go
// gover2/internal/bridge/bridge_test.go
package bridge_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gover2/internal/database"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, database.InitSchema(db))
	t.Cleanup(func() { db.Close() })
	return db
}

func TestBridge_ListComputers_NilDB(t *testing.T) {
	app := &App{}
	_, err := app.ListComputers()
	assert.Error(t, err)
	assert.EqualError(t, err, "no database open")
}

func TestBridge_ListComputers_WithData(t *testing.T) {
	db := openTestDB(t)
	db.Exec(`INSERT INTO computers (name, model, status, created_at, updated_at) VALUES ('PC-Test','Model X','active','2024-01-01','2024-01-01')`)

	app := &App{db: db}
	rows, err := app.ListComputers()
	require.NoError(t, err)
	assert.Len(t, rows, 1)
	assert.Equal(t, "PC-Test", rows[0].Name)
}

func TestBridge_ListWindowsKeys_NilDB(t *testing.T) {
	app := &App{}
	_, err := app.ListWindowsKeys()
	assert.EqualError(t, err, "no database open")
}

func TestBridge_GetAlerts_NilDB(t *testing.T) {
	app := &App{}
	_, err := app.GetAlerts()
	assert.EqualError(t, err, "no database open")
}

func TestBridge_GetConfig_ReturnsDefaults(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	app := &App{}
	cfg, err := app.GetConfig()
	require.NoError(t, err)
	assert.Equal(t, "comfortable", cfg.Density)
}

func TestBridge_NewDatabase_OpensDB(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/test.db"
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	app := &App{}
	err := app.NewDatabase(path)
	require.NoError(t, err)
	assert.NotNil(t, app.db)
	assert.Equal(t, path, app.cfg.DBPath)
}

func TestBridge_OpenDatabase_FileNotFound(t *testing.T) {
	app := &App{}
	err := app.OpenDatabase("/nonexistent/path/test.db")
	assert.Error(t, err)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run from `gover2/`:
```bash
go test ./internal/bridge/...
```

Expected: FAIL — `App{}` has no `db` field, `ListComputers` does not return `"no database open"` error.

- [ ] **Step 3: Implement bridge.go**

Replace the entire contents of `gover2/internal/bridge/bridge.go`:

```go
package bridge

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"gover2/internal/config"
	"gover2/internal/database"
	"gover2/internal/models"
	"gover2/internal/services"
)

// App is the Wails application struct. Fields are populated in startup().
type App struct {
	ctx context.Context
	db  *sql.DB
	cfg models.AppConfig
}

func NewApp() *App {
	return &App{}
}

// startup is called by Wails OnStartup. It loads config and opens the database.
// If either fails, the app continues with a nil db — views show empty state.
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

// Startup is the exported startup method for Wails wiring.
func (a *App) Startup(ctx context.Context) {
	a.startup(ctx)
}

func (a *App) ListComputers() ([]models.Computer, error) {
	if a.db == nil {
		return nil, errors.New("no database open")
	}
	return services.ListComputers(a.db)
}

func (a *App) ListSmartphones() ([]models.Smartphone, error) {
	if a.db == nil {
		return nil, errors.New("no database open")
	}
	return services.ListSmartphones(a.db)
}

func (a *App) ListTablets() ([]models.Tablet, error) {
	if a.db == nil {
		return nil, errors.New("no database open")
	}
	return services.ListTablets(a.db)
}

func (a *App) ListWindowsKeys() ([]models.WindowsKey, error) {
	if a.db == nil {
		return nil, errors.New("no database open")
	}
	return services.ListWindowsKeys(a.db)
}

func (a *App) ListAntivirus() ([]models.Antivirus, error) {
	if a.db == nil {
		return nil, errors.New("no database open")
	}
	return services.ListAntivirus(a.db)
}

func (a *App) ListOtherSoftware() ([]models.OtherSoftware, error) {
	if a.db == nil {
		return nil, errors.New("no database open")
	}
	return services.ListOtherSoftware(a.db)
}

func (a *App) ListUsers() ([]models.User, error) {
	if a.db == nil {
		return nil, errors.New("no database open")
	}
	return services.ListUsers(a.db)
}

func (a *App) GetAlerts() ([]models.Alert, error) {
	if a.db == nil {
		return nil, errors.New("no database open")
	}
	days := a.cfg.ExpiryWarningDays
	if days == 0 {
		days = 30
	}
	return services.GetAlerts(a.db, days)
}

func (a *App) GetDatabasePath() (string, error) {
	return a.cfg.DBPath, nil
}

func (a *App) GetConfig() (models.AppConfig, error) {
	return config.Load()
}

func (a *App) SetConfig(cfg models.AppConfig) error {
	a.cfg = cfg
	return config.Save(cfg)
}

func (a *App) NewDatabase(path string) error {
	db, err := database.Open(path)
	if err != nil {
		return err
	}
	database.InitSchema(db)
	if a.db != nil {
		a.db.Close()
	}
	a.db = db
	a.cfg.DBPath = path
	return config.Save(a.cfg)
}

func (a *App) OpenDatabase(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("file not found: %s", path)
	}
	return a.NewDatabase(path)
}

func (a *App) ExportCSV(category string, destPath string) error {
	return nil // Phase 5
}
```

- [ ] **Step 4: Update app.go to export startup**

Replace the entire contents of `gover2/app.go`:

```go
package main

import (
	"context"

	"gover2/internal/bridge"
)

func NewApp() *bridge.App {
	return bridge.NewApp()
}

// startup is called at application startup — wired in main.go via OnStartup.
func startup(app *bridge.App) func(ctx context.Context) {
	return app.Startup
}
```

- [ ] **Step 5: Update main.go to wire OnStartup**

Replace the entire contents of `gover2/main.go`:

```go
package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "gover2",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        startup(app),
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
```

- [ ] **Step 6: Run test to verify it passes**

Run from `gover2/`:
```bash
go test -cover ./internal/bridge/...
```

Expected: `ok gover2/internal/bridge` with coverage >= 50%.

- [ ] **Step 7: Verify all Go packages compile**

Run from `gover2/`:
```bash
go build ./...
```

Expected: exit 0, no errors.

- [ ] **Step 8: Commit**

```bash
git add gover2/internal/bridge/bridge.go gover2/internal/bridge/bridge_test.go gover2/app.go gover2/main.go
git commit -m "feat(gover2): stateful bridge App with OnStartup DB wiring and all List methods"
```

---

### Task 5: Setup Frontend Test Infrastructure

**Files:**
- Modify: `gover2/frontend/package.json` (add vitest, @vue/test-utils, jsdom dev deps; add test script)
- Modify: `gover2/frontend/vite.config.ts` (add `test` block for jsdom + globals)

**Interfaces:**
- Produces: `pnpm --prefix gover2/frontend run test` command works and runs `.test.ts` files

- [x] **Step 1: Add test dependencies**

Run from `gover2/frontend/`:
```bash
pnpm add -D vitest@^2.0 @vue/test-utils@^2 jsdom@^25
```

Expected: `pnpm-lock.yaml` updated; packages installed in `node_modules`.

- [x] **Step 2: Add test script to package.json**

In `gover2/frontend/package.json`, add `"test": "vitest run"` to the `scripts` section:

```json
{
  "name": "gover2-frontend",
  "private": true,
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc && vite build",
    "typecheck": "vue-tsc --noEmit",
    "preview": "vite preview",
    "clean": "rm -rf dist",
    "test": "vitest run"
  },
  "dependencies": {
    "pinia": "^2.1.7",
    "vue": "^3.4.0",
    "vue-router": "^4.3.0"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.0.0",
    "@vue/test-utils": "^2",
    "autoprefixer": "^10.4.0",
    "jsdom": "^25",
    "postcss": "^8.4.0",
    "tailwindcss": "3",
    "typescript": "^5.4.0",
    "vite": "^5.0.0",
    "vitest": "^2.0",
    "vue-tsc": "^2.0.0"
  }
}
```

- [x] **Step 3: Add test config to vite.config.ts**

Replace the entire contents of `gover2/frontend/vite.config.ts`:

```ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
  },
})
```

- [x] **Step 4: Verify test infrastructure works (no test files yet = zero tests pass)**

Run from project root:
```bash
pnpm --prefix gover2/frontend run test
```

Expected: `No test files found` or `0 tests` — exit 0. (Not a test failure.)

- [x] **Step 5: Commit**

```bash
git add gover2/frontend/package.json gover2/frontend/pnpm-lock.yaml gover2/frontend/vite.config.ts
git commit -m "chore(gover2): add vitest + @vue/test-utils + jsdom for frontend testing"
```

---

### Task 6: DataTable.vue — Generic Column/Row Component + Tests

**Files:**
- Modify: `gover2/frontend/src/components/DataTable.vue` (replace hard-coded stub with generic columns/rows/loading props)
- Create: `gover2/frontend/src/components/DataTable.test.ts`

**Interfaces:**
- Produces: `<DataTable :columns="ColDef[]" :rows="Record<string,unknown>[]" :loading="boolean" />`
  - `ColDef = { key: string; label: string }`
  - Renders `<th>` for each column, `<tr>` for each row, empty-state `<td>` when `rows.length === 0 && !loading`, loading row when `loading === true`
- Note: All view files that currently pass `:items="[]"` will be updated in Task 7.

- [x] **Step 1: Write the failing test**

```ts
// gover2/frontend/src/components/DataTable.test.ts
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import DataTable from './DataTable.vue'

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'status', label: 'Status' },
  { key: 'model', label: 'Model' },
]

const rows = [
  { name: 'PC-001', status: 'active', model: 'Dell' },
  { name: 'PC-002', status: 'retired', model: 'HP' },
  { name: 'PC-003', status: 'spare', model: 'Lenovo' },
]

describe('DataTable', () => {
  it('renders correct number of header columns', () => {
    const wrapper = mount(DataTable, { props: { columns, rows, loading: false } })
    const ths = wrapper.findAll('thead th')
    expect(ths).toHaveLength(3)
    expect(ths[0].text()).toBe('Name')
    expect(ths[1].text()).toBe('Status')
    expect(ths[2].text()).toBe('Model')
  })

  it('renders correct number of body rows', () => {
    const wrapper = mount(DataTable, { props: { columns, rows, loading: false } })
    const trs = wrapper.findAll('tbody tr')
    expect(trs).toHaveLength(3)
  })

  it('renders correct cell values', () => {
    const wrapper = mount(DataTable, { props: { columns, rows, loading: false } })
    const cells = wrapper.findAll('tbody tr:first-child td')
    expect(cells[0].text()).toBe('PC-001')
    expect(cells[1].text()).toBe('active')
    expect(cells[2].text()).toBe('Dell')
  })

  it('shows empty state when rows is empty and not loading', () => {
    const wrapper = mount(DataTable, { props: { columns, rows: [], loading: false } })
    expect(wrapper.find('[data-testid="empty-state"]').exists()).toBe(true)
    expect(wrapper.findAll('tbody tr')).toHaveLength(1) // the empty-state row
  })

  it('shows loading indicator when loading is true', () => {
    const wrapper = mount(DataTable, { props: { columns, rows: [], loading: true } })
    expect(wrapper.find('[data-testid="loading-state"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="empty-state"]').exists()).toBe(false)
  })
})
```

- [x] **Step 2: Run test to verify it fails**

Run from project root:
```bash
pnpm --prefix gover2/frontend run test
```

Expected: FAIL — current `DataTable.vue` uses `:items` prop (not `columns`/`rows`/`loading`), `[data-testid]` selectors not present.

- [x] **Step 3: Implement DataTable.vue**

Replace the entire contents of `gover2/frontend/src/components/DataTable.vue`:

```vue
<template>
  <div class="bg-surface border border-border rounded-none overflow-hidden">
    <table class="w-full text-sm">
      <thead class="bg-th-bg">
        <tr>
          <th
            v-for="col in columns"
            :key="col.key"
            class="text-left px-4 py-2 text-th-text font-medium"
          >
            {{ col.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="loading" data-testid="loading-state">
          <td :colspan="columns.length" class="px-4 py-8 text-center text-text-secondary">
            Loading…
          </td>
        </tr>
        <tr v-else-if="rows.length === 0" data-testid="empty-state">
          <td :colspan="columns.length" class="px-4 py-8 text-center text-text-secondary">
            No items to display.
          </td>
        </tr>
        <template v-else>
          <tr
            v-for="(row, i) in rows"
            :key="i"
            class="border-t border-border hover:bg-surface-hover"
          >
            <td
              v-for="col in columns"
              :key="col.key"
              class="px-4 py-2 text-text-primary"
            >
              {{ row[col.key] ?? '' }}
            </td>
          </tr>
        </template>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
export interface ColDef {
  key: string
  label: string
}

defineProps<{
  columns: ColDef[]
  rows: Record<string, unknown>[]
  loading: boolean
}>()
</script>
```

- [x] **Step 4: Run test to verify it passes**

Run from project root:
```bash
pnpm --prefix gover2/frontend run test
```

Expected: `5 tests pass` in DataTable.test.ts.

- [x] **Step 5: Commit**

```bash
git add gover2/frontend/src/components/DataTable.vue gover2/frontend/src/components/DataTable.test.ts
git commit -m "feat(gover2): generic DataTable component with columns/rows/loading props"
```

---

### Task 7: StatusBadge Test + All 8 Feature Views

**Files:**
- Create: `gover2/frontend/src/components/StatusBadge.test.ts`
- Modify: `gover2/frontend/src/features/assets/ComputersView.vue`
- Modify: `gover2/frontend/src/features/assets/SmartphonesView.vue`
- Modify: `gover2/frontend/src/features/assets/TabletsView.vue`
- Modify: `gover2/frontend/src/features/assets/AllAssetsView.vue`
- Modify: `gover2/frontend/src/features/licenses/WindowsKeysView.vue`
- Modify: `gover2/frontend/src/features/licenses/AntivirusView.vue`
- Modify: `gover2/frontend/src/features/licenses/OtherSoftwareView.vue`
- Modify: `gover2/frontend/src/features/users/UsersView.vue`

**Interfaces:**
- Consumes: `<DataTable :columns="ColDef[]" :rows="Record<string,unknown>[]" :loading="boolean" />` from Task 6; api functions from Task 8 (wired through `@/lib/api`)
- Produces: each view loads data on `onMounted`, passes it to `<DataTable>`; `AllAssetsView` unions computers + smartphones + tablets with a `category` field

- [x] **Step 1: Write the StatusBadge test**

```ts
// gover2/frontend/src/components/StatusBadge.test.ts
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StatusBadge from './StatusBadge.vue'

describe('StatusBadge', () => {
  const cases = [
    ['active', 'bg-s-active'],
    ['repair', 'bg-s-repair'],
    ['spare', 'bg-s-spare'],
    ['retired', 'bg-s-retired'],
    ['missing', 'bg-s-missing'],
    ['expiring', 'bg-s-expiring'],
    ['expired', 'bg-s-expired'],
  ] as const

  for (const [status, expectedClass] of cases) {
    it(`status "${status}" renders class "${expectedClass}"`, () => {
      const wrapper = mount(StatusBadge, { props: { status } })
      expect(wrapper.find('span').classes()).toContain(expectedClass)
    })
  }

  it('unknown status falls back to bg-s-retired', () => {
    const wrapper = mount(StatusBadge, { props: { status: 'unknown-xyz' } })
    expect(wrapper.find('span').classes()).toContain('bg-s-retired')
  })
})
```

- [x] **Step 2: Run StatusBadge test to verify it passes (StatusBadge.vue already correct)**

Run from project root:
```bash
pnpm --prefix gover2/frontend run test
```

Expected: StatusBadge tests pass (8 tests); DataTable tests still pass.

- [x] **Step 3: Implement ComputersView.vue**

Replace the entire contents of `gover2/frontend/src/features/assets/ComputersView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Computers</h1>
    <DataTable :columns="columns" :rows="rows as Record<string, unknown>[]" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '@/components/DataTable.vue'
import { listComputers } from '@/lib/api'
import type { Computer } from '@/lib/types'

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'model', label: 'Model' },
  { key: 'userId', label: 'User' },
  { key: 'status', label: 'Status' },
  { key: 'purchaseDate', label: 'Purchase Date' },
  { key: 'warrantyExpiry', label: 'Warranty Expiry' },
]

const rows = ref<Computer[]>([])
const loading = ref(true)

onMounted(async () => {
  rows.value = await listComputers()
  loading.value = false
})
</script>
```

- [x] **Step 4: Implement SmartphonesView.vue**

Replace the entire contents of `gover2/frontend/src/features/assets/SmartphonesView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Smartphones</h1>
    <DataTable :columns="columns" :rows="rows as Record<string, unknown>[]" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '@/components/DataTable.vue'
import { listSmartphones } from '@/lib/api'
import type { Smartphone } from '@/lib/types'

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'model', label: 'Model' },
  { key: 'userId', label: 'User' },
  { key: 'status', label: 'Status' },
  { key: 'purchaseDate', label: 'Purchase Date' },
  { key: 'warrantyExpiry', label: 'Warranty Expiry' },
]

const rows = ref<Smartphone[]>([])
const loading = ref(true)

onMounted(async () => {
  rows.value = await listSmartphones()
  loading.value = false
})
</script>
```

- [x] **Step 5: Implement TabletsView.vue**

Replace the entire contents of `gover2/frontend/src/features/assets/TabletsView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Tablets</h1>
    <DataTable :columns="columns" :rows="rows as Record<string, unknown>[]" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '@/components/DataTable.vue'
import { listTablets } from '@/lib/api'
import type { Tablet } from '@/lib/types'

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'model', label: 'Model' },
  { key: 'userId', label: 'User' },
  { key: 'status', label: 'Status' },
  { key: 'purchaseDate', label: 'Purchase Date' },
  { key: 'warrantyExpiry', label: 'Warranty Expiry' },
]

const rows = ref<Tablet[]>([])
const loading = ref(true)

onMounted(async () => {
  rows.value = await listTablets()
  loading.value = false
})
</script>
```

- [x] **Step 6: Implement WindowsKeysView.vue**

Replace the entire contents of `gover2/frontend/src/features/licenses/WindowsKeysView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Windows Keys</h1>
    <DataTable :columns="columns" :rows="rows as Record<string, unknown>[]" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '@/components/DataTable.vue'
import { listWindowsKeys } from '@/lib/api'
import type { WindowsKey } from '@/lib/types'

const columns = [
  { key: 'licenseKey', label: 'License Key' },
  { key: 'computerId', label: 'Computer' },
  { key: 'status', label: 'Status' },
  { key: 'notes', label: 'Notes' },
]

const rows = ref<WindowsKey[]>([])
const loading = ref(true)

onMounted(async () => {
  rows.value = await listWindowsKeys()
  loading.value = false
})
</script>
```

- [x] **Step 7: Implement AntivirusView.vue**

Replace the entire contents of `gover2/frontend/src/features/licenses/AntivirusView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Antivirus</h1>
    <DataTable :columns="columns" :rows="rows as Record<string, unknown>[]" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '@/components/DataTable.vue'
import { listAntivirus } from '@/lib/api'
import type { Antivirus } from '@/lib/types'

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'licenseKey', label: 'License Key' },
  { key: 'computerId', label: 'Computer' },
  { key: 'status', label: 'Status' },
  { key: 'expiryDate', label: 'Expiry Date' },
]

const rows = ref<Antivirus[]>([])
const loading = ref(true)

onMounted(async () => {
  rows.value = await listAntivirus()
  loading.value = false
})
</script>
```

- [x] **Step 8: Implement OtherSoftwareView.vue**

Replace the entire contents of `gover2/frontend/src/features/licenses/OtherSoftwareView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Other Software</h1>
    <DataTable :columns="columns" :rows="rows as Record<string, unknown>[]" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '@/components/DataTable.vue'
import { listOtherSoftware } from '@/lib/api'
import type { OtherSoftware } from '@/lib/types'

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'licenseKey', label: 'License Key' },
  { key: 'computerId', label: 'Computer' },
  { key: 'status', label: 'Status' },
  { key: 'expiryDate', label: 'Expiry Date' },
]

const rows = ref<OtherSoftware[]>([])
const loading = ref(true)

onMounted(async () => {
  rows.value = await listOtherSoftware()
  loading.value = false
})
</script>
```

- [x] **Step 9: Implement UsersView.vue**

Replace the entire contents of `gover2/frontend/src/features/users/UsersView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Users</h1>
    <DataTable :columns="columns" :rows="rows as Record<string, unknown>[]" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '@/components/DataTable.vue'
import { listUsers } from '@/lib/api'
import type { User } from '@/lib/types'

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'surname', label: 'Surname' },
  { key: 'status', label: 'Status' },
  { key: 'notes', label: 'Notes' },
]

const rows = ref<User[]>([])
const loading = ref(true)

onMounted(async () => {
  rows.value = await listUsers()
  loading.value = false
})
</script>
```

- [x] **Step 10: Implement AllAssetsView.vue**

Replace the entire contents of `gover2/frontend/src/features/assets/AllAssetsView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">All Assets</h1>
    <DataTable :columns="columns" :rows="rows" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '@/components/DataTable.vue'
import { listComputers, listSmartphones, listTablets } from '@/lib/api'

interface AssetRow {
  category: string
  name: string
  model: string
  userId: number | null
  status: string
  purchaseDate: string | null
  warrantyExpiry: string | null
  [key: string]: unknown
}

const columns = [
  { key: 'category', label: 'Category' },
  { key: 'name', label: 'Name' },
  { key: 'model', label: 'Model' },
  { key: 'userId', label: 'User' },
  { key: 'status', label: 'Status' },
  { key: 'purchaseDate', label: 'Purchase Date' },
  { key: 'warrantyExpiry', label: 'Warranty Expiry' },
]

const rows = ref<AssetRow[]>([])
const loading = ref(true)

onMounted(async () => {
  const [computers, smartphones, tablets] = await Promise.all([
    listComputers(),
    listSmartphones(),
    listTablets(),
  ])
  rows.value = [
    ...computers.map((c) => ({ ...c, category: 'Computer' })),
    ...smartphones.map((s) => ({ ...s, category: 'Smartphone' })),
    ...tablets.map((t) => ({ ...t, category: 'Tablet' })),
  ]
  loading.value = false
})
</script>
```

- [x] **Step 11: Run all frontend tests**

Run from project root:
```bash
pnpm --prefix gover2/frontend run test
```

Expected: all 13 tests pass (5 DataTable + 8 StatusBadge).

- [x] **Step 12: Commit**

```bash
git add \
  gover2/frontend/src/components/StatusBadge.test.ts \
  gover2/frontend/src/features/assets/ComputersView.vue \
  gover2/frontend/src/features/assets/SmartphonesView.vue \
  gover2/frontend/src/features/assets/TabletsView.vue \
  gover2/frontend/src/features/assets/AllAssetsView.vue \
  gover2/frontend/src/features/licenses/WindowsKeysView.vue \
  gover2/frontend/src/features/licenses/AntivirusView.vue \
  gover2/frontend/src/features/licenses/OtherSoftwareView.vue \
  gover2/frontend/src/features/users/UsersView.vue
git commit -m "feat(gover2): wire all 8 feature views to DataTable with real api calls"
```

---

### Task 8: api/index.ts — Replace Stubs with Real Wails Bindings

**Files:**
- Modify: `gover2/frontend/src/lib/api/index.ts` (replace all `Promise.resolve` stubs with real Wails bridge calls)

**Interfaces:**
- Consumes: generated `wailsjs/go/bridge/App.js` + `App.d.ts` (already present; regenerated by `wails build` in Task 9)
- Note: The `wailsjs/` files are already committed from the scaffold. The existing `App.d.ts` does NOT include `GetAlerts`. After `wails build` in Task 9 it will be regenerated to include it. For now, call the bridge methods that already exist and add `getAlerts` as a direct `window.go` call as a fallback until the bindings regenerate.

- [x] **Step 1: Verify wailsjs bindings already present**

Run from project root:
```bash
ls gover2/frontend/wailsjs/go/bridge/
```

Expected: `App.d.ts` and `App.js` present.

- [x] **Step 2: Replace api/index.ts with real bridge calls**

Replace the entire contents of `gover2/frontend/src/lib/api/index.ts`:

```ts
// All Wails bridge calls go through this file.
// Components MUST NOT call window.go.* directly.
// wailsjs/ is regenerated by `wails build` — do not hand-edit those files.

import {
  ListComputers,
  ListSmartphones,
  ListTablets,
  ListWindowsKeys,
  ListAntivirus,
  ListOtherSoftware,
  ListUsers,
  GetConfig,
  SetConfig,
  GetDatabasePath,
  OpenDatabase,
  NewDatabase,
  ExportCSV,
} from '../../wailsjs/go/bridge/App'
import type { models } from '../../wailsjs/go/models'
import type { Alert } from '@/lib/types'

export type { models }

const defaultConfig: models.AppConfig = {
  dbPath: '',
  density: 'comfortable',
  darkMode: false,
  expiryWarningDays: 30,
}

export const listComputers = (): Promise<models.Computer[]> =>
  ListComputers().catch(() => [])

export const listSmartphones = (): Promise<models.Smartphone[]> =>
  ListSmartphones().catch(() => [])

export const listTablets = (): Promise<models.Tablet[]> =>
  ListTablets().catch(() => [])

export const listWindowsKeys = (): Promise<models.WindowsKey[]> =>
  ListWindowsKeys().catch(() => [])

export const listAntivirus = (): Promise<models.Antivirus[]> =>
  ListAntivirus().catch(() => [])

export const listOtherSoftware = (): Promise<models.OtherSoftware[]> =>
  ListOtherSoftware().catch(() => [])

export const listUsers = (): Promise<models.User[]> =>
  ListUsers().catch(() => [])

// GetAlerts is added to the bridge in this change; wailsjs bindings regenerate on `wails build`.
// Until then, call the bridge method directly via window.go.
export const getAlerts = (): Promise<Alert[]> =>
  (window as unknown as Record<string, Record<string, Record<string, () => Promise<Alert[]>>>>)
    ?.go?.bridge?.App?.GetAlerts?.()
    .catch(() => []) ?? Promise.resolve([])

export const getConfig = (): Promise<models.AppConfig> =>
  GetConfig().catch(() => defaultConfig)

export const setConfig = (cfg: models.AppConfig): Promise<void> =>
  SetConfig(cfg)

export const getDatabasePath = (): Promise<string> =>
  GetDatabasePath().catch(() => '')

export const openDatabase = (path: string): Promise<void> =>
  OpenDatabase(path)

export const newDatabase = (path: string): Promise<void> =>
  NewDatabase(path)

export const exportCSV = (category: string, destPath: string): Promise<void> =>
  ExportCSV(category, destPath)
```

- [x] **Step 3: Run TypeScript typecheck**

Run from project root:
```bash
pnpm --prefix gover2/frontend run typecheck
```

Expected: exit 0. (If `models.Alert` is not in `models.ts`, that is expected — it will appear after `wails build` in Task 9. The `getAlerts` fallback using `window.go` avoids the type dependency.)

- [x] **Step 4: Commit**

```bash
git add gover2/frontend/src/lib/api/index.ts
git commit -m "feat(gover2): replace api stubs with real Wails bridge bindings"
```

---

### Task 9: Build Verification + wailsjs Regeneration

**Files:**
- Modified by build: `gover2/frontend/wailsjs/go/bridge/App.d.ts`, `gover2/frontend/wailsjs/go/bridge/App.js`, `gover2/frontend/wailsjs/go/models.ts` (regenerated by `wails build`)

**Interfaces:**
- After this task: `App.d.ts` includes `GetAlerts`; `models.ts` includes `Alert` class; `api/index.ts` can import `GetAlerts` from the generated bindings

- [x] **Step 1: Run go build to check for import cycles**

Run from `gover2/`:
```bash
go build ./...
```

Expected: exit 0.

- [x] **Step 2: Run all Go tests**

Run from `gover2/`:
```bash
go test -cover ./...
```

Expected: all packages pass; `internal/database` >= 70%, `internal/services` >= 70%, `internal/bridge` >= 50%.

- [x] **Step 3: Run wails build to regenerate wailsjs bindings**

Run from `gover2/`:
```bash
wails build
```

Expected: exit 0; binary produced at `gover2/build/bin/gover2`; `wailsjs/go/bridge/App.d.ts` now includes `GetAlerts`.

- [x] **Step 4: Update api/index.ts to import GetAlerts from generated bindings**

Now that `wails build` has regenerated `wailsjs/go/bridge/App.d.ts` with `GetAlerts`, replace the `getAlerts` function in `gover2/frontend/src/lib/api/index.ts`:

Change the import line at the top to add `GetAlerts`:

```ts
import {
  ListComputers,
  ListSmartphones,
  ListTablets,
  ListWindowsKeys,
  ListAntivirus,
  ListOtherSoftware,
  ListUsers,
  GetAlerts,
  GetConfig,
  SetConfig,
  GetDatabasePath,
  OpenDatabase,
  NewDatabase,
  ExportCSV,
} from '../../wailsjs/go/bridge/App'
```

And replace the `getAlerts` export with the clean version:

```ts
export const getAlerts = (): Promise<models.Alert[]> =>
  GetAlerts().catch(() => [])
```

And remove the `import type { Alert } from '@/lib/types'` line (now using `models.Alert`).

- [x] **Step 5: Run typecheck again**

Run from project root:
```bash
pnpm --prefix gover2/frontend run typecheck
```

Expected: exit 0.

- [x] **Step 6: Run all frontend tests**

Run from project root:
```bash
pnpm --prefix gover2/frontend run test
```

Expected: all 13 tests pass.

- [x] **Step 7: Commit**

```bash
git add \
  gover2/frontend/wailsjs/go/bridge/App.d.ts \
  gover2/frontend/wailsjs/go/bridge/App.js \
  gover2/frontend/wailsjs/go/models.ts \
  gover2/frontend/src/lib/api/index.ts
git commit -m "chore(gover2): regenerate wailsjs bindings after adding GetAlerts to bridge"
```

---

### Task 10: Integration Smoke Test

**Files:** none — manual verification only.

**Goal:** Confirm the compiled app opens `inventory.db`, displays real data in all 7 views, and handles the empty/no-DB state gracefully.

- [ ] **Step 1: Launch the app in dev mode**

Run from `gover2/`:
```bash
wails dev
```

Expected: App window opens. Browser console shows no errors. Sidebar displays 7 category links.

- [ ] **Step 2: Verify Computers view**

Click "Computers" in the sidebar.

Expected: Table renders with columns Name / Model / User / Status / Purchase Date / Warranty Expiry. Rows match data in `inventory.db`. StatusBadge shows correct color per status value.

- [ ] **Step 3: Verify remaining 6 views**

Navigate to Smartphones, Tablets, Windows Keys, Antivirus, Other Software, Users in turn.

Expected for each: rows appear (or empty-state "No items to display." if no data — no crash). Columns match the column definitions from Task 7.

- [ ] **Step 4: Verify All Assets view**

Click "All Assets" in the sidebar.

Expected: Category column shows "Computer" / "Smartphone" / "Tablet" values. All three asset types appear as a flat list.

- [ ] **Step 5: Verify no-DB fallback**

Temporarily rename `~/.local/share/inventory/inventory.db` and restart `wails dev`.

Expected: All views show empty-state "No items to display." — no crash, no error dialog.

Restore the file after verifying.

- [ ] **Step 6: Final build**

Run from `gover2/`:
```bash
wails build
```

Expected: exit 0; `gover2/build/bin/gover2` binary produced.

- [ ] **Step 7: Commit if any files changed**

```bash
git status
# Only commit if wails build produced changes to tracked files (e.g., wailsjs/)
git add gover2/build/bin/gover2 2>/dev/null || true
git commit -m "chore(gover2): final build verification for gover2-readonly" || true
```
