---
change: gover-phase1-architecture
design-doc: docs/superpowers/specs/2026-06-17-gover-phase1-architecture-design.md
base-ref: bf8ad333e08fce05f6660502d12bcc6b8da57e62
---

# Gover Phase 1 Architecture — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Mark the Phase 1 architecture as resolved in `gover/ARCHITECTURE_NOTES.md` and confirm no application code leaked into the diff.

**Architecture:** This is a DOCUMENTATION-ONLY change. The canonical spec (`openspec/changes/gover-phase1-architecture/specs/architecture/spec.md`) and the Design Doc (`docs/superpowers/specs/2026-06-17-gover-phase1-architecture-design.md`) are already written and are the primary deliverables. The only remaining work is to annotate `gover/ARCHITECTURE_NOTES.md` with the resolved status and to verify the diff contains no application code.

**Tech Stack:** Markdown, git diff inspection. No Go, Vue, TypeScript, or Wails code is written in this change.

## Global Constraints

- This change produces documentation only — no Go, Vue, TypeScript, Python, or Wails code.
- `gover2/` must not be created or modified.
- `app/`, `db/`, and all Python source files must be untouched.
- Annotation style must be consistent with the existing resolved annotation in `gover/ARCHITECTURE_NOTES.md` (see the `## UI Strategy Options To Evaluate` section for the pattern to follow).
- Spec link in the annotation must use a repo-root-relative path.

---

## Already-Complete Tasks (tasks 1.1 – 5.3)

The following task groups from `openspec/changes/gover-phase1-architecture/tasks.md` are **already complete** — their deliverables exist in the repo:

| Group | Tasks | Deliverable |
|---|---|---|
| 1. Go Package Architecture | 1.1 – 1.7 | `openspec/changes/gover-phase1-architecture/specs/architecture/spec.md` §1, §2 |
| 2. Vue + TypeScript App Structure | 2.1 – 2.5 | spec §3 |
| 3. Tailwind Configuration | 3.1 – 3.4 | spec §4 |
| 4. Test Strategy | 4.1 – 4.3 | spec §5 |
| 5. Taskfile Commands | 5.1 – 5.3 | spec §6 |

No work is needed for these. Do not re-create or re-write the spec or design doc.

---

## File Map

| File | Action | Responsibility |
|---|---|---|
| `gover/ARCHITECTURE_NOTES.md` | **Modify** | Add resolved annotation to the Phase 1 section |

No files are created. No other files are modified.

---

### Task 1: Annotate ARCHITECTURE_NOTES.md with Phase 1 resolved status

**Files:**
- Modify: `gover/ARCHITECTURE_NOTES.md`

**Interfaces:**
- Consumes: existing resolved annotation pattern from `## UI Strategy Options To Evaluate` (line 52–57 in `gover/ARCHITECTURE_NOTES.md`)
- Produces: a `> **Resolved …**` blockquote above each Phase 1 candidate section that was settled by this change

**Context:** `gover/ARCHITECTURE_NOTES.md` is a living planning doc. It contains several "Candidate …" and "… Options To Evaluate" sections that record open questions. When a question is settled by a spec change, a resolved blockquote is inserted immediately above the relevant section body, following the established pattern visible in the `## UI Strategy Options To Evaluate` section:

```markdown
> **Resolved (YYYY-MM-DD, change `change-name`):** <one sentence summary>.
> Full decisions: [`path/to/spec.md`](../path/to/spec.md)
```

- [x] **Step 1: Read the current file to confirm context and line numbers**

Read `gover/ARCHITECTURE_NOTES.md` in full. Identify the sections that Phase 1 architecture resolves:
- `## Candidate Runtime Boundaries` — resolved: Go/Vue boundary confirmed
- `## Candidate Go Areas` — resolved: 6-package `internal/` layout locked
- `## Candidate Frontend Areas` — resolved: `features/`, `lib/`, `components/` structure locked
- `## Data Strategy Options To Evaluate` — resolved: Option 1 (reuse current schema) adopted

Note the line numbers of each section header so the annotation insertion is precise.

- [x] **Step 2: Insert resolved annotation above `## Candidate Runtime Boundaries`**

Insert the following blockquote immediately after the `## Candidate Runtime Boundaries` heading line and before the first bullet:

```markdown
> **Resolved (2026-06-17, change `gover-phase1-architecture`):** Go owns all data access and business logic via a 6-package `internal/` layout; Vue owns presentation; TypeScript API wrappers in `src/lib/api/` mediate the Wails bridge.
> Full decisions: [`openspec/changes/gover-phase1-architecture/specs/architecture/spec.md`](../openspec/changes/gover-phase1-architecture/specs/architecture/spec.md)
```

- [x] **Step 3: Insert resolved annotation above `## Candidate Go Areas`**

Insert the following blockquote immediately after the `## Candidate Go Areas` heading line and before the first bullet:

```markdown
> **Resolved (2026-06-17, change `gover-phase1-architecture`):** Six packages confirmed: `internal/database`, `internal/models`, `internal/services`, `internal/bridge`, `internal/config`, `internal/backup`. Dependency order locked; `bridge` is the sole package that imports from all others.
> Full decisions: [`openspec/changes/gover-phase1-architecture/specs/architecture/spec.md`](../openspec/changes/gover-phase1-architecture/specs/architecture/spec.md)
```

- [x] **Step 4: Insert resolved annotation above `## Candidate Frontend Areas`**

Insert the following blockquote immediately after the `## Candidate Frontend Areas` heading line and before the first bullet:

