# Go-Default Cutover Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the Wails/Go app the canonical `inventory` product at the repo root, delete the PyQt6 Python app, and rework Taskfile/CI/release/packaging for the Go stack — shipping as `v1.0.0`.

**Architecture:** Six ordered phases on one feature branch, each ending green: (1) promote `gover/` → root + rename module, (2) Taskfile + delete Python, (3) packaging pipeline, (4) CI, (5) release, (6) docs + cleanup. The Go app and Python app can coexist at the root between Tasks 1–2 (no filename collisions: `main.go`≠`main.py`, `internal/`≠`app/`), so nothing breaks mid-flight.

**Tech Stack:** Go 1.25, Wails v2.12, Vue 3 + TypeScript, Vite, pnpm, Tailwind, SQLite (`modernc.org/sqlite`), Task (go-task), GitHub Actions.

## Global Constraints

- Product/binary name: **`inventory`** (was `gover`). Go module path: **`inventory`** (was `gover`).
- Go module: `module inventory`, `go 1.25.0`. Build tag: **`webkit2_41`** (matches `wails.json` `build:tags`). CGO required for any package importing Wails (e.g. `internal/bridge`).
- Starting version: **`1.0.0`** in `.app-version` (source of truth), mirrored into `wails.json` `info.productVersion` and `frontend/package.json` `version`.
- Release targets: Linux `.deb` + `.rpm` + `.tar.gz`, Windows `.exe`. **No macOS. No NSIS.**
- Wails Linux binaries are **not** self-contained: deb must `Depends: libwebkit2gtk-4.1-0, libgtk-3-0`; rpm must `Requires: webkit2gtk4.1, gtk3`.
- After cutover: **zero** Python source, manifests, or caches remain anywhere in the tree.
- Frontend has no ESLint — `lint` is Go-only (`golangci-lint`); the frontend is gated by `vue-tsc` typecheck.

---

### Task 1: Promote the Go app to the repo root and rename the module to `inventory`

**Files:**
- Move (filesystem): `gover/{main.go, app.go, go.mod, go.sum, wails.json, internal/, frontend/, build/}` → repo root
- Move (git-tracked): `gover/docs/{DEVELOPMENT_PLAN.md, ARCHITECTURE_NOTES.md, DECISIONS.md, UI_DESIGN_REFERENCE.html}` → `docs/`
- Stage for later merge: `gover/README.md` → `README.gover.md` (temp; consumed in Task 6)
- Modify: `go.mod` (module path), all `*.go` (import paths)
- Replace: `.gitignore`
- Delete: `gover/.gitignore`, `gover/gover2` (stray 5.2M binary), `gover/` (empty dir), `build/bin/*` (stray 12.2M `gover2`)

**Interfaces:**
- Produces: a buildable Go module `inventory` rooted at the repo top — `main.go`, `app.go`, `internal/`, `frontend/`, `wails.json`, `go.mod` all at root. `wails build` succeeds. Later tasks assume these root paths.
- Note: Python files (`main.py`, `app/`, `db/`, etc.) remain untouched and still present after this task — they are deleted in Task 2.

- [ ] **Step 1: Create the feature branch**

```bash
git checkout -b feat/go-default-cutover
```

- [ ] **Step 2: Move the Go source, config, and assets to the root**

`gover/build/` is currently *untracked* (hidden by the old `.gitignore` `build/` rule), so use plain `mv`, not `git mv`. `git add -A` later detects renames.

```bash
rm -rf build                       # clear any transient PyInstaller build dir (gitignored)
mv gover/main.go gover/app.go gover/go.mod gover/go.sum gover/wails.json .
mv gover/internal internal
mv gover/frontend frontend
mv gover/build build
rm -rf build/bin                   # drop the stray 12.2M gover2 binary; wails regenerates build/bin
```

- [ ] **Step 3: Move the architecture docs and stage the gover README**

```bash
mv gover/docs/DEVELOPMENT_PLAN.md gover/docs/ARCHITECTURE_NOTES.md gover/docs/DECISIONS.md gover/docs/UI_DESIGN_REFERENCE.html docs/
mv gover/README.md README.gover.md
```

- [ ] **Step 4: Delete the leftover gover directory and stray binary**

```bash
rm -f gover/gover2 gover/.gitignore
rmdir gover/docs gover 2>/dev/null || true
ls gover 2>/dev/null && echo "STILL NOT EMPTY — investigate" || echo "gover/ removed"
```

Expected: `gover/ removed`.

