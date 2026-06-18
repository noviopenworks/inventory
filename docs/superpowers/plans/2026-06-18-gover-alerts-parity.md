# Gover Alerts Parity Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Bring gover's expiry/warranty alert query to behavioral parity with the PyQt6 reference (`db/alerts.py`) so computers, smartphones, and tablets contribute warranty-expiry alerts alongside the existing software licence-expiry alerts, with Retired items excluded.

**Architecture:** Table-driven `UNION ALL` query in `services.GetAlerts` (no API change). Frontend gets a summary line and prettified severity labels in `AlertsModal.vue`. `useUiStore` exposes `expiryWarningDays` as a computed so `AppShell` can bind it to the modal.

**Tech Stack:** Go 1.25, `database/sql` + `modernc.org/sqlite`, `stretchr/testify`; Vue 3 + Pinia + TypeScript + Vitest + `@vue/test-utils`.

**Spec:** [`docs/superpowers/specs/2026-06-18-gover-alerts-parity-design.md`](../specs/2026-06-18-gover-alerts-parity-design.md)

---

## File Structure

| File | Responsibility | Change |
|------|----------------|--------|
| `gover/internal/services/services.go` | Read queries incl. `GetAlerts` | Rewrite `GetAlerts` body (table-driven `UNION ALL`); add `strings` import and `alertSource` type+var |
| `gover/internal/services/services_test.go` | Service-layer Go tests | Add 4 new tests for hardware warranty + Retired exclusion |
| `gover/frontend/src/stores/ui.ts` | Pinia UI store | Expose `expiryWarningDays` computed |
| `gover/frontend/src/stores/ui.test.ts` | UI store tests | Add one assertion for `expiryWarningDays` read-through |
| `gover/frontend/src/components/AlertsModal.vue` | Alerts modal | Add optional `warningDays` prop, summary line, prettified badge labels |
| `gover/frontend/src/components/AlertsModal.test.ts` | Modal tests | Add 2 tests: summary counts, prettified labels |
| `gover/frontend/src/layouts/AppShell.vue` | App shell | One-line binding: pass `ui.expiryWarningDays` to the modal |

---

## Task 1: Backend — extend `GetAlerts` to cover warranty expiry on hardware

**Files:**
- Modify: `gover/internal/services/services.go` (imports + new `alertSource` type/var + rewrite `GetAlerts` at lines 156–188)
- Test: `gover/internal/services/services_test.go` (append 4 new tests after the existing `TestGetAlerts_AlreadyExpired` at line 151)

- [ ] **Step 1.1: Write the 4 new failing tests**

Append to `gover/internal/services/services_test.go` immediately after `TestGetAlerts_AlreadyExpired` (currently ending at line 151):

```go
func TestGetAlerts_ComputerWarrantyExpiring(t *testing.T) {
	db := openTestDB(t)
	warranty := time.Now().AddDate(0, 0, 15).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO computers (name, status, warranty_expiry) VALUES (?, ?, ?)`, "Dell XPS", "Active", warranty)
	require.NoError(t, err)

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Computer", result[0].Category)
	assert.Equal(t, "expiring", result[0].Severity)
	assert.Equal(t, "Dell XPS", result[0].Name)
}

func TestGetAlerts_SmartphoneWarrantyExpired(t *testing.T) {
	db := openTestDB(t)
	warranty := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO smartphones (name, status, warranty_expiry) VALUES (?, ?, ?)`, "iPhone", "Active", warranty)
	require.NoError(t, err)

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Smartphone", result[0].Category)
	assert.Equal(t, "expired", result[0].Severity)
}

func TestGetAlerts_TabletWarrantyCovered(t *testing.T) {
	db := openTestDB(t)
	warranty := time.Now().AddDate(0, 0, 5).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO tablets (name, status, warranty_expiry) VALUES (?, ?, ?)`, "iPad", "Active", warranty)
	require.NoError(t, err)

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Tablet", result[0].Category)
	assert.Equal(t, "expiring", result[0].Severity)
}

