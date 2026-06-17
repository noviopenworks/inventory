---
change: gover-ui-spec
design-doc: docs/superpowers/specs/2026-06-17-gover-ui-spec-design.md
base-ref: 5d93ada840ad4c35c1d44add95d41b3783fe2141
---

# Gover v2 UI Spec Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Finalise the UI decisions spec, annotate the open-questions document, add three mockup views (compact density, dark mode, alerts overlay) to the HTML design reference, and link everything from the architecture notes — all documentation changes, zero application code.

**Architecture:** This is a pure documentation change. The canonical spec lives at `openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md`; the visual reference lives at `gover/UI_DESIGN_REFERENCE.html`. Supporting documents (`gover/QUESTIONS.md`, `gover/ARCHITECTURE_NOTES.md`) are annotated to point readers to the resolved decisions.

**Tech Stack:** Markdown, HTML/CSS (static mockup, no JavaScript required for the new views), plain text editing.

## Global Constraints

- No application code: the git diff must contain only `.md` and `.html` file changes.
- All colour values, token names, and component names must match the values already in `openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md` verbatim (e.g. `#0E1520`, `--sidebar-bg`, `33px`/`24px` row heights).
- The HTML file must remain a single self-contained file (no external CSS or JS imports).
- The spec's Native Menu Bar section must already exist in `specs/ui-decisions/spec.md` before any annotation work references it — Task 1 ensures this.
- Every task ends with a `git commit` scoped to that task's files only.

---

## 1. Finalize the UI Decisions Spec

**Files:**
- Modify: `openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md`

**Interfaces:**
- Produces: A complete spec that covers all 9 UI decisions documented in the design doc, including the Native Menu Bar section that Tasks 2 and 3 reference.

- [ ] **Step 1.1: Verify the Native Menu Bar section exists**

Open `openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md` and confirm the section heading `## Native Menu Bar` is present with the `wails/v2/pkg/menu` menu structure.

If the section is absent, add it at the bottom of the spec (before the Component Checklist) with exactly this content:

```markdown
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
```

- [ ] **Step 1.2: Verify spec coverage against all QUESTIONS.md UI/UX items**

Read `gover/QUESTIONS.md`, section `## UI And UX`. The 7 questions there are:

1. "Should the current tab model be preserved?" → spec section `## Navigation Model`
2. "Should the new UI use a sidebar, dashboard, command palette, or another navigation model?" → spec section `## Navigation Model`
3. "Should the first screen be an overview dashboard or the asset table?" → spec section `## First Screen`
4. "Should add/edit use modal dialogs, side panels, or full pages?" → spec section `## Add / Edit UX Pattern`
5. "How important is keyboard-first operation?" → spec section `## Keyboard-First Operation`
6. "Should the UI support compact and comfortable density modes?" → spec section `## Density Modes`
7. "Should dark mode follow the system theme or remain a manual toggle?" → spec section `## Dark Mode`

Confirm each spec section exists. If any is missing, add it following the existing Q/A format used by the other sections.

Also confirm the `## "All" Overview Tab Scope` section exists (the implicit question resolved during design), and the `## Asset Categories` section covering the fixed-vs-configurable decision.

- [ ] **Step 1.3: Check design token completeness**

In `spec.md`, under `## Design Token Reference`, confirm these CSS custom properties are listed:

Light mode: `--sidebar-bg`, `--sidebar-active`, `--sidebar-text`, `--sidebar-text-hi`, `--sidebar-accent`, `--sidebar-w`, `--topbar-h`, `--page`, `--surface`, `--border`, `--text`, `--text-2`, `--text-3`, `--s-active`, `--s-repair`, `--s-spare`, `--s-retired`, `--s-missing`, `--s-expiring`, `--s-expired`, `--accent`, `--accent-d`, `--th-bg`, `--th-text`, `--font-ui`, `--font-mono`.

If any token from `gover/UI_DESIGN_REFERENCE.html`'s `:root` block is absent from the spec, add it with the matching hex value.

- [ ] **Step 1.4: Commit the finalized spec**

```bash
git add openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md
git commit -m "docs(gover-ui-spec): finalize ui-decisions spec with native menu bar and full token coverage"
```

Expected: 1 file changed.

---

## 2. Annotate gover/QUESTIONS.md

**Files:**
- Modify: `gover/QUESTIONS.md`

**Interfaces:**
- Consumes: The finalized spec from Task 1 (section names must exist before being referenced).
- Produces: `gover/QUESTIONS.md` with the UI/UX section marked resolved, each question pointing to the spec.

- [ ] **Step 2.1: Add a resolution header to the UI/UX section**

In `gover/QUESTIONS.md`, find the line:

```
## UI And UX
```

Replace it with:

