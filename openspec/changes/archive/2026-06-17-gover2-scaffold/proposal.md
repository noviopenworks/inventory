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
