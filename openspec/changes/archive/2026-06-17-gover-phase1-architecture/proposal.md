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