```markdown
## UI And UX

> **Resolved** — all questions answered in `openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md`
> (change: `gover-ui-spec`, 2026-06-17). See also `gover/UI_DESIGN_REFERENCE.html`.
```

- [ ] **Step 2.2: Annotate each UI/UX question with its spec reference**

For each of the 7 questions below, add a one-line annotation directly after the question bullet. Use this exact format:

```
  → *Resolved:* spec § Navigation Model
```

Full annotation map:

| Question text (first ~6 words) | Annotation |
|---|---|
| Should the current tab model… | `→ *Resolved:* spec § Navigation Model — replaced with persistent left sidebar (200 px, dark bg)` |
| Should the new UI use a sidebar… | `→ *Resolved:* spec § Navigation Model — sidebar chosen` |
| Should the first screen be an overview… | `→ *Resolved:* spec § First Screen — Computers table (dashboard deferred to Phase 4)` |
| Should add/edit use modal dialogs… | `→ *Resolved:* spec § Add / Edit UX Pattern — slide-in right panel (300 px)` |
| How important is keyboard-first operation? | `→ *Resolved:* spec § Keyboard-First Operation — all PyQt6 shortcuts preserved, panel shortcuts added` |
| Should the UI support compact and comfortable… | `→ *Resolved:* spec § Density Modes — Comfortable 33 px default, Compact 24 px toggle` |
| Should dark mode follow the system theme… | `→ *Resolved:* spec § Dark Mode — manual toggle, Light default, system-follow deferred to Phase 4` |

**How to do this:** For each question, the bullet in `QUESTIONS.md` currently looks like:

```
- Should the current tab model be preserved?
```

Change it to:

```
- Should the current tab model be preserved?
  → *Resolved:* spec § Navigation Model — replaced with persistent left sidebar (200 px, dark bg)
```

- [ ] **Step 2.3: Commit the annotation**

```bash
git add gover/QUESTIONS.md
git commit -m "docs(gover-ui-spec): annotate QUESTIONS.md UI/UX section with spec resolutions"
```

Expected: 1 file changed.

---

## 3. Add Mockup Views to UI_DESIGN_REFERENCE.html

**Files:**
- Modify: `gover/UI_DESIGN_REFERENCE.html`

**Interfaces:**
- Consumes: Design token values from Task 1's spec (all hex values must match `:root` block already in the HTML file).
- Produces: Three additional view sections appended before `</body>`: compact density view, dark mode skeleton, alerts overlay view.

The existing HTML file ends at line ~977 with:
```html
</div><!-- /shell -->

</body>
</html>
```

All three new views are appended as sibling `<div class="shell">` blocks (or appropriately labelled `<section>` wrappers) before `</body>`.

**Important:** All colour values must exactly match the existing `:root` CSS variables already in the file — do not introduce new hex values for colours that already have a token.

- [ ] **Step 3.1: Add a section separator comment style**

Before adding new views, confirm the existing separator comment pattern used in the file (search for `═══`). New sections must use the same comment style for consistency.

The existing pattern is:
```html
<!-- ═══════════════════════════════════════════════════════
     SECTION NAME
════════════════════════════════════════════════════════ -->
```

- [ ] **Step 3.2: Add the Compact Density view**

Find the closing `</div><!-- /shell -->` followed by the `</body>` tag. Insert the following block immediately before `</body>`:

