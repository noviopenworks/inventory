---
comet_change: gover-phase1-architecture
role: technical-design
canonical_spec: openspec
---

# Gover v2 Phase 1 — Technical Architecture Design Doc

## Context

The `gover-phase1-architecture` change produces the locked technical blueprint for the Gover v2 Wails rewrite. Phase 0 UI/UX decisions are complete (`gover-ui-spec`, archived). This change resolves the remaining Phase 1 decisions so that Phase 2 (read-only prototype scaffold) can begin without ad-hoc choices.

Canonical spec: `openspec/changes/gover-phase1-architecture/specs/architecture/spec.md`

---

## Goals / Non-Goals

**Goals:**
- Lock Go package boundaries, module path, and project location
- Define Wails bridge method signatures for Phase 2 (read-only) and establish naming pattern for Phase 3 (CRUD)
- Define Vue app directory structure, routes, stores, and API wrapper shape
- Map all locked design tokens to `tailwind.config.ts` entries
- Define test tooling conventions and coverage targets
- Write the 9 `gover:` Taskfile commands

**Non-Goals:**
- Writing any Go, Vue, TypeScript, or Wails application code
- Scaffolding the `gover2/` directory
- CRUD behavior spec (Phase 3)
- Packaging or migration spec (Phase 5)
- Windows platform setup (deferred)

---

## Decisions

### 1. Project Location & Module

The Wails project scaffolds at `gover2/` inside the current `inventory/` git repo. This keeps the Python app and Go app in the same git history during the transition period and allows the root `Taskfile.yml` to coordinate both.

```
inventory/
├── app/           ← Python app (untouched)
├── db/            ← Python db layer (untouched)
├── gover/         ← planning docs
├── gover2/        ← NEW: Wails project
│   ├── go.mod     ← module gover2
│   ├── main.go
│   ├── app.go
│   ├── internal/
│   └── frontend/
└── Taskfile.yml   ← gover: namespace, dir: gover2
```

Go module: `module gover2`, minimum Go version `1.22`. SQLite driver: `modernc.org/sqlite` (pure Go, no CGo, cross-compiles without a C toolchain).

### 2. Go Package Layout

All application code under `internal/` (unexported from module root). `bridge` is the only package that imports from all others.

| Package | Responsibility |
|---|---|
| `internal/database` | Connection open, WAL pragma (`PRAGMA journal_mode=WAL`), schema init — carries the 7-table DDL verbatim from `db/schema.py` |
| `internal/models` | Typed Go structs, one per table, column-to-field mapping (snake_case → CamelCase) |
| `internal/services` | Business operations: list, get, create, update, delete per category; alert calculation |
| `internal/bridge` | `App` struct receiver, all Wails-exported methods, thin delegation to services |
| `internal/config` | `AppConfig` struct, JSON persistence at `$XDG_CONFIG_HOME/inventory/config.json` |
| `internal/backup` | `BackupDB(srcPath, destDir string) error` — copies DB file; triggered manually and pre-write |

Dependency rule: `models` → nothing internal. `database` → `models`. `services` → `database`, `models`. `bridge` → `services`, `config`, `backup`. `config`, `backup` → `models`.

### 3. Wails Bridge Method Signatures

`bridge.App` is the Wails context receiver. All exported methods are callable from TypeScript via Wails-generated bindings, but Vue components always go through `src/lib/api/`.

**Phase 2 (read-only — implement in Phase 2):**

