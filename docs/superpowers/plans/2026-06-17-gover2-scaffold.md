---
change: gover2-scaffold
design-doc: docs/superpowers/specs/2026-06-17-gover2-scaffold-design.md
base-ref: 39fb9f362c810c56a3214f33fce474c0fb739db5
archived-with: 2026-06-17-gover2-scaffold
---

# Gover2 Scaffold Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Scaffold the `gover2/` Wails v2 + Vue 3 project with all 6 Go package stubs, complete TypeScript frontend structure, Tailwind v3 design tokens, and `gover:` Taskfile commands — zero business logic, all stubs return empty/nil.

**Architecture:** A Wails v2 desktop app living at `gover2/` inside the existing `inventory/` monorepo. Go backend is split into 6 layered `internal/` packages (`models` → `database`/`config`/`backup` → `services` → `bridge`). Vue 3 frontend uses feature-folder layout with Pinia + Vue Router hash mode; all bridge calls are routed through `lib/api/index.ts`, never `window.go.*` directly.

**Tech Stack:** Go 1.22+, Wails v2, modernc.org/sqlite (pure-Go, no CGo), Vue 3, Vite, TypeScript (strict), Pinia, Vue Router 4, Tailwind CSS v3 (pinned — v4 is incompatible), pnpm, vue-tsc.

## Global Constraints

- All commands run from `gover2/` unless otherwise stated.
- Go module name: `module gover2`, `go 1.22`.
- pnpm is mandatory — npm must never be invoked in the frontend.
- `wails.json` must set pnpm overrides BEFORE any `wails build` or `wails dev` call.
- Tailwind v3 only — `tailwindcss@3` pinned explicitly; `@tailwindcss/vite` is v4-only and must NOT be used.
- `lib/api/index.ts` returns `Promise.resolve([])` stubs — no `wailsjs/` imports in scaffold phase.
- No `window.go.*` calls anywhere in `frontend/src/features/` or `frontend/src/components/`.
- No `any` types in component or API files — `tsconfig.json` must have `"strict": true`.
- `internal/bridge` is the only package that imports from all others; no other package may import `bridge`.
- `internal/database` must keep `_ "modernc.org/sqlite"` blank import even in stub form to prevent `go mod tidy` from removing the dep.
- No unit tests in this change (added in `gover2-readonly`).
- No SQLite connections, no real data, no packaging scripts (Phase 5).

archived-with: 2026-06-17-gover2-scaffold
---

## File Map

Files created or modified by this change, grouped by task:

**Task 0 — Toolchain**
- (system-level installs, no repo files)

**Task 1 — Wails init + go.mod**
- Create: `gover2/go.mod`
- Create: `gover2/go.sum`
- Create: `gover2/main.go`
- Create: `gover2/app.go` (generated, then replaced in Task 2.7)
- Create: `gover2/wails.json`
- Create: `gover2/frontend/` (generated scaffolding, replaced in Task 3)

**Task 2 — Go package stubs**
- Create: `gover2/internal/models/models.go`
- Create: `gover2/internal/database/database.go`
- Create: `gover2/internal/config/config.go`
- Create: `gover2/internal/backup/backup.go`
- Create: `gover2/internal/services/services.go`
- Create: `gover2/internal/bridge/bridge.go`
- Modify: `gover2/app.go`

**Task 3 — Frontend package + TypeScript types + API layer**
- Modify: `gover2/frontend/package.json`
- Create: `gover2/frontend/vite.config.ts`
- Create: `gover2/frontend/tsconfig.json`
- Create: `gover2/frontend/index.html`
- Create: `gover2/frontend/src/lib/types/index.ts`
- Create: `gover2/frontend/src/lib/api/index.ts`

**Task 4 — Pinia stores**
- Create: `gover2/frontend/src/stores/assets.ts`
- Create: `gover2/frontend/src/stores/ui.ts`

**Task 5 — Vue Router + App shell**
- Create: `gover2/frontend/src/router/index.ts`
- Create: `gover2/frontend/src/App.vue`
- Create: `gover2/frontend/src/main.ts`
- Create: `gover2/frontend/src/layouts/AppShell.vue`

**Task 6 — Feature view stubs**
- Create: `gover2/frontend/src/features/assets/ComputersView.vue`
- Create: `gover2/frontend/src/features/assets/SmartphonesView.vue`
- Create: `gover2/frontend/src/features/assets/TabletsView.vue`
- Create: `gover2/frontend/src/features/assets/AllAssetsView.vue`
- Create: `gover2/frontend/src/features/licenses/WindowsKeysView.vue`
- Create: `gover2/frontend/src/features/licenses/AntivirusView.vue`
- Create: `gover2/frontend/src/features/licenses/OtherSoftwareView.vue`
- Create: `gover2/frontend/src/features/users/UsersView.vue`
- Create: `gover2/frontend/src/features/alerts/AlertsView.vue`

**Task 7 — Shared component stubs**
- Create: `gover2/frontend/src/components/Sidebar.vue`
- Create: `gover2/frontend/src/components/SidebarItem.vue`
- Create: `gover2/frontend/src/components/Topbar.vue`
- Create: `gover2/frontend/src/components/DataTable.vue`
- Create: `gover2/frontend/src/components/StatusBadge.vue`
- Create: `gover2/frontend/src/components/EditPanel.vue`
- Create: `gover2/frontend/src/components/FormField.vue`
- Create: `gover2/frontend/src/components/AlertsModal.vue`
- Create: `gover2/frontend/src/components/Statusbar.vue`

**Task 8 — Tailwind CSS v3**
- Create: `gover2/frontend/tailwind.config.ts`
- Create: `gover2/frontend/postcss.config.js`
- Create: `gover2/frontend/src/style.css`

**Task 9 — Taskfile integration**
- Modify: `Taskfile.yml` (repo root)

**Task 10 — Verification**
- (no new files; smoke-test sequence)

archived-with: 2026-06-17-gover2-scaffold
---

### Task 0: Toolchain Prerequisites

**Files:** none (system-level installs)

**Interfaces:**
- Produces: `wails` binary on PATH, `pnpm` binary on PATH

- [x] **Step 1: Install Wails v2**

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Expected: installs without error. The binary lands in `$(go env GOPATH)/bin/`.

- [x] **Step 2: Verify wails is on PATH**

```bash
wails version
```

Expected output (version may differ):
```
Wails CLI v2.9.x
```

If `command not found`, add `$(go env GOPATH)/bin` to your `PATH` and reload the shell.

- [x] **Step 3: Enable pnpm via corepack**

```bash
corepack enable pnpm
```

Expected: no error. Corepack is bundled with Node 16.9+. If Node is not installed, install it first (`nvm install --lts` or system package manager), then re-run.

- [x] **Step 4: Verify pnpm is available**

```bash
pnpm --version
```

