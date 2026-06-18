---
change: gover2-parity
design-doc: docs/superpowers/specs/2026-06-18-gover2-parity-design.md
base-ref: 2059dd97252d640743a44e967951fb2fa3124136
---

# gover2 Parity Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Bring the gover2 Wails app to feature parity with the PyQt inventory app — CSV export, working dark mode + density, startup expiry alerts, client-side search, an action-row topbar, and database new/open + About/License modals.

**Architecture:** Reuse existing bridge methods; the only new Go code is CSV export (a pure `services/export.go` plus a thin dialog wrapper in `bridge`). On the frontend, dark mode flips a `.dark` class on `<html>` against CSS-variable Tailwind tokens, search is coordinated through a tiny Pinia store + a pure `matchesQuery` helper, and `AppShell.vue` owns the alerts modal so it can auto-open on startup. See the design doc for the full *how*.

**Tech Stack:** Go 1.x + `database/sql` + `modernc.org/sqlite` + Wails v2.12.0 (`github.com/wailsapp/wails/v2/pkg/runtime`); Vue 3 + Pinia + vue-router + Tailwind; tests via Go `testing`/testify and Vitest + `@vue/test-utils` (jsdom).

## Global Constraints

- All bridge calls go through `frontend/src/lib/api/index.ts` (`call<T>` proxy). Components MUST NOT call `window.go.*` directly. — design doc "Architecture".
- Go dependency chain stays `models ← database ← services ← bridge`. No schema change, no `AppConfig` schema change. — design doc "Non-Goals".
- Go test coverage targets: services ≥ 70%, bridge ≥ 50% (existing thresholds). — design doc "Testing Strategy".
- In-memory SQLite for Go tests via `database.Open(":memory:")` + `database.InitSchema(db)`. — `internal/bridge/bridge_test.go`.
- Tailwind token **names** stay identical when converted to CSS variables — no per-component class churn. — design doc "Dark Mode".
- Theme-agnostic colors stay hardcoded: modal overlays (`bg-black/40`), the red delete button (`bg-red-600`), status badge colors. — design doc "Dark Mode".
- CSV export is strictly per-category; the export action is hidden on `/all`. Search still works on `/all`. — design doc "Topbar Action Row".
- Working directory for all Go commands: `/home/mg/inventory/gover2`. For all frontend commands: `/home/mg/inventory/gover2/frontend`.
- Frontend commands: tests `pnpm test` (alias of `vitest run`); typecheck `pnpm run typecheck` (`vue-tsc --noEmit`).

---

## File Structure

**Go (under `gover2/internal/`)**
- Create `services/export.go` — pure `BuildCSV` + exported `WriteCSV`.
- Create `services/export_test.go` — table-driven per-category CSV tests.
- Modify `bridge/bridge.go:172-174` — reshape `ExportCSV` to dialog + write.
- Modify `bridge/bridge_test.go` — add `BuildCSV`/`WriteCSV`/`ExportCSV` guard tests (or a new `services/export_test.go` for the pure parts).

**Frontend (under `gover2/frontend/src/`)**
- Modify `frontend/tailwind.config.ts` — themeable tokens become `var(--c-*)`; add `darkMode: 'class'`.
- Modify `style.css` — `:root` (light) + `.dark` (dark) variable sets.
- Modify `stores/ui.ts` — load from config, persist on toggle/density, apply `.dark`.
- Create `stores/search.ts` — `query` ref + `setQuery`/`clear`.
- Create `lib/filter.ts` — pure `matchesQuery`.
- Create `lib/routeCategory.ts` — route→category map + helper.
- Modify `components/Topbar.vue` — action row.
- Modify `components/DataTable.vue` — density + zebra.
- Modify `components/AlertsModal.vue` — render real alert rows.
- Create `components/AboutModal.vue`, `components/LicenseModal.vue`.
- Modify `layouts/AppShell.vue` — own alerts modal, auto-open, apply theme, clear search on route change.
- Modify all 8 views (`features/assets/*`, `features/licenses/*`, `features/users/UsersView.vue`) — wire `filteredRows`.
- Modify `lib/api/index.ts` — `exportCSV(category)` (drop `destPath`).
- Create test files: `lib/filter.test.ts`, `stores/ui.test.ts`, `components/AlertsModal.test.ts`, `components/AboutModal.test.ts`, `components/LicenseModal.test.ts`; extend `components/DataTable.test.ts`.

---

## Task 1: Go CSV export — pure `BuildCSV` + `WriteCSV`

Implements design doc "CSV Export — Go-driven" (the pure, testable half) and tasks.md 1.1.

**Files:**
- Create: `gover2/internal/services/export.go`
- Test: `gover2/internal/services/export_test.go`

**Interfaces:**
- Consumes: existing `services.List*` functions (`internal/services/services.go`) and `database.Open`/`database.InitSchema` (for tests).
- Produces:
  - `func BuildCSV(db *sql.DB, category string) ([][]string, error)` — header row + data rows in canonical column order; unknown category → error.
  - `func WriteCSV(path string, rows [][]string) error` — writes rows with `encoding/csv`.
  - Canonical categories and column order (must match the view columns):
    - `computers` / `smartphones` / `tablets`: `["Name","Model","Status","Purchase Date","Warranty Expiry"]`
    - `windowskeys`: `["License Key","Status"]`
    - `antivirus` / `othersoftware`: `["Name","License Key","Status","Expiry Date"]`
    - `users`: `["Name","Surname","Status"]`

- [x] **Step 1: Write the failing test**

Create `gover2/internal/services/export_test.go`:

```go
package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gover2/internal/database"
	"gover2/internal/models"
	"gover2/internal/services"
)

func newDB(t *testing.T) *sql.DB { // placeholder; replaced below
	return nil
}
```

Replace the whole file body with the real version (uses `database/sql`):

```go
package services_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gover2/internal/database"
	"gover2/internal/models"
	"gover2/internal/services"
)

func newExportDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, database.InitSchema(db))
	t.Cleanup(func() { db.Close() })
	return db
}

func TestBuildCSV_Computers(t *testing.T) {
	db := newExportDB(t)
	_, err := services.InsertComputer(db, models.ComputerInput{Name: "PC1", Model: "Dell", Status: "active"})
	require.NoError(t, err)

	rows, err := services.BuildCSV(db, "computers")
	require.NoError(t, err)
	require.Len(t, rows, 2) // header + 1 data row
	assert.Equal(t, []string{"Name", "Model", "Status", "Purchase Date", "Warranty Expiry"}, rows[0])
	assert.Equal(t, "PC1", rows[1][0])
	assert.Equal(t, "Dell", rows[1][1])
	assert.Equal(t, "active", rows[1][2])
}

func TestBuildCSV_Users(t *testing.T) {
	db := newExportDB(t)
	_, err := services.InsertUser(db, models.UserInput{Name: "Alice", Status: "active"})
	require.NoError(t, err)

	rows, err := services.BuildCSV(db, "users")
	require.NoError(t, err)
	assert.Equal(t, []string{"Name", "Surname", "Status"}, rows[0])
	assert.Equal(t, "Alice", rows[1][0])
}

func TestBuildCSV_Antivirus(t *testing.T) {
	db := newExportDB(t)
	_, err := services.InsertAntivirus(db, models.AntivirusInput{Name: "Norton", LicenseKey: "K1", Status: "active"})
	require.NoError(t, err)

	rows, err := services.BuildCSV(db, "antivirus")
	require.NoError(t, err)
	assert.Equal(t, []string{"Name", "License Key", "Status", "Expiry Date"}, rows[0])
	assert.Equal(t, "Norton", rows[1][0])
	assert.Equal(t, "K1", rows[1][1])
}

func TestBuildCSV_UnknownCategory(t *testing.T) {
	db := newExportDB(t)
	_, err := services.BuildCSV(db, "nope")
	assert.Error(t, err)
}

func TestWriteCSV_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.csv")
	rows := [][]string{{"A", "B"}, {"1", "2"}}
	require.NoError(t, services.WriteCSV(path, rows))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "A,B\n1,2\n", string(data))
}
```