- [ ] **Step 5: Rename the Go module and rewrite import paths**

```bash
sed -i 's|^module gover$|module inventory|' go.mod
grep -rl '"gover/internal' --include='*.go' . | xargs --no-run-if-empty sed -i 's|"gover/internal|"inventory/internal|g'
gofmt -w .
head -1 go.mod
grep -rn '"gover/' --include='*.go' . || echo "no stray gover/ imports"
```

Expected: `module inventory`, then `no stray gover/ imports`.

- [ ] **Step 6: Remove the dead `replace` line in `go.mod`**

Delete the trailing commented line `// replace github.com/wailsapp/wails/v2 v2.12.0 => /home/mg/go/pkg/mod` if present. Then:

```bash
go mod tidy
```

- [ ] **Step 7: Update `wails.json` name + version fields**

Set `outputfilename` to `inventory`, `info.productName` to `Inventory`, `info.productVersion` to `1.0.0`:

```json
{
  "name": "inventory",
  "outputfilename": "inventory",
  "frontend:install": "pnpm install",
  "frontend:build": "pnpm run build",
  "frontend:dev:watcher": "pnpm run dev",
  "frontend:dev:serverUrl": "auto",
  "frontend:dir": "frontend",
  "info": {
    "productName": "Inventory",
    "productVersion": "1.0.0",
    "companyName": "",
    "copyright": "",
    "comments": ""
  },
  "bindings": {
    "ts_generation_dir": "frontend/wailsjs"
  },
  "build:tags": "webkit2_41"
}
```

- [ ] **Step 8: Set the version source of truth**

```bash
printf '1.0.0\n' > .app-version
```

- [ ] **Step 9: Replace `.gitignore` so `build/` assets are trackable but `build/bin/` is not**

```gitignore
# ── Build output ──
build/bin/
dist/
/inventory
/inventory.exe

# ── Node / frontend ──
node_modules/
frontend/node_modules/
frontend/dist/

# ── Go ──
*.test
*.out

# ── Task (go-task) ──
.task/

# ── Packaging artefacts ──
*.deb
*.rpm
*.tar.gz
*.tar.bz2
*.tar.xz
*.zip

# ── SQLite databases (runtime data, not source) ──
*.db
*.db-shm
*.db-wal
*.sqlite
*.sqlite3

# ── Editors ──
.vscode/launch.json
.vscode/tasks.json
.idea/
*.iml
*.swp
*.swo
*~

# ── OS ──
.DS_Store
._*
Thumbs.db
desktop.ini

# ── Secrets & local config ──
.env
.env.*
!.env.example
*.pem
*.key

# ── Logs ──
*.log
logs/
```

- [ ] **Step 10: Build the app to verify the move + rename (authoritative gate)**

Run: `wails build -tags webkit2_41`
Expected: completes, prints `Built '.../build/bin/inventory'` (no compile/import errors).

If `wails` is unavailable in the environment, fall back to: `CGO_ENABLED=1 go build -tags webkit2_41 ./...` (requires `frontend/dist` to exist on disk for the `//go:embed` in `main.go`; run `pnpm --prefix frontend run build` first if needed).

- [ ] **Step 11: Stage everything (including the now-trackable build assets) and commit**

```bash
git add -A
git status --short | grep -E '^A.*build/(appicon|windows|darwin)' && echo "build assets tracked" || echo "check build asset tracking"
git commit -m "refactor(cutover): promote gover to repo root, rename module to inventory"
```

---

### Task 2: Rework the Taskfile for the Go toolchain and delete the Python app

**Files:**
- Replace: `Taskfile.yml`
- Create: `.golangci.yml`
- Delete: `main.py`, `app/`, `db/`, `tests/`, `pyproject.toml`, `uv.lock`, `inventory.spec`, `.mypy_cache/`, `.ruff_cache/`, `.pytest_cache/`, `.venv/`

**Interfaces:**
- Consumes: root-level Go module `inventory` and `.app-version` (`1.0.0`) from Task 1.
- Produces: top-level tasks `install`, `dev`, `build`, `lint`, `typecheck`, `test`, `check`, `version:bump`, `version:sync`, plus packaging tasks `build:linux_base`, `build:linux`, `build:deb`, `build:tar`, `build:rpm`, `build:windows` (templates wired in Task 3). `task check` is the CI gate.

- [ ] **Step 1: Replace `Taskfile.yml` wholesale**

