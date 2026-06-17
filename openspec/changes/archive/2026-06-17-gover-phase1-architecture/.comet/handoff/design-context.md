# Comet Design Handoff

- Change: gover-phase1-architecture
- Phase: design
- Mode: compact
- Context hash: 9ffe5a99d6fd93c29651726eda93e507afa9eb8a9eb726940734b55adba09e4a

Generated-by: comet-handoff.sh

OpenSpec remains the canonical capability spec. This handoff is a deterministic, source-traceable context pack, not an agent-authored summary.

## openspec/changes/gover-phase1-architecture/proposal.md

- Source: openspec/changes/gover-phase1-architecture/proposal.md
- Lines: 1-38
- SHA256: bfd7f928851b210c26fc7a13ae81b96ee9e046c671165b3d7fae0f752745b9e4

```md
## Why

The IT Asset Inventory app is currently a single-file PyQt6 desktop application. The existing code mixes UI, data access, and business logic, making it difficult to extend safely. The rewrite targets Wails + Go + Vue 3 + TypeScript + Tailwind CSS — a stack that separates frontend and backend concerns cleanly and produces a proper desktop executable without a Python runtime dependency.

Phase 0 (Discovery) is complete for UI/UX: all decisions are locked in `openspec/changes/archive/2026-06-17-gover-ui-spec/`. The visual design reference is `gover/UI_DESIGN_REFERENCE.html`. The remaining prerequisite for Phase 2 (read-only prototype) is a locked technical architecture document that any developer or agent can use to scaffold the project without making ad-hoc decisions.

## What Changes

- Produces `specs/architecture/spec.md`: Go package layout, Vue app structure, Wails bridge method signatures, Tailwind config strategy mapping the locked design tokens, test tooling choices, and Taskfile `gover:` command names.
- Updates `gover/ARCHITECTURE_NOTES.md` with a resolved status linking to this spec.
- No application code is written. The current Python app is not modified.

## Capabilities

### New Capabilities

- `architecture`: Complete technical blueprint for the Gover v2 Wails app. Covers Go package boundaries, Vue feature structure, TypeScript API client shape, Tailwind config, test tooling, and Taskfile command names. Sufficient to begin Phase 2 scaffolding without further architectural decisions.

### Modified Capabilities

_(none)_

## Impact

- Unblocks: Phase 2 (read-only Wails prototype scaffold)
- Does not affect: current `app/` Python code, `db/` layer, packaging
- Platform target: Linux first (`.deb`, `.rpm`); Windows deferred

## Key Decisions Already Made

| Decision | Choice |
|---|---|
| Database strategy | Reuse current SQLite schema directly (`SCHEMA_VERSION=1`) |
| Go SQLite driver | `modernc.org/sqlite` (pure Go, no CGo) |
| Frontend package manager | pnpm |
| Taskfile | Root `Taskfile.yml`, `gover:` namespace |
| First platform | Linux; Windows deferred |
| UI/UX | All decisions locked in `gover-ui-spec` |
```

## openspec/changes/gover-phase1-architecture/design.md

- Source: openspec/changes/gover-phase1-architecture/design.md
- Lines: 1-74
- SHA256: 2d224b4acff42783b70e95be9642a638adaed32cd33cfcb3ddcb945bce7c4cf7

```md
## Approach

Produce a single capability spec (`specs/architecture/spec.md`) that documents every technical decision a developer needs to scaffold Phase 2. The spec is written in Q&A format, mirroring the UI decisions spec from `gover-ui-spec`. It is the authoritative reference for package names, bridge method signatures, Taskfile command names, and Tailwind config structure.

A second small update annotates `gover/ARCHITECTURE_NOTES.md` to mark the architecture section as resolved, parallel to how `gover/QUESTIONS.md` was annotated in `gover-ui-spec`.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      Wails Desktop App                      │
├────────────────────────┬────────────────────────────────────┤
│     Go Backend         │        Vue 3 Frontend              │
│  ─────────────────     │  ──────────────────────────────    │
│  internal/database     │  src/layouts/AppShell.vue          │
│  internal/models       │  src/features/assets/              │
│  internal/services     │  src/features/users/               │
│  internal/bridge       │  src/features/alerts/              │
│  internal/config       │  src/lib/api/          (TS)        │
│  internal/backup       │  src/lib/types/        (TS)        │
│                        │  src/components/                   │
└────────────────────────┴────────────────────────────────────┘
         ▲                          ▲
         │  Wails generated         │  pnpm + Vite
         │  bindings                │
         └──────────── bridge ──────┘
```

## Go Package Design

