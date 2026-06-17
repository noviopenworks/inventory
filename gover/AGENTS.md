# Agent Instructions For Gover

These instructions apply to future coding agents working inside `gover/`.

## Current Status

`gover/` is a planning-only workspace. Do not create application code, Wails scaffolding, package manifests, generated bindings, or build configuration until implementation is explicitly approved.

## Intended Stack

- Wails
- Go
- Vue
- TypeScript
- Tailwind CSS
- SQLite
- Taskfile

## Source Reference

Use the current PyQt6 app as the behavioral reference:

- `app/main_window.py` for current UI orchestration
- `app/dialogs.py` for form behavior
- `app/models.py` for table behavior
- `db/schema.py` for schema details
- `db/config.py` for categories, columns, labels, and status options
- `db/queries.py` for CRUD/list behavior
- `db/alerts.py` for expiry and warranty alerts
- `tests/` for expected behavior

## Working Rules

- Preserve user data as a first-class requirement.
- Do not change the existing Python app unless explicitly asked.
- Do not introduce backward compatibility code without a concrete compatibility requirement.
- Prefer small, testable backend services over large bridge methods.
- Keep Wails bridge methods stable and typed.
- Keep Vue components focused on presentation and interaction, not database rules.
- Add tests with implementation work.
- Update documentation when decisions change.

## Required Before Implementation

- Review `QUESTIONS.md`.
- Confirm the selected migration strategy.
- Confirm supported platforms.
- Confirm first milestone scope.
- Confirm task commands and quality gates.

## Do Not Do Yet

- Do not run `wails init`.
- Do not create `go.mod`.
- Do not create `package.json`.
- Do not install dependencies.
- Do not create Vue components.
- Do not create Go services.
- Do not modify CI.