Expected: prints a version string like `9.x.x`. If this fails after `corepack enable`, run `corepack prepare pnpm@latest --activate`.

archived-with: 2026-06-17-gover2-scaffold
---

### Task 1: Wails Project Init + pnpm Wiring

**Files:**
- Create: `gover2/go.mod`, `gover2/go.sum`, `gover2/main.go`, `gover2/app.go`, `gover2/wails.json`
- Create: `gover2/frontend/` (generated stub — replaced in Task 3)

**Interfaces:**
- Produces: compilable `gover2/` Go module with `module gover2`; `wails.json` configured for pnpm

- [x] **Step 1: Run wails init from inventory root**

```bash
cd /home/mg/inventory
wails init -n gover2 -t vue -d gover2 -g
```

Flags: `-n gover2` (app name), `-t vue` (Vue template), `-d gover2` (output directory), `-g` (skip git init — we're already in a git repo).

Expected: directory `gover2/` is created with `go.mod`, `main.go`, `app.go`, `wails.json`, and `frontend/`.

- [x] **Step 2: Verify module name**

```bash
head -1 /home/mg/inventory/gover2/go.mod
```

Expected:
```
module gover2
```

If it says something else (e.g. `module github.com/...`), edit `go.mod` so line 1 is exactly `module gover2`.

- [x] **Step 3: Verify generated files exist**

```bash
ls /home/mg/inventory/gover2/
```

Expected to include: `go.mod`, `main.go`, `app.go`, `wails.json`, `frontend/`

- [x] **Step 4: Add modernc.org/sqlite dependency**

```bash
cd /home/mg/inventory/gover2
go get modernc.org/sqlite
```

Expected: `go.mod` and `go.sum` updated; no error output.

- [x] **Step 5: Patch wails.json for pnpm**

Open `gover2/wails.json`. The generated file uses npm. Replace it entirely with:

```json
{
  "name": "gover2",
  "outputfilename": "gover2",
  "frontend:install": "pnpm install",
  "frontend:build": "pnpm run build",
  "frontend:dev:watcher": "pnpm run dev",
  "frontend:dev:serverUrl": "auto",
  "frontend:dir": "frontend",
  "info": {
    "productName": "Gover2",
    "productVersion": "0.1.0",
    "companyName": "",
    "copyright": "",
    "comments": ""
  },
  "bindings": {
    "ts_generation_dir": "frontend/wailsjs"
  }
}
```

Save the file.

- [x] **Step 6: Verify Go compiles with no errors**

```bash
cd /home/mg/inventory/gover2
go build ./...
```

Expected: exits 0, no output. (The generated `app.go` stubs compile cleanly; we replace them in Task 2.)

- [x] **Step 7: Commit**

```bash
cd /home/mg/inventory
git add gover2/go.mod gover2/go.sum gover2/main.go gover2/app.go gover2/wails.json gover2/frontend/
git commit -m "feat(gover2): wails init with pnpm wiring and sqlite dep"
```

archived-with: 2026-06-17-gover2-scaffold
---

### Task 2: Go Package Stubs

**Files:**
- Create: `gover2/internal/models/models.go`
- Create: `gover2/internal/database/database.go`
- Create: `gover2/internal/config/config.go`
- Create: `gover2/internal/backup/backup.go`
- Create: `gover2/internal/services/services.go`
- Create: `gover2/internal/bridge/bridge.go`
- Modify: `gover2/app.go`

**Interfaces:**
- Produces:
  - `models.Computer`, `models.Smartphone`, `models.Tablet`, `models.WindowsKey`, `models.Antivirus`, `models.OtherSoftware`, `models.User`, `models.AppConfig`, `models.Alert` structs
  - `database.Open(path string) (*sql.DB, error)`
  - `config.Load() (models.AppConfig, error)`, `config.Save(cfg models.AppConfig) error`
  - `backup.BackupDB(srcPath, destDir string) error`
  - `services.ListComputers(db *sql.DB) ([]models.Computer, error)` and 6 equivalent List* functions
  - `bridge.App` struct with `NewApp() *bridge.App`; methods: `ListComputers()`, `ListSmartphones()`, `ListTablets()`, `ListWindowsKeys()`, `ListAntivirus()`, `ListOtherSoftware()`, `ListUsers()`, `GetDatabasePath()`, `GetConfig()`, `SetConfig()`, `NewDatabase()`, `OpenDatabase()`, `ExportCSV()`
  - `app.go` `NewApp() *bridge.App` function

- [x] **Step 1: Create models package**

Create file `gover2/internal/models/models.go`:

```go
package models

// Computer maps the computers table.
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

// Smartphone maps the smartphones table (identical shape to Computer).
type Smartphone struct {
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

// Tablet maps the tablets table (identical shape to Computer).
type Tablet struct {
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

// WindowsKey maps the windows_keys table.
type WindowsKey struct {
	ID         int     `json:"id"`
	LicenseKey string  `json:"licenseKey"`
	ComputerID *int    `json:"computerId"`
	Status     string  `json:"status"`
	Notes      *string `json:"notes"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

// Antivirus maps the antivirus table.
type Antivirus struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	LicenseKey   string  `json:"licenseKey"`
	ComputerID   *int    `json:"computerId"`
	SmartphoneID *int    `json:"smartphoneId"`
	TabletID     *int    `json:"tabletId"`
	Status       string  `json:"status"`
	ExpiryDate   *string `json:"expiryDate"`
	Notes        *string `json:"notes"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

// OtherSoftware maps the other_software table (identical shape to Antivirus).
type OtherSoftware struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	LicenseKey   string  `json:"licenseKey"`
	ComputerID   *int    `json:"computerId"`
	SmartphoneID *int    `json:"smartphoneId"`
	TabletID     *int    `json:"tabletId"`
	Status       string  `json:"status"`
	ExpiryDate   *string `json:"expiryDate"`
	Notes        *string `json:"notes"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

// User maps the users table.
type User struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Surname   *string `json:"surname"`
	Status    string  `json:"status"`
	Notes     *string `json:"notes"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

// AppConfig holds application-level settings persisted to JSON.
type AppConfig struct {
	DBPath            string `json:"dbPath"`
	Density           string `json:"density"`
	DarkMode          bool   `json:"darkMode"`
	ExpiryWarningDays int    `json:"expiryWarningDays"`
}

// Alert represents an item approaching or past its expiry date.
type Alert struct {
	Category     string `json:"category"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ExpiryDate   string `json:"expiryDate"`
	DaysRemaining int   `json:"daysRemaining"`
	Severity     string `json:"severity"`
}
```

- [x] **Step 2: Create database package stub**

Create file `gover2/internal/database/database.go`:

```go
package database

import (
	"database/sql"

	_ "modernc.org/sqlite" // register sqlite3 driver
)