```html

<!-- ═══════════════════════════════════════════════════════
     VIEW: COMPACT DENSITY (24 px rows)
════════════════════════════════════════════════════════ -->
<p style="font-family:var(--ui);font-size:11px;color:var(--text-2);padding:8px 16px;background:var(--page);border-top:1px solid var(--border);margin-top:24px;">
  VIEW: Compact density — 24 px row height (matches current PyQt6 app)
</p>
<div class="shell" style="height:320px;margin-top:0;">
  <nav class="sidebar">
    <!-- (same sidebar markup as primary view — copy sidebar <nav> content here) -->
    <div class="sidebar-head"><div class="logo"><svg width="20" height="20" viewBox="0 0 20 20" fill="none"><rect x="2" y="2" width="7" height="7" rx="1" fill="#2563EB"/><rect x="11" y="2" width="7" height="7" rx="1" fill="#2563EB" opacity=".5"/><rect x="2" y="11" width="7" height="7" rx="1" fill="#2563EB" opacity=".5"/><rect x="11" y="11" width="7" height="7" rx="1" fill="#2563EB" opacity=".3"/></svg></div><span class="logo-text">IT Inventory</span></div>
    <div class="sidebar-section"><div class="section-label">Overview</div>
      <div class="sidebar-item"><span style="opacity:.5">⊞</span><span class="item-label">All Assets</span><span class="item-count">42</span></div>
    </div>
    <div class="sidebar-section"><div class="section-label">Hardware</div>
      <div class="sidebar-item active"><span style="opacity:.5">⬜</span><span class="item-label">Computers</span><span class="item-count">18</span></div>
      <div class="sidebar-item"><span style="opacity:.5">📱</span><span class="item-label">Smartphones</span><span class="item-count">9</span></div>
      <div class="sidebar-item"><span style="opacity:.5">⬜</span><span class="item-label">Tablets</span><span class="item-count">5</span></div>
    </div>
    <div class="sidebar-footer"><span style="color:var(--sidebar-text);font-size:11px;">v2.0.0 · inventory.db</span></div>
  </nav>
  <div class="main">
    <div class="topbar"><div class="search-wrap"><input class="search" type="text" placeholder="Search computers…" style="width:220px;"></div><div style="display:flex;gap:8px;margin-left:auto;align-items:center;"><button class="btn-outline">Export CSV</button><button class="btn-outline" style="position:relative;">Alerts<span style="position:absolute;top:-4px;right:-4px;background:var(--s-expired);color:#fff;border-radius:10px;font-size:10px;padding:0 5px;">2</span></button></div></div>
    <div class="content" style="padding:12px 16px;">
      <div class="table-wrap">
        <table class="data-table" style="--row-h:24px;">
          <thead><tr>
            <th style="width:4px;padding:0;background:var(--th-bg);"></th>
            <th style="background:var(--th-bg);color:var(--th-text);font-size:11px;padding:0 8px;height:28px;">ID</th>
            <th style="background:var(--th-bg);color:var(--th-text);font-size:11px;padding:0 8px;">Name</th>
            <th style="background:var(--th-bg);color:var(--th-text);font-size:11px;padding:0 8px;">Model</th>
            <th style="background:var(--th-bg);color:var(--th-text);font-size:11px;padding:0 8px;">User</th>
            <th style="background:var(--th-bg);color:var(--th-text);font-size:11px;padding:0 8px;">Status</th>
            <th style="background:var(--th-bg);color:var(--th-text);font-size:11px;padding:0 8px;">Warranty Expiry</th>
          </tr></thead>
          <tbody>
            <tr style="background:#EBF1FB;height:24px;">
              <td style="width:4px;padding:0;background:var(--s-active);"></td>
              <td style="padding:0 8px;font-family:var(--mono);font-size:11px;height:24px;">C-001</td>
              <td style="padding:0 8px;font-size:12px;height:24px;font-weight:500;">WKSTN-MGLOW</td>
              <td style="padding:0 8px;font-size:12px;height:24px;color:var(--text-2);">MacBook Pro 16"</td>
              <td style="padding:0 8px;font-size:12px;height:24px;color:var(--text-2);">Głowacki, M.</td>
              <td style="padding:0 8px;height:24px;"><span class="status-badge" style="background:rgba(37,99,235,.12);color:var(--s-active);font-size:10px;padding:1px 6px;border-radius:2px;">● Active</span></td>
              <td style="padding:0 8px;font-family:var(--mono);font-size:11px;height:24px;color:var(--text-2);">2026-04-12</td>
            </tr>
            <tr style="background:#FFFBEB;height:24px;">
              <td style="width:4px;padding:0;background:var(--s-expiring);"></td>
              <td style="padding:0 8px;font-family:var(--mono);font-size:11px;height:24px;">C-002</td>
              <td style="padding:0 8px;font-size:12px;height:24px;font-weight:500;">LAPTOP-SALES</td>
              <td style="padding:0 8px;font-size:12px;height:24px;color:var(--text-2);">Dell XPS 15</td>
              <td style="padding:0 8px;font-size:12px;height:24px;color:var(--text-2);">Nowak, A.</td>
              <td style="padding:0 8px;height:24px;"><span class="status-badge" style="background:rgba(217,119,6,.12);color:var(--s-expiring);font-size:10px;padding:1px 6px;border-radius:2px;">● Expiring</span></td>
              <td style="padding:0 8px;font-family:var(--mono);font-size:11px;height:24px;color:var(--s-expiring);">2026-07-03</td>
            </tr>
            <tr style="background:var(--surface);height:24px;">
              <td style="width:4px;padding:0;background:var(--s-spare);"></td>
              <td style="padding:0 8px;font-family:var(--mono);font-size:11px;height:24px;">C-003</td>
              <td style="padding:0 8px;font-size:12px;height:24px;font-weight:500;">SPARE-01</td>
              <td style="padding:0 8px;font-size:12px;height:24px;color:var(--text-2);">Lenovo ThinkPad</td>
              <td style="padding:0 8px;font-size:12px;height:24px;color:var(--text-2);">—</td>
              <td style="padding:0 8px;height:24px;"><span class="status-badge" style="background:rgba(124,58,237,.12);color:var(--s-spare);font-size:10px;padding:1px 6px;border-radius:2px;">● Spare</span></td>
              <td style="padding:0 8px;font-family:var(--mono);font-size:11px;height:24px;color:var(--text-2);">2027-01-15</td>
            </tr>
            <tr style="background:var(--surface);height:24px;">
              <td style="width:4px;padding:0;background:var(--s-retired);"></td>
              <td style="padding:0 8px;font-family:var(--mono);font-size:11px;height:24px;">C-004</td>
              <td style="padding:0 8px;font-size:12px;height:24px;font-weight:500;">OLD-MGMT-01</td>
              <td style="padding:0 8px;font-size:12px;height:24px;color:var(--text-2);">HP EliteBook</td>
              <td style="padding:0 8px;font-size:12px;height:24px;color:var(--text-2);">—</td>
              <td style="padding:0 8px;height:24px;"><span class="status-badge" style="background:rgba(156,163,175,.15);color:var(--s-retired);font-size:10px;padding:1px 6px;border-radius:2px;">● Retired</span></td>
              <td style="padding:0 8px;font-family:var(--mono);font-size:11px;height:24px;color:var(--text-2);">2023-06-01</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <div class="statusbar"><span class="sb-item"><strong>18</strong> computers</span><div class="sb-rule"></div><span class="sb-item">sorted by <strong>id</strong> desc</span><div class="sb-rule"></div><span class="sb-item">density: <strong>compact</strong></span></div>
  </div>
</div>
```

