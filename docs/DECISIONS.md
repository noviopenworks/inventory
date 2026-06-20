# Decisions

This is the canonical decision log for gover. Each entry references the
OpenSpec change that locked the decision. Open questions still needing
answers are listed at the bottom.

## Resolved Decisions

### Architecture

| Decision | Resolution | Source |
|----------|------------|--------|
| Runtime boundaries | Go owns data + business logic; Vue owns presentation; TS wrappers mediate the Wails bridge | `openspec/changes/archive/2026-06-17-gover-phase1-architecture` |
| Go package layout | 6 packages: `database`, `models`, `services`, `bridge`, `config`, `backup`. `bridge` is the sole importer of all others | same |
| Frontend layout | Feature-folder structure; Pinia stores; hash-mode Vue Router with 9 routes | same |
| SQLite driver | `modernc.org/sqlite` (pure Go, no CGo) | same |

### Data

| Decision | Resolution | Source |
|----------|------------|--------|
| Schema strategy | Reuse the PyQt6 schema verbatim (7 tables). No migration required | `gover-phase1-architecture` |
| Database location | Reuses the PyQt6 default (`~/inventory.db`); user-configurable via File → Open Database | `gover2-parity` |

### UI / UX

| Decision | Resolution | Source |
|----------|------------|--------|
| Navigation model | Persistent left sidebar (200 px, themed), not tabs | `openspec/changes/archive/2026-06-17-gover-ui-spec` |
| First screen | Computers table (dashboard deferred) | same |
| Add / Edit UX | Slide-in right panel (300 px) | same |
| Density modes | Comfortable 33 px default, Compact 24 px toggle | same |
| Dark mode | Manual toggle, Light default; system-follow deferred | same |
| Keyboard | All PyQt6 shortcuts preserved plus panel shortcuts | same |

### Desktop integration

| Decision | Resolution | Source |
|----------|------------|--------|
| Wails bridge style | Generic asset methods where the schema allows; category-specific where the columns differ | `gover2-crud` |
| Backup / restore | Exposed via `internal/backup`; reachable from settings | `gover-phase1-architecture` |

### Tooling

| Decision | Resolution | Source |
|----------|------------|--------|
| Taskfile namespace | `gover:*` tasks in the root `Taskfile.yml` | root `Taskfile.yml` |
| Frontend package manager | pnpm | `gover/wails.json`, `gover/frontend/package.json` |
| Generated Wails TS bindings | Committed under `frontend/wailsjs/` | repo |
| Versioning | gover currently uses `wails.json` `info.productVersion`; root `.app-version` sync is part of Phase 5 | `gover/wails.json` |

## Open Questions (Phase 5+)

These were not blocking the rewrite but need answers before the first
packaged release:

- **Asset model evolution:** should categories remain fixed, or become configurable? Should software licenses be modeled separately from installed software? Should devices gain serial numbers, brands, locations, ownership, or tags?
- **Status options:** remain fixed or configurable?
- **Audit history / activity logs:** needed before release, or out of scope?
- **Import beyond CSV:** required for first release?
- **Target data volume:** what is the largest realistic inventory we must handle smoothly?
- **Platform priorities for Phase 5:** which of `.deb` / `.rpm` / `.tar.gz` / `.exe` ship first, and does GitHub Actions build all of them?
- **Old PyQt6 app:** continue receiving fixes during the rewrite, freeze, or archive on first gover release?