(Delete the placeholder `newDB`/first import block — only the second full version remains.)

- [x] **Step 2: Run the test to verify it fails**

Run: `cd /home/mg/inventory/gover2 && go test ./internal/services/ -run 'BuildCSV|WriteCSV' -v`
Expected: FAIL — `undefined: services.BuildCSV` / `services.WriteCSV`.

- [x] **Step 3: Write the minimal implementation**

Create `gover2/internal/services/export.go`. Use the existing `List*` services and dereference nullable string pointers with a tiny helper:

```go
package services

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
)

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// BuildCSV returns header + data rows for a category in canonical column order.
func BuildCSV(db *sql.DB, category string) ([][]string, error) {
	switch category {
	case "computers":
		items, err := ListComputers(db)
		if err != nil {
			return nil, err
		}
		out := [][]string{{"Name", "Model", "Status", "Purchase Date", "Warranty Expiry"}}
		for _, c := range items {
			out = append(out, []string{c.Name, c.Model, c.Status, deref(c.PurchaseDate), deref(c.WarrantyExpiry)})
		}
		return out, nil
	case "smartphones":
		items, err := ListSmartphones(db)
		if err != nil {
			return nil, err
		}
		out := [][]string{{"Name", "Model", "Status", "Purchase Date", "Warranty Expiry"}}
		for _, c := range items {
			out = append(out, []string{c.Name, c.Model, c.Status, deref(c.PurchaseDate), deref(c.WarrantyExpiry)})
		}
		return out, nil
	case "tablets":
		items, err := ListTablets(db)
		if err != nil {
			return nil, err
		}
		out := [][]string{{"Name", "Model", "Status", "Purchase Date", "Warranty Expiry"}}
		for _, c := range items {
			out = append(out, []string{c.Name, c.Model, c.Status, deref(c.PurchaseDate), deref(c.WarrantyExpiry)})
		}
		return out, nil
	case "windowskeys":
		items, err := ListWindowsKeys(db)
		if err != nil {
			return nil, err
		}
		out := [][]string{{"License Key", "Status"}}
		for _, w := range items {
			out = append(out, []string{w.LicenseKey, w.Status})
		}
		return out, nil
	case "antivirus":
		items, err := ListAntivirus(db)
		if err != nil {
			return nil, err
		}
		out := [][]string{{"Name", "License Key", "Status", "Expiry Date"}}
		for _, a := range items {
			out = append(out, []string{a.Name, a.LicenseKey, a.Status, deref(a.ExpiryDate)})
		}
		return out, nil
	case "othersoftware":
		items, err := ListOtherSoftware(db)
		if err != nil {
			return nil, err
		}
		out := [][]string{{"Name", "License Key", "Status", "Expiry Date"}}
		for _, s := range items {
			out = append(out, []string{s.Name, s.LicenseKey, s.Status, deref(s.ExpiryDate)})
		}
		return out, nil
	case "users":
		items, err := ListUsers(db)
		if err != nil {
			return nil, err
		}
		out := [][]string{{"Name", "Surname", "Status"}}
		for _, u := range items {
			out = append(out, []string{u.Name, deref(u.Surname), u.Status})
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unknown category: %s", category)
	}
}

// WriteCSV writes rows to path using encoding/csv.
func WriteCSV(path string, rows [][]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	if err := w.WriteAll(rows); err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}
```

- [x] **Step 4: Run the test to verify it passes**

Run: `cd /home/mg/inventory/gover2 && go test ./internal/services/ -run 'BuildCSV|WriteCSV' -v`
Expected: PASS (all 5 tests).

- [x] **Step 5: Verify package + coverage still green**

Run: `cd /home/mg/inventory/gover2 && go test ./internal/services/ -cover`
Expected: PASS, coverage ≥ 70%.

- [x] **Step 6: Commit**

```bash
cd /home/mg/inventory && git add gover2/internal/services/export.go gover2/internal/services/export_test.go
git commit -m "feat(gover2): add pure BuildCSV/WriteCSV services"
```

---

## Task 2: Go bridge — reshape `ExportCSV(category)` with save dialog

Implements design doc "CSV Export — Go-driven" (bridge half) and tasks.md 1.2 + 1.3.

**Files:**
- Modify: `gover2/internal/bridge/bridge.go:172-174` (the `ExportCSV` stub) and the import block at `:3-15`
- Test: `gover2/internal/bridge/bridge_test.go`

**Interfaces:**
- Consumes: `services.BuildCSV`, `services.WriteCSV` (Task 1); `runtime.SaveFileDialog` from `github.com/wailsapp/wails/v2/pkg/runtime`.
- Produces: `func (a *App) ExportCSV(category string) error` — nil-db guard, builds rows, opens save dialog, returns nil on cancel, writes file otherwise.

- [x] **Step 1: Write the failing tests**

Append to `gover2/internal/bridge/bridge_test.go`:

```go
func TestExportCSV_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.ExportCSV("computers")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no database")
}

func TestExportCSV_BuildRows_Computers(t *testing.T) {
	// Verifies BuildCSV wiring at the services boundary used by ExportCSV.
	app := newBridgeWithDB(t)
	require.NoError(t, app.AddComputer(models.ComputerInput{Name: "PC1", Model: "Dell", Status: "active"}))

	rows, err := services.BuildCSV(dbFromApp(t, app), "computers")
	require.NoError(t, err)
	require.Len(t, rows, 2)
	assert.Equal(t, "PC1", rows[1][0])
}
```

The `ExportCSV` dialog path itself cannot run headless, so it is left to the smoke test (Task 14). To keep the second test independent of unexported fields, prefer asserting through a DB the test owns instead of reaching into `app`. Replace the second test with this self-contained form:

```go
func TestExportCSV_BuildRows_Computers(t *testing.T) {
	db, err := database.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, database.InitSchema(db))
	t.Cleanup(func() { db.Close() })

	_, err = services.InsertComputer(db, models.ComputerInput{Name: "PC1", Model: "Dell", Status: "active"})
	require.NoError(t, err)

	rows, err := services.BuildCSV(db, "computers")
	require.NoError(t, err)
	require.Len(t, rows, 2)
	assert.Equal(t, "PC1", rows[1][0])
}
```

Add `"gover2/internal/services"` to the test imports.

- [x] **Step 2: Run the tests to verify they fail**

Run: `cd /home/mg/inventory/gover2 && go test ./internal/bridge/ -run 'ExportCSV' -v`
Expected: FAIL — `app.ExportCSV` still has 2-arg signature (`too many arguments`) / compile error.

- [x] **Step 3: Update the bridge import block**

In `gover2/internal/bridge/bridge.go`, add the runtime import to the block at lines 3-15:

```go
import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gover2/internal/backup"
	"gover2/internal/config"
	"gover2/internal/database"
	"gover2/internal/models"
	"gover2/internal/services"
)
```

- [x] **Step 4: Replace the `ExportCSV` stub**

Replace lines 172-174:

