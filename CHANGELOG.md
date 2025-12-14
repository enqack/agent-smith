# Changelog

All notable changes to this project will be documented in this file.

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
