# Development Plan

This plan describes how to prepare and later build a new version of IT Asset Inventory using Wails, Go, Vue, TypeScript, and Tailwind CSS. It is a roadmap, not an implementation.

## Objectives

- Modernize the desktop experience while preserving the useful behavior of the current PyQt6 app.
- Keep the app local-first, with no server requirement and no internet dependency for normal use.
- Preserve or migrate existing SQLite data safely.
- Improve maintainability by separating backend domain logic, database access, frontend views, and desktop packaging concerns.
- Define quality gates before feature implementation begins.

## Non-Goals For This Planning Phase

- Do not scaffold a Wails project.
- Do not create Go, Vue, TypeScript, Tailwind, or package manager files.
- Do not modify the current Python application.
- Do not change the SQLite schema yet.
- Do not decide final UI layouts before requirements are reviewed.

## Phase 0: Discovery And Decisions

Purpose: decide what the new version must preserve, improve, or remove.

Work items:

- Review current user workflows: asset list, add, edit, delete, export, alerts, dark mode, and settings.
- Inventory the current database schema and query behavior.
- Decide whether the new version uses the existing schema directly or introduces a new schema plus migration.
- Define supported platforms and packaging targets.
- Define the minimum viable first release.
- Resolve the questions in `QUESTIONS.md`.

Exit criteria:

- Requirements are approved.
- Migration strategy is chosen.
- First release scope is defined.
- Quality gates are agreed.

## Phase 1: Technical Foundation Plan

Purpose: prepare the Wails/Vue/Go architecture before writing app features.

Planned decisions:

- Wails project layout and naming.
- Go package boundaries for database, domain services, app config, and desktop bridge methods.
- Vue application structure for routes, layout, components, stores, and API clients.
- Tailwind configuration and design token strategy.
- Testing tools for Go and frontend code.
- Taskfile commands for install, dev, test, lint, type-check, build, package, and clean.

Exit criteria:

- Architecture notes are accepted.
- Task command names are agreed.
- Test strategy is documented.

## Phase 2: Read-Only Prototype

Purpose: prove the stack can open the local database and display inventory data without changing records.

Future implementation scope:

- Wails app shell.
- Go SQLite connection and read queries.
- TypeScript API bindings.
- Vue layout with read-only tables.
- Basic error handling for missing or invalid database files.

Exit criteria:

- Existing inventory records can be viewed.
- No write operations are exposed.
- Build and test commands run locally.

## Phase 3: Core CRUD

Purpose: replace the current add, edit, and delete flows.

Future implementation scope:

- Create, update, and delete users and asset records.
- Form validation.
- Relationship selectors for assigned users and linked devices.
- Confirmation flows for destructive actions.
- Data refresh after changes.

Exit criteria:

- Current core CRUD behavior is covered.
- Database writes are tested.
- UI states for loading, empty, validation, and failure are handled.

## Phase 4: Alerts, Search, Export, And Settings

Purpose: restore supporting features and improve daily usability.

Future implementation scope:

- Expiry and warranty alerts.
- Search, filtering, sorting, and saved table preferences.
- CSV export.
- Dark mode and persisted appearance settings.
- Configurable warning window if approved.

Exit criteria:

- Users can find expiring assets quickly.
- CSV export matches agreed requirements.
- Appearance and table preferences persist.

## Phase 5: Packaging And Migration

Purpose: make the new app installable and safe for existing users.

Future implementation scope:

- Linux packages: `.deb`, `.rpm`, and `.tar.gz` if still required.
- Windows `.exe` build if still required.
- Version sync and release automation.
- Database backup before migrations.
- Import or migration path from the current Python app database.

Exit criteria:

- Packages can be built reproducibly.
- Upgrade path is documented.
- Existing data can be preserved or imported.

## Phase 6: Polish And Release Candidate

Purpose: make the modern version stable enough to replace the PyQt6 app.

Future implementation scope:

- UI polish and accessibility pass.
- Keyboard navigation review.
- Performance review with realistic data volumes.
- Documentation updates.
- Regression testing against current app behavior.

Exit criteria:

- Release candidate is approved.
- Known limitations are documented.
- Current app maintenance plan is defined.

## Suggested First Milestone

The first development milestone should be a read-only Wails/Vue shell that opens the existing SQLite database and displays the current asset categories. This reduces migration risk before building forms and write behavior.