```go
func (a *App) ExportCSV(category string) error {
	if a.db == nil {
		return errNoDB
	}
	rows, err := services.BuildCSV(a.db, category)
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
}
```

- [x] **Step 5: Run the tests to verify they pass**

Run: `cd /home/mg/inventory/gover2 && go test ./internal/bridge/ -run 'ExportCSV' -v`
Expected: PASS.

- [x] **Step 6: Build + full Go test + coverage**

Run: `cd /home/mg/inventory/gover2 && go build ./... && go test ./... -cover`
Expected: PASS; bridge ≥ 50%, services ≥ 70%.

- [x] **Step 7: Commit**

```bash
cd /home/mg/inventory && git add gover2/internal/bridge/bridge.go gover2/internal/bridge/bridge_test.go
git commit -m "feat(gover2): reshape ExportCSV to category-only with save dialog"
```

---

## Task 3: Tailwind CSS-variable tokens + dark palette

Implements design doc "Dark Mode — CSS-variable tokens" and tasks.md 2.2 (token migration half).

**Files:**
- Modify: `gover2/frontend/tailwind.config.ts`
- Modify: `gover2/frontend/src/style.css`

**Interfaces:**
- Produces: themeable Tailwind tokens resolve to `var(--c-*)`; `:root` holds light values, `.dark` holds dark values; `darkMode: 'class'` enabled. Token **names** unchanged so no component edits are required. New token `row-alt` → `var(--c-row-alt)` for zebra (consumed by Task 8).

- [x] **Step 1: Enable class dark mode and convert tokens in `tailwind.config.ts`**

Set `darkMode: 'class'` at config top level and change the themeable color values to variable references. Sidebar palette and the semantic status colors (`s-*`) stay as literal hex (theme-agnostic per the design doc). Replace the `colors` block and add `darkMode`:

```ts
export default {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // Sidebar palette (constant across themes)
        'sidebar':        '#0E1520',
        'sidebar-hover':  '#131E2E',
        'sidebar-active': '#182337',
        'sidebar-border': '#141D2B',
        'sidebar-text':   '#6B82A0',
        'sidebar-hi':     '#DDE8F5',
        'sidebar-accent': '#2563EB',
        // Page palette (themed)
        'page':    'var(--c-page)',
        'surface': 'var(--c-surface)',
        'border':  'var(--c-border)',
        // Status colors (semantic, constant)
        's-active':   '#2563EB',
        's-repair':   '#EA580C',
        's-spare':    '#7C3AED',
        's-retired':  '#9CA3AF',
        's-missing':  '#1E293B',
        's-expiring': '#D97706',
        's-expired':  '#DC2626',
        // Table header (themed)
        'th-bg':   'var(--c-th-bg)',
        'th-text': 'var(--c-th-text)',
        // Typography (themed)
        'text-primary':   'var(--c-text-primary)',
        'text-secondary': 'var(--c-text-secondary)',
        'text-tertiary':  'var(--c-text-tertiary)',
        // Interactive
        'accent-hover': '#1D4ED8',
        'accent':       '#2563EB',
        // Page structure (themed)
        'border-light': 'var(--c-border-light)',
        // Zebra row (themed)
        'row-alt': 'var(--c-row-alt)',
        // Status row highlight backgrounds (themed)
        'row-expiring': 'var(--c-row-expiring)',
        'row-expired':  'var(--c-row-expired)',
      },
      // fontFamily, borderRadius, width, height blocks unchanged
      fontFamily: {
        ui:   ['-apple-system', 'BlinkMacSystemFont', '"Segoe UI"', 'system-ui', 'sans-serif'],
        mono: ['"Cascadia Code"', '"Fira Code"', 'Consolas', '"Courier New"', 'monospace'],
      },
      borderRadius: { badge: '2px' },
      width: { sidebar: '200px' },
      height: { topbar: '48px' },
    },
  },
  plugins: [],
} satisfies Config
```

(Note: `accent` is added because `ComputersView.vue` already uses `bg-accent`; it was previously missing. Keeping it constant.)

- [x] **Step 2: Define the variable sets in `style.css`**

Replace the whole file:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

:root {
  --c-page: #EEF1F7;
  --c-surface: #FFFFFF;
  --c-border: #DDE1EC;
  --c-border-light: #EAECF4;
  --c-th-bg: #0E1520;
  --c-th-text: #8FA5BF;
  --c-text-primary: #0F172A;
  --c-text-secondary: #64748B;
  --c-text-tertiary: #94A3B8;
  --c-row-alt: #F6F8FC;
  --c-row-expiring: #FFFBEB;
  --c-row-expired: #FFF5F5;
}

.dark {
  --c-page: #0A0F18;
  --c-surface: #0E1520;
  --c-border: #1E293B;
  --c-border-light: #16202F;
  --c-th-bg: #060A12;
  --c-th-text: #6B82A0;
  --c-text-primary: #DDE8F5;
  --c-text-secondary: #8FA5BF;
  --c-text-tertiary: #64748B;
  --c-row-alt: #131C2A;
  --c-row-expiring: #2A2410;
  --c-row-expired: #2A1414;
}
```

- [x] **Step 3: Typecheck (config compiles)**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm run typecheck`
Expected: PASS (no type errors in `tailwind.config.ts`).

- [x] **Step 4: Commit**

```bash
cd /home/mg/inventory && git add gover2/frontend/tailwind.config.ts gover2/frontend/src/style.css
git commit -m "feat(gover2): convert Tailwind tokens to CSS variables with dark palette"
```

---

## Task 4: `stores/ui.ts` — load, toggle, persist, apply `.dark`

Implements design doc "Dark Mode — `stores/ui.ts`" and tasks.md 2.1 + 2.4 (ui-store half).

**Files:**
- Modify: `gover2/frontend/src/stores/ui.ts`
- Test: `gover2/frontend/src/stores/ui.test.ts` (create)

**Interfaces:**
- Consumes: `getConfig`, `setConfig`, type `AppConfig` from `@/lib/api`.
- Produces (store `ui`):
  - state `density: 'comfortable' | 'compact'`, `darkMode: boolean`, `sidebarCollapsed: boolean`.
  - `applyTheme()` — adds/removes `dark` class on `document.documentElement`.
  - `loadFromConfig(): Promise<void>` — reads `getConfig()`, sets `density`/`darkMode`, applies theme.
  - `toggleDarkMode()` — flips `darkMode`, applies theme, persists via `setConfig`.
  - `setDensity(d)` — sets `density`, persists via `setConfig`.
  - `persist()` — internal: `setConfig({ ...currentConfig, density, darkMode })`.

- [x] **Step 1: Write the failing test**

Create `gover2/frontend/src/stores/ui.test.ts`:

```ts
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

const getConfig = vi.fn()
const setConfig = vi.fn().mockResolvedValue(undefined)

vi.mock('@/lib/api', () => ({
  getConfig: () => getConfig(),
  setConfig: (cfg: unknown) => setConfig(cfg),
}))

import { useUiStore } from './ui'

describe('ui store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    document.documentElement.classList.remove('dark')
    getConfig.mockReset().mockResolvedValue({
      dbPath: '/x', density: 'compact', darkMode: true, expiryWarningDays: 30,
    })
    setConfig.mockClear()
  })

  it('loadFromConfig applies darkMode + density and sets the dark class', async () => {
    const ui = useUiStore()
    await ui.loadFromConfig()
    expect(ui.darkMode).toBe(true)
    expect(ui.density).toBe('compact')
    expect(document.documentElement.classList.contains('dark')).toBe(true)
  })

  it('toggleDarkMode flips the ref, toggles the class, and persists', async () => {
    const ui = useUiStore()
    await ui.loadFromConfig()
    await ui.toggleDarkMode()
    expect(ui.darkMode).toBe(false)
    expect(document.documentElement.classList.contains('dark')).toBe(false)
    expect(setConfig).toHaveBeenCalledWith(
      expect.objectContaining({ darkMode: false, density: 'compact' }),
    )
  })

  it('setDensity updates and persists', async () => {
    const ui = useUiStore()
    await ui.loadFromConfig()
    await ui.setDensity('comfortable')
    expect(ui.density).toBe('comfortable')
    expect(setConfig).toHaveBeenLastCalledWith(
      expect.objectContaining({ density: 'comfortable' }),
    )
  })
})
```