- Module path: `github.com/user/inventory` (or local `inventory` module — decided in spec)
- All internal packages under `internal/` (not importable by external tools)
- Package names match ARCHITECTURE_NOTES.md candidate areas: `database`, `models`, `services`, `bridge`, `config`, `backup`
- `database` owns schema open/init and carries forward the current 7-table DDL verbatim
- `models` defines typed Go structs matching every current SQLite column
- `services` implements all operations; `bridge` exposes only what Vue needs

## Vue App Design

- Vite + Vue 3 Composition API + TypeScript strict mode
- Feature-folder structure under `src/features/` (not component-type folders)
- `src/lib/api/` wraps Wails-generated bindings — Vue components never call `window.go.*` directly
- Pinia for state management
- Vue Router for sidebar navigation (hash mode, no server needed)

## Tailwind Config Strategy

- Design tokens from `gover-ui-spec` map directly to `tailwind.config.ts` `theme.extend`
- No `@import` of external fonts — system font stacks only
- All CSS custom properties in the spec become Tailwind color/spacing utilities
- `tailwind.config.ts` is the single source of truth; no inline style fallbacks in components

## Test Tooling

- Go: standard `testing` package + `testify/assert` for assertions; table-driven tests
- Frontend unit tests: Vitest (Vite-native, no separate webpack config)
- Frontend component tests: Vue Test Utils + Vitest
- No E2E tests in Phase 1 or 2; E2E deferred to Phase 4 (Playwright candidate)

## Taskfile Commands (`gover:` namespace)

All commands live in root `Taskfile.yml` under `gover:` namespace:

| Command | Purpose |
|---|---|
| `gover:install` | Install Go deps + pnpm install |
| `gover:dev` | `wails dev` (hot-reload) |
| `gover:build` | `wails build` (production binary) |
| `gover:test` | Go tests + frontend tests |
| `gover:lint` | `golangci-lint` + `eslint` |
| `gover:typecheck` | `vue-tsc --noEmit` |
| `gover:clean` | Remove build artifacts |
| `gover:package:deb` | Build `.deb` package (Linux) |
| `gover:package:rpm` | Build `.rpm` package (Linux) |
```

## openspec/changes/gover-phase1-architecture/tasks.md

- Source: openspec/changes/gover-phase1-architecture/tasks.md
- Lines: 1-41
- SHA256: 75f883fec6a8a22924661e90cfceaccd30279a6223edc8b99cc7acb11ecf7031

```md
## 1. Go Package Architecture

- [ ] 1.1 Document Go module path and `internal/` package layout
- [ ] 1.2 Specify `internal/database` — connection open, WAL pragma, schema init (DDL carried from current schema)
- [ ] 1.3 Specify `internal/models` — Go structs for all 7 tables with column-to-field mapping
- [ ] 1.4 Specify `internal/services` — method signatures for list, get, create, update, delete per category
- [ ] 1.5 Specify `internal/bridge` — Wails-exposed method names and signatures (the public API surface)
- [ ] 1.6 Specify `internal/config` — app config struct, file location on Linux (`$XDG_CONFIG_HOME/inventory/`)
- [ ] 1.7 Specify `internal/backup` — backup file path convention and trigger (pre-write, manual)

## 2. Vue + TypeScript App Structure

- [ ] 2.1 Document `src/` directory layout (features/, layouts/, lib/, components/)
- [ ] 2.2 Specify Vue Router routes (one route per sidebar category + All overview)
- [ ] 2.3 Specify Pinia store shape (assets store, ui store)
- [ ] 2.4 Specify `src/lib/api/` TypeScript wrapper shape — mirror of bridge methods
- [ ] 2.5 Specify `src/lib/types/` — TypeScript interfaces for all 7 domain entities

## 3. Tailwind Configuration

- [ ] 3.1 Map all design tokens from `gover-ui-spec` to `tailwind.config.ts` color palette entries
- [ ] 3.2 Specify font family entries (UI stack, mono stack) in Tailwind config
- [ ] 3.3 Specify border-radius scale (0 for structural, 2px for badges → Tailwind custom values)
- [ ] 3.4 Specify any custom spacing/sizing tokens (sidebar-w: 200px, topbar-h: 48px)

## 4. Test Strategy

- [ ] 4.1 Document Go test conventions (table-driven, testify/assert, coverage target)
- [ ] 4.2 Document frontend test setup (Vitest config, Vue Test Utils, test file co-location)
- [ ] 4.3 Specify which current Python tests have Go equivalents to write in Phase 2

## 5. Taskfile Commands

- [ ] 5.1 Write the `gover:` namespace task definitions (all 9 commands from design.md)
- [ ] 5.2 Verify each command is runnable on Linux (check tool availability requirements)
- [ ] 5.3 Add `gover:` commands to root `Taskfile.yml`

## 6. Cross-Reference Update

