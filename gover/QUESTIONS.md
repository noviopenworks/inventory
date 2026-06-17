# Questions Before Development

These questions should be answered before creating the Wails project or writing application code.

## Product Scope

- Should the new version be a full replacement for the PyQt6 app or a parallel preview version first?
- What is the minimum acceptable first release?
- Which current features are mandatory on day one?
- Which current features can be delayed?
- Should any current feature be removed or redesigned?

## Users And Workflows

- Who uses the app: one IT admin, multiple admins on the same machine, or non-technical staff?
- Is the app still strictly local-only?
- Is multi-user collaboration out of scope?
- Are audit history, deleted-item recovery, or activity logs needed?

## Data And Migration

- Should the new app use the current SQLite schema directly?
- Should a new schema be introduced with a migration step?
- Should the app automatically discover the existing database location?
- Should the app create a backup before any migration?
- Should import/export support be expanded beyond CSV?
- What data volume should the app comfortably handle?

## Asset Model

- Should asset categories remain fixed: computers, smartphones, tablets, Windows keys, antivirus, other software, users?
- Should categories become configurable?
- Should software licenses be modeled separately from installed software?
- Should devices have serial numbers, brands, locations, ownership, or tags?
- Should status options remain fixed or become configurable?

## UI And UX

- Should the current tab model be preserved?
- Should the new UI use a sidebar, dashboard, command palette, or another navigation model?
- Should the first screen be an overview dashboard or the asset table?
- Should add/edit use modal dialogs, side panels, or full pages?
- How important is keyboard-first operation?
- Should the UI support compact and comfortable density modes?
- Should dark mode follow the system theme or remain a manual toggle?

## Backend And Desktop Integration

- Which Go SQLite driver should be used?
- Where should app configuration live on Linux and Windows?
- Should database path selection be user-configurable?
- Should the app expose a backup and restore workflow?
- Should the Wails bridge expose category-specific methods or generic asset methods?

## Testing And Quality

- What minimum Go test coverage is expected?
- What frontend test tools should be used?
- Should end-to-end tests be required before release?
- Which current Python tests should be translated into Go or frontend tests?
- What manual QA checklist is required before packaging?

## Packaging And Release

- Which platforms are required first: Linux, Windows, or both?
- Should packaging mirror the current `.deb`, `.rpm`, `.tar.gz`, and `.exe` outputs?
- Should GitHub Actions build all packages?
- Should versioning continue to use `.app-version`?
- Should the old Python app continue receiving fixes during the rewrite?

## Internal Tooling

- Should the rewrite use the root `Taskfile.yml` or a separate `gover/Taskfile.yml` later?
- Which package manager should be used for frontend dependencies?
- Should generated Wails bindings be committed?
- Should agents be allowed to scaffold the project automatically after approval?