Verify: rows must be visually shorter than the primary view. Table header row height 28 px, data rows 24 px (matching current PyQt6 default). All CSS classes (`sidebar`, `main`, `topbar`, `content`, `statusbar`, `sb-item`, `sb-rule`, `data-table`) must already exist in the file's `<style>` block — do not add new class definitions.

- [ ] **Step 3.3: Add the Dark Mode skeleton view**

Append the following block immediately after the compact density block (still before `</body>`):

```html

<!-- ═══════════════════════════════════════════════════════
     VIEW: DARK MODE SKELETON
════════════════════════════════════════════════════════ -->
<p style="font-family:var(--ui);font-size:11px;color:var(--text-2);padding:8px 16px;background:var(--page);border-top:1px solid var(--border);margin-top:24px;">
  VIEW: Dark mode — manual toggle, Light default. System-follow deferred to Phase 4.
</p>
<div class="shell" style="height:320px;margin-top:0;background:#0F172A;">
  <nav class="sidebar" style="background:#060D18;border-right-color:#0A1120;">
    <div class="sidebar-head"><div class="logo"><svg width="20" height="20" viewBox="0 0 20 20" fill="none"><rect x="2" y="2" width="7" height="7" rx="1" fill="#2563EB"/><rect x="11" y="2" width="7" height="7" rx="1" fill="#2563EB" opacity=".5"/><rect x="2" y="11" width="7" height="7" rx="1" fill="#2563EB" opacity=".5"/><rect x="11" y="11" width="7" height="7" rx="1" fill="#2563EB" opacity=".3"/></svg></div><span class="logo-text" style="color:#DDE8F5;">IT Inventory</span></div>
    <div class="sidebar-section"><div class="section-label" style="color:#3A4F6A;">Overview</div>
      <div class="sidebar-item"><span style="opacity:.4;color:#6B82A0;">⊞</span><span class="item-label" style="color:#6B82A0;">All Assets</span><span class="item-count">42</span></div>
    </div>
    <div class="sidebar-section"><div class="section-label" style="color:#3A4F6A;">Hardware</div>
      <div class="sidebar-item active" style="background:#182337;border-left-color:#2563EB;"><span style="opacity:.6;color:#DDE8F5;">⬜</span><span class="item-label" style="color:#DDE8F5;">Computers</span><span class="item-count" style="color:#4A6480;">18</span></div>
      <div class="sidebar-item"><span style="opacity:.4;color:#6B82A0;">📱</span><span class="item-label" style="color:#6B82A0;">Smartphones</span><span class="item-count">9</span></div>
      <div class="sidebar-item"><span style="opacity:.4;color:#6B82A0;">⬜</span><span class="item-label" style="color:#6B82A0;">Tablets</span><span class="item-count">5</span></div>
    </div>
    <div class="sidebar-footer"><span style="color:#3A4F6A;font-size:11px;">v2.0.0 · inventory.db</span></div>
  </nav>
  <div class="main" style="background:#0F172A;">
    <div class="topbar" style="background:#1E293B;border-bottom-color:#0A1120;">
      <div class="search-wrap"><input class="search" type="text" placeholder="Search computers…" style="width:220px;background:#0F172A;color:#DDE8F5;border-color:#334155;"></div>
      <div style="display:flex;gap:8px;margin-left:auto;align-items:center;">
        <button class="btn-outline" style="color:#94A3B8;border-color:#334155;background:transparent;">Export CSV</button>
        <button class="btn-outline" style="color:#94A3B8;border-color:#334155;background:transparent;position:relative;">Alerts<span style="position:absolute;top:-4px;right:-4px;background:var(--s-expired);color:#fff;border-radius:10px;font-size:10px;padding:0 5px;">2</span></button>
        <button class="btn-outline" style="color:#F1C40F;border-color:#334155;background:transparent;font-size:16px;line-height:1;" title="Switch to Light mode">☀</button>
      </div>
    </div>
    <div class="content" style="padding:12px 16px;background:#0F172A;">
      <div class="table-wrap" style="border-color:#1E293B;">
        <table class="data-table">
          <thead><tr>
            <th style="width:4px;padding:0;background:#060D18;"></th>
            <th style="background:#060D18;color:#4A6480;font-size:11px;padding:0 8px;height:33px;">ID</th>
            <th style="background:#060D18;color:#4A6480;font-size:11px;padding:0 8px;">Name</th>
            <th style="background:#060D18;color:#4A6480;font-size:11px;padding:0 8px;">Model</th>
            <th style="background:#060D18;color:#4A6480;font-size:11px;padding:0 8px;">User</th>
            <th style="background:#060D18;color:#4A6480;font-size:11px;padding:0 8px;">Status</th>
            <th style="background:#060D18;color:#4A6480;font-size:11px;padding:0 8px;">Warranty Expiry</th>
          </tr></thead>
          <tbody>
            <tr style="background:#1A2947;height:33px;border-bottom:1px solid #1E293B;">
              <td style="width:4px;padding:0;background:var(--s-active);"></td>
              <td style="padding:0 8px;font-family:var(--mono);font-size:11px;color:#94A3B8;">C-001</td>
              <td style="padding:0 8px;font-size:12px;font-weight:500;color:#DDE8F5;">WKSTN-MGLOW</td>
              <td style="padding:0 8px;font-size:12px;color:#64748B;">MacBook Pro 16"</td>
              <td style="padding:0 8px;font-size:12px;color:#64748B;">Głowacki, M.</td>
              <td style="padding:0 8px;"><span style="background:rgba(37,99,235,.2);color:#60A5FA;font-size:10px;padding:1px 6px;border-radius:2px;">● Active</span></td>
              <td style="padding:0 8px;font-family:var(--mono);font-size:11px;color:#64748B;">2026-04-12</td>
            </tr>
            <tr style="background:#1C1A10;height:33px;border-bottom:1px solid #1E293B;">
              <td style="width:4px;padding:0;background:var(--s-expiring);"></td>
              <td style="padding:0 8px;font-family:var(--mono);font-size:11px;color:#94A3B8;">C-002</td>
              <td style="padding:0 8px;font-size:12px;font-weight:500;color:#DDE8F5;">LAPTOP-SALES</td>
              <td style="padding:0 8px;font-size:12px;color:#64748B;">Dell XPS 15</td>
              <td style="padding:0 8px;font-size:12px;color:#64748B;">Nowak, A.</td>
              <td style="padding:0 8px;"><span style="background:rgba(217,119,6,.2);color:#FBBF24;font-size:10px;padding:1px 6px;border-radius:2px;">● Expiring</span></td>
              <td style="padding:0 8px;font-family:var(--mono);font-size:11px;color:var(--s-expiring);">2026-07-03</td>
            </tr>
            <tr style="background:#0F172A;height:33px;border-bottom:1px solid #1E293B;">
              <td style="width:4px;padding:0;background:var(--s-spare);"></td>
              <td style="padding:0 8px;font-family:var(--mono);font-size:11px;color:#94A3B8;">C-003</td>
              <td style="padding:0 8px;font-size:12px;font-weight:500;color:#DDE8F5;">SPARE-01</td>
              <td style="padding:0 8px;font-size:12px;color:#64748B;">Lenovo ThinkPad</td>
              <td style="padding:0 8px;font-size:12px;color:#64748B;">—</td>
              <td style="padding:0 8px;"><span style="background:rgba(124,58,237,.2);color:#A78BFA;font-size:10px;padding:1px 6px;border-radius:2px;">● Spare</span></td>
              <td style="padding:0 8px;font-family:var(--mono);font-size:11px;color:#64748B;">2027-01-15</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <div class="statusbar" style="background:#1E293B;border-top-color:#0A1120;color:#64748B;">
      <span class="sb-item"><strong style="color:#94A3B8;">18</strong> computers</span>
      <div class="sb-rule" style="background:#334155;"></div>
      <span class="sb-item">dark mode <strong style="color:#94A3B8;">on</strong></span>
    </div>
  </div>
</div>
```

