---
comet_change: gover-all-assets-parity
role: technical-design
canonical_spec: openspec
status: final
---

# Gover All Assets View Parity Design

> Gover's "All" sidebar route currently renders 4 columns (category/name/model/
> status) by merging three separate API calls in the frontend. The PyQt6
> reference (`db/queries.py:118-165` `_fetch_all`) shows a single backend
> `UNION ALL` query across the 3 hardware tables with a `LEFT JOIN users` to
> resolve `user_name`, returning 9 columns ordered `type ASC, id DESC`. This
> design closes that gap.

## Goal

Bring gover's "All Assets" view to behavioral parity with the PyQt6 reference
so users see the same 9 columns at the same ordering, with user names resolved
server-side, in a single round-trip.

## Scope

- Add a new backend `services.ListAllHardware(db)` that builds a `UNION ALL`
  across `computers`, `smartphones`, and `tablets`, with `LEFT JOIN users` to
  resolve `user_name`, returning the new `models.HardwareRow` type with 9
  fields.
- Add `bridge.App.ListAllHardware()` method.
- Regenerate Wails TypeScript bindings (one `wails generate module` call).
- Add `HardwareRow` TS type and `listAllHardware()` wrapper in `lib/api`.
- Rewrite `AllAssetsView.vue` to call the single new endpoint and render 9
  columns matching PyQt6's `ALL_DISPLAY_COLS` order.

## Explicit Non-Goals

- Do **not** change per-category views (`ComputersView`, `SmartphonesView`,
  `TabletsView`). They also currently lack a User column, but that's a
  separate parity cycle.
- Do **not** change existing `List*` signatures or `models.Computer`/
  `Smartphone`/`Tablet` types.
