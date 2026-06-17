# Future Task Backlog

This backlog is a planning template. Items should become actionable implementation tasks only after the relevant questions are answered.

## Planning Tasks

- [ ] Confirm the new version scope.
- [ ] Decide whether the app is a replacement or preview build.
- [ ] Decide platform priorities.
- [ ] Decide database compatibility and migration strategy.
- [ ] Decide first-release feature set.
- [ ] Decide UI navigation direction.
- [ ] Decide frontend package manager and test tooling.
- [ ] Decide release and packaging requirements.

## Foundation Tasks

- [ ] Create Wails project after approval.
- [ ] Add Vue, TypeScript, and Tailwind setup.
- [ ] Define Go package structure.
- [ ] Define TypeScript type strategy.
- [ ] Define Taskfile commands.
- [ ] Define lint, format, type-check, and test gates.
- [ ] Add initial CI workflow after local commands stabilize.

## Data Tasks

- [ ] Document the current SQLite schema.
- [ ] Decide schema reuse or migration.
- [ ] Add database backup plan.
- [ ] Add database connection layer.
- [ ] Add read queries for users and asset categories.
- [ ] Add migration tests if a new schema is approved.

## Frontend Tasks

- [ ] Define main layout.
- [ ] Define navigation model.
- [ ] Define table component requirements.
- [ ] Define form component requirements.
- [ ] Define theme tokens.
- [ ] Define empty, loading, validation, and error states.

## Feature Tasks

- [ ] Read-only asset list.
- [ ] Read-only users list.
- [ ] Add/edit/delete users.
- [ ] Add/edit/delete hardware assets.
- [ ] Add/edit/delete licenses and software.
- [ ] Expiry and warranty alerts.
- [ ] Search, filtering, and sorting.
- [ ] CSV export.
- [ ] Settings and preferences.

## Packaging Tasks

- [ ] Build Linux binary.
- [ ] Build Windows binary.
- [ ] Package `.deb` if required.
- [ ] Package `.rpm` if required.
- [ ] Package `.tar.gz` if required.
- [ ] Publish release artifacts through CI if required.

## Review Gates

- [ ] Planning approved.
- [ ] Architecture approved.
- [ ] Prototype approved.
- [ ] CRUD approved.
- [ ] Migration approved.
- [ ] Release candidate approved.