// Open returns a stub nil DB. Real implementation in gover2-readonly.
func Open(path string) (*sql.DB, error) {
	return nil, nil
}
```

Note: the blank import `_ "modernc.org/sqlite"` is required even in stub form to keep the dep in `go.sum` so `go mod tidy` does not remove it.

- [x] **Step 3: Create config package stub**

Create file `gover2/internal/config/config.go`:

```go
package config

import "gover2/internal/models"

// Load returns a stub zero-value AppConfig. Real implementation in gover2-readonly.
func Load() (models.AppConfig, error) {
	return models.AppConfig{
		Density:           "comfortable",
		DarkMode:          false,
		ExpiryWarningDays: 30,
	}, nil
}

// Save is a no-op stub. Real implementation in gover2-readonly.
func Save(cfg models.AppConfig) error {
	return nil
}
```

- [x] **Step 4: Create backup package stub**

Create file `gover2/internal/backup/backup.go`:

```go
package backup

// BackupDB is a no-op stub. Real implementation in gover2-readonly.
func BackupDB(srcPath, destDir string) error {
	return nil
}
```

- [x] **Step 5: Create services package stub**

Create file `gover2/internal/services/services.go`:

```go
package services

import (
	"database/sql"

	"gover2/internal/models"
)

// ListComputers returns an empty slice. Real implementation in gover2-readonly.
func ListComputers(db *sql.DB) ([]models.Computer, error) {
	return []models.Computer{}, nil
}

// ListSmartphones returns an empty slice. Real implementation in gover2-readonly.
func ListSmartphones(db *sql.DB) ([]models.Smartphone, error) {
	return []models.Smartphone{}, nil
}

// ListTablets returns an empty slice. Real implementation in gover2-readonly.
func ListTablets(db *sql.DB) ([]models.Tablet, error) {
	return []models.Tablet{}, nil
}

// ListWindowsKeys returns an empty slice. Real implementation in gover2-readonly.
func ListWindowsKeys(db *sql.DB) ([]models.WindowsKey, error) {
	return []models.WindowsKey{}, nil
}

// ListAntivirus returns an empty slice. Real implementation in gover2-readonly.
func ListAntivirus(db *sql.DB) ([]models.Antivirus, error) {
	return []models.Antivirus{}, nil
}

// ListOtherSoftware returns an empty slice. Real implementation in gover2-readonly.
func ListOtherSoftware(db *sql.DB) ([]models.OtherSoftware, error) {
	return []models.OtherSoftware{}, nil
}

// ListUsers returns an empty slice. Real implementation in gover2-readonly.
func ListUsers(db *sql.DB) ([]models.User, error) {
	return []models.User{}, nil
}
```

- [x] **Step 6: Create bridge package stub**

Create file `gover2/internal/bridge/bridge.go`:

```go
package bridge

import (
	"gover2/internal/backup"
	"gover2/internal/config"
	"gover2/internal/models"
	"gover2/internal/services"
)

// App is the Wails context struct. All exported methods are callable from TypeScript.
// In scaffold phase, db is nil — all list methods return empty slices.
type App struct{}

// NewApp creates a new App instance for Wails binding.
func NewApp() *App {
	return &App{}
}

// -- Read-only methods (Phase 2) --

func (a *App) ListComputers() ([]models.Computer, error) {
	return services.ListComputers(nil)
}

func (a *App) ListSmartphones() ([]models.Smartphone, error) {
	return services.ListSmartphones(nil)
}

func (a *App) ListTablets() ([]models.Tablet, error) {
	return services.ListTablets(nil)
}

func (a *App) ListWindowsKeys() ([]models.WindowsKey, error) {
	return services.ListWindowsKeys(nil)
}

func (a *App) ListAntivirus() ([]models.Antivirus, error) {
	return services.ListAntivirus(nil)
}

func (a *App) ListOtherSoftware() ([]models.OtherSoftware, error) {
	return services.ListOtherSoftware(nil)
}

func (a *App) ListUsers() ([]models.User, error) {
	return services.ListUsers(nil)
}

func (a *App) GetDatabasePath() (string, error) {
	return "", nil
}

func (a *App) GetConfig() (models.AppConfig, error) {
	return config.Load()
}

func (a *App) SetConfig(cfg models.AppConfig) error {
	return config.Save(cfg)
}

// -- File dialog methods (Phase 2) --

func (a *App) NewDatabase(path string) error {
	return nil
}

func (a *App) OpenDatabase(path string) error {
	return nil
}

func (a *App) ExportCSV(category string, destPath string) error {
	return nil
}

// keep backup imported to satisfy dependency graph (used in gover2-readonly)
var _ = backup.BackupDB
```

- [x] **Step 7: Replace app.go to wire bridge**

Replace the contents of `gover2/app.go` entirely:

```go
package main

import "gover2/internal/bridge"

// NewApp returns the Wails context object bound to the frontend.
func NewApp() *bridge.App {
	return bridge.NewApp()
}
```

- [x] **Step 8: Verify no import cycles**

```bash
cd /home/mg/inventory/gover2
go build ./...
```

Expected: exits 0, no output. Any import cycle error here means a package in `internal/` (other than `bridge`) is importing `bridge` — check the failing package and remove the cycle.

- [x] **Step 9: Commit**

```bash
cd /home/mg/inventory
git add gover2/internal/ gover2/app.go
git commit -m "feat(gover2): add 6 Go package stubs with dependency-ordered layout"
```

archived-with: 2026-06-17-gover2-scaffold
---

### Task 3: Frontend Package Setup + TypeScript Types + API Layer

**Files:**
- Modify: `gover2/frontend/package.json`
- Create: `gover2/frontend/vite.config.ts`
- Create: `gover2/frontend/tsconfig.json`
- Create: `gover2/frontend/index.html`
- Create: `gover2/frontend/src/lib/types/index.ts`
- Create: `gover2/frontend/src/lib/api/index.ts`

**Interfaces:**
- Consumes: nothing from earlier tasks (pure TypeScript)
- Produces:
  - All 9 domain interfaces + `AppConfig` + `Alert` exported from `frontend/src/lib/types/index.ts`
  - All bridge wrapper functions exported from `frontend/src/lib/api/index.ts`, all returning `Promise.resolve([])`/`Promise.resolve('')`/`Promise.resolve()`
  - Runnable `pnpm run typecheck` command

- [x] **Step 1: Replace package.json with pnpm + Vue 3 stack**

Replace `gover2/frontend/package.json` entirely:

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
    "clean": "rm -rf dist"
  },
  "dependencies": {
    "pinia": "^2.1.7",
    "vue": "^3.4.0",
    "vue-router": "^4.3.0"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.0.0",
    "autoprefixer": "^10.4.0",
    "postcss": "^8.4.0",
    "tailwindcss": "3",
    "typescript": "^5.4.0",
    "vite": "^5.0.0",
    "vue-tsc": "^2.0.0"
  }
}
```

- [x] **Step 2: Create vite.config.ts**

