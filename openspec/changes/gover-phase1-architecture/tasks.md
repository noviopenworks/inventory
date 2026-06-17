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
