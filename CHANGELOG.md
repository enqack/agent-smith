# Changelog

All notable changes to this project will be documented in this file.

## [0.3.0] - 2025-12-14

### Documentation System Overhaul
- **Read the Docs Integration**: Migrated documentation system to Sphinx + MyST + Furo theme.
- **Documentation Map**: Created `docs/index.md` as the unified entry point for guides and reference manual.
- **Canonical Source Linking**: Linked root documents (`CONTRIBUTING.md`, `CHANGELOG.md`, `LICENSE`) into `docs/` to avoid duplication.
- **Build System**: Updated `flake.nix` and `magefile.go` to support Sphinx builds.
- **CI/CD**: Added `Build Docs` step to GitHub Actions.
- **Cleanup**: Removed deprecated `mkdocs.yml`.

### Refactoring & Architecture

- **Canonical Target as Source of Truth**: The main symlink (default `$XDG_CONFIG_HOME/agents/AGENTS.md`) is now the single source of truth for the active persona.
- **Reconcile**: Now respects manual changes to the main symlink (drift) and propagates them to secondary targets.
- **Use**: Optimized to calculate targets once and fail fast.
- **Unuse**: Now explicitly clears internal state to prevent "resurrection" of old targets.
- **Robustness**: 
    - Fixed potential resource leaks (`defer` in loop).
    - Added `chmod 0644` to copied files.
    - Downgraded `go.mod` to stable 1.23.0.
    - **State Management**: `SaveState` now **replaces** targets for a persona instead of merging/accumulating, preventing stale targets from persisting.
    - **Drop Safety**: Added guard to prevent `drop` from removing directories.
- **Modularization**: Restructured codebase to Domain-Driven Design layout. Moved CLI logic to `internal/cli`, state management to `internal/state`, configuration to `internal/config`, and operations to `internal/ops`.
- **ApplyPersona**: Improved reliability and error reporting in `internal/ops/apply.go`.
- **Environment Handling**: robust mapping of environment variables to configuration (e.g., `AGENTS_TARGET_FILE`).

### Features

- **Multi-Persona Status**: `status` command now lists all tracked personas in the state file, indicating active status for each ([internal/cli/status.go](file:///home/sysop/Projects/agent-smith/internal/cli/status.go)).
- **Multi-Target Accumulation**: `use` command now merges targets for the same persona instead of overwriting, enabling complex setups with multiple links/copies per persona.
- **Enhanced Error Reporting**: CLI now consistently returns non-zero exit codes on failure and provides clearer error messages for permissions and missing files.

### Testing

- **End-to-End Suite**: Implemented comprehensive E2E tests in `tests/e2e/` covering:
    - **Flows**: Happy path for `list`, `use`, `status`, `reconcile`.
    - **Configuration Precedence**: Verified Flag > Env > Config File priority.
    - **Error Recovery**: Automated verification of recovery from missing files, drift, and permission errors.
- **Unit Tests**: Added dedicated unit tests for `config`, `ops`, and `state` packages.

### CI/CD

- **GitHub Actions**: Added `.github/workflows/ci.yml` to run `nix flake check`, unit tests, and E2E tests on detailed PRs and main branch pushes.
- **Contribution Guide**: Added `CONTRIBUTING.md` with development instructions.

### Documentation

- **Refactoring**: Separated `README.md` to focus on high-level overview. Technical reference moved to `docs/man/`.
- **Man Pages**: Added comprehensive man page suite:
    - `agents(1)`: Main CLI reference.
    - `agents-config(5)`: Configuration file format.
    - `agents-status(5)`: Internal status file format.
    - `agents-format(7)`: AGENTS.md concepts and specifications.

## [0.2.0] - 2025-12-13

### Core Features & Logic

- **XDG Compliance**: implemented XDG Base Directory support (`cmd/xdg.go`).
  - Config is now looked for in `$XDG_CONFIG_HOME/agent-smith` (default `~/.config/agent-smith`).
  - Personas are now stored in `$XDG_DATA_HOME/agent-smith/personas` (default `~/.local/share/agent-smith/personas`).
  - State is stored in `$XDG_STATE_HOME/agent-smith` (default `~/.local/state/agent-smith`).
- **Directory Structure**: Renamed the default definition storage folder from `agents` to `personas`.
- **New Command**: Added `version` command (`cmd/version.go`) reading from a new `VERSION` file.
- **Build System**: Updated `flake.nix` and `magefile.go` to inject the version string during build.
- **Man Pages**: Added `Docs` target to `magefile.go` and configured `flake.nix` to generate and install man pages using `pandoc`.

### CLI Enhancements

- **Output Refinement**:
    - `use`: Swapped output messages and changed terminology ("Persona switched").
    - `status`: Removed verbose path output, keeping the user-friendly "Active Persona" message.
- **Terminology**: Unified terminology to use "Persona" instead of "Agent" across help text and output (`cmd/list.go`, etc).

### Documentation

- Updated all paths to reflect XDG standards in `README.md`.
- Added documentation for the `version` command.
- Standardized terminology to "Persona".
- Fixed typos.
- **Man Page**: Added automatic generation of `agents.1` man page from `README.md`.

### Testing

- **Integration Fix**: Updated `cmd/cli_test.go` to match the actual output of the `use` command.
- **Test Updates**: Updated `cmd/root_test.go` and `cmd/state_test.go` to verify XDG path resolution and the new `personas` directory default.
