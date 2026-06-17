# Comet Design Handoff

- Change: gover-ui-spec
- Phase: design
- Mode: compact
- Context hash: eded2a4130c4b02f8fc050338cbe0f9c2a320b6e4282f2cba9120945766a3ceb

Generated-by: comet-handoff.sh

OpenSpec remains the canonical capability spec. This handoff is a deterministic, source-traceable context pack, not an agent-authored summary.

## openspec/changes/gover-ui-spec/proposal.md

- Source: openspec/changes/gover-ui-spec/proposal.md
- Lines: 1-29
- SHA256: 397bd8a4ed33c3989e5fef56a77fe6e8183bddbe99b7e9327c74292402ae74aa

```md
## Why

The Gover v2 rewrite plan (`gover/QUESTIONS.md`) leaves the "UI And UX" section unresolved, blocking Phase 1 architecture work. The visual design reference (`gover/UI_DESIGN_REFERENCE.html`) was produced during the planning session and contains implicit decisions that must be made explicit before the Vue/Tailwind frontend can be specified or implemented.

## What Changes

- Documents definitive answers to every question in the "UI And UX" section of `gover/QUESTIONS.md`
- Formally adopts the design system shown in the HTML reference: dark sidebar, cobalt accent, status-strip table rows, slide-in edit panel
- Answers adjacent asset-model questions whose answers are constrained by the chosen UI (fixed vs. configurable categories)
- Produces a UI design spec (`specs/ui-decisions/spec.md`) that Phase 1 technical planning can treat as the UI contract

No application code is written. The current Python app is not modified.

## Capabilities

### New Capabilities

- `ui-decisions`: Resolved answers to all UI/UX questions from `gover/QUESTIONS.md`, including navigation model, first-screen choice, add/edit UX pattern, keyboard-first posture, density modes, dark mode strategy, and any asset-model questions that directly constrain the UI.

### Modified Capabilities

_(none — no existing specs to modify)_

## Impact

- `gover/QUESTIONS.md`: UI/UX section items become answered (document updated or annotated)
- `gover/UI_DESIGN_REFERENCE.html`: referenced as primary evidence; no changes to the file
- Unblocks: Phase 1 technical foundation planning (architecture, Tailwind config strategy, component naming)
- Does not affect: current `app/` Python code, `db/` layer, packaging, CI
```

## openspec/changes/gover-ui-spec/design.md

- Source: openspec/changes/gover-ui-spec/design.md
- Lines: 1-132
- SHA256: 7aadae266a24895e29e8e83a8d047ac80af62bbd3431a68308904fccecf98110

[TRUNCATED]

