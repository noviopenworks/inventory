# Comet Design Handoff

- Change: gover2-scaffold
- Phase: design
- Mode: compact
- Context hash: ae7c9a46e66b726f9a30c3445ad53fd4e17b95e6a058c265146e2ccef182c052

Generated-by: comet-handoff.sh

OpenSpec remains the canonical capability spec. This handoff is a deterministic, source-traceable context pack, not an agent-authored summary.

## openspec/changes/gover2-scaffold/proposal.md

- Source: openspec/changes/gover2-scaffold/proposal.md
- Lines: 1-46
- SHA256: 31817f4f8673e8a62e7abbbd211f168a2486c393f5255314da3e3b4b72bbfb03

```md
## Why

Phase 1 produced a locked technical architecture for the Gover v2 Wails rewrite (`openspec/changes/archive/2026-06-17-gover-phase1-architecture/`). Phase 2 begins by scaffolding the `gover2/` project skeleton: installing the required toolchain, initialising the Wails project, creating all 6 Go internal packages as stubs, setting up the Vue 3 + TypeScript + Tailwind CSS frontend with the locked directory structure, and adding the 9 `gover:` Taskfile commands to the root `Taskfile.yml`.

The result is a compiling, runnable Wails window that shows the full sidebar and AppShell UI with no real data — a verified skeleton ready for `gover2-readonly` to wire in database access.

## What Changes

- Creates `gover2/` directory from `wails init` with `module gover2`
- Implements 6 Go package stubs: `internal/database`, `internal/models`, `internal/services`, `internal/bridge`, `internal/config`, `internal/backup` — each compiles but contains no real business logic
- Implements full Vue frontend skeleton: `AppShell.vue`, `Sidebar.vue`, 9 hash-mode routes with empty data table stubs, 2 Pinia stores (assets + ui), `lib/api/` returning hardcoded empty slices, `lib/types/` with all 7 TypeScript interfaces
- Adds `tailwind.config.ts` with all 25 design tokens from `gover-ui-spec`
- Adds 9 `gover:` tasks to root `Taskfile.yml`

## Capabilities

### New Capabilities

- `gover2-runnable-skeleton`: The `gover2/` Wails project compiles (`go build ./...`), builds (`wails build`), and launches (`wails dev`) showing an empty but correctly structured sidebar + topbar + main area with zero real data.

## Impact

- Unblocks: `gover2-readonly` (data access and display)
- Does not affect: current `app/` Python code, `db/` layer, packaging
- Platform target: Linux (same as architecture spec)

## Key Decisions Already Made (from gover-phase1-architecture)

| Decision | Choice |
|---|---|
| Project location | `gover2/` subdirectory in `inventory/` repo |
| Go module | `module gover2`, Go 1.22+ |
| SQLite driver | `modernc.org/sqlite` (added to go.mod now; used in gover2-readonly) |
| Frontend package manager | pnpm |
| Vue Router | Hash mode |
| Tailwind design tokens | 25 tokens from gover-ui-spec |
| Taskfile namespace | `gover:` in root `Taskfile.yml` |

## Non-Goals

- No database connection or SQLite access (gover2-readonly)
- No bridge methods returning real data (stubs return empty slices)
- No CRUD operations (Phase 3)
- No alerts or CSV export
- No Windows packaging or Flatpak (Phase 5)
- No production packaging
```

## openspec/changes/gover2-scaffold/design.md

- Source: openspec/changes/gover2-scaffold/design.md
- Lines: 1-122
- SHA256: 218a71b9bd715d129091b46e09a73d916174dc0d1ffc708e41b7308f0a11ccbd

[TRUNCATED]

