---
comet_change: gover-alerts-parity
role: technical-design
canonical_spec: openspec
status: final
---

# Gover Alerts Parity Design

> Gover's `services.GetAlerts` currently checks only `expiry_date` on
> `antivirus` and `other_software`. The PyQt6 reference (`db/alerts.py`) also
> checks `warranty_expiry` on `computers`, `smartphones`, and `tablets`, and
> excludes `Retired` items. This design closes that gap.

## Goal

Bring gover's expiry/warranty alert query to behavioral parity with the PyQt6
reference so users see every alert the Python app would have surfaced — no
silent data gaps.

## Scope

- Extend `services.GetAlerts` to cover all 5 non-User tables that have an
  expiry/warranty column.
- Add the `status != 'Retired'` filter that PyQt6 applies.
- Preserve the current single-method API; severity (`expired` vs `expiring`)
  stays computed from `days_remaining` sign.
- Surface a summary line in `AlertsModal.vue` matching PyQt6's "X expired · Y
  expiring within N days" header.
- Prettify the severity badge label (`expired` → `Expired`, `expiring` →
  `Expiring soon`) for readability without changing the underlying value.

## Explicit Non-Goals

- Do **not** add a `Kind` field distinguishing warranty from license expiry.
  PyQt6 does not expose this; `Category` already conveys it.
- Do **not** include `windows_keys` in alerts. PyQt6 excludes them because the
  table has no expiry/warranty column.
- Do **not** change the `GetAlerts` signature, bridge method, or generated
  Wails bindings.
- Do **not** change the database schema.
- Do **not** re-route the auto-open-on-startup behavior in `AppShell.vue`; it
  is already wired and continues to work because the API is unchanged.
- Do **not** introduce sorting in the frontend; the SQL `ORDER BY` is
  authoritative.

## Reference Behavior (PyQt6)

Source: `db/alerts.py:11-50`. Two functions, `get_expiring_soon(days)` and
`get_expired()`, iterate `TABLE_CONFIG` and skip `Users`. For each remaining
category, they pick the date column (`expiry_date` if present, else
`warranty_expiry`) and skip if neither exists. Both apply `status != 'Retired'`.

- `get_expiring_soon(N)`: rows where `today <= date <= today + N days`, ordered
  by date ascending.
- `get_expired()`: rows where `date < today`, ordered by date ascending.

UI (`app/main_window.py:561-607`) merges both lists with status badges (🔴
Expired / 🟡 Expiring soon), shows a summary header (`X expired · Y expiring
within N days`), and a table with columns Status | Name | Type | Date.

## Approach

### Backend (`gover/internal/services/services.go`)

Replace the current hand-written two-table `UNION ALL` with a table-driven
generator:

```go
type alertSource struct{ table, category, dateCol string }

var alertSources = []alertSource{
    {"computers",       "Computer",       "warranty_expiry"},
    {"smartphones",     "Smartphone",     "warranty_expiry"},
    {"tablets",         "Tablet",         "warranty_expiry"},
    {"antivirus",       "Antivirus",      "expiry_date"},
    {"other_software",  "Other Software", "expiry_date"},
}
```

For each source, emit a `SELECT` of the form:

```sql
SELECT '<category>'  AS category,
       id,
       COALESCE(name, '') AS name,
       <dateCol>      AS expiry_date,
       CAST(julianday(<dateCol>) - julianday('now') AS INTEGER) AS days_remaining
FROM   <table>
WHERE  <dateCol> IS NOT NULL
  AND  date(<dateCol>) <= date('now', ?)   -- ? = '+N days'
  AND  status != 'Retired'
```

Join with `UNION ALL`; terminate with `ORDER BY days_remaining ASC` so the most
overdue row appears first. Pass the threshold as a parameter once per source
(N placeholders for N sources). After scanning, set `Severity` in Go:
`daysRemaining < 0 → "expired"`, otherwise `"expiring"`.

