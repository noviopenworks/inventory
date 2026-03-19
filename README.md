# IT Asset Inventory

A desktop GUI application for tracking IT assets — computers, smartphones, tablets, software licences, Windows keys, antivirus subscriptions — and the users they are assigned to. Built with PyQt6 and a local SQLite database; no server or internet connection required.

> **Note:** This project was mostly vibe-coded with GitHub Copilot as an AI pair programmer. It works, it has tests, and the code is reasonably clean — but don't expect every architectural decision to hold up to deep scrutiny.

---

## Features

- **Asset tracking** across six categories: Computers, Smartphones, Tablets, Windows Keys, Antivirus, Other Software
- **User management** — assign devices to people, track user status
- **Unified "All" view** — see every hardware asset in one table
- **Expiry/warranty alerts** — highlights items expiring within 30 days on startup
- **Inline editing** — add, edit, delete records through modal dialogs
- **CSV export** — export any tab's data to a CSV file
- **Dark mode** — toggle via the View menu; preference is persisted
- **Local SQLite database** — stored in your home directory, no cloud, no account

---

## Screenshots

_TODO_

---

## Installation

### Pre-built binaries (recommended)

Download the latest release for your platform from the [Releases](../../releases) page:

| Platform | File |
|----------|------|
| Debian/Ubuntu | `inventory_<version>_amd64.deb` |
| Red Hat/Fedora | `inventory-<version>-1.x86_64.rpm` |
| Other Linux | `inventory-<version>-linux-x86_64.tar.gz` |
| Windows | `inventory-<version>-windows-x86_64.exe` |

**Debian/Ubuntu:**
```sh
sudo dpkg -i inventory_<version>_amd64.deb
```

**RPM:**
```sh
sudo rpm -i inventory-<version>-1.x86_64.rpm
```

**tar.gz:**
```sh
tar xzf inventory-<version>-linux-x86_64.tar.gz
sudo ./inventory-<version>-linux-x86_64/install.sh
```

**Windows:** Run the `.exe` directly — no installer needed.

---

## Running from source

### Requirements

- Python 3.14+
- [uv](https://docs.astral.sh/uv/)
- [Task](https://taskfile.dev)

### Setup

```sh
git clone https://github.com/noviopenworks/inventory.git
cd inventory
task install       # uv sync --all-groups
task run           # launch the app
```

---

## Development tasks

```sh
task               # list all tasks
task run           # run the app
task test          # run the test suite
task test:cov      # tests with coverage report (target ≥ 85 %)
task lint          # ruff check
task lint:fix      # ruff check --fix
task fmt           # ruff format
task type:check    # mypy
task check         # lint + type-check + tests (CI-style)
task version:bump VERSION=1.2.3   # bump version everywhere
task build:linux   # build .deb + .tar.gz + .rpm (Linux only)
task build:windows # build .exe (Windows only)
task clean         # remove dist/, build/ and caches
```

---

## Project structure

```
main.py              # entry point
app/
  main_window.py     # main window, menus, tab bar
  dialogs.py         # add/edit dialogs
  models.py          # Qt table model
  style.py           # theme and stylesheet helpers
db/
  schema.py          # SQLite DDL and migrations
  queries.py         # read queries (All view, expiry alerts)
  connection.py      # connection factory
  config.py          # per-table column metadata and constants
  alerts.py          # expiry/warranty alert queries
tests/               # pytest test suite
packaging/
  deb/               # Debian packaging templates
  rpm/               # RPM spec template
```

---

## Building packages

Packages are built with PyInstaller and wrapped by platform-specific tools.
CI runs `task build:linux` (deb + rpm + tar.gz) and `task build:windows` (exe) in parallel and publishes everything to the GitHub release.

---

## License

GNU General Public License v3 or later — see [LICENSE](LICENSE).
