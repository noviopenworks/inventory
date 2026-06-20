# Brainstorm Summary

- Change: gover2-readonly
- Date: 2026-06-18

## Confirmed Technical Approach

Layer-by-layer replacement of gover2-scaffold stubs with real implementations: database → models → config → services → bridge → frontend.

**Bridge startup pattern (confirmed: Option B — Wails OnStartup callback)**:
- `bridge.App` holds `*sql.DB`, `models.AppConfig`, `context.Context`
- `NewApp()` returns empty struct (no DB open)
- `startup(ctx context.Context)` wired via `main.go` `OnStartup:` — loads config, attempts `database.Open(cfg.DBPath)`, stores result; if DB fails to open, `a.db` stays nil
- All List methods guard with `if a.db == nil { return nil, errors.New("no database open") }`

**Config layer**: `config.Load()` reads `$XDG_CONFIG_HOME/inventory/config.json`; default DBPath = `~/.local/share/inventory/inventory.db`. `config.Save()` writes JSON atomically.

**Database layer**: `database.Open()` adds `PRAGMA journal_mode=WAL`. `InitSchema()` runs all 7 CREATE TABLE IF NOT EXISTS DDLs verbatim from `db/schema.py`.

**Services layer**: 7 `List*` functions + `GetAlerts` (expiry ≤30 days from antivirus + other_software). Simple SELECT queries, scan into model slices.

**Frontend**: `lib/api/index.ts` imports from `wailsjs/go/bridge/App` (generated bindings), wraps each call with `.catch(() => [])`. `DataTable.vue` gets real `columns`/`rows`/`loading` props. Each view fetches on mount. `AllAssetsView.vue` uses `Promise.all` for 3 tables + prepends Category column.

## Key Trade-offs and Risks

| Risk | Mitigation |
|---|---|
| WAL mode changes Python app's DB permanently | Document: apps must not run simultaneously; WAL is strictly better |
| `a.db == nil` at startup (DB path wrong/missing) | All methods return error → Vue shows empty state; no crash |
| Wails bindings throw JS exceptions on error | api wrapper catches → `[]`; silent failure acceptable for read-only phase |
| `wailsjs/go/models.ts` must be regenerated after real structs | `wails build` in verification step regenerates bindings |
| testify + vitest not yet in deps | Added in Task 0 (setup step before task 1) |

## Testing Strategy

- Go: table-driven tests with in-memory SQLite (`:memory:`) for database, services, bridge
- Target: ≥70% services coverage, ≥50% bridge coverage
- `GetAlerts`: test expiry within 30 days appears, beyond 30 days does not
- Frontend: Vitest + Vue Test Utils; DataTable renders columns/rows/empty-state; StatusBadge maps 7 status values to correct CSS classes
- No E2E; integration smoke test is manual (wails dev + real DB)

## Spec Patches

None — the OpenSpec tasks.md and design.md are complete and consistent with confirmed design.