func TestGetAlerts_RetiredExcluded(t *testing.T) {
	db := openTestDB(t)
	expired := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO antivirus (name, status, expiry_date) VALUES (?, ?, ?)`, "Old AV", "Retired", expired)
	require.NoError(t, err)

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	assert.Empty(t, result)
}
```

- [ ] **Step 1.2: Run the new tests — confirm they fail**

Run: `cd gover && go test ./internal/services/ -run TestGetAlerts -v`

Expected: 4 FAIL. `TestGetAlerts_ComputerWarrantyExpiring`, `_SmartphoneWarrantyExpired`, `_TabletWarrantyCovered` fail with `expected 1 alert, got 0` (current `GetAlerts` doesn't query hardware tables). `TestGetAlerts_RetiredExcluded` fails with `expected empty, got 1` (current query has no `status != 'Retired'` filter). The 3 pre-existing tests (`ExpiringWithin30Days`, `NotExpiringSoon`, `AlreadyExpired`) still PASS.

- [ ] **Step 1.3: Rewrite `GetAlerts` and update imports**

In `gover/internal/services/services.go`:

1. Replace the import block at lines 3–8 with:

```go
import (
	"database/sql"
	"fmt"
	"strings"

	"gover/internal/models"
)
```

2. Replace the entire `GetAlerts` function (currently lines 156–188) with:

```go
// alertSources lists every (table, date column) pair that contributes to
// expiry/warranty alerts. Mirrors db/alerts.py: Windows Key is intentionally
// omitted (no expiry column); Users is irrelevant.
type alertSource struct {
	table    string
	category string
	dateCol  string
}

var alertSources = []alertSource{
	{"computers", "Computer", "warranty_expiry"},
	{"smartphones", "Smartphone", "warranty_expiry"},
	{"tablets", "Tablet", "warranty_expiry"},
	{"antivirus", "Antivirus", "expiry_date"},
	{"other_software", "Other Software", "expiry_date"},
}

func GetAlerts(db *sql.DB, warningDays int) ([]models.Alert, error) {
	threshold := fmt.Sprintf("+%d days", warningDays)
	parts := make([]string, 0, len(alertSources))
	params := make([]any, 0, len(alertSources))
	for _, src := range alertSources {
		parts = append(parts, fmt.Sprintf(`
		SELECT '%s' AS category, id, COALESCE(name, '') AS name, %s AS expiry_date,
		       CAST(julianday(%s) - julianday('now') AS INTEGER) AS days_remaining
		FROM %s
		WHERE %s IS NOT NULL
		  AND date(%s) <= date('now', ?)
		  AND status != 'Retired'`,
			src.category, src.dateCol, src.dateCol, src.table, src.dateCol, src.dateCol))
		params = append(params, threshold)
	}
	query := strings.Join(parts, "\n UNION ALL \n") + "\n ORDER BY days_remaining ASC"

	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Alert
	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(&a.Category, &a.ID, &a.Name, &a.ExpiryDate, &a.DaysRemaining); err != nil {
			return nil, err
		}
		if a.DaysRemaining < 0 {
			a.Severity = "expired"
		} else {
			a.Severity = "expiring"
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
```

- [ ] **Step 1.4: Run the full services test suite — confirm everything passes**

Run: `cd gover && go test ./internal/services/ -v`

Expected: PASS. All 7 `TestGetAlerts_*` tests pass (3 pre-existing + 4 new), and no other service tests regress.

- [ ] **Step 1.5: Run the full Go suite + build — confirm no wider breakage**

Run: `cd gover && go test ./... && go build ./...`

Expected: PASS. `81 + 4 = 85` tests pass across 7 packages; build succeeds.

- [ ] **Step 1.6: Commit**

```bash
git add gover/internal/services/services.go gover/internal/services/services_test.go
git commit -m "feat(gover): alerts parity — cover hardware warranty + exclude Retired

GetAlerts now builds a UNION ALL across all 5 (table, date_col) pairs
(computers/smartphones/tablets on warranty_expiry, antivirus/other_software
on expiry_date), filters out Retired items, and orders by days_remaining
ascending so the most overdue row appears first. Mirrors db/alerts.py."
```

---

## Task 2: Frontend — expose `expiryWarningDays` from the UI store

**Files:**
- Modify: `gover/frontend/src/stores/ui.ts` (imports + computed + return)
- Test: `gover/frontend/src/stores/ui.test.ts` (append one assertion to existing `loadFromConfig` test)

- [ ] **Step 2.1: Add a failing assertion to `ui.test.ts`**

In `gover/frontend/src/stores/ui.test.ts`, find the existing test `'loadFromConfig applies darkMode + density and sets the dark class'` (lines 24–30) and extend it to also assert `expiryWarningDays`. Replace the test body with:

```go
  it('loadFromConfig applies darkMode + density and sets the dark class', async () => {
    const ui = useUiStore()
    await ui.loadFromConfig()
    expect(ui.darkMode).toBe(true)
    expect(ui.density).toBe('compact')
    expect(ui.expiryWarningDays).toBe(30)
    expect(document.documentElement.classList.contains('dark')).toBe(true)
  })
```

- [ ] **Step 2.2: Run the test — confirm it fails**

Run: `cd gover/frontend && pnpm run test -- --run ui.test`

Expected: FAIL with `TypeError: ui.expiryWarningDays is undefined` (or similar — the store doesn't currently expose this property).

- [ ] **Step 2.3: Expose `expiryWarningDays` from the store**

In `gover/frontend/src/stores/ui.ts`:

1. Update the import on line 2 from `import { ref } from 'vue'` to:

```ts
import { ref, computed } from 'vue'
```

2. Immediately after the `config` ref declaration (currently ending at line 14), add:

```ts
  const expiryWarningDays = computed(() => config.value.expiryWarningDays)
```

3. In the returned object (currently lines 41–52), add `expiryWarningDays,` to the list. The final return block should read:

```ts
  return {
    density,
    darkMode,
    sidebarCollapsed,
    dbVersion,
    expiryWarningDays,
    applyTheme,
    loadFromConfig,
    toggleDarkMode,
    toggleSidebar,
    setDensity,
    bumpDbVersion,
  }
```

- [ ] **Step 2.4: Run the test — confirm it passes**

Run: `cd gover/frontend && pnpm run test -- --run ui.test`

Expected: PASS. All 4 ui store tests pass including the new assertion.

- [ ] **Step 2.5: Commit**

```bash
git add gover/frontend/src/stores/ui.ts gover/frontend/src/stores/ui.test.ts
git commit -m "feat(gover): expose expiryWarningDays from useUiStore

Read-through computed so consumers (AlertsModal via AppShell) can render the
configured warning window without coupling to the config object shape."
```

---

## Task 3: Frontend — `AlertsModal` summary line + prettified labels + `warningDays` prop

**Files:**
- Modify: `gover/frontend/src/components/AlertsModal.vue` (template + script)
- Test: `gover/frontend/src/components/AlertsModal.test.ts` (append 2 new tests)

- [ ] **Step 3.1: Add 2 new failing tests to `AlertsModal.test.ts`**

Append to `gover/frontend/src/components/AlertsModal.test.ts` (after the existing `'emits close when Close is clicked'` test at line 33):

```ts
  it('shows expired/expiring summary counts and the warning window', () => {
    const wrapper = mount(AlertsModal, {
      props: { open: true, alerts, warningDays: 30 },
    })
    // 1 expired (Norton, severity: 'expired') + 1 expiring (Office, severity: 'expiring')
    expect(wrapper.text()).toMatch(/1\s+expired/)
    expect(wrapper.text()).toMatch(/1\s+expiring within 30 days/)
  })

  it('prettifies the severity badge labels', () => {
    const wrapper = mount(AlertsModal, { props: { open: true, alerts } })
    const cells = wrapper.findAll('tbody td:nth-child(5)')
    expect(cells).toHaveLength(2)
    expect(cells[0].text()).toBe('Expired')
    expect(cells[1].text()).toBe('Expiring soon')
  })
```

- [ ] **Step 3.2: Run the new tests — confirm they fail**

Run: `cd gover/frontend && pnpm run test -- --run AlertsModal`

Expected: 2 FAIL. Summary test fails (`does not match /1 expired/` — no summary line yet). Badge-label test fails (`cells[0].text()` returns `"expired"`, not `"Expired"`). The 4 pre-existing tests still PASS.

- [ ] **Step 3.3: Rewrite `AlertsModal.vue` with the new template + script**

Replace the entire contents of `gover/frontend/src/components/AlertsModal.vue` with:

```vue
<template>
  <div v-if="open" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
    <div class="bg-surface border border-border p-6 w-[560px] max-h-[70vh] overflow-y-auto">
      <h2 class="text-base font-semibold text-text-primary mb-2">Alerts</h2>

      <p v-if="alerts.length > 0" class="text-sm text-text-secondary mb-4">
        <b class="text-text-primary">{{ expiredCount }}</b> expired
        &nbsp;·&nbsp;
        <b class="text-text-primary">{{ expiringCount }}</b> expiring within {{ warningDays }} days
      </p>

      <p v-else class="text-text-secondary text-sm">No alerts.</p>

      <table v-if="alerts.length > 0" class="w-full text-sm">
        <thead class="bg-th-bg">
          <tr>
            <th class="text-left px-3 py-2 text-th-text font-medium">Category</th>
            <th class="text-left px-3 py-2 text-th-text font-medium">Name</th>
            <th class="text-left px-3 py-2 text-th-text font-medium">Expiry</th>
            <th class="text-left px-3 py-2 text-th-text font-medium">Days</th>
            <th class="text-left px-3 py-2 text-th-text font-medium">Severity</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="a in alerts" :key="a.category + '-' + a.id" class="border-t border-border">
            <td class="px-3 py-2 text-text-primary">{{ a.category }}</td>
            <td class="px-3 py-2 text-text-primary">{{ a.name }}</td>
            <td class="px-3 py-2 text-text-primary">{{ a.expiryDate }}</td>
            <td class="px-3 py-2 text-text-primary">{{ a.daysRemaining }}</td>
            <td class="px-3 py-2">
              <span
                class="px-2 py-0.5 rounded-badge text-xs text-white"
                :class="a.severity === 'expired' ? 'bg-s-expired' : 'bg-s-expiring'"
              >
                {{ severityLabel[a.severity] ?? a.severity }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>

      <button
        class="mt-4 px-3 py-1.5 text-sm bg-s-active text-white hover:bg-accent-hover"
        @click="$emit('close')"
      >
        Close
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Alert } from '@/lib/api'

const props = withDefaults(defineProps<{
  open: boolean
  alerts: Alert[]
  warningDays?: number
}>(), { warningDays: 30 })

defineEmits<{ (e: 'close'): void }>()

const severityLabel: Record<string, string> = {
  expired: 'Expired',
  expiring: 'Expiring soon',
}

const expiredCount = computed(() => props.alerts.filter(a => a.severity === 'expired').length)
const expiringCount = computed(() => props.alerts.length - expiredCount.value)
</script>
```

- [ ] **Step 3.4: Run the modal tests — confirm all pass**

Run: `cd gover/frontend && pnpm run test -- --run AlertsModal`

Expected: PASS. All 6 tests pass (4 pre-existing + 2 new). Existing tests:
- `'renders a row per alert when open'` — still finds `Norton` and `Office`, 2 rows. (Modal now also shows summary but that doesn't break this assertion.)
- `'shows empty message when alerts is empty'` — text `'No alerts.'` still rendered.
- `'renders nothing when closed'` — `v-if="open"` still gates the `.fixed` wrapper.
- `'emits close when Close is clicked'` — Close button still emits.

- [ ] **Step 3.5: Commit**

```bash
git add gover/frontend/src/components/AlertsModal.vue gover/frontend/src/components/AlertsModal.test.ts
git commit -m "feat(gover): AlertsModal summary line + prettified severity labels

Adds a 'X expired · Y expiring within N days' header (matches PyQt6
_show_alerts_dialog) and maps raw severity values to 'Expired' / 'Expiring
soon' for display. New optional warningDays prop defaults to 30."
```

---

## Task 4: Wire `AppShell` to pass `warningDays` to the modal

**Files:**
- Modify: `gover/frontend/src/layouts/AppShell.vue` (line 17 — one-line prop binding)

- [ ] **Step 4.1: Bind `:warning-days` to the modal**

In `gover/frontend/src/layouts/AppShell.vue` line 17, change:

```vue
    <AlertsModal :open="alertsOpen" :alerts="alerts" @close="alertsOpen = false" />
```

to:

```vue
    <AlertsModal
      :open="alertsOpen"
      :alerts="alerts"
      :warning-days="ui.expiryWarningDays"
      @close="alertsOpen = false"
    />
```

`ui` is already declared on line 37 (`const ui = useUiStore()`), and Task 2 made `ui.expiryWarningDays` available. No other changes needed.

- [ ] **Step 4.2: Run typecheck — confirm the binding typechecks**

Run: `cd gover/frontend && pnpm run typecheck`

Expected: PASS with no errors. `ui.expiryWarningDays` is now a `number` (from the Task 2 computed), matching the modal's optional `warningDays?: number` prop.

- [ ] **Step 4.3: Run the full frontend test suite — confirm nothing regresses**

Run: `cd gover/frontend && pnpm run test`

Expected: PASS. All frontend tests pass, including the existing `AppShell.test.ts` (which stubs `AlertsModal` and won't be affected by the new prop binding).

- [ ] **Step 4.4: Commit**

```bash
git add gover/frontend/src/layouts/AppShell.vue
git commit -m "feat(gover): pass expiryWarningDays from useUiStore to AlertsModal

One-line binding so the modal's summary line shows the configured warning
window (currently always 30) instead of the hardcoded default."
```

---

## Task 5: Final end-to-end verification

**Files:** None modified — verification only.

- [ ] **Step 5.1: Run the full Go suite + build**

Run: `cd gover && go test ./... && go build ./...`

Expected: PASS. 85 tests across 7 packages; build succeeds.

- [ ] **Step 5.2: Run the full frontend suite + typecheck**

Run: `cd gover/frontend && pnpm run test && pnpm run typecheck`

Expected: PASS. All Vitest tests green; `vue-tsc --noEmit` reports no errors.

- [ ] **Step 5.3: Confirm the four commit boundaries**

Run: `cd /home/mg/inventory && git log --oneline -4`

Expected: the four commits from Tasks 1–4 appear in order:
1. `feat(gover): alerts parity — cover hardware warranty + exclude Retired`
2. `feat(gover): expose expiryWarningDays from useUiStore`
3. `feat(gover): AlertsModal summary line + prettified severity labels`
4. `feat(gover): pass expiryWarningDays from useUiStore to AlertsModal`

No additional commit is produced by this task — verification only.

---

## Self-Review Notes

**Spec coverage check** (against `docs/superpowers/specs/2026-06-18-gover-alerts-parity-design.md`):

- ✅ Backend `UNION ALL` over 5 sources with `status != 'Retired'` → Task 1
- ✅ No model changes (`models.Alert` unchanged) → no task needed (confirmed in Task 1.3 scan: `&a.Category, &a.ID, &a.Name, &a.ExpiryDate, &a.DaysRemaining` all map to existing fields)
- ✅ No bridge changes → no task needed (Task 1.5 verifies the bridge still compiles via `go build ./...`)
- ✅ Frontend summary line + badge labels + `warningDays` prop → Task 3
- ✅ Expose `expiryWarningDays` from `useUiStore` → Task 2 (added during planning after spec erratum)
- ✅ Wire `AppShell` to pass `warningDays` → Task 4
- ✅ 4 new Go tests → Task 1.1
- ✅ 1 new Vitest for summary counts → Task 3.1 (plus the badge-label test, which the spec mentioned as "if cheap")
- ✅ Verification gates (`go test`, `go build`, `pnpm test`, `pnpm typecheck`) → Tasks 1.5, 2.4, 3.4, 4.2, 4.3, 5.1, 5.2

**Placeholder scan:** none. All code blocks contain real, copy-pasteable Go/Vue/TS. No "TBD", "TODO", "implement later", or "similar to Task N".

**Type/name consistency:**
- `severityLabel` map keys (`expired`, `expiring`) match the strings emitted by Go (`services.go` lines `"expired"`, `"expiring"`) and the existing `Alert.severity` values in `AlertsModal.test.ts` test data.
- `warningDays` prop name is consistent across `AlertsModal.vue` (prop definition), `AppShell.vue` binding (`:warning-days` — Vue's kebab-case attribute form auto-maps to camelCase prop), and the test (`warningDays: 30`).
- `expiryWarningDays` exposed name matches the existing `AppConfig.expiryWarningDays` field name (`lib/types/index.ts:92`).
- `alertSource` struct field names (`table`, `category`, `dateCol`) match the format-string placeholders in `GetAlerts`.

**Scope check:** single focused change, ~5 commits, single implementation session. No decomposition needed.
