# Development Plan

Roadmap for the Wails + Go + Vue + TypeScript rewrite of IT Asset Inventory.
This document started as a planning artifact; phases are now tracked here as
they are delivered. The current PyQt6 app remains the behavior reference until
gover reaches feature parity.

## Objectives

- Modernize the desktop experience while preserving the useful behavior of the current PyQt6 app.
- Keep the app local-first, with no server requirement and no internet dependency for normal use.
- Preserve or migrate existing SQLite data safely.
- Improve maintainability by separating backend domain logic, database access, frontend views, and desktop packaging concerns.
- Define quality gates before feature implementation begins.

## Phase Status

| Phase | Status | Delivered By |
|-------|--------|--------------|
| 0 — Discovery And Decisions | ✅ Done | `openspec/changes/archive/2026-06-17-gover-phase1-architecture`, `2026-06-17-gover-ui-spec` |
| 1 — Technical Foundation | ✅ Done | `openspec/changes/archive/2026-06-17-gover2-scaffold` |
| 2 — Read-Only Prototype | ✅ Done | `openspec/changes/archive/2026-06-18-gover2-readonly` |
| 3 — Core CRUD | ✅ Done | `openspec/changes/archive/2026-06-18-gover2-crud` |
| 4 — Alerts, Search, Export, Settings | ✅ Done | `openspec/changes/archive/2026-06-18-gover2-parity` |
| 5 — Packaging And Migration | ⏳ Pending | `task gover:package:deb` / `:rpm` are stubs |
| 6 — Polish And Release Candidate | ⏳ Pending | — |

## Phase 0: Discovery And Decisions

Purpose: decide what the new version must preserve, improve, or remove.

Work items:

- Review current user workflows: asset list, add, edit, delete, export, alerts, dark mode, and settings.
- Inventory the current database schema and query behavior.
- Decide whether the new version uses the existing schema directly or introduces a new schema plus migration. → **Decision:** reuse current schema verbatim (see `DECISIONS.md`).
- Define supported platforms and packaging targets.
- Define the minimum viable first release.
- Resolve the open questions (captured in `DECISIONS.md`).

## Phase 1: Technical Foundation

Purpose: prepare the Wails/Vue/Go architecture before writing app features.

Decisions (locked):

- Wails project layout: root `gover/` with `internal/` Go packages and `frontend/` Vue app.
- Go package boundaries: `database`, `models`, `services`, `bridge`, `config`, `backup`. `bridge` is the sole package that imports from all others.
- Vue application structure: feature folders (`features/assets`, `features/users`, `features/alerts`, `features/settings`), shared `components/`, `lib/api/`, `lib/types/`, Pinia stores, hash-mode Vue Router.
- Tailwind configuration: CSS variables (`--c-*`) for theme tokens, with `:root` (light) and `.dark` overrides.
- Testing tools: Go `testing` + `stretchr/testify`; frontend Vitest + `@vue/test-utils`.
- Taskfile commands: see root `Taskfile.yml` under the `gover:*` namespace.

## Phase 2: Read-Only Prototype

Purpose: prove the stack can open the local database and display inventory data without changing records.

Delivered:

- Wails app shell (`main.go`, `app.go`, `internal/bridge`).
- Go SQLite connection and read queries (`internal/database`, `internal/services`).
- TypeScript API bindings (`frontend/wailsjs/`).
- Vue layout with read-only tables (`frontend/src/layouts`, `features/assets`).
- Basic error handling for missing or invalid database files.

## Phase 3: Core CRUD

Purpose: replace the current add, edit, and delete flows.

Delivered:

- Create, update, and delete users and asset records (`internal/services/write.go`).
- Form validation.
- Relationship selectors for assigned users and linked devices.
- Confirmation flows for destructive actions.
- Data refresh after changes.

## Phase 4: Alerts, Search, Export, And Settings

Purpose: restore supporting features and improve daily usability.

Delivered:

- Expiry and warranty alerts (`internal/services` alert queries + `features/alerts` modal).
- Client-side search filtering across all views.
- CSV export (`internal/services/export.go`).
- Dark mode (manual toggle, themed sidebar + table header) and persisted appearance settings.
- Action-row topbar with new/open database dialogs.
- About and License modals.

Exit criteria:

- Users can find expiring assets quickly.
- CSV export matches the PyQt6 app's output.
- Appearance and table preferences persist.

## Phase 5: Packaging And Migration

Purpose: make the new app installable and safe for existing users. **Not yet started.**

Future work:

- Linux packages: `.deb`, `.rpm`, and `.tar.gz` mirroring the PyQt6 outputs.
- Windows `.exe` build via `wails build -platform windows/amd64`.
- Version sync with `.app-version` and release automation.
- Database backup before migrations.
- Document the upgrade path from the PyQt6 database (schema is reused as-is, so this is primarily a "point gover at the existing `.db` file" workflow).

Exit criteria:

- Packages can be built reproducibly.
- Upgrade path is documented.
- Existing data can be preserved or imported.

## Phase 6: Polish And Release Candidate

Purpose: make gover stable enough to replace the PyQt6 app. **Not yet started.**

Future work:

- UI polish and accessibility pass.
- Keyboard navigation review (all PyQt6 shortcuts should be preserved).
- Performance review with realistic data volumes.
- Regression testing against current PyQt6 app behavior.
- Documentation updates and release notes.

Exit criteria:

- Release candidate is approved.
- Known limitations are documented.
- PyQt6 app maintenance plan is defined (freeze, archive, or continue patching).

## First Milestone (Completed)

The first development milestone was a read-only Wails/Vue shell that opens the
existing SQLite database and displays the current asset categories. This was
delivered in Phase 2 and de-risked migration before building forms and write
behavior.