Verify: the dark shell must use a dark page background (`#0F172A`), sidebar must be darker than main area (`#060D18`), table header must use `#060D18` (matching dark sidebar), and the ☀ (sun) icon must be visible in the topbar representing the toggle to light mode.

Note from design doc risk: "On Linux WSL, `system-ui` falls back to DejaVu or Noto — check both platforms when adding compact density to the HTML reference." The same applies here for dark mode. No action needed in the HTML file, but verify visually.

- [ ] **Step 3.4: Add the Alerts Overlay view**

Append immediately after the dark mode block (still before `</body>`):

```html

<!-- ═══════════════════════════════════════════════════════
     VIEW: ALERTS OVERLAY (F5 / Alerts button)
════════════════════════════════════════════════════════ -->
<p style="font-family:var(--ui);font-size:11px;color:var(--text-2);padding:8px 16px;background:var(--page);border-top:1px solid var(--border);margin-top:24px;">
  VIEW: Alerts overlay — triggered by F5 or Alerts button. Shows warranty/expiry alerts. Closes on Esc or click-outside.
</p>
<div style="position:relative;height:420px;background:var(--page);overflow:hidden;margin-top:0;border-top:1px solid var(--border);">
  <!-- Dimmed shell behind overlay -->
  <div class="shell" style="height:420px;opacity:.35;pointer-events:none;">
    <nav class="sidebar">
      <div class="sidebar-head"><div class="logo"><svg width="20" height="20" viewBox="0 0 20 20" fill="none"><rect x="2" y="2" width="7" height="7" rx="1" fill="#2563EB"/><rect x="11" y="2" width="7" height="7" rx="1" fill="#2563EB" opacity=".5"/><rect x="2" y="11" width="7" height="7" rx="1" fill="#2563EB" opacity=".5"/><rect x="11" y="11" width="7" height="7" rx="1" fill="#2563EB" opacity=".3"/></svg></div><span class="logo-text">IT Inventory</span></div>
    </nav>
    <div class="main">
      <div class="topbar"></div>
    </div>
  </div>
  <!-- Overlay backdrop -->
  <div style="position:absolute;inset:0;background:rgba(15,23,42,.45);z-index:10;"></div>
  <!-- Alerts modal panel -->
  <div style="position:absolute;top:50%;left:50%;transform:translate(-50%,-50%);z-index:20;background:var(--surface);border:1px solid var(--border);border-radius:0;width:640px;max-height:320px;display:flex;flex-direction:column;box-shadow:0 8px 32px rgba(15,23,42,.18);">
    <!-- Modal header -->
    <div style="display:flex;align-items:center;justify-content:space-between;padding:12px 16px;border-bottom:1px solid var(--border);background:var(--th-bg);">
      <span style="font-size:13px;font-weight:600;color:#DDE8F5;letter-spacing:.02em;">ALERTS — Expiring &amp; Expired</span>
      <button style="background:none;border:none;cursor:pointer;color:#6B82A0;font-size:16px;line-height:1;padding:2px 6px;" title="Close (Esc)">×</button>
    </div>
    <!-- Modal table -->
    <div style="overflow-y:auto;flex:1;">
      <table style="width:100%;border-collapse:collapse;font-family:var(--ui);font-size:12px;">
        <thead>
          <tr style="background:var(--page);">
            <th style="padding:6px 10px;text-align:left;color:var(--text-2);font-weight:600;font-size:11px;border-bottom:1px solid var(--border);">Type</th>
            <th style="padding:6px 10px;text-align:left;color:var(--text-2);font-weight:600;font-size:11px;border-bottom:1px solid var(--border);">ID</th>
            <th style="padding:6px 10px;text-align:left;color:var(--text-2);font-weight:600;font-size:11px;border-bottom:1px solid var(--border);">Name</th>
            <th style="padding:6px 10px;text-align:left;color:var(--text-2);font-weight:600;font-size:11px;border-bottom:1px solid var(--border);">Warranty Expiry</th>
            <th style="padding:6px 10px;text-align:left;color:var(--text-2);font-weight:600;font-size:11px;border-bottom:1px solid var(--border);">Days</th>
          </tr>
        </thead>
        <tbody>
          <tr style="background:#FFF5F5;border-bottom:1px solid #FFE4E4;">
            <td style="padding:6px 10px;color:var(--text-2);">Computer</td>
            <td style="padding:6px 10px;font-family:var(--mono);font-size:11px;color:var(--text-2);">C-007</td>
            <td style="padding:6px 10px;font-weight:500;color:var(--text);">SERVER-BACK</td>
            <td style="padding:6px 10px;font-family:var(--mono);font-size:11px;color:var(--s-expired);">2026-04-01</td>
            <td style="padding:6px 10px;"><span style="background:rgba(220,38,38,.12);color:var(--s-expired);font-size:10px;padding:1px 6px;border-radius:2px;font-weight:600;">EXPIRED</span></td>
          </tr>
          <tr style="background:#FFFBEB;border-bottom:1px solid #FEF3C7;">
            <td style="padding:6px 10px;color:var(--text-2);">Computer</td>
            <td style="padding:6px 10px;font-family:var(--mono);font-size:11px;color:var(--text-2);">C-002</td>
            <td style="padding:6px 10px;font-weight:500;color:var(--text);">LAPTOP-SALES</td>
            <td style="padding:6px 10px;font-family:var(--mono);font-size:11px;color:var(--s-expiring);">2026-07-03</td>
            <td style="padding:6px 10px;"><span style="background:rgba(217,119,6,.12);color:var(--s-expiring);font-size:10px;padding:1px 6px;border-radius:2px;font-weight:600;">16 days</span></td>
          </tr>
          <tr style="background:#FFFBEB;border-bottom:1px solid #FEF3C7;">
            <td style="padding:6px 10px;color:var(--text-2);">Smartphone</td>
            <td style="padding:6px 10px;font-family:var(--mono);font-size:11px;color:var(--text-2);">S-004</td>
            <td style="padding:6px 10px;font-weight:500;color:var(--text);">PHONE-HR-02</td>
            <td style="padding:6px 10px;font-family:var(--mono);font-size:11px;color:var(--s-expiring);">2026-07-14</td>
            <td style="padding:6px 10px;"><span style="background:rgba(217,119,6,.12);color:var(--s-expiring);font-size:10px;padding:1px 6px;border-radius:2px;font-weight:600;">27 days</span></td>
          </tr>
        </tbody>
      </table>
    </div>
    <!-- Modal footer -->
    <div style="padding:8px 16px;border-top:1px solid var(--border);display:flex;justify-content:flex-end;gap:8px;background:var(--page);">
      <span style="font-size:11px;color:var(--text-3);align-self:center;">1 expired · 2 expiring within 30 days</span>
      <button class="btn-outline" style="font-size:12px;padding:4px 12px;">Close</button>
    </div>
  </div>
</div>
```

