# Go Quality Pass Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Raise the Go app's quality — stricter linting, a refactored `bridge.go`, a real SQLite backup feature wired to the UI, and an enforced coverage floor — without changing existing behavior.

**Architecture:** Five ordered tasks on the `feat/go-quality` branch (stacked on `feat/go-default-cutover`), lint-first so refactored and new code conforms to the tighter rules: (1) enable linters + fix, (2) refactor bridge boilerplate with generics, (3) backup backend, (4) backup frontend, (5) coverage + CI floor.

**Tech Stack:** Go 1.25 (local 1.26), Wails v2.12, `modernc.org/sqlite`, golangci-lint v2.12.2, Vue 3 + TypeScript, Vitest, Task.

## Global Constraints

- Module path `inventory`; build tag **`webkit2_41`**; CGO required for packages importing Wails (`internal/bridge`).
- Lint runs via `task lint` (`golangci-lint run --build-tags webkit2_41 ./...`). golangci-lint must be on PATH: `export PATH="$(go env GOPATH)/bin:$PATH"`.
- Tests run via `task test` (`CGO_ENABLED=1 go test -tags webkit2_41 ./...` + vitest). `task check` = lint + typecheck + test.
- Coverage floor: **≥ 75% on `./internal/...`** (excludes Wails-entry `main`/`app.go`).
- Backup mechanism: SQLite **`VACUUM INTO`** on the live connection.
- Bridge **public** method signatures must stay unchanged so `frontend/wailsjs/` bindings and frontend callers are unaffected (the only new public methods are the two backup methods).
- Commits use `--no-verify` (pre-commit hooks require network-installed hook repos; CI is the real gate).

---

### Task 1: Enable stricter linters and fix all findings

**Files:**
- Modify: `.golangci.yml`
- Modify (fixes): any Go file the linters flag

**Interfaces:**
- Consumes: existing `.golangci.yml` (v2, default standard linters).
- Produces: a strict lint config that all later tasks' code must satisfy; `task lint` reports 0 issues.

- [ ] **Step 1: Replace `.golangci.yml` with the stricter set**

```yaml
version: "2"

run:
  build-tags:
    - webkit2_41

linters:
  enable:
    - errorlint
    - revive
    - gocritic
    - misspell
    - unconvert
    - gocyclo
  settings:
    gocyclo:
      min-complexity: 15

issues:
  # Report every finding (defaults cap duplicates at 3 / 50).
  max-same-issues: 0
  max-issues-per-linter: 0
```

(The default standard linters — errcheck, govet, ineffassign, staticcheck, unused — remain enabled implicitly under v2; `enable` adds to them.)

- [ ] **Step 2: Run lint and capture the findings**

Run: `export PATH="$(go env GOPATH)/bin:$PATH" && task lint 2>&1 | tee /tmp/lint-findings.txt`
Expected: a list of findings (or "0 issues"). If 0 issues, skip to Step 4.

- [ ] **Step 3: Fix each finding, then re-run until clean**

Apply the canonical fix per linter, re-running `task lint` after each batch until it reports `0 issues`:
- **errorlint** — comparing/wrapping errors: replace `err == someErr` with `errors.Is(err, someErr)`; replace `fmt.Errorf("...: %v", err)` with `%w`.
- **revive** — usually exported-symbol doc comments or naming: add the missing doc comment in the form `// Name ...`, or rename per the message.
- **gocritic** — apply the suggested rewrite (e.g. `if-else` → `switch`, `append` assignment, single-case switch).
- **misspell** — accept the suggested spelling.
- **unconvert** — delete the redundant conversion.
- **gocyclo** — if a function exceeds complexity 15, extract a helper; if it is an inherently flat query-builder, leave it and add `//nolint:gocyclo // flat query assembly` with that justification on the function.

Re-run: `export PATH="$(go env GOPATH)/bin:$PATH" && task lint`
Expected (final): `0 issues.`

- [ ] **Step 4: Confirm tests still pass**