```yaml
# https://taskfile.dev
version: '3'

set: [pipefail]

vars:
  APP_NAME: inventory
  APP_VERSION:
    sh: cat .app-version
  ARCH:
    sh: dpkg --print-architecture 2>/dev/null || echo "amd64"
  RPM_ARCH:
    sh: rpm --eval '%{_arch}' 2>/dev/null || uname -m 2>/dev/null || echo "x86_64"
  WAILS_TAGS: webkit2_41
  WAILS_BIN: build/bin/{{.APP_NAME}}
  TAR_NAME: "{{.APP_NAME}}-{{.APP_VERSION}}-linux-{{.RPM_ARCH}}"
  TAR_OUT: "dist/{{.APP_NAME}}-{{.APP_VERSION}}-linux-{{.RPM_ARCH}}.tar.gz"
  DEB_STAGE: dist/deb/{{.APP_NAME}}_{{.APP_VERSION}}_{{.ARCH}}
  DEB_OUT: dist/{{.APP_NAME}}_{{.APP_VERSION}}_{{.ARCH}}.deb
  RPM_OUT: "dist/{{.APP_NAME}}-{{.APP_VERSION}}-1.{{.RPM_ARCH}}.rpm"

tasks:
  default:
    desc: List all available tasks
    silent: true
    cmds:
      - task --list

  # ──────────────────────────────── version ───────────────────────────────────
  version:bump:
    desc: "Set a new version and sync everywhere (usage: task version:bump VERSION=1.2.3)"
    aliases: [bump]
    preconditions:
      - sh: 'printf "%s" "{{.VERSION}}" | grep -Eq "^[0-9]+\.[0-9]+\.[0-9]+$"'
        msg: "VERSION must follow semver — e.g. task version:bump VERSION=1.2.3"
    cmds:
      - printf '%s\n' "{{.VERSION}}" > .app-version
      - task: version:sync
      - echo "Bumped to {{.VERSION}}"

  version:sync:
    desc: "Sync version from .app-version into wails.json and frontend/package.json"
    aliases: [vsync]
    cmds:
      - sed -i -E 's/("productVersion"[[:space:]]*:[[:space:]]*")[^"]*(")/\1{{.APP_VERSION}}\2/' wails.json
      - sed -i -E '0,/("version"[[:space:]]*:[[:space:]]*")[^"]*(")/s//\1{{.APP_VERSION}}\2/' frontend/package.json
      - echo "Synced {{.APP_VERSION}} to wails.json and frontend/package.json"

  # ──────────────────────────────── dev ───────────────────────────────────────
  install:
    desc: Download Go deps and install frontend packages
    cmds:
      - go mod download
      - pnpm install --prefix frontend

  dev:
    desc: Launch the app in development mode (hot reload)
    cmds:
      - wails dev

  build:
    desc: Build the desktop binary via wails
    cmds:
      - wails build -tags {{.WAILS_TAGS}}

  # ──────────────────────────────── quality ───────────────────────────────────
  lint:
    desc: Lint Go code with golangci-lint
    cmds:
      - golangci-lint run --build-tags {{.WAILS_TAGS}} ./...

  typecheck:
    desc: Type-check the frontend (vue-tsc)
    aliases: [tc]
    cmds:
      - pnpm --prefix frontend run typecheck

  test:
    desc: Run Go and frontend tests
    cmds:
      - CGO_ENABLED=1 go test -tags {{.WAILS_TAGS}} ./...
      - pnpm --prefix frontend run test

  check:
    desc: Run lint + type-check + tests (CI-style)
    cmds:
      - task: lint
      - task: typecheck
      - task: test

  # ──────────────────────────────── build / packaging ─────────────────────────
  build:linux_base:
    desc: Build the Linux Wails binary
    platforms: [linux]
    cmds:
      - wails build -platform linux/amd64 -tags {{.WAILS_TAGS}}
      - echo "Binary → {{.WAILS_BIN}}"

  build:windows:
    desc: Build a standalone Windows .exe (run on Windows)
    platforms: [windows]
    cmds:
      - wails build -platform windows/amd64
      - echo Binary → {{.WAILS_BIN}}.exe

  build:deb:
    desc: Package the Linux binary as a .deb (requires build:linux_base)
    platforms: [linux]
    deps: [build:linux_base]
    vars:
      INSTALLED_SIZE:
        sh: du -sk {{.WAILS_BIN}} | cut -f1
    cmds:
      - rm -rf "{{.DEB_STAGE}}"
      - mkdir -p "{{.DEB_STAGE}}/DEBIAN" "{{.DEB_STAGE}}/usr/local/bin" "{{.DEB_STAGE}}/usr/share/applications"
      - cp "{{.WAILS_BIN}}" "{{.DEB_STAGE}}/usr/local/bin/{{.APP_NAME}}"
      - chmod 755 "{{.DEB_STAGE}}/usr/local/bin/{{.APP_NAME}}"
      - cp packaging/deb/usr/share/applications/inventory.desktop "{{.DEB_STAGE}}/usr/share/applications/{{.APP_NAME}}.desktop"
      - |
        sed \
          -e 's/__APP_NAME__/{{.APP_NAME}}/g' \
          -e 's/__APP_VERSION__/{{.APP_VERSION}}/g' \
          -e 's/__ARCH__/{{.ARCH}}/g' \
          -e 's/__INSTALLED_SIZE__/{{.INSTALLED_SIZE}}/g' \
          packaging/deb/DEBIAN/control.tpl \
          > "{{.DEB_STAGE}}/DEBIAN/control"
      - dpkg-deb --root-owner-group --build "{{.DEB_STAGE}}" "{{.DEB_OUT}}"
      - echo "Package → {{.DEB_OUT}}"

  build:tar:
    desc: Package the Linux binary as a .tar.gz archive
    platforms: [linux]
    deps: [build:linux_base]
    cmds:
      - rm -rf "dist/{{.TAR_NAME}}"
      - mkdir -p "dist/{{.TAR_NAME}}/bin" "dist/{{.TAR_NAME}}/share/applications"
      - cp "{{.WAILS_BIN}}" "dist/{{.TAR_NAME}}/bin/{{.APP_NAME}}"
      - chmod 755 "dist/{{.TAR_NAME}}/bin/{{.APP_NAME}}"
      - cp packaging/deb/usr/share/applications/inventory.desktop "dist/{{.TAR_NAME}}/share/applications/{{.APP_NAME}}.desktop"
      - |
        cat > "dist/{{.TAR_NAME}}/install.sh" << 'EOF'
        #!/usr/bin/env sh
        # Install inventory system-wide (run as root)
        set -e
        SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
        install -D -m 0755 "$SCRIPT_DIR/bin/inventory" /usr/local/bin/inventory
        install -D -m 0644 "$SCRIPT_DIR/share/applications/inventory.desktop" \
            /usr/share/applications/inventory.desktop
        echo "Installed to /usr/local/bin/inventory"
        EOF
      - chmod 755 "dist/{{.TAR_NAME}}/install.sh"
      - tar -C dist -czf "{{.TAR_OUT}}" "{{.TAR_NAME}}"
      - rm -rf "dist/{{.TAR_NAME}}"
      - echo "Archive → {{.TAR_OUT}}"

  build:rpm:
    desc: Package the Linux binary as an .rpm (requires rpmbuild)
    platforms: [linux]
    deps: [build:linux_base]
    vars:
      RPM_DATE:
        sh: date '+%a %b %d %Y'
      RPM_BUILD: "{{.HOME}}/rpmbuild"
    cmds:
      - mkdir -p "{{.RPM_BUILD}}/SOURCES" "{{.RPM_BUILD}}/SPECS" "{{.RPM_BUILD}}/RPMS" "{{.RPM_BUILD}}/SRPMS" "{{.RPM_BUILD}}/BUILD"
      - cp "{{.WAILS_BIN}}" "{{.RPM_BUILD}}/SOURCES/{{.APP_NAME}}"
      - cp packaging/deb/usr/share/applications/inventory.desktop "{{.RPM_BUILD}}/SOURCES/inventory.desktop"
      - |
        sed \
          -e 's/__APP_NAME__/{{.APP_NAME}}/g' \
          -e 's/__APP_VERSION__/{{.APP_VERSION}}/g' \
          -e 's/__RPM_ARCH__/{{.RPM_ARCH}}/g' \
          -e 's/__DATE__/{{.RPM_DATE}}/g' \
          packaging/rpm/inventory.spec.tpl \
          > "{{.RPM_BUILD}}/SPECS/{{.APP_NAME}}.spec"
      - rpmbuild -bb "{{.RPM_BUILD}}/SPECS/{{.APP_NAME}}.spec"
      - mkdir -p dist
      - cp "{{.RPM_BUILD}}/RPMS/{{.RPM_ARCH}}/{{.APP_NAME}}-{{.APP_VERSION}}-1.{{.RPM_ARCH}}.rpm" "{{.RPM_OUT}}"
      - echo "Package → {{.RPM_OUT}}"

  build:linux:
    desc: Build Linux binary + .deb + .tar.gz + .rpm packages
    platforms: [linux]
    cmds:
      - task: build:linux_base
      - task: build:deb
      - task: build:tar
      - task: build:rpm

  # ──────────────────────────────── clean ─────────────────────────────────────
  clean:
    desc: Remove build artefacts and caches
    prompt: "Remove build/bin, dist/, frontend/dist, and .task caches?"
    cmds:
      - rm -rf build/bin dist/ .task/
      - pnpm --prefix frontend run clean
      - echo "Clean done."
```