- Do **not** expand CSV export columns (separate gap, currently 5 columns vs
  PyQt6's 8/9).
- Do **not** introduce a new search model — the existing client-side
  `matchesQuery` over the 9 columns stays.
- Do **not** change the database schema.
- Do **not** touch the auto-open alerts behavior or any other view.

## Reference Behavior (PyQt6)

Source: `db/queries.py:118-165` `_fetch_all` + `db/config.py:165-175`
`ALL_DISPLAY_COLS`.

`_fetch_all` iterates `ALL_TAB_CATS = {"Computer", "Smartphone", "Tablet"}`
and emits one SELECT per category:

```sql
SELECT t.id, '<Category>' AS type,
       t.name AS name, t.model AS model,
       u.name AS user_name,
       t.status, t.purchase_date, t.warranty_expiry, t.notes
FROM <table> t
LEFT JOIN users u ON t.user_id = u.id
[optional WHERE clause for search]
```

Joined with `UNION ALL`, terminated with `ORDER BY type, id DESC`. Columns
returned (in `ALL_DISPLAY_COLS` order):

1. `type` — `"Computer"` / `"Smartphone"` / `"Tablet"`
2. `id`
3. `name`
4. `model`
5. `user_name` — resolved from the JOIN; NULL if no user assigned
6. `status`
7. `purchase_date`
8. `warranty_expiry`
9. `notes`

## Approach

### Backend (`gover/internal/services/services.go`)

Add a table-driven `UNION ALL` builder mirroring the `alertSources` pattern
landed in the alerts parity work:

```go
type hardwareSource struct {
    table string
    typ   string // displayed in the Type column
}

var hardwareSources = []hardwareSource{
    {"computers",   "Computer"},
    {"smartphones", "Smartphone"},
    {"tablets",     "Tablet"},
}
```

Per source, emit a SELECT with `COALESCE` on every nullable column so the Go
struct can use plain `string` fields (no `*string` juggling, no frontend
null-checks — matches the existing `Alert` pattern from `services.go`):

```sql
SELECT t.id,
       '<typ>'  AS type,
       COALESCE(t.name, '')            AS name,
       COALESCE(t.model, '')           AS model,
       COALESCE(u.name, '')            AS user_name,
       COALESCE(t.status, '')          AS status,
       COALESCE(t.purchase_date, '')   AS purchase_date,
       COALESCE(t.warranty_expiry, '') AS warranty_expiry,
       COALESCE(t.notes, '')           AS notes
FROM   <table> t
LEFT JOIN users u ON t.user_id = u.id
```

Join with `UNION ALL`; terminate with `ORDER BY type ASC, id DESC`. Build
query with `strings.Join`; run with `db.Query`. Scan into
`[]models.HardwareRow`. No user input participates in identifier position —
all `table`/`typ` values come from the package-level hardcoded slice.

The WHERE clause for search is **not** added server-side. Search stays
client-side via the existing `matchesQuery` filter in `AllAssetsView.vue`
(gover convention established before this change).

### Models (`gover/internal/models/models.go`)

New struct:

```go
// HardwareRow is a unified row across the 3 hardware tables, for the "All"
// view. Mirrors the field set of PyQt6's _fetch_all result.
type HardwareRow struct {
    ID              int    `json:"id"`
    Type            string `json:"type"`
    Name            string `json:"name"`
    Model           string `json:"model"`
    UserName        string `json:"userName"`
    Status          string `json:"status"`
    PurchaseDate    string `json:"purchaseDate"`
    WarrantyExpiry  string `json:"warrantyExpiry"`
    Notes           string `json:"notes"`
}
```

All plain strings — matches `models.Alert` style and lets the frontend render
without null-coalescing boilerplate.

### Bridge (`gover/internal/bridge/bridge.go`)

New method, following the established `List*` pattern in the file:

```go
func (a *App) ListAllHardware() ([]models.HardwareRow, error) {
    if a.db == nil {
        return nil, errNoDB
    }
    out, err := services.ListAllHardware(a.db)
    if out == nil {
        out = []models.HardwareRow{}
    }
    return out, err
}
```

### Wails Bindings Regeneration

From `gover/` directory: `wails generate module`. This updates:

- `gover/frontend/wailsjs/go/bridge/App.js` — adds `ListAllHardware` JS stub.
- `gover/frontend/wailsjs/go/bridge/App.d.ts` — adds the TS signature.

Both files are committed (already the convention in this repo).

### Frontend Types + API

`gover/frontend/src/lib/types/index.ts` — add:

```ts
export interface HardwareRow {
  id: number
  type: string
  name: string
  model: string
  userName: string
  status: string
  purchaseDate: string
  warrantyExpiry: string
  notes: string
}
```

`gover/frontend/src/lib/api/index.ts`:

- Add `HardwareRow` to the `export type { ... } from '@/lib/types'` block.
- Add `HardwareRow` to the `import type { ... }` block.
- Add wrapper:

```ts
export const listAllHardware = (): Promise<HardwareRow[]> =>
  call<HardwareRow[]>('ListAllHardware').catch(() => [])
```

### AllAssetsView.vue

Replace the 3-call merge with a single call. The new columns (in PyQt6
`ALL_DISPLAY_COLS` order):

```ts
const columns = [
  { key: 'type',           label: 'Type' },
  { key: 'id',             label: 'ID' },
  { key: 'name',           label: 'Name' },
  { key: 'model',          label: 'Model' },
  { key: 'userName',       label: 'User' },
  { key: 'status',         label: 'Status' },
  { key: 'purchaseDate',   label: 'Purchase Date' },
  { key: 'warrantyExpiry', label: 'Warranty Expiry' },
  { key: 'notes',          label: 'Notes' },
]
```

`onMounted` becomes a single `await listAllHardware()`. The `useSearchStore`
+ `matchesQuery` filter continues to work — it iterates `columns` and the new
columns are all string fields that `matchesQuery` can search across.

Drop the now-unused `listComputers`/`listSmartphones`/`listTablets` imports
from this file.

### Tests

Backend (`gover/internal/services/services_test.go`) — add 5 tests:

- `TestListAllHardware_Empty` — no rows → empty slice.
- `TestListAllHardware_AllThreeTables` — insert 1 row in each table → 3 rows
  in order: Computer, Smartphone, Tablet (matches `ORDER BY type ASC`).
- `TestListAllHardware_OrderByIDDesc` — 2 computers with sequential IDs →
  newest first (matches `ORDER BY id DESC` within type).
- `TestListAllHardware_ResolvesUserName` — computer with `user_id` pointing
  at a user → `UserName` is the user's `name`.
- `TestListAllHardware_NullUser` — computer with `user_id = NULL` →
  `UserName == ""` (via `COALESCE(u.name, '')`).

Frontend (`gover/frontend/src/features/assets/AllAssetsView.test.ts` — does
not exist yet; create):

- `'renders 9 columns with the expected labels'` — mount, mock
  `listAllHardware` to return a known row, assert the rendered header labels.
- `'calls listAllHardware once on mount'` — mock the API, assert it was
  called exactly once.

## Verification

Target gates:

- `cd gover && go test ./internal/services/...` — all old + 5 new tests pass.
- `cd gover && go test ./...` — full Go suite green (85 + 5 = 90 tests).
- `cd gover && go build ./...` — compiles cleanly (includes the new bridge
  method).
- `cd gover/frontend && pnpm run test` — Vitest green including the new
  `AllAssetsView.test.ts`.
- `cd gover/frontend && pnpm run typecheck` — no type drift from the new
  type and wrapper.
- `wails generate module` from `gover/` — produces the expected diff in
  `frontend/wailsjs/go/bridge/App.{js,d.ts}` (additive only).

## Risks / Trade-offs

- **Wails binding regen required** — adding a bridge method means the
  generated `frontend/wailsjs/` files change. Mitigated by committing the
  regenerated bindings; `wails generate module` is deterministic.
- **NULL vs empty-string distinction lost** — `COALESCE(col, '')` collapses
  SQL NULL to `''` for nullable columns (`purchase_date`, `warranty_expiry`,
  `notes`). PyQt6 returns NULLs and the Python frontend renders them as
  empty cells; net behavior is identical, but the JSON shape differs (no
  `null`s). Acceptable: display is the same and the frontend's
  `row[col.key] ?? ''` already normalizes both.
- **Binding regeneration order** — the frontend test that imports
  `listAllHardware` won't compile until both the regenerated bindings AND the
  `lib/api/index.ts` wrapper exist. The implementation plan must order:
  backend → bridge → bindings regen → TS type + wrapper → frontend view +
  test.
- **Per-category views still lack User column** — explicitly out of scope.
  Listed as a follow-up parity cycle.
- **`type ASC` ordering is alphabetical** — puts "Computer" < "Smartphone" <
  "Tablet", which matches PyQt6 verbatim. No custom ordering introduced.

## Files Touched

- `gover/internal/services/services.go` — add `hardwareSource` type+var and
  `ListAllHardware` function; add `strings` import (already present from the
  alerts work).
- `gover/internal/services/services_test.go` — 5 new tests.
- `gover/internal/models/models.go` — add `HardwareRow`.
- `gover/internal/bridge/bridge.go` — add `ListAllHardware` method.
- `gover/frontend/wailsjs/go/bridge/App.js` — regenerated.
- `gover/frontend/wailsjs/go/bridge/App.d.ts` — regenerated.
- `gover/frontend/src/lib/types/index.ts` — add `HardwareRow`.
- `gover/frontend/src/lib/api/index.ts` — add `listAllHardware` + type
  re-export.
- `gover/frontend/src/features/assets/AllAssetsView.vue` — single-call +
  9 columns.
- `gover/frontend/src/features/assets/AllAssetsView.test.ts` — new file.
