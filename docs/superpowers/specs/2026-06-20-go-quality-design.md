# Design: Go quality pass

**Date:** 2026-06-20
**Status:** Approved (brainstorming) — pending implementation plan
**Branch:** `feat/go-quality`, stacked on `feat/go-default-cutover` (PR #7), its own PR.
**Topic:** Raise the quality of the Go/Wails app along four axes — stricter linting, bridge refactor, a real database backup feature, and test coverage with a CI floor.

---

## 1. Goal

After the cutover, the Go app is the sole `inventory` product. This pass improves its internal quality without changing observable behavior (except adding a backup action):

- **Stricter static analysis** — go beyond golangci-lint's default linters and fix what they surface.
- **Refactor `bridge.go`** — collapse ~40 repetitions of two boilerplate patterns.
- **Implement real backup** — replace the dead `backup.BackupDB` no-op stub with a working, UI-wired database backup.
- **Test coverage** — fill the gaps and enforce a floor in CI.

## 2. Decisions

| # | Decision | Choice |
|---|----------|--------|
| 1 | Scope | All four areas: linting, coverage, bridge refactor, backup. |
| 2 | Backup | Implement for real, **backup-only** (restore is already covered by *File → Open Database…*). |
| 3 | Branch | `feat/go-quality` stacked on `feat/go-default-cutover`; separate PR. |
| 4 | Sequencing | Lint-first, one branch, four ordered phases. |
| 5 | Coverage floor | 75% on `./internal/...` (excludes Wails-entry `main`/`app.go`). |
| 6 | Backup mechanism | SQLite `VACUUM INTO` on the live connection. |

## 3. Approach

**Lint-first, four ordered phases on one branch**, each ending green:

1. Stricter lint config + fix existing findings
2. Refactor `bridge.go` boilerplate
3. Implement backup (Go + binding + UI + tests)
4. Coverage tests + CI coverage floor

Lint-first so the refactored and new code conforms to the tighter rules from the start. Rejected: feature-first (backup written before lint tightens → rework), and splitting into two PRs (more overhead than requested).

## 4. Detailed changes by area

### 4.1 Stricter linting

Extend `.golangci.yml` (v2 schema) beyond the default standard set (errcheck, govet, ineffassign, staticcheck, unused) with:

- **`errorlint`** — correct error wrapping/comparison (`errors.Is`/`%w`)
- **`revive`** — general style (golint successor)
- **`gocritic`** — opinionated diagnostics
- **`misspell`** — typos in comments/strings
- **`unconvert`** — redundant type conversions
- **`gocyclo`** — cyclomatic complexity, threshold 15

`gosec` is intentionally **excluded** — low value for a local single-user desktop app and it conflicts with the backup SQL (§4.3); can be added later.

Then fix every finding these surface across the ~2,700 lines of Go. The exact findings are discovered during implementation; the plan handles them as a fix-until-clean loop.

### 4.2 Refactor `bridge.go` (439 → ~200 lines)

`bridge.go` repeats two patterns ~40 times. Collapse with Go generics:

```go
// list runs a read query, normalizing a nil slice to an empty one so the
// frontend always receives [] rather than null, and guarding the closed-db case.
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

// withDB guards the closed-db case for write/side-effecting methods.
func (a *App) withDB(fn func(*sql.DB) error) error {
	if a.db == nil {
		return errNoDB
	}
	return fn(a.db)
}
```

- The 8 list/dropdown methods become one-liners, e.g.
  `func (a *App) ListComputers() ([]models.Computer, error) { return list(a, services.ListComputers) }`
- The 21 Add/Update/Delete methods become `withDB` closures, e.g.
  `func (a *App) DeleteComputer(id int) error { return a.withDB(func(db *sql.DB) error { return services.DeleteComputer(db, id) }) }`
  Add methods discard the returned id inside the closure.
- Delete both smell lines: `var _ = backup.BackupDB` (replaced by a real call in §4.3) and `var _ = database.Open` (dead — `database.Open` is already used in `Startup`/`NewDatabase`; the sqlite driver remains imported transitively through the `database` package).
- **Public method signatures are unchanged**, so the generated Wails bindings (`frontend/wailsjs/`) and all frontend callers are untouched. `ExportCSV` keeps its current dialog body (it does more than a plain service call).

### 4.3 Implement backup

**`internal/backup/backup.go`:**

```go
// BackupDB writes a consistent standalone copy of the live database to
// destPath using SQLite's VACUUM INTO, which works regardless of WAL state.
func BackupDB(db *sql.DB, destPath string) error
```

Primary implementation: `db.Exec("VACUUM INTO ?", destPath)`. If the driver (`modernc.org/sqlite`) rejects a bound parameter in `VACUUM INTO`, fall back to a literal with single quotes doubled: `VACUUM INTO '<escaped>'`. The plan verifies which form works.

**Bridge (`internal/bridge/bridge.go`):** mirror the existing `NewDatabase`/`NewDatabaseDialog` split:

```go
func (a *App) BackupDatabase(destPath string) error      // testable core (guards db==nil, calls backup.BackupDB)
func (a *App) BackupDatabaseDialog() error               // SaveFileDialog, default inventory-backup-YYYY-MM-DD.db, then BackupDatabase
```

`BackupDatabaseDialog` returns `nil` on user-cancel (empty path), consistent with `NewDatabaseDialog`/`ExportCSV`.

**Frontend:**
- `frontend/src/lib/api/index.ts`: `export const backupDatabaseDialog = (): Promise<void> => call<void>('BackupDatabaseDialog')`
- `frontend/src/components/Topbar.vue`: add a "Backup database…" item to the existing DB menu (next to New/Open), with an `onBackupDatabase` handler that calls `backupDatabaseDialog()` (same pattern as `onNewDatabase`).
- `frontend/src/components/Topbar.test.ts`: mock `backupDatabaseDialog` and assert the menu item triggers it.

### 4.4 Test coverage + CI floor

Add Go tests for the gaps:

- **`internal/backup`** (0% → covered): back up a temp-file DB seeded with rows, then open the destination and assert the schema + row counts match.
- **`internal/bridge`** (60%): `BackupDatabase` happy path + `errNoDB` path; `GetConfig`/`SetConfig`; any uncovered list/write paths after the refactor.
- **`internal/database`** (67%): `Open` error path (unwritable path), pragma/schema paths.
- **`internal/config`** (70%): `Load`/`Save` round-trip and error paths.

**Coverage gate:**

```yaml
# Taskfile
test:cov:
  desc: Run Go tests with coverage and enforce the floor
  cmds:
    - CGO_ENABLED=1 go test -tags webkit2_41 -coverpkg=./internal/... -covermode=atomic -coverprofile=coverage.out ./internal/...
    - go tool cover -func=coverage.out | awk '/^total:/ {print; if ($3+0 < 75.0) {print "coverage below 75%"; exit 1}}'
```

Add `coverage.out` to `.gitignore`. Wire `task test:cov` into `ci.yml` as a step after `task test`. The 75% floor is a ratchet — raise it as coverage grows.

## 5. Acceptance criteria

- `.golangci.yml` enables the §4.1 linter set; `task lint` passes with **0 issues**.
- `bridge.go` is under ~250 lines, with `list`/`withDB` helpers and no `var _ =` lines; public method signatures unchanged; `frontend/wailsjs/` unchanged.
- `BackupDB` produces a valid, openable copy; `BackupDatabaseDialog` is exposed and wired to a "Backup database…" menu item; `task build` succeeds.
- `task test:cov` reports `./internal/...` coverage ≥ 75% and is enforced in CI.
- `task check` (lint + typecheck + test) and the frontend tests all pass.

## 6. Out of scope

- Database **restore** UI (covered by existing *Open Database…*).
- `gosec` and other security linters.
- Backup scheduling/automation, cloud backup, or multiple backup slots.
- Any change to bridge **public** method signatures or to UI behavior beyond the new backup menu item.
- Testing `main`/`app.go` (Wails entry; excluded from the coverage floor).