```markdown
> **Resolved (2026-06-17, change `gover-phase1-architecture`):** Feature-folder structure adopted: `features/` (assets, licenses, users, alerts stub), `components/` (10 shared components), `lib/api/index.ts`, `lib/types/index.ts`, `stores/assets.ts`, `stores/ui.ts`. Hash-mode Vue Router with 9 routes.
> Full decisions: [`openspec/changes/gover-phase1-architecture/specs/architecture/spec.md`](../openspec/changes/gover-phase1-architecture/specs/architecture/spec.md)
```

- [x] **Step 5: Insert resolved annotation above `## Data Strategy Options To Evaluate`**

Insert the following blockquote immediately after the `## Data Strategy Options To Evaluate` heading line and before the first numbered item:

```markdown
> **Resolved (2026-06-17, change `gover-phase1-architecture`):** Option 1 adopted — reuse current schema directly. Go DDL in `internal/database` carries the 7-table schema verbatim from `db/schema.py`. SQLite driver: `modernc.org/sqlite` (pure Go, no CGo).
> Full decisions: [`openspec/changes/gover-phase1-architecture/specs/architecture/spec.md`](../openspec/changes/gover-phase1-architecture/specs/architecture/spec.md)
```

- [x] **Step 6: Verify the file reads correctly**

Read the modified `gover/ARCHITECTURE_NOTES.md` and confirm:
- Four blockquotes are present, one per resolved section.
- Each blockquote contains a link using the path `../openspec/changes/gover-phase1-architecture/specs/architecture/spec.md` (relative to `gover/`).
- The date in every blockquote is `2026-06-17`.
- The change name in every blockquote is `gover-phase1-architecture`.
- No existing content was removed or reordered.

- [x] **Step 7: Commit**

```bash
git add gover/ARCHITECTURE_NOTES.md
git commit -m "docs: annotate ARCHITECTURE_NOTES.md Phase 1 sections as resolved

Marks Candidate Runtime Boundaries, Candidate Go Areas, Candidate
Frontend Areas, and Data Strategy Options as resolved by the
gover-phase1-architecture change. Links to the locked spec."
```

---

### Task 2: Verify no application code in diff

**Files:**
- No files are written in this task — read-only verification.

**Interfaces:**
- Consumes: git diff output against `bf8ad333e08fce05f6660502d12bcc6b8da57e62`
- Produces: confirmation (or list of violations) that no Go, Vue, TypeScript, Python, or Wails application code was introduced

- [x] **Step 1: Collect the full diff against the base ref**

```bash
git diff bf8ad333e08fce05f6660502d12bcc6b8da57e62 --name-only
```

Expected output: only documentation and spec files. Example of an acceptable output:
```
docs/superpowers/plans/2026-06-17-gover-phase1-architecture.md
docs/superpowers/specs/2026-06-17-gover-phase1-architecture-design.md
gover/ARCHITECTURE_NOTES.md
openspec/changes/gover-phase1-architecture/specs/architecture/spec.md
openspec/changes/gover-phase1-architecture/tasks.md
```
(Other openspec change metadata files such as `.comet.yaml`, `proposal.md`, `design.md` are also acceptable.)

- [x] **Step 2: Check for prohibited file extensions**

```bash
git diff bf8ad333e08fce05f6660502d12bcc6b8da57e62 --name-only | grep -E '\.(go|vue|ts|py|mod|sum|yaml|toml|html)$' | grep -v '^openspec/' | grep -v '^docs/'
```

Expected output: empty. If any file appears here that is outside `openspec/` or `docs/`, it is a violation — list it and stop.

Permitted exceptions (these extensions are allowed):
- Any file under `openspec/changes/gover-phase1-architecture/` — these are spec/metadata files, not application code.
- Any file under `docs/superpowers/` — these are planning documents.
- `gover/ARCHITECTURE_NOTES.md` — planning document.

- [x] **Step 3: Check for application directories**

```bash
git diff bf8ad333e08fce05f6660502d12bcc6b8da57e62 --name-only | grep -E '^gover2/'
```

Expected output: empty. If any `gover2/` path appears, it is a violation — the `gover2/` scaffold must not be created in this change.

```bash
git diff bf8ad333e08fce05f6660502d12bcc6b8da57e62 --name-only | grep -E '^(app/|db/|tests/)'
```

Expected output: empty. Python app directories must be untouched.

- [x] **Step 4: Record verification result**

If all three checks above produced empty output, the diff is clean. No further action needed.

If any violation was found, list it explicitly as:
```
VIOLATION: <file path> — reason: <application code / scaffold / Python modification>
```
Then stop and report to the user before proceeding.

---

## Self-Review Checklist

**Spec coverage:**

| Spec requirement | Covered by |
|---|---|
| Annotate `gover/ARCHITECTURE_NOTES.md` with resolved status | Task 1 |
| Verify no application code in diff | Task 2 |
| Spec and Design Doc already written (tasks 1.1–5.3) | Pre-complete (noted above) |

**Placeholder scan:** No "TBD", "TODO", or "implement later" text in this plan. All steps show exact commands, exact file paths, and exact text to insert.

**Consistency:** The resolved blockquote format in steps 2–5 of Task 1 matches the existing pattern in `gover/ARCHITECTURE_NOTES.md` lines 52–57. The spec link path `../openspec/changes/gover-phase1-architecture/specs/architecture/spec.md` is relative to the `gover/` directory where `ARCHITECTURE_NOTES.md` lives, matching the existing `../openspec/changes/gover-ui-spec/specs/ui-decisions/spec.md` link pattern.