- [ ] **Step 2: Create `.golangci.yml` for deterministic linting**

```yaml
run:
  build-tags:
    - webkit2_41
  timeout: 5m

linters:
  enable:
    - govet
    - ineffassign
    - staticcheck
    - unused
    - errcheck
    - gosimple

issues:
  exclude-dirs:
    - frontend
```

- [ ] **Step 3: Delete the Python application and its caches**

```bash
git rm -r --quiet main.py app db tests pyproject.toml uv.lock inventory.spec
rm -rf .mypy_cache .ruff_cache .pytest_cache .venv .python-version
```

Note: `.python-version` is untracked/gitignored; `rm -rf` (not `git rm`) is correct for it.

- [ ] **Step 4: Re-sync the version into the JSON files**

```bash
task version:sync
grep -n '"productVersion"' wails.json
grep -m1 -n '"version"' frontend/package.json
```

Expected: both show `1.0.0`.

- [ ] **Step 5: Run the full quality gate**

Run: `task check`
Expected: `golangci-lint` passes, `vue-tsc` reports no errors, `go test` and `vitest` all PASS.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "build(cutover): Go-stack Taskfile + golangci config, delete Python app"
```

---

### Task 3: Update the packaging templates for the Wails runtime

**Files:**
- Modify: `packaging/deb/DEBIAN/control.tpl`
- Modify: `packaging/rpm/inventory.spec.tpl`
- Modify: `packaging/deb/usr/share/applications/inventory.desktop`

**Interfaces:**
- Consumes: `build:deb`/`build:tar`/`build:rpm` tasks and `{{.WAILS_BIN}}` from Task 2.
- Produces: `dist/inventory_1.0.0_<arch>.deb`, `dist/inventory-1.0.0-linux-<arch>.tar.gz`, `dist/inventory-1.0.0-1.<arch>.rpm`, each declaring webkit2gtk/gtk3 runtime deps. These artifact names are consumed by the release workflow in Task 5.

- [ ] **Step 1: Rewrite `packaging/deb/DEBIAN/control.tpl` with Wails runtime deps**

```
Package: __APP_NAME__
Version: __APP_VERSION__
Architecture: __ARCH__
Maintainer: IT Department <it@example.com>
Installed-Size: __INSTALLED_SIZE__
Depends: libwebkit2gtk-4.1-0, libgtk-3-0
Section: utils
Priority: optional
Description: IT Asset Inventory Manager
 A desktop GUI tool to track computers, phones, tablets,
 software licences, Windows keys, antivirus and subscriptions.
 Built with Go, Wails and Vue; uses a local SQLite database,
 no server required.
