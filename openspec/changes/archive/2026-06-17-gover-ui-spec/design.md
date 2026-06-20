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
| `--sidebar` | `#0E1520` | Sidebar background, table header background |
| `--page` | `#EEF1F7` | Page background |
| `--surface` | `#FFFFFF` | Cards, panels, topbar |
| `--accent` | `#2563EB` | Primary interactive, active nav border |
| `--warn` | `#D97706` | Expiring-soon highlights |
| `--danger` | `#DC2626` | Expired and destructive actions |

**Status strip colors:**

| Status | Color |
|---|---|
| Active | `#2563EB` |
| Expiring soon | `#D97706` |
| Expired / In Repair | `#DC2626` |
| Spare | `#7C3AED` |
| Retired | `#9CA3AF` |
| Missing | `#1E293B` |

**Typography:**

- UI chrome: system font stack (`-apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif`)
- Data values (IDs, license keys, dates): monospace stack (`'Cascadia Code', 'Fira Code', 'Consolas', monospace`)
- No external font imports; this app is local-only with no internet requirement

**Signature constraint:** Zero `border-radius` on structural containers (sidebar, topbar, table, edit panel). `border-radius: 2px` only on inline badges.

### Component surface (minimum required)

| Component | Notes |
|---|---|
| `AppShell` | Layout: sidebar + main area |
| `Sidebar` | Nav items with icon, label, count; section labels; footer |
| `Topbar` | Search input, Export CSV, Alerts badge, Dark mode toggle |
| `DataTable` | Status strip column, header, sortable columns, row hover/select, expiring/expired row tints |
| `StatusBadge` | Dot + label, 7 status variants |
| `EditPanel` | Slide-in right panel, header, form body, footer actions |
| `FormField` | Label + input (text, date, select, textarea) |
| `AlertsModal` | Expiry alerts table with status + name + type + date |
| `Statusbar` | Record count, sort state, DB path |

## Risks / Trade-offs

**Slide-in panel + table layout at small widths** → The edit panel occupies 300 px. At 1280 px (the minimum window size from the current app), the table has 780 px after sidebar and panel. This is acceptable for typical column counts but may cause horizontal scroll on the Windows Keys or Antivirus categories with license key columns. Mitigation: the panel can be toggled closed and the window minimum width kept at 1280 px per the current app.

**No font imports** → System font stacks render differently on Windows (Segoe UI) vs Linux (DejaVu/Noto). The design was built with system fonts intentionally, but spacing and weight rendering may vary slightly. Mitigation: use `font-synthesis: none` and test on both platforms.

**Fixed categories** → The sidebar is hard-coded to 7 categories. If configurable categories are added later, the nav items must become data-driven. The design anticipates this (section labels and counts are parameterised in the HTML reference), but the Vue component will need to be written with that future-proofing in mind.

## Open Questions

- Should the "All" overview tab show only hardware (Computer, Smartphone, Tablet) as in the current app, or include license categories? The HTML reference shows hardware-only in the sidebar overview. This decision should be confirmed before `DataTable` is implemented.
- Keyboard shortcut for toggling the edit panel — no current equivalent in PyQt6. Suggest `Ctrl+E` (consistent with "Edit" menu item), but needs confirmation to avoid clash with Export CSV (`Ctrl+E` in the current app's menu bar).
