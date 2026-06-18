# Gover — IT Asset Inventory (Wails rewrite)

Gover is the in-development successor to the [PyQt6 IT Asset Inventory](../README.md)
app in this repository. It is a local-first desktop application for tracking
computers, smartphones, tablets, Windows keys, antivirus subscriptions, and
other software licenses, plus the users they are assigned to.

Built with:

- [Wails v2](https://wails.io/) — desktop shell and Go ↔ TypeScript bridge
- [Go 1.25](https://go.dev/) — backend services and SQLite access (`modernc.org/sqlite`, pure Go — no CGo)
- [Vue 3](https://vuejs.org/) + [TypeScript](https://www.typescriptlang.org/) — frontend
- [Tailwind CSS](https://tailwindcss.com/) — styling, with CSS-variable theme tokens
- [Pinia](https://pinia.vuejs.org/) — state management
- SQLite — local data store; the schema is reused verbatim from the PyQt6 app, so the same `.db` file works in both

> **Status:** Phases 0–4 of the [development plan](docs/DEVELOPMENT_PLAN.md) are complete (discovery, foundation, read-only prototype, CRUD, and supporting features). Packaging (Phase 5) and polish (Phase 6) are pending, so there are no packaged releases yet — run from source.

## Project layout

```
gover/
├── main.go              # Wails entry point
├── app.go               # App factory (delegates to internal/bridge)
├── go.mod               # module gover
├── wails.json           # Wails project config
├── internal/
│   ├── backup/          # database backup and restore
│   ├── bridge/          # Wails-exposed methods (sole importer of all other internal pkgs)
│   ├── config/          # app paths, DB location, preferences
│   ├── database/        # SQLite connection, pragmas, schema
│   ├── models/          # typed records for users, devices, licenses, alerts, settings
│   └── services/        # list/create/update/delete, CSV export, alert queries
├── frontend/
│   ├── src/
│   │   ├── features/    # assets, users, alerts, settings views
│   │   ├── components/  # shared buttons, fields, modals, tables
│   │   ├── layouts/     # AppShell, sidebar, topbar
│   │   ├── lib/api/     # TypeScript wrappers around Wails bindings
│   │   ├── lib/types/   # frontend-facing TS types
│   │   ├── stores/      # Pinia stores (assets, ui)
│   │   └── router/      # hash-mode Vue Router (9 routes)
│   ├── wailsjs/         # generated Wails bindings (committed)
│   └── package.json
└── docs/
    ├── DEVELOPMENT_PLAN.md       # phased roadmap + status
    ├── ARCHITECTURE_NOTES.md     # resolved architecture decisions
    ├── DECISIONS.md              # canonical decision log
    └── UI_DESIGN_REFERENCE.html  # visual reference for the UI
```

## Requirements

- [Go](https://go.dev/doc/install) 1.25+
- [Node.js](https://nodejs.org/) 20+ with [pnpm](https://pnpm.io/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2
- Linux: `webkit2gtk-4.1` deps (Debian: `libwebkit2gtk-4.1-dev gtk3-devel`; Fedora: `webkit2gtk4.1-devel`)

## Setup & run

All commands go through the root [`Taskfile.yml`](../Taskfile.yml) under the
`gover:*` namespace:

```sh
task gover:install     # go mod download + pnpm install
task gover:dev         # wails dev — hot reload, opens the app window
task gover:build       # wails build — produces a release binary in build/bin/
```

## Quality gates

```sh
task gover:test        # go test ./... + pnpm run test (Vitest)
task gover:lint        # golangci-lint + pnpm run lint
task gover:typecheck   # vue-tsc --noEmit
task gover:clean       # remove build/ and frontend/dist
```

## Data compatibility with the PyQt6 app

Gover reuses the PyQt6 SQLite schema verbatim. To point gover at an existing
inventory, use **File → Open Database…** and select the `.db` file written by
the PyQt6 app (default location: `~/inventory.db`). No migration is required.

## Documentation

- [`docs/DEVELOPMENT_PLAN.md`](docs/DEVELOPMENT_PLAN.md) — phased roadmap and current status
- [`docs/ARCHITECTURE_NOTES.md`](docs/ARCHITECTURE_NOTES.md) — runtime boundaries, package layout, data and UI strategy
- [`docs/DECISIONS.md`](docs/DECISIONS.md) — canonical decision log with links to source OpenSpec changes
- [`docs/UI_DESIGN_REFERENCE.html`](docs/UI_DESIGN_REFERENCE.html) — static HTML visual reference (open in a browser)
- Root [`openspec/specs/`](../openspec/specs/) — published specifications (alerts, search, export, preferences, database management, repository hygiene)
- Root [`openspec/changes/archive/`](../openspec/changes/archive/) — frozen history of how gover was built change-by-change

## License

GNU General Public License v3 or later — see root [`LICENSE`](../LICENSE).