```

- [ ] **Step 2: Rewrite `packaging/rpm/inventory.spec.tpl` with Wails runtime deps**

```
Name:           __APP_NAME__
Version:        __APP_VERSION__
Release:        1%{?dist}
Summary:        IT Asset Inventory Manager

License:        GPL-3.0-or-later
BuildArch:      __RPM_ARCH__

Requires:       webkit2gtk4.1, gtk3

%description
A desktop GUI tool to track computers, phones, tablets,
software licences, Windows keys, antivirus and subscriptions.
Built with Go, Wails and Vue; uses a local SQLite database,
no server required.

%install
install -D -m 0755 %{_sourcedir}/inventory \
    %{buildroot}/usr/local/bin/inventory
install -D -m 0644 %{_sourcedir}/inventory.desktop \
    %{buildroot}/usr/share/applications/inventory.desktop

%files
/usr/local/bin/inventory
/usr/share/applications/inventory.desktop

%changelog
* __DATE__ IT Department <it@example.com> __APP_VERSION__-1
- Automated build via Taskfile
```

- [ ] **Step 3: Refresh the `.desktop` comment (Exec/Name already correct)**

Update only the `Comment` line in `packaging/deb/usr/share/applications/inventory.desktop` for accuracy; leave `Exec=/usr/local/bin/inventory`, `Name=IT Asset Inventory`, `Icon=computer` as-is:

```
Comment=Track computers, phones, tablets, licences and subscriptions
```

- [ ] **Step 4: Build the full Linux package set**

Run: `task build:linux`
Expected: prints `Package → dist/inventory_1.0.0_<arch>.deb`, `Archive → dist/inventory-1.0.0-linux-<arch>.tar.gz`, `Package → dist/inventory-1.0.0-1.<arch>.rpm`.

- [ ] **Step 5: Verify the deb declares the correct runtime dependency**

```bash
dpkg-deb -f dist/inventory_1.0.0_*.deb Depends
```

Expected: `libwebkit2gtk-4.1-0, libgtk-3-0`

- [ ] **Step 6: Commit**

```bash
git add packaging/
git commit -m "build(cutover): packaging templates declare webkit2gtk/gtk3 runtime deps"
```

---

### Task 4: Rework the CI workflow for the Go + frontend stack

**Files:**
- Replace: `.github/workflows/ci.yml`

**Interfaces:**
- Consumes: `task install`, `task lint`, `task typecheck`, `task test`, `task build` from Task 2.
- Produces: a CI job that runs the Go + frontend gates and a `wails build` smoke on every PR/push to `main`.

- [ ] **Step 1: Replace `.github/workflows/ci.yml`**

```yaml
name: CI

