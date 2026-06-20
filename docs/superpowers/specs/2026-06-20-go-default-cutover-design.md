# Design: Cut over to Go, retire Python, rework CI/CD

**Date:** 2026-06-20
**Status:** Approved (brainstorming) — pending implementation plan
**Topic:** Make the Wails/Go `gover` app the canonical `inventory` product, delete the PyQt6 Python app, and rework all tooling (Taskfile, CI, release, packaging, docs).

---

## 1. Goal

Today the repository ships two apps:

- **Python (root):** PyQt6 desktop app, `inventory` v0.5.1, managed with `uv`/`ruff`/`mypy`/`pytest`, packaged with PyInstaller → `.deb`/`.rpm`/`.tar.gz`/`.exe`. This is the **only thing CI tests and the only thing that ships.**
- **Go (`gover/`):** Wails v2 + Go 1.25 + Vue 3 + TypeScript rewrite. Phases 0–4 complete (CRUD, alerts, search, export, settings). **Phase 5 packaging is stubbed; gover is not in CI at all.**

The goal is a **hard cutover**: make the Go app the one and only `inventory` product, build a real packaging + CI/CD pipeline for it, and delete the Python app entirely.

## 2. Decisions

| # | Decision | Choice |
|---|----------|--------|
| 1 | Cutover strategy | **Hard cutover now** — build full packaging + Go CI, delete all Python in this work, point releases at the Go app. |
| 2 | Repo layout | **Promote `gover/` to repo root** — the Go app becomes the repo. |
| 3 | Product name | **Rename to `inventory`** — binary, packages, artifacts, `.desktop`, and the Go module path (`gover` → `inventory`). |
| 4 | Release targets | **Parity** — Linux `.deb`/`.rpm`/`.tar.gz` + Windows `.exe`. No macOS. |
| 5 | Starting version | **`1.0.0`** — marks the Go rewrite as GA. |

## 3. Execution approach

**One feature branch, ordered phases so the tree never goes red**, merged as a single PR. Phase order:

1. Restructure + rename (move `gover/*` to root, rename module/product, delete Python)
2. Taskfile rework
3. Packaging pipeline (deb/rpm/tar/exe)
4. CI rework
5. Release rework
6. Docs + cleanup

Each phase ends green (build + existing tests pass) before the next begins.

**Rejected alternatives:**
- *Two PRs* (dev cutover, then packaging) — hard cutover means a window with no releasable artifact between the two PRs.
- *Big-bang single commit* — unreviewable and impossible to bisect.

## 4. Detailed changes by area

### 4.1 Restructure & rename

- Move to repo root: `gover/{main.go, app.go, go.mod, go.sum, wails.json, internal/, frontend/, build/}`.
- `go.mod`: change `module gover` → `module inventory`. Rewrite every `gover/internal/...` import to `inventory/internal/...` (mechanical `sed` across `*.go`, then `gofmt`). Remove the dead commented `replace github.com/wailsapp/wails/v2 ...` line.
- `wails.json`: `outputfilename` `gover` → `inventory`; `info.productName` `Gover` → `Inventory`; `info.productVersion` → from `.app-version`.
- **Verify Wails bindings unaffected:** `frontend/wailsjs/go/<package>/...` are keyed by Go *package* names (e.g. `bridge`), not the module path, so the module rename should not change them. Confirm by regenerating bindings (`wails generate module`) or a build after rename; if anything references `gover`, fix it.
- **Delete Python:** `main.py`, `app/`, `db/`, `tests/`, `pyproject.toml`, `uv.lock`, `inventory.spec`, `.python-version`, and tool caches `.mypy_cache/`, `.ruff_cache/`, `.pytest_cache/`, `.venv/`.

### 4.2 Taskfile

- **Remove** Python tasks: `install` (uv), `run` (python main.py), `lint`/`lint:fix`/`fmt`/`fmt:check` (ruff), `type:check` (mypy), `test`/`test:verbose`/`test:cov` (pytest), `build:linux_base`/`build:windows`/`build:deb`/`build:tar`/`build:rpm`/`build:linux` (PyInstaller), `hooks:*` (pre-commit), Python-based `version:sync`.
- **Promote** the `gover:*` tasks to top-level (drop the `gover:` prefix and `dir: gover`):
  - `install` → `go mod download` + `pnpm install --prefix frontend`
  - `dev` → `wails dev`
  - `build` → `wails build`
  - `test` → `go test ./...` + `pnpm --prefix frontend run test`
  - `lint` → `golangci-lint run ./...` (frontend lint dropped — see §5)
  - `typecheck` → `pnpm --prefix frontend run typecheck` (vue-tsc)
  - `clean` → `rm -rf build/` + `pnpm --prefix frontend run clean`
  - `check` → `lint` + `typecheck` + `test`
- **Version tasks:** `version:bump VERSION=x.y.z` writes `.app-version` then `version:sync`. `version:sync` reads `.app-version` and writes `wails.json` `info.productVersion` and `frontend/package.json` `version` (and exposes it to the binary via Go ldflags `-X main.version=`). Replace the Python one-liners with a small portable script (Go program, `task`-invoked, or `jq`/`sed`).
- **New packaging tasks:**
  - `build:linux` (platforms: linux) → `wails build -platform linux/amd64` then `build:deb` + `build:tar` + `build:rpm`.
  - `build:deb`/`build:tar`/`build:rpm` → reuse existing staging logic; binary source becomes `build/bin/inventory`.
  - `build:windows` (platforms: windows) → `wails build -platform windows/amd64` → `build/bin/inventory.exe`.

