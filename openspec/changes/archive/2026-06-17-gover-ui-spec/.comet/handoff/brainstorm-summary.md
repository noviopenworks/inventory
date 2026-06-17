# Brainstorm Summary

- Change: gover-ui-spec
- Date: 2026-06-17

## Confirmed Technical Approach

Spec-first (Approach A):
1. Apply Spec Patch: add "Native Menu Bar" section to `specs/ui-decisions/spec.md`
2. Annotate `gover/QUESTIONS.md` — mark each UI/UX question resolved with spec reference
3. Add HTML views to `UI_DESIGN_REFERENCE.html`: compact density, dark mode skeleton, alerts overlay
4. Update `gover/ARCHITECTURE_NOTES.md` with links to spec and reference

## Key Trade-offs and Risks

- Native menu bar requires `wails/v2/pkg/menu` API surface — deferred to build phase, but the decision is now locked in the spec
- Ctrl+E shortcut conflict resolved: keep as Export CSV, no panel toggle shortcut needed (panel managed by row selection + Enter/Esc)
- All tab scope confirmed: hardware only (matching current ALL_TAB_CATS)

## Testing Strategy

Documentation completeness check:
- Every question in QUESTIONS.md "UI And UX" section maps 1:1 to a spec decision
- HTML reference demonstrates: primary view, compact density, dark mode, alerts overlay, edit panel
- No application code in git diff

## Spec Patches

Add to `specs/ui-decisions/spec.md`: "Native Menu Bar" section documenting that the Wails app uses `wails/v2/pkg/menu` for File (New/Open Database, Export CSV, Quit) and Help (About, License) actions, matching the current PyQt6 app's menu structure.
