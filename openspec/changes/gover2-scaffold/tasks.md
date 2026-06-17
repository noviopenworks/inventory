## 0. Toolchain Prerequisites

- [x] 0.1 Install Wails v2: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- [x] 0.2 Install pnpm: `corepack enable pnpm`
- [x] 0.3 Verify: `wails version` and `pnpm --version` both print

## 1. Wails Project Init

- [x] 1.1 Run `wails init -n gover2 -t vue -d gover2 -g` (or equivalent) from inventory root
- [x] 1.2 Verify `gover2/go.mod` contains `module gover2`
- [x] 1.3 Verify `gover2/main.go` and `gover2/app.go` exist
- [x] 1.4 Add `modernc.org/sqlite` to `gover2/go.mod` via `go get modernc.org/sqlite`
- [x] 1.5 Verify `go build ./...` passes in `gover2/`

## 2. Go Package Stubs

- [x] 2.1 Create `internal/models/models.go` with 7 domain structs + AppConfig + Alert (exact field mapping from architecture spec §3)
- [x] 2.2 Create `internal/database/database.go` with `Open(path string) (*sql.DB, error)` stub returning nil
- [x] 2.3 Create `internal/config/config.go` with `AppConfig` struct and `Load()/Save()` stubs
- [x] 2.4 Create `internal/backup/backup.go` with `BackupDB(srcPath, destDir string) error` stub returning nil
- [x] 2.5 Create `internal/services/services.go` with `ListComputers`, `ListSmartphones`, `ListTablets`, `ListWindowsKeys`, `ListAntivirus`, `ListOtherSoftware`, `ListUsers` — all return `([]models.X, nil)`
- [x] 2.6 Create `internal/bridge/bridge.go` with `App` struct and all 10 Phase-2 bridge methods + 3 file-dialog stubs, delegating to services
- [x] 2.7 Update `app.go` to use `bridge.App` as the Wails context struct
- [x] 2.8 Verify `go build ./...` passes with no import cycles

## 3. Vue Frontend Structure

- [x] 3.1 Replace generated `frontend/` with pnpm-based Vite + Vue 3 + TypeScript project
- [x] 3.2 Create `frontend/src/lib/types/index.ts` with all 7 domain interfaces (Computer, Smartphone, Tablet, WindowsKey, Antivirus, OtherSoftware, User) + AppConfig + Alert
- [x] 3.3 Create `frontend/src/lib/api/index.ts` wrapping all 10 bridge methods (returning `Promise.resolve([])` for list methods in scaffold)
- [x] 3.4 Create `frontend/src/stores/assets.ts` (currentCategory, items, selectedId, filters)
- [x] 3.5 Create `frontend/src/stores/ui.ts` (density, darkMode, sidebarCollapsed)
- [ ] 3.6 Create `frontend/src/layouts/AppShell.vue` with sidebar + topbar + main content slot
- [x] 3.7 Create `frontend/src/components/` — 9 stub components (Sidebar, SidebarItem, Topbar, DataTable, StatusBadge, EditPanel, FormField, AlertsModal, Statusbar)
- [x] 3.8 Create `frontend/src/features/` — 4 feature folders (assets, licenses, users, alerts) with stub Index.vue per category view (Computers, Smartphones, Tablets, All, WindowsKeys, Antivirus, OtherSoftware, Users)
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
