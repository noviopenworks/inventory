# Brainstorm Summary

- Change: gover2-crud
- Date: 2026-06-18

## Confirmed Technical Approach

Layer-by-layer write operations following the same dependency chain as Phase 2:
`services → bridge → api → components → views`.

**Go backend:**
- 7 Input structs (`ComputerInput` etc.) + `DropdownItem` + `DeviceDropdownItem` in `internal/models/models.go`
- `Insert*/Update*/Delete*` per category in `internal/services/` (existing files extended)
- `ListUsersForDropdown` and `ListDevicesForDropdown` in `internal/services/dropdowns.go`
  - `ListDevicesForDropdown`: 3-way UNION ALL (computers, smartphones, tablets) with literal kind column
- `updated_at` set via `datetime('now','localtime')` in SQL, never from client
- 23 new bridge methods (21 write + 2 dropdown helpers) with nil-db guard; bridge does NOT refetch after writes
- `AddComputer` (not `InsertComputer`) naming convention at bridge layer

**Frontend:**
- `EditPanel.vue`: existing w-80 stub replaced with full implementation
  - `FIELD_CONFIGS` map drives fields per category (mirrors Python TABLE_CONFIG)
  - `FormField.vue` (existing) used as field wrapper
  - Device FK: single mixed `<select>` from `devices` prop, `"Name (type)"` labels
  - On save: extract kind+id from selection → set correct FK column in payload
  - Required field: `name` for all categories except WindowsKeys (required: `license_key`)
  - Date fields: text input, YYYY-MM-DD placeholder, validated on save
  - Escape key → emits `cancelled`
  - 200ms ease-in-out slide-in transition
- `ConfirmDialog.vue`: new centered modal overlay, `message` prop, emits `confirmed`/`cancelled`
- All 7 views: load users+devices dropdowns on mount; row-click → edit; "Add" → create; delete button → ConfirmDialog → hard delete; after `saved` refetch list

## Key Trade-offs and Risks

- **Generic vs per-category panels**: chose generic EditPanel with FIELD_CONFIGS map (same as Python app). Risk: FIELD_CONFIGS map is verbose but keeps view-wiring code DRY.
- **Device FK mapping**: single mixed dropdown resolves to correct FK column on save. Edge case: if a device is deleted after the panel opens, the payload will send a stale FK — acceptable for Phase 3 scope.
- **No optimistic updates**: view refetches after save rather than mutating local state. Simpler and avoids stale-data bugs.
- **Bridge does not refetch**: frontend is responsible for all list refreshes. Consistent with existing read pattern.

## Testing Strategy

**Go (table-driven, in-memory SQLite):**
- `services_test.go`: for each category — insert, verify in list; update, verify field changed; delete, verify list empty
- `bridge_test.go`: `AddComputer` round-trip; `DeleteComputer` removes record; nil-db guard returns error
- Coverage: ≥70% services, ≥50% bridge

**Frontend (Vitest + Vue Test Utils):**
- `EditPanel.test.ts`: edit mode pre-fills fields; create mode fields empty; Cancel emits `cancelled`; Escape emits `cancelled`
- `ConfirmDialog.test.ts`: Confirm emits `confirmed`; Cancel emits `cancelled`

## Spec Patches

None — OpenSpec delta spec does not require supplementation.