Run: `export PATH="$(go env GOPATH)/bin:$PATH" && task test`
Expected: all Go packages `ok`, vitest all pass.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit --no-verify -m "style(quality): enable errorlint/revive/gocritic/misspell/unconvert/gocyclo and fix findings"
```

---

### Task 2: Refactor `bridge.go` boilerplate with generics

**Files:**
- Modify: `internal/bridge/bridge.go`
- Test (regression guard, must keep passing): `internal/bridge/bridge_test.go`

**Interfaces:**
- Consumes: `services.List*`/`Insert*`/`Update*`/`Delete*` functions (unchanged); `errNoDB`.
- Produces: `list[T any](a *App, fn func(*sql.DB) ([]T, error)) ([]T, error)` and `func (a *App) withDB(fn func(*sql.DB) error) error`. **No public method signatures change.**

- [ ] **Step 1: Confirm the bridge tests pass (baseline)**

Run: `export PATH="$(go env GOPATH)/bin:$PATH" && CGO_ENABLED=1 go test -tags webkit2_41 ./internal/bridge/...`
Expected: `ok inventory/internal/bridge`. These tests are the regression guard for this refactor.

- [ ] **Step 2: Add the two helpers near the top of `bridge.go` (after `errNoDB`)**

```go
// list runs a read query, guarding the closed-db case and normalizing a nil
// slice to an empty one so the frontend always receives [] rather than null.
func list[T any](a *App, fn func(*sql.DB) ([]T, error)) ([]T, error) {
	if a.db == nil {
		return nil, errNoDB
	}
	out, err := fn(a.db)
	if out == nil {
		out = []T{}
	}
	return out, err
}

// withDB guards the closed-db case for write and side-effecting methods.
func (a *App) withDB(fn func(*sql.DB) error) error {
	if a.db == nil {
		return errNoDB
	}
	return fn(a.db)
}
```

- [ ] **Step 3: Rewrite the read methods using `list`**

Replace the 7 list methods, 1 alerts method, and 2 dropdown methods with these one-liners (delete the old multi-line bodies):

```go
func (a *App) ListComputers() ([]models.Computer, error) { return list(a, services.ListComputers) }
func (a *App) ListSmartphones() ([]models.Smartphone, error) { return list(a, services.ListSmartphones) }
func (a *App) ListTablets() ([]models.Tablet, error) { return list(a, services.ListTablets) }
func (a *App) ListWindowsKeys() ([]models.WindowsKey, error) { return list(a, services.ListWindowsKeys) }
func (a *App) ListAntivirus() ([]models.Antivirus, error) { return list(a, services.ListAntivirus) }
func (a *App) ListOtherSoftware() ([]models.OtherSoftware, error) { return list(a, services.ListOtherSoftware) }
func (a *App) ListUsers() ([]models.User, error) { return list(a, services.ListUsers) }
func (a *App) ListUsersForDropdown() ([]models.DropdownItem, error) { return list(a, services.ListUsersForDropdown) }
func (a *App) ListDevicesForDropdown() ([]models.DeviceDropdownItem, error) { return list(a, services.ListDevicesForDropdown) }

func (a *App) GetAlerts() ([]models.Alert, error) {
	return list(a, func(db *sql.DB) ([]models.Alert, error) {
		return services.GetAlerts(db, a.cfg.ExpiryWarningDays)
	})
}
```

(`GetAlerts` now also normalizes nil→`[]` — a harmless improvement; the frontend already renders an empty list as "no alerts".)

- [ ] **Step 4: Rewrite the write methods using `withDB`**

Replace the 21 Add/Update/Delete methods and `ExportCSV`. Add methods discard the inserted id inside the closure:

```go
func (a *App) AddComputer(data models.ComputerInput) error {
	return a.withDB(func(db *sql.DB) error { _, err := services.InsertComputer(db, data); return err })
}
func (a *App) UpdateComputer(id int, data models.ComputerInput) error {
	return a.withDB(func(db *sql.DB) error { return services.UpdateComputer(db, id, data) })
}
func (a *App) DeleteComputer(id int) error {
	return a.withDB(func(db *sql.DB) error { return services.DeleteComputer(db, id) })
}

