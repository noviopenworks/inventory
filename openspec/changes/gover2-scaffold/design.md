## Approach

Follow the locked architecture from `gover-phase1-architecture` exactly. This change is pure scaffolding — no business logic, no database, no real data. Every implementation decision is already made; this change executes them.

## Architecture Overview

```
gover2/
├── go.mod              module gover2, go 1.22, wails v2 + modernc.org/sqlite deps
├── main.go             wails.Run(app.NewApp())
├── app.go              App struct — placeholder, no real bridge methods yet
├── wails.json          app config
├── internal/
│   ├── database/       Open() stub — returns nil DB for now
│   ├── models/         7 structs (Computer, Smartphone, Tablet, WindowsKey, Antivirus, OtherSoftware, User) + AppConfig + Alert
│   ├── services/       ListComputers() etc — return ([]models.X, nil) stubs
│   ├── bridge/         App struct with List* methods delegating to services stubs
│   ├── config/         AppConfig JSON load/save — reads $XDG_CONFIG_HOME/inventory/config.json
│   └── backup/         BackupDB() stub — returns nil
└── frontend/
    ├── package.json    pnpm workspace, Vue 3 + Vite + TypeScript
    ├── tailwind.config.ts  25 design tokens from gover-ui-spec
    ├── tsconfig.json   strict: true
    └── src/
        ├── main.ts
        ├── App.vue
        ├── layouts/AppShell.vue    sidebar + topbar + main content slot
        ├── features/               8 feature folders, each with an Index.vue stub
        ├── components/             9 shared components (stubs)
        ├── lib/api/index.ts        bridge wrappers returning Promise<[]>
        ├── lib/types/index.ts      7 domain interfaces + AppConfig + Alert
        └── stores/
            ├── assets.ts           Pinia store — empty items state
            └── ui.ts               density, darkMode, sidebarCollapsed
```

## Wails Init Strategy

Wails v2 `wails init` generates a project with a bundled Go+Vue template. We use the vanilla Vue template and immediately replace the generated frontend with our locked structure. The init command creates:
- `go.mod`, `main.go`, `app.go`, `wails.json`
- `frontend/` with Vite + vanilla Vue

We then:
1. Replace `frontend/` content with our locked Vue structure
2. Add `modernc.org/sqlite` to `go.mod` (needed in gover2-readonly, good to have early)
3. Add all 6 `internal/` packages as stubs
4. Wire `app.go` to expose the bridge struct

## Toolchain Prerequisites

Install before scaffolding (one-time):
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
corepack enable pnpm
```

## Taskfile Integration

Add to root `Taskfile.yml` at repo root (not inside `gover2/`):

```yaml
tasks:
  gover:install:
    dir: gover2
    cmds:
      - go mod download
      - pnpm install --prefix frontend
  gover:dev:
    dir: gover2
    cmds: [wails dev]
  gover:build:
    dir: gover2
    cmds: [wails build]
  gover:test:
    dir: gover2
    cmds:
      - go test ./...
      - pnpm --prefix frontend run test
  gover:lint:
    dir: gover2
    cmds:
      - golangci-lint run ./...
      - pnpm --prefix frontend run lint
  gover:typecheck:
    dir: gover2
    cmds: [pnpm --prefix frontend run typecheck]
  gover:clean:
    dir: gover2
    cmds:
      - rm -rf build/
      - pnpm --prefix frontend run clean
  gover:package:deb:
    dir: gover2
    cmds: [wails build -platform linux/amd64 -o gover2.bin && fpm ...]
  gover:package:rpm:
    dir: gover2
    cmds: [wails build -platform linux/amd64 -o gover2.bin && fpm ...]
```

## Verification

```
go build ./...           → exit 0, no import cycles
wails build              → exit 0, binary in gover2/build/
gover:dev                → window opens, sidebar shows 7 categories
pnpm typecheck           → exit 0, no TypeScript errors
```