Create `gover2/frontend/vite.config.ts`:

```ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': '/src',
    },
  },
})
```

Note: do NOT use `@tailwindcss/vite` — that is the v4 plugin and incompatible with the v3 config. PostCSS handles Tailwind separately.

- [x] **Step 3: Create tsconfig.json with strict mode**

Create `gover2/frontend/tsconfig.json`:

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "module": "ESNext",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "preserve",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src/**/*.ts", "src/**/*.tsx", "src/**/*.vue"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

Create `gover2/frontend/tsconfig.node.json`:

```json
{
  "compilerOptions": {
    "composite": true,
    "skipLibCheck": true,
    "module": "ESNext",
    "moduleResolution": "bundler",
    "allowSyntheticDefaultImports": true,
    "strict": true
  },
  "include": ["vite.config.ts"]
}
```

- [x] **Step 4: Create index.html**

Create `gover2/frontend/index.html`:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Gover2</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.ts"></script>
  </body>
</html>
```

- [x] **Step 5: Install frontend dependencies**

```bash
cd /home/mg/inventory/gover2/frontend
pnpm install
```

Expected: `node_modules/` created, `pnpm-lock.yaml` written. No npm warnings. Tailwind version installed should be 3.x.x — confirm with `pnpm list tailwindcss`.

- [x] **Step 6: Create TypeScript domain interfaces**

Create `gover2/frontend/src/lib/types/index.ts`:

```ts
// All domain interfaces for Gover2.
// Wails auto-converts Go snake_case JSON tags to TypeScript camelCase.

export interface Computer {
  id: number
  name: string
  model: string
  userId: number | null
  status: string
  purchaseDate: string | null
  warrantyExpiry: string | null
  notes: string | null
  createdAt: string
  updatedAt: string
}

// Smartphone has identical shape to Computer.
export interface Smartphone {
  id: number
  name: string
  model: string
  userId: number | null
  status: string
  purchaseDate: string | null
  warrantyExpiry: string | null
  notes: string | null
  createdAt: string
  updatedAt: string
}

// Tablet has identical shape to Computer.
export interface Tablet {
  id: number
  name: string
  model: string
  userId: number | null
  status: string
  purchaseDate: string | null
  warrantyExpiry: string | null
  notes: string | null
  createdAt: string
  updatedAt: string
}

export interface WindowsKey {
  id: number
  licenseKey: string
  computerId: number | null
  status: string
  notes: string | null
  createdAt: string
  updatedAt: string
}

export interface Antivirus {
  id: number
  name: string
  licenseKey: string
  computerId: number | null
  smartphoneId: number | null
  tabletId: number | null
  status: string
  expiryDate: string | null
  notes: string | null
  createdAt: string
  updatedAt: string
}

// OtherSoftware has identical shape to Antivirus.
export interface OtherSoftware {
  id: number
  name: string
  licenseKey: string
  computerId: number | null
  smartphoneId: number | null
  tabletId: number | null
  status: string
  expiryDate: string | null
  notes: string | null
  createdAt: string
  updatedAt: string
}

export interface User {
  id: number
  name: string
  surname: string | null
  status: string
  notes: string | null
  createdAt: string
  updatedAt: string
}

export interface AppConfig {
  dbPath: string
  density: 'comfortable' | 'compact'
  darkMode: boolean
  expiryWarningDays: number
}

export interface Alert {
  category: string
  id: number
  name: string
  expiryDate: string
  daysRemaining: number
  severity: 'expired' | 'expiring'
}
```

- [x] **Step 7: Create API bridge wrapper layer**

Create `gover2/frontend/src/lib/api/index.ts`:

```ts
// All Wails bridge calls go through this file.
// In scaffold phase all methods return empty/stub values via Promise.resolve().
// DO NOT import from wailsjs/ here — those bindings are generated by wails build
// and do not exist yet. Real imports are added in gover2-readonly.
// Components MUST NOT call window.go.* directly.

import type {
  Antivirus,
  AppConfig,
  Computer,
  OtherSoftware,
  Smartphone,
  Tablet,
  User,
  WindowsKey,
} from '@/lib/types'

export function listComputers(): Promise<Computer[]> {
  return Promise.resolve([])
}

export function listSmartphones(): Promise<Smartphone[]> {
  return Promise.resolve([])
}

export function listTablets(): Promise<Tablet[]> {
  return Promise.resolve([])
}

export function listWindowsKeys(): Promise<WindowsKey[]> {
  return Promise.resolve([])
}

export function listAntivirus(): Promise<Antivirus[]> {
  return Promise.resolve([])
}

export function listOtherSoftware(): Promise<OtherSoftware[]> {
  return Promise.resolve([])
}

export function listUsers(): Promise<User[]> {
  return Promise.resolve([])
}

export function getDatabasePath(): Promise<string> {
  return Promise.resolve('')
}

export function getConfig(): Promise<AppConfig> {
  return Promise.resolve({
    dbPath: '',
    density: 'comfortable',
    darkMode: false,
    expiryWarningDays: 30,
  })
}

export function setConfig(_cfg: AppConfig): Promise<void> {
  return Promise.resolve()
}

export function newDatabase(_path: string): Promise<void> {
  return Promise.resolve()
}

export function openDatabase(_path: string): Promise<void> {
  return Promise.resolve()
}

export function exportCSV(_category: string, _destPath: string): Promise<void> {
  return Promise.resolve()
}
```

- [x] **Step 8: Commit**

```bash
cd /home/mg/inventory
git add gover2/frontend/
git commit -m "feat(gover2): frontend package setup with TS types and API stub layer"
```

archived-with: 2026-06-17-gover2-scaffold
---

### Task 4: Pinia Stores

**Files:**
- Create: `gover2/frontend/src/stores/assets.ts`
- Create: `gover2/frontend/src/stores/ui.ts`

**Interfaces:**
- Consumes: `Computer`, `Smartphone`, `Tablet`, `WindowsKey`, `Antivirus`, `OtherSoftware`, `User` from `@/lib/types`
- Produces:
  - `useAssetsStore()` — exposes `currentCategory`, `items`, `selectedId`, `searchQuery`, `statusFilter`
  - `useUiStore()` — exposes `density`, `darkMode`, `sidebarCollapsed`

- [x] **Step 1: Create assets store**

Create `gover2/frontend/src/stores/assets.ts`:

```ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type {
  Antivirus,
  Computer,
  OtherSoftware,
  Smartphone,
  Tablet,
  User,
  WindowsKey,
} from '@/lib/types'

export type AssetCategory =
  | 'computers'
  | 'smartphones'
  | 'tablets'
  | 'all'
  | 'windows-keys'
  | 'antivirus'
  | 'other-software'
  | 'users'

export type AssetItem =
  | Computer
  | Smartphone
  | Tablet
  | WindowsKey
  | Antivirus
  | OtherSoftware
  | User

