# Agent Smith (agents)

`agents` is a CLI tool for managing `AGENTS.md` symlinks, allowing users to switch between different agent personas easily.
It aligns with the [AGENTS.md](https://agents.md/) specification by allowing you to dynamically switch the active context for AI agents.

## The Dynamic Context Pattern

The [AGENTS.md](https://agents.md/) specification provides a standard location for agent instructions. However, a single static file can become unwieldy when switching between different types of work (e.g., coding, architecture, documentation).

`agents` solves this by turning `AGENTS.md` into a **dynamic interface**:
1.  **Define Personas**: Create focused context files like `AGENTS.coder.md` or `AGENTS.architect.md`.
2.  **Switch Context**: Run `agents use coder`.
3.  **Standard Compliance**: The tool symlinks your chosen persona to `AGENTS.md`.

This ensures any AI tool looking for `AGENTS.md` automatically receives the specific instructions relevant to your current task.

## Features

- **Profile Management**: Manage multiple agent profiles (personas).
- **Symlink Switching**: Instantly switch the `AGENTS.md` symlink to point to the desired persona.
- **Configurable**: extensive configuration via YAML files or environment variables.

## Installation

### Using Nix (Flakes)

This is the recommended way to install and run `agents`.

#### Run directly

```bash
nix run github:enqack/agent-smith -- list
```

#### Build locally

```bash
nix build .
./result/bin/agents list
```

#### FHS Environment

If you need a standard Linux filesystem hierarchy (simulated), use the `fhs` output:

```bash
nix run .#fhs -- list
```

### Manual Build

Requries Go 1.21+.

```bash
go build -o bin/agents .
./bin/agents help
```

## Usage

### Configuration

By default, `agents` looks for agent personas in the following order:
1. `$HOME/.config/agents`
2. `$HOME/.config/agent-smith/agents`
3. `/usr/share/agent-smith/agents`

You can override this by configuring `agents_dir` in:
- `$HOME/.config/agents/config.yaml`
- `/etc/agents/config.yaml`

Example `config.yaml`:

```yaml
agents_dir: 
  - "/path/to/primary/agents"
  - "/path/to/fallback/agents"
target_file: "~/.config/agents/AGENTS.md"
```

### Commands

#### List available agents

```bash
agents list
```

#### Switch to an agent

```bash
agents use coder
```
This will link `AGENTS.md` to `AGENTS.coder.md` found in your `agents_dir`.

#### Check status

```bash
agents status
```
Shows which agent is currently active.

## Directory Structure

Expected structure for `agents_dir`:

```
/path/to/agents/
  ├── AGENTS.coder.md
  ├── AGENTS.writer.md
  └── AGENTS.manager.md
```