### 4.3 Packaging

The existing deb/rpm/tar staging logic in the Taskfile only stages a single binary plus a `.desktop` file, so it is largely reusable. Changes:

- **Binary source:** Wails output `build/bin/inventory` (instead of PyInstaller `dist/linux/inventory`).
- **Runtime dependencies** (Wails apps are *not* self-contained like PyInstaller bundles):
  - `packaging/deb/DEBIAN/control.tpl`: `Depends: libwebkit2gtk-4.1-0, libgtk-3-0`
  - `packaging/rpm/inventory.spec.tpl`: `Requires: webkit2gtk4.1, gtk3`
- **`.desktop`:** `Exec=inventory`, product name/description updated; ship the Wails app icon (`build/appicon.png`) into `usr/share/icons` or alongside the desktop entry.
- Description strings updated from "PyQt6" to the Wails/Go stack.

### 4.4 CI (`.github/workflows/ci.yml`)

Replace the Python job with a Go + frontend job on `ubuntu-latest`:

- **Setup:** Go 1.25, Node 20 + pnpm, Task.
- **System deps:** `libwebkit2gtk-4.1-dev libgtk-3-dev` (plus existing apt update). **Required because `go test ./...` compiles the `internal/bridge` package, which imports the Wails runtime → needs `CGO_ENABLED=1` and webkit.** Build tag `webkit2_41` (matches `wails.json` `build:tags`).
- **Steps:** `task install` → `task lint` (golangci-lint) → `task typecheck` (vue-tsc) → `task test` (go test + vitest) → `wails build` smoke (catch build breaks without packaging).
- Keep triggers: `pull_request` + `push` on `main`. Keep `FORCE_JAVASCRIPT_ACTIONS_TO_NODE24` if still needed by actions.
- CI stays Linux-only for tests; Windows is covered by the release workflow.

### 4.5 Release (`.github/workflows/release.yml`)

- **`tag` job:** unchanged — read `.app-version`, create + push `v<version>` tag.
- **`build-linux` (ubuntu-latest):** checkout → setup Go + Node + pnpm + Task → install `libwebkit2gtk-4.1-dev libgtk-3-dev rpm` → install Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0`) → `task build:linux` → upload `dist/*.deb`, `dist/*.tar.gz`, `dist/*.rpm`.
- **`build-windows` (windows-latest):** checkout → setup Go + Node + pnpm + Task → install Wails CLI → `task build:windows` → upload `build/bin/inventory.exe`. Bare `.exe` (parity; no NSIS). WebView2 runtime is present on modern Windows.
- **`publish` job:** download artifacts, rename Windows exe to `inventory-<version>-windows-x86_64.exe`, `gh release create` with all artifacts. Adapt existing logic.

### 4.6 Docs & cleanup

- **README:** merge `gover/README.md` into root `README.md`. Drop the "rewrite in progress" framing — it is *the* app now. Rewrite Installation / Running from source / Development tasks / Building sections for the Go/Task/Wails toolchain. Remove all Python references.
- **Move** `gover/docs/{DEVELOPMENT_PLAN.md, ARCHITECTURE_NOTES.md, DECISIONS.md, UI_DESIGN_REFERENCE.html}` → `docs/`. Update `DEVELOPMENT_PLAN.md` Phase 5/6 status (packaging now done; default cutover done).
- **`.pre-commit-config.yaml`:** rework from ruff/mypy hooks to `gofmt`/`golangci-lint` (Go). Frontend hooks optional.
- **`.gitignore`:** merge `gover/.gitignore` rules into root; ensure `build/bin/` and built binaries are ignored.
- **Delete stray binaries:** `gover/gover2` (5.2M) and `gover/build/bin/gover2` (12.2M). Confirm whether tracked in git (`git rm` if so).
- **OpenSpec:** specs under `openspec/specs/` are already gover-focused; keep. Update the repository-hygiene spec if it references Python tooling.

## 5. Issues found (fixed as part of this work)

- **`task gover:lint` is currently broken:** it runs `pnpm --prefix frontend run lint`, but `frontend/package.json` has no `lint` script and no eslint dependency. Resolution: Go uses `golangci-lint`; the frontend relies on `vue-tsc` typecheck. The new top-level `lint` task drops the frontend `pnpm run lint` call. (Adding eslint is optional and out of scope unless requested.)
- **No `.golangci.yml`:** add a minimal config so `golangci-lint run ./...` is deterministic in CI.

## 6. Acceptance criteria

- Repo root contains the Go app (`main.go`, `internal/`, `frontend/`, `wails.json`); no Python source, manifests, or caches remain.
- `go.mod` module is `inventory`; `go build ./...` and `wails build` succeed.
- `task check` (lint + typecheck + test) passes locally and in CI.
- CI workflow runs Go + frontend gates on every PR/push and goes green.
- Release workflow produces `inventory` `.deb`, `.rpm`, `.tar.gz` (Linux) and `.exe` (Windows), attached to a `v1.0.0` GitHub release; deb/rpm declare webkit2gtk/gtk3 runtime deps.
- `README.md` documents only the Go/Wails app; the "rewrite in progress" framing is gone.
- Stray `gover2` binaries removed; `.gitignore` prevents recurrence.

## 7. Out of scope

- macOS packaging (`.dmg`/`.app`).
- Windows NSIS installer (bare `.exe` only).
- Adding ESLint to the frontend.
- Any new product features — this is a tooling/cutover change only; gover behavior is frozen at its current Phase 0–4 state.
