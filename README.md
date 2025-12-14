# Agent Smith (agents)

**Agent Smith** (`agents`) is a specialized CLI tool for managing AI context files. It allows you to define multiple "personas" (confgurations of `AGENTS.md`) and instantly switch between them using symbolic links.

It implements the [AGENTS.md](https://agents.md/) specification pattern dynamically.

## Quick Start

### 1. Install (via Nix)

```bash
nix run github:enqack/agent-smith -- --help
```

### 2. Define Personas

Create your context files in `~/.local/share/agent-smith/personas/`:

```bash
mkdir -p ~/.local/share/agent-smith/personas
echo "You happen to be a senior Go engineer." > ~/.local/share/agent-smith/personas/AGENTS.coder.md
echo "You are a technical writer." > ~/.local/share/agent-smith/personas/AGENTS.writer.md
```

### 3. Switch Context

```bash
agents use coder
```

This creates a symlink at `~/.config/agents/AGENTS.md` pointing to your coder persona. Any AI tool configured to read that path will now see the Coder instructions.

## Documentation

* **[Man Page (CLI)](docs/man/agents.1.md)**: Detailed command reference.
* **[Configuration](docs/man/agents-config.5.md)**: `config.yaml` reference.
* **[Status File](docs/man/agents-status.5.md)**: `status.yaml` reference.
* **[Concepts](docs/man/agents-format.7.md)**: AGENTS.md format and specifications.
* **[Contributing](CONTRIBUTING.md)**: Development setup and guide.

## Features

* **Dynamic Context**: Switch personas instantly.
* **Source of Truth**: The active symlink determines the state; the tool reconciles everything else to match it.
* **Drift Detection**: `agents status` shows if your files have drifted from the active persona.
* **XDG Compliant**: Follows standard Linux directory specs.

## Build

```bash
nix build .
```
