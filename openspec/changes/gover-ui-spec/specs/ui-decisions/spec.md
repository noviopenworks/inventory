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

---

## Keyboard-First Operation

**Q: How important is keyboard-first operation?**
**A: Important — preserve existing shortcuts and extend at panel level.**

Keyboard shortcuts (baseline, matching current PyQt6 app):

| Key | Action |
|---|---|
| `Ins` | Add new record (current category) |
| `Enter` / `Return` | Edit selected record (opens panel) |
| `Del` | Delete selected record |
| `Ctrl+F` | Focus search input |
| `F5` | Check alerts |
| `Esc` | Clear search field / close panel |
| `↑` / `↓` | Move row selection |

Additional shortcuts at panel level:

| Key | Action |
|---|---|
| `Tab` | Move between form fields |
| `Ctrl+S` | Save form |
| `Esc` | Close panel without saving |

Command palette (`Cmd+K` / `Ctrl+K`) is deferred to Phase 4.

---

## Density Modes

**Q: Should the UI support compact and comfortable density modes?**
**A: Yes, two modes: Comfortable (default) and Compact.**

| Mode | Row height | When to use |
|---|---|---|
| Comfortable | 33 px | Default; better for wide monitors |
| Compact | 24 px | Matches current app; useful at smaller resolutions |

Controlled by a toggle in Settings. Persisted in app config. No live preview required — change applies on next navigation or page refresh within the app.

---

## Dark Mode

**Q: Should dark mode follow the system theme or remain a manual toggle?**
**A: Manual toggle.** Default is Light. The toggle is in the topbar and persists to app config.

System-theme following (Wails `window.RegisterThemeChangedHook`) is deferred to Phase 4. The manual toggle matches the current app behavior, so no user habits break on migration.

---

## Asset Categories (UI-constrained decision)

**Q: Should asset categories remain fixed or become configurable?**
**A: Fixed for Gover v2.** The 7 categories (`Computer`, `Smartphone`, `Tablet`, `Windows Key`, `Antivirus`, `Other Software`, `Users`) are hard-coded in the sidebar for this version.

The sidebar nav items must be rendered from a data source (not hard-coded strings in the Vue template) so that configurable categories can be introduced later with a schema migration. The decision to make them configurable is deferred to a future change.

---

## "All" Overview Tab Scope

**Q: (Implicit) What does the All Assets view show?**
**A: Hardware only** (Computer, Smartphone, Tablet) — matching the current `ALL_TAB_CATS` set. License categories (Windows Keys, Antivirus, Other Software) are excluded because they have different column schemas (no `user_id`, no `warranty_expiry`).

Columns shown in All Assets: Type, ID, Name, Model, User, Status, Purchase Date, Warranty Expiry, Notes.

---

## Design Token Reference

### Colors

```css
/* Sidebar */
--sidebar-bg:      #0E1520;
--sidebar-hover:   #131E2E;
--sidebar-active:  #182337;
--sidebar-border:  #141D2B;
--sidebar-text:    #6B82A0;
--sidebar-text-hi: #DDE8F5;
--sidebar-accent:  #2563EB;

/* Page structure */
--page:         #EEF1F7;
--surface:      #FFFFFF;
--border:       #DDE1EC;
--border-light: #EAECF4;

/* Text */
--text:   #0F172A;
--text-2: #64748B;
--text-3: #94A3B8;

/* Status strip / badge colors */
--s-active:   #2563EB;
--s-repair:   #EA580C;
--s-spare:    #7C3AED;
--s-retired:  #9CA3AF;
--s-missing:  #1E293B;
--s-expiring: #D97706;  /* row background: #FFFBEB */
--s-expired:  #DC2626;  /* row background: #FFF5F5 */

/* Interactive */
--accent:   #2563EB;
--accent-d: #1D4ED8;    /* hover state */

/* Signature: table header uses sidebar color */
--th-bg:   #0E1520;
--th-text: #8FA5BF;
```

### Typography

```css
--font-ui:   -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
--font-mono: 'Cascadia Code', 'Fira Code', 'Consolas', 'Courier New', monospace;
```

Monospace font applies to: IDs, license keys, date columns, the DB path in the statusbar.

### Border radius

```
Structural containers: 0 (sidebar, topbar, table, edit panel)
Badges: 2px
Buttons: 0
```

### Sidebar width and topbar height

```
--sidebar-w: 200px
--topbar-h:  48px
```

---

## Native Menu Bar

**Q: Where do file-level operations (Open/New Database, About, License) live?**
**A: Native OS menu bar via `wails/v2/pkg/menu`.** Matches the current PyQt6 app's menu structure so user habits transfer directly.

Menu structure:
```
File
  New Database…       Ctrl+N
  Open Database…      Ctrl+O
  ────────────────────────────
  Export CSV…         Ctrl+E
  ────────────────────────────
  Quit                Ctrl+Q

Help
  About
  ────────────────────────────
  License
```

The topbar Export CSV button is a secondary in-app access point. `Ctrl+E` stays as Export CSV — no keyboard shortcut conflict with the edit panel (panel is managed by row selection / Enter / Esc).

Note for Phase 1: native menu callbacks must be registered in Go at app startup. Bridge method names for file dialog operations (`NewDatabase`, `OpenDatabase`, `ExportCSV`) should be defined in Phase 1 architecture.

---

## Component Checklist (minimum for Phase 2 prototype)

| Component | Description |
|---|---|
| `AppShell.vue` | Root layout: sidebar + main area flex row |
| `Sidebar.vue` | Nav items, section labels, counts, footer |
| `SidebarItem.vue` | Single nav item with icon, label, count, active state |
| `Topbar.vue` | Search, Export CSV, Alerts badge, Dark mode toggle |
| `DataTable.vue` | Status strip col, sortable headers, row states, keyboard nav |
| `StatusBadge.vue` | Dot + label, 7 variants |
| `EditPanel.vue` | Slide-in panel with header, scrollable body, footer |
| `FormField.vue` | Label + control wrapper (text, date, select, textarea) |
| `AlertsModal.vue` | Overlay with expiry/warranty alert table |
| `Statusbar.vue` | Record count, sort state, DB path |