func (a *App) AddSmartphone(data models.SmartphoneInput) error {
	return a.withDB(func(db *sql.DB) error { _, err := services.InsertSmartphone(db, data); return err })
}
func (a *App) UpdateSmartphone(id int, data models.SmartphoneInput) error {
	return a.withDB(func(db *sql.DB) error { return services.UpdateSmartphone(db, id, data) })
}
func (a *App) DeleteSmartphone(id int) error {
	return a.withDB(func(db *sql.DB) error { return services.DeleteSmartphone(db, id) })
}

func (a *App) AddTablet(data models.TabletInput) error {
	return a.withDB(func(db *sql.DB) error { _, err := services.InsertTablet(db, data); return err })
}
func (a *App) UpdateTablet(id int, data models.TabletInput) error {
	return a.withDB(func(db *sql.DB) error { return services.UpdateTablet(db, id, data) })
}
func (a *App) DeleteTablet(id int) error {
	return a.withDB(func(db *sql.DB) error { return services.DeleteTablet(db, id) })
}

func (a *App) AddWindowsKey(data models.WindowsKeyInput) error {
	return a.withDB(func(db *sql.DB) error { _, err := services.InsertWindowsKey(db, data); return err })
}
func (a *App) UpdateWindowsKey(id int, data models.WindowsKeyInput) error {
	return a.withDB(func(db *sql.DB) error { return services.UpdateWindowsKey(db, id, data) })
}
func (a *App) DeleteWindowsKey(id int) error {
	return a.withDB(func(db *sql.DB) error { return services.DeleteWindowsKey(db, id) })
}

func (a *App) AddAntivirus(data models.AntivirusInput) error {
	return a.withDB(func(db *sql.DB) error { _, err := services.InsertAntivirus(db, data); return err })
}
func (a *App) UpdateAntivirus(id int, data models.AntivirusInput) error {
	return a.withDB(func(db *sql.DB) error { return services.UpdateAntivirus(db, id, data) })
}
func (a *App) DeleteAntivirus(id int) error {
	return a.withDB(func(db *sql.DB) error { return services.DeleteAntivirus(db, id) })
}

func (a *App) AddOtherSoftware(data models.OtherSoftwareInput) error {
	return a.withDB(func(db *sql.DB) error { _, err := services.InsertOtherSoftware(db, data); return err })
}
func (a *App) UpdateOtherSoftware(id int, data models.OtherSoftwareInput) error {
	return a.withDB(func(db *sql.DB) error { return services.UpdateOtherSoftware(db, id, data) })
}
func (a *App) DeleteOtherSoftware(id int) error {
	return a.withDB(func(db *sql.DB) error { return services.DeleteOtherSoftware(db, id) })
}

func (a *App) AddUser(data models.UserInput) error {
	return a.withDB(func(db *sql.DB) error { _, err := services.InsertUser(db, data); return err })
}
func (a *App) UpdateUser(id int, data models.UserInput) error {
	return a.withDB(func(db *sql.DB) error { return services.UpdateUser(db, id, data) })
}
func (a *App) DeleteUser(id int) error {
	return a.withDB(func(db *sql.DB) error { return services.DeleteUser(db, id) })
}