export const useAssetsStore = defineStore('assets', () => {
  const currentCategory = ref<AssetCategory>('computers')
  const items = ref<AssetItem[]>([])
  const selectedId = ref<number | null>(null)
  const searchQuery = ref('')
  const statusFilter = ref('')

  function setCategory(category: AssetCategory) {
    currentCategory.value = category
    selectedId.value = null
    items.value = []
  }

  function setItems(newItems: AssetItem[]) {
    items.value = newItems
  }

  function setSelectedId(id: number | null) {
    selectedId.value = id
  }

  return {
    currentCategory,
    items,
    selectedId,
    searchQuery,
    statusFilter,
    setCategory,
    setItems,
    setSelectedId,
  }
})
```

- [x] **Step 2: Create UI store**

Create `gover2/frontend/src/stores/ui.ts`:

```ts
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUiStore = defineStore('ui', () => {
  const density = ref<'comfortable' | 'compact'>('comfortable')
  const darkMode = ref(false)
  const sidebarCollapsed = ref(false)

  function toggleDarkMode() {
    darkMode.value = !darkMode.value
  }

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setDensity(d: 'comfortable' | 'compact') {
    density.value = d
  }

  return {
    density,
    darkMode,
    sidebarCollapsed,
    toggleDarkMode,
    toggleSidebar,
    setDensity,
  }
})
```

- [x] **Step 3: Commit**

```bash
cd /home/mg/inventory
git add gover2/frontend/src/stores/
git commit -m "feat(gover2): add Pinia assets and ui stores"
```

archived-with: 2026-06-17-gover2-scaffold
---

### Task 5: Vue Router + App Shell

**Files:**
- Create: `gover2/frontend/src/router/index.ts`
- Create: `gover2/frontend/src/App.vue`
- Create: `gover2/frontend/src/main.ts`
- Create: `gover2/frontend/src/layouts/AppShell.vue`

**Interfaces:**
- Consumes: `useAssetsStore` from `@/stores/assets`, `useUiStore` from `@/stores/ui`; view components from `@/features/` (Tasks 6)
- Produces: mounted Vue app with 9 routes and AppShell layout wrapping all views

Note: Task 5 is written before Task 6 (view stubs). The router imports views by path — if typecheck is run before Task 6, it will fail. Run typecheck after Task 6 is complete.

- [x] **Step 1: Create router**

Create `gover2/frontend/src/router/index.ts`:

```ts
import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  { path: '/', redirect: '/computers' },
  {
    path: '/computers',
    component: () => import('@/features/assets/ComputersView.vue'),
  },
  {
    path: '/smartphones',
    component: () => import('@/features/assets/SmartphonesView.vue'),
  },
  {
    path: '/tablets',
    component: () => import('@/features/assets/TabletsView.vue'),
  },
  {
    path: '/all',
    component: () => import('@/features/assets/AllAssetsView.vue'),
  },
  {
    path: '/windows-keys',
    component: () => import('@/features/licenses/WindowsKeysView.vue'),
  },
  {
    path: '/antivirus',
    component: () => import('@/features/licenses/AntivirusView.vue'),
  },
  {
    path: '/other-software',
    component: () => import('@/features/licenses/OtherSoftwareView.vue'),
  },
  {
    path: '/users',
    component: () => import('@/features/users/UsersView.vue'),
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
```

- [x] **Step 2: Create App.vue**

Create `gover2/frontend/src/App.vue`:

```vue
<template>
  <AppShell>
    <router-view />
  </AppShell>
</template>

<script setup lang="ts">
import AppShell from '@/layouts/AppShell.vue'
</script>
```

- [x] **Step 3: Create main.ts entry point**

Create `gover2/frontend/src/main.ts`:

```ts
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from '@/router'
import App from '@/App.vue'
import '@/style.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
```

Note: `style.css` is created in Task 8. The import here will cause a typecheck error until that file exists. Create a placeholder file now if needed:

```bash
touch /home/mg/inventory/gover2/frontend/src/style.css
```

- [x] **Step 4: Create AppShell layout**

Create `gover2/frontend/src/layouts/AppShell.vue`:

```vue
<template>
  <div class="flex h-screen bg-page font-ui">
    <!-- Sidebar -->
    <Sidebar />

    <!-- Main content area -->
    <div class="flex flex-col flex-1 overflow-hidden">
      <!-- Top bar -->
      <Topbar />

      <!-- Page content -->
      <main class="flex-1 overflow-auto p-4">
        <slot />
      </main>

      <!-- Status bar -->
      <Statusbar />
    </div>
  </div>
</template>

<script setup lang="ts">
import Sidebar from '@/components/Sidebar.vue'
import Topbar from '@/components/Topbar.vue'
import Statusbar from '@/components/Statusbar.vue'
</script>
```

- [x] **Step 5: Commit**

```bash
cd /home/mg/inventory
git add gover2/frontend/src/router/ gover2/frontend/src/App.vue gover2/frontend/src/main.ts gover2/frontend/src/layouts/
git commit -m "feat(gover2): add Vue Router hash-mode routes and AppShell layout"
```

archived-with: 2026-06-17-gover2-scaffold
---

### Task 6: Feature View Stubs

**Files:**
- Create: `gover2/frontend/src/features/assets/ComputersView.vue`
- Create: `gover2/frontend/src/features/assets/SmartphonesView.vue`
- Create: `gover2/frontend/src/features/assets/TabletsView.vue`
- Create: `gover2/frontend/src/features/assets/AllAssetsView.vue`
- Create: `gover2/frontend/src/features/licenses/WindowsKeysView.vue`
- Create: `gover2/frontend/src/features/licenses/AntivirusView.vue`
- Create: `gover2/frontend/src/features/licenses/OtherSoftwareView.vue`
- Create: `gover2/frontend/src/features/users/UsersView.vue`
- Create: `gover2/frontend/src/features/alerts/AlertsView.vue` (Phase 4 placeholder — not routed)

**Interfaces:**
- Consumes: `useAssetsStore` from `@/stores/assets`
- Produces: 9 `<template>` stubs that render a placeholder `<div>` with the category name; no `window.go.*` calls

- [x] **Step 1: Create ComputersView.vue**

Create `gover2/frontend/src/features/assets/ComputersView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Computers</h1>
    <DataTable :items="[]" />
  </div>
</template>

<script setup lang="ts">
import DataTable from '@/components/DataTable.vue'
</script>
```

- [x] **Step 2: Create SmartphonesView.vue**

Create `gover2/frontend/src/features/assets/SmartphonesView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Smartphones</h1>
    <DataTable :items="[]" />
  </div>
</template>

<script setup lang="ts">
import DataTable from '@/components/DataTable.vue'
</script>
```

- [x] **Step 3: Create TabletsView.vue**

Create `gover2/frontend/src/features/assets/TabletsView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Tablets</h1>
    <DataTable :items="[]" />
  </div>
</template>

<script setup lang="ts">
import DataTable from '@/components/DataTable.vue'
</script>
```

- [x] **Step 4: Create AllAssetsView.vue**

Create `gover2/frontend/src/features/assets/AllAssetsView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">All Assets</h1>
    <DataTable :items="[]" />
  </div>
</template>

<script setup lang="ts">
import DataTable from '@/components/DataTable.vue'
</script>
```

- [x] **Step 5: Create WindowsKeysView.vue**

Create `gover2/frontend/src/features/licenses/WindowsKeysView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Windows Keys</h1>
    <DataTable :items="[]" />
  </div>
</template>

<script setup lang="ts">
import DataTable from '@/components/DataTable.vue'
</script>
```

- [x] **Step 6: Create AntivirusView.vue**

Create `gover2/frontend/src/features/licenses/AntivirusView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Antivirus</h1>
    <DataTable :items="[]" />
  </div>
</template>

<script setup lang="ts">
import DataTable from '@/components/DataTable.vue'
</script>
```

- [x] **Step 7: Create OtherSoftwareView.vue**

Create `gover2/frontend/src/features/licenses/OtherSoftwareView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Other Software</h1>
    <DataTable :items="[]" />
  </div>
</template>

<script setup lang="ts">
import DataTable from '@/components/DataTable.vue'
</script>
```

- [x] **Step 8: Create UsersView.vue**

Create `gover2/frontend/src/features/users/UsersView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Users</h1>
    <DataTable :items="[]" />
  </div>
</template>

<script setup lang="ts">
import DataTable from '@/components/DataTable.vue'
</script>
```

- [x] **Step 9: Create AlertsView.vue (Phase 4 placeholder — not routed)**

Create `gover2/frontend/src/features/alerts/AlertsView.vue`:

```vue
<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Alerts</h1>
    <p class="text-text-secondary">Alerts view — Phase 4 placeholder.</p>
  </div>
</template>

<script setup lang="ts">
// Phase 4 placeholder — no logic yet
</script>
```

- [x] **Step 10: Commit**

```bash
cd /home/mg/inventory
git add gover2/frontend/src/features/
git commit -m "feat(gover2): add 9 feature view stubs (scaffold placeholders)"
```

archived-with: 2026-06-17-gover2-scaffold
---

### Task 7: Shared Component Stubs

**Files:**
- Create: `gover2/frontend/src/components/Sidebar.vue`
- Create: `gover2/frontend/src/components/SidebarItem.vue`
- Create: `gover2/frontend/src/components/Topbar.vue`
- Create: `gover2/frontend/src/components/DataTable.vue`
- Create: `gover2/frontend/src/components/StatusBadge.vue`
- Create: `gover2/frontend/src/components/EditPanel.vue`
- Create: `gover2/frontend/src/components/FormField.vue`
- Create: `gover2/frontend/src/components/AlertsModal.vue`
- Create: `gover2/frontend/src/components/Statusbar.vue`

**Interfaces:**
- Consumes: Tailwind classes from `tailwind.config.ts` (Task 8); `useUiStore` from `@/stores/ui`; router-link for sidebar nav
- Produces: 9 stub components; `DataTable` accepts `items: unknown[]` prop; `StatusBadge` accepts `status: string` prop

- [x] **Step 1: Create Sidebar.vue**

Create `gover2/frontend/src/components/Sidebar.vue`:

```vue
<template>
  <nav
    class="flex flex-col bg-sidebar border-r border-sidebar-border"
    :style="{ width: '200px' }"
  >
    <div class="p-4 text-sidebar-hi font-semibold text-sm border-b border-sidebar-border">
      Inventory
    </div>
    <div class="flex flex-col gap-0.5 p-2 flex-1">
      <SidebarItem to="/all" label="All Assets" />
      <SidebarItem to="/computers" label="Computers" />
      <SidebarItem to="/smartphones" label="Smartphones" />
      <SidebarItem to="/tablets" label="Tablets" />
      <div class="my-2 border-t border-sidebar-border" />
      <SidebarItem to="/windows-keys" label="Windows Keys" />
      <SidebarItem to="/antivirus" label="Antivirus" />
      <SidebarItem to="/other-software" label="Other Software" />
      <div class="my-2 border-t border-sidebar-border" />
      <SidebarItem to="/users" label="Users" />
    </div>
  </nav>
</template>

<script setup lang="ts">
import SidebarItem from '@/components/SidebarItem.vue'
</script>
```

- [x] **Step 2: Create SidebarItem.vue**

Create `gover2/frontend/src/components/SidebarItem.vue`:

```vue
<template>
  <router-link
    :to="to"
    class="block px-3 py-1.5 rounded text-sm text-sidebar-text hover:bg-sidebar-hover hover:text-sidebar-hi transition-colors"
    active-class="bg-sidebar-active text-sidebar-hi"
  >
    {{ label }}
  </router-link>
</template>

<script setup lang="ts">
defineProps<{
  to: string
  label: string
}>()
</script>
```

- [x] **Step 3: Create Topbar.vue**

Create `gover2/frontend/src/components/Topbar.vue`:

```vue
<template>
  <header class="flex items-center px-4 bg-surface border-b border-border" style="height: 48px;">
    <span class="text-sm text-text-secondary">{{ title }}</span>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const title = computed(() => String(route.path).replace('/', '').replace('-', ' ') || 'Inventory')
</script>
```

- [x] **Step 4: Create DataTable.vue**

Create `gover2/frontend/src/components/DataTable.vue`:

```vue
<template>
  <div class="bg-surface border border-border rounded-none overflow-hidden">
    <table class="w-full text-sm">
      <thead class="bg-th-bg">
        <tr>
          <th class="text-left px-4 py-2 text-th-text font-medium">Name</th>
          <th class="text-left px-4 py-2 text-th-text font-medium">Status</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="items.length === 0">
          <td colspan="2" class="px-4 py-8 text-center text-text-secondary">
            No items to display.
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  items: unknown[]
}>()
</script>
```

- [x] **Step 5: Create StatusBadge.vue**

Create `gover2/frontend/src/components/StatusBadge.vue`:

```vue
<template>
  <span
    class="inline-block px-1.5 py-0.5 text-xs font-medium text-white rounded-badge"
    :class="badgeClass"
  >
    {{ status }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  status: string
}>()