| Method | Signature | TS wrapper |
|---|---|---|
| ListComputers | `() ([]models.Computer, error)` | `api.listComputers()` |
| ListSmartphones | `() ([]models.Smartphone, error)` | `api.listSmartphones()` |
| ListTablets | `() ([]models.Tablet, error)` | `api.listTablets()` |
| ListWindowsKeys | `() ([]models.WindowsKey, error)` | `api.listWindowsKeys()` |
| ListAntivirus | `() ([]models.Antivirus, error)` | `api.listAntivirus()` |
| ListOtherSoftware | `() ([]models.OtherSoftware, error)` | `api.listOtherSoftware()` |
| ListUsers | `() ([]models.User, error)` | `api.listUsers()` |
| GetDatabasePath | `() (string, error)` | `api.getDatabasePath()` |
| GetConfig | `() (config.AppConfig, error)` | `api.getConfig()` |
| SetConfig | `(cfg config.AppConfig) error` | `api.setConfig(cfg)` |

**Phase 3 naming pattern (establish now, implement later):**
```
Get<Category>(id int) (models.<Category>, error)
Create<Category>(input models.Create<Category>Input) (models.<Category>, error)
Update<Category>(id int, input models.Update<Category>Input) (models.<Category>, error)
Delete<Category>(id int) error
```

**Phase 4:**
```
GetAlerts() ([]models.Alert, error)
```

Rule: `bridge` methods validate input and delegate to `services`. No SQL in `bridge`.

### 4. Vue App Structure & TypeScript Types

Vite + Vue 3 Composition API + TypeScript strict mode. Feature-folder structure. Vue Router hash mode (no server needed in Wails WebView). Pinia for state.

**`frontend/src/` layout:**
```
src/
├── main.ts
├── App.vue
├── layouts/AppShell.vue
├── features/
│   ├── assets/        ← computers, smartphones, tablets, all-view
│   ├── licenses/      ← windows keys, antivirus, other software
│   ├── users/
│   └── alerts/        ← Phase 4 stub
├── components/        ← 9 shared components from gover-ui-spec
│   ├── Sidebar.vue, SidebarItem.vue, Topbar.vue
│   ├── DataTable.vue, StatusBadge.vue
│   ├── EditPanel.vue, FormField.vue
│   ├── AlertsModal.vue, Statusbar.vue
├── lib/
│   ├── api/index.ts   ← all bridge wrappers, typed
│   └── types/index.ts ← all domain interfaces
└── stores/
    ├── assets.ts      ← current list, selected row, filters
    └── ui.ts          ← density, dark mode, sidebar collapse state
```

**Routes (hash mode):**

| Path | View | Sidebar item |
|---|---|---|
| `/` | redirect → `/computers` | — |
| `/all` | All Assets (hardware only) | All Assets |
| `/computers` | Computers | Computers |
| `/smartphones` | Smartphones | Smartphones |
| `/tablets` | Tablets | Tablets |
| `/windows-keys` | Windows Keys | Windows Keys |
| `/antivirus` | Antivirus | Antivirus |
| `/other-software` | Other Software | Other Software |
| `/users` | Users | Users |

**TypeScript interfaces (`lib/types/index.ts`):**

Wails auto-converts Go `snake_case` JSON tags to TypeScript `camelCase`.

```ts
interface Computer {
  id: number; name: string; model: string;
  userId: number | null; status: string;
  purchaseDate: string | null; warrantyExpiry: string | null;
  notes: string | null; createdAt: string; updatedAt: string;
}
// Smartphone, Tablet: identical shape to Computer

interface WindowsKey {
  id: number; licenseKey: string; computerId: number | null;
  status: string; notes: string | null; createdAt: string; updatedAt: string;
}

interface Antivirus {
  id: number; name: string; licenseKey: string;
  computerId: number | null; smartphoneId: number | null; tabletId: number | null;
  status: string; expiryDate: string | null;
  notes: string | null; createdAt: string; updatedAt: string;
}
// OtherSoftware: identical shape to Antivirus

interface User {
  id: number; name: string; surname: string | null;
  status: string; notes: string | null; createdAt: string; updatedAt: string;
}

interface AppConfig {
  dbPath: string; density: 'comfortable' | 'compact';
  darkMode: boolean; expiryWarningDays: number;
}

interface Alert {
  category: string; id: number; name: string;
  expiryDate: string; daysRemaining: number;
  severity: 'expired' | 'expiring';
}
```

