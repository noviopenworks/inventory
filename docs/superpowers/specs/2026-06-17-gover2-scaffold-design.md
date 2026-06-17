---
comet_change: gover2-scaffold
role: technical-design
canonical_spec: openspec
---

# Technical Design: gover2-scaffold

## Context

The architecture is locked in `openspec/changes/archive/2026-06-17-gover-phase1-architecture/specs/architecture/spec.md`. This change executes that blueprint — no design decisions remain open. The Design Doc records the implementation-level details that translate the locked spec into concrete build steps.

## Implementation Approach

Pure scaffolding in dependency order:

1. Install toolchain (Wails v2, pnpm)
2. `wails init` to create the Go project skeleton
3. Add Go package stubs (`internal/`) in dependency order
4. Replace generated frontend with the locked Vue structure
5. Configure Tailwind CSS v3 with all 25 design tokens
6. Add `gover:` Taskfile commands to root `Taskfile.yml`
7. Verify with `go build ./...` + `pnpm run typecheck` + `wails build`

## Wails + pnpm Integration

Wails v2 defaults to `npm`. Using pnpm requires explicit overrides in `gover2/wails.json`:

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

**Critical**: set these fields before any `wails build` or `wails dev` invocation. Without them, Wails calls `npm install` and the build fails because npm is not configured.

## Go Package Order and Stubs

Packages must be created in dependency order (each layer imports only from layers below it):

| Order | Package | Key export | Imports |
|---|---|---|---|
| 1 | `internal/models` | 7 domain structs + `AppConfig` + `Alert` | nothing |
| 2 | `internal/database` | `Open(path string) (*sql.DB, error)` | `models`, `_ "modernc.org/sqlite"` |
| 3 | `internal/config` | `AppConfig`, `Load()`, `Save()` | `models` |
| 4 | `internal/backup` | `BackupDB(srcPath, destDir string) error` | `models` |
| 5 | `internal/services` | `ListComputers(db) ([]models.Computer, error)` × 7 | `database`, `models` |
| 6 | `internal/bridge` | `App` struct with all bridge methods | `services`, `config`, `backup` |

### modernc.org/sqlite blank import

The `database` stub must include `_ "modernc.org/sqlite"` even though it uses no SQLite calls yet. Without a real import, `go mod tidy` removes the dependency from `go.sum`, breaking `gover2-readonly`.

```go
import (
    "database/sql"
    _ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) { return nil, nil }
```

### Bridge stub in app.go

Wails generates `app.go` with a placeholder `App` struct. We replace it to delegate to `bridge.App`:

```go
// app.go
package main

import "gover2/internal/bridge"

func NewApp() *bridge.App { return bridge.NewApp() }
```

`main.go` calls `wails.Run(app.NewOptions{Ctx: ctx, Bind: []interface{}{NewApp()}})`.

## Frontend Structure

`wails init -t vue` generates a minimal Vite + Vue frontend. We keep the generated scaffolding (`index.html`, `vite.config.ts`) and replace `src/` entirely with the locked structure:

```
frontend/src/
├── main.ts                  # Pinia + Router + App mount
├── App.vue                  # <router-view> wrapper
├── layouts/
│   └── AppShell.vue         # sidebar + topbar slot + <router-view>
├── features/
│   ├── assets/              # ComputersView, SmartphonesView, TabletsView, AllAssetsView stubs
│   ├── licenses/            # WindowsKeysView, AntivirusView, OtherSoftwareView stubs
│   ├── users/               # UsersView stub
│   └── alerts/              # AlertsView stub (Phase 4 placeholder)
├── components/              # 9 shared components (all stubs)
├── lib/
│   ├── api/index.ts         # All bridge wrappers returning Promise.resolve([])
│   └── types/index.ts       # 7 domain interfaces + AppConfig + Alert
└── stores/
    ├── assets.ts            # currentCategory, items, selectedId, filters
    └── ui.ts                # density, darkMode, sidebarCollapsed
```

### TypeScript bindings deferred

Wails generates `frontend/wailsjs/go/bridge/App.js` (and `.d.ts`) during `wails build`. In scaffold phase these files do not exist yet. `lib/api/index.ts` must **not** import from `wailsjs/` — it returns `Promise.resolve([])` stubs directly. The real import path is added in `gover2-readonly` after the first successful `wails build` generates the bindings.

## Tailwind CSS v3 Setup

Pin `tailwindcss@3` explicitly. Tailwind v4 (released 2025) switched to CSS `@theme` directives and is incompatible with `tailwind.config.ts`. The locked 25-token config requires v3 format.

```bash
pnpm add -D tailwindcss@3 postcss autoprefixer @vitejs/plugin-vue vue-tsc
```

`vite.config.ts` uses the standard Vue plugin (not `@tailwindcss/vite` which is v4-only). PostCSS config:

```js
// postcss.config.js
export default {
  plugins: { tailwindcss: {}, autoprefixer: {} }
}
```

`src/style.css` entry point:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

`tailwind.config.ts` content array must include `./src/**/*.{vue,ts}` to prevent purging of dynamic class strings like `bg-s-active` used in `StatusBadge`.

## Taskfile Integration

9 `gover:` tasks appended to root `Taskfile.yml` under the existing task definitions. Each task has `dir: gover2` so it runs in the Wails project directory regardless of where `task` is invoked.

Packaging tasks (`gover:package:deb`, `gover:package:rpm`) are stubs in the scaffold — the `fpm` packaging command is filled in during Phase 5.

## Verification Sequence

In order (each step must exit 0 before the next):

1. `go build ./...` — all 6 packages compile, no import cycles
2. `pnpm --prefix frontend run typecheck` — vue-tsc strict, 0 errors
3. `wails build` — binary produced at `gover2/build/bin/gover2`
4. Manual: `wails dev` → window opens, sidebar shows 7 categories, no console errors

## Non-Goals for This Change

- No SQLite connection (all DB calls return nil in stubs)
- No real bridge data (all list methods return empty slices)
- No unit tests (added in gover2-readonly with testify + Vitest)
- No packaging scripts (Phase 5)