const badgeClass = computed(() => {
  const map: Record<string, string> = {
    active: 'bg-s-active',
    repair: 'bg-s-repair',
    spare: 'bg-s-spare',
    retired: 'bg-s-retired',
    missing: 'bg-s-missing',
    expiring: 'bg-s-expiring',
    expired: 'bg-s-expired',
  }
  return map[props.status.toLowerCase()] ?? 'bg-s-retired'
})
</script>
```

- [x] **Step 6: Create EditPanel.vue**

Create `gover2/frontend/src/components/EditPanel.vue`:

```vue
<template>
  <aside class="w-80 bg-surface border-l border-border p-4 overflow-y-auto">
    <p class="text-text-secondary text-sm">Select an item to edit.</p>
  </aside>
</template>

<script setup lang="ts">
// Stub — real implementation in gover2-crud (Phase 3)
</script>
```

- [x] **Step 7: Create FormField.vue**

Create `gover2/frontend/src/components/FormField.vue`:

```vue
<template>
  <div class="flex flex-col gap-1">
    <label class="text-xs font-medium text-text-secondary">{{ label }}</label>
    <slot />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  label: string
}>()
</script>
```

- [x] **Step 8: Create AlertsModal.vue**

Create `gover2/frontend/src/components/AlertsModal.vue`:

```vue
<template>
  <div v-if="open" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
    <div class="bg-surface border border-border rounded-none p-6 w-[480px] max-h-[70vh] overflow-y-auto">
      <h2 class="text-base font-semibold text-text-primary mb-4">Alerts</h2>
      <p class="text-text-secondary text-sm">No alerts. (Phase 4 placeholder)</p>
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

