# Design: "Keep small files" architecture pass

**Date:** 2026-06-20
**Status:** Approved (brainstorming) — pending implementation plan
**Branch:** committed onto `feat/go-default-cutover` (PR #7).
**Topic:** Break the codebase's large files into smaller, focused, per-entity files — within the existing layered packages, with no behavior change.

---

## 1. Goal

Several files have grown to span all seven asset/entity types at once, making them large and unfocused. Split the clear offenders (files **clearly over ~200 lines**) into smaller per-entity files, keeping the existing layered package structure (`models/`, `services/`, `bridge/`, frontend `components/`). This is a pure structural refactor: **no behavior changes**, the existing test suite is the regression net.

## 2. Decisions

| # | Decision | Choice |
|---|----------|--------|
| 1 | Decomposition strategy | Split **within the existing layers, by entity** — no package restructure, no import-path churn. |
| 2 | Threshold | **Light touch:** split only files clearly over ~200 lines; leave smaller files alone. |
| 3 | Scope | Go source, Go tests, and frontend — limited to the >200 offenders. |
| 4 | `models.go` (exactly 200) | **Include** the split (keeps the package consistent with the new `services/` layout). |
| 5 | Barrels (`api/index.ts` 190, `types/index.ts` 179) | **Leave as-is** — under threshold and cohesive single-purpose modules. |
| 6 | Branch | Commit onto `feat/go-default-cutover` (PR #7). |
| 7 | Behavior | Pure refactor; no logic edits; tests must pass unchanged. |

## 3. Files in scope and their target split

### 3.1 `internal/models/models.go` (200) → per-entity files
- `computer.go`, `smartphone.go`, `tablet.go`, `windowskey.go`, `antivirus.go`, `othersoftware.go`, `user.go` — each holds that entity's record struct + its `*Input` DTO.
- `shared.go` — cross-entity types: `Alert`, `AppConfig`, `DropdownItem`, `DeviceDropdownItem`, and any shared enums/helpers.
- Package stays `models`; no API change.

### 3.2 `internal/services/{services.go (221), write.go (238)}` → per-entity files
- Merge each entity's read query + insert/update/delete into one file: `computer.go`, `smartphone.go`, `tablet.go`, `windowskey.go`, `antivirus.go`, `othersoftware.go`, `user.go` (~50 lines each).
- `alerts.go` — the cross-entity `GetAlerts` (and any alert helpers).
- **Keep** `export.go` (110) and `dropdowns.go` (52) as-is (under threshold, already concern-focused).
- If shared row-scan/query helpers exist, put them in `helpers.go`.
- `services.go` and `write.go` are deleted once emptied.

### 3.3 `internal/bridge/bridge.go` (386) → per-entity + core
- `app.go` — `App` struct, `NewApp`, `NewAppWithDB`, `Startup`, `ErrNoDB`, and the `list[T]`/`withDB` helpers (the package core).
- `computer.go` … `user.go` — each entity's bridge methods (`List*`, `Add*`, `Update*`, `Delete*`).
- `database.go` — DB lifecycle/config methods: `GetDatabasePath`, `GetConfig`, `SetConfig`, `NewDatabase`, `OpenDatabase`, `NewDatabaseDialog`, `OpenDatabaseDialog`, `BackupDatabase`, `BackupDatabaseDialog`.
- `queries.go` — `GetAlerts`, `ListUsersForDropdown`, `ListDevicesForDropdown`, and `ExportCSV` (or `ExportCSV` in its own `export.go` if `queries.go` would exceed ~150).
- Package stays `bridge`; **all public method signatures unchanged**, so generated `frontend/wailsjs/` bindings do not change.

### 3.4 `frontend/src/components/EditPanel.vue` (460) → shell + 3 modules
EditPanel's per-entity variation is already data-driven via maps, so its natural split is **by concern**:
- `EditPanel.vue` — keeps the template + thin orchestration (props, computed, lifecycle, `onSave` wiring). Target ~150 lines.
- `components/editpanel/fieldConfigs.ts` — `FieldConfig` interface, `FIELD_CONFIGS`, status constants + `STATUS_MAP`, `CATEGORY_LABELS`.
- `components/editpanel/payload.ts` — `buildPayload`, `deriveDeviceSelect`, `deriveComputerSelect`, `nullOrStr`, `nullOrNum`, `initForm` (pure functions taking category/form/row).
- `components/editpanel/save.ts` — `callApi` (the per-category add/update dispatch).
- `EditPanel.vue` imports from these; behavior and the rendered DOM/`data-testid`s are unchanged so `EditPanel.test.ts` passes untouched.

### 3.5 Test files mirror the source split
- `internal/services/services_test.go` (516) → per-entity `computer_test.go` … `user_test.go` + `alerts_test.go`, with shared helpers (`openTestDB`) moved to `helpers_test.go`. `export_test.go` (118) stays.
- `internal/bridge/bridge_test.go` (585) → per-entity `*_test.go` + `database_test.go` (DB-mgmt/backup tests) + shared helpers (`newBridgeWithDB`) in `helpers_test.go`.
- Test names and assertions are unchanged — only their file location moves.

## 4. Explicitly NOT split (under the light-touch threshold)
`api/index.ts` (190), `types/index.ts` (179), `Topbar.vue` (150), `database.go` (121), `export.go` (110), the per-view `*View.vue` files (~100), `config.go` (75), and all other files under ~200 lines stay as they are.

## 5. Approach

Mechanical, **one package/area at a time**, in dependency order so each step builds on a green tree:

1. `models` split
2. `services` split + its test split
3. `bridge` split + its test split
4. `EditPanel.vue` split

Each step: move code into the new files (no logic edits), `gofmt`/build, run the **unchanged** test suite (all green = behavior preserved), `task lint` (0 issues), commit. For the frontend step, `pnpm typecheck` + the Vitest suite are the gates.

## 6. Acceptance criteria

- No source file in scope exceeds ~200 lines after the split (test files may sit slightly higher but are materially smaller than before).
- Every package keeps its public API: Go exported identifiers and bridge method signatures unchanged; `frontend/wailsjs/` bindings unchanged; `EditPanel.vue`'s props/emits and rendered `data-testid`s unchanged.
- `task check` (lint + typecheck + test) and `task test:cov` (≥75%) pass unchanged.
- No `services.go`/`write.go`/old monolithic test files remain (deleted once emptied).

## 7. Out of scope
- Any behavior, logic, or API change.
- Splitting files under ~200 lines (barrels, views, Topbar, etc.).
- Restructuring into per-domain packages (vertical slices) — rejected in favor of split-within-layers.
- New tests or coverage changes (the 75% floor from the prior pass must simply continue to hold).