Verify: expired rows must use `#FFF5F5` background and `var(--s-expired)` (`#DC2626`) accents. Expiring rows must use `#FFFBEB` background and `var(--s-expiring)` (`#D97706`) accents. The modal header must use `var(--th-bg)` (`#0E1520`) — the signature dark header matching the table header in the primary view.

- [ ] **Step 3.5: Verify the HTML renders correctly in Chromium**

Open the file in Chromium on WSL:

```bash
chromium-browser /home/mg/inventory/gover/UI_DESIGN_REFERENCE.html 2>/dev/null &
# or, if that doesn't work:
wslview /home/mg/inventory/gover/UI_DESIGN_REFERENCE.html
```

Scroll to the bottom and confirm:
1. Compact density view: rows visibly shorter than the primary view above.
2. Dark mode view: dark backgrounds throughout, ☀ icon visible in topbar.
3. Alerts overlay: modal centred over dimmed shell, expired row red-tinted, expiring rows amber-tinted, header dark.

If Chromium is unavailable, verify by checking the HTML is well-formed:

```bash
grep -c "</div>" /home/mg/inventory/gover/UI_DESIGN_REFERENCE.html
# Note the count before and after edits — it should increase by the expected number of new divs
```

- [ ] **Step 3.6: Commit the HTML changes**