### 5. Tailwind Configuration

Design tokens from `gover-ui-spec` map to `tailwind.config.ts` `theme.extend`. No external font imports — system stacks only. `tailwind.config.ts` is the single source of truth; no inline style fallbacks in Vue components.

```ts
// tailwind.config.ts
theme: {
  extend: {
    colors: {
      'sidebar':        '#0E1520',
      'sidebar-hover':  '#131E2E',
      'sidebar-active': '#182337',
      'sidebar-border': '#141D2B',
      'sidebar-text':   '#6B82A0',
      'sidebar-hi':     '#DDE8F5',
      'sidebar-accent': '#2563EB',
      'page':    '#EEF1F7',
      'surface': '#FFFFFF',
      'border':  '#DDE1EC',
      's-active':   '#2563EB',
      's-repair':   '#EA580C',
      's-spare':    '#7C3AED',
      's-retired':  '#9CA3AF',
      's-missing':  '#1E293B',
      's-expiring': '#D97706',
      's-expired':  '#DC2626',
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
    borderRadius: { badge: '2px' },
    width:  { sidebar: '200px' },
    height: { topbar: '48px' },
  }
}
```

Border radius: structural containers use `rounded-none` (Tailwind default 0). Only badges use `rounded-badge`.

### 6. Test Tooling

| Layer | Tool | Convention |
|---|---|---|
| Go unit tests | `testing` + `testify/assert` | Table-driven; files co-located (`database_test.go`) |
| Go coverage | `go test -cover ./...` | Target: 70% for `services`, 50% for `bridge` |
| Frontend unit | Vitest | Co-located with component (`DataTable.test.ts`) |
| Frontend component | Vue Test Utils + Vitest | Render → interact → assert |
| E2E | Playwright (deferred to Phase 4) | — |

Python tests with Go equivalents to write in Phase 2:
- `tests/test_alerts.py` → `internal/services` alert calculation tests
- `tests/test_models.py` → `internal/models` struct round-trip tests

### 7. Taskfile Commands (`gover:` namespace)

All 9 commands in root `Taskfile.yml`, each with `dir: gover2`:

| Task | Command | Tools required |
|---|---|---|
| `gover:install` | `go mod download` + `pnpm install --prefix frontend` | go, pnpm |
| `gover:dev` | `wails dev` | wails, node |
| `gover:build` | `wails build` | wails |
| `gover:test` | `go test ./...` + `pnpm --prefix frontend run test` | go, pnpm |
| `gover:lint` | `golangci-lint run ./...` + `pnpm --prefix frontend run lint` | golangci-lint, pnpm |
| `gover:typecheck` | `pnpm --prefix frontend run typecheck` | pnpm, vue-tsc |
| `gover:clean` | `rm -rf build/` + `pnpm --prefix frontend run clean` | — |
| `gover:package:deb` | `wails build -platform linux/amd64` | wails, dpkg tools |
| `gover:package:rpm` | `wails build -platform linux/amd64` | wails, rpm tools |

---

## Risks / Trade-offs

**`modernc.org/sqlite` performance:** Pure Go driver is ~10–20% slower than CGo at high query volumes. For a local IT inventory (hundreds of records), this is irrelevant.

**`gover2/` co-location:** Python and Go toolchains share the repo. Risk: accidental cross-contamination of build artifacts or CI steps. Mitigated by Taskfile namespace isolation and separate `go.mod`.

**Wails-generated bindings:** `wails dev` regenerates TypeScript bindings on Go changes. If `src/lib/api/` doesn't match the regenerated bindings, type errors surface immediately — this is the feature, not a bug.

**Tailwind purge scope:** `tailwind.config.ts` must include `./src/**/*.{vue,ts}` in `content` to avoid purging classes used only in dynamic strings. Strings like `bg-s-active` must appear in source, not be computed.
