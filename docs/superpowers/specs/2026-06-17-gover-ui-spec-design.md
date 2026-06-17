---
comet_change: gover-ui-spec
role: technical-design
canonical_spec: openspec
---

# Gover v2 UI Design — Technical Design Doc

## Context

The `gover-ui-spec` change resolves all "UI And UX" open questions from `gover/QUESTIONS.md` so that Phase 1 (Wails/Vue/Tailwind technical foundation) can proceed with a confirmed UI contract.

The visual design reference (`gover/UI_DESIGN_REFERENCE.html`) was produced in the same session and shows the primary table view with edit panel. This document formalises the decisions implicit in that reference, resolves two previously-open questions, and adds one new decision (native menu bar) discovered during design.

Canonical spec: `openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md`

---

## Goals / Non-Goals

**Goals:**
- Resolve every UI/UX question in `gover/QUESTIONS.md` with a documented decision and rationale
- Add a "Native Menu Bar" decision to the spec (discovered gap)
- Extend `UI_DESIGN_REFERENCE.html` with compact density, dark mode, and alerts overlay views
- Annotate `QUESTIONS.md` and link from `ARCHITECTURE_NOTES.md`

**Non-Goals:**
- Writing Vue, TypeScript, Go, or Tailwind application code
- Answering backend, packaging, testing, or data-migration questions
- Pixel-perfect wireframes for every possible state

---

## Decisions

### 1. Navigation: persistent left sidebar

Replaces the `QTabBar`. Width 200 px, dark background (`#0E1520`), collapses to icon-only at ≤ 1024 px. Rationale: sidebar scales beyond 10 categories, shows record counts, and frees the full horizontal space for data columns.

### 2. First screen: Computers table

Default landing view is Computers (the most frequently accessed category). Dashboard deferred to Phase 4 — the sidebar counts provide ambient overview.

### 3. Add/edit: slide-in right panel

Width 300 px. Opens on row selection (edit mode) or Add button (create mode). Table remains visible. Delete confirmation inline in panel footer. No modal dialogs. Rationale: modals lose table context; panel keeps the selected row in view.

### 4. Keyboard shortcuts: preserve current + panel extensions

All PyQt6 shortcuts preserved: `Ins` (Add), `Enter` (Edit/open panel), `Del` (Delete), `Ctrl+F` (Search), `F5` (Alerts), `Esc` (Clear search / close panel). Panel adds: `Tab` (next field), `Ctrl+S` (Save), `Esc` (close).

No shortcut for "toggle edit panel" — row selection handles open, `Esc` handles close. `Ctrl+E` remains Export CSV (unchanged from current app).

### 5. Density: comfortable (default) + compact

Comfortable = 33 px row height (default). Compact = 24 px (matches current PyQt6 app). Controlled by toggle in Settings. Persists to app config.

### 6. Dark mode: manual toggle, Light default

Toggle in topbar, persists to app config. System-theme following (`wails/v2` `RegisterThemeChangedHook`) deferred to Phase 4.

### 7. Asset categories: fixed for v2

7 categories hard-coded (Computer, Smartphone, Tablet, Windows Key, Antivirus, Other Software, Users). Sidebar nav items rendered from a data source (not Vue template literals) to allow future configurability without a component rewrite.

### 8. Native OS menu bar via Wails

**New decision (not in the original question set).** The Wails app uses `wails/v2/pkg/menu` to attach a native OS menu bar matching the current PyQt6 app's structure:

```
File
  New Database…       Ctrl+N
  Open Database…      Ctrl+O
  ────────────────────────
  Export CSV…         Ctrl+E
  ────────────────────────
  Quit                Ctrl+Q

Help
  About
  ────────────────────────
  License
```

Rationale: these actions involve OS-level file dialogs (`wailsv2.OpenFileDialog`, `SaveFileDialog`) which feel most natural invoked from a native menu. The in-app topbar carries Export CSV as a secondary access point for discoverability.

### 9. "All" overview tab scope

Hardware only: Computer, Smartphone, Tablet — matching the current `ALL_TAB_CATS`. License categories excluded (different column schema, no `user_id` / `warranty_expiry`). Columns: Type, ID, Name, Model, User, Status, Purchase Date, Warranty Expiry, Notes.

---

## Design Token Summary

Full token values live in `specs/ui-decisions/spec.md`. Key entries:

| Token | Hex | Use |
|---|---|---|
| `--sidebar-bg` | `#0E1520` | Sidebar + table header background |
| `--page` | `#EEF1F7` | Page background |
| `--surface` | `#FFFFFF` | Cards, panels, topbar |
| `--accent` | `#2563EB` | Active nav, primary button, active status strip |
| `--warn` | `#D97706` | Expiring-soon strip and row tint (`#FFFBEB`) |
| `--danger` | `#DC2626` | Expired strip and row tint (`#FFF5F5`), destructive |

Table header uses `--sidebar-bg` — the signature element that makes the header feel structurally joined to the sidebar.

Typography: system UI stack for chrome, monospace stack for data values (IDs, license keys, dates). No external font imports.

Border radius: 0 on structural containers, 2 px on badges only.

---

## Component Surface (Phase 2 minimum)

| Component | Role |
|---|---|
| `AppShell.vue` | Root layout flex: sidebar + main |
| `Sidebar.vue` | Nav items, section labels, counts, footer |
| `SidebarItem.vue` | Icon + label + count, active state |
| `Topbar.vue` | Search, Export CSV, Alerts badge, Dark mode toggle |
| `DataTable.vue` | Status strip, sortable headers, row states, keyboard nav |
| `StatusBadge.vue` | Dot + label, 7 status variants |
| `EditPanel.vue` | Slide-in panel, form body, footer (Save/Cancel/Delete) |
| `FormField.vue` | Label + control (text, date, select, textarea) |
| `AlertsModal.vue` | Expiry alerts overlay table |
| `Statusbar.vue` | Record count, sort state, DB path |

---

## Risks / Trade-offs

**Wails native menu API** — `wails/v2/pkg/menu` requires Go menu registration at startup. Menu callbacks invoke Go bridge methods which must be defined before the menu is registered. This is a Phase 1 architecture concern (bridge method naming and registration order). No impact on Vue component design.

**Slide-in panel at minimum window width** — At 1280 px (app minimum), sidebar (200 px) + panel (300 px) leaves 780 px for the table. Adequate for most category columns, but license key columns (Windows Keys, Antivirus) may need horizontal scroll or column truncation. Mitigated by collapsing the sidebar to icon-only mode, reclaiming ~140 px.

**System font rendering** — No font imports. On Linux WSL, `system-ui` falls back to DejaVu or Noto, which has slightly different metrics than Segoe UI (Windows). Spacing tested with 33 px row height; compact 24 px may be tighter. Annotation: check both platforms when adding compact density to the HTML reference.

---

## Execution Order (Approach A)

1. Apply Spec Patch: add native menu bar section to `specs/ui-decisions/spec.md`
2. Verify spec covers all QUESTIONS.md UI/UX items (checklist in tasks.md)
3. Annotate `gover/QUESTIONS.md` — resolved items reference spec section
4. Add HTML views to `UI_DESIGN_REFERENCE.html`: compact, dark mode, alerts
5. Verify HTML renders in Chromium (WSL)
6. Link from `gover/ARCHITECTURE_NOTES.md`
7. Confirm no application code in git diff
