---
title: AGENTS
section: 1
header: Agent Smith Manual
date: 2025-12-14
---

# agents(1)

## NAME

agents - switch between different agent personas using symlinks

## SYNOPSIS

**agents** [command] [flags]

## DESCRIPTION

**agents** is a CLI tool for managing `AGENTS.md` symlinks. It allows you to define multiple "personas" (e.g., `AGENTS.coder.md`, `AGENTS.writer.md`) and quickly switch the active context for AI agents by updating a canonical symlink.

It adheres to the **Canonical Target as Source of Truth** pattern: the state of the symlink determines the active persona.

## COMMANDS

### list

List all available agent personas found in the configured `agents_dir`.

### use [persona]

Switch the active persona to the specified one.

**Example:**
`agents use coder`

**Behavior:**
1. Updates the canonical `target_file` (symlink) to point to `AGENTS.coder.md`.
2. Updates any other configured targets (copies/links) to match the new persona.
3. Saves the state.

**Flags:**
* `--target-file`: Specify an additional target to apply/track for this operation.

### status

Show the current status of the agent system.

* Identifies the **Active Persona** based on where the canonical symlink points.
* Lists all managed targets and their status vs the active persona:
    * `[OK]`: Matches active persona.
    * `[DRIFT]`: Points to a different persona.
    * `[MISSING]`: File does not exist.

### reconcile

Re-apply the currently active persona to all configured targets.

**Purpose:**
Fixes drift or restores missing files.

**Drift Handling:**
If you manually change the canonical symlink (e.g., `ln -sf ...`), `reconcile` accepts this change as the new truth and updates all other targets to match it.

### unuse

Deactivate the current persona.

* Removes the canonical target symlink/file.
* Removes all other tracked targets.
* Clears the active persona from state.

### drop [persona]

Stop tracking a specific target or an entire persona.

**Example:**
`agents drop coder --target-file ./local_copy.md`

### version

Print the version number.

## CONFIGURATION

The tool follows the XDG Base Directory Specification.

### Files

* **Configuration**: `$XDG_CONFIG_HOME/agent-smith/config.yaml` (default: `~/.config/agent-smith/config.yaml`)
* **Personas**: `$XDG_DATA_HOME/agent-smith/personas` (default: `~/.local/share/agent-smith/personas`)
* **State**: `$XDG_STATE_HOME/agent-smith/status.yaml` (default: `~/.local/state/agent-smith/status.yaml`)

### Canonical Target

The **active persona link** (Source of Truth) defaults to:
`$XDG_CONFIG_HOME/agents/AGENTS.md`

### Environment Variables

* `AGENTS_TARGET_FILE`: Override the path to the canonical symlink.
* `AGENTS_AGENTS_DIR`: Override the directory to search for personas.

## EXAMPLES

**List available personas:**
```bash
agents list
```

**Switch to 'coder':**
```bash
agents use coder
```

**Reconcile state after manual changes:**
```bash
agents reconcile
```