- [x] **Step 2: Run the test to verify it fails**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm test src/stores/ui.test.ts`
Expected: FAIL — `loadFromConfig` / `applyTheme` not exported.

- [x] **Step 3: Rewrite `stores/ui.ts`**

```ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getConfig, setConfig } from '@/lib/api'
import type { AppConfig } from '@/lib/api'

export const useUiStore = defineStore('ui', () => {
  const density = ref<'comfortable' | 'compact'>('comfortable')
  const darkMode = ref(false)
  const sidebarCollapsed = ref(false)
  const config = ref<AppConfig>({
    dbPath: '', density: 'comfortable', darkMode: false, expiryWarningDays: 30,
  })

  function applyTheme() {
    const root = document.documentElement
    if (darkMode.value) root.classList.add('dark')
    else root.classList.remove('dark')
  }

  async function persist() {
    config.value = { ...config.value, density: density.value, darkMode: darkMode.value }
    await setConfig(config.value)
  }

  async function loadFromConfig() {
    const cfg = await getConfig()
    config.value = cfg
    density.value = cfg.density === 'compact' ? 'compact' : 'comfortable'
    darkMode.value = cfg.darkMode
    applyTheme()
  }

  async function toggleDarkMode() {
    darkMode.value = !darkMode.value
    applyTheme()
    await persist()
  }

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  async function setDensity(d: 'comfortable' | 'compact') {
    density.value = d
    await persist()
  }

  return {
    density,
    darkMode,
    sidebarCollapsed,
    applyTheme,
    loadFromConfig,
    toggleDarkMode,
    toggleSidebar,
    setDensity,
  }
})
```

- [x] **Step 4: Run the test to verify it passes**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm test src/stores/ui.test.ts`
Expected: PASS (3 tests).

- [x] **Step 5: Typecheck**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm run typecheck`
Expected: PASS.

- [x] **Step 6: Commit**

```bash
cd /home/mg/inventory && git add gover2/frontend/src/stores/ui.ts gover2/frontend/src/stores/ui.test.ts
git commit -m "feat(gover2): ui store loads/persists theme + density and applies dark class"
```

---

## Task 5: Call `loadFromConfig` on app startup

Implements design doc "Dark Mode — single source of truth on app start" and tasks.md 2.1.

**Files:**
- Modify: `gover2/frontend/src/layouts/AppShell.vue`

**Interfaces:**
- Consumes: `useUiStore().loadFromConfig` (Task 4).
- Produces: theme applied on mount. (AppShell gains more in Tasks 11–13; this step only adds the ui-store wiring.)

- [x] **Step 1: Add the onMounted hook to `AppShell.vue`**

Replace the `<script setup>` block of `gover2/frontend/src/layouts/AppShell.vue`:

```ts
import { onMounted } from 'vue'
import Sidebar from '@/components/Sidebar.vue'
import Topbar from '@/components/Topbar.vue'
import Statusbar from '@/components/Statusbar.vue'
import { useUiStore } from '@/stores/ui'

const ui = useUiStore()

onMounted(async () => {
  await ui.loadFromConfig()
})
```

- [x] **Step 2: Typecheck + full frontend test run**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm run typecheck && pnpm test`
Expected: PASS.

- [x] **Step 3: Commit**

```bash
cd /home/mg/inventory && git add gover2/frontend/src/layouts/AppShell.vue
git commit -m "feat(gover2): apply persisted theme on app startup"
```

---

## Task 6: `lib/filter.ts` — pure `matchesQuery`

Implements design doc "Search — `lib/filter.ts`" and tasks.md 5.1 + 5.3.

**Files:**
- Create: `gover2/frontend/src/lib/filter.ts`
- Test: `gover2/frontend/src/lib/filter.test.ts`

**Interfaces:**
- Produces: `export function matchesQuery(row: Record<string, unknown>, columns: { key: string }[], query: string): boolean` — empty/whitespace query → `true`; otherwise case-insensitive substring match against any column's displayed value (`row[col.key] ?? ''`).

- [x] **Step 1: Write the failing test**

Create `gover2/frontend/src/lib/filter.test.ts`:

```ts
import { describe, it, expect } from 'vitest'
import { matchesQuery } from './filter'

const columns = [{ key: 'name' }, { key: 'model' }, { key: 'status' }]
const row = { name: 'ThinkPad', model: 'X1 Carbon', status: 'active', notes: 'secret' }

describe('matchesQuery', () => {
  it('returns true for empty query', () => {
    expect(matchesQuery(row, columns, '')).toBe(true)
    expect(matchesQuery(row, columns, '   ')).toBe(true)
  })

  it('matches case-insensitively across visible columns', () => {
    expect(matchesQuery(row, columns, 'thinkpad')).toBe(true)
    expect(matchesQuery(row, columns, 'CARBON')).toBe(true)
    expect(matchesQuery(row, columns, 'active')).toBe(true)
  })

  it('returns false when no visible column matches', () => {
    expect(matchesQuery(row, columns, 'dell')).toBe(false)
  })

  it('ignores columns not in the columns list', () => {
    expect(matchesQuery(row, columns, 'secret')).toBe(false)
  })

  it('treats null/undefined cell values as empty string', () => {
    expect(matchesQuery({ name: null }, [{ key: 'name' }], 'x')).toBe(false)
  })
})
```

- [x] **Step 2: Run the test to verify it fails**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm test src/lib/filter.test.ts`
Expected: FAIL — cannot resolve `./filter`.

- [x] **Step 3: Write the implementation**

Create `gover2/frontend/src/lib/filter.ts`:

```ts
export function matchesQuery(
  row: Record<string, unknown>,
  columns: { key: string }[],
  query: string,
): boolean {
  const q = query.trim().toLowerCase()
  if (q === '') return true
  return columns.some((col) => {
    const value = row[col.key] ?? ''
    return String(value).toLowerCase().includes(q)
  })
}
```

- [x] **Step 4: Run the test to verify it passes**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm test src/lib/filter.test.ts`
Expected: PASS (5 tests).

- [x] **Step 5: Commit**

```bash
cd /home/mg/inventory && git add gover2/frontend/src/lib/filter.ts gover2/frontend/src/lib/filter.test.ts
git commit -m "feat(gover2): add pure matchesQuery filter helper"
```

---

## Task 7: `stores/search.ts` + `lib/routeCategory.ts`

Implements design doc "Search — `stores/search.ts`" and "Topbar — `routeCategory` map".

**Files:**
- Create: `gover2/frontend/src/stores/search.ts`
- Create: `gover2/frontend/src/lib/routeCategory.ts`

