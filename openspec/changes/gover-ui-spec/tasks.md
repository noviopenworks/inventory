## 1. Answer UI/UX Questions

- [x] 1.1 Read every question in the "UI And UX" section of `gover/QUESTIONS.md` and verify each has a matching decision in `specs/ui-decisions/spec.md`
- [x] 1.2 Annotate `gover/QUESTIONS.md` — mark each UI/UX question as resolved with a reference to the spec
- [x] 1.3 Confirm the two open questions in `design.md` (All tab scope, Ctrl+E conflict) are resolved or deferred with a documented rationale

## 2. Validate Design Token System

- [x] 2.1 Review all color values in `specs/ui-decisions/spec.md` against `gover/UI_DESIGN_REFERENCE.html` for consistency
- [x] 2.2 Verify status strip colors cover all `STATUS_OPTIONS` from `db/config.py` (Active, In Repair, Spare, Retired, Missing) plus expiring/expired states
- [x] 2.3 Confirm font stacks render acceptably on both Linux (WSL/Ubuntu) and Windows — note any adjustments needed

## 3. Resolve Adjacent Asset-Model Questions

- [x] 3.1 Document the decision on fixed vs. configurable categories in `specs/ui-decisions/spec.md` (already drafted — confirm with stakeholder)
- [x] 3.2 Document the "All" tab scope decision (hardware only vs. all categories)
- [x] 3.3 Review remaining asset-model questions in `gover/QUESTIONS.md` and flag any that directly affect the UI (serial numbers, locations, tags, etc.)

## 4. Finalize the Design Reference

- [x] 4.1 Add a Compact density mockup to `gover/UI_DESIGN_REFERENCE.html` (33 px → 24 px rows)
- [x] 4.2 Add a Dark mode skeleton to `gover/UI_DESIGN_REFERENCE.html` (sidebar accent colors, surface/page inversion)
- [x] 4.3 Add the Alerts overlay view to `gover/UI_DESIGN_REFERENCE.html`
- [x] 4.4 Verify the HTML reference renders correctly in Chromium (the Wails WebView engine on Linux)

## 5. Documentation Hand-Off

- [x] 5.1 Update `gover/QUESTIONS.md` header to note that the UI/UX section has been resolved via this change
- [x] 5.2 Add a link to `specs/ui-decisions/spec.md` and `gover/UI_DESIGN_REFERENCE.html` from `gover/ARCHITECTURE_NOTES.md` under "UI Strategy"
- [x] 5.3 Confirm no application code was written during this change (verify git diff shows only docs/html changes)