```bash
git add gover/UI_DESIGN_REFERENCE.html
git commit -m "docs(gover-ui-spec): add compact density, dark mode, and alerts overlay views to HTML reference"
```

Expected: 1 file changed, significant insertions.

---

## 4. Update Cross-Reference Document

**Files:**
- Modify: `gover/ARCHITECTURE_NOTES.md`

**Interfaces:**
- Consumes: The finalized spec path (`openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md`) and HTML reference path (`gover/UI_DESIGN_REFERENCE.html`) from Tasks 1 and 3.
- Produces: `gover/ARCHITECTURE_NOTES.md` with a "UI Strategy" section that links to both resolved documents, replacing or augmenting the existing "UI Strategy Options To Evaluate" prose.

- [ ] **Step 4.1: Locate the UI Strategy section in ARCHITECTURE_NOTES.md**

Read `gover/ARCHITECTURE_NOTES.md`. Find the section heading:

```
## UI Strategy Options To Evaluate
```

This section currently lists three options (Familiar table-first, Dashboard plus sidebar, Command-center) as unresolved candidates.

- [ ] **Step 4.2: Add a resolution note above the options**

Insert the following block immediately after the `## UI Strategy Options To Evaluate` heading line (before the numbered list):

```markdown
> **Resolved (2026-06-17, change `gover-ui-spec`):** Option 2 (Dashboard plus sidebar) adopted with modifications.
> Full decisions: [`openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md`](../openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md)
> Visual reference: [`gover/UI_DESIGN_REFERENCE.html`](UI_DESIGN_REFERENCE.html)
>
> Key decisions: persistent left sidebar (200 px, `#0E1520`), Computers as first screen, slide-in right panel for add/edit (300 px), Comfortable/Compact density toggle, manual dark mode toggle (Light default), 7 fixed categories, native OS menu bar via `wails/v2/pkg/menu`.

