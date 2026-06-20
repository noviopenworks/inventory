# IT Asset Inventory

A local-first desktop application for tracking IT assets — computers, smartphones,
tablets, software licences, Windows keys, antivirus subscriptions — and the users
they are assigned to. No server or internet connection required.

Built with:

- [Wails v2](https://wails.io/) — desktop shell and Go ↔ TypeScript bridge
- [Go 1.25](https://go.dev/) — backend services and SQLite access (`modernc.org/sqlite`, pure Go, no CGo for the DB)
- [Vue 3](https://vuejs.org/) + [TypeScript](https://www.typescriptlang.org/) — frontend
- [Tailwind CSS](https://tailwindcss.com/) — styling with CSS-variable theme tokens
- [Pinia](https://pinia.vuejs.org/) — state management
- SQLite — local data store

---

## Features

- **Asset tracking** across six categories: Computers, Smartphones, Tablets, Windows Keys, Antivirus, Other Software
- **User management** — assign devices to people, track user status
- **Unified "All" view** — every hardware asset in one table
- **Expiry/warranty alerts** — highlights items expiring soon
- **Inline editing** — add, edit, delete records through modal dialogs
- **CSV export** — export any view's data to a CSV file
- **Dark mode** — toggle in the UI; preference is persisted
- **Local SQLite database** — stored in your home directory, no cloud, no account

---

## Installation

Download the latest release for your platform from the [Releases](../../releases) page:

| Platform | File |
|----------|------|
| Debian/Ubuntu | `inventory_<version>_amd64.deb` |
| Red Hat/Fedora | `inventory-<version>-1.x86_64.rpm` |
| Other Linux | `inventory-<version>-linux-x86_64.tar.gz` |
| Windows | `inventory-<version>-windows-x86_64.exe` |

**Debian/Ubuntu:** `sudo dpkg -i inventory_<version>_amd64.deb`
**RPM:** `sudo rpm -i inventory-<version>-1.x86_64.rpm`
**tar.gz:** `tar xzf inventory-<version>-linux-x86_64.tar.gz && sudo ./inventory-<version>-linux-x86_64/install.sh`
**Windows:** run the `.exe` directly.

> Linux packages depend on `webkit2gtk-4.1` and `gtk3` at runtime (declared in the deb/rpm metadata).

---

## Running from source

### Requirements

- [Go](https://go.dev/doc/install) 1.25+
- [Node.js](https://nodejs.org/) 20+ with [pnpm](https://pnpm.io/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2
- [Task](https://taskfile.dev)
- Linux: `libwebkit2gtk-4.1-dev` and `libgtk-3-dev`

### Setup

```sh
git clone https://github.com/noviopenworks/inventory.git
cd inventory
task install       # go mod download + pnpm install
task dev           # launch the app with hot reload
```

---

## Development tasks

```sh
task                 # list all tasks
task dev             # run the app (hot reload)
task build           # build a release binary → build/bin/inventory
task test            # go test + vitest
task lint            # golangci-lint
task typecheck       # vue-tsc
task check           # lint + typecheck + test (CI-style)
task version:bump VERSION=1.2.3   # bump version everywhere
task build:linux     # build .deb + .tar.gz + .rpm (Linux only)
task build:windows   # build .exe (Windows only)
task clean           # remove build/bin, dist/, caches
```

---

## Project structure

```
main.go              # Wails entry point
app.go               # App factory (delegates to internal/bridge)
wails.json           # Wails project config
internal/
  backup/            # database backup and restore
  bridge/            # Wails-exposed methods
  config/            # app paths, DB location, preferences
  database/          # SQLite connection, pragmas, schema
  models/            # typed records
  services/          # list/create/update/delete, CSV export, alerts
frontend/            # Vue 3 + TypeScript + Tailwind UI
packaging/
  deb/               # Debian packaging templates
  rpm/               # RPM spec template
docs/                # architecture notes, decisions, development plan
```

---

## Documentation

- [`docs/DEVELOPMENT_PLAN.md`](docs/DEVELOPMENT_PLAN.md) — phased roadmap and status
- [`docs/ARCHITECTURE_NOTES.md`](docs/ARCHITECTURE_NOTES.md) — runtime boundaries and package layout
- [`docs/DECISIONS.md`](docs/DECISIONS.md) — canonical decision log
- [`openspec/specs/`](openspec/specs/) — published specifications

---

## License

GNU General Public License v3 or later — see [LICENSE](LICENSE).
