# Claude Instructions For Gover

This directory exists to plan the next major version of IT Asset Inventory. Treat it as a planning and governance area until the user explicitly approves implementation.

## Collaboration Style

- Keep questions focused and ask one at a time when decisions are needed.
- Prefer concrete options with trade-offs.
- Separate decisions from implementation steps.
- Do not scaffold the app unless the user explicitly asks for development to begin.

## Current Planning Choice

The selected approach is planning-only governance. The desired output is documentation under `gover/`, not application code.

## Future Implementation Guardrails

- Use Wails, Go, Vue, TypeScript, Tailwind CSS, and SQLite unless the plan changes.
- Keep the current PyQt6 app as the behavior reference.
- Validate database compatibility before write operations.
- Start with a read-only prototype before CRUD.
- Add repeatable Taskfile commands before relying on manual commands.
- Verify work with tests and documented commands before claiming completion.

## Useful Current Files

- `README.md`
- `Taskfile.yml`
- `pyproject.toml`
- `db/schema.py`
- `db/config.py`
- `db/queries.py`
- `db/alerts.py`
- `app/main_window.py`
- `app/dialogs.py`
- `app/models.py`
- `tests/`

## Planning Files

- `gover/README.md`
- `gover/DEVELOPMENT_PLAN.md`
- `gover/QUESTIONS.md`
- `gover/ARCHITECTURE_NOTES.md`
- `gover/TASKS.md`
- `gover/internal-tools.md`
- `gover/AGENTS.md`
