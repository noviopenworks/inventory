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