- [x] **Step 9: Create Statusbar.vue**

Create `gover2/frontend/src/components/Statusbar.vue`:

```vue
<template>
  <footer class="flex items-center px-4 h-6 bg-sidebar border-t border-sidebar-border text-sidebar-text text-xs">
    <span>Ready</span>
  </footer>
</template>

<script setup lang="ts">
// Stub — real content added in gover2-readonly
</script>
```

- [x] **Step 10: Commit**

```bash
cd /home/mg/inventory
git add gover2/frontend/src/components/
git commit -m "feat(gover2): add 9 shared component stubs"
```

archived-with: 2026-06-17-gover2-scaffold
---

### Task 8: Tailwind CSS v3 Configuration

**Files:**
- Create: `gover2/frontend/tailwind.config.ts`
- Create: `gover2/frontend/postcss.config.js`
- Modify: `gover2/frontend/src/style.css` (replace placeholder from Task 5)

**Interfaces:**
- Produces: all 25 design tokens available as Tailwind utility classes (e.g. `bg-sidebar`, `text-text-primary`, `bg-s-active`, `rounded-badge`, `font-ui`)

- [x] **Step 1: Create tailwind.config.ts**

Create `gover2/frontend/tailwind.config.ts`:

```ts
import type { Config } from 'tailwindcss'

export default {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {
      colors: {
        // Sidebar palette
        'sidebar':        '#0E1520',
        'sidebar-hover':  '#131E2E',
        'sidebar-active': '#182337',
        'sidebar-border': '#141D2B',
        'sidebar-text':   '#6B82A0',
        'sidebar-hi':     '#DDE8F5',
        'sidebar-accent': '#2563EB',
        // Page palette
        'page':    '#EEF1F7',
        'surface': '#FFFFFF',
        'border':  '#DDE1EC',
        // Status colors
        's-active':   '#2563EB',
        's-repair':   '#EA580C',
        's-spare':    '#7C3AED',
        's-retired':  '#9CA3AF',
        's-missing':  '#1E293B',
        's-expiring': '#D97706',
        's-expired':  '#DC2626',
        // Table header
        'th-bg':   '#0E1520',
        'th-text': '#8FA5BF',
        // Typography
        'text-primary':   '#0F172A',
        'text-secondary': '#64748B',
        'text-tertiary':  '#94A3B8',
        // Interactive
        'accent-hover': '#1D4ED8',
        // Page structure
        'border-light': '#EAECF4',
        // Status row highlight backgrounds
        'row-expiring': '#FFFBEB',
        'row-expired':  '#FFF5F5',
      },
      fontFamily: {
        ui:   ['-apple-system', 'BlinkMacSystemFont', '"Segoe UI"', 'system-ui', 'sans-serif'],
        mono: ['"Cascadia Code"', '"Fira Code"', 'Consolas', '"Courier New"', 'monospace'],
      },
      borderRadius: {
        badge: '2px',
      },
      width: {
        sidebar: '200px',
      },
      height: {
        topbar: '48px',
      },
    },
  },
  plugins: [],
} satisfies Config
```

The `content` array includes `./src/**/*.{vue,ts}` so dynamic class strings like `bg-s-active` are not purged by Tailwind's JIT engine.

- [x] **Step 2: Create postcss.config.js**

Create `gover2/frontend/postcss.config.js`:

```js
export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
```

Note: do NOT use `@tailwindcss/postcss` (that is the v4 PostCSS plugin). The standard `tailwindcss` key is correct for v3.

- [x] **Step 3: Replace style.css with Tailwind directives**

Replace `gover2/frontend/src/style.css` entirely:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

- [x] **Step 4: Commit**

```bash
cd /home/mg/inventory
git add gover2/frontend/tailwind.config.ts gover2/frontend/postcss.config.js gover2/frontend/src/style.css
git commit -m "feat(gover2): configure Tailwind CSS v3 with all 25 design tokens"
```

archived-with: 2026-06-17-gover2-scaffold
---

### Task 9: Taskfile Integration

**Files:**
- Modify: `Taskfile.yml` (repo root, append 9 `gover:` tasks)

**Interfaces:**
- Produces: `task gover:install`, `task gover:dev`, `task gover:build`, `task gover:test`, `task gover:lint`, `task gover:typecheck`, `task gover:clean`, `task gover:package:deb`, `task gover:package:rpm` — all runnable from repo root

- [x] **Step 1: Append gover: tasks to root Taskfile.yml**

Open `/home/mg/inventory/Taskfile.yml`. Append the following block after the last existing task (`clean:` around line 271), before the end of the file:

```yaml

# ──────────────────────────────── gover2 ─────────────────────────────────────
  gover:install:
    desc: Download Go deps and install frontend packages for gover2
    dir: gover2
    cmds:
      - go mod download
      - pnpm install --prefix frontend

  gover:dev:
    desc: Launch gover2 in development mode (hot reload)
    dir: gover2
    cmds:
      - wails dev

  gover:build:
    desc: Build the gover2 desktop binary via wails
    dir: gover2
    cmds:
      - wails build

  gover:test:
    desc: Run Go and frontend tests for gover2
    dir: gover2
    cmds:
      - go test ./...
      - pnpm --prefix frontend run test

  gover:lint:
    desc: Lint Go and frontend code for gover2
    dir: gover2
    cmds:
      - golangci-lint run ./...
      - pnpm --prefix frontend run lint

  gover:typecheck:
    desc: Run vue-tsc type check for gover2 frontend
    dir: gover2
    cmds:
      - pnpm --prefix frontend run typecheck

  gover:clean:
    desc: Remove gover2 build artifacts
    dir: gover2
    cmds:
      - rm -rf build/
      - pnpm --prefix frontend run clean

  gover:package:deb:
    desc: "Build gover2 .deb package (stub — fpm command filled in Phase 5)"
    dir: gover2
    cmds:
      - wails build -platform linux/amd64

  gover:package:rpm:
    desc: "Build gover2 .rpm package (stub — fpm command filled in Phase 5)"
    dir: gover2
    cmds:
      - wails build -platform linux/amd64
```

