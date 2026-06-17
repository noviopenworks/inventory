# Gover

Planning and governance workspace for the next major version of IT Asset Inventory.

This directory is intentionally documentation-only. It does not scaffold Wails, Vue, Go, TypeScript, Tailwind, or any runtime code yet. Future implementation should begin only after the open questions are answered and the development plan is approved.

## Goal

Prepare a controlled rewrite plan for a modern desktop application using:

- Wails for the desktop shell and Go/TypeScript bridge
- Go for backend application services and SQLite access
- Vue and TypeScript for the frontend
- Tailwind CSS for styling
- Taskfile for repeatable developer workflows
- SQLite as the local data store unless a later decision changes this

## Source Of Truth

The current PyQt6 application remains the feature and behavior reference until replacement behavior is explicitly designed.

Important current modules:

- `app/main_window.py`: current UI orchestration
- `app/dialogs.py`: add/edit form behavior
- `app/models.py`: table display behavior
- `db/schema.py`: current SQLite schema
- `db/config.py`: asset categories, columns, labels, and status options
- `db/queries.py`: CRUD and list queries
- `db/alerts.py`: expiry and warranty alert logic

## Files

- `DEVELOPMENT_PLAN.md`: phased roadmap for the rewrite
- `QUESTIONS.md`: decisions to resolve before implementation
- `ARCHITECTURE_NOTES.md`: architecture areas to investigate
- `TASKS.md`: future task backlog template
- `internal-tools.md`: planned local tooling and task commands
- `AGENTS.md`: instructions for coding agents working in this area
- `CLAUDE.md`: Claude-specific collaboration instructions

## Rule

Do not implement application code from this directory until the planning questions are reviewed and an implementation phase is explicitly approved.