- [ ] 6.1 Annotate `gover/ARCHITECTURE_NOTES.md` Phase 1 section with resolved status and spec link
- [ ] 6.2 Verify no application code in diff
```

## openspec/changes/gover-phase1-architecture/specs/architecture/spec.md

- Source: openspec/changes/gover-phase1-architecture/specs/architecture/spec.md
- Lines: 1-315
- SHA256: 6411e27b3bb678ec2b93cf4d5533c5b15ba8722b2b9cba11537845cec192be68

[TRUNCATED]

```md
## ADDED Capabilities

### Capability: architecture

Complete technical blueprint for the Gover v2 Wails application. Covers Go package boundaries and module layout, Wails bridge method signatures for Phase 2 (read-only) and naming conventions for Phase 3 (CRUD), Vue 3 + TypeScript frontend structure, Tailwind CSS configuration mapping all locked design tokens, test tooling conventions and coverage targets, and root Taskfile `gover:` command definitions. Sufficient to begin Phase 2 scaffolding without making ad-hoc technical decisions.

---

## Section 1: Project Structure & Go Package Layout

The Wails project lives at `gover2/` inside the current `inventory/` repository. The Python app (`app/`, `db/`) is untouched during Phase 1 — the two live side by side.

```
inventory/
├── app/                  ← Python app (unchanged)
├── db/                   ← Python db layer (unchanged)
├── gover/                ← planning docs
├── gover2/               ← Wails project root
│   ├── go.mod            ← module gover2, go 1.22+
│   ├── main.go
│   ├── app.go
│   ├── internal/
│   │   ├── database/
│   │   ├── models/
│   │   ├── services/
│   │   ├── bridge/
│   │   ├── config/
│   │   └── backup/
│   └── frontend/
│       └── src/
└── Taskfile.yml          ← gover: namespace, dir: gover2
```

**Go module:** `module gover2`, `go 1.22`. All application code is under `internal/` — nothing is exported from the module root for external import.

**SQLite driver:** `modernc.org/sqlite` — pure Go, no CGo, cross-compiles without a C toolchain.

### Go Package Responsibilities

| Package | Responsibility |
|---|---|
| `internal/database` | Open connection, set `PRAGMA journal_mode=WAL`, run schema init DDL verbatim from `db/schema.py` |
| `internal/models` | Typed Go structs for all 7 tables; column-to-field mapping (`snake_case` → `CamelCase`) |
| `internal/services` | Business logic: list, get, create, update, delete per category; alert calculation (`EXPIRY_WARNING_DAYS=30`) |
| `internal/bridge` | `App` struct; all Wails-exported methods; thin delegation to `services` — no SQL here |
| `internal/config` | `AppConfig` struct; JSON persistence at `$XDG_CONFIG_HOME/inventory/config.json` |
| `internal/backup` | `BackupDB(srcPath, destDir string) error`; copies DB file; triggered pre-write and manually |

**Dependency rule:** `models` → nothing internal. `database` → `models`. `services` → `database`, `models`. `bridge` → `services`, `config`, `backup`. `config`, `backup` → `models`.

### Requirement: go-package-layout

The `internal/` package graph must satisfy the dependency rule above. `bridge` is the only package that imports from all others. No package may import `bridge`.

#### Scenario: circular import detected

- Given a `go build ./...` run on `gover2/`
- When any package other than `bridge` imports `bridge`
- Then the build fails with an import cycle error

---

## Section 2: Wails Bridge Method Signatures

`bridge.App` is the Wails context receiver. All exported methods are callable from TypeScript via Wails-generated bindings. Vue components **never** call `window.go.*` directly — all calls go through `src/lib/api/`.

### Phase 2 — Read-only methods (implement in Phase 2)

| Method | Go signature | TS wrapper |
|---|---|---|
| `ListComputers` | `() ([]models.Computer, error)` | `api.listComputers()` |
| `ListSmartphones` | `() ([]models.Smartphone, error)` | `api.listSmartphones()` |
| `ListTablets` | `() ([]models.Tablet, error)` | `api.listTablets()` |
| `ListWindowsKeys` | `() ([]models.WindowsKey, error)` | `api.listWindowsKeys()` |
| `ListAntivirus` | `() ([]models.Antivirus, error)` | `api.listAntivirus()` |
| `ListOtherSoftware` | `() ([]models.OtherSoftware, error)` | `api.listOtherSoftware()` |
| `ListUsers` | `() ([]models.User, error)` | `api.listUsers()` |
| `GetDatabasePath` | `() (string, error)` | `api.getDatabasePath()` |
| `GetConfig` | `() (models.AppConfig, error)` | `api.getConfig()` |
| `SetConfig` | `(cfg models.AppConfig) error` | `api.setConfig(cfg)` |
```

Full source: openspec/changes/gover-phase1-architecture/specs/architecture/spec.md
