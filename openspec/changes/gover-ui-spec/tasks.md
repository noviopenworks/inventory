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