```md
## Approach

Follow the locked architecture from `gover-phase1-architecture` exactly. This change is pure scaffolding — no business logic, no database, no real data. Every implementation decision is already made; this change executes them.

## Architecture Overview

```
gover2/
├── go.mod              module gover2, go 1.22, wails v2 + modernc.org/sqlite deps
├── main.go             wails.Run(app.NewApp())
├── app.go              App struct — placeholder, no real bridge methods yet
├── wails.json          app config
├── internal/
│   ├── database/       Open() stub — returns nil DB for now
│   ├── models/         7 structs (Computer, Smartphone, Tablet, WindowsKey, Antivirus, OtherSoftware, User) + AppConfig + Alert
│   ├── services/       ListComputers() etc — return ([]models.X, nil) stubs
│   ├── bridge/         App struct with List* methods delegating to services stubs
│   ├── config/         AppConfig JSON load/save — reads $XDG_CONFIG_HOME/inventory/config.json
│   └── backup/         BackupDB() stub — returns nil
└── frontend/
    ├── package.json    pnpm workspace, Vue 3 + Vite + TypeScript
    ├── tailwind.config.ts  25 design tokens from gover-ui-spec
    ├── tsconfig.json   strict: true
    └── src/
        ├── main.ts
        ├── App.vue
        ├── layouts/AppShell.vue    sidebar + topbar + main content slot
        ├── features/               8 feature folders, each with an Index.vue stub
        ├── components/             9 shared components (stubs)
        ├── lib/api/index.ts        bridge wrappers returning Promise<[]>
        ├── lib/types/index.ts      7 domain interfaces + AppConfig + Alert
        └── stores/
            ├── assets.ts           Pinia store — empty items state
            └── ui.ts               density, darkMode, sidebarCollapsed
```

## Wails Init Strategy

Wails v2 `wails init` generates a project with a bundled Go+Vue template. We use the vanilla Vue template and immediately replace the generated frontend with our locked structure. The init command creates:
- `go.mod`, `main.go`, `app.go`, `wails.json`
- `frontend/` with Vite + vanilla Vue

We then:
1. Replace `frontend/` content with our locked Vue structure
2. Add `modernc.org/sqlite` to `go.mod` (needed in gover2-readonly, good to have early)
3. Add all 6 `internal/` packages as stubs
4. Wire `app.go` to expose the bridge struct

## wails.json pnpm Configuration

Wails v2 defaults to `npm`. Using pnpm requires these fields in `wails.json` before any build:

```json
{
  "frontend:install": "pnpm install",
  "frontend:build": "pnpm run build",
  "frontend:dev:watcher": "pnpm run dev",
  "frontend:dev:serverUrl": "auto"
}
```

Without this, `wails build` calls `npm install` and fails.

## Toolchain Prerequisites

Install before scaffolding (one-time):
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
corepack enable pnpm
```

## Taskfile Integration

Add to root `Taskfile.yml` at repo root (not inside `gover2/`):

```yaml
tasks:
  gover:install:
    dir: gover2
    cmds:
```

Full source: openspec/changes/gover2-scaffold/design.md

## openspec/changes/gover2-scaffold/tasks.md

- Source: openspec/changes/gover2-scaffold/tasks.md
- Lines: 1-62
- SHA256: cf8113e1016f7325cd221341a4621128a4be31d65c3334b02b6f4dcd32eff44a