func (a *App) ExportCSV(category string) error {
	return a.withDB(func(db *sql.DB) error {
		rows, err := services.BuildCSV(db, category)
		if err != nil {
			return err
		}
		path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
			DefaultFilename: category + ".csv",
			Filters:         []runtime.FileFilter{{DisplayName: "CSV", Pattern: "*.csv"}},
		})
		if err != nil {
			return err
		}
		if path == "" {
			return nil // user cancelled
		}
		return services.WriteCSV(path, rows)
	})
}
```

- [ ] **Step 5: Remove the two `var _ =` smell lines at the bottom of `bridge.go`**

Delete:
```go
var _ = backup.BackupDB
var _ = database.Open // keeps modernc.org/sqlite in go.mod
```
The `backup` import stays (used in Task 3); `database` is still used by `Startup`/`NewDatabase`. If `gofmt`/compile reports `backup` as unused at this point (Task 3 not yet done), temporarily keep `var _ = backup.BackupDB` and remove it in Task 3 Step 4 instead.

- [ ] **Step 6: Format, lint, and run the bridge tests**

Run: `export PATH="$(go env GOPATH)/bin:$PATH" && gofmt -w internal/bridge/bridge.go && task lint && CGO_ENABLED=1 go test -tags webkit2_41 ./internal/bridge/...`
Expected: lint `0 issues.`, bridge tests `ok`.

- [ ] **Step 7: Confirm the generated bindings are unchanged**

Run: `wails build -tags webkit2_41 >/dev/null 2>&1; git status --porcelain frontend/wailsjs/`
Expected: no output (bindings unchanged — proves public signatures didn't change). If `models.ts` shows churn, investigate before committing.

- [ ] **Step 8: Commit**

```bash
git add internal/bridge/bridge.go
git commit --no-verify -m "refactor(bridge): collapse db-guard + nil-slice boilerplate into list/withDB helpers"
```

---

### Task 3: Backup backend — real `BackupDB` + bridge methods

**Files:**
- Rewrite: `internal/backup/backup.go`
- Create: `internal/backup/backup_test.go`
- Modify: `internal/bridge/bridge.go` (add two methods)
- Modify: `internal/bridge/bridge_test.go` (add tests)

**Interfaces:**
- Consumes: `*sql.DB`, `errNoDB`, `runtime.SaveFileDialog`.
- Produces: `backup.BackupDB(db *sql.DB, destPath string) error`; bridge `func (a *App) BackupDatabase(destPath string) error` and `func (a *App) BackupDatabaseDialog() error`.

- [ ] **Step 1: Write the failing test for `backup.BackupDB`**

Create `internal/backup/backup_test.go`:

```go
package backup_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"inventory/internal/backup"
	"inventory/internal/database"
)

func TestBackupDB_ProducesOpenableCopyWithSameRows(t *testing.T) {
	src, err := database.Open(filepath.Join(t.TempDir(), "src.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = src.Close() })
	require.NoError(t, database.InitSchema(src))
	_, err = src.Exec(`INSERT INTO computers (name, model, status) VALUES (?, ?, ?)`, "PC1", "Dell", "Active")
	require.NoError(t, err)
	_, err = src.Exec(`INSERT INTO computers (name, model, status) VALUES (?, ?, ?)`, "PC2", "HP", "Active")
	require.NoError(t, err)

	dest := filepath.Join(t.TempDir(), "backup.db")
	require.NoError(t, backup.BackupDB(src, dest))

	copyDB, err := database.Open(dest)
	require.NoError(t, err)
	t.Cleanup(func() { _ = copyDB.Close() })
	var count int
	require.NoError(t, copyDB.QueryRow(`SELECT COUNT(*) FROM computers`).Scan(&count))
	require.Equal(t, 2, count)
}
```

- [ ] **Step 2: Run it to confirm it fails**

Run: `export PATH="$(go env GOPATH)/bin:$PATH" && CGO_ENABLED=1 go test -tags webkit2_41 ./internal/backup/...`
Expected: FAIL — the current `BackupDB(srcPath, destDir string)` signature won't compile against the new call, or the no-op produces no file.

- [ ] **Step 3: Rewrite `internal/backup/backup.go`**

```go
// Package backup creates standalone copies of the live SQLite database.
package backup

import (
	"database/sql"
	"fmt"
	"strings"
)

// BackupDB writes a consistent standalone copy of db to destPath using SQLite's
// VACUUM INTO, which produces a defragmented snapshot regardless of WAL state.
// destPath must not already exist (VACUUM INTO refuses to overwrite).
func BackupDB(db *sql.DB, destPath string) error {
	// VACUUM INTO does not accept a bound parameter for its destination, so the
	// path is inlined with single quotes escaped by doubling.
	escaped := strings.ReplaceAll(destPath, "'", "''")
	if _, err := db.Exec(fmt.Sprintf("VACUUM INTO '%s'", escaped)); err != nil {
		return fmt.Errorf("backup database: %w", err)
	}
	return nil
}
```

- [ ] **Step 4: Add the bridge methods and remove the backup smell line**

In `internal/bridge/bridge.go`, add (near the other DB-management methods) and remove `var _ = backup.BackupDB` if it still exists:

```go
// BackupDatabase writes a standalone copy of the open database to destPath.
func (a *App) BackupDatabase(destPath string) error {
	return a.withDB(func(db *sql.DB) error { return backup.BackupDB(db, destPath) })
}

// BackupDatabaseDialog prompts for a destination and backs up the open
// database there. Returns nil if the user cancels.
func (a *App) BackupDatabaseDialog() error {
	if a.db == nil {
		return errNoDB
	}
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "inventory-backup-" + time.Now().Format("2006-01-02") + ".db",
		Filters:         []runtime.FileFilter{{DisplayName: "SQLite database", Pattern: "*.db"}},
	})
	if err != nil {
		return err
	}
	if path == "" {
		return nil // user cancelled
	}
	return a.BackupDatabase(path)
}
```

Add `"time"` to the bridge imports.

- [ ] **Step 5: Add bridge tests for the backup methods**

Append to `internal/bridge/bridge_test.go`:

```go
func TestApp_BackupDatabase_NoDB(t *testing.T) {
	app := bridge.NewApp()
	require.ErrorIs(t, app.BackupDatabase(filepath.Join(t.TempDir(), "x.db")), bridge.ErrNoDB)
}