**Interfaces:**
- Produces (store `search`): state `query: string`; `setQuery(q: string)`; `clear()`.
- Produces (`routeCategory.ts`):
  - `export const routeCategory: Record<string, string>` mapping `/computers→computers`, `/smartphones→smartphones`, `/tablets→tablets`, `/windows-keys→windowskeys`, `/antivirus→antivirus`, `/other-software→othersoftware`, `/users→users`.
  - `export function categoryForRoute(path: string): string | null` — returns the category or `null` (e.g. `/all`).

- [x] **Step 1: Create `stores/search.ts`**

```ts
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useSearchStore = defineStore('search', () => {
  const query = ref('')
  function setQuery(q: string) {
    query.value = q
  }
  function clear() {
    query.value = ''
  }
  return { query, setQuery, clear }
})
```

- [x] **Step 2: Create `lib/routeCategory.ts`**

```ts
export const routeCategory: Record<string, string> = {
  '/computers': 'computers',
  '/smartphones': 'smartphones',
  '/tablets': 'tablets',
  '/windows-keys': 'windowskeys',
  '/antivirus': 'antivirus',
  '/other-software': 'othersoftware',
  '/users': 'users',
}

export function categoryForRoute(path: string): string | null {
  return routeCategory[path] ?? null
}
```

- [x] **Step 3: Typecheck**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm run typecheck`
Expected: PASS.

- [x] **Step 4: Commit**

```bash
cd /home/mg/inventory && git add gover2/frontend/src/stores/search.ts gover2/frontend/src/lib/routeCategory.ts
git commit -m "feat(gover2): add search store and route-to-category map"
```

---

## Task 8: `DataTable.vue` — density + theme-aware zebra

Implements design doc "Density + Zebra" and tasks.md 2.3 + 2.4 (DataTable half).

**Files:**
- Modify: `gover2/frontend/src/components/DataTable.vue`
- Test: `gover2/frontend/src/components/DataTable.test.ts` (extend)

**Interfaces:**
- Consumes: `useUiStore().density` (Task 4); `row-alt` token (Task 3).
- Produces: cell vertical padding `py-2` (comfortable) / `py-1` (compact); alternating row backgrounds via `even:bg-row-alt`. Existing props/emits unchanged.

- [x] **Step 1: Add failing tests for density + zebra**

Append to `gover2/frontend/src/components/DataTable.test.ts` (it already imports `mount`, `describe`, `it`, `expect`). Add Pinia setup and two tests:

```ts
import { setActivePinia, createPinia } from 'pinia'
import { useUiStore } from '@/stores/ui'

describe('DataTable density + zebra', () => {
  beforeEach(() => setActivePinia(createPinia()))

  it('uses py-2 cells in comfortable density', () => {
    const wrapper = mount(DataTable, { props: { columns, rows, loading: false } })
    const cell = wrapper.find('tbody td')
    expect(cell.classes()).toContain('py-2')
  })

  it('uses py-1 cells in compact density', () => {
    const ui = useUiStore()
    ui.density = 'compact'
    const wrapper = mount(DataTable, { props: { columns, rows, loading: false } })
    const cell = wrapper.find('tbody td')
    expect(cell.classes()).toContain('py-1')
  })

  it('applies an alternating-row background class', () => {
    const wrapper = mount(DataTable, { props: { columns, rows, loading: false } })
    const bodyRows = wrapper.findAll('tbody tr')
    expect(bodyRows[0].classes()).toContain('even:bg-row-alt')
  })
})
```

Add `beforeEach` to the top import: `import { describe, it, expect, beforeEach } from 'vitest'`.

- [x] **Step 2: Run the tests to verify they fail**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm test src/components/DataTable.test.ts`
Expected: FAIL — cells have no `py-1`/`even:bg-row-alt` and the store import is unused by the component yet.

- [x] **Step 3: Update `DataTable.vue`**

Add the ui store and a density-padding computed, apply the zebra class on body rows, and use the padding class on cells. Replace `<script setup>`:

```ts
import { computed } from 'vue'
import { useUiStore } from '@/stores/ui'

defineProps<{
  columns: { key: string; label: string }[]
  rows: Record<string, unknown>[]
  loading: boolean
}>()

defineEmits<{
  'row-click': [row: Record<string, unknown>]
  'delete': [row: Record<string, unknown>]
}>()

const ui = useUiStore()
const cellPad = computed(() => (ui.density === 'compact' ? 'py-1' : 'py-2'))
```

In the template:
- Body data row `<tr>` (the `v-for` at lines 27-33): add `even:bg-row-alt` to its class list, keeping the existing `border-t border-border hover:bg-sidebar/10 cursor-pointer`.
- The data `<td>` (`v-for="col in columns"`): change `class="px-4 py-2 text-text-primary"` to `:class="['px-4', cellPad, 'text-text-primary']"`.
- The delete `<td>`: change `class="px-4 py-2"` to `:class="['px-4', cellPad]"`.

(Header cells and the loading/empty cells keep `py-2`.)

- [x] **Step 4: Run the tests to verify they pass**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm test src/components/DataTable.test.ts`
Expected: PASS (existing 6 + new 3).

- [x] **Step 5: Typecheck**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm run typecheck`
Expected: PASS.

- [x] **Step 6: Commit**

```bash
cd /home/mg/inventory && git add gover2/frontend/src/components/DataTable.vue gover2/frontend/src/components/DataTable.test.ts
git commit -m "feat(gover2): density-aware cell padding and zebra rows in DataTable"
```

---

## Task 9: `AlertsModal.vue` — render real alert rows

Implements design doc "Alerts at App-Shell Level" (modal rendering half) and tasks.md 4.1 + 4.4 (render half).

**Files:**
- Modify: `gover2/frontend/src/components/AlertsModal.vue`
- Test: `gover2/frontend/src/components/AlertsModal.test.ts` (create)

**Interfaces:**
- Consumes: `Alert` type from `@/lib/api`.
- Produces: props `{ open: boolean; alerts: Alert[] }`; emits `close`. Renders a table of `category, name, expiryDate, daysRemaining` with a severity-colored badge (`expired` → `bg-s-expired`, `expiring` → `bg-s-expiring`); empty `alerts` shows a "No alerts" message.

- [x] **Step 1: Write the failing test**

Create `gover2/frontend/src/components/AlertsModal.test.ts`:

```ts
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AlertsModal from './AlertsModal.vue'

const alerts = [
  { category: 'antivirus', id: 1, name: 'Norton', expiryDate: '2026-01-01', daysRemaining: -10, severity: 'expired' },
  { category: 'other_software', id: 2, name: 'Office', expiryDate: '2026-07-01', daysRemaining: 13, severity: 'expiring' },
]

describe('AlertsModal', () => {
  it('renders a row per alert when open', () => {
    const wrapper = mount(AlertsModal, { props: { open: true, alerts } })
    expect(wrapper.text()).toContain('Norton')
    expect(wrapper.text()).toContain('Office')
    expect(wrapper.findAll('tbody tr')).toHaveLength(2)
  })

  it('shows empty message when alerts is empty', () => {
    const wrapper = mount(AlertsModal, { props: { open: true, alerts: [] } })
    expect(wrapper.text()).toContain('No alerts')
  })

  it('renders nothing when closed', () => {
    const wrapper = mount(AlertsModal, { props: { open: false, alerts } })
    expect(wrapper.find('.fixed').exists()).toBe(false)
  })

  it('emits close when Close is clicked', async () => {
    const wrapper = mount(AlertsModal, { props: { open: true, alerts } })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })
})
```

