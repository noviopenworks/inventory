## 1. Go Package Architecture

- [x] 1.1 Document Go module path and `internal/` package layout
- [x] 1.2 Specify `internal/database` — connection open, WAL pragma, schema init (DDL carried from current schema)
- [x] 1.3 Specify `internal/models` — Go structs for all 7 tables with column-to-field mapping
- [x] 1.4 Specify `internal/services` — method signatures for list, get, create, update, delete per category
- [x] 1.5 Specify `internal/bridge` — Wails-exposed method names and signatures (the public API surface)
- [x] 1.6 Specify `internal/config` — app config struct, file location on Linux (`$XDG_CONFIG_HOME/inventory/`)
- [x] 1.7 Specify `internal/backup` — backup file path convention and trigger (pre-write, manual)

## 2. Vue + TypeScript App Structure

- [x] 2.1 Document `src/` directory layout (features/, layouts/, lib/, components/)
- [x] 2.2 Specify Vue Router routes (one route per sidebar category + All overview)
- [x] 2.3 Specify Pinia store shape (assets store, ui store)
- [x] 2.4 Specify `src/lib/api/` TypeScript wrapper shape — mirror of bridge methods
- [x] 2.5 Specify `src/lib/types/` — TypeScript interfaces for all 7 domain entities

## 3. Tailwind Configuration

- [x] 3.1 Map all design tokens from `gover-ui-spec` to `tailwind.config.ts` color palette entries
- [x] 3.2 Specify font family entries (UI stack, mono stack) in Tailwind config
- [x] 3.3 Specify border-radius scale (0 for structural, 2px for badges → Tailwind custom values)
- [x] 3.4 Specify any custom spacing/sizing tokens (sidebar-w: 200px, topbar-h: 48px)

## 4. Test Strategy

- [x] 4.1 Document Go test conventions (table-driven, testify/assert, coverage target)
- [x] 4.2 Document frontend test setup (Vitest config, Vue Test Utils, test file co-location)
- [x] 4.3 Specify which current Python tests have Go equivalents to write in Phase 2

## 5. Taskfile Commands

- [x] 5.1 Write the `gover:` namespace task definitions (all 9 commands from design.md)
- [x] 5.2 Verify each command is runnable on Linux (check tool availability requirements)
- [x] 5.3 Document `gover:` namespace command definitions in spec §6 (actual Taskfile.yml entries added in Phase 2 scaffolding)

## 6. Cross-Reference Update

- [x] 6.1 Annotate `gover/ARCHITECTURE_NOTES.md` Phase 1 section with resolved status and spec link
- [x] 6.2 Verify no application code in diff