on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

env:
  FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true

jobs:
  test:
    name: Lint, type-check, test (Go 1.25 / Linux)
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.25"
          cache-dependency-path: go.sum

      - name: Set up pnpm
        uses: pnpm/action-setup@v4
        with:
          version: 9

      - name: Set up Node
        uses: actions/setup-node@v4
        with:
          node-version: "20"
          cache: pnpm
          cache-dependency-path: frontend/pnpm-lock.yaml

      - name: Install Task
        run: sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

      - name: Install Wails system dependencies
        run: |
          sudo apt-get update -qq
          sudo apt-get install -y --no-install-recommends \
            libgtk-3-dev \
            libwebkit2gtk-4.1-dev \
            pkg-config \
            build-essential

      - name: Install Wails CLI
        run: go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0

      - name: Install golangci-lint
        run: |
          curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
            | sh -s -- -b "$(go env GOPATH)/bin" v1.64.8

      - name: Install dependencies
        run: task install

      - name: Lint (golangci-lint)
        run: task lint

      - name: Type-check (vue-tsc)
        run: task typecheck

      - name: Test (go test + vitest)
        run: task test

      - name: Build smoke (wails build)
        run: task build
```

- [ ] **Step 2: Validate the workflow YAML**

```bash
python3 -c "import yaml,sys; yaml.safe_load(open('.github/workflows/ci.yml')); print('ci.yml: valid YAML')"
```

Expected: `ci.yml: valid YAML`. (If `python3` is unavailable, use any YAML linter, e.g. `yamllint .github/workflows/ci.yml`.)

- [ ] **Step 3: Confirm every CI step maps to a real task**

```bash
grep -E 'task (install|lint|typecheck|test|build)' .github/workflows/ci.yml
task --list | grep -E '^\* (install|lint|typecheck|test|build):'
```

Expected: each referenced task exists in `task --list`.

- [ ] **Step 4: Commit**

```bash
git add .github/workflows/ci.yml
git commit -m "ci(cutover): Go + frontend CI (lint, typecheck, test, build smoke)"
```

---

### Task 5: Rework the release workflow to ship the Go packages

**Files:**
- Replace: `.github/workflows/release.yml`

**Interfaces:**
- Consumes: `task build:linux` (→ `dist/*.deb`, `dist/*.tar.gz`, `dist/*.rpm`) and `task build:windows` (→ `build/bin/inventory.exe`) from Tasks 2–3; `.app-version` = `1.0.0`.
- Produces: a `v1.0.0` GitHub release with all four artifacts named `inventory-...`.

- [ ] **Step 1: Replace `.github/workflows/release.yml`**

```yaml
name: Release

on:
  workflow_dispatch:

permissions:
  contents: write

env:
  FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true

jobs:
  # ── 1. Tag the commit ──────────────────────────────────────────────────────
  tag:
    name: Create git tag
    runs-on: ubuntu-latest
    outputs:
      version: ${{ steps.ver.outputs.version }}
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
          token: ${{ secrets.GITHUB_TOKEN }}

      - name: Read version
        id: ver
        run: echo "version=$(cat .app-version | tr -d '[:space:]')" >> "$GITHUB_OUTPUT"

      - name: Create and push tag
        run: |
          git config user.name  "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git tag "v${{ steps.ver.outputs.version }}"
          git push origin "v${{ steps.ver.outputs.version }}"

  # ── 2a. Linux binary + packages ───────────────────────────────────────────
  build-linux:
    name: Build Linux packages
    needs: tag
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.25"
          cache-dependency-path: go.sum

      - name: Set up pnpm
        uses: pnpm/action-setup@v4
        with:
          version: 9

      - name: Set up Node
        uses: actions/setup-node@v4
        with:
          node-version: "20"
          cache: pnpm
          cache-dependency-path: frontend/pnpm-lock.yaml

      - name: Install build + packaging tools
        run: |
          sudo apt-get update -qq
          sudo apt-get install -y --no-install-recommends \
            libgtk-3-dev \
            libwebkit2gtk-4.1-dev \
            pkg-config \
            build-essential \
            rpm

      - name: Install Wails CLI
        run: go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0

      - name: Install Task
        run: sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

      - name: Build Linux packages
        run: task build:linux

      - name: Upload Linux artifacts
        uses: actions/upload-artifact@v4
        with:
          name: linux-packages
          path: |
            dist/*.deb
            dist/*.tar.gz
            dist/*.rpm

  # ── 2b. Windows .exe ──────────────────────────────────────────────────────
  build-windows:
    name: Build Windows binary
    needs: tag
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.25"
          cache-dependency-path: go.sum

      - name: Set up pnpm
        uses: pnpm/action-setup@v4
        with:
          version: 9

      - name: Set up Node
        uses: actions/setup-node@v4
        with:
          node-version: "20"
          cache: pnpm
          cache-dependency-path: frontend/pnpm-lock.yaml

      - name: Install Wails CLI
        run: go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0

      - name: Install Task
        uses: arduino/setup-task@v2
        with:
          version: 3.x
          repo-token: ${{ secrets.GITHUB_TOKEN }}

      - name: Build Windows binary
        run: task build:windows

      - name: Upload Windows artifact
        uses: actions/upload-artifact@v4
        with:
          name: windows-binary
          path: build/bin/inventory.exe

  # ── 3. Publish GitHub release ─────────────────────────────────────────────
  publish:
    name: Publish GitHub release
    needs: [tag, build-linux, build-windows]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Download Linux artifacts
        uses: actions/download-artifact@v4
        with:
          name: linux-packages
          path: dist/linux

      - name: Download Windows artifact
        uses: actions/download-artifact@v4
        with:
          name: windows-binary
          path: dist/windows

      - name: Create GitHub release
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          VERSION="${{ needs.tag.outputs.version }}"
          mv dist/windows/inventory.exe "dist/windows/inventory-${VERSION}-windows-x86_64.exe"
          gh release create "v${VERSION}" \
            --title "Release v${VERSION}" \
            --notes "Release v${VERSION}" \
            dist/linux/*.deb \
            dist/linux/*.tar.gz \
            dist/linux/*.rpm \
            "dist/windows/inventory-${VERSION}-windows-x86_64.exe"
```

- [ ] **Step 2: Validate the workflow YAML**

```bash
python3 -c "import yaml,sys; yaml.safe_load(open('.github/workflows/release.yml')); print('release.yml: valid YAML')"
```

Expected: `release.yml: valid YAML`.

- [ ] **Step 3: Commit**

```bash
git add .github/workflows/release.yml
git commit -m "ci(cutover): release workflow builds Wails deb/rpm/tar/exe as inventory v1.0.0"
```

---

### Task 6: Docs, pre-commit, and final cleanup

**Files:**
- Replace: `README.md`
- Delete: `README.gover.md` (the temp staged in Task 1)
- Modify: `docs/DEVELOPMENT_PLAN.md` (Phase 5/6 status)
- Replace: `.pre-commit-config.yaml`

**Interfaces:**
- Consumes: the final repo shape from Tasks 1–5.
- Produces: a Go-only README and pre-commit config; a tree with no Python or `gover` references and no stray binaries.

- [ ] **Step 1: Write the new root `README.md`**

```markdown
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
```

- [ ] **Step 2: Delete the temp gover README**

```bash
rm -f README.gover.md
```

- [ ] **Step 3: Update `docs/DEVELOPMENT_PLAN.md` Phase 5/6 status**

In the Phase Status table, change the Phase 5 row from `⏳ Pending` / "`task gover:package:deb` / `:rpm` are stubs" to reflect completion, e.g.:

```
| 5 — Packaging And Migration | ✅ Complete | `task build:linux` (deb/rpm/tar) + `task build:windows` (exe); Go app promoted to repo root as the default `inventory` product |
```

Add a one-line note under the table that the PyQt6 app has been removed and the Go app is now the sole product.

- [ ] **Step 4: Replace `.pre-commit-config.yaml` with Go hooks**

```yaml
repos:
  # Go formatting + vet via local toolchain
  - repo: https://github.com/dnephin/pre-commit-golang
    rev: v0.5.1
    hooks:
      - id: go-fmt
      - id: go-vet-mod
        args: [-tags=webkit2_41]

  # golangci-lint (uses .golangci.yml)
  - repo: local
    hooks:
      - id: golangci-lint
        name: golangci-lint
        entry: golangci-lint run --build-tags webkit2_41
        language: system
        types: [go]
        pass_filenames: false

  # Version sync — keep wails.json + frontend/package.json in step with .app-version
  - repo: local
    hooks:
      - id: version-sync
        name: version sync
        entry: task version:sync
        language: system
        pass_filenames: false
        files: ^\.app-version$

  # General file hygiene
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v6.0.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-json
      - id: check-merge-conflict
```

- [ ] **Step 5: Final sweep — no Python or `gover` references, no stray binaries**

```bash
echo "== python/pyqt references ==";  grep -rniE 'pyqt|pyinstaller|\buv sync\b|ruff|mypy|pytest' --include='*.md' --include='*.yml' --include='*.yaml' --include='*.toml' . | grep -v docs/superpowers/ || echo "none"
echo "== stray gover token ==";       grep -rni '\bgover\b' --include='*.go' --include='*.json' --include='*.yml' --include='*.yaml' . | grep -v wailsjs || echo "none"
echo "== stray binaries ==";          find . -maxdepth 2 -name 'gover2' -o -maxdepth 2 -name 'gover'  | grep -v '\.go$' || echo "none"
echo "== python files ==";            find . -name '*.py' -not -path './.git/*' || echo "none"
```

Expected: each section prints `none` (matches inside `docs/superpowers/` history specs are acceptable and excluded above).

- [ ] **Step 6: Run the full gate one more time**

Run: `task check`
Expected: all gates PASS.

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "docs(cutover): Go-only README, dev plan status, Go pre-commit hooks"
```

- [ ] **Step 8: Open the PR**

```bash
git push -u origin feat/go-default-cutover
gh pr create --title "Cut over to Go: retire Python, ship inventory v1.0.0" \
  --body "Promotes the Wails/Go app to the repo root as the canonical \`inventory\` product, deletes the PyQt6 app, and reworks Taskfile/CI/release/packaging for the Go stack. Ships as v1.0.0. See docs/superpowers/specs/2026-06-20-go-default-cutover-design.md."
```

---

## Self-Review

**Spec coverage** (design §4 → tasks):
- §4.1 restructure & rename → Task 1 ✓
- §4.2 Taskfile → Task 2 ✓
- §4.3 packaging → Task 3 ✓ (templates) + Task 2 ✓ (tasks)
- §4.4 CI → Task 4 ✓
- §4.5 release → Task 5 ✓
- §4.6 docs & cleanup → Task 6 ✓ (+ docs move in Task 1)
- §5 issues (broken frontend lint, no golangci config) → Task 2 ✓ (`.golangci.yml`, Go-only `lint` task)
- §6 acceptance criteria → covered across Tasks 1–6; final sweep in Task 6 Step 5
- §5 decision: `1.0.0` GA → Task 1 Step 8; module rename → Task 1 Step 5

**Placeholder scan:** No TBD/TODO; all file contents are complete. The only intentionally descriptive step is Task 6 Step 3 (editing an existing doc table whose exact current wording is shown in the design) — the replacement line is given verbatim.

**Type/name consistency:** `inventory` used uniformly as module path, binary, artifact prefix, and `wails.json` outputfilename. Build tag `webkit2_41` consistent across Taskfile, `.golangci.yml`, CI, pre-commit. Artifact names in Task 3 (`dist/inventory_1.0.0_*.deb` etc.) match the release upload globs in Task 5. `{{.WAILS_BIN}}` = `build/bin/inventory` consistent across all packaging tasks.