```md
## 0. Toolchain Prerequisites

- [ ] 0.1 Install Wails v2: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- [ ] 0.2 Install pnpm: `corepack enable pnpm`
- [ ] 0.3 Verify: `wails version` and `pnpm --version` both print

## 1. Wails Project Init

- [ ] 1.1 Run `wails init -n gover2 -t vue -d gover2 -g` (or equivalent) from inventory root
- [ ] 1.2 Verify `gover2/go.mod` contains `module gover2`
- [ ] 1.3 Verify `gover2/main.go` and `gover2/app.go` exist
- [ ] 1.4 Add `modernc.org/sqlite` to `gover2/go.mod` via `go get modernc.org/sqlite`
- [ ] 1.5 Verify `go build ./...` passes in `gover2/`

## 2. Go Package Stubs

- [ ] 2.1 Create `internal/models/models.go` with 7 domain structs + AppConfig + Alert (exact field mapping from architecture spec §3)
- [ ] 2.2 Create `internal/database/database.go` with `Open(path string) (*sql.DB, error)` stub returning nil
- [ ] 2.3 Create `internal/config/config.go` with `AppConfig` struct and `Load()/Save()` stubs
- [ ] 2.4 Create `internal/backup/backup.go` with `BackupDB(srcPath, destDir string) error` stub returning nil
- [ ] 2.5 Create `internal/services/services.go` with `ListComputers`, `ListSmartphones`, `ListTablets`, `ListWindowsKeys`, `ListAntivirus`, `ListOtherSoftware`, `ListUsers` — all return `([]models.X, nil)`
- [ ] 2.6 Create `internal/bridge/bridge.go` with `App` struct and all 10 Phase-2 bridge methods + 3 file-dialog stubs, delegating to services
- [ ] 2.7 Update `app.go` to use `bridge.App` as the Wails context struct
- [ ] 2.8 Verify `go build ./...` passes with no import cycles

## 3. Vue Frontend Structure

- [ ] 3.1 Replace generated `frontend/` with pnpm-based Vite + Vue 3 + TypeScript project
- [ ] 3.2 Create `frontend/src/lib/types/index.ts` with all 7 domain interfaces (Computer, Smartphone, Tablet, WindowsKey, Antivirus, OtherSoftware, User) + AppConfig + Alert
- [ ] 3.3 Create `frontend/src/lib/api/index.ts` wrapping all 10 bridge methods (returning `Promise.resolve([])` for list methods in scaffold)
- [ ] 3.4 Create `frontend/src/stores/assets.ts` (currentCategory, items, selectedId, filters)
- [ ] 3.5 Create `frontend/src/stores/ui.ts` (density, darkMode, sidebarCollapsed)
- [ ] 3.6 Create `frontend/src/layouts/AppShell.vue` with sidebar + topbar + main content slot
- [ ] 3.7 Create `frontend/src/components/` — 9 stub components (Sidebar, SidebarItem, Topbar, DataTable, StatusBadge, EditPanel, FormField, AlertsModal, Statusbar)
- [ ] 3.8 Create `frontend/src/features/` — 4 feature folders (assets, licenses, users, alerts) with stub Index.vue per category view (Computers, Smartphones, Tablets, All, WindowsKeys, Antivirus, OtherSoftware, Users)
- [ ] 3.9 Set up `frontend/src/router/index.ts` with 9 hash-mode routes (/ → /computers redirect + 8 views)
- [ ] 3.10 Wire `frontend/src/main.ts` with Pinia + Router + App mount

## 4. Tailwind Configuration

- [ ] 4.1 Install Tailwind CSS + Vite plugin in frontend
- [ ] 4.2 Create `frontend/tailwind.config.ts` with all 25 design tokens (from architecture spec §4)
- [ ] 4.3 Create `frontend/src/style.css` with `@tailwind` directives
- [ ] 4.4 Verify no inline style fallbacks in any component

## 5. Taskfile Commands

- [ ] 5.1 Add 9 `gover:` tasks to root `Taskfile.yml` (from architecture spec §6)
- [ ] 5.2 Verify `task gover:install` succeeds (downloads Go deps + pnpm install)
- [ ] 5.3 Verify `task gover:build` exits 0 (wails build produces binary)

## 6. TypeScript Verification

- [ ] 6.1 Ensure `frontend/tsconfig.json` has `"strict": true`
- [ ] 6.2 Run `pnpm --prefix frontend run typecheck` (vue-tsc) — must exit 0
- [ ] 6.3 Confirm no `window.go.*` calls in any feature or component file

## 7. Smoke Test

- [ ] 7.1 `go build ./...` exits 0 with no import cycle errors
- [ ] 7.2 `wails build` exits 0 — binary produced in `gover2/build/`
- [ ] 7.3 `wails dev` opens window — sidebar shows 7 categories, main area shows empty table, no console errors
```