func TestApp_BackupDatabase_WritesCopy(t *testing.T) {
	db, err := database.Open(filepath.Join(t.TempDir(), "live.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, database.InitSchema(db))
	_, err = db.Exec(`INSERT INTO users (name, surname) VALUES (?, ?)`, "Alice", "Smith")
	require.NoError(t, err)

	app := bridge.NewAppWithDB(db)
	dest := filepath.Join(t.TempDir(), "copy.db")
	require.NoError(t, app.BackupDatabase(dest))

	copyDB, err := database.Open(dest)
	require.NoError(t, err)
	t.Cleanup(func() { _ = copyDB.Close() })
	var count int
	require.NoError(t, copyDB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count))
	require.Equal(t, 1, count)
}
```

This references `bridge.ErrNoDB`. If `errNoDB` is currently unexported, export it: rename `errNoDB` → `ErrNoDB` in `bridge.go` (all in-package uses update with it). This is an internal package, so no external impact.

- [ ] **Step 6: Run backup + bridge tests**

Run: `export PATH="$(go env GOPATH)/bin:$PATH" && CGO_ENABLED=1 go test -tags webkit2_41 ./internal/backup/... ./internal/bridge/...`
Expected: both `ok`.

- [ ] **Step 7: Lint and build (regenerate bindings for the new methods)**

Run: `export PATH="$(go env GOPATH)/bin:$PATH" && gofmt -w internal/ && task lint && wails build -tags webkit2_41 >/dev/null && git status --porcelain frontend/wailsjs/`
Expected: lint `0 issues.`; `frontend/wailsjs/` now shows the two new methods in the generated bindings (expected churn — commit it).

- [ ] **Step 8: Commit**

```bash
git add internal/backup/ internal/bridge/ frontend/wailsjs/
git commit --no-verify -m "feat(backup): real SQLite backup via VACUUM INTO + BackupDatabase bridge methods"
```

---

### Task 4: Backup frontend — API wrapper + Topbar menu item

**Files:**
- Modify: `frontend/src/lib/api/index.ts`
- Modify: `frontend/src/components/Topbar.vue`
- Modify: `frontend/src/components/Topbar.test.ts`

**Interfaces:**
- Consumes: bridge method `BackupDatabaseDialog` (Task 3).
- Produces: `backupDatabaseDialog(): Promise<void>` and a "Backup database…" menu action.

- [ ] **Step 1: Add the API wrapper**

In `frontend/src/lib/api/index.ts`, next to the `newDatabaseDialog`/`openDatabaseDialog` exports:

```ts
export const backupDatabaseDialog = (): Promise<void> =>
  call<void>('BackupDatabaseDialog')
```

- [ ] **Step 2: Write the failing Topbar test**

In `frontend/src/components/Topbar.test.ts`, add a mock and a test mirroring the existing New/Open ones. Add to the `vi.mock('@/lib/api', ...)` factory:

```ts
const backupDatabaseDialog = vi.fn().mockResolvedValue(undefined)
// inside the mocked module object:
  backupDatabaseDialog: () => backupDatabaseDialog(),
```
And add a test:

```ts
it('backs up the database when the menu item is clicked', async () => {
  const wrapper = mountTopbar()
  await openDbMenu(wrapper) // use the same menu-open helper the New/Open tests use
  await wrapper.get('[data-testid="backup-database"]').trigger('click')
  expect(backupDatabaseDialog).toHaveBeenCalled()
})
```

Match the existing test's helpers for opening the menu and locating items (read the New/Open tests at the top of the file and copy their exact pattern; if they locate by text rather than `data-testid`, use `findByText`/role the same way).

- [ ] **Step 3: Run it to confirm it fails**

Run: `pnpm --prefix frontend run test -- Topbar`
Expected: FAIL — no backup menu item / handler yet.

- [ ] **Step 4: Add the menu item and handler in `Topbar.vue`**

Import the wrapper and add a handler next to `onNewDatabase`/`onOpenDatabase`:

```ts
import { backupDatabaseDialog, exportCSV, newDatabaseDialog, openDatabaseDialog } from '@/lib/api'

async function onBackupDatabase() {
  closeDbMenu() // call whatever the New/Open handlers call to close the menu
  await backupDatabaseDialog()
}
```

Add the menu entry in the DB menu template, next to the New/Open entries, matching their markup exactly (same button/anchor classes and the `data-testid` used in the test):

```html
<button class="<same classes as New/Open>" data-testid="backup-database" @click="onBackupDatabase">
  Backup database…
</button>
```

- [ ] **Step 5: Run the test and typecheck**

Run: `pnpm --prefix frontend run test -- Topbar && pnpm --prefix frontend run typecheck`
Expected: Topbar tests PASS; `vue-tsc` reports no errors.

- [ ] **Step 6: Commit**

```bash
git add frontend/src/lib/api/index.ts frontend/src/components/Topbar.vue frontend/src/components/Topbar.test.ts
git commit --no-verify -m "feat(backup): Backup database… menu action wired to BackupDatabaseDialog"
```

---

### Task 5: Coverage tests + enforced CI floor

**Files:**
- Modify: `internal/database/database_test.go`, `internal/config/config_test.go`, `internal/bridge/bridge_test.go` (fill gaps)
- Modify: `Taskfile.yml` (add `test:cov`)
- Modify: `.gitignore` (ignore `coverage.out`)
- Modify: `.github/workflows/ci.yml` (add coverage gate)

**Interfaces:**
- Consumes: all packages under `./internal/...`.
- Produces: `task test:cov` enforcing ≥ 75% coverage on `./internal/...`, wired into CI.

- [ ] **Step 1: Add the coverage task to `Taskfile.yml`**

After the `test:` task:

```yaml
  test:cov:
    desc: Run Go tests with coverage and enforce the floor (>= 75% of ./internal/...)
    cmds:
      - CGO_ENABLED=1 go test -tags {{.WAILS_TAGS}} -coverpkg=./internal/... -covermode=atomic -coverprofile=coverage.out ./internal/...
      - go tool cover -func=coverage.out | awk '/^total:/ {pct=$3+0; print "total coverage:", $3; if (pct < 75.0) {print "FAIL: coverage below 75%"; exit 1}}'
```

- [ ] **Step 2: Ignore the coverage profile**

Add to `.gitignore` under the Go section:
```gitignore
coverage.out
```

- [ ] **Step 3: Measure the current floor**

Run: `export PATH="$(go env GOPATH)/bin:$PATH" && task test:cov`
Expected: prints `total coverage: NN.N%`. If ≥ 75%, skip to Step 6. If < 75%, continue.

- [ ] **Step 4: Add targeted tests for the lowest-covered paths**

Add the following (each is a real gap identified in the design):

In `internal/config/config_test.go`:
```go
func TestSaveLoad_RoundTrip(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	in := models.AppConfig{DBPath: "/tmp/x.db", DarkMode: true, Density: "compact", ExpiryWarningDays: 14}
	require.NoError(t, config.Save(in))
	got, err := config.Load()
	require.NoError(t, err)
	require.Equal(t, in, got)
}
```
(Adjust the env var to whatever `config.go` uses to locate its file — read `config.go` first and match it; if it uses `os.UserConfigDir`, set `XDG_CONFIG_HOME` on Linux.)

In `internal/database/database_test.go`:
```go
func TestOpen_BadPath(t *testing.T) {
	_, err := database.Open("/nonexistent-dir/sub/dir/x.db")
	require.Error(t, err)
}
```

In `internal/bridge/bridge_test.go`:
```go
func TestApp_GetSetConfig_RoundTrip(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	app := bridge.NewApp()
	cfg := models.AppConfig{Density: "compact", DarkMode: true, ExpiryWarningDays: 7}
	require.NoError(t, app.SetConfig(cfg))
	got, err := app.GetConfig()
	require.NoError(t, err)
	require.Equal(t, cfg, got)
}
```

- [ ] **Step 5: Re-run the coverage gate until it passes**

Run: `export PATH="$(go env GOPATH)/bin:$PATH" && task test:cov`
Expected: `total coverage: NN.N%` with `NN.N ≥ 75` and exit 0. If still short, add a test for the next-lowest uncovered function shown by `go tool cover -func=coverage.out` (target `internal/bridge` write/list paths or `internal/database` pragma paths) and re-run.

- [ ] **Step 6: Wire the coverage gate into CI**

In `.github/workflows/ci.yml`, add a step after the `Test (go test + vitest)` step:

```yaml
      - name: Coverage floor (>= 75% of internal)
        run: task test:cov
```

- [ ] **Step 7: Validate CI YAML and run the full local gate**

Run: `python3 -c "import yaml; yaml.safe_load(open('.github/workflows/ci.yml')); print('ci.yml ok')" && export PATH="$(go env GOPATH)/bin:$PATH" && task check && task test:cov`
Expected: `ci.yml ok`; `task check` all green; `task test:cov` ≥ 75%.

- [ ] **Step 8: Commit**

```bash
git add -A
git commit --no-verify -m "test(quality): fill coverage gaps and enforce a 75% floor on ./internal in CI"
```

---

### Task 6: Open the PR

- [ ] **Step 1: Push and open the PR (base = the cutover branch)**

```bash
git push -u origin feat/go-quality
gh pr create --base feat/go-default-cutover \
  --title "Go quality pass: stricter lint, bridge refactor, real backup, coverage floor" \
  --body "Stricter golangci-lint set, bridge.go boilerplate collapsed into list/withDB helpers, real SQLite backup (VACUUM INTO) wired to a Backup database… menu item, and a 75% coverage floor on ./internal enforced in CI. Stacked on #7. See docs/superpowers/specs/2026-06-20-go-quality-design.md."
```

(Once #7 merges, retarget this PR's base to `main`.)

---

## Self-Review

**Spec coverage** (design §4 → tasks):
- §4.1 stricter linting → Task 1 ✓
- §4.2 bridge refactor (list/withDB, remove `var _ =`) → Task 2 ✓
- §4.3 backup (backup.go + bridge methods + frontend) → Task 3 (backend) ✓ + Task 4 (frontend) ✓
- §4.4 coverage + CI floor → Task 5 ✓
- §5 acceptance: 0 lint issues (T1/T2/T3), bridge <250 lines & signatures unchanged (T2 Step 7), backup openable + wired (T3/T4), `task test:cov` ≥75% in CI (T5) ✓

**Placeholder scan:** Task 1's fix loop is tool-output-driven (the linter output is the exact instruction) with canonical fixes per linter — not a vague placeholder. Task 4 references "match the existing New/Open test pattern" because the test-helper names live in a file the implementer will read; the exact assertion and data-testid are given. No "TODO"/"implement later".

**Type/name consistency:** `BackupDB(db *sql.DB, destPath string)` consistent across §4.3, Task 3 Step 3/Step 1. Bridge methods `BackupDatabase`/`BackupDatabaseDialog` consistent across Tasks 3–4 and the API wrapper `backupDatabaseDialog`. `errNoDB`→`ErrNoDB` export noted in Task 3 Step 5 (used by Task 5 bridge test). `list`/`withDB` helper signatures consistent between Task 2 Step 2 and their call sites in Steps 3–4. Build tag `webkit2_41`, coverage floor 75%, and `coverpkg=./internal/...` consistent across Task 5 and Global Constraints.
