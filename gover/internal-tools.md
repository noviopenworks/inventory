# Internal Tools Plan

This file lists tools and commands to define later. It does not install or configure them yet.

## Expected Tools

- Go toolchain for backend code and tests.
- Wails CLI for desktop app development and packaging.
- Node.js for frontend tooling.
- Vue and TypeScript for the frontend application.
- Tailwind CSS for styling.
- SQLite tooling for database inspection and troubleshooting.
- Taskfile for consistent local commands.

## Candidate Task Commands

The future rewrite should expose commands similar to the current Python project while matching Wails conventions.

- `task gover:install`: install frontend and backend dependencies.
- `task gover:dev`: run the Wails development server.
- `task gover:test`: run Go and frontend tests.
- `task gover:test:go`: run Go tests.
- `task gover:test:web`: run frontend tests.
- `task gover:lint`: run Go and frontend linters.
- `task gover:fmt`: format Go and frontend code.
- `task gover:type:check`: run TypeScript checks.
- `task gover:check`: run lint, type-check, and tests.
- `task gover:build`: build the desktop application.
- `task gover:package:linux`: build Linux packages if required.
- `task gover:package:windows`: build Windows package if required.
- `task gover:clean`: remove generated build output.

## Quality Gates To Define

- Go unit tests for services, database access, migrations, and export logic.
- Frontend tests for forms, table behavior, filters, and state transitions.
- TypeScript strictness level.
- Linting rules for Go, Vue, TypeScript, and CSS.
- Formatting rules.
- Manual QA checklist for packaged builds.

## Generated Files Policy To Decide

- Whether Wails-generated TypeScript bindings are committed.
- Whether build artifacts are ignored.
- Whether packaged release output follows the existing `dist/` layout.
- Whether the new version uses the root `.app-version` file or its own version source.

## Environment Questions

- Which Node package manager should be used?
- Which Go version should be targeted?
- Which Wails major version should be used?
- Should local development require CGO because of the SQLite driver?
- Should package builds happen only in CI or also locally?