```md
## Context

The Gover v2 rewrite will use Wails + Go + Vue 3 + TypeScript + Tailwind CSS. Before writing any application code, the UI/UX questions in `gover/QUESTIONS.md` must be answered so that the Vue component tree, routing, and Tailwind config can be designed correctly in Phase 1.

`gover/UI_DESIGN_REFERENCE.html` was produced as part of planning and encodes a complete visual design system. This design document formalises the decisions implicit in that reference and fills in gaps (density, dark mode strategy) not shown in the HTML mockup.

## Goals / Non-Goals

**Goals:**

- Answer every question in the "UI And UX" section of `gover/QUESTIONS.md` with a definitive choice and rationale
- Answer adjacent asset-model questions whose outcomes constrain the UI
- Specify the design token system (colors, typography, spacing) as Tailwind CSS variables
- Define the component surface area at a high level (which components must exist, not their implementation)

**Non-Goals:**

- Writing Vue, TypeScript, Go, or Tailwind code
- Answering backend, packaging, testing, or data-migration questions
- Pixel-perfect wireframes for every view (the HTML reference covers the primary view)
- Defining Wails bridge method signatures

## Decisions

### Navigation model: sidebar replaces tab bar

**Decision:** Use a persistent left sidebar (200 px, collapsible to icon-only at ≤ 1024 px) instead of the current QTabBar.

**Rationale:** The current tab bar works for 7–8 categories but is a horizontal strip that competes with column space on wide tables. A sidebar lets category navigation grow (configurable categories become possible later), displays record counts per section, and frees the horizontal space entirely for data columns.

**Alternative considered:** Keep tabs across the top. Rejected because the tab bar disables "Add" on the "All" tab by convention, has no count indicators, and does not scale past ~10 categories.

### First screen: Computers tab (not a dashboard)

**Decision:** The default landing view is the Computers table — the most frequently accessed category for a typical IT admin.

**Rationale:** A dashboard adds design complexity and a new code surface before any CRUD is stable. The sidebar shows counts for all categories, giving the ambient overview a dashboard would provide. This can be revisited in Phase 4.

**Alternative considered:** An overview dashboard with counts and alert summary. Deferred: will add value once Phase 3 (CRUD) is stable and alert data is populated.

### Add / edit UX: slide-in panel, not modal dialogs

**Decision:** Selecting a row opens a slide-in panel on the right (300 px). "Add" opens the same panel in an empty state. No modal dialogs.

**Rationale:** Modals block the table completely and require the user to remember context from a record they can no longer see. The panel keeps the table visible and row selected, which is essential when editing to verify context (e.g., "is this the right model?"). The HTML reference shows this layout.

**Alternative considered:** Full-page edit route. Rejected for a local desktop app where context switching between list and form is jarring and routing adds complexity for no network benefit.

**Alternative considered:** Keep modal dialogs matching the current PyQt6 behavior. Rejected: a web-rendered panel is a net improvement for all editing tasks, and modals in web UIs are harder to make keyboard-navigable.

### Keyboard-first posture: maintain current shortcuts, extend later

**Decision:** Preserve all existing keyboard shortcuts (Ins = Add, Enter/Return = Edit, Del = Delete, Ctrl+F = focus search, F5 = check alerts, Esc = clear search). Add Tab/arrow navigation within the edit panel. Treat this as a baseline, not a full command-palette implementation.

**Rationale:** The current PyQt6 app has well-established shortcuts that power users depend on. Maintaining them reduces the transition cost. A full command palette (Cmd+K) is a Phase 4 stretch goal.

### Density: comfortable by default, compact as a setting

**Decision:** Default row height is 33 px (comfortable). A "Compact" toggle in Settings reduces row height to 24 px (matching the current PyQt6 default). The setting persists in app config.

**Rationale:** The design reference uses 33 px rows. Most IT admins work at 1920×1080 and benefit from the breathing room. Users who prefer the current app's density can switch.

### Dark mode: manual toggle, system default Light

**Decision:** Dark mode is a manual toggle (button in topbar, persisted in app config). Default is Light. System-theme following is a future enhancement.

**Rationale:** The current app already has a manual toggle and users know where it is. Implementing system-theme following correctly requires Wails platform event integration that is better addressed in Phase 4 alongside other OS-level features. The manual toggle has zero implementation risk.

### Asset categories: fixed for this change

**Decision:** Categories remain the fixed set defined in `db/config.py`: Computer, Smartphone, Tablet, Windows Key, Antivirus, Other Software, Users. Configurable categories are out of scope for this change.

**Rationale:** Configurable categories require a different data model (categories table + foreign keys instead of typed tables) which is a schema migration decision. That decision should be made in data/migration planning, not UI planning. The sidebar is designed to accommodate it later (section labels + nav items are dynamically rendered).

### Design token strategy

**Color palette (6 named values):**

| Token | Hex | Role |
|---|---|---|
```

Full source: openspec/changes/gover-ui-spec/design.md

## openspec/changes/gover-ui-spec/tasks.md

- Source: openspec/changes/gover-ui-spec/tasks.md
- Lines: 1-30
- SHA256: 5c28b9a154747ab2625da77c5b5111c74a89c3b61eb33c6f83b5b8bcfaa42ef8