- [x] **Step 2: Run the test to verify it fails**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm test src/components/AlertsModal.test.ts`
Expected: FAIL — modal still renders the placeholder text and accepts no `alerts` prop.

- [x] **Step 3: Rewrite `AlertsModal.vue`**

```vue
<template>
  <div v-if="open" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
    <div class="bg-surface border border-border p-6 w-[560px] max-h-[70vh] overflow-y-auto">
      <h2 class="text-base font-semibold text-text-primary mb-4">Alerts</h2>

      <p v-if="alerts.length === 0" class="text-text-secondary text-sm">No alerts.</p>

      <table v-else class="w-full text-sm">
        <thead class="bg-th-bg">
          <tr>
            <th class="text-left px-3 py-2 text-th-text font-medium">Category</th>
            <th class="text-left px-3 py-2 text-th-text font-medium">Name</th>
            <th class="text-left px-3 py-2 text-th-text font-medium">Expiry</th>
            <th class="text-left px-3 py-2 text-th-text font-medium">Days</th>
            <th class="text-left px-3 py-2 text-th-text font-medium">Severity</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="a in alerts" :key="a.category + '-' + a.id" class="border-t border-border">
            <td class="px-3 py-2 text-text-primary">{{ a.category }}</td>
            <td class="px-3 py-2 text-text-primary">{{ a.name }}</td>
            <td class="px-3 py-2 text-text-primary">{{ a.expiryDate }}</td>
            <td class="px-3 py-2 text-text-primary">{{ a.daysRemaining }}</td>
            <td class="px-3 py-2">
              <span
                class="px-2 py-0.5 rounded-badge text-xs text-white"
                :class="a.severity === 'expired' ? 'bg-s-expired' : 'bg-s-expiring'"
              >
                {{ a.severity }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>

      <button
        class="mt-4 px-3 py-1.5 text-sm bg-s-active text-white hover:bg-accent-hover"
        @click="$emit('close')"
      >
        Close
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Alert } from '@/lib/api'

defineProps<{ open: boolean; alerts: Alert[] }>()
defineEmits<{ (e: 'close'): void }>()
</script>
```

- [x] **Step 4: Run the test to verify it passes**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm test src/components/AlertsModal.test.ts`
Expected: PASS (4 tests).

- [x] **Step 5: Commit**

```bash
cd /home/mg/inventory && git add gover2/frontend/src/components/AlertsModal.vue gover2/frontend/src/components/AlertsModal.test.ts
git commit -m "feat(gover2): render real alert rows with severity badges"
```

---

## Task 10: `AboutModal.vue` + `LicenseModal.vue`

Implements design doc "Topbar — ⋯ overflow: About/License open local modals" and tasks.md 6.2 + 6.3.

**Files:**
- Create: `gover2/frontend/src/components/AboutModal.vue`
- Create: `gover2/frontend/src/components/LicenseModal.vue`
- Test: `gover2/frontend/src/components/AboutModal.test.ts`, `gover2/frontend/src/components/LicenseModal.test.ts`

**Interfaces:**
- Produces (each modal): prop `{ open: boolean }`; emits `close`. Renders static informational content; hidden when `open` is false.

- [ ] **Step 1: Write the failing tests**

Create `gover2/frontend/src/components/AboutModal.test.ts`:

```ts
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AboutModal from './AboutModal.vue'

describe('AboutModal', () => {
  it('renders content when open', () => {
    const wrapper = mount(AboutModal, { props: { open: true } })
    expect(wrapper.text()).toContain('About')
    expect(wrapper.find('.fixed').exists()).toBe(true)
  })

  it('renders nothing when closed', () => {
    const wrapper = mount(AboutModal, { props: { open: false } })
    expect(wrapper.find('.fixed').exists()).toBe(false)
  })

  it('emits close when Close is clicked', async () => {
    const wrapper = mount(AboutModal, { props: { open: true } })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })
})
```

Create `gover2/frontend/src/components/LicenseModal.test.ts` (identical but importing `LicenseModal` and asserting `wrapper.text()` contains `'License'`):

```ts
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import LicenseModal from './LicenseModal.vue'

describe('LicenseModal', () => {
  it('renders content when open', () => {
    const wrapper = mount(LicenseModal, { props: { open: true } })
    expect(wrapper.text()).toContain('License')
    expect(wrapper.find('.fixed').exists()).toBe(true)
  })

  it('renders nothing when closed', () => {
    const wrapper = mount(LicenseModal, { props: { open: false } })
    expect(wrapper.find('.fixed').exists()).toBe(false)
  })

  it('emits close when Close is clicked', async () => {
    const wrapper = mount(LicenseModal, { props: { open: true } })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })
})
```

- [ ] **Step 2: Run the tests to verify they fail**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm test src/components/AboutModal.test.ts src/components/LicenseModal.test.ts`
Expected: FAIL — components do not exist.

- [ ] **Step 3: Create `AboutModal.vue`**

```vue
<template>
  <div v-if="open" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
    <div class="bg-surface border border-border p-6 w-[420px]">
      <h2 class="text-base font-semibold text-text-primary mb-2">About</h2>
      <p class="text-sm text-text-primary">Inventory</p>
      <p class="text-sm text-text-secondary mt-1">A desktop inventory manager built with Wails and Vue.</p>
      <button
        class="mt-4 px-3 py-1.5 text-sm bg-s-active text-white hover:bg-accent-hover"
        @click="$emit('close')"
      >
        Close
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ open: boolean }>()
defineEmits<{ (e: 'close'): void }>()
</script>
```

- [ ] **Step 4: Create `LicenseModal.vue`**

```vue
<template>
  <div v-if="open" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
    <div class="bg-surface border border-border p-6 w-[480px] max-h-[70vh] overflow-y-auto">
      <h2 class="text-base font-semibold text-text-primary mb-2">License</h2>
      <p class="text-sm text-text-secondary whitespace-pre-line">{{ licenseText }}</p>
      <button
        class="mt-4 px-3 py-1.5 text-sm bg-s-active text-white hover:bg-accent-hover"
        @click="$emit('close')"
      >
        Close
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ open: boolean }>()
defineEmits<{ (e: 'close'): void }>()

const licenseText = 'This software is distributed under the project license. See the LICENSE file in the repository for the full terms.'
</script>
```

- [ ] **Step 5: Run the tests to verify they pass**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm test src/components/AboutModal.test.ts src/components/LicenseModal.test.ts`
Expected: PASS (3 each).

- [ ] **Step 6: Commit**

```bash
cd /home/mg/inventory && git add gover2/frontend/src/components/AboutModal.vue gover2/frontend/src/components/LicenseModal.vue gover2/frontend/src/components/AboutModal.test.ts gover2/frontend/src/components/LicenseModal.test.ts
git commit -m "feat(gover2): add About and License modals"
```

---

## Task 11: api wrapper — `exportCSV(category)` signature change

Implements design doc "CSV Export — api" and tasks.md 3.3 (api half).

**Files:**
- Modify: `gover2/frontend/src/lib/api/index.ts:103-104`

**Interfaces:**
- Produces: `export const exportCSV = (category: string): Promise<void>` — calls `call<void>('ExportCSV', category)` (drops `destPath`).

- [ ] **Step 1: Update the wrapper**

Replace lines 103-104 of `gover2/frontend/src/lib/api/index.ts`:

```ts
export const exportCSV = (category: string): Promise<void> =>
  call<void>('ExportCSV', category)
```