- [x] **Step 2: Verify task list includes gover: tasks**

```bash
cd /home/mg/inventory
task --list | grep gover
```

Expected to show all 9 `gover:*` tasks. If Taskfile syntax is invalid (e.g. indentation error), `task --list` will print a parse error — fix the YAML indentation before continuing.

- [x] **Step 3: Verify gover:install runs**

```bash
cd /home/mg/inventory
task gover:install
```

Expected: `go mod download` completes, then `pnpm install --prefix frontend` installs packages in `gover2/frontend/node_modules/`. Exit 0.

- [x] **Step 4: Commit**

```bash
cd /home/mg/inventory
git add Taskfile.yml
git commit -m "feat(gover2): add 9 gover: Taskfile commands"
```

archived-with: 2026-06-17-gover2-scaffold
---

### Task 10: Verification

**Files:** none new

**Interfaces:**
- Consumes: all tasks above

- [x] **Step 1: Go build — all packages, no import cycles**

```bash
cd /home/mg/inventory/gover2
go build ./...
```

Expected: exits 0, no output. If you see `import cycle not allowed`, check which package is importing `bridge` — only `bridge` is allowed to import all others.

- [x] **Step 2: TypeScript typecheck — strict, 0 errors**

```bash
cd /home/mg/inventory/gover2
pnpm --prefix frontend run typecheck
```

Expected: exits 0, no output or only file-list output. Any `error TS` lines must be fixed before proceeding.

Common issues:
- `Cannot find module '@/...'` — check `tsconfig.json` `paths` and `vite.config.ts` `alias`
- `Property 'X' does not exist on type 'never'` — a store or prop is typed `never`; check the generic type parameter
- `Argument of type 'unknown[]' is not assignable` — DataTable prop type mismatch; `items: unknown[]` must be consistent

- [x] **Step 3: Verify no window.go.* calls**

```bash
grep -r "window\.go\." /home/mg/inventory/gover2/frontend/src/features/ /home/mg/inventory/gover2/frontend/src/components/ || echo "CLEAN"
```

Expected output: `CLEAN`. Any match is a violation of the `api-wrapper` requirement.

- [x] **Step 4: wails build — produces binary**

```bash
cd /home/mg/inventory/gover2
wails build -tags webkit2_41
```

Expected: exits 0. Binary produced at `gover2/build/bin/gover2`. This step:
1. Calls `pnpm install` in `frontend/` (via `wails.json` `frontend:install`)
2. Calls `pnpm run build` to produce `frontend/dist/`
3. Embeds the frontend and compiles the Go binary

Note: On Debian/Ubuntu systems with WebKit2GTK 4.1 (instead of 4.0), pass `-tags webkit2_41`.
The `wails.json` has been updated with `"build:tags": "webkit2_41"` for future builds.

- [x] **Step 5: Confirm binary exists**

```bash
ls -lh /home/mg/inventory/gover2/build/bin/gover2
```

Expected: a non-zero-size ELF binary. (Actual: 11.9M)

- [x] **Step 6: Manual smoke test — wails dev** (PENDING — requires display)

```bash
cd /home/mg/inventory/gover2
wails dev
```

Expected behaviour:
- A desktop window opens (Wails WebView).
- Sidebar shows 9 navigation items: All Assets, Computers, Smartphones, Tablets, Windows Keys, Antivirus, Other Software, Users — 7 categories + All Assets.
- Clicking a sidebar item routes to the corresponding view (showing "No items to display.").
- No console errors in the browser devtools (F12).

This step is manual — press Ctrl+C to exit `wails dev` when done.

- [x] **Step 7: Final commit**

```bash
cd /home/mg/inventory
git add -p  # review and stage any unstaged cleanup changes
git commit -m "feat(gover2): scaffold complete — go build, typecheck, and wails build pass"
```

archived-with: 2026-06-17-gover2-scaffold
---

## Self-Review Checklist

**Spec coverage:**

| Spec requirement | Covered by |
|---|---|
| Toolchain install (Wails v2, pnpm) | Task 0 |
| `wails init -t vue` | Task 1 |
| `wails.json` pnpm override (critical per design doc) | Task 1, Step 5 |
| `modernc.org/sqlite` blank import in database stub | Task 2, Step 2 |
| 6 Go packages in dependency order | Task 2, Steps 1–6 |
| `app.go` delegates to `bridge.App` | Task 2, Step 7 |
| `go build ./...` exits 0 | Task 2, Step 8 + Task 10, Step 1 |
| `frontend/src/lib/types/index.ts` with 9 interfaces | Task 3, Step 6 |
| `lib/api/index.ts` with 13 stubs returning `Promise.resolve([])` | Task 3, Step 7 |
| No `window.go.*` calls anywhere | Task 10, Step 3 |
| Pinia stores (assets + ui) | Task 4 |
| Vue Router hash mode, 9 routes | Task 5, Step 1 |
| `AppShell.vue` layout | Task 5, Step 4 |
| 8 routed feature view stubs + AlertsView placeholder | Task 6 |
| 9 shared component stubs | Task 7 |
| Tailwind v3 pinned, not v4 | Task 3, Step 1 (devDeps) + Task 8, Step 1 |
| All 25 design tokens in `tailwind.config.ts` | Task 8, Step 1 |
| PostCSS config (not `@tailwindcss/vite`) | Task 8, Step 2 |
| `@tailwind` directives in `style.css` | Task 8, Step 3 |
| `"strict": true` in tsconfig | Task 3, Step 3 |
| `pnpm run typecheck` exits 0 | Task 10, Step 2 |
| 9 `gover:` Taskfile tasks | Task 9 |
| `task gover:install` runs | Task 9, Step 3 |
| `wails build` exits 0 with binary | Task 10, Step 4 |
| `wails dev` smoke test | Task 10, Step 6 |
| No unit tests in scaffold | (no test files created — by design) |
| No SQL connections | (all stubs return nil DB) |
| No packaging scripts (Phase 5) | `gover:package:*` are stubs |

**Placeholder scan:** No TBDs, no "implement later" notes. All code blocks are complete. The only intentional placeholders are the `gover:package:deb` and `gover:package:rpm` stub commands (packaging is explicitly deferred to Phase 5 per the design doc).

**Type consistency:** `DataTable` accepts `items: unknown[]` in the component definition and callers pass `items="[]"` (literal empty array). `StatusBadge` accepts `status: string` consistently. `AppConfig` shape matches between Go struct JSON tags and TypeScript interface (camelCase). `AssetItem` union in `stores/assets.ts` covers all 7 domain types matching the 7 models structs.