`strings.Join` builds the query; `db.Query(query, params...)` runs it. Returns
`[]models.Alert` (`nil` → `[]models.Alert{}` is the bridge's job, not the
service's).

### Models

No changes. `models.Alert` already has the right fields:
`Category, ID, Name, ExpiryDate, DaysRemaining, Severity`.

### Bridge

No changes. `bridge.App.GetAlerts()` already calls
`services.GetAlerts(a.db, a.cfg.ExpiryWarningDays)` and returns the slice; it
transparently picks up the new behavior on recompile.

### Frontend

`gover/frontend/src/components/AlertsModal.vue`:

1. Add a `<p>` summary line above the table, derived from `props.alerts`:
   - `expiredCount = alerts.filter(a => a.severity === 'expired').length`
   - `expiringCount = alerts.length - expiredCount`
   - Render: `<b>{expiredCount}</b> expired · <b>{expiringCount}</b> expiring within N days`
   - `N` is not currently a prop; pass it as a new optional prop
     `warningDays?: number` (default 30) to avoid coupling the modal to the
     config store.
2. Prettify the severity badge label via a small lookup:
   `{ expired: 'Expired', expiring: 'Expiring soon' }[a.severity] ?? a.severity`.
3. Pass `warningDays` from `AppShell.vue` when binding `<AlertsModal>`:
   `:warning-days="ui.expiryWarningDays"`. This requires exposing
   `expiryWarningDays` from `useUiStore` as a computed derived from the existing
   `config` ref (the value is already loaded by `loadFromConfig` but not
   currently in the store's return — see `gover/frontend/src/stores/ui.ts`).
   `AppShell.vue` already imports `useUiStore`, so once the store exposes the
   field the binding is one line.
4. No sorting logic — backend already returns expired-first by ascending
   overdue.

`AlertsView.vue` (sidebar route) inherits the richer dataset for free; its
current columns (`category`, `name`, `expiryDate`, `daysRemaining`, `severity`)
remain valid. No changes required.

### Tests

Backend (`gover/internal/services/services_test.go`) — add:

- `TestGetAlerts_ComputerWarrantyExpiring` — insert computer with
  `warranty_expiry = today + 15d`, status `Active`; expect 1 alert with
  `Category == "Computer"` and `Severity == "expiring"`.
- `TestGetAlerts_SmartphoneWarrantyExpired` — insert smartphone with
  `warranty_expiry = today - 1d`, status `Active`; expect `Severity == "expired"`.
- `TestGetAlerts_TabletWarrantyCovered` — insert tablet with
  `warranty_expiry = today + 5d`; expect 1 alert with `Category == "Tablet"`.
- `TestGetAlerts_RetiredExcluded` — insert antivirus with `status = 'Retired'`
  and an expired `expiry_date`; expect 0 alerts.

Existing tests (`ExpiringWithin30Days`, `NotExpiringSoon`, `AlreadyExpired`)
continue to pass — they only touch antivirus/other_software rows, which remain
in scope.

Frontend (`gover/frontend/src/components/AlertsModal.test.ts`) — add:

- `'shows expired/expiring summary counts'` — pass a known mix of alerts,
  assert the summary text contains both counts.

## Verification

Target gates:

- `cd gover && go test ./internal/services/...` — all old + new tests pass.
- `cd gover && go test ./...` — full Go suite green.
- `cd gover/frontend && pnpm run test` — Vitest green including new summary
  assertion.
- `cd gover/frontend && pnpm run typecheck` — no type drift from the new
  optional prop.
- `cd gover && go build ./...` — compiles cleanly.

No DB migration to verify (schema unchanged). No smoke test needed (alerts
auto-open behavior already covered by `AppShell.test.ts`).

## Risks / Trade-offs

- **Existing tests break** → they won't: the three current tests only insert
  into `antivirus`/`other_software`, both still in `alertSources`. Mitigated by
  keeping the API signature identical.
- **Parameter count mismatch in SQL** → mitigated by building `params` slice
  in lockstep with the per-source SELECT fragments in the same loop.
- **`windows_keys` accidentally included** → can't happen: the source list is
  explicit; PyQt6 also excludes them.
- **`status` NULL handling** → schema has `DEFAULT 'Active'`, so NULL doesn't
  occur in practice. `!= 'Retired'` matches PyQt6 verbatim, including its
  NULL-exclusion semantics; no special handling needed.
- **Frontend coupling to `warningDays`** → mitigated by making the prop
  optional with a default; `useUiStore` already owns the value, so the binding
  is one line.

## Files Touched

- `gover/internal/services/services.go` — rewrite `GetAlerts`.
- `gover/internal/services/services_test.go` — 4 new tests.
- `gover/frontend/src/components/AlertsModal.vue` — summary line, badge labels,
  optional `warningDays` prop.
- `gover/frontend/src/components/AlertsModal.test.ts` — 1 new test for summary.
- `gover/frontend/src/stores/ui.ts` — expose `expiryWarningDays` as a computed.
- `gover/frontend/src/stores/ui.test.ts` — assert `ui.expiryWarningDays` reads
  through from `config`.
- `gover/frontend/src/layouts/AppShell.vue` — pass `warningDays` to the modal
  (one-line binding).