- [ ] **Step 2: Typecheck**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm run typecheck`
Expected: PASS (no remaining 2-arg callers — the stub had none).

- [ ] **Step 3: Commit**

```bash
cd /home/mg/inventory && git add gover2/frontend/src/lib/api/index.ts
git commit -m "feat(gover2): exportCSV api wrapper drops destPath"
```

---

## Task 12: `Topbar.vue` — action row

Implements design doc "Topbar Action Row" and tasks.md 3.1 + 3.2 + 3.3 + 3.4 (UI wiring; modal open events bubble to AppShell in Task 13).

**Files:**
- Modify: `gover2/frontend/src/components/Topbar.vue`

**Interfaces:**
- Consumes: `useSearchStore` (Task 7), `useUiStore` (Task 4), `categoryForRoute` (Task 7), `exportCSV` (Task 11), `newDatabase`/`openDatabase` from `@/lib/api`.
- Produces: emits `open-alerts`, `open-about`, `open-license` (handled by AppShell in Task 13). Search input `v-model`s `search.query`. ↓ CSV hidden when `categoryForRoute(route.path)` is `null` (i.e. `/all`). 🌙 calls `ui.toggleDarkMode()`. ⋯ overflow toggles a local `menuOpen` ref.

- [ ] **Step 1: Rewrite `Topbar.vue`**

```vue
<template>
  <header class="flex items-center gap-3 px-4 bg-surface border-b border-border" style="height: 48px">
    <span class="text-sm text-text-secondary">{{ title }}</span>

    <input
      v-model="search.query"
      type="search"
      placeholder="Search..."
      class="ml-2 px-2 py-1 text-sm bg-page text-text-primary border border-border rounded w-56"
    />

    <div class="ml-auto flex items-center gap-2">
      <button
        class="px-2 py-1 text-sm text-text-primary hover:bg-page rounded"
        title="Alerts"
        @click="$emit('open-alerts')"
      >
        ⚠ Alerts
      </button>

      <button
        v-if="exportCategory"
        class="px-2 py-1 text-sm text-text-primary hover:bg-page rounded"
        title="Export CSV"
        @click="onExport"
      >
        ↓ CSV
      </button>

      <button
        class="px-2 py-1 text-sm text-text-primary hover:bg-page rounded"
        :title="ui.darkMode ? 'Switch to light' : 'Switch to dark'"
        @click="ui.toggleDarkMode()"
      >
        {{ ui.darkMode ? '☀' : '🌙' }}
      </button>

      <div class="relative">
        <button
          class="px-2 py-1 text-sm text-text-primary hover:bg-page rounded"
          title="More"
          @click="menuOpen = !menuOpen"
        >
          ⋯
        </button>
        <div
          v-if="menuOpen"
          class="absolute right-0 mt-1 w-44 bg-surface border border-border shadow z-50 text-sm"
          @click="menuOpen = false"
        >
          <button class="block w-full text-left px-3 py-2 text-text-primary hover:bg-page" @click="onNewDatabase">New Database</button>
          <button class="block w-full text-left px-3 py-2 text-text-primary hover:bg-page" @click="onOpenDatabase">Open Database</button>
          <div class="border-t border-border"></div>
          <button class="block w-full text-left px-3 py-2 text-text-primary hover:bg-page" @click="$emit('open-about')">About</button>
          <button class="block w-full text-left px-3 py-2 text-text-primary hover:bg-page" @click="$emit('open-license')">License</button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useSearchStore } from '@/stores/search'
import { useUiStore } from '@/stores/ui'
import { categoryForRoute } from '@/lib/routeCategory'
import { exportCSV, newDatabase, openDatabase } from '@/lib/api'

const emit = defineEmits<{
  (e: 'open-alerts'): void
  (e: 'open-about'): void
  (e: 'open-license'): void
  (e: 'db-changed'): void
}>()

const route = useRoute()
const search = useSearchStore()
const ui = useUiStore()
const menuOpen = ref(false)

const title = computed(() =>
  String(route.path).replace('/', '').replace(/-/g, ' ') || 'Inventory',
)
const exportCategory = computed(() => categoryForRoute(route.path))

async function onExport() {
  if (exportCategory.value) {
    await exportCSV(exportCategory.value)
  }
}

async function onNewDatabase() {
  await newDatabase('')
  emit('db-changed')
}

async function onOpenDatabase() {
  await openDatabase('')
  emit('db-changed')
}
</script>
```

Note: `newDatabase`/`openDatabase` currently take a path string; the design defers native-dialog path selection to the Go side. Passing `''` is a placeholder until Task 13 finalizes DB-management semantics — see Task 13 Step 1 for the bridge-side decision. If the bridge still requires a path argument at this point, leave the call as `newDatabase('')` and the smoke test (Task 14) will confirm behavior; do not invent a frontend file dialog.

- [ ] **Step 2: Typecheck**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm run typecheck`
Expected: PASS.

- [ ] **Step 3: Commit**

```bash
cd /home/mg/inventory && git add gover2/frontend/src/components/Topbar.vue
git commit -m "feat(gover2): topbar action row with search, alerts, CSV, dark, overflow"
```

---

## Task 13: `AppShell.vue` — own modals, auto-open alerts, clear search on navigation

Implements design doc "Alerts at App-Shell Level" + "Topbar overflow About/License" and tasks.md 4.2 + 4.3 + 6.1 + 6.2 (ownership) + 5.2 (route-change clear half).

**Files:**
- Modify: `gover2/frontend/src/layouts/AppShell.vue`

**Interfaces:**
- Consumes: `getAlerts` from `@/lib/api`; `AlertsModal`/`AboutModal`/`LicenseModal` (Tasks 9–10); `useUiStore` (already wired in Task 5); `useSearchStore` (Task 7); Topbar emits `open-alerts`/`open-about`/`open-license`/`db-changed` (Task 12).
- Produces: AppShell holds `alertsOpen`/`aboutOpen`/`licenseOpen` refs and `alerts` list; auto-opens alerts on startup when `getAlerts()` non-empty; ⚠ action re-fetches + opens; clears search on route change.

- [ ] **Step 1: Rewrite `AppShell.vue`**

```vue
<template>
  <div class="flex h-screen bg-page font-ui">
    <Sidebar />
    <div class="flex flex-col flex-1 overflow-hidden">
      <Topbar
        @open-alerts="onOpenAlerts"
        @open-about="aboutOpen = true"
        @open-license="licenseOpen = true"
        @db-changed="onDbChanged"
      />
      <main class="flex-1 overflow-auto p-4">
        <slot />
      </main>
      <Statusbar />
    </div>

    <AlertsModal :open="alertsOpen" :alerts="alerts" @close="alertsOpen = false" />
    <AboutModal :open="aboutOpen" @close="aboutOpen = false" />
    <LicenseModal :open="licenseOpen" @close="licenseOpen = false" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import Topbar from '@/components/Topbar.vue'
import Statusbar from '@/components/Statusbar.vue'
import AlertsModal from '@/components/AlertsModal.vue'
import AboutModal from '@/components/AboutModal.vue'
import LicenseModal from '@/components/LicenseModal.vue'
import { useUiStore } from '@/stores/ui'
import { useSearchStore } from '@/stores/search'
import { getAlerts } from '@/lib/api'
import type { Alert } from '@/lib/api'

const ui = useUiStore()
const search = useSearchStore()
const route = useRoute()

const alerts = ref<Alert[]>([])
const alertsOpen = ref(false)
const aboutOpen = ref(false)
const licenseOpen = ref(false)

onMounted(async () => {
  await ui.loadFromConfig()
  alerts.value = await getAlerts()
  if (alerts.value.length > 0) alertsOpen.value = true
})

async function onOpenAlerts() {
  alerts.value = await getAlerts()
  alertsOpen.value = true
}

function onDbChanged() {
  // views refetch on mount; nothing else needed at shell level for now
}

watch(() => route.path, () => search.clear())
</script>
```

