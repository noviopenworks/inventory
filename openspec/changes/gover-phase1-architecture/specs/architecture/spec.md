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

### Phase 3 — CRUD naming pattern (implement in Phase 3)

```
Get<Category>(id int) (models.<Category>, error)
Create<Category>(input models.Create<Category>Input) (models.<Category>, error)
Update<Category>(id int, input models.Update<Category>Input) (models.<Category>, error)
Delete<Category>(id int) error
```

### Phase 4 — Alerts (implement in Phase 4)

```
GetAlerts() ([]models.Alert, error)
```

### Requirement: bridge-no-sql

`internal/bridge` methods must contain no SQL statements. All data access is delegated to `internal/services`.

### Requirement: api-wrapper

Vue components must not call `window.go.*` directly. All Wails bridge calls are made through `src/lib/api/index.ts`.

#### Scenario: component calls bridge directly

- Given a grep of `frontend/src/features/` and `frontend/src/components/` for `window.go.`
- When any `.vue` or `.ts` file contains a direct `window.go.*` call
- Then it is a violation of the api-wrapper requirement

---

## Section 3: Vue App Structure & TypeScript Types

Vite + Vue 3 Composition API + TypeScript strict mode. Feature-folder structure. Vue Router hash mode (no server needed in Wails WebView). Pinia for state management.

### `frontend/src/` directory layout

```
src/
├── main.ts
├── App.vue
├── layouts/
│   └── AppShell.vue
├── features/
│   ├── assets/         ← computers, smartphones, tablets, all-view
│   ├── licenses/       ← windows keys, antivirus, other software
│   ├── users/
│   └── alerts/         ← Phase 4 stub
├── components/         ← 10 shared components from gover-ui-spec
│   ├── Sidebar.vue
│   ├── SidebarItem.vue
│   ├── Topbar.vue
│   ├── DataTable.vue
│   ├── StatusBadge.vue
│   ├── EditPanel.vue
│   ├── FormField.vue
│   ├── AlertsModal.vue
│   ├── Statusbar.vue
│   └── AppShell.vue
├── lib/
│   ├── api/
│   │   └── index.ts    ← all bridge wrappers, typed
│   └── types/
│       └── index.ts    ← all domain interfaces
└── stores/
    ├── assets.ts       ← current list, selected row, filters
    └── ui.ts           ← density, dark mode, sidebar collapse state
```

### Vue Router routes (hash mode)

| Hash path | View | Sidebar label |
|---|---|---|
| `/` | redirect → `/computers` | — |
| `/all` | All Assets (hardware: Computer + Smartphone + Tablet) | All Assets |
| `/computers` | Computers | Computers |
| `/smartphones` | Smartphones | Smartphones |
| `/tablets` | Tablets | Tablets |
| `/windows-keys` | Windows Keys | Windows Keys |
| `/antivirus` | Antivirus | Antivirus |
| `/other-software` | Other Software | Other Software |
| `/users` | Users | Users |

### TypeScript interfaces (`lib/types/index.ts`)

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

### Pinia store shapes

**`stores/assets.ts`**: `currentCategory`, `items` (typed union), `selectedId`, `filters` (search string, status filter).

**`stores/ui.ts`**: `density` (`'comfortable' | 'compact'`), `darkMode` (`boolean`), `sidebarCollapsed` (`boolean`).

### Requirement: typescript-strict

`frontend/tsconfig.json` must include `"strict": true`. All domain types must be defined in `lib/types/index.ts`; no inline `any` types in component or API files.

---

## Section 4: Tailwind Configuration

Design tokens from `gover-ui-spec` map to `tailwind.config.ts` `theme.extend`. No external font imports — system font stacks only. `tailwind.config.ts` is the single source of truth; no inline style fallbacks in Vue components.

```ts
// tailwind.config.ts (theme.extend excerpt)
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
},
fontFamily: {
  ui:   ['-apple-system', 'BlinkMacSystemFont', '"Segoe UI"', 'system-ui', 'sans-serif'],
  mono: ['"Cascadia Code"', '"Fira Code"', 'Consolas', '"Courier New"', 'monospace'],
},
borderRadius: { badge: '2px' },
width:  { sidebar: '200px' },
height: { topbar: '48px' },
```

Structural containers use `rounded-none`. Only `StatusBadge` uses `rounded-badge`. The `content` array must include `./src/**/*.{vue,ts}` so dynamic class strings like `bg-s-active` survive purging.

### Requirement: tailwind-tokens

All colors, font families, border-radius, width, and height values used in Vue components must be derived from `tailwind.config.ts` `theme.extend` entries. No `style=""` inline color values in component templates.

---

## Section 5: Test Tooling

### Go tests

- Framework: `testing` + `testify/assert`
- Style: table-driven; test files co-located with source (`database_test.go`)
- Coverage target: 70% for `internal/services`, 50% for `internal/bridge`
- Command: `go test -cover ./...`

### Frontend tests

- Framework: Vitest (Vite-native)
- Component tests: Vue Test Utils + Vitest
- Test files: co-located with component (`DataTable.test.ts` beside `DataTable.vue`)
- E2E: deferred to Phase 4 (Playwright candidate)

### Python equivalents

| Python test file | Go equivalent location |
|---|---|
| `tests/test_alerts.py` | `internal/services` — alert calculation tests |
| `tests/test_models.py` | `internal/models` — struct round-trip tests |

### Requirement: go-test-coverage

Phase 2 must include tests for all `internal/services` list operations reaching ≥70% line coverage in that package, verified by `go test -cover ./internal/services/`.

---

## Section 6: Taskfile Commands

All 9 `gover:` commands live in the root `Taskfile.yml` with `dir: gover2` on each task.

| Task key | Command(s) | Required tools |
|---|---|---|
| `gover:install` | `go mod download` + `pnpm install --prefix frontend` | go, pnpm |
| `gover:dev` | `wails dev` | wails, node |
| `gover:build` | `wails build` | wails |
| `gover:test` | `go test ./...` + `pnpm --prefix frontend run test` | go, pnpm |
| `gover:lint` | `golangci-lint run ./...` + `pnpm --prefix frontend run lint` | golangci-lint, pnpm |
| `gover:typecheck` | `pnpm --prefix frontend run typecheck` | pnpm, vue-tsc |
| `gover:clean` | `rm -rf build/` + `pnpm --prefix frontend run clean` | — |
| `gover:package:deb` | `wails build -platform linux/amd64` + deb packaging step | wails, dpkg tools |
| `gover:package:rpm` | `wails build -platform linux/amd64` + rpm packaging step | wails, rpm tools |

### Requirement: taskfile-namespace

All Gover build operations must be invokable from the repo root via `task gover:<subcommand>`. No `gover:` task may have a side effect on the Python app directory.