```md
## 1. Answer UI/UX Questions

- [ ] 1.1 Read every question in the "UI And UX" section of `gover/QUESTIONS.md` and verify each has a matching decision in `specs/ui-decisions/spec.md`
- [ ] 1.2 Annotate `gover/QUESTIONS.md` — mark each UI/UX question as resolved with a reference to the spec
- [ ] 1.3 Confirm the two open questions in `design.md` (All tab scope, Ctrl+E conflict) are resolved or deferred with a documented rationale

## 2. Validate Design Token System

- [ ] 2.1 Review all color values in `specs/ui-decisions/spec.md` against `gover/UI_DESIGN_REFERENCE.html` for consistency
- [ ] 2.2 Verify status strip colors cover all `STATUS_OPTIONS` from `db/config.py` (Active, In Repair, Spare, Retired, Missing) plus expiring/expired states
- [ ] 2.3 Confirm font stacks render acceptably on both Linux (WSL/Ubuntu) and Windows — note any adjustments needed

## 3. Resolve Adjacent Asset-Model Questions

- [ ] 3.1 Document the decision on fixed vs. configurable categories in `specs/ui-decisions/spec.md` (already drafted — confirm with stakeholder)
- [ ] 3.2 Document the "All" tab scope decision (hardware only vs. all categories)
- [ ] 3.3 Review remaining asset-model questions in `gover/QUESTIONS.md` and flag any that directly affect the UI (serial numbers, locations, tags, etc.)

## 4. Finalize the Design Reference

- [ ] 4.1 Add a Compact density mockup to `gover/UI_DESIGN_REFERENCE.html` (33 px → 24 px rows)
- [ ] 4.2 Add a Dark mode skeleton to `gover/UI_DESIGN_REFERENCE.html` (sidebar accent colors, surface/page inversion)
- [ ] 4.3 Add the Alerts overlay view to `gover/UI_DESIGN_REFERENCE.html`
- [ ] 4.4 Verify the HTML reference renders correctly in Chromium (the Wails WebView engine on Linux)

## 5. Documentation Hand-Off

- [ ] 5.1 Update `gover/QUESTIONS.md` header to note that the UI/UX section has been resolved via this change
- [ ] 5.2 Add a link to `specs/ui-decisions/spec.md` and `gover/UI_DESIGN_REFERENCE.html` from `gover/ARCHITECTURE_NOTES.md` under "UI Strategy"
- [ ] 5.3 Confirm no application code was written during this change (verify git diff shows only docs/html changes)
```

## openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md

- Source: openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md
- Lines: 1-260
- SHA256: b8bb5a73a4e73802a1ae38fd48e4b33cdc7117dc6ebccc047f678c6bddeb40e3

[TRUNCATED]

```md
# UI Decisions Spec — Gover v2

This spec documents the resolved answers to the "UI And UX" section of `gover/QUESTIONS.md` and adjacent asset-model questions whose answers are constrained by the UI. All decisions reference `gover/UI_DESIGN_REFERENCE.html` as the visual source of truth.

---

## Navigation Model

**Q: Should the current tab model be preserved?**
**A: No.** Replace the horizontal `QTabBar` with a persistent left sidebar.

**Q: Should the new UI use a sidebar, dashboard, command palette, or another navigation model?**
**A: Sidebar.** Fixed width 200 px, dark background (`#0E1520`), with category nav items showing icon + label + record count. Collapses to icon-only at window widths ≤ 1024 px.

Sidebar structure:
```
OVERVIEW
  All Assets           [count]

HARDWARE
  Computers            [count]
  Smartphones          [count]
  Tablets              [count]

LICENSES
  Windows Keys         [count]
  Antivirus            [count]
  Other Software       [count]

─────────────────────────────
  Users                [count]

─────────────────────────────
  [version] · [db filename]
```

Active item has a 2 px left-edge accent (`#2563EB`) and a slightly lighter background (`#182337`). Inactive items use muted text (`#6B82A0`).

---

## First Screen

**Q: Should the first screen be an overview dashboard or the asset table?**
**A: Asset table (Computers).** The default landing view is the Computers category. The sidebar counts provide ambient overview of all categories without requiring a dedicated dashboard view. A dashboard may be added in Phase 4 once CRUD is stable.

---

## Add / Edit UX Pattern

**Q: Should add/edit use modal dialogs, side panels, or full pages?**
**A: Slide-in right panel.** Width 300 px, attached to the right edge of the window. Selecting a row opens the panel in edit mode. Clicking "Add" opens the panel in create mode with empty fields.

The panel remains open and the selected row stays highlighted while editing, so the user can verify context from the table. The panel closes on Save, Cancel, or pressing Escape.

Panel anatomy:
```
┌─────────────────────────────┐
│ Edit Computer          [×]  │  ← header: category + close
│─────────────────────────────│
│ NAME                        │
│ [WKSTN-MGLOW              ] │
│ MODEL                       │
│ [MacBook Pro 16"          ] │
│ ASSIGNED USER               │
│ [Głowacki, M.       ▾    ] │
│ STATUS                      │
│ [Active              ▾    ] │
│ PURCHASE DATE               │
│ [2023-04-12               ] │
│ WARRANTY EXPIRY             │
│ [2026-04-12               ] │
│ NOTES                       │
│ [Main dev machine         ] │
│                             │
│─────────────────────────────│
│ [   Save   ] [Cancel] [Del] │  ← footer
└─────────────────────────────┘
```

Delete in the panel footer shows a confirmation inline (not a separate modal). Destructive state: button turns `#DC2626`, label changes to "Confirm delete?".
```

Full source: openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md
