# Architecture Notes

These notes describe areas to investigate. They do not lock final implementation decisions.

## Proposed Direction

Use Wails as the desktop runtime, Go as the backend application layer, Vue and TypeScript as the frontend, Tailwind CSS for styling, and SQLite for local persistence.

## Candidate Runtime Boundaries

> **Resolved (2026-06-17, change `gover-phase1-architecture`):** Go owns all data access and business logic via a 6-package `internal/` layout; Vue owns presentation; TypeScript API wrappers in `src/lib/api/` mediate the Wails bridge.
> Full decisions: [`openspec/changes/gover-phase1-architecture/specs/architecture/spec.md`](../openspec/changes/gover-phase1-architecture/specs/architecture/spec.md)

- Go owns database access, migrations, validation that protects data integrity, file system paths, backup, restore, import, export, and platform-specific behavior.
- Vue owns layout, navigation, presentation state, forms, table interactions, filters, and user-facing feedback.
- TypeScript API clients wrap Wails-generated bindings so frontend components do not depend directly on low-level bridge details.

## Candidate Go Areas

> **Resolved (2026-06-17, change `gover-phase1-architecture`):** Six packages confirmed: `internal/database`, `internal/models`, `internal/services`, `internal/bridge`, `internal/config`, `internal/backup`. Dependency order locked; `bridge` is the sole package that imports from all others.
> Full decisions: [`openspec/changes/gover-phase1-architecture/specs/architecture/spec.md`](../openspec/changes/gover-phase1-architecture/specs/architecture/spec.md)

- `database`: connection management, SQLite pragmas, transactions, migrations.
- `models`: typed records for users, devices, licenses, software, alerts, and settings.
- `services`: application operations such as list assets, create record, update record, delete record, export CSV, check alerts.
- `bridge`: Wails-exposed methods called by the frontend.
- `config`: app paths, database location, preferences, and version metadata.
- `backup`: database backup and restore utilities if approved.

## Candidate Frontend Areas

> **Resolved (2026-06-17, change `gover-phase1-architecture`):** Feature-folder structure adopted: `features/` (assets, licenses, users, alerts stub), `components/` (10 shared components), `lib/api/index.ts`, `lib/types/index.ts`, `stores/assets.ts`, `stores/ui.ts`. Hash-mode Vue Router with 9 routes.
> Full decisions: [`openspec/changes/gover-phase1-architecture/specs/architecture/spec.md`](../openspec/changes/gover-phase1-architecture/specs/architecture/spec.md)

- `app`: Vue app setup and global providers.
- `layouts`: main shell, navigation, top bar, and responsive structure.
- `features/assets`: asset tables, forms, filters, and category views.
- `features/users`: user list and user forms.
- `features/alerts`: expiry and warranty alert views.
- `features/settings`: appearance, database path, export, backup, and restore.
- `lib/api`: TypeScript wrappers around Wails bindings.
- `lib/types`: frontend-facing TypeScript types.
- `components`: reusable buttons, fields, modals, panels, tables, badges, and empty states.

## Data Strategy Options To Evaluate

> **Resolved (2026-06-17, change `gover-phase1-architecture`):** Option 1 adopted — reuse current schema directly. Go DDL in `internal/database` carries the 7-table schema verbatim from `db/schema.py`. SQLite driver: `modernc.org/sqlite` (pure Go, no CGo).
> Full decisions: [`openspec/changes/gover-phase1-architecture/specs/architecture/spec.md`](../openspec/changes/gover-phase1-architecture/specs/architecture/spec.md)

1. Reuse current schema directly.
   Pros: fastest path and less migration risk.
   Cons: carries current modeling limitations forward.

2. Create a new schema with migration.
   Pros: better long-term model and more room for richer features.
   Cons: requires careful migration, backup, and compatibility testing.

3. Hybrid approach.
   Pros: start compatible, migrate only when needed.
   Cons: can become confusing if compatibility rules are not strict.

## UI Strategy Options To Evaluate

> **Resolved (2026-06-17, change `gover-ui-spec`):** Option 2 (Dashboard plus sidebar) adopted with modifications.
> Full decisions: [`openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md`](../openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md)
> Visual reference: [`gover/UI_DESIGN_REFERENCE.html`](UI_DESIGN_REFERENCE.html)
>
> Key decisions: persistent left sidebar (200 px, `#0E1520`), Computers as first screen, slide-in right panel for add/edit (300 px), Comfortable/Compact density toggle, manual dark mode toggle (Light default), 7 fixed categories, native OS menu bar via `wails/v2/pkg/menu`.

1. Familiar table-first layout.
   Pros: easy transition from current app.
   Cons: less modern if not improved carefully.

2. Dashboard plus sidebar layout.
   Pros: more modern and better overview of alerts and counts.
   Cons: more design work and more decisions upfront.

3. Command-center layout.
   Pros: efficient for power users and search-heavy workflows.
   Cons: may be overbuilt for a simple local asset tracker.

## Risk Areas

- Data migration from the current SQLite database.
- Recreating current table and form behavior without losing small usability details.
- Keeping Wails bindings clean instead of exposing database details directly to Vue.
- Avoiding a frontend that looks modern but makes bulk data work slower.
- Packaging differences between PyInstaller and Wails.

## Recommended First Architecture Validation

Build a future read-only proof of concept that opens the existing database and displays users and hardware assets. Do not start write paths until data access, app shell, table rendering, and packaging assumptions are validated.