```

- [ ] **Step 4.3: Commit the cross-reference update**

```bash
git add gover/ARCHITECTURE_NOTES.md
git commit -m "docs(gover-ui-spec): link resolved UI strategy to spec and HTML reference in ARCHITECTURE_NOTES"
```

Expected: 1 file changed.

---

## 5. Final Verification

**Files:**
- Read-only: git diff, all modified files

**Interfaces:**
- Consumes: All four committed files from Tasks 1–4.
- Produces: Confirmed clean state — no application code in diff, all cross-references valid.

- [ ] **Step 5.1: Verify no application code in the diff**

```bash
git diff 5d93ada840ad4c35c1d44add95d41b3783fe2141..HEAD --name-only
```

Expected output: only these file paths (order may vary):
```
openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md
gover/QUESTIONS.md
gover/UI_DESIGN_REFERENCE.html
gover/ARCHITECTURE_NOTES.md
```

If any `.py`, `.go`, `.ts`, `.vue`, or other source file appears, stop and investigate before proceeding.

- [ ] **Step 5.2: Verify all internal links resolve**

Check each file reference added during the change:

```bash
# Spec exists
test -f /home/mg/inventory/openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md && echo "spec OK"

# HTML reference exists
test -f /home/mg/inventory/gover/UI_DESIGN_REFERENCE.html && echo "html OK"

# QUESTIONS.md annotation present
grep -c "Resolved.*spec §" /home/mg/inventory/gover/QUESTIONS.md

# ARCHITECTURE_NOTES.md link present
grep -c "gover-ui-spec" /home/mg/inventory/gover/ARCHITECTURE_NOTES.md
```

Expected: all `echo` lines print, `grep -c` returns ≥ 7 for QUESTIONS.md (one per UI/UX question) and ≥ 1 for ARCHITECTURE_NOTES.md.

- [ ] **Step 5.3: Verify tasks.md items map to plan tasks**

Read `openspec/changes/gover-ui-spec/tasks.md`. Confirm each checkbox item has a corresponding task/step in this plan:

| tasks.md item | Plan task |
|---|---|
| 1.1 Read UI/UX questions, verify spec coverage | Task 1, Step 1.2 |
| 1.2 Annotate QUESTIONS.md | Task 2, Steps 2.1–2.2 |
| 1.3 Confirm open questions resolved | Task 1, Steps 1.1–1.2 |
| 2.1 Review color values for consistency | Task 1, Step 1.3 |
| 2.2 Verify status strip colors cover all STATUS_OPTIONS | Task 1, Step 1.3 |
| 2.3 Confirm font stacks render acceptably | Task 3, Step 3.5 note |
| 3.1 Document fixed vs configurable categories | Task 1, Step 1.2 (§ Asset Categories) |
| 3.2 Document All tab scope | Task 1, Step 1.2 (§ All Overview Tab Scope) |
| 3.3 Review remaining asset-model questions | (verification only — flag any UI-affecting items found) |
| 4.1 Add Compact density mockup | Task 3, Step 3.2 |
| 4.2 Add Dark mode skeleton | Task 3, Step 3.3 |
| 4.3 Add Alerts overlay | Task 3, Step 3.4 |
| 4.4 Verify HTML renders in Chromium | Task 3, Step 3.5 |
| 5.1 Update QUESTIONS.md header | Task 2, Step 2.1 |
| 5.2 Link from ARCHITECTURE_NOTES.md | Task 4, Steps 4.1–4.2 |
| 5.3 Confirm no application code in git diff | Task 5, Step 5.1 |

If tasks.md item 3.3 reveals any asset-model question that directly constrains the UI (serial numbers, locations, or tags affecting column layout), document the finding as a comment in `spec.md` under a new `## Deferred / Flagged Questions` section rather than resolving it here.
