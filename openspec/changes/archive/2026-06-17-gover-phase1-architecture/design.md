## Approach

Produce a single capability spec (`specs/architecture/spec.md`) that documents every technical decision a developer needs to scaffold Phase 2. The spec is written in Q&A format, mirroring the UI decisions spec from `gover-ui-spec`. It is the authoritative reference for package names, bridge method signatures, Taskfile command names, and Tailwind config structure.

A second small update annotates `gover/ARCHITECTURE_NOTES.md` to mark the architecture section as resolved, parallel to how `gover/QUESTIONS.md` was annotated in `gover-ui-spec`.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      Wails Desktop App                      │
├────────────────────────┬────────────────────────────────────┤
│     Go Backend         │        Vue 3 Frontend              │
│  ─────────────────     │  ──────────────────────────────    │
│  internal/database     │  src/layouts/AppShell.vue          │
│  internal/models       │  src/features/assets/              │
│  internal/services     │  src/features/users/               │
│  internal/bridge       │  src/features/alerts/              │
│  internal/config       │  src/lib/api/          (TS)        │
│  internal/backup       │  src/lib/types/        (TS)        │
│                        │  src/components/                   │
└────────────────────────┴────────────────────────────────────┘
         ▲                          ▲
         │  Wails generated         │  pnpm + Vite
         │  bindings                │
         └──────────── bridge ──────┘
```

## Go Package Design

- Module path: `github.com/user/inventory` (or local `inventory` module — decided in spec)
- All internal packages under `internal/` (not importable by external tools)
- Package names match ARCHITECTURE_NOTES.md candidate areas: `database`, `models`, `services`, `bridge`, `config`, `backup`
- `database` owns schema open/init and carries forward the current 7-table DDL verbatim
- `models` defines typed Go structs matching every current SQLite column
- `services` implements all operations; `bridge` exposes only what Vue needs

## Vue App Design

- Vite + Vue 3 Composition API + TypeScript strict mode
- Feature-folder structure under `src/features/` (not component-type folders)
- `src/lib/api/` wraps Wails-generated bindings — Vue components never call `window.go.*` directly
- Pinia for state management
- Vue Router for sidebar navigation (hash mode, no server needed)

## Tailwind Config Strategy

- Design tokens from `gover-ui-spec` map directly to `tailwind.config.ts` `theme.extend`
- No `@import` of external fonts — system font stacks only
- All CSS custom properties in the spec become Tailwind color/spacing utilities
- `tailwind.config.ts` is the single source of truth; no inline style fallbacks in components

## Test Tooling

- Go: standard `testing` package + `testify/assert` for assertions; table-driven tests
- Frontend unit tests: Vitest (Vite-native, no separate webpack config)
- Frontend component tests: Vue Test Utils + Vitest
- No E2E tests in Phase 1 or 2; E2E deferred to Phase 4 (Playwright candidate)

## Taskfile Commands (`gover:` namespace)

All commands live in root `Taskfile.yml` under `gover:` namespace:

| Command | Purpose |
|---|---|
| `gover:install` | Install Go deps + pnpm install |
| `gover:dev` | `wails dev` (hot-reload) |
| `gover:build` | `wails build` (production binary) |
| `gover:test` | Go tests + frontend tests |
| `gover:lint` | `golangci-lint` + `eslint` |
| `gover:typecheck` | `vue-tsc --noEmit` |
| `gover:clean` | Remove build artifacts |
| `gover:package:deb` | Build `.deb` package (Linux) |
| `gover:package:rpm` | Build `.rpm` package (Linux) |