- [ ] **Step 2: Typecheck + full frontend test run**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm run typecheck && pnpm test`
Expected: PASS.

- [ ] **Step 3: Commit**

```bash
cd /home/mg/inventory && git add gover2/frontend/src/layouts/AppShell.vue
git commit -m "feat(gover2): AppShell owns alerts/about/license modals, auto-opens alerts, clears search on nav"
```

---

## Task 14: Wire search filtering into all 8 views

Implements design doc "Search — each view computes `filteredRows`" and tasks.md 5.2.

**Files:**
- Modify: `gover2/frontend/src/features/assets/ComputersView.vue`
- Modify: `gover2/frontend/src/features/assets/SmartphonesView.vue`
- Modify: `gover2/frontend/src/features/assets/TabletsView.vue`
- Modify: `gover2/frontend/src/features/assets/AllAssetsView.vue`
- Modify: `gover2/frontend/src/features/licenses/WindowsKeysView.vue`
- Modify: `gover2/frontend/src/features/licenses/AntivirusView.vue`
- Modify: `gover2/frontend/src/features/licenses/OtherSoftwareView.vue`
- Modify: `gover2/frontend/src/features/users/UsersView.vue`

**Interfaces:**
- Consumes: `useSearchStore().query` (Task 7), `matchesQuery` (Task 6).
- Produces: each view passes a `filteredRows` computed to `DataTable` instead of the raw `rows`.

- [ ] **Step 1: Apply the same pattern to each view**

For each view, in `<script setup>` add:

```ts
import { computed } from 'vue'
import { useSearchStore } from '@/stores/search'
import { matchesQuery } from '@/lib/filter'

const search = useSearchStore()
const filteredRows = computed(() =>
  (rows.value as Record<string, unknown>[]).filter((r) => matchesQuery(r, columns, search.query)),
)
```

(If the view already imports `computed` from `vue` — e.g. via an existing import line — merge rather than duplicate.)

Then change the `DataTable` binding from `:rows="rows ..."` to `:rows="filteredRows"`.

Concretely, in `ComputersView.vue` replace:

```vue
    <DataTable
      :columns="columns"
      :rows="rows as Record<string, unknown>[]"
      :loading="loading"
```

with:

```vue
    <DataTable
      :columns="columns"
      :rows="filteredRows"
      :loading="loading"
```

`AllAssetsView.vue` already declares `rows` as `Record<string, unknown>[]`; its filter is:

```ts
const filteredRows = computed(() =>
  rows.value.filter((r) => matchesQuery(r, columns, search.query)),
)
```

and change `<DataTable :columns="columns" :rows="rows" ... />` to `:rows="filteredRows"`.

Apply the structurally identical change to `SmartphonesView.vue`, `TabletsView.vue`, `WindowsKeysView.vue`, `AntivirusView.vue`, `OtherSoftwareView.vue`, and `UsersView.vue`, using each view's own `rows` ref and `columns` array.

- [ ] **Step 2: Typecheck**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm run typecheck`
Expected: PASS.

- [ ] **Step 3: Full frontend test run**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm test`
Expected: PASS.

- [ ] **Step 4: Commit**

```bash
cd /home/mg/inventory && git add gover2/frontend/src/features
git commit -m "feat(gover2): client-side search filtering across all views"
```

---

## Task 15: Regenerate wailsjs bindings + full verification

Implements design doc "Testing Strategy" + "api: wailsjs bindings regenerated by `wails build`" and tasks.md 7.1 + 7.2 + 7.3.

**Files:**
- Generated: `gover2/frontend/wailsjs/**` (machine-generated; do not hand-edit).

**Interfaces:**
- Consumes: everything above. Produces a clean build with the new `ExportCSV(category)` binding.

- [ ] **Step 1: Regenerate bindings via wails build**

Run: `cd /home/mg/inventory/gover2 && wails build`
Expected: build succeeds; `frontend/wailsjs/go/bridge/App.d.ts` now shows `ExportCSV(arg1: string): Promise<void>` (single arg).

If `wails` is not installed in the environment, run the binding-generation step the project uses instead (`wails generate module`) and note it; do not hand-edit `wailsjs/`.

- [ ] **Step 2: Go build + full Go test with coverage**

Run: `cd /home/mg/inventory/gover2 && go build ./... && go test ./... -cover`
Expected: PASS; services ≥ 70%, bridge ≥ 50%.

- [ ] **Step 3: Frontend typecheck + full test run**

Run: `cd /home/mg/inventory/gover2/frontend && pnpm run typecheck && pnpm test`
Expected: both PASS.

- [ ] **Step 4: Manual smoke test (record results)**

Run: `cd /home/mg/inventory/gover2 && wails dev`
Then verify each item:
- Startup alerts modal auto-opens when expiry alerts exist (seed an antivirus row with a past `expiry_date` first).
- ↓ CSV on a category view opens a save dialog and writes a correct file; hidden on `/all`.
- 🌙 toggles dark mode; restart the app and confirm the theme persisted.
- Search box filters the visible table live and resets when navigating to another view.
- ⋯ → New Database / Open Database switch the DB and views refetch; About / License modals open and close.

- [ ] **Step 5: Commit generated bindings**

```bash
cd /home/mg/inventory && git add gover2/frontend/wailsjs
git commit -m "chore(gover2): regenerate wailsjs bindings for ExportCSV signature"
```

---

## Self-Review

**Spec coverage (design doc + tasks.md):**
- Go CSV export (1.1–1.3): Tasks 1–2. ✓
- Theme/preferences (2.1–2.4): Tasks 3 (tokens), 4 (store), 5 (startup), 8 (DataTable). ✓
- Topbar (3.1–3.4): Tasks 11 (api), 12 (action row). ✓
- Alerts modal (4.1–4.4): Tasks 9 (render), 13 (ownership/auto-open/reopen). ✓
- Search (5.1–5.3): Tasks 6 (helper), 7 (store), 13 (clear on nav), 14 (views). ✓
- DB management + dialogs (6.1–6.3): Tasks 10 (modals), 12 (overflow wiring), 13 (db-changed). ✓
- Integration & verification (7.1–7.3): Task 15. ✓

**Placeholder scan:** Each code step contains full code. The one open decision — whether `NewDatabase`/`OpenDatabase` keep a path argument or move the dialog Go-side — is flagged explicitly in Task 12 Step 1 with a concrete fallback (`newDatabase('')`) rather than a hand-waved "handle dialog"; resolving it is part of the smoke test in Task 15.

**Type consistency:** `matchesQuery(row, columns, query)` signature is identical in Tasks 6 and 14. `BuildCSV`/`WriteCSV` names match between Tasks 1 and 2. `exportCSV(category)` matches between Tasks 11 and 12. ui-store members (`loadFromConfig`, `toggleDarkMode`, `setDensity`, `applyTheme`, `density`, `darkMode`) are consistent across Tasks 4, 5, 8, 12, 13. `AlertsModal` props `{ open, alerts }` match between Tasks 9 and 13. Category strings (`windowskeys`, `othersoftware`) match between Task 1 (`BuildCSV` switch) and Task 7 (`routeCategory` map).

## Execution Handoff

Plan complete. Two execution options:

1. **Subagent-Driven (recommended)** — dispatch a fresh subagent per task, review between tasks (REQUIRED SUB-SKILL: superpowers:subagent-driven-development).
2. **Inline Execution** — execute tasks in this session with checkpoints (REQUIRED SUB-SKILL: superpowers:executing-plans).
